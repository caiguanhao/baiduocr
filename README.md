BaiduOCR
========

Read Chinese and English text from JPEG/PNG image with Baidu OCR services.

[![CircleCI](https://circleci.com/gh/caiguanhao/baiduocr.svg?style=svg)](https://circleci.com/gh/caiguanhao/baiduocr)

```go
ocr := baiduocr.OCR{APIKey: APIKey}
results, err := ocr.ParseImageFile("test/fixtures/chinese/hanzi.jpg")
```

See [docs](https://godoc.org/github.com/caiguanhao/baiduocr) for usage and examples.

LICENSE: MIT

Copyright (C) 2016 Cai Guanhao (Choi Goon-ho)
