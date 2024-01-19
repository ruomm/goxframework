/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/19 16:38
 * @version 1.0
 */
package loggerx

type LogConfigs struct {
	Level        string `yaml:"level"`
	StdOut       bool   `yaml:"stdOut"`
	MaxSize      int    `yaml:"maxSize"`
	MaxBackups   int    `yaml:"maxBackups"`
	MaxAges      int    `yaml:"maxAges"`
	Compress     bool   `yaml:"compress"`
	ServiceName  string `yaml:"serviceName"`
	InstanceName string `yaml:"instanceName"`
}
