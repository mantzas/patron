package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
		resp *httptest.ResponseRecorder
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{"default middleware success", args{next: testHandle, resp: httptest.NewRecorder()}, 202},
		{"default middleware panic string", args{next: testPanicHandleString, resp: httptest.NewRecorder()}, 500},
		{"default middleware panic error", args{next: testPanicHandleError, resp: httptest.NewRecorder()}, 500},
		{"default middleware panic other", args{next: testPanicHandleInt, resp: httptest.NewRecorder()}, 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DefaultMiddleware("path", tt.args.next)(tt.args.resp, r)
			assert.Equal(tt.expectedCode, tt.args.resp.Code)
		})
	}
}

func TestResponseWriter(t *testing.T) {
	assert := assert.New(t)
	rc := httptest.NewRecorder()
	rw := newResponseWriter(rc)

	_, err := rw.Write([]byte("test"))
	assert.NoError(err)
	rw.WriteHeader(202)

	assert.Equal(202, rw.status, "status expected 202 but got %d", rw.status)
	assert.Len(rw.Header(), 1, "header count expected to be 1")
	assert.True(rw.statusHeaderWritten, "expected to be true")
	assert.Equal("test", rc.Body.String(), "body expected to be test but was %s", rc.Body.String())
}
