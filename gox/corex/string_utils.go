/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/24 11:23
 * @version 1.0
 */
package corex

import (
	"regexp"
	"strings"
)

// 判断字符串是否匹配，支持正则表达式，支持开头结尾*通配符，支持严格匹配
func MatchStringCommon(pattern string, s string) bool {
	if len(s) <= 0 {
		return false
	}
	if strings.HasPrefix(pattern, "^") || strings.HasSuffix(pattern, "$") {
		//使用正则校验
		isMatch, err := regexp.MatchString(pattern, s)
		if err != nil {
			return false
		} else {
			return isMatch
		}
	} else if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
		//使用头尾通配符匹配模式
		realPattern := pattern[1 : len(pattern)-1]
		if len(realPattern) <= 0 {
			return false
		} else if strings.Contains(s, realPattern) {
			return true
		} else {
			return false
		}
	} else if strings.HasSuffix(pattern, "*") {
		//使用尾部通配符匹配模式
		realPattern := pattern[0 : len(pattern)-1]
		if len(realPattern) <= 0 {
			return false
		} else if strings.HasPrefix(s, realPattern) {
			return true
		} else {
			return false
		}
	} else if strings.HasPrefix(pattern, "*") {
		//使用头部通配符匹配模式
		realPattern := pattern[1:len(pattern)]
		if len(realPattern) <= 0 {
			return false
		} else if strings.HasSuffix(s, realPattern) {
			return true
		} else {
			return false
		}
	} else if pattern == s {
		//使用严格匹配模式
		return true
	} else {
		return false
	}
}
