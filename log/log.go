package log

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
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
	Level() Level
	Fields() map[string]interface{}
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

// Factory interface for creating loggers.
type Factory interface {
	Create(map[string]interface{}) Logger
}

// MsgfFunc defines a logging function type with formatting.
type MsgfFunc func(msg string, args ...interface{})

// MsgFunc defines a logging function type.
type MsgFunc func(msg string)

var (
	// NilMsg instance of a nil logging function.
	NilMsg MsgFunc = func(msg string) {}
	// NilMsgf instance of a nil logging function with formating.
	NilMsgf MsgfFunc = func(msg string, args ...interface{}) {}
	factory Factory
	fields  = make(map[string]interface{})
)

// Setup logging by providing a logger factory.
func Setup(f Factory, fls map[string]interface{}) error {
	if f == nil {
		return errors.New("factory is nil")
	}

	if fls == nil {
		fls = make(map[string]interface{})
	}

	factory = f
	fields = fls
	return nil
}

// Create returns a new logger with all fields inherited and with source file mapping.
func Create() Logger {
	if factory == nil {
		return nil
	}
	fields = make(map[string]interface{})

	if key, val, ok := sourceFields(); ok {
		fields[key] = val
	}

	return factory.Create(fields)
}

func sourceFields() (key string, src string, ok bool) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}

	src = getSource(file, line)
	if src == "" {
		return
	}
	key = "src"
	ok = true
	return
}

func getSource(file string, line int) (src string) {
	if file == "" {
		return
	}
	d, f := filepath.Split(file)
	d = path.Base(d)
	if d == "." || d == "" {
		src = fmt.Sprintf("%s:%d", f, line)
	} else {
		src = fmt.Sprintf("%s/%s:%d", d, f, line)
	}
	return
}
