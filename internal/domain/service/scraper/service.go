package scraper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/scraper"
	"anilibrary-scraper/internal/scraper/client"
)

type Service struct {
	repository repository.AnimeRepository
}

func NewScraperService(repository repository.AnimeRepository) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) scrape(url string) (*entity.Anime, error) {
	switch true {
	case strings.Contains(url, "animego.org"):
		instance := scraper.New(url, client.DefaultClient())
		return instance.Scrape(scraper.NewAnimeGo())
	case strings.Contains(url, "animevost.org"):
		instance := scraper.New(url, client.DefaultClient())
		return instance.Scrape(scraper.NewAnimeVost())
	default:
		return nil, errors.New("undefined scraper")
	}
}

func (s Service) Process(dto dto.RequestDTO) (*entity.Anime, error) {
	if dto.FromCache {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		anime, _ := s.repository.FindByUrl(ctx, dto.Url)
		if anime != nil {
			return anime, nil
		}
	}

	anime, err := s.scrape(dto.Url)
	if err != nil {
		return nil, fmt.Errorf("while scraping: %w", err)
	}

	if dto.FromCache {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = s.repository.Create(ctx, dto.Url, *anime)
	}

	return anime, nil
}
