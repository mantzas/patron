package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError(t *testing.T) {
	v := NewValidationError("TEST", "payload")
	assert.Equal(t, "TEST", v.Error())
	assert.Equal(t, 400, v.Code())
	assert.Equal(t, "payload", v.Payload())
}
func TestUnauthorizedError(t *testing.T) {
	v := NewUnauthorizedError("TEST", "payload")
	assert.Equal(t, "TEST", v.Error())
	assert.Equal(t, 401, v.Code())
	assert.Equal(t, "payload", v.Payload())
}

func TestForbiddenError(t *testing.T) {
	v := NewForbiddenError("TEST", "payload")
	assert.Equal(t, "TEST", v.Error())
	assert.Equal(t, 403, v.Code())
	assert.Equal(t, "payload", v.Payload())
}

func TestNotFoundError(t *testing.T) {
	v := NewNotFoundError("TEST", "payload")
	assert.Equal(t, "TEST", v.Error())
	assert.Equal(t, 404, v.Code())
	assert.Equal(t, "payload", v.Payload())
}

func TestServiceUnavailableError(t *testing.T) {
	v := NewServiceUnavailableError("TEST", "payload")
	assert.Equal(t, "TEST", v.Error())
	assert.Equal(t, 503, v.Code())
	assert.Equal(t, "payload", v.Payload())
}
