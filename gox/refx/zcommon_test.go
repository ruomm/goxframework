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

type OrigV struct {
	Orig int
}

type DestV struct {
	Vint     int       `xref:"Orig" json:""`
	Vint8    int8      `xref:"Orig" json:""`
	Vint16   int16     `xref:"Orig" json:""`
	Vint32   int32     `xref:"Orig" json:""`
	Vint64   int64     `xref:"Orig" json:""`
	Vuint    uint      `xref:"Orig" json:""`
	Vuint8   uint8     `xref:"Orig" json:""`
	Vuint16  uint16    `xref:"Orig" json:""`
	Vuint32  uint32    `xref:"Orig" json:""`
	Vuint64  uint64    `xref:"Orig" json:""`
	Vfloat32 float32   `xref:"Orig" json:""`
	Vfloat64 float64   `xref:"Orig" json:""`
	Vstring  string    `xref:"Orig" json:""`
	Vbool    bool      `xref:"Orig" json:""`
	VTime    time.Time `xref:"Orig" json:""`
}

type DestPV struct {
	Vint     *int      `xref:"Orig" json:""`
	Vint8    *int8     `xref:"Orig" json:""`
	Vint16   *int16    `xref:"Orig" json:""`
	Vint32   *int32    `xref:"Orig" json:""`
	Vint64   *int64    `xref:"Orig" json:""`
	Vuint    *uint     `xref:"Orig" json:""`
	Vuint8   *uint8    `xref:"Orig" json:""`
	Vuint16  *uint16   `xref:"Orig" json:""`
	Vuint32  *uint32   `xref:"Orig" json:""`
	Vuint64  *uint64   `xref:"Orig" json:""`
	Vfloat32 *float32  `xref:"Orig" json:""`
	Vfloat64 *float64  `xref:"Orig" json:""`
	Vstring  *string   `xref:"Orig" json:""`
	Vbool    *bool     `xref:"Orig" json:""`
	VTime    time.Time `xref:"Orig" json:""`
}

func TestIntCase(t *testing.T) {
	orgiV := OrigV{
		Orig: 456456,
	}
	destV := DestV{}
	XReflectCopy(orgiV, &destV)

	jsonStr, _ := corex.JsonToString(destV)
	fmt.Println(jsonStr)

}
