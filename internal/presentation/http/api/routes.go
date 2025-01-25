package api

import (
	_ "github.com/VampireAotD/anilibrary-scraper/docs" // generated swagger docs
	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/config"
	"github.com/VampireAotD/anilibrary-scraper/internal/presentation/http/api/v1/anime"
	"github.com/VampireAotD/anilibrary-scraper/internal/presentation/http/middleware"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Router          fiber.Router `name:"api-server"`
	Metrics         *fiberprometheus.FiberPrometheus
	AnimeController anime.Controller
	AppConfig       config.App
	JWTConfig       config.JWT
}

func RegisterAPIRoutes(params Params) {
	if !params.AppConfig.Env.Production() {
		params.Router.Get("/swagger/*", fiberSwagger.WrapHandler)
	}

	api := params.Router.Group("/api")
	v1 := api.Group("/v1")
	animeGroup := v1.Group("/anime")

	animeGroup.Use(middleware.NewJWTAuth(params.JWTConfig))
	animeGroup.Use(otelfiber.Middleware())
	animeGroup.Use(params.Metrics.Middleware)

	animeGroup.Post("/scrape", params.AnimeController.Scrape)
}
