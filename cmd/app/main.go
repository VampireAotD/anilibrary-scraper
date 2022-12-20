package main

import (
	"anilibrary-scraper/internal/app"
)

func main() {
	application, cleanup := app.Bootstrap()
	defer cleanup()

	application.Run()
}
