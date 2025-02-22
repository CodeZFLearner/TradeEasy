package east

import (
	"HelloGolang/internal/platform"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const StockDetail = "https://push2.eastmoney.com/api/qt/ulist.np/get?fltt=2&fields=f14,f148,f3,f12,f2,f13,f29&secids=%s"

type StockStrategy struct{}

func (s *StockStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	codes := query.Get("codes")
	if len(codes) == 0 {
		return nil, nil
	}
	var parseCodes []string
	for _, code := range strings.Split(codes, ",") {
		if code[0] == '6' {
			parseCodes = append(parseCodes, "1."+code)
		} else {
			parseCodes = append(parseCodes, "0."+code)
		}
	}
	url := fmt.Sprintf(StockDetail, strings.Join(parseCodes, ","))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (s *StockStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return platform.ParseEastMoneyResponseData(responseBody)
}
