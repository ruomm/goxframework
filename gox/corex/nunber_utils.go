package corex

import "strconv"

// 填充数字为多少位数

func FillStringLeft(input string, padding string, len_total int) string {
	return fillStringByMode(input, padding, len_total, true)
}

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

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StrToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func Uint64ToStr(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func StrToUInt64(str string) uint64 {
	i, _ := strconv.ParseUint(str, 10, 64)
	return i
}
