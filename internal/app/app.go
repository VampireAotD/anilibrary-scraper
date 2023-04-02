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

// Bootstrap method creates new App, setting up all dependencies and initializes traces
func Bootstrap() (*App, func(), error) {
	cfg, err := config.New()
	if err != nil {
		return nil, nil, fmt.Errorf("config: %w", err)
	}

	logger, loggerCloser, err := providers.NewLoggerProvider()
	if err != nil {
		return nil, nil, fmt.Errorf("logger: %w", err)
	}

	redisConnection, redisCloser, err := providers.NewRedisProvider(cfg.Redis, logger)
	if err != nil {
		loggerCloser()
		return nil, nil, fmt.Errorf("redis: %w", err)
	}

	dependencies := SetupDependencies(logger, redisConnection)
	app := New(cfg, dependencies)

	jaegerCloser, err := providers.NewJaegerTracerProvider(
		app.config.Jaeger.TraceEndpoint,
		app.config.App.Name,
		string(app.config.App.Env),
		app.dependencies.logger,
	)
	if err != nil {
		redisCloser()
		loggerCloser()
		return nil, nil, fmt.Errorf("jaeger: %w", err)
	}

	return app, func() {
		jaegerCloser()
		redisCloser()
		loggerCloser()
	}, nil
}

func (app *App) Run() {
	address := fmt.Sprintf("%s:%d", app.config.HTTP.Addr, app.config.HTTP.Port)
	httpServer := server.NewHTTPServer(
		address,
		router.NewRouter(
			&router.Config{
				EnableProfiling: app.config.App.Env == config.Local,
				Logger:          app.dependencies.logger.Named("api/http"),
				RedisConnection: app.dependencies.redisConnection,
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
