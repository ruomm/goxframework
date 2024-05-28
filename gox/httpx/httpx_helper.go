/**
 * @copyright 像衍科技-idr.ai
 * @author 牛牛-研发部-www.ruomm.com
 * @create 2024/2/2 09:58
 * @version 1.0
 */
package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// 判断是否200成功
func Success200(httpxResponse *HttpxResponse) bool {
	if nil == httpxResponse {
		return false
	}
	if httpxResponse.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

// 判断是否2xx成功
func Success2xx(httpxResponse *HttpxResponse) bool {
	if nil == httpxResponse {
		return false
	}
	if httpxResponse.StatusCode >= 200 && httpxResponse.StatusCode < 300 {
		return true
	} else {
		return false
	}
}

// 转换响应请求
func xToHttpxResponse(resp *http.Response) (*HttpxResponse, error) {
	if nil == resp {
		return nil, errors.New("Http Response is Empty")
	}
	defer resp.Body.Close()
	httpxResponse := HttpxResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Proto:      resp.Proto,
		ProtoMajor: resp.ProtoMajor,
		ProtoMinor: resp.ProtoMinor,

		// Header maps header keys to values. If the response had multiple
		// headers with the same key, they may be concatenated, with comma
		// delimiters.  (RFC 7230, section 3.2.2 requires that multiple headers
		// be semantically equivalent to a comma-delimited sequence.) When
		// Header values are duplicated by other fields in this struct (e.g.,
		// ContentLength, TransferEncoding, Trailer), the field values are
		// authoritative.
		//
		// Keys in the map are canonicalized (see CanonicalHeaderKey).
		Header: resp.Header,

		// ContentLength records the length of the associated content. The
		// value -1 indicates that the length is unknown. Unless Request.Method
		// is "HEAD", values >= 0 indicate that the given number of bytes may
		// be read from Body.
		ContentLength: resp.ContentLength,
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Errorf("Http Read Response Body err:%v", err)
		return &httpxResponse, fmt.Errorf("Http Read Response Body err:%v", err)
	}
	httpxResponse.Body = body
	if !Success2xx(&httpxResponse) {
		return &httpxResponse, errors.New("http请求失败，响应码为：" + strconv.Itoa(httpxResponse.StatusCode) + "，响应信息为：" + httpxResponse.Status)
	}
	return &httpxResponse, nil
}
func xToHttpxResponseJson(resp *http.Response, resultObjs ...interface{}) (*HttpxResponse, error) {
	pHttpxResponse, err := xToHttpxResponse(resp)
	if err != nil {
		return pHttpxResponse, err
	}
	if resultObjs == nil || len(resultObjs) <= 0 {
		return pHttpxResponse, err
	}
	if len(pHttpxResponse.Body) <= 0 {
		return pHttpxResponse, errors.New("Http Response Body is Empty,Can Not Unmarshal Response By JSON")
	}
	str := string(pHttpxResponse.Body)
	fmt.Println(str)
	// 只有第一个主对象析失败才会返回JSON解析错误，其他对象解析失败，不返回JSON解析错误
	var errResult error = nil
	resultLen := len(resultObjs)
	for i := 0; i < resultLen; i++ {
		resultObj := resultObjs[i]
		if nil == resultObj {
			continue
		}
		errTemp := json.Unmarshal(pHttpxResponse.Body, resultObj)
		if i == 0 {
			errResult = errTemp
		}
	}
	return pHttpxResponse, errResult

}

func xParseRealGetUrl(urlForGet string, urlQueryStr string) string {
	if len(urlQueryStr) == 0 {
		return urlForGet
	} else if len(urlForGet) == 0 {
		return urlQueryStr
	} else if strings.HasSuffix(urlForGet, "?") {
		return urlForGet + urlQueryStr
	} else if strings.HasSuffix(urlForGet, "&") {
		return urlForGet + urlQueryStr
	} else if strings.Contains(urlForGet, "?") {
		return urlForGet + "&" + urlQueryStr
	} else {
		return urlForGet + "?" + urlQueryStr
	}
}

func xParseRequestUrl(reqUrl string, reqParamStr string, reqQueryStr string) string {
	if len(reqUrl) <= 0 {
		if len(reqParamStr) <= 0 && len(reqQueryStr) <= 0 {
			return reqUrl
		} else if len(reqParamStr) <= 0 {
			return reqQueryStr
		} else if len(reqQueryStr) <= 0 {
			return reqParamStr
		} else {
			return reqParamStr + "?" + reqQueryStr
		}
	} else {
		lenUrl := len(reqUrl)
		indexWH := strings.LastIndex(reqUrl, "?")
		url1 := ""
		url2 := ""
		if indexWH < 0 {
			url1 = reqUrl
			url2 = ""
		} else if indexWH == 0 {
			url1 = ""
			url2 = reqUrl[1:]
		} else if indexWH < lenUrl-1 {
			url1 = reqUrl[0:indexWH]
			url2 = reqUrl[indexWH+1:]
		} else {
			url1 = reqUrl[0:indexWH]
			url2 = ""
		}
		if len(reqParamStr) > 0 {
			if strings.HasSuffix(url1, "/") {
				url1 = url1 + reqParamStr[1:]
			} else {
				url1 = url1 + reqParamStr
			}
		}
		if len(reqQueryStr) > 0 {
			if len(url2) > 0 {
				url2 = url2 + "&" + reqQueryStr
			} else {
				url2 = reqQueryStr
			}
		}
		if len(url2) > 0 {
			return url1 + "?" + url2
		} else {
			return url1
		}
	}
}

func xParseRequestMime(postStr string) string {
	if xIsXMLString(postStr) {
		index := strings.Index(postStr, "?>")
		if index > 0 {
			return "application/xml"
		} else {
			return "text/xml"
		}
	} else if xIsJsonString(postStr) {
		return "application/json"
	} else if xIsUrlString(postStr) {
		return "application/x-www-form-urlencoded"
	} else {
		return "text/plain"
	}
}

func xIsXMLString(postStr string) bool {
	if len(postStr) <= 0 {
		return false
	} else if strings.HasPrefix(postStr, "<") && strings.HasSuffix(postStr, ">") {
		return true
	} else {
		return false
	}
}

func xIsJsonString(postStr string) bool {
	if len(postStr) <= 0 {
		return false
	} else if strings.HasPrefix(postStr, "{") && strings.HasSuffix(postStr, "}") {
		return true
	} else if strings.HasPrefix(postStr, "[") && strings.HasSuffix(postStr, "]") {
		return true
	} else {
		return false
	}
}

func xIsUrlString(postStr string) bool {
	if len(postStr) <= 0 {
		return false
	}
	datas := strings.Split(postStr, "&")
	total := len(datas)
	count := 0
	for _, tmp := range datas {
		if len(tmp) <= 0 {
			continue
		}
		strings.LastIndex(tmp, "=")
		firstIndex := strings.Index(tmp, "=")
		if firstIndex <= 0 {
			continue
		}
		lastIndex := strings.LastIndex(tmp, "=")

		if firstIndex != lastIndex {
			continue
		}
		count++
	}
	return count == total
}
