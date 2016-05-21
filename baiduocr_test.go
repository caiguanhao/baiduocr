package baiduocr_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/caiguanhao/baiduocr"
)

var APIKey string = os.Getenv("BAIDUOCR_APIKEY")

func Example_solveSimpleCaptcha() {
	ocr := baiduocr.OCR{APIKey: APIKey}
	results, err := ocr.ParsePNGFile("test/fixtures/simple-captcha/3560.png", baiduocr.SetLangTypeENG())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strings.Join(results, ", "))
	// Output:
	// 3560
}

func Example_parseChineseText() {
	ocr := baiduocr.OCR{APIKey: APIKey}
	results, err := ocr.ParseJPEGFile("test/fixtures/chinese/hanzi.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strings.Join(results, ", "))
	// Output:
	// 漢字
}

func Example_parseVerticalChineseTextWithTransparentBackground() {
	ocr := baiduocr.OCR{APIKey: APIKey}
	// png file with a transparent background
	results, err := ocr.ParseImageFile("test/fixtures/chinese/vertical.png", baiduocr.SetPNGBackgroundColorRGBA(255, 255, 255, 255))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strings.Join(results, ", "))
	// Output:
	// 中文語法
}

func Example_parseHTTPResponse() {
	resp, err := http.Get("https://upload.wikimedia.org/wikipedia/commons/thumb/c/c9/Twemoji_1f21a.svg/200px-Twemoji_1f21a.svg.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	ocr := baiduocr.OCR{APIKey: APIKey}
	results, err := ocr.ParsePNG(body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strings.Join(results, ", "))
	// Output:
	// 無
}
