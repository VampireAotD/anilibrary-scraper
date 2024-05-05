package scraper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/scraper/model"
	"anilibrary-scraper/internal/scraper/parsers"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-playground/validator/v10"
)

var ErrSiteNotSupported = errors.New("site is not supported for scraping")

// HTTPClient must be implemented by any client that will be used for scraping.
type HTTPClient interface {
	// Image scrapes response and returns its as string.
	Image(ctx context.Context, url string) (string, error)

	// HTML scrapes response and returns its as goquery.Document.
	HTML(ctx context.Context, url string) (*goquery.Document, error)
}

// Parser must be implemented by any parser that will be used for scraping particular site.
type Parser interface {
	// ImageURL scrapes and returns the URL of an anime's promotional image or cover art.
	ImageURL() string

	// Title scrapes and returns the title of the anime.
	Title() string

	// Status scrapes and returns the current status of the anime (e.g., ongoing, completed).
	Status() model.Status

	// Rating scrapes and returns the current rating of the anime from a predetermined source.
	Rating() float32

	// Episodes scrapes and returns the total number of episodes for the anime.
	Episodes() string

	// Genres scrapes and returns all genres associated with the anime.
	Genres() []string

	// VoiceActing scrapes and returns the list of voice actors associated with the anime.
	VoiceActing() []string

	// Synonyms scrapes and returns alternative names or titles for the anime.
	Synonyms() []string

	// Year scrapes and returns the year the anime was released.
	Year() int

	// Type scrapes and returns the format type of the anime (e.g., TV series, movie).
	Type() model.Type
}

type Scraper struct {
	client       HTTPClient
	validator    *validator.Validate
	panicHandler func()
}

func New(client HTTPClient, validate *validator.Validate) Scraper {
	return Scraper{
		client:    client,
		validator: validate,
		panicHandler: func() {
			if err := recover(); err != nil {
				metrics.IncrPanicCounter()
			}
		},
	}
}

func (s Scraper) ScrapeAnime(ctx context.Context, url string) (entity.Anime, error) {
	parser, err := s.scrape(ctx, url)
	if err != nil {
		return entity.Anime{}, err
	}

	anime, err := s.parse(ctx, parser)
	if err != nil {
		return entity.Anime{}, err
	}

	return anime.MapToDomainEntity(), nil
}

func (s Scraper) scrape(ctx context.Context, url string) (Parser, error) {
	switch {
	case strings.Contains(url, parsers.AnimeGoURL):
		document, err := s.client.HTML(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("fetching HTML: %w", err)
		}

		return parsers.NewAnimeGo(document), nil
	case strings.Contains(url, parsers.AnimeVostURL):
		document, err := s.client.HTML(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("fetching HTML: %w", err)
		}

		return parsers.NewAnimeVost(document), nil
	default:
		return nil, ErrSiteNotSupported
	}
}

func (s Scraper) parse(ctx context.Context, parser Parser) (model.Anime, error) {
	var anime model.Anime

	imageCh := make(chan struct{})
	parseHTML := func(extractor func()) {
		defer s.panicHandler()
		extractor()
	}

	go parseHTML(func() {
		defer close(imageCh)

		if url := parser.ImageURL(); url != "" {
			imageCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			image, err := s.client.Image(imageCtx, url)
			if err != nil {
				return
			}

			anime.Image = image
		}
	})

	parseHTML(func() { anime.Title = parser.Title() })
	parseHTML(func() { anime.Status = parser.Status() })
	parseHTML(func() { anime.Rating = parser.Rating() })
	parseHTML(func() { anime.Episodes = parser.Episodes() })
	parseHTML(func() { anime.Genres = parser.Genres() })
	parseHTML(func() { anime.VoiceActing = parser.VoiceActing() })
	parseHTML(func() { anime.Synonyms = parser.Synonyms() })
	parseHTML(func() { anime.Year = parser.Year() })
	parseHTML(func() { anime.Type = parser.Type() })

	<-imageCh

	if err := anime.Validate(s.validator); err != nil {
		return model.Anime{}, fmt.Errorf("validating response data: %w", err)
	}

	return anime, nil
}
