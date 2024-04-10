/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/3/28 12:58
 * @version 1.0
 */
package refx

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	xRef_key_unsigned = "unsigned"
)

// 泛型函数，转换一个类型的Slice到另一个类型的Slice，使用xref库进行转换。
func XSliceCopy(srcSlice interface{}, destSlice interface{}, options ...XrefOption) error {
	if nil == destSlice {
		return errors.New("destSlice must not nil")
	}
	destSliceValue := reflect.ValueOf(destSlice)
	if destSliceValue.Kind() != reflect.Ptr && destSliceValue.Elem().Kind() != reflect.Slice {
		return errors.New("destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destSliceValue.Elem()
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return errors.New("srcSlice must be a slice")
	}
	// 如果目标切片容量不足，则扩展其容量
	if destSliceElem.Cap() < srcSliceValue.Len() {
		destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), srcSliceValue.Len(), srcSliceValue.Len()))
	} else if srcSliceValue.Len() == 0 {
		destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), 0, 0))
	}
	destSliceElem.SetLen(srcSliceValue.Len())
	var errG error = nil
	for i := 0; i < srcSliceValue.Len(); i++ {
		destValue := destSliceElem.Index(i)
		err, _ := XRefValueCopy(srcSliceValue.Index(i).Interface(), destValue, options...)
		if err != nil {
			errG = err
		}
	}
	return errG
}

// 泛型函数，转换一个类型的Slice到另一个类型的Slice，使用xref库进行转换。
func XSliceCopyByKey(srcSlice interface{}, destSlice interface{}, key string, options ...XrefOption) error {
	if nil == destSlice {
		return errors.New("destSlice must not nil")
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	//isTidy := xTagContainKey(cpOpt, xRef_key_tidy)
	checkUnsigned := xTagContainKey(do.copyOption, xRef_key_unsigned)
	destSliceValue := reflect.ValueOf(destSlice)
	if destSliceValue.Kind() != reflect.Ptr && destSliceValue.Elem().Kind() != reflect.Slice {
		return errors.New("destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destSliceValue.Elem()
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return errors.New("srcSlice must be a slice")
	}
	// 如果目标切片容量不足，则扩展其容量
	if destSliceElem.Cap() < srcSliceValue.Len() {
		destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), srcSliceValue.Len(), srcSliceValue.Len()))
	} else if srcSliceValue.Len() == 0 {
		destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), 0, 0))
	}
	destSliceElem.SetLen(srcSliceValue.Len())
	var errG error = nil
	for i := 0; i < srcSliceValue.Len(); i++ {
		destValue := destSliceElem.Index(i)
		srcItemValue := srcSliceValue.Index(i).FieldByName(key)
		origValue := srcItemValue.Interface()
		rtVal, transOk, transErr := xRef_transOrigToDestValue(key, do.copyOption, origValue, destValue, checkUnsigned)
		if transErr != nil {
			errG = transErr
		}
		if !transOk {
			errG = errors.New("XSliceCopyByKey excute failed")
		}
		if rtVal == nil {
			continue
		}
		// 赋值
		if destValue.CanSet() {
			kind := destValue.Type().Kind()
			if xIsIntegerKind(kind) {
				if xIsUnsignedIntegerKind(kind) {
					rtConvert := uint64(rtVal.(int64))
					destValue.SetUint(rtConvert)
				} else {
					rtConvert := rtVal.(int64)
					destValue.SetInt(rtConvert)
				}
			} else if xIsStringKind(kind) {
				rtConvert := rtVal.(string)
				destValue.SetString(rtConvert)
			} else if xIsFloatKind(kind) {
				rtConvert := rtVal.(float64)
				destValue.SetFloat(rtConvert)
			} else if kind == reflect.Bool {
				rtConvert := rtVal.(bool)
				destValue.SetBool(rtConvert)
			} else {
				destValue.Set(reflect.ValueOf(rtVal))
			}
		}
	}
	return errG
}

// 泛型函数，转换一个类型的Slice到另一个类型的Slice，使用xref库进行转换。
func oldXSliceCopyByKey(srcSlice interface{}, destSlice interface{}, key string) error {
	if nil == destSlice {
		return errors.New("destSlice must not nil")
	}
	destSliceValue := reflect.ValueOf(destSlice)
	if destSliceValue.Kind() != reflect.Ptr && destSliceValue.Elem().Kind() != reflect.Slice {
		return errors.New("destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destSliceValue.Elem()
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return errors.New("srcSlice must be a slice")
	}
	// 如果目标切片容量不足，则扩展其容量
	if destSliceElem.Cap() < srcSliceValue.Len() {
		destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), srcSliceValue.Len(), srcSliceValue.Len()))
	}
	destSliceElem.SetLen(srcSliceValue.Len())
	var errG error = nil
	for i := 0; i < srcSliceValue.Len(); i++ {
		destValue := destSliceElem.Index(i)
		srcItemValue := srcSliceValue.Index(i).FieldByName(key)
		origKind := srcItemValue.Type().Kind()
		destTypeOf := destValue.Type()
		destKind := destTypeOf.Kind()
		if origKind == reflect.Int {
			if xIsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(int)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(int)))
			}
		} else if origKind == reflect.Int8 {
			if xIsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(int8)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(int8)))
			}
		} else if origKind == reflect.Int16 {
			if xIsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(int16)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(int16)))
			}

		} else if origKind == reflect.Int32 {
			if xIsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(int32)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(int32)))
			}
		} else if origKind == reflect.Int64 {
			if xIsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(int64)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(int64)))
			}

		} else if origKind == reflect.Uint {
			if xIsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uint)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uint)))
			}

		} else if origKind == reflect.Uint8 {
			if xIsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uint8)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uint8)))
			}
		} else if origKind == reflect.Uint16 {
			if xIsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uint16)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uint16)))
			}

		} else if origKind == reflect.Uint32 {
			if xIsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uint32)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uint32)))
			}

		} else if origKind == reflect.Uint64 {
			if xIsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uint64)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uint64)))
			}
		} else if origKind == reflect.Float32 {
			destValue.SetFloat(float64(srcItemValue.Interface().(float32)))
		} else if origKind == reflect.Float64 {
			destValue.SetFloat(float64(srcItemValue.Interface().(float64)))
		} else if origKind == reflect.String {
			destValue.SetString(srcItemValue.Interface().(string))
		} else if origKind == reflect.Bool {
			destValue.SetBool(srcItemValue.Interface().(bool))
		} else {
			destValue.Set(srcItemValue)
		}
	}
	return errG
}

// 泛型函数，转换一个类型的Slice到另一个类型的Map，使用xref库进行转换。
func XSliceCopyToMap(srcSlice interface{}, destMap interface{}, keyTag string, valueTag string, options ...XrefOption) error {
	if nil == destMap {
		return errors.New("destSlice must not nil")
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	//isTidy := xTagContainKey(cpOpt, xRef_key_tidy)
	checkUnsigned := xTagContainKey(do.copyOption, xRef_key_unsigned)
	destMapValue := reflect.ValueOf(destMap)
	if destMapValue.Kind() != reflect.Ptr && destMapValue.Elem().Kind() != reflect.Map {
		return errors.New("destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destMapElem := destMapValue.Elem()
	// 获取map的键(key)类型

	keyType := reflect.TypeOf(destMapElem.Interface()).Key().Kind()
	keyTypeName := reflect.TypeOf(destMapElem.Interface()).Key().Name()

	// 获取map的值(value)类型
	valueType := reflect.TypeOf(destMapElem.Interface()).Elem().Kind()
	valueTypeName := reflect.TypeOf(destMapElem.Interface()).Elem().String()

	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return errors.New("srcSlice must be a slice")
	}
	destMapElem.Set(reflect.MakeMap(destMapElem.Type()))

	var errG error = nil
	for i := 0; i < srcSliceValue.Len(); i++ {
		keyItem := srcSliceValue.Index(i).FieldByName(keyTag)
		keyItemValue := keyItem.Interface()
		keyRtValue, keyTransOk, keyTransErr := xRefMap_transOrigToDestValue(valueTag, do.copyOption, keyItemValue, keyType, keyTypeName, checkUnsigned)
		if keyTransErr != nil {
			errG = keyTransErr
		}
		if !keyTransOk {
			errG = errors.New("XSliceCopyByKey excute failed")
		}
		if keyRtValue == nil {
			continue
		}
		if len(do.mapKeyAppend) > 0 && keyType == reflect.String {
			keyRtValue = do.mapKeyAppend + keyRtValue.(string)
		}
		valItem := srcSliceValue.Index(i).FieldByName(valueTag)
		valItemValue := valItem.Interface()
		valRtValue, valueTransOk, valueTransErr := xRefMap_transOrigToDestValue(valueTag, do.copyOption, valItemValue, valueType, valueTypeName, checkUnsigned)
		if valueTransErr != nil {
			errG = valueTransErr
		}
		if !valueTransOk {
			errG = errors.New("XSliceCopyByKey excute failed")
		}
		// 赋值
		destMapElem.SetMapIndex(reflect.ValueOf(xCopyToMapTypeTrans(keyType, keyRtValue)), reflect.ValueOf(xCopyToMapTypeTrans(valueType, valRtValue)))
	}
	return errG
}

// 解析来源字段值为目标待赋值字段
func xRefMap_transOrigToDestValue(key string, cpOpt string, origValue interface{}, destActualTypeKind reflect.Kind, destTypeName string, checkUnsigned bool) (interface{}, bool, error) {

	isTidy := xTagContainKey(cpOpt, xRef_key_tidy)
	//if xRef_log {
	//	fmt.Println(fmt.Sprintf("来源类型:%d-%s,目标类型:%d-%s,Tidy:%t", origKind, origType, destKind, destTypeName, isTidy))
	//}
	if xIsIntegerKind(destActualTypeKind) {
		return xParseToInt(key, origValue, destTypeName, destActualTypeKind, cpOpt, isTidy, checkUnsigned)
	} else if xIsFloatKind(destActualTypeKind) {
		parseVal, parseFlag := xParseToFloat(key, origValue, destTypeName, destActualTypeKind, cpOpt, isTidy)
		return parseVal, parseFlag, nil
	} else if destActualTypeKind == reflect.Bool {
		parseVal, parseFlag := xParseToBool(key, origValue, destTypeName, destActualTypeKind, cpOpt, isTidy)
		return parseVal, parseFlag, nil
	} else if destActualTypeKind == reflect.String {
		parseVal, parseFlag := xParseToString(key, origValue, destTypeName, destActualTypeKind, cpOpt, isTidy)
		return parseVal, parseFlag, nil
	} else if xIsTimeType(destTypeName) {
		parseVal, parseFlag := xParseToTime(key, origValue, destTypeName, cpOpt, isTidy)
		return parseVal, parseFlag, nil
	} else {
		// 目标是切片数组，来源是字符串时候解析字符串为数组
		if destActualTypeKind == reflect.Slice {
			// 获取真实的数值
			actualValue := reflect.ValueOf(origValue)
			if actualValue.Kind() == reflect.Pointer || actualValue.Kind() == reflect.Interface {
				if actualValue.IsNil() {
					return nil, true, nil
				}
				actualValue = actualValue.Elem()
			}
			actualKind := actualValue.Kind()
			if actualKind == reflect.String {
				stringType := reflect.TypeOf("")
				if stringType != actualValue.Type() {
					actualValue = actualValue.Convert(stringType)
				}
				viString := actualValue.Interface().(string)
				parseVal, parseFlag := xParseStringToSlice(key, viString, destTypeName, destActualTypeKind, cpOpt)
				return parseVal, parseFlag, nil
			}
		}

		origTypeOf := reflect.TypeOf(origValue)
		origKind := origTypeOf.Kind()
		origType := origTypeOf.String()
		if origKind != destActualTypeKind {
			if xRef_log {
				fmt.Println(key + "字段无法赋值，切片错误，目标和来源切片类型不同")
			}
			return nil, false, nil
		} else if origType != destTypeName {
			if xRef_log {
				fmt.Println(key + "字段无法赋值，结构错误，目标和来源结构类型不同")
			}
			return nil, false, nil
		} else {
			return origValue, true, nil
		}
	}
}

func xCopyToMapTypeTrans(destActualTypeKind reflect.Kind, vi interface{}) interface{} {
	if nil == vi {
		return nil
	}
	if destActualTypeKind == reflect.Int {
		viInt64 := vi.(int64)
		return int(viInt64)
	} else if destActualTypeKind == reflect.Int8 {
		viInt64 := vi.(int64)
		return int8(viInt64)
	} else if destActualTypeKind == reflect.Int16 {
		viInt64 := vi.(int64)
		return int16(viInt64)
	} else if destActualTypeKind == reflect.Int32 {
		viInt64 := vi.(int64)
		return int32(viInt64)
	} else if destActualTypeKind == reflect.Int64 {
		viInt64 := vi.(int64)
		return viInt64
	} else if destActualTypeKind == reflect.Uint {
		viInt64 := vi.(int64)
		return uint(viInt64)
	} else if destActualTypeKind == reflect.Uint8 {
		viInt64 := vi.(int64)
		return uint8(viInt64)
	} else if destActualTypeKind == reflect.Uint16 {
		viInt64 := vi.(int64)
		return uint16(viInt64)
	} else if destActualTypeKind == reflect.Uint32 {
		viInt64 := vi.(int64)
		return uint32(viInt64)
	} else if destActualTypeKind == reflect.Uint64 {
		viInt64 := vi.(int64)
		return uint64(viInt64)
	} else {
		return vi
	}

}
