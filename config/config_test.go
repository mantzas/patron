package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		c       Config
		wantErr bool
	}{
		{"failed with nil config", nil, true},
		{"success", NewMapConfig(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := Setup(tt.c)

			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	c := NewMapConfig()
	assert.NoError(Setup(c))
	key := "key"
	value := "value"
	c.Set(key, value)
	v := Get(key)
	assert.Equal(value, v)
}

func TestBool(t *testing.T) {
	assert := assert.New(t)
	c := NewMapConfig()
	assert.NoError(Setup(c))
	key := "key"
	value := true
	c.Set(key, value)
	v := GetBool(key)
	assert.Equal(value, v)
}

func TestInt(t *testing.T) {
	assert := assert.New(t)
	c := NewMapConfig()
	assert.NoError(Setup(c))
	key := "key"
	value := 1
	c.Set(key, value)
	v := GetInt(key)
	assert.Equal(value, v)
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	c := NewMapConfig()
	assert.NoError(Setup(c))
	key := "key"
	value := "value"
	c.Set(key, value)
	v := GetString(key)
	assert.Equal(value, v)
}

func TestFloat32(t *testing.T) {
	assert := assert.New(t)
	c := NewMapConfig()
	assert.NoError(Setup(c))
	key := "key"
	value := float32(3.2)
	c.Set(key, value)
	v := GetFloat32(key)
	assert.Equal(value, v)
}

func TestFloat64(t *testing.T) {
	assert := assert.New(t)
	c := NewMapConfig()
	assert.NoError(Setup(c))
	key := "key"
	value := 3.2
	c.Set(key, value)
	v := GetFloat64(key)
	assert.Equal(value, v)
}
