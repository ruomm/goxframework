/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:24
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

func xParseToString(key string, origVal interface{}, destTypeName string, destActualTypeKind reflect.Kind, cpOpt string, isTidy bool) (interface{}, bool) {
	vi := ParseToString(origVal, cpOpt)
	if vi == nil {
		if xRef_log {
			fmt.Println(key + "字段无法赋值，来源字段值无法解析或者为nil。")
		}
		return nil, false
	}
	viString := vi.(string)
	if isTidy && viString == "" {
		if xRef_log {
			fmt.Println(key + "字段无需赋值，来源字段值解析后的值为空字符串。")
		}
		return nil, true
	}
	if !strings.HasPrefix(destTypeName, "*") {
		return vi, true
	}
	if destActualTypeKind == reflect.String {
		return &viString, true
	} else {
		if xRef_log {
			fmt.Println(key + "字段无需赋值，目标指针类型未知。")
		}
		return nil, false
	}
}

// 转换各种类型为字符串
func ParseToString(origVal interface{}, cpOpt string) interface{} {
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
	if xIsIntegerKind(actualKind) {
		if xIsUnsignedIntegerKind(actualKind) {
			uint64Type := reflect.TypeOf(uint64(0))
			if uint64Type != actualValue.Type() {
				actualValue = actualValue.Convert(uint64Type)
			}
			vi = strconv.FormatUint(actualValue.Interface().(uint64), 10)
		} else {
			int64Type := reflect.TypeOf(int64(0))
			if int64Type != actualValue.Type() {
				actualValue = actualValue.Convert(int64Type)
			}
			vi = strconv.FormatInt(actualValue.Interface().(int64), 10)
		}
	} else if xIsFloatKind(actualKind) {
		float64Type := reflect.TypeOf(float64(0))
		if float64Type != actualValue.Type() {
			actualValue = actualValue.Convert(float64Type)
		}
		viFloat64 := actualValue.Interface().(float64)
		optStr := xTagFindValueByKey(cpOpt, xRef_key_number_point)
		prec := -1
		if len(optStr) > 0 {
			prec64, errPrec := strconv.ParseInt(optStr, 10, 64)
			if errPrec != nil {
				prec = -1
			} else if prec64 >= -1 && prec64 <= 10 {
				prec = int(prec64)
			} else {
				prec = -1
			}
		}
		return strconv.FormatFloat(viFloat64, 'f', prec, 64)
	} else if actualKind == reflect.Bool {
		boolType := reflect.TypeOf(true)
		if boolType != actualValue.Type() {
			actualValue = actualValue.Convert(boolType)
		}
		viBool := actualValue.Interface().(bool)
		if viBool {
			vi = "true"
		} else {
			vi = "false"
		}
	} else if xIsStringKind(actualKind) {
		stringType := reflect.TypeOf("")
		if stringType != actualValue.Type() {
			actualValue = actualValue.Convert(stringType)
		}
		vi = actualValue.Interface().(string)
	} else if xIsStructType(actualKind) {
		origFieldVT := reflect.TypeOf(origVal).String()
		if xIsStructType(actualKind) && xIsTimeType(origFieldVT) {
			optStr := xTagFindValueByKey(cpOpt, xRef_key_time_tf)
			viTimeValue := actualValue.Interface().(time.Time)
			if viTimeValue.Unix() == xRef_AD_Zero_Second {
				vi = nil
			} else {
				vi = xFormatTimeToString(&viTimeValue, optStr)
			}
		} else {
			vi = nil
		}
	} else {
		vi = nil
	}
	return vi
}
