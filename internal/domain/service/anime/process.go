package anime

import (
	"context"
	"errors"
	"strings"
	"time"

	"anilibrary-request-parser/internal/adapter/client"
	"anilibrary-request-parser/internal/adapter/scraper"
	"anilibrary-request-parser/internal/adapter/scraper/animego"
	"anilibrary-request-parser/internal/adapter/scraper/animevost"
	"anilibrary-request-parser/internal/adapter/scraper/contract"
	"anilibrary-request-parser/internal/domain/dto"
	"anilibrary-request-parser/internal/domain/entity"
)

func (s *ScraperService) Process(dto dto.ParseDTO) (*entity.Anime, error) {
	base := scraper.New(dto.Url, client.DefaultClient())
	var instance contract.Scraper

	switch true {
	case strings.Contains(dto.Url, "animego.org"):
		instance = animego.New(base)
	case strings.Contains(dto.Url, "animevost.org"):
		instance = animevost.New(base)
	default:
		return nil, errors.New("undefined scraper")
	}

	base.Scraper = instance
	s.scraper = base

	if dto.FromCache {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		anime, _ := s.repository.FindByUrl(ctx, dto.Url)

		if anime != nil {
			return anime, nil
		}
	}

	anime, err := s.scraper.Process()

	if err != nil {
		return anime, err
	}

	if dto.FromCache {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = s.repository.Create(ctx, dto.Url, *anime)
	}

	return anime, nil
}
