package scraper

import (
	"context"
	"fmt"

	"anilibrary-scraper/internal/entity"
)

type (
	Service interface {
		Process(ctx context.Context, url string) (*entity.Anime, error)
	}

	EventService interface {
		Send(ctx context.Context, url string) error
	}
)

type UseCase struct {
	scraperService Service
	eventService   EventService
}

func NewUseCase(scraperService Service, eventService EventService) UseCase {
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
