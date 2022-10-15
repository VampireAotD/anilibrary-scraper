package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var ErrInvalidUrl = errors.New("invalid url")

type RequestDTO struct {
	Url       string `json:"url" validate:"required,url"`
	FromCache bool
}

func (dto RequestDTO) Validate() error {
	if validator.New().Struct(dto) != nil {
		return ErrInvalidUrl
	}

	return nil
}
