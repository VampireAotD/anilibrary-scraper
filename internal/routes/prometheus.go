package routes

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PrometheusRoutes() http.Handler {
	return promhttp.Handler()
}
