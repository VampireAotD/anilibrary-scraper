package animego

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"anilibrary-request-parser/app/internal/domain/entity"
	animeEnum "anilibrary-request-parser/app/internal/domain/enum/anime"
	"anilibrary-request-parser/app/internal/infrastructure/scraper"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type AnimeGo struct {
	*scraper.Scrapper // should this be embedded?
}

func New(scrapper *scraper.Scrapper) *AnimeGo {
	return &AnimeGo{Scrapper: scrapper}
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
			return 0.0
		}

		return float32(value)
	}

	return 0.0
}

func (a AnimeGo) Episodes(document *goquery.Document) string {
	if episodesText := document.Find(".anime-info .row dt:contains(Эпизоды) + dd").First(); episodesText != nil {
		return episodesText.Text()
	}

	return `0 / ?`
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

func (a AnimeGo) Process() (*entity.Anime, error) {
	a.Logger.Info("Scraping", zap.String("url", a.Url))

	response, err := a.Client.Request(a.Url)

	if err != nil {
		return nil, fmt.Errorf("sending request %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code %d", response.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return nil, fmt.Errorf("creating document %v", err)
	}

	var anime entity.Anime

	anime.Title = a.Title(document)
	anime.Status = a.Status(document)
	anime.Rating = a.Rating(document)
	anime.Episodes = a.Episodes(document)
	anime.Genres = a.Genres(document)
	anime.VoiceActing = a.VoiceActing(document)

	return &anime, err
}
