/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/24 13:13
 * @version 1.0
 */
package corex

import (
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

/*
*
判断是否闰年
*/
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

// 解析查询的开始时间
func QueryDayParseStart(queryStart string) string {
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

// 解析查询的结束时间
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
	} else {
		return ""
	}
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
