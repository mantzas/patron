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
	c.Set(key, value)
	v := c.Get(key)
	assert.Equal(value, v)
}

func TestBool(t *testing.T) {
	assert := assert.New(t)
	key := "key"
	value := true
	c := New()
	c.Set(key, value)
	v := c.GetBool(key)
	assert.Equal(value, v)
}

func TestInt(t *testing.T) {
	assert := assert.New(t)
	key := "key"
	value := 1
	c := New()
	c.Set(key, value)
	v := c.GetInt(key)
	assert.Equal(value, v)
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	key := "key"
	value := "value"
	c := New()
	c.Set(key, value)
	v := c.GetString(key)
	assert.Equal(value, v)
}

func TestFloat64(t *testing.T) {
	assert := assert.New(t)
	key := "key"
	value := 3.2
	c := New()
	c.Set(key, value)
	v := c.GetFloat64(key)
	assert.Equal(value, v)
}
