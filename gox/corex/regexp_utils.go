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
