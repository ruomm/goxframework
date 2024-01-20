/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:38
 * @version 1.0
 */
package refxstandard

import (
	"fmt"
	"github.com/ruomm/goxframework/gox/corex"
	"path"
)

const (
	//TEST_JSON_STRING = "{\"Vint\":123456789,\"Vint8\":21,\"Vint16\":-13035,\"Vint32\":123456789,\"Vint64\":123456789,\"Vuint\":123456789,\"Vuint8\":21,\"Vuint16\":52501,\"Vuint32\":123456789,\"Vuint64\":123456789,\"Vfloat32\":123456790,\"Vfloat64\":123456789,\"Vstring\":\"123456789\",\"Vbool\":true,\"VTime\":\"1970-01-02T18:17:36.789+08:00\"}"
	//TEST_JSON_STRING        = "{\"Vint\":123456789,\"Vint8\":21,\"Vint16\":-13035,\"Vint32\":123456789,\"Vint64\":123456789,\"Vuint\":123456789,\"Vuint8\":21,\"Vuint16\":52501,\"Vuint32\":123456789,\"Vuint64\":123456789,\"Vfloat32\":123456790,\"Vfloat64\":123456789,\"Vstring\":\"123456789\",\"Vbool\":true,\"VTime\":\"1970-01-02T18:17:36.789+08:00\"}"
	PATH_JSON_FIlE = "02.txt"
)

func readJson(index int) string {
	dir := "D:\\temp\\json"
	fileName := fmt.Sprintf("%02d", index) + ".txt"
	jsonFilePath := path.Join(dir, fileName)
	jsonStr, _ := corex.ReadFile(jsonFilePath)
	return jsonStr
}
