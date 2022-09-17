package client

import "net/http"

func (c Client) Request(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	c.DefaultHeaders(request)

	return c.base.Do(request)
}
