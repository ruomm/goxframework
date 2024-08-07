/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:02
 * @version 1.0
 */
package refxstandard

import (
	"errors"
	"fmt"
	"github.com/morrisxyang/xreflect"
	"github.com/ruomm/goxframework/gox/corex"
	"reflect"
	"strings"
)

/**
xref可以来源控制
	`xref:"Name,User:UserName,Role:RoleName"`表示如是来源模型以User开始或结束则从UserName来赋值，来源模型以Role开始或结束则从RoleName来赋值，其他来源从Name来赋值。
其他控制参数：
bs：字符串转换为int类型时候按照存储空间模式计算，可以转换kb、mb、gb、tb的单位。
tns：字符串转换为int类型时候以秒为单位按照时间模式计算，可以转换ms、s、m、h、d、w、mon、y的单位。
tnm：字符串转换为int类型时候以毫秒为单位按照时间模式计算，可以转换ms、s、m、h、d、w、mon、y的单位。
t：时间类型和int、float类型相互转换时候的单位设置，默认毫秒，可选参数sec、min、hour、day、mil、mic、nano。
tf：字符串和时间类型相互转换时候的格式化设置，默认：yyyy-MM-dd HH:mm:ss格式。
p：Float类型转换成字符串时候保留小数位数，默认不设置。
snb：字符串转换成int类型时候，true解析为1，false解析为0，字符串转换成boolean类型时候，大于0的解析为true，小于0的解析为false。
z8：字符串转为数字类型时候，以0开头的字符串以8进制进行解析。0x固定以16进制解析。
.bymt：定义字段使用方法复制，可以使用复制方法复制
bymv：定义字段使用方法复制时候是否使用value模式
tomt：定义字段转换方法，可以使用转换方法赋值
tomv：定义字段转换方法时候是否使用value模式

完整示例如下：
`xref:"Name,User:UserName,Role:ParseRoleName.bymt;bs,tns,tnm,t:sec,tf:2006-01-02 15:04:05,p:2,snb,z8,bymv,tomt:TransMethodInt,tomv"`
*/

const (
	xRef_AD_Zero_Second = int64(-62135596800)
	xRef_log            = false
	//xRef_tag_cp_xreft   = "xref"

	//var xReflect_location, _ = time.LoadLocation("Asia/Shanghai")

	xRef_time_layout = "2006-01-02 15:04:05"

	xRef_key_origvalue_by_method     = ".bymt"
	xRef_len_origvalue_by_method     = len(xRef_key_origvalue_by_method)
	xRef_key_getvalue_by_method_mode = "bymv"
	xRef_key_time_t                  = "t"
	xRef_key_bytesize                = "bs"
	xRef_key_timenumber_seconds      = "tns"
	xRef_key_timenumber_millis       = "tnm"
	xRef_key_zero_to_8               = "z8"
	xRef_key_string_bool_number      = "snb"
	xRef_key_time_tf                 = "tf"
	xRef_key_number_point            = "p"
	xRef_key_method_trans            = "tomt"
	xRef_key_method_trans_value_mode = "tomv"

	// 如是omitempty参数存在，来源的数字类型的0、bool类型的false、字符串类型的空、时间类型的0或负数不会赋值的目标里面
	xRef_key_tidy        = "tidy"
	xRef_key_slice_split = "split"
)

// key目标字段的key值，origKey源字段的key值
// 返回需要往目标里面注入的值和时候有错误发生
type XrefHandler func(origKey string, key string, cpOpt string) (interface{}, error)

/*
*
origO：源结构体
destO：目标切片，不可以传入结构体
*/
// TransferObj 将origO对象的属性值转成destO对象的属性值，属性对应关系和控制指令通过`xref`标签指定
// 无标签的如果再按属性名匹配
func XRefStructCopy(origO interface{}, destO interface{}, options ...XrefOption) (error, []string) {
	if nil == origO {
		return errors.New("XRefStructCopy error,source interface is nil"), nil
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	xRef_tag_cp_xreft := xParseRefTagName(do.optTag)
	// 获取origO的类名称
	origT := reflect.TypeOf(origO)
	var origNameSpace string
	if origT.Kind() == reflect.Array || origT.Kind() == reflect.Chan || origT.Kind() == reflect.Map || origT.Kind() == reflect.Pointer || origT.Kind() == reflect.Slice {
		origNameSpace = origT.Elem().String()
	} else {
		origNameSpace = origT.String()
	}
	origNameSpace = xParseRefNameSpace(do.optNameSpace, origNameSpace)
	resOpt := make(map[string]string)
	resOrig := make(map[string]string)
	reflectValueMap, errG := xreflect.SelectFieldsDeep(destO, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagXreft, okXreft := field.Tag.Lookup(xRef_tag_cp_xreft)
		if !okXreft {
			tagOrigByDo, okCanXrefByDo := xParseOrigCopyKeyByXrefOptions(s, &field, &do)
			if !okCanXrefByDo {
				return false
			}
			resOrig[s] = tagOrigByDo
			// 解析属性控制
			tagOptByDo := do.copyOption
			resOpt[s] = tagOptByDo
			if xRef_log {
				fmt.Println("解析复制字段，目标：" + s + "，来源：" + tagOrigByDo + "，控制协议：" + tagOptByDo)
			}
			return true
		}
		// 开始分割目标控制和属性控制
		tagOrigVal, tagOpt := corex.ParseTagToNameOptionFenHao(tagXreft)
		// 判断是否可以进行复制
		tagOrig, okCanXref := xReflect_canXCopy(tagOrigVal, origNameSpace)
		if !okCanXref {
			return false
		}
		resOrig[s] = tagOrig
		// 解析属性控制
		resOpt[s] = string(tagOpt)
		if xRef_log {
			fmt.Println("解析复制字段，目标：" + s + "，来源：" + tagOrig + "，控制协议：" + string(tagOpt))
		}
		return true
	})
	if errG != nil {
		return errG, nil
	}
	var transFailsKeys []string = nil
	for key, value := range reflectValueMap {
		var origKey string
		if resOrig[key] != "" {
			origKey = resOrig[key]
		} else {
			origKey = key
		}
		//origValue, tmpErr01 := xreflect.EmbedFieldValue(origO, origKey)
		cpOpt := resOpt[key]
		var origValue interface{} = nil
		var tmpErr01 error = nil
		if strings.HasSuffix(origKey, xRef_key_origvalue_by_method) && len(origKey) > xRef_len_origvalue_by_method {
			origValue, tmpErr01 = xGetOrigValueByMethod(origO, origKey[0:len(origKey)-xRef_len_origvalue_by_method], cpOpt)
		} else {
			origValue, tmpErr01 = xreflect.EmbedFieldValue(origO, origKey)
		}
		if tmpErr01 != nil {
			transFailsKeys = append(transFailsKeys, key)
			errG = tmpErr01
			continue
		}
		if origValue == nil {
			if xRef_log {
				fmt.Println(key + "字段无需赋值，来源字段值为nil。")
			}
			continue
		}
		//cpOpt := resOpt[key]
		method_trans := xTagFindValueByKey(cpOpt, xRef_key_method_trans)
		if len(method_trans) > 0 {
			origValueByMethod, errByMethod := xParseOrigValueByMethod(method_trans, cpOpt, origValue, destO)
			if errByMethod != nil {
				if xRef_log {
					fmt.Println(key + errByMethod.Error())
				}
				continue
			} else {
				origValue = origValueByMethod
			}
		}
		rtVal, transOk, transErr := xRef_transOrigToDestValue(key, cpOpt, origValue, value, do.checkUnsigned)
		if transErr != nil {
			errG = transErr
		}
		if !transOk {
			transFailsKeys = append(transFailsKeys, key)
		}
		if rtVal == nil {
			continue
		}
		tmpErr02 := xreflect.SetEmbedField(destO, key, rtVal)
		if tmpErr02 != nil {
			transFailsKeys = append(transFailsKeys, key)
			errG = tmpErr02
		}

	}
	return errG, transFailsKeys
}

/*
*
origMap：源map数据
destO：目标切片，不可以传入结构体
*/
// TransferObj 将origO对象的属性值转成destO对象的属性值，属性对应关系和控制指令通过`xref`标签指定
// 无标签的如果再按属性名匹配
func XRefMapCopy(origMap map[string]string, destO interface{}, options ...XrefOption) (error, []string) {
	if nil == origMap {
		return errors.New("XRefStructCopy error,source map is nil"), nil
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	xRef_tag_cp_xreft := xParseRefTagName(do.optTag)
	origNameSpace := xParseRefNameSpace(do.optNameSpace, "")
	resOpt := make(map[string]string)
	resOrig := make(map[string]string)
	reflectValueMap, errG := xreflect.SelectFieldsDeep(destO, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagXreft, okXreft := field.Tag.Lookup(xRef_tag_cp_xreft)
		if !okXreft {
			tagOrigByDo, okCanXrefByDo := xParseOrigCopyKeyByXrefOptions(s, &field, &do)
			if !okCanXrefByDo {
				return false
			}
			resOrig[s] = tagOrigByDo
			// 解析属性控制
			tagOptByDo := do.copyOption
			resOpt[s] = tagOptByDo
			if xRef_log {
				fmt.Println("解析复制字段，目标：" + s + "，来源：" + tagOrigByDo + "，控制协议：" + tagOptByDo)
			}
			return true
		}
		// 开始分割目标控制和属性控制
		tagOrigVal, tagOpt := corex.ParseTagToNameOptionFenHao(tagXreft)
		// 判断是否可以进行复制
		tagOrig, okCanXref := xReflect_canXCopy(tagOrigVal, origNameSpace)
		if !okCanXref {
			return false
		}
		resOrig[s] = tagOrig
		// 解析属性控制
		resOpt[s] = string(tagOpt)
		if xRef_log {
			fmt.Println("解析复制字段，目标：" + s + "，来源：" + tagOrig + "，控制协议：" + string(tagOpt))
		}
		return true
	})
	if errG != nil {
		return errG, nil
	}
	var transFailsKeys []string = nil
	for key, value := range reflectValueMap {
		var origKey string
		if resOrig[key] != "" {
			origKey = resOrig[key]
		} else {
			origKey = key
		}

		origValueStr, ok := origMap[origKey]
		if !ok {
			if xRef_log {
				fmt.Println(key + "字段无需赋值，来源字段值为空。")
			}
			continue
		}
		cpOpt := resOpt[key]
		var origValue interface{}
		origValue = origValueStr
		method_trans := xTagFindValueByKey(cpOpt, xRef_key_method_trans)
		if len(method_trans) > 0 {
			origValueByMethod, errByMethod := xParseOrigValueByMethod(method_trans, cpOpt, origValue, destO)
			if errByMethod != nil {
				if xRef_log {
					fmt.Println(key + errByMethod.Error())
				}
				continue
			} else {
				origValue = origValueByMethod
			}
		}
		rtVal, transOk, transErr := xRef_transOrigToDestValue(key, cpOpt, origValue, value, do.checkUnsigned)
		if transErr != nil {
			errG = transErr
		}
		if !transOk {
			transFailsKeys = append(transFailsKeys, key)
		}
		if rtVal == nil {
			continue
		}
		tmpErr02 := xreflect.SetEmbedField(destO, key, rtVal)
		if tmpErr02 != nil {
			transFailsKeys = append(transFailsKeys, key)
			errG = tmpErr02
		}

	}
	return errG, transFailsKeys
}

/*
*
origMap：源map数据
destO：目标切片，不可以传入结构体
*/
// TransferObj 将origO对象的属性值转成destO对象的属性值，属性对应关系和控制指令通过`xref`标签指定
// 无标签的如果再按属性名匹配
func XRefHandlerCopy(xrefOrigHandler XrefHandler, destO interface{}, options ...XrefOption) (error, []string) {
	if nil == xrefOrigHandler {
		return errors.New("XRefStructCopy error,xrefOrigHandler is nil"), nil
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	xRef_tag_cp_xreft := xParseRefTagName(do.optTag)
	origNameSpace := xParseRefNameSpace(do.optNameSpace, "")
	resOpt := make(map[string]string)
	resOrig := make(map[string]string)
	reflectValueMap, errG := xreflect.SelectFieldsDeep(destO, func(s string, field reflect.StructField, value reflect.Value) bool {
		tagXreft, okXreft := field.Tag.Lookup(xRef_tag_cp_xreft)
		if !okXreft {
			tagOrigByDo, okCanXrefByDo := xParseOrigCopyKeyByXrefOptions(s, &field, &do)
			if !okCanXrefByDo {
				return false
			}
			resOrig[s] = tagOrigByDo
			// 解析属性控制
			tagOptByDo := do.copyOption
			resOpt[s] = tagOptByDo
			if xRef_log {
				fmt.Println("解析复制字段，目标：" + s + "，来源：" + tagOrigByDo + "，控制协议：" + tagOptByDo)
			}
			return true
		}
		// 开始分割目标控制和属性控制
		tagOrigVal, tagOpt := corex.ParseTagToNameOptionFenHao(tagXreft)
		// 判断是否可以进行复制
		tagOrig, okCanXref := xReflect_canXCopy(tagOrigVal, origNameSpace)
		if !okCanXref {
			return false
		}
		resOrig[s] = tagOrig
		// 解析属性控制
		resOpt[s] = string(tagOpt)
		if xRef_log {
			fmt.Println("解析复制字段，目标：" + s + "，来源：" + tagOrig + "，控制协议：" + string(tagOpt))
		}
		return true
	})
	if errG != nil {
		return errG, nil
	}
	var transFailsKeys []string = nil
	for key, value := range reflectValueMap {
		var origKey string
		if resOrig[key] != "" {
			origKey = resOrig[key]
		} else {
			origKey = key
		}

		origValue, tmpErr01 := xrefOrigHandler(origKey, key, resOpt[key])
		if tmpErr01 != nil {
			transFailsKeys = append(transFailsKeys, key)
			errG = tmpErr01
			continue
		}
		if origValue == nil {
			if xRef_log {
				fmt.Println(key + "字段无需赋值，来源字段值为nil。")
			}
			continue
		}
		cpOpt := resOpt[key]
		method_trans := xTagFindValueByKey(cpOpt, xRef_key_method_trans)
		if len(method_trans) > 0 {
			origValueByMethod, errByMethod := xParseOrigValueByMethod(method_trans, cpOpt, origValue, destO)
			if errByMethod != nil {
				if xRef_log {
					fmt.Println(key + errByMethod.Error())
				}
				continue
			} else {
				origValue = origValueByMethod
			}
		}
		rtVal, transOk, transErr := xRef_transOrigToDestValue(key, cpOpt, origValue, value, do.checkUnsigned)
		if transErr != nil {
			errG = transErr
		}
		if !transOk {
			transFailsKeys = append(transFailsKeys, key)
		}
		if rtVal == nil {
			continue
		}
		tmpErr02 := xreflect.SetEmbedField(destO, key, rtVal)
		if tmpErr02 != nil {
			transFailsKeys = append(transFailsKeys, key)
			errG = tmpErr02
		}

	}
	return errG, transFailsKeys
}

/*
*
origO：源结构体
refValue：目标数据
*/
// TransferObj 将origO对象的属性值转成destO对象的属性值，属性对应关系和控制指令通过`xref`标签指定
// 无标签的如果再按属性名匹配
func XRefValueCopy(origO interface{}, refValue reflect.Value, options ...XrefOption) (error, []string) {
	if nil == origO {
		return errors.New("XRefStructCopy error,source interface is nil"), nil
	}
	do := xrefOptions{}
	for _, option := range options {
		option.f(&do)
	}
	xRef_tag_cp_xreft := xParseRefTagName(do.optTag)
	// 获取origO的类名称
	origT := reflect.TypeOf(origO)
	var origNameSpace string
	if origT.Kind() == reflect.Array || origT.Kind() == reflect.Chan || origT.Kind() == reflect.Map || origT.Kind() == reflect.Pointer || origT.Kind() == reflect.Slice {
		origNameSpace = origT.Elem().String()
	} else {
		origNameSpace = origT.String()
	}
	origNameSpace = xParseRefNameSpace(do.optNameSpace, origNameSpace)
	resOpt := make(map[string]string)
	resOrig := make(map[string]string)
	reflectValueMap, errG := xreflect.SelectFieldsDeep(refValue.Interface(), func(s string, field reflect.StructField, value reflect.Value) bool {
		tagXreft, okXreft := field.Tag.Lookup(xRef_tag_cp_xreft)
		if !okXreft {
			tagOrigByDo, okCanXrefByDo := xParseOrigCopyKeyByXrefOptions(s, &field, &do)
			if !okCanXrefByDo {
				return false
			}
			resOrig[s] = tagOrigByDo
			// 解析属性控制
			tagOptByDo := do.copyOption
			resOpt[s] = tagOptByDo
			if xRef_log {
				fmt.Println("解析复制字段，目标：" + s + "，来源：" + tagOrigByDo + "，控制协议：" + tagOptByDo)
			}
			return true
		}
		// 开始分割目标控制和属性控制
		tagOrigVal, tagOpt := corex.ParseTagToNameOptionFenHao(tagXreft)
		// 判断是否可以进行复制
		tagOrig, okCanXref := xReflect_canXCopy(tagOrigVal, origNameSpace)
		if !okCanXref {
			return false
		}
		resOrig[s] = tagOrig
		// 解析属性控制
		resOpt[s] = string(tagOpt)
		if xRef_log {
			fmt.Println("解析复制字段，目标：" + s + "，来源：" + tagOrig + "，控制协议：" + string(tagOpt))
		}
		return true
	})
	if errG != nil {
		return errG, nil
	}
	var transFailsKeys []string = nil
	for key, value := range reflectValueMap {
		var origKey string
		if resOrig[key] != "" {
			origKey = resOrig[key]
		} else {
			origKey = key
		}
		//origValue, tmpErr01 := xreflect.EmbedFieldValue(origO, origKey)
		cpOpt := resOpt[key]
		var origValue interface{} = nil
		var tmpErr01 error = nil
		if strings.HasSuffix(origKey, xRef_key_origvalue_by_method) && len(origKey) > xRef_len_origvalue_by_method {
			origValue, tmpErr01 = xGetOrigValueByMethod(origO, origKey[0:len(origKey)-xRef_len_origvalue_by_method], cpOpt)
		} else {
			origValue, tmpErr01 = xreflect.EmbedFieldValue(origO, origKey)
		}
		if tmpErr01 != nil {
			transFailsKeys = append(transFailsKeys, key)
			errG = tmpErr01
			continue
		}
		if origValue == nil {
			if xRef_log {
				fmt.Println(key + "字段无需赋值，来源字段值为nil。")
			}
			continue
		}
		//cpOpt := resOpt[key]
		method_trans := xTagFindValueByKey(cpOpt, xRef_key_method_trans)
		if len(method_trans) > 0 {
			origValueByMethod, errByMethod := xParseOrigValueByMethod(method_trans, cpOpt, origValue, refValue.Interface())
			if errByMethod != nil {
				if xRef_log {
					fmt.Println(key + errByMethod.Error())
				}
				continue
			} else {
				origValue = origValueByMethod
			}
		}
		rtVal, transOk, transErr := xRef_transOrigToDestValue(key, cpOpt, origValue, value, do.checkUnsigned)
		if transErr != nil {
			errG = transErr
		}
		if !transOk {
			transFailsKeys = append(transFailsKeys, key)
		}
		if rtVal == nil {
			continue
		}
		// 获取字段Field1的reflect.Value
		//field := refValue.FieldByName(xParseRefValueKey(key))
		//field, errFindField := xreflect.Field(refValue, xParseRefValueKey(key))
		field, errFindField := xreflect.EmbedField(refValue, key)
		if errFindField != nil {
			errG = errFindField
		}
		if field.CanSet() {
			kind := field.Type().Kind()
			if IsIntegerKind(kind) {
				if IsUnsignedIntegerKind(kind) {
					rtConvert := uint64(rtVal.(int64))
					field.SetUint(rtConvert)
				} else {
					rtConvert := rtVal.(int64)
					field.SetInt(rtConvert)
				}
			} else if IsStringKind(kind) {
				rtConvert := rtVal.(string)
				field.SetString(rtConvert)
			} else if IsFloatKind(kind) {
				rtConvert := rtVal.(float64)
				field.SetFloat(rtConvert)
			} else if kind == reflect.Bool {
				rtConvert := rtVal.(bool)
				field.SetBool(rtConvert)
			} else {
				field.Set(reflect.ValueOf(rtVal))
			}
		}
	}
	return errG, transFailsKeys
}

// 字段是否需要XReflect复制
func xReflect_canXCopy(tagOrigVal string, origNameSpace string) (string, bool) {
	if tagOrigVal == "" {
		return "", true
	}
	cpEnable := false
	cpOrigKey := ""
	tagOriglist := strings.Split(tagOrigVal, ",")
	for _, tagOrigItem := range tagOriglist {
		if tagOrigItem == "" {
			continue
		}
		subVList := strings.Split(tagOrigItem, ":")
		lenVList := len(subVList)
		if lenVList == 0 {
			continue
		} else if lenVList == 1 {
			if len(subVList[0]) > 0 {
				cpOrigKey = subVList[0]
				cpEnable = true
			}
		} else if lenVList == 2 && len(origNameSpace) > 0 {
			//if len(subVList[0]) > 0 && (strings.HasSuffix(origNameSpace, subVList[0]) || strings.HasPrefix(origNameSpace, subVList[0])) {
			if len(subVList[0]) > 0 && strings.Contains(origNameSpace, subVList[0]) {
				cpEnable = true
				cpOrigKey = subVList[1]
				break
			}
		}
	}
	return cpOrigKey, cpEnable
}

// 解析来源字段值为目标待赋值字段
func xRef_transOrigToDestValue(key string, cpOpt string, origValue interface{}, destValue reflect.Value, checkUnsigned bool) (interface{}, bool, error) {
	destTypeOf := destValue.Type()
	destKind := destTypeOf.Kind()
	destTypeName := destTypeOf.String()
	destActualTypeKind := reflect.Invalid
	destActualTypeOf := destTypeOf
	if destKind == reflect.Pointer {
		destActualTypeOf = destTypeOf.Elem()
		destActualTypeKind = destTypeOf.Elem().Kind()
	} else {
		destActualTypeOf = destTypeOf
		destActualTypeKind = destKind
	}

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
	} else if IsTimeTypeByName(destTypeOf.String()) {
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
				parseVal, parseFlag := xParseStringToSlice(key, viString, destTypeName, destActualTypeOf.Elem().Kind(), cpOpt)
				return parseVal, parseFlag, nil
			}
		}

		origTypeOf := reflect.TypeOf(origValue)
		origKind := origTypeOf.Kind()
		origType := origTypeOf.String()
		if origKind != destKind {
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

func xTagContainKey(tagValue string, key string) bool {
	return corex.TagOptions(tagValue).Contains(key)
}

func xTagFindValueByKey(tagValue string, key string) string {
	return corex.TagOptions(tagValue).OptionValue(key)
}

func IsIntegerKind(kind reflect.Kind) bool {
	if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 ||
		kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 || kind == reflect.Uintptr {
		return true
	} else {
		return false
	}
}
func IsUnsignedIntegerKind(kind reflect.Kind) bool {
	if kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 || kind == reflect.Uintptr {
		return true
	} else {
		return false
	}
}
func IsFloatKind(kind reflect.Kind) bool {
	if kind == reflect.Float64 || kind == reflect.Float32 {
		return true
	} else {
		return false
	}
}

func IsStringKind(kind reflect.Kind) bool {
	return kind == reflect.String
}

func xIsIntegerPointer(typeName string) bool {
	if typeName == "*int" || typeName == "*int8" || typeName == "*int16" || typeName == "*int32" || typeName == "*int64" ||
		typeName == "*uint" || typeName == "*uint8" || typeName == "*uint16" || typeName == "*uint32" || typeName == "*uint64" || typeName == "*uintptr" {
		return true
	} else {
		return false
	}
}
func xIsUnsignedIntegerPointer(typeName string) bool {
	if typeName == "*uint" || typeName == "*uint8" || typeName == "*uint16" || typeName == "*uint32" || typeName == "*uint64" || typeName == "*uintptr" {
		return true
	} else {
		return false
	}
}
func xIsFloatPointer(typeName string) bool {
	if typeName == "*float64" || typeName == "*float32" {
		return true
	} else {
		return false
	}
}

//func xIsStringPointer(typeName string) bool {
//	if typeName == "*string" {
//		return true
//	} else {
//		return false
//	}
//}

func IsStructKind(kind reflect.Kind) bool {
	return kind == reflect.Struct
}
func xIsStringTypeByName(typeName string) bool {
	return typeName == "string" || typeName == "*string"
}
func IsTimeTypeByName(typeName string) bool {
	return typeName == "time.Time" || typeName == "Time" || typeName == "*time.Time"
}

func xIsPointor(typeName string) bool {
	return strings.HasPrefix(typeName, "*")
}
