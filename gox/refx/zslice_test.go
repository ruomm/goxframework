/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:38
 * @version 1.0
 */
package refx

import (
	"fmt"
	"math/rand"
	rand2 "math/rand/v2"
	"testing"
	"time"
)

type SType int
type SRType int8
type SBOrig struct {
	//Orig *int
	UserType SType
	//Orig *float64
	//Orig *bool
	//Orig *time.Time
}
type SBOrigExt struct {
	SBOrig
	//Orig *int
	Age SRType
	//Orig *float64
	//Orig *bool
	//Orig *time.Time
}

type SBDest struct {
	RoleType SRType `xref:"UserType;tidy"`
}

type SBDestExt struct {
	SBDest
	Role2 SType  `xref:"UserType;tidy"`
	Name  string `xref:"Age;tidy"`
}

func GenerateOrigStuct() SBOrigExt {
	sbOrig := SBOrigExt{}
	sbOrig.UserType = SType(rand.Intn(600000000))
	sbOrig.Age = SRType(rand2.UintN(10000))
	return sbOrig
}
func GenerateOrigSlice() []SBOrigExt {
	var sbSlice []SBOrigExt
	for i := 0; i < 10; i++ {
		sbOrig := SBOrigExt{}
		sbOrig.UserType = SType(rand.Intn(600000000))
		sbOrig.Age = SRType(rand2.UintN(10000))
		sbSlice = append(sbSlice, sbOrig)
	}
	return sbSlice
}

func Test0001(t *testing.T) {
	sbOrig := GenerateOrigStuct()
	sbDest := SBDestExt{}
	XRefStructCopy(sbOrig, &sbDest)
	fmt.Println(sbDest)
}

func Test0002(t *testing.T) {
	sbOrigSlice := GenerateOrigSlice()
	var sbDestSlice []SBDestExt
	XSliceCopy(sbOrigSlice, &sbDestSlice)
	fmt.Println(sbDestSlice)
}

func Test0003(t *testing.T) {
	sbOrigSlice := GenerateOrigSlice()
	var sbDestSlice []time.Time
	XSliceCopyByKey(sbOrigSlice, &sbDestSlice, "UserType")
	fmt.Println(sbDestSlice)
}