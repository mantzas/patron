package http

import (
	"errors"
	"net/http"

	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace"
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
func DefaultMiddleware(path string, next http.HandlerFunc) http.HandlerFunc {
	return TracingMiddleware(path, RecoveryMiddleware(next))
}

// TracingMiddleware for handling tracing and metrics
func TracingMiddleware(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sp := trace.StartHTTPSpan(path, r)
		lw := newResponseWriter(w)
		next(lw, r)
		trace.FinishHTTPSpan(sp, lw.Status())
	}
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
