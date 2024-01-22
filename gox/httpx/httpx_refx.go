package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/morrisxyang/xreflect"
	"github.com/ruomm/goxframework/gox/corex"
	"github.com/ruomm/goxframework/gox/refx"
	"net/url"
	"reflect"
	"strings"
)

const xRef_log = false

const xRef_tag_key_xurl_param = "xurl_param"

/*
*
srcO：源结构体
destO：目标切片，不可以传入结构体
*/
// TransferObj 将srcO对象的属性值转成destO对象的属性值，属性对应关系和控制指令通过`xref`标签指定
// 无标签的如果再按属性名匹配
func ParseToJSONByte(jsonObj interface{}, transBasePointer ...bool) ([]byte, error) {
	if nil == jsonObj {
		return nil, errors.New("ParseToJSONByte error,source interface is nil")
	}
	srcTypeOf := reflect.TypeOf(jsonObj)
	srcName := srcTypeOf.String()
	// 目标是字符串类型
	if srcName == "string" {
		return []byte(jsonObj.(string)), nil
	} else {
		jsonStr, err := json.Marshal(jsonObj)
		return jsonStr, err
	}
}

/*
*
srcO：源结构体
destO：目标切片，不可以传入结构体
*/
// TransferObj 将srcO对象的属性值转成destO对象的属性值，属性对应关系和控制指令通过`xref`标签指定
// 无标签的如果再按属性名匹配
func ParseToUrlEncodeString(origO interface{}) (string, error) {
	if nil == origO {
		//return "", errors.New("XReflectCopy error,source interface is nil")
		return "", nil
	}
	origTypeOf := reflect.TypeOf(origO)
	origTypeName := origTypeOf.String()
	// 目标是字符串类型
	if origTypeName == "string" {
		return origO.(string), nil
	} else if origTypeName == "*string" {
		if nil == origO.(*string) {
			return "", nil
		} else {
			return *(origO.(*string)), nil
		}

	} else if origTypeName == "url.Values" {
		v := origO.(url.Values)
		if len(v) <= 0 {
			return "", nil
		} else {
			return v.Encode(), nil
		}
	} else if origTypeName == "*url.Values" {
		v := origO.(*url.Values)
		if len(*v) <= 0 {
			return "", nil
		} else {
			return v.Encode(), nil
		}
	}
	resOpt := make(map[string]string)
	resOrig := make(map[string]string)
	reflectValueMap, errG := xreflect.SelectFieldsDeep(origO, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagXreft, okXreft := field.Tag.Lookup(xRef_tag_key_xurl_param)
		if !okXreft {
			return false
		}
		// 开始分割目标控制和属性控制
		subTags := corex.ParseToSubTag(tagXreft)
		// 解析目标控制
		urlKey := ""
		if len(subTags) > 0 {
			urlKey = subTags[0]
		}
		if urlKey == "-" {
			return false
		}
		if urlKey == "" {
			urlKey = strings.ToLower(s[0:1]) + s[1:len(s)]
		}
		resOrig[s] = urlKey
		// 解析属性控制
		tagOpt := ""
		if len(subTags) > 1 {
			tagOpt = subTags[1]
		}
		resOpt[s] = tagOpt
		if xRef_log {
			fmt.Println("解析URL参数字段，目标：" + urlKey + "，来源：" + s + "，控制协议：" + tagOpt)
		}
		return true
	})
	if errG != nil {
		return "", errG
	}
	v := url.Values{}
	for key, _ := range reflectValueMap {
		var srcKey string
		if resOrig[key] != "" {
			srcKey = resOrig[key]
		} else {
			srcKey = key
		}
		srcValue, err := xreflect.EmbedFieldValue(origO, key)
		if err != nil {
			errG = err
			continue
		}
		if srcValue == nil {
			continue
		}
		cpOpt := resOpt[key]
		rtVal := refx.ParseToString(srcValue, cpOpt)

		if rtVal == nil {
			continue
		} else {
			rtString := rtVal.(string)
			if refx.XrefTagTidy(cpOpt) && len(rtString) <= 0 {
				continue
			} else {
				v.Set(srcKey, rtVal.(string))
			}
		}
	}
	if len(v) <= 0 {
		return "", errG
	} else {
		return v.Encode(), errG
	}
}
