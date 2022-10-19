//go:build !windows
// +build !windows

package console

import (
	"io"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func DecoderReader(types string, r io.Reader) *transform.Reader {
	var decoder *transform.Reader
	types = strings.ToUpper(types)
	switch types {
	case "UTF-8":
		decoder = transform.NewReader(r, unicode.UTF8.NewDecoder())
	case "GBK":
		decoder = transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
	case "BIG5":
		decoder = transform.NewReader(r, traditionalchinese.Big5.NewDecoder())
	case "SHIFTJIS":
		decoder = transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	case "KS_C_5601":
		decoder = transform.NewReader(r, korean.EUCKR.NewDecoder())
	case "GB2312":
		decoder = transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	case "GB18030":
		decoder = transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	case "BIG5-HKSCS":
		decoder = transform.NewReader(r, traditionalchinese.Big5.NewDecoder())
	case "UTF-16":
		decoder = transform.NewReader(r, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder())
	default:
		decoder = transform.NewReader(r, unicode.UTF8.NewDecoder())
	}
	return decoder
}

func DecoderWriter(types string, r io.Writer) *transform.Writer {
	var decoder *transform.Writer
	types = strings.ToUpper(types)
	switch types {
	case "UTF-8":
		decoder = transform.NewWriter(r, unicode.UTF8.NewDecoder())
	case "GBK":
		decoder = transform.NewWriter(r, simplifiedchinese.GBK.NewDecoder())
	case "BIG5":
		decoder = transform.NewWriter(r, traditionalchinese.Big5.NewDecoder())
	case "ShiftJIS":
		decoder = transform.NewWriter(r, japanese.ShiftJIS.NewDecoder())
	case "KS_C_5601":
		decoder = transform.NewWriter(r, korean.EUCKR.NewDecoder())
	case "GB2312":
		decoder = transform.NewWriter(r, simplifiedchinese.GB18030.NewDecoder())
	case "GB18030":
		decoder = transform.NewWriter(r, simplifiedchinese.GB18030.NewDecoder())
	case "Big5-HKSCS":
		decoder = transform.NewWriter(r, traditionalchinese.Big5.NewDecoder())
	case "UTF-16":
		decoder = transform.NewWriter(r, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder())
	default:
		decoder = transform.NewWriter(r, unicode.UTF8.NewDecoder())
	}
	return decoder
}

func EncoderReader(types string, r io.Reader) *transform.Reader {
	var encoder *transform.Reader
	types = strings.ToUpper(types)
	switch types {
	case "UTF-8":
		encoder = transform.NewReader(r, unicode.UTF8.NewEncoder())
	case "GBK":
		encoder = transform.NewReader(r, simplifiedchinese.GBK.NewEncoder())
	case "BIG5":
		encoder = transform.NewReader(r, traditionalchinese.Big5.NewEncoder())
	case "SHIFTJIS":
		encoder = transform.NewReader(r, japanese.ShiftJIS.NewEncoder())
	case "KS_C_5601":
		encoder = transform.NewReader(r, korean.EUCKR.NewEncoder())
	case "GB2312":
		encoder = transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())
	case "GB18030":
		encoder = transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())
	case "BIG5-HKSCS":
		encoder = transform.NewReader(r, traditionalchinese.Big5.NewEncoder())
	case "UTF-16":
		encoder = transform.NewReader(r, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder())
	default:
		encoder = transform.NewReader(r, unicode.UTF8.NewEncoder())
	}
	return encoder
}

func EncoderWriter(types string, r io.Writer) *transform.Writer {
	var encoder *transform.Writer
	types = strings.ToUpper(types)
	switch types {
	case "UTF-8":
		encoder = transform.NewWriter(r, unicode.UTF8.NewEncoder())
	case "GBK":
		encoder = transform.NewWriter(r, simplifiedchinese.GBK.NewEncoder())
	case "BIG5":
		encoder = transform.NewWriter(r, traditionalchinese.Big5.NewEncoder())
	case "SHIFTJIS":
		encoder = transform.NewWriter(r, japanese.ShiftJIS.NewEncoder())
	case "KS_C_5601":
		encoder = transform.NewWriter(r, korean.EUCKR.NewEncoder())
	case "GB2312":
		encoder = transform.NewWriter(r, simplifiedchinese.GB18030.NewEncoder())
	case "GB18030":
		encoder = transform.NewWriter(r, simplifiedchinese.GB18030.NewEncoder())
	case "BIG5-HKSCS":
		encoder = transform.NewWriter(r, traditionalchinese.Big5.NewEncoder())
	case "UTF-16":
		encoder = transform.NewWriter(r, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder())
	default:
		encoder = transform.NewWriter(r, unicode.UTF8.NewEncoder())
	}
	return encoder
}
