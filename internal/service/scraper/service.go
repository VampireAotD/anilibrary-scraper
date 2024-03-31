package scraper

import (
	"context"
	"errors"
	"fmt"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/repository/model"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

//go:generate mockgen -source=service.go -destination=./mocks.go -package=scraper
type (
	Scraper interface {
		ScrapeAnime(ctx context.Context, url string) (entity.Anime, error)
	}

	AnimeRepository interface {
		FindByURL(ctx context.Context, url string) (entity.Anime, error)
		Create(ctx context.Context, anime model.Anime) error
	}
)

type Service struct {
	scraper         Scraper
	cacheRepository AnimeRepository
}

func NewScraperService(scraper Scraper, repository AnimeRepository) Service {
	return Service{
		scraper:         scraper,
		cacheRepository: repository,
	}
}

func (s Service) Process(ctx context.Context, url string) (entity.Anime, error) {
	ctx, span := otel.Tracer("ScraperService").Start(ctx, "Process")
	defer span.End()

	span.AddEvent("Searching for anime in cache")

	anime, err := s.cacheRepository.FindByURL(ctx, url)
	if err != nil {
		if errors.Is(err, entity.ErrAnimeNotFound) {
			// No need to record error if anime is not found in cache
			metrics.IncrCacheMissCounter()
		} else {
			span.SetStatus(codes.Error, "failed to get anime from cache")
			span.RecordError(err)
		}
	} else {
		metrics.IncrCacheHitCounter()
		span.SetStatus(codes.Ok, "anime has been fetched from cache")
		return anime, nil
	}

	span.AddEvent("Scraping data from url")

	anime, err = s.scraper.ScrapeAnime(ctx, url)
	if err != nil {
		span.SetStatus(codes.Error, "failed to scrape anime")
		span.RecordError(err)
		return entity.Anime{}, fmt.Errorf("scraping anime: %w", err)
	}

	span.AddEvent("Creating anime record in cache")

	err = s.cacheRepository.Create(ctx, model.Anime{
		URL:         url,
		Image:       anime.Image,
		Title:       anime.Title,
		Status:      anime.Status,
		Episodes:    anime.Episodes,
		Genres:      anime.Genres,
		VoiceActing: anime.VoiceActing,
		Synonyms:    anime.Synonyms,
		Rating:      anime.Rating,
		Year:        anime.Year,
		Type:        anime.Type,
	})
	if err != nil {
		span.SetStatus(codes.Error, "failed to create anime record in cache")
		span.RecordError(err)
	}

	return anime, nil
}
