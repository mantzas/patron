package async

import (
	"testing"

	"github.com/mantzas/patron/encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(NewMessage(nil, nil))
}

func TestMessage_Decode(t *testing.T) {
	assert := assert.New(t)

	j, err := json.Encode("string")
	assert.NoError(err)

	req := NewMessage(j, json.DecodeRaw)
	assert.NotNil(req)

	var data string

	err = req.Decode(&data)
	assert.NoError(err)

	assert.Equal("string", data)
}
