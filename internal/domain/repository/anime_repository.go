package repository

import (
	"context"

	"anilibrary-request-parser/internal/domain/entity"
)

//go:generate mockgen -source=anime_repository.go -destination=../../adapter/db/redis/repository/anime_repository_mock.go -package=repository

type AnimeRepository interface {
	FindByUrl(ctx context.Context, url string) (*entity.Anime, error)
	Create(ctx context.Context, key string, anime entity.Anime) error
}
