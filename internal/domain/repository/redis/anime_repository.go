package redis

import (
	"context"
	"time"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository"
	"github.com/go-redis/redis/v9"
)

const sevenDaysInHours string = "168h"

var _ repository.AnimeRepository = (*AnimeRepository)(nil)

type AnimeRepository struct {
	client *redis.Client
}

func NewAnimeRepository(client *redis.Client) repository.AnimeRepository {
	return AnimeRepository{client: client}
}

func (a AnimeRepository) FindByUrl(ctx context.Context, url string) (*entity.Anime, error) {
	res, err := a.client.Get(ctx, url).Bytes()
	if err != nil {
		return nil, err
	}

	var anime entity.Anime
	return anime.FromJSON(res)
}

func (a AnimeRepository) Create(ctx context.Context, key string, anime entity.Anime) error {
	expire, _ := time.ParseDuration(sevenDaysInHours)
	data, _ := anime.ToJSON()

	return a.client.Set(ctx, key, data, expire).Err()
}
