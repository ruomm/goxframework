/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:38
 * @version 1.0
 */
package refx

import (
	"fmt"
	"testing"
	"time"
)

type Z02t int
type Z02A struct {
	Age string
	//Age2  uint64
	Age3  float32
	Vtime time.Time
}
type Z02B struct {
	Year float32 `xref:"Age3;t.nano"`
}

func Test_Z00201(t *testing.T) {
	z02A := Z02A{
		Age: "18446744073709551550",
		//Age2:  18446744073709551539,
		Age3:  12185.123456,
		Vtime: time.Now(),
	}
	//a := int64(64)
	z02B := Z02B{
		Year: 00,
	}
	XReflectCopy(z02A, &z02B)
	fmt.Printf("%.10f\n", z02A.Age3)
	fmt.Printf("%.10f\n", z02B.Year)

}
