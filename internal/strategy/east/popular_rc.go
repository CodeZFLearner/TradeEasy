package east

import (
	"HelloGolang/internal/platform"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type PopularRcStrategy struct{}

const PopularRc = "https://emappdata.eastmoney.com/stockrank/getAllHisRcList"

func (s *PopularRcStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	url := PopularRc
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

func (s *PopularRcStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return platform.ParseEastMoneyResponseData(responseBody)
}
