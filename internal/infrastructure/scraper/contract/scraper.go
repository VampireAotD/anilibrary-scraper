package contract

import (
	animeEnum "anilibrary-request-parser/internal/domain/enum/anime"
	"github.com/PuerkitoBio/goquery"
)

type Scraper interface {
	Title(document *goquery.Document) string
	Status(document *goquery.Document) animeEnum.Status
	Rating(document *goquery.Document) float32
	Episodes(document *goquery.Document) string
	Genres(document *goquery.Document) []string
	VoiceActing(document *goquery.Document) []string
}
