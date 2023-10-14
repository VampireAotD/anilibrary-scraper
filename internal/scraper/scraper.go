package scraper

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/scraper/client"
	"anilibrary-scraper/internal/scraper/parsers"
	"anilibrary-scraper/internal/scraper/parsers/model"

	"github.com/PuerkitoBio/goquery"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Scraper is a basically a factory of all parsers that can resolve parser for current url
// and scrape all data concurrently
type Scraper struct {
	client client.TLSClient
}

func New() Scraper {
	return Scraper{
		client: client.NewTLSClient(10),
	}
}

func (s Scraper) resolveParser(url string) (parsers.Contract, error) {
	switch {
	case strings.Contains(url, parsers.AnimeGoURL):
		return parsers.NewAnimeGo(), nil
	case strings.Contains(url, parsers.AnimeVostURL):
		return parsers.NewAnimeVost(), nil
	default:
		return nil, ErrUnsupportedScraper
	}
}

func (s Scraper) recover() {
	if err := recover(); err != nil {
		metrics.IncrPanicCounter()
	}
}

func (s Scraper) process(parser parsers.Contract, document *goquery.Document) *model.Anime {
	var (
		wg    sync.WaitGroup
		anime = new(model.Anime)
	)

	parse := func(processor func()) {
		defer s.recover()
		defer wg.Done()
		processor()
	}

	wg.Add(8)

	go parse(func() {
		response, err := s.client.FetchResponseBody(parser.Image(document))
		if err != nil {
			return
		}

		anime.Image = fmt.Sprintf(
			"data:%s;base64,%s",
			http.DetectContentType(response),
			base64.StdEncoding.EncodeToString(response),
		)
	})

	go parse(func() {
		anime.Title = parser.Title(document)
	})

	go parse(func() {
		anime.Status = parser.Status(document)
	})

	go parse(func() {
		anime.Rating = parser.Rating(document)
	})

	go parse(func() {
		anime.Episodes = parser.Episodes(document)
	})

	go parse(func() {
		anime.Genres = parser.Genres(document)
	})

	go parse(func() {
		anime.VoiceActing = parser.VoiceActing(document)
	})

	go parse(func() {
		anime.Synonyms = parser.Synonyms(document)
	})

	wg.Wait()

	return anime
}

func (s Scraper) Scrape(ctx context.Context, url string) (*entity.Anime, error) {
	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("Scraper").Start(ctx, "Scrape")
	defer span.End()

	parser, err := s.resolveParser(url)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	document, err := s.client.FetchDocument(url)
	if err != nil {
		return nil, err
	}

	anime := s.process(parser, document)
	if err = anime.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return anime.MapToDomainEntity(), nil
}
