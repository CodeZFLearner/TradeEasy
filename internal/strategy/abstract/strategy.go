package abstract

import (
	"net/http"
	"net/url"
)

type RequestStrategy interface {
	BuildRequest(query url.Values) (*http.Request, error)
	ParseResponse(resp *http.Response) (map[string]interface{}, error)
}
