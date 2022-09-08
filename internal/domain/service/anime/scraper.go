package anime

import (
	"anilibrary-request-parser/internal/adapter/scraper"
	"anilibrary-request-parser/internal/domain/repository"
)

type ScraperService struct {
	scraper    *scraper.Scraper
	repository repository.AnimeRepository
}

func NewScraperService(repository repository.AnimeRepository) *ScraperService {
	return &ScraperService{
		repository: repository,
	}
}
