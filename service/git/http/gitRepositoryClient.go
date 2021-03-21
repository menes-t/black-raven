package http

import (
	"github.com/valyala/fasthttp"
)

type httpClient struct {
	client *fasthttp.Client
}

type Client interface {
	Get(url string) ([]byte, error)
}

func NewGitRepositoryHTTPClient() Client {
	return &httpClient{
		&fasthttp.Client{},
	}
}

func (httpClient *httpClient) Get(url string) ([]byte, error) {
	res := fasthttp.AcquireResponse()
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseResponse(res)
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	err := httpClient.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	result, err := res.BodyGunzip()
	if err != nil {
		return nil, err
	}

	return result, nil
}
