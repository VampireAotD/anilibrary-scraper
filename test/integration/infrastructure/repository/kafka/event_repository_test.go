//go:build integration

package kafka

import (
	"testing"
	"time"

	kafkaRepository "github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/kafka"
	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/model"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type EventRepositorySuite struct {
	suite.Suite

	client          *kgo.Client
	container       testcontainers.Container
	eventRepository kafkaRepository.EventRepository
}

func TestEventRepositorySuite(t *testing.T) {
	suite.Run(t, new(EventRepositorySuite))
}

func (s *EventRepositorySuite) SetupSuite() {
	kafkaContainer, err := testcontainers.GenericContainer(s.T().Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "bitnami/kafka:3.9.0",
			ExposedPorts: []string{
				"9092:9092/tcp",
			},
			Env: map[string]string{
				"KAFKA_CFG_NODE_ID":                        "0",
				"KAFKA_KRAFT_CLUSTER_ID":                   "anilibrary-scraper-test",
				"KAFKA_CFG_PROCESS_ROLES":                  "broker,controller",
				"KAFKA_CFG_CONTROLLER_QUORUM_VOTERS":       "0@localhost:9093",
				"ALLOW_PLAINTEXT_LISTENER":                 "yes",
				"KAFKA_CFG_LISTENERS":                      "INTERNAL://:9092,CONTROLLER://:9093",
				"KAFKA_CFG_ADVERTISED_LISTENERS":           "INTERNAL://localhost:9092",
				"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP": "INTERNAL:PLAINTEXT,CONTROLLER:PLAINTEXT",
				"KAFKA_CFG_INTER_BROKER_LISTENER_NAME":     "INTERNAL",
				"KAFKA_CFG_CONTROLLER_LISTENER_NAMES":      "CONTROLLER",
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

	client, err := kgo.NewClient(
		kgo.SeedBrokers(ip+":9092"),
		kgo.ClientID("test"),
		kgo.DefaultProduceTopic("test_topic"),
		kgo.ProducerBatchMaxBytes(1024*1024),
		kgo.ProducerBatchCompression(kgo.ZstdCompression()),
		kgo.AllowAutoTopicCreation(),
	)
	s.Require().NoError(err)

	s.client = client
	s.container = kafkaContainer
	s.eventRepository = kafkaRepository.NewEventRepository(s.client)
}

func (s *EventRepositorySuite) TearDownSuite() {
	s.client.Close()
	s.Require().NoError(s.container.Terminate(s.T().Context()))
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
