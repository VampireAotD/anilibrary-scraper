package api

import (
	"anilibrary-scraper/internal/domain/service"
	"anilibrary-scraper/internal/domain/service/scraper"
	"anilibrary-scraper/internal/handler/http/api/anime"
	"github.com/google/wire"
)

var AnimeControllerProviderSet = wire.NewSet(
	service.ScraperProviderSet,
	wire.Bind(new(service.ScraperService), new(scraper.Service)),
	anime.NewController,
)
