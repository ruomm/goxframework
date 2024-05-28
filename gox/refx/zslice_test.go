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
	randv2 "math/rand/v2"
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
type SCOrigExt struct {
	SBOrig
	//Orig *int
	Age2 SRType
	//Orig *float64
	//Orig *bool
	//Orig *time.Time
}

type SBDest struct {
	RoleType SRType `xref:"SCOrigExt:UserType;tidy"`
}

type SBDestExt struct {
	SBDest
	Role2 SType  `xref:"SBOrigExt:UserType"`
	Age2  SType  `xref:"SCOrigExt:Age2"`
	Name  string `xref:"SCOrigExt:Age2"`
}

func (d SBDestExt) TransMethodExaple(orig SType) int {
	return 546
}

func GenerateOrigStuct() SBOrigExt {
	sbOrig := SBOrigExt{}
	sbOrig.UserType = SType(rand.Intn(600000000))
	sbOrig.Age = SRType(randv2.UintN(10000))
	return sbOrig
}
func GenerateOrigSlice() []SBOrigExt {
	var sbSlice []SBOrigExt
	for i := 0; i < 5; i++ {
		sbOrig := SBOrigExt{}
		sbOrig.UserType = SType(rand.Intn(600000000))
		sbOrig.Age = SRType(randv2.UintN(10000))
		sbSlice = append(sbSlice, sbOrig)
	}
	return sbSlice
}
func GenerateOrigCSlice() []SCOrigExt {
	var scSlice []SCOrigExt
	for i := 0; i < 10; i++ {
		scOrig := SCOrigExt{}
		//scOrig.UserType = SType(rand.Intn(600000000))
		scOrig.Age2 = SRType(randv2.UintN(10000))
		scSlice = append(scSlice, scOrig)
	}
	return scSlice
}

func GenerateOrigSliceP() []*SBOrigExt {
	var sbSlice []*SBOrigExt
	for i := 0; i < 10; i++ {
		sbOrig := SBOrigExt{}
		sbOrig.UserType = SType(rand.Intn(600000000))
		sbOrig.Age = SRType(randv2.UintN(10000))
		sbSlice = append(sbSlice, &sbOrig)
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
	scOrigSlice := GenerateOrigCSlice()
	var sbDestSlice []SBDestExt
	XSliceCopy(sbOrigSlice, &sbDestSlice)
	XSliceQcopy(sbOrigSlice, &sbDestSlice)
	err := XSliceQcopy(scOrigSlice, &sbDestSlice)
	print(err)
	fmt.Println(sbDestSlice)
}

func Test0003(t *testing.T) {
	sbOrigSlice := GenerateOrigSlice()
	var sbDestSlice []time.Time
	XSliceCopyByKey(sbOrigSlice, &sbDestSlice, "UserType", XrefOptCopyOption("tidy,tomt:TransMethodExaple"))
	fmt.Println(sbDestSlice)
}
func Test0005(t *testing.T) {
	sbOrigSlice := GenerateOrigSliceP()
	var sbDestMap map[string]*SBDestExt
	XSliceCopyToMap(sbOrigSlice, &sbDestMap, "UserType", "", XrefOptMapKeyAppend("Key-"))
	fmt.Println(sbDestMap)
}

func Test0006(t *testing.T) {
	//str := ""
	str := "123,456789"
	var sbDestSlice []int
	//XSliceCopy(sbOrigSlice, &sbDestSlice, XrefOptNameSpace("dBOrigfasd"))
	err := XStringToSlice(str, "", true, &sbDestSlice)
	fmt.Println(err)
	fmt.Println(sbDestSlice)
}
