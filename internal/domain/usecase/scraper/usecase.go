package scraper

import (
	"context"
	"fmt"

	"anilibrary-scraper/internal/domain/entity"
	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/domain/usecase"
)

var _ usecase.ScraperUseCase = (*UseCase)(nil)

type UseCase struct {
	scraperService service.ScraperService
	eventService   service.EventService
}

func NewUseCase(scraperService service.ScraperService, eventService service.EventService) UseCase {
	return UseCase{
		scraperService: scraperService,
		eventService:   eventService,
	}
}

func (u UseCase) Scrape(ctx context.Context, url string) (*entity.Anime, error) {
	if err := u.eventService.Send(ctx, url); err != nil {
		return nil, fmt.Errorf("event service: %w", err)
	}

	return u.scraperService.Process(ctx, url)
}
