package scraper

import (
	"context"
	"fmt"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/scraper"
)

type Service struct {
	repository repository.AnimeRepository
}

func NewScraperService(repository repository.AnimeRepository) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) Process(dto dto.RequestDTO) (*entity.Anime, error) {
	if dto.FromCache {
		anime, _ := s.repository.FindByUrl(context.Background(), dto.Url)
		if anime != nil {
			return anime, nil
		}
	}

	anime, err := scraper.Scrape(dto.Url)
	if err != nil {
		return nil, fmt.Errorf("while scraping: %w", err)
	}

	if dto.FromCache {
		_ = s.repository.Create(context.Background(), dto.Url, anime)
	}

	return anime, nil
}
