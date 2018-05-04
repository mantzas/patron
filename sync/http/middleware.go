package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/mantzas/patron/log"
)

// DefaultMiddleware which handles Logging and Recover middleware
func DefaultMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return LoggingMetricMiddleware(RecoveryMiddleware(next))
}

// LoggingMetricMiddleware for handling logging and metrics
func LoggingMetricMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lw := NewResponseWriter(w)
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
