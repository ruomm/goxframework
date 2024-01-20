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

type OrigS2P struct {
	Orig int
}

type DestS2P struct {
	Vint     *int       `xref:"Orig" json:""`
	Vint8    *int8      `xref:"Orig" json:""`
	Vint16   *int16     `xref:"Orig" json:""`
	Vint32   *int32     `xref:"Orig" json:""`
	Vint64   *int64     `xref:"Orig" json:""`
	Vuint    *uint      `xref:"Orig" json:""`
	Vuint8   *uint8     `xref:"Orig" json:""`
	Vuint16  *uint16    `xref:"Orig" json:""`
	Vuint32  *uint32    `xref:"Orig" json:""`
	Vuint64  *uint64    `xref:"Orig" json:""`
	Vfloat32 *float32   `xref:"Orig" json:""`
	Vfloat64 *float64   `xref:"Orig" json:""`
	Vstring  *string    `xref:"Orig" json:""`
	Vbool    *bool      `xref:"Orig" json:""`
	VTime    *time.Time `xref:"Orig;t.day" json:""`
}

func TestS2P(t *testing.T) {
	orgiP := OrigS2P{
		Orig: 456456,
	}
	destP := DestS2P{
		//Vint: &a,
	}
	XReflectCopy(orgiP, &destP)

	jsonStr, _ := corex.JsonToString(destP)
	fmt.Println(jsonStr)

}
