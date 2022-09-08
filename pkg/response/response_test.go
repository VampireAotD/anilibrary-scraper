package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	response := New(recorder)

	testCase := struct {
		Message string `json:"message"`
	}{
		Message: "test",
	}
	err := response.JSON(http.StatusOK, testCase)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, recorder.Header().Get("Content-Type"), "application/json")
}

func TestErrorJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	response := New(recorder)

	err := response.ErrorJSON(http.StatusUnprocessableEntity, errors.New("error"))

	require.NoError(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	require.JSONEq(t, `{"message": "error"}`, recorder.Body.String())
}
