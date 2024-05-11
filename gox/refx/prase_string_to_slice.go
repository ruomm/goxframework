/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/21 19:53
 * @version 1.0
 */
package refx

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func xParseStringToSlice(key string, origVal string, destTypeName string, destSliceKind reflect.Kind, cpOpt string) (interface{}, bool) {

	splitTag := xTagFindValueByKey(cpOpt, xRef_key_slice_split)
	if len(splitTag) == 0 {
		splitTag = ","
	}
	vi, err := ParseStringToSlice(key, origVal, splitTag, destSliceKind, destTypeName)
	if err != nil {
		return nil, false
	} else {
		return vi, true
	}
}

func ParseStringToSlice(key string, viString string, splitTag string, destSliceKind reflect.Kind, destTypeName string) (interface{}, error) {
	destStrArr := strings.Split(viString, splitTag)
	if destSliceKind == reflect.Int {
		var viArr []int = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseInt(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, int(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Int8 {
		var viArr []int8 = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseInt(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, int8(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Int16 {
		var viArr []int16 = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseInt(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, int16(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Int32 {
		var viArr []int32 = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseInt(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, int32(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Int64 {
		var viArr []int64 = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseInt(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, int64(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Uint {
		var viArr []uint = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseUint(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, uint(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Uint8 {
		var viArr []uint8 = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseUint(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, uint8(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Uint16 {
		var viArr []uint16 = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseUint(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, uint16(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Uint32 {
		var viArr []uint32 = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseUint(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, uint32(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Uint64 {
		var viArr []uint64 = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseUint(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, uint64(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Uintptr {
		var viArr []uintptr = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseUint(tmpStr, 10, 64)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, uintptr(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Float32 {
		var viArr []float32 = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseFloat(tmpStr, 10)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, float32(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Float64 {
		var viArr []float64 = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseFloat(tmpStr, 10)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, float64(tmpVal))
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.Bool {
		var viArr []bool = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			tmpVal, tmpErr := strconv.ParseBool(tmpStr)
			if tmpErr != nil {
				return nil, tmpErr
			} else {
				viArr = append(viArr, tmpVal)
			}
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else if destSliceKind == reflect.String {
		var viArr []string = nil
		for _, tmpStr := range destStrArr {
			if len(tmpStr) <= 0 {
				continue
			}
			viArr = append(viArr, tmpStr)
		}
		if strings.HasPrefix(destTypeName, "*") {
			if nil == viArr {
				return nil, nil
			} else {
				return &viArr, nil
			}
		} else {
			return viArr, nil
		}
	} else {
		return nil, errors.New(fmt.Sprintf("%s字段无法赋值，无法转换字符串为%s类型", key, destTypeName))
	}
}
