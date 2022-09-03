package repository

import (
	"context"

	"anilibrary-request-parser/internal/domain/entity"
)

type AnimeRepositoryInterface interface {
	FindByUrl(ctx context.Context, url string) (*entity.Anime, error)
	Create(ctx context.Context, key string, anime entity.Anime) error
}
