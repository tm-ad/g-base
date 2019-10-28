package util

import (
	"strings"
)

// SubStr 截取字符串 start 起点下标 length 需要截取的长度
func SubStr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

// SubString 截取字符串 start 起点下标 end 终点下标(不包括)
func SubString(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func ToWebSafeBase64(str string) string {
	// 替换 URL 不安全的 + 到 -，/ 到 _
	// 去掉多余的 padding characters
	str = strings.ReplaceAll(str, `+`, `-`)
	str = strings.ReplaceAll(str, `/`, `_`)
	// 从末尾去掉 =
	str = strings.TrimRightFunc(str, func(r rune) bool {
		if string(r) == "=" {
			return true
		}
		return false
	})

	return str
}

func FromWebSafeBase64(str string) string {
	// 替换 URL 安全的 - 到 +，_ 到 /
	str = strings.ReplaceAll(str, `-`, `+`)
	str = strings.ReplaceAll(str, `_`, `/`)
	// 补齐 =
	missingPaddingChars := 0
	switch len(str) % 4 {
	case 3:
		missingPaddingChars = 1
		break
	case 2:
		missingPaddingChars = 2
		break
	case 0:
		missingPaddingChars = 0
	default:
		panic(`invalid web safe base64 format`)
	}
	for i := 0; i < missingPaddingChars; i++ {
		str = str + `=`
	}

	return str
}
