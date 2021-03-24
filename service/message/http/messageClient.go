package http

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type httpClient struct {
	client *fasthttp.Client
}

type Client interface {
	Get(host string, request interface{}) ([]byte, error)
}

func NewMessageHTTPClient() Client {
	return &httpClient{
		&fasthttp.Client{},
	}
}

func (httpClient *httpClient) Get(host string, request interface{}) ([]byte, error) {
	res := fasthttp.AcquireResponse()
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseResponse(res)
	defer fasthttp.ReleaseRequest(req)

	body, _ := jsoniter.Marshal(request)

	req.SetRequestURI(host)
	req.SetBody(body)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	err := httpClient.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return res.Body(), nil
}
