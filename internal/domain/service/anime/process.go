package anime

import (
	"context"
	"errors"
	"strings"

	"anilibrary-request-parser/internal/domain/dto"
	"anilibrary-request-parser/internal/domain/entity"
	"anilibrary-request-parser/internal/infrastructure/client"
	"anilibrary-request-parser/internal/infrastructure/scraper"
	"anilibrary-request-parser/internal/infrastructure/scraper/animego"
	"anilibrary-request-parser/internal/infrastructure/scraper/animevost"
	"anilibrary-request-parser/internal/infrastructure/scraper/contract"
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
		anime, _ := s.repository.FindByUrl(context.Background(), dto.Url)

		if anime != nil {
			return anime, nil
		}
	}

	anime, err := s.scraper.Process()

	if err != nil {
		return anime, err
	}

	if dto.FromCache {
		_ = s.repository.Create(context.Background(), dto.Url, *anime)
	}

	return anime, nil
}
