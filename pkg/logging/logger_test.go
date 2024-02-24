package logging

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	const testMessage string = "test"

	var buff bytes.Buffer

	t.Run("Get logger", func(t *testing.T) {
		require.NotNil(t, Get())
	})

	t.Run("Console encoder", func(t *testing.T) {
		defer buff.Reset()

		logger := New(WithOutput(&buff))
		logger.Info(testMessage)

		log := buff.String()
		require.NotZero(t, log)
		require.Contains(t, log, testMessage)
	})

	t.Run("JSON encoder", func(t *testing.T) {
		defer buff.Reset()

		logger := New(WithOutput(&buff), ConvertToJSON())
		logger.Info(testMessage)

		testLog := struct {
			Message string `json:"message"`
		}{}

		require.NoError(t, json.NewDecoder(&buff).Decode(&testLog))
		require.NotZero(t, testLog.Message)
		require.Equal(t, testMessage, testLog.Message)
	})

	t.Run("ECS support", func(t *testing.T) {
		var buff bytes.Buffer
		defer buff.Reset()

		logger := New(WithOutput(&buff), ConvertToJSON(), ECSCompatible())
		logger.Info(testMessage)

		testLog := struct {
			Message    string `json:"message"`
			ECSVersion string `json:"ecs.version"`
		}{}

		require.NoError(t, json.NewDecoder(&buff).Decode(&testLog))
		require.Equal(t, testMessage, testLog.Message)
		require.NotZero(t, testLog.ECSVersion)
	})
}
