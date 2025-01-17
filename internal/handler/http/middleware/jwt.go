package middleware

import (
	"anilibrary-scraper/internal/config"

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
