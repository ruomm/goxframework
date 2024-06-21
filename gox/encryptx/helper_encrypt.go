/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/17 21:23
 * @version 1.0
 */
package encryptx

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/rand"
	"strings"
)

type MODE_KEY int
type MODE_ENCODE int
type MODE_PADDING string

const (
	MODE_KEY_BASE64       MODE_KEY     = 1
	MODE_KEY_HEX_LOWER    MODE_KEY     = 2
	MODE_KEY_HEX_UPPER    MODE_KEY     = 3
	MODE_KEY_STRING       MODE_KEY     = 4
	MODE_KEY_PEM          MODE_KEY     = 5
	MODE_ENCODE_BASE64    MODE_ENCODE  = 1
	MODE_ENCODE_HEX_LOWER MODE_ENCODE  = 2
	MODE_ENCODE_HEX_UPPER MODE_ENCODE  = 3
	MODE_PADDING_PKCS7    MODE_PADDING = "pkcs7"
	MODE_PADDING_PKCS5    MODE_PADDING = "pkcs5"
	FILE_BUFFER_SIZE                   = 1024
)

type EncryptHelper interface {
	// key模式
	ModeOfKey() MODE_KEY
	// 字节转字符串编码方案
	ModeOfEncode() MODE_ENCODE
	// Padding的模式
	ModeOfPadding() MODE_PADDING

	// 设置key字节数组
	SetKeyData(key []byte)
	// 设置iv字节数组
	SetIVData(iv []byte)
	// 设置key字符串
	SetKeyString(keyStr string) error
	// 设置iv字符串
	SetIVString(ivStr string) error
	// 获取key字节数组
	GetKeyData() []byte
	// 设置iv字节数组
	GetIVData() []byte
	// 获取key字符串
	GetKeyString() (string, error)
	// 获取iv字符串
	GetIVString() (string, error)
	// 设置Blocksize
	SetBlockSize(blockSize int)
	// 获取Blocksize
	GetBlockSize() (int, error)
	// 生成key或iv字节数组
	GenKeyIvData(len int) ([]byte, error)
	// 生成key或iv字符串
	GenKeyIvString(len int) (string, error)

	//// 生成iv字节数组
	//GenIVData() ([]byte, error)
	//// 生成iv字符串
	//GenIVString() (string, error)
	// 还原key或iv字符串为key或iv字节数组
	RestoreKeyIV(keyStr string) ([]byte, error)
	// ECB加密字节数组
	EncDataECB(data []byte) ([]byte, error)
	// CBC加密字节数组
	EncDataCBC(data []byte) ([]byte, error)
	// ECB解密字节数组
	DecDataECB(dataEnc []byte) ([]byte, error)
	// CBC解密字节数组
	DecDataCBC(dataEnc []byte) ([]byte, error)
	// ECB加密字符串
	EncStringECB(str string) (string, error)
	// CBC加密字符串
	EncStringCBC(str string) (string, error)
	// ECB解密字符串
	DecStringECB(strEnc string) (string, error)
	// CBC解密字符串
	DecStringCBC(strEnc string) (string, error)
	// ECB加密字符串
	EncFileECB(fileSrc string, fileEnc string) error
	// CBC加密字符串
	EncFileCBC(fileSrc string, fileEnc string) error
	// ECB解密字符串
	DecFileECB(fileEnc string, fileDest string) error
	// CBC解密字符串
	DecFileCBC(fileEnc string, fileDest string) error
	// Padding&UnPadding
	Padding(data []byte, blockSize int) []byte
	UnPadding(data []byte, blockSize int) []byte
}

func ParseKeyMode(keyMode MODE_KEY) MODE_KEY {
	if keyMode == MODE_KEY_BASE64 || keyMode == MODE_KEY_HEX_LOWER || keyMode == MODE_KEY_HEX_UPPER || keyMode == MODE_KEY_STRING || keyMode == MODE_KEY_PEM {
		return keyMode
	} else {
		return MODE_KEY_STRING
	}
}
func ParseEncodeMode(encodeMode MODE_ENCODE) MODE_ENCODE {
	if encodeMode == MODE_ENCODE_BASE64 || encodeMode == MODE_ENCODE_HEX_LOWER || encodeMode == MODE_ENCODE_HEX_UPPER {
		return encodeMode
	} else {
		return MODE_ENCODE_BASE64
	}
}
func ParsePaddingMode(paddingMode MODE_PADDING) MODE_PADDING {
	if paddingMode == MODE_PADDING_PKCS7 || paddingMode == MODE_PADDING_PKCS5 {
		return paddingMode
	} else {
		return MODE_PADDING_PKCS7
	}
}
func EncodingToString(encodeMode MODE_ENCODE, data []byte) (string, error) {
	//if len(data) <= 0 {
	//	return "", errors.New("lenght of data([]byte) must greater than 0")
	//}
	if encodeMode == MODE_ENCODE_BASE64 {
		return base64.StdEncoding.EncodeToString(data), nil
	} else if encodeMode == MODE_ENCODE_HEX_LOWER {
		return strings.ToLower(hex.EncodeToString(data)), nil
	} else if encodeMode == MODE_ENCODE_HEX_UPPER {
		return strings.ToUpper(hex.EncodeToString(data)), nil
	} else {
		return base64.StdEncoding.EncodeToString(data), nil
	}
}
func DecodingToByte(encodeMode MODE_ENCODE, dataStr string) ([]byte, error) {
	//if len(dataStr) <= 0 {
	//	return []byte{}, errors.New("lenght of dataStr(string) must greater than 0")
	//}
	if encodeMode == MODE_ENCODE_BASE64 {
		return base64.StdEncoding.DecodeString(dataStr)
	} else if encodeMode == MODE_ENCODE_HEX_LOWER {
		return hex.DecodeString(strings.ToLower(dataStr))
	} else if encodeMode == MODE_ENCODE_HEX_UPPER {
		return hex.DecodeString(strings.ToLower(dataStr))
	} else {
		return base64.StdEncoding.DecodeString(dataStr)
	}
}

func KeyIVStringToByte(keyMode MODE_KEY, keyStr string) ([]byte, error) {
	if len(keyStr) <= 0 {
		return []byte{}, errors.New("lenght of keyStr(string) must greater than 0")
	}
	if keyMode == MODE_KEY_BASE64 {
		return base64.StdEncoding.DecodeString(keyStr)
	} else if keyMode == MODE_KEY_HEX_LOWER {
		return hex.DecodeString(strings.ToLower(keyStr))
	} else if keyMode == MODE_KEY_HEX_UPPER {
		return hex.DecodeString(strings.ToLower(keyStr))
	} else if keyMode == MODE_KEY_STRING {
		return []byte(keyStr), nil
	} else if keyMode == MODE_KEY_PEM {
		return ReadFormatKey(keyStr)
	} else {
		return []byte(keyStr), nil
	}
}
func KeyIVByteToString(keyMode MODE_KEY, keyData []byte) (string, error) {
	if len(keyData) <= 0 {
		return "", errors.New("lenght of keyData([]byte) must greater than 0")
	}
	if keyMode == MODE_KEY_BASE64 {
		return base64.StdEncoding.EncodeToString(keyData), nil
	} else if keyMode == MODE_KEY_HEX_LOWER {
		return strings.ToLower(hex.EncodeToString(keyData)), nil
	} else if keyMode == MODE_KEY_HEX_UPPER {
		return strings.ToUpper(hex.EncodeToString(keyData)), nil
	} else if keyMode == MODE_KEY_STRING {
		return string(keyData), nil
	} else if keyMode == MODE_KEY_PEM {
		tag := "AES KEY OR IV"
		return FormatKeyByData(keyData, tag)
	} else {
		return string(keyData), nil
	}
}

func GenKeyString(keyMode MODE_KEY, keyLen int) (string, error) {
	if keyLen <= 0 {
		return "", errors.New("KeyLen must greater than 0")
	}
	if keyMode == MODE_KEY_BASE64 {
		keyBytes := GenKeyData(keyLen)
		return base64.StdEncoding.EncodeToString(keyBytes), nil
	} else if keyMode == MODE_KEY_HEX_LOWER {
		keyBytes := GenKeyData(keyLen)
		return strings.ToLower(hex.EncodeToString(keyBytes)), nil
	} else if keyMode == MODE_KEY_HEX_UPPER {
		keyBytes := GenKeyData(keyLen)
		return strings.ToUpper(hex.EncodeToString(keyBytes)), nil
	} else if keyMode == MODE_KEY_STRING {
		return genTokenString(keyLen), nil
	} else if keyMode == MODE_KEY_PEM {
		keyBytes := GenKeyData(keyLen)
		tag := "AES KEY OR IV"
		return FormatKeyByData(keyBytes, tag)
	} else {
		return genTokenString(keyLen), nil
	}
}

func GenKeyData(keyLen int) []byte {
	if keyLen <= 0 {
		return []byte{}
	}
	var build bytes.Buffer
	for i := 0; i < keyLen; i++ {
		randIndex := rand.Intn(0xFF)
		build.WriteByte(byte(randIndex))
	}
	return build.Bytes()
}

func genTokenString(keyLen int) string {
	if keyLen <= 0 {
		return ""
	}
	dicts := "0123456789abcdefghijklmnopqrstuvwxyz"
	dictsLen := len(dicts)
	var build strings.Builder
	for i := 0; i < keyLen; i++ {
		randIndex := rand.Intn(dictsLen)
		build.WriteString(dicts[randIndex : randIndex+1])
	}
	return build.String()
}

// 动态计算blockSize
//func calBlockSize(sizeHelper int, sizeChiper int, sizeKey int, blockSizeByKey bool) int {
//	if sizeHelper > 0 {
//		return sizeHelper
//	}
//	if sizeChiper <= 0 && sizeKey <= 0 {
//		return 16
//	} else if sizeKey <= 0 {
//		return sizeChiper
//	} else if sizeChiper <= 0 {
//		count := sizeKey / 16
//		if count <= 0 {
//			return 16
//		} else if sizeKey%16 == 0 {
//			return count * 16
//		} else {
//			return (count + 1) * 16
//		}
//	} else {
//		if blockSizeByKey {
//			count := sizeKey / 16
//			if count <= 0 {
//				return 16
//			} else if sizeKey%16 == 0 {
//				return count * 16
//			} else {
//				return (count + 1) * 16
//			}
//		} else {
//			return sizeChiper
//		}
//	}
//}
