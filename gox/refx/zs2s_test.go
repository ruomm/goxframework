/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:38
 * @version 1.0
 */
package refx

import (
	"fmt"
	"github.com/ruomm/goxframework/gox/corex"
	"testing"
	"time"
)

type OrigS2S struct {
	//Orig int
	//Orig string
	//Orig float64
	//Orig bool
	Orig time.Time
}

type DestS2S struct {
	Vint     int       `xref:"Orig;tidy" json:""`
	Vint8    int8      `xref:"Orig;tidy" json:""`
	Vint16   int16     `xref:"Orig;tidy" json:""`
	Vint32   int32     `xref:"Orig;tidy" json:""`
	Vint64   int64     `xref:"Orig;tidy" json:""`
	Vuint    uint      `xref:"Orig;tidy" json:""`
	Vuint8   uint8     `xref:"Orig;tidy" json:""`
	Vuint16  uint16    `xref:"Orig;tidy" json:""`
	Vuint32  uint32    `xref:"Orig;tidy" json:""`
	Vuint64  uint64    `xref:"Orig;tidy" json:""`
	Vfloat32 float32   `xref:"Orig;tidy" json:""`
	Vfloat64 float64   `xref:"Orig;tidy" json:""`
	Vstring  string    `xref:"Orig;tidy" json:""`
	Vbool    bool      `xref:"Orig;tidy" json:""`
	VTime    time.Time `xref:"Orig;tidy" json:""`
}

func TestS2S(t *testing.T) {
	//a := 123456.567
	//a := true
	//a := time.Now()
	a := time.Time{}
	fmt.Println(a)
	orgi := OrigS2S{
		Orig: a,
	}
	dest := DestS2S{}
	println(dest.VTime.UnixMilli())
	XRefCopy(orgi, &dest)

	jsonStr, _ := corex.JsonToString(dest)
	fmt.Println(jsonStr)
	if jsonStr != readJson(3) {
		t.Error("测试失败")
	}

}
