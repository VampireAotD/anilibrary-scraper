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

	"anilibrary-request-parser/internal/handler/http/server"
	"go.uber.org/zap"
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

		app.logger.Info("Starting server at", zap.String("addr", httpServer.Addr))

		err := httpServer.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("while closing server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	app.logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		app.logger.Error("error while shutting down server", zap.Error(err))
	}
}
