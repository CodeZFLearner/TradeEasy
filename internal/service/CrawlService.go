package service

import (
	strategy2 "HelloGolang/internal/strategy"
	"fmt"
	"net/http"
	"net/url"
)

func CrawlByCategory(query url.Values) (map[string]interface{}, error) {
	factory := &strategy2.StrategyFactory{}
	category := query.Get("category")
	strategy, err := factory.GetStrategy(category)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req, err := strategy.BuildRequest(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	result, err := strategy.ParseResponse(resp)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}
