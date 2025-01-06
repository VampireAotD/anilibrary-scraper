package redis

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/infrastructure/repository/model"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

const testURL string = "https://animego.org/anime/naruto-uragannye-hroniki-103"

type AnimeRepositorySuite struct {
	suite.Suite

	redisServer   *miniredis.Miniredis
	repository    AnimeRepository
	expectedAnime model.Anime
}

func TestAnimeRepositorySuite(t *testing.T) {
	suite.Run(t, new(AnimeRepositorySuite))
}

func (ars *AnimeRepositorySuite) SetupSuite() {
	ctrl := gomock.NewController(ars.T())
	defer ctrl.Finish()

	ars.redisServer = miniredis.RunT(ars.T())
	ars.repository = NewAnimeRepository(
		redis.NewClient(&redis.Options{
			Addr: ars.redisServer.Addr(),
		}),
	)
	ars.expectedAnime = model.Anime{
		URL:      testURL,
		Image:    base64.StdEncoding.EncodeToString([]byte("data:image/png;base64,image")),
		Title:    "test",
		Status:   entity.Ready,
		Type:     entity.Show,
		Episodes: 120,
		Year:     time.Now().Year(),
	}
}

func (ars *AnimeRepositorySuite) TearDownTest() {
	ars.redisServer.Del(testURL)
}

func (ars *AnimeRepositorySuite) TearDownSuite() {
	ars.redisServer.Close()
}

func (ars *AnimeRepositorySuite) TestFindByURL() {
	var (
		t       = ars.T()
		require = ars.Require()
	)

	t.Run("Not found in cache", func(_ *testing.T) {
		anime, err := ars.repository.FindByURL(context.Background(), testURL)
		require.Error(err)
		require.Zero(anime)
	})

	t.Run("Found in cache", func(_ *testing.T) {
		err := ars.repository.Create(context.Background(), ars.expectedAnime)
		require.NoError(err)

		anime, err := ars.repository.FindByURL(context.Background(), testURL)
		require.NoError(err)
		require.NotZero(anime)
		require.Equal(ars.expectedAnime.MapToDomainEntity(), anime)
	})
}

func (ars *AnimeRepositorySuite) TestCreate() {
	var (
		t       = ars.T()
		require = ars.Require()
	)

	t.Run("With required TTL", func(_ *testing.T) {
		err := ars.repository.Create(context.Background(), ars.expectedAnime)
		require.NoError(err)

		ttl := ars.redisServer.TTL(testURL)
		expectedTTL, err := time.ParseDuration(sevenDaysInHours)
		require.NoError(err)
		require.Equal(expectedTTL, ttl)
	})
}
