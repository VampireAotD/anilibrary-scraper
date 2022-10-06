package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func ProfilerRoutes() http.Handler {
	return middleware.Profiler()
}
