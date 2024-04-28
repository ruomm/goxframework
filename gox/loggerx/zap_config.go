/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 16:38
 * @version 1.0
 */
package loggerx

type LogConfigs struct {
	Level        string `xref:"Level"`
	StdOut       bool   `xref:"StdOut"`
	MaxSize      int    `xref:"MaxSize"`
	MaxBackups   int    `xref:"MaxBackups"`
	MaxAges      int    `xref:"MaxAges"`
	Compress     bool   `xref:"Compress"`
	ServiceName  string `xref:"ServiceName"`
	InstanceName string `xref:"InstanceName"`
	TextMode     bool   `xref:"TextMode;tidy"`
}
