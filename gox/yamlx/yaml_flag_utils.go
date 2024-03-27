package yamlx

import (
	"flag"
	"fmt"
	"github.com/ruomm/goxframework/gox/refx"
)

type CommondConfigs struct {
	ConfYaml       *string
	WebBindPort    *int
	ProfileActive  *string
	WebContextPath *string
	AuthGlobalTag  *string
}

func initCMDConfig() *CommondConfigs {
	cmdConfigs := CommondConfigs{
		ConfYaml:       flag.String("c", "config/conf.yaml", "config file name, yaml format."),
		WebBindPort:    flag.Int("p", 0, "config server port.\ndefault is define in config file."),
		ProfileActive:  flag.String("env", "", "config run environment, will load environment config file. \ndefault is define in config file."),
		WebContextPath: flag.String("uri", "", "web context path，will left append to api uri. \ndefault is define in config file."),
		AuthGlobalTag:  flag.String("au", "", "change this value，all users authorization will invalid. \ndefault is define in config file."),
	}
	flag.Parse()
	return &cmdConfigs
}

func initCMDConfigLite() *CommondConfigs {
	cmdConfigs := CommondConfigs{
		ConfYaml:      flag.String("c", "config/conf.yaml", "config file name, yaml format."),
		ProfileActive: flag.String("env", "", "config run environment, will load environment config file. \ndefault is define in config file."),
	}
	flag.Parse()
	return &cmdConfigs
}

func initCMDConfigYsy() *CommondConfigs {
	cmdConfigs := CommondConfigs{
		ConfYaml:      flag.String("c", "", "config file name, yaml format.\nnot set use conf.yaml as default"),
		ProfileActive: flag.String("env", "", "config run environment, will load environment config file. \ndefault is define in config file."),
		WebBindPort:   flag.Int("p", 0, "config server port.\ndefault is define in config file."),
	}
	flag.Parse()
	return &cmdConfigs
}

/*
* 解析yaml配置文件为对象
* 依据envKey激活环境配置文件，读取环境配置文件和指定的yaml配置文件，环境配置文件的值会覆盖指定的yaml配置文件的值。
 */
func ParseYamlFileByFlag(obj interface{}, options ...YamlxOption) error {
	cmdCofig := initCMDConfig()
	err := ParseYamlFileByEnv(*cmdCofig.ConfYaml, "profileActive", *cmdCofig.ProfileActive, obj, options...)
	if err != nil {
		panic(fmt.Sprintf("config load error, cause by ParseYamlFileByFlag: %v", err))
	}
	errXcp, _ := refx.XRefStructCopy(cmdCofig, obj)
	if errXcp != nil {
		panic(fmt.Sprintf("config load error, cause by XReflectCopy:", err))
	} else {
		return nil
	}
}

/*
* 解析yaml配置文件为对象
* 依据envKey激活环境配置文件，读取环境配置文件和指定的yaml配置文件，环境配置文件的值会覆盖指定的yaml配置文件的值。
 */
func ParseYamlFileByFlagLite(obj interface{}, options ...YamlxOption) error {
	cmdCofig := initCMDConfigLite()
	err := ParseYamlFileByEnv(*cmdCofig.ConfYaml, "profileActive", *cmdCofig.ProfileActive, obj, options...)
	if err != nil {
		panic(fmt.Sprintf("config load error, cause by ParseYamlFileByFlag: %v", err))
	}
	errXcp, _ := refx.XRefStructCopy(cmdCofig, obj)
	if errXcp != nil {
		panic(fmt.Sprintf("config load error, cause by XReflectCopy:", err))
	} else {
		return nil
	}
}

/*
* 解析yaml配置文件为对象
* 依据envKey激活环境配置文件，读取环境配置文件和指定的yaml配置文件，环境配置文件的值会覆盖指定的yaml配置文件的值。
 */
func ParseYamlFileByFlagWithPath(filePath string, obj interface{}, options ...YamlxOption) error {
	cmdCofig := initCMDConfigYsy()
	confFilePath := *cmdCofig.ConfYaml
	if len(confFilePath) <= 0 {
		confFilePath = filePath
	}
	if len(confFilePath) <= 0 {
		confFilePath = "conf.yaml"
	}
	err := ParseYamlFileByEnv(confFilePath, "profileActive", *cmdCofig.ProfileActive, obj, options...)
	if err != nil {
		panic(fmt.Sprintf("config load error, cause by ParseYamlFileByFlag: %v", err))
	}
	errXcp, _ := refx.XRefStructCopy(cmdCofig, obj)
	if errXcp != nil {
		panic(fmt.Sprintf("config load error, cause by XReflectCopy:", err))
	} else {
		return nil
	}
}
