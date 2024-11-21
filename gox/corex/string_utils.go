/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/24 11:23
 * @version 1.0
 */
package corex

import (
	"bufio"
	"regexp"
	"strings"
)

// 判断字符串是否匹配，支持正则表达式，支持开头结尾*通配符，支持严格匹配，空字符串不匹配
func MatchStringCommon(pattern string, s string) bool {
	if len(pattern) <= 0 {
		return false
	}
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

// 判断字符串是否匹配，支持正则表达式，支持开头结尾*通配符，支持严格匹配，空字符串匹配
func MatchStringEmptyPatternPass(pattern string, s string) bool {
	if len(pattern) <= 0 {
		return true
	}
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

/**
 * 字符串转换为NO BOM格式
 *
 * @param str 原始字符串
 * @return NO BOM字符串
 */
func StringToNoBom(str string) string {
	lenStr := len(str)
	if lenStr <= 0 {
		return str
	}
	isBom := false
	chars := []rune(str)
	charCode := chars[0]
	if charCode == rune(65279) {
		isBom = true
	}
	if isBom {
		return str[1:]
	} else {
		return str
	}
}

// 逐行转换string为slice，emptyRetain=true表示空的行也保留
func StringReadByLine(str string, emptyRetain bool) []string {
	lenStr := len(str)
	if lenStr <= 0 {
		return nil
	}
	var resultSlice []string
	// 创建一个Scanner对象
	scanner := bufio.NewScanner(strings.NewReader(str))
	// 设置扫描器的分隔函数为ScanLines
	scanner.Split(bufio.ScanLines)
	//firstFlag := true
	// 遍历读取每一行
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行的文本
		//if firstFlag {
		//	firstFlag = false
		//	if firsrBomRemove {
		//		line = StringToNoBom(line)
		//	}
		//}
		if emptyRetain {
			resultSlice = append(resultSlice, line)
		} else if len(line) > 0 {
			resultSlice = append(resultSlice, line)
		}
	}
	return resultSlice
}

// 逐行转换string为slice，首行去除bom，emptyRetain=true表示空的行也保留
func StringReadByLineNoBom(str string, emptyRetain bool) []string {
	lenStr := len(str)
	if lenStr <= 0 {
		return nil
	}
	var resultSlice []string
	// 创建一个Scanner对象
	scanner := bufio.NewScanner(strings.NewReader(str))
	// 设置扫描器的分隔函数为ScanLines
	scanner.Split(bufio.ScanLines)
	firstFlag := true
	// 遍历读取每一行
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行的文本
		if firstFlag {
			firstFlag = false
			line = StringToNoBom(line)
		}
		if emptyRetain {
			resultSlice = append(resultSlice, line)
		} else if len(line) > 0 {
			resultSlice = append(resultSlice, line)
		}
	}
	return resultSlice
}
