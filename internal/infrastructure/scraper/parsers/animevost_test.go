package parsers

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"anilibrary-scraper/internal/infrastructure/scraper/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func TestAnimeVostShow(t *testing.T) {
	html, err := os.Open(filepath.Join("..", "testdata", "animevost", "show.html"))
	require.NoError(t, err)

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(t, err)
	require.NoError(t, html.Close())

	parser := NewAnimeVost(document)

	expected := model.Anime{
		Image:       "https://animevost.org/uploads/posts/2016-06/1464842897_1.jpg",
		Title:       "Наруто Ураганные Хроники",
		Status:      model.Ready,
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
		Type:        parser.Type(),
		Episodes:    parser.Episodes(),
		Genres:      parser.Genres(),
		VoiceActing: parser.VoiceActing(),
		Synonyms:    parser.Synonyms(),
		Rating:      parser.Rating(),
		Year:        parser.Year(),
	}

	require.Equal(t, expected, actual)

	_, err = url.Parse(actual.Image)
	require.NoError(t, err)
}

// Sometimes animevost renders a mobile layout
func TestAnimeVostMobileGrid(t *testing.T) {
	html, err := os.Open(filepath.Join("..", "testdata", "animevost", "mobile.html"))
	require.NoError(t, err)

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(t, err)
	require.NoError(t, html.Close())

	parser := NewAnimeVost(document)

	expected := model.Anime{
		Image:       "https://animevost.org/uploads/posts/2014-08/1409038345_1.jpg",
		Title:       "Убийца Акаме!",
		Status:      model.Ready,
		Episodes:    "24",
		Genres:      []string{"Приключения", "Фэнтези"},
		VoiceActing: []string{"AnimeVost"},
		Synonyms:    []string{"Akame ga Kill!"},
		Rating:      MinimalAnimeRating,
		Year:        2014,
		Type:        model.Show,
	}

	actual := model.Anime{
		Image:       parser.ImageURL(),
		Title:       parser.Title(),
		Status:      parser.Status(),
		Type:        parser.Type(),
		Episodes:    parser.Episodes(),
		Genres:      parser.Genres(),
		VoiceActing: parser.VoiceActing(),
		Synonyms:    parser.Synonyms(),
		Rating:      parser.Rating(),
		Year:        parser.Year(),
	}

	require.Equal(t, expected, actual)

	_, err = url.Parse(actual.Image)
	require.NoError(t, err)
}

func TestAnimeVostMovie(t *testing.T) {
	html, err := os.Open(filepath.Join("..", "testdata", "animevost", "movie.html"))
	require.NoError(t, err)

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(t, err)
	require.NoError(t, html.Close())

	parser := NewAnimeVost(document)

	expected := model.Anime{
		Image:       "https://animevost.org/uploads/posts/2020-02/1581173195_1.jpg",
		Title:       "Ван Пис: Бегство",
		Status:      model.Ready,
		Type:        model.Movie,
		Episodes:    "1",
		Genres:      []string{"Приключения", "Комедия", "Драма", "Фэнтези"},
		VoiceActing: []string{"AnimeVost"},
		Synonyms:    []string{"Gekijouban One Piece: Stampede"},
		Rating:      8,
		Year:        2019,
	}

	actual := model.Anime{
		Image:       parser.ImageURL(),
		Title:       parser.Title(),
		Status:      parser.Status(),
		Type:        parser.Type(),
		Episodes:    parser.Episodes(),
		Genres:      parser.Genres(),
		VoiceActing: parser.VoiceActing(),
		Synonyms:    parser.Synonyms(),
		Rating:      parser.Rating(),
		Year:        parser.Year(),
	}

	require.Equal(t, expected, actual)

	_, err = url.Parse(actual.Image)
	require.NoError(t, err)
}
