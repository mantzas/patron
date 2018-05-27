package http

import (
	"net/http"

	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/sync"
	"github.com/pkg/errors"
)

const (
	// AcceptHeader HTTP constant
	AcceptHeader string = "Accept"
	// ContentTypeHeader HTTP constant
	ContentTypeHeader string = "Content-Type"

	// JSONContentType JSON definition
	JSONContentType string = "application/json"
	// JSONContentTypeCharset JSON definition with charset
	JSONContentTypeCharset string = "application/json; charset=utf-8"
)

func handler(hnd sync.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		h := extractHeaders(r)

		dec, enc, err := determineEncoding(h)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		req := sync.NewRequest(h, extractFields(r), r.Body, dec)

		rsp, err := hnd.Handle(r.Context(), req)
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

func extractHeaders(r *http.Request) map[string]string {

	h := make(map[string]string)

	for name, hdr := range r.Header {
		h[name] = hdr[0]
	}
	return h
}

func extractFields(r *http.Request) map[string]string {

	f := make(map[string]string)

	for name, values := range r.URL.Query() {
		f[name] = values[0]
	}
	return f
}

func determineEncoding(hdr map[string]string) (encoding.Decode, encoding.Encode, error) {

	c, err := determineContentType(hdr)
	if err != nil {
		return nil, nil, err
	}

	switch c {
	case JSONContentType, JSONContentTypeCharset:
		return json.Decode, json.Encode, nil

	}
	return nil, nil, errors.Errorf("accept header %s is unsupported", c)
}

func determineContentType(hdr map[string]string) (string, error) {
	ah, aOk := hdr[AcceptHeader]
	ch, cOk := hdr[ContentTypeHeader]
	if !aOk && !cOk {
		return "", errors.New("accept and content type header is missing")

	}

	if (aOk && cOk) && (ah != ch) {
		return "", errors.New("accept and content type header are different")
	}

	if ah != "" {
		return ah, nil
	}
	return ch, nil
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

	w.Write(p)
	return nil
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
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
