package client

import (
	"net/http"
	"time"
)

type Client struct {
	http.Client // better to do like this or `client http.Client?`
}

func New(client http.Client) *Client {
	return &Client{Client: client}
}

func DefaultClient() *Client {
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
