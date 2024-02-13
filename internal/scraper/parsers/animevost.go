package parsers

import (
	"regexp"
	"strconv"
	"strings"

	"anilibrary-scraper/internal/scraper/model"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const amountToMakeFloat int = 10

var _ Parser = (*AnimeVost)(nil)

type AnimeVost struct{}

func NewAnimeVost() AnimeVost {
	return AnimeVost{}
}

func (a AnimeVost) Title(document *goquery.Document) string {
	if title := document.Find(".shortstoryHead h1, .infoContent h3").First().Text(); title != "" {
		raw := strings.TrimSpace(title)
		return raw[:strings.Index(raw, " /")]
	}

	return ""
}

func (a AnimeVost) Status(document *goquery.Document) model.Status {
	if status := document.Find("#nexttime").Text(); status != "" {
		return model.Ongoing
	}

	return model.Ready
}

func (a AnimeVost) Rating(document *goquery.Document) float32 {
	if rating := document.Find(".current-rating").Text(); rating != "" {
		value, err := strconv.Atoi(rating)
		if err != nil {
			return MinimalAnimeRating
		}

		return float32(value / amountToMakeFloat)
	}

	return MinimalAnimeRating
}

func (a AnimeVost) Episodes(document *goquery.Document) string {
	if episodesText := document.Find("p strong:contains(Количество)").Parent().Children().Remove().End().Text(); episodesText != "" {
		regex := regexp.MustCompile(`^\d+`)
		return regex.FindString(episodesText)
	}

	return MinimalAnimeEpisodes
}

func (a AnimeVost) Genres(document *goquery.Document) []string {
	if genres := document.Find("p strong:contains(Жанр)").Parent().Text(); genres != "" {
		genres = strings.Replace(genres, "Жанр: ", "", 1)

		return strings.Split(cases.Title(language.Russian).String(genres), ", ")
	}

	return nil
}

func (a AnimeVost) VoiceActing(_ *goquery.Document) []string {
	return []string{"AnimeVost"}
}

func (a AnimeVost) Synonyms(document *goquery.Document) []string {
	if title := document.Find(".shortstoryHead h1, .infoContent h3").First().Text(); title != "" {
		regex := regexp.MustCompile(`/\s*(.*?)\s*\[`)
		// If there is a synonym, than the correct one without any symbols will be at index 1
		if entries := regex.FindStringSubmatch(title); len(entries) > 1 {
			return []string{entries[1]}
		}

		return nil
	}

	return nil
}

func (a AnimeVost) Image(document *goquery.Document) string {
	if attr, exists := document.Find(".imgRadius, .infoContent img").First().Attr("src"); exists {
		return AnimeVostURL + attr
	}

	return ""
}
