package httpx

import (
	"bytes"
	"net/http"
)

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

func DoHttpGet(urlOfGet string, data interface{}, result interface{}) (*HttpxResponse, error) {
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
