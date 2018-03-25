package pprof

import (
	_ "net/http/pprof"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	s := New(10001)
	assert.NotNil(s.srv)
	assert.Equal(":10001", s.srv.Addr)
}

func TestServer_SetPort(t *testing.T) {
	assert := assert.New(t)
	s := New(10001)
	s.SetPort(10002)
	assert.NotNil(s.srv)
	assert.Equal(":10002", s.srv.Addr)
}
