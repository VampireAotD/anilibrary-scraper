package scraper

import (
	"context"
	"fmt"

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

func (s Service) Process(ctx context.Context, url string) (*entity.Anime, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("ScraperService").Start(ctx, "Process")
	defer span.End()

	anime, _ := s.repository.FindByURL(ctx, url)
	if anime != nil {
		metrics.IncrCacheHitCounter()
		return anime, nil
	}

	anime, err := scraper.Scrape(url)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("while scraping: %w", err)
	}

	_ = s.repository.Create(ctx, url, anime)

	return anime, nil
}
