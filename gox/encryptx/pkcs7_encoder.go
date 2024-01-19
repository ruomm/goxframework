/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/1/17 22:39
 * @version 1.0
 */
package encryptx

import "bytes"

/**
 * 获得对明文进行补位填充的字节.
 *
 * @param count 需要进行填充补位操作的明文字节个数
 * @return 补齐用的字节数组
 */

func Pkcs7Padding(data []byte, blockSize int) []byte {
	count := len(data)
	// 计算需要填充的位数
	amountToPad := blockSize - (count % blockSize)
	if amountToPad == 0 {
		amountToPad = blockSize
	}
	// 获得补位所用的字符
	padChr := Pkcs7Byte(amountToPad)
	var buffer bytes.Buffer
	buffer.Write(data)
	for index := 0; index < amountToPad; index++ {
		buffer.WriteByte(padChr)
	}
	return buffer.Bytes()
}

//func encodepadding(count int, blockSize int) []byte {
//	// 计算需要填充的位数
//	amountToPad := blockSize - (count % blockSize)
//	if amountToPad == 0 {
//		amountToPad = blockSize
//	}
//	// 获得补位所用的字符
//	padChr := Pkcs7Byte(amountToPad)
//	var buffer bytes.Buffer
//	for index := 0; index < amountToPad; index++ {
//		buffer.WriteByte(padChr)
//	}
//	return buffer.Bytes()
//}

/**
 * 删除解密后明文的补位字符
 *
 * @param decrypted 解密后的明文
 * @return 删除补位字符后的明文
 */
func Pkcs7UnPadding(decrypted []byte) []byte {
	pad := int(decrypted[len(decrypted)-1])
	return decrypted[0 : len(decrypted)-pad]
}

/**
 * 将数字转化成ASCII码对应的字符，用于对明文进行补码
 *
 * @param a 需要转化的数字
 * @return 转化得到的字符
 */
func Pkcs7Byte(a int) byte {
	return (byte)(a & 0xFF)
}
