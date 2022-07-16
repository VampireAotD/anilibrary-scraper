package animevost

import (
	"anilibrary-request-parser/app/internal/domain/entity"
	animeEnum "anilibrary-request-parser/app/internal/domain/enum/anime"
	"github.com/PuerkitoBio/goquery"
)

type AnimeVost struct {
}

func New() *AnimeVost {
	return &AnimeVost{}
}

func (a AnimeVost) GetTitle(document *goquery.Document) string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetStatus(document *goquery.Document) animeEnum.Status {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetRating(document *goquery.Document) float32 {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetEpisodes(document *goquery.Document) string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetGenres(document *goquery.Document) []string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetVoiceActing(document *goquery.Document) []string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) GetAnime(document *goquery.Document) entity.Anime {
	//TODO implement me
	panic("implement me")
}
