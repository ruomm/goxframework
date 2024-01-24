package corex

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	TIME_PATTERN_STANDARD       string = "2006-01-02 15:04:05"
	TIME_PATTERN_STANDARD_DAY   string = "2006-01-02"
	TIME_PATTERN_STANDARD_MILLS string = "2006-01-02 15:04:05.999"
	TIME_PATTERN_NOSPACE        string = "20060102150405"
	TIME_PATTERN_NOSPACE_DAY    string = "20060102"
	TIME_PATTERN_NOSPACE_MILLS  string = "20060102150405999"
)

// var Time_Location_ShangHai, _ = time.LoadLocation("Asia/Shanghai")
const timezone_location = "Asia/Shanghai"
const timezone_offset int = 8 * 3600

//var Time_Location_ShangHai = ToShanghaiLocation()

// 转换时间为本地时间
func ToTimeLocation() *time.Location {
	// missing Location in call to Date
	location, err := time.LoadLocation(timezone_location)
	if err != nil {
		location = time.FixedZone("CST", 8*3600) //替换上海时区方式
	}
	return location
}

// 转换时间为特定时区时间
func ToTimeLocationPattern(timezoneValue string, timezoneOffset int) *time.Location {
	// missing Location in call to Date
	location, err := time.LoadLocation(timezoneValue)
	if err != nil {
		location = time.FixedZone("CST", timezoneOffset)
	}
	return location
}

var (
	//TIME_LOCATION_CN, _         = time.LoadLocation("Asia/Shanghai")
	xvalid_matchNonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	xvalid_matchFirstCap        = regexp.MustCompile("(.)([A-Z][a-z]+)")
	xvalid_matchAllCap          = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ServerTrace(tracePort int) {
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(tracePort), nil)
}

func GetMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func GetMd5WithSlat(data, slat string) string {
	var realSlat string
	h := md5.New()
	h.Write([]byte(data + realSlat))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

// 格式化时间为字符串
func TimeFormatByString(timeLayout string, t *time.Time) string {
	//TimeFormatByString
	//"America/Adak" "Asia/Shanghai"
	var realTimeLayout string
	if len(timeLayout) > 0 {
		realTimeLayout = timeLayout
	} else {
		realTimeLayout = TIME_PATTERN_STANDARD
	}
	return t.In(ToTimeLocation()).Format(realTimeLayout)
}

// TimeParseByString
// 解析字符串为时间
func TimeParseByString(timeLayout string, sTime string) (*time.Time, error) {
	var realTimeLayout string
	if len(timeLayout) > 0 {
		realTimeLayout = timeLayout
	} else {
		realTimeLayout = TIME_PATTERN_STANDARD
	}
	timeStamp, err := time.ParseInLocation(realTimeLayout, sTime, ToTimeLocation())
	if err != nil {
		return nil, err
	} else {
		return &timeStamp, nil
	}
}

func TimeAfterCurrent(timeLayout string, timeStr string) (bool, error) {
	cTime, err := TimeParseByString(timeLayout, timeStr)
	if err != nil {
		return false, err
	}
	return cTime.After(time.Now()), nil
}

func TimeBeforceCurrent(timeLayout string, timeStr string) (bool, error) {
	cTime, err := TimeParseByString(timeLayout, timeStr)
	if err != nil {
		return false, err
	}
	return cTime.Before(time.Now()), nil
}

// 驼峰转下划线工具
func ToSnakeCase(str string) string {
	str = xvalid_matchNonAlphaNumeric.ReplaceAllString(str, "_")     //非常规字符转化为 _
	snake := xvalid_matchFirstCap.ReplaceAllString(str, "${1}_${2}") //拆分出连续大写
	snake = xvalid_matchAllCap.ReplaceAllString(snake, "${1}_${2}")  //拆分单词
	return strings.ToLower(snake)                                    //全部转小写
}
func JsonParseByString(str string, v any) error {
	if str == "" {
		return errors.New("json Unmarshal not support empty string")
	}
	err := json.Unmarshal([]byte(str), v)
	return err
}

func JsonFormatByString(v any) (string, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	if len(jsonData) == 0 {
		return "", errors.New("json Marshal not support this object")
	}
	return string(jsonData), err
}
