package yamlx

import (
	"github.com/ruomm/goxframework/gox/corex"
	"github.com/ruomm/goxframework/gox/mergox"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// 属性空值设置 for refx.
type YamlxOption struct {
	f func(*yamlxOptions)
}

type yamlxOptions struct {
	envFileCanEmpty bool //环境配置文件非强制加载
}

// 设置关联tag的名称，不设置默认为xref
func YamlxEnvFileCanEmpty(canEmpty bool) YamlxOption {
	return YamlxOption{func(do *yamlxOptions) {
		do.envFileCanEmpty = canEmpty
	}}
}

/*
* 解析yaml配置文件为对象
* 只能读取指定的配置文件，不会激活环境配置文件
 */
func ParseYamlFile(filePath string, obj interface{}) error {
	byteData, err := ioutil.ReadFile(corex.GetAbsDir(filePath))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(byteData, obj)
}

/*
* 解析yaml配置文件为对象
* filePath：配置文件路径
* envKey：配置文件中的环境配置key
* envValue：指定特定的环境配置值，不为空时候envKey失效
* obj：解析后的配置对象
* # envKey和envValue用来指定环境配置值。
* # envValue有值激活envValue值的环境配置值，没有值则依据{envKey}去指定的配置文件中查找环境配置值。
* # 环境配置值存在时候，会额外加载环境配置文件。环境配置名称为指定的配置文件名称拼接上"-环境配置值"。
* # 如配置文件为：config/conf.yaml，环境配置值为：dev"，则激活的配置文件为：config/conf-dev.yaml。
* # 如是指定了环境配置值，激活的配置文件必须存在且可以解析，否则程序异常报错。
 */
func ParseYamlFileByEnv(filePath string, envKey string, envValue string, obj interface{}, options ...YamlxOption) error {
	mapConfigs, configEnv, err := ParseYamlFileToMapByEnv(filePath, envKey, envValue, options...)
	if err != nil {
		return err
	}
	// 将map转换成YAML格式的二进制数据
	yamlData, err := yaml.Marshal(mapConfigs)

	if err != nil {
		return err
	}
	if len(envKey) > 0 && len(configEnv) > 0 {
		yamlStr := string(yamlData)
		//if configEnv != "" {
		//	yamlStr = strings.Replace(yamlStr, "${env}", configEnv, -1)
		//	yamlStr = strings.Replace(yamlStr, "${"+envKey+"}", configEnv, -1)
		//} else {
		//	yamlStr = strings.Replace(yamlStr, "${env}", "dev", -1)
		//	yamlStr = strings.Replace(yamlStr, "${"+envKey+"}", "dev", -1)
		//}
		yamlStr = strings.Replace(yamlStr, "${env}", configEnv, -1)
		yamlStr = strings.Replace(yamlStr, "${"+envKey+"}", configEnv, -1)
		yamlData = []byte(yamlStr)

	}
	return yaml.Unmarshal(yamlData, obj)
}

/*
* 解析yaml配置文件为map
* 只能读取指定的配置文件，不会激活环境配置文件
 */
func ParseYamlFileToMap(filePath string) (*map[string]interface{}, error) {
	yamlFilePath := corex.GetAbsDir(filePath)
	_, errFile := corex.IsFileExitWithErr(yamlFilePath)
	if errFile != nil {
		return nil, errFile
	}
	yamlFile, err := os.Open(yamlFilePath)
	if err != nil {
		//fmt.Println("open file err = ", err)
		return nil, err
	}
	decode := yaml.NewDecoder(yamlFile)
	fileMap := make(map[string]interface{})
	err = decode.Decode(&fileMap)
	if err != nil {
		//fmt.Println("onekit yaml decode has error:%v", err)
		return nil, err
	} else {
		return &fileMap, err
	}
}

/*
* 解析yaml配置文件为map
* 依据envKey激活环境配置文件，读取环境配置文件和指定的yaml配置文件，环境配置文件的值会覆盖指定的yaml配置文件的值。
 */
func ParseYamlFileToMapByEnv(filePath string, envKey string, envValue string, options ...YamlxOption) (*map[string]interface{}, string, error) {
	do := yamlxOptions{}
	for _, option := range options {
		option.f(&do)
	}
	// 读取全局配置文件
	mapGlob, errGlob := ParseYamlFileToMap(filePath)
	if errGlob != nil {
		return nil, envValue, errGlob
	}
	// 获取全局配置文件中的环境参数值
	var configEnv string
	if len(envKey) <= 0 {
		configEnv = envValue
	} else if len(envValue) > 0 {
		configEnv = envValue
		if len(envKey) > 0 {
			(*mapGlob)[envKey] = envValue
		}
	} else {
		configActiveVal := (*mapGlob)[envKey]
		if nil == configActiveVal {
			configEnv = ""
		} else {
			configActiveTyp := reflect.TypeOf(configActiveVal)
			if configActiveTyp == nil {
				configEnv = ""
			} else if configActiveTyp.Kind() == reflect.String {
				configEnv = configActiveVal.(string)
			} else if configActiveTyp.Kind() == reflect.Int {
				configEnvInt := configActiveVal.(int)
				configEnv = strconv.Itoa(configEnvInt)
			} else {
				configEnv = ""
			}
		}
	}
	var filePathEnv string
	if len(configEnv) > 0 {
		indexSep := strings.LastIndex(filePath, ".")
		if indexSep <= 0 {
			filePathEnv = ""
		} else {
			var build strings.Builder
			build.WriteString(filePath[0:indexSep])
			build.WriteString("-")
			build.WriteString(configEnv)
			build.WriteString(filePath[indexSep:])
			filePathEnv = build.String()
		}
	}
	if len(filePathEnv) > 0 {
		mapEnv, errEnv := ParseYamlFileToMap(filePathEnv)
		if errEnv != nil {
			//fmt.Println("open file err = ", errEnv)
			if do.envFileCanEmpty && os.IsNotExist(errEnv) {
				return mapGlob, configEnv, nil
			} else {
				return nil, configEnv, errEnv
			}
		}
		errMergo := mergox.Map(mapGlob, mapEnv, mergox.WithOverride)
		if errMergo != nil {
			//fmt.Println("open file err = ", errMergo)
			return nil, configEnv, errMergo
		} else {
			return mapGlob, configEnv, nil
		}
	} else {
		return mapGlob, configEnv, nil
	}
}
