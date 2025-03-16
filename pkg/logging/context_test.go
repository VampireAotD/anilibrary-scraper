package logging

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextWithLogger(t *testing.T) {
	var buf bytes.Buffer
	defer buf.Reset()

	logger := New(WithOutput(&buf))
	ctx := ContextWithLogger(t.Context(), logger.With(String("logger-test", "test")))

	loggerFromCtx := FromContext(ctx)
	require.NotNil(t, loggerFromCtx)

	loggerFromCtx.Info("test")
	require.Contains(t, buf.String(), "logger-test", "Logger must be set in context with fields")
}

func TestFromContext(t *testing.T) {
	t.Run("Logger wasn't set in context", func(t *testing.T) {
		logger := FromContext(t.Context())
		require.Equal(t, Get(), logger, "Default logger will be provided")
	})

	t.Run("Logger was set in context", func(t *testing.T) {
		logger := New(WithOutput(io.Discard))
		ctx := ContextWithLogger(t.Context(), logger.With(String("logger-test", "test")))

		loggerFromCtx := FromContext(ctx)
		require.NotNil(t, loggerFromCtx)
		require.NotEqual(t, Get(), loggerFromCtx, "Provided logger will be used")
	})
}
