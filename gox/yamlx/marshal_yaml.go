/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/23 21:06
 * @version 1.0
 */
package yamlx

import "github.com/ruomm/goxframework/gox/corex"

type StoreSize int64
type TimeSeconds int64
type TimeMilliSeconds int64

func (e *StoreSize) UnmarshalYAML(unmarshal func(interface{}) error) error {
	ymalString := ""
	err := unmarshal(&ymalString)
	if err != nil {
		return err
	}
	byteSize, err := corex.ByteSizeParse(ymalString)
	if err != nil {
		return err
	}
	*e = StoreSize(byteSize)
	return nil
}

func (e *TimeSeconds) UnmarshalYAML(unmarshal func(interface{}) error) error {
	ymalString := ""
	err := unmarshal(&ymalString)
	if err != nil {
		return err
	}
	timeNumber, err := corex.TimeNumberParse(ymalString, true)
	if err != nil {
		return err
	}
	*e = TimeSeconds(timeNumber)
	return nil
}

func (e *TimeMilliSeconds) UnmarshalYAML(unmarshal func(interface{}) error) error {
	ymalString := ""
	err := unmarshal(&ymalString)
	if err != nil {
		return err
	}
	timeNumber, err := corex.TimeNumberParse(ymalString, false)
	if err != nil {
		return err
	}
	*e = TimeMilliSeconds(timeNumber)
	return nil
}
