package parsers

import (
	"os"
	"testing"

	"anilibrary-scraper/internal/scraper/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/suite"
)

type AnimeGoSuite struct {
	suite.Suite

	parser AnimeGo
}

func TestAnimeGoSuite(t *testing.T) {
	suite.Run(t, new(AnimeGoSuite))
}

func (suite *AnimeGoSuite) SetupSuite() {
	suite.parser = NewAnimeGo()
}

func (suite *AnimeGoSuite) TestFullHTML() {
	require := suite.Require()

	html, err := os.Open("testdata/animego/full.html")
	require.NoError(err)
	defer func() {
		require.NoError(html.Close())
	}()

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(err)

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

	var actual model.Anime

	actual.Image = suite.parser.Image(document)
	actual.Title = suite.parser.Title(document)
	actual.Status = suite.parser.Status(document)
	actual.Episodes = suite.parser.Episodes(document)
	actual.Genres = suite.parser.Genres(document)
	actual.VoiceActing = suite.parser.VoiceActing(document)
	actual.Synonyms = suite.parser.Synonyms(document)
	actual.Rating = suite.parser.Rating(document)

	require.Equal(expected, actual)
}

func (suite *AnimeGoSuite) TestPartialHTML() {
	require := suite.Require()

	html, err := os.Open("testdata/animego/partial.html")
	require.NoError(err)
	defer func() {
		require.NoError(html.Close())
	}()

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(err)

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

	var actual model.Anime

	actual.Image = suite.parser.Image(document)
	actual.Title = suite.parser.Title(document)
	actual.Status = suite.parser.Status(document)
	actual.Episodes = suite.parser.Episodes(document)
	actual.Genres = suite.parser.Genres(document)
	actual.VoiceActing = suite.parser.VoiceActing(document)
	actual.Synonyms = suite.parser.Synonyms(document)
	actual.Rating = suite.parser.Rating(document)

	require.Equal(expected, actual)
}
