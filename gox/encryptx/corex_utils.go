/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/6/21 13:25
 * @version 1.0
 */
package encryptx

import (
	"bufio"
	"math/rand"
	"strings"
)

// 转换string为slice
func stringReadByLineNoBom(str string, emptyRetain bool) []string {
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
			line = stringToNoBom(line)
		}
		if emptyRetain {
			resultSlice = append(resultSlice, line)
		} else if len(line) > 0 {
			resultSlice = append(resultSlice, line)
		}
	}
	return resultSlice
}

/**
 * 字符串转换为NO BOM格式
 *
 * @param str 原始字符串
 * @return NO BOM字符串
 */
func stringToNoBom(str string) string {
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

func generateToken(token_len int) string {
	return xGenerateToken(token_len, false)
}

func generateTokenNoStartWithZero(token_len int) string {
	return xGenerateToken(token_len, true)
}

//func xRandBool() bool {
//	if rand.Intn(2) == 1 {
//		return true
//	} else {
//		return false
//	}
//}

func xGenerateToken(token_len int, noZeroStart bool) string {
	realLen := 0
	if token_len > 0 {
		realLen = token_len
	} else {
		realLen = 4
	}
	realDicts := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	tokenresult := ""
	for i := 0; i < realLen; i++ {
		if i == 0 && noZeroStart {
			tokenresult = tokenresult + xGenerateNoZeroStr(realDicts)
		} else {
			tmpIndx := rand.Intn(len(realDicts))
			tokenresult = tokenresult + realDicts[tmpIndx:tmpIndx+1]
		}
	}
	return tokenresult
}

// 生成首字符串不为0的随机字符串
func xGenerateNoZeroStr(realDicts string) string {
	resultStr := ""
	for {
		tmpIndx := rand.Intn(len(realDicts))
		resultStr = realDicts[tmpIndx : tmpIndx+1]
		if resultStr == "0" {
			continue
		} else {
			break
		}
	}
	return resultStr
}
