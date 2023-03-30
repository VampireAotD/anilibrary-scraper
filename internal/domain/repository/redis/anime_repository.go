package redis

import (
	"context"
	"fmt"
	"time"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

const sevenDaysInHours string = "168h"

var _ repository.AnimeRepository = (*AnimeRepository)(nil)

type AnimeRepository struct {
	client *redis.Client
}

func NewAnimeRepository(client *redis.Client) AnimeRepository {
	return AnimeRepository{
		client: client,
	}
}

func (a AnimeRepository) FindByURL(ctx context.Context, url string) (*entity.Anime, error) {
	_, span := otel.Tracer("AnimeRepository").Start(ctx, "FindByURL")
	defer span.End()

	res, err := a.client.Get(ctx, url).Bytes()
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("while fetching from redis: %w", err)
	}

	var anime entity.Anime
	unmarshalled, err := anime.FromJSON(res)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("while converting from bytes: %w", err)
	}

	return unmarshalled, nil
}

func (a AnimeRepository) Create(ctx context.Context, key string, entity *entity.Anime) error {
	expire, _ := time.ParseDuration(sevenDaysInHours)
	data, err := entity.ToJSON()
	if err != nil {
		return fmt.Errorf("while converting to json: %w", err)
	}

	return a.client.Set(ctx, key, data, expire).Err()
}
