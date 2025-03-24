//go:build integration

package redis

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/VampireAotD/anilibrary-scraper/internal/domain/entity"
	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/model"
	redisRepository "github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type AnimeRepositorySuite struct {
	suite.Suite

	client     redis.UniversalClient
	container  testcontainers.Container
	repository redisRepository.AnimeRepository
}

func TestAnimeRepositorySuite(t *testing.T) {
	suite.Run(t, new(AnimeRepositorySuite))
}

func (s *AnimeRepositorySuite) SetupSuite() {
	request := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "redis:7.4.2",
			ExposedPorts: []string{
				"6379:6379/tcp",
			},
			WaitingFor: wait.ForListeningPort("6379/tcp").
				WithPollInterval(250 * time.Millisecond).
				WithStartupTimeout(time.Minute),
		},
		Started: true,
	}

	redisContainer, err := testcontainers.GenericContainer(s.T().Context(), request)
	s.Require().NoError(err)

	ip, err := redisContainer.ContainerIP(s.T().Context())
	s.Require().NoError(err)

	port, err := redisContainer.MappedPort(s.T().Context(), "6379/tcp")
	s.Require().NoError(err)

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{ip + ":" + port.Port()},
	})

	s.Require().NoError(client.Ping(s.T().Context()).Err())

	s.client = client
	s.container = redisContainer
	s.repository = redisRepository.NewAnimeRepository(s.client)
}

func (s *AnimeRepositorySuite) TearDownSuite() {
	s.Require().NoError(s.client.Close())
	s.Require().NoError(s.container.Terminate(s.T().Context()))
}

func (s *AnimeRepositorySuite) TestCreate() {
	anime := model.Anime{
		URL:      "https://animego.org/anime/naruto-uragannye-hroniki-103",
		Image:    base64.StdEncoding.EncodeToString([]byte("data:image/png;base64,image")),
		Title:    "test",
		Status:   entity.Ready,
		Type:     entity.Show,
		Episodes: 120,
		Year:     time.Now().Year(),
	}

	err := s.repository.Create(s.T().Context(), anime)
	s.Require().NoError(err)

	ttl, err := s.client.TTL(s.T().Context(), anime.URL).Result()
	s.Require().NoError(err)

	expectedTTL, err := time.ParseDuration("168h")
	s.Require().NoError(err)
	s.Require().Equal(expectedTTL, ttl)
}

func (s *AnimeRepositorySuite) TestFindByURL() {
	anime := model.Anime{
		URL:      "https://animego.org/anime/naruto-uragannye-hroniki-103",
		Image:    base64.StdEncoding.EncodeToString([]byte("data:image/png;base64,image")),
		Title:    "test",
		Status:   entity.Ongoing,
		Type:     entity.Movie,
		Episodes: 20,
		Year:     time.Now().Year(),
	}

	err := s.repository.Create(s.T().Context(), anime)
	s.Require().NoError(err)

	found, err := s.repository.FindByURL(s.T().Context(), anime.URL)
	s.Require().NoError(err)
	s.Require().Equal(anime.MapToDomainEntity(), found)
}
