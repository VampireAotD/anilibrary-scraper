//go:build integration

package kafka

import (
	"context"
	"testing"
	"time"

	kafkaRepository "anilibrary-scraper/internal/infrastructure/repository/kafka"
	"anilibrary-scraper/internal/infrastructure/repository/model"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/mock/gomock"
)

type EventRepositorySuite struct {
	suite.Suite

	kafkaContainer  testcontainers.Container
	eventRepository kafkaRepository.EventRepository
}

func TestEventRepositorySuite(t *testing.T) {
	suite.Run(t, new(EventRepositorySuite))
}

func (ers *EventRepositorySuite) SetupSuite() {
	ctrl := gomock.NewController(ers.T())
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
	ers.Require().NoError(err)

	ip, err := kafkaContainer.ContainerIP(ctx)
	ers.Require().NoError(err)

	conn, err := kafka.DialLeader(ctx, "tcp", ip+":9092", "test_topic", 0)
	ers.Require().NoError(err)

	ers.kafkaContainer = kafkaContainer
	ers.eventRepository = kafkaRepository.NewEventRepository(conn)
}

func (ers *EventRepositorySuite) TearDownSuite() {
	ers.Require().NoError(ers.kafkaContainer.Stop(context.Background(), nil))
	ers.Require().NoError(ers.kafkaContainer.Terminate(context.Background()))
}

func (ers *EventRepositorySuite) TestSend() {
	const testURL string = "https://google.com/"

	err := ers.eventRepository.Send(context.Background(), model.Event{
		URL:       testURL,
		Timestamp: time.Now().Unix(),
		IP:        "0.0.0.0",
		UserAgent: "Mozilla/5.0",
	})
	ers.Require().NoError(err)
}
