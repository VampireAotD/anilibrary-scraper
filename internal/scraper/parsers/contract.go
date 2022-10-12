package parsers

import (
	"anilibrary-scraper/internal/domain/entity"
	"github.com/PuerkitoBio/goquery"
)

type Contract interface {
	Title(document *goquery.Document) string
	Status(document *goquery.Document) entity.Status
	Rating(document *goquery.Document) float32
	Episodes(document *goquery.Document) string
	Genres(document *goquery.Document) []string
	VoiceActing(document *goquery.Document) []string
	Image(document *goquery.Document) string
}
