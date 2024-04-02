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
	"sort"
	"strconv"
	"strings"
)

const (
	xRef_log                 = false
	xRequest_Parse_Param     = "xreq_param"
	xRequest_Parse_Query     = "xreq_query"
	xRequest_Parse_Header    = "xreq_header"
	xRequest_Option_Order    = "order"
	xRequest_Option_Response = "resp"
)

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
		tagXreft, okXreft := field.Tag.Lookup(xRequest_Parse_Query)
		if !okXreft {
			return false
		}
		// 开始分割目标控制和属性控制
		urlKey, tagOpt := corex.ParseTagToNameOptionFenHao(tagXreft)
		if urlKey == "-" {
			return false
		}
		if urlKey == "" {
			urlKey = strings.ToLower(s[0:1]) + s[1:len(s)]
		}
		resOrig[s] = urlKey
		// 解析属性控制
		resOpt[s] = string(tagOpt)
		if xRef_log {
			fmt.Println("解析URL参数字段，目标：" + urlKey + "，来源：" + s + "，控制协议：" + string(tagOpt))
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
			// if refx.XrefTagTidy(cpOpt) && len(rtString) <= 0 {
			// 修订为空字符串跳过参数获取
			if len(rtString) <= 0 {
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

func ParseToRequest(reqObj interface{}) (string, []byte, string, string, map[string]string, error) {
	if nil == reqObj {
		//return "", errors.New("XReflectCopy error,source interface is nil")
		return "", nil, "", "", nil, errors.New("ParseToRequest error,reqObj interface is nil")
	}
	reqString, reqStrOk := xParseReqToString(reqObj)
	if reqStrOk {
		if len(reqString) <= 0 {
			return "GET", nil, "", "", nil, nil
		} else if xIsJsonString(reqString) {
			return "POST", []byte(reqString), "", "", nil, nil
		} else if xIsUrlString(reqString) {
			return "GET", nil, "", reqString, nil, nil
		} else {
			if strings.HasPrefix(reqString, "/") {
				return "GET", nil, reqString, "", nil, nil
			} else {
				return "GET", nil, "/" + reqString, "", nil, nil
			}
		}
	}
	reqMethod, _ := xParseHttpxMethod(reqObj)
	reqBody, err := xParseReqJson(reqObj)
	if err != nil {
		return "", nil, "", "", nil, err
	}

	reqParam, err := xParseReqParam(reqObj)
	if err != nil {
		return "", nil, "", "", nil, err
	}
	reqQuery, err := xParseReqQuery(reqObj)
	if err != nil {
		return "", nil, "", "", nil, err
	}
	reqHeaderMap, err := xParseReqHeaderMap(reqObj)
	if err != nil {
		return "", nil, "", "", nil, err
	}
	return reqMethod, reqBody, reqParam, reqQuery, reqHeaderMap, nil

}

type HttpxPararmValule struct {
	Order      string
	ParamKey   string
	ParamValue string
}

// 解析请求参数
func xParseHttpxMethod(reqObj interface{}) (string, error) {
	refVals, err := xreflect.CallMethod(reqObj, "HttpxMethod")
	if err != nil {
		//return "POST", err
		return "POST", nil
	}
	if len(refVals) <= 0 {
		return "POST", nil
	}
	httpxMethod := ""
	for _, origVal := range refVals {
		if origVal.Kind() == reflect.String {
			httpxMethod = origVal.String()
		}
		if len(httpxMethod) > 0 {
			break
		}
	}
	if len(httpxMethod) <= 0 {
		return "POST", nil
	} else {
		return strings.ToUpper(httpxMethod), nil
	}

}

// 解析为JSON请求体字符串
func xParseReqJson(reqObj interface{}) ([]byte, error) {
	if reqObj == nil {
		return nil, nil
	}
	jsonData, err := json.Marshal(reqObj)
	if err != nil {
		return nil, err
	}
	jsonDataStr := string(jsonData)
	if len(jsonDataStr) <= 2 {
		return nil, nil
	} else {
		trimJsonStr := strings.TrimSpace(jsonDataStr[1 : len(jsonDataStr)-1])
		if (len(trimJsonStr)) <= 0 {
			return nil, nil
		} else {
			return jsonData, nil
		}
	}

}

// 解析为URI路径字符串
func xParseReqParam(reqObj interface{}) (string, error) {
	resOpt := make(map[string]string)
	resOrig := make(map[string]string)
	reflectValueMap, errG := xreflect.SelectFieldsDeep(reqObj, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagXreft, okXreft := field.Tag.Lookup(xRequest_Parse_Param)
		if !okXreft {
			return false
		}
		// 开始分割目标控制和属性控制
		urlKey, tagOpt := corex.ParseTagToNameOptionFenHao(tagXreft)
		if urlKey == "-" {
			return false
		}
		if urlKey == "" {
			urlKey = strings.ToLower(s[0:1]) + s[1:len(s)]
		}
		resOrig[s] = urlKey
		// 解析属性控制
		resOpt[s] = string(tagOpt)
		if xRef_log {
			fmt.Println("解析URL参数字段，目标：" + urlKey + "，来源：" + s + "，控制协议：" + string(tagOpt))
		}
		return true
	})
	if errG != nil {
		return "", errG
	}
	var paramList []HttpxPararmValule
	orderInt := 0
	for key, _ := range reflectValueMap {
		var srcKey string
		if resOrig[key] != "" {
			srcKey = resOrig[key]
		} else {
			srcKey = key
		}
		srcValue, err := xreflect.EmbedFieldValue(reqObj, key)
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
			// if refx.XrefTagTidy(cpOpt) && len(rtString) <= 0 {
			// 修订为空字符串跳过参数获取
			if len(rtString) <= 0 {
				continue
			} else {
				orderString := xTagFindValueByKey(cpOpt, xRequest_Option_Order)
				if len(orderString) <= 0 {
					orderString = strconv.Itoa(orderInt)
					orderInt = orderInt + 1
				}

				paramList = append(paramList, HttpxPararmValule{
					Order:      orderString,
					ParamKey:   srcKey,
					ParamValue: rtVal.(string),
				})
			}
		}
	}
	if paramList == nil || len(paramList) <= 0 {
		return "", nil
	}
	sort.SliceIsSorted(paramList, func(i, j int) bool {
		iOrder := paramList[i].Order
		jOrder := paramList[j].Order
		orderIndex := strings.Compare(iOrder, jOrder)
		return orderIndex < 0
	})
	paramString := ""
	for _, tmpPararmValule := range paramList {
		paramString = paramString + "/" + tmpPararmValule.ParamValue
	}
	return paramString, nil
}

// 解析为Query请求字符串
func xParseReqQuery(reqObj interface{}) (string, error) {
	resOpt := make(map[string]string)
	resOrig := make(map[string]string)
	reflectValueMap, errG := xreflect.SelectFieldsDeep(reqObj, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagXreft, okXreft := field.Tag.Lookup(xRequest_Parse_Query)
		if !okXreft {
			return false
		}
		// 开始分割目标控制和属性控制
		urlKey, tagOpt := corex.ParseTagToNameOptionFenHao(tagXreft)
		if urlKey == "-" {
			return false
		}
		if urlKey == "" {
			urlKey = strings.ToLower(s[0:1]) + s[1:len(s)]
		}
		resOrig[s] = urlKey
		// 解析属性控制
		resOpt[s] = string(tagOpt)
		if xRef_log {
			fmt.Println("解析URL参数字段，目标：" + urlKey + "，来源：" + s + "，控制协议：" + string(tagOpt))
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
		srcValue, err := xreflect.EmbedFieldValue(reqObj, key)
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
			// if refx.XrefTagTidy(cpOpt) && len(rtString) <= 0 {
			// 修订为空字符串跳过参数获取
			if len(rtString) <= 0 {
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

// 解析为HeaderMap数据
func xParseReqHeaderMap(reqObj interface{}) (map[string]string, error) {
	resOpt := make(map[string]string)
	resOrig := make(map[string]string)
	reflectValueMap, errG := xreflect.SelectFieldsDeep(reqObj, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagXreft, okXreft := field.Tag.Lookup(xRequest_Parse_Header)
		if !okXreft {
			return false
		}
		// 开始分割目标控制和属性控制
		urlKey, tagOpt := corex.ParseTagToNameOptionFenHao(tagXreft)
		if urlKey == "-" {
			return false
		}
		if xTagContainKey(string(tagOpt), xRequest_Option_Response) {
			// 如是响应的字段信息则不设置请求返回
			return false
		}
		if urlKey == "" {
			urlKey = strings.ToLower(s[0:1]) + s[1:len(s)]
		}
		resOrig[s] = urlKey
		// 解析属性控制
		resOpt[s] = string(tagOpt)
		if xRef_log {
			fmt.Println("解析URL参数字段，目标：" + urlKey + "，来源：" + s + "，控制协议：" + string(tagOpt))
		}
		return true
	})
	if errG != nil {
		return nil, errG
	}
	paramMap := make(map[string]string)
	for key, _ := range reflectValueMap {
		var srcKey string
		if resOrig[key] != "" {
			srcKey = resOrig[key]
		} else {
			srcKey = key
		}
		srcValue, err := xreflect.EmbedFieldValue(reqObj, key)
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
			// if refx.XrefTagTidy(cpOpt) && len(rtString) <= 0 {
			// 修订为空字符串跳过参数获取
			if len(rtString) <= 0 {
				continue
			} else {
				paramMap[srcKey] = rtVal.(string)
			}
		}
	}
	if len(paramMap) <= 0 {
		return nil, errG
	} else {
		return paramMap, errG
	}
}

func xTagContainKey(tagValue string, key string) bool {
	return corex.TagOptions(tagValue).Contains(key)
}

func xTagFindValueByKey(tagValue string, key string) string {
	return corex.TagOptions(tagValue).OptionValue(key)
}

func xParseReqToString(reqObj interface{}) (string, bool) {
	origTypeOf := reflect.TypeOf(reqObj)
	origTypeName := origTypeOf.String()
	// 目标是字符串类型
	if origTypeName == "string" {
		return reqObj.(string), true
	} else if origTypeName == "*string" {
		return *(reqObj.(*string)), true
	} else if origTypeName == "url.Values" {
		v := reqObj.(url.Values)
		if len(v) <= 0 {
			return "", true
		} else {
			return v.Encode(), true
		}
	} else if origTypeName == "*url.Values" {
		v := reqObj.(*url.Values)
		if len(*v) <= 0 {
			return "", true
		} else {
			return v.Encode(), true
		}
	} else {
		return "", false
	}
}
