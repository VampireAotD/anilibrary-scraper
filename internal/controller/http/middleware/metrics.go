package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func ResponseMetrics(next http.Handler) http.Handler {
	buckets := []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

	responseHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http_api",
		Name:      "parser_request_duration",
		Help:      "Handler response time in seconds",
		Buckets:   buckets,
	}, []string{"route", "method"})

	prometheus.MustRegister(responseHistogram)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next.ServeHTTP(writer, request)

		duration := time.Since(start)
		route := routePattern(request)

		responseHistogram.WithLabelValues(route, request.Method).Observe(duration.Seconds())
	})
}

func routePattern(r *http.Request) string {
	ctx := chi.RouteContext(r.Context())

	if pattern := ctx.RoutePattern(); pattern != "" {
		return pattern
	}

	return ""
}
