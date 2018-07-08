package config

import (
	"errors"
)

// Config interface defining of a config implementation.
type Config interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	GetBool(key string) (bool, error)
	GetInt64(key string) (int64, error)
	GetString(key string) (string, error)
	GetFloat64(key string) (float64, error)
}

var config Config

// Setup a new config to the global state.
func Setup(c Config) error {
	if c == nil {
		return errors.New("config is nil")
	}

	config = c
	return nil
}

// Set a key and a value to config.
func Set(key string, value interface{}) error {
	return config.Set(key, value)
}

// Get returns the value of the key.
func Get(key string) (interface{}, error) {
	return config.Get(key)
}

// GetBool returns the bool value of the key.
func GetBool(key string) (bool, error) {
	return config.GetBool(key)
}

// GetInt64 returns the int64 value of the key.
func GetInt64(key string) (int64, error) {
	return config.GetInt64(key)
}

// GetString returns the string value of the key.
func GetString(key string) (string, error) {
	return config.GetString(key)
}

// GetFloat64 returns the float64 value of the key.
func GetFloat64(key string) (float64, error) {
	return config.GetFloat64(key)
}
