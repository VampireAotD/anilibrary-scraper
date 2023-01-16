package service

import (
	"context"

	"anilibrary-scraper/internal/domain/entity"
)

//go:generate mockgen -source=contract.go -destination=./mocks.go -package=service

type (
	ScraperService interface {
		// Process method scraping all data from given url
		Process(ctx context.Context, url string) (*entity.Anime, error)
	}
)
