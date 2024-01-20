/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:23
 * @version 1.0
 */
package refx

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func xParseToTime(key string, origVal interface{}, destTypeName string, cpOpt string, isTidy bool) (interface{}, bool) {
	vi := ParseToTime(origVal, cpOpt)
	if vi == nil {
		if xRef_log {
			fmt.Println(key + "字段无法赋值，来源字段值无法解析或者为nil。")
		}
		return nil, false
	}
	viPTime := vi.(*time.Time)
	if isTidy && (viPTime.Unix() == xRef_AD_Zero_Second || viPTime.UnixMilli() == 0) {
		if xRef_log {
			fmt.Println(key + "字段无需赋值，来源字段值解析后的时间值为空。")
		}
		return nil, true
	}
	if !strings.HasPrefix(destTypeName, "*") {
		return *viPTime, true
	}
	if destTypeName == "*time.Time" {
		return viPTime, true
	} else {
		if xRef_log {
			fmt.Println(key + "字段无需赋值，目标指针类型未知。")
		}
		return nil, false
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
		optStr := xTagFindValueByKey(cpOpt, xRef_key_time_t)
		vi = xTransInt64ToTime(viInt64, optStr)
	} else if xIsStringKind(actualKind) {
		stringType := reflect.TypeOf("")
		if stringType != actualValue.Type() {
			actualValue = actualValue.Convert(stringType)
		}
		viStr := actualValue.Interface().(string)
		optStr := xTagFindValueByKey(cpOpt, xRef_key_time_tf)
		pViTime := xTransStringToTime(viStr, optStr)
		if pViTime == nil {
			if xRef_log {
				fmt.Println("字段无法赋值，转换错误，string->time.Time")
			}
			return nil
		} else {
			vi = pViTime
		}
	} else if xIsStructType(actualKind) {
		origFieldVT := reflect.TypeOf(origVal).String()
		if xIsStructType(actualKind) && xIsTimeType(origFieldVT) {
			viTimeValue := actualValue.Interface().(time.Time)
			if viTimeValue.Unix() == xRef_AD_Zero_Second {
				vi = nil
			} else {
				vi = &viTimeValue
			}
		} else {
			vi = nil
		}
	} else {
		vi = nil
	}
	return vi
}
