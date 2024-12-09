package corex

import (
	"math/rand"
	"strings"
)

// 随机字符串结构体-支持中文
type TokenHelper struct {
	// 默认随机字符串token生成的长度
	TokenLen int
	// 随机字符串token生成的字典
	Dicts string
}

// 以默认长度生成随机字符串，默认长度为0时候生成4位随机字符串
func (t *TokenHelper) GenTokenDefault() string {
	return t.genTokenCommon(0, false, false)
}

// 以默认长度生成非0开头的随机字符串，默认长度为0时候生成4位随机字符串
func (t *TokenHelper) GenTokenDefaultNoZeroStart() string {
	return t.genTokenCommon(0, false, false)
}

// 以特定长度生成随机字符串，特定长度为0时候生成4位随机字符串
func (t *TokenHelper) GenToken(token_len int) string {
	return t.genTokenCommon(token_len, false, false)
}

// 以特定长度生成非0开头的随机字符串，特定长度为0时候生成4位随机字符串
func (t *TokenHelper) GenTokenNoZeroStart(token_len int) string {
	return t.genTokenCommon(token_len, true, false)
}

// 以特定长度生成非数字开头的随机字符串，特定长度为0时候生成4位随机字符串
func (t *TokenHelper) GenTokenNoNumberStart(token_len int) string {
	return t.genTokenCommon(token_len, true, true)
}

func (t *TokenHelper) genTokenCommon(token_len int, noZeroStart bool, noNumStart bool) string {
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
	dictsLen := len(realDicts)
	tokenresult := ""
	for i := 0; i < realLen; i++ {
		if i == 0 {
			if noZeroStart || noNumStart {
				tokenresult = tokenresult + t.generateNoNumbnerStr(realDicts, noZeroStart, noNumStart)
			} else {
				tmpIndx := rand.Intn(dictsLen)
				tokenresult = tokenresult + realDicts[tmpIndx:tmpIndx+1]
			}
		} else {
			tmpIndx := rand.Intn(dictsLen)
			tokenresult = tokenresult + realDicts[tmpIndx:tmpIndx+1]
		}
	}
	return tokenresult
}

// 生成首字符串不为数字或不为0的随机字符串
func (t *TokenHelper) generateNoNumbnerStr(realDicts string, noZeroStart bool, noNumStart bool) string {
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
	letterLen := len(letterDicts)
	if letterLen > 0 {
		tmpIndx := rand.Intn(letterLen)
		resultStr := letterDicts[tmpIndx : tmpIndx+1]
		return resultStr
	} else {
		dictsLen := len(realDicts)
		tmpIndx := rand.Intn(dictsLen)
		resultStr := realDicts[tmpIndx : tmpIndx+1]
		return resultStr
	}
}
