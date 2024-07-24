package corex

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	file_extension_separator = "."
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/**
 * 文件路径是否以文件夹连接符号(File.separator)结尾
 *
 * @param filePath 文件路径
 * @return 文件路径是否以文件夹连接符号(File.separator)结尾
 */
func IsEndWithFileSeparator(filePath string) bool {
	if len(filePath) <= 0 {
		return false
	}
	flag := strings.HasSuffix(filePath, "\\")
	if !flag {
		flag = strings.HasSuffix(filePath, "/")
	}
	return flag
}

/**
 * 文件路径是否以文件夹连接符号(File.separator)开始
 *
 * @param filePath 文件路径
 * @return 文件路径是否以文件夹连接符号(File.separator)开始
 */
func IsStartWithFileSeparator(filePath string) bool {
	if len(filePath) <= 0 {
		return false
	}
	flag := strings.HasPrefix(filePath, "\\")
	if !flag {
		flag = strings.HasPrefix(filePath, "/")
	}
	return flag
}

/**
 * 查找文件分割符号最后的位置
 *
 * @param filePath 文件路径
 * @return 返回分割符号最后的位置，空返回-99，没有返回-1，
 */
func LastIndexOfFileSeparator(filePath string) int {
	if len(filePath) <= 0 {
		return -1
	}
	filePosi1 := strings.LastIndex(filePath, "\\")
	filePosi2 := strings.LastIndex(filePath, "/")
	if filePosi1 < 0 && filePosi2 < 0 {
		return -1
	} else if filePosi1 > filePosi2 {
		return filePosi1
	} else {
		return filePosi2
	}
}

/**
 * 查找文件分割符号第一个的位置
 *
 * @param filePath 文件路径
 * @return 返回分割符号第一个的位置位置，空返回-99，没有返回-1，
 */
func FirstIndexOfFileSeparator(filePath string) int {
	if len(filePath) <= 0 {
		return -99
	}
	filePosi1 := strings.Index(filePath, "\\")
	filePosi2 := strings.Index(filePath, "/")
	if filePosi1 < 0 && filePosi2 < 0 {
		return -1
	} else if filePosi1 < 0 {
		return filePosi2
	} else if filePosi2 < 0 {
		return filePosi1
	} else if filePosi1 < filePosi2 {
		return filePosi1
	} else {
		return filePosi2
	}
}

/**
 * 依据文文件或文件夹路径智能获取文件夹连接符号(File.separator)。
 *
 * @param filePath 文件或文件夹路径
 * @return 文件夹连接符号(File.separator)
 */
func ParseFileSeparator(filePath string) string {
	if len(filePath) <= 0 {
		return string(os.PathSeparator)
	}
	filePosi1 := strings.Index(filePath, "\\")
	filePosi2 := strings.Index(filePath, "/")
	if filePosi1 >= 0 && filePosi2 >= 0 {
		if filePosi1 > filePosi2 {
			return "\\"
		} else {
			return "/"
		}
	} else if filePosi1 >= 0 {
		return "\\"
	} else if filePosi2 >= 0 {
		return "/"
	} else {
		return string(os.PathSeparator)
	}
}

/**
 * 拼接2个目录或1个目录1个文件为完整文件路径
 *
 * @param filePath            目录路径
 * @param fileName        文件名称或目录路径
 * @param separatorRevise 是否修订名称中的文件分隔路径
 * @return 完整文件路径
 */
func ConcatFilePath(filePath string, fileName string, separatorRevise bool) string {
	if len(filePath) == 0 && len(fileName) == 0 {
		return ""
	} else if len(filePath) == 0 {
		return fileName
	} else if len(fileName) == 0 {
		return filePath
	} else {
		realFileName := ""
		if len(fileName) == 0 {
			realFileName = ""
		} else if IsStartWithFileSeparator(fileName) {
			realFileName = fileName[1:]
		} else {
			realFileName = fileName
		}
		if len(realFileName) == 0 {
			return filePath
		}
		tmpFileSeparator := ParseFileSeparator(filePath)
		if separatorRevise {
			realFileName = strings.ReplaceAll(realFileName, "\\", tmpFileSeparator)
			realFileName = strings.ReplaceAll(realFileName, "/", tmpFileSeparator)
		}
		if IsEndWithFileSeparator(filePath) {
			return filePath + realFileName
		} else {
			return filePath + tmpFileSeparator + realFileName
		}
	}
}

/**
 * 拼接2个目录完整目录路径
 *
 * @param filePath            目录路径
 * @param subPath         次级目录路径
 * @param separatorRevise 是否修订名称中的文件分隔路径
 * @return 完整文件路径
 */
func ConcatDirPath(filePath string, subPath string, separatorRevise bool) string {
	if len(filePath) == 0 && len(subPath) == 0 {
		return ""
	} else if len(filePath) == 0 {
		return subPath
	} else if len(subPath) == 0 {
		return filePath
	} else {
		realSubPath := ""
		if len(subPath) == 0 {
			realSubPath = ""
		} else if IsStartWithFileSeparator(subPath) {
			//realSubPath = subPath[0 : len(subPath)-1]
			realSubPath = subPath[1:]
		} else {
			realSubPath = subPath
		}
		tmpFileSeparator := ParseFileSeparator(filePath)
		if len(realSubPath) == 0 {
			if IsEndWithFileSeparator(filePath) {
				return filePath
			} else {
				return filePath + tmpFileSeparator
			}
		}
		if separatorRevise {
			realSubPath = strings.ReplaceAll(realSubPath, "\\", tmpFileSeparator)
			realSubPath = strings.ReplaceAll(realSubPath, "/", tmpFileSeparator)
		}
		resultPath := ""
		if IsEndWithFileSeparator(filePath) {
			resultPath = filePath + realSubPath
		} else {
			resultPath = filePath + tmpFileSeparator + realSubPath
		}
		if IsEndWithFileSeparator(resultPath) {
			return resultPath
		} else {
			return resultPath + tmpFileSeparator
		}
	}
}

// 获取当前的执行路径
// C:\Users\Vic\AppData\Local\Temp\
func GetCurrentPath() string {
	//s, err := exec.LookPath(os.Args[0])
	//checkErr(err)
	//i := strings.LastIndex(s, "\\")
	//path := string(s[0 : i+1])
	//return path
	return filepath.Dir(os.Args[0])
}

func IsAbsDir(path string) bool {
	if len(path) <= 0 {
		return false
	} else if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "\\") {
		return true
	} else {
		indexSpec := strings.Index(path, ":")
		if indexSpec >= 0 && indexSpec < 12 {
			return true
		} else {
			return false
		}
	}
}

func GetAbsDir(relativePath string) string {
	if len(relativePath) <= 0 {
		return filepath.Dir(os.Args[0])
	} else if IsAbsDir(relativePath) {
		return relativePath
	} else {
		dir := filepath.Dir(os.Args[0])
		return path.Join(dir, relativePath)
	}
}

func IsFileExit(file_path string) (bool, error) {
	_, err := os.Stat(file_path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func IsFileExitWithErr(file_path string) (bool, error) {
	_, err := os.Stat(file_path)
	if err == nil {
		return true, nil
	}
	return false, err
}

func MkdirAll(file_path string) (bool, error) {
	fileExit, errExit := IsFileExit(file_path)
	if errExit != nil {
		fmt.Println(errExit)
		return false, errExit
	}
	if fileExit {
		return true, nil
	}
	errDir := os.MkdirAll(file_path, 0777)
	if errDir != nil {
		fmt.Println(errDir)
		return false, errDir
	} else {
		fmt.Println("Successfully created directories")
		return true, nil
	}
}

func ReadFile(path string) (string, error) {
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	contentString := string(contentBytes)
	return contentString, nil
}

func ReadFileBuffer(path string) (string, error) {
	fi, errOpen := os.Open(path)
	if errOpen != nil {
		fmt.Println("open file error: ", errOpen)
		return "", errOpen
	}
	defer fi.Close()
	reader := bufio.NewReader(fi)
	var build strings.Builder
	//hasRead := false
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("read file to string fail", err)
			return "", err
		}
		if line != "" {
			build.WriteString(line)
		}
		if err == io.EOF {
			break
		}
	}
	return build.String(), nil
}

func WriteFile(path string, content string) error {
	// 文件不存在时创建，存在时清空，以读写方式打开
	// 如果不清空，再次打开写入时默认是从最开始写入的
	fi, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		fmt.Println("open file error: ", err)
		return err
	}
	defer fi.Close()
	_, err = fi.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func WriteFileBuffer(path string, content string) error {
	// 文件不存在时创建，存在时清空，以读写方式打开
	// 如果不清空，再次打开写入时默认是从最开始写入的
	fi, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		fmt.Println("open file error: ", err)
		return err
	}
	defer fi.Close()
	//bufio.NewWriterSize()
	bfd := bufio.NewWriter(fi) // 调用的也是NewWriterSize，默认写入大小是4096
	_, err = bfd.WriteString(content)
	if err != nil {
		return err
	}
	// 清空缓冲区，进行缓冲区内容落到磁盘
	err = bfd.Flush()
	if err != nil {
		return err
	}
	return nil
}

func CopyFile(pathSrc string, pathDest string) error {
	fiR, errOpenR := os.Open(pathSrc)
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
	bufferSize := 1024
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
			break
		}
		_, errW := writer.Write(bRead[0:nR])
		if errW != nil {
			return errW
		}
		if nR < bufferSize {
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

/**
 * 从路径获取文件名，包括后缀
 *
 * <pre>
 *      getFileName(null)               =   null
 *      getFileName("")                 =   ""
 *      getFileName("   ")              =   "   "
 *      getFileName("a.mp3")            =   "a.mp3"
 *      getFileName("a.b.rmvb")         =   "a.b.rmvb"
 *      getFileName("abc")              =   "abc"
 *      getFileName("c:\\")              =   ""
 *      getFileName("c:\\a")             =   "a"
 *      getFileName("c:\\a.b")           =   "a.b"
 *      getFileName("c:a.txt\\a")        =   "a"
 *      getFileName("/home/admin")      =   "admin"
 *      getFileName("/home/admin/a.txt/b.mp3")  =   "b.mp3"
 * </pre>
 *
 * @param filePath 文件路径 -
 * @return 从路径文件名，包括后缀
 */
func GetFileName(filePath string) string {
	if len(filePath) <= 0 {
		return filePath
	}

	filePosi := LastIndexOfFileSeparator(filePath)
	if filePosi < 0 {
		return filePath
	} else {
		return filePath[filePosi+1:]
	}
}

/**
 * 从路径获取文件名，不包括后缀
 *
 * <pre>
 *      getFileNameWithoutExtension(null)               =   null
 *      getFileNameWithoutExtension("")                 =   ""
 *      getFileNameWithoutExtension("   ")              =   "   "
 *      getFileNameWithoutExtension("abc")              =   "abc"
 *      getFileNameWithoutExtension("a.mp3")            =   "a"
 *      getFileNameWithoutExtension("a.b.rmvb")         =   "a.b"
 *      getFileNameWithoutExtension("c:\\")              =   ""
 *      getFileNameWithoutExtension("c:\\a")             =   "a"
 *      getFileNameWithoutExtension("c:\\a.b")           =   "a"
 *      getFileNameWithoutExtension("c:a.txt\\a")        =   "a"
 *      getFileNameWithoutExtension("/home/admin")      =   "admin"
 *      getFileNameWithoutExtension("/home/admin/a.txt/b.mp3")  =   "b"
 * </pre>
 *
 * @param filePath 文件路径 -
 * @return 从路径文件名，不包括后缀
 * @see
 */
func GetFileNameWithoutExtension(filePath string) string {
	if len(filePath) <= 0 {
		return filePath
	}
	extenPosi := strings.LastIndex(filePath, file_extension_separator)
	filePosi := LastIndexOfFileSeparator(filePath)
	if filePosi < 0 {
		if extenPosi < 0 {
			return filePath
		} else {
			return filePath[0:extenPosi]
		}
	}
	if extenPosi < 0 {
		return filePath[filePosi+1:]
	} else {
		if filePosi < extenPosi {
			return filePath[filePosi+1 : extenPosi]
		} else {
			return filePath[filePosi+1:]
		}
	}
}

/**
 * 获取文件的后缀从路径
 *
 * <pre>
 *      getFileExtension(null)               =   ""
 *      getFileExtension("")                 =   ""
 *      getFileExtension("   ")              =   "   "
 *      getFileExtension("a.mp3")            =   "mp3"
 *      getFileExtension("a.b.rmvb")         =   "rmvb"
 *      getFileExtension("abc")              =   ""
 *      getFileExtension("c:\\")              =   ""
 *      getFileExtension("c:\\a")             =   ""
 *      getFileExtension("c:\\a.b")           =   "b"
 *      getFileExtension("c:a.txt\\a")        =   ""
 *      getFileExtension("/home/admin")      =   ""
 *      getFileExtension("/home/admin/a.txt/b")  =   ""
 *      getFileExtension("/home/admin/a.txt/b.mp3")  =   "mp3"
 * </pre>
 *
 * @param filePath
 * @return
 */
func GetFileExtension(filePath string) string {
	if len(filePath) <= 0 {
		return ""
	}
	extenPosi := strings.LastIndex(filePath, file_extension_separator)
	filePosi := LastIndexOfFileSeparator(filePath)
	if extenPosi < 0 {
		return ""
	} else if filePosi >= extenPosi {
		return ""
	} else {
		return filePath[extenPosi+1:]
	}
}
