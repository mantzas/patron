package log

import (
	"errors"
)

var factory Factory
var logger Logger
var fields = make(map[string]interface{})

// Setup logging by providing a logger factory.
func Setup(f Factory) error {
	if f == nil {
		return errors.New("factory is nil")
	}
	factory = f
	logger = f.Create(fields)
	return nil
}

// AppendField appends a field to the global logger.
func AppendField(key string, value interface{}) {
	if factory == nil {
		return
	}
	fields[key] = value
	logger = factory.Create(fields)
}

// Sub returns a new sub logger with all fields inherited.
func Sub(fields map[string]interface{}) Logger {
	if factory == nil || logger == nil {
		return nil
	}
	return factory.CreateSub(logger, fields)
}

// Panic logging.
func Panic(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Panic(args...)
}

// Panicf logging.
func Panicf(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Panicf(msg, args...)
}

// Fatal logging.
func Fatal(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Fatal(args...)
}

// Fatalf logging.
func Fatalf(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Fatalf(msg, args...)
}

// Error logging.
func Error(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Error(args...)
}

// Errorf logging.
func Errorf(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Errorf(msg, args...)
}

// Warn logging.
func Warn(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Warn(args...)
}

// Warnf logging.
func Warnf(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Warnf(msg, args...)
}

// Info logging.
func Info(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Info(args...)
}

// Infof logging.
func Infof(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Infof(msg, args...)
}

// Debug logging.
func Debug(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Debug(args...)
}

// Debugf logging.
func Debugf(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Debugf(msg, args...)
}
