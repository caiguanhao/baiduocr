// Read Chinese, Japanese and English text from JPEG/PNG image with Baidu OCR services.
// PNG image will be converted to JPEG on the fly because Baidu OCR recognizes only JPEG image files.
package baiduocr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	OCR struct {
		// Set API key
		APIKey string
		// Set API entrypoint path, default is http://apis.baidu.com/apistore/idlocr/ocr
		APIPath string
		// Set request timeout in milliseconds (ms), default is 5000, set to -1 means no timeout
		TimeoutInMilliseconds int64
	}

	BaiduOCROption struct {
		f func(*baiduOCROption)
	}

	baiduOCROption struct {
		languageType string

		pngBackgroundColor color.Color
	}

	baiduOCRRet struct {
		ErrMsg  string `json:"errMsg"`
		RetData []struct {
			Rect struct {
				Height string `json:"height"`
				Left   string `json:"left"`
				Top    string `json:"top"`
				Width  string `json:"width"`
			} `json:"rect"`
			Word string `json:"word"`
		} `json:"retData"`
	}
)

const (
	DEFAULT_LANG = "CHN_ENG"

	CHINESE  = "CHN_ENG"
	ENGLISH  = "ENG"
	JAPANESE = "JAP"
)

// Option to set OCR language type to Chinese (and English). This is the default option for language type.
func SetLanaguageTypeToChinese() BaiduOCROption {
	return BaiduOCROption{func(option *baiduOCROption) { option.languageType = CHINESE }}
}

// Option to set OCR language type to English.
func SetLanguageTypeToEnglish() BaiduOCROption {
	return BaiduOCROption{func(option *baiduOCROption) { option.languageType = ENGLISH }}
}

// Option to set OCR language type to Japanese.
func SetLanguageTypeToJapanese() BaiduOCROption {
	return BaiduOCROption{func(option *baiduOCROption) { option.languageType = JAPANESE }}
}

// If the image is a PNG with transparent background, use this option to set the background color.
func SetPNGBackgroundColor(bgColor color.Color) BaiduOCROption {
	return BaiduOCROption{func(option *baiduOCROption) { option.pngBackgroundColor = bgColor }}
}

// Set the PNG background color with RGBA values.
func SetPNGBackgroundColorRGBA(r, g, b, a uint8) BaiduOCROption {
	return BaiduOCROption{func(option *baiduOCROption) { option.pngBackgroundColor = color.RGBA{r, g, b, a} }}
}

func (ocr OCR) ParseImage(imageBytes []byte, options ...BaiduOCROption) (results []string, err error) {
	switch http.DetectContentType(imageBytes) {
	case "image/png":
		results, err = ocr.ParsePNG(imageBytes, options...)
	case "image/jpeg":
		results, err = ocr.ParseJPEG(imageBytes, options...)
	default:
		err = errors.New("unrecognized image file format")
	}
	return
}

func (ocr OCR) ParseJPEG(imageBytes []byte, options ...BaiduOCROption) (results []string, err error) {
	opts := baiduOCROption{
		languageType: DEFAULT_LANG,
	}
	for _, option := range options {
		option.f(&opts)
	}

	reqBody := strings.NewReader(url.Values{
		"fromdevice":   {"pc"},
		"clientip":     {"10.10.10.0"},
		"detecttype":   {"LocateRecognize"},
		"languagetype": {opts.languageType},
		"imagetype":    {"1"},
		"image":        {base64.StdEncoding.EncodeToString(imageBytes)},
		"version":      {"v1"},
		"sizetype":     {"small"},
	}.Encode())

	path := ocr.APIPath
	if len(path) == 0 {
		path = "http://apis.baidu.com/apistore/idlocr/ocr"
	}

	var req *http.Request
	req, err = http.NewRequest("POST", path, reqBody)
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("apikey", ocr.APIKey)

	var timeout time.Duration
	ms := ocr.TimeoutInMilliseconds
	if ms < -1 {
		panic("TimeoutInMilliseconds must not be less than -1")
	} else if ms > -1 {
		if ms == 0 {
			ms = 5000
		}
		timeout = time.Duration(ms) * time.Millisecond
	}
	client := &http.Client{
		Timeout: timeout,
	}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var ret baiduOCRRet
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return
	}

	if len(ret.RetData) == 0 {
		msg := "BaiduOCR failed to recognize any text in the image."
		if ret.ErrMsg != "" {
			msg += fmt.Sprintf(" reason: %s", ret.ErrMsg)
		}
		err = errors.New(msg)
		return
	}
	for _, data := range ret.RetData {
		results = append(results, data.Word)
	}
	return
}

func (ocr OCR) ParsePNG(imageBytes []byte, options ...BaiduOCROption) (results []string, err error) {
	opts := baiduOCROption{}
	for _, option := range options {
		option.f(&opts)
	}

	var buffer *bytes.Buffer
	buffer, err = pngTojpeg(bytes.NewReader(imageBytes), opts.pngBackgroundColor)
	if err != nil {
		return
	}
	results, err = ocr.ParseJPEG((*buffer).Bytes(), options...)
	return
}

// Read text from image file of unknown type.
func (ocr OCR) ParseImageFile(filename string, options ...BaiduOCROption) (results []string, err error) {
	var file []byte
	file, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	results, err = ocr.ParseImage(file, options...)
	return
}

// Read text from JPEG image file.
func (ocr OCR) ParseJPEGFile(filename string, options ...BaiduOCROption) (results []string, err error) {
	var file []byte
	file, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	results, err = ocr.ParseJPEG(file, options...)
	return
}

// Read text from PNG image file. PNG image will be converted to JPEG image on the fly.
// By default, transparent background of PNG image will become black.
// You can add an option to specify the background color for better OCR results.
func (ocr OCR) ParsePNGFile(filename string, options ...BaiduOCROption) (results []string, err error) {
	var file []byte
	file, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	results, err = ocr.ParsePNG(file, options...)
	return
}

func pngTojpeg(reader io.Reader, pngBackgroundColor color.Color) (buffer *bytes.Buffer, err error) {
	var img image.Image
	img, err = png.Decode(reader)
	if err != nil {
		return
	}
	if pngBackgroundColor != nil {
		bounds := img.Bounds()
		newImg := image.NewRGBA(bounds)
		draw.Draw(newImg, bounds, &image.Uniform{pngBackgroundColor}, image.ZP, draw.Src)
		draw.Draw(newImg, bounds, img, image.ZP, draw.Over)
		img = newImg
	}
	buffer = new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, &jpeg.Options{100})
	return
}
