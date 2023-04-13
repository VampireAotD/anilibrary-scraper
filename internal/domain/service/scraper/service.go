package scraper

import (
	"context"
	"fmt"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/scraper"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var _ service.ScraperService = (*Service)(nil)

type Service struct {
	repository repository.AnimeRepository
	scraper    scraper.Contract
}

func NewScraperService(repository repository.AnimeRepository, scraper scraper.Contract) Service {
	return Service{
		repository: repository,
		scraper:    scraper,
	}
}

func (s Service) Process(ctx context.Context, url string) (*entity.Anime, error) {
	ctx, span := otel.Tracer("ScraperService").Start(ctx, "Process")
	defer span.End()

	span.AddEvent("Searching for scraped data in cache")

	anime, _ := s.repository.FindByURL(ctx, url)
	if anime != nil {
		metrics.IncrCacheHitCounter()
		span.SetStatus(codes.Ok, "anime was fetched from cache")
		return anime, nil
	}

	span.AddEvent("Scraping data from url")

	anime, err := s.scraper.Scrape(ctx, url)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, fmt.Errorf("scraping : %w", err)
	}

	span.AddEvent("Creating cache")

	_ = s.repository.Create(ctx, url, anime)

	return anime, nil
}
