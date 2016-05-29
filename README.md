BaiduOCR
========

Read Chinese, Japanese and English text from JPEG/PNG image with Baidu OCR services.

[![CircleCI](https://circleci.com/gh/caiguanhao/baiduocr.svg?style=svg)](https://circleci.com/gh/caiguanhao/baiduocr)

```go
results, err := baiduocr.OCR{
	APIKey: "", // your Baidu OCR API Key
}.ParseImageFile("test/fixtures/chinese/hanzi.jpg")
if err == nil {
	fmt.Println(results[0] == "漢字")
}
```

See [docs](https://godoc.org/github.com/caiguanhao/baiduocr) for usage and examples.

LICENSE: MIT

Copyright (C) 2016 Cai Guanhao (Choi Goon-ho)
