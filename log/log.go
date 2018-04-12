package log

import (
	"errors"
)

var factory Factory
var logger Logger
var fields = make(map[string]interface{})

// Setup set's up a new factory to the global state
func Setup(f Factory) error {
	if f == nil {
		return errors.New("factory is nil")
	}
	factory = f
	logger = f.Create(fields)
	return nil
}

// AppendField appends a field
func AppendField(key string, value interface{}) {
	fields[key] = value
	logger = factory.Create(fields)
}

// Sub returns a new sub logger
func Sub(fields map[string]interface{}) Logger {
	return factory.CreateSub(logger, fields)
}

// Panic logging
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Panicf logging
func Panicf(msg string, args ...interface{}) {
	logger.Panicf(msg, args...)
}

// Fatal logging
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf logging
func Fatalf(msg string, args ...interface{}) {
	logger.Fatalf(msg, args...)
}

// Error logging
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf logging
func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

// Warn logging
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf logging
func Warnf(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

// Info logging
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof logging
func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

// Debug logging
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf logging
func Debugf(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}
