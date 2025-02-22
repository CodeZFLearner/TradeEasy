package strategy

import (
	"HelloGolang/internal/platform"
	"HelloGolang/pkg/common"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const ETF = "https://stock.xueqiu.com/v5/stock/screener/fund/list.json?page=1&size=1080&order=desc&order_by=percent&type=18&parent_type=1"

type EtfStrategy struct{}

func (s *EtfStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	req, err := http.NewRequest("GET", ETF, nil)
	if err != nil {
		return nil, err
	}
	output := platform.XqId()
	parts := strings.Split(output, ";")
	req.AddCookie(&http.Cookie{
		Name:     "xqat",
		Value:    parts[1],
		Path:     "/",
		Domain:   ".xueqiu.com",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return req, nil
}

func (s *EtfStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return handleEtf(responseBody)
}

var fields = []string{
	"symbol", // 股票代码
	"amount", // 成交额
	"chg",    //
	"acc_unit_nav",
	"expiration_time",
	"rating",
	"etf_types",
	"type",
	"percent", //涨幅
	"tick_size",
	"volume",               // 成交量
	"current",              // 现价
	"current_year_percent", // 当年涨幅
	"etf_parent_type",
	"followers",
	"name",           // 名字
	"market_capital", // 净值
	"premium_rate",   // 溢价
	"lot_size",
	"unit_nav",
	"total_shares", // 总数
	"timestamp",
}

var etfFields = []string{
	"symbol",               // 股票代码
	"amount",               // 成交额
	"percent",              //涨幅
	"volume",               // 成交量
	"current",              // 现价
	"current_year_percent", // 当年涨幅
	"name",                 // 名字
	"market_capital",       // 净值
	"premium_rate",         // 溢价
}

func handleEtf(responseBody []byte) (map[string]interface{}, error) {
	var resultList []map[string]interface{}
	var originMap map[string]interface{}
	// 将JSON字节切片反序列化为map
	err := json.Unmarshal(responseBody, &originMap)
	if err != nil {
		return nil, err
	}
	var jdataList []interface{}

	pmh := &common.ParseMapHelper{}
	if listSlice, err := pmh.GetTypedNestedValue(originMap, common.TypeSlice, "data", "list"); err == nil {
		jdataList = listSlice.([]interface{})
	} else {
		return nil, err
	}

	if jdataList != nil && len(jdataList) > 0 {
		fmt.Println(len(jdataList))
		sort.Slice(jdataList, func(i, j int) bool {
			itemI := jdataList[i].(map[string]interface{})
			itemJ := jdataList[j].(map[string]interface{})
			return itemI["premium_rate"].(float64) < itemJ["premium_rate"].(float64)
		})
		for _, item := range jdataList {
			tmpMap := make(map[string]interface{})

			// 将item断言为map[string]interface{}类型
			if itemMap, ok := item.(map[string]interface{}); ok {
				for _, field := range etfFields {
					if v, e := itemMap[field]; e {
						tmpMap[field] = v
					}
				}
			}
			resultList = append(resultList, tmpMap)
		}

	} else {
		log.Printf("响应不存在data ")
		return nil, nil
	}
	data := map[string]interface{}{
		"data": resultList,
	}
	return data, nil
}
