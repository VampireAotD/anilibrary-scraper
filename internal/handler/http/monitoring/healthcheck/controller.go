package healthcheck

import (
	"context"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Controller struct {
	redisConnection *redis.Client
}

func NewController(redisConnection *redis.Client) Controller {
	return Controller{
		redisConnection: redisConnection,
	}
}

func (c Controller) Healthcheck(w http.ResponseWriter, _ *http.Request) {
	if err := c.redisConnection.Ping(context.Background()).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
