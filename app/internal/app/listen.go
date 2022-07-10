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

	"go.uber.org/zap"
)

func (a *App) Listen() {
	router, err := a.Router()

	if err != nil {
		a.logger.Error("error while creating router", zap.Error(err))
		a.closer.Close()

		os.Exit(1)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", a.config.Http.Addr, a.config.Http.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Info("Starting server at", zap.String("addr", a.config.Http.Addr))

		err = server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("error from server", zap.Error(err))
		}

		stop()
	}()

	<-ctx.Done()

	a.logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		a.logger.Error("error while shutting down server", zap.Error(err))
	}
}
