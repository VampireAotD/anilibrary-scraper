package client

import (
	"net/http"
	"time"
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
