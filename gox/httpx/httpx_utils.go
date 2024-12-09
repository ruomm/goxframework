package httpx

import (
	"bytes"
	"errors"
	"github.com/ruomm/goxframework/gox/corex"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// JSON请求自动封装
// httpClient不传入则使用默认的http.Client
// reqUrl:请求路径
// httpxMethod:请求的方法，reqObj定义有HttpxMethod则此参数可以不填写，使用reqObj定义的请求方法
// httpxHeaders自定义的请求Header列表
// reqObj:请求的实体类，通过json、xreq_param、xreq_query、xreq_header注解定义请求参数
func DoHttpToResponse(httpClient *http.Client, reqUrl string, httpxMethod string, httpxHeaders map[string]string, reqOjb interface{}) (*http.Response, error) {
	reqMethod, reqBody, reqParam, reqQuery, reqHeaders, err := ParseToRequest(httpxMethod, reqOjb)
	if err != nil {
		return nil, err
	}
	// headerMap值修订
	reqHeaderMap := make(map[string]string, 0)
	for headerKey, HeaderVal := range httpxHeaders {
		if _, ok := reqHeaderMap[headerKey]; !ok {
			reqHeaderMap[headerKey] = HeaderVal
		}
	}
	for headerKey, HeaderVal := range reqHeaders {
		if _, ok := reqHeaderMap[headerKey]; !ok {
			reqHeaderMap[headerKey] = HeaderVal
		}
	}
	// 请求路径获取
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
	var client *http.Client
	if nil == httpClient {
		client = &http.Client{}
	} else {
		client = httpClient
	}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Execute Request  err: %v", contentType, err.Error()))
		return resp, err
	}
	if nil == resp {
		return nil, errors.New("Http Response is Empty")
	}
	return resp, nil
}

// JSON请求自动封装和解封，
// httpClient不传入则使用默认的http.Client
// reqUrl:请求路径
// httpxMethod:请求的方法，reqObj定义有HttpxMethod则此参数可以不填写，使用reqObj定义的请求方法
// httpxHeaders自定义的请求Header列表
// reqObj:请求的实体类，通过json、xreq_param、xreq_query、xreq_header注解定义请求参数，
// resultObjects:请求结果自动解析为对象列表，多个参数只有第一个解析错误才返回error。resultObjects中间有字符串的，则视为需要解析的JSON节点。
// 如：DoHttpToJson(nil,"http://xxxx","GET",nil,reqObj,jsonA,"data.info",jsonB),jsonB解析的数据来源为data.info节点的数据。
func DoHttpToJson(httpClient *http.Client, reqUrl string, httpxMethod string, httpxHeaders map[string]string, reqOjb interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	reqMethod, reqBody, reqParam, reqQuery, reqHeaders, err := ParseToRequest(httpxMethod, reqOjb)
	if err != nil {
		return nil, err
	}
	// headerMap值修订
	reqHeaderMap := make(map[string]string, 0)
	for headerKey, HeaderVal := range httpxHeaders {
		if _, ok := reqHeaderMap[headerKey]; !ok {
			reqHeaderMap[headerKey] = HeaderVal
		}
	}
	for headerKey, HeaderVal := range reqHeaders {
		if _, ok := reqHeaderMap[headerKey]; !ok {
			reqHeaderMap[headerKey] = HeaderVal
		}
	}
	// 请求路径获取
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
	var client *http.Client
	if nil == httpClient {
		client = &http.Client{}
	} else {
		client = httpClient
	}
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

// http请求-字符串模式
// httpClient不传入则使用默认的http.Client
// httpxMethod:请求的方法，reqObj定义有HttpxMethod则此参数可以不填写，使用reqObj定义的请求方法
// reqUrl:请求路径
// postContentType请求的Content-Type，如实不填写，则依据postStr自动判断
// postStr:请求体文本
func DoHttp(httpClient *http.Client, httpxMethod string, reqUrl string, postContentType string, postStr string) (*HttpxResponse, error) {
	reqMethodVerify := xReqMethodVerify(httpxMethod)
	if !reqMethodVerify {
		return nil, errors.New("Request Method invalid error")
	}
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
	req, err := http.NewRequest(httpxMethod, reqUrl, reqIo)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Build NewRequest err: %v", contentType, err.Error()))
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	var client *http.Client
	if nil == httpClient {
		client = &http.Client{}
	} else {
		client = httpClient
	}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Execute Request  err: %v", contentType, err.Error()))
		return nil, err
	}
	return xToHttpxResponse(resp)
}

// http请求-JSON模式
// httpClient不传入则使用默认的http.Client
// httpxMethod:请求的方法，reqObj定义有HttpxMethod则此参数可以不填写，使用reqObj定义的请求方法
// reqUrl:请求路径
// data:请求的JSON实体
// resultObjects:请求结果自动解析为对象列表，多个参数只有第一个解析错误才返回error。resultObjects中间有字符串的，则视为需要解析的JSON节点。
func DoHttpJson(httpClient *http.Client, httpxMethod string, reqUrl string, data interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	reqMethodVerify := xReqMethodVerify(httpxMethod)
	if !reqMethodVerify {
		return nil, errors.New("Request Method invalid error")
	}
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
	req, err := http.NewRequest(httpxMethod, reqUrl, reqIo)
	if err != nil {
		// logger.Error("Http Post by application/json Build NewRequest err:" + err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	var client *http.Client
	if nil == httpClient {
		client = &http.Client{}
	} else {
		client = httpClient
	}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error("Http Post by application/json Execute Request err:" + err.Error())
		return nil, err
	}
	return xToHttpxResponseJson(resp, resultObjs...)
}

// http请求-POST-字符串模式
// httpClient不传入则使用默认的http.Client
// httpxMethod:请求的方法，reqObj定义有HttpxMethod则此参数可以不填写，使用reqObj定义的请求方法
// reqUrl:请求路径
// data:请求的JSON实体
// resultObjects:请求结果自动解析为对象列表，多个参数只有第一个解析错误才返回error。resultObjects中间有字符串的，则视为需要解析的JSON节点。
func DoPost(httpClient *http.Client, reqUrl string, postContentType string, postStr string) (*HttpxResponse, error) {
	return DoHttp(httpClient, http.MethodPost, reqUrl, postContentType, postStr)
}

// http请求-POST-JSON模式
// httpClient不传入则使用默认的http.Client
// httpxMethod:请求的方法，reqObj定义有HttpxMethod则此参数可以不填写，使用reqObj定义的请求方法
// reqUrl:请求路径
// data:请求的JSON实体
// resultObjects:请求结果自动解析为对象列表，多个参数只有第一个解析错误才返回error。resultObjects中间有字符串的，则视为需要解析的JSON节点。
func DoPostJson(httpClient *http.Client, reqUrl string, data interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	return DoHttpJson(httpClient, http.MethodPost, reqUrl, data, resultObjs...)
}

// http请求-PUT-字符串模式
// httpClient不传入则使用默认的http.Client
// httpxMethod:请求的方法，reqObj定义有HttpxMethod则此参数可以不填写，使用reqObj定义的请求方法
// reqUrl:请求路径
// data:请求的JSON实体
// resultObjects:请求结果自动解析为对象列表，多个参数只有第一个解析错误才返回error。resultObjects中间有字符串的，则视为需要解析的JSON节点。
func DoPut(httpClient *http.Client, reqUrl string, postContentType string, postStr string) (*HttpxResponse, error) {
	return DoHttp(httpClient, http.MethodPut, reqUrl, postContentType, postStr)
}

// http请求-PUT-JSON模式
// httpClient不传入则使用默认的http.Client
// httpxMethod:请求的方法，reqObj定义有HttpxMethod则此参数可以不填写，使用reqObj定义的请求方法
// reqUrl:请求路径
// data:请求的JSON实体
// resultObjects:请求结果自动解析为对象列表，多个参数只有第一个解析错误才返回error。resultObjects中间有字符串的，则视为需要解析的JSON节点。
func DoPutJson(httpClient *http.Client, reqUrl string, data interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	return DoHttpJson(httpClient, http.MethodPut, reqUrl, data, resultObjs...)
}

// http请求-DELETE-字符串模式
// httpClient不传入则使用默认的http.Client
// httpxMethod:请求的方法，reqObj定义有HttpxMethod则此参数可以不填写，使用reqObj定义的请求方法
// reqUrl:请求路径
// data:请求的JSON实体
// resultObjects:请求结果自动解析为对象列表，多个参数只有第一个解析错误才返回error。resultObjects中间有字符串的，则视为需要解析的JSON节点。
func DoDelete(httpClient *http.Client, reqUrl string, postContentType string, postStr string) (*HttpxResponse, error) {
	return DoHttp(httpClient, http.MethodDelete, reqUrl, postContentType, postStr)
}

// http请求-DELETE-JSON模式
// httpClient不传入则使用默认的http.Client
// httpxMethod:请求的方法，reqObj定义有HttpxMethod则此参数可以不填写，使用reqObj定义的请求方法
// reqUrl:请求路径
// data:请求的JSON实体
// resultObjects:请求结果自动解析为对象列表，多个参数只有第一个解析错误才返回error。resultObjects中间有字符串的，则视为需要解析的JSON节点。
func DoDeleteJson(httpClient *http.Client, reqUrl string, data interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	return DoHttpJson(httpClient, http.MethodDelete, reqUrl, data, resultObjs...)
}

// http请求-GETT-string、url.Values、xreq_query模式
// httpClient不传入则使用默认的http.Client
// urlOfGet:请求路径
// data:请求实体，支持stirng、url.Values、xreq_query模式
func DoGet(httpClient *http.Client, urlOfGet string, data interface{}) (*HttpxResponse, error) {
	urlData, err := ParseToUrlEncodeString(data)
	if err != nil {
		// logger.Error("Http Get Encode Request Data err:" + err.Error())
		return nil, err
	}
	var resp *http.Response
	// logger.Debug("Http Get Encode Request Data ok:" + urlData)
	if nil == httpClient {
		resp, err = http.Get(xParseRealGetUrl(urlOfGet, urlData))
	} else {
		resp, err = httpClient.Get(xParseRealGetUrl(urlOfGet, urlData))
	}
	if err != nil {
		// logger.Error("Http Get Do Request err:" + err.Error())
		return nil, err
	}
	return xToHttpxResponse(resp)
}

// http请求-GET-string、url.Values、xreq_query模式
// httpClient不传入则使用默认的http.Client
// urlOfGet:请求路径
// data:请求实体，支持stirng、url.Values、xreq_query模式
// resultObjects:请求结果自动解析为对象列表，多个参数只有第一个解析错误才返回error。resultObjects中间有字符串的，则视为需要解析的JSON节点。
func DoGetJson(httpClient *http.Client, urlOfGet string, data interface{}, resultObjs ...interface{}) (*HttpxResponse, error) {
	urlData, err := ParseToUrlEncodeString(data)
	if err != nil {
		// logger.Error("Http Get Encode Request Data err:" + err.Error())
		return nil, err
	}
	var resp *http.Response
	// logger.Debug("Http Get Encode Request Data ok:" + urlData)
	if nil == httpClient {
		resp, err = http.Get(xParseRealGetUrl(urlOfGet, urlData))
	} else {
		resp, err = httpClient.Get(xParseRealGetUrl(urlOfGet, urlData))
	}
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
