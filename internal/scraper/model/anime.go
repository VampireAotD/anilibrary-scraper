package model

import (
	"errors"

	"anilibrary-scraper/internal/entity"

	"github.com/go-playground/validator/v10"
)

var (
	ErrNotEnoughData = errors.New("entity wasn't filled with required data")
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
	Image       string   `validate:"required,url"`
	Title       string   `validate:"required"`
	Status      Status   `validate:"required,oneof=Онгоинг Вышел"`
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
		// TODO add more informative error
		return ErrNotEnoughData
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
