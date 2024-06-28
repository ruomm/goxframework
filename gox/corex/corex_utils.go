package corex

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
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

// field名称简化
func FieldNameToSimply(fieldName string) string {
	lenFieldName := len(fieldName)
	if lenFieldName <= 0 {
		return fieldName
	}
	lastIndex := strings.LastIndex(fieldName, ".")
	if lastIndex >= 0 && lastIndex < lenFieldName-1 {
		return fieldName[lastIndex+1:]
	} else {
		return fieldName
	}
}

// 驼峰转下划线工具
func ToSnakeCase(str string) string {
	str = xvalid_matchNonAlphaNumeric.ReplaceAllString(str, "_")     //非常规字符转化为 _
	snake := xvalid_matchFirstCap.ReplaceAllString(str, "${1}_${2}") //拆分出连续大写
	snake = xvalid_matchAllCap.ReplaceAllString(snake, "${1}_${2}")  //拆分单词
	return strings.ToLower(snake)                                    //全部转小写
}
func JsonUnmarshal(str string, v any) error {
	if str == "" {
		return errors.New("json Unmarshal not support empty string")
	}
	err := json.Unmarshal([]byte(str), v)
	return err
}

func JsonMarshal(v any) (string, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	if len(jsonData) == 0 {
		return "", errors.New("json Marshal not support this object")
	}
	return string(jsonData), err
}

func JsonMarshalIndent(v any, prefix, indent string) (string, error) {
	jsonData, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return "", err
	}
	if len(jsonData) == 0 {
		return "", errors.New("json Marshal not support this object")
	}
	return string(jsonData), err
}

func JsonMarshalPretty(v any) (string, error) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	if len(jsonData) == 0 {
		return "", errors.New("json Marshal not support this object")
	}
	return string(jsonData), err
}

func GetLenForAny(i interface{}) int {
	if i == nil {
		return 0
	}
	vi := reflect.ValueOf(i)
	viKind := vi.Kind()
	//chan, func, interface, map, pointer, or slice
	if viKind == reflect.Slice || viKind == reflect.Pointer || viKind == reflect.Map || viKind == reflect.Interface || viKind == reflect.Func || viKind == reflect.Chan || viKind == reflect.Ptr {
		if vi.IsNil() {
			return 0
		} else {
			return vi.Len()
		}
	} else {
		return vi.Len()
	}
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	viKind := vi.Kind()
	//chan, func, interface, map, pointer, or slice
	if viKind == reflect.Slice || viKind == reflect.Pointer || viKind == reflect.Map || viKind == reflect.Interface || viKind == reflect.Func || viKind == reflect.Chan || viKind == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
