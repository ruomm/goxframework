/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/3/28 12:58
 * @version 1.0
 */
package refxstandard

import (
	"errors"
	"fmt"
	"github.com/ruomm/goxframework/gox/corex"
	"reflect"
)

const (
	xRef_key_unsigned = "unsigned"
)

// 泛型函数，转换字符串为特定类型的Slice
func XStringToSlice(str string, sep string, emptyRetain bool, destSlice interface{}, options ...XrefOption) error {
	tag := "XStringToSlice:"
	if nil == destSlice {
		return errors.New(tag + "destSlice must not nil")
	}
	srcSlice := corex.StringToSlice(str, sep, emptyRetain)
	lenSrcSlice := len(srcSlice)
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	//isTidy := xTagContainKey(cpOpt, xRef_key_tidy)
	checkUnsigned := xTagContainKey(do.copyOption, xRef_key_unsigned)
	destSliceValue := reflect.ValueOf(destSlice)
	if destSliceValue.Kind() != reflect.Ptr && destSliceValue.Elem().Kind() != reflect.Slice {
		return errors.New(tag + "destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destSliceValue.Elem()
	// 如果目标切片容量不足，则扩展其容量
	if destSliceElem.Cap() < lenSrcSlice {
		destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), lenSrcSlice, lenSrcSlice))
	} else if lenSrcSlice == 0 {
		destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), 0, 0))
	}
	destSliceElem.SetLen(lenSrcSlice)
	var errG error = nil
	for i := 0; i < lenSrcSlice; i++ {
		destValue := destSliceElem.Index(i)
		origValue := srcSlice[i]
		key := origValue
		if len(key) <= 0 {
			key = "EmptyString"
		}
		rtVal, transOk, transErr := xRef_transOrigToDestValue(key, do.copyOption, origValue, destValue, checkUnsigned)
		if transErr != nil {
			errG = transErr
		}
		if !transOk {
			errG = errors.New(tag + "excute failed")
		}
		if rtVal == nil {
			continue
		}
		// 赋值
		if destValue.CanSet() {
			kind := destValue.Type().Kind()
			if IsIntegerKind(kind) {
				if IsUnsignedIntegerKind(kind) {
					rtConvert := uint64(rtVal.(int64))
					destValue.SetUint(rtConvert)
				} else {
					rtConvert := rtVal.(int64)
					destValue.SetInt(rtConvert)
				}
			} else if IsStringKind(kind) {
				rtConvert := rtVal.(string)
				destValue.SetString(rtConvert)
			} else if IsFloatKind(kind) {
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
func XSliceCopy(srcSlice interface{}, destSlice interface{}, options ...XrefOption) error {
	tag := "XSliceCopy:"
	if nil == destSlice {
		return errors.New(tag + "destSlice must not nil")
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	destSliceValue := reflect.ValueOf(destSlice)
	if destSliceValue.Kind() != reflect.Ptr && destSliceValue.Elem().Kind() != reflect.Slice {
		return errors.New(tag + "destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destSliceValue.Elem()
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return errors.New(tag + "srcSlice must be a slice")
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
	tag := "XSliceCopyByKey:"
	if nil == destSlice {
		return errors.New(tag + "destSlice must not nil")
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	//isTidy := xTagContainKey(cpOpt, xRef_key_tidy)
	checkUnsigned := xTagContainKey(do.copyOption, xRef_key_unsigned)
	destSliceValue := reflect.ValueOf(destSlice)
	if destSliceValue.Kind() != reflect.Ptr && destSliceValue.Elem().Kind() != reflect.Slice {
		return errors.New(tag + "destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destSliceValue.Elem()
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return errors.New(tag + "srcSlice must be a slice")
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
		valItemActual := srcSliceValue.Index(i)
		if valItemActual.Kind() == reflect.Pointer || valItemActual.Kind() == reflect.Interface {
			if valItemActual.IsNil() {
				continue
			}
			valItemActual = valItemActual.Elem()
		}
		srcItemValue := valItemActual.FieldByName(key)

		origValue := srcItemValue.Interface()
		rtVal, transOk, transErr := xRef_transOrigToDestValue(key, do.copyOption, origValue, destValue, checkUnsigned)
		if transErr != nil {
			errG = transErr
		}
		if !transOk {
			errG = errors.New(tag + "excute failed")
		}
		if rtVal == nil {
			continue
		}
		// 赋值
		if destValue.CanSet() {
			kind := destValue.Type().Kind()
			if IsIntegerKind(kind) {
				if IsUnsignedIntegerKind(kind) {
					rtConvert := uint64(rtVal.(int64))
					destValue.SetUint(rtConvert)
				} else {
					rtConvert := rtVal.(int64)
					destValue.SetInt(rtConvert)
				}
			} else if IsStringKind(kind) {
				rtConvert := rtVal.(string)
				destValue.SetString(rtConvert)
			} else if IsFloatKind(kind) {
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

// 泛型函数，浅复制一个类型的Slice到另一个类型的Slice，使用xref库进行转换，目标slice不进行扩容、缩容操作。
func XSliceQcopy(srcSlice interface{}, destSlice interface{}, options ...XrefOption) error {
	tag := "XSliceQcopy:"
	if nil == destSlice {
		return errors.New(tag + "destSlice must not nil")
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	destSliceValue := reflect.ValueOf(destSlice)
	if destSliceValue.Kind() != reflect.Ptr && destSliceValue.Elem().Kind() != reflect.Slice {
		return errors.New(tag + "destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destSliceValue.Elem()
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return errors.New(tag + "srcSlice must be a slice")
	}
	//// 如果目标切片容量不足，则扩展其容量
	//if destSliceElem.Cap() < srcSliceValue.Len() {
	//	destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), srcSliceValue.Len(), srcSliceValue.Len()))
	//} else if srcSliceValue.Len() == 0 {
	//	destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), 0, 0))
	//}
	//destSliceElem.SetLen(srcSliceValue.Len())
	var errG error = nil
	for i := 0; i < srcSliceValue.Len() && i < destSliceElem.Len(); i++ {
		destValue := destSliceElem.Index(i)
		err, _ := XRefValueCopy(srcSliceValue.Index(i).Interface(), destValue, options...)
		if err != nil {
			errG = err
		}
	}
	return errG
}

// 泛型函数，浅复制一个类型的Slice到另一个类型的Slice，使用xref库进行转换，目标slice不进行扩容、缩容操作。
func XSliceQcopyByKey(srcSlice interface{}, destSlice interface{}, key string, options ...XrefOption) error {
	tag := "XSliceQcopyByKey:"
	if nil == destSlice {
		return errors.New(tag + "destSlice must not nil")
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	//isTidy := xTagContainKey(cpOpt, xRef_key_tidy)
	checkUnsigned := xTagContainKey(do.copyOption, xRef_key_unsigned)
	destSliceValue := reflect.ValueOf(destSlice)
	if destSliceValue.Kind() != reflect.Ptr && destSliceValue.Elem().Kind() != reflect.Slice {
		return errors.New(tag + "destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destSliceValue.Elem()
	srcSliceValue := reflect.ValueOf(srcSlice)
	//if srcSliceValue.Kind() != reflect.Slice {
	//	return errors.New(tag + "srcSlice must be a slice")
	//}
	//// 如果目标切片容量不足，则扩展其容量
	//if destSliceElem.Cap() < srcSliceValue.Len() {
	//	destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), srcSliceValue.Len(), srcSliceValue.Len()))
	//} else if srcSliceValue.Len() == 0 {
	//	destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), 0, 0))
	//}
	//destSliceElem.SetLen(srcSliceValue.Len())
	var errG error = nil
	for i := 0; i < srcSliceValue.Len() && i < destSliceElem.Len(); i++ {
		destValue := destSliceElem.Index(i)
		valItemActual := srcSliceValue.Index(i)
		if valItemActual.Kind() == reflect.Pointer || valItemActual.Kind() == reflect.Interface {
			if valItemActual.IsNil() {
				continue
			}
			valItemActual = valItemActual.Elem()
		}
		srcItemValue := valItemActual.FieldByName(key)

		origValue := srcItemValue.Interface()
		rtVal, transOk, transErr := xRef_transOrigToDestValue(key, do.copyOption, origValue, destValue, checkUnsigned)
		if transErr != nil {
			errG = transErr
		}
		if !transOk {
			errG = errors.New(tag + "excute failed")
		}
		if rtVal == nil {
			continue
		}
		// 赋值
		if destValue.CanSet() {
			kind := destValue.Type().Kind()
			if IsIntegerKind(kind) {
				if IsUnsignedIntegerKind(kind) {
					rtConvert := uint64(rtVal.(int64))
					destValue.SetUint(rtConvert)
				} else {
					rtConvert := rtVal.(int64)
					destValue.SetInt(rtConvert)
				}
			} else if IsStringKind(kind) {
				rtConvert := rtVal.(string)
				destValue.SetString(rtConvert)
			} else if IsFloatKind(kind) {
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
	tag := "oldXSliceCopyByKey:"
	if nil == destSlice {
		return errors.New(tag + "destSlice must not nil")
	}
	destSliceValue := reflect.ValueOf(destSlice)
	if destSliceValue.Kind() != reflect.Ptr && destSliceValue.Elem().Kind() != reflect.Slice {
		return errors.New(tag + "destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destSliceValue.Elem()
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return errors.New(tag + "srcSlice must be a slice")
	}
	// 如果目标切片容量不足，则扩展其容量
	if destSliceElem.Cap() < srcSliceValue.Len() {
		destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), srcSliceValue.Len(), srcSliceValue.Len()))
	}
	destSliceElem.SetLen(srcSliceValue.Len())
	var errG error = nil
	for i := 0; i < srcSliceValue.Len(); i++ {
		destValue := destSliceElem.Index(i)
		valItemActual := srcSliceValue.Index(i)
		if valItemActual.Kind() == reflect.Pointer || valItemActual.Kind() == reflect.Interface {
			if valItemActual.IsNil() {
				continue
			}
			valItemActual = valItemActual.Elem()
		}
		srcItemValue := valItemActual.FieldByName(key)
		origKind := srcItemValue.Type().Kind()
		destTypeOf := destValue.Type()
		destKind := destTypeOf.Kind()
		if origKind == reflect.Int {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(int)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(int)))
			}
		} else if origKind == reflect.Int8 {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(int8)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(int8)))
			}
		} else if origKind == reflect.Int16 {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(int16)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(int16)))
			}

		} else if origKind == reflect.Int32 {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(int32)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(int32)))
			}
		} else if origKind == reflect.Int64 {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(int64)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(int64)))
			}

		} else if origKind == reflect.Uint {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uint)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uint)))
			}

		} else if origKind == reflect.Uint8 {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uint8)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uint8)))
			}
		} else if origKind == reflect.Uint16 {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uint16)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uint16)))
			}

		} else if origKind == reflect.Uint32 {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uint32)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uint32)))
			}

		} else if origKind == reflect.Uint64 {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uint64)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uint64)))
			}
		} else if origKind == reflect.Uintptr {
			if IsUnsignedIntegerKind(destKind) {
				destValue.SetUint(uint64(srcItemValue.Interface().(uintptr)))
			} else {
				destValue.SetInt(int64(srcItemValue.Interface().(uintptr)))
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
	return XSliceCopyToMapCommon(srcSlice, destMap, keyTag, valueTag, false, options...)
}

// 泛型函数，转换一个类型的Slice到另一个类型的Map，使用xref库进行转换。quoteSameType:slice和map相同类型时候是否直接应用而不是使用refx复制。
func XSliceCopyToMapCommon(srcSlice interface{}, destMap interface{}, keyTag string, valueTag string, sameTypeQuote bool, options ...XrefOption) error {
	tag := ":XSliceCopyToMapCommon"
	if nil == destMap {
		return errors.New(tag + "destMap must not nil")
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	//isTidy := xTagContainKey(cpOpt, xRef_key_tidy)
	checkUnsigned := xTagContainKey(do.copyOption, xRef_key_unsigned)
	destMapValue := reflect.ValueOf(destMap)
	if destMapValue.Kind() != reflect.Ptr && destMapValue.Elem().Kind() != reflect.Map {
		return errors.New(tag + "destMap must be a map pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destMapElem := destMapValue.Elem()
	// 获取map的键(key)类型
	keyType := reflect.TypeOf(destMapElem.Interface()).Key()
	keyKind := keyType.Kind()
	keyTypeName := keyType.Name()

	// 获取map的值(value)类型
	valueType := reflect.TypeOf(destMapElem.Interface()).Elem()
	valueKind := valueType.Kind()
	valueTypeName := valueType.String()
	valueIsPrt := false
	var valueActualType reflect.Type = nil
	if valueKind == reflect.Pointer || valueKind == reflect.Interface {
		valueIsPrt = true
		valueActualType = valueType.Elem()
	} else {
		valueActualType = valueType
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return errors.New(tag + "srcSlice must be a slice")
	}
	destMapElem.Set(reflect.MakeMap(destMapElem.Type()))

	var errG error = nil
	for i := 0; i < srcSliceValue.Len(); i++ {
		valItemActual := srcSliceValue.Index(i)
		if valItemActual.Kind() == reflect.Pointer || valItemActual.Kind() == reflect.Interface {
			if valItemActual.IsNil() {
				continue
			}
			valItemActual = valItemActual.Elem()
		}
		keyItem := valItemActual.FieldByName(keyTag)
		keyItemValue := keyItem.Interface()
		keyRtValue, keyTransOk, keyTransErr := xRefMap_transOrigToDestValue(valueTag, do.copyOption, keyItemValue, keyKind, keyTypeName, checkUnsigned)
		if keyTransErr != nil {
			errG = keyTransErr
		}
		if !keyTransOk {
			errG = errors.New(tag + "excute failed")
		}
		if keyRtValue == nil {
			continue
		}
		if len(do.mapKeyAppend) > 0 && keyKind == reflect.String {
			keyRtValue = do.mapKeyAppend + keyRtValue.(string)
		}
		if len(valueTag) > 0 {
			valItem := valItemActual.FieldByName(valueTag)
			if sameTypeQuote && valItem.Type() == valueType {
				// 赋值
				destMapElem.SetMapIndex(reflect.ValueOf(xCopyToMapTypeTrans(keyKind, keyRtValue)), valItem)
			} else {
				valItemValue := valItem.Interface()
				valRtValue, valueTransOk, valueTransErr := xRefMap_transOrigToDestValue(valueTag, do.copyOption, valItemValue, valueKind, valueTypeName, checkUnsigned)
				if valueTransErr != nil {
					errG = valueTransErr
				}
				if !valueTransOk {
					errG = errors.New(tag + "excute failed")
				}
				// 赋值
				destMapElem.SetMapIndex(reflect.ValueOf(xCopyToMapTypeTrans(keyKind, keyRtValue)), reflect.ValueOf(xCopyToMapTypeTrans(valueKind, valRtValue)))
			}
		} else {
			valItem := srcSliceValue.Index(i)
			if sameTypeQuote && valueType == valItem.Type() {
				// 赋值
				destMapElem.SetMapIndex(reflect.ValueOf(xCopyToMapTypeTrans(keyKind, keyRtValue)), valItem)
			} else {
				rtMapActualValue := reflect.New(valueActualType)
				fmt.Println(rtMapActualValue.String())
				err, _ := XRefValueCopy(valItem.Interface(), rtMapActualValue.Elem(), options...)
				if err != nil {
					errG = err
				}
				// 赋值
				if valueIsPrt {
					destMapElem.SetMapIndex(reflect.ValueOf(xCopyToMapTypeTrans(keyKind, keyRtValue)), rtMapActualValue)
				} else {
					destMapElem.SetMapIndex(reflect.ValueOf(xCopyToMapTypeTrans(keyKind, keyRtValue)), rtMapActualValue.Elem())
				}

			}
		}

	}
	return errG
}

// 解析来源字段值为目标待赋值字段
func xRefMap_transOrigToDestValue(key string, cpOpt string, origValue interface{}, destActualTypeKind reflect.Kind, destTypeName string, checkUnsigned bool) (interface{}, bool, error) {

	isTidy := xTagContainKey(cpOpt, xRef_key_tidy)
	//if xRef_log {
	//	fmt.Println(fmt.Sprintf("来源类型:%d-%s,目标类型:%d-%s,Tidy:%t", origKind, origType, destKind, destTypeName, isTidy))
	//}
	if IsIntegerKind(destActualTypeKind) {
		return xParseToInt(key, origValue, destTypeName, destActualTypeKind, cpOpt, isTidy, checkUnsigned)
	} else if IsFloatKind(destActualTypeKind) {
		parseVal, parseFlag := xParseToFloat(key, origValue, destTypeName, destActualTypeKind, cpOpt, isTidy)
		return parseVal, parseFlag, nil
	} else if destActualTypeKind == reflect.Bool {
		parseVal, parseFlag := xParseToBool(key, origValue, destTypeName, destActualTypeKind, cpOpt, isTidy)
		return parseVal, parseFlag, nil
	} else if destActualTypeKind == reflect.String {
		parseVal, parseFlag := xParseToString(key, origValue, destTypeName, destActualTypeKind, cpOpt, isTidy)
		return parseVal, parseFlag, nil
	} else if IsTimeTypeByName(destTypeName) {
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
	} else if destActualTypeKind == reflect.Uintptr {
		viInt64 := vi.(int64)
		return uintptr(viInt64)
	} else {
		return vi
	}

}
