package repository

import (
	"context"
	"encoding/json"
	"time"

	"anilibrary-request-parser/app/internal/domain/entity"
	"github.com/go-redis/redis/v9"
)

const sevenDaysInHours string = "168h"

type AnimeRepository struct {
	client *redis.Client
}

func NewAnimeRepository(client *redis.Client) *AnimeRepository {
	return &AnimeRepository{client: client}
}

func (a *AnimeRepository) FindByUrl(ctx context.Context, url string) (*entity.Anime, error) {
	res, err := a.client.Get(ctx, url).Bytes()

	if err != nil {
		return nil, err
	}

	var anime entity.Anime

	json.Unmarshal(res, &anime)

	return &anime, nil
}

func (a *AnimeRepository) Create(ctx context.Context, key string, anime entity.Anime) error {
	expire, _ := time.ParseDuration(sevenDaysInHours)
	data, _ := anime.ToJson()

	return a.client.Set(ctx, key, data, expire).Err()
}
