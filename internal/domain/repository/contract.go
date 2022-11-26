package repository

import (
	"context"

	"anilibrary-scraper/internal/domain/entity"
)

//go:generate mockgen -source=contract.go -destination=./mocks/repository_mock.go -package=mocks

type (
	AnimeRepository interface {
		// FindByUrl method searching cached/stored anime and returns nil if not found
		FindByUrl(ctx context.Context, url string) (*entity.Anime, error)

		// Create method creates anime cache/record
		Create(ctx context.Context, key string, entity *entity.Anime) error
	}
)
