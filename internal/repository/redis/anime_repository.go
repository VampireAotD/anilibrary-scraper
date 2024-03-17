package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/repository/model"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
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
	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("AnimeRepository").Start(ctx, "FindByURL")
	defer span.End()

	bytes, err := a.client.Get(ctx, url).Bytes()
	if err != nil {
		span.RecordError(err)
		return entity.Anime{}, fmt.Errorf("while fetching from redis: %w", err)
	}

	var anime model.Anime
	if err = json.Unmarshal(bytes, &anime); err != nil {
		span.RecordError(err)
		return entity.Anime{}, fmt.Errorf("while converting from bytes: %w", err)
	}

	return anime.MapToDomainEntity(), nil
}

func (a AnimeRepository) Create(ctx context.Context, anime model.Anime) error {
	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("AnimeRepository").Start(ctx, "Create")
	defer span.End()

	// Error is not checked here because the only way error can occur
	// is when sevenDaysInHours will have invalid data
	expire, _ := time.ParseDuration(sevenDaysInHours)
	data, err := json.Marshal(anime)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("while converting to json: %w", err)
	}

	return a.client.Set(ctx, anime.URL, data, expire).Err()
}
