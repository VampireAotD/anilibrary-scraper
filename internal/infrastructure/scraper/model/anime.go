package model

import (
	"errors"
	"fmt"
	"strings"

	"anilibrary-scraper/internal/domain/entity"

	"github.com/go-playground/validator/v10"
)

type Status = entity.Status
type Type = entity.Type

const (
	Ongoing = entity.Ongoing
	Ready   = entity.Ready

	Show  = entity.Show
	Movie = entity.Movie
)

type Anime struct {
	Image       string `validate:"required"`
	Title       string `validate:"required"`
	Status      Status `validate:"required,oneof=Анонс Онгоинг Вышел"`
	Type        Type   `validate:"required,oneof='ТВ Сериал' Фильм"`
	Genres      []string
	VoiceActing []string
	Synonyms    []string
	Episodes    int
	Year        int     `validate:"required,gt=0"`
	Rating      float32 `validate:"omitempty,gte=0,lte=10"`
}

func (a *Anime) Validate(validate *validator.Validate) error {
	if err := validate.Struct(a); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			var builder strings.Builder

			// Total length of validation message will be 87 if all validations fail
			builder.Grow(87)

			for i := range errs {
				builder.WriteString(errs[i].Field())
				builder.WriteString(" - ")
				builder.WriteString(errs[i].Tag())
				if i != len(errs)-1 {
					builder.WriteString("; ")
				}
			}

			return errors.New(builder.String())
		}

		return fmt.Errorf("validator: %w", err)
	}

	return nil
}

func (a *Anime) MapToDomainEntity() entity.Anime {
	return entity.Anime{
		Image:       a.Image,
		Title:       a.Title,
		Status:      a.Status,
		Episodes:    a.Episodes,
		Genres:      a.Genres,
		VoiceActing: a.VoiceActing,
		Synonyms:    a.Synonyms,
		Rating:      a.Rating,
		Year:        a.Year,
		Type:        a.Type,
	}
}
