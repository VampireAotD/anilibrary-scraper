package anime

import (
	"context"
	"fmt"
	"time"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/entity"
)

func (s ScraperService) Process(dto dto.RequestDTO) (*entity.Anime, error) {
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
