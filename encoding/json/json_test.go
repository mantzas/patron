package json

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	j, err := Encode("string")
	assert.NoError(t, err)
	b := bytes.NewBuffer(j)
	var data string
	err = Decode(b, &data)
	assert.NoError(t, err)
	assert.Equal(t, "string", data)
	err = DecodeRaw(j, &data)
	assert.NoError(t, err)
	assert.Equal(t, "string", data)
}
