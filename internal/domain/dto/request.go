package dto

import (
	"github.com/go-playground/validator/v10"
)

type RequestDTO struct {
	Url       string `json:"url" validate:"url"`
	FromCache bool
}

func (dto RequestDTO) Validate() error {
	return validator.New().Struct(dto)
}
