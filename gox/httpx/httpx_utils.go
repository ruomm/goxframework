package httpx

import (
	"bytes"
	"net/http"
	"strings"
)

// JSON请求自动封装和解封
func DoHttpJson(reqUrl string, httpxMethod string, reqOjb interface{}, result interface{}) (*HttpxResponse, error) {
	reqMethod, reqBody, reqParam, reqQuery, reqHeaderMap, err := ParseToRequest(reqOjb)
	if err != nil {
		return nil, err
	}
	if len(httpxMethod) > 0 {
		reqMethod = strings.ToUpper(httpxMethod)
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
	return xToHttpxResponseJson(resp, result)
}

func DoHttpPost(postUrl string, postContentType string, postStr string) (*HttpxResponse, error) {
	contentType := ""
	if len(postContentType) > 0 {
		contentType = postContentType
	} else {
		contentType = xParseRequestMime(postStr)
	}
	// logger.Debug(fmt.Sprintf("Http Post by %v Message of request:", contentType))
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer([]byte(postStr)))
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

func DoHttpPostJson(postUrl string, data interface{}, result interface{}) (*HttpxResponse, error) {
	jsonData, err := ParseToJSONByte(data)
	if err != nil {
		// logger.Error("Http Post by application/json Marshal Request Data err:" + err.Error())
		return nil, err
	}
	// logger.Debug("Http Post by application/json Message of request:" + string(jsonData))
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonData))
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
	return xToHttpxResponseJson(resp, result)
}

func DoHttpPut(postUrl string, postContentType string, postStr string) (*HttpxResponse, error) {
	contentType := ""
	if len(postContentType) > 0 {
		contentType = postContentType
	} else {
		contentType = xParseRequestMime(postStr)
	}
	// logger.Debug(fmt.Sprintf("Http Post by %v Message of request:", contentType))
	req, err := http.NewRequest("PUT", postUrl, bytes.NewBuffer([]byte(postStr)))
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

func DoHttpPutJson(postUrl string, data interface{}, result interface{}) (*HttpxResponse, error) {
	jsonData, err := ParseToJSONByte(data)
	if err != nil {
		// logger.Error("Http Post by application/json Marshal Request Data err:" + err.Error())
		return nil, err
	}
	// logger.Debug("Http Post by application/json Message of request:" + string(jsonData))
	req, err := http.NewRequest("PUT", postUrl, bytes.NewBuffer(jsonData))
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
	return xToHttpxResponseJson(resp, result)
}

func DoHttpDelete(postUrl string, postContentType string, postStr string) (*HttpxResponse, error) {
	contentType := ""
	if len(postContentType) > 0 {
		contentType = postContentType
	} else {
		contentType = xParseRequestMime(postStr)
	}
	// logger.Debug(fmt.Sprintf("Http Post by %v Message of request:", contentType))
	req, err := http.NewRequest("DELETE", postUrl, bytes.NewBuffer([]byte(postStr)))
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

func DoHttpDeleteJson(postUrl string, data interface{}, result interface{}) (*HttpxResponse, error) {
	jsonData, err := ParseToJSONByte(data)
	if err != nil {
		// logger.Error("Http Post by application/json Marshal Request Data err:" + err.Error())
		return nil, err
	}
	// logger.Debug("Http Post by application/json Message of request:" + string(jsonData))
	req, err := http.NewRequest("DELETE", postUrl, bytes.NewBuffer(jsonData))
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
	return xToHttpxResponseJson(resp, result)
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

func DoHttpGetJson(urlOfGet string, data interface{}, result interface{}) (*HttpxResponse, error) {
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
	return xToHttpxResponseJson(resp, result)
}
