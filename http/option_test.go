package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPorts(t *testing.T) {

	assert := assert.New(t)
	s, err := New(Ports(40000, 40001))
	assert.NoError(err)
	assert.NotNil(s)
	assert.Equal(":40000", s.srv.Addr)
}

func TestPortsSamePorts(t *testing.T) {

	assert := assert.New(t)
	_, err := New(Ports(40000, 40000))
	assert.Error(err)
}
