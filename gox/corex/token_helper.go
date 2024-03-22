package corex

import "math/rand"

type TokenHelper struct {
	TokenLen int
	Dicts    string
}

func (t *TokenHelper) GenTokenDefault() string {
	return t.genTokenCommon(0, false)
}

func (t *TokenHelper) GenTokenDefaultNoZeroStart() string {
	return t.genTokenCommon(0, false)
}

func (t *TokenHelper) GenToken(token_len int) string {
	return t.genTokenCommon(token_len, false)
}

func (t *TokenHelper) GenTokenNoZeroStart(token_len int) string {
	return t.genTokenCommon(token_len, true)
}

func (t *TokenHelper) genTokenCommon(token_len int, noZeroStart bool) string {
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
	tokenresult := ""
	for i := 0; i < realLen; i++ {
		if i == 0 && noZeroStart {
			tokenresult = tokenresult + generateNoZeroStr(realDicts)
		} else {
			tmpIndx := rand.Intn(len(realDicts))
			tokenresult = tokenresult + realDicts[tmpIndx:tmpIndx+1]
		}
	}
	return tokenresult
}

// 生成首字符串不为0的随机字符串
func generateNoZeroStr(realDicts string) string {
	resultStr := ""
	for {
		tmpIndx := rand.Intn(len(realDicts))
		resultStr = realDicts[tmpIndx : tmpIndx+1]
		if resultStr == "0" {
			continue
		} else {
			break
		}
	}
	return resultStr
}
