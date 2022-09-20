package middleware

import (
	"net/http"
	"time"

	"anilibrary-scraper/internal/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func ResponseMetrics(next http.Handler) http.Handler {
	prometheus.MustRegister(metrics.ResponseHistogram)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next.ServeHTTP(writer, request)

		duration := time.Since(start)
		route := routePattern(request)

		metrics.ResponseHistogram.WithLabelValues(route, request.Method).Observe(duration.Seconds())
	})
}

func routePattern(r *http.Request) string {
	ctx := chi.RouteContext(r.Context())

	if pattern := ctx.RoutePattern(); pattern != "" {
		return pattern
	}

	return ""
}
