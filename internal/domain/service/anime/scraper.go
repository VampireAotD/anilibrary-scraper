package anime

import (
	"anilibrary-request-parser/internal/domain/repository"
	"anilibrary-request-parser/internal/infrastructure/scraper"
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
