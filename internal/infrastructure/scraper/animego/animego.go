package animego

import (
	"regexp"
	"strconv"
	"strings"

	animeEnum "anilibrary-request-parser/internal/domain/enum/anime"
	"anilibrary-request-parser/internal/infrastructure/scraper"
	"github.com/PuerkitoBio/goquery"
)

type AnimeGo struct {
	*scraper.Scraper // should this be embedded?
}

func New(Scraper *scraper.Scraper) *AnimeGo {
	return &AnimeGo{Scraper: Scraper}
}

func (a AnimeGo) Title(document *goquery.Document) string {
	if title := document.Find(".anime-title div h1").First(); title != nil {
		return title.Text()
	}

	return ""
}

func (a AnimeGo) Status(document *goquery.Document) animeEnum.Status {
	if status := document.Find(".anime-info .row dt:contains(Статус) + dd").First(); status != nil {
		return animeEnum.Status(status.Text())
	}

	return animeEnum.Ready
}

func (a AnimeGo) Rating(document *goquery.Document) float32 {
	if rating := document.Find(".rating-value").First(); rating != nil {
		value, err := strconv.ParseFloat(strings.Replace(rating.Text(), ",", ".", 1), 64)

		if err != nil {
			return scraper.MinimalAnimeRating
		}

		return float32(value)
	}

	return scraper.MinimalAnimeRating
}

func (a AnimeGo) Episodes(document *goquery.Document) string {
	if episodesText := document.Find(".anime-info .row dt:contains(Эпизоды) + dd").First(); episodesText != nil {
		return episodesText.Text()
	}

	return scraper.MinimalAnimeEpisodes
}

func (a AnimeGo) Genres(document *goquery.Document) []string {
	if genresText := document.Find(".anime-info .row dt:contains(Жанр) + dd").First(); genresText != nil {
		genres := strings.Split(genresText.Text(), ",")

		for i := range genres {
			genres[i] = strings.TrimSpace(genres[i])
		}

		return genres
	}

	return nil
}

func (a AnimeGo) VoiceActing(document *goquery.Document) []string {
	if voiceActingText := document.Find(".anime-info .row dt:contains(Озвучка) + dd").First(); voiceActingText != nil {
		regex := regexp.MustCompile(`,\s+`)

		return strings.Split(regex.ReplaceAllString(voiceActingText.Text(), ","), ",")
	}

	return nil
}
