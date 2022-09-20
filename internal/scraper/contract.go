package scraper

import (
	"anilibrary-scraper/internal/domain/enum"
	"github.com/PuerkitoBio/goquery"
)

type Contract interface {
	Title(document *goquery.Document) string
	Status(document *goquery.Document) enum.Status
	Rating(document *goquery.Document) float32
	Episodes(document *goquery.Document) string
	Genres(document *goquery.Document) []string
	VoiceActing(document *goquery.Document) []string
}
