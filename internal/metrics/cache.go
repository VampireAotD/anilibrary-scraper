package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cacheHitCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_counter",
			Help: "Counter of cache hits",
		},
	)
)

func IncrCacheHitCounter() {
	cacheHitCounter.Inc()
}
