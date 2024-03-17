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

const (
	delimiter int = 10

	animeVostMovie      string = "полнометражный фильм"
	animeVostShortMovie string = "короткометражный фильм"
)

type AnimeVost struct {
	document *goquery.Document
}

func NewAnimeVost(document *goquery.Document) AnimeVost {
	return AnimeVost{
		document: document,
	}
}

func (a AnimeVost) ImageURL() string {
	if attr, exists := a.document.Find(".imgRadius, .infoContent img").First().Attr("src"); exists {
		return AnimeVostURL + attr
	}

	return ""
}

func (a AnimeVost) Title() string {
	if title := a.document.Find(".shortstoryHead h1, .infoContent h3").First().Text(); title != "" {
		raw := strings.TrimSpace(title)
		return raw[0:strings.Index(raw, " /")]
	}

	return ""
}

func (a AnimeVost) Status() model.Status {
	if status := a.document.Find("#nexttime").Text(); status != "" {
		return model.Ongoing
	}

	return model.Ready
}

func (a AnimeVost) Rating() float32 {
	if ratingText := a.document.Find(".current-rating").Text(); ratingText != "" {
		rating, err := strconv.Atoi(ratingText)
		if err != nil {
			return MinimalAnimeRating
		}

		return float32(rating / delimiter)
	}

	return MinimalAnimeRating
}

func (a AnimeVost) Episodes() string {
	if episodesText := a.document.Find("p strong:contains(Количество)").Parent().Children().Remove().End().Text(); episodesText != "" {
		regex := regexp.MustCompile(`^\d+`)
		return regex.FindString(episodesText)
	}

	return MinimalAnimeEpisodes
}

func (a AnimeVost) Genres() []string {
	if genres := a.document.Find("p strong:contains(Жанр)").Parent().Text(); genres != "" {
		genres = strings.Replace(genres, "Жанр: ", "", 1)

		return strings.Split(cases.Title(language.Russian).String(genres), ", ")
	}

	return nil
}

func (a AnimeVost) VoiceActing() []string {
	return []string{"AnimeVost"}
}

func (a AnimeVost) Synonyms() []string {
	if title := a.document.Find(".shortstoryHead h1, .infoContent h3").First().Text(); title != "" {
		regex := regexp.MustCompile(`/\s*(.*?)\s*\[`)
		// If there is a synonym, then the correct one without any symbols will be at index 1
		if entries := regex.FindStringSubmatch(title); len(entries) > 1 {
			return []string{entries[1]}
		}

		return nil
	}

	return nil
}

func (a AnimeVost) Year() int {
	if yearText := a.document.Find("p strong:contains(Год)").Parent().Text(); yearText != "" {
		regex := regexp.MustCompile(`\d{4}`)
		year, err := strconv.Atoi(regex.FindString(yearText))
		if err != nil {
			return 0
		}

		return year
	}

	return 0
}

func (a AnimeVost) Type() model.Type {
	if typeText := a.document.Find("p strong:contains(Тип)").Parent().Text(); typeText != "" {
		animeType := strings.Replace(typeText, "Тип: ", "", 1)

		switch animeType {
		case animeVostMovie, animeVostShortMovie:
			return model.Movie
		default:
			return model.Show
		}
	}

	return model.Show
}
