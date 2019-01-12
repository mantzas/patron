package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/sync/http/auth"
	"github.com/stretchr/testify/assert"
)

func testHandle(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(202)
}

func testPanicHandleString(_ http.ResponseWriter, _ *http.Request) {
	panic("test")
}

func testPanicHandleError(_ http.ResponseWriter, _ *http.Request) {
	panic(errors.New("TEST"))
}

func testPanicHandleInt(_ http.ResponseWriter, _ *http.Request) {
	panic(1000)
}

func TestMiddleware(t *testing.T) {
	r, err := http.NewRequest("POST", "/test", nil)
	assert.NoError(t, err)

	type args struct {
		next  http.HandlerFunc
		trace bool
		auth  auth.Authenticator
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{
			name:         "middleware success",
			args:         args{next: testHandle, trace: false, auth: &MockAuthenticator{success: true}},
			expectedCode: 202,
		},
		{
			name:         "middleware trace success",
			args:         args{next: testHandle, trace: true, auth: &MockAuthenticator{success: true}},
			expectedCode: 202,
		},
		{
			name:         "middleware panic string",
			args:         args{next: testPanicHandleString, trace: true, auth: &MockAuthenticator{success: true}},
			expectedCode: 500,
		},
		{
			name:         "middleware panic error",
			args:         args{next: testPanicHandleError, trace: true, auth: &MockAuthenticator{success: true}},
			expectedCode: 500,
		},
		{
			name:         "middleware panic other",
			args:         args{next: testPanicHandleInt, trace: true, auth: &MockAuthenticator{success: true}},
			expectedCode: 500,
		},
		{
			name:         "middleware auth error",
			args:         args{next: testPanicHandleInt, trace: true, auth: &MockAuthenticator{err: errors.New("TEST")}},
			expectedCode: 500,
		},
		{
			name:         "middleware auth failure",
			args:         args{next: testPanicHandleInt, trace: true, auth: &MockAuthenticator{success: false}},
			expectedCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			Middleware(tt.args.trace, tt.args.auth, "path", tt.args.next)(resp, r)
			assert.Equal(t, tt.expectedCode, resp.Code)
		})
	}
}

func TestResponseWriter(t *testing.T) {
	rc := httptest.NewRecorder()
	rw := newResponseWriter(rc)

	_, err := rw.Write([]byte("test"))
	assert.NoError(t, err)
	rw.WriteHeader(202)

	assert.Equal(t, 202, rw.status, "status expected 202 but got %d", rw.status)
	assert.Len(t, rw.Header(), 1, "header count expected to be 1")
	assert.True(t, rw.statusHeaderWritten, "expected to be true")
	assert.Equal(t, "test", rc.Body.String(), "body expected to be test but was %s", rc.Body.String())
}
