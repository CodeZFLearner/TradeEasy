package east

import (
	"HelloGolang/internal/platform"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// PopularStrategy 处理 Popular 类型的请求
type PopularStrategy struct{}

const Popular = "https://emappdata.eastmoney.com/stockrank/getAllCurrentList"

func (s *PopularStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	url := Popular
	body := bytes.NewBuffer(nil)
	param := map[string]interface{}{"appId": "appId01", "globalId": "786e4c21-70dc-435a-93bb-38", "marketType": "", "pageNo": 1, "pageSize": 100}
	err := json.NewEncoder(body).Encode(param)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (s *PopularStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return platform.ParseEastMoneyResponseData(responseBody)
}

// SinaStrategy 处理 Sina 类型的请求

// 其他策略类可以类似地实现
