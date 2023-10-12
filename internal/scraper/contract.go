package scraper

import (
	"context"
	"errors"

	"anilibrary-scraper/internal/entity"
)

var ErrUnsupportedScraper = errors.New("unsupported scraper")

//go:generate mockgen -source=contract.go -destination=./mocks.go -package=scraper

type Contract interface {
	// Scrape method resolves parser for supported url and scrape all data
	Scrape(ctx context.Context, url string) (*entity.Anime, error)
}
