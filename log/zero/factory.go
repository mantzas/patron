package zero

import (
	"os"

	"github.com/mantzas/patron/log"
	"github.com/rs/zerolog"
)

// Factory of the zero logger
type Factory struct {
	l *zerolog.Logger
}

// NewFactory returns a new zero logger factory
func NewFactory(l *zerolog.Logger) log.Factory {
	return &Factory{l}
}

// DefaultFactory returns a zero logger factory with default settings
func DefaultFactory(l zerolog.Level) log.Factory {
	zerolog.SetGlobalLevel(l)
	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"
	zerolog.TimestampFieldName = "ts"
	zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return NewFactory(&zl)
}

// Create a zero logger
func (zf *Factory) Create(f map[string]interface{}) log.Logger {
	return NewLogger(zf.l, f)
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

	l := zf.l.With().Fields(fields).Logger()
	return NewLogger(&l, all)
}
