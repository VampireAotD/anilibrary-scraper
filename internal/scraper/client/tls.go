package client

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/corpix/uarand"
)

type TLSClient struct {
	client tls_client.HttpClient
}

func NewTLSClient(timeout int) TLSClient {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(timeout),
		tls_client.WithClientProfile(tls_client.Chrome_110),
		tls_client.WithCookieJar(jar),
		tls_client.WithInsecureSkipVerify(),
	}

	client, _ := tls_client.NewHttpClient(tls_client.NewLogger(), options...)

	return TLSClient{
		client: client,
	}
}

func (c TLSClient) fetch(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	request.Header = http.Header{
		"accept":          {"*/*"},
		"accept-language": {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"user-agent":      {uarand.GetRandom()},
		"origin":          {"https://www.google.com"},
		http.HeaderOrderKey: {
			"accept",
			"accept-language",
			"user-agent",
			"origin",
		},
	}

	response, err := c.client.Do(request)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response: %w:%d", err, response.StatusCode)
	}

	return response, nil
}

func (c TLSClient) FetchDocument(url string) (document *goquery.Document, err error) {
	response, err := c.fetch(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = errors.Join(err, response.Body.Close())
	}()

	return goquery.NewDocumentFromReader(response.Body)
}

func (c TLSClient) FetchResponseBody(url string) (body []byte, err error) {
	response, err := c.fetch(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = errors.Join(err, response.Body.Close())
	}()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
