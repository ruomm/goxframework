/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:38
 * @version 1.0
 */
package utils

import (
	"fmt"
	"github.com/ruomm/goxframework/gox/corex"
	"testing"
)

const (
	//TEST_JSON_STRING = "{\"Vint\":123456789,\"Vint8\":21,\"Vint16\":-13035,\"Vint32\":123456789,\"Vint64\":123456789,\"Vuint\":123456789,\"Vuint8\":21,\"Vuint16\":52501,\"Vuint32\":123456789,\"Vuint64\":123456789,\"Vfloat32\":123456790,\"Vfloat64\":123456789,\"Vstring\":\"123456789\",\"Vbool\":true,\"VTime\":\"1970-01-02T18:17:36.789+08:00\"}"
	//TEST_JSON_STRING        = "{\"Vint\":123456789,\"Vint8\":21,\"Vint16\":-13035,\"Vint32\":123456789,\"Vint64\":123456789,\"Vuint\":123456789,\"Vuint8\":21,\"Vuint16\":52501,\"Vuint32\":123456789,\"Vuint64\":123456789,\"Vfloat32\":123456790,\"Vfloat64\":123456789,\"Vstring\":\"123456789\",\"Vbool\":true,\"VTime\":\"1970-01-02T18:17:36.789+08:00\"}"
	PATH_JSON_FIlE = "02.txt"
)

func TestParseRemain(t *testing.T) {
	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	tStart, _ := corex.TimeParseByString(corex.TIME_PATTERN_STANDARD, "2022-12-10 10:50:11")
	tEnd, _ := corex.TimeParseByString(corex.TIME_PATTERN_STANDARD, "2023-12-11 00:50:00")
	months, seconds, err := ParseRemainMonthSenconds(tStart, tEnd)
	if err != nil {
		fmt.Print("错误：")
		fmt.Println(err)
	}
	days := float64(seconds) / float64(3600*24)
	fmt.Print("月份：")
	fmt.Println(months)
	fmt.Print("天数：")
	fmt.Println(days)
	fmt.Print("秒数：")
	fmt.Println(seconds)
	formatShow, err := ParseRemainFormat(tStart, tEnd, true)
	if err != nil {
		fmt.Print("错误：")
		fmt.Println(err)
	}
	fmt.Print("格式化：")
	fmt.Println(formatShow)
}
