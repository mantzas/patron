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
func (c *Config) Set(key string, value interface{}) {
	viper.Set(key, value)
}

// Get returns the value of the key
func (c *Config) Get(key string) interface{} {
	return viper.Get(key)
}

// GetBool returns the bool value of the key
func (c *Config) GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetInt returns the int value of the key
func (c *Config) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetString returns the string value of the key
func (c *Config) GetString(key string) string {
	return viper.GetString(key)
}

// GetFloat64 returns the float64 value of the key
func (c *Config) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}
