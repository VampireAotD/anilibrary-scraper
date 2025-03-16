package redis

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/VampireAotD/anilibrary-scraper/internal/domain/entity"
	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/model"

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

func (s *AnimeRepositorySuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.redisServer = miniredis.RunT(s.T())
	s.repository = NewAnimeRepository(
		redis.NewClient(&redis.Options{
			Addr: s.redisServer.Addr(),
		}),
	)
	s.expectedAnime = model.Anime{
		URL:      testURL,
		Image:    base64.StdEncoding.EncodeToString([]byte("data:image/png;base64,image")),
		Title:    "test",
		Status:   entity.Ready,
		Type:     entity.Show,
		Episodes: 120,
		Year:     time.Now().Year(),
	}
}

func (s *AnimeRepositorySuite) TearDownTest() {
	s.redisServer.Del(testURL)
}

func (s *AnimeRepositorySuite) TearDownSuite() {
	s.redisServer.Close()
}

func (s *AnimeRepositorySuite) TestFindByURL() {
	var require = s.Require()

	s.T().Run("Not found in cache", func(t *testing.T) {
		anime, err := s.repository.FindByURL(t.Context(), testURL)
		require.Error(err)
		require.Zero(anime)
	})

	s.T().Run("Found in cache", func(t *testing.T) {
		err := s.repository.Create(t.Context(), s.expectedAnime)
		require.NoError(err)

		anime, err := s.repository.FindByURL(t.Context(), testURL)
		require.NoError(err)
		require.NotZero(anime)
		require.Equal(s.expectedAnime.MapToDomainEntity(), anime)
	})
}

func (s *AnimeRepositorySuite) TestCreate() {
	var require = s.Require()

	s.T().Run("With required TTL", func(t *testing.T) {
		err := s.repository.Create(t.Context(), s.expectedAnime)
		require.NoError(err)

		ttl := s.redisServer.TTL(testURL)
		expectedTTL, err := time.ParseDuration(sevenDaysInHours)
		require.NoError(err)
		require.Equal(expectedTTL, ttl)
	})
}
