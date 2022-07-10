package animego

import (
	"anilibrary-request-parser/app/internal/domain/entity"
)

type AnimeGo struct {
}

func New() *AnimeGo {
	return &AnimeGo{}
}

func (a AnimeGo) GetTitle() string {
	return "zzzzzzzz"
}

func (a AnimeGo) GetStatus() string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeGo) GetRating() float32 {
	//TODO implement me
	panic("implement me")
}

func (a AnimeGo) GetEpisodes() uint16 {
	//TODO implement me
	panic("implement me")
}

func (a AnimeGo) GetGenres() []string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeGo) GetVoiceActing() []string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeGo) GetAnime() entity.Anime {
	//TODO implement me
	panic("implement me")
}
