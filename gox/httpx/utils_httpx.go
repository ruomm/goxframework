/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/11/14 16:42
 * @version 1.0
 */
package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// 获取请求的真实RemoteAddr，包含host:port
func GetRealRemoteAddr(reqest *http.Request, proxyIpHeaderKey ...string) string {
	//X-Forwarded-For：Squid 服务代理
	//Proxy-Client-IP：apache 服务代理
	//WL-Proxy-Client-IP：weblogic 服务代理
	//HTTP_CLIENT_IP：有些代理服务器
	//HTTP_X_FORWARDED_FOR：有些代理服务器
	//X-Real-IP：nginx服务代理
	ipAddresses := ""
	localIP := "127.0.0.1"
	var realIpHeaderKyes []string
	if len(proxyIpHeaderKey) > 0 {
		realIpHeaderKyes = append(realIpHeaderKyes, proxyIpHeaderKey...)
	} else {
		//realIpHeaderKyes = []string{"X-Forwarded-For", "Proxy-Client-IP", "WL-Proxy-Client-IP", "HTTP_CLIENT_IP", "X-Real-IP", "HTTP_X_FORWARDED_FOR"}
		realIpHeaderKyes = []string{"X-Real-IP", "X-Forwarded-For", "Proxy-Client-IP", "WL-Proxy-Client-IP", "HTTP_CLIENT_IP", "HTTP_X_FORWARDED_FOR"}
	}
	for _, ipHeaderKey := range realIpHeaderKyes {
		if ipHeaderKey == "" {
			continue
		}
		tmpIpAddresses := reqest.Header.Get(ipHeaderKey)
		if len(tmpIpAddresses) <= 0 {
			continue
		}
		tmpIpAddressesLower := strings.ToLower(tmpIpAddresses)
		if "unknown" == tmpIpAddressesLower || localIP == tmpIpAddressesLower {
			continue
		}
		ipAddresses = tmpIpAddresses
		break
	}
	ip := ""
	//有些网络通过多层代理，那么获取到的ip就会有多个，一般都是通过逗号（,）分割开来，并且第一个ip为客户端的真实IP
	if len(ipAddresses) > 0 {
		ip = strings.Split(ipAddresses, ",")[0]
	}
	//还是不能获取到，reqest.RemoteAddr
	if len(ip) == 0 {
		ip = reqest.RemoteAddr
	}
	return ip
}

// 获取常用协议的默认端口，不支持的协议返回0
func GetProtocolDefaultPort(protocol string) int {
	port := 0
	protocolLower := strings.ToLower(protocol)
	if protocolLower == "http" {
		port = 80
	} else if protocolLower == "https" {
		port = 443
	} else if protocolLower == "ftp" {
		port = 21
	} else if protocolLower == "tftp" {
		port = 69
	} else if protocolLower == "smtp" {
		port = 25
	} else if protocolLower == "pop3" {
		port = 110
	} else if protocolLower == "imap" {
		port = 143
	} else if protocolLower == "dns" {
		port = 53
	} else if protocolLower == "ssh" {
		port = 22
	} else if protocolLower == "sftp" {
		port = 22
	} else if protocolLower == "telnet" {
		port = 23
	} else if protocolLower == "ntp" {
		port = 123
	}
	return port
}

// 解析地址端口，返回host、port
func SplitHostPort(hostPort string) (string, int, error) {
	if len(hostPort) == 0 {
		return "", 0, errors.New("host port string is empty")
	}
	hostPorts := strings.Split(hostPort, ":")
	hostPortLen := len(hostPorts)
	if hostPortLen <= 0 {
		return "", 0, errors.New("host port string is invalid")
	}
	hostStr := ""
	portStr := ""
	port := 0
	var errMsgs []string
	if hostPortLen == 1 {
		hostStr = hostPorts[0]
	} else if hostPortLen == 2 {
		hostStr = hostPorts[0]
		portStr = hostPorts[1]
	} else {
		hostStr = hostPorts[0]
		portStr = hostPorts[1]
		//errMsgs = append(errMsgs, "host port string after split is invalid")
	}
	if len(portStr) > 0 {
		intVal, intValErr := strconv.ParseInt(portStr, 10, 64)
		if intValErr != nil {
			errMsgs = append(errMsgs, intValErr.Error())
		} else {
			portTmp := int(intVal)
			if portTmp < 0 {
				errMsgs = append(errMsgs, "port after parse is invalid")
			} else {
				port = portTmp
			}
		}
	}
	if len(errMsgs) > 0 {
		return hostStr, port, errors.New(strings.Join(errMsgs, "\n"))
	} else {
		return hostStr, port, nil
	}
	return hostStr, port, nil
}

// 解析URL单元，返回协议、host、port、相对路径
func ParseUrlElement(urlstr string) (string, string, int, string, error) {
	url := strings.ReplaceAll(urlstr, "\\", "/")
	protocol := ""
	contentStr := ""
	hostPortStr := ""
	protocolIndex := strings.Index(url, ":/")
	if protocolIndex >= 0 {
		protocol = url[0:protocolIndex]
		contentStr = url[protocolIndex+3:]
	} else {
		protocol = ""
		contentStr = url
	}
	uriItems := strings.Split(contentStr, "/")
	for _, uriItem := range uriItems {
		if len(uriItem) > 0 {
			hostPortStr = uriItem
			break
		}
	}
	relativeUrl := ""
	if len(hostPortStr) > 0 {
		relativeUrlIndex := strings.Index(contentStr, hostPortStr)
		if relativeUrlIndex >= 0 {
			relativeUrl = contentStr[relativeUrlIndex+len(hostPortStr):]
		} else {
			relativeUrl = ""
		}
	}
	if relativeUrl == "/" {
		relativeUrl = ""
	}
	host, port, err := SplitHostPort(hostPortStr)
	if err == nil && port == 0 && len(protocol) > 0 {
		port = GetProtocolDefaultPort(protocol)
	}
	return protocol, host, port, relativeUrl, err
}
