//go:build integration

package healthcheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VampireAotD/anilibrary-scraper/internal/presentation/http/monitoring/healthcheck"

	"github.com/alicebob/miniredis/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const endpoint string = "/healthcheck"

type HealthcheckControllerSuite struct {
	suite.Suite

	redisServer    *miniredis.Miniredis
	kafkaContainer testcontainers.Container
	controller     healthcheck.Controller
	router         *fiber.App
}

func TestHealthcheckControllerSuite(t *testing.T) {
	suite.Run(t, new(HealthcheckControllerSuite))
}

func (s *HealthcheckControllerSuite) initRedis() {
	s.redisServer = miniredis.RunT(s.T())
}

func (s *HealthcheckControllerSuite) initKafka() {
	kafkaCtx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	kafkaContainer, err := testcontainers.GenericContainer(kafkaCtx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "bitnami/kafka",
			ExposedPorts: []string{
				"9095:9095/tcp",
			},
			Env: map[string]string{
				"KAFKA_CFG_NODE_ID":                        "1",
				"KAFKA_CFG_PROCESS_ROLES":                  "controller,broker",
				"ALLOW_PLAINTEXT_LISTENER":                 "yes",
				"KAFKA_CFG_CONTROLLER_LISTENER_NAMES":      "CONTROLLER",
				"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP": "TEST:PLAINTEXT,PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT",
				"KAFKA_CFG_CONTROLLER_QUORUM_VOTERS":       "1@localhost:9093",
				"KAFKA_CFG_LISTENERS":                      "TEST://:9095,PLAINTEXT://:9092,CONTROLLER://:9093",
				"KAFKA_CFG_ADVERTISED_LISTENERS":           "TEST://localhost:9095,PLAINTEXT://localhost:9092",
			},
			WaitingFor: wait.ForListeningPort("9095/tcp").
				WithPollInterval(time.Millisecond * 100).
				WithStartupTimeout(time.Minute),
		},
		Started: true,
	})
	s.Require().NoError(err)

	s.kafkaContainer = kafkaContainer
}

func (s *HealthcheckControllerSuite) SetupSuite() {
	s.initRedis()
	s.initKafka()
	s.router = fiber.New()

	// Setup Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: s.redisServer.Addr(),
	})

	kafkaIP, err := s.kafkaContainer.ContainerIP(context.Background())
	s.Require().NoError(err)

	// Dial Kafka connection
	kafkaConnection, err := kafka.DialLeader(context.Background(), "tcp", kafkaIP+":9095", "test-topic", 0)
	s.Require().NoError(err)

	s.controller = healthcheck.NewController(redisClient, kafkaConnection)
	s.router.Get(endpoint, s.controller.Healthcheck)
}

func (s *HealthcheckControllerSuite) SetupTest() {
	if s.redisServer.Addr() == "" {
		s.Require().NoError(s.redisServer.Start())
	}

	if !s.kafkaContainer.IsRunning() {
		s.Require().NoError(s.kafkaContainer.Start(context.Background()))
	}
}

func (s *HealthcheckControllerSuite) TearDownSuite() {
	if s.kafkaContainer.IsRunning() {
		s.Require().NoError(s.kafkaContainer.Terminate(context.Background()))
	}

	s.redisServer.Close()
}

func (s *HealthcheckControllerSuite) sendHealthcheckRequest() *http.Response {
	request := httptest.NewRequest(http.MethodGet, endpoint, http.NoBody)

	response, err := s.router.Test(request, -1)
	s.Require().NoError(err)

	return response
}

func (s *HealthcheckControllerSuite) TestHealthcheck() {
	var require = s.Require()

	s.T().Run("Redis", func(t *testing.T) {
		t.Run("Redis up", func(_ *testing.T) {
			response := s.sendHealthcheckRequest()
			require.Equal(http.StatusOK, response.StatusCode)
			require.NoError(response.Body.Close())
		})

		t.Run("Redis down", func(_ *testing.T) {
			s.redisServer.Close()

			response := s.sendHealthcheckRequest()
			require.Equal(http.StatusInternalServerError, response.StatusCode)
			require.NoError(s.redisServer.Start())
			require.NoError(response.Body.Close())
		})
	})

	s.T().Run("Kafka", func(t *testing.T) {
		t.Run("Kafka up", func(_ *testing.T) {
			response := s.sendHealthcheckRequest()
			require.Equal(http.StatusOK, response.StatusCode)
			require.NoError(response.Body.Close())
		})

		t.Run("Kafka down", func(t *testing.T) {
			require.NoError(s.kafkaContainer.Stop(t.Context(), nil))

			response := s.sendHealthcheckRequest()
			require.Equal(http.StatusInternalServerError, response.StatusCode)
			require.NoError(response.Body.Close())
		})
	})
}
