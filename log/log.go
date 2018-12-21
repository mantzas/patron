package log

import (
	"context"

	"github.com/mantzas/patron/errors"
)

// The Level type definition.
type Level string

const (
	// DebugLevel level.
	DebugLevel Level = "debug"
	// InfoLevel level.
	InfoLevel Level = "info"
	// WarnLevel level.
	WarnLevel Level = "warn"
	// ErrorLevel level.
	ErrorLevel Level = "error"
	// FatalLevel level.
	FatalLevel Level = "fatal"
	// PanicLevel level.
	PanicLevel Level = "panic"
	// NoLevel level.
	NoLevel Level = ""
)

// Logger interface definition of a logger.
type Logger interface {
	Sub(map[string]interface{}) Logger
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Panic(...interface{})
	Panicf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
}

type ctxKey struct{}

// FactoryFunc function type for creating loggers.
type FactoryFunc func(map[string]interface{}) Logger

var logger Logger = &nilLogger{}

// Setup logging by providing a logger factory.
func Setup(f FactoryFunc, fls map[string]interface{}) error {
	if f == nil {
		return errors.New("factory is nil")
	}

	if fls == nil {
		fls = make(map[string]interface{})
	}

	logger = f(fls)
	return nil
}

// FromContext returns the logger in the context or a nil logger.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return l
	}
	return &nilLogger{}
}

// WithContext associates a logger with a context for later reuse.
func WithContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// Sub returns a sub logger with new fields attached.
func Sub(ff map[string]interface{}) Logger {
	return logger.Sub(ff)
}

// Panic logging.
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Panicf logging.
func Panicf(msg string, args ...interface{}) {
	logger.Panicf(msg, args...)
}

// Fatal logging.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf logging.
func Fatalf(msg string, args ...interface{}) {
	logger.Fatalf(msg, args...)
}

// Error logging.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf logging.
func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

// Warn logging.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf logging.
func Warnf(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

// Info logging.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof logging.
func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

// Debug logging.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf logging.
func Debugf(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}

type nilLogger struct {
}

func (nl *nilLogger) Sub(map[string]interface{}) Logger {
	return nl
}

func (nl *nilLogger) Panic(args ...interface{}) {
}

func (nl *nilLogger) Panicf(msg string, args ...interface{}) {
}

func (nl *nilLogger) Fatal(args ...interface{}) {
}

func (nl *nilLogger) Fatalf(msg string, args ...interface{}) {
}

func (nl *nilLogger) Error(args ...interface{}) {
}

func (nl *nilLogger) Errorf(msg string, args ...interface{}) {
}

func (nl *nilLogger) Warn(args ...interface{}) {
}

func (nl *nilLogger) Warnf(msg string, args ...interface{}) {
}

func (nl *nilLogger) Info(args ...interface{}) {
}

func (nl *nilLogger) Infof(msg string, args ...interface{}) {
}

func (nl *nilLogger) Debug(args ...interface{}) {
}

func (nl *nilLogger) Debugf(msg string, args ...interface{}) {
}
