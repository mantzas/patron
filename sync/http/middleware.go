package http

import (
	"errors"
	"net/http"

	"github.com/mantzas/patron/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type responseWriter struct {
	status              int
	statusHeaderWritten bool
	w                   http.ResponseWriter
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{-1, false, w}
}

// Status returns the http response status
func (w *responseWriter) Status() int {
	return w.status
}

// Header returns the header
func (w *responseWriter) Header() http.Header {
	return w.w.Header()
}

// Write to the internal ResponseWriter and sets the status if not set already
func (w *responseWriter) Write(d []byte) (int, error) {

	value, err := w.w.Write(d)
	if err != nil {
		return value, err
	}

	if !w.statusHeaderWritten {
		w.status = http.StatusOK
		w.statusHeaderWritten = true
	}

	return value, err
}

// WriteHeader writes the internal header and saves the status for retrieval
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.w.WriteHeader(code)
	w.statusHeaderWritten = true
}

// DefaultMiddleware which handles Logging and Recover middleware
func DefaultMiddleware(tr opentracing.Tracer, path string, next http.HandlerFunc) http.HandlerFunc {
	return TracingMiddleware(tr, path, RecoveryMiddleware(next))
}

// TracingMiddleware for handling tracing and metrics
func TracingMiddleware(tr opentracing.Tracer, path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, _ := tr.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		sp := tr.StartSpan(opName(r.Method, path), ext.RPCServerOption(ctx))
		ext.HTTPMethod.Set(sp, r.Method)
		ext.HTTPUrl.Set(sp, r.URL.String())
		ext.Component.Set(sp, "http")
		r = r.WithContext(opentracing.ContextWithSpan(r.Context(), sp))
		lw := newResponseWriter(w)

		next(lw, r)

		ext.HTTPStatusCode.Set(sp, uint16(lw.Status()))
		sp.Finish()
	}
}

func opName(method, path string) string {
	return "HTTP " + method + " " + path
}

// RecoveryMiddleware for recovering from failed requests
func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("unknown panic")
				}
				log.Errorf("recovering from an error %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}
