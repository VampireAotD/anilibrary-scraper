package scraper

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/scraper/model"
	"anilibrary-scraper/internal/scraper/parsers"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrUnsupportedScraper = errors.New("unsupported scraper")
)

type Scraper struct {
	config Config
}

func New() Scraper {
	return Scraper{
		config: NewDefaultConfig(),
	}
}

func (s Scraper) ScrapeAnime(ctx context.Context, url string) (entity.Anime, error) {
	parser, err := s.resolveParser(url)
	if err != nil {
		return entity.Anime{}, fmt.Errorf("resolving parser %s: %w", url, err)
	}

	document, err := s.config.client.HTMLDocument(ctx, url)
	if err != nil {
		return entity.Anime{}, fmt.Errorf("scraping %s: %w", url, err)
	}

	return s.extractData(ctx, parser, document).MapToDomainEntity(), nil
}

func (s Scraper) resolveParser(url string) (parsers.Parser, error) {
	switch {
	case strings.Contains(url, parsers.AnimeGoURL):
		return parsers.NewAnimeGo(), nil
	case strings.Contains(url, parsers.AnimeVostURL):
		return parsers.NewAnimeVost(), nil
	default:
		return nil, ErrUnsupportedScraper
	}
}

func (s Scraper) extractData(ctx context.Context, parser parsers.Parser, document *goquery.Document) model.Anime {
	var (
		anime model.Anime
		wg    sync.WaitGroup
	)

	parseHTML := func(extractor func()) {
		defer s.config.panicHandler()
		defer wg.Done()
		extractor()
	}

	wg.Add(8)

	go parseHTML(func() {
		response, err := s.config.client.Response(ctx, parser.Image(document))
		if err != nil {
			return
		}

		anime.Image = fmt.Sprintf(
			"data:%s;base64,%s",
			http.DetectContentType(response),
			base64.StdEncoding.EncodeToString(response),
		)
	})

	go parseHTML(func() {
		anime.Title = parser.Title(document)
	})

	go parseHTML(func() {
		anime.Status = parser.Status(document)
	})

	go parseHTML(func() {
		anime.Rating = parser.Rating(document)
	})

	go parseHTML(func() {
		anime.Episodes = parser.Episodes(document)
	})

	go parseHTML(func() {
		anime.Genres = parser.Genres(document)
	})

	go parseHTML(func() {
		anime.VoiceActing = parser.VoiceActing(document)
	})

	go parseHTML(func() {
		anime.Synonyms = parser.Synonyms(document)
	})

	wg.Wait()

	return anime
}
