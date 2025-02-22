/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-01-05 19:18:30
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-02-14 01:10:43
 * @FilePath: \HelloGolang\internal\platform\eastmoney.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package platform

import (
	"HelloGolang/pkg/common"
	"encoding/json"
	"log"
)

//func parseResponseData(query url.Values, responseBody []byte) (map[string]interface{}, error) {
//	if query.Get("category") == "Sina" {
//		return handleSina(responseBody)
//	}
//	return parseEastMoneyResponseData(responseBody)
//}

var EastMoneyFields = [...]string{
	"sc",            // 热榜代码
	"f14",           //stock 名称
	"f12",           //stock code
	"f2",            //stock 价格
	"f3",            // 涨幅
	"f24",           // 60 日涨跌幅
	"CHANGE_DATE",   //股权变动
	"CHANGE_AMOUNT", //金额(万)
	"POSITION_NAME", //职位
	"AVERAGE_PRICE", //股权交易均价
	"artTitle",      //资讯标题
	"artCode",       //资讯code
	"summary",       // 资讯描述
	"mediaName",     // 新闻来源
	"showTime",      // 展示时间
	"pinglunNum",
	"clickNum",
	"shareNum",
	"showTime", //重要事件开始时间
	"name",     //事件标题
	"code",     // 文章code地址
	//"updateTime",//
	"title",
	// "digest", // 摘要
	"realSort",
}

/*
*
 */
func ParseEastMoneyResponseData(responseBody []byte) (map[string]interface{}, error) {
	var resultList []map[string]interface{}
	var originMap map[string]interface{}
	// 将JSON字节切片反序列化为map
	err := json.Unmarshal(responseBody, &originMap)
	if err != nil {
		return nil, err
	}
	var jdataList []interface{}

	pmh := &common.ParseMapHelper{}

	if listSlice, err := pmh.GetTypedNestedValue(originMap, common.TypeSlice, "data"); err == nil {
		jdataList = listSlice.([]interface{})
	} else if listSlice, err := pmh.GetTypedNestedValue(originMap, common.TypeSlice, "data", "importantEventList"); err == nil {
		jdataList = listSlice.([]interface{})
	} else if listSlice, err := pmh.GetTypedNestedValue(originMap, common.TypeSlice, "data", "diff"); err == nil {
		jdataList = listSlice.([]interface{})
	} else if listSlice, err := pmh.GetTypedNestedValue(originMap, common.TypeSlice, "data", "items"); err == nil {
		jdataList = listSlice.([]interface{})
	} else if listSlice, err := pmh.GetTypedNestedValue(originMap, common.TypeSlice, "result", "data"); err == nil {
		jdataList = listSlice.([]interface{})
	} else {
		return nil, err
	}

	// 解析
	if jdataList != nil && len(jdataList) > 0 {
		for _, item := range jdataList {
			tmpMap := make(map[string]interface{})

			// 将item断言为map[string]interface{}类型
			if itemMap, ok := item.(map[string]interface{}); ok {
				for _, field := range EastMoneyFields {
					if v, e := itemMap[field]; e {
						tmpMap[field] = v
					}
				}
			}
			resultList = append(resultList, tmpMap)
		}
	} else {
		log.Printf("[东方财富]:响应不存在data ")
		return nil, nil
	}

	data := map[string]interface{}{
		"data": resultList,
	}
	return data, nil
}
