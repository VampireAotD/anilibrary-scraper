package scraper

import (
	"context"
	"fmt"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/repository"
	"anilibrary-scraper/internal/repository/model"
	"anilibrary-scraper/internal/scraper"
	"anilibrary-scraper/internal/service"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("ScraperService").Start(ctx, "Process")
	defer span.End()

	span.AddEvent("Searching for scraped data in cache")

	anime, _ := s.repository.FindByURL(ctx, url)
	if anime != nil {
		metrics.IncrCacheHitCounter()
		span.SetStatus(codes.Ok, "anime was fetched from cache")
		return anime, nil
	}

	span.AddEvent("Scraping data from url")

	anime, err := s.scraper.ScrapeAnime(ctx, url)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("scraping : %w", err)
	}

	span.AddEvent("Creating cache")

	_ = s.repository.Create(ctx, model.Anime{
		URL:         url,
		Image:       anime.Image,
		Title:       anime.Title,
		Status:      anime.Status,
		Episodes:    anime.Episodes,
		Genres:      anime.Genres,
		VoiceActing: anime.VoiceActing,
		Synonyms:    anime.Synonyms,
		Rating:      anime.Rating,
	})

	return anime, nil
}
