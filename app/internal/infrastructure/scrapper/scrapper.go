package scrapper

import (
	"fmt"
	"net/http"
	"strings"

	"anilibrary-request-parser/app/internal/domain/contract"
	"anilibrary-request-parser/app/internal/domain/entity"
	"anilibrary-request-parser/app/internal/infrastructure/scrapper/animego"
	"anilibrary-request-parser/app/internal/infrastructure/scrapper/animevost"
	"github.com/PuerkitoBio/goquery"
)

type Scrapper struct {
	url      string
	instance contract.Parser
}

func New(url string) (*Scrapper, error) {
	switch true {
	case strings.Contains(url, "animego.org"):
		instance := animego.New()
		return &Scrapper{url, instance}, nil
	case strings.Contains(url, "animevost.org"):
		instance := animevost.New()
		return &Scrapper{url, instance}, nil
	default:
		return nil, nil
	}
}

func (s Scrapper) Process() (*entity.Anime, error) {
	response, err := http.Get(s.url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code %d", response.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return nil, err
	}

	fmt.Println(document.Html())

	return nil, err
}
