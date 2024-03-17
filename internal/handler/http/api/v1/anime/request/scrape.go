package request

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

func (request *ScrapeRequest) MapAndValidate(c *fiber.Ctx, validate *validator.Validate) error {
	if err := c.BodyParser(&request); err != nil {
		return ErrUnableToDecodeRequest
	}

	if validate.Struct(request) != nil {
		return ErrInvalidURL
	}

	return nil
}
