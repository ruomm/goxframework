/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/24 13:13
 * @version 1.0
 */
package corex

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	TIME_PATTERN_STANDARD       string = "2006-01-02 15:04:05"
	TIME_PATTERN_STANDARD_DAY   string = "2006-01-02"
	TIME_PATTERN_STANDARD_MONTH string = "2006-01"
	TIME_PATTERN_STANDARD_MILLS string = "2006-01-02 15:04:05.999"
	TIME_PATTERN_NOSPACE        string = "20060102150405"
	TIME_PATTERN_NOSPACE_DAY    string = "20060102"
	TIME_PATTERN_NOSPACE_MONTH  string = "200601"
	TIME_PATTERN_NOSPACE_MILLS  string = "20060102150405999"
)

// var Time_Location_ShangHai, _ = time.LoadLocation("Asia/Shanghai")
const timezone_location = "Asia/Shanghai"
const timezone_offset int = 8 * 3600

//var Time_Location_ShangHai = ToShanghaiLocation()

// time.Location转换时间为本地时间
func ToTimeLocation() *time.Location {
	// missing Location in call to Date
	location, err := time.LoadLocation(timezone_location)
	if err != nil {
		location = time.FixedZone("CST", 8*3600) //替换上海时区方式
	}
	return location
}

// time.Location转换时间为特定时区时间
func ToTimeLocationPattern(timezoneValue string, timezoneOffset int) *time.Location {
	// missing Location in call to Date
	location, err := time.LoadLocation(timezoneValue)
	if err != nil {
		location = time.FixedZone("CST", timezoneOffset)
	}
	return location
}

// 格式化时间为字符串，timeLayout:格式化模板，mesc:毫秒时间戳
func TimeFormatByMilliSeconds(timeLayout string, msec int64) string {
	//TimeFormatByString
	//"America/Adak" "Asia/Shanghai"
	var realTimeLayout string
	if len(timeLayout) > 0 {
		realTimeLayout = timeLayout
	} else {
		realTimeLayout = TIME_PATTERN_STANDARD
	}
	t := time.UnixMilli(msec)
	return t.In(ToTimeLocation()).Format(realTimeLayout)
}

// 格式化时间为字符串，timeLayout:格式化模板
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
// 解析字符串为时间，timeLayout:格式化模板
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

// 判断时间是否在当前时间之后，timeLayout:格式化模板
func TimeAfterCurrent(timeLayout string, timeStr string) (bool, error) {
	cTime, err := TimeParseByString(timeLayout, timeStr)
	if err != nil {
		return false, err
	}
	return cTime.After(time.Now()), nil
}

// 判断时间是否在当前时间之前，timeLayout:格式化模板
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

// 判断是否闰年
func IsLeapyear(year int) bool {
	if year%4 == 0 && year%100 != 0 {
		return true
	} else if year%400 == 0 {
		return true
	} else {
		return false
	}
}

// 获取一个月有多少天
func GetDayCountByMonth(year int, month int) int {
	daysnumberinmonth := 30
	if month == 1 || month == 3 || month == 5 || month == 7 || month == 8 || month == 10 || month == 12 {
		daysnumberinmonth = 31
	} else if month == 4 || month == 6 || month == 9 || month == 11 {
		daysnumberinmonth = 30
	} else if month == 2 {
		if IsLeapyear(year) {
			daysnumberinmonth = 29
		} else {
			daysnumberinmonth = 28
		}
	}
	return daysnumberinmonth
}

// 判断是否一个月的最右一天
func IsLastDayInMonth(year int, month int, day int) bool {
	daysnumberinmonth := GetDayCountByMonth(year, month)
	if day >= daysnumberinmonth {
		return true
	} else {
		return false
	}
}

// 解析查询的开始时间，支持yyyy-MM,yyyy-MM-dd,yyyy-MM-dd HH:mm:ss格式
// yyyy-MM解析为当月1号0点0分0秒,yyyy-MM-dd解析为当天的0点0分0秒，yyyy-MM-dd HH:mm:ss实际返回
func QueryDayParseStart(queryStart string) string {
	lenQuery := len(queryStart)
	if lenQuery <= 0 {
		return ""
	} else if lenQuery == 7 {
		return queryStart + "-01" + " 00:00:00"
	} else if lenQuery == 10 {
		return queryStart + " 00:00:00"
	} else if lenQuery == 19 {
		return queryStart
	} else {
		return ""
	}
}

// 解析查询的结束时间，支持yyyy-MM,yyyy-MM-dd,yyyy-MM-dd HH:mm:ss格式
// yyyy-MM解析为下月1号0点0分0秒,yyyy-MM-dd解析为下一天的0点0分0秒，yyyy-MM-dd HH:mm:ss实际返回
func QueryDayParseEnd(queryEnd string) string {
	lenQuery := len(queryEnd)
	if lenQuery <= 0 {
		return ""
	} else if lenQuery == 7 {
		timeArr := strings.Split(queryEnd, "-")
		year := StrToInt64(timeArr[0])
		month := StrToInt64(timeArr[1])
		if month >= 12 {
			year = year + 1
			month = 1
		} else {
			month = month + 1
		}
		return Int64ToStrFill(year, 4) + "-" + Int64ToStrFill(month, 2) + "-01" + " 00:00:00"
	} else if lenQuery == 10 {
		timeArr := strings.Split(queryEnd, "-")
		year := StrToInt64(timeArr[0])
		month := StrToInt64(timeArr[1])
		day := StrToInt64(timeArr[2])
		isLastDayInMonth := IsLastDayInMonth(int(year), int(month), int(day))
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
		return Int64ToStrFill(year, 4) + "-" + Int64ToStrFill(month, 2) + "-" + Int64ToStrFill(day, 2) + " 00:00:00"
	} else if lenQuery == 19 {
		return queryEnd
	} else {
		return ""
	}
}

// 获取当前天
func TimeCurrentDay(currentDay *time.Time) string {
	if currentDay == nil {
		return ""
	}
	currentDayString := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, currentDay)
	return currentDayString
}

// 获取前一天的时间
func TimePreDay(currentDay *time.Time) string {
	if currentDay == nil {
		return ""
	}
	currentDayString := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, currentDay)
	timeArr := strings.Split(currentDayString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	day := StrToInt64(timeArr[2])
	if day <= 1 {
		if month <= 1 {
			month = 12
			year = year - 1
		} else {
			month = month - 1
		}
		day = int64(GetDayCountByMonth(int(year), int(month)))
	} else {
		day = day - 1
	}
	preDayTimeString := Int64ToStrFill(year, 4) + "-" + Int64ToStrFill(month, 2) + "-" + Int64ToStrFill(day, 2)
	return preDayTimeString
}

// 获取后一天的时间
func TimeNextDay(currentDay *time.Time) string {
	if currentDay == nil {
		return ""
	}
	currentDayString := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, currentDay)
	timeArr := strings.Split(currentDayString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	day := StrToInt64(timeArr[2])
	isLastDayInMonth := IsLastDayInMonth(int(year), int(month), int(day))
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
	nextDayTimeString := Int64ToStrFill(year, 4) + "-" + Int64ToStrFill(month, 2) + "-" + Int64ToStrFill(day, 2)
	return nextDayTimeString
}

// 获取间隔特定天数的时间
func TimeOffsetDay(currentDay *time.Time, offsetDay int) string {
	if currentDay == nil {
		return ""
	}
	monthStr := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, currentDay)
	timeOffsetString := TimeOffsetDayByString(monthStr, offsetDay)
	return timeOffsetString
	//timeOffset := (*currentDay).AddDate(0, 0, offsetDay)
	//timeOffsetString := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, &timeOffset)
	//return timeOffsetString
}

// 获取前一天的时间
func TimePreDayByString(dayString string) string {
	if !TimeValidDayString(dayString) {
		return ""
	}
	timeArr := strings.Split(dayString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	day := StrToInt64(timeArr[2])
	if day <= 1 {
		if month <= 1 {
			month = 12
			year = year - 1
		} else {
			month = month - 1
		}
		day = int64(GetDayCountByMonth(int(year), int(month)))
	} else {
		day = day - 1
	}
	preDayTimeString := Int64ToStrFill(year, 4) + "-" + Int64ToStrFill(month, 2) + "-" + Int64ToStrFill(day, 2)
	return preDayTimeString
}

// 获取后一天的时间
func TimeNextDayByString(dayString string) string {
	if !TimeValidDayString(dayString) {
		return ""
	}
	timeArr := strings.Split(dayString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	day := StrToInt64(timeArr[2])
	isLastDayInMonth := IsLastDayInMonth(int(year), int(month), int(day))
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
	nextDayTimeString := Int64ToStrFill(year, 4) + "-" + Int64ToStrFill(month, 2) + "-" + Int64ToStrFill(day, 2)
	return nextDayTimeString
}

// 获取间隔特定天数的时间
func TimeOffsetDayByString(dayString string, offsetDay int) string {
	if !TimeValidDayString(dayString) {
		return ""
	}
	currentDay, err := TimeParseByString(TIME_PATTERN_STANDARD, dayString+" 12:00:00")
	if err != nil {
		return ""
	}
	timeOffset := (*currentDay).AddDate(0, 0, offsetDay)
	timeOffsetString := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, &timeOffset)
	return timeOffsetString
}

// 获取当前月
func TimeCurrentMonth(currentDay *time.Time) string {
	if currentDay == nil {
		return ""
	}
	currentDayString := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, currentDay)
	return currentDayString[0:7]
}

// 获取前一月的时间
func TimePreMonth(currentDay *time.Time) string {
	if currentDay == nil {
		return ""
	}
	currentDayString := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, currentDay)
	timeArr := strings.Split(currentDayString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	if month <= 1 {
		month = 12
		year = year - 1
	} else {
		month = month - 1
	}
	return Int64ToStrFill(year, 4) + "-" + Int64ToStrFill(month, 2)
}

// 获取后一月的时间
func TimeNextMonth(currentDay *time.Time) string {
	if currentDay == nil {
		return ""
	}
	currentDayString := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, currentDay)
	timeArr := strings.Split(currentDayString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	if month >= 12 {
		year = year + 1
		month = 1
	} else {
		month = month + 1
	}
	return Int64ToStrFill(year, 4) + "-" + Int64ToStrFill(month, 2)
}

// 获取间隔特定月数的时间
func TimeOffsetMonth(currentDay *time.Time, offsetMonth int) string {
	if currentDay == nil {
		return ""
	}
	//timeOffset := (*currentDay).AddDate(0, offsetMonth, 0)
	//return TimeFormatByString(TIME_PATTERN_STANDARD_MONTH, &timeOffset)
	monthStr := TimeFormatByString(TIME_PATTERN_STANDARD_MONTH, currentDay)
	timeOffsetString := TimeOffsetMonthByString(monthStr, offsetMonth)
	return timeOffsetString
}

// 获取前一月的时间
func TimePreMonthByString(monthString string) string {
	if !TimeValidMonthString(monthString) && !TimeValidDayString(monthString) {
		return ""
	}
	timeArr := strings.Split(monthString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	if month <= 1 {
		month = 12
		year = year - 1
	} else {
		month = month - 1
	}
	return Int64ToStrFill(year, 4) + "-" + Int64ToStrFill(month, 2)
}

// 获取后一月的时间
func TimeNextMonthByString(monthString string) string {
	if !TimeValidMonthString(monthString) && !TimeValidDayString(monthString) {
		return ""
	}
	timeArr := strings.Split(monthString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	if month >= 12 {
		year = year + 1
		month = 1
	} else {
		month = month + 1
	}
	return Int64ToStrFill(year, 4) + "-" + Int64ToStrFill(month, 2)
}

// 获取间隔特定月数的时间
func TimeOffsetMonthByString(monthString string, offsetMonth int) string {
	dayString := ""
	if TimeValidMonthString(monthString) {
		dayString = monthString + "-15"
	} else if TimeValidDayString(monthString) {
		dayString = monthString[0:7] + "-15"
	} else {
		return ""
	}
	currentDay, err := TimeParseByString(TIME_PATTERN_STANDARD_DAY, dayString)
	if err != nil {
		return ""
	}
	timeOffset := (*currentDay).AddDate(0, offsetMonth, 0)
	timeOffsetString := TimeFormatByString(TIME_PATTERN_STANDARD_MONTH, &timeOffset)
	return timeOffsetString
}

// 格式化年月日转换为int类型的年月日，支持yyyy-MM、yyyy-MM-dd，无法转换返回0、0、0
// 2024-06转换后为2024、6、0，2024-06-07转换后为2024、6、7
func TimeToYearMonthDay(dayString string) (int64, int64, int64) {
	if !TimeValidMonthString(dayString) && !TimeValidDayString(dayString) {
		return 0, 0, 0
	}
	timeArr := strings.Split(dayString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	if len(timeArr) > 2 {
		day := StrToInt64(timeArr[2])
		return year, month, day
	} else {
		return year, month, 0
	}
}

// 判断字符串是否符合yyyy-MM-dd日期格式
func TimeValidDayString(dayString string) bool {
	re, err := regexp.Compile("^\\d{4}-\\d{2}-\\d{2}$")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if !re.MatchString(dayString) {
		return false
	}
	timeArr := strings.Split(dayString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	day := StrToInt64(timeArr[2])
	if year < 0 || year >= 3000 {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}
	dayCountByMonth := GetDayCountByMonth(int(year), int(month))
	if day < 1 || day > int64(dayCountByMonth) {
		return false
	}
	return true
}

// 判断字符串是否符合yyyy-MM月份格式
func TimeValidMonthString(monthString string) bool {
	re, err := regexp.Compile("^\\d{4}-\\d{2}$")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if !re.MatchString(monthString) {
		return false
	}
	timeArr := strings.Split(monthString, "-")
	year := StrToInt64(timeArr[0])
	month := StrToInt64(timeArr[1])
	if year < 0 || year >= 3000 {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}
	return true
}

// 依据开始时间，结束时间查询一共有多少天
func TimeTotalDaysForQuery(timeStart *time.Time, timeEnd *time.Time) int {
	if nil == timeStart || nil == timeEnd {
		return 0
	}
	startDayStr := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, timeStart)
	endDayStr := TimeFormatByString(TIME_PATTERN_STANDARD_DAY, timeEnd)
	resultStart, _ := TimeParseByString(TIME_PATTERN_STANDARD, startDayStr+" 00:00:00")
	resultEnd, _ := TimeParseByString(TIME_PATTERN_STANDARD, endDayStr+" 00:00:00")
	if nil == resultStart || nil == resultEnd {
		return 0
	}
	if resultEnd.Before(*timeStart) {
		return int(resultStart.Sub(*resultEnd).Hours()/24) + 1
	} else {
		return int(resultEnd.Sub(*resultStart).Hours()/24) + 1
	}

}
