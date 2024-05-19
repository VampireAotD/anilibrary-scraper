package healthcheck

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Controller struct {
	redisConnection redis.UniversalClient
	kafkaConnection *kafka.Conn
}

func NewController(redisConnection redis.UniversalClient, kafkaConnection *kafka.Conn) Controller {
	return Controller{
		redisConnection: redisConnection,
		kafkaConnection: kafkaConnection,
	}
}

func (c Controller) Healthcheck(ctx *fiber.Ctx) error {
	if err := c.redisConnection.Ping(ctx.UserContext()).Err(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if _, err := c.kafkaConnection.ReadPartitions(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusOK)
}
