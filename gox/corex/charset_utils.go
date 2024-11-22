package corex

import (
	"strings"
	"unicode"
	"unicode/utf8"
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

// 判断字符是否含有中文标点符号，IsChinesePunctuation("你好，晴天")=true,IsChinesePunctuation("hello，world。")=true,IsChinesePunctuation("你好,晴天.")=false
func IsChinesePunctuation(str string) bool {
	var count int
	for _, v := range str {
		if isChinesePunctuationByRune(v) {
			count++
		}
	}
	return count > 0
}

// 判断字符是否含有中文字符，IsChinese("你好,晴天.")=true,IsChinese("hello，world。")=false,IsChinese("hello,world.")=false
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

// 判断字符是否含有中文字符或中文标点符号，IsChineseWithPunctuation("你好,晴天.")=true,IsChineseWithPunctuation("hello，world。")=true,IsChineseWithPunctuation("hello,world.")=false
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

// 判断字符串是否全部是中文字符，IsAllChinese("你好晴天")=true,IsAllChinese("你好，晴天。")=false,IsAllChinese("hello,world.")=false
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

// 判断字符串是否全部是中文字符或中文标点符号，IsAllChineseWithPunctuation("你好晴天")=true,IsAllChineseWithPunctuation("你好，晴天。")=true,IsAllChineseWithPunctuation("hello,world.")=false
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

// 字符串掩码函数，掩码数量至少字符串长度的向上取整的一半，可显示字符前面最大8字符，后面最大6字符，prettyMode:true的时候，过长字符掩码时候中间*会缩短为固定长度
// 如：MaskString("张三丰")="张**",MaskString("1234567890")="1234567890"
func MaskString(str string, prettyMode bool) string {
	lenRunes := utf8.RuneCountInString(str)
	countClear := lenRunes / 2
	countEnd := countClear / 2
	countStart := countClear - countEnd
	if countStart > 8 {
		countStart = 8
	}
	if countEnd > 6 {
		countEnd = 6
	}
	sb := strings.Builder{}
	i := 0
	for _, v := range str {
		if i < countStart || i >= lenRunes-countEnd {
			sb.WriteRune(v)
		} else if prettyMode {
			if i < countStart+8 || i >= lenRunes-countEnd-6 {
				sb.WriteString("*")
			}
		} else {
			sb.WriteString("*")
		}
		i++
	}
	return sb.String()
}
