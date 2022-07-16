package scraper

import (
	"fmt"
	"net/http"

	"anilibrary-request-parser/app/internal/domain/contract"
	"anilibrary-request-parser/app/internal/domain/entity"
	"anilibrary-request-parser/app/pkg/logger"
	"github.com/PuerkitoBio/goquery"
)

type Scrapper struct {
	url      string
	instance contract.Scraper
	logger   logger.Logger
}

func New(url string, instance contract.Scraper, logger logger.Logger) *Scrapper {
	return &Scrapper{url: url, instance: instance, logger: logger}
}

func (s Scrapper) Process() (entity.Anime, error) {
	var anime entity.Anime

	response, err := http.Get(s.url)

	if err != nil {
		return anime, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return anime, fmt.Errorf("bad status code %d", response.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return anime, err
	}

	anime.Title = s.instance.GetTitle(document)
	anime.Status = s.instance.GetStatus(document)
	anime.Rating = s.instance.GetRating(document)
	anime.Episodes = s.instance.GetEpisodes(document)
	anime.Genres = s.instance.GetGenres(document)
	anime.VoiceActing = s.instance.GetVoiceActing(document)

	return anime, err
}
