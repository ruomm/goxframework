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
	"reflect"
	"strings"
)

const (
	xRef_AD_Zero_Second = int64(-62135596800)
	xRef_log            = true

	// const xReflect_tag_cp_opt = "xref_opt"
	// const xReflect_tag_cp_src = "xref"
	xRef_tag_cp_xreft = "xref"

	//var xReflect_location, _ = time.LoadLocation("Asia/Shanghai")

	xRef_time_layout = "2006-01-02 15:04:05"

	xRef_key_time_t             = "t"
	xRef_key_zero_to_8          = "z8"
	xRef_key_string_bool_number = "snb"
	xRef_key_time_tf            = "tf"
	xRef_key_number_point       = "p"

	// 如是omitempty参数存在，来源的数字类型的0、bool类型的false、字符串类型的空、时间类型的0或负数不会赋值的目标里面
	xRef_key_tidy = "tidy"
)

/*
*
srcO：源结构体
destO：目标切片，不可以传入结构体
*/
// TransferObj 将srcO对象的属性值转成destO对象的属性值，属性对应关系和控制指令通过`xref`标签指定
// 无标签的如果再按属性名匹配
func XRefCopy(srcO interface{}, destO interface{}) error {
	if nil == srcO {
		return errors.New("XRefCopy error,source interface is nil")
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
		tagXreft, okXreft := field.Tag.Lookup(xRef_tag_cp_xreft)
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
		if xRef_log {
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
		if srcValue == nil {
			if xRef_log {
				fmt.Println(key + "字段无需赋值，来源字段值为nil。")
			}
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

// 解析来源字段值为目标待赋值字段
func xReflect_transSrcToDestValue(key string, cpOpt string, srcValue interface{}, destValue reflect.Value) interface{} {

	//srcTypeOf := reflect.TypeOf(srcValue)
	//srcKind := srcTypeOf.Kind()
	//srcType := srcTypeOf.String()
	//m := xParseToInt64(srcValue, cpOpt)
	//println(m)
	destTypeOf := destValue.Type()
	//isPointor := false
	//if destValue.Kind() == reflect.Pointer || destValue.Kind() == reflect.Interface {
	//	isPointor = true
	//}
	destKind := destTypeOf.Kind()
	destTypeName := destTypeOf.String()
	destActualTypeKind := reflect.Invalid
	if destKind == reflect.Pointer {
		destActualTypeKind = destTypeOf.Elem().Kind()
	} else {
		destActualTypeKind = destKind
	}

	isTidy := xTagContainKey(cpOpt, xRef_key_tidy)
	//if xRef_log {
	//	fmt.Println(fmt.Sprintf("来源类型:%d-%s,目标类型:%d-%s,Tidy:%t", srcKind, srcType, destKind, destType, isTidy))
	//}
	if xIsIntegerKind(destActualTypeKind) {
		return xParseToInt(key, srcValue, destTypeName, destActualTypeKind, cpOpt, isTidy)
	} else if xIsFloatKind(destActualTypeKind) {
		return xParseToFloat(key, srcValue, destTypeName, destActualTypeKind, cpOpt, isTidy)
	} else if destActualTypeKind == reflect.Bool {
		return xParseToBool(key, srcValue, destTypeName, destActualTypeKind, cpOpt, isTidy)
	} else if destActualTypeKind == reflect.String {
		return xParseToString(key, srcValue, destTypeName, destActualTypeKind, cpOpt, isTidy)
	} else if xIsTimeType(destTypeOf.String()) {
		return xParseToTime(key, srcValue, destTypeName, cpOpt, isTidy)
	} else {
		srcTypeOf := reflect.TypeOf(srcValue)
		srcKind := srcTypeOf.Kind()
		srcType := srcTypeOf.String()
		if srcKind != destKind {
			if xRef_log {
				fmt.Println(key + "字段无法赋值，切片错误，目标和来源切片类型不同")
			}
			return nil
		} else if srcType != destTypeName {
			if xRef_log {
				fmt.Println(key + "字段无法赋值，结构错误，目标和来源结构类型不同")
			}
			return nil
		} else {
			return srcValue
		}
	}
}

func xTagContainKey(tagValue string, key string) bool {
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
func xTagFindValueByKey(tagValue string, key string) string {
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
func xIsUnsignedIntegerKind(kind reflect.Kind) bool {
	if kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
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

func xIsStringKind(kind reflect.Kind) bool {
	return kind == reflect.String
}

func xIsIntegerPointer(typeName string) bool {
	if typeName == "*int" || typeName == "*int8" || typeName == "*int16" || typeName == "*int32" || typeName == "*int64" ||
		typeName == "*uint" || typeName == "*uint8" || typeName == "*uint16" || typeName == "*uint32" || typeName == "*uint64" {
		return true
	} else {
		return false
	}
}
func xIsUnsignedIntegerPointer(typeName string) bool {
	if typeName == "*uint" || typeName == "*uint8" || typeName == "*uint16" || typeName == "*uint32" || typeName == "*uint64" {
		return true
	} else {
		return false
	}
}
func xIsFloatPointer(typeName string) bool {
	if typeName == "*float64" || typeName == "*float32" {
		return true
	} else {
		return false
	}
}

//func xIsStringPointer(typeName string) bool {
//	if typeName == "*string" {
//		return true
//	} else {
//		return false
//	}
//}

func xIsStructType(kind reflect.Kind) bool {
	return kind == reflect.Struct
}
func xIsStringTypeByName(typeName string) bool {
	return typeName == "string" || typeName == "*string"
}
func xIsTimeType(typeName string) bool {
	return typeName == "time.Time" || typeName == "Time" || typeName == "*time.Time"
}

func xIsPointor(typeName string) bool {
	return strings.HasPrefix(typeName, "*")
}
