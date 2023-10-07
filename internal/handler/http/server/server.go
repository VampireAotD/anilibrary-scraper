package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"anilibrary-scraper/config"
	"anilibrary-scraper/internal/handler/http/router"
	"anilibrary-scraper/pkg/logging"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 15 * time.Second
	defaultIdleTimeout  = 15 * time.Second
)

type Server struct {
	router *router.Router
	cfg    config.HTTP
}

func New(router *router.Router, cfg config.HTTP) Server {
	return Server{
		router: router,
		cfg:    cfg,
	}
}

func (s Server) Start(lifecycle fx.Lifecycle) {
	address := net.JoinHostPort(s.cfg.Addr, strconv.Itoa(s.cfg.Port))

	server := &http.Server{
		Addr:         address,
		Handler:      s.router.WithMetrics().Routes(),
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logging.Get().Info("HTTP server started at", zap.String("addr", address))

			go func() {
				err := server.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logging.Get().Error("HTTP server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
