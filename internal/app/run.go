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

	"anilibrary-scraper/internal/handler/http/server"
	"anilibrary-scraper/pkg/logger"
)

func (app *App) Run() {
	defer app.closer.Close(app.logger)

	httpServer := server.NewHTTPServer(
		fmt.Sprintf("%s:%d", app.config.HTTP.Addr, app.config.HTTP.Port),
		app.Router(),
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
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
