/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:22
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

func xParseToBool(key string, origVal interface{}, destTypeName string, destActualTypeKind reflect.Kind, cpOpt string, isTidy bool) (interface{}, bool) {
	vi := ParseToBool(origVal, cpOpt)
	if vi == nil {
		if xRef_log {
			fmt.Println(key + "字段无法赋值，来源字段值无法解析或者为nil。")
		}
		return nil, false
	}
	viBool := vi.(bool)
	if isTidy && !viBool {
		if xRef_log {
			fmt.Println(key + "字段无需赋值，来源字段值解析后的值为false。")
		}
		return nil, true
	}
	if !strings.HasPrefix(destTypeName, "*") {
		return vi, true
	}
	if destActualTypeKind == reflect.Bool {
		return &viBool, true
	} else {
		if xRef_log {
			fmt.Println(key + "字段无需赋值，目标指针类型未知。")
		}
		return nil, false
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
	if IsIntegerKind(actualKind) {
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
	} else if IsFloatKind(actualKind) {
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
	} else if IsStringKind(actualKind) {
		stringType := reflect.TypeOf("")
		if stringType != actualValue.Type() {
			actualValue = actualValue.Convert(stringType)
		}
		viString := actualValue.Interface().(string)
		viBool, err := strconv.ParseBool(viString)
		if err != nil && xTagContainKey(cpOpt, xRef_key_string_bool_number) {
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

	} else if IsStructKind(actualKind) {
		origFieldVT := reflect.TypeOf(origVal).String()
		if IsStructKind(actualKind) && IsTimeTypeByName(origFieldVT) {
			optStr := xTagFindValueByKey(cpOpt, xRef_key_time_t)
			viTimeValue := actualValue.Interface().(time.Time)
			if viTimeValue.Unix() == xRef_AD_Zero_Second {
				vi = nil
			} else {
				viInt64 := xTransTimeToInt64(&viTimeValue, optStr)
				if viInt64 > 0 {
					vi = true
				} else {
					vi = false
				}
			}

		} else {
			vi = nil
		}
	} else {
		vi = nil
	}
	return vi
}
