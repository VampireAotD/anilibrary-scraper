package api

import (
	_ "anilibrary-scraper/docs" // generated swagger docs
	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/internal/handler/http/api/v1/anime"
	"anilibrary-scraper/internal/handler/http/middleware"

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
	JWTConfig       config.JWT
}

func RegisterAPIRoutes(params Params) {
	params.Router.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := params.Router.Group("/api")
	v1 := api.Group("/v1")
	animeGroup := v1.Group("/anime")

	animeGroup.Use(middleware.NewJWTAuth(params.JWTConfig))
	animeGroup.Use(otelfiber.Middleware())
	animeGroup.Use(params.Metrics.Middleware)

	animeGroup.Post("/parse", params.AnimeController.Parse)
}