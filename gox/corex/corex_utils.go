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

// 驼峰转蛇形工具
func ToSnakeCase(str string) string {
	str = xvalid_matchNonAlphaNumeric.ReplaceAllString(str, "_")     //非常规字符转化为 _
	snake := xvalid_matchFirstCap.ReplaceAllString(str, "${1}_${2}") //拆分出连续大写
	snake = xvalid_matchAllCap.ReplaceAllString(snake, "${1}_${2}")  //拆分单词
	return strings.ToLower(snake)                                    //全部转小写
}

// 蛇形转驼峰工具
func ToCamelCase(str string) string {
	sb := strings.Builder{}
	upFlag := false
	for i := 0; i < len(str); i++ {
		if str[i] == '_' {
			upFlag = true
			continue
		} else {
			if upFlag {
				upFlag = false
				sb.WriteString(strings.ToUpper(str[i : i+1]))
			} else {
				sb.WriteString(str[i : i+1])
			}
		}
	}
	return sb.String()
}

// 首字母大小
func FirstLetterToUpper(str string) string {
	length := len(str)
	if length == 0 {
		return str
	} else if length == 1 {
		return strings.ToUpper(str)
	} else {
		return strings.ToUpper(str[0:1]) + str[1:length]
	}
}

// 首字母小写
func FirstLetterToLower(str string) string {
	length := len(str)
	if length == 0 {
		return str
	} else if length == 1 {
		return strings.ToLower(str)
	} else {
		return strings.ToLower(str[0:1]) + str[1:length]
	}
}

func JsonUnmarshal(str string, v any, nodes ...string) error {
	if str == "" {
		return errors.New("json Unmarshal not support empty string")
	}
	err := JsonUnmarshalByBytes([]byte(str), v, nodes...)
	return err
}

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

// 解析Http强求range头的第一个请求范围
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
