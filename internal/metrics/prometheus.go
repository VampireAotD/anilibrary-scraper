package metrics

import "github.com/prometheus/client_golang/prometheus"

var ResponseHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "http_api",
	Name:      "parser_request_duration",
	Help:      "Handler response time in seconds",
	Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
}, []string{"route", "method"})
