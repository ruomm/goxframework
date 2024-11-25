/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/10/28 16:33
 * @version 1.0
 */
package corex

import (
	"errors"
	"strings"
)

// 判断是否Hex字符串，strictMatch:true时候字符串长度必须为偶数才判断
func IsHexString(str string, strictMatch bool) bool {
	s := ""
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		s = str[2:]
	} else {
		s = str
	}
	if len(s) <= 0 {
		return false
	}
	if strictMatch && len(s)%2 != 0 {
		return false
	}
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}

// 解析hex字符串为正常字符串，strictMatch:true时候字符串长度必须为偶数才解析
func ToHexString(str string, strictMatch bool) (string, error) {
	s := ""
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		s = str[2:]
	} else {
		s = str
	}
	if len(s) <= 0 {
		return s, errors.New("invalid hex string，string is empty")
	}
	if strictMatch && len(s)%2 != 0 {
		return "", errors.New("invalid hex string，string length is not even number")
	}
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return "", errors.New("invalid hex string，string contain no hex letter")
		}
	}
	return s, nil
}
