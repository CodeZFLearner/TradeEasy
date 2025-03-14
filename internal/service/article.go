/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-01-19 20:31:48
 * @LastEditors: zff 2059577798@qq.com
 * @LastEditTime: 2025-03-02 17:17:24
 * @FilePath: \HelloGolang\internal\service\day.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const EastMoneyArticle = "https://finance.eastmoney.com/a/%s.html"
const EastMoneyNews = "https://caifuhao.eastmoney.com/news/%s"

// 202501173300826597
func EastmoneyArticleService(artCode []string) (map[string]interface{}, error) {
	var configs []RequestConfig
	for _, code := range artCode {
		var config RequestConfig

		if len(code) == 18 {
			config = RequestConfig{
				ID:              code,
				URL:             fmt.Sprintf(EastMoneyArticle, code),
				Method:          "GET",
				Headers:         nil,
				Cookies:         nil,
				Body:            nil,
				ResponseHandler: eastMoneyResponseHandler,
			}
		} else {
			config = RequestConfig{
				ID:              code,
				URL:             fmt.Sprintf(EastMoneyNews, code),
				Method:          "GET",
				Headers:         nil,
				Cookies:         nil,
				Body:            nil,
				ResponseHandler: eastMoneyNew,
			}
		}
		configs = append(configs, config)
	}
	return Handle(configs), nil
}

func eastMoneyResponseHandler(response *http.Response) (interface{}, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	re := regexp.MustCompile(`文章主体-->([\n\w\W]{0,})文尾部其它信息`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("no match found")
	}
	content := matches[1]
	// Remove HTML tags
	re = regexp.MustCompile(`<[^>]*>`)
	content = re.ReplaceAllString(content, "")
	// Remove special characters like newlines and carriage returns
	re = regexp.MustCompile(`[\r\n\s]+`)
	content = re.ReplaceAllString(content, " ")
	return content, nil
}
func eastMoneyNew(response *http.Response) (interface{}, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	re := regexp.MustCompile(`articleTxt =(.*)`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("no match found")
	}
	content := matches[1]
	re = regexp.MustCompile(`<.*?>`)
	content = re.ReplaceAllString(content, "")
	return content, nil
}
