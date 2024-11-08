/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:38
 * @version 1.0
 */
package corex

import (
	"fmt"
	url2 "net/url"
	"testing"
)

const (
	//TEST_JSON_STRING = "{\"Vint\":123456789,\"Vint8\":21,\"Vint16\":-13035,\"Vint32\":123456789,\"Vint64\":123456789,\"Vuint\":123456789,\"Vuint8\":21,\"Vuint16\":52501,\"Vuint32\":123456789,\"Vuint64\":123456789,\"Vfloat32\":123456790,\"Vfloat64\":123456789,\"Vstring\":\"123456789\",\"Vbool\":true,\"VTime\":\"1970-01-02T18:17:36.789+08:00\"}"
	//TEST_JSON_STRING        = "{\"Vint\":123456789,\"Vint8\":21,\"Vint16\":-13035,\"Vint32\":123456789,\"Vint64\":123456789,\"Vuint\":123456789,\"Vuint8\":21,\"Vuint16\":52501,\"Vuint32\":123456789,\"Vuint64\":123456789,\"Vfloat32\":123456790,\"Vfloat64\":123456789,\"Vstring\":\"123456789\",\"Vbool\":true,\"VTime\":\"1970-01-02T18:17:36.789+08:00\"}"
	PATH_JSON_FIlE = "02.txt"
)

func TestNextAndPreDayMonth(t *testing.T) {
	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-12-31 00:50:11")
	//fmt.Println(TimePreDay(time))
	//fmt.Println(TimeNextDay(time))
	//fmt.Println(TimePreMonth(time))
	//fmt.Println(TimeNextMonth(time))
	urlStr, _ := url2.JoinPath("http://10.0.100.18/", "/api/v1", "update/1")
	fmt.Println(urlStr)
}

func TestFileNameAndExtension(t *testing.T) {
	file_path := "/a.txt"
	fmt.Println(GetFileName(file_path))
	fmt.Println(GetFileNameWithoutExtension(file_path))
	fmt.Println(GetFileExtension(file_path))
}

type SliceDuplicatesByKeyTest struct {
	Name string
	Age  int
}

func TestContainsDuplicates(t *testing.T) {
	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{1, 2, 3, 3, 5, 6}
	fmt.Println(SliceDuplicates(slice1)) // 输出：false
	fmt.Println(SliceDuplicates(slice2))
	fmt.Println(SliceContains(slice1, 6, 5))
	fmt.Println(SliceOnlyContains(slice2, 1, 2, 3, 4, 5, 6))
}

func TestContainsDuplicatesByKey(t *testing.T) {
	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	slice1 := []SliceDuplicatesByKeyTest{SliceDuplicatesByKeyTest{Name: "张三", Age: 10},
		SliceDuplicatesByKeyTest{Name: "李四", Age: 11}, SliceDuplicatesByKeyTest{Name: "王五", Age: 12},
		SliceDuplicatesByKeyTest{Name: "赵六", Age: 13}, SliceDuplicatesByKeyTest{Name: "张三", Age: 10}}
	fmt.Println(SliceDuplicatesByKey(slice1, ""))
	fmt.Println(SliceDuplicatesByKey(slice1, "Name", "Age"))
	fmt.Println(SliceContainsByKey(slice1, "Name", "张三2", "王五"))
	fmt.Println(SliceOnlyContainsByKey(slice1, "Name", "张三", "王五", "李四", "王五", "赵六"))
	// 输出：false
	//fmt.Println(SliceDuplicates(slice2))
	//fmt.Println(SliceContains(slice1, 6))
}

func TestTimeValidDayString(t *testing.T) {
	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	validResult := TimeNextMonthByString("2024-03-05")
	fmt.Println(validResult)
	// 输出：false
	//fmt.Println(SliceDuplicates(slice2))
	//fmt.Println(SliceContains(slice1, 6))
}

func TestTimeOffsetFunctions(t *testing.T) {
	time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-31 00:50:11")
	timeResultString := TimeOffsetDay(time, -5)
	fmt.Println(timeResultString)
	timeResultString = TimeOffsetDayByString("2024-03-05", -5)
	fmt.Println(timeResultString)
	timeResultString = TimeOffsetMonth(time, -4)
	fmt.Println(timeResultString)
	timeResultString = TimeOffsetMonthByString("2024-03-12", -4)
	fmt.Println(timeResultString)
}
func TestIsRegexMatch(t *testing.T) {
	match := IsRegexMatch("^[a-zA-Z0-9-_\\.]{1,255}$", "a_-.123")
	fmt.Println(match)
	match = IsRegexMatch("^[a-z]$", "a_-.123")
	fmt.Println(match)
}

func TestIsString(t *testing.T) {

	fmt.Println(ToCamelCase("d_d@DSa"))
	fmt.Println(ToSnakeCase("ADSSdas"))
	fmt.Println(FirstLetterToLower("ADSSdas"))
	fmt.Println(FirstLetterToUpper("adssDAS"))
	//match = IsRegexMatch("^[a-z]$", "a_-.123")
	//fmt.Println(match)
}

func TestIsChinese(t *testing.T) {
	str := "、我是河南人住在黄河边上"
	//index := strings.LastIndex(str, "黄")
	//fmt.Println(index)
	//strResult := str[index+len("黄"):]
	//fmt.Println(strResult)
	result := IsAllChineseWithPunctuation(str)
	fmt.Println(result)
}
