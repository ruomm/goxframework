/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/3/28 12:58
 * @version 1.0
 */
package refx

import (
	"errors"
	"reflect"
)

// 泛型函数，转换一个类型的Slice到另一个类型的Slice，使用xref库进行转换。
func XSliceCopy(srcSlice interface{}, destSlice interface{}) error {
	if nil == destSlice {
		return errors.New("destSlice must not nil")
	}
	destValue := reflect.ValueOf(destSlice)
	if destValue.Kind() != reflect.Ptr && destValue.Elem().Kind() != reflect.Slice {
		return errors.New("destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destValue.Elem()
	srcValue := reflect.ValueOf(srcSlice)
	if srcValue.Kind() != reflect.Slice {
		return errors.New("srcSlice must be a slice")
	}
	// 如果目标切片容量不足，则扩展其容量
	if destSliceElem.Cap() < srcValue.Len() {
		destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), srcValue.Len(), srcValue.Len()))
	}
	destSliceElem.SetLen(srcValue.Len())
	var errG error = nil
	for i := 0; i < srcValue.Len(); i++ {
		itemValue := destSliceElem.Index(i)
		err, _ := XRefValueCopy(srcValue.Index(i).Interface(), itemValue)
		if err != nil {
			errG = err
		}
	}
	return errG
}

// 泛型函数，转换一个类型的Slice到另一个类型的Slice，使用xref库进行转换。
func XSliceCopyByKey(srcSlice interface{}, destSlice interface{}, key string) error {
	if nil == destSlice {
		return errors.New("destSlice must not nil")
	}
	destValue := reflect.ValueOf(destSlice)
	if destValue.Kind() != reflect.Ptr && destValue.Elem().Kind() != reflect.Slice {
		return errors.New("destSlice must be a slice pointer")
	}
	if nil == srcSlice {
		return nil
	}
	destSliceElem := destValue.Elem()
	srcValue := reflect.ValueOf(srcSlice)
	if srcValue.Kind() != reflect.Slice {
		return errors.New("srcSlice must be a slice")
	}
	// 如果目标切片容量不足，则扩展其容量
	if destSliceElem.Cap() < srcValue.Len() {
		destSliceElem.Set(reflect.MakeSlice(destSliceElem.Type(), srcValue.Len(), srcValue.Len()))
	}
	destSliceElem.SetLen(srcValue.Len())
	var errG error = nil
	for i := 0; i < srcValue.Len(); i++ {
		itemValue := destSliceElem.Index(i)
		srcItemValue := srcValue.Index(i).FieldByName(key)
		kindSrc := srcItemValue.Type().Kind()
		kindDest := itemValue.Type().Kind()
		if kindSrc == reflect.Int {
			if xIsUnsignedIntegerKind(kindDest) {
				itemValue.SetUint(uint64(srcItemValue.Interface().(int)))
			} else {
				itemValue.SetInt(int64(srcItemValue.Interface().(int)))
			}
		} else if kindSrc == reflect.Int8 {
			if xIsUnsignedIntegerKind(kindDest) {
				itemValue.SetUint(uint64(srcItemValue.Interface().(int8)))
			} else {
				itemValue.SetInt(int64(srcItemValue.Interface().(int8)))
			}
		} else if kindSrc == reflect.Int16 {
			if xIsUnsignedIntegerKind(kindDest) {
				itemValue.SetUint(uint64(srcItemValue.Interface().(int16)))
			} else {
				itemValue.SetInt(int64(srcItemValue.Interface().(int16)))
			}

		} else if kindSrc == reflect.Int32 {
			if xIsUnsignedIntegerKind(kindDest) {
				itemValue.SetUint(uint64(srcItemValue.Interface().(int32)))
			} else {
				itemValue.SetInt(int64(srcItemValue.Interface().(int32)))
			}
		} else if kindSrc == reflect.Int64 {
			if xIsUnsignedIntegerKind(kindDest) {
				itemValue.SetUint(uint64(srcItemValue.Interface().(int64)))
			} else {
				itemValue.SetInt(int64(srcItemValue.Interface().(int64)))
			}

		} else if kindSrc == reflect.Uint {
			if xIsUnsignedIntegerKind(kindDest) {
				itemValue.SetUint(uint64(srcItemValue.Interface().(uint)))
			} else {
				itemValue.SetInt(int64(srcItemValue.Interface().(uint)))
			}

		} else if kindSrc == reflect.Uint8 {
			if xIsUnsignedIntegerKind(kindDest) {
				itemValue.SetUint(uint64(srcItemValue.Interface().(uint8)))
			} else {
				itemValue.SetInt(int64(srcItemValue.Interface().(uint8)))
			}
		} else if kindSrc == reflect.Uint16 {
			if xIsUnsignedIntegerKind(kindDest) {
				itemValue.SetUint(uint64(srcItemValue.Interface().(uint16)))
			} else {
				itemValue.SetInt(int64(srcItemValue.Interface().(uint16)))
			}

		} else if kindSrc == reflect.Uint32 {
			if xIsUnsignedIntegerKind(kindDest) {
				itemValue.SetUint(uint64(srcItemValue.Interface().(uint32)))
			} else {
				itemValue.SetInt(int64(srcItemValue.Interface().(uint32)))
			}

		} else if kindSrc == reflect.Uint64 {
			if xIsUnsignedIntegerKind(kindDest) {
				itemValue.SetUint(uint64(srcItemValue.Interface().(uint64)))
			} else {
				itemValue.SetInt(int64(srcItemValue.Interface().(uint64)))
			}
		} else if kindSrc == reflect.Float32 {
			itemValue.SetFloat(float64(srcItemValue.Interface().(float32)))
		} else if kindSrc == reflect.Float64 {
			itemValue.SetFloat(float64(srcItemValue.Interface().(float64)))
		} else if kindSrc == reflect.String {
			itemValue.SetString(srcItemValue.Interface().(string))
		} else if kindSrc == reflect.Bool {
			itemValue.SetBool(srcItemValue.Interface().(bool))
		} else {
			itemValue.Set(srcItemValue)
		}
	}
	return errG
}
