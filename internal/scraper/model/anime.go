package model

import (
	"errors"

	"anilibrary-scraper/internal/entity"
)

type Status string

const (
	Ongoing Status = "Онгоинг"
	Ready   Status = "Вышел"
)

var (
	ErrNotEnoughData = errors.New("entity wasn't filled with required data")
)

type Anime struct {
	Image       string
	Title       string
	Status      Status
	Episodes    string
	Genres      []string
	VoiceActing []string
	Synonyms    []string
	Rating      float32
}

func (a Anime) Validate() error {
	if a.Image == "" || a.Title == "" {
		return ErrNotEnoughData
	}

	return nil
}

func (a Anime) MapToDomainEntity() entity.Anime {
	return entity.Anime{
		Image:       a.Image,
		Title:       a.Title,
		Status:      string(a.Status),
		Episodes:    a.Episodes,
		Genres:      a.Genres,
		VoiceActing: a.VoiceActing,
		Synonyms:    a.Synonyms,
		Rating:      a.Rating,
	}
}
