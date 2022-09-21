package logger

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTypes(t *testing.T) {
	const key string = "key"

	t.Run("Bool", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, Bool(key, true), zap.Bool(key, true))
	})

	t.Run("Int64", func(t *testing.T) {
		t.Parallel()

		value := rand.Int63()
		require.Equal(t, Int64(key, value), zap.Int64(key, value))
	})

	t.Run("Int", func(t *testing.T) {
		t.Parallel()

		value := rand.Int()
		require.Equal(t, Int(key, value), zap.Int(key, value))
	})

	t.Run("Any", func(t *testing.T) {
		t.Parallel()

		value := any("rand")
		require.Equal(t, Any(key, value), zap.Any(key, value))
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		value := errors.New("rand")
		require.Equal(t, Error(value), zap.Error(value))
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, String(key, key), zap.String(key, key))
	})
}
