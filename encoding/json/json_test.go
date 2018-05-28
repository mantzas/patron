package json

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	assert := assert.New(t)
	j, err := Encode("string")
	assert.NoError(err)
	b := bytes.NewBuffer(j)
	var data string
	err = Decode(b, &data)
	assert.NoError(err)
	assert.Equal("string", data)

}
