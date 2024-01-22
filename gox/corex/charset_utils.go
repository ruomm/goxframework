package corex

import (
	"regexp"
	"unicode"
)

// 判断是否含有中文方法1
func IsChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

// 判断是否含有中文方法2
func IsChineseChar(str string) bool {
	var count int
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			count++
			break
		}
	}
	return count > 0
}

// 判断是全部是中文方法1
func IsAllChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
		}
	}
	return count == len(str)
}

// 判断是全部是中文方法2
func IsAllChineseChar(str string) bool {
	var count int
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			count++
		}
	}
	return count == len(str)
}
