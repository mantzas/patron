package http

import (
	"net/http"

	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/info"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(encoding.ContentTypeHeader, json.TypeCharset)

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
