package utils

import (
	"bytes"
	"io"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func Decoder(types string, r io.Reader) *transform.Reader {
	var decoder *transform.Reader
	switch types {
	case "UTF-8":
		decoder = transform.NewReader(r, unicode.UTF8.NewDecoder())
	case "GBK":
		decoder = transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
	case "GB2312":
		decoder = transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	case "GB18030":
		decoder = transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	case "BIG5":
		decoder = transform.NewReader(r, traditionalchinese.Big5.NewDecoder())
	}
	return decoder
}

func Encoder(types string, r io.Reader) *transform.Reader {
	var encoder *transform.Reader
	switch types {
	case "UTF-8":
		encoder = transform.NewReader(r, unicode.UTF8.NewEncoder())
	case "GBK":
		encoder = transform.NewReader(r, simplifiedchinese.GBK.NewEncoder())
	case "GB2312":
		encoder = transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())
	case "GB18030":
		encoder = transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())
	case "BIG5":
		encoder = transform.NewReader(r, traditionalchinese.Big5.NewEncoder())
	}
	return encoder
}

func Decode(types string, msg *[]byte) *[]byte {
	var decoder *transform.Reader
	switch types {
	case "UTF-8":
		decoder = transform.NewReader(bytes.NewReader(*msg), unicode.UTF8.NewDecoder())
	case "GBK":
		decoder = transform.NewReader(bytes.NewReader(*msg), simplifiedchinese.GBK.NewDecoder())
	case "GB2312":
		decoder = transform.NewReader(bytes.NewReader(*msg), simplifiedchinese.GB18030.NewDecoder())
	case "GB18030":
		decoder = transform.NewReader(bytes.NewReader(*msg), simplifiedchinese.GB18030.NewDecoder())
	case "BIG5":
		decoder = transform.NewReader(bytes.NewReader(*msg), traditionalchinese.Big5.NewDecoder())
	}
	content, _ := ioutil.ReadAll(decoder)
	return &content
}

func Encode(types string, msg *[]byte) *[]byte {
	var decoder *transform.Reader
	switch types {
	case "UTF-8":
		decoder = transform.NewReader(bytes.NewReader(*msg), unicode.UTF8.NewEncoder())
	case "GBK":
		decoder = transform.NewReader(bytes.NewReader(*msg), simplifiedchinese.GBK.NewEncoder())
	case "GB2312":
		decoder = transform.NewReader(bytes.NewReader(*msg), simplifiedchinese.GB18030.NewEncoder())
	case "GB18030":
		decoder = transform.NewReader(bytes.NewReader(*msg), simplifiedchinese.GB18030.NewEncoder())
	case "BIG5":
		decoder = transform.NewReader(bytes.NewReader(*msg), traditionalchinese.Big5.NewEncoder())
	}
	content, _ := ioutil.ReadAll(decoder)
	return &content
}
