package anime

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/go-playground/validator/v10"
)

var (
	ErrUnableToDecodeRequest = errors.New("unable to decode incoming request")
	ErrInvalidURL            = errors.New("invalid url")
)

type ScrapeRequest struct {
	URL string `json:"url" validate:"required,url"`
}

func (request *ScrapeRequest) MapAndValidate(content io.ReadCloser) error {
	decoder := json.NewDecoder(content)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(request); err != nil {
		return ErrUnableToDecodeRequest
	}

	if validator.New().Struct(request) != nil {
		return ErrInvalidURL
	}

	return nil
}
