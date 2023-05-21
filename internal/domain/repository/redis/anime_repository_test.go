package redis

import (
	"context"
	"encoding/base64"
	"testing"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/repository"

	"github.com/alicebob/miniredis/v2"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

const testURL string = "https://animego.org/anime/naruto-uragannye-hroniki-103"

var expectedAnime = &entity.Anime{
	Image:       base64.StdEncoding.EncodeToString([]byte("data:image/jpeg;base64,random")),
	Title:       "Наруто: Ураганные хроники",
	Status:      "Вышел",
	Episodes:    "500",
	Genres:      []string{"Боевые искусства", "Комедия", "Сёнэн", "Супер сила", "Экшен"},
	VoiceActing: []string{"AniDUB", "AniLibria", "SHIZA Project", "2x2"},
	Synonyms:    []string{"Naruto: Shippuden", "ナルト- 疾風伝", "Naruto Hurricane Chronicles"},
	Rating:      9.5,
}

type AnimeRepositorySuite struct {
	suite.Suite

	redisServer     *miniredis.Miniredis
	animeRepository repository.AnimeRepository
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
}

func (suite *AnimeRepositorySuite) TearDownTest() {
	suite.redisServer.Del(testURL)
}

func (suite *AnimeRepositorySuite) TearDownSuite() {
	suite.redisServer.Close()
}

func (suite *AnimeRepositorySuite) TestCreate() {
	t := suite.T()
	require := suite.Require()

	t.Run("Invalid cases", func(t *testing.T) {
		t.Run("Incorrect data", func(t *testing.T) {
			err := suite.animeRepository.Create(context.Background(), testURL, new(entity.Anime))
			require.ErrorIs(err, entity.ErrNotEnoughData)
		})
	})

	t.Run("Correct data", func(t *testing.T) {
		err := suite.animeRepository.Create(context.Background(), testURL, expectedAnime)
		require.NoError(err)
	})

}

func (suite *AnimeRepositorySuite) TestFindByURL() {
	t := suite.T()
	require := suite.Require()

	t.Run("Not found in cache", func(t *testing.T) {
		anime, err := suite.animeRepository.FindByURL(context.Background(), testURL)
		require.Error(err)
		require.Nil(anime)
	})

	t.Run("Found in cache", func(t *testing.T) {
		err := suite.animeRepository.Create(context.Background(), testURL, expectedAnime)
		require.NoError(err)

		anime, err := suite.animeRepository.FindByURL(context.Background(), testURL)
		require.NoError(err)
		require.NotNil(anime)
		require.Equal(expectedAnime, anime)
	})
}
