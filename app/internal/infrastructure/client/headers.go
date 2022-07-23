package client

import "net/http"

func (c Client) DefaultHeaders(request *http.Request) *http.Request {
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("Origin", "https://www.google.com")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")

	return request
}
