//go:build integration

package kafka

import (
	"context"
	"testing"
	"time"

	kafka2 "anilibrary-scraper/internal/domain/repository/kafka"
	"anilibrary-scraper/internal/domain/repository/models"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/mock/gomock"
)

type EventRepositorySuite struct {
	suite.Suite

	kafkaContainer  testcontainers.Container
	eventRepository kafka2.EventRepository
}

func TestEventRepositorySuite(t *testing.T) {
	suite.Run(t, new(EventRepositorySuite))
}

func (suite *EventRepositorySuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	kafkaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "bitnami/kafka",
			ExposedPorts: []string{
				"9092:9092/tcp",
			},
			Env: map[string]string{
				"KAFKA_CFG_NODE_ID":                        "1",
				"KAFKA_CFG_PROCESS_ROLES":                  "controller,broker",
				"ALLOW_PLAINTEXT_LISTENER":                 "yes",
				"KAFKA_CFG_CONTROLLER_LISTENER_NAMES":      "CONTROLLER",
				"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP": "PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT",
				"KAFKA_CFG_CONTROLLER_QUORUM_VOTERS":       "1@localhost:9093",
				"KAFKA_CFG_LISTENERS":                      "PLAINTEXT://:9092,CONTROLLER://:9093",
				"KAFKA_CFG_ADVERTISED_LISTENERS":           "PLAINTEXT://localhost:9092",
			},
			WaitingFor: wait.ForListeningPort("9092/tcp").
				WithPollInterval(time.Millisecond * 100).
				WithStartupTimeout(time.Minute),
		},
		Started: true,
	})
	suite.Require().NoError(err)

	ip, err := kafkaContainer.ContainerIP(ctx)
	suite.Require().NoError(err)

	conn, err := kafka.DialLeader(ctx, "tcp", ip+":9092", "test_topic", 0)
	suite.Require().NoError(err)

	suite.kafkaContainer = kafkaContainer
	suite.eventRepository = kafka2.NewEventRepository(conn)
}

func (suite *EventRepositorySuite) TearDownSuite() {
	suite.Require().NoError(suite.kafkaContainer.Stop(context.Background(), nil))
	suite.Require().NoError(suite.kafkaContainer.Terminate(context.Background()))
}

func (suite *EventRepositorySuite) TestSend() {
	const testURL string = "https://google.com/"

	err := suite.eventRepository.Send(context.Background(), models.Event{
		URL:  testURL,
		Date: time.Now().Unix(),
	})
	suite.Require().NoError(err)
}
