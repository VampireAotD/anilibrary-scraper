package parsers

import (
	"os"
	"testing"

	"anilibrary-scraper/internal/scraper/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/suite"
)

type AnimeVostSuite struct {
	suite.Suite

	parser AnimeVost
}

func TestAnimeVostSuiteSuite(t *testing.T) {
	suite.Run(t, new(AnimeVostSuite))
}

func (suite *AnimeVostSuite) SetupSuite() {
	suite.parser = NewAnimeVost()
}

func (suite *AnimeVostSuite) TestFullHTML() {
	require := suite.Require()

	html, err := os.Open("testdata/animevost/full.html")
	require.NoError(err)
	defer func() {
		require.NoError(html.Close())
	}()

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(err)

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

func (suite *AnimeVostSuite) TestPartialHTML() {
	require := suite.Require()

	html, err := os.Open("testdata/animevost/partial.html")
	require.NoError(err)
	defer func() {
		require.NoError(html.Close())
	}()

	document, err := goquery.NewDocumentFromReader(html)
	require.NoError(err)

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
