package healthcheck

import (
	"github.com/gofiber/fiber/v2"
)

type Controller struct{}

func NewController() Controller {
	return Controller{}
}

func (c Controller) Healthcheck(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
