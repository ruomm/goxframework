package httpx

import (
	"bytes"
	"github.com/ruomm/goxframework/gox/corex"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// JSON请求自动封装和解封
func DoHttpJson(reqUrl string, httpxMethod string, reqOjb interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	reqMethod, reqBody, reqParam, reqQuery, reqHeaderMap, err := ParseToRequest(httpxMethod, reqOjb)
	if err != nil {
		return nil, err
	}
	requestUrl := xParseRequestUrl(reqUrl, reqParam, reqQuery)
	var req *http.Request = nil
	if nil == reqBody {
		req, err = http.NewRequest(reqMethod, requestUrl, nil)
	} else {
		req, err = http.NewRequest(reqMethod, requestUrl, bytes.NewBuffer(reqBody))
	}
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Build NewRequest err: %v", contentType, err.Error()))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if nil != reqHeaderMap && len(reqHeaderMap) > 0 {
		for headerKey, HeaderVal := range reqHeaderMap {
			req.Header.Set(headerKey, HeaderVal)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Execute Request  err: %v", contentType, err.Error()))
		return nil, err
	}
	if resultObjs == nil || len(resultObjs) <= 0 {
		return xToHttpxResponse(resp)
	}
	return xToHttpxResponseJson(resp, resultObjs...)
}

func DoHttpPost(reqUrl string, postContentType string, postStr string) (*HttpxResponse, error) {
	contentType := ""
	if len(postContentType) > 0 {
		contentType = postContentType
	} else {
		contentType = xParseRequestMime(postStr)
	}
	// logger.Debug(fmt.Sprintf("Http Post by %v Message of request:", contentType))
	var reqIo io.Reader = nil
	if len(postStr) > 0 {
		reqIo = bytes.NewBuffer([]byte(postStr))
	}
	req, err := http.NewRequest("POST", reqUrl, reqIo)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Build NewRequest err: %v", contentType, err.Error()))
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Execute Request  err: %v", contentType, err.Error()))
		return nil, err
	}
	return xToHttpxResponse(resp)
}

func DoHttpPostJson(reqUrl string, data interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	jsonData, err := xParseReqJson(data)
	if err != nil {
		// logger.Error("Http Post by application/json Marshal Request Data err:" + err.Error())
		return nil, err
	}
	// logger.Debug("Http Post by application/json Message of request:" + string(jsonData))
	var reqIo io.Reader = nil
	if nil != jsonData {
		reqIo = bytes.NewBuffer(jsonData)
	}
	req, err := http.NewRequest("POST", reqUrl, reqIo)
	if err != nil {
		// logger.Error("Http Post by application/json Build NewRequest err:" + err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error("Http Post by application/json Execute Request err:" + err.Error())
		return nil, err
	}
	return xToHttpxResponseJson(resp, resultObjs...)
}

func DoHttpPut(reqUrl string, postContentType string, postStr string) (*HttpxResponse, error) {
	contentType := ""
	if len(postContentType) > 0 {
		contentType = postContentType
	} else {
		contentType = xParseRequestMime(postStr)
	}
	// logger.Debug(fmt.Sprintf("Http Post by %v Message of request:", contentType))
	var reqIo io.Reader = nil
	if len(postStr) > 0 {
		reqIo = bytes.NewBuffer([]byte(postStr))
	}
	req, err := http.NewRequest("PUT", reqUrl, reqIo)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Build NewRequest err: %v", contentType, err.Error()))
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Execute Request  err: %v", contentType, err.Error()))
		return nil, err
	}
	return xToHttpxResponse(resp)
}

func DoHttpPutJson(reqUrl string, data interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	jsonData, err := xParseReqJson(data)
	if err != nil {
		// logger.Error("Http Post by application/json Marshal Request Data err:" + err.Error())
		return nil, err
	}
	// logger.Debug("Http Post by application/json Message of request:" + string(jsonData))
	var reqIo io.Reader = nil
	if nil != jsonData {
		reqIo = bytes.NewBuffer(jsonData)
	}
	req, err := http.NewRequest("PUT", reqUrl, reqIo)
	if err != nil {
		// logger.Error("Http Post by application/json Build NewRequest err:" + err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error("Http Post by application/json Execute Request err:" + err.Error())
		return nil, err
	}
	return xToHttpxResponseJson(resp, resultObjs...)
}

func DoHttpDelete(reqUrl string, postContentType string, postStr string) (*HttpxResponse, error) {
	contentType := ""
	if len(postContentType) > 0 {
		contentType = postContentType
	} else {
		contentType = xParseRequestMime(postStr)
	}
	// logger.Debug(fmt.Sprintf("Http Post by %v Message of request:", contentType))
	var reqIo io.Reader = nil
	if len(postStr) > 0 {
		reqIo = bytes.NewBuffer([]byte(postStr))
	}
	req, err := http.NewRequest("DELETE", reqUrl, reqIo)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Build NewRequest err: %v", contentType, err.Error()))
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Execute Request  err: %v", contentType, err.Error()))
		return nil, err
	}
	return xToHttpxResponse(resp)
}

func DoHttpDeleteJson(reqUrl string, data interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	jsonData, err := xParseReqJson(data)
	if err != nil {
		// logger.Error("Http Post by application/json Marshal Request Data err:" + err.Error())
		return nil, err
	}
	// logger.Debug("Http Post by application/json Message of request:" + string(jsonData))
	var reqIo io.Reader = nil
	if nil != jsonData {
		reqIo = bytes.NewBuffer(jsonData)
	}
	req, err := http.NewRequest("DELETE", reqUrl, reqIo)
	if err != nil {
		// logger.Error("Http Post by application/json Build NewRequest err:" + err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error("Http Post by application/json Execute Request err:" + err.Error())
		return nil, err
	}
	return xToHttpxResponseJson(resp, resultObjs...)
}

func DoHttpGet(urlOfGet string, data interface{}) (*HttpxResponse, error) {
	urlData, err := ParseToUrlEncodeString(data)
	if err != nil {
		// logger.Error("Http Get Encode Request Data err:" + err.Error())
		return nil, err
	}
	// logger.Debug("Http Get Encode Request Data ok:" + urlData)
	resp, err := http.Get(xParseRealGetUrl(urlOfGet, urlData))
	if err != nil {
		// logger.Error("Http Get Do Request err:" + err.Error())
		return nil, err
	}
	return xToHttpxResponse(resp)
}

func DoHttpGetJson(urlOfGet string, data interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	urlData, err := ParseToUrlEncodeString(data)
	if err != nil {
		// logger.Error("Http Get Encode Request Data err:" + err.Error())
		return nil, err
	}
	// logger.Debug("Http Get Encode Request Data ok:" + urlData)
	resp, err := http.Get(xParseRealGetUrl(urlOfGet, urlData))
	if err != nil {
		// logger.Error("Http Get Do Request err:" + err.Error())
		return nil, err
	}
	return xToHttpxResponseJson(resp, resultObjs...)
}

// 解析Http强求range头的第一个请求范围
func ParseHttpRangeFirst(rangeHeader string) (int64, int64) {
	if len(rangeHeader) <= 0 {
		return -1, -1
	}
	rangePrefix := "bytes="
	rangeLower := strings.ToLower(rangeHeader)
	if strings.HasPrefix(rangeLower, rangePrefix) && len(rangeLower) > len(rangePrefix) {
		rangeStr := rangeLower[len(rangePrefix):]
		rangeByDh := corex.StringToSlice(rangeStr, ",", false)
		if len(rangeByDh) <= 0 {
			return -1, -1
		}
		rangeAreas := corex.StringToSlice(rangeByDh[0], "-", true)
		if len(rangeAreas) == 1 {
			start, err := strconv.ParseInt(rangeAreas[0], 10, 64)
			if err != nil {
				return -1, -1
			} else {
				return start, -1
			}
		} else if len(rangeAreas) == 2 {
			start, errStart := strconv.ParseInt(rangeAreas[0], 10, 64)
			end, errEnd := strconv.ParseInt(rangeAreas[1], 10, 64)
			if errStart != nil && errEnd != nil {
				return -1, -1
			} else if errStart != nil {
				return -1, end
			} else if errEnd != nil {
				return start, -1
			} else if end >= start {
				return start, end
			} else {
				return -1, -1
			}
		} else {
			return -1, -1
		}
	} else {
		return -1, -1
	}
}
