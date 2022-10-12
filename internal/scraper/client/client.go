package client

import (
	"net/http"
	"time"

	"github.com/corpix/uarand"
)

type Client struct {
	base http.Client
}

func New(client http.Client) Client {
	return Client{base: client}
}

func DefaultClient() Client {
	t := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
	}

	client := http.Client{
		Transport: t,
		Timeout:   15 * time.Second,
	}

	return New(client)
}

func (c Client) DefaultHeaders(request *http.Request) *http.Request {
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("Origin", "https://www.google.com")
	request.Header.Set("User-Agent", uarand.GetRandom())

	return request
}

func (c Client) Request(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	c.DefaultHeaders(request)

	return c.base.Do(request)
}
