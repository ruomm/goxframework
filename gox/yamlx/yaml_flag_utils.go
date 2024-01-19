package xconfutils

import (
	"flag"
	"fmt"
	"github.com/ruomm/goxframework/gox/reflectx"
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

/*
* 解析yaml配置文件为对象
* 依据envKey激活环境配置文件，读取环境配置文件和指定的yaml配置文件，环境配置文件的值会覆盖指定的yaml配置文件的值。
 */
func ParseYamlFileByFlag(obj interface{}) error {
	cmdCofig := initCMDConfig()
	err := ParseYamlFileByEnv(*cmdCofig.ConfYaml, "profileActive", *cmdCofig.ProfileActive, obj)
	if err != nil {
		panic(fmt.Sprintf("config load error, cause by ParseYamlFileByFlag: %v", err))
	}
	errXcp := reflectx.XReflectCopy(cmdCofig, obj, true)
	if errXcp != nil {
		panic(fmt.Sprintf("config load error, cause by XReflectCopy:", err))
	} else {
		return nil
	}
}
