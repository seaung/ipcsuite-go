package http

import (
	"net/http"
	"net/url"
	"time"
)

func NewHttpClient(thread int, proxyUri string, timeout time.Duration) error {
	return nil
}

func SendRequest(request *http.Request, allowRedirect bool) error {
	return nil
}

func ParseRequest(request *http.Request) error {
	return nil
}

func ParseUri(uri *url.URL) {}

func ParseResponse(response *http.Response) error {
	return nil
}

func getResponseBody(response *http.Response) ([]byte, error) {
	return nil, nil
}
