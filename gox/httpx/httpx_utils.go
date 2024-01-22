package httpx

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func ConfigLogger() {

}
func DoHttpPostJson(postUrl string, data interface{}, result interface{}) error {
	jsonData, err := ParseToJSONByte(data)
	if err != nil {
		// logger.Error("Http Post by application/json Marshal Request Data err:" + err.Error())
		return err
	}
	// logger.Debug("Http Post by application/json Message of request:" + string(jsonData))
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		// logger.Error("Http Post by application/json Build NewRequest err:" + err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Error("Http Post by application/json Execute Request err:" + err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// logger.Error("Http Post by application/json Read Response Body err:" + err.Error())
		return err
	}
	err = json.Unmarshal(body, result)

	if err != nil {
		// logger.Debug("Http Post by application/json Message of response:" + string(body))
		// logger.Error("Http Post by application/json Unmarshal Response Data err:" + err.Error())
		return err
	} else {
		// logger.Debug("Http Post by application/json Message of response:" + string(body))
		return nil
	}
}

func DoHttpPost(postUrl string, postContentType string, postStr string) ([]byte, error) {
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
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// logger.Error(fmt.Sprintf("Http Post by %v Read Response Body err: %v", contentType, err.Error()))
		return body, err
	} else {
		// logger.Debug(fmt.Sprintf("Http Post by %v Message of response: %v", contentType, string(body)))
		return body, err
	}
}

func DoHttpGet(urlOfGet string, data interface{}, result interface{}) error {
	urlData, err := ParseToUrlEncodeString(data)
	if err != nil {
		// logger.Error("Http Get Encode Request Data err:" + err.Error())
		return err
	}
	// logger.Debug("Http Get Encode Request Data ok:" + urlData)
	resp, err := http.Get(xParseRealGetUrl(urlOfGet, urlData))
	if err != nil {
		// logger.Error("Http Get Do Request err:" + err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// logger.Error("Http Get Read Response Body err:" + err.Error())
		return err
	}
	err = json.Unmarshal(body, result)

	if err != nil {
		// logger.Error("Http Get Unmarshal Response Data err:" + err.Error())
		return err
	} else {
		// logger.Debug("Http Get Message of response:" + string(body))
		return nil
	}
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
