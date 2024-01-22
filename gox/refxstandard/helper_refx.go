/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/20 11:26
 * @version 1.0
 */
package refxstandard

import (
	"errors"
	"github.com/ruomm/goxframework/gox/corex"
	"strconv"
	"strings"
	"time"
)

// 是否tidy
func XrefTagTidy(tagValue string) bool {
	return xTagContainKey(tagValue, xRef_key_tidy)
}

// 解析真实的tag名称
func xParseRefTagName(optTag string) string {
	if len(optTag) <= 0 {
		//real_tag_name = xRef_tag_cp_xreft
		return "xref"
	} else {
		return optTag
	}

}

// 解析真实的NameSpace
func xParseRefNameSpace(optNameSpace string, origNameSpace string) string {
	if len(optNameSpace) <= 0 {
		return origNameSpace
	} else {
		return optNameSpace
	}
}

// 属性空值设置 for refx.
type XrefOption struct {
	f func(*xrefOptions)
}

type xrefOptions struct {
	optTag       string //关联的tag的名称
	optNameSpace string //关联源的nameSpace
}

// 设置关联tag的名称，不设置默认为xref
func XrefOptTag(tag string) XrefOption {
	return XrefOption{func(do *xrefOptions) {
		do.optTag = tag
	}}
}

// 设置关联源的nameSpace，不设置则取源对象的类型名称
func XrefOptNameSpace(nameSpace string) XrefOption {
	return XrefOption{func(do *xrefOptions) {
		do.optNameSpace = nameSpace
	}}
}

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

// // 字符串转换为bool
func xTransStringIntToBool(viString string, cpOpt string) (bool, error) {
	if viString == "0" {
		return false, nil
	} else if viString == "1" {
		return true, nil
	} else {
		return false, errors.New("string int to bool must be \"0\" or \"1\",\"0\"->false,\"1\"->true")
	}
	//
	//numBase := 10
	//if strings.HasPrefix(viString, "0x") || strings.HasPrefix(viString, "0X") {
	//	numBase = 16
	//	viString = viString[2:]
	//} else if strings.HasPrefix(viString, "-0x") || strings.HasPrefix(viString, "-0X") {
	//	numBase = 16
	//	viString = "-" + viString[3:]
	//} else if (strings.HasPrefix(viString, "0") || strings.HasPrefix(viString, "-0")) && xTagContainKey(cpOpt, xRef_key_zero_to_8) {
	//	numBase = 8
	//}
	//if strings.HasPrefix(viString, "-") {
	//	viInt64, err := strconv.ParseInt(viString, numBase, 64)
	//	if err != nil {
	//		return false, err
	//	} else {
	//		return viInt64 > 0, nil
	//	}
	//} else {
	//	viUint64, err := strconv.ParseUint(viString, numBase, 64)
	//	if err != nil {
	//		return false, err
	//	} else {
	//		return viUint64 > 0, nil
	//	}
	//}
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

func xTransInt64ToTime(origVal int64, optStr string) *time.Time {
	var timeDest *time.Time
	if len(optStr) <= 0 {
		timeValue := time.UnixMilli(origVal)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "sec") {
		timeValue := time.UnixMilli(origVal * 1000)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "min") {
		timeValue := time.UnixMilli(origVal * 1000 * 60)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "hour") {
		timeValue := time.UnixMilli(origVal * 1000 * 60 * 60)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "day") {
		timeValue := time.UnixMilli(origVal * 1000 * 60 * 60 * 24)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "mil") {
		timeValue := time.UnixMilli(origVal)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "mic") {
		timeValue := time.UnixMicro(origVal / 1e3)
		timeDest = &timeValue
	} else if strings.Contains(optStr, "nano") {
		timeValue := time.UnixMilli(origVal / 1e6)
		timeDest = &timeValue
	} else {
		timeValue := time.UnixMilli(origVal)
		timeDest = &timeValue
	}
	return timeDest
}
func xTransTimeToInt64(pTime *time.Time, optStr string) int64 {
	if len(optStr) <= 0 {
		return pTime.UnixMilli()
	} else if strings.Contains(optStr, "sec") {
		return pTime.UnixMilli() / 1000
	} else if strings.Contains(optStr, "min") {
		return pTime.UnixMilli() / (1000 * 60)
	} else if strings.Contains(optStr, "hour") {
		return pTime.UnixMilli() / (1000 * 60 * 60)
	} else if strings.Contains(optStr, "day") {
		return pTime.UnixMilli() / (1000 * 60 * 60 * 24)
	} else if strings.Contains(optStr, "mil") {
		return pTime.UnixMilli()
	} else if strings.Contains(optStr, "mic") {
		return pTime.UnixMicro()
	} else if strings.Contains(optStr, "nano") {
		return pTime.UnixNano()
	} else {
		return pTime.UnixMilli()
	}
}
