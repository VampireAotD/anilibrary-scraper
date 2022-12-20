package v1

import (
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/domain/repository/redis"
	"anilibrary-scraper/internal/domain/service"
	scraperService "anilibrary-scraper/internal/domain/service/scraper"
	"anilibrary-scraper/internal/handler/http/v1/anime"
	"anilibrary-scraper/internal/scraper"
	"github.com/google/wire"
)

var AnimeControllerProviderSet = wire.NewSet(
	redis.NewAnimeRepository,
	wire.Bind(new(repository.AnimeRepository), new(redis.AnimeRepository)),
	scraper.New,
	wire.Bind(new(scraper.Contract), new(scraper.Scraper)),
	scraperService.NewScraperService,
	wire.Bind(new(service.ScraperService), new(scraperService.Service)),
	anime.NewController,
)
