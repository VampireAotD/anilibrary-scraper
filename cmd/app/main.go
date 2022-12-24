package main

import (
	"anilibrary-scraper/internal/app"
)

//	@title			Anilibrary-scraper
//	@version		1.0
//	@description	Microservice for scraping anime data
//	@termsOfService	http://swagger.io/terms/

//	@license.name	Boost Software License, Version 1.0
//	@license.url	https://www.boost.org/LICENSE_1_0.txt

//	@host		localhost:8080
//	@BasePath	/api/v1
func main() {
	application, cleanup := app.Bootstrap()
	defer cleanup()

	application.Run()
}
