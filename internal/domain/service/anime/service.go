package anime

import (
	"errors"
	"strings"

	"anilibrary-request-parser/internal/domain/entity"
	"anilibrary-request-parser/internal/domain/repository"
	"anilibrary-request-parser/internal/scraper"
	"anilibrary-request-parser/internal/scraper/client"
	"github.com/PuerkitoBio/goquery"
)

type ScraperService struct {
	scraper    scraper.Contract
	repository repository.AnimeRepository
	client     client.Client
}

func NewScraperService(repository repository.AnimeRepository) *ScraperService {
	return &ScraperService{
		repository: repository,
		client:     client.DefaultClient(),
	}
}

func (s *ScraperService) composeScraper(url string) (scraper.Contract, error) {
	switch true {
	case strings.Contains(url, "animego.org"):
		return scraper.NewAnimeGo(), nil
	case strings.Contains(url, "animevost.org"):
		return scraper.NewAnimeVost(), nil
	default:
		return nil, errors.New("undefined scraper")
	}
}

func (s *ScraperService) scrape(document *goquery.Document) entity.Anime {
	var anime entity.Anime

	anime.Title = s.scraper.Title(document)
	anime.Status = s.scraper.Status(document)
	anime.Rating = s.scraper.Rating(document)
	anime.Episodes = s.scraper.Episodes(document)
	anime.Genres = s.scraper.Genres(document)
	anime.VoiceActing = s.scraper.VoiceActing(document)

	return anime
}
