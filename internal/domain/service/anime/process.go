package anime

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"anilibrary-scraper/internal/domain/dto"
	"anilibrary-scraper/internal/domain/entity"
	"github.com/PuerkitoBio/goquery"
)

func (s *ScraperService) Process(dto dto.ParseDTO) (*entity.Anime, error) {
	if dto.FromCache {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		anime, _ := s.repository.FindByUrl(ctx, dto.Url)
		if anime != nil {
			return anime, nil
		}
	}

	scraper, err := s.composeScraper(dto.Url)
	if err != nil {
		return nil, err
	}

	s.scraper = scraper

	response, err := s.client.Request(dto.Url)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sending request %v", err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("creating document %v", err)
	}

	anime := s.scrape(document)

	if dto.FromCache {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = s.repository.Create(ctx, dto.Url, anime)
	}

	return &anime, nil
}
