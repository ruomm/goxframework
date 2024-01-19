/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:02
 * @version 1.0
 */
package refx

import (
	"errors"
	"fmt"
	"github.com/morrisxyang/xreflect"
	"github.com/ruomm/goxframework/gox/corex"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	xReflect_log = true

	// const xReflect_tag_cp_opt = "xref_opt"
	// const xReflect_tag_cp_src = "xref"
	xReflect_tag_cp_xreft = "xref"

	//var xReflect_location, _ = time.LoadLocation("Asia/Shanghai")

	xReflect_time_layout = "2006-01-02 15:04:05"

	xReflect_key_time_t             = "t"
	xReflect_key_zero_to_8          = "z8"
	xReflect_key_string_bool_number = "snb"
	xReflect_key_time_tf            = "tf"
	xReflect_key_number_point       = "p"

	// 如是omitempty参数存在，来源的数字类型的0、bool类型的false、字符串类型的空、时间类型的0或负数不会赋值的目标里面
	xReflect_key_tidy = "tidy"
)

/*
*
srcO：源结构体
destO：目标切片，不可以传入结构体
*/
// TransferObj 将srcO对象的属性值转成destO对象的属性值，属性对应关系和控制指令通过`xref`标签指定
// 无标签的如果再按属性名匹配
func XReflectCopy(srcO interface{}, destO interface{}, transBasePointer ...bool) error {
	if nil == srcO {
		return errors.New("XReflectCopy error,source interface is nil")
	}
	var transBasePointerFlag = false
	if nil != transBasePointer && len(transBasePointer) > 0 && transBasePointer[0] {
		transBasePointerFlag = transBasePointer[0]
	}
	// 获取srcO的类名称
	srcT := reflect.TypeOf(srcO)
	var srcNameStr string
	if srcT.Kind() == reflect.Array || srcT.Kind() == reflect.Chan || srcT.Kind() == reflect.Map || srcT.Kind() == reflect.Pointer || srcT.Kind() == reflect.Slice {
		srcNameStr = srcT.Elem().String()
	} else {
		srcNameStr = srcT.String()
	}
	resOpt := make(map[string]string)
	resSrc := make(map[string]string)
	reflectValueMap, err := xreflect.SelectFieldsDeep(destO, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagXreft, okXreft := field.Tag.Lookup(xReflect_tag_cp_xreft)
		if !okXreft {
			return false
		}
		// 开始分割目标控制和属性控制
		subTags := corex.ParseToSubTag(tagXreft)
		// 解析目标控制
		tagSrcVal := ""
		if len(subTags) > 0 {
			tagSrcVal = subTags[0]
		}
		tagSrc, okCanXref := xReflect_canXCopy(tagSrcVal, srcNameStr)
		if !okCanXref {
			return false
		}
		resSrc[s] = tagSrc
		// 解析属性控制
		tagOpt := ""
		if len(subTags) > 1 {
			tagOpt = subTags[1]
		}
		resOpt[s] = tagOpt
		if xReflect_log {
			fmt.Println("解析复制字段，目标：" + s + "，来源：" + tagSrc + "，控制协议：" + tagOpt)
		}
		return true
	})
	if err != nil {
		return err
	}
	for key, value := range reflectValueMap {
		var srcKey string
		if resSrc[key] != "" {
			srcKey = resSrc[key]
		} else {
			srcKey = key
		}
		srcValue, err := xreflect.EmbedFieldValue(srcO, srcKey)
		if err != nil {
			continue
		}
		if transBasePointerFlag {
			//srcValue = transBasePointerToValue(srcValue)
			//srcValue = transBasePointerToValue(srcValue)
			fmt.Println("transBasePointerToValue")
		}
		if srcValue == nil {
			continue
		}
		cpOpt := resOpt[key]
		rtVal := xReflect_transSrcToDestValue(key, cpOpt, srcValue, value)
		if rtVal == nil {
			continue
		}
		tmpErr := xreflect.SetEmbedField(destO, key, rtVal)
		if tmpErr != nil {
			err = tmpErr
		}

	}
	return err
}

// 字段是否需要XReflect复制
func xReflect_canXCopy(tagSrcVal string, srcNameStr string) (string, bool) {
	if tagSrcVal == "" {
		return "", true
	}
	cpEnable := false
	cpSrcId := ""
	srcVlist := strings.Split(tagSrcVal, ",")
	for _, srcV := range srcVlist {
		if srcV == "" {
			continue
		}
		subVList := strings.Split(srcV, "-")
		lenVList := len(subVList)
		if lenVList == 0 {
			continue
		} else if lenVList == 1 {
			if len(subVList[0]) > 0 {
				cpSrcId = subVList[0]
				cpEnable = true
			}
		} else if lenVList == 2 {
			if len(subVList[0]) > 0 && (strings.HasSuffix(srcNameStr, subVList[0]) || strings.HasPrefix(srcNameStr, subVList[0])) {
				cpEnable = true
				cpSrcId = subVList[1]
				break
			}
		}
	}
	return cpSrcId, cpEnable
}

// 解析来源字段为目标待赋值字段
func xReflect_transSrcToDestValue(key string, cpOpt string, srcValue interface{}, destValue reflect.Value) interface{} {

	//srcTypeOf := reflect.TypeOf(srcValue)
	//srcKind := srcTypeOf.Kind()
	//srcType := srcTypeOf.String()
	//m := xReflect_transToInt64(srcValue, cpOpt)
	//println(m)
	destTypeOf := destValue.Type()
	destKind := destTypeOf.Kind()
	//destType := destTypeOf.String()
	//isTidy := xReflect_tagContainKey(cpOpt, xReflect_key_tidy)
	//if xReflect_log {
	//	fmt.Println(fmt.Sprintf("来源类型:%d-%s,目标类型:%d-%s,Tidy:%t", srcKind, srcType, destKind, destType, isTidy))
	//}
	if xIsIntegerKind(destKind) {
		return xReflect_transToInt64(srcValue, cpOpt)
	} else if xIsFloatKind(destKind) {
		return xReflect_transToFloat64(srcValue, cpOpt)
	} else if destKind == reflect.Bool {
		return xReflect_transToBool(srcValue,cpOpt)
	} else if destKind == reflect.String {
		return nil
	}
	return nil
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
func xIsIntegerKind(kind reflect.Kind) bool {
	if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 ||
		kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
		return true
	} else {
		return false
	}
}
func xIsFloatKind(kind reflect.Kind) bool {
	if kind == reflect.Float64 || kind == reflect.Float32 {
		return true
	} else {
		return false
	}
}

func xIsStringType(kind reflect.Kind) bool {
	return kind == reflect.String
}

func xIsStructType(kind reflect.Kind) bool {
	return kind == reflect.Struct
}
func xIsStringTypeByName(typeName string) bool {
	return typeName == "string" || typeName == "*string"
}
func xIsTimeType(typeName string) bool {
	return typeName == "time.Time" || typeName == "Time" || typeName == "*time.Time"
}

// 转换各种类型为int64，浮点型进行math.Round，字符串进行格式化，时间类型取得毫秒时间戳
func xReflect_transToInt64(origVal interface{}, cpOpt string) interface{} {
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
			optStr := xReflect_findTagValueByKey(cpOpt, xReflect_key_time_t)
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

// 转换各种类型为浮点型，整形进行转换，字符串进行格式化，时间类型取得毫秒时间戳
func xReflect_transToFloat64(origVal interface{}, cpOpt string) interface{} {
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
		vi = float64(actualValue.Interface().(int64))
	} else if xIsFloatKind(actualKind) {
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
	} else if xIsStringType(actualKind) {
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
	} else if xIsStructType(actualKind) {
		srcFieldVT := reflect.TypeOf(origVal).String()
		if xIsStructType(actualKind) && xIsTimeType(srcFieldVT) {
			optStr := xReflect_findTagValueByKey(cpOpt, xReflect_key_time_t)
			viTimeValue := actualValue.Interface().(time.Time)
			vi = float64(xTransTimeToInt64(&viTimeValue, optStr))
		} else {
			vi = nil
		}
	} else {
		vi = nil
	}
	return vi
}

// 转换各种类型为bool类型，整形进行转换，字符串进行格式转换
func xReflect_transToBool(origVal interface{}, cpOpt string) interface{} {
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
	} else if xIsStringType(actualKind) {
		stringType := reflect.TypeOf("")
		if stringType != actualValue.Type() {
			actualValue = actualValue.Convert(stringType)
		}
		viString := actualValue.Interface().(string)
		viBool, err := strconv.ParseBool(viString)
		if err != nil && xReflect_tagContainKey(cpOpt, xReflect_key_string_bool_number) {
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
			optStr := xReflect_findTagValueByKey(cpOpt, xReflect_key_time_t)
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

// 字符串转换为int64
func xTransStringToInt64(viString string, cpOpt string) (int64, error) {
	numBase := 10
	if strings.HasPrefix(viString, "0x") || strings.HasPrefix(viString, "0X") {
		numBase = 16
		viString = viString[2:]
	} else if strings.HasPrefix(viString, "-0x") || strings.HasPrefix(viString, "-0X") {
		numBase = 16
		viString = "-" + viString[3:]
	} else if (strings.HasPrefix(viString, "0") || strings.HasPrefix(viString, "-0")) && xReflect_tagContainKey(cpOpt, xReflect_key_zero_to_8) {
		numBase = 8
	}
	if strings.HasPrefix(viString, "-") {
		viInt64, err := strconv.ParseInt(viString, numBase, 64)
		if err != nil {
			return 0, err
		} else {
			return viInt64, nil
		}
	} else {
		viUint64, err := strconv.ParseUint(viString, numBase, 64)
		if err != nil {
			if xReflect_tagContainKey(cpOpt, xReflect_key_string_bool_number) {
				viBool, errB := strconv.ParseBool(viString)
				if errB != nil {
					return 0, errB
				} else if viBool {
					return 1, nil
				} else {
					return 0, nil
				}
			} else {
				return 0, err
			}
		} else {
			return int64(viUint64), nil
		}
	}
}

// 字符串转换为int64
func xTransStringIntToBool(viString string, cpOpt string) (bool, error) {
	numBase := 10
	if strings.HasPrefix(viString, "0x") || strings.HasPrefix(viString, "0X") {
		numBase = 16
		viString = viString[2:]
	} else if strings.HasPrefix(viString, "-0x") || strings.HasPrefix(viString, "-0X") {
		numBase = 16
		viString = "-" + viString[3:]
	} else if (strings.HasPrefix(viString, "0") || strings.HasPrefix(viString, "-0")) && xReflect_tagContainKey(cpOpt, xReflect_key_zero_to_8) {
		numBase = 8
	}
	if strings.HasPrefix(viString, "-") {
		viInt64, err := strconv.ParseInt(viString, numBase, 64)
		if err != nil {
			return false, err
		} else {
			return viInt64 > 0, nil
		}
	} else {
		viUint64, err := strconv.ParseUint(viString, numBase, 64)
		if err != nil {
			return false, err
		} else {
			return viUint64 > 0, nil
		}
	}
}

func xTransInt64ToTime(srcVal int64, optStr1 string) *time.Time {
	var timeDest *time.Time
	if len(optStr1) <= 0 {
		timeValue := time.UnixMilli(srcVal)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "sec") {
		timeValue := time.UnixMilli(srcVal * 1000)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "min") {
		timeValue := time.UnixMilli(srcVal * 1000 * 60)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "hour") {
		timeValue := time.UnixMilli(srcVal * 1000 * 60 * 60)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "day") {
		timeValue := time.UnixMilli(srcVal * 1000 * 60 * 60 * 24)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "mil") {
		timeValue := time.UnixMilli(srcVal)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "mic") {
		timeValue := time.UnixMicro(srcVal / 1e3)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "nano") {
		timeValue := time.UnixMilli(srcVal / 1e6)
		timeDest = &timeValue
	} else {
		timeValue := time.UnixMilli(srcVal)
		timeDest = &timeValue
	}
	return timeDest
}
func xTransTimeToInt64(pTime *time.Time, optStr1 string) int64 {
	if len(optStr1) <= 0 {
		return pTime.UnixMilli()
	} else if strings.Contains(optStr1, "sec") {
		return pTime.UnixMilli() / 1000
	} else if strings.Contains(optStr1, "min") {
		return pTime.UnixMilli() / (1000 * 60)
	} else if strings.Contains(optStr1, "hour") {
		return pTime.UnixMilli() / (1000 * 60 * 60)
	} else if strings.Contains(optStr1, "day") {
		return pTime.UnixMilli() / (1000 * 60 * 60 * 24)
	} else if strings.Contains(optStr1, "mil") {
		return pTime.UnixMilli()
	} else if strings.Contains(optStr1, "mic") {
		return pTime.UnixMicro()
	} else if strings.Contains(optStr1, "nano") {
		return pTime.UnixNano()
	} else {
		return pTime.UnixMilli()
	}
}
