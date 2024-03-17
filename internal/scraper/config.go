package scraper

import (
	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-playground/validator/v10"
)

// HTTPClient interface must be implemented by all clients that will be scraping
type HTTPClient interface {
	// HTMLDocument returns response body as *goquery.Document
	HTMLDocument(ctx context.Context, url string) (*goquery.Document, error)

	// Response returns response body as []byte
	Response(ctx context.Context, url string) ([]byte, error)
}

type config struct {
	client       HTTPClient
	validator    *validator.Validate
	panicHandler func()
}

type Option func(cfg *config)

func WithHTTPClient(client HTTPClient) Option {
	return func(cfg *config) {
		cfg.client = client
	}
}

func WithValidator(validate *validator.Validate) Option {
	return func(cfg *config) {
		cfg.validator = validate
	}
}

func WithPanicHandler(handler func()) Option {
	return func(cfg *config) {
		cfg.panicHandler = handler
	}
}
