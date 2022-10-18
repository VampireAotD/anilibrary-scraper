package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"anilibrary-scraper/internal/app/providers"
	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/internal/handler/http/router"
	"anilibrary-scraper/internal/handler/http/server"
	"anilibrary-scraper/pkg/logger"
	"github.com/go-redis/redis/v9"
)

type App struct {
	logger     logger.Contract
	config     config.Config
	connection *redis.Client
}

func New(logger logger.Contract, config config.Config, connection *redis.Client) *App {
	return &App{
		logger:     logger,
		config:     config,
		connection: connection,
	}
}

func Bootstrap() (*App, func()) {
	app, cleanup, err := WireApp()
	if err != nil {
		app.stopOnError("boostrap app", err)
	}

	err = providers.NewTimezoneProvider(app.config.App.Timezone, app.logger)
	if err != nil {
		app.stopOnError("timezone", err)
	}
	err = providers.NewJaegerTracerProvider(
		app.config.Jaeger.TraceEndpoint,
		app.config.App.Name,
		string(app.config.App.Env),
	)
	if err != nil {
		app.stopOnError("tracer", err)
	}

	return app, cleanup
}

func (app *App) stopOnError(info string, err error) {
	app.logger.Error(info, logger.Error(err))
	os.Exit(1)
}

func (app *App) Run() {
	address := fmt.Sprintf("%s:%d", app.config.HTTP.Addr, app.config.HTTP.Port)
	httpServer := server.NewHTTPServer(
		address,
		router.NewRouter(
			&router.Config{
				Url:             address,
				EnableProfiling: app.config.App.Env == config.Local,
				Logger:          app.logger.Named("api/http"),
				Handler:         WireAnimeController(app.connection),
			},
		),
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		os.Kill,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		defer stop()

		app.logger.Info("Starting server at", logger.String("addr", httpServer.Addr))

		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("while closing server", logger.Error(err))
		}
	}()

	<-ctx.Done()

	app.logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		app.logger.Error("error while shutting down server", logger.Error(err))
	}
}
