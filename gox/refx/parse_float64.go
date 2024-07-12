/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:21
 * @version 1.0
 */
package refx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func xParseToFloat(key string, origVal interface{}, destTypeName string, destActualTypeKind reflect.Kind, cpOpt string, isTidy bool) (interface{}, bool) {
	vi := ParseToFloat64(origVal, cpOpt)
	if vi == nil {
		if xRef_log {
			fmt.Println(key + "字段无法赋值，来源字段值无法解析或者为nil。")
		}
		return nil, false
	}
	viFloat64 := vi.(float64)
	if isTidy && viFloat64 >= -0.0000000001 && viFloat64 <= 0.0000000001 {
		if xRef_log {
			fmt.Println(key + "字段无需赋值，来源字段值解析后的值为0。")
		}
		return nil, true
	}
	if !strings.HasPrefix(destTypeName, "*") {
		return vi, true
	}
	if destActualTypeKind == reflect.Float32 {
		rtVal := float32(viFloat64)
		return &rtVal, true
	} else if destActualTypeKind == reflect.Float64 {
		rtVal := viFloat64
		return &rtVal, true
	} else {
		if xRef_log {
			fmt.Println(key + "字段无需赋值，目标指针类型未知。")
		}
		return nil, false
	}
}

// 转换各种类型为浮点型，整形进行转换，字符串进行格式化，时间类型取得毫秒时间戳
func ParseToFloat64(origVal interface{}, cpOpt string) interface{} {
	// 获取真实的数值
	actualValue := reflect.ValueOf(origVal)
	if actualValue.Kind() == reflect.Pointer || actualValue.Kind() == reflect.Interface {
		if actualValue.IsNil() {
			return nil
		}
		actualValue = actualValue.Elem()
	}
	actualKind := actualValue.Kind()
	var vi interface{} = nil
	// 判断类型并转换
	if IsIntegerKind(actualKind) {
		int64Type := reflect.TypeOf(int64(0))
		if int64Type != actualValue.Type() {
			actualValue = actualValue.Convert(int64Type)
		}
		vi = float64(actualValue.Interface().(int64))
	} else if IsFloatKind(actualKind) {
		float64Type := reflect.TypeOf(float64(0))
		if float64Type != actualValue.Type() {
			actualValue = actualValue.Convert(float64Type)
		}
		vi = actualValue.Interface().(float64)
	} else if actualKind == reflect.Bool {
		boolType := reflect.TypeOf(true)
		if boolType != actualValue.Type() {
			actualValue = actualValue.Convert(boolType)
		}
		viBool := actualValue.Interface().(bool)
		if viBool {
			vi = float64(1)
		} else {
			vi = float64(0)
		}
	} else if IsStringKind(actualKind) {
		stringType := reflect.TypeOf("")
		if stringType != actualValue.Type() {
			actualValue = actualValue.Convert(stringType)
		}
		viString := actualValue.Interface().(string)
		viFloat64, err := strconv.ParseFloat(viString, 64)
		if err != nil {
			viInt64, errI := xTransStringToInt64(viString, cpOpt)
			if errI != nil {
				vi = nil
			} else {
				vi = float64(viInt64)
			}
		} else {
			vi = viFloat64
		}
	} else if IsStructKind(actualKind) {
		origFieldVT := reflect.TypeOf(origVal).String()
		if IsStructKind(actualKind) && IsTimeTypeByName(origFieldVT) {
			optStr := xTagFindValueByKey(cpOpt, xRef_key_time_t)
			viTimeValue := actualValue.Interface().(time.Time)
			if viTimeValue.Unix() == xRef_AD_Zero_Second {
				vi = nil
			} else {
				vi = float64(xTransTimeToInt64(&viTimeValue, optStr))
			}
		} else {
			vi = nil
		}
	} else {
		vi = nil
	}
	return vi
}
