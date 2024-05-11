package parsers

import (
	"strconv"
	"strings"

	"anilibrary-scraper/internal/scraper/model"

	"github.com/PuerkitoBio/goquery"
)

const (
	animeGoMovie string = "Фильм"
)

type AnimeGo struct {
	document *goquery.Document
}

func NewAnimeGo(document *goquery.Document) AnimeGo {
	return AnimeGo{
		document: document,
	}
}

func (a AnimeGo) ImageURL() string {
	if attr, exists := a.document.Find(".anime-poster img").First().Attr("src"); exists {
		return strings.Replace(attr, "/media/cache/thumbs_250x350", "", 1)
	}

	return ""
}

func (a AnimeGo) Title() string {
	if title := a.document.Find(".anime-title div h1").Text(); title != "" {
		return title
	}

	return ""
}

func (a AnimeGo) Status() model.Status {
	if status := a.document.Find(".anime-info .row dt:contains(Статус) + dd").Text(); status != "" {
		return model.Status(status)
	}

	return model.Ready
}

func (a AnimeGo) Rating() float32 {
	if rating := a.document.Find(".rating-value").Text(); rating != "" {
		value, err := strconv.ParseFloat(strings.Replace(rating, ",", ".", 1), 64)
		if err != nil {
			return MinimalAnimeRating
		}

		return float32(value)
	}

	return MinimalAnimeRating
}

func (a AnimeGo) Episodes() string {
	if episodesText := a.document.Find(".anime-info .row dt:contains(Эпизоды) + dd").Text(); episodesText != "" {
		return episodesText
	}

	return MinimalAnimeEpisodes
}

func (a AnimeGo) Genres() []string {
	if genresText := a.document.Find(".anime-info .row dt:contains(Жанр) + dd").Text(); genresText != "" {
		genres := strings.Split(genresText, ",")

		for i := range genres {
			genres[i] = strings.TrimSpace(genres[i])
		}

		return genres
	}

	return nil
}

func (a AnimeGo) VoiceActing() []string {
	if voiceActingText := a.document.Find(".anime-info .row dt:contains(Озвучка) + dd").Text(); voiceActingText != "" {
		voiceActing := strings.Split(voiceActingText, ",")

		for i := range voiceActing {
			voiceActing[i] = strings.TrimSpace(voiceActing[i])
		}

		return voiceActing
	}

	return nil
}

func (a AnimeGo) Synonyms() []string {
	if synonymsList := a.document.Find(".synonyms ul li"); synonymsList.Length() != 0 {
		synonyms := make([]string, 0, synonymsList.Length())

		synonymsList.Each(func(_ int, selection *goquery.Selection) {
			synonyms = append(synonyms, selection.First().Text())
		})

		return synonyms
	}

	return nil
}

func (a AnimeGo) Year() int {
	if yearText := a.document.Find(".anime-info .row dt:contains(Сезон) + dd").Text(); yearText != "" {
		year, err := strconv.Atoi(yearRegex.FindString(yearText))
		if err != nil {
			return 0
		}

		return year
	}

	return 0
}

func (a AnimeGo) Type() model.Type {
	if typeText := a.document.Find(".anime-info .row dt:contains(Тип) + dd").Text(); typeText != "" {
		switch typeText {
		case animeGoMovie:
			return model.Movie
		default:
			return model.Show
		}
	}

	return model.Show
}
