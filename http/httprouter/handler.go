package httprouter

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/patron/http/middleware"
	"github.com/mantzas/patron/http/route"
	"github.com/mantzas/patron/log"
)

// CreateHandler creates a router
func CreateHandler(routes []route.Route) (http.Handler, error) {
	if len(routes) == 0 {
		return nil, errors.New("no routes defined")
	}
	log.Infof("adding %d routes", len(routes))

	router := httprouter.New()

	for _, route := range routes {
		log.Infof("adding route %s %s", route.Method, route.Pattern)
		h := middleware.DefaultMiddleware(route.Handler)
		router.HandlerFunc(route.Method, route.Pattern, h)
	}

	return router, nil
}
