/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/5/14 10:11
 * @version 1.0
 */
package corex

import (
	"fmt"
	"regexp"
)

// 验证字符是否符合字符串规则，主要用于正则表达式验证
func IsRegexMatch(reg_pattern string, val string) bool {
	re, err := regexp.Compile(reg_pattern)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if !re.MatchString(val) {
		return false
	} else {
		return true
	}
}
