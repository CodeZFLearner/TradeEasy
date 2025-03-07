/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-02-07 12:39:44
 * @LastEditors: zff 2059577798@qq.com
 * @LastEditTime: 2025-03-08 00:58:15
 * @FilePath: \HelloGolang\pkg\common\DeepSeek.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

const (
	ds_url = "https://api.deepseek.com/chat/completions"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Payload struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	MaxTokens   int16     `json:"max_tokens"`
	Temperature float32   `json:"temperature"`
}

type Usage struct {
	PromptTokens        int `json:"prompt_tokens"`
	CompletionTokens    int `json:"completion_tokens"`
	TotalTokens         int `json:"total_tokens"`
	PromptTokensDetails struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"prompt_tokens_details"`
	PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
}
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Logprobs     *string `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}
type Response struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

func QDeepSeek(messages []Message) (string, error) {
	// Get the API key from environment variable
	api_key := os.Getenv("DEEPSEEK_API_KEY")
	if api_key == "" {
		return "", errors.New("DEEPSEEK_API_KEY environment variable not set")
	}

	// Create a new POST request with the payload
	payload := Payload{
		Model:       "deepseek-chat",
		Messages:    messages,
		Stream:      false,
		MaxTokens:   4096,
		Temperature: 0.6,
	}

	// fmt.Println(payload)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", ds_url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api_key)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Parse the JSON response
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	// Extract the value for the specified key
	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "", nil
}
