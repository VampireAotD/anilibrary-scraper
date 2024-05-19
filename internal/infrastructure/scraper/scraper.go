package scraper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/infrastructure/metrics"
	"anilibrary-scraper/internal/infrastructure/scraper/model"
	"anilibrary-scraper/internal/infrastructure/scraper/parsers"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-playground/validator/v10"
	"github.com/phuslu/lru"
)

var (
	ErrSiteNotSupported = errors.New("site is not supported for scraping")
	ErrCouldNotScrape   = errors.New("could not scrape anime")
)

//go:generate mockgen -source=scraper.go -destination=./mocks.go -package=scraper
type (
	// HTTPClient must be implemented by any client that will be used for scraping.
	HTTPClient interface {
		// Image scrapes response and returns its as string.
		Image(ctx context.Context, url string) (string, error)

		// HTML scrapes response and returns its as goquery.Document.
		HTML(ctx context.Context, url string) (*goquery.Document, error)
	}

	// Parser must be implemented by any parser that will be used for scraping particular site.
	Parser interface {
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
)

type Scraper struct {
	client    HTTPClient
	validator *validator.Validate
	cond      *sync.Cond
	queue     *lru.TTLCache[string, struct{}]
	cache     *lru.TTLCache[string, entity.Anime]
}

func New(client HTTPClient, validate *validator.Validate) Scraper {
	return Scraper{
		client:    client,
		validator: validate,
		cond:      sync.NewCond(&sync.Mutex{}),
		queue:     lru.NewTTLCache[string, struct{}](100),
		cache:     lru.NewTTLCache[string, entity.Anime](100),
	}
}

func (s Scraper) ScrapeAnime(ctx context.Context, url string) (entity.Anime, error) {
	// TODO Make an interface for queue and cache to be able to switch implementations, from memory to Redis.

	// Before making any request check if anime is already cached
	// if it is, return it.
	if anime, cached := s.cache.Get(url); cached {
		return anime, nil
	}

	// If anime is not cached, then current url must be added to queue.
	// This is made to prevent multiple HTTP requests for the same url.
	// When receiving multiple exact same urls - make only one of them
	// to perform HTTP request, other same urls will be waiting for
	// results, if after wait there is a cached anime - return it.
	// And if even after waiting there is no anime in cache and
	// no this url not in queue - return error.
	s.cond.L.Lock()

	_, scraping := s.queue.Get(url)
	for scraping {
		s.cond.Wait()

		// If anime has been cache after waiting, return it.
		if anime, cached := s.cache.Get(url); cached {
			s.cond.Signal()
			s.cond.L.Unlock()
			return anime, nil
		}

		// The only way now that the anime is not cached after the wait
		// is that the request returned the error.
		_, scraping = s.queue.Get(url)
		if !scraping {
			s.cond.Signal()
			s.cond.L.Unlock()
			return entity.Anime{}, ErrCouldNotScrape
		}
	}

	// Add anime url to queue, this prevents making multiple HTTP requests per similar url.
	s.queue.Set(url, struct{}{}, 1*time.Minute)
	s.cond.L.Unlock()

	parser, err := s.scrape(ctx, url)
	if err != nil {
		// Delete url from queue and signal other waiting goroutines
		// for this url that the request failed.
		s.cond.L.Lock()
		s.queue.Delete(url)
		s.cond.Signal()
		s.cond.L.Unlock()

		metrics.IncrScraperFailedRequestCounter()
		return entity.Anime{}, err
	}

	anime, err := s.parse(ctx, parser)
	if err != nil {
		// Delete url from queue and signal other waiting goroutines
		// for this url that the request failed.
		s.cond.L.Lock()
		s.queue.Delete(url)
		s.cond.Signal()
		s.cond.L.Unlock()

		return entity.Anime{}, err
	}

	ent := anime.MapToDomainEntity()

	// Delete url from queue and signal other waiting goroutines
	// for this url that anime has been scraped and cached.
	s.cond.L.Lock()
	s.cache.Set(url, ent, 1*time.Minute)
	s.queue.Delete(url)
	s.cond.Signal()
	s.cond.L.Unlock()

	return ent, nil
}

func (s Scraper) scrape(ctx context.Context, url string) (Parser, error) {
	switch {
	case strings.HasPrefix(url, parsers.AnimeGoURL):
		document, err := s.client.HTML(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("fetching HTML: %w", err)
		}

		metrics.IncrScraperRequestCounter()
		return parsers.NewAnimeGo(document), nil
	case strings.HasPrefix(url, parsers.AnimeVostURL):
		document, err := s.client.HTML(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("fetching HTML: %w", err)
		}

		metrics.IncrScraperRequestCounter()
		return parsers.NewAnimeVost(document), nil
	default:
		return nil, ErrSiteNotSupported
	}
}

func (s Scraper) parse(ctx context.Context, parser Parser) (model.Anime, error) {
	var anime model.Anime

	imageCh := make(chan struct{})
	parseHTML := func(extractor func()) {
		defer func() {
			if err := recover(); err != nil {
				metrics.IncrPanicCounter()
			}
		}()

		extractor()
	}

	go parseHTML(func() {
		defer close(imageCh)

		if url := parser.ImageURL(); url != "" {
			imageCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			image, err := s.client.Image(imageCtx, url)
			if err != nil {
				metrics.IncrScraperFailedImageScrapeCounter()
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
