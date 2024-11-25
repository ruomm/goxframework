/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/11/22 13:49
 * @version 1.0
 */
package corex

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// uft8模式：计算字符串长度
func Uft8Len(str string) int {
	if str == "" {
		return 0
	}
	lenRunes := utf8.RuneCountInString(str)
	return lenRunes
}

// uft8模式：获取一个字符串特定位置的字符串，Utf8At("hello，你好！", 6) = "你"
func Utf8At(str string, index int) string {
	lenRunes := utf8.RuneCountInString(str)
	if index < 0 || index > lenRunes {
		errStr := fmt.Errorf("invalid argument: index %d out of bounds [0:%d]", index, lenRunes)
		panic(errStr)
	}
	sb := strings.Builder{}
	runeIndex := 0
	for _, r := range str {
		if index == runeIndex {
			sb.WriteRune(r)
		}
		runeIndex++
	}
	return sb.String()
}

// uft8模式：获取一个字符串特定开始和结束位置的中间字符串，startIndex:索引开始位置(包含)，endIndex:索引结束位置(不包含)
// Utf8Sub("hello，你好！", 6, -1)="你好！"，Utf8Sub("hello，你好！", 6, 8)"=你好"
func Utf8Sub(str string, startIndex int, endIndex int) string {
	lenRunes := utf8.RuneCountInString(str)
	if startIndex < 0 || startIndex > lenRunes {
		errStr := fmt.Errorf("invalid argument: startIndex %d out of bounds [0:%d]", startIndex, lenRunes)
		panic(errStr)
	}
	if endIndex != -1 {
		if endIndex < 0 || endIndex > lenRunes {
			errStr := fmt.Errorf("invalid argument: endIndex %d out of bounds [0:%d]", endIndex, lenRunes)
			panic(errStr)
		}
		if endIndex < startIndex {
			errStr := fmt.Errorf("invalid argument: endIndex %d is less than startIndex %d", endIndex, startIndex)
			panic(errStr)
		}
	} else {
		endIndex = lenRunes
	}

	sb := strings.Builder{}
	runeIndex := 0
	for _, r := range str {
		if runeIndex >= startIndex && runeIndex < endIndex {
			sb.WriteRune(r)
		}
		runeIndex++
	}
	return sb.String()
}

// uft8模式：获取字符串substr在字符串s中的首位索引位置，没有则返回-1
// 如：Utf8Index("hello，你好啊！", "好啊") = 7
func Utf8Index(s string, substr string) int {
	index := strings.Index(s, substr)
	if index <= 0 {
		return index
	}
	prefixStr := s[0:index]
	lenRunes := utf8.RuneCountInString(prefixStr)
	return lenRunes
}
