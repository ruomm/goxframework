package corex

import "math/rand"

type TokenHelper struct {
	TokenLen int
	Dicts    string
}

func (t *TokenHelper) GenTokenDefault() string {
	return t.GenToken(0)
}
func (t *TokenHelper) GenToken(token_len int) string {
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
		tmpIndx := rand.Intn(len(realDicts))
		tokenresult = tokenresult + realDicts[tmpIndx:tmpIndx+1]
	}
	return tokenresult
}
