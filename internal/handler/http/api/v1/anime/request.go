package anime

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrUnableToDecodeRequest = errors.New("unable to decode incoming request")
	ErrInvalidURL            = errors.New("invalid url")
)

type ScrapeRequest struct {
	URL string `json:"url" validate:"required,url"`
}

func (request *ScrapeRequest) MapAndValidate(c *fiber.Ctx) error {
	if err := c.BodyParser(&request); err != nil {
		return ErrUnableToDecodeRequest
	}

	if validator.New().Struct(request) != nil {
		return ErrInvalidURL
	}

	return nil
}
