package http

import (
	"fmt"
	"net/http"
	"time"
)

// CreateHTTPServer returns a new HTTP server on a specific port
func CreateHTTPServer(port int, sm *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      sm,
	}
}
