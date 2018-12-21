package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/encoding/protobuf"
	"github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/sync"
)

func handler(hnd sync.ProcessorFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ct, dec, enc, err := determineEncoding(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}
		prepareResponse(w, ct)

		f := extractFields(r)
		for k, v := range extractParams(r) {
			f[k] = v
		}

		ctx := log.WithContext(r.Context(), log.Sub(map[string]interface{}{"requestID": uuid.New().String()}))
		req := sync.NewRequest(f, r.Body, dec)
		rsp, err := hnd(ctx, req)
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

func determineEncoding(r *http.Request) (string, encoding.DecodeFunc, encoding.EncodeFunc, error) {
	cth, cok := r.Header[encoding.ContentTypeHeader]
	ach, aok := r.Header[encoding.AcceptHeader]

	// No headers default to JSON
	if !cok && !aok {
		return json.TypeCharset, json.Decode, json.Encode, nil
	}

	var enc encoding.EncodeFunc
	var dec encoding.DecodeFunc
	var ct string

	if cok {
		switch cth[0] {
		case "*/*", json.Type, json.TypeCharset:
			enc = json.Encode
			dec = json.Decode
			ct = json.TypeCharset
		case protobuf.Type, protobuf.TypeGoogle:
			enc = protobuf.Encode
			dec = protobuf.Decode
			ct = protobuf.Type
		default:
			return "", nil, nil, errors.New("content type header not supported")
		}
	}

	if aok {
		switch ach[0] {
		case "*/*", json.Type, json.TypeCharset:
			enc = json.Encode
			if dec == nil {
				dec = json.Decode
			}
			ct = json.TypeCharset
		case protobuf.Type, protobuf.TypeGoogle:
			enc = protobuf.Encode
			if dec == nil {
				dec = protobuf.Decode
			}
			ct = protobuf.Type
		default:
			return "", nil, nil, errors.New("accept header not supported")
		}
	}

	return ct, dec, enc, nil
}

func extractFields(r *http.Request) map[string]string {
	f := make(map[string]string)

	for name, values := range r.URL.Query() {
		f[name] = values[0]
	}
	return f
}

func handleSuccess(w http.ResponseWriter, r *http.Request, rsp *sync.Response, enc encoding.EncodeFunc) error {
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
		http.Error(w, "", http.StatusNotFound)
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
