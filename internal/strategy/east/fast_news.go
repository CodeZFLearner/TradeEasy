package east

import (
	"HelloGolang/internal/platform"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type FastNewsStrategy struct{}

const FastNews_URL = "https://np-seclist.eastmoney.com/sec/getFastNews"

func (s *FastNewsStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	pageSize := query.Get("pageSize")
	if pageSize==""{
		pageSize = "200"
	}
	param := map[string]interface{}{
		"biz":"sec_724",
		"client":"sec_android",
		"h24ColumnCode":"102",
		"order":2,
		"pageSize":pageSize,
		"timestamp":query.Get("timestamp"),
		"trace":"fd189d7e-02b7-456e-ac47-6ac93ee1484b",
	}
	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(param)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", FastNews_URL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (s *FastNewsStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return platform.ParseEastMoneyResponseData(responseBody)
}