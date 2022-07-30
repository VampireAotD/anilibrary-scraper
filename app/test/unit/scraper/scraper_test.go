package scraper

import (
	"testing"

	"anilibrary-request-parser/app/internal/domain/service/anime"
	"github.com/stretchr/testify/require"
)

func TestScraperService(t *testing.T) {
	testCase := "https://google.com"

	_, err := anime.NewScrapperService(testCase)

	require.Error(t, err, "resolving scraper")
}

func TestAnimeGoScraper(t *testing.T) {
	testCase := "https://animego.org/anime/naruto-102"

	service, _ := anime.NewScrapperService(testCase)
	_, err := service.Process()

	require.NoError(t, err, "scraping animego")
}

func TestAnimeVostScraper(t *testing.T) {
	testCase := "https://animevost.org/tip/tv/5-naruto-shippuuden12.html"

	service, _ := anime.NewScrapperService(testCase)
	_, err := service.Process()

	require.NoError(t, err, "scraping animevost")
}
