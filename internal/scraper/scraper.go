package scraper

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/scraper/client"
	"anilibrary-scraper/internal/scraper/parsers"
	"anilibrary-scraper/internal/scraper/parsers/model"
)

var ErrUnsupportedScraper = errors.New("unsupported scraper")

// Scraper is a basically a factory of all parsers that can resolve parser for current url
// and scrape all data concurrently
type Scraper[I parsers.Contract] struct {
	client client.ChromeDp
	wg     *sync.WaitGroup
	anime  *model.Anime
	url    string
}

// Scrape method resolves parser for supported url and scrape all data
// throws error if url is not supported
func Scrape(url string) (*entity.Anime, error) {
	scraper := Scraper[parsers.Contract]{
		url:    url,
		client: client.NewChromeDpClient(),
		wg:     new(sync.WaitGroup),
		anime:  new(model.Anime),
	}

	switch true {
	case strings.Contains(url, parsers.AnimeGoURL):
		return scraper.process(parsers.NewAnimeGo())
	case strings.Contains(url, parsers.AnimeVostURL):
		return scraper.process(parsers.NewAnimeVost())
	default:
		return nil, ErrUnsupportedScraper
	}
}

type processor func(anime *model.Anime)

func (s Scraper[I]) parse(callback processor) {
	defer s.recover()
	defer s.wg.Done()

	callback(s.anime)
}

func (s Scraper[I]) recover() {
	if err := recover(); err != nil {
		metrics.IncrPanicCounter()
	}
}

func (s Scraper[I]) process(instance I) (*entity.Anime, error) {
	document, err := s.client.FetchDocument(60*time.Second, s.url)
	if err != nil {
		return nil, fmt.Errorf("scraper: %w", err)
	}

	s.wg.Add(8)

	go s.parse(func(anime *model.Anime) {
		responseBody, err := s.client.FetchResponseBody(60*time.Second, instance.Image(document))
		if err != nil {
			// TODO: handle error
			return
		}

		anime.Image = fmt.Sprintf(
			"data:%s;base64,%s",
			http.DetectContentType(responseBody),
			base64.StdEncoding.EncodeToString(responseBody),
		)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Title = instance.Title(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Status = instance.Status(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Rating = instance.Rating(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Episodes = instance.Episodes(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Genres = instance.Genres(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.VoiceActing = instance.VoiceActing(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Synonyms = instance.Synonyms(document)
	})

	s.wg.Wait()

	return s.anime.ToEntity(), nil
}
