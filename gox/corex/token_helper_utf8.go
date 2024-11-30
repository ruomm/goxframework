package corex

import (
	"math/rand"
	"strings"
)

// 随机字符串结构体
type TokenHelperUtf8 struct {
	// 默认随机字符串token生成的长度
	TokenLen int
	// 随机字符串token生成的字典
	Dicts string
}

// 以默认长度生成随机字符串，默认长度为0时候生成4位随机字符串
func (t *TokenHelperUtf8) GenTokenDefault() string {
	return t.genTokenCommon(0, false, false)
}

// 以默认长度生成非0开头的随机字符串，默认长度为0时候生成4位随机字符串
func (t *TokenHelperUtf8) GenTokenDefaultNoZeroStart() string {
	return t.genTokenCommon(0, false, false)
}

// 以特定长度生成随机字符串，特定长度为0时候生成4位随机字符串
func (t *TokenHelperUtf8) GenToken(token_len int) string {
	return t.genTokenCommon(token_len, false, false)
}

// 以特定长度生成非0开头的随机字符串，特定长度为0时候生成4位随机字符串
func (t *TokenHelperUtf8) GenTokenNoZeroStart(token_len int) string {
	return t.genTokenCommon(token_len, true, false)
}

// 以特定长度生成非数字开头的随机字符串，特定长度为0时候生成4位随机字符串
func (t *TokenHelperUtf8) GenTokenNoNumberStart(token_len int) string {
	return t.genTokenCommon(token_len, true, true)
}

func (t *TokenHelperUtf8) genTokenCommon(token_len int, noZeroStart bool, noNumStart bool) string {
	realLen := 0
	if token_len > 0 {
		realLen = token_len
	} else if t.TokenLen > 0 {
		realLen = t.TokenLen
	} else {
		realLen = 4
	}
	realDicts := ""
	if len(t.Dicts) == 0 {
		realDicts = "0123456789"
	} else {
		realDicts = t.Dicts
	}
	utf8Len := Utf8Len(realDicts)
	tokenresult := ""
	for i := 0; i < realLen; i++ {
		if i == 0 {
			if noZeroStart || noNumStart {
				tokenresult = tokenresult + t.generateNoNumbnerStr(realDicts, noZeroStart, noNumStart)
			} else {
				tmpIndx := rand.Intn(utf8Len)
				tokenresult = tokenresult + Utf8At(realDicts, tmpIndx)
			}
		} else {
			tmpIndx := rand.Intn(utf8Len)
			tokenresult = tokenresult + Utf8At(realDicts, tmpIndx)
		}
	}
	return tokenresult
}

// 生成首字符串不为数字或不为0的随机字符串
func (t *TokenHelperUtf8) generateNoNumbnerStr(realDicts string, noZeroStart bool, noNumStart bool) string {
	sb := strings.Builder{}
	for _, r := range realDicts {
		if noNumStart && r >= 48 && r <= 57 {
			continue
		} else if noZeroStart && r == 48 {
			continue
		} else {
			sb.WriteRune(r)
		}
	}
	letterDicts := sb.String()
	letterLen := Utf8Len(letterDicts)
	if letterLen > 0 {
		tmpIndx := rand.Intn(letterLen)
		resultStr := Utf8At(letterDicts, tmpIndx)
		return resultStr
	} else {
		utf8Len := Utf8Len(realDicts)
		tmpIndx := rand.Intn(utf8Len)
		resultStr := Utf8At(realDicts, tmpIndx)
		return resultStr
	}
}
