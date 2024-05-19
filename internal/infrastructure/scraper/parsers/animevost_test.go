package parsers

import (
	"os"
	"path/filepath"
	"testing"

	"anilibrary-scraper/internal/infrastructure/scraper/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func TestAnimeVost_FullHTML(t *testing.T) {
	html, err := os.Open(filepath.Join("..", "testdata", "animevost", "full.html"))
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
		Genres:      []string{"Приключения", "Боевые искусства", "Сёнэн"},
		VoiceActing: []string{"AnimeVost"},
		Synonyms:    []string{"Naruto Shippuuden"},
		Rating:      10,
		Year:        2007,
		Type:        model.Show,
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
		Year:        parser.Year(),
		Type:        parser.Type(),
	}

	require.Equal(t, expected, actual)
}

func TestAnimeVost_PartialHTML(t *testing.T) {
	html, err := os.Open(filepath.Join("..", "testdata", "animevost", "partial.html"))
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
		VoiceActing: []string{"AnimeVost"},
		Synonyms:    []string{"Naruto Shippuuden"},
		Rating:      MinimalAnimeRating,
		Year:        2007,
		Type:        model.Show,
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
		Year:        parser.Year(),
		Type:        parser.Type(),
	}

	require.Equal(t, expected, actual)
}
