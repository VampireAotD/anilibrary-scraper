package scraper

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/metrics"
	"anilibrary-scraper/internal/scraper/client"
	"anilibrary-scraper/internal/scraper/parsers"
	"anilibrary-scraper/internal/scraper/parsers/model"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Scraper is a basically a factory of all parsers that can resolve parser for current url
// and scrape all data concurrently
type Scraper struct {
	client client.TLSClient
	wg     *sync.WaitGroup
	anime  *model.Anime
}

func New() Scraper {
	return Scraper{
		client: client.NewTLSClient(10),
		wg:     new(sync.WaitGroup),
		anime:  new(model.Anime),
	}
}

type processor func(anime *model.Anime)

func (s Scraper) recover() {
	if err := recover(); err != nil {
		metrics.IncrPanicCounter()
	}
}

func (s Scraper) parse(callback processor) {
	defer s.recover()
	defer s.wg.Done()

	callback(s.anime)
}

func (s Scraper) Scrape(ctx context.Context, url string) (*entity.Anime, error) {
	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("Scraper").Start(ctx, "Scrape")
	defer span.End()

	var parser parsers.Contract

	switch {
	case strings.Contains(url, parsers.AnimeGoURL):
		parser = parsers.NewAnimeGo()
	case strings.Contains(url, parsers.AnimeVostURL):
		parser = parsers.NewAnimeVost()
	default:
		return nil, ErrUnsupportedScraper
	}

	document, err := s.client.FetchDocument(url)
	if err != nil {
		return nil, err
	}

	s.wg.Add(8)

	go s.parse(func(anime *model.Anime) {
		response, err := s.client.FetchResponseBody(parser.Image(document))
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return
		}

		anime.Image = fmt.Sprintf(
			"data:%s;base64,%s",
			http.DetectContentType(response),
			base64.StdEncoding.EncodeToString(response),
		)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Title = parser.Title(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Status = parser.Status(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Rating = parser.Rating(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Episodes = parser.Episodes(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Genres = parser.Genres(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.VoiceActing = parser.VoiceActing(document)
	})

	go s.parse(func(anime *model.Anime) {
		anime.Synonyms = parser.Synonyms(document)
	})

	s.wg.Wait()

	if !s.anime.IsValid() {
		return nil, model.ErrInvalidParsedData
	}

	return s.anime.ToEntity(), nil
}
