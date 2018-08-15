package zerolog

import (
	"fmt"

	"github.com/mantzas/patron/log"
	"github.com/rs/zerolog"
)

var levelMap = map[log.Level]zerolog.Level{
	log.NoLevel:    zerolog.NoLevel,
	log.DebugLevel: zerolog.DebugLevel,
	log.InfoLevel:  zerolog.InfoLevel,
	log.WarnLevel:  zerolog.WarnLevel,
	log.ErrorLevel: zerolog.ErrorLevel,
	log.FatalLevel: zerolog.FatalLevel,
	log.PanicLevel: zerolog.PanicLevel,
}

// Logger abstraction based on zerolog.
type Logger struct {
	logger *zerolog.Logger
}

// NewLogger creates a new logger.
func NewLogger(l *zerolog.Logger, lvl log.Level, f map[string]interface{}) log.Logger {
	if len(f) == 0 {
		f = make(map[string]interface{})
	}
	zl := l.Level(levelMap[lvl]).With().Fields(f).Logger()
	return &Logger{logger: &zl}
}

// Panic logging.
func (zl Logger) Panic(args ...interface{}) {
	zl.logger.Panic().Msg(fmt.Sprint(args...))
}

// Panicf logging.
func (zl Logger) Panicf(msg string, args ...interface{}) {
	zl.logger.Panic().Msgf(msg, args...)
}

// Fatal logging.
func (zl Logger) Fatal(args ...interface{}) {
	zl.logger.Fatal().Msg(fmt.Sprint(args...))
}

// Fatalf logging.
func (zl Logger) Fatalf(msg string, args ...interface{}) {
	zl.logger.Fatal().Msgf(msg, args...)
}

// Error logging.
func (zl Logger) Error(args ...interface{}) {
	zl.logger.Error().Msg(fmt.Sprint(args...))
}

// Errorf logging.
func (zl Logger) Errorf(msg string, args ...interface{}) {
	zl.logger.Error().Msgf(msg, args...)
}

// Warn logging.
func (zl Logger) Warn(args ...interface{}) {
	zl.logger.Warn().Msg(fmt.Sprint(args...))
}

// Warnf logging.
func (zl Logger) Warnf(msg string, args ...interface{}) {
	zl.logger.Warn().Msgf(msg, args...)
}

// Info logging.
func (zl Logger) Info(args ...interface{}) {
	zl.logger.Info().Msg(fmt.Sprint(args...))
}

// Infof logging.
func (zl Logger) Infof(msg string, args ...interface{}) {
	zl.logger.Info().Msgf(msg, args...)
}

// Debug logging.
func (zl Logger) Debug(args ...interface{}) {
	zl.logger.Debug().Msg(fmt.Sprint(args...))
}

// Debugf logging.
func (zl Logger) Debugf(msg string, args ...interface{}) {
	zl.logger.Debug().Msgf(msg, args...)
}
