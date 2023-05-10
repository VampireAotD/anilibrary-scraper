package container

import (
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/domain/repository/kafka"
	"anilibrary-scraper/internal/domain/repository/redis"
	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/domain/service/event"
	scraperService "anilibrary-scraper/internal/domain/service/scraper"
	"anilibrary-scraper/internal/domain/usecase"
	scraperUseCase "anilibrary-scraper/internal/domain/usecase/scraper"
	"anilibrary-scraper/internal/handler/http/api/v1/anime"
	"anilibrary-scraper/internal/scraper"

	"github.com/google/wire"
)

// Scraper

var scraperBinding = wire.NewSet(
	scraper.New,
	wire.Bind(new(scraper.Contract), new(scraper.Scraper)),
)

// Repositories

var redisAnimeRepositoryBinding = wire.NewSet(
	redis.NewAnimeRepository,
	wire.Bind(new(repository.AnimeRepository), new(redis.AnimeRepository)),
)

var kafkaEventRepositoryBinding = wire.NewSet(
	kafka.NewEventRepository,
	wire.Bind(new(repository.EventRepository), new(kafka.EventRepository)),
)

// Services

var scraperServiceBinding = wire.NewSet(
	redisAnimeRepositoryBinding,
	scraperBinding,
	scraperService.NewScraperService,
	wire.Bind(new(service.ScraperService), new(scraperService.Service)),
)

var eventServiceBinding = wire.NewSet(
	kafkaEventRepositoryBinding,
	event.NewService,
	wire.Bind(new(service.EventService), new(event.Service)),
)

// UseCases

var scraperUseCaseBinding = wire.NewSet(
	scraperServiceBinding,
	eventServiceBinding,
	scraperUseCase.NewUseCase,
	wire.Bind(new(usecase.ScraperUseCase), new(scraperUseCase.UseCase)),
)

// Handlers

var HTTPAnimeHandlerSet = wire.NewSet(
	scraperUseCaseBinding,
	anime.NewController,
)
