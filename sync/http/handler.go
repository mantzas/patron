package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/sync"
	"github.com/pkg/errors"
)

func handler(hnd sync.ProcessorFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ct, dec, enc, err := determineEncoding(r.Header)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}
		prepareResponse(w, ct)

		f := extractFields(r)
		for k, v := range extractParams(r) {
			f[k] = v
		}

		req := sync.NewRequest(f, r.Body, dec)
		rsp, err := hnd(r.Context(), req)
		if err != nil {
			handleError(w, err)
			return
		}

		err = handleSuccess(w, r, rsp, enc)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func extractFields(r *http.Request) map[string]string {

	f := make(map[string]string)

	for name, values := range r.URL.Query() {
		f[name] = values[0]
	}
	return f
}

func determineEncoding(hdr http.Header) (string, encoding.Decode, encoding.Encode, error) {

	c, err := determineContentType(hdr)
	if err != nil {
		return "", nil, nil, err
	}

	switch c {
	case json.ContentType, json.ContentTypeCharset:
		return c, json.Decode, json.Encode, nil
	}
	return "", nil, nil, errors.Errorf("accept header %s is unsupported", c)
}

func determineContentType(hdr http.Header) (string, error) {
	h, ok := hdr[encoding.ContentTypeHeader]
	if !ok {
		return "", errors.New("accept and content type header is missing")

	}
	return h[0], nil
}

func handleSuccess(w http.ResponseWriter, r *http.Request, rsp *sync.Response, enc encoding.Encode) error {

	if rsp == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	p, err := enc(rsp.Payload)
	if err != nil {
		return err
	}

	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusCreated)
	}

	_, err = w.Write(p)
	return err
}

func handleError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *sync.ValidationError:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	case *sync.UnauthorizedError:
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	case *sync.ForbiddenError:
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	case *sync.NotFoundError:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	case *sync.ServiceUnavailableError:
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func prepareResponse(w http.ResponseWriter, ct string) {
	w.Header().Set(encoding.ContentTypeHeader, ct)
}

func extractParams(r *http.Request) map[string]string {
	par := httprouter.ParamsFromContext(r.Context())
	if len(par) == 0 {
		return make(map[string]string)
	}
	p := make(map[string]string)
	for _, v := range par {
		p[v.Key] = v.Value
	}
	return p
}
