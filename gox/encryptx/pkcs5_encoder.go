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

func Pkcs5Padding(data []byte, blockSize int) []byte {
	//填充
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

/**
 * 删除解密后明文的补位字符
 *
 * @param decrypted 解密后的明文
 * @return 删除补位字符后的明文
 */
func Pkcs5UnPadding(decrypted []byte) []byte {
	length := len(decrypted)
	unpadding := int(decrypted[length-1])
	if unpadding < 0 || unpadding > length {
		unpadding = 0
	}
	return decrypted[:(length - unpadding)]
}
