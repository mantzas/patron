package pprof

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Package blank import for supporting pprof
	"time"
)

// Server defines a HTTP server with pprof handler enabled
type Server struct {
	srv *http.Server
}

// New returns a new pprof HTTP server
func New(port int) *Server {

	s := Server{
		srv: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      http.DefaultServeMux,
		},
	}

	return &s
}

// SetPort sets the port of the server
func (s *Server) SetPort(port int) {
	s.srv.Addr = fmt.Sprintf(":%d", port)
}

// GetAddr gets the address of the service
func (s *Server) GetAddr() string {
	return s.srv.Addr
}

// ListenAndServe starts up the pprof server, listens and serves requests
func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

// Shutdown shuts down pprof server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
