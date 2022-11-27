package scraper

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/scraper/client"
	"anilibrary-scraper/internal/scraper/parsers"
	"anilibrary-scraper/internal/scraper/parsers/model"
	"github.com/PuerkitoBio/goquery"
)

var ErrUnsupportedScraper = errors.New("unsupported scraper")

// Scraper is a basically a factory of all parsers that can resolve parser for current url
// and scrape all data concurrently
type Scraper[I parsers.Contract] struct {
	url    string
	client client.Client
	wg     *sync.WaitGroup
	anime  *model.Anime
}

// Scrape method resolves parser for supported url and scrape all data
// throws error if url is not supported
func Scrape(url string) (*entity.Anime, error) {
	scraper := Scraper[parsers.Contract]{
		url:    url,
		client: client.DefaultClient(),
		wg:     &sync.WaitGroup{},
		anime:  &model.Anime{},
	}

	switch true {
	case strings.Contains(url, parsers.AnimeGoUrl):
		return scraper.process(parsers.NewAnimeGo())
	case strings.Contains(url, parsers.AnimeVostUrl):
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
	response, err := s.client.Request(s.url)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sending request %w, status code %d", err, response.StatusCode)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("creating document %v", err)
	}

	s.wg.Add(7)

	go s.parse(func(anime *model.Anime) {
		img, err := s.client.Request(instance.Image(document))
		if err != nil {
			return
		}
		defer img.Body.Close()

		var buff bytes.Buffer
		defer buff.Reset()

		_, err = buff.ReadFrom(img.Body)
		if err != nil {
			return
		}

		anime.Image = fmt.Sprintf(
			"data:%s;base64,%s",
			http.DetectContentType(buff.Bytes()),
			base64.StdEncoding.EncodeToString(buff.Bytes()),
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

	s.wg.Wait()

	return s.anime.ToEntity(), nil
}
