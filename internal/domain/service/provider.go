package service

import (
	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/domain/repository/redis"
	"anilibrary-scraper/internal/domain/service/scraper"
	"github.com/google/wire"
)

var ScraperProviderSet = wire.NewSet(
	redis.NewAnimeRepository,
	wire.Bind(new(repository.AnimeRepository), new(redis.AnimeRepository)),
	scraper.NewScraperService,
)
