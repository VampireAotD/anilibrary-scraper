//go:build integration

package healthcheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"anilibrary-scraper/internal/handler/http/monitoring/healthcheck"

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

func (suite *HealthcheckControllerSuite) initRedis() {
	suite.redisServer = miniredis.RunT(suite.T())
}

func (suite *HealthcheckControllerSuite) initKafka() {
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
	suite.Require().NoError(err)

	suite.kafkaContainer = kafkaContainer
}

func (suite *HealthcheckControllerSuite) SetupSuite() {
	suite.initRedis()
	suite.initKafka()
	suite.router = fiber.New()

	// Setup Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: suite.redisServer.Addr(),
	})

	kafkaIP, err := suite.kafkaContainer.ContainerIP(context.Background())
	suite.Require().NoError(err)

	// Dial Kafka connection
	kafkaConnection, err := kafka.DialLeader(context.Background(), "tcp", kafkaIP+":9095", "test-topic", 0)
	suite.Require().NoError(err)

	suite.controller = healthcheck.NewController(redisClient, kafkaConnection)
	suite.router.Get(endpoint, suite.controller.Healthcheck)
}

func (suite *HealthcheckControllerSuite) SetupTest() {
	if suite.redisServer.Addr() == "" {
		suite.Require().NoError(suite.redisServer.Start())
	}

	if !suite.kafkaContainer.IsRunning() {
		suite.Require().NoError(suite.kafkaContainer.Start(context.Background()))
	}
}

func (suite *HealthcheckControllerSuite) TearDownSuite() {
	if suite.kafkaContainer.IsRunning() {
		suite.Require().NoError(suite.kafkaContainer.Terminate(context.Background()))
	}

	suite.redisServer.Close()
}

func (suite *HealthcheckControllerSuite) sendHealthcheckRequest() *http.Response {
	request := httptest.NewRequest(http.MethodGet, endpoint, nil)

	response, err := suite.router.Test(request, -1)
	suite.Require().NoError(err)

	return response
}

func (suite *HealthcheckControllerSuite) TestHealthcheck() {
	var (
		t       = suite.T()
		require = suite.Require()
	)

	t.Run("Redis", func(t *testing.T) {
		t.Run("Redis up", func(t *testing.T) {
			response := suite.sendHealthcheckRequest()
			require.Equal(http.StatusOK, response.StatusCode)
		})

		t.Run("Redis down", func(t *testing.T) {
			suite.redisServer.Close()

			response := suite.sendHealthcheckRequest()
			require.Equal(http.StatusInternalServerError, response.StatusCode)

			require.NoError(suite.redisServer.Start())
		})
	})

	t.Run("Kafka", func(t *testing.T) {
		t.Run("Kafka up", func(t *testing.T) {
			response := suite.sendHealthcheckRequest()
			require.Equal(http.StatusOK, response.StatusCode)
		})

		t.Run("Kafka down", func(t *testing.T) {
			require.NoError(suite.kafkaContainer.Stop(context.Background(), nil))

			response := suite.sendHealthcheckRequest()
			require.Equal(http.StatusInternalServerError, response.StatusCode)
		})
	})
}
