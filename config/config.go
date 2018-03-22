package config

import (
	"errors"
)

// Config interface which has to be implemented in order to be used here
type Config interface {
	Get(key string) interface{}
	GetBool(key string) bool
	GetInt(key string) int
	GetString(key string) string
	GetFloat32(key string) float32
	GetFloat64(key string) float64
}

var config Config

func init() {
	config = NewMapConfig()
}

// Setup set's up a new config to the global state
func Setup(c Config) error {
	if c == nil {
		return errors.New("config is nil")
	}

	config = c
	return nil
}

// Get returns the value of the key
func Get(key string) interface{} {
	return config.Get(key)
}

// GetBool returns the bool value of the key
func GetBool(key string) bool {
	return config.GetBool(key)
}

// GetInt returns the int value of the key
func GetInt(key string) int {
	return config.GetInt(key)
}

// GetString returns the string value of the key
func GetString(key string) string {
	return config.GetString(key)
}

// GetFloat32 returns the float32 value of the key
func GetFloat32(key string) float32 {
	return config.GetFloat32(key)
}

// GetFloat64 returns the float64 value of the key
func GetFloat64(key string) float64 {
	return config.GetFloat64(key)
}
