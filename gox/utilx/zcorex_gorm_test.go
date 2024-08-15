/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:38
 * @version 1.0
 */
package utilx

import (
	"fmt"
	"testing"
	"time"
)

type GromModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type DemoUserName struct {
	GromModel
	UserId   int    `xorderby:"4;table=userName,opt=deSc"`
	UserName string `xorderby:"5;table=-,"`
	AGE      int64  `xorderby:"6;table=,"`
}

func TestGormParseOrderBy(t *testing.T) {
	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	orderByMap := xGormParseOrderByTag(&DemoUserName{}, "nihadoDDdd")
	fmt.Println(orderByMap)
	orderByList := []XOrderBy{
		{SortField: 8, SortDesc: false},
		{SortField: 1, SortDesc: true},
		{SortField: 6, SortDesc: false},
	}
	//orderByString := GormParseOrderBy(&DemoUserName{}, "demoUserOrg", orderByList, &XOrderBy{SortField: 3, SortDesc: false})
	orderByString := GormParseOrderBy(&DemoUserName{}, "demoUserOrg", orderByList, &XOrderBy{SortField: 3, SortDesc: false}, XOrderColumn{8, "t.fav_count"})
	//orderByString := GormParseOrderByCreatedAt(&DemoUserName{}, "", orderByList, true)

	fmt.Println(orderByString)
}
