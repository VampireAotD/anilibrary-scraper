package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_counter",
		Help: "Counter of http requests",
	})

	httpRequestErrorsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_errors_counter",
		Help: "Counter of http errors",
	})

	httpRequestSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_success_counter",
		Help: "Counter of http successful requests",
	})

	httpResponseTimeHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http_api",
		Name:      "request_duration",
		Help:      "Handler response time in seconds",
		Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"route", "method"})
)

func IncrHTTPRequestsCounter() {
	httpRequestsCounter.Inc()
}

func IncrHTTPErrorsCounter() {
	httpRequestErrorsCounter.Inc()
}

func IncrHTTPSuccessCounter() {
	httpRequestSuccessCounter.Inc()
}

func RecordHTTPResponseTime(route, method string, duration float64) {
	httpResponseTimeHistogram.WithLabelValues(route, method).Observe(duration)
}
