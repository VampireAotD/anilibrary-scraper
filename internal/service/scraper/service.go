package scraper

import (
	"context"
	"fmt"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/repository/model"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

//go:generate mockgen -source=service.go -destination=./mocks.go -package=scraper
type (
	AnimeRepository interface {
		FindByURL(ctx context.Context, url string) (entity.Anime, error)
		Create(ctx context.Context, anime model.Anime) error
	}

	Scraper interface {
		ScrapeAnime(ctx context.Context, url string) (entity.Anime, error)
	}
)

type Service struct {
	repository AnimeRepository
	scraper    Scraper
}

func NewScraperService(repository AnimeRepository, scraper Scraper) Service {
	return Service{
		repository: repository,
		scraper:    scraper,
	}
}

func (s Service) Process(ctx context.Context, url string) (entity.Anime, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("ScraperService").Start(ctx, "Process")
	defer span.End()

	span.AddEvent("Searching for scraped data in cache")

	anime, _ := s.repository.FindByURL(ctx, url)
	if anime.Acceptable() {
		metrics.IncrCacheHitCounter()
		span.SetStatus(codes.Ok, "anime was fetched from cache")
		return anime, nil
	}

	span.AddEvent("Scraping data from url")

	anime, err := s.scraper.ScrapeAnime(ctx, url)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return entity.Anime{}, fmt.Errorf("scraping : %w", err)
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
