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
	"anilibrary-scraper/internal/scraper/client"
	"anilibrary-scraper/internal/scraper/parsers"
	"github.com/PuerkitoBio/goquery"
)

var ErrUndefinedScraper = errors.New("undefined scraper")

type Scraper[I parsers.Contract] struct {
	url    string
	client client.Client
}

func Scrape(url string) (*entity.Anime, error) {
	scraper := Scraper[parsers.Contract]{
		url:    url,
		client: client.DefaultClient(),
	}

	switch true {
	case strings.Contains(url, parsers.AnimeGoUrl):
		return scraper.process(parsers.NewAnimeGo())
	case strings.Contains(url, parsers.AnimeVostUrl):
		return scraper.process(parsers.NewAnimeVost())
	default:
		return nil, ErrUndefinedScraper
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

	anime := &entity.Anime{}
	var wg sync.WaitGroup

	wg.Add(7)

	go func(anime *entity.Anime) {
		defer wg.Done()

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
	}(anime)
	go func(anime *entity.Anime) {
		defer wg.Done()
		anime.Title = instance.Title(document)
	}(anime)
	go func(anime *entity.Anime) {
		defer wg.Done()
		anime.Status = instance.Status(document)
	}(anime)
	go func(anime *entity.Anime) {
		defer wg.Done()
		anime.Rating = instance.Rating(document)
	}(anime)
	go func(anime *entity.Anime) {
		defer wg.Done()
		anime.Episodes = instance.Episodes(document)
	}(anime)
	go func(anime *entity.Anime) {
		defer wg.Done()
		anime.Genres = instance.Genres(document)
	}(anime)
	go func(anime *entity.Anime) {
		defer wg.Done()
		anime.VoiceActing = instance.VoiceActing(document)
	}(anime)

	wg.Wait()

	return anime, nil
}
