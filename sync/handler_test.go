package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_Error(t *testing.T) {
	assert := assert.New(t)
	v := ValidationError{"TEST"}
	assert.Equal("TEST", v.Error())
}
func TestUnauthorizedError_Error(t *testing.T) {
	assert := assert.New(t)
	v := UnauthorizedError{"TEST"}
	assert.Equal("TEST", v.Error())
}

func TestForbiddenError_Error(t *testing.T) {
	assert := assert.New(t)
	v := ForbiddenError{"TEST"}
	assert.Equal("TEST", v.Error())
}

func TestNotFoundError_Error(t *testing.T) {
	assert := assert.New(t)
	v := NotFoundError{"TEST"}
	assert.Equal("TEST", v.Error())
}

func TestNewRequest(t *testing.T) {
	assert := assert.New(t)
	req := NewRequest(nil, nil, nil, nil)
	assert.NotNil(req)
}

func TestNewResponse(t *testing.T) {
	assert := assert.New(t)
	rsp := NewResponse("test")
	assert.NotNil(rsp)
	assert.IsType("test", rsp.Payload)
}
