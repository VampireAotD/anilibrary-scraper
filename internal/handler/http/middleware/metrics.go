package middleware

import (
	"net/http"
	"time"

	"anilibrary-scraper/internal/metrics"

	"github.com/go-chi/chi/v5"
)

func ResponseMetrics(next http.Handler) http.Handler {
	metrics.IncrHTTPRequestsCounter()

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next.ServeHTTP(writer, request)

		duration := time.Since(start)

		metrics.RecordHTTPResponseTime(chi.RouteContext(request.Context()).RoutePattern(), request.Method, duration.Seconds())
	})
}
