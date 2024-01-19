/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:38
 * @version 1.0
 */
package refx

import (
	"testing"
	"time"
)

type Z01t int
type Z01A struct {
	Age string
	//Age2  uint64
	Age3  uint64
	Vtime time.Time
}
type Z01B struct {
	Year uint64 `xref:"Age;t.nano"`
}

func Test_Z01001(t *testing.T) {
	a := int(0274)
	println(a)
	z01A := Z01A{
		Age: "18446744073709551550",
		//Age2:  18446744073709551539,
		Age3:  18446744073709551550,
		Vtime: time.Now(),
	}
	//a := int64(64)
	z01B := Z01B{
		Year: 00,
	}
	XReflectCopy(z01A, &z01B)
	println(z01A.Age3, z01B.Year)

}
