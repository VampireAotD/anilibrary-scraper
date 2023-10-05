package router

import (
	"anilibrary-scraper/internal/handler/http/api/v1/anime"
	"anilibrary-scraper/internal/handler/http/monitoring/healthcheck"
)

type Handlers struct {
	anime       anime.Controller
	healthcheck healthcheck.Controller
}

func NewHandlers(animeHandler anime.Controller, hcHandler healthcheck.Controller) Handlers {
	return Handlers{anime: animeHandler, healthcheck: hcHandler}
}
