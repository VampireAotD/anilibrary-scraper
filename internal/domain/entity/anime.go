package entity

import "errors"

var ErrAnimeNotFound = errors.New("anime not found")

type Status string
type Type string

const (
	_       Status = "Анонс"
	Ongoing Status = "Онгоинг"
	Ready   Status = "Вышел"

	Show  Type = "ТВ Сериал"
	Movie Type = "Фильм"
)

type Anime struct {
	Image       string
	Title       string
	Status      Status
	Type        Type
	Episodes    string
	Genres      []string
	VoiceActing []string
	Synonyms    []string
	Rating      float32
	Year        int
}
