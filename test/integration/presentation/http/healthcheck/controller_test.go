//go:build integration

package healthcheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"anilibrary-scraper/internal/presentation/http/monitoring/healthcheck"

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

func (hcs *HealthcheckControllerSuite) initRedis() {
	hcs.redisServer = miniredis.RunT(hcs.T())
}

func (hcs *HealthcheckControllerSuite) initKafka() {
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
	hcs.Require().NoError(err)

	hcs.kafkaContainer = kafkaContainer
}

func (hcs *HealthcheckControllerSuite) SetupSuite() {
	hcs.initRedis()
	hcs.initKafka()
	hcs.router = fiber.New()

	// Setup Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: hcs.redisServer.Addr(),
	})

	kafkaIP, err := hcs.kafkaContainer.ContainerIP(context.Background())
	hcs.Require().NoError(err)

	// Dial Kafka connection
	kafkaConnection, err := kafka.DialLeader(context.Background(), "tcp", kafkaIP+":9095", "test-topic", 0)
	hcs.Require().NoError(err)

	hcs.controller = healthcheck.NewController(redisClient, kafkaConnection)
	hcs.router.Get(endpoint, hcs.controller.Healthcheck)
}

func (hcs *HealthcheckControllerSuite) SetupTest() {
	if hcs.redisServer.Addr() == "" {
		hcs.Require().NoError(hcs.redisServer.Start())
	}

	if !hcs.kafkaContainer.IsRunning() {
		hcs.Require().NoError(hcs.kafkaContainer.Start(context.Background()))
	}
}

func (hcs *HealthcheckControllerSuite) TearDownSuite() {
	if hcs.kafkaContainer.IsRunning() {
		hcs.Require().NoError(hcs.kafkaContainer.Terminate(context.Background()))
	}

	hcs.redisServer.Close()
}

func (hcs *HealthcheckControllerSuite) sendHealthcheckRequest() *http.Response {
	request := httptest.NewRequest(http.MethodGet, endpoint, http.NoBody)

	response, err := hcs.router.Test(request, -1)
	hcs.Require().NoError(err)

	return response
}

func (hcs *HealthcheckControllerSuite) TestHealthcheck() {
	var (
		t       = hcs.T()
		require = hcs.Require()
	)

	t.Run("Redis", func(t *testing.T) {
		t.Run("Redis up", func(_ *testing.T) {
			response := hcs.sendHealthcheckRequest()
			require.Equal(http.StatusOK, response.StatusCode)
			require.NoError(response.Body.Close())
		})

		t.Run("Redis down", func(_ *testing.T) {
			hcs.redisServer.Close()

			response := hcs.sendHealthcheckRequest()
			require.Equal(http.StatusInternalServerError, response.StatusCode)
			require.NoError(hcs.redisServer.Start())
			require.NoError(response.Body.Close())
		})
	})

	t.Run("Kafka", func(t *testing.T) {
		t.Run("Kafka up", func(_ *testing.T) {
			response := hcs.sendHealthcheckRequest()
			require.Equal(http.StatusOK, response.StatusCode)
			require.NoError(response.Body.Close())
		})

		t.Run("Kafka down", func(_ *testing.T) {
			require.NoError(hcs.kafkaContainer.Stop(context.Background(), nil))

			response := hcs.sendHealthcheckRequest()
			require.Equal(http.StatusInternalServerError, response.StatusCode)
			require.NoError(response.Body.Close())
		})
	})
}
