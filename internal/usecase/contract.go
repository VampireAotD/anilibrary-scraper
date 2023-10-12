package usecase

import (
	"context"

	"anilibrary-scraper/internal/entity"
)

//go:generate mockgen -source=contract.go -destination=./mocks.go -package=usecase

type (
	ScraperUseCase interface {
		Scrape(ctx context.Context, url string) (*entity.Anime, error)
	}
)
