package anime

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var ErrInvalidURL = errors.New("invalid url")

type ScrapeRequest struct {
	URL string `json:"url" validate:"required,url"`
}

func (request ScrapeRequest) Validate() error {
	if validator.New().Struct(request) != nil {
		return ErrInvalidURL
	}

	return nil
}
