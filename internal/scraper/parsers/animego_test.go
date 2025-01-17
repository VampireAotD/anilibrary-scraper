package parsers

import (
	"os"
	"testing"

	"anilibrary-scraper/internal/scraper/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func TestAnimeGo_FullHTML(t *testing.T) {
	html, err := os.Open("testdata/animego/full.html")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, html.Close())
	}()

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(t, err)

	parser := NewAnimeGo(document)

	expected := model.Anime{
		Image:       "https://animego.org/upload/anime/images/5a3ff73e8bd5b.jpg",
		Title:       "Наруто: Ураганные хроники",
		Status:      "Вышел",
		Episodes:    "500",
		Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн", "Супер сила", "Экшен"},
		VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
		Synonyms:    []string{"Naruto: Shippuuden", "Naruto: Shippuden", "ナルト- 疾風伝", "Naruto Hurricane Chronicles"},
		Rating:      9.5,
	}

	actual := model.Anime{
		Image:       parser.ImageURL(),
		Title:       parser.Title(),
		Status:      parser.Status(),
		Rating:      parser.Rating(),
		Episodes:    parser.Episodes(),
		Genres:      parser.Genres(),
		VoiceActing: parser.VoiceActing(),
		Synonyms:    parser.Synonyms(),
	}

	require.Equal(t, expected, actual)
}

func TestAnimeGo_PartialHTML(t *testing.T) {
	html, err := os.Open("testdata/animego/partial.html")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, html.Close())
	}()

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(t, err)

	parser := NewAnimeGo(document)

	expected := model.Anime{
		Image:       "https://animego.org/upload/anime/images/5a3ff73e8bd5b.jpg",
		Title:       "Наруто: Ураганные хроники",
		Status:      model.Ready,
		Episodes:    MinimalAnimeEpisodes,
		Genres:      nil,
		VoiceActing: nil,
		Synonyms:    nil,
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
