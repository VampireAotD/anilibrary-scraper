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

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type App struct {
	logger          logging.Contract
	redisConnection *redis.Client
	kafkaConnection *kafka.Conn
	config          config.Config
}

// New creates new App instance, setting up all dependencies and initializes traces, also returns cleanup function
// with all closers and an error if any occurs
func New() (*App, func(), error) {
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

	jaegerCloser, err := providers.NewJaegerTracerProvider(
		cfg.App.Name,
		string(cfg.App.Env),
		logger,
	)
	if err != nil {
		redisCloser()
		loggerCloser()
		return nil, nil, fmt.Errorf("jaeger: %w", err)
	}

	kafkaConnection, kafkaCleanup, err := providers.NewKafkaProvider(cfg.Kafka, logger)
	if err != nil {
		jaegerCloser()
		redisCloser()
		loggerCloser()
		return nil, nil, fmt.Errorf("kafka: %w", err)
	}

	cleanup := func() {
		kafkaCleanup()
		jaegerCloser()
		redisCloser()
		loggerCloser()
	}

	return &App{
		logger:          logger,
		redisConnection: redisConnection,
		kafkaConnection: kafkaConnection,
		config:          cfg,
	}, cleanup, nil
}

func (app *App) Run() {
	address := fmt.Sprintf("%s:%d", app.config.HTTP.Addr, app.config.HTTP.Port)
	httpServer := server.NewHTTPServer(
		address,
		router.NewRouter(
			router.WithProfilingRoutes(app.config.App.Env == config.Local),
			router.WithLogger(app.logger.Named("http")),
			router.WithRedisConnection(app.redisConnection),
			router.WithKafkaConnection(app.kafkaConnection),
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

		app.logger.Info("Starting server at", logging.String("addr", httpServer.Addr))

		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("while closing server", logging.Error(err))
		}
	}()

	<-ctx.Done()

	app.logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		app.logger.Error("error while shutting down server", logging.Error(err))
	}
}
