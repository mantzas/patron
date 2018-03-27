package route

import "net/http"

// Route definition
type Route struct {
	Pattern string
	Method  string
	Handler http.HandlerFunc
}

// New returns a new route
func New(p string, m string, h http.HandlerFunc) *Route {
	return &Route{p, m, h}
}
