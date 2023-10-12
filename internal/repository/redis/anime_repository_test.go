package redis

import (
	"context"
	"encoding/base64"
	"testing"

	"anilibrary-scraper/internal/repository"
	"anilibrary-scraper/internal/repository/models"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

const testURL string = "https://animego.org/anime/naruto-uragannye-hroniki-103"

type AnimeRepositorySuite struct {
	suite.Suite

	redisServer     *miniredis.Miniredis
	animeRepository repository.AnimeRepository
	expectedAnime   models.Anime
}

func TestAnimeRepositorySuite(t *testing.T) {
	suite.Run(t, new(AnimeRepositorySuite))
}

func (suite *AnimeRepositorySuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.redisServer = miniredis.RunT(suite.T())
	suite.animeRepository = NewAnimeRepository(redis.NewClient(&redis.Options{
		Addr: suite.redisServer.Addr(),
	}))
	suite.expectedAnime = models.Anime{
		URL:         testURL,
		Image:       base64.StdEncoding.EncodeToString([]byte("random")),
		Title:       "random",
		Status:      "Вышел",
		Episodes:    "120",
		Genres:      nil,
		VoiceActing: nil,
		Synonyms:    nil,
		Rating:      0,
	}
}

func (suite *AnimeRepositorySuite) TearDownTest() {
	suite.redisServer.Del(testURL)
}

func (suite *AnimeRepositorySuite) TearDownSuite() {
	suite.redisServer.Close()
}

func (suite *AnimeRepositorySuite) TestFindByURL() {
	var (
		t       = suite.T()
		require = suite.Require()
	)

	t.Run("Not found in cache", func(t *testing.T) {
		anime, err := suite.animeRepository.FindByURL(context.Background(), testURL)
		require.Error(err)
		require.Zero(anime)
	})

	t.Run("Found in cache", func(t *testing.T) {
		err := suite.animeRepository.Create(context.Background(), suite.expectedAnime)
		require.NoError(err)

		anime, err := suite.animeRepository.FindByURL(context.Background(), testURL)
		require.NoError(err)
		require.NotZero(anime)
		require.Equal(suite.expectedAnime.MapToDomainEntity(), anime)
	})
}

func (suite *AnimeRepositorySuite) TestCreate() {
	var (
		t       = suite.T()
		require = suite.Require()
	)

	t.Run("Invalid cases", func(t *testing.T) {
		t.Run("Missing image", func(t *testing.T) {
			err := suite.animeRepository.Create(context.Background(), models.Anime{
				URL:   testURL,
				Title: "random",
			})
			require.ErrorIs(err, models.ErrInvalidData)
		})

		t.Run("Missing title", func(t *testing.T) {
			err := suite.animeRepository.Create(context.Background(), models.Anime{
				URL:   testURL,
				Image: base64.StdEncoding.EncodeToString([]byte("random")),
			})
			require.ErrorIs(err, models.ErrInvalidData)
		})
	})

	t.Run("Correct data", func(t *testing.T) {
		err := suite.animeRepository.Create(context.Background(), suite.expectedAnime)
		require.NoError(err)
	})
}
