package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/corpix/uarand"
)

const defaultTimeoutInSeconds int = 10

type TLSClient struct {
	client tlsclient.HttpClient
}

func NewTLSClient() TLSClient {
	jar := tlsclient.NewCookieJar()
	options := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(defaultTimeoutInSeconds),
		tlsclient.WithClientProfile(profiles.Chrome_117),
		tlsclient.WithCookieJar(jar),
		tlsclient.WithInsecureSkipVerify(),
	}

	// Error is ignored because method validateConfig in tlsclient package always returns nil
	client, _ := tlsclient.NewHttpClient(tlsclient.NewLogger(), options...)

	return TLSClient{
		client: client,
	}
}

func (c TLSClient) HTMLDocument(ctx context.Context, url string) (*goquery.Document, error) {
	response, err := c.fetch(ctx, url)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = errors.Join(err, response.Body.Close())
	}()

	return goquery.NewDocumentFromReader(response.Body)
}

func (c TLSClient) Response(ctx context.Context, url string) ([]byte, error) {
	response, err := c.fetch(ctx, url)
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

func (c TLSClient) fetch(ctx context.Context, url string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
