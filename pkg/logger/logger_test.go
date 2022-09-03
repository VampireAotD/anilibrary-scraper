package logger

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateLogger(t *testing.T) {
	logger, err := New(Config{
		ConsoleOutput: io.Discard,
	})

	require.NoError(t, err)
	require.NotNil(t, logger)
}
