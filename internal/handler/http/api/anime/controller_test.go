package anime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/service/mocks"
	"anilibrary-scraper/internal/handler/http/middleware"
	"anilibrary-scraper/pkg/logger"
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

	ctx := middleware.WithLogger(request.Context(), logger.NewLogger(io.Discard))
	handler(recorder, request.WithContext(ctx))

	return recorder
}

func composeDTO(url string) dto.RequestDTO {
	return dto.RequestDTO{
		Url:       url,
		FromCache: true,
	}
}

func (suite *AnimeControllerSuite) TestParse() {
	t := suite.T()

	expected := &entity.Anime{
		Title:    "Блич: Тысячелетняя кровавая война",
		Status:   "Онгоинг",
		Episodes: "1 / ?",
		Rating:   9.7,
	}

	testCases := []struct {
		name, url    string
		statusCode   int
		requireError bool
	}{
		{
			name:         "Invalid url",
			url:          "",
			statusCode:   http.StatusUnprocessableEntity,
			requireError: true,
		},
		{
			name:         "Unsupported url",
			url:          "https://www.google.com/",
			statusCode:   http.StatusUnprocessableEntity,
			requireError: true,
		},
		{
			name:         "Supported url",
			url:          "https://animego.org/anime/blich-tysyacheletnyaya-krovavaya-voyna-2129",
			statusCode:   http.StatusOK,
			requireError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.requireError {
				suite.serviceMock.Process(composeDTO(testCase.url)).Return(nil, errors.New("expected error"))
			} else {
				suite.serviceMock.Process(composeDTO(testCase.url)).Return(expected, nil)
			}
			resp := suite.sendParseRequest(testCase.url)

			decoder := json.NewDecoder(resp.Body)
			decoder.DisallowUnknownFields()

			require.Equal(t, testCase.statusCode, resp.Code)

			if testCase.requireError {
				var err response.Error
				require.NoError(t, decoder.Decode(&err))
			} else {
				var anime *entity.Anime
				require.NoError(t, decoder.Decode(&anime))
				require.Equal(t, expected, anime)
			}
		})
	}
}
