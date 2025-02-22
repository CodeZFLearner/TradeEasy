package strategy

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// %s f_560610,f_511990,f_511850
const Sina = "https://hq.sinajs.cn/?list=%s"

type SinaStrategy struct{}

func (s *SinaStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	codes := query.Get("codes")
	if len(codes) == 0 {
		return nil, fmt.Errorf("codes is empty")
	}
	url := fmt.Sprintf(Sina, codes)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://stock.finance.sina.com.cn/")
	return req, nil
}

func (s *SinaStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(transform.NewReader(resp.Body, simplifiedchinese.GB18030.NewDecoder()))
	if err != nil {
		return nil, err
	}
	return handleSina(responseBody)
}

func handleSina(responseBody []byte) (map[string]interface{}, error) {

	content := string(responseBody)
	re := regexp.MustCompile(".*(f_\\d+)=\"(.*)\"")
	matches := re.FindAllStringSubmatch(content, -1)

	var result []map[string]interface{}

	fmt.Println(content)

	for _, match := range matches {
		parts := strings.Split(match[2], ",")
		if len(parts) == 6 {
			result = append(result, map[string]interface{}{
				"code":    match[1],
				"name":    parts[0],
				"wprofit": parts[1],
				"lv":      parts[2],
				"date":    parts[4],
				"amount":  parts[5],
			})
		} else {
			fmt.Printf("ID: %s, Value: %s\n", match[1], match[2])
		}
	}

	data := map[string]interface{}{
		"data": result,
	}

	return data, nil
}
