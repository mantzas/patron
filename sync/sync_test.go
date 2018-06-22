package sync

import (
	"bytes"
	"testing"

	"github.com/mantzas/patron/encoding/json"
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

func TestNewRequest(t *testing.T) {
	assert := assert.New(t)
	req := NewRequest(nil, nil, nil)
	assert.NotNil(req)
}

func TestRequest_Decode(t *testing.T) {
	assert := assert.New(t)
	j, err := json.Encode("string")
	assert.NoError(err)
	b := bytes.NewBuffer(j)
	req := NewRequest(nil, b, json.Decode)
	assert.NotNil(req)
	var data string
	err = req.Decode(&data)
	assert.NoError(err)
	assert.Equal("string", data)
}

func TestNewResponse(t *testing.T) {
	assert := assert.New(t)
	rsp := NewResponse("test")
	assert.NotNil(rsp)
	assert.IsType("test", rsp.Payload)
}
