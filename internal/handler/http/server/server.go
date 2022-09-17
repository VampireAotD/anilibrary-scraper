package server

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 5 * time.Second
	defaultIdleTimeout  = 15 * time.Second
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(address string, router http.Handler) *HTTPServer {
	return &HTTPServer{server: composeServer(address, router)}
}

func composeServer(address string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}
}

func (s HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s HTTPServer) Address() string {
	return s.server.Addr
}

func (s HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
