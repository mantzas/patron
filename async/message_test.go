package async

import (
	"bytes"
	"testing"

	"github.com/mantzas/patron/encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(NewMessage(nil, nil, nil))
}

func TestMessage_Decode(t *testing.T) {
	assert := assert.New(t)

	j, err := json.Encode("string")
	assert.NoError(err)

	b := bytes.NewBuffer(j)

	req := NewMessage(nil, b, json.Decode)
	assert.NotNil(req)

	var data string

	err = req.Decode(&data)
	assert.NoError(err)

	assert.Equal("string", data)
}
