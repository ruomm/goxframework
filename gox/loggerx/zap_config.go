/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 16:38
 * @version 1.0
 */
package loggerx

import "go.uber.org/zap"

type LogConfigs struct {
	LogPath         string `xref:"LogPath"`
	Level           string `xref:"Level"`
	StdOut          bool   `xref:"StdOut"`
	MaxSize         int    `xref:"MaxSize"`
	MaxBackups      int    `xref:"MaxBackups"`
	MaxAges         int    `xref:"MaxAges"`
	Compress        bool   `xref:"Compress"`
	ServiceName     string `xref:"ServiceName"`
	InstanceName    string `xref:"InstanceName"`
	TextMode        bool   `xref:"TextMode;tidy"`
	StatsTimeEnable bool   `xref:"StatsTimeEnable;tidy"` // 普通日志是否打开耗时统计功能
	//Branch          string                 `xref:"Branch;tidy"`          // Git的分支名称，如：dev/v1.5.0
	//Version         string                 `xref:"Version;tidy"`         // Git的commit的hash值，如：297b1b7c039e918d3b006d954bf27e415ae5599d
	ZapFields []zap.Field `xref:"ZapFields;tidy"`
}
