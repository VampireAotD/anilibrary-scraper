package redis

import (
	"context"
	"encoding/base64"
	"testing"

	"anilibrary-scraper/internal/repository/model"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

const testURL string = "https://animego.org/anime/naruto-uragannye-hroniki-103"

type AnimeRepositorySuite struct {
	suite.Suite

	redisServer     *miniredis.Miniredis
	animeRepository AnimeRepository
	expectedAnime   model.Anime
}

func TestAnimeRepositorySuite(t *testing.T) {
	suite.Run(t, new(AnimeRepositorySuite))
}

func (ars *AnimeRepositorySuite) SetupSuite() {
	ctrl := gomock.NewController(ars.T())
	defer ctrl.Finish()

	ars.redisServer = miniredis.RunT(ars.T())
	ars.animeRepository = NewAnimeRepository(redis.NewClient(&redis.Options{
		Addr: ars.redisServer.Addr(),
	}))
	ars.expectedAnime = model.Anime{
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
		anime, err := ars.animeRepository.FindByURL(context.Background(), testURL)
		require.Error(err)
		require.Zero(anime)
	})

	t.Run("Found in cache", func(_ *testing.T) {
		err := ars.animeRepository.Create(context.Background(), ars.expectedAnime)
		require.NoError(err)

		anime, err := ars.animeRepository.FindByURL(context.Background(), testURL)
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

	t.Run("Invalid cases", func(t *testing.T) {
		t.Run("Missing image", func(_ *testing.T) {
			err := ars.animeRepository.Create(context.Background(), model.Anime{
				URL:   testURL,
				Title: "random",
			})
			require.ErrorIs(err, model.ErrInvalidData)
		})

		t.Run("Missing title", func(_ *testing.T) {
			err := ars.animeRepository.Create(context.Background(), model.Anime{
				URL:   testURL,
				Image: base64.StdEncoding.EncodeToString([]byte("random")),
			})
			require.ErrorIs(err, model.ErrInvalidData)
		})
	})

	t.Run("Correct data", func(_ *testing.T) {
		err := ars.animeRepository.Create(context.Background(), ars.expectedAnime)
		require.NoError(err)
	})
}
