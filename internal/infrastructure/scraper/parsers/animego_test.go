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

func TestAnimeGoShow(t *testing.T) {
	html, err := os.Open(filepath.Join("..", "testdata", "animego", "show.html"))
	require.NoError(t, err)

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(t, err)
	require.NoError(t, html.Close())

	parser := NewAnimeGo(document)

	expected := model.Anime{
		Image:       "https://animego.org/upload/anime/images/5a3ff73e8bd5b.jpg",
		Title:       "Наруто: Ураганные хроники",
		Status:      model.Ready,
		Type:        model.Show,
		Episodes:    500,
		Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн", "Супер сила", "Экшен"},
		VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
		Synonyms:    []string{"Naruto: Shippuuden", "Naruto: Shippuden", "ナルト- 疾風伝", "Naruto Hurricane Chronicles"},
		Year:        2007,
		Rating:      9.5,
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

func TestAnimeGoMovie(t *testing.T) {
	html, err := os.Open(filepath.Join("..", "testdata", "animego", "movie.html"))
	require.NoError(t, err)

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(t, err)
	require.NoError(t, html.Close())

	parser := NewAnimeGo(document)

	expected := model.Anime{
		Image:       "https://animego.org/upload/anime/images/5a3d312d9404a.jpg",
		Title:       "Ван-Пис: Золото",
		Status:      model.Ready,
		Type:        model.Movie,
		Episodes:    MinimalAnimeEpisodes,
		Genres:      []string{"Драма", "Комедия", "Приключения", "Сёнэн", "Фэнтези", "Экшен"},
		VoiceActing: []string{"Persona99", "AniMaunt", "AniMedia"},
		Synonyms:    []string{"One Piece Film: Gold", "ONE PIECE FILM GOLD", "One Piece Movie 13", "Ван-Пис фильм 13"},
		Rating:      8.4,
		Year:        2016,
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
