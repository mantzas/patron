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
func NewRoute(p string, m string, pr sync.Processor) Route {
	return Route{p, m, handler(pr)}
}

// NewRouteRaw returns a new route from a HTTP handler
func NewRouteRaw(p string, m string, h http.HandlerFunc) Route {
	return Route{p, m, h}
}
