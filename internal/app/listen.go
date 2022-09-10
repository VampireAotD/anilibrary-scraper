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

	"anilibrary-request-parser/pkg/logger"
)

func (app *App) Listen() {
	router, err := app.Router()

	if err != nil {
		app.logger.Error("error while creating router", logger.Error(err))
		app.closer.Close()

		os.Exit(1)
	}

	defer app.closer.Close()

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.config.HTTP.Addr, app.config.HTTP.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		app.logger.Info("Starting server at", logger.String("addr", server.Addr))

		err = server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("error from server", logger.Error(err))
		}

		stop()
	}()

	<-ctx.Done()

	app.logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		app.logger.Error("error while shutting down server", logger.Error(err))
	}
}
