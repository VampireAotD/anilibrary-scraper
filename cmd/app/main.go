package main

import (
	"anilibrary-request-parser/internal/app"
)

func main() {
	application := app.Init()

	application.Run()
}
