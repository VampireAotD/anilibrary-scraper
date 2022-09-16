package main

import (
	"anilibrary-request-parser/internal/app"
)

func main() {
	application := app.Bootstrap()

	application.Run()
}
