package http

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/seaung/ipcsuite-go/internal/protos"
)

var (
	client      *http.Client
	cloneClient *http.Client
)

func NewHttpClient(threadNumber int, proxyUri string, timeout time.Duration) error {
	dial := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 15 * time.Second,
	}

	transport := &http.Transport{
		DialContext:         dial.DialContext,
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: threadNumber * 2,
		IdleConnTimeout:     5 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 5 * time.Second,
	}

	if proxyUri != "" {
		u, err := url.Parse(proxyUri)
		if err != nil {
			return err
		}
		transport.Proxy = http.ProxyURL(u)
	}

	client = &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	cloneClient = &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	cloneClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return nil
}

func SendRequest(request *http.Request, allowRedirect bool) (*protos.Response, error) {
	var response *http.Response
	var err error

	if request.Body == nil || request.Body == http.NoBody {
		// this not request
	} else {
		request.Header.Set("Content-Length", strconv.Itoa(int(request.ContentLength)))
		if request.Header.Get("Content-Type") == "" {
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	if allowRedirect {
		response, err = client.Do(request)
	} else {
		response, err = cloneClient.Do(request)
	}

	if response != nil {
		defer response.Body.Close()
	}

	resp, err := ParseResponse(response)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ParseRequest(request *http.Request) (*protos.Request, error) {
	req := &protos.Request{}
	header := make(map[string]string)

	req.Method = request.Method
	req.Url = ParseUri(request.URL)

	for key := range request.Header {
		header[key] = request.Header.Get(key)
	}

	req.Headers = header
	req.ContentType = request.Header.Get("Content-Type")

	if request.Body == nil || request.Body == http.NoBody {
		// not process
	} else {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		req.Body = body
		request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	return req, nil
}

func ParseUri(uri *url.URL) *protos.UrlType {
	u := &protos.UrlType{}
	u.Scheme = uri.Scheme
	u.Domain = uri.Hostname()
	u.Host = uri.Host
	u.Port = uri.Port()
	u.Path = uri.EscapedPath()
	u.Query = uri.RawQuery
	u.Fragment = uri.Fragment
	return u
}

func ParseResponse(response *http.Response) (*protos.Response, error) {
	var resp protos.Response
	header := make(map[string]string)

	resp.Status = int32(response.StatusCode)
	resp.Url = ParseUri(response.Request.URL)

	for key := range response.Header {
		header[key] = response.Header.Get(key)
	}

	resp.Headers = header
	resp.ContentType = response.Header.Get("Content-Type")

	body, err := getResponseBody(response)
	if err != nil {
		return nil, err
	}

	resp.Body = body

	return &resp, nil
}

func getResponseBody(response *http.Response) ([]byte, error) {
	var body []byte
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, _ := gzip.NewReader(response.Body)
		defer reader.Close()

		for {
			buffer := make([]byte, 1024)
			nbyte, err := reader.Read(buffer)
			if err != nil && err != io.EOF {
				return nil, err
			}

			if nbyte == 0 {
				break
			}
			body = append(body, buffer...)
		}
	} else {
		rawContent, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		body = rawContent
	}
	return body, nil
}
