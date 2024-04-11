/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/4/11 08:53
 * @version 1.0
 */
package timex

import (
	"errors"
	"github.com/ruomm/goxframework/gox/corex"
	"math"
	"strings"
	"time"
)

// 获取剩余时间，返回总计月份和秒数
func ParseRemainMonthSenconds(t1 *time.Time, t2 *time.Time) (int64, int64, error) {
	if t1 == nil || t2 == nil {
		return 0, 0, errors.New("time start is nil or time end is nil")
	}
	var tStart *time.Time = nil
	var tEnd *time.Time = nil
	if t1.Unix() > t2.Unix() {
		tStart = t2
		tEnd = t1
	} else {
		tStart = t1
		tEnd = t2
	}
	// 解析开始的日
	tStartString := corex.TimeFormatByString(corex.TIME_PATTERN_STANDARD_DAY, tStart)
	tEndString := corex.TimeFormatByString(corex.TIME_PATTERN_STANDARD_DAY, tEnd)
	tStartZero, err := corex.TimeParseByString(corex.TIME_PATTERN_STANDARD, tStartString[0:7]+"-01 00:00:00")
	if err != nil {
		return 0, 0, err
	}
	tEndZero, err := corex.TimeParseByString(corex.TIME_PATTERN_STANDARD, tEndString[0:7]+"-01 00:00:00")
	if err != nil {
		return 0, 0, err
	}
	tStartMonths := xGetMonthsForParseRemain(tStartString)
	tEndMonths := xGetMonthsForParseRemain(tEndString)
	offsetSecondStart := tStart.Unix() - tStartZero.Unix()
	offsetSecondEnd := tEnd.Unix() - tEndZero.Unix()
	// 结束的月总计秒数大于开始总计秒数直接返回
	if offsetSecondEnd >= offsetSecondStart {
		return tEndMonths - tStartMonths, offsetSecondEnd - offsetSecondStart, nil
	}
	tStartNextMonth, err := corex.TimeParseByString(corex.TIME_PATTERN_STANDARD, corex.TimeNextMonth(tStart)+"-01 00:00:00")
	if err != nil {
		return 0, 0, err
	}
	offsetSecondStart = tStartNextMonth.Unix() - tStart.Unix()
	return tEndMonths - tStartMonths - 1, offsetSecondStart + offsetSecondEnd, nil
}

func ParseRemainFormat(t1 *time.Time, t2 *time.Time, zeroShow bool) (string, error) {
	months, seconds, err := ParseRemainMonthSenconds(t1, t2)
	if err != nil {
		return "", err
	}
	year, month, day, hour, minute := xGetParseRemainShow(months, seconds)
	var build strings.Builder
	//build.WriteString(gpuBrand)
	//build.WriteString("_")
	//build.WriteString(gpuModel)
	//build.WriteString("_")
	//build.WriteString(strconv.Itoa(cardMemory))
	//build.WriteString("G")
	//return build.String()
	if year <= 0 && month <= 0 && day <= 0 {
		if hour <= 0 && minute <= 0 {
			build.WriteString("0" + "分钟")
		} else {
			if hour > 0 {
				build.WriteString(corex.Int64ToStr(hour) + "小时")
			} else {
				if zeroShow && build.Len() > 0 {
					build.WriteString(corex.Int64ToStr(hour) + "小时")
				}
			}

			if minute > 0 {
				build.WriteString(corex.Int64ToStr(minute) + "分钟")
			} else {
				if zeroShow && build.Len() > 0 {
					build.WriteString(corex.Int64ToStr(minute) + "分钟")
				}
			}
		}
	} else {
		if year > 0 {
			build.WriteString(corex.Int64ToStr(year) + "年")
		}
		if month > 0 {
			build.WriteString(corex.Int64ToStr(month) + "月")
		} else {
			if zeroShow && build.Len() > 0 {
				build.WriteString(corex.Int64ToStr(month) + "月")
			}
		}
		if day > 0 {
			build.WriteString(corex.Int64ToStr(day) + "天")
		} else {
			if zeroShow && build.Len() > 0 {
				build.WriteString(corex.Int64ToStr(day) + "天")
			}
		}
		if hour > 0 {
			build.WriteString(corex.Int64ToStr(hour) + "小时")
		} else {
			if zeroShow && build.Len() > 0 {
				build.WriteString(corex.Int64ToStr(hour) + "小时")
			}
		}
	}
	return build.String(), nil
}

func xGetMonthsForParseRemain(timeFormatString string) int64 {
	timeArr := strings.Split(timeFormatString, "-")
	year := corex.StrToInt64(timeArr[0])
	month := corex.StrToInt64(timeArr[1])
	return year*12 + month
}

func xGetParseRemainShow(months int64, seconds int64) (int64, int64, int64, int64, int64) {
	year := months / 12
	month := months % 12
	day := seconds / (3600 * 24)
	seconds = seconds - day*3600*24
	hour := seconds / 3600
	seconds = seconds - hour*3600
	minute := int64(math.Round(float64(seconds) / float64(60)))
	// 不足一天的时间和分钟显示
	if year <= 0 && month <= 0 && day <= 0 {
		if minute >= 60 {
			hour = hour + 1
			minute = 0
		}
		if hour >= 24 {
			day = day + 1
			hour = 0
		}
		return year, month, day, hour, minute
	}
	// 大于一天的年月日时间显示
	if minute > 30 {
		hour = hour + 1
	}
	if hour >= 24 {
		day = day + 1
		hour = 0
	}
	return year, month, day, hour, 0
}
