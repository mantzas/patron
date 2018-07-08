package env

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
)

// Config implementation for handling environment vars.
type Config struct {
}

// New creates a new config.
// By providing a reader, which might contain environment variables coming for a file, you can set env vars.
// This is useful for development.
func New(r io.Reader) (*Config, error) {
	err := initialize(r)
	if err != nil {
		return nil, err
	}
	return &Config{}, nil
}

// Set a environment var.
func (c Config) Set(key string, value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errors.New("failed to type assert value to string")
	}
	if _, ok = os.LookupEnv(key); ok {
		log.Warnf("overwrite existing env var %s", key)
	}
	return errors.Wrap(os.Setenv(key, v), "failed to set env var")
}

// Get a env var.
func (c Config) Get(key string) (interface{}, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return nil, fmt.Errorf("failed to find env var with name %s", key)
	}
	return v, nil
}

// GetBool returns a bool env var.
func (c Config) GetBool(key string) (bool, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return false, fmt.Errorf("failed to find env var with name %s", key)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, errors.Wrapf(err, "env var with key %s is not of type boolean", key)
	}
	return b, nil
}

// GetInt64 returns a int64 env var.
func (c Config) GetInt64(key string) (int64, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return 0, fmt.Errorf("failed to find env var with name %s", key)
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "env var with key %s is not of type int", key)
	}
	return i, nil
}

// GetString returns a string env var.
func (c Config) GetString(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("failed to find env var with name %s", key)
	}
	return v, nil
}

// GetFloat64 returns a float64 env var.
func (c Config) GetFloat64(key string) (float64, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return 0.0, fmt.Errorf("failed to find env var with name %s", key)
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0.0, errors.Wrapf(err, "env var with key %s is not of type float64", key)
	}
	return f, nil
}

func initialize(r io.Reader) error {
	if r == nil {
		return nil
	}

	var vars = make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) != 2 {
			return errors.New("line should contain key and value separated by a equal symbol")
		}
		vars[line[0]] = line[1]
	}

	for k, v := range vars {
		if _, ok := os.LookupEnv(k); ok {
			log.Warnf("env var %s is already defined, skipping", k)
			continue
		}
		log.Infof("setting env var %s", k)
		err := os.Setenv(k, v)
		if err != nil {
			return errors.Wrapf(err, "failed to set env var %s", k)
		}
	}

	return nil
}
