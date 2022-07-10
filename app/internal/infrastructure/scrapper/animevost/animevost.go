package animevost

import "anilibrary-request-parser/app/internal/domain/entity"

type AnimeVost struct {
}

func New() *AnimeVost {
	return &AnimeVost{}
}

func (a AnimeVost) GetTitle() string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetStatus() string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetRating() float32 {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetEpisodes() uint16 {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetGenres() []string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetVoiceActing() []string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetAnime() entity.Anime {
	//TODO implement me
	panic("implement me")
}
