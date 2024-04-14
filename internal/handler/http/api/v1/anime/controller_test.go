package anime

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"anilibrary-scraper/internal/entity"
	"anilibrary-scraper/internal/handler/http/api/v1/anime/response"
	"anilibrary-scraper/internal/scraper"
	scraperUseCase "anilibrary-scraper/internal/usecase/scraper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

const endpoint string = "/api/v1/anime/parse"

type AnimeControllerSuite struct {
	suite.Suite

	useCaseMock *MockScraperUseCaseMockRecorder
	controller  Controller
	router      *fiber.App
}

func TestAnimeControllerSuite(t *testing.T) {
	suite.Run(t, new(AnimeControllerSuite))
}

func (s *AnimeControllerSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	useCaseMock := NewMockScraperUseCase(ctrl)

	s.useCaseMock = useCaseMock.EXPECT()
	s.controller = NewController(useCaseMock, validator.New())
	s.router = fiber.New()
	s.router.Post(endpoint, s.controller.Scrape)
}

func (s *AnimeControllerSuite) sendRequest(url string) *http.Response {
	req := httptest.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewBufferString(fmt.Sprintf(`{"url":%q}`, url)),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.router.Test(req, -1)
	s.Require().NoError(err)

	return resp
}

func (s *AnimeControllerSuite) TestParse() {
	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("Bad request", func(_ *testing.T) {
		testCases := []struct {
			name            string
			URL             string
			responseMessage string
			statusCode      int
		}{
			{
				name:            "Invalid url",
				URL:             "",
				responseMessage: "Invalid URL",
				statusCode:      http.StatusUnprocessableEntity,
			},
			{
				name:            "Unsupported url",
				URL:             "https://www.google.com",
				responseMessage: scraper.ErrSiteNotSupported.Error(),
				statusCode:      http.StatusUnprocessableEntity,
			},
		}

		for _, testCase := range testCases {
			dto := scraperUseCase.DTO{URL: testCase.URL, IP: "0.0.0.0"}
			s.useCaseMock.Scrape(gomock.Any(), dto).Return(entity.Anime{}, scraper.ErrSiteNotSupported)

			resp := s.sendRequest(testCase.URL)

			decoder := json.NewDecoder(resp.Body)
			decoder.DisallowUnknownFields()

			var err response.ScrapeErrorResponse
			require.NoError(decoder.Decode(&err))
			require.Equal(testCase.statusCode, resp.StatusCode)
			require.Equal(testCase.responseMessage, err.Message)
			require.NoError(resp.Body.Close())
		}
	})

	t.Run("Supported urls", func(_ *testing.T) {
		dto := scraperUseCase.DTO{
			URL: "https://animego.org/anime/naruto-uragannye-hroniki-103",
			IP:  "0.0.0.0",
		}

		expectedEntity := entity.Anime{
			Image:       base64.StdEncoding.EncodeToString([]byte("data:image/jpeg;base64,image")),
			Title:       "Наруто: Ураганные хроники",
			Status:      entity.Ready,
			Episodes:    "500",
			Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн", "Супер сила", "Экшен"},
			VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
			Synonyms:    []string{"Naruto: Shippuden", "ナルト- 疾風伝", "Naruto Hurricane Chronicles"},
			Rating:      9.5,
			Year:        2007,
			Type:        entity.Show,
		}

		expectedResponse := response.ScrapeResponse{
			Image:    base64.StdEncoding.EncodeToString([]byte("data:image/jpeg;base64,image")),
			Title:    "Наруто: Ураганные хроники",
			Status:   string(entity.Ready),
			Episodes: "500",
			Genres: []response.Entry{
				{Name: "Боевые искусства"},
				{Name: "Комедия"},
				{Name: "Сёнэн"},
				{Name: "Супер сила"},
				{Name: "Экшен"},
			},
			VoiceActing: []response.Entry{
				{Name: "AniDUB"},
				{Name: "AniLibria"},
				{Name: "SHIZA Project"},
				{Name: "2x2"},
			},
			Synonyms: []response.Entry{
				{Name: "Naruto: Shippuden"},
				{Name: "ナルト- 疾風伝"},
				{Name: "Naruto Hurricane Chronicles"},
			},
			Rating: 9.5,
			Year:   2007,
			Type:   string(entity.Show),
		}

		s.useCaseMock.Scrape(gomock.Any(), dto).Return(expectedEntity, nil)
		resp := s.sendRequest(dto.URL)
		defer func() {
			require.NoError(resp.Body.Close())
		}()

		decoder := json.NewDecoder(resp.Body)
		decoder.DisallowUnknownFields()

		var scrapeResponse response.ScrapeResponse

		require.NoError(decoder.Decode(&scrapeResponse))
		require.Equal(http.StatusOK, resp.StatusCode)
		require.Equal(expectedResponse, scrapeResponse)
	})
}
