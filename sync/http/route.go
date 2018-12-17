package http

import (
	"net/http"

	"github.com/mantzas/patron/sync"
)

// Route definition of a HTTP route.
type Route struct {
	Pattern string
	Method  string
	Handler http.HandlerFunc
	Trace   bool
	Auth    Authenticator
}

// NewGetRoute creates a new GET route from a generic handler.
func NewGetRoute(p string, pr sync.ProcessorFunc, trace bool, auth Authenticator) Route {
	return NewRoute(p, http.MethodGet, pr, trace, auth)
}

// NewPostRoute creates a new POST route from a generic handler.
func NewPostRoute(p string, pr sync.ProcessorFunc, trace bool, auth Authenticator) Route {
	return NewRoute(p, http.MethodPost, pr, trace, auth)
}

// NewPutRoute creates a new PUT route from a generic handler.
func NewPutRoute(p string, pr sync.ProcessorFunc, trace bool, auth Authenticator) Route {
	return NewRoute(p, http.MethodPut, pr, trace, auth)
}

// NewDeleteRoute creates a new DELETE route from a generic handler.
func NewDeleteRoute(p string, pr sync.ProcessorFunc, trace bool, auth Authenticator) Route {
	return NewRoute(p, http.MethodDelete, pr, trace, auth)
}

// NewRoute creates a new route from a generic handler.
func NewRoute(p string, m string, pr sync.ProcessorFunc, trace bool, auth Authenticator) Route {
	return Route{Pattern: p, Method: m, Handler: handler(pr), Trace: trace, Auth: auth}
}

// NewRouteRaw creates a new route from a HTTP handler.
func NewRouteRaw(p string, m string, h http.HandlerFunc, trace bool, auth Authenticator) Route {
	return Route{Pattern: p, Method: m, Handler: h, Trace: trace, Auth: auth}
}
