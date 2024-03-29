/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/3/26 10:33
 * @version 1.0
 */
package corex

import (
	"github.com/morrisxyang/xreflect"
	"reflect"
	"strings"
	"time"
)
import "errors"

/*
* 转换gorm模型为MAP对象，不包含gorm模型内置的id和时间相关字段，如是selectKeys有值则只会转换选定的key值字段
 */
func ToGormMap(gormModel interface{}, selectKeys ...string) (map[string]interface{}, error) {
	if nil == gormModel {
		return nil, errors.New("ToGormMap error,gormModel is nil")
	}
	//mapReflectValue, err := xreflect.Fields(destO)
	mapReflectValue, err := xreflect.SelectFields(gormModel, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagGorm, okGorm := field.Tag.Lookup("gorm")
		if !okGorm {
			return false
		}
		// 开始分割目标控制和属性控制
		subTags := ParseToSubTag(tagGorm)
		ignoreFlag := false
		if subTags != nil && len(subTags) > 0 {
			for _, subTag := range subTags {
				if len(subTag) > 0 && strings.HasPrefix(subTag, "-") {
					ignoreFlag = true
					break
				}
			}
		}
		if ignoreFlag {
			return false
		}
		// 判断是否需要选定特定字段
		if nil != selectKeys && len(selectKeys) > 0 {
			if isContainGormKey(field.Name, selectKeys...) {
				return true
			} else {
				return false
			}
		} else {
			return true
		}

	})
	if err != nil {
		return nil, errors.New("To GormMap error,xreflect parse gormModel error")
	}
	if len(mapReflectValue) == 0 {
		return nil, errors.New("To GormMap error,xreflect parse gormModel empty")
	}
	mapresult := make(map[string]interface{})
	for key, value := range mapReflectValue {
		mapresult[key] = value.Interface()
	}
	mapresult["UpdatedAt"] = time.Now()
	return mapresult, nil
}

/*
* 转换gorm模型为MAP对象，不包含gorm模型内置的id和时间相关字段，如是ignorekeys有值则忽略转换选定的key值字段
 */
func ToGormMapIgnoreMode(gormModel interface{}, ignoreKeys ...string) (map[string]interface{}, error) {
	if nil == gormModel {
		return nil, errors.New("ToGormMap error,gormModel is nil")
	}
	//mapReflectValue, err := xreflect.Fields(destO)
	mapReflectValue, err := xreflect.SelectFields(gormModel, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagGorm, okGorm := field.Tag.Lookup("gorm")
		if !okGorm {
			return false
		}
		// 开始分割目标控制和属性控制
		subTags := ParseToSubTag(tagGorm)
		ignoreFlag := false
		if subTags != nil && len(subTags) > 0 {
			for _, subTag := range subTags {
				if len(subTag) > 0 && strings.HasPrefix(subTag, "-") {
					ignoreFlag = true
					break
				}
			}
		}
		if ignoreFlag {
			return false
		}
		// 排除部分字段
		if isContainGormKey(field.Name, ignoreKeys...) {
			return false
		}
		return true

	})
	if err != nil {
		return nil, errors.New("To GormMap error,xreflect parse gormModel error")
	}
	if len(mapReflectValue) == 0 {
		return nil, errors.New("To GormMap error,xreflect parse gormModel empty")
	}
	mapresult := make(map[string]interface{})
	for key, value := range mapReflectValue {
		mapresult[key] = value.Interface()
	}
	mapresult["UpdatedAt"] = time.Now()
	return mapresult, nil
}

func isContainGormKey(fieldName string, fieldKeys ...string) bool {
	if fieldKeys == nil || len(fieldKeys) <= 0 {
		return false
	}
	if len(fieldName) <= 0 {
		return false
	}
	fieldNameLower := strings.ToLower(fieldName)
	fieldNameSnake := ToSnakeCase(fieldName)
	containFlag := false
	for _, key := range fieldKeys {
		if len(key) <= 0 {
			continue
		}
		if fieldNameLower == strings.ToLower(key) {
			containFlag = true
			break
		}
		if fieldNameSnake == strings.ToLower(key) {
			containFlag = true
			break
		}
	}
	return containFlag
}
