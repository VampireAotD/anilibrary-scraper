package client

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"maps"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

var tlsClientRequestHeaders = http.Header{
	"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	"accept-encoding":           {"gzip"},
	"accept-language":           {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
	"cache-control":             {"max-age=0"},
	"sec-ch-ua":                 {`"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`},
	"sec-ch-ua-mobile":          {"?0"},
	"sec-ch-ua-platform":        {`"macOS"`},
	"sec-fetch-dest":            {"document"},
	"sec-fetch-mode":            {"navigate"},
	"sec-fetch-site":            {"none"},
	"sec-fetch-user":            {"?1"},
	"upgrade-insecure-requests": {"1"},
	"user-agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"},
	http.HeaderOrderKey: {
		"accept",
		"accept-encoding",
		"accept-language",
		"cache-control",
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-platform",
		"sec-fetch-dest",
		"sec-fetch-mode",
		"sec-fetch-site",
		"sec-fetch-user",
		"upgrade-insecure-requests",
		"user-agent",
	},
}

type TLSClient struct {
	client tlsclient.HttpClient
	pool   *sync.Pool
}

func NewTLSClient() (TLSClient, error) {
	client, err := tlsclient.NewHttpClient(
		tlsclient.NewLogger(),
		tlsclient.WithTimeoutSeconds(tlsclient.DefaultTimeoutSeconds),
		tlsclient.WithClientProfile(profiles.DefaultClientProfile),
		tlsclient.WithRandomTLSExtensionOrder(),
		tlsclient.WithInsecureSkipVerify(),
		tlsclient.WithCookieJar(tlsclient.NewCookieJar()),
	)
	if err != nil {
		return TLSClient{}, err
	}

	return TLSClient{
		client: client,
		pool: &sync.Pool{
			New: func() any {
				var builder strings.Builder
				builder.Grow(1024)
				return &builder
			},
		},
	}, nil
}

func (c TLSClient) Image(ctx context.Context, url string) (_ string, err error) {
	response, err := c.fetch(ctx, url)
	if err != nil {
		return "", err
	}

	defer func() {
		err = errors.Join(err, response.Body.Close())
	}()

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("invalid content type: %s, expected 'image/*'", contentType)
	}

	builder, _ := c.pool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		c.pool.Put(builder)
	}()

	encodedLen := base64.StdEncoding.EncodedLen(int(response.ContentLength))
	builder.Grow(5 + len(contentType) + 8 + encodedLen) // 5 - "data:", 8 - ";base64,"

	builder.WriteString("data:")
	builder.WriteString(contentType)
	builder.WriteString(";base64,")

	encoder := base64.NewEncoder(base64.StdEncoding, builder)
	defer func() {
		err = errors.Join(err, encoder.Close())
	}()

	_, err = io.Copy(encoder, response.Body)
	if err != nil {
		return "", fmt.Errorf("could not encode image: %w", err)
	}

	return builder.String(), nil
}

func (c TLSClient) HTML(ctx context.Context, url string) (_ *goquery.Document, err error) {
	response, err := c.fetch(ctx, url)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = errors.Join(err, response.Body.Close())
	}()

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("invalid content type: %s, expected 'text/html'", contentType)
	}

	return goquery.NewDocumentFromReader(response.Body)
}

func (c TLSClient) fetch(ctx context.Context, url string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	request.Header = maps.Clone(tlsClientRequestHeaders)
	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response: %d", response.StatusCode)
	}

	return response, nil
}
