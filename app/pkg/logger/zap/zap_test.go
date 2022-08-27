package zap

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateLogger(t *testing.T) {
	logger := NewLogger(io.Discard)
	require.NotNil(t, logger)
}
