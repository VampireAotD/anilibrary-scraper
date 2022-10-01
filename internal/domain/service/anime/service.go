package anime

import (
	"errors"
	"strings"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/scraper"
	"anilibrary-scraper/internal/scraper/client"
)

type ScraperService struct {
	repository repository.AnimeRepository
}

func NewScraperService(repository repository.AnimeRepository) ScraperService {
	return ScraperService{
		repository: repository,
	}
}

func (s ScraperService) scrape(url string) (entity.Anime, error) {
	switch true {
	case strings.Contains(url, "animego.org"):
		instance := scraper.New(url, client.DefaultClient())
		return instance.Scrape(scraper.NewAnimeGo())
	case strings.Contains(url, "animevost.org"):
		instance := scraper.New(url, client.DefaultClient())
		return instance.Scrape(scraper.NewAnimeVost())
	default:
		return entity.Anime{}, errors.New("undefined scraper")
	}
}
