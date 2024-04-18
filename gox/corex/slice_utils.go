/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/4/18 09:59
 * @version 1.0
 */
package corex

import (
	"reflect"
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
			if tmpItemValueI == tmpItemValueJ {
				itemDuplicate = true
				break
			}
		}
	}
	return itemDuplicate
}
