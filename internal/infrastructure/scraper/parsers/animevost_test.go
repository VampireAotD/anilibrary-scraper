package parsers

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/scraper/model"

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
		Title:       "Наруто Ураганные Хроники",
		Status:      model.Ready,
		Episodes:    500,
		Genres:      []string{"Приключения", "Боевые искусства", "Сёнэн"},
		VoiceActing: []string{"AnimeVost"},
		Synonyms:    []string{"Naruto Shippuuden"},
		Rating:      10,
		Year:        2007,
		Type:        model.Show,
	}

	_, err = url.Parse(parser.ImageURL())
	require.NoError(t, err)
	require.Equal(t, expected, parser.Parse())
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
		Title:       "Убийца Акаме!",
		Status:      model.Ready,
		Episodes:    24,
		Genres:      []string{"Приключения", "Фэнтези"},
		VoiceActing: []string{"AnimeVost"},
		Synonyms:    []string{"Akame ga Kill!"},
		Rating:      MinimalAnimeRating,
		Year:        2014,
		Type:        model.Show,
	}

	_, err = url.Parse(parser.ImageURL())
	require.NoError(t, err)
	require.Equal(t, expected, parser.Parse())
}

func TestAnimeVostMovie(t *testing.T) {
	html, err := os.Open(filepath.Join("..", "testdata", "animevost", "movie.html"))
	require.NoError(t, err)

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(t, err)
	require.NoError(t, html.Close())

	parser := NewAnimeVost(document)

	expected := model.Anime{
		Title:       "Ван Пис: Бегство",
		Status:      model.Ready,
		Type:        model.Movie,
		Episodes:    1,
		Genres:      []string{"Приключения", "Комедия", "Драма", "Фэнтези"},
		VoiceActing: []string{"AnimeVost"},
		Synonyms:    []string{"Gekijouban One Piece: Stampede"},
		Rating:      8,
		Year:        2019,
	}

	_, err = url.Parse(parser.ImageURL())
	require.NoError(t, err)
	require.Equal(t, expected, parser.Parse())
}
