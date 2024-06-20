/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/6/20 09:47
 * @version 1.0
 */
package encryptx

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"github.com/ruomm/goxframework/gox/corex"
	"strings"
)

type RsaHelper interface {
	// 生成秘钥对
	GenrateKeyPair(bits int) error
	// 字节转字符串编码方案
	ModeOfEncode() MODE_ENCODE
	// Padding的模式
	ModeOfPadding() MODE_PADDING
	// 设置公钥
	SetPubicKey(pubKey []byte) error
	// 设置私钥
	SetPrivateKey(priKey []byte) error
	// 加载公钥
	LoadPulicKey(modeOfKey MODE_KEY, pubKeyStr string) error
	// 加载私钥
	LoadPrivateKey(modeOfKey MODE_KEY, priKeyStr string) error
	// 格式化公钥
	FormatPublicKey(modeOfKey MODE_KEY) (string, error)
	// 格式化私钥
	FormatPrivateKey(modeOfKey MODE_KEY) (string, error)
	// 公钥长度
	SizeOfPublicKey() int
	// 秘钥长度
	SizeOfPrivateKey() int
	// 使用公钥进行PKCS1v15加密，待加密信息长度不能超过秘钥模长-11
	EncryptPKCS1v15(origMsg []byte) ([]byte, error)
	// 使用私钥进行PKCS1v15解密，解密后信息长度不超过秘钥模长-11
	DecryptPKCS1v15(encMsg []byte) ([]byte, error)
	// 使用公钥进行PKCS1v15加密，待加密信息长度超过秘钥模长-11则使用分段加密
	EncryptPKCS1v15Big(origMsg []byte) ([]byte, error)
	// 使用私钥进行PKCS1v15解密，待解密后信息长度超过秘钥模长则使用分段解密
	DecryptPKCS1v15Big(encMsg []byte) ([]byte, error)
	// 使用公钥进行PKCS1v15加密，待加密信息长度不能超过秘钥模长-11
	EncryptPKCS1v15String(encodeMode MODE_ENCODE, origStr string) (string, error)
	// 使用私钥进行PKCS1v15解密，解密后信息长度不超过秘钥模长-11
	DecryptPKCS1v15String(encodeMode MODE_ENCODE, encStr string) (string, error)
	// 使用公钥进行PKCS1v15加密，待加密信息长度超过秘钥模长-11则使用分段加密
	EncryptPKCS1v15StringBig(encodeMode MODE_ENCODE, origStr string) (string, error)
	// 使用私钥进行PKCS1v15解密，待解密后信息长度超过秘钥模长则使用分段解密
	DecryptPKCS1v15StringBig(encodeMode MODE_ENCODE, encStr string) (string, error)
	// 使用私钥进行签名-字节数组模式
	SignPSS(hash crypto.Hash, origData []byte, opts *rsa.PSSOptions) ([]byte, error)
	// 使用公钥验证签名-字节数组模式
	VerifyPSS(hash crypto.Hash, origData []byte, sig []byte, opts *rsa.PSSOptions) error
	// 使用私钥进行签名-字符串模式
	SignPSSByString(encodeMode MODE_ENCODE, hash crypto.Hash, origStr string, opts *rsa.PSSOptions) (string, error)
	// 使用公钥验证签名-字符串模式
	VerifyPSSByString(encodeMode MODE_ENCODE, hash crypto.Hash, origStr string, sigStr string, opts *rsa.PSSOptions) error
}

// 格式化密钥
func FormatKeyByData(keyData []byte, tag string) (string, error) {
	if len(keyData) == 0 {
		return "", errors.New("format key error,key byte is empty")
	}
	keyStr := base64.StdEncoding.EncodeToString(keyData)
	return FormatKeyByString(keyStr, tag)
}

// 格式化密钥
func FormatKeyByString(keyStr string, tag string) (string, error) {
	if len(keyStr) <= 0 {
		return "", errors.New("format key error,key string is empty")
	}
	build := strings.Builder{}
	build.WriteString("-----BEGIN " + tag + "-----")
	build.WriteString("\r\n")
	lengthKey := len(keyStr)
	size := lengthKey / 64
	for i := 0; i < size; i++ {
		build.WriteString(keyStr[i*64 : (i+1)*64])
		build.WriteString("\r\n")
	}
	if lengthKey%64 != 0 {
		build.WriteString(keyStr[size*64 : lengthKey])
		build.WriteString("\r\n")
	}
	build.WriteString("-----END " + tag + "-----")
	return build.String(), nil
}

// 读取格式化密钥
func ReadFormatKey(str string) ([]byte, error) {
	strSlice := corex.StringReadByLineNoBom(str, false)
	if strSlice == nil || len(strSlice) == 0 {
		return nil, errors.New("read format key error,key string is empty")
	}
	build := strings.Builder{}
	for _, line := range strSlice {
		if strings.HasPrefix(line, "----") {
			if build.Len() > 0 {
				break
			}
		} else {
			build.WriteString(line)
		}
	}
	if build.Len() <= 0 {
		return nil, errors.New("read format key error,key string is empty")
	}
	return base64.StdEncoding.DecodeString(build.String())
}

func RsaKeyStringToByte(keyMode MODE_KEY, keyStr string) ([]byte, error) {
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
		return ReadFormatKey(keyStr)
	} else if keyMode == MODE_KEY_PEM {
		block, _ := pem.Decode([]byte(keyStr))
		if block == nil || len(block.Bytes) == 0 {
			return nil, errors.New("decode pem error!")
		}
		return block.Bytes, nil
	} else {
		return ReadFormatKey(keyStr)
	}
}
func RsaKeyByteToString(keyMode MODE_KEY, keyData []byte, public bool) (string, error) {
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
		tag := "RSA PRIVATE KEY"
		if public {
			tag = "RSA PUBLIC KEY"
		}
		return FormatKeyByData(keyData, tag)
	} else if keyMode == MODE_KEY_PEM {
		tag := "RSA PRIVATE KEY"
		if public {
			tag = "RSA PUBLIC KEY"
		}
		pemBlock := &pem.Block{
			Type:  tag,
			Bytes: keyData,
		}
		return string(pem.EncodeToMemory(pemBlock)), nil
	} else {
		tag := "RSA PRIVATE KEY"
		if public {
			tag = "RSA PUBLIC KEY"
		}
		return FormatKeyByData(keyData, tag)
	}
}
