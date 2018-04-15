package httprouter

import (
	"net/http"
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/patron"
	"github.com/mantzas/patron/http/middleware"
	"github.com/mantzas/patron/log"
)

// CreateHandler creates a router
func CreateHandler(routes []patron.Route) http.Handler {
	routes = append(routes, profilingRoutes()...)
	log.Infof("adding %d routes", len(routes))

	router := httprouter.New()

	for _, route := range routes {
		h := middleware.DefaultMiddleware(route.Handler)
		router.HandlerFunc(route.Method, route.Pattern, h)
		log.Infof("added route %s %s", route.Method, route.Pattern)
	}
	return router
}

func profilingRoutes() []patron.Route {

	return []patron.Route{
		patron.NewRoute("/debug/pprof/", http.MethodGet, middleware.DefaultMiddleware(index)),
		patron.NewRoute("/debug/pprof/cmdline/", http.MethodGet, middleware.DefaultMiddleware(cmdline)),
		patron.NewRoute("/debug/pprof/profile/", http.MethodGet, middleware.DefaultMiddleware(profile)),
		patron.NewRoute("/debug/pprof/symbol/", http.MethodGet, middleware.DefaultMiddleware(symbol)),
		patron.NewRoute("/debug/pprof/trace/", http.MethodGet, middleware.DefaultMiddleware(trace)),
		patron.NewRoute("/debug/pprof/heap/", http.MethodGet, middleware.DefaultMiddleware(heap)),
		patron.NewRoute("/debug/pprof/goroutine/", http.MethodGet, middleware.DefaultMiddleware(goroutine)),
		patron.NewRoute("/debug/pprof/block/", http.MethodGet, middleware.DefaultMiddleware(block)),
		patron.NewRoute("/debug/pprof/threadcreate/", http.MethodGet, middleware.DefaultMiddleware(threadcreate)),
		patron.NewRoute("/debug/pprof/mutex/", http.MethodGet, middleware.DefaultMiddleware(mutex)),
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	pprof.Index(w, r)
}

func cmdline(w http.ResponseWriter, r *http.Request) {
	pprof.Cmdline(w, r)
}

func profile(w http.ResponseWriter, r *http.Request) {
	pprof.Profile(w, r)
}

func symbol(w http.ResponseWriter, r *http.Request) {
	pprof.Symbol(w, r)
}

func trace(w http.ResponseWriter, r *http.Request) {
	pprof.Trace(w, r)
}

func heap(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("heap").ServeHTTP(w, r)
}

func goroutine(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("goroutine").ServeHTTP(w, r)
}

func block(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("block").ServeHTTP(w, r)
}

func threadcreate(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("threadcreate").ServeHTTP(w, r)
}

func mutex(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("mutex").ServeHTTP(w, r)
}
