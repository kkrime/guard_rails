package client

import (
	"net/http"
)

type httpClient struct {
}

func NewHttpCleint() HttpClient {
	return &httpClient{}
}

func (hc *httpClient) IsUrlReachable(url string) bool {
	resp, err := http.Get(url)

	if err != nil {
		return false
	}

	if resp.Status[:3] != "200" {
		return false
	}

	return true
}
