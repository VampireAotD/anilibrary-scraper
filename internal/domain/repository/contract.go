package repository

import (
	"context"

	"anilibrary-scraper/internal/domain/entity"
)

//go:generate mockgen -source=contract.go -destination=./mock/repository_mock.go -package=mock

type (
	AnimeRepository interface {
		FindByUrl(ctx context.Context, url string) (*entity.Anime, error)
		Create(ctx context.Context, key string, anime entity.Anime) error
	}
)
