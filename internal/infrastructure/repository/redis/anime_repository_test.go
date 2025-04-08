package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTTL(t *testing.T) {
	t.Run("Has valid TTL", func(t *testing.T) {
		_, err := time.ParseDuration(sevenDaysInHours)
		require.NoError(t, err)
	})
}
