package http

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseWriter(t *testing.T) {
	assert := assert.New(t)
	rc := httptest.NewRecorder()
	rw := NewResponseWriter(rc)

	rw.Write([]byte("test"))
	rw.WriteHeader(202)

	assert.Equal(202, rw.status, "status expected 202 but got %d", rw.status)
	assert.Len(rw.Header(), 1, "header count expected to be 1")
	assert.True(rw.statusHeaderWritten, "expected to be true")
	assert.Equal("test", rc.Body.String(), "body expected to be test but was %s", rc.Body.String())
}
