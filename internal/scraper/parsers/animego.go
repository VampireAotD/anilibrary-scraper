package parsers

import (
	"regexp"
	"strconv"
	"strings"

	"anilibrary-scraper/internal/scraper/parsers/model"
	"github.com/PuerkitoBio/goquery"
)

var _ Contract = (*AnimeGo)(nil)

type AnimeGo struct{}

func NewAnimeGo() AnimeGo {
	return AnimeGo{}
}

func (a AnimeGo) Title(document *goquery.Document) string {
	if title := document.Find(".anime-title div h1").Text(); title != "" {
		return title
	}

	return ""
}

func (a AnimeGo) Status(document *goquery.Document) model.Status {
	if status := document.Find(".anime-info .row dt:contains(Статус) + dd").Text(); status != "" {
		return model.Status(status)
	}

	return model.Ready
}

func (a AnimeGo) Rating(document *goquery.Document) float32 {
	if rating := document.Find(".rating-value").Text(); rating != "" {
		value, err := strconv.ParseFloat(strings.Replace(rating, ",", ".", 1), 64)
		if err != nil {
			return MinimalAnimeRating
		}

		return float32(value)
	}

	return MinimalAnimeRating
}

func (a AnimeGo) Episodes(document *goquery.Document) string {
	if episodesText := document.Find(".anime-info .row dt:contains(Эпизоды) + dd").Text(); episodesText != "" {
		return episodesText
	}

	return MinimalAnimeEpisodes
}

func (a AnimeGo) Genres(document *goquery.Document) []string {
	if genresText := document.Find(".anime-info .row dt:contains(Жанр) + dd").Text(); genresText != "" {
		genres := strings.Split(genresText, ",")

		for i := range genres {
			genres[i] = strings.TrimSpace(genres[i])
		}

		return genres
	}

	return nil
}

func (a AnimeGo) VoiceActing(document *goquery.Document) []string {
	if voiceActingText := document.Find(".anime-info .row dt:contains(Озвучка) + dd").Text(); voiceActingText != "" {
		regex := regexp.MustCompile(`,\s+`)
		return strings.Split(regex.ReplaceAllString(voiceActingText, ","), ",")
	}

	return nil
}

func (a AnimeGo) Synonyms(document *goquery.Document) []string {
	if synonymsList := document.Find(".synonyms ul li"); synonymsList.Length() != 0 {
		synonyms := make([]string, 0, synonymsList.Length())

		synonymsList.Each(func(i int, selection *goquery.Selection) {
			synonyms = append(synonyms, selection.First().Text())
		})

		return synonyms
	}

	return nil
}

func (a AnimeGo) Image(document *goquery.Document) string {
	if attr, exists := document.Find(".anime-poster img").First().Attr("src"); exists {
		return strings.Replace(attr, "/media/cache/thumbs_250x350", "", 1)
	}

	return ""
}
