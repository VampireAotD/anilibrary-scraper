package main

import (
	"anilibrary-scraper/internal/app"
)

func main() {
	application := app.Bootstrap()

	application.Run()
}
