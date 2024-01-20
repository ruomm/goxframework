/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:21
 * @version 1.0
 */
package refx

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func xParseToInt(key string, origVal interface{}, destTypeName string, destActualTypeKind reflect.Kind, cpOpt string, isTidy bool) interface{} {
	vi := ParseToInt64(origVal, cpOpt)
	if vi == nil {
		if xRef_log {
			fmt.Println(key + "字段无法赋值，来源字段值无法解析或者为nil。")
		}
		return nil
	}
	viInt64 := vi.(int64)
	if isTidy && viInt64 == 0 {
		if xRef_log {
			fmt.Println(key + "字段无需赋值，来源字段值解析后的值为0。")
		}
		return nil
	}
	if !strings.HasPrefix(destTypeName, "*") {
		return vi
	}
	if destActualTypeKind == reflect.Int {
		rtVal := int(viInt64)
		return &rtVal
	} else if destActualTypeKind == reflect.Int8 {
		rtVal := int8(viInt64)
		return &rtVal
	} else if destActualTypeKind == reflect.Int16 {
		rtVal := int16(viInt64)
		return &rtVal
	} else if destActualTypeKind == reflect.Int32 {
		rtVal := int32(viInt64)
		return &rtVal
	} else if destActualTypeKind == reflect.Int64 {
		rtVal := viInt64
		return &rtVal
	} else if destActualTypeKind == reflect.Uint {
		rtVal := uint(viInt64)
		return &rtVal
	} else if destActualTypeKind == reflect.Uint8 {
		rtVal := uint8(viInt64)
		return &rtVal
	} else if destActualTypeKind == reflect.Uint16 {
		rtVal := uint16(viInt64)
		return &rtVal
	} else if destActualTypeKind == reflect.Uint32 {
		rtVal := uint32(viInt64)
		return &rtVal
	} else if destActualTypeKind == reflect.Uint64 {
		rtVal := uint64(viInt64)
		return &rtVal
	} else {
		if xRef_log {
			fmt.Println(key + "字段无需赋值，目标指针类型未知。")
		}
		return nil
	}

}

// 转换各种类型为int64，浮点型进行math.Round，字符串进行格式化，时间类型取得毫秒时间戳
func ParseToInt64(origVal interface{}, cpOpt string) interface{} {
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
	} else if xIsStringKind(actualKind) {
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
			optStr := xTagFindValueByKey(cpOpt, xRef_key_time_t)
			viTimeValue := actualValue.Interface().(time.Time)
			if viTimeValue.Unix() == xRef_AD_Zero_Second {
				vi = nil
			} else {
				vi = xTransTimeToInt64(&viTimeValue, optStr)
			}
		} else {
			vi = nil
		}
	} else {
		vi = nil
	}
	return vi
}
