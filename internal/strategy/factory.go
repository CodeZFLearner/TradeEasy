/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-01-05 18:59:33
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-02-11 23:07:33
 * @FilePath: \HelloGolang\internal\strategy\factory.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package strategy

import (
	"HelloGolang/internal/strategy/east"
	"fmt"
	"net/http"
	"net/url"
)

const (
	Org = "https://datacenter.eastmoney.com/securities/api/data/v1/get" +
		"?reportName=RPT_BILLBOARD_LISTDAILY" +
		"&columns=TRADE_DATE,OPERATEDEPT_NAME,OPERATEDEPT_NAME_NEW,ORG_NAME_ABBR,OPERATEDEPT_CODE,BUYER_APPEAR_NUM,SELLER_APPEAR_NUM,TOTAL_BUYAMT,TOTAL_SELLAMT,TOTAL_NETAMT" +
		"&filter=(TRADE_DATE='2024-12-18')&source=SECURITIES&pageNumber=1&pageSize=5"
	// %s:2024-12-24
	SeniorIncrease = "https://datacenter.eastmoney.com/securities/api/data/get" +
		"?type=RPT_EXECUTIVE_CHANGEHOLD_RANK&sty=SECURITY_CODE,AMOUNT,SHARES_NUM,CHANGE_AVG_PRICE" +
		"&extraCols=f14,f2&filter=(CHANGE_TYPE+=+\"01\")((DATE_TYPE+=+\"1\"))&p=1&ps=5&sr=-1&st=AMOUNT"
	SeniorReduce = "https://datacenter.eastmoney.com/securities/api/data/get" +
		"?type=RPT_EXECUTIVE_CHANGEHOLD_RANK&sty=SECURITY_CODE,AMOUNT,SHARES_NUM,CHANGE_AVG_PRICE" +
		"&extraCols=f14,f2&filter=(CHANGE_TYPE+=+\"02\")((DATE_TYPE+=+\"1\"))&p=1&ps=5&sr=1&st=AMOUNT"

	// %s 0.000759,0.002085,
)

type RequestStrategy interface {
	BuildRequest(query url.Values) (*http.Request, error)
	ParseResponse(resp *http.Response) (map[string]interface{}, error)
}
type StrategyFactory struct{}

func (f *StrategyFactory) GetStrategy(category string) (RequestStrategy, error) {
	switch category {
	case "Popular":
		return &east.PopularStrategy{}, nil
	case "PopularRc":
		return &east.PopularRcStrategy{}, nil
	case "HotNews":
		return &east.HotNewsStrategy{}, nil
	case "Sina":
		return &SinaStrategy{}, nil
	case "LastestChangeShares":
		return &east.ShareStrategy{}, nil
	case "StockDetail":
		return &east.StockStrategy{}, nil
	case "Event":
		return &east.EventStrategy{}, nil
	case "FastNews":
		return &east.FastNewsStrategy{}, nil
	case "Etf":
		return &EtfStrategy{}, nil
	case "lv60":
		return &east.TopChangeStrategy{}, nil

	// 添加其他策略
	default:
		return nil, fmt.Errorf("unknown category: %s", category)
	}
}
