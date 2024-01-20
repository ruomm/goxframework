/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:24
 * @version 1.0
 */
package refx

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

func xParseToString(origVal interface{}, destTypeName string, cpOpt string, isTidy bool) interface{} {
	vi := ParseToString(origVal, cpOpt)
	if vi == nil {
		return nil
	}
	viString := vi.(string)
	if isTidy && viString == "" {
		return nil
	}
	if !strings.HasPrefix(destTypeName, "*") {
		return vi
	}
	if destTypeName == "*string" {
		return &viString
	} else {
		return nil
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
		optStr := xTagFindValueByKey(cpOpt, xReflect_key_number_point)
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
		srcFieldVT := reflect.TypeOf(origVal).String()
		if xIsStructType(actualKind) && xIsTimeType(srcFieldVT) {
			optStr := xTagFindValueByKey(cpOpt, xReflect_key_time_tf)
			viTimeValue := actualValue.Interface().(time.Time)
			vi = xFormatTimeToString(&viTimeValue, optStr)
		} else {
			vi = nil
		}
	} else {
		vi = nil
	}
	return vi
}
