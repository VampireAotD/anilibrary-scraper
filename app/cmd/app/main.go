package main

import (
	"anilibrary-request-parser/app/internal/app"
)

func main() {
	application := app.Init()

	application.Run()
}
