package amqp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	assert := assert.New(t)
	m := NewMessage("xxx", []byte("test"))
	assert.Equal("xxx", m.contentType)
	assert.Equal([]byte("test"), m.body)
}

func TestNewJSONMessage(t *testing.T) {
	assert := assert.New(t)
	m, err := NewJSONMessage("xxx")
	assert.NoError(err)
	assert.Equal("application/json", m.contentType)
	assert.Equal([]byte(`"xxx"`), m.body)
	_, err = NewJSONMessage(make(chan bool))
	assert.Error(err)
}
