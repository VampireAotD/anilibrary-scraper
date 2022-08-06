package dto

import (
	"github.com/go-playground/validator/v10"
)

type ParseDTO struct {
	Url string `json:"url" validate:"url"`
}

func (dto *ParseDTO) Validate() error {
	return validator.New().Struct(dto)
}
