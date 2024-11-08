package corex

import (
	"unicode"
)

// isChinesePunctuation 判断字符是否是中文标点符号
func isChinesePunctuationByRune(ch rune) bool {
	// 中文标点符号的Unicode区间
	isPunctuation := (ch >= 0x3000 && ch <= 0x303F) || // 中文符号和标点
		(ch >= 0xFF00 && ch <= 0xFFEF) || // 全角拉丁字母、日文平假名、片假名、韩文兼容字母及全角符号
		unicode.Is(unicode.Sc, ch) || // 公共标点符号
		unicode.Is(unicode.Pd, ch) // 连接符号

	return isPunctuation
}

func IsChinesePunctuation(str string) bool {
	var count int
	for _, v := range str {
		if isChinesePunctuationByRune(v) {
			count++
		}
	}
	return count > 0
}

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
func IsChineseWithPunctuation(str string) bool {
	var count int
	for _, v := range str {
		// 标点符号规则1：(regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(v)))
		// 标点符号规则2：isChinesePunctuationByRune(v)
		if unicode.Is(unicode.Han, v) || isChinesePunctuationByRune(v) {
			count++
			break
		}
	}
	return count > 0
}

// 判断是全部是中文方法1
func IsAllChinese(str string) bool {
	var count int
	var lenRunes int
	for _, v := range str {
		lenRunes++
		if unicode.Is(unicode.Han, v) {
			count++
		}
	}
	return count == lenRunes
}

// 判断是全部是中文方法2
func IsAllChineseWithPunctuation(str string) bool {
	var count int
	var lenRunes int
	for _, v := range str {
		lenRunes++
		// 标点符号规则1：(regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(v)))
		// 标点符号规则2：isChinesePunctuationByRune(v)
		if unicode.Is(unicode.Han, v) || isChinesePunctuationByRune(v) {
			count++
		}
	}
	return count == lenRunes
}
