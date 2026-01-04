package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

// HttpRequest 函数用于发送 HTTP 请求
func HttpRequest(urlStr string, method string, headers map[string]string, params map[string]string,
	data any) (resp *http.Response, err error) {
	// 创建 URL 对象
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// 向 URL 添加查询参数
	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()

	// 将请求体数据（如果有）编码成JSON格式
	buf := new(bytes.Buffer)
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	// 创建 HTTP 请求对象
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 如果请求体存在，将 Content-Type 设置为 application/json
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送 HTTP 请求并获取响应
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
