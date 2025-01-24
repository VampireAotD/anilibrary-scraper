package middleware

import (
	"github.com/VampireAotD/anilibrary-scraper/pkg/logging"

	"github.com/gofiber/fiber/v2"
)

func NewLogger(logger *logging.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctxLogger := logger.With(
			logging.String("type", "http"),
			logging.String("method", ctx.Method()),
			logging.String("path", ctx.Path()),
			logging.String("protocol", ctx.Protocol()),
		)

		ctx.SetUserContext(logging.ContextWithLogger(ctx.UserContext(), ctxLogger))

		return ctx.Next()
	}
}
