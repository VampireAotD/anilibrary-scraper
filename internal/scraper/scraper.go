package scraper

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/scraper/model"
	"anilibrary-scraper/internal/scraper/parsers"
)

var (
	ErrUnsupportedScraper = errors.New("unsupported scraper")
)

type Parser interface {
	// Title method scraping anime title and returns empty string if none found
	Title() string

	// Status method scraping current anime status
	Status() model.Status

	// Rating method scraping current anime rating and returns parsers.MinimalAnimeRating if none found
	Rating() float32

	// Episodes method scraping amount of anime episodes and returns parsers.MinimalAnimeEpisodes if none found
	Episodes() string

	// Genres method scraping all anime genres
	Genres() []string

	// VoiceActing method scraping all anime voice acting
	VoiceActing() []string

	// Synonyms method scraping all similar anime names
	Synonyms() []string

	// ImageURL method scraping image url returns empty string if none found
	ImageURL() string
}

type Scraper struct {
	config Config
}

func New() Scraper {
	return Scraper{
		config: NewDefaultConfig(),
	}
}

func (s Scraper) ScrapeAnime(ctx context.Context, url string) (entity.Anime, error) {
	parser, err := s.resolveParser(ctx, url)
	if err != nil {
		return entity.Anime{}, fmt.Errorf("resolving parser %s: %w", url, err)
	}

	return s.extractData(ctx, parser).MapToDomainEntity(), nil
}

func (s Scraper) resolveParser(ctx context.Context, url string) (Parser, error) {
	switch {
	case strings.Contains(url, parsers.AnimeGoURL):
		document, err := s.config.client.HTMLDocument(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("scraping %s: %w", url, err)
		}

		return parsers.NewAnimeGo(document), nil
	case strings.Contains(url, parsers.AnimeVostURL):
		document, err := s.config.client.HTMLDocument(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("scraping %s: %w", url, err)
		}

		return parsers.NewAnimeVost(document), nil
	default:
		return nil, ErrUnsupportedScraper
	}
}

func (s Scraper) extractData(ctx context.Context, parser Parser) model.Anime {
	var anime model.Anime

	imageCh := make(chan struct{})
	parseHTML := func(extractor func()) {
		defer s.config.panicHandler()
		extractor()
	}

	go parseHTML(func() {
		defer close(imageCh)

		if url := parser.ImageURL(); url != "" {
			response, err := s.config.client.Response(ctx, url)
			if err != nil {
				return
			}

			anime.Image = fmt.Sprintf(
				"data:%s;base64,%s",
				http.DetectContentType(response),
				base64.StdEncoding.EncodeToString(response),
			)
		}
	})

	parseHTML(func() {
		anime.Title = parser.Title()
	})

	parseHTML(func() {
		anime.Status = parser.Status()
	})

	parseHTML(func() {
		anime.Rating = parser.Rating()
	})

	parseHTML(func() {
		anime.Episodes = parser.Episodes()
	})

	parseHTML(func() {
		anime.Genres = parser.Genres()
	})

	parseHTML(func() {
		anime.VoiceActing = parser.VoiceActing()
	})

	parseHTML(func() {
		anime.Synonyms = parser.Synonyms()
	})

	<-imageCh

	return anime
}
