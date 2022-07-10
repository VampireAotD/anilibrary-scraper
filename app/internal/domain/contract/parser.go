package contract

import "anilibrary-request-parser/app/internal/domain/entity"

type Parser interface {
	GetTitle() string
	GetStatus() string
	GetRating() float32
	GetEpisodes() uint16
	GetGenres() []string
	GetVoiceActing() []string
	GetAnime() entity.Anime
}
