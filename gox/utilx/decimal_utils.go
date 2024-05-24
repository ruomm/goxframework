/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/4/10 13:33
 * @version 1.0
 */
package utilx

import (
	"github.com/shopspring/decimal"
	"math"
)

type DecimalHelper struct {
	Prec int32 // 保留的小数位数
}

// float64类型转换为Decimal类型
func (u *DecimalHelper) ToDecimal(floatVal float64) decimal.Decimal {
	return decimal.NewFromFloat(floatVal).Round(u.Prec)
}

// int64类型转换为Decimal类型
func (u *DecimalHelper) ToDecimalByInt(intVal int64) decimal.Decimal {
	return decimal.NewFromFloat(float64(intVal)).Round(u.Prec)
}

// Decimal类型转换为float64类型
func (u *DecimalHelper) ToFloat(d1 decimal.Decimal) float64 {
	//return d1.InexactFloat64()
	floatVal, _ := d1.Float64()
	return floatVal
}

// Decimal类型转换为int64类型
func (u *DecimalHelper) ToInt(d1 decimal.Decimal) int64 {
	//return d1.InexactFloat64()
	floatVal, _ := d1.Float64()
	intVal := int64(math.Round(floatVal))
	return intVal
}

// 使用Decimal格式化float64值
func (u *DecimalHelper) FormatFloat(v float64) float64 {
	//return d1.InexactFloat64()
	d1 := u.ToDecimal(v)
	floatVal, _ := d1.Float64()
	return floatVal
}

// 使用Decimal格式化int64值
func (u *DecimalHelper) FormatInt(v int64) int64 {
	//return d1.InexactFloat64()
	d1 := u.ToDecimalByInt(v)
	floatVal, _ := d1.Float64()
	intVal := int64(math.Round(floatVal))
	return intVal
}

// 比较大小
func (u *DecimalHelper) Compare(d1 decimal.Decimal, d2 decimal.Decimal) int {
	return d1.Cmp(d2)
}

// 和0比较大小
func (u *DecimalHelper) CompareByZero(d1 decimal.Decimal) int {
	dZero := u.ToDecimal(0)
	return d1.Cmp(dZero)
}

// 相加
func (u *DecimalHelper) Add(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Add(d2).Round(u.Prec)
}

// 相减
func (u *DecimalHelper) Sub(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Sub(d2).Round(u.Prec)
}

// 乘以
func (u *DecimalHelper) Mul(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Mul(d2).Round(u.Prec)
}

// 除以
func (u *DecimalHelper) Div(d1 decimal.Decimal, d2 decimal.Decimal) decimal.Decimal {
	return d1.Div(d2).Round(u.Prec)
}

// 相反
func (u *DecimalHelper) Reverse(d1 decimal.Decimal) decimal.Decimal {
	dZero := u.ToDecimal(0)
	return dZero.Sub(d1).Round(u.Prec)
}

// 绝对值
func (u *DecimalHelper) Abs(d1 decimal.Decimal) decimal.Decimal {
	return d1.Abs().Round(u.Prec)
}

// 是否相等
func (u *DecimalHelper) Equal(d1 decimal.Decimal, d2 decimal.Decimal) bool {
	return d1.Equal(d2)
}

// 和0是否相等
func (u *DecimalHelper) EqualByZero(d1 decimal.Decimal) bool {
	dZero := u.ToDecimal(0)
	return d1.Equal(dZero)
}

// 是否小于
func (u *DecimalHelper) LessThan(d1 decimal.Decimal, d2 decimal.Decimal) bool {
	return d1.LessThan(d2)
}

// 和0是否小于
func (u *DecimalHelper) LessThanByZero(d1 decimal.Decimal) bool {
	dZero := u.ToDecimal(0)
	return d1.LessThan(dZero)
}

// 是否小于等于
func (u *DecimalHelper) LessThanOrEqual(d1 decimal.Decimal, d2 decimal.Decimal) bool {
	return d1.LessThanOrEqual(d2)
}

// 和0是否小于等于
func (u *DecimalHelper) LessThanOrEqualByZero(d1 decimal.Decimal) bool {
	dZero := u.ToDecimal(0)
	return d1.LessThanOrEqual(dZero)
}

// 是否大于
func (u *DecimalHelper) GreaterThan(d1 decimal.Decimal, d2 decimal.Decimal) bool {
	return d1.GreaterThan(d2)
}

// 和0是否大于
func (u *DecimalHelper) GreaterThanByZero(d1 decimal.Decimal) bool {
	dZero := u.ToDecimal(0)
	return d1.GreaterThan(dZero)
}

// 是否大于等于
func (u *DecimalHelper) GreaterThanOrEqual(d1 decimal.Decimal, d2 decimal.Decimal) bool {
	return d1.GreaterThanOrEqual(d2)
}

// 和0是否大于等于
func (u *DecimalHelper) GreaterThanOrEqualByZero(d1 decimal.Decimal) bool {
	dZero := u.ToDecimal(0)
	return d1.GreaterThanOrEqual(dZero)
}
