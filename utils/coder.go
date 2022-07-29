package utils

import (
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
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
	case "BIG5":
		decoder = transform.NewReader(r, traditionalchinese.Big5.NewDecoder())
	case "ShiftJIS":
		decoder = transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	case "KS_C_5601":
		decoder = transform.NewReader(r, korean.EUCKR.NewDecoder())
	case "GB2312":
		decoder = transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	case "GB18030":
		decoder = transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	case "Big5-HKSCS":
		decoder = transform.NewReader(r, traditionalchinese.Big5.NewDecoder())
	case "UTF-16":
		decoder = transform.NewReader(r, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder())
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
	case "BIG5":
		encoder = transform.NewReader(r, traditionalchinese.Big5.NewEncoder())
	case "ShiftJIS":
		encoder = transform.NewReader(r, japanese.ShiftJIS.NewEncoder())
	case "KS_C_5601":
		encoder = transform.NewReader(r, korean.EUCKR.NewEncoder())
	case "GB2312":
		encoder = transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())
	case "GB18030":
		encoder = transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())
	case "Big5-HKSCS":
		encoder = transform.NewReader(r, traditionalchinese.Big5.NewEncoder())
	case "UTF-16":
		encoder = transform.NewReader(r, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder())
	}
	return encoder
}
