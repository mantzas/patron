package patron

import (
	"os"

	"github.com/mantzas/patron/config"
	"github.com/mantzas/patron/config/env"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/pkg/errors"
)

// Config defines configuration properties of the framework.
type Config struct {
	Name    string
	Version string
}

// Configure set's up configuration and logging.
func Configure(name, version string) (*Config, error) {

	cfg := Config{Name: name, Version: version}

	if err := setupConfig(); err != nil {
		return nil, err
	}

	if err := setupLogging(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setupLogging(cfg Config) error {

	lvl, err := config.GetString("PATRON_LOG_LEVEL")
	if err != nil {
		lvl = string(log.InfoLevel)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "failed to get hostname")
	}

	f := map[string]interface{}{
		"srv":  cfg.Name,
		"ver":  cfg.Version,
		"host": hostname,
	}

	err = log.Setup(zerolog.DefaultFactory(log.Level(lvl)), f)
	if err != nil {
		return errors.Wrap(err, "failed to setup logging")
	}

	return nil
}

func setupConfig() error {
	f, err := os.Open(".env")
	if err != nil {
		f = nil
	}

	cfg, err := env.New(f)
	if err != nil {
		return err
	}

	return config.Setup(cfg)
}
