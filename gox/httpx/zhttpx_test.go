/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 21:38
 * @version 1.0
 */
package httpx

import (
	"fmt"
	"testing"
)

type ConfigGpuSpecDeleteReq struct {
	GpuSpecId uint `json:"-" xreq_param:"id;order=66'" validate:"min=1" xvalid_error:"GPU规格编号必须填写，且必须是正整数"`
	//GpuSpecName string `json:"gpuSpecName" validate:"min=1,max=64" xvalid_error:"GPU规格名称必须填写，长度1-64位字符"`
	GpuBrand   string `json:"-" xreq_query:"gpuBrand" validate:"min=1,max=32" xvalid_error:"GPU品牌必须填写，长度1-32位字符"`
	GpuModel   string `json:"-" xreq_header:"gpuModel" validate:"min=1,max=32" xvalid_error:"GPU型号必须填写，长度1-32位字符"`
	CardMemory int    `json:"-" xreq_param:"cardMemory;order=58" validate:"min=1,max=10000" xvalid_error:"显存大小必须填写，范围1-10000G"`
	Memory     int    `json:"-" xreq_query:"memory" validate:"min=1,max=10000" xvalid_error:"内存大小必须填写，范围1-10000G"`
	Memory2    int    `json:"-" xreq_header:"memory2" validate:"min=1,max=10000" xvalid_error:"内存大小必须填写，范围1-10000G"`
}

func (u ConfigGpuSpecDeleteReq) HttpxMethod() string {
	return "DELETE"
}

func (u CommonResult) HttpxMethod() string {
	return "DELETE"
}

type CommonResult struct {
	TraceId string                 `json:"traceId,omitempty" newtag:"traceId"`
	Code    int                    `json:"code" newtag:"code"`
	UserMsg string                 `json:"msg,omitempty" newtag:"msg"`           //用户查看的信息，可读性更强
	LogKV   map[string]interface{} `json:"errorMsg,omitempty" newtag:"errorMsg"` //打印日志的信息，携带错误详情，便于追查问题
	Data    interface{}            `json:"data,omitempty" newtag:"-"`
}

func TestS2S(t *testing.T) {
	//a := 123456.567
	//a := true
	req := ConfigGpuSpecDeleteReq{

		GpuSpecId:  456,
		GpuBrand:   "nvidia",
		GpuModel:   "A800",
		CardMemory: 123,
		Memory:     456,
	}
	result := CommonResult{}
	DoHttpJson("http://localhost:8010/api/v1/configspec/gpu/delete", "DELETE", req, &result)
	fmt.Print(result)

}
