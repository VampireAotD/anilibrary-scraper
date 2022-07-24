package app

import (
	"flag"
	"runtime"

	"anilibrary-request-parser/app/internal/config"
)

type flags struct {
	logPath, envPath string
	prom, pprof      bool
}

func Init() *App {
	var (
		logPath = flag.String("log", config.DefaultLoggerFileLocation, "Define log file path")
		envPath = flag.String("env", config.DefaultEnvLocation, "Define path to .env")
		prom    = flag.Bool("prom", false, "Enable Prometheus")
		pprof   = flag.Bool("pprof", false, "Enable pprof")
	)

	flag.Parse()

	defer runtime.GC()

	var app App

	app.flags = flags{
		logPath: *logPath,
		envPath: *envPath,
		prom:    *prom,
		pprof:   *pprof,
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
