package utils

import (
	"io"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type CoderType int

const (
	T_Auto CoderType = iota
	T_UTF8
	T_GBK
	T_Big5
	T_ShiftJIS
	T_EUCKR
	T_GB18030
	T_UTF16_L
	T_UTF16_B
)

var chcp = map[CoderType]string{
	T_UTF8: "65001", T_Auto: "65001",
	T_UTF16_L: "1200", T_UTF16_B: "1200",
	T_GBK:      "936",
	T_GB18030:  "54936",
	T_Big5:     "950",
	T_EUCKR:    "949",
	T_ShiftJIS: "932",
}

func CodePage(types CoderType) string {
	if cp, ok := chcp[types]; ok {
		return cp
	} else {
		return "65001"
	}
}

func CoderToType(types string) CoderType {
	types = strings.ToUpper(types)
	switch types {
	case "GBK":
		return T_GBK
	case "BIG5", "BIG5-HKSCS":
		return T_Big5
	case "SHIFTJIS":
		return T_ShiftJIS
	case "KS_C_5601":
		return T_EUCKR
	case "GB18030", "GB2312":
		return T_GB18030
	case "UTF-16", "UTF-16-L":
		return T_UTF16_L
	case "UTF-16-B":
		return T_UTF16_B
	case "AUTO":
		return T_Auto
	default:
		return T_UTF8
	}
}

func DecoderReader(types CoderType, r io.Reader) *transform.Reader {
	return transform.NewReader(r, newDeCoder(types))
}

func DecoderWriter(types CoderType, r io.Writer) *transform.Writer {
	return transform.NewWriter(r, newDeCoder(types))
}

func EncoderReader(types CoderType, r io.Reader) *transform.Reader {
	return transform.NewReader(r, newEeCoder(types))
}

func EncoderWriter(types CoderType, r io.Writer) *transform.Writer {
	return transform.NewWriter(r, newEeCoder(types))
}

func newDeCoder(coder CoderType) *encoding.Decoder {
	var decoder *encoding.Decoder
	switch coder {
	case T_GBK:
		decoder = simplifiedchinese.GBK.NewDecoder()
	case T_Big5:
		decoder = traditionalchinese.Big5.NewDecoder()
	case T_ShiftJIS:
		decoder = japanese.ShiftJIS.NewDecoder()
	case T_EUCKR:
		decoder = korean.EUCKR.NewDecoder()
	case T_GB18030:
		decoder = simplifiedchinese.GB18030.NewDecoder()
	case T_UTF16_L:
		decoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	case T_UTF16_B:
		decoder = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
	default:
		decoder = unicode.UTF8.NewDecoder()
	}
	return decoder
}

func newEeCoder(coder CoderType) *encoding.Encoder {
	var decoder *encoding.Encoder
	switch coder {
	case T_GBK:
		decoder = simplifiedchinese.GBK.NewEncoder()
	case T_Big5:
		decoder = traditionalchinese.Big5.NewEncoder()
	case T_ShiftJIS:
		decoder = japanese.ShiftJIS.NewEncoder()
	case T_EUCKR:
		decoder = korean.EUCKR.NewEncoder()
	case T_GB18030:
		decoder = simplifiedchinese.GB18030.NewEncoder()
	case T_UTF16_L:
		decoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	case T_UTF16_B:
		decoder = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()
	default:
		decoder = unicode.UTF8.NewEncoder()
	}
	return decoder
}
