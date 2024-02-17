package parsers

import (
	"os"
	"testing"

	"anilibrary-scraper/internal/scraper/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func TestAnimeVost_FullHTML(t *testing.T) {
	html, err := os.Open("testdata/animevost/full.html")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, html.Close())
	}()

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(t, err)

	parser := NewAnimeVost(document)

	expected := model.Anime{
		Image:       "https://animevost.org/uploads/posts/2021-11/1636403661_1.png",
		Title:       "Наруто Ураганные Хроники",
		Status:      "Вышел",
		Episodes:    "500",
		Genres:      []string{"Приключения", "Боевые Искусства", "Сёнэн"},
		VoiceActing: []string{"AnimeVost"},
		Synonyms:    []string{"Naruto Shippuuden"},
		Rating:      10,
	}

	actual := model.Anime{
		Image:       parser.ImageURL(),
		Title:       parser.Title(),
		Status:      parser.Status(),
		Episodes:    parser.Episodes(),
		Genres:      parser.Genres(),
		VoiceActing: parser.VoiceActing(),
		Synonyms:    parser.Synonyms(),
		Rating:      parser.Rating(),
	}

	require.Equal(t, expected, actual)
}

func TestAnimeVost_PartialHTML(t *testing.T) {
	html, err := os.Open("testdata/animevost/partial.html")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, html.Close())
	}()

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(t, err)

	parser := NewAnimeVost(document)

	expected := model.Anime{
		Image:       "https://animevost.org/uploads/posts/2021-11/1636403661_1.png",
		Title:       "Наруто Ураганные Хроники",
		Status:      model.Ready,
		Episodes:    MinimalAnimeEpisodes,
		Genres:      nil,
		VoiceActing: []string{"AnimeVost"},
		Synonyms:    []string{"Naruto Shippuuden"},
		Rating:      MinimalAnimeRating,
	}

	actual := model.Anime{
		Image:       parser.ImageURL(),
		Title:       parser.Title(),
		Status:      parser.Status(),
		Episodes:    parser.Episodes(),
		Genres:      parser.Genres(),
		VoiceActing: parser.VoiceActing(),
		Synonyms:    parser.Synonyms(),
		Rating:      parser.Rating(),
	}

	require.Equal(t, expected, actual)
}
