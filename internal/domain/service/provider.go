package service

import (
	"anilibrary-request-parser/internal/domain/repository/redis"
	"anilibrary-request-parser/internal/domain/service/anime"
	"github.com/google/wire"
)

var ScraperServiceProvider = wire.NewSet(
	redis.NewAnimeRepository,
	anime.NewScraperService,
)
