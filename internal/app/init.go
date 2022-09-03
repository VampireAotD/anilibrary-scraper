package app

import (
	"flag"
	"runtime"
)

type flags struct {
	prom, pprof bool
}

func Init() *App {
	var (
		prom  = flag.Bool("prom", false, "Enable Prometheus")
		pprof = flag.Bool("pprof", false, "Enable pprof")
	)

	flag.Parse()

	defer runtime.GC()

	var app App

	app.flags = flags{
		prom:  *prom,
		pprof: *pprof,
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
