/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:38
 * @version 1.0
 */
package corex

import (
	"fmt"
	"testing"
)

const (
	//TEST_JSON_STRING = "{\"Vint\":123456789,\"Vint8\":21,\"Vint16\":-13035,\"Vint32\":123456789,\"Vint64\":123456789,\"Vuint\":123456789,\"Vuint8\":21,\"Vuint16\":52501,\"Vuint32\":123456789,\"Vuint64\":123456789,\"Vfloat32\":123456790,\"Vfloat64\":123456789,\"Vstring\":\"123456789\",\"Vbool\":true,\"VTime\":\"1970-01-02T18:17:36.789+08:00\"}"
	//TEST_JSON_STRING        = "{\"Vint\":123456789,\"Vint8\":21,\"Vint16\":-13035,\"Vint32\":123456789,\"Vint64\":123456789,\"Vuint\":123456789,\"Vuint8\":21,\"Vuint16\":52501,\"Vuint32\":123456789,\"Vuint64\":123456789,\"Vfloat32\":123456790,\"Vfloat64\":123456789,\"Vstring\":\"123456789\",\"Vbool\":true,\"VTime\":\"1970-01-02T18:17:36.789+08:00\"}"
	PATH_JSON_FIlE = "02.txt"
)

func TestNextAndPreDayMonth(t *testing.T) {
	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-12-31 00:50:11")
	fmt.Println(TimePreDay(time))
	fmt.Println(TimeNextDay(time))
	fmt.Println(TimePreMonth(time))
	fmt.Println(TimeNextMonth(time))
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
	slice2 := []int{1, 2, 3, 3, 5}
	fmt.Println(SliceDuplicates(slice1)) // 输出：false
	fmt.Println(SliceDuplicates(slice2))
	fmt.Println(SliceContains(slice1, 6))
}

func TestContainsDuplicatesByKey(t *testing.T) {
	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	slice1 := []SliceDuplicatesByKeyTest{SliceDuplicatesByKeyTest{Name: "张三"},
		SliceDuplicatesByKeyTest{Name: "李四"}, SliceDuplicatesByKeyTest{Name: "王五"},
		SliceDuplicatesByKeyTest{Name: "赵六"}, SliceDuplicatesByKeyTest{Name: "张三2"}}
	fmt.Println(SliceDuplicatesByKey(slice1, "Age")) // 输出：false
	//fmt.Println(SliceDuplicates(slice2))
	//fmt.Println(SliceContains(slice1, 6))
}
