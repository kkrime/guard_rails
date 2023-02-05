package client

type HttpClient interface {
	IsUrlReachable(url string) bool
}
