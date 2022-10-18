package api

import (
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/domain/repository/redis"
	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/domain/service/scraper"
	"anilibrary-scraper/internal/handler/http/api/anime"
	"github.com/google/wire"
)

var AnimeControllerProviderSet = wire.NewSet(
	redis.NewAnimeRepository,
	wire.Bind(new(repository.AnimeRepository), new(redis.AnimeRepository)),
	scraper.NewScraperService,
	wire.Bind(new(service.ScraperService), new(scraper.Service)),
	anime.NewController,
)
