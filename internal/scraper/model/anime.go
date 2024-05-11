package model

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"anilibrary-scraper/internal/entity"

	"github.com/go-playground/validator/v10"
)

var validationPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

type Status = entity.Status
type Type = entity.Type

const (
	Ongoing = entity.Ongoing
	Ready   = entity.Ready

	Show  = entity.Show
	Movie = entity.Movie
)

type Anime struct {
	Image       string   `validate:"required"`
	Title       string   `validate:"required"`
	Status      Status   `validate:"required,oneof=Анонс Онгоинг Вышел"`
	Type        Type     `validate:"required,oneof='ТВ Сериал' Фильм"`
	Episodes    string   // validation not required
	Genres      []string // validation not required
	VoiceActing []string // validation not required
	Synonyms    []string // validation not required
	Rating      float32  `validate:"omitempty,gte=0,lte=10"`
	Year        int      `validate:"required,gt=0"`
}

func (a Anime) Validate(validate *validator.Validate) error {
	if err := validate.Struct(a); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			buf, _ := validationPool.Get().(*bytes.Buffer)
			defer func() {
				buf.Reset()
				validationPool.Put(buf)
			}()

			for i := range errs {
				buf.WriteString(errs[i].Field())
				buf.WriteString(" - ")
				buf.WriteString(errs[i].Tag())
				if i != len(errs)-1 {
					buf.WriteString("; ")
				}
			}

			return errors.New(buf.String())
		}

		return fmt.Errorf("validator: %w", err)
	}

	return nil
}

func (a Anime) MapToDomainEntity() entity.Anime {
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
