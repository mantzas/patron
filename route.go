package patron

import "net/http"

// Route definition
type Route struct {
	Pattern string
	Method  string
	Handler http.HandlerFunc
}

// NewRoute returns a new route
func NewRoute(p string, m string, h http.HandlerFunc) Route {
	return Route{p, m, h}
}
