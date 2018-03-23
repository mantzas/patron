package config

// MapConfig defines a struct for handling configuration
// mapped in a map structure
type MapConfig struct {
	store map[string]interface{}
}

// NewMapConfig returns a new map config
func NewMapConfig() *MapConfig {
	return &MapConfig{make(map[string]interface{})}
}

// Set the key and value to the store
func (mc *MapConfig) Set(key string, value interface{}) {
	mc.store[key] = value
}

// Get returns the value of the key
func (mc *MapConfig) Get(key string) interface{} {
	return mc.store[key].(string)
}

// GetBool returns the bool value of the key
func (mc *MapConfig) GetBool(key string) bool {
	return mc.store[key].(bool)
}

// GetInt returns the int value of the key
func (mc *MapConfig) GetInt(key string) int {
	return mc.store[key].(int)
}

// GetString returns the string value of the key
func (mc *MapConfig) GetString(key string) string {
	return mc.store[key].(string)
}

// GetFloat64 returns the float64 value of the key
func (mc *MapConfig) GetFloat64(key string) float64 {
	return mc.store[key].(float64)
}
