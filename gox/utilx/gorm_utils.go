/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/3/26 10:33
 * @version 1.0
 */
package utilx

import (
	"github.com/morrisxyang/xreflect"
	"github.com/ruomm/goxframework/gox/corex"
	"github.com/ruomm/goxframework/gox/refx"
	"reflect"
	"strconv"
	"strings"
	"time"
)
import "errors"

const (
	db_args_separator_space       = " ?"
	db_args_in_space              = " (?)"
	db_args_like_space            = " concat('%',?,'%')"
	db_args_equal_space           = " = ?"
	xRef_key_grom_order_by        = "xorderby"
	xRef_key_grom_order_by_table  = "table"
	xRef_key_grom_order_by_option = "opt"
)

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
		subTag, _ := corex.ParseTagToNameOptionFenHao(tagGorm)
		if len(subTag) > 0 && strings.HasPrefix(subTag, "-") {
			return false
		}
		// 判断是否需要选定特定字段
		if nil != selectKeys && len(selectKeys) > 0 {
			if xGormIsContainKey(field.Name, selectKeys...) {
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
		if key == "Version" || key == "version" {
			vi, _ := refx.ParseToInt64(value.Interface(), "")
			if vi == nil {
				continue
			}
			viInt64 := vi.(int64)
			if viInt64 <= 0 {
				continue
			}
			mapresult[key] = viInt64 + 1
		} else {
			mapresult[key] = value.Interface()
		}
	}
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
		subTag, _ := corex.ParseTagToNameOptionFenHao(tagGorm)
		if len(subTag) > 0 && strings.HasPrefix(subTag, "-") {
			return false
		}
		// 排除部分字段
		if xGormIsContainKey(field.Name, ignoreKeys...) {
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
		if key == "Version" || key == "version" {
			vi, _ := refx.ParseToInt64(value.Interface(), "")
			if vi == nil {
				continue
			}
			viInt64 := vi.(int64)
			if viInt64 <= 0 {
				continue
			}
			mapresult[key] = viInt64 + 1
		} else {
			mapresult[key] = value.Interface()
		}
	}
	return mapresult, nil
}

func xGormIsContainKey(fieldName string, fieldKeys ...string) bool {
	if fieldKeys == nil || len(fieldKeys) <= 0 {
		return false
	}
	if len(fieldName) <= 0 {
		return false
	}
	fieldNameLower := strings.ToLower(fieldName)
	fieldNameSnake := corex.ToSnakeCase(fieldName)
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

func ParseConditionMap(conditionMap map[string]interface{}) (string, []interface{}) {
	return ParseConditionMapWithTable(conditionMap, "", false)
}

func ParseConditionMapWithTable(conditionMap map[string]interface{}, tableName string, deleteAtNotNull bool) (string, []interface{}) {
	conditionKey := ""
	var conditionArgs []interface{}
	for keyTemp, value := range conditionMap {
		key := strings.TrimSpace(keyTemp)
		countDbArgs := xGromParseDbArgsCount(key)
		if len(key) > 0 {
			if len(conditionKey) > 0 {
				conditionKey = conditionKey + " and "
			}
			if countDbArgs == 1 {
				conditionKey = conditionKey + xGromParseKeyName(key, tableName)
				conditionArgs = append(conditionArgs, value)
			} else if countDbArgs > 1 {
				conditionKey = conditionKey + xGromParseKeyName(key, tableName)
				sliceLen, _ := xGromParseSliceValueLen(value)
				if sliceLen <= 0 {
					conditionArgs = append(conditionArgs, value)
				} else {
					srcSliceValue := reflect.ValueOf(value)
					for i := 0; i < sliceLen; i++ {
						conditionArgs = append(conditionArgs, srcSliceValue.Index(i).Interface())
					}
				}
			} else if strings.HasSuffix(key, "=") || strings.HasSuffix(key, ">") || strings.HasSuffix(key, "<") {
				conditionKey = conditionKey + xGromParseKeyName(key, tableName) + db_args_separator_space
				conditionArgs = append(conditionArgs, value)
			} else if xGromIsConditionFuncName(key, "is") {
				conditionKey = conditionKey + xGromParseKeyName(key, tableName) + db_args_separator_space
				conditionArgs = append(conditionArgs, value)
			} else if xGromIsConditionFuncName(key, "not") {
				conditionKey = conditionKey + xGromParseKeyName(key, tableName) + db_args_separator_space
				conditionArgs = append(conditionArgs, value)
			} else if xGromIsConditionFuncName(key, "in") {
				conditionKey = conditionKey + xGromParseKeyName(key, tableName) + db_args_in_space
				conditionArgs = append(conditionArgs, value)
			} else if xGromIsConditionFuncName(key, "like") {
				conditionKey = conditionKey + xGromParseKeyName(key, tableName) + db_args_like_space
				conditionArgs = append(conditionArgs, value)
			} else if xGromIsConditionBaseKey(key) {
				conditionKey = conditionKey + xGromParseKeyName(key, tableName) + db_args_equal_space
				conditionArgs = append(conditionArgs, value)
			} else {
				conditionKey = conditionKey + key
			}
		}
	}
	if deleteAtNotNull {
		if len(conditionKey) > 0 {
			conditionKey = conditionKey + " and "
		}
		conditionKey = conditionKey + xGromParseKeyName("deleted_at", tableName) + " is NULL"
	}
	return conditionKey, conditionArgs
}

func xGromParseSliceValueLen(srcSlice interface{}) (int, error) {
	if nil == srcSlice {
		return -1, errors.New("srcSlice is nil")
	}
	srcSliceValue := reflect.ValueOf(srcSlice)
	if srcSliceValue.Kind() != reflect.Slice {
		return -1, errors.New("srcSlice must be a slice")
	}
	return srcSliceValue.Len(), nil
}

func xGromParseDbArgsCount(key string) int {
	characterCount := 0
	for _, char := range key {
		if char == '?' {
			characterCount++
		}
	}
	return characterCount
}

func xGromParseKeyName(key string, tableName string) string {
	if len(tableName) <= 0 {
		return key
	}
	if strings.HasPrefix(key, "(") {
		return key
	}
	indexArgs := strings.Index(key, "?")
	if indexArgs > 0 {
		return key
	}
	index := strings.Index(key, ".")
	if index > 0 {
		return key
	} else {
		return tableName + "." + key
	}
}

func xGromIsConditionFuncName(key string, funcNameUpper string) bool {
	keyLower := strings.ToLower(key)
	lenKeyLower := len(keyLower)
	lenFunc := len(funcNameUpper)
	if lenFunc <= 0 {
		return false
	}
	if lenKeyLower <= lenFunc+1 {
		return false
	}
	funcName := strings.ToLower(funcNameUpper)
	funKey := keyLower[lenKeyLower-lenFunc-1:]
	if strings.HasSuffix(funKey, funcName) {
		if strings.TrimSpace(funKey) == funcName {
			return true
		}
	}
	return false
}

func xGromIsConditionBaseKey(key string) bool {
	keyLower := strings.ToLower(key)
	lenKeyLower := len(keyLower)
	if lenKeyLower <= 0 {
		return false
	}
	_, _, foundSpace := strings.Cut(key, " ")
	if foundSpace {
		return false
	}
	_, _, foundTab := strings.Cut(key, "\t")
	if foundTab {
		return false
	}
	return true
}

// 获取数据库记录列表记录的MD5值-依据更新时间
func GetResultListMD5ByCreatedAt(modelList interface{}) string {
	// 获取全部更新字段转换成version版本
	var updateVersionArray []int
	refx.XSliceCopyByKey(modelList, &updateVersionArray, "CreatedAt")
	updateVersionStr := ""
	for _, updateVersion := range updateVersionArray {
		updateVersionStr = updateVersionStr + "," + strconv.Itoa(updateVersion)
	}
	if len(updateVersionStr) <= 0 {
		return ""
	} else {
		updateVersionMd5 := corex.GetMd5(updateVersionStr)
		return updateVersionMd5
	}
}

// 获取数据库记录列表记录的MD5值-依据Version版本
func GetResultListMD5ByVersion(modelList interface{}) string {
	// 获取全部更新字段转换成version版本
	var updateVersionArray []int
	refx.XSliceCopyByKey(modelList, &updateVersionArray, "Version")
	updateVersionStr := ""
	for _, updateVersion := range updateVersionArray {
		updateVersionStr = updateVersionStr + "," + strconv.Itoa(updateVersion)
	}
	if len(updateVersionStr) <= 0 {
		return ""
	} else {
		updateVersionMd5 := corex.GetMd5(updateVersionStr)
		return updateVersionMd5
	}
}

func IsContainsId(idList []uint, id uint) bool {
	if nil == idList {
		return false
	}
	containsIdFlag := false
	for _, tmpId := range idList {
		if tmpId == id {
			containsIdFlag = true
			break
		}
	}
	return containsIdFlag
}

// 获取数据库记录列表记录的MD5值
func GetResultListMD5(modelList interface{}, versionKey string) string {
	// 获取全部更新字段转换成version版本
	var updateVersionArray []int64
	refx.XSliceCopyByKey(modelList, &updateVersionArray, versionKey)
	updateVersionStr := ""
	for _, updateVersion := range updateVersionArray {
		updateVersionStr = updateVersionStr + "," + strconv.FormatInt(updateVersion, 10)
	}
	if len(updateVersionStr) <= 0 {
		return ""
	} else {
		updateVersionMd5 := corex.GetMd5(updateVersionStr)
		return updateVersionMd5
	}
}

func GormParseQueryDay(time *time.Time) string {
	stateDateString := corex.TimeFormatByString(corex.TIME_PATTERN_STANDARD_DAY, time)
	return stateDateString
}

func GormParseQueryDayFirstInMonth(time *time.Time) string {
	stateDateString := corex.TimeFormatByString(corex.TIME_PATTERN_STANDARD_MONTH, time)
	return stateDateString + "-01"
}

func GormParseQueryZeroTimeInDay(stateDateString string) *time.Time {
	chargingTime, _ := corex.TimeParseByString(corex.TIME_PATTERN_STANDARD, stateDateString+" 00:00:00")
	return chargingTime
}

func GormParseQueryStart(queryStart string) string {
	lenQuery := len(queryStart)
	if lenQuery <= 0 {
		return ""
	} else if lenQuery == 7 {
		return queryStart + "-01" + " 00:00:00"
	} else if lenQuery == 10 {
		return queryStart + " 00:00:00"
	} else {
		return ""
	}
}
func GormParseQueryEnd(queryEnd string) string {
	lenQuery := len(queryEnd)
	if lenQuery <= 0 {
		return ""
	} else if lenQuery == 7 {
		timeArr := strings.Split(queryEnd, "-")
		year := corex.StrToInt64(timeArr[0])
		month := corex.StrToInt64(timeArr[1])
		if month >= 12 {
			year = year + 1
			month = 1
		} else {
			month = month + 1
		}
		return corex.Int64ToStrFill(year, 4) + "-" + corex.Int64ToStrFill(month, 2) + "-01" + " 00:00:00"
	} else if lenQuery == 10 {
		timeArr := strings.Split(queryEnd, "-")
		year := corex.StrToInt64(timeArr[0])
		month := corex.StrToInt64(timeArr[1])
		day := corex.StrToInt64(timeArr[2])
		isLastDayInMonth := corex.IsLastDayInMonth(int(year), int(month), int(day))
		if isLastDayInMonth {
			if month >= 12 {
				year = year + 1
				month = 1
				day = 1
			} else {
				month = month + 1
				day = 1
			}
		} else {
			day = day + 1
		}
		return corex.Int64ToStrFill(year, 4) + "-" + corex.Int64ToStrFill(month, 2) + "-" + corex.Int64ToStrFill(day, 2) + " 00:00:00"
	} else {
		return ""
	}
}

type xGromOrderByTag struct {
	Index     int
	FieldName string
	Table     string
	Opt       string
}

type XOrderBy struct {
	SortField int  // 排序字段索引 1.编号(ID)排序 2.创建时间(CreatedAt)排序 3.更新时间(UpdatedAt)排序 >=4.其他自定义字段排序，参考说明中的排序编号说明
	SortDesc  bool // 是否降序排序 true：降序 false：升序
}

// 解析gorm排序规则,如是tableName传入"-"则依据model解析tableName，注解含有table:=-则依据model解析tableName
func GormParseOrderByID(model interface{}, tableName string, xOrderByList []XOrderBy, sortDesc bool) string {
	return GormParseOrderBy(model, tableName, xOrderByList, &XOrderBy{SortField: 1, SortDesc: sortDesc})
}

// 解析gorm排序规则,如是tableName传入"-"则依据model解析tableName，注解含有table:=-则依据model解析tableName
func GormParseOrderByCreatedAt(model interface{}, tableName string, xOrderByList []XOrderBy, sortDesc bool) string {
	return GormParseOrderBy(model, tableName, xOrderByList, &XOrderBy{SortField: 2, SortDesc: sortDesc})
}

// 解析gorm排序规则,如是tableName传入"-"则依据model解析tableName，注解含有table:=-则依据model解析tableName
func GormParseOrderByUpdatedAt(model interface{}, tableName string, xOrderByList []XOrderBy, sortDesc bool) string {
	return GormParseOrderBy(model, tableName, xOrderByList, &XOrderBy{SortField: 3, SortDesc: sortDesc})
}

// 解析gorm排序规则,如是tableName传入"-"则依据model解析tableName，注解含有table:=-则依据model解析tableName
func GormParseOrderBy(model interface{}, tableName string, xOrderByList []XOrderBy, xOrderByDefault *XOrderBy) string {
	orderByList := xOrderByList
	if nil != xOrderByDefault && xOrderByDefault.SortField > 0 {
		if !corex.SliceContainsByKey(orderByList, "SortField", xOrderByDefault.SortField) {
			orderByList = append(orderByList, *xOrderByDefault)
		}
	}
	if nil == orderByList || len(orderByList) <= 0 {
		return ""
	}
	// 解析排序tag
	orderByMap := xGormParseOrderByTag(model, tableName)
	if nil == orderByMap || len(orderByMap) <= 0 {
		return ""
	}
	var build strings.Builder
	for _, xOrderBy := range orderByList {
		orderByItem := parseOrderByItem(&orderByMap, &xOrderBy)
		if len(orderByItem) <= 0 {
			continue
		}
		if build.Len() > 0 {
			build.WriteString(",")
		}
		build.WriteString(orderByItem)
	}
	return build.String()
}

// 解析gorm排序规则,如是tableName传入"-"则依据model解析tableName，注解含有table:=-则依据model解析tableName
func xGormParseOrderByTag(model interface{}, tableName string) map[int]xGromOrderByTag {
	tableNameDefault := ""
	if tableName == "-" {
		tableNameDefault = corex.ToSnakeCase(XreflectNameToSimply(reflect.TypeOf(model).String()))
	} else {
		tableNameDefault = corex.ToSnakeCase(tableName)
	}
	orderByMap := make(map[int]xGromOrderByTag)
	orderByMapDefault := make(map[int]xGromOrderByTag)
	xreflect.SelectFieldsDeep(model, func(s string, field reflect.StructField, value reflect.Value) bool {
		fieldName := XreflectNameToSimply(field.Name)
		if len(fieldName) <= 0 {
			return false
		}
		if fieldName == "ID" {
			gromOrderBy := xGromOrderByTag{
				Index:     1,
				Table:     tableNameDefault,
				FieldName: fieldName,
				Opt:       "",
			}
			orderByMapDefault[gromOrderBy.Index] = gromOrderBy
			return true
		} else if fieldName == "CreatedAt" {
			gromOrderBy := xGromOrderByTag{
				Index:     2,
				Table:     tableNameDefault,
				FieldName: fieldName,
				Opt:       "",
			}
			orderByMapDefault[gromOrderBy.Index] = gromOrderBy
			return true
		} else if fieldName == "UpdatedAt" {
			gromOrderBy := xGromOrderByTag{
				Index:     3,
				Table:     tableNameDefault,
				FieldName: fieldName,
				Opt:       "",
			}
			orderByMapDefault[gromOrderBy.Index] = gromOrderBy
			return true
		}
		tagXreft, okXreft := field.Tag.Lookup(xRef_key_grom_order_by)
		if !okXreft {
			return false
		}
		// 开始分割目标控制和属性控制
		// 开始分割目标控制和属性控制
		tagOrigVal, tagOpt := corex.ParseTagToNameOptionFenHao(tagXreft)
		if len(tagOrigVal) <= 0 {
			return false
		}
		i, err := strconv.ParseInt(tagOrigVal, 10, 64)
		if err != nil {
			return false
		}
		if i < 0 {
			return false
		}
		tagOption := corex.TagOptions(tagOpt)
		tagOrderByTable := tagOption.OptionValue(xRef_key_grom_order_by_table)
		if tagOrderByTable == "-" {
			tagOrderByTable = tableNameDefault
		} else if len(tagOrderByTable) > 0 {
			tagOrderByTable = corex.ToSnakeCase(tagOrderByTable)
		}
		tagOrderByOpt := strings.ToLower(tagOption.OptionValue(xRef_key_grom_order_by_option))
		if "desc" != tagOrderByOpt && "asc" != tagOrderByOpt {
			tagOrderByOpt = ""
		}
		gromOrderBy := xGromOrderByTag{
			Index:     int(i),
			Table:     tagOrderByTable,
			FieldName: fieldName,
			Opt:       tagOrderByOpt,
		}
		orderByMap[gromOrderBy.Index] = gromOrderBy
		return true
	})

	for key, value := range orderByMapDefault {
		_, exitOk := orderByMap[key]
		if exitOk {
			continue
		} else {
			orderByMap[key] = value
		}
	}
	return orderByMap
}

func parseOrderByItem(pOrderByMap *map[int]xGromOrderByTag, xOrderBy *XOrderBy) string {
	orderByItem := ""
	if nil == pOrderByMap || nil == xOrderBy || xOrderBy.SortField <= 0 {
		return orderByItem
	}
	orderByMap := *pOrderByMap
	_, exitOk := orderByMap[xOrderBy.SortField]
	if !exitOk {
		return orderByItem
	}
	xOrderByTag := orderByMap[xOrderBy.SortField]
	if len(xOrderByTag.FieldName) <= 0 {
		return orderByItem
	}
	if len(xOrderByTag.Table) > 0 {
		if strings.HasSuffix(xOrderByTag.Table, ".") {
			orderByItem = xOrderByTag.Table + corex.ToSnakeCase(xOrderByTag.FieldName)
		} else {
			orderByItem = xOrderByTag.Table + "." + corex.ToSnakeCase(xOrderByTag.FieldName)
		}
	} else {
		orderByItem = corex.ToSnakeCase(xOrderByTag.FieldName)
	}
	if len(xOrderByTag.Opt) > 0 {
		orderByItem = orderByItem + " " + xOrderByTag.Opt
	} else if xOrderBy.SortDesc {
		orderByItem = orderByItem + " " + "desc"
	} else {
		orderByItem = orderByItem + " " + "asc"
	}
	return orderByItem

}

// field名称简化
func XreflectNameToSimply(reflectName string) string {
	lenFieldName := len(reflectName)
	if lenFieldName <= 0 {
		return reflectName
	}
	lastIndex := strings.LastIndex(reflectName, ".")
	if lastIndex >= 0 && lastIndex < lenFieldName-1 {
		return reflectName[lastIndex+1:]
	} else {
		return reflectName
	}
}
