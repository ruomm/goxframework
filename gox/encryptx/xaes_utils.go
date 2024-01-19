/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/17 22:01
 * @version 1.0
 */
package encryptx

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
)

type XAes struct {
	ModeKey     MODE_KEY
	ModeEncode  MODE_ENCODE
	ModePadding MODE_PADDING
	Key         []byte
	Iv          []byte
	//KeyLen          int
	//IvLen           int
	BlockSize       int
	BlockSizeByKey  bool
	PaddingHelper   func(data []byte, blockSize int) []byte
	UnPaddingHelper func(data []byte, blockSize int) []byte
}

// key模式
func (x *XAes) ModeOfKey() MODE_KEY {
	return ParseKeyMode(x.ModeKey)
}

// 字节转字符串编码方案
func (x *XAes) ModeOfEncode() MODE_ENCODE {
	return ParseEncodeMode(x.ModeEncode)
}

// Padding的模式
func (x *XAes) ModeOfPadding() MODE_PADDING {
	return ParsePaddingMode(x.ModePadding)
}

// 设置key字节数组
func (x *XAes) SetKeyData(key []byte) {
	x.Key = key
	//x.KeyLen = len(key)
}

// 设置iv字节数组
func (x *XAes) SetIVData(iv []byte) {
	x.Iv = iv
	//x.IvLen = len(iv)
}

// 设置key字符串
func (x *XAes) SetKeyString(keyStr string) error {
	key, err := x.RestoreKeyIV(keyStr)
	if err == nil {
		x.Key = key
		//x.KeyLen = len(key)
	}
	return err
}

// 获取key字节数组
func (x *XAes) GetKeyData() []byte {
	return x.Key
}

// 设置iv字节数组
func (x *XAes) GetIVData() []byte {
	return x.Iv
}

// 获取key字符串
func (x *XAes) GetKeyString() (string, error) {
	return KeyIVByteToString(x.ModeOfKey(), x.Key)
}

// 获取iv字符串
func (x *XAes) GetIVString() (string, error) {
	return KeyIVByteToString(x.ModeOfKey(), x.Iv)
}

// 设置iv字符串
func (x *XAes) SetIVString(ivStr string) error {
	iv, err := x.RestoreKeyIV(ivStr)
	if err == nil {
		x.Iv = iv
		//x.IvLen = len(iv)
	}
	return err
}

// 设置Blocksize
func (x *XAes) SetBlockSize(blockSize int) {
	if blockSize%8 == 0 && blockSize > 0 {
		x.BlockSize = blockSize
	}
}

// 设置BlockSize依据Key长度自动适配
func (x *XAes) SetBlockSizeByKey(blockSizeByKey bool) {
	x.BlockSizeByKey = blockSizeByKey
}

// 生成key或iv字节数组
func (x *XAes) GenKeyIvData(len int) ([]byte, error) {
	return GenKeyData(len), nil
}

// 生成key或iv字符串
func (x *XAes) GenKeyIvString(len int) (string, error) {
	return GenKeyString(x.ModeOfKey(), len)
}

// 还原key或iv字符串为key或iv字节数组
func (x *XAes) RestoreKeyIV(keyStr string) ([]byte, error) {
	return KeyIVStringToByte(x.ModeOfKey(), keyStr)
}

// ECB加密字节数组
func (x *XAes) EncDataECB(data []byte) ([]byte, error) {
	// 创建新的AES cipher对象
	block, err := aes.NewCipher(x.Key)
	if err != nil {
		return nil, err
	}
	// 补齐明文长度为blockSize字节（AES block size）的倍数
	blockSize := calBlockSize(x.BlockSize, block.BlockSize(), len(x.Key), x.BlockSizeByKey)
	ecbBlockSize := block.BlockSize()
	plaintext := x.Padding(data, blockSize)
	// 加密
	crypted := make([]byte, len(plaintext))
	// 分组分块加密
	for bs := 0; bs < len(plaintext); bs = bs + ecbBlockSize {
		block.Encrypt(crypted[bs:bs+ecbBlockSize], plaintext[bs:bs+ecbBlockSize])
	}
	// IV + Encrypted Data
	return crypted, nil
}

// CBC加密字节数组
func (x *XAes) EncDataCBC(data []byte) ([]byte, error) {
	// 创建新的AES cipher对象
	block, err := aes.NewCipher(x.Key)
	if err != nil {
		return nil, err
	}
	// 补齐明文长度为blockSize字节（AES block size）的倍数
	blockSize := calBlockSize(x.BlockSize, block.BlockSize(), len(x.Key), x.BlockSizeByKey)
	plaintext := x.Padding(data, blockSize)
	// 加密
	blockMode := cipher.NewCBCEncrypter(block, x.Iv)
	crypted := make([]byte, len(plaintext))
	blockMode.CryptBlocks(crypted, plaintext)
	// IV + Encrypted Data
	return crypted, nil
}

// ECB解密字节数组
func (x *XAes) DecDataECB(dataEnc []byte) ([]byte, error) {
	// 创建新的AES cipher对象
	block, err := aes.NewCipher(x.Key)
	if err != nil {
		return nil, err
	}
	blockSize := calBlockSize(x.BlockSize, block.BlockSize(), len(x.Key), x.BlockSizeByKey)
	ecbBlockSize := block.BlockSize()
	// 解密
	decrypted := make([]byte, len(dataEnc))
	block.Decrypt(decrypted, dataEnc)
	// 分组分块解密
	for bs := 0; bs < len(dataEnc); bs = bs + ecbBlockSize {
		block.Decrypt(decrypted[bs:bs+ecbBlockSize], dataEnc[bs:bs+ecbBlockSize])
	}
	unpaddingText := x.UnPadding(decrypted, blockSize)
	// 返回去除padding的数据
	return unpaddingText, nil
}

// CBC解密字节数组
func (x *XAes) DecDataCBC(dataEnc []byte) ([]byte, error) {
	// 创建新的AES cipher对象
	block, err := aes.NewCipher(x.Key)
	if err != nil {
		return nil, err
	}
	blockSize := calBlockSize(x.BlockSize, block.BlockSize(), len(x.Key), x.BlockSizeByKey)
	// 解密
	decrypted := make([]byte, len(dataEnc))
	blockMode := cipher.NewCBCDecrypter(block, x.Iv)
	blockMode.CryptBlocks(decrypted, dataEnc)
	unpaddingText := x.UnPadding(decrypted, blockSize)
	// 返回去除padding的数据
	return unpaddingText, nil
}

// ECB加密字符串
func (x *XAes) EncStringECB(str string) (string, error) {
	crypted, err := x.EncDataECB([]byte(str))
	if err != nil {
		return "", err
	}
	return EncodingToString(x.ModeOfEncode(), crypted)
}

// CBC加密字符串
func (x *XAes) EncStringCBC(str string) (string, error) {
	crypted, err := x.EncDataCBC([]byte(str))
	if err != nil {
		return "", err
	}
	return EncodingToString(x.ModeOfEncode(), crypted)
}

// ECB解密字符串
func (x *XAes) DecStringECB(strEnc string) (string, error) {
	dataEnc, err := DecodingToByte(x.ModeOfEncode(), strEnc)
	if err != nil {
		return "", err
	}
	decrypted, err := x.DecDataECB(dataEnc)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// CBC解密字符串
func (x *XAes) DecStringCBC(strEnc string) (string, error) {
	dataEnc, err := DecodingToByte(x.ModeOfEncode(), strEnc)
	if err != nil {
		return "", err
	}
	decrypted, err := x.DecDataCBC(dataEnc)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// ECB加密字符串
func (x *XAes) EncFileECB(pathSrc string, pathEnc string) error {
	// 打开读写文件
	fiR, errOpenR := os.Open(pathSrc)
	if errOpenR != nil {
		fmt.Println("open file error: ", errOpenR)
		return errOpenR
	}
	defer fiR.Close()
	fiW, errOpenW := os.OpenFile(pathEnc, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if errOpenW != nil {
		fmt.Println("open file error: ", errOpenW)
		return errOpenW
	}
	defer fiW.Close()
	// 创建新的AES cipher对象
	block, err := aes.NewCipher(x.Key)
	if err != nil {
		return err
	}
	// 补齐明文长度为blockSize字节（AES block size）的倍数
	blockSize := calBlockSize(x.BlockSize, block.BlockSize(), len(x.Key), x.BlockSizeByKey)
	ecbBlockSize := block.BlockSize()
	bufferSize := (FILE_BUFFER_SIZE / blockSize) * blockSize
	if bufferSize < 2*blockSize {
		bufferSize = 2 * blockSize
	}
	// 开始文件加密
	reader := bufio.NewReader(fiR)
	writer := bufio.NewWriter(fiW)
	var bRead = make([]byte, bufferSize)
	//hasRead := false
	for {
		nR, errR := reader.Read(bRead)
		if errR != nil && errR != io.EOF {
			return errR
		}
		if nR == 0 {
			plaintext := x.Padding([]byte{}, blockSize)
			// 加密
			crypted := make([]byte, len(plaintext))
			// 分组分块加密
			for bs := 0; bs < len(plaintext); bs = bs + ecbBlockSize {
				block.Encrypt(crypted[bs:bs+ecbBlockSize], plaintext[bs:bs+ecbBlockSize])
			}
			_, errW := writer.Write(crypted)
			if errW != nil {
				return errW
			}
			break
		}
		if nR < bufferSize {
			plaintext := x.Padding(bRead[0:nR], blockSize)
			// 加密
			crypted := make([]byte, len(plaintext))
			// 分组分块加密
			for bs := 0; bs < len(plaintext); bs = bs + ecbBlockSize {
				block.Encrypt(crypted[bs:bs+ecbBlockSize], plaintext[bs:bs+ecbBlockSize])
			}
			_, errW := writer.Write(crypted)
			if errW != nil {
				return errW
			}
			break
		} else {
			// 加密
			crypted := make([]byte, bufferSize)
			// 分组分块加密
			for bs := 0; bs < len(bRead); bs = bs + ecbBlockSize {
				block.Encrypt(crypted[bs:bs+ecbBlockSize], bRead[bs:bs+ecbBlockSize])
			}
			_, errW := writer.Write(crypted)
			if errW != nil {
				return errW
			}
		}

	}
	// 清空缓冲区，进行缓冲区内容落到磁盘
	errFlush := writer.Flush()
	if errFlush != nil {
		return errFlush
	}
	return nil
}

// CBC加密字符串
func (x *XAes) EncFileCBC(pathSrc string, pathEnc string) error {
	// 打开读写文件
	fiR, errOpenR := os.Open(pathSrc)
	if errOpenR != nil {
		fmt.Println("open file error: ", errOpenR)
		return errOpenR
	}
	defer fiR.Close()
	fiW, errOpenW := os.OpenFile(pathEnc, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if errOpenW != nil {
		fmt.Println("open file error: ", errOpenW)
		return errOpenW
	}
	defer fiW.Close()
	// 创建新的AES cipher对象
	block, err := aes.NewCipher(x.Key)
	if err != nil {
		return err
	}
	// 补齐明文长度为blockSize字节（AES block size）的倍数
	blockSize := calBlockSize(x.BlockSize, block.BlockSize(), len(x.Key), x.BlockSizeByKey)
	blockMode := cipher.NewCBCEncrypter(block, x.Iv)
	bufferSize := (FILE_BUFFER_SIZE / blockSize) * blockSize
	if bufferSize < 2*blockSize {
		bufferSize = 2 * blockSize
	}
	// 开始文件加密
	reader := bufio.NewReader(fiR)
	writer := bufio.NewWriter(fiW)
	var bRead = make([]byte, bufferSize)
	//hasRead := false
	for {
		nR, errR := reader.Read(bRead)
		if errR != nil && errR != io.EOF {
			return errR
		}
		if nR == 0 {
			plaintext := x.Padding([]byte{}, blockSize)
			// 加密
			crypted := make([]byte, len(plaintext))
			blockMode.CryptBlocks(crypted, plaintext)
			_, errW := writer.Write(crypted)
			if errW != nil {
				return errW
			}
			break
		}
		if nR < bufferSize {
			plaintext := x.Padding(bRead[0:nR], blockSize)
			// 加密
			crypted := make([]byte, len(plaintext))
			blockMode.CryptBlocks(crypted, plaintext)
			_, errW := writer.Write(crypted)
			if errW != nil {
				return errW
			}
			break
		} else {
			// 加密
			crypted := make([]byte, bufferSize)
			blockMode.CryptBlocks(crypted, bRead)
			_, errW := writer.Write(crypted)
			if errW != nil {
				return errW
			}
		}

	}
	// 清空缓冲区，进行缓冲区内容落到磁盘
	errFlush := writer.Flush()
	if errFlush != nil {
		return errFlush
	}
	return nil
}

// ECB解密字符串
func (x *XAes) DecFileECB(pathEnc string, pathDest string) error {
	// 打开读写文件
	fiR, errOpenR := os.Open(pathEnc)
	if errOpenR != nil {
		fmt.Println("open file error: ", errOpenR)
		return errOpenR
	}
	defer fiR.Close()
	fiW, errOpenW := os.OpenFile(pathDest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if errOpenW != nil {
		fmt.Println("open file error: ", errOpenW)
		return errOpenW
	}
	defer fiW.Close()
	// 创建新的AES cipher对象
	block, err := aes.NewCipher(x.Key)
	if err != nil {
		return err
	}
	// 补齐明文长度为blockSize字节（AES block size）的倍数
	blockSize := calBlockSize(x.BlockSize, block.BlockSize(), len(x.Key), x.BlockSizeByKey)
	ecbBlockSize := block.BlockSize()
	bufferSize := (FILE_BUFFER_SIZE / blockSize) * blockSize
	if bufferSize < 2*blockSize {
		bufferSize = 2 * blockSize
	}
	// 开始文件解密
	reader := bufio.NewReader(fiR)
	writer := bufio.NewWriter(fiW)
	var bRead = make([]byte, bufferSize)
	var preBuffer []byte = nil
	for {
		nR, errR := reader.Read(bRead)
		if errR != nil && errR != io.EOF {
			return errR
		}
		if nR > 0 {
			if nil != preBuffer && len(preBuffer) > 0 {
				//解密上一次数据并写入
				// 解密
				decrypted := make([]byte, len(preBuffer))
				// 分组分块解密
				for bs := 0; bs < len(preBuffer); bs = bs + ecbBlockSize {
					block.Decrypt(decrypted[bs:bs+ecbBlockSize], preBuffer[bs:bs+ecbBlockSize])
				}
				_, errW := writer.Write(decrypted)
				if errW != nil {
					return errW
				}
			}
			preBuffer = make([]byte, nR)
			copy(preBuffer, bRead[0:nR])
		}
		if nR < bufferSize {
			break
		}

	}
	if nil != preBuffer && len(preBuffer) > 0 {
		//解密上一次数据并写入
		// 解密
		decrypted := make([]byte, len(preBuffer))
		// 分组分块解密
		for bs := 0; bs < len(preBuffer); bs = bs + ecbBlockSize {
			block.Decrypt(decrypted[bs:bs+ecbBlockSize], preBuffer[bs:bs+ecbBlockSize])
		}
		plaintText := x.UnPadding(decrypted, blockSize)
		if len(plaintText) > 0 {
			_, errW := writer.Write(plaintText)
			if errW != nil {
				return errW
			}
		}

	}
	// 清空缓冲区，进行缓冲区内容落到磁盘
	errFlush := writer.Flush()
	if errFlush != nil {
		return errFlush
	}
	return nil
}

// CBC解密字符串
func (x *XAes) DecFileCBC(pathEnc string, pathDest string) error {
	// 打开读写文件
	fiR, errOpenR := os.Open(pathEnc)
	if errOpenR != nil {
		fmt.Println("open file error: ", errOpenR)
		return errOpenR
	}
	defer fiR.Close()
	fiW, errOpenW := os.OpenFile(pathDest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if errOpenW != nil {
		fmt.Println("open file error: ", errOpenW)
		return errOpenW
	}
	defer fiW.Close()
	// 创建新的AES cipher对象
	block, err := aes.NewCipher(x.Key)
	if err != nil {
		return err
	}
	// 补齐明文长度为blockSize字节（AES block size）的倍数
	blockSize := calBlockSize(x.BlockSize, block.BlockSize(), len(x.Key), x.BlockSizeByKey)
	blockMode := cipher.NewCBCDecrypter(block, x.Iv)
	bufferSize := (FILE_BUFFER_SIZE / blockSize) * blockSize
	if bufferSize < 2*blockSize {
		bufferSize = 2 * blockSize
	}
	// 开始文件解密
	reader := bufio.NewReader(fiR)
	writer := bufio.NewWriter(fiW)
	var bRead = make([]byte, bufferSize)
	var preBuffer []byte = nil
	for {
		nR, errR := reader.Read(bRead)
		if errR != nil && errR != io.EOF {
			return errR
		}
		if nR > 0 {
			if nil != preBuffer && len(preBuffer) > 0 {
				//解密上一次数据并写入
				// 解密
				decrypted := make([]byte, len(preBuffer))
				// 分组分块解密
				blockMode.CryptBlocks(decrypted, preBuffer)
				_, errW := writer.Write(decrypted)
				if errW != nil {
					return errW
				}
			}
			preBuffer = make([]byte, nR)
			copy(preBuffer, bRead[0:nR])
		}
		if nR < bufferSize {
			break
		}

	}
	if nil != preBuffer && len(preBuffer) > 0 {
		//解密上一次数据并写入
		// 解密
		decrypted := make([]byte, len(preBuffer))
		blockMode.CryptBlocks(decrypted, preBuffer)
		plaintText := x.UnPadding(decrypted, blockSize)
		if len(plaintText) > 0 {
			_, errW := writer.Write(plaintText)
			if errW != nil {
				return errW
			}
		}

	}
	// 清空缓冲区，进行缓冲区内容落到磁盘
	errFlush := writer.Flush()
	if errFlush != nil {
		return errFlush
	}
	return nil
}

// Padding&UnPadding
func (x *XAes) Padding(data []byte, blockSize int) []byte {
	if nil != x.PaddingHelper {
		return x.Padding(data, blockSize)
	}
	modeOfPadding := x.ModeOfPadding()
	if modeOfPadding == MODE_PADDING_PKCS7 {
		return Pkcs7Padding(data, blockSize)
	} else if modeOfPadding == MODE_PADDING_PKCS5 {
		return Pkcs5Padding(data, blockSize)
	} else {
		return Pkcs7Padding(data, blockSize)
	}
}
func (x *XAes) UnPadding(data []byte, blockSize int) []byte {
	if len(data) <= 0 {
		return data
	}
	if nil != x.UnPaddingHelper {
		return x.UnPaddingHelper(data, blockSize)
	}
	modeOfPadding := x.ModeOfPadding()
	if modeOfPadding == MODE_PADDING_PKCS7 {
		return Pkcs7UnPadding(data)
	} else if modeOfPadding == MODE_PADDING_PKCS5 {
		return Pkcs5UnPadding(data)
	} else {
		return Pkcs7UnPadding(data)
	}
}
