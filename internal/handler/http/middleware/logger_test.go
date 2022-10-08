package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"anilibrary-scraper/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	var buff bytes.Buffer
	log := logger.NewLogger(io.Discard, &buff)

	t.Run("WithLogger", func(t *testing.T) {
		ctx := WithLogger(context.Background(), log)
		require.NotNil(t, ctx)
	})

	t.Run("GetLogger", func(t *testing.T) {
		ctx := WithLogger(context.Background(), log)
		loggerFromCtx := GetLogger(ctx)
		require.NotNil(t, loggerFromCtx)
	})

	t.Run("MustGetLogger", func(t *testing.T) {
		t.Run("NoError", func(t *testing.T) {
			ctx := WithLogger(context.Background(), log)
			loggerFromCtx := MustGetLogger(ctx)
			require.NotNil(t, loggerFromCtx)
		})

		t.Run("WithError", func(t *testing.T) {
			defer func() {
				require.Equal(t, ErrNoLogger, recover())
			}()

			ctx := WithLogger(context.Background(), nil)
			loggerFromCtx := MustGetLogger(ctx)
			require.NotNil(t, loggerFromCtx)
		})
	})

	t.Run("WriteLogs", func(t *testing.T) {
		const testLog string = "test"
		router := chi.NewRouter()

		router.Use(Logger(log))
		router.Get("/test", func(writer http.ResponseWriter, request *http.Request) {
			loggerFromContext := MustGetLogger(request.Context())
			loggerFromContext.Info(testLog)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(httptest.NewRecorder(), req)

		expected := struct {
			Log string `json:"log"`
		}{}

		require.NoError(t, json.NewDecoder(&buff).Decode(&expected))
		require.Equal(t, testLog, expected.Log)
	})
}
