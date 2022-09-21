package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateLogger(t *testing.T) {
	logger := NewLogger(io.Discard)
	require.NotNil(t, logger)
}

func TestMethods(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(io.Discard, &buf)

	testCases := []struct {
		Name string
		Func func(msg string, fields ...Field)
	}{
		{
			Name: "Debug",
			Func: logger.Debug,
		},
		{
			Name: "Info",
			Func: logger.Info,
		},
		{
			Name: "Warn",
			Func: logger.Warn,
		},
		{
			Name: "Error",
			Func: logger.Error,
		},
	}

	for _, testCase := range testCases {
		buf.Reset()

		t.Run(testCase.Name, func(t *testing.T) {
			defer buf.Reset()

			testCase.Func("test")
			require.NoError(t, logger.Sync())
			require.NoError(t, json.NewDecoder(&buf).Decode(&struct{}{}))
		})
	}
}
