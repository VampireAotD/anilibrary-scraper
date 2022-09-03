package anime

import (
	"anilibrary-request-parser/internal/domain/repository"
	"anilibrary-request-parser/internal/infrastructure/scraper"
)

type ScraperService struct {
	scraper    *scraper.Scraper
	repository repository.AnimeRepositoryInterface
}

func NewScraperService(repository repository.AnimeRepositoryInterface) *ScraperService {
	return &ScraperService{
		repository: repository,
	}
}
