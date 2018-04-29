package viper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	c := New()
	assert.NotNil(c)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	key := "key"
	value := "value"
	c := New()
	err := c.Set(key, value)
	assert.NoError(err)
	v, err := c.Get(key)
	assert.Equal(value, v)
	assert.NoError(err)
}

func TestBool(t *testing.T) {
	assert := assert.New(t)
	key := "key"
	value := true
	c := New()
	err := c.Set(key, value)
	assert.NoError(err)
	v, err := c.GetBool(key)
	assert.Equal(value, v)
	assert.NoError(err)
}

func TestInt64(t *testing.T) {
	assert := assert.New(t)
	key := "key"
	value := int64(1)
	c := New()
	err := c.Set(key, value)
	assert.NoError(err)
	v, err := c.GetInt64(key)
	assert.Equal(value, v)
	assert.NoError(err)
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	key := "key"
	value := "value"
	c := New()
	err := c.Set(key, value)
	assert.NoError(err)
	v, err := c.GetString(key)
	assert.Equal(value, v)
	assert.NoError(err)
}

func TestFloat64(t *testing.T) {
	assert := assert.New(t)
	key := "key"
	value := 3.2
	c := New()
	err := c.Set(key, value)
	assert.NoError(err)
	v, err := c.GetFloat64(key)
	assert.Equal(value, v)
	assert.NoError(err)
}
