package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cacheHitCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_counter",
		Help: "Counter of cache hits",
	})

	cacheMissCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_counter",
		Help: "Counter of cache misses",
	})
)

func IncrCacheHitCounter() {
	cacheHitCounter.Inc()
}

func IncrCacheMissCounter() {
	cacheMissCounter.Inc()
}
