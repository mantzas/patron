package zerolog

import (
	"os"

	"github.com/mantzas/patron/log"
	"github.com/rs/zerolog"
)

// Factory of the zero logger
type Factory struct {
	logger *zerolog.Logger
	lvl    log.Level
}

// NewFactory returns a new zero logger factory
func NewFactory(l *zerolog.Logger, lvl log.Level) log.Factory {
	return &Factory{logger: l, lvl: lvl}
}

// DefaultFactory returns a zero logger factory with default settings
func DefaultFactory(lvl log.Level) log.Factory {
	zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return NewFactory(&zl, lvl)
}

// Create a zero logger
func (zf *Factory) Create(f map[string]interface{}) log.Logger {
	return NewLogger(zf.logger, zf.lvl, f)
}

// CreateSub a zero sub logger with defined fields
func (zf *Factory) CreateSub(logger log.Logger, fields map[string]interface{}) log.Logger {

	if len(fields) == 0 {
		return logger
	}

	all := logger.Fields()

	for k, v := range fields {
		all[k] = v
	}

	l := zf.logger.With().Fields(fields).Logger()
	return NewLogger(&l, zf.lvl, all)
}
