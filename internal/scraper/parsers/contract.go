package parsers

import (
	"anilibrary-scraper/internal/scraper/parsers/model"
	"github.com/PuerkitoBio/goquery"
)

// Urls to scrape
const (
	AnimeGoUrl   string = "https://animego.org"
	AnimeVostUrl string = "https://animevost.org"
)

// Default values
const (
	MinimalAnimeRating   float32 = 0
	MinimalAnimeEpisodes string  = "0 / ?"
)

type Contract interface {
	// Title method scraping anime title and returns empty string if none found
	Title(document *goquery.Document) string

	// Status method scraping current anime status
	Status(document *goquery.Document) model.Status

	// Rating method scraping current anime rating and returns parsers.MinimalAnimeRating if none found
	Rating(document *goquery.Document) float32

	// Episodes method scraping amount of anime episodes and returns parsers.MinimalAnimeEpisodes if none found
	Episodes(document *goquery.Document) string

	// Genres method scraping all anime genres
	Genres(document *goquery.Document) []string

	// VoiceActing method scraping all anime voice acting
	VoiceActing(document *goquery.Document) []string

	// Image method scraping image url returns empty string if none found
	Image(document *goquery.Document) string
}
