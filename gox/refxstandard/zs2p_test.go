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
	"testing"
	"time"
)

type OrigS2P struct {
	//Orig int
	//Orig string
	//Orig float64
	//Orig bool
	Orig time.Time
}

type DestS2P struct {
	Vint     *int       `xref:"Orig;tidy" json:""`
	Vint8    *int8      `xref:"Orig;tidy,tomt:TransMethodExaple" json:""`
	Vint16   *int16     `xref:"Orig;tidy" json:""`
	Vint32   *int32     `xref:"Orig;tidy" json:""`
	Vint64   *int64     `xref:"Orig;tidy" json:""`
	Vuint    *uint      `xref:"Orig;tidy" json:""`
	Vuint8   *uint8     `xref:"Orig;tidy" json:""`
	Vuint16  *uint16    `xref:"Orig;tidy" json:""`
	Vuint32  *uint32    `xref:"Orig;tidy" json:""`
	Vuint64  *uint64    `xref:"Orig;tidy" json:""`
	Vfloat32 *float32   `xref:"Orig;tidy" json:""`
	Vfloat64 *float64   `xref:"Orig;tidy" json:""`
	Vstring  *string    `xref:"Orig;tidy" json:""`
	Vbool    *bool      `xref:"Orig;tidy" json:""`
	VTime    *time.Time `xref:"Orig;tidy" json:""`
}

func (d DestS2P) TransMethodExaple(orig time.Time) int {
	return 456
}

func TestS2P(t *testing.T) {
	//a := 123456.567
	//a := true
	a := time.Now()
	//a := time.Time{}
	fmt.Println(a)
	orig := OrigS2P{
		Orig: a,
	}
	dest := DestS2P{
		//Vint: &a,
	}
	XRefStructCopy(orig, &dest)

	jsonStr, _ := corex.JsonFormatByString(dest)
	fmt.Println(jsonStr)
	if jsonStr != readJson(3) {
		t.Error("测试失败")
	}

}
