package parsers

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/scraper/model"

	"github.com/PuerkitoBio/goquery"
)

const (
	delimiter int = 10

	animeVostMovie      string = "полнометражный фильм"
	animeVostShortMovie string = "короткометражный фильм"
)

var (
	animeVostEpisodesRegex = regexp.MustCompile(`\d+`)
	animeVostSynonymsRegex = regexp.MustCompile(`/\s*(.*?)\s*\[`)
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
	if attr, exists := a.document.Find(".imgRadius, .infoContent img").Attr("src"); exists {
		return AnimeVostURL + attr
	}

	return ""
}

func (a AnimeVost) Parse() model.Anime {
	return model.Anime{
		Title:       a.title(),
		Status:      a.status(),
		Type:        a.animeType(),
		Episodes:    a.episodes(),
		Genres:      a.genres(),
		VoiceActing: a.voiceActing(),
		Synonyms:    a.synonyms(),
		Rating:      a.rating(),
		Year:        a.year(),
	}
}

func (a AnimeVost) title() string {
	if title := a.document.Find(".shortstoryHead h1, .infoTitle h1").Text(); title != "" {
		title = strings.TrimSpace(title)
		slashIndex := strings.Index(title, " /")
		if slashIndex < 0 {
			return ""
		}

		return title[0:slashIndex]
	}

	return ""
}

func (a AnimeVost) status() model.Status {
	if a.document.Find("#nexttime").Length() > 0 {
		return model.Ongoing
	}

	return model.Ready
}

func (a AnimeVost) rating() float32 {
	if ratingText := a.document.Find(".current-rating").Text(); ratingText != "" {
		rating, err := strconv.Atoi(ratingText)
		if err != nil {
			return MinimalAnimeRating
		}

		return float32(rating / delimiter)
	}

	return MinimalAnimeRating
}

func (a AnimeVost) episodes() int {
	if episodesText := a.document.Find("p:contains(Количество)").Text(); episodesText != "" {
		episodes, err := strconv.Atoi(animeVostEpisodesRegex.FindString(episodesText))
		if err != nil {
			return MinimalAnimeEpisodes
		}

		return episodes
	}

	return MinimalAnimeEpisodes
}

func (a AnimeVost) genres() []string {
	if genresText := a.document.Find("p:contains(Жанр)").Text(); genresText != "" {
		raw := strings.TrimPrefix(genresText, "Жанр: ")

		genres := strings.Split(raw, ", ")

		for i := range genres {
			// Genres on AnimeVost are always in lowercase, need to make first rune in uppercase
			letter, size := utf8.DecodeRuneInString(genres[i])
			letterToUpper := unicode.ToUpper(letter)
			genres[i] = string(letterToUpper) + genres[i][size:]
		}

		return genres
	}

	return nil
}

func (a AnimeVost) voiceActing() []string {
	return []string{"AnimeVost"}
}

func (a AnimeVost) synonyms() []string {
	if title := a.document.Find(".shortstoryHead h1, .infoTitle h1").Text(); title != "" {
		// If there is a synonym, then the correct one without any symbols will be at index 1
		if entries := animeVostSynonymsRegex.FindStringSubmatch(title); len(entries) > 1 {
			return []string{entries[1]}
		}

		return nil
	}

	return nil
}

func (a AnimeVost) year() int {
	if yearText := a.document.Find("p:contains(Год)").Text(); yearText != "" {
		year, err := strconv.Atoi(yearRegex.FindString(yearText))
		if err != nil {
			return 0
		}

		return year
	}

	return 0
}

func (a AnimeVost) animeType() model.Type {
	if typeText := a.document.Find("p:contains(Тип)").Text(); typeText != "" {
		animeType := strings.TrimPrefix(typeText, "Тип: ")

		switch animeType {
		case animeVostMovie, animeVostShortMovie:
			return model.Movie
		default:
			return model.Show
		}
	}

	return model.Show
}
