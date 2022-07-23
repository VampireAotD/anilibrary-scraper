package app

import (
	"flag"
	"runtime"

	"anilibrary-request-parser/app/internal/config"
)

type flags struct {
	logPath, envPath string
}

func Init() *App {
	var (
		logPath = flag.String("log", config.DefaultLoggerFileLocation, "Define log file path")
		envPath = flag.String("env", config.DefaultEnvLocation, "Define path to .env")
	)

	flag.Parse()

	defer runtime.GC()

	var app App

	app.flags = flags{
		logPath: *logPath,
		envPath: *envPath,
	}

	app.SetCloser()
	app.SetLogger()

	app.logger.Info("Initializing logger")

	app.ReadConfig()

	app.logger.Info("Reading config")

	app.logger.Info("Setting timezone")

	app.SetTimezone()

	return &app
}
