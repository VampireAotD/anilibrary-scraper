package animevost

import (
	"anilibrary-request-parser/app/internal/domain/entity"
	animeEnum "anilibrary-request-parser/app/internal/domain/enum/anime"
	"anilibrary-request-parser/app/internal/infrastructure/scraper"
	"github.com/PuerkitoBio/goquery"
)

type AnimeVost struct {
	*scraper.Scrapper
}

func New(scrapper *scraper.Scrapper) *AnimeVost {
	return &AnimeVost{Scrapper: scrapper}
}

func (a AnimeVost) Title(document *goquery.Document) string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) Status(document *goquery.Document) animeEnum.Status {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) Rating(document *goquery.Document) float32 {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) Episodes(document *goquery.Document) string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) Genres(document *goquery.Document) []string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) VoiceActing(document *goquery.Document) []string {
	//TODO implement me
	panic("implement me")
}

func (a AnimeVost) Process() (*entity.Anime, error) {
	//TODO implement me
	panic("implement me")
}
