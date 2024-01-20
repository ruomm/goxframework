/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:23
 * @version 1.0
 */
package refx

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

func xParseToTime(origVal interface{}, cpOpt string, isPointer bool, isTidy bool) interface{} {
	vi := ParseToTime(origVal, cpOpt)
	if vi == nil {
		return nil
	}
	rtVal := vi.(*time.Time)
	if isTidy && rtVal.UnixMilli() == 0 {
		return nil
	} else if isPointer {
		return rtVal
	} else {
		return *rtVal
	}
}

// 转换各种类型为时间类型
func ParseToTime(origVal interface{}, cpOpt string) interface{} {
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
		optStr := xTagFindValueByKey(cpOpt, xReflect_key_time_t)
		vi = xTransInt64ToTime(viInt64, optStr)
	} else if xIsFloatKind(actualKind) {
		float64Type := reflect.TypeOf(float64(0))
		if float64Type != actualValue.Type() {
			actualValue = actualValue.Convert(float64Type)
		}
		viFloat64 := actualValue.Interface().(float64)
		viInt64 := int64(math.Round(viFloat64))
		optStr := xTagFindValueByKey(cpOpt, xReflect_key_time_t)
		vi = xTransInt64ToTime(viInt64, optStr)
	} else if actualKind == reflect.Bool {
		vi = nil
	} else if xIsStringKind(actualKind) {
		stringType := reflect.TypeOf("")
		if stringType != actualValue.Type() {
			actualValue = actualValue.Convert(stringType)
		}
		viStr := actualValue.Interface().(string)
		optStr := xTagFindValueByKey(cpOpt, xReflect_key_time_tf)
		pViTime := xTransStringToTime(viStr, optStr)
		if pViTime == nil {
			if xReflect_log {
				fmt.Println("字段无法赋值，转换错误，string->time.Time")
			}
			return nil
		} else {
			vi = pViTime
		}
	} else if xIsStructType(actualKind) {
		srcFieldVT := reflect.TypeOf(origVal).String()
		if xIsStructType(actualKind) && xIsTimeType(srcFieldVT) {
			viTimeValue := actualValue.Interface().(time.Time)
			vi = &viTimeValue
		} else {
			vi = nil
		}
	} else {
		vi = nil
	}
	return vi
}
