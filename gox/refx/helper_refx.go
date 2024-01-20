/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:26
 * @version 1.0
 */
package refx

import (
	"github.com/ruomm/goxframework/gox/corex"
	"strconv"
	"strings"
	"time"
)

// 字符串转换为int64
func xTransStringToInt64(viString string, cpOpt string) (int64, error) {
	numBase := 10
	if strings.HasPrefix(viString, "0x") || strings.HasPrefix(viString, "0X") {
		numBase = 16
		viString = viString[2:]
	} else if strings.HasPrefix(viString, "-0x") || strings.HasPrefix(viString, "-0X") {
		numBase = 16
		viString = "-" + viString[3:]
	} else if (strings.HasPrefix(viString, "0") || strings.HasPrefix(viString, "-0")) && xTagContainKey(cpOpt, xRef_key_zero_to_8) {
		numBase = 8
	}
	if strings.HasPrefix(viString, "-") {
		viInt64, err := strconv.ParseInt(viString, numBase, 64)
		if err != nil {
			return 0, err
		} else {
			return viInt64, nil
		}
	} else {
		viUint64, err := strconv.ParseUint(viString, numBase, 64)
		if err != nil {
			if xTagContainKey(cpOpt, xRef_key_string_bool_number) {
				viBool, errB := strconv.ParseBool(viString)
				if errB != nil {
					return 0, errB
				} else if viBool {
					return 1, nil
				} else {
					return 0, nil
				}
			} else {
				return 0, err
			}
		} else {
			return int64(viUint64), nil
		}
	}
}

// 字符串转换为bool
func xTransStringIntToBool(viString string, cpOpt string) (bool, error) {
	numBase := 10
	if strings.HasPrefix(viString, "0x") || strings.HasPrefix(viString, "0X") {
		numBase = 16
		viString = viString[2:]
	} else if strings.HasPrefix(viString, "-0x") || strings.HasPrefix(viString, "-0X") {
		numBase = 16
		viString = "-" + viString[3:]
	} else if (strings.HasPrefix(viString, "0") || strings.HasPrefix(viString, "-0")) && xTagContainKey(cpOpt, xRef_key_zero_to_8) {
		numBase = 8
	}
	if strings.HasPrefix(viString, "-") {
		viInt64, err := strconv.ParseInt(viString, numBase, 64)
		if err != nil {
			return false, err
		} else {
			return viInt64 > 0, nil
		}
	} else {
		viUint64, err := strconv.ParseUint(viString, numBase, 64)
		if err != nil {
			return false, err
		} else {
			return viUint64 > 0, nil
		}
	}
}

// 格式化时间为字符串
func xFormatTimeToString(t *time.Time, timeLayout string) string {
	//"America/Adak" "Asia/Shanghai"
	var realTimeLayout string
	if len(timeLayout) > 0 {
		realTimeLayout = timeLayout
	} else {
		realTimeLayout = xRef_time_layout
	}
	return t.In(corex.ToTimeLocation()).Format(realTimeLayout)
}

// 解析字符串为时间
func xTransStringToTime(sTime string, timeLayout string) *time.Time {
	var realTimeLayout string
	if len(timeLayout) > 0 {
		realTimeLayout = timeLayout
	} else {
		realTimeLayout = xRef_time_layout
	}
	timeStamp, err := time.ParseInLocation(realTimeLayout, sTime, corex.ToTimeLocation())
	if err != nil {
		return nil
	}
	return &timeStamp
}

func xTransInt64ToTime(srcVal int64, optStr1 string) *time.Time {
	var timeDest *time.Time
	if len(optStr1) <= 0 {
		timeValue := time.UnixMilli(srcVal)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "sec") {
		timeValue := time.UnixMilli(srcVal * 1000)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "min") {
		timeValue := time.UnixMilli(srcVal * 1000 * 60)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "hour") {
		timeValue := time.UnixMilli(srcVal * 1000 * 60 * 60)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "day") {
		timeValue := time.UnixMilli(srcVal * 1000 * 60 * 60 * 24)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "mil") {
		timeValue := time.UnixMilli(srcVal)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "mic") {
		timeValue := time.UnixMicro(srcVal / 1e3)
		timeDest = &timeValue
	} else if strings.Contains(optStr1, "nano") {
		timeValue := time.UnixMilli(srcVal / 1e6)
		timeDest = &timeValue
	} else {
		timeValue := time.UnixMilli(srcVal)
		timeDest = &timeValue
	}
	return timeDest
}
func xTransTimeToInt64(pTime *time.Time, optStr1 string) int64 {
	if len(optStr1) <= 0 {
		return pTime.UnixMilli()
	} else if strings.Contains(optStr1, "sec") {
		return pTime.UnixMilli() / 1000
	} else if strings.Contains(optStr1, "min") {
		return pTime.UnixMilli() / (1000 * 60)
	} else if strings.Contains(optStr1, "hour") {
		return pTime.UnixMilli() / (1000 * 60 * 60)
	} else if strings.Contains(optStr1, "day") {
		return pTime.UnixMilli() / (1000 * 60 * 60 * 24)
	} else if strings.Contains(optStr1, "mil") {
		return pTime.UnixMilli()
	} else if strings.Contains(optStr1, "mic") {
		return pTime.UnixMicro()
	} else if strings.Contains(optStr1, "nano") {
		return pTime.UnixNano()
	} else {
		return pTime.UnixMilli()
	}
}
