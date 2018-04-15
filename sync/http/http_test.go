package http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_createHTTPServer(t *testing.T) {
	assert := assert.New(t)
	s := createHTTPServer(10000, nil)
	assert.NotNil(s)
	assert.Equal(":10000", s.Addr)
	assert.Equal(5*time.Second, s.ReadTimeout)
	assert.Equal(60*time.Second, s.WriteTimeout)
	assert.Equal(120*time.Second, s.IdleTimeout)
}
