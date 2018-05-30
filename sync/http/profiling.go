package http

import (
	"net/http"
	"net/http/pprof"
)

func profilingRoutes() []Route {
	return []Route{
		NewRouteRaw("/debug/pprof/", http.MethodGet, index),
		NewRouteRaw("/debug/pprof/cmdline/", http.MethodGet, cmdline),
		NewRouteRaw("/debug/pprof/profile/", http.MethodGet, profile),
		NewRouteRaw("/debug/pprof/symbol/", http.MethodGet, symbol),
		NewRouteRaw("/debug/pprof/trace/", http.MethodGet, trace),
		NewRouteRaw("/debug/pprof/heap/", http.MethodGet, heap),
		NewRouteRaw("/debug/pprof/goroutine/", http.MethodGet, goroutine),
		NewRouteRaw("/debug/pprof/block/", http.MethodGet, block),
		NewRouteRaw("/debug/pprof/threadcreate/", http.MethodGet, threadcreate),
		NewRouteRaw("/debug/pprof/mutex/", http.MethodGet, mutex),
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
