/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/11/14 16:42
 * @version 1.0
 */
package httpx

import (
	"net/http"
	"strings"
)

func GetRequestIP(reqest *http.Request, proxyIpHeaderKey ...string) string {
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
