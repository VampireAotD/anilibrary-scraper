package main

import (
	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/internal/di"

	"go.uber.org/fx"
)

//	@title			Anilibrary Scraper
//	@version		1.0
//	@description	Microservice for scraping anime data
//	@termsOfService	http://swagger.io/terms/

//	@license.name	Boost Software License, Version 1.0
//	@license.url	https://github.com/VampireAotD/anilibrary-scraper/blob/main/LICENSE

//	@host		localhost:8080
//	@BasePath	/api/v1
func main() {
	fx.New(createApp()).Run()
}

func createApp() fx.Option {
	return fx.Options(
		fx.Provide(
			config.New,
		),
		di.ProviderModule,
		di.RepositoryModule,
		di.ServiceModule,
		di.UseCaseModule,
		di.HTTPModule,
	)
}
