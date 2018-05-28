package http

import (
	"net/http"

	"github.com/mantzas/patron/sync"
)

// Route definition
type Route struct {
	Pattern string
	Method  string
	Handler http.HandlerFunc
}

// NewRoute returns a new route from a generic handler
func NewRoute(p string, m string, h sync.Handler) Route {
	return Route{p, m, DefaultMiddleware(handler(h))}
}

// NewRouteRaw returns a new route from a HTTP handler
func NewRouteRaw(p string, m string, h http.HandlerFunc) Route {
	return Route{p, m, DefaultMiddleware(h)}
}
