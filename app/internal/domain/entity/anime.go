package entity

import (
	"anilibrary-request-parser/app/internal/domain/enum/anime"
)

type Anime struct {
	Title       string       `json:"title"`
	Status      anime.Status `json:"status"`
	Rating      float32      `json:"rating"`
	Episodes    string       `json:"episodes"`
	Genres      []string     `json:"genres"`
	VoiceActing []string     `json:"voice_acting"`
}
