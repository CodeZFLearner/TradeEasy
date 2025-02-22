/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-02-01 17:32:14
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-02-01 18:04:58
 * @FilePath: \HelloGolang\internal\strategy\east\top_change.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package east

import (
	"HelloGolang/internal/platform"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TopChangeStrategy struct{}

const prefix_url = "https://push2.eastmoney.com"
const fields = "f12,f13,f14,f1,f2,f4,f3,f152,f5,f6,f7,f15,f18,f16,f17,f10,f8,f9,f24"
const sort_field = "f24"
const incre = "1"
const pn = "1"
const pz = "20"
const TopChange = "%s/api/qt/clist/get?fields=%s&fid=%s&cb&np=1&fltt=2&pn=%s&pz=%s&po=%s&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23,m:0+t:81+s:2048"

func (s *TopChangeStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	po := query.Get("incre")
	if len(po)==0{
		po = incre
	}
	url := fmt.Sprintf(TopChange, prefix_url, fields, sort_field, pn, pz, po)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (s *TopChangeStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return platform.ParseEastMoneyResponseData(responseBody)
}
