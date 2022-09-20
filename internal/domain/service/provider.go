package service

import (
	"anilibrary-scraper/internal/domain/repository/redis"
	"anilibrary-scraper/internal/domain/service/anime"
	"github.com/google/wire"
)

var ScraperProviderSet = wire.NewSet(redis.NewAnimeRepository, anime.NewScraperService)
