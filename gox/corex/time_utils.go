/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/24 13:13
 * @version 1.0
 */
package corex

import "time"

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

/**
 * 依据时间生成文件名称
 * @param pTime      时间
 * @param filehead      文件头
 * @param filetype      文件结尾
 * @param outTimeFormat SimpleDateFormat的格式，默认yyyyMMdd_HHmmss格式
 * @return 生成的文件名称
 */
func TimeToFileName(pTime *time.Time, filehead string, filetype string, fileTimeLayout string) string {
	var timeFile *time.Time = nil
	if pTime == nil {
		timeNow := time.Now()
		timeFile = &timeNow
	} else {
		timeFile = pTime
	}
	timeLayout := fileTimeLayout
	if timeLayout == "" {
		timeLayout = "20060102_150405"
	}
	dataStr := TimeFormatByString(timeLayout, timeFile)
	return filehead + dataStr + filetype
}
