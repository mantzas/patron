package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/mantzas/patron/log"
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
func DefaultMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return LoggingMetricMiddleware(RecoveryMiddleware(next))
}

// LoggingMetricMiddleware for handling logging and metrics
func LoggingMetricMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lw := newResponseWriter(w)
		st := time.Now()
		next(lw, r)
		latency := float64(time.Since(st)) / float64(time.Millisecond)
		status := strconv.Itoa(lw.Status())
		log.Infof("method=%s route=%s status=%s time=%f", r.Method, r.URL.Path, status, latency)
		recordMetric(r.Context(), r.URL.Host, r.Method, r.URL.Path, status, latency)
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
