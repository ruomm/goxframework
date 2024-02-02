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
	"strings"
)

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
		fmt.Errorf("Http Read Response Body err:" + err.Error())
		return &httpxResponse, err
	}
	httpxResponse.Body = body
	return &httpxResponse, nil
}
func xToHttpxResponseJson(resp *http.Response, result interface{}) (*HttpxResponse, error) {
	pHttpxResponse, err := xToHttpxResponse(resp)
	if err != nil {
		return pHttpxResponse, err
	}
	if len(pHttpxResponse.Body) <= 0 {
		return pHttpxResponse, errors.New("Http Response Body is Empty,Can Not Unmarshal Response By JSON")
	}
	err = json.Unmarshal(pHttpxResponse.Body, result)
	return pHttpxResponse, err

}

func xParseRealGetUrl(urlForGet string, urlParams string) string {
	if len(urlParams) == 0 {
		return urlForGet
	} else if len(urlForGet) == 0 {
		return urlParams
	} else if strings.HasSuffix(urlForGet, "?") {
		return urlForGet + urlParams
	} else if strings.HasSuffix(urlForGet, "&") {
		return urlForGet + urlParams
	} else if strings.Contains(urlForGet, "?") {
		return urlForGet + "&" + urlParams
	} else {
		return urlForGet + "?" + urlParams
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
