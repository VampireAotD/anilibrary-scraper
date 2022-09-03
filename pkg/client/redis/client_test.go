package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/require"
)

func TestRedisClient(t *testing.T) {
	mr, err := miniredis.Run()

	require.NoError(t, err)
	require.NotNil(t, mr)

	client, err := New(context.Background(), Config{
		Address: mr.Addr(),
	})

	require.NoError(t, err)
	require.NotNil(t, client)
}
