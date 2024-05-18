package main

import (
	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/internal/di"

	"go.uber.org/fx"
)

//	@title			Anilibrary Scraper
//	@version		1.0
//	@description	Microservice for scraping anime data
//	@termsOfService	https://swagger.io/terms/

//	@license.name	Boost Software License, Version 1.0
//	@license.url	https://github.com/VampireAotD/anilibrary-scraper/blob/main/LICENSE

//	@host						localhost:8080
//	@BasePath					/api/v1
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
func main() {
	fx.New(createApp()).Run()
}

func createApp() fx.Option {
	return fx.Options(
		// Resolve config
		fx.Provide(
			config.New,
		),

		// Provide dependencies
		di.ProviderModule,

		// Resolve services
		di.ServiceModule,

		// Resolve use cases
		di.UseCaseModule,

		// Start API and Monitoring servers
		di.HTTPModule,
	)
}
