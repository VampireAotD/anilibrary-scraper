package scraper

import (
	"fmt"
	"net/http"
	"net/url"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/scraper/client"
	"github.com/PuerkitoBio/goquery"
)

type Contract interface {
	Title(document *goquery.Document) string
	Status(document *goquery.Document) entity.Status
	Rating(document *goquery.Document) float32
	Episodes(document *goquery.Document) string
	Genres(document *goquery.Document) []string
	VoiceActing(document *goquery.Document) []string
	Image(document *goquery.Document) string
}

const (
	MinimalAnimeRating   float32 = 0
	MinimalAnimeEpisodes string  = "0 / ?"
)

type Scraper[I Contract] struct {
	url    string
	client client.Client
}

func New(url string, client client.Client) Scraper[Contract] {
	return Scraper[Contract]{url: url, client: client}
}

func (s Scraper[I]) Scrape(instance I) (entity.Anime, error) {
	response, err := s.client.Request(s.url)
	if err != nil {
		return entity.Anime{}, fmt.Errorf("sending request %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return entity.Anime{}, fmt.Errorf("bad status code: %d", response.StatusCode)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return entity.Anime{}, fmt.Errorf("creating document %v", err)
	}

	var anime entity.Anime

	parse, _ := url.Parse(s.url)

	anime.Image = fmt.Sprintf("%s://%s%s", parse.Scheme, parse.Host, instance.Image(document)) // TODO fix this
	anime.Title = instance.Title(document)
	anime.Status = instance.Status(document)
	anime.Rating = instance.Rating(document)
	anime.Episodes = instance.Episodes(document)
	anime.Genres = instance.Genres(document)
	anime.VoiceActing = instance.VoiceActing(document)

	return anime, nil
}
