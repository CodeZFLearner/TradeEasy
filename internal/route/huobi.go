package route

import (
	"HelloGolang/internal/platform"
	"HelloGolang/internal/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var etfs = map[string]string{
	"交易货币ETF":     "SH511690",
	"场内货币ETF":     "SH511700",
	"财富宝ETF":       "SH511850",
	"货币ETF建信添益": "SH511660",
	"理财金货币ETF":   "SH511810",
	"华宝添益ETF":     "SH511990",
	"招商快线ETF":     "SZ159003",
	"货币ETF":         "SZ159001",
}

func Huobi(w http.ResponseWriter, r *http.Request) {
	var configs []service.RequestConfig
	output := platform.XqId()
	parts := strings.Split(output, ";")
	for name, symbol := range etfs {
		config := service.RequestConfig{
			ID:      name,
			URL:     fmt.Sprintf("https://stock.xueqiu.com/v5/stock/realtime/pankou.json?symbol=%s", symbol),
			Method:  "GET",
			Headers: map[string]string{"Content-Type": "application/json"},
			Cookies: []*http.Cookie{
				{Name: "xqat", Value: parts[1]},
			},
			Body:            nil,
			ResponseHandler: convertResponseToMap,
		}
		configs = append(configs, config)

		url := fmt.Sprintf("https://stock.xueqiu.com/v5/stock/chart/kline.json?symbol=%s&begin=%d000&period=60m&type=before&count=-284&indicator=kline", symbol, time.Now().Unix())
		fmt.Println(url)

		config = service.RequestConfig{
			ID:      fmt.Sprintf("%s%s", name, ".k"),
			URL:     url,
			Method:  "GET",
			Headers: map[string]string{"Content-Type": "application/json"},
			Cookies: []*http.Cookie{
				{Name: "u", Value: parts[0]},
				{Name: "xqat", Value: parts[1]},
			},
			Body:            nil,
			ResponseHandler: convertResponseToMap,
		}
		configs = append(configs, config)
	}
	result := service.Handle(configs)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func convertResponseToMap(response *http.Response) (interface{}, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// Parse the body into a map
	var result map[string]interface{}
	fmt.Println(string(body))

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
