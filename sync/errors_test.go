package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_Error(t *testing.T) {
	assert := assert.New(t)
	v := NewValidationError("TEST")
	assert.Equal("TEST", v.Error())
}
func TestUnauthorizedError_Error(t *testing.T) {
	assert := assert.New(t)
	v := NewUnauthorizedError("TEST")
	assert.Equal("TEST", v.Error())
}

func TestForbiddenError_Error(t *testing.T) {
	assert := assert.New(t)
	v := NewForbiddenError("TEST")
	assert.Equal("TEST", v.Error())
}

func TestNotFoundError_Error(t *testing.T) {
	assert := assert.New(t)
	v := NewNotFoundError("TEST")
	assert.Equal("TEST", v.Error())
}

func TestServiceUnavailableError_Error(t *testing.T) {
	assert := assert.New(t)
	v := NewServiceUnavailableError("TEST")
	assert.Equal("TEST", v.Error())
}
