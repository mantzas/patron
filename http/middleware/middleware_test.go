package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Setup(zerolog.DefaultFactory(log.DebugLevel))
}

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

			DefaultMiddleware(tt.args.next)(tt.args.w, r)

			assert.Equal(tt.expectedCode, tt.args.w.Code, "default middleware expected %d but got %d", tt.expectedCode, tt.args.w.Code)
		})
	}
}
