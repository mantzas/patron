package viper

import (
	"github.com/mantzas/patron/config"
	"github.com/spf13/viper"
)

// Config defines a viper backed config
type Config struct {
}

// New returns a new viper config
func New() config.Config {
	return &Config{}
}

// Set the value under the key in config
func (c *Config) Set(key string, value interface{}) error {
	viper.Set(key, value)
	return nil
}

// Get returns the value of the key
func (c *Config) Get(key string) (interface{}, error) {
	return viper.Get(key), nil
}

// GetBool returns the bool value of the key
func (c *Config) GetBool(key string) (bool, error) {
	return viper.GetBool(key), nil
}

// GetInt64 returns the int64 value of the key
func (c *Config) GetInt64(key string) (int64, error) {
	return viper.GetInt64(key), nil
}

// GetString returns the string value of the key
func (c *Config) GetString(key string) (string, error) {
	return viper.GetString(key), nil
}

// GetFloat64 returns the float64 value of the key
func (c *Config) GetFloat64(key string) (float64, error) {
	return viper.GetFloat64(key), nil
}
