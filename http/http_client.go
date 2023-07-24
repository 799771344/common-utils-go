package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Http struct {
	Request  *http.Request
	Response *http.Response
}

type NewRequest struct{}

func (r *NewRequest) GET(reqUrl string, headers map[string]string, params map[string]string) (*Http, error) {
	// 创建请求参数
	paramsNew := url.Values{}
	for k, v := range params {
		paramsNew.Add(k, v)
	}
	reqUrl = fmt.Sprintf("%s?%s", reqUrl, paramsNew.Encode())
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return &Http{Request: req}, nil
}

func (r *NewRequest) POST(reqUrl string, headers map[string]string, data map[string]interface{}) (*Http, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return &Http{Request: req}, nil
}

func (h *Http) ClientDo() (*Http, error) {
	client := &http.Client{}
	resp, err := client.Do(h.Request)
	if err != nil {
		return nil, err
	}
	return &Http{Response: resp}, nil
}

func (h *Http) GetResponseJson(v any) {
	respBytes, _ := ioutil.ReadAll(h.Response.Body)
	defer h.Response.Body.Close()
	json.Unmarshal(respBytes, v)
}

func (h *Http) GetResponseText() string {
	respBytes, _ := ioutil.ReadAll(h.Response.Body)
	defer h.Response.Body.Close()
	return string(respBytes)
}

func main() {
	req := &NewRequest{}
	//h, _ := req.POST("http://www.baidu.com", map[string]string{}, map[string]interface{}{})
	//resp, _ := h.ClientDo()
	//results := resp.GetResponseText()
	//fmt.Println(results)

	h1, _ := req.GET("http://www.baidu.com", map[string]string{}, map[string]string{})
	resp1, _ := h1.ClientDo()
	results1 := resp1.GetResponseText()
	fmt.Println(results1)
}
