// go:build test

package scraper

import (
	"testing"

	"anilibrary-request-parser/internal/adapter/db/redis/repository"
	"anilibrary-request-parser/internal/domain/dto"
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

	repo := repository.NewMockAnimeRepository(ctrl)
	return anime.NewScraperService(repo)
}

func TestScraperService(t *testing.T) {
	testCase := "https://google.com"

	service := composeService(t)
	_, err := service.Process(composeDto(testCase))

	require.Error(t, err, "resolving scraper")
}

func TestAnimeGoScraper(t *testing.T) {
	testCase := "https://animego.org/anime/naruto-102"

	service := composeService(t)
	_, err := service.Process(composeDto(testCase))

	require.NoError(t, err, "scraping animego")
}

func TestAnimeVostScraper(t *testing.T) {
	testCase := "https://animevost.org/tip/tv/5-naruto-shippuuden12.html"

	service := composeService(t)
	_, err := service.Process(composeDto(testCase))

	require.NoError(t, err, "scraping animevost")
}
