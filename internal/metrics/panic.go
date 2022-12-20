package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	panicCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "panics_counter",
			Help: "Counter of app panics",
		},
	)
)

func IncrPanicCounter() {
	panicCounter.Inc()
}
