package corex

import (
	"net/url"
	"strings"
)

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
