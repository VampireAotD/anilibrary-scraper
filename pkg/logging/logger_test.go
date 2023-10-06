package logging

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("Get logger", func(t *testing.T) {
		require.NotNil(t, Get())
	})

	t.Run("ECS support", func(t *testing.T) {
		var buff bytes.Buffer
		defer buff.Reset()

		logger := New(WithOutput(io.Discard), WithLogFiles(&buff), ECSCompatible())

		testLog := struct {
			Message    string `json:"message"`
			ECSVersion string `json:"ecs.version"`
		}{}

		const log string = "test"
		logger.Info(log)

		require.NoError(t, json.NewDecoder(&buff).Decode(&testLog))
		require.Equal(t, log, testLog.Message)
		require.NotZero(t, testLog.ECSVersion)
	})
}
