package corex

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	//TIME_LOCATION_CN, _         = time.LoadLocation("Asia/Shanghai")
	xvalid_matchNonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	xvalid_matchFirstCap        = regexp.MustCompile("(.)([A-Z][a-z]+)")
	xvalid_matchAllCap          = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ServerTrace(tracePort int) {
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(tracePort), nil)
}

func GetMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func GetMd5WithSlat(data, slat string) string {
	var realSlat string
	h := md5.New()
	h.Write([]byte(data + realSlat))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

// 驼峰转下划线工具
func ToSnakeCase(str string) string {
	str = xvalid_matchNonAlphaNumeric.ReplaceAllString(str, "_")     //非常规字符转化为 _
	snake := xvalid_matchFirstCap.ReplaceAllString(str, "${1}_${2}") //拆分出连续大写
	snake = xvalid_matchAllCap.ReplaceAllString(snake, "${1}_${2}")  //拆分单词
	return strings.ToLower(snake)                                    //全部转小写
}
func JsonParseByString(str string, v any) error {
	if str == "" {
		return errors.New("json Unmarshal not support empty string")
	}
	err := json.Unmarshal([]byte(str), v)
	return err
}

func JsonFormatByString(v any) (string, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	if len(jsonData) == 0 {
		return "", errors.New("json Marshal not support this object")
	}
	return string(jsonData), err
}
