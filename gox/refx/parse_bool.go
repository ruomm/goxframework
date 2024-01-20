/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:22
 * @version 1.0
 */
package refx

import (
	"math"
	"reflect"
	"strconv"
	"time"
)

func xParseToBool(origVal interface{}, cpOpt string, isPointer bool, isTidy bool) interface{} {
	vi := ParseToBool(origVal, cpOpt)
	if vi == nil {
		return nil
	}
	rtVal := vi.(bool)
	if isTidy && !rtVal {
		return nil
	} else if isPointer {
		return &vi
	} else {
		return vi
	}
}

// 转换各种类型为bool类型，整形进行转换，字符串进行格式转换
func ParseToBool(origVal interface{}, cpOpt string) interface{} {
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
		int64Type := reflect.TypeOf(int64(0))
		if int64Type != actualValue.Type() {
			actualValue = actualValue.Convert(int64Type)
		}
		viInt64 := actualValue.Interface().(int64)
		if viInt64 > 0 {
			vi = true
		} else {
			vi = false
		}
	} else if xIsFloatKind(actualKind) {
		float64Type := reflect.TypeOf(float64(0))
		if float64Type != actualValue.Type() {
			actualValue = actualValue.Convert(float64Type)
		}
		viFloat64 := actualValue.Interface().(float64)
		viInt64 := int64(math.Round(viFloat64))
		if viInt64 > 0 {
			vi = true
		} else {
			vi = false
		}
	} else if actualKind == reflect.Bool {
		boolType := reflect.TypeOf(true)
		if boolType != actualValue.Type() {
			actualValue = actualValue.Convert(boolType)
		}
		vi = actualValue.Interface().(bool)
	} else if xIsStringKind(actualKind) {
		stringType := reflect.TypeOf("")
		if stringType != actualValue.Type() {
			actualValue = actualValue.Convert(stringType)
		}
		viString := actualValue.Interface().(string)
		viBool, err := strconv.ParseBool(viString)
		if err != nil && xTagContainKey(cpOpt, xReflect_key_string_bool_number) {
			viBoolByInt, errB := xTransStringIntToBool(viString, cpOpt)
			if errB != nil {
				viFloat64, errF := strconv.ParseFloat(viString, 64)
				if errF != nil {
					vi = nil
				} else {
					if int64(math.Round(viFloat64)) > 0 {
						vi = true
					} else {
						vi = false
					}
				}
			} else {
				vi = viBoolByInt
			}
		} else if err != nil {
			vi = nil
		} else {
			vi = viBool
		}

	} else if xIsStructType(actualKind) {
		srcFieldVT := reflect.TypeOf(origVal).String()
		if xIsStructType(actualKind) && xIsTimeType(srcFieldVT) {
			optStr := xTagFindValueByKey(cpOpt, xReflect_key_time_t)
			viTimeValue := actualValue.Interface().(time.Time)
			viInt64 := xTransTimeToInt64(&viTimeValue, optStr)
			if viInt64 > 0 {
				vi = true
			} else {
				vi = false
			}
		} else {
			vi = nil
		}
	} else {
		vi = nil
	}
	return vi
}
