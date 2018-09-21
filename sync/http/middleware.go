package http

import (
	"net/http"

	"github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace"
)

type responseWriter struct {
	status              int
	statusHeaderWritten bool
	writer              http.ResponseWriter
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{status: -1, statusHeaderWritten: false, writer: w}
}

// Status returns the http response status.
func (w *responseWriter) Status() int {
	return w.status
}

// Header returns the header.
func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write to the internal ResponseWriter and sets the status if not set already.
func (w *responseWriter) Write(d []byte) (int, error) {

	value, err := w.writer.Write(d)
	if err != nil {
		return value, err
	}

	if !w.statusHeaderWritten {
		w.status = http.StatusOK
		w.statusHeaderWritten = true
	}

	return value, err
}

// WriteHeader writes the internal header and saves the status for retrieval.
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.writer.WriteHeader(code)
	w.statusHeaderWritten = true
}

// DefaultMiddleware which handles tracing and recovery.
func DefaultMiddleware(path string, next http.HandlerFunc) http.HandlerFunc {
	return TracingMiddleware(path, RecoveryMiddleware(next))
}

// TracingMiddleware for handling tracing and metrics.
func TracingMiddleware(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sp, r := trace.HTTPSpan(path, r)
		lw := newResponseWriter(w)
		next(lw, r)
		trace.FinishHTTPSpan(sp, lw.Status())
	}
}

// RecoveryMiddleware for recovering from failed requests.
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
				_ = err
				log.Errorf("recovering from an error %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}
