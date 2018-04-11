package http

import (
	_ "net/http/pprof" // Setting up default server mux pprof

	"net/http"

	"github.com/mantzas/patron/http/middleware"
	"github.com/mantzas/patron/log"
)

// HandlerGen type for implementing handler generator
type HandlerGen func([]Route) http.Handler

// CreateHandler creates a new handler
func CreateHandler(routes []Route) http.Handler {
	log.Infof("adding %d routes with pprof enabled", len(routes))

	for _, route := range routes {
		h := middleware.DefaultMiddleware(route.Handler)
		http.DefaultServeMux.Handle(route.Pattern, h)
		log.Infof("added route %s %s", route.Method, route.Pattern)
	}
	return http.DefaultServeMux
}
