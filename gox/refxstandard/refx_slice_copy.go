/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/3/28 12:58
 * @version 1.0
 */
package refxstandard

import (
	"errors"
	"reflect"
)

const (
	xRef_key_unsigned = "unsigned"
)

// 泛型函数，转换一个类型的Slice到另一个类型的Slice，使用xref库进行转换。
func XSliceCopy(srcSlice interface{}, destSlice interface{}) error {
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
		err, _ := XRefValueCopy(srcSliceValue.Index(i).Interface(), destValue)
		if err != nil {
			errG = err
		}
	}
	return errG
}

// 泛型函数，转换一个类型的Slice到另一个类型的Slice，使用xref库进行转换。
func XSliceCopyByKey(srcSlice interface{}, destSlice interface{}, key string, optionTags ...string) error {
	if nil == destSlice {
		return errors.New("destSlice must not nil")
	}
	cpOpt := ""
	if len(optionTags) > 0 {
		for _, tmp := range optionTags {
			if len(tmp) <= 0 {
				continue
			}
			if len(cpOpt) > 0 {
				cpOpt = cpOpt + ","
			}
			cpOpt = cpOpt + tmp
		}
	}
	//isTidy := xTagContainKey(cpOpt, xRef_key_tidy)
	checkUnsigned := xTagContainKey(cpOpt, xRef_key_unsigned)
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
		origValue := srcItemValue.Interface()
		rtVal, transOk, transErr := xRef_transOrigToDestValue(key, cpOpt, origValue, destValue, checkUnsigned)
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
