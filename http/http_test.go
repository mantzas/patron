package http

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateHTTPServer(t *testing.T) {
	assert := assert.New(t)
	s := CreateHTTPServer(10000, http.DefaultServeMux)
	assert.Equal(":10000", s.Addr)
	assert.Equal(5*time.Second, s.ReadTimeout)
	assert.Equal(10*time.Second, s.WriteTimeout)
	assert.Equal(120*time.Second, s.IdleTimeout)
	assert.Equal(s.Handler, http.DefaultServeMux)
}
