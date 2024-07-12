/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:26
 * @version 1.0
 */
package refx

import (
	"errors"
	"github.com/morrisxyang/xreflect"
	"github.com/ruomm/goxframework/gox/corex"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 解析XRefValueCopy时候的key值
func xParseRefValueKey(key string) string {
	lenKey := len(key)
	if lenKey <= 0 {
		return key
	}
	indexSpec := strings.LastIndex(key, ".")
	if indexSpec >= 0 && indexSpec < lenKey-1 {
		return key[indexSpec+1:]
	} else {
		return key
	}
}

// 是否tidy
func XrefTagTidy(tagValue string) bool {
	return xTagContainKey(tagValue, xRef_key_tidy)
}

// 解析真实的tag名称
func xParseRefTagName(optTag string) string {
	if len(optTag) <= 0 {
		//real_tag_name = xRef_tag_cp_xreft
		return "xref"
	} else {
		return optTag
	}

}

// 解析真实的NameSpace
func xParseRefNameSpace(optNameSpace string, origNameSpace string) string {
	if len(optNameSpace) <= 0 {
		return origNameSpace
	} else {
		return optNameSpace
	}
}

// 属性空值设置 for refx.
type XrefOption struct {
	f func(*xrefOptions)
}

type xrefOptions struct {
	optTag        string            //关联的tag的名称
	optNameSpace  string            //关联源的nameSpace
	checkUnsigned bool              //无符号字符串是否严格匹配校验
	copyOption    string            //复制时候的控制属性
	mapKeyAppend  string            //设置Map赋值时候的拼接Key值
	copyDefault   bool              //设置没有注解时候是否默认复制
	copyMap       map[string]string //设置没有注解时候复制的来源目标Map数据
	copyList      []string          //设置没有注解时候复制的M来源目标List数据
}

// 设置关联tag的名称，不设置默认为xref
func XrefOptTag(tag string) XrefOption {
	return XrefOption{func(do *xrefOptions) {
		do.optTag = tag
	}}
}

// 设置关联源的nameSpace，不设置则取源对象的类型名称
func XrefOptNameSpace(nameSpace string) XrefOption {
	return XrefOption{func(do *xrefOptions) {
		do.optNameSpace = nameSpace
	}}
}

// 设置无符号字符串是否严格匹配校验
func XrefOptCheckUnsigned(checkUnsigned bool) XrefOption {
	return XrefOption{func(do *xrefOptions) {
		do.checkUnsigned = checkUnsigned
	}}
}

// 设置复制时候的控制属性
func XrefOptCopyOption(copyOption string) XrefOption {
	return XrefOption{func(do *xrefOptions) {
		do.copyOption = copyOption
	}}
}

// 设置没有注解时候是否默认复制
func XrefOptCopyDefault(copyDefault bool) XrefOption {
	return XrefOption{func(do *xrefOptions) {
		do.copyDefault = copyDefault
	}}
}

// 设置没有注解时候复制的来源目标Map数据
func XrefOptCopyMap(copyMap map[string]string) XrefOption {
	return XrefOption{func(do *xrefOptions) {
		do.copyMap = copyMap
	}}
}

// 设置没有注解时候复制的来源目标List数据
func XrefOptCopyList(copyList []string) XrefOption {
	return XrefOption{func(do *xrefOptions) {
		do.copyList = copyList
	}}
}

// 设置Map赋值时候的拼接Key值
func XrefOptMapKeyAppend(mapKeyAppend string) XrefOption {
	return XrefOption{func(do *xrefOptions) {
		do.mapKeyAppend = mapKeyAppend
	}}
}

// 字符串转换为int64
func xTransStringToInt64(viString string, cpOpt string) (int64, error) {
	if xTagContainKey(cpOpt, xRef_key_bytesize) {
		return corex.StoreSizeParse(viString)
	} else if xTagContainKey(cpOpt, xRef_key_timenumber_millis) {
		return corex.TimeNumberParse(viString, false)
	} else if xTagContainKey(cpOpt, xRef_key_timenumber_seconds) {
		return corex.TimeNumberParse(viString, true)
	}
	numBase := 10
	if strings.HasPrefix(viString, "0x") || strings.HasPrefix(viString, "0X") {
		numBase = 16
		viString = viString[2:]
	} else if strings.HasPrefix(viString, "-0x") || strings.HasPrefix(viString, "-0X") {
		numBase = 16
		viString = "-" + viString[3:]
	} else if (strings.HasPrefix(viString, "0") || strings.HasPrefix(viString, "-0")) && xTagContainKey(cpOpt, xRef_key_zero_to_8) {
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
			if xTagContainKey(cpOpt, xRef_key_string_bool_number) {
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

// 字符串转换为bool
func xTransStringIntToBool(viString string, cpOpt string) (bool, error) {
	numBase := 10
	if strings.HasPrefix(viString, "0x") || strings.HasPrefix(viString, "0X") {
		numBase = 16
		viString = viString[2:]
	} else if strings.HasPrefix(viString, "-0x") || strings.HasPrefix(viString, "-0X") {
		numBase = 16
		viString = "-" + viString[3:]
	} else if (strings.HasPrefix(viString, "0") || strings.HasPrefix(viString, "-0")) && xTagContainKey(cpOpt, xRef_key_zero_to_8) {
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

// 格式化时间为字符串
func xFormatTimeToString(t *time.Time, timeLayout string) string {
	//"America/Adak" "Asia/Shanghai"
	var realTimeLayout string
	if len(timeLayout) > 0 {
		realTimeLayout = timeLayout
	} else {
		realTimeLayout = xRef_time_layout
	}
	return t.In(corex.ToTimeLocation()).Format(realTimeLayout)
}

// 解析字符串为时间
func xTransStringToTime(sTime string, timeLayout string) *time.Time {
	var realTimeLayout string
	if len(timeLayout) > 0 {
		realTimeLayout = timeLayout
	} else {
		realTimeLayout = xRef_time_layout
	}
	timeStamp, err := time.ParseInLocation(realTimeLayout, sTime, corex.ToTimeLocation())
	if err != nil {
		return nil
	}
	return &timeStamp
}

func xTransInt64ToTime(origVal int64, optStr string) *time.Time {
	var timeDest *time.Time
	if len(optStr) <= 0 {
		timeValue := time.UnixMilli(origVal)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "sec") {
		timeValue := time.UnixMilli(origVal * 1000)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "min") {
		timeValue := time.UnixMilli(origVal * 1000 * 60)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "hour") {
		timeValue := time.UnixMilli(origVal * 1000 * 60 * 60)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "day") {
		timeValue := time.UnixMilli(origVal * 1000 * 60 * 60 * 24)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "mil") {
		timeValue := time.UnixMilli(origVal)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "mic") {
		timeValue := time.UnixMicro(origVal / 1e3)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "nano") {
		timeValue := time.UnixMilli(origVal / 1e6)
		timeDest = &timeValue
	} else {
		timeValue := time.UnixMilli(origVal)
		timeDest = &timeValue
	}
	return timeDest
}
func xTransTimeToInt64(pTime *time.Time, optStr string) int64 {
	if len(optStr) <= 0 {
		return pTime.UnixMilli()
	} else if strings.Contains(optStr, "sec") {
		return pTime.UnixMilli() / 1000
	} else if strings.Contains(optStr, "min") {
		return pTime.UnixMilli() / (1000 * 60)
	} else if strings.Contains(optStr, "hour") {
		return pTime.UnixMilli() / (1000 * 60 * 60)
	} else if strings.Contains(optStr, "day") {
		return pTime.UnixMilli() / (1000 * 60 * 60 * 24)
	} else if strings.Contains(optStr, "mil") {
		return pTime.UnixMilli()
	} else if strings.Contains(optStr, "mic") {
		return pTime.UnixMicro()
	} else if strings.Contains(optStr, "nano") {
		return pTime.UnixNano()
	} else {
		return pTime.UnixMilli()
	}
}

// 来源方法转换赋值
func xParseOrigValueByMethod(method_trans string, cpOpt string, origVal interface{}, destO interface{}) (interface{}, error) {
	actualValue := reflect.ValueOf(origVal)
	if xTagContainKey(cpOpt, xRef_key_method_trans_value_mode) {
		if actualValue.Kind() == reflect.Pointer || actualValue.Kind() == reflect.Interface {
			if actualValue.IsNil() {
				return nil, errors.New("字段无需赋值，来源字段值为nil，来源字段转换方法无法执行。")
			}
			actualValue = actualValue.Elem()
		}
	}
	callResultValues, err := xreflect.CallMethod(destO, method_trans, actualValue.Interface())
	if err != nil {
		return nil, errors.New("字段无法赋值，来源字段转换方法执行错误。")
	}
	if callResultValues == nil || len(callResultValues) <= 0 {
		return nil, errors.New("字段无需赋值，来源字段转换方法无有效值返回。")
	}
	actualDestValue := callResultValues[0]
	if actualDestValue.Kind() == reflect.Pointer || actualDestValue.Kind() == reflect.Interface {
		if actualDestValue.IsNil() {
			return nil, errors.New("字段无需赋值，来源字段值为nil，来源字段转换方法无法执行。")
		}
		actualDestValue = actualDestValue.Elem()
	}
	return actualDestValue.Interface(), nil
}

// 来源复制方法转换为复制的值
func xGetOrigValueByMethod(origO interface{}, origKey string, cpOpt string) (interface{}, error) {
	actualValue := reflect.ValueOf(origO)
	if xTagContainKey(cpOpt, xRef_key_getvalue_by_method_mode) {
		if actualValue.Kind() == reflect.Pointer || actualValue.Kind() == reflect.Interface {
			if actualValue.IsNil() {
				return nil, errors.New("字段无需赋值，来源对象为nil，来源复制方法无法执行。")
			}
			actualValue = actualValue.Elem()
		}
	}
	callResultValues, err := xreflect.CallMethod(actualValue.Interface(), origKey)
	if err != nil {
		return nil, errors.New("字段无法赋值，来源复制方法执行错误。")
	}
	if callResultValues == nil || len(callResultValues) <= 0 {
		return nil, errors.New("字段无需赋值，来源复制方法无有效值返回。")
	}
	actualDestValue := callResultValues[0]
	if actualDestValue.Kind() == reflect.Pointer || actualDestValue.Kind() == reflect.Interface {
		if actualDestValue.IsNil() {
			return nil, errors.New("字段无需赋值，来源复制方法返回值为nil。")
		}
		actualDestValue = actualDestValue.Elem()
	}
	return actualDestValue.Interface(), nil
}

// 依据copyDefault、copyMap、copyMapList获取复制的Key值
func xParseOrigCopyKeyByXrefOptions(s string, pField *reflect.StructField, do *xrefOptions) (string, bool) {
	fileNameSimply := corex.FieldNameToSimply(s)
	if len(fileNameSimply) <= 0 {
		return "", false
	}
	//tagOrig, okCanXref :=
	tagOrig := ""
	if len(do.copyMap) > 0 {
		for itemKey, itemValue := range do.copyMap {
			if fileNameSimply == itemKey {
				tagOrig = itemValue
				break
			}
		}
	}
	if len(tagOrig) > 0 {
		return tagOrig, true
	}
	if len(do.copyList) > 0 {
		for _, itemKey := range do.copyList {
			if fileNameSimply == itemKey {
				tagOrig = itemKey
				break
			}
		}
	}
	if len(tagOrig) > 0 {
		return tagOrig, true
	}
	if do.copyDefault && nil != pField {
		actualType := pField.Type
		if actualType.Kind() == reflect.Ptr {
			actualType = actualType.Elem()
		}
		actualKind := actualType.Kind()
		canCopy := false
		switch actualKind {
		case reflect.Bool:
			canCopy = true
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			canCopy = true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			canCopy = true
		case reflect.Float32, reflect.Float64:
			canCopy = true
		case reflect.String:
			canCopy = false
		case reflect.Array, reflect.Map, reflect.Slice:
			canCopy = false
		case reflect.Struct:
			canCopy = IsTimeTypeByName(actualType.String())
		case reflect.Interface, reflect.Ptr:
			canCopy = false
		case reflect.Func:
			canCopy = false
		case reflect.Invalid:
			canCopy = false
		default:
			canCopy = false
		}
		if canCopy {
			tagOrig = fileNameSimply
		}
	}
	if len(tagOrig) > 0 {
		return tagOrig, true
	} else {
		return "", false
	}
}
