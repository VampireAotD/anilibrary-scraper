// go:build test

package scraper

import (
	"testing"

	"anilibrary-request-parser/app/internal/domain/dto"
	"anilibrary-request-parser/app/internal/domain/service/anime"
	"github.com/stretchr/testify/require"
)

var (
	service = anime.NewScraperService(nil)
)

func composeDto(testCase string) dto.ParseDTO {
	return dto.ParseDTO{
		Url:       testCase,
		FromCache: false,
	}
}

func TestScraperService(t *testing.T) {
	testCase := "https://google.com"

	_, err := service.Process(composeDto(testCase))

	require.Error(t, err, "resolving scraper")
}

func TestAnimeGoScraper(t *testing.T) {
	testCase := "https://animego.org/anime/naruto-102"

	_, err := service.Process(composeDto(testCase))

	require.NoError(t, err, "scraping animego")
}

func TestAnimeVostScraper(t *testing.T) {
	testCase := "https://animevost.org/tip/tv/5-naruto-shippuuden12.html"

	_, err := service.Process(composeDto(testCase))

	require.NoError(t, err, "scraping animevost")
}
