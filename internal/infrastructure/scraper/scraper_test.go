package scraper

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/infrastructure/scraper/model"
	"anilibrary-scraper/internal/infrastructure/scraper/parsers"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

const (
	animego   string = "https://animego.org/anime/naruto-uragannye-hroniki-103"
	animevost string = "https://animevost.org/tip/tv/5-naruto-shippuuden12.html"
)

type ScraperSuite struct {
	suite.Suite

	controller *gomock.Controller
	clientMock *MockHTTPClientMockRecorder
	parserMock *MockParserMockRecorder

	scraper  Scraper
	testData map[string]*goquery.Document
}

func TestScraperSuite(t *testing.T) {
	suite.Run(t, new(ScraperSuite))
}

func (s *ScraperSuite) SetupSuite() {
	s.controller = gomock.NewController(s.T())

	client := NewMockHTTPClient(s.controller)

	s.clientMock = client.EXPECT()
	s.parserMock = NewMockParser(s.controller).EXPECT()

	s.scraper = New(client, validator.New())

	s.testData = map[string]*goquery.Document{
		animego:   s.loadHTML(filepath.Join("testdata", "animego", "full.html")),
		animevost: s.loadHTML(filepath.Join("testdata", "animevost", "full.html")),
	}
}

func (s *ScraperSuite) loadHTML(path string) *goquery.Document {
	file, err := os.Open(path)
	s.Require().NoError(err)

	html, err := goquery.NewDocumentFromReader(file)
	s.Require().NoError(err)
	s.Require().NoError(file.Close())

	return html
}

func (s *ScraperSuite) TearDownSubTest() {
	s.scraper.cache.Delete(animego)
	s.scraper.cache.Delete(animevost)
}

func (s *ScraperSuite) TearDownSuite() {
	s.controller.Finish()
}

func (s *ScraperSuite) TestScrapeAnime() {
	var require = s.Require()

	s.Run("Cache hit", func() {
		cached := entity.Anime{
			Image:       base64.StdEncoding.EncodeToString([]byte("data:image/png;base64,image")),
			Title:       "Наруто: Ураганные хроники",
			Status:      model.Ready,
			Type:        model.Show,
			Episodes:    "500",
			Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн", "Супер сила", "Экшен"},
			VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
			Synonyms:    []string{"Naruto: Shippuden", "ナルト- 疾風伝", "Naruto Hurricane Chronicles"},
			Rating:      9.5,
			Year:        2007,
		}

		s.scraper.cache.Set(animego, cached, 1*time.Minute)

		anime, err := s.scraper.ScrapeAnime(context.Background(), animego)
		require.NoError(err)
		require.Equal(cached, anime)
	})

	s.Run("Multiple requests", func() {
		s.Run("Same URL", func() {
			s.clientMock.HTML(gomock.Any(), animego).Return(s.testData[animego], nil).Times(1)
			s.clientMock.Image(gomock.Any(), gomock.Any()).Return("data:image/png;base64,image", nil).Times(1)

			var wg sync.WaitGroup

			wg.Add(3)

			for i := 0; i < 3; i++ {
				go func() {
					defer wg.Done()

					if i == 0 {
						_, exists := s.scraper.cache.Get(animego)
						require.False(exists, "cache will not contain any data at first request")
					}

					anime, err := s.scraper.ScrapeAnime(context.Background(), animego)
					require.NoError(err)
					require.NotEmpty(anime)

					_, exists := s.scraper.cache.Get(animego)
					require.True(exists)
				}()
			}

			wg.Wait()
		})

		s.Run("Different URLs", func() {
			s.clientMock.HTML(gomock.Any(), animego).Return(s.testData[animego], nil).Times(1)
			s.clientMock.HTML(gomock.Any(), animevost).Return(s.testData[animevost], nil).Times(1)
			s.clientMock.Image(gomock.Any(), gomock.Any()).Return("data:image/png;base64,image", nil).Times(2)

			var wg sync.WaitGroup

			wg.Add(10)

			testCases := []string{animego, animevost}

			for i := 0; i < 10; i++ {
				go func() {
					defer wg.Done()

					anime, err := s.scraper.ScrapeAnime(context.Background(), testCases[i%2])
					require.NoError(err)
					require.NotEmpty(anime)

					_, exists := s.scraper.cache.Get(testCases[i%2])
					require.True(exists)
				}()
			}

			wg.Wait()
		})
	})
}

func (s *ScraperSuite) TestScrape() {
	var require = s.Require()

	s.Run("Unsupported site", func() {
		parser, err := s.scraper.scrape(context.Background(), "https://google.com")
		require.Error(err)
		require.Nil(parser)
	})

	s.Run("Supported site", func() {
		s.Run("AnimeGo", func() {
			s.clientMock.HTML(context.Background(), animego).Return(&goquery.Document{}, nil)

			parser, err := s.scraper.scrape(context.Background(), animego)
			require.NoError(err)
			require.NotNil(parser)
			require.IsType(parsers.AnimeGo{}, parser)
		})

		s.Run("AnimeVost", func() {
			s.clientMock.HTML(context.Background(), animevost).Return(&goquery.Document{}, nil)

			parser, err := s.scraper.scrape(context.Background(), animevost)
			require.NoError(err)
			require.NotNil(parser)
			require.IsType(parsers.AnimeVost{}, parser)
		})
	})
}

func (s *ScraperSuite) TestParse() {
	var require = s.Require()

	s.Run("No image", func() {
		s.clientMock.Image(gomock.Any(), gomock.Any()).Return("", errors.New("no image")).Times(1)
		anime, err := s.scraper.parse(context.Background(), parsers.NewAnimeGo(s.testData[animego]))

		require.Error(err)
		require.Empty(anime)
	})

	s.Run("Parse", func() {
		s.Run("AnimeGo", func() {
			s.clientMock.Image(gomock.Any(), gomock.Any()).Return("data:image/png;base64,image", nil).Times(1)
			anime, err := s.scraper.parse(context.Background(), parsers.NewAnimeGo(s.testData[animego]))

			require.NoError(err)
			require.NotEmpty(anime)
		})

		s.Run("AnimeVost", func() {
			s.clientMock.Image(gomock.Any(), gomock.Any()).Return("data:image/png;base64,image", nil).Times(1)
			anime, err := s.scraper.parse(context.Background(), parsers.NewAnimeVost(s.testData[animevost]))

			require.NoError(err)
			require.NotEmpty(anime)
		})
	})
}
