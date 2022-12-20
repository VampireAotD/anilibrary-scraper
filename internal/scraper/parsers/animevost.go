package parsers

import (
	"strconv"
	"strings"

	"anilibrary-scraper/internal/scraper/parsers/model"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var _ Contract = (*AnimeVost)(nil)

type AnimeVost struct{}

const amountToMakeFloat int = 10

func NewAnimeVost() AnimeVost {
	return AnimeVost{}
}

func (a AnimeVost) Title(document *goquery.Document) string {
	if title := document.Find(".shortstoryHead h1, .infoContent h3").First().Text(); title != "" {
		raw := strings.TrimSpace(title)
		return raw[0:strings.Index(raw, " /")]
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
	if episodes := document.Find("p strong:contains(Количество)").Parent().Children().Remove().End().Text(); episodes != "" {
		episodes = strings.Replace(episodes, "+", "", 1)
		end := strings.Index(episodes, " (")
		if end < 0 {
			end = len(episodes)
		}

		return episodes[0:end]
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

func (a AnimeVost) VoiceActing(document *goquery.Document) []string {
	return []string{"AnimeVost"}
}

func (a AnimeVost) Synonyms(document *goquery.Document) []string {
	if title := document.Find(".shortstoryHead h1, .infoContent h3").First().Text(); title != "" {
		text := strings.TrimSpace(title)
		start := strings.Index(text, "/ ")
		end := strings.Index(text, " [")

		if start < 0 || end < 0 {
			return nil
		}

		synonym := strings.TrimLeft(text[start:end], "/ ")

		return []string{synonym}
	}

	return nil
}

func (a AnimeVost) Image(document *goquery.Document) string {
	if attr, exists := document.Find(".imgRadius, .infoContent img").First().Attr("src"); exists {
		return AnimeVostURL + attr
	}

	return ""
}
