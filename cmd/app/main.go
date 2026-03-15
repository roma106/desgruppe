package main

import (
	"desgruppe/internal/app"
	"desgruppe/internal/logger"
)

func main() {
	app, err := app.New()
	if err != nil {
		logger.Error("Failed to build app: ", err)
		return
	}
	err = app.Run()
	if err != nil {
		logger.Error("Failed to start app: ", err)
	}
}
