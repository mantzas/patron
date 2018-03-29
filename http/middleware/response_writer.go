package middleware

import "net/http"

// ResponseWriter wrapper around the ResponseWriter to expose status
type ResponseWriter struct {
	status              int
	statusHeaderWritten bool
	w                   http.ResponseWriter
}

// NewResponseWriter constructor
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{-1, false, w}
}

// Status returns the http response status
func (w *ResponseWriter) Status() int {
	return w.status
}

// Header returns the header
func (w *ResponseWriter) Header() http.Header {
	return w.w.Header()
}

// Write to the internal ResponseWriter and sets the status if not set already
func (w *ResponseWriter) Write(d []byte) (int, error) {

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
func (w *ResponseWriter) WriteHeader(code int) {
	w.status = code
	w.w.WriteHeader(code)
	w.statusHeaderWritten = true
}
