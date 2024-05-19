package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/infrastructure/repository/model"

	"github.com/redis/go-redis/v9"
)

const sevenDaysInHours string = "168h"

type AnimeRepository struct {
	client redis.UniversalClient
}

func NewAnimeRepository(client redis.UniversalClient) AnimeRepository {
	return AnimeRepository{
		client: client,
	}
}

func (a AnimeRepository) FindByURL(ctx context.Context, url string) (entity.Anime, error) {
	bytes, err := a.client.Get(ctx, url).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return entity.Anime{}, entity.ErrAnimeNotFound
		}

		return entity.Anime{}, fmt.Errorf("could not get anime from Redis: %w", err)
	}

	var anime model.Anime
	if err = json.Unmarshal(bytes, &anime); err != nil {
		return entity.Anime{}, fmt.Errorf("unmarshal 'anime' model from Redis data: %w", err)
	}

	return anime.MapToDomainEntity(), nil
}

func (a AnimeRepository) Create(ctx context.Context, anime model.Anime) error {
	// Error is not checked here because the only way error can occur
	// is when sevenDaysInHours will have invalid data
	expire, _ := time.ParseDuration(sevenDaysInHours)
	data, err := json.Marshal(anime)
	if err != nil {
		return fmt.Errorf("marshal 'anime' model for Redis: %w", err)
	}

	return a.client.Set(ctx, anime.URL, data, expire).Err()
}
