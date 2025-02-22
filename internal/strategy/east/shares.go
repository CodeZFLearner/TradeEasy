/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-01-05 19:38:13
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-01-20 00:47:55
 * @FilePath: \HelloGolang\internal\strategy\east\shares.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package east

import (
	"HelloGolang/internal/platform"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
)

const LastestChangeShares = "https://datacenter.eastmoney.com/securities/api/data/get" +
	"?type=RPT_EXECUTIVE_HOLD_DETAILS&sty=CHANGE_DATE,SECURITY_CODE,PERSON_NAME,CHANGE_SHARES,AVERAGE_PRICE,CHANGE_AMOUNT,CHANGE_RATIO,CHANGE_AFTER_HOLDNUM,HOLD_TYPE,DSE_PERSON_NAME,POSITION_NAME,PERSON_DSE_RELATION&sr=-1,-1" +
	"&st=CHANGE_DATE,CHANGE_AMOUNT&extraCols=f2,f12,f14&p=1&ps=1000"

type ShareStrategy struct{}

func (s *ShareStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	url := LastestChangeShares
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (s *ShareStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data, err := platform.ParseEastMoneyResponseData(responseBody)
	if err != nil {
		return data, err
	}
	if shares, ok := data["data"]; ok {

		if shareList, ok := shares.([]map[string]interface{}); ok {
			increGroup := map[string][]interface{}{}
			decreGroup := map[string][]interface{}{}
			for _, share := range shareList {
				if changeAmount, ok := share["CHANGE_AMOUNT"].(float64); ok && changeAmount > 0 {
					f12 := share["f12"].(string)
					increGroup[f12] = append(increGroup[f12], share)
				} else {
					f12 := share["f12"].(string)
					decreGroup[f12] = append(decreGroup[f12], share)
				}
			}
			increList := merageShares(increGroup)
			decreList := merageShares(decreGroup)

			if len(increList) > 10 {
				increList = increList[:10]
			}
			if len(decreList) > 10 {
				decreList = decreList[:10]
			}
			return map[string]interface{}{"increList": increList, "decreList": decreList}, nil

		}

	}
	return nil, nil

}

// 合并 一段区间 一只标的总增/减持，将所有标的 按照金额排序
func merageShares(group map[string][]interface{}) []interface{} {
	var result []interface{}
	for _, shares := range group {
		var totalAmount float64
		var weightedSum float64
		for _, share := range shares {
			if shareMap, ok := share.(map[string]interface{}); ok {
				changeAmount := shareMap["CHANGE_AMOUNT"].(float64)
				averagePrice := shareMap["AVERAGE_PRICE"].(float64)
				totalAmount += changeAmount
				weightedSum += float64(changeAmount) * averagePrice
			}
		}
		averagePrice := weightedSum / float64(totalAmount)
		obj := map[string]interface{}{
			"CHANGE_AMOUNT": math.Round(totalAmount*100) / 100,
			"AVERAGE_PRICE": math.Round(averagePrice*100) / 100,
			"f14":           shares[0].(map[string]interface{})["f14"],
			"f12":           shares[0].(map[string]interface{})["f12"],
			"f2":           shares[0].(map[string]interface{})["f2"],
		}
		result = append(result, obj)
	}

	sort.Slice(result, func(i, j int) bool {
		return math.Abs(result[i].(map[string]interface{})["CHANGE_AMOUNT"].(float64)) >
			math.Abs(result[j].(map[string]interface{})["CHANGE_AMOUNT"].(float64))
	})

	return result
}
