package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	key = "key"
)

func TestSetup(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		c       Config
		wantErr bool
	}{
		{"failed with nil config", nil, true},
		{"success", newTestConfig(), false},
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
	err := Setup(newTestConfig())
	assert.NoError(err)
	value := "value"
	err = Set(key, value)
	assert.NoError(err)
	v, err := Get(key)
	assert.NoError(err)
	assert.Equal(value, v)
}

func TestBool(t *testing.T) {
	assert := assert.New(t)
	err := Setup(newTestConfig())
	assert.NoError(err)
	value := true
	err = Set(key, value)
	assert.NoError(err)
	v, err := GetBool(key)
	assert.NoError(err)
	assert.Equal(value, v)
}

func TestInt(t *testing.T) {
	assert := assert.New(t)
	err := Setup(newTestConfig())
	assert.NoError(err)
	value := int64(1)
	err = Set(key, value)
	assert.NoError(err)
	v, err := GetInt64(key)
	assert.NoError(err)
	assert.Equal(value, v)
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	err := Setup(newTestConfig())
	assert.NoError(err)
	value := "value"
	err = Set(key, value)
	assert.NoError(err)
	v, err := GetString(key)
	assert.NoError(err)
	assert.Equal(value, v)
}

func TestFloat64(t *testing.T) {
	assert := assert.New(t)
	err := Setup(newTestConfig())
	assert.NoError(err)
	value := 3.2
	err = Set(key, value)
	assert.NoError(err)
	v, err := GetFloat64(key)
	assert.NoError(err)
	assert.Equal(value, v)
}

type testConfig struct {
	store map[string]interface{}
}

func newTestConfig() *testConfig {
	return &testConfig{store: make(map[string]interface{})}
}

func (tc *testConfig) Set(key string, value interface{}) error {
	tc.store[key] = value
	return nil
}

// Get returns the value of the key
func (tc *testConfig) Get(key string) (interface{}, error) {
	return tc.store[key].(string), nil
}

// GetBool returns the bool value of the key
func (tc *testConfig) GetBool(key string) (bool, error) {
	return tc.store[key].(bool), nil
}

// GetInt returns the int value of the key
func (tc *testConfig) GetInt64(key string) (int64, error) {
	return tc.store[key].(int64), nil
}

// GetString returns the string value of the key
func (tc *testConfig) GetString(key string) (string, error) {
	return tc.store[key].(string), nil
}

// GetFloat64 returns the float64 value of the key
func (tc *testConfig) GetFloat64(key string) (float64, error) {
	return tc.store[key].(float64), nil
}
