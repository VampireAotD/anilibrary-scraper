package anime

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

var ErrInvalidUrl = errors.New("invalid url")

type ScrapeRequest struct {
	Url string `json:"url" validate:"required,url"`
}

func (request ScrapeRequest) Validate() error {
	if validator.New().Struct(request) != nil {
		return ErrInvalidUrl
	}

	return nil
}
