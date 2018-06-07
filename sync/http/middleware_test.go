package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-client-go"
)

func testHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(202)
}

func testPanicHandleString(w http.ResponseWriter, r *http.Request) {
	panic("test")
}

func testPanicHandleError(w http.ResponseWriter, r *http.Request) {
	panic(errors.New("TEST"))
}

func testPanicHandleInt(w http.ResponseWriter, r *http.Request) {
	panic(1000)
}

func TestMiddleware(t *testing.T) {
	assert := assert.New(t)
	r, _ := http.NewRequest("POST", "/test", nil)

	type args struct {
		next http.HandlerFunc
		w    *httptest.ResponseRecorder
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{"default middleware success", args{testHandle, httptest.NewRecorder()}, 202},
		{"default middleware panic string", args{testPanicHandleString, httptest.NewRecorder()}, 500},
		{"default middleware panic error", args{testPanicHandleError, httptest.NewRecorder()}, 500},
		{"default middleware panic other", args{testPanicHandleInt, httptest.NewRecorder()}, 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reporter := jaeger.NewInMemoryReporter()
			tr, trCloser := jaeger.NewTracer("test", jaeger.NewConstSampler(true), reporter)
			defer trCloser.Close()

			DefaultMiddleware(tr, "path", tt.args.next)(tt.args.w, r)

			assert.Equal(tt.expectedCode, tt.args.w.Code, "default middleware expected %d but got %d", tt.expectedCode, tt.args.w.Code)
		})
	}
}

func TestResponseWriter(t *testing.T) {
	assert := assert.New(t)
	rc := httptest.NewRecorder()
	rw := newResponseWriter(rc)

	rw.Write([]byte("test"))
	rw.WriteHeader(202)

	assert.Equal(202, rw.status, "status expected 202 but got %d", rw.status)
	assert.Len(rw.Header(), 1, "header count expected to be 1")
	assert.True(rw.statusHeaderWritten, "expected to be true")
	assert.Equal("test", rc.Body.String(), "body expected to be test but was %s", rc.Body.String())
}
