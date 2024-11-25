package corex

import (
	"net/url"
	"strings"
)

// 一个elem的uri路径拼接上contextPath路径返回，elem以/开头，则返回以/开头，emel不以/开头，，则返回不以/开头
// 如：ParseUrlWithContextPath("web/","/api") = "/web/api",ParseUrlWithContextPath("/web/","api") = "web/api"
func ParseUrlWithContextPath(contextPath string, elem string) (string, error) {
	if len(contextPath) <= 0 {
		return elem, nil
	} else if len(elem) <= 0 {
		return contextPath, nil
	}
	resultUrl, err := url.JoinPath(contextPath, elem)
	if err != nil {
		return resultUrl, err
	}
	if strings.HasPrefix(elem, "/") && !strings.HasPrefix(resultUrl, "/") {
		return "/" + resultUrl, nil
	} else if !strings.HasPrefix(elem, "/") && strings.HasPrefix(resultUrl, "/") {
		return resultUrl[1:], nil
	} else {
		return resultUrl, nil
	}
}
