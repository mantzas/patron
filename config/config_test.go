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
	Setup(newTestConfig())
	key := "key"
	value := "value"
	Set(key, value)
	v := Get(key)
	assert.Equal(value, v)
}

func TestBool(t *testing.T) {
	assert := assert.New(t)
	Setup(newTestConfig())
	key := "key"
	value := true
	Set(key, value)
	v := GetBool(key)
	assert.Equal(value, v)
}

func TestInt(t *testing.T) {
	assert := assert.New(t)
	Setup(newTestConfig())
	key := "key"
	value := 1
	Set(key, value)
	v := GetInt(key)
	assert.Equal(value, v)
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	Setup(newTestConfig())
	key := "key"
	value := "value"
	Set(key, value)
	v := GetString(key)
	assert.Equal(value, v)
}

func TestFloat64(t *testing.T) {
	assert := assert.New(t)
	Setup(newTestConfig())
	key := "key"
	value := 3.2
	Set(key, value)
	v := GetFloat64(key)
	assert.Equal(value, v)
}

type testConfig struct {
	store map[string]interface{}
}

func newTestConfig() *testConfig {
	return &testConfig{make(map[string]interface{}, 0)}
}

func (tc *testConfig) Set(key string, value interface{}) {
	tc.store[key] = value
}

// Get returns the value of the key
func (tc *testConfig) Get(key string) interface{} {
	return tc.store[key].(string)
}

// GetBool returns the bool value of the key
func (tc *testConfig) GetBool(key string) bool {
	return tc.store[key].(bool)
}

// GetInt returns the int value of the key
func (tc *testConfig) GetInt(key string) int {
	return tc.store[key].(int)
}

// GetString returns the string value of the key
func (tc *testConfig) GetString(key string) string {
	return tc.store[key].(string)
}

// GetFloat64 returns the float64 value of the key
func (tc *testConfig) GetFloat64(key string) float64 {
	return tc.store[key].(float64)
}
