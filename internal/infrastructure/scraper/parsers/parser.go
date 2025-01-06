package parsers

import "regexp"

// Urls to scrape
const (
	AnimeGoURL   string = "https://animego.org"
	AnimeVostURL string = "https://animevost.org"
)

// Default values
const (
	MinimalAnimeEpisodes int     = 0
	MinimalAnimeRating   float32 = 0.0
)

// Common regexes for parsers
var (
	yearRegex = regexp.MustCompile(`\d{4}`)
)
