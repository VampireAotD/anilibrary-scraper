package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTracerMiddleware(t *testing.T) {
	t.Run("WithTracer", func(t *testing.T) {
		withTracer := WithTracer(context.Background())
		require.NotNil(t, withTracer)
	})

	t.Run("MushGetTracer", func(t *testing.T) {
		t.Run("With panic", func(t *testing.T) {
			defer func() {
				require.Equal(t, ErrNoTracer, recover())
			}()
			tracer := MustGetTracer(context.Background())
			require.Nil(t, tracer)
		})

		t.Run("Without panic", func(t *testing.T) {
			withTracer := WithTracer(context.Background())
			tracer := MustGetTracer(withTracer)
			require.NotNil(t, tracer)
		})
	})
}
