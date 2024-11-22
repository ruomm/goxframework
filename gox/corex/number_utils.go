package corex

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// 左侧填充字符串到特定长度，input:字符串,paddin:填充的字符,len_total:总计长度
// 如：FillStringLeft("abc", "x", 6)="xxxabc",FillStringLeft("123456789", "x", 6)="123456789"
func FillStringLeft(input string, padding string, len_total int) string {
	return fillStringByMode(input, padding, len_total, true)
}

// 右侧侧填充字符串到特定长度，input:字符串,paddin:填充的字符,len_total:总计长度
// 如：FillStringRight("abc", "x", 6)="abcxxx",FillStringLeft("123456789", "x", 6)="123456789"
func FillStringRight(input string, padding string, len_total int) string {
	return fillStringByMode(input, padding, len_total, false)
}

func fillStringByMode(input string, padding string, len_total int, leftFillMode bool) string {
	fillLen := len_total - len(input)
	if fillLen <= 0 {
		return input
	}
	if len(padding) <= 0 {
		padding = " "
	}
	for i := 0; i < fillLen; i++ {
		if leftFillMode {
			input = padding + input
		} else {
			input += padding
		}

	}
	return input
}

// int64转string
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

// int64转string，左侧填充0到特定长度
// 如：Int64ToStrFill(15,4) = 0015,Int64ToStrFill(-15,4) = -0015
func Int64ToStrFill(i int64, len_total int) string {
	if i >= 0 {
		return FillStringLeft(strconv.FormatInt(i, 10), "0", len_total)
	} else {
		return "-" + FillStringLeft(strconv.FormatInt(-i, 10), "0", len_total)
	}
}

// string转int64，无法转换时候返回0
func StrToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

// uint64转string
func Uint64ToStr(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// uint64转string，左侧填充0到特定长度
func Uint64ToStrFill(i uint64, len_total int) string {
	return FillStringLeft(strconv.FormatUint(i, 10), "0", len_total)
}

// string转uint64，无法转换时候返回0
func StrToUInt64(str string) uint64 {
	i, _ := strconv.ParseUint(str, 10, 64)
	return i
}

// float64转string
func Float64ToStr(floatVal float64, prec int) string {
	realPrec := -1
	if prec >= 0 {
		realPrec = prec
	}
	return strconv.FormatFloat(floatVal, 'f', realPrec, 64)
}

// float64转换为特定小数位数的float
// 如：Float64Format(5.6781, 1) = 5.68
func Float64Format(floatVal float64, prec int) float64 {
	if prec < 0 {
		return floatVal
	}
	resultVal, _ := strconv.ParseFloat(Float64ToStr(floatVal, prec), 64)
	return resultVal
}

// string转float64，无法转换时候返回0
func StrToFloat64(str string) float64 {
	numFloat64, _ := strconv.ParseFloat(str, 64)
	return numFloat64
}

// byteSize存储大小格式化为KB、MB、GB、TB存储大小显示，prec为保留的小数位数
// 如：512*1024*1024 = 512MB，1000*1024*1024 = 0.9766GB
func StoreSizeFormat(byteSize int64, prec int) string {
	realPrec := -1
	if prec >= 0 {
		realPrec = prec
	}
	if byteSize >= 0 {
		if byteSize < int64(1000) {
			return strconv.FormatInt(byteSize, 10)
		} else if byteSize < int64(1024*1000) {
			mbSize := float64(byteSize) / 1024
			return strconv.FormatFloat(mbSize, 'f', realPrec, 64) + "KB"
		} else if byteSize < int64(1024*1024*1000) {
			mbSize := float64(byteSize) / (1024 * 1024)
			return strconv.FormatFloat(mbSize, 'f', realPrec, 64) + "MB"
		} else if byteSize < int64(1024*1024*1024*1000) {
			mbSize := float64(byteSize) / (1024 * 1024 * 1024)
			return strconv.FormatFloat(mbSize, 'f', realPrec, 64) + "GB"
		} else {
			mbSize := float64(byteSize) / (1024 * 1024 * 1024 * 1024)
			return strconv.FormatFloat(mbSize, 'f', realPrec, 64) + "TB"
		}
	} else {
		if byteSize > int64(-1000) {
			return strconv.FormatInt(byteSize, 10)
		} else if byteSize > int64(-1024*1000) {
			mbSize := float64(byteSize) / 1024
			return strconv.FormatFloat(mbSize, 'f', realPrec, 64) + "KB"
		} else if byteSize > int64(-1024*1024*1000) {
			mbSize := float64(byteSize) / (1024 * 1024)
			return strconv.FormatFloat(mbSize, 'f', realPrec, 64) + "MB"
		} else if byteSize > int64(-1024*1024*1024*1000) {
			mbSize := float64(byteSize) / (1024 * 1024 * 1024)
			return strconv.FormatFloat(mbSize, 'f', realPrec, 64) + "GB"
		} else {
			mbSize := float64(byteSize) / (1024 * 1024 * 1024 * 1024)
			return strconv.FormatFloat(mbSize, 'f', realPrec, 64) + "TB"
		}
	}
}

// KB、MB、GB、TB存储大小解析为byteSize存储大小
// 如： 512MB、512M = 512*1024*1024，0.9766GB、0.9766G = 1000*1024*1024
func StoreSizeParse(bytesizeStr string) (int64, error) {
	if len(bytesizeStr) == 0 {
		return 0, errors.New("parse bytesize err,input string is empty")
	}
	sizeString := strings.ToLower(bytesizeStr)
	sizeStrLen := len(sizeString)
	numRate := int64(1)
	numStr := ""
	if strings.HasSuffix(sizeString, "kb") {
		numStr = strings.TrimSpace(sizeString[0 : sizeStrLen-2])
		numRate = int64(1024)
	} else if strings.HasSuffix(sizeString, "mb") {
		numStr = strings.TrimSpace(sizeString[0 : sizeStrLen-2])
		numRate = int64(1024 * 1024)
	} else if strings.HasSuffix(sizeString, "gb") {
		numStr = strings.TrimSpace(sizeString[0 : sizeStrLen-2])
		numRate = int64(1024 * 1024 * 1024)
	} else if strings.HasSuffix(sizeString, "tb") {
		numStr = strings.TrimSpace(sizeString[0 : sizeStrLen-2])
		numRate = int64(1024 * 1024 * 1024 * 1024)
	} else if strings.HasSuffix(sizeString, "k") {
		numStr = strings.TrimSpace(sizeString[0 : sizeStrLen-1])
		numRate = int64(1024)
	} else if strings.HasSuffix(sizeString, "m") {
		numStr = strings.TrimSpace(sizeString[0 : sizeStrLen-1])
		numRate = int64(1024 * 1024)
	} else if strings.HasSuffix(sizeString, "g") {
		numStr = strings.TrimSpace(sizeString[0 : sizeStrLen-1])
		numRate = int64(1024 * 1024 * 1024)
	} else if strings.HasSuffix(sizeString, "t") {
		numStr = strings.TrimSpace(sizeString[0 : sizeStrLen-1])
		numRate = int64(1024 * 1024 * 1024 * 1024)
	} else {
		numStr = strings.TrimSpace(sizeString)
		numRate = int64(1)
	}
	if len(numStr) == 0 {
		return 0, errors.New("parse bytesize err,input string is not bytesize format")
	}
	numFloat64, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, errors.New("parse bytesize err,input string is not bytesize format")
	}
	value := int64(math.Round(numFloat64 * float64(numRate)))
	return value, nil
}

// 格式化时间戳为ms、s、m、h、d(毫秒、秒、分钟、小时、天)显示的时间，timeNumberValue:时间戳,prce:保留小数位数,secondMode:false-毫秒模式,true-秒模式
// TimeNumberFormat(150,2,true) = 2.50m,TimeNumberFormat(150,2,false) = 150ms
func TimeNumberFormat(timeNumberValue int64, prec int, secondMode bool) string {
	realPrec := -1
	if prec >= 0 {
		realPrec = prec
	}
	timeMillSeconds := timeNumberValue
	if secondMode {
		timeMillSeconds = timeNumberValue * 1000
	}
	if timeMillSeconds >= 0 {
		if timeMillSeconds == 0 {
			return "0s"
		} else if timeMillSeconds < int64(1000) {
			// 小于1秒，毫秒显示
			return strconv.FormatInt(timeMillSeconds, 10) + "ms"
		} else if timeMillSeconds < int64(1000*120) {
			// 小于2分钟，秒显示
			floatVal := float64(timeMillSeconds) / 1000
			return strconv.FormatFloat(floatVal, 'f', realPrec, 64) + "s"
		} else if timeMillSeconds < int64(1000*60*60*2) {
			// 小于2小时，分显示
			floatVal := float64(timeMillSeconds) / (1000 * 60)
			return strconv.FormatFloat(floatVal, 'f', realPrec, 64) + "m"
		} else if timeMillSeconds < int64(1000*60*60*24*5) {
			// 小于5天，时显示
			floatVal := float64(timeMillSeconds) / (1000 * 60 * 60)
			return strconv.FormatFloat(floatVal, 'f', realPrec, 64) + "h"
		} else {
			// 其他，天显示
			floatVal := float64(timeMillSeconds) / (1000 * 60 * 60 * 24)
			return strconv.FormatFloat(floatVal, 'f', realPrec, 64) + "d"
		}
	} else {
		if timeMillSeconds > int64(-1000) {
			// 小于1秒，毫秒显示
			return strconv.FormatInt(timeMillSeconds, 10) + "ms"
		} else if timeMillSeconds > int64(-1000*120) {
			// 小于2分钟，秒显示
			floatVal := float64(timeMillSeconds) / 1000
			return strconv.FormatFloat(floatVal, 'f', realPrec, 64) + "s"
		} else if timeMillSeconds > int64(-1000*60*60*2) {
			// 小于2小时，分显示
			floatVal := float64(timeMillSeconds) / (1000 * 60)
			return strconv.FormatFloat(floatVal, 'f', realPrec, 64) + "m"
		} else if timeMillSeconds > int64(-1000*60*60*24*5) {
			// 小于5天，时显示
			floatVal := float64(timeMillSeconds) / (1000 * 60 * 60)
			return strconv.FormatFloat(floatVal, 'f', realPrec, 64) + "h"
		} else {
			// 其他，天显示
			floatVal := float64(timeMillSeconds) / (1000 * 60 * 60 * 24)
			return strconv.FormatFloat(floatVal, 'f', realPrec, 64) + "d"
		}
	}
}

// ms、s、m、h、d(毫秒、秒、分钟、小时、天)显示的时间解析为时间戳，timenumStr:格式化的时间戳,secondMode:false-毫秒模式,true-秒模式
// TimeNumberParse(150s,true) = 150,TimeNumberFormat(150s,false) = 150*1000
// TimeNumberParse(150,true) = 150,TimeNumberFormat(150,false) = 150
func TimeNumberParse(timenumStr string, secondMode bool) (int64, error) {
	if len(timenumStr) == 0 {
		if secondMode {
			return 0, errors.New("parse time seconds err,input string is empty")
		} else {
			return 0, errors.New("parse time milliseconds err,input string is empty")
		}
	}
	if secondMode {
		timeString := strings.ToLower(timenumStr)
		sizeStrLen := len(timeString)
		numRate := int64(1)
		numStr := ""
		isMs := false
		if strings.HasSuffix(timeString, "ms") {
			numRate = 1
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-2])
			isMs = true
		} else if strings.HasSuffix(timeString, "s") {
			numRate = 1
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
		} else if strings.HasSuffix(timeString, "m") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
			numRate = 60
		} else if strings.HasSuffix(timeString, "h") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
			numRate = 3600
		} else if strings.HasSuffix(timeString, "d") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
			numRate = 3600 * 24
		} else if strings.HasSuffix(timeString, "w") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
			numRate = 3600 * 24 * 7
		} else if strings.HasSuffix(timeString, "mon") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-3])
			numRate = 3600 * 24 * 30
		} else if strings.HasSuffix(timeString, "y") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
			numRate = 3600 * 24 * 365
		} else {
			numRate = 1
			numStr = strings.TrimSpace(timeString)
		}
		if len(numStr) == 0 {
			return 0, errors.New("parse time seconds err,input string is not time number format")
		}
		numFloat64, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, errors.New("parse time seconds err,input string is not time number format")
		}
		if isMs {
			value := int64(math.Round(numFloat64 * float64(numRate) / 1000))
			return value, nil
		} else {
			value := int64(math.Round(numFloat64 * float64(numRate)))
			return value, nil
		}
	} else {
		timeString := strings.ToLower(timenumStr)
		sizeStrLen := len(timeString)
		numRate := int64(1)
		numStr := ""
		if strings.HasSuffix(timeString, "ms") {
			numRate = 1
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-2])
		} else if strings.HasSuffix(timeString, "s") {
			numRate = 1000
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
		} else if strings.HasSuffix(timeString, "m") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
			numRate = 1000 * 60
		} else if strings.HasSuffix(timeString, "h") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
			numRate = 1000 * 3600
		} else if strings.HasSuffix(timeString, "d") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
			numRate = 1000 * 3600 * 24
		} else if strings.HasSuffix(timeString, "w") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
			numRate = 1000 * 3600 * 24 * 7
		} else if strings.HasSuffix(timeString, "mon") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-3])
			numRate = 1000 * 3600 * 24 * 30
		} else if strings.HasSuffix(timeString, "y") {
			numStr = strings.TrimSpace(timeString[0 : sizeStrLen-1])
			numRate = 1000 * 3600 * 24 * 365
		} else {
			numRate = 1
			numStr = strings.TrimSpace(timeString)
		}
		if len(numStr) == 0 {
			return 0, errors.New("parse time milliseconds err,input string is not time number format")

		}
		numFloat64, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, errors.New("parse time milliseconds err,input string is not time number format")
		}
		value := int64(math.Round(numFloat64 * float64(numRate)))
		return value, nil
	}
}
