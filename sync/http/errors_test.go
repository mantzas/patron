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

func TestValidationErrorWithPayload(t *testing.T) {
	v := NewValidationErrorWithPayload("test")
	assert.Equal(t, "test", v.Error())
	assert.Equal(t, 400, v.code)
}
func TestUnauthorizedError(t *testing.T) {
	v := NewUnauthorizedError()
	assert.Equal(t, "Unauthorized", v.Error())
	assert.Equal(t, 401, v.code)
}

func TestUnauthorizedErrorWithPayload(t *testing.T) {
	v := NewUnauthorizedErrorWithPayload("test")
	assert.Equal(t, "test", v.Error())
	assert.Equal(t, 401, v.code)
}

func TestForbiddenError(t *testing.T) {
	v := NewForbiddenError()
	assert.Equal(t, "Forbidden", v.Error())
	assert.Equal(t, 403, v.code)
}

func TestForbiddenErrorWithPayload(t *testing.T) {
	v := NewForbiddenErrorWithPayload("test")
	assert.Equal(t, "test", v.Error())
	assert.Equal(t, 403, v.code)
}

func TestNotFoundError(t *testing.T) {
	v := NewNotFoundError()
	assert.Equal(t, "Not Found", v.Error())
	assert.Equal(t, 404, v.code)
}

func TestNotFoundErrorWithPayload(t *testing.T) {
	v := NewNotFoundErrorWithPayload("test")
	assert.Equal(t, "test", v.Error())
	assert.Equal(t, 404, v.code)
}

func TestServiceUnavailableError(t *testing.T) {
	v := NewServiceUnavailableError()
	assert.Equal(t, "Service Unavailable", v.Error())
	assert.Equal(t, 503, v.code)
}

func TestServiceUnavailableErrorWithPayload(t *testing.T) {
	v := NewServiceUnavailableErrorWithPayload("test")
	assert.Equal(t, "test", v.Error())
	assert.Equal(t, 503, v.code)
}

func TestNewError(t *testing.T) {
	v := NewError()
	assert.Equal(t, "Internal Server Error", v.Error())
	assert.Equal(t, 500, v.code)
}

func TestNewErrorWithCodeAndPayload(t *testing.T) {
	v := NewErrorWithCodeAndPayload(409, "Conflict")
	assert.Equal(t, "Conflict", v.Error())
	assert.Equal(t, 409, v.code)
}
