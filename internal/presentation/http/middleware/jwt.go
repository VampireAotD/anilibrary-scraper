package middleware

import (
	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func NewJWTAuth(cfg config.JWT) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: cfg.Secret,
		},
	})
}
