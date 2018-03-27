package http

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResponseWriter(t *testing.T) {

	req := require.New(t)
	rc := httptest.NewRecorder()
	rw := NewResponseWriter(rc)

	rw.Write([]byte("test"))
	rw.WriteHeader(202)

	req.Equal(202, rw.status, "status expected 202 but got %d", rw.status)
	req.Len(rw.Header(), 1, "header count expected to be 1")
	req.True(rw.statusHeaderWritten, "expected to be true")
	req.Equal("test", rc.Body.String(), "body expected to be test but was %s", rc.Body.String())
}
