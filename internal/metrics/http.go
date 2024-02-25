package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestErrorsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_errors_counter",
		Help: "Counter of http errors",
	})

	httpRequestSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_success_counter",
		Help: "Counter of http successful requests",
	})
)

func IncrHTTPErrorsCounter() {
	httpRequestErrorsCounter.Inc()
}

func IncrHTTPSuccessCounter() {
	httpRequestSuccessCounter.Inc()
}
