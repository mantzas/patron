package http

import (
	"net/http"

	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/info"
	"github.com/prometheus/client_golang/prometheus"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(encoding.ContentTypeHeader, json.TypeCharset)

	mm, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, m := range mm {
		info.AddMetric(m.GetName(), m.GetHelp())
	}

	body, err := info.Marshal()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func infoRoute() Route {
	return NewRouteRaw("/info", http.MethodGet, infoHandler)
}
