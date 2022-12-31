package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func sendRequest(r *http.Request) *httptest.ResponseRecorder {
	testHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {})
	response := httptest.NewRecorder()

	testHandler(response, r)
	middleware := JWTAuth(testHandler)

	middleware.ServeHTTP(response, r)

	return response
}

func TestJWTAuthMiddleware(t *testing.T) {
	const testRequestURL string = "http://test/request"
	const testJWTSecret string = "test"

	t.Run("Error", func(t *testing.T) {
		t.Run("Without token in request", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, testRequestURL, nil)
			response := sendRequest(request)

			require.Equal(t, http.StatusUnauthorized, response.Code)
		})

		t.Run("Invalid signature", func(t *testing.T) {
			const invalidJWTSecret string = "random"

			request := httptest.NewRequest(http.MethodPatch, testRequestURL, nil)
			token, err := jwt.New(jwt.SigningMethodHS512).SignedString([]byte(invalidJWTSecret))
			require.NoError(t, err)

			request.Header.Set("Authorization", "Bearer "+token)
			response := sendRequest(request)

			require.Equal(t, http.StatusUnauthorized, response.Code)
		})

		t.Run("Invalid signing method", func(t *testing.T) {
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPatch, testRequestURL, nil)
			token, err := jwt.New(jwt.SigningMethodRS512).SignedString(privateKey)
			require.NoError(t, err)

			request.Header.Set("Authorization", "Bearer "+token)
			response := sendRequest(request)

			require.Equal(t, http.StatusUnauthorized, response.Code)
		})
	})

	t.Run("Valid token", func(t *testing.T) {
		require.NoError(t, os.Setenv("JWT_SECRET", testJWTSecret))

		request := httptest.NewRequest(http.MethodPut, testRequestURL, nil)
		token, err := jwt.New(jwt.SigningMethodHS512).SignedString([]byte(testJWTSecret))
		require.NoError(t, err)

		request.Header.Set("Authorization", "Bearer "+token)
		response := sendRequest(request)

		require.Equal(t, http.StatusOK, response.Code)
	})
}
