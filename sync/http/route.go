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
	Trace   bool
}

// NewRoute returns a new route from a generic handler
func NewRoute(p string, m string, pr sync.ProcessorFunc, trace bool) Route {
	return Route{Pattern: p, Method: m, Handler: handler(pr), Trace: trace}
}

// NewRouteRaw returns a new route from a HTTP handler
func NewRouteRaw(p string, m string, h http.HandlerFunc) Route {
	return Route{Pattern: p, Method: m, Handler: h, Trace: false}
}
