package anime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/service/mocks"
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/internal/scraper"
	"anilibrary-scraper/pkg/logging"
	"anilibrary-scraper/pkg/response"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AnimeControllerSuite struct {
	suite.Suite

	serviceMock *mocks.MockScraperServiceMockRecorder
	controller  Controller
}

func TestAnimeControllerSuite(t *testing.T) {
	suite.Run(t, new(AnimeControllerSuite))
}

func (suite *AnimeControllerSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	service := mocks.NewMockScraperService(ctrl)

	suite.serviceMock = service.EXPECT()
	suite.controller = NewController(service)
}

func (suite *AnimeControllerSuite) sendParseRequest(url string) *httptest.ResponseRecorder {
	handler := suite.controller.Parse
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/anime/parse",
		bytes.NewBufferString(fmt.Sprintf(`{"url":"%s"}`, url)),
	)

	ctx := middleware.WithLogger(request.Context(), logging.NewLogger(io.Discard))
	ctx = middleware.WithTracer(ctx)
	handler(recorder, request.WithContext(ctx))

	return recorder
}

func (suite *AnimeControllerSuite) TestParse() {
	t := suite.T()

	t.Run("Bad request", func(t *testing.T) {
		testCases := []struct {
			name, url  string
			statusCode int
			err        error
		}{
			{
				name:       "Invalid url",
				url:        "",
				statusCode: http.StatusUnprocessableEntity,
				err:        ErrInvalidUrl,
			},
			{
				name:       "Unsupported url",
				url:        "https://www.google.com/",
				statusCode: http.StatusUnprocessableEntity,
				err:        scraper.ErrUnsupportedScraper,
			},
		}

		for _, testCase := range testCases {
			suite.serviceMock.Process(gomock.Any(), testCase.url).Return(nil, testCase.err)

			resp := suite.sendParseRequest(testCase.url)

			decoder := json.NewDecoder(resp.Body)
			decoder.DisallowUnknownFields()

			var err response.Error
			require.NoError(t, decoder.Decode(&err))
			require.Equal(t, testCase.statusCode, resp.Code)
			require.Equal(t, err.Message, testCase.err.Error())
		}
	})

	t.Run("Supported urls", func(t *testing.T) {
		const url string = "https://animego.org/anime/naruto-uragannye-hroniki-103"
		expected := &entity.Anime{
			Title:       "Наруто: Ураганные хроники",
			Status:      "Вышел",
			Episodes:    "500",
			Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн", "Супер сила", "Экшен"},
			VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
			Rating:      9.5,
		}

		suite.serviceMock.Process(gomock.Any(), url).Return(expected, nil)
		resp := suite.sendParseRequest(url)

		decoder := json.NewDecoder(resp.Body)
		decoder.DisallowUnknownFields()

		var anime *entity.Anime

		require.NoError(t, decoder.Decode(&anime))
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, expected, anime)
	})
}
