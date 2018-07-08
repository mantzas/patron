package sync

import (
	"bytes"
	"testing"

	"github.com/mantzas/patron/encoding/json"
	"github.com/stretchr/testify/assert"
)

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
