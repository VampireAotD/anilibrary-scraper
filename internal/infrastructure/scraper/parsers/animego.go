package parsers

import (
	"strconv"
	"strings"

	"anilibrary-scraper/internal/infrastructure/scraper/model"

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
	if attr, exists := a.document.Find(".anime-poster img").Attr("src"); exists {
		return strings.Replace(attr, "/media/cache/thumbs_250x350", "", 1)
	}

	return ""
}

func (a AnimeGo) Parse() model.Anime {
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

func (a AnimeGo) title() string {
	if title := a.document.Find(".anime-title div h1").Text(); title != "" {
		return title
	}

	return ""
}

func (a AnimeGo) status() model.Status {
	if status := a.document.Find(".anime-info .row dt:contains(Статус) + dd").Text(); status != "" {
		return model.Status(status)
	}

	return model.Ready
}

func (a AnimeGo) rating() float32 {
	if rating := a.document.Find(".rating-value").Text(); rating != "" {
		value, err := strconv.ParseFloat(strings.Replace(rating, ",", ".", 1), 64)
		if err != nil {
			return MinimalAnimeRating
		}

		return float32(value)
	}

	return MinimalAnimeRating
}

func (a AnimeGo) episodes() int {
	if episodesText := a.document.Find(".anime-info .row dt:contains(Эпизоды) + dd").Text(); episodesText != "" {
		episodes, err := strconv.Atoi(episodesText)
		if err != nil {
			return MinimalAnimeEpisodes
		}

		return episodes
	}

	return MinimalAnimeEpisodes
}

func (a AnimeGo) genres() []string {
	if genresText := a.document.Find(".anime-info .row dt:contains(Жанр) + dd").Text(); genresText != "" {
		genres := strings.Split(genresText, ",")

		for i := range genres {
			genres[i] = strings.TrimSpace(genres[i])
		}

		return genres
	}

	return nil
}

func (a AnimeGo) voiceActing() []string {
	if voiceActingText := a.document.Find(".anime-info .row dt:contains(Озвучка) + dd").Text(); voiceActingText != "" {
		voiceActing := strings.Split(voiceActingText, ",")

		for i := range voiceActing {
			voiceActing[i] = strings.TrimSpace(voiceActing[i])
		}

		return voiceActing
	}

	return nil
}

func (a AnimeGo) synonyms() []string {
	if synonymsList := a.document.Find(".synonyms ul li"); synonymsList.Length() != 0 {
		synonyms := make([]string, 0, synonymsList.Length())

		synonymsList.Each(func(_ int, selection *goquery.Selection) {
			synonyms = append(synonyms, selection.Text())
		})

		return synonyms
	}

	return nil
}

func (a AnimeGo) year() int {
	if yearText := a.document.Find(".anime-info .row dt:contains(Выпуск) + dd").Text(); yearText != "" {
		year, err := strconv.Atoi(yearRegex.FindString(yearText))
		if err != nil {
			return 0
		}

		return year
	}

	return 0
}

func (a AnimeGo) animeType() model.Type {
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
