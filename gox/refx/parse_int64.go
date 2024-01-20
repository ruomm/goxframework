/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:21
 * @version 1.0
 */
package refx

import (
	"math"
	"reflect"
	"strconv"
	"time"
)

// 转换各种类型为int64，浮点型进行math.Round，字符串进行格式化，时间类型取得毫秒时间戳
func xParseToInt64(origVal interface{}, cpOpt string) interface{} {
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
		vi = actualValue.Interface()
	} else if xIsFloatKind(actualKind) {
		float64Type := reflect.TypeOf(float64(0))
		if float64Type != actualValue.Type() {
			actualValue = actualValue.Convert(float64Type)
		}
		viFloat64 := actualValue.Interface().(float64)
		vi = int64(math.Round(viFloat64))
	} else if actualKind == reflect.Bool {
		boolType := reflect.TypeOf(true)
		if boolType != actualValue.Type() {
			actualValue = actualValue.Convert(boolType)
		}
		viBool := actualValue.Interface().(bool)
		if viBool {
			vi = int64(1)
		} else {
			vi = int64(0)
		}
	} else if xIsStringType(actualKind) {
		stringType := reflect.TypeOf("")
		if stringType != actualValue.Type() {
			actualValue = actualValue.Convert(stringType)
		}
		viString := actualValue.Interface().(string)
		viInt64, err := xTransStringToInt64(viString, cpOpt)
		if err != nil {
			viFloat64, errF := strconv.ParseFloat(viString, 64)
			if errF != nil {
				vi = nil
			} else {
				vi = int64(math.Round(viFloat64))
			}
		} else {
			vi = viInt64
		}

	} else if xIsStructType(actualKind) {
		srcFieldVT := reflect.TypeOf(origVal).String()
		if xIsStructType(actualKind) && xIsTimeType(srcFieldVT) {
			optStr := xTagFindValueByKey(cpOpt, xReflect_key_time_t)
			viTimeValue := actualValue.Interface().(time.Time)
			vi = xTransTimeToInt64(&viTimeValue, optStr)
		} else {
			vi = nil
		}
	} else {
		vi = nil
	}
	return vi
}
