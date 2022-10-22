package scraper

import (
	"context"
	"fmt"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/scraper"
	"go.opentelemetry.io/otel/trace"
)

var _ service.ScraperService = (*Service)(nil)

type Service struct {
	repository repository.AnimeRepository
}

func NewScraperService(repository repository.AnimeRepository) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) Process(ctx context.Context, dto dto.RequestDTO) (*entity.Anime, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().
		Tracer("ScraperService").
		Start(ctx, "Process")
	defer span.End()

	if dto.FromCache {
		anime, _ := s.repository.FindByUrl(ctx, dto.Url)
		if anime != nil {
			metrics.IncrCacheHitCounter()
			return anime, nil
		}
	}

	anime, err := scraper.Scrape(dto.Url)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("while scraping: %w", err)
	}

	if dto.FromCache {
		_ = s.repository.Create(ctx, dto.Url, anime)
	}

	return anime, nil
}
