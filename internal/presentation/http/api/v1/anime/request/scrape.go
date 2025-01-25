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

func (r *ScrapeRequest) Validate(c *fiber.Ctx, validate *validator.Validate) error {
	if err := c.BodyParser(&r); err != nil {
		return ErrUnableToDecodeRequest
	}

	if validate.Struct(r) != nil {
		return ErrInvalidURL
	}

	return nil
}
