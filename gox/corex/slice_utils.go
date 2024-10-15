/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/4/18 09:59
 * @version 1.0
 */
package corex

import (
	"reflect"
	"strings"
)

// 判断Slice是否含有重复元素
func SliceDuplicates(srcSlice interface{}) bool {

	if nil == srcSlice {
		return false
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return false
	}
	// 开始判断重复
	lenSlice := srcSliceValue.Len()
	if lenSlice <= 1 {
		return false
	}
	itemDuplicate := false
	for i := 0; i < lenSlice; i++ {
		tmpItemValueI := srcSliceValue.Index(i).Interface()
		for j := i + 1; j < lenSlice; j++ {
			tmpItemValueJ := srcSliceValue.Index(j).Interface()
			tmpEqual := false
			if tmpItemValueI == nil && tmpItemValueJ == nil {
				tmpEqual = true
			} else if tmpItemValueI == nil || tmpItemValueJ == nil {
				tmpEqual = false
			} else if tmpItemValueI == tmpItemValueJ {
				tmpEqual = true
			} else {
				tmpEqual = false
			}
			if tmpEqual {
				itemDuplicate = true
				break
			}
		}
	}
	return itemDuplicate
}

// 判断Slice是否含有重复元素
func SliceDuplicatesByKey(srcSlice interface{}, keys ...string) bool {

	if nil == srcSlice {
		return false
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return false
	}
	// 开始判断重复
	lenSlice := srcSliceValue.Len()
	if lenSlice <= 1 {
		return false
	}
	if len(keys) <= 0 || (len(keys) == 1 && len(keys[0]) <= 0) {
		itemDuplicate := false
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).Interface()
			for j := i + 1; j < lenSlice; j++ {
				tmpItemValueJ := srcSliceValue.Index(j).Interface()
				tmpEqual := false
				if tmpItemValueI == nil && tmpItemValueJ == nil {
					tmpEqual = true
				} else if tmpItemValueI == nil || tmpItemValueJ == nil {
					tmpEqual = false
				} else if tmpItemValueI == tmpItemValueJ {
					tmpEqual = true
				} else {
					tmpEqual = false
				}
				if tmpEqual {
					itemDuplicate = true
					break
				}
			}
		}
		return itemDuplicate
	} else {
		itemDuplicate := false
		for i := 0; i < lenSlice; i++ {
			mapI := make(map[string]any, 0)
			for _, key := range keys {
				if len(key) <= 0 {
					continue
				}
				tmpItemValueI := srcSliceValue.Index(i).FieldByName(key).Interface()
				mapI[key] = tmpItemValueI
			}
			for j := i + 1; j < lenSlice; j++ {
				tmpDuplicate := true
				for _, key := range keys {
					if len(key) <= 0 {
						continue
					}
					tmpItemValueI := srcSliceValue.Index(i).FieldByName(key).Interface()
					tmpItemValueJ := srcSliceValue.Index(j).FieldByName(key).Interface()
					tmpEqual := false
					if tmpItemValueI == nil && tmpItemValueJ == nil {
						tmpEqual = true
					} else if tmpItemValueI == nil || tmpItemValueJ == nil {
						tmpEqual = false
					} else if tmpItemValueI == tmpItemValueJ {
						tmpEqual = true
					} else {
						tmpEqual = false
					}
					if !tmpEqual {
						tmpDuplicate = false
					}
				}
				if tmpDuplicate {
					itemDuplicate = true
					break
				}
			}
		}
		return itemDuplicate
	}
}

// 判断Slice是否含有特定元素
func SliceContains(srcSlice interface{}, values ...interface{}) bool {
	if nil == srcSlice {
		return false
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return false
	}
	// 开始判断重复
	lenSlice := srcSliceValue.Len()
	if lenSlice <= 0 {
		return false
	}
	itemContains := false
	for i := 0; i < lenSlice; i++ {
		tmpItemValueI := srcSliceValue.Index(i).Interface()
		lenValues := len(values)
		for j := 0; j < lenValues; j++ {
			value := values[j]
			if tmpItemValueI == value {
				itemContains = true
				break
			}
		}

	}
	return itemContains
}

// 判断Slice是否含有特定元素
func SliceContainsByKey(srcSlice interface{}, key string, values ...interface{}) bool {
	if nil == srcSlice {
		return false
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return false
	}
	// 开始判断重复
	lenSlice := srcSliceValue.Len()
	if lenSlice <= 0 {
		return false
	}
	itemContains := false
	if len(key) <= 0 {
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).Interface()
			lenValues := len(values)
			for j := 0; j < lenValues; j++ {
				value := values[j]
				if tmpItemValueI == value {
					itemContains = true
					break
				}
			}
		}
	} else {
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).FieldByName(key).Interface()
			lenValues := len(values)
			for j := 0; j < lenValues; j++ {
				value := values[j]
				if tmpItemValueI == value {
					itemContains = true
					break
				}
			}
		}
	}
	return itemContains
}

// 判断Slice是否含有特定元素
func SliceIndexOf(srcSlice interface{}, values ...interface{}) int {
	if nil == srcSlice {
		return -1
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return -1
	}
	// 开始判断重复
	lenSlice := srcSliceValue.Len()
	if lenSlice <= 0 {
		return -1
	}
	itemIndex := -1
	for i := 0; i < lenSlice; i++ {
		tmpItemValueI := srcSliceValue.Index(i).Interface()
		lenValues := len(values)
		for j := 0; j < lenValues; j++ {
			value := values[j]
			if tmpItemValueI == value {
				itemIndex = i
				break
			}
		}

	}
	return itemIndex
}

// 判断Slice是否含有特定元素
func SliceIndexOfByKey(srcSlice interface{}, key string, values ...interface{}) int {
	if nil == srcSlice {
		return -1
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return -1
	}
	// 开始判断重复
	lenSlice := srcSliceValue.Len()
	if lenSlice <= 0 {
		return -1
	}
	itemIndex := -1
	if len(key) <= 0 {
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).Interface()
			lenValues := len(values)
			for j := 0; j < lenValues; j++ {
				value := values[j]
				if tmpItemValueI == value {
					itemIndex = i
					break
				}
			}
		}
	} else {
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).FieldByName(key).Interface()
			lenValues := len(values)
			for j := 0; j < lenValues; j++ {
				value := values[j]
				if tmpItemValueI == value {
					itemIndex = i
					break
				}
			}
		}
	}
	return itemIndex
}

// 获取Slice中的特定元素
func SliceFindValue(srcSlice interface{}, key string, values ...interface{}) interface{} {
	if nil == srcSlice {
		return nil
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return nil
	}
	// 开始判断重复
	lenSlice := srcSliceValue.Len()
	if lenSlice <= 0 {
		return nil
	}
	var itemContain interface{} = nil
	if len(key) <= 0 {
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).Interface()
			lenValues := len(values)
			for j := 0; j < lenValues; j++ {
				value := values[j]
				if tmpItemValueI == value {
					itemContain = tmpItemValueI
					break
				}
			}
		}
	} else {
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).FieldByName(key).Interface()
			lenValues := len(values)
			for j := 0; j < lenValues; j++ {
				value := values[j]
				if tmpItemValueI == value {
					itemContain = srcSliceValue.Index(i).Interface()
					break
				}
			}
		}
	}
	return itemContain
}

// 获取Slice中的特定元素
func SliceFindFieldValue(srcSlice interface{}, key string, destKey string, values ...interface{}) interface{} {
	if nil == srcSlice {
		return nil
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return nil
	}
	// 开始判断重复
	lenSlice := srcSliceValue.Len()
	if lenSlice <= 0 {
		return nil
	}
	var itemContain interface{} = nil
	if len(key) <= 0 {
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).Interface()
			lenValues := len(values)
			for j := 0; j < lenValues; j++ {
				value := values[j]
				if tmpItemValueI == value {
					if len(destKey) > 0 {
						itemContain = srcSliceValue.Index(i).FieldByName(destKey).Interface()
					} else {
						itemContain = tmpItemValueI
					}
					break
				}
			}
		}
	} else {
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).FieldByName(key).Interface()
			lenValues := len(values)
			for j := 0; j < lenValues; j++ {
				value := values[j]
				if tmpItemValueI == value {
					if len(destKey) > 0 {
						itemContain = srcSliceValue.Index(i).FieldByName(destKey).Interface()
					} else {
						itemContain = srcSliceValue.Index(i).Interface()
					}
					break
				}
			}
		}
	}
	return itemContain
}

// 判断Slice是否仅仅含有特定元素
func SliceOnlyContains(srcSlice interface{}, values ...interface{}) bool {
	if nil == srcSlice {
		return false
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return false
	}
	// 开始判断重复
	lenSlice := srcSliceValue.Len()
	if lenSlice <= 0 {
		return true
	}
	for i := 0; i < lenSlice; i++ {
		itemContains := false
		tmpItemValueI := srcSliceValue.Index(i).Interface()
		lenValues := len(values)
		for j := 0; j < lenValues; j++ {
			value := values[j]
			if tmpItemValueI == value {
				itemContains = true
			}
		}
		if !itemContains {
			return false
		}
	}
	return true
}

// 判断Slice是否仅仅含有特定元素
func SliceOnlyContainsByKey(srcSlice interface{}, key string, values ...interface{}) bool {
	if nil == srcSlice {
		return false
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return false
	}
	// 开始判断重复
	lenSlice := srcSliceValue.Len()
	if lenSlice <= 0 {
		return false
	}
	if len(key) <= 0 {
		itemContains := false
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).Interface()
			lenValues := len(values)
			for j := 0; j < lenValues; j++ {
				value := values[j]
				if tmpItemValueI == value {
					itemContains = true
				}
			}
			if !itemContains {
				return false
			}
		}
	} else {
		itemContains := false
		for i := 0; i < lenSlice; i++ {
			tmpItemValueI := srcSliceValue.Index(i).FieldByName(key).Interface()
			lenValues := len(values)
			for j := 0; j < lenValues; j++ {
				value := values[j]
				if tmpItemValueI == value {
					itemContains = true
				}
			}
			if !itemContains {
				return false
			}
		}
	}
	return true
}

// 转换string为slice
func StringToSlice(str string, sep string, emptyRetain bool) []string {
	srcSlice := strings.Split(str, sep)
	if emptyRetain {
		return srcSlice
	}
	var resultSlice []string
	for _, tmp := range srcSlice {
		if len(tmp) > 0 {
			resultSlice = append(resultSlice, tmp)
		}
	}
	if resultSlice == nil {
		resultSlice = make([]string, 0)
	}
	return resultSlice

}
