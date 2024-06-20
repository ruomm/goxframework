/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/6/20 17:52
 * @version 1.0
 */
package encryptx

import (
	"bufio"
	"crypto"
	"errors"
	"fmt"
	"io"
	"os"
)

func Sum(hash crypto.Hash, sha []byte) ([]byte, error) {
	if nil == sha {
		return nil, errors.New("invalid sha")
	}
	h := hash.New()
	_, errSum := h.Write(sha)
	if errSum != nil {
		return nil, errSum
	}
	sum := h.Sum(nil)
	return sum, nil
}

func SumByString(encodeMode MODE_ENCODE, hash crypto.Hash, str string) (string, error) {
	h := hash.New()
	_, errSum := h.Write([]byte(str))
	if errSum != nil {
		return "", errSum
	}
	sum := h.Sum(nil)
	return EncodingToString(ParseEncodeMode(encodeMode), sum)
}

func SumFile(hash crypto.Hash, origFile string) ([]byte, error) {
	if len(origFile) <= 0 {
		return nil, errors.New("origFile path is empty")
	}
	// 打开读写文件
	fiR, errOpenR := os.Open(origFile)
	if errOpenR != nil {
		fmt.Println("open file error: ", errOpenR)
		return nil, errOpenR
	}
	defer fiR.Close()
	h := hash.New()
	blockSize := FILE_BUFFER_SIZE
	// 开始文件加密
	reader := bufio.NewReader(fiR)
	var bRead = make([]byte, blockSize)
	emptyShaFlag := true
	for {
		nR, errR := reader.Read(bRead)
		if errR != nil && errR != io.EOF {
			return nil, errR
		}
		if nR == 0 {
			if emptyShaFlag {
				emptyShaFlag = false
				// sum缓存写入
				_, errSum := h.Write(make([]byte, 0))
				if errSum != nil {
					return nil, errSum
				}
			} else {
				emptyShaFlag = false
			}
			break
		} else if nR > 0 {
			emptyShaFlag = false
			// sum缓存写入
			_, errSum := h.Write(bRead[0:nR])
			if errSum != nil {
				return nil, errSum
			}
			if nR < blockSize {
				break
			}
		} else {
			emptyShaFlag = false
			break
		}
	}
	sum := h.Sum(nil)
	return sum, nil
}

func SumFileByString(encodeMode MODE_ENCODE, hash crypto.Hash, origFile string) (string, error) {
	sum, err := SumFile(hash, origFile)
	if err != nil {
		return "", err
	}
	return EncodingToString(ParseEncodeMode(encodeMode), sum)
}
