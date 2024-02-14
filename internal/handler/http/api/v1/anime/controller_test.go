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
	request := httptest.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewBufferString(fmt.Sprintf(`{"url":"%s"}`, url)),
	)

	request.Header.Set("Content-Type", "application/json")

	response, err := suite.router.Test(request, -1)
	suite.Require().NoError(err)

	return response
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
				err:        ErrInvalidURL,
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

			response := suite.sendRequest(testCase.dto.URL)

			decoder := json.NewDecoder(response.Body)
			decoder.DisallowUnknownFields()

			var err ErrorResponse
			require.NoError(decoder.Decode(&err))
			require.Equal(testCase.statusCode, response.StatusCode)
			require.Equal(testCase.err.Error(), err.Message)
			require.NoError(response.Body.Close())
		}
	})

	t.Run("Supported urls", func(_ *testing.T) {
		dto := scraperUseCase.DTO{
			URL: "https://animego.org/anime/naruto-uragannye-hroniki-103",
			IP:  "0.0.0.0",
		}

		expected := entity.Anime{
			Image:       base64.StdEncoding.EncodeToString([]byte("data:image/jpeg;base64,random")),
			Title:       "Наруто: Ураганные хроники",
			Status:      "Вышел",
			Episodes:    "500",
			Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн", "Супер сила", "Экшен"},
			VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
			Synonyms:    []string{"Naruto: Shippuden", "ナルト- 疾風伝", "Naruto Hurricane Chronicles"},
			Rating:      9.5,
		}

		suite.useCaseMock.Scrape(gomock.Any(), dto).Return(expected, nil)
		response := suite.sendRequest(dto.URL)
		defer func() {
			require.NoError(response.Body.Close())
		}()

		decoder := json.NewDecoder(response.Body)
		decoder.DisallowUnknownFields()

		var anime entity.Anime

		require.NoError(decoder.Decode(&anime))
		require.Equal(http.StatusOK, response.StatusCode)
		require.Equal(expected, anime)
	})
}
