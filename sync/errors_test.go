package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_Error(t *testing.T) {
	v := NewValidationError("TEST")
	assert.Equal(t, "TEST", v.Error())
}
func TestUnauthorizedError_Error(t *testing.T) {
	v := NewUnauthorizedError("TEST")
	assert.Equal(t, "TEST", v.Error())
}

func TestForbiddenError_Error(t *testing.T) {
	v := NewForbiddenError("TEST")
	assert.Equal(t, "TEST", v.Error())
}

func TestNotFoundError_Error(t *testing.T) {
	v := NewNotFoundError("TEST")
	assert.Equal(t, "TEST", v.Error())
}

func TestServiceUnavailableError_Error(t *testing.T) {
	v := NewServiceUnavailableError("TEST")
	assert.Equal(t, "TEST", v.Error())
}
