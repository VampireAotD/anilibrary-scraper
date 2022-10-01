package scraper

import (
	"strconv"
	"strings"

	"anilibrary-scraper/internal/domain/entity"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var _ Contract = (*AnimeVost)(nil)

type AnimeVost struct {
}

const amountToMakeFloat int = 10

func NewAnimeVost() AnimeVost {
	return AnimeVost{}
}

func (a AnimeVost) Title(document *goquery.Document) string {
	if title := document.Find(".shortstoryHead h1").First(); title != nil {
		raw := title.Text()
		return strings.TrimSpace(raw[0:strings.Index(raw, " /")])
	}

	return ""
}

func (a AnimeVost) Status(document *goquery.Document) entity.Status {
	if status := document.Find("#nexttime").First(); status != nil {
		return entity.Ongoing
	}

	return entity.Ready
}

func (a AnimeVost) Rating(document *goquery.Document) float32 {
	if rating := document.Find(".current-rating").First(); rating != nil {
		value, err := strconv.Atoi(rating.Text())

		if err != nil {
			return MinimalAnimeRating
		}

		return float32(value / amountToMakeFloat)
	}

	return MinimalAnimeRating
}

func (a AnimeVost) Episodes(document *goquery.Document) string {
	if episodes := document.Find("p strong:contains(Количество)").First(); episodes != nil {
		raw := episodes.Parent().Text()
		raw = strings.Replace(raw, "Количество серий: ", "", 1)
		raw = strings.Replace(raw, "+", "", 1)
		end := strings.Index(raw, " (")

		return raw[0:end]
	}

	return MinimalAnimeEpisodes
}

func (a AnimeVost) Genres(document *goquery.Document) []string {
	if genres := document.Find("p strong:contains(Жанр)").First(); genres != nil {
		raw := genres.Parent().Text()
		raw = strings.Replace(raw, "Жанр: ", "", 1)

		return strings.Split(cases.Title(language.Russian).String(raw), ", ")
	}

	return nil
}

func (a AnimeVost) VoiceActing(document *goquery.Document) []string {
	return []string{"AnimeVost"}
}

func (a AnimeVost) Image(document *goquery.Document) string {
	if attr, exists := document.Find(".imgRadius").First().Attr("src"); exists {
		return attr
	}

	return ""
}
