package scraper_test

import (
	"testing"

	"anilibrary-request-parser/internal/domain/dto"
	"anilibrary-request-parser/internal/domain/repository/mock"
	"anilibrary-request-parser/internal/domain/service/anime"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func composeDto(testCase string) dto.ParseDTO {
	return dto.ParseDTO{
		Url:       testCase,
		FromCache: false,
	}
}

func composeService(t *testing.T) *anime.ScraperService {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockAnimeRepository(ctrl)
	return anime.NewScraperService(repo)
}

func TestScraperService(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		requireError bool
	}{
		{
			name:         "Random url",
			url:          "https://google.com",
			requireError: true,
		},
		{
			name:         "AnimeGo",
			url:          "https://animego.org/anime/naruto-102",
			requireError: false,
		},
		{
			name:         "AnimeVostOrg",
			url:          "https://animevost.org/tip/tv/5-naruto-shippuuden12.html",
			requireError: false,
		},
	}

	t.Run("Scraper tests", func(t *testing.T) {
		service := composeService(t)

		for _, testCase := range cases {
			t.Run(testCase.name, func(t *testing.T) {
				result, err := service.Process(composeDto(testCase.url))

				if testCase.requireError {
					require.Error(t, err, testCase.name)
					require.Nil(t, result)
				} else {
					require.NoError(t, err, testCase.name)
					require.NotNil(t, result)
				}
			})
		}
	})
}
