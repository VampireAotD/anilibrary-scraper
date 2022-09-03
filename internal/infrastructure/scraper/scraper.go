package scraper

import (
	"fmt"
	"net/http"

	"anilibrary-request-parser/internal/domain/entity"
	"anilibrary-request-parser/internal/infrastructure/client"
	"anilibrary-request-parser/internal/infrastructure/scraper/contract"
	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	contract.Scraper

	Url    string
	Client *client.Client
}

const (
	MinimalAnimeRating   float32 = 0
	MinimalAnimeEpisodes string  = "0 / ?"
)

func New(url string, client *client.Client) *Scraper {
	return &Scraper{Url: url, Client: client}
}

func (s Scraper) Process() (*entity.Anime, error) {
	response, err := s.Client.Request(s.Url)

	if err != nil {
		return nil, fmt.Errorf("sending request %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code %d", response.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return nil, fmt.Errorf("creating document %v", err)
	}

	var anime entity.Anime

	anime.Title = s.Title(document)
	anime.Status = s.Status(document)
	anime.Rating = s.Rating(document)
	anime.Episodes = s.Episodes(document)
	anime.Genres = s.Genres(document)
	anime.VoiceActing = s.VoiceActing(document)

	return &anime, err
}
