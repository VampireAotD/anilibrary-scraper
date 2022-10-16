package app

import (
	"os"
	"time"

	"anilibrary-scraper/internal/app/providers"
	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/pkg/closer"
	"anilibrary-scraper/pkg/logger"
	"github.com/go-redis/redis/v9"
)

type App struct {
	logger     logger.Contract
	config     config.Config
	closer     closer.Closers
	connection *redis.Client
}

func (app *App) stopOnError(info string, err error) {
	app.logger.Error(info, logger.Error(err))
	os.Exit(1)
}

func (app *App) ReadConfig() {
	cfg, err := config.New()
	if err != nil {
		app.stopOnError("error while reading config", err)
	}

	app.config = cfg
}

func (app *App) SetTimezone() {
	location, err := time.LoadLocation(app.config.App.Timezone)
	if err != nil {
		app.stopOnError("error while setting timezone", err)
	}

	time.Local = location
}

func (app *App) SetRedisConnection() {
	client, err := providers.NewRedisProvider(app.config.Redis)
	if err != nil {
		app.stopOnError("error while connecting to redis", err)
	}

	app.connection = client
	app.closer.Add("redis", app.connection.Close)
}

func (app *App) JaegerTracer() {
	err := providers.NewJaegerTracerProvider(
		app.config.Jaeger.TraceEndpoint,
		app.config.App.Name,
		string(app.config.App.Env),
	)
	if err != nil {
		app.stopOnError("jaeger tracing", err)
	}
}
