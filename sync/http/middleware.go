package http

import (
	"errors"
	"net/http"

	"time"

	"github.com/mantzas/patron/log"
)

// DefaultMiddleware which handles Logging and Recover middleware
func DefaultMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return LoggingMiddleware(RecoveryMiddleware(next))
}

// LoggingMiddleware for recovering from failed requests
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		lw := NewResponseWriter(w)
		startTime := time.Now()
		next(lw, r)
		log.Infof("method=%s route=%s status=%d time=%s", r.Method, r.URL.String(), lw.Status(), time.Since(startTime))
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

				log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}
