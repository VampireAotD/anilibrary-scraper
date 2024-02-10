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
	"anilibrary-scraper/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

const endpoint string = "/api/v1/anime/parse"

type AnimeControllerSuite struct {
	suite.Suite

	useCaseMock *usecase.MockScraperUseCaseMockRecorder
	controller  Controller
	router      *fiber.App
}

func TestAnimeControllerSuite(t *testing.T) {
	suite.Run(t, new(AnimeControllerSuite))
}

func (suite *AnimeControllerSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	useCaseMock := usecase.NewMockScraperUseCase(ctrl)

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
			name, url  string
			statusCode int
			err        error
		}{
			{
				name:       "Invalid url",
				url:        "",
				statusCode: http.StatusUnprocessableEntity,
				err:        ErrInvalidURL,
			},
			{
				name:       "Unsupported url",
				url:        "https://www.google.com/",
				statusCode: http.StatusUnprocessableEntity,
				err:        scraper.ErrUnsupportedScraper,
			},
		}

		for _, testCase := range testCases {
			suite.useCaseMock.Scrape(gomock.Any(), testCase.url).Return(nil, testCase.err)

			response := suite.sendRequest(testCase.url)

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
		const url string = "https://animego.org/anime/naruto-uragannye-hroniki-103"
		expected := &entity.Anime{
			Image:       base64.StdEncoding.EncodeToString([]byte("data:image/jpeg;base64,random")),
			Title:       "Наруто: Ураганные хроники",
			Status:      "Вышел",
			Episodes:    "500",
			Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн", "Супер сила", "Экшен"},
			VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
			Synonyms:    []string{"Naruto: Shippuden", "ナルト- 疾風伝", "Naruto Hurricane Chronicles"},
			Rating:      9.5,
		}

		suite.useCaseMock.Scrape(gomock.Any(), url).Return(expected, nil)
		response := suite.sendRequest(url)
		defer func() {
			require.NoError(response.Body.Close())
		}()

		decoder := json.NewDecoder(response.Body)
		decoder.DisallowUnknownFields()

		var anime *entity.Anime

		require.NoError(decoder.Decode(&anime))
		require.Equal(http.StatusOK, response.StatusCode)
		require.Equal(expected, anime)
	})
}
