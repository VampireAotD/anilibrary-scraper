package providers

import (
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/domain/repository/redis"
	"anilibrary-scraper/internal/domain/service"
	scraperService "anilibrary-scraper/internal/domain/service/scraper"
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

// Services

var scraperServiceBinding = wire.NewSet(
	scraperService.NewScraperService,
	wire.Bind(new(service.ScraperService), new(scraperService.Service)),
)
