package logging

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLevel(t *testing.T) {
	require.Equal(t, DebugLevel, zapcore.DebugLevel)
	require.Equal(t, InfoLevel, zapcore.InfoLevel)
	require.Equal(t, WarnLevel, zapcore.WarnLevel)
	require.Equal(t, ErrorLevel, zapcore.ErrorLevel)
}

func TestFields(t *testing.T) {
	require.Equal(t, zap.String("test", "test"), String("test", "test"))
	require.Equal(t, zap.Error(errors.New("test")), Error(errors.New("test")))
	require.Equal(t, zap.Any("test", "test"), Any("test", "test"))
	require.Equal(t, zap.Int("test", 1), Int("test", 1))
	require.Equal(t, zap.Bool("test", true), Bool("test", true))
	require.Equal(t, zap.Float64("test", 1.0), Float64("test", 1.0))
}
