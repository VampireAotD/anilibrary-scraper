//go:build integration

package kafka

import (
	"testing"
	"time"

	kafkaRepository "github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/kafka"
	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/model"

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

func (s *EventRepositorySuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	kafkaContainer, err := testcontainers.GenericContainer(s.T().Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "bitnami/kafka:3.9.0",
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
				"KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE":      "true",
			},
			WaitingFor: wait.ForListeningPort("9092/tcp").
				WithPollInterval(250 * time.Millisecond).
				WithStartupTimeout(time.Minute),
		},
		Started: true,
	})
	s.Require().NoError(err)

	ip, err := kafkaContainer.ContainerIP(s.T().Context())
	s.Require().NoError(err)

	topicName := "test_topic"

	// To fix https://github.com/segmentio/kafka-go/issues/683
	conn, err := kafka.Dial("tcp", ip+":9092")
	s.Require().NoError(err)
	defer func() {
		s.Require().NoError(conn.Close())
	}()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topicName,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = conn.CreateTopics(topicConfigs...)
	s.Require().NoError(err)

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(ip + ":9092"),
		Topic:                  topicName,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	s.kafkaContainer = kafkaContainer
	s.eventRepository = kafkaRepository.NewEventRepository(writer)
}

func (s *EventRepositorySuite) TearDownSuite() {
	s.Require().NoError(s.kafkaContainer.Stop(s.T().Context(), nil))
	s.Require().NoError(s.kafkaContainer.Terminate(s.T().Context()))
}

func (s *EventRepositorySuite) TestSend() {
	const testURL string = "https://google.com/"

	err := s.eventRepository.Send(s.T().Context(), model.Event{
		URL:       testURL,
		Timestamp: time.Now().Unix(),
		IP:        "0.0.0.0",
		UserAgent: "Mozilla/5.0",
	})
	s.Require().NoError(err)
}
