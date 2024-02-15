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
	"anilibrary-scraper/internal/handler/http/api/v1/anime/request"
	"anilibrary-scraper/internal/handler/http/api/v1/anime/response"
	"anilibrary-scraper/internal/scraper"
	scraperUseCase "anilibrary-scraper/internal/usecase/scraper"

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

func (suite *AnimeControllerSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	useCaseMock := NewMockScraperUseCase(ctrl)

	suite.useCaseMock = useCaseMock.EXPECT()
	suite.controller = NewController(useCaseMock)
	suite.router = fiber.New()
	suite.router.Post(endpoint, suite.controller.Parse)
}

func (suite *AnimeControllerSuite) sendRequest(url string) *http.Response {
	req := httptest.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewBufferString(fmt.Sprintf(`{"url":"%s"}`, url)),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.router.Test(req, -1)
	suite.Require().NoError(err)

	return resp
}

func (suite *AnimeControllerSuite) TestParse() {
	var (
		t       = suite.T()
		require = suite.Require()
	)

	t.Run("Bad request", func(_ *testing.T) {
		testCases := []struct {
			name       string
			dto        scraperUseCase.DTO
			statusCode int
			err        error
		}{
			{
				name: "Invalid url",
				dto: scraperUseCase.DTO{
					URL: "",
					IP:  "0.0.0.0",
				},
				statusCode: http.StatusUnprocessableEntity,
				err:        request.ErrInvalidURL,
			},
			{
				name: "Unsupported url",
				dto: scraperUseCase.DTO{
					URL: "https://www.google.com",
					IP:  "0.0.0.0",
				},
				statusCode: http.StatusUnprocessableEntity,
				err:        scraper.ErrUnsupportedScraper,
			},
		}

		for _, testCase := range testCases {
			suite.useCaseMock.Scrape(gomock.Any(), testCase.dto).Return(entity.Anime{}, testCase.err)

			resp := suite.sendRequest(testCase.dto.URL)

			decoder := json.NewDecoder(resp.Body)
			decoder.DisallowUnknownFields()

			var err response.ErrorResponse
			require.NoError(decoder.Decode(&err))
			require.Equal(testCase.statusCode, resp.StatusCode)
			require.Equal(testCase.err.Error(), err.Message)
			require.NoError(resp.Body.Close())
		}
	})

	t.Run("Supported urls", func(_ *testing.T) {
		dto := scraperUseCase.DTO{
			URL: "https://animego.org/anime/naruto-uragannye-hroniki-103",
			IP:  "0.0.0.0",
		}

		expectedEntity := entity.Anime{
			Image:       base64.StdEncoding.EncodeToString([]byte("data:image/jpeg;base64,random")),
			Title:       "Наруто: Ураганные хроники",
			Status:      "Вышел",
			Episodes:    "500",
			Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн", "Супер сила", "Экшен"},
			VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
			Synonyms:    []string{"Naruto: Shippuden", "ナルト- 疾風伝", "Naruto Hurricane Chronicles"},
			Rating:      9.5,
		}

		expectedResponse := response.ScrapeResponse{
			Image:    base64.StdEncoding.EncodeToString([]byte("data:image/jpeg;base64,random")),
			Title:    "Наруто: Ураганные хроники",
			Status:   "Вышел",
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
		}

		suite.useCaseMock.Scrape(gomock.Any(), dto).Return(expectedEntity, nil)
		resp := suite.sendRequest(dto.URL)
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
