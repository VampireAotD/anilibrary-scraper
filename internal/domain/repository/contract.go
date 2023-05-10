package repository

import (
	"context"

	"anilibrary-scraper/internal/domain/entity"
)

//go:generate mockgen -source=contract.go -destination=./mocks.go -package=repository

type (
	AnimeRepository interface {
		// FindByURL method searching cached/stored anime and returns nil if not found
		FindByURL(ctx context.Context, url string) (*entity.Anime, error)

		// Create method creates anime cache/record
		Create(ctx context.Context, key string, anime *entity.Anime) error
	}

	EventRepository interface {
		// Send method sends event data
		Send(ctx context.Context, event *entity.Event) error
	}
)
