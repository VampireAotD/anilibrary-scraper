package logging

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetLevel(t *testing.T) {
	logger := New()

	t.Run("Switch level", func(t *testing.T) {
		SetLevel("error")
		require.Equal(t, ErrorLevel, logger.Level().String())

		SetLevel("panic")
		require.Equal(t, PanicLevel, logger.Level().String())
	})

	t.Run("Unsupported level", func(t *testing.T) {
		SetLevel("random")
		require.Equal(t, DebugLevel, logger.Level().String())
	})
}
