package service

import (
	"context"

	"anilibrary-scraper/internal/entity"
)

//go:generate mockgen -source=contract.go -destination=./mocks.go -package=service

type (
	ScraperService interface {
		// Process method scraping all data from given url
		Process(ctx context.Context, url string) (*entity.Anime, error)
	}

	EventService interface {
		// Send method sends event data for given url
		Send(ctx context.Context, url string) error
	}
)
