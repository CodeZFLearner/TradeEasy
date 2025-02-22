package service

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type ResponseHandler func(*http.Response) (interface{}, error)
type RequestConfig struct {
	ID              string
	URL             string
	Method          string
	Headers         map[string]string
	Cookies         []*http.Cookie
	Body            io.Reader
	ResponseHandler ResponseHandler
}

func fetchData(config RequestConfig, wg *sync.WaitGroup, results chan<- map[string]interface{}) {
	defer wg.Done()

	req, err := http.NewRequest(config.Method, config.URL, config.Body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")

	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	for _, cookie := range config.Cookies {
		req.AddCookie(cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()
	data, err := config.ResponseHandler(resp)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return
	}
	// Assuming the response is a JSON object
	var result map[string]interface{}

	result = map[string]interface{}{
		"id":   config.ID,
		"data": data,
	}

	results <- result
}

func Handle(configs []RequestConfig) map[string]interface{} {
	var wg sync.WaitGroup

	results := make(chan map[string]interface{}, 10)
	finalResults := make(map[string]interface{})

	for _, config := range configs {
		wg.Add(1)
		go fetchData(config, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		finalResults[result["id"].(string)] = result["data"]
	}

	return finalResults
}
