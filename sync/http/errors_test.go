package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError(t *testing.T) {
	v := NewValidationError()
	assert.Equal(t, "Bad Request", v.Error())
	assert.Equal(t, 400, v.code)
}
func TestUnauthorizedError(t *testing.T) {
	v := NewUnauthorizedError()
	assert.Equal(t, "Unauthorized", v.Error())
	assert.Equal(t, 401, v.code)
}

func TestForbiddenError(t *testing.T) {
	v := NewForbiddenError()
	assert.Equal(t, "Forbidden", v.Error())
	assert.Equal(t, 403, v.code)
}

func TestNotFoundError(t *testing.T) {
	v := NewNotFoundError()
	assert.Equal(t, "Not Found", v.Error())
	assert.Equal(t, 404, v.code)
}

func TestServiceUnavailableError(t *testing.T) {
	v := NewServiceUnavailableError()
	assert.Equal(t, "Service Unavailable", v.Error())
	assert.Equal(t, 503, v.code)
}
