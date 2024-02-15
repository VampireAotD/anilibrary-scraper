package scraper

import (
	"context"

	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/scraper/client"

	"github.com/PuerkitoBio/goquery"
)

// Client interface must be implemented by all clients that will be scraping
type Client interface {
	// HTMLDocument returns response body as *goquery.Document
	HTMLDocument(ctx context.Context, url string) (*goquery.Document, error)

	// Response returns response body as []byte
	Response(ctx context.Context, url string) ([]byte, error)
}

type Config struct {
	client       Client
	panicHandler func()
}

func NewDefaultConfig() Config {
	return Config{
		client: client.NewTLSClient(10),
		panicHandler: func() {
			if err := recover(); err != nil {
				metrics.IncrPanicCounter()
			}
		},
	}
}
