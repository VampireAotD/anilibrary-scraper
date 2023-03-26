package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"anilibrary-scraper/internal/config"
	"anilibrary-scraper/internal/handler/http/router"
	"anilibrary-scraper/internal/handler/http/server"
	"anilibrary-scraper/internal/providers"
	"anilibrary-scraper/pkg/logging"
)

type App struct {
	dependencies Dependencies
	config       config.Config
}

func New(config config.Config, dependencies Dependencies) *App {
	return &App{
		dependencies: dependencies,
		config:       config,
	}
}

func (app *App) abort(info string, err error) {
	log.Fatalln(info, err)
}

func Bootstrap() (*App, func()) {
	app, closers, err := WireApp()
	if err != nil {
		app.abort("boostrap app", err)
	}

	jaegerCloser, err := providers.NewJaegerTracerProvider(
		app.config.Jaeger.TraceEndpoint,
		app.config.App.Name,
		string(app.config.App.Env),
		app.dependencies.logger,
	)
	if err != nil {
		app.abort("tracer", err)
	}

	cleanup := func() {
		jaegerCloser()
		closers()
	}

	return app, cleanup
}

func (app *App) Run() {
	address := fmt.Sprintf("%s:%d", app.config.HTTP.Addr, app.config.HTTP.Port)
	httpServer := server.NewHTTPServer(
		address,
		router.NewRouter(
			&router.Config{
				URL:             address,
				EnableProfiling: app.config.App.Env == config.Local,
				Logger:          app.dependencies.logger.Named("api/http"),
				Handler:         WireAnimeController(app.dependencies.redisConnection),
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

		app.dependencies.logger.Info("Starting server at", logging.String("addr", httpServer.Addr))

		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.dependencies.logger.Error("while closing server", logging.Error(err))
		}
	}()

	<-ctx.Done()

	app.dependencies.logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		app.dependencies.logger.Error("error while shutting down server", logging.Error(err))
	}
}
