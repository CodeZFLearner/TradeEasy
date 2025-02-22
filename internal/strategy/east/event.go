package east

import (
	"HelloGolang/internal/platform"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type EventStrategy struct{}

const Event = "https://np-listapi.eastmoney.com/sec/economiccalendar/findEventByDateRange?startDate=%s&eventType=001&pageNum=1&pageSize=15&req_trace="

func (s *EventStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	url := fmt.Sprintf(Event, time.Now().Format("2006-01-02"))
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (s *EventStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return platform.ParseEastMoneyResponseData(responseBody)
}
