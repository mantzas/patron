package zerolog

import (
	"os"

	"github.com/mantzas/patron/log"
	"github.com/rs/zerolog"
)

// Factory implementation of zerolog.
type Factory struct {
	logger *zerolog.Logger
	lvl    log.Level
}

// NewFactory creates a new zerolog factory.
func NewFactory(l *zerolog.Logger, lvl log.Level) log.Factory {
	return &Factory{logger: l, lvl: lvl}
}

// DefaultFactory creates a zerolog factory with default settings.
func DefaultFactory(lvl log.Level) log.Factory {
	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"
	zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return NewFactory(&zl, lvl)
}

// Create a new logger.
func (zf *Factory) Create(f map[string]interface{}) log.Logger {
	return NewLogger(zf.logger, zf.lvl, f)
}

// CreateSub creates a logger with inherited fields.
func (zf *Factory) CreateSub(logger log.Logger, fields map[string]interface{}) log.Logger {

	if len(fields) == 0 {
		return logger
	}

	all := logger.Fields()

	for k, v := range fields {
		all[k] = v
	}

	return NewLogger(zf.logger, zf.lvl, all)
}
