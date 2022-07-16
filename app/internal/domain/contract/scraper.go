package contract

import (
	"anilibrary-request-parser/app/internal/domain/entity"
	animeEnum "anilibrary-request-parser/app/internal/domain/enum/anime"
	"github.com/PuerkitoBio/goquery"
)

type Scraper interface {
	GetTitle(document *goquery.Document) string
	GetStatus(document *goquery.Document) animeEnum.Status
	GetRating(document *goquery.Document) float32
	GetEpisodes(document *goquery.Document) string
	GetGenres(document *goquery.Document) []string
	GetVoiceActing(document *goquery.Document) []string
	GetAnime(document *goquery.Document) entity.Anime
}
