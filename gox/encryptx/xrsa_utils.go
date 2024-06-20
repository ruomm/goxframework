/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/6/20 09:50
 * @version 1.0
 */
package encryptx

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"os"
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

// 公钥长度
func (x *XRsa) SizeOfPublicKey() int {
	if nil == x.PublicKey {
		return 0
	} else {
		return x.PublicKey.Size()
	}
}

// 秘钥长度
func (x *XRsa) SizeOfPrivateKey() int {
	if nil == x.PrivateKey {
		return 0
	} else {
		return x.PrivateKey.Size()
	}
}

// 使用公钥进行PKCS1v15加密，待加密信息长度不能超过秘钥模长-11
func (x *XRsa) EncryptPKCS1v15(origMsg []byte) ([]byte, error) {
	if nil == x.PublicKey {
		return nil, errors.New("XRsa.PublicKey is nil")
	}
	if nil == origMsg {
		return nil, errors.New("origMsg is nil")
	}
	return rsa.EncryptPKCS1v15(rand.Reader, x.PublicKey, origMsg)
}

// 使用私钥进行PKCS1v15解密，解密后信息长度不超过秘钥模长-11
func (x *XRsa) DecryptPKCS1v15(encMsg []byte) ([]byte, error) {
	if nil == x.PrivateKey {
		return nil, errors.New("XRsa.PrivateKey is nil")
	}
	if nil == encMsg {
		return nil, errors.New("encMsg is nil")
	}
	return rsa.DecryptPKCS1v15(rand.Reader, x.PrivateKey, encMsg)
}

// 使用公钥进行PKCS1v15加密，待加密信息长度超过秘钥模长-11则使用分段加密
func (x *XRsa) EncryptPKCS1v15Big(origMsg []byte) ([]byte, error) {
	if nil == x.PublicKey {
		return nil, errors.New("XRsa.PublicKey is nil")
	}
	if nil == origMsg {
		return nil, errors.New("origMsg is nil")
	}
	keySize := x.SizeOfPublicKey()
	blockSize := keySize - 11
	if blockSize <= 0 {
		return nil, errors.New("XRsa.PublicKey blockSize is too small")
	}
	lenData := len(origMsg)
	if lenData <= blockSize {
		return rsa.EncryptPKCS1v15(rand.Reader, x.PublicKey, origMsg)
	}
	buffer := bytes.Buffer{}
	size := lenData / blockSize
	for i := 0; i < size; i++ {
		tmpData := make([]byte, blockSize)
		copy(tmpData, origMsg[i*blockSize:(i+1)*blockSize])
		tmpEnc, tmpErr := rsa.EncryptPKCS1v15(rand.Reader, x.PublicKey, tmpData)
		if nil != tmpErr {
			return nil, tmpErr
		}
		buffer.Write(tmpEnc)
	}
	if lenData%blockSize != 0 {
		tmpData := make([]byte, lenData-size*blockSize)
		copy(tmpData, origMsg[size*blockSize:lenData])
		tmpEnc, tmpErr := rsa.EncryptPKCS1v15(rand.Reader, x.PublicKey, tmpData)
		if nil != tmpErr {
			return nil, tmpErr
		}
		buffer.Write(tmpEnc)
	}
	return buffer.Bytes(), nil
}

// 使用私钥进行PKCS1v15解密，待解密后信息长度超过秘钥模长则使用分段解密
func (x *XRsa) DecryptPKCS1v15Big(encMsg []byte) ([]byte, error) {
	if nil == x.PrivateKey {
		return nil, errors.New("XRsa.PrivateKey is nil")
	}
	if nil == encMsg {
		return nil, errors.New("encMsg is nil")
	}
	keySize := x.SizeOfPrivateKey()
	blockSize := keySize
	if blockSize <= 0 {
		return nil, errors.New("XRsa.PrivateKey blockSize is too small")
	}
	lenData := len(encMsg)
	if lenData <= blockSize {
		return rsa.DecryptPKCS1v15(rand.Reader, x.PrivateKey, encMsg)
	}
	if lenData%blockSize != 0 {
		return nil, errors.New("lenData is not a multiple of the block size")
	}
	buffer := bytes.Buffer{}
	size := lenData / blockSize
	for i := 0; i < size; i++ {
		tmpData := make([]byte, blockSize)
		copy(tmpData, encMsg[i*blockSize:(i+1)*blockSize])
		tmpDec, tmpErr := rsa.DecryptPKCS1v15(rand.Reader, x.PrivateKey, tmpData)
		if nil != tmpErr {
			return nil, tmpErr
		}
		buffer.Write(tmpDec)
	}
	return buffer.Bytes(), nil
}

// 使用公钥进行PKCS1v15加密，待加密信息长度不能超过秘钥模长-11
func (x *XRsa) EncryptPKCS1v15String(encodeMode MODE_ENCODE, origStr string) (string, error) {
	encMsg, err := x.EncryptPKCS1v15([]byte(origStr))
	if nil != err {
		return "", err
	}
	return EncodingToString(ParseEncodeMode(encodeMode), encMsg)
}

// 使用私钥进行PKCS1v15解密，解密后信息长度不超过秘钥模长-11
func (x *XRsa) DecryptPKCS1v15String(encodeMode MODE_ENCODE, encStr string) (string, error) {
	encMsg, err := DecodingToByte(ParseEncodeMode(encodeMode), encStr)
	if nil != err {
		return "", err
	}
	decMsg, err := x.DecryptPKCS1v15(encMsg)
	if nil != err {
		return "", err
	}
	return string(decMsg), nil
}

// 使用公钥进行PKCS1v15加密，待加密信息长度超过秘钥模长-11则使用分段加密
func (x *XRsa) EncryptPKCS1v15StringBig(encodeMode MODE_ENCODE, origStr string) (string, error) {
	encMsg, err := x.EncryptPKCS1v15Big([]byte(origStr))
	if nil != err {
		return "", err
	}
	return EncodingToString(ParseEncodeMode(encodeMode), encMsg)
}

// 使用私钥进行PKCS1v15解密，待解密后信息长度超过秘钥模长则使用分段解密
func (x *XRsa) DecryptPKCS1v15StringBig(encodeMode MODE_ENCODE, encStr string) (string, error) {
	encMsg, err := DecodingToByte(ParseEncodeMode(encodeMode), encStr)
	if nil != err {
		return "", err
	}
	decMsg, err := x.DecryptPKCS1v15Big(encMsg)
	if nil != err {
		return "", err
	}
	return string(decMsg), nil
}

// 使用私钥进行签名-字节数组模式
func (x *XRsa) SignPSS(hash crypto.Hash, origData []byte, opts *rsa.PSSOptions) ([]byte, error) {
	if nil == x.PrivateKey {
		return nil, errors.New("XRsa.PrivateKey is nil")
	}
	if nil == origData {
		return nil, errors.New("SignPSS err,origData is nil")
	}
	h := hash.New()
	h.Write(origData)
	digest := h.Sum(nil)
	return rsa.SignPSS(rand.Reader, x.PrivateKey, hash, digest, opts)
}

// 使用公钥验证签名-字节数组模式
func (x *XRsa) VerifyPSS(hash crypto.Hash, origData []byte, sig []byte, opts *rsa.PSSOptions) error {
	if nil == x.PrivateKey {
		return errors.New("XRsa.PrivateKey is nil")
	}
	if nil == origData {
		return errors.New("VerifyPSS err,origData is nil")
	}
	if nil == sig {
		return errors.New("VerifyPSS err,sig is nil")
	}
	h := hash.New()
	h.Write(origData)
	digest := h.Sum(nil)
	return rsa.VerifyPSS(x.PublicKey, hash, digest, sig, opts)
}

// 使用私钥进行签名-字符串模式
func (x *XRsa) SignPSSByString(encodeMode MODE_ENCODE, hash crypto.Hash, origStr string, opts *rsa.PSSOptions) (string, error) {
	if nil == x.PrivateKey {
		return "", errors.New("XRsa.PrivateKey is nil")
	}
	h := hash.New()
	h.Write([]byte(origStr))
	digest := h.Sum(nil)
	sig, err := rsa.SignPSS(rand.Reader, x.PrivateKey, hash, digest, opts)
	if err != nil {
		return "", err
	}
	return EncodingToString(ParseEncodeMode(encodeMode), sig)
}

// 使用公钥验证签名-字符串模式
func (x *XRsa) VerifyPSSByString(encodeMode MODE_ENCODE, hash crypto.Hash, origStr string, sigStr string, opts *rsa.PSSOptions) error {
	if nil == x.PrivateKey {
		return errors.New("XRsa.PrivateKey is nil")
	}
	h := hash.New()
	h.Write([]byte(origStr))
	digest := h.Sum(nil)
	sig, err := DecodingToByte(encodeMode, sigStr)
	if err != nil {
		return err
	}
	return rsa.VerifyPSS(x.PublicKey, hash, digest, sig, opts)
}

// 使用公钥进行PKCS1v15文件加密
func (x *XRsa) EncryptPKCS1v15File(origFile string, encFile string, emptyEncrypt bool) error {
	if nil == x.PublicKey {
		return errors.New("XRsa.PublicKey is nil")
	}
	if len(origFile) <= 0 {
		return errors.New("origFile path is empty")
	}
	if len(encFile) <= 0 {
		return errors.New("encFile path is empty")
	}
	keySize := x.SizeOfPublicKey()
	blockSize := keySize - 11
	if blockSize <= 0 {
		return errors.New("XRsa.PublicKey blockSize is too small")
	}
	// 打开读写文件
	fiR, errOpenR := os.Open(origFile)
	if errOpenR != nil {
		fmt.Println("open file error: ", errOpenR)
		return errOpenR
	}
	defer fiR.Close()
	fiW, errOpenW := os.OpenFile(encFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if errOpenW != nil {
		fmt.Println("open file error: ", errOpenW)
		return errOpenW
	}
	defer fiW.Close()
	// 开始文件加密
	reader := bufio.NewReader(fiR)
	writer := bufio.NewWriter(fiW)
	var bRead = make([]byte, blockSize)
	emptyEcnFlag := false
	if emptyEncrypt {
		emptyEcnFlag = true
	}
	for {
		nR, errR := reader.Read(bRead)
		if errR != nil && errR != io.EOF {
			return errR
		}
		if nR == 0 {
			if emptyEcnFlag {
				emptyEcnFlag = false
				// 加密
				encMsg, errSub := x.EncryptPKCS1v15(make([]byte, 0))
				if errSub != nil {
					return errSub
				}
				_, errW := writer.Write(encMsg)
				if errW != nil {
					return errW
				}
			} else {
				emptyEcnFlag = false
			}
			break

		} else if nR > 0 {
			emptyEcnFlag = false
			// 加密
			encMsg, errSub := x.EncryptPKCS1v15(bRead[0:nR])
			if errSub != nil {
				return errSub
			}
			if len(encMsg) <= 0 {
				continue
			}
			_, errW := writer.Write(encMsg)
			if errW != nil {
				return errW
			}
			if nR < blockSize {
				break
			}
		} else {
			emptyEcnFlag = false
			break
		}
	}
	// 清空缓冲区，进行缓冲区内容落到磁盘
	errFlush := writer.Flush()
	if errFlush != nil {
		return errFlush
	}
	return nil
}

// 使用私钥进行PKCS1v15文件解密
func (x *XRsa) DecryptPKCS1v15File(encFile string, decFile string) error {
	if nil == x.PrivateKey {
		return errors.New("XRsa.PrivateKey is nil")
	}
	if len(encFile) <= 0 {
		return errors.New("encFile path is empty")
	}
	if len(decFile) <= 0 {
		return errors.New("decFile path is empty")
	}
	keySize := x.SizeOfPrivateKey()
	blockSize := keySize
	if blockSize <= 0 {
		return errors.New("XRsa.PrivateKey blockSize is too small")
	}
	// 打开读写文件
	fiR, errOpenR := os.Open(encFile)
	if errOpenR != nil {
		fmt.Println("open file error: ", errOpenR)
		return errOpenR
	}
	defer fiR.Close()
	fiW, errOpenW := os.OpenFile(decFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if errOpenW != nil {
		fmt.Println("open file error: ", errOpenW)
		return errOpenW
	}
	defer fiW.Close()

	// 开始文件解密
	reader := bufio.NewReader(fiR)
	writer := bufio.NewWriter(fiW)
	var bRead = make([]byte, blockSize)
	for {
		nR, errR := reader.Read(bRead)
		if errR != nil && errR != io.EOF {
			return errR
		}
		if nR == blockSize {
			// 解密
			decMsg, errSub := x.DecryptPKCS1v15(bRead[0:nR])
			if errSub != nil {
				return errSub
			}
			if len(decMsg) <= 0 {
				continue
			}
			_, errW := writer.Write(decMsg)
			if errW != nil {
				return errW
			}
		} else {
			break
		}
	}
	// 清空缓冲区，进行缓冲区内容落到磁盘
	errFlush := writer.Flush()
	if errFlush != nil {
		return errFlush
	}
	return nil
}

//// 使用私钥进行签名-字节数组模式
//func (x *XRsa) SignPSSFile(hash crypto.Hash, origFile string, opts *rsa.PSSOptions) ([]byte, error) {
//	if nil == x.PrivateKey {
//		return nil, errors.New("XRsa.PrivateKey is nil")
//	}
//	if len(origFile) <= 0 {
//		return nil, errors.New("origFile path is empty")
//	}
//	// 打开读写文件
//	fiR, errOpenR := os.Open(origFile)
//	if errOpenR != nil {
//		fmt.Println("open file error: ", errOpenR)
//		return nil, errOpenR
//	}
//	defer fiR.Close()
//	// 开始文件解密
//	reader := bufio.NewReader(fiR)
//	h := hash.New()
//	h.Write()
//	h.Write(origData)
//	digest := h.Sum(nil)
//	return rsa.SignPSS(rand.Reader, x.PrivateKey, hash, digest, opts)
//}
//
//// 使用公钥验证签名-字节数组模式
//func (x *XRsa) VerifyPSSFile(hash crypto.Hash, origFile string, sig []byte, opts *rsa.PSSOptions) error {
//	if nil == x.PrivateKey {
//		return errors.New("XRsa.PrivateKey is nil")
//	}
//	if len(origFile) <= 0 {
//		return errors.New("origFile path is empty")
//	}
//	if nil == sig {
//		return errors.New("VerifyPSS err,sig is nil")
//	}
//	h := hash.New()
//	h.Write(origData)
//	digest := h.Sum(nil)
//	return rsa.VerifyPSS(x.PublicKey, hash, digest, sig, opts)
//}
//
//func (x *XRsa) VerifyPSS(hash crypto.Hash, origFile string, sig []byte) error {}
