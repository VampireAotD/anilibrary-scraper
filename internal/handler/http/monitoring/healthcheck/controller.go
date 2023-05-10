package healthcheck

import (
	"context"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Controller struct {
	redisConnection *redis.Client
	kafkaConnection *kafka.Conn
}

func NewController(redisConnection *redis.Client, kafkaConnection *kafka.Conn) Controller {
	return Controller{
		redisConnection: redisConnection,
		kafkaConnection: kafkaConnection,
	}
}

func (c Controller) Healthcheck(w http.ResponseWriter, _ *http.Request) {
	if err := c.redisConnection.Ping(context.Background()).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := c.kafkaConnection.ReadPartitions(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
