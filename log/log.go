package log

import (
	"errors"
	"sync"
)

var factory Factory
var logger Logger
var fields = make(map[string]interface{})
var m = sync.Mutex{}

// Setup set's up a new factory to the global state
func Setup(f Factory) error {
	if f == nil {
		return errors.New("factory is nil")
	}
	m.Lock()
	defer m.Unlock()
	factory = f
	logger = f.Create(fields)
	return nil
}

// AppendField appends a field
func AppendField(key string, value interface{}) {
	if factory == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	fields[key] = value
	logger = factory.Create(fields)
}

// Sub returns a new sub logger
func Sub(fields map[string]interface{}) Logger {
	if factory == nil || logger == nil {
		return nil
	}
	m.Lock()
	defer m.Unlock()
	return factory.CreateSub(logger, fields)
}

// Panic logging
func Panic(args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Panic(args...)
}

// Panicf logging
func Panicf(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Panicf(msg, args...)
}

// Fatal logging
func Fatal(args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Fatal(args...)
}

// Fatalf logging
func Fatalf(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Fatalf(msg, args...)
}

// Error logging
func Error(args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Error(args...)
}

// Errorf logging
func Errorf(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Errorf(msg, args...)
}

// Warn logging
func Warn(args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Warn(args...)
}

// Warnf logging
func Warnf(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Warnf(msg, args...)
}

// Info logging
func Info(args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Info(args...)
}

// Infof logging
func Infof(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Infof(msg, args...)
}

// Debug logging
func Debug(args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Debug(args...)
}

// Debugf logging
func Debugf(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	logger.Debugf(msg, args...)
}
