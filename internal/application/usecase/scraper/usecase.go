package scraper

import (
	"context"
	"fmt"
	"time"

	"anilibrary-scraper/internal/application/service/event"
	"anilibrary-scraper/internal/domain/entity"
)

type (
	Service interface {
		Process(ctx context.Context, url string) (entity.Anime, error)
	}

	EventService interface {
		Send(ctx context.Context, dto event.DTO) error
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

func (u UseCase) Scrape(ctx context.Context, dto DTO) (entity.Anime, error) {
	if err := u.eventService.Send(ctx, event.DTO{
		URL:       dto.URL,
		Time:      time.Now(),
		IP:        dto.IP,
		UserAgent: dto.UserAgent,
	}); err != nil {
		return entity.Anime{}, fmt.Errorf("event service: %w", err)
	}

	return u.scraperService.Process(ctx, dto.URL)
}
