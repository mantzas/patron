package zerolog

import (
	"os"
	"time"

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
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zl := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	return NewFactory(&zl, lvl)
}

// Create a new logger.
func (zf *Factory) Create(f map[string]interface{}) log.Logger {
	return NewLogger(zf.logger, zf.lvl, f)
}
