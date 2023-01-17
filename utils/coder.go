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

// 先判断是否是UTF8再判断是否是其它编码才有意义
func isUtf8(data []byte) (bool, CoderType) {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false, T_UTF8
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false, T_UTF8
		}
	}
	return true, T_UTF8
}

func isGBK(data []byte) (bool, CoderType) {
	length := len(data)
	var i int = 0
	for i < length && i+1 < length {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0x7f {
				i += 2
				continue
			} else {
				return false, T_GBK
			}
		}
	}
	return true, T_GBK
}

func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}
