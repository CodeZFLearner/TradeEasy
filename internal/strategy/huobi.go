package strategy

import (
	"HelloGolang/internal/platform"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HuoBiStrategy struct{}

func (s *HuoBiStrategy) BuildRequest(query url.Values) (*http.Request, error) {
	req, err := http.NewRequest("GET", ETF, nil)
	if err != nil {
		return nil, err
	}
	output := platform.XqId()
	parts := strings.Split(output, ";")
	req.AddCookie(&http.Cookie{
		Name:     "xqat",
		Value:    parts[1],
		Path:     "/",
		Domain:   ".xueqiu.com",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return req, nil
}

func (s *HuoBiStrategy) ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return handleEtf(responseBody)
}
