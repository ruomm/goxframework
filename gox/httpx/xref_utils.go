package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/morrisxyang/xreflect"
	"github.com/ruomm/goxframework/gox/corex"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const xReflect_log = false

// const xReflect_tag_cp_opt = "xref_opt"
// const xReflect_tag_cp_src = "xref"
const xReflect_tag_cp_xurlget = "xurl_param"

//var xReflect_location, _ = time.LoadLocation("Asia/Shanghai")

const xReflect_time_layout = "2006-01-02 15:04:05"

const xReflect_key_time_t = "t"
const xReflect_key_time_tf = "tf"
const xReflect_key_number_point = "p"

// 如是omitempty参数存在，来源的数字类型的0、bool类型的false、字符串类型的空、时间类型的0或负数不会赋值的目标里面
const xReflect_key_tidy = "tidy"

/*
*
srcO：源结构体
destO：目标切片，不可以传入结构体
*/
// TransferObj 将srcO对象的属性值转成destO对象的属性值，属性对应关系和控制指令通过`xref`标签指定
// 无标签的如果再按属性名匹配
func ParseToJSONByte(jsonObj interface{}, transBasePointer ...bool) ([]byte, error) {
	if nil == jsonObj {
		return nil, errors.New("ParseToJSONByte error,source interface is nil")
	}
	srcTypeOf := reflect.TypeOf(jsonObj)
	srcName := srcTypeOf.String()
	// 目标是字符串类型
	if srcName == "string" {
		return []byte(jsonObj.(string)), nil
	} else {
		jsonStr, err := json.Marshal(jsonObj)
		return jsonStr, err
	}
}

/*
*
srcO：源结构体
destO：目标切片，不可以传入结构体
*/
// TransferObj 将srcO对象的属性值转成destO对象的属性值，属性对应关系和控制指令通过`xref`标签指定
// 无标签的如果再按属性名匹配
func ParseToUrlEncodeString(srcO interface{}, transBasePointer ...bool) (string, error) {
	if nil == srcO {
		//return "", errors.New("XReflectCopy error,source interface is nil")
		return "", nil
	}
	srcTypeOf := reflect.TypeOf(srcO)
	srcName := srcTypeOf.String()
	// 目标是字符串类型
	if srcName == "string" {
		return srcO.(string), nil
	} else if srcName == "url.Values" {
		v := srcO.(url.Values)
		if len(v) <= 0 {
			return "", nil
		} else {
			return v.Encode(), nil
		}
	} else if srcName == "*url.Values" {
		v := srcO.(*url.Values)
		if len(*v) <= 0 {
			return "", nil
		} else {
			return v.Encode(), nil
		}
	}

	var transBasePointerFlag = false
	if nil != transBasePointer && len(transBasePointer) > 0 && transBasePointer[0] {
		transBasePointerFlag = transBasePointer[0]
	}
	resOpt := make(map[string]string)
	resSrc := make(map[string]string)
	reflectValueMap, err := xreflect.SelectFieldsDeep(srcO, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagXreft, okXreft := field.Tag.Lookup(xReflect_tag_cp_xurlget)
		if !okXreft {
			return false
		}
		// 开始分割目标控制和属性控制
		subTags := corex.ParseToSubTag(tagXreft)
		// 解析目标控制
		urlKey := ""
		if len(subTags) > 0 {
			urlKey = subTags[0]
		}
		if urlKey == "-" {
			return false
		}
		if urlKey == "" {
			urlKey = strings.ToLower(s[0:1]) + s[1:len(s)]
		}
		resSrc[s] = urlKey
		// 解析属性控制
		tagOpt := ""
		if len(subTags) > 1 {
			tagOpt = subTags[1]
		}
		resOpt[s] = tagOpt
		if xReflect_log {
			fmt.Println("解析URL参数字段，目标：" + urlKey + "，来源：" + s + "，控制协议：" + tagOpt)
		}
		return true
	})
	if err != nil {
		return "", err
	}
	v := url.Values{}
	for key, _ := range reflectValueMap {
		var srcKey string
		if resSrc[key] != "" {
			srcKey = resSrc[key]
		} else {
			srcKey = key
		}
		srcValue, err := xreflect.EmbedFieldValue(srcO, key)
		if err != nil {
			continue
		}
		if transBasePointerFlag {
			srcValue = transBasePointerToValue(srcValue)
		}
		if srcValue == nil {
			continue
		}
		cpOpt := resOpt[key]
		rtVal := xReflect_transSrcToDestString(key, cpOpt, srcValue)
		if rtVal == nil {
			continue
		}
		v.Set(srcKey, rtVal.(string))
	}
	if len(v) <= 0 {
		return "", nil
	} else {
		return v.Encode(), nil
	}

}

// 解析来源字段为目标待赋值字段
func xReflect_transSrcToDestString(key string, cpOpt string, srcValue interface{}) interface{} {
	srcTypeOf := reflect.TypeOf(srcValue)
	srcKind := srcTypeOf.Kind()
	srcType := srcTypeOf.String()
	isTidy := xReflect_tagContainKey(cpOpt, xReflect_key_tidy)
	// 目标是字符串类型
	if xReflect_isNumberType(srcKind, true) {
		rtVal := xReflect_transToInt64(srcValue, srcKind, srcType, isTidy)
		if rtVal == nil {
			if xReflect_log {
				fmt.Println(key + "字段无法赋值，转换错误，int64->string")
			}
			return nil
		} else {
			return strconv.FormatInt(rtVal.(int64), 10)
		}
	} else if xReflect_isFloatType(srcKind) {
		rtVal := xReflect_transToFloat64(srcValue, srcKind, srcType, isTidy)
		if rtVal == nil {
			if xReflect_log {
				fmt.Println(key + "字段无法赋值，转换错误，float64->string")
			}
			return nil
		} else {
			optStr1 := xReflect_findTagValueByKey(cpOpt, xReflect_key_number_point)
			prec := -1
			if len(optStr1) > 0 {
				prec64, error := strconv.ParseInt(optStr1, 10, 64)
				if error != nil {
					prec = -1
				} else if prec64 >= -1 && prec64 <= 10 {
					prec = int(prec64)
				} else {
					prec = -1
				}
			}
			return strconv.FormatFloat(rtVal.(float64), 'f', prec, 64)
		}

	} else if xReflect_isBoolType(srcKind) {
		// 来源数字类型-转换赋值
		rtVal := xReflect_transToBool(srcValue, srcKind, srcType, isTidy)
		if rtVal == nil {
			if xReflect_log {
				fmt.Println(key + "字段无法赋值，转换错误，bool->string")
			}
			return nil
		} else {
			if rtVal.(bool) {
				return "true"
			} else {
				return "false"
			}
		}
	} else if xReflect_isStringType(srcKind) {
		if isTidy {
			rtVal := srcValue.(string)
			if len(rtVal) <= 0 {
				return nil
			} else {
				return srcValue
			}
		} else {
			return srcValue
		}

	} else if xReflect_isTimeTypeByName(srcType) {
		var srcConv time.Time
		if xReflect_isPointor(srcType) {
			timePointor := srcValue.(*time.Time)
			if nil == timePointor {
				return nil
			}
			srcConv = *timePointor
		} else {
			srcConv = srcValue.(time.Time)
		}
		optStr1 := xReflect_findTagValueByKey(cpOpt, xReflect_key_time_tf)
		rtVal := xReflect_formatTimeToString(&srcConv, optStr1, isTidy)
		if len(rtVal) <= 0 {
			if xReflect_log {
				fmt.Println(key + "字段无法赋值，转换错误，time.Time->string")
			}
			return nil
		} else {
			return rtVal
		}
	} else {
		if xReflect_log {
			fmt.Println(key + "字段无法赋值，非数字、字符串、时间类型无法转换成字符串类型")
		}
		return nil
	}
}
func xReflect_isNumberType(kind reflect.Kind, onlyNatural bool) bool {
	if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 ||
		kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
		return true
	} else if !onlyNatural {
		if kind == reflect.Float32 || kind == reflect.Float64 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
func xReflect_isFloatType(kind reflect.Kind) bool {
	if kind == reflect.Float64 || kind == reflect.Float32 {
		return true
	} else {
		return false
	}
}

func xReflect_isBoolType(kind reflect.Kind) bool {
	if kind == reflect.Bool {
		return true
	} else {
		return false
	}
}

func xReflect_isStringType(kind reflect.Kind) bool {
	return kind == reflect.String
}

func xReflect_tagContainKey(tagValue string, key string) bool {
	tagsOptions := corex.ParseTagToOptions(tagValue)
	if len(tagsOptions) == 0 {
		return false
	}
	keyContain := false
	for _, tmpOption := range tagsOptions {
		if tmpOption.Contains(key) {
			keyContain = true
			break
		}
	}
	return keyContain
}

func xReflect_findTagValueByKey(tagValue string, key string) string {
	tagsOptions := corex.ParseTagToOptions(tagValue)
	if len(tagsOptions) == 0 {
		return ""
	}
	var keyVal string
	for _, tmpOption := range tagsOptions {
		if tmpOption.Contains(key) {
			keyVal = tmpOption.OptionValue(key)
		}
	}
	return keyVal
}
func xReflect_isTimeTypeByName(typeName string) bool {
	return typeName == "time.Time" || typeName == "Time" || typeName == "*time.Time"
}

func xReflect_isPointor(typeName string) bool {
	return strings.HasPrefix(typeName, "*")
}

// 转换各种数字类型为int64，浮点型进行math.Round，时间类型取得毫秒时间戳
func xReflect_transToInt64(srcFieldV interface{}, srcKind reflect.Kind, srcFieldVT string, isTidy bool) interface{} {
	var vi interface{} = nil
	if srcKind == reflect.Int {
		if srcFieldVT == "int" {
			vi = int64(srcFieldV.(int))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Int8 {
		if srcFieldVT == "int8" {
			vi = int64(srcFieldV.(int8))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Int16 {
		if srcFieldVT == "int16" {
			vi = int64(srcFieldV.(int16))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Int32 {
		if srcFieldVT == "int32" {
			vi = int64(srcFieldV.(int32))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Int64 {
		if srcFieldVT == "int64" {
			vi = int64(srcFieldV.(int64))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Uint {
		if srcFieldVT == "uint" {
			vi = int64(srcFieldV.(uint))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Uint8 {
		if srcFieldVT == "uint8" {
			vi = int64(srcFieldV.(uint8))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Uint16 {
		if srcFieldVT == "uint16" {
			vi = int64(srcFieldV.(uint16))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Uint32 {
		if srcFieldVT == "uint32" {
			vi = int64(srcFieldV.(uint32))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Uint64 {
		if srcFieldVT == "uint64" {
			vi = int64(srcFieldV.(uint64))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Float64 {
		if srcFieldVT == "float64" {
			vi = int64(math.Round(srcFieldV.(float64)))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Float32 {
		if srcFieldVT == "float32" {
			vi = int64(math.Round(float64(srcFieldV.(float32))))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Bool {
		if srcFieldVT == "bool" {
			if srcFieldV.(bool) {
				vi = int64(1)
			} else {
				vi = int64(0)
			}
		} else {
			vi = nil
		}
	} else if "Time" == srcFieldVT {
		t := srcFieldV.(time.Time)
		vi = t.UnixMilli()
	}
	if vi == nil {
		return vi
	} else if isTidy {
		if reflect.TypeOf(vi).Name() == "int64" {
			if vi.(int64) == 0 {
				return nil
			} else {
				return vi
			}
		} else {
			return vi
		}

	} else {
		return vi
	}

}

// 转换各种数字类型为int64，浮点型进行math.Round，时间类型取得毫秒时间戳
func xReflect_transToBool(srcFieldV interface{}, srcKind reflect.Kind, srcFieldVT string, isTidy bool) interface{} {
	var vi interface{} = nil
	if srcKind == reflect.Int {
		if srcFieldVT == "int" {
			if int64(srcFieldV.(int)) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Int8 {
		if srcFieldVT == "int8" {
			if int64(srcFieldV.(int8)) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Int16 {
		if srcFieldVT == "int16" {
			if int64(srcFieldV.(int16)) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Int32 {
		if srcFieldVT == "int32" {
			if int64(srcFieldV.(int32)) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Int64 {
		if srcFieldVT == "int64" {
			if int64(srcFieldV.(int64)) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Uint {
		if srcFieldVT == "uint" {
			if int64(srcFieldV.(uint)) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Uint8 {
		if srcFieldVT == "uint8" {
			if int64(srcFieldV.(uint8)) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Uint16 {
		if srcFieldVT == "uint16" {
			if int64(srcFieldV.(uint16)) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Uint32 {
		if srcFieldVT == "uint32" {
			if int64(srcFieldV.(uint32)) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Uint64 {
		if srcFieldVT == "uint64" {
			if int64(srcFieldV.(uint64)) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Float64 {
		if srcFieldVT == "float64" {
			if int64(math.Round(srcFieldV.(float64))) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Float32 {
		if srcFieldVT == "float32" {
			if int64(math.Round(float64(srcFieldV.(float32)))) > 0 {
				vi = true
			} else {
				vi = false
			}
		}
	} else if srcKind == reflect.Bool {
		if srcFieldVT == "bool" {
			vi = srcFieldV.(bool)
		}
	}
	if vi == nil {
		return vi
	} else if isTidy {
		if reflect.TypeOf(vi).Name() == "bool" {
			if !vi.(bool) {
				return nil
			} else {
				return vi
			}
		} else {
			return vi
		}

	} else {
		return vi
	}

}

// 转换各种数字类型为float64，浮点型进行math.Round，时间类型取得毫秒时间戳
func xReflect_transToFloat64(srcFieldV interface{}, srcKind reflect.Kind, srcFieldVT string, isTidy bool) interface{} {
	var vi interface{} = nil
	if srcKind == reflect.Int {
		if srcFieldVT == "int" {
			vi = float64(srcFieldV.(int))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Int8 {
		if srcFieldVT == "int8" {
			vi = float64(srcFieldV.(int8))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Int16 {
		if srcFieldVT == "int16" {
			vi = float64(srcFieldV.(int16))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Int32 {
		if srcFieldVT == "int32" {
			vi = float64(srcFieldV.(int32))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Int64 {
		if srcFieldVT == "int64" {
			vi = float64(srcFieldV.(int64))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Uint {
		if srcFieldVT == "uint" {
			vi = float64(srcFieldV.(uint))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Uint8 {
		if srcFieldVT == "uint8" {
			vi = float64(srcFieldV.(uint8))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Uint16 {
		if srcFieldVT == "uint16" {
			vi = float64(srcFieldV.(uint16))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Uint32 {
		if srcFieldVT == "uint32" {
			vi = float64(srcFieldV.(uint32))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Uint64 {
		if srcFieldVT == "uint64" {
			vi = float64(srcFieldV.(uint64))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Float64 {
		if srcFieldVT == "float64" {
			vi = float64(srcFieldV.(float64))
		} else {
			vi = srcFieldV
		}
	} else if srcKind == reflect.Bool {
		if srcFieldVT == "bool" {
			if srcFieldV.(bool) {
				vi = float64(1)
			} else {
				vi = float64(0)
			}
		} else {
			vi = nil
		}
	} else if srcKind == reflect.Float32 {
		if srcFieldVT == "float32" {
			vi = float64(srcFieldV.(float32))
		} else {
			vi = srcFieldV
		}
	} else if "Time" == srcFieldVT {
		t := srcFieldV.(time.Time)
		vi = float64(t.UnixMilli())
	}
	if vi == nil {
		return vi
	} else if isTidy {
		if vi.(float64) >= -0.0000000001 && vi.(float64) <= 0.0000000001 {
			return nil
		} else {
			return vi
		}
	} else {
		return vi
	}

}

// 格式化时间为字符串
func xReflect_formatTimeToString(t *time.Time, timeLayout string, isTidy bool) string {
	if isTidy && t.UnixMilli() <= 0 {
		return ""
	}
	//"America/Adak" "Asia/Shanghai"
	var realTimeLayout string
	if len(timeLayout) > 0 {
		realTimeLayout = timeLayout
	} else {
		realTimeLayout = xReflect_time_layout
	}
	return t.In(corex.ToTimeLocation()).Format(realTimeLayout)
}

func transBasePointerToValue(srcValue interface{}) interface{} {
	if nil == srcValue {
		return nil
	}
	srcTypeOf := reflect.TypeOf(srcValue)
	srcType := srcTypeOf.String()
	if srcType == "*int" {
		conValP := srcValue.(*int)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*int8" {
		conValP := srcValue.(*int8)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*int16" {
		conValP := srcValue.(*int16)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*int32" {
		conValP := srcValue.(*int32)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*int64" {
		conValP := srcValue.(*int64)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*uint" {
		conValP := srcValue.(*uint)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*uint8" {
		conValP := srcValue.(*uint8)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*uint16" {
		conValP := srcValue.(*uint16)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*uint32" {
		conValP := srcValue.(*uint32)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*uint64" {
		conValP := srcValue.(*uint64)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*string" {
		conValP := srcValue.(*string)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*float32" {
		conValP := srcValue.(*float32)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*float64" {
		conValP := srcValue.(*float64)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else if srcType == "*bool" {
		conValP := srcValue.(*bool)
		if nil == conValP {
			return nil
		} else {
			var conVal = *conValP
			return conVal
		}
	} else {
		return srcValue
	}
}
