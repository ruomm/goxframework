package corex

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
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

// 获取字符串MD5值
func GetMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

// 获取字符串结尾拼接盐值的MD5值
func GetMd5WithSlat(data, slat string) string {
	//var realSlat string
	h := md5.New()
	h.Write([]byte(data + slat))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

// field完整名称简化为字段名称,FieldNameToSimply("vo.UserInfo.UserName")="UserName"
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

// 驼峰转蛇形工具,ToSnakeCase("DemoUserInfo") = "demo_user_info"
func ToSnakeCase(str string) string {
	str = xvalid_matchNonAlphaNumeric.ReplaceAllString(str, "_")     //非常规字符转化为 _
	snake := xvalid_matchFirstCap.ReplaceAllString(str, "${1}_${2}") //拆分出连续大写
	snake = xvalid_matchAllCap.ReplaceAllString(snake, "${1}_${2}")  //拆分单词
	return strings.ToLower(snake)                                    //全部转小写
}

// 蛇形转驼峰工具,ToCamelCase("demo_user_info") = "demoUserInfo"
func ToCamelCase(str string) string {
	sb := strings.Builder{}
	upFlag := false
	for _, v := range str {
		if v == '_' {
			upFlag = true
			continue
		} else {
			if upFlag {
				upFlag = false
				sb.WriteString(strings.ToUpper(string(v)))
			} else {
				sb.WriteString(string(v))
			}
		}
	}
	return sb.String()
}

// 首字母大小
func FirstLetterToUpper(str string) string {
	lenRunes := utf8.RuneCountInString(str)
	if lenRunes <= 0 {
		return str
	} else if lenRunes == 1 {
		return strings.ToUpper(str)
	} else {
		firstLetter := ""
		secondIndex := -1
		for i, r := range str {
			if i == 0 {
				firstLetter = strings.ToUpper(string(r))
			} else {
				secondIndex = i
				break
			}
		}
		return firstLetter + str[secondIndex:]
	}
}

// 首字母小写
func FirstLetterToLower(str string) string {
	lenRunes := utf8.RuneCountInString(str)
	if lenRunes <= 0 {
		return str
	} else if lenRunes == 1 {
		return strings.ToLower(str)
	} else {
		firstLetter := ""
		secondIndex := -1
		for i, r := range str {
			if i == 0 {
				firstLetter = strings.ToLower(string(r))
			} else {
				secondIndex = i
				break
			}
		}
		return firstLetter + str[secondIndex:]
	}
}

// 判断字符串str是否JSON字符串，trimMode:true时候，字符串str经过trim后判断
func IsJsonString(str string, trimMode bool) bool {
	json := ""
	if trimMode {
		json = strings.TrimSpace(str)
	} else {
		json = str
	}
	if len(json) <= 0 {
		return false
	} else if strings.HasPrefix(json, "{") && strings.HasSuffix(json, "}") {
		return true
	} else if strings.HasPrefix(json, "[") && strings.HasSuffix(json, "]") {
		return true
	} else {
		return false
	}
}

// JSON解析，nodes有值会解析节点的值，如JsonUnmarshal(str, v, data, user)会解析str字符串的data.user节点
func JsonUnmarshal(str string, v any, nodes ...string) error {
	if str == "" {
		return errors.New("json Unmarshal not support empty string")
	}
	err := JsonUnmarshalByBytes([]byte(str), v, nodes...)
	return err
}

// JSON解析，nodes有值会解析节点的值，如JsonUnmarshal(jsonByte, v, data, user)会解析jsonByte字节的data.user节点
func JsonUnmarshalByBytes(jsonByte []byte, v any, nodes ...string) error {
	if len(jsonByte) <= 0 {
		return errors.New("json Unmarshal not support empty byte")
	}
	if len(nodes) <= 0 {
		err := json.Unmarshal(jsonByte, v)
		return err
	}
	var mapResult map[string]interface{}
	err := json.Unmarshal(jsonByte, &mapResult)
	var nodeMapResult map[string]interface{}
	nodeMapResult = mapResult
	var realResult interface{} = nil
	for i, node := range nodes {
		tmpResult, tmpResultOk := nodeMapResult[node]
		if !tmpResultOk {
			nodeMapResult = nil
			break
		}
		if i == len(nodes)-1 {
			// 结果结束
			realResult = tmpResult
			break
		}
		tmpMap, tmpMapOk := tmpResult.(map[string]interface{})
		if !tmpMapOk {
			nodeMapResult = nil
			return errors.New(fmt.Sprintf("json Unmarshal %s node fail,unable convert to map[string]interface{}", node))
		} else {
			nodeMapResult = tmpMap
		}
	}
	if realResult == nil {
		return nil
	}
	jsonData, err := json.Marshal(realResult)
	if err != nil {
		return err
	}
	if len(jsonData) == 0 {
		return errors.New("json Marshal not support this object")
	}
	err = json.Unmarshal(jsonData, v)
	return err
}

// JSON编码为字符串
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

// JSON编码为字符串，prefix、indent设置编码的美化方式
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

// JSON编码为字符串，默认换行和缩进进行美化
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

// 获取任意对象的长度
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

// 判断任意对象是否为空
func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	viKind := vi.Kind()
	//chan, func, interface, map, pointer, or slice
	if viKind == reflect.Slice || viKind == reflect.Pointer || viKind == reflect.Map || viKind == reflect.Interface || viKind == reflect.Func || viKind == reflect.Chan || viKind == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

// 解析Http强求range头的第一个请求范围，主要是分片上传和分片下载时候使用
func ParseHttpRangeFirst(rangeHeader string) (int64, int64) {
	if len(rangeHeader) <= 0 {
		return -1, -1
	}
	rangePrefix := "bytes="
	rangeLower := strings.ToLower(rangeHeader)
	if strings.HasPrefix(rangeLower, rangePrefix) && len(rangeLower) > len(rangePrefix) {
		rangeStr := rangeLower[len(rangePrefix):]
		rangeByDh := StringToSlice(rangeStr, ",", false)
		if len(rangeByDh) <= 0 {
			return -1, -1
		}
		rangeAreas := StringToSlice(rangeByDh[0], "-", true)
		if len(rangeAreas) == 1 {
			start, err := strconv.ParseInt(rangeAreas[0], 10, 64)
			if err != nil {
				return -1, -1
			} else {
				return start, -1
			}
		} else if len(rangeAreas) == 2 {
			start, errStart := strconv.ParseInt(rangeAreas[0], 10, 64)
			end, errEnd := strconv.ParseInt(rangeAreas[1], 10, 64)
			if errStart != nil && errEnd != nil {
				return -1, -1
			} else if errStart != nil {
				return -1, end
			} else if errEnd != nil {
				return start, -1
			} else if end >= start {
				return start, end
			} else {
				return -1, -1
			}
		} else {
			return -1, -1
		}
	} else {
		return -1, -1
	}
}
