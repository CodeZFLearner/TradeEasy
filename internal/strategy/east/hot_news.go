package east

import (
	"HelloGolang/internal/platform"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const HotNews = "https://np-listapi.eastmoney.com/sec/hotNews?biz=sec_hotnews&client=sec_all&needInteractData=0&req_trace=11&size=10"

type HotNewsStrategy struct{}

func (s *HotNewsStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	url := HotNews
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (s *HotNewsStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return platform.ParseEastMoneyResponseData(responseBody)
}
