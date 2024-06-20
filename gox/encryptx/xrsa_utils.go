/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/6/20 09:50
 * @version 1.0
 */
package encryptx

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

type XRsa struct {
	//ModeKey         MODE_KEY
	ModeEncode      MODE_ENCODE
	ModePadding     MODE_PADDING
	PublicKey       *rsa.PublicKey
	PrivateKey      *rsa.PrivateKey
	PaddingHelper   func(data []byte, blockSize int) []byte
	UnPaddingHelper func(data []byte, blockSize int) []byte
}

// key模式
//func (x *XRsa) ModeOfKey() MODE_KEY {
//	return ParseKeyMode(x.ModeKey)
//}

// 生成秘钥对
func (x *XRsa) GenrateKeyPair(bits int) error {
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	x.PrivateKey = priKey
	x.PublicKey = &priKey.PublicKey
	return nil
}

// 字节转字符串编码方案
func (x *XRsa) ModeOfEncode() MODE_ENCODE {
	return ParseEncodeMode(x.ModeEncode)
}

// Padding的模式
func (x *XRsa) ModeOfPadding() MODE_PADDING {
	return ParsePaddingMode(x.ModePadding)
}

// 设置公钥
func (x *XRsa) SetPubicKey(pubKey []byte) error {
	parsePub, err := x509.ParsePKCS1PublicKey(pubKey)
	if err != nil {
		return err
	}
	x.PublicKey = parsePub
	return nil
}

// 设置私钥
func (x *XRsa) SetPrivateKey(priKey []byte) error {
	parsePri, err := x509.ParsePKCS1PrivateKey(priKey)
	if err != nil {
		return err
	}
	x.PrivateKey = parsePri
	return nil
}

// 加载公钥
func (x *XRsa) LoadPulicKey(modeOfKey MODE_KEY, pubKeyStr string) error {
	keyByte, err := RsaKeyStringToByte(ParseKeyMode(modeOfKey), pubKeyStr)
	if err != nil {
		return err
	}
	parsePub, err := x509.ParsePKCS1PublicKey(keyByte)
	if err != nil {
		return err
	}
	x.PublicKey = parsePub
	return nil
}

// 加载私钥
func (x *XRsa) LoadPrivateKey(modeOfKey MODE_KEY, priKeyStr string) error {
	keyByte, err := RsaKeyStringToByte(ParseKeyMode(modeOfKey), priKeyStr)
	if err != nil {
		return err
	}
	parsePri, err := x509.ParsePKCS1PrivateKey(keyByte)
	if err != nil {
		return err
	}
	x.PrivateKey = parsePri
	return nil
}

// 格式化公钥
func (x *XRsa) FormatPublicKey(modeOfKey MODE_KEY) (string, error) {
	keyData := x509.MarshalPKCS1PublicKey(x.PublicKey)
	return RsaKeyByteToString(ParseKeyMode(modeOfKey), keyData, true)
}

// 格式化私钥
func (x *XRsa) FormatPrivateKey(modeOfKey MODE_KEY) (string, error) {
	keyData := x509.MarshalPKCS1PrivateKey(x.PrivateKey)
	return RsaKeyByteToString(ParseKeyMode(modeOfKey), keyData, false)
}
