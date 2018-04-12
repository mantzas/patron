package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	panic   = "PANIC"
	fatal   = "FATAL"
	err     = "ERROR"
	warning = "WARN"
	info    = "INFO"
	debug   = "DEBUG"
)

// StdLogger implements logging with std log.
type StdLogger struct {
	w        io.Writer
	fields   map[string]interface{}
	lvlOrder LevelOrder
}

// NewStdLogger returns a new logger. If a writer is not provided, stdout will be used.
func NewStdLogger(w io.Writer, lvl Level, f map[string]interface{}) Logger {

	if f == nil {
		f = make(map[string]interface{})
	}

	if w == nil {
		w = os.Stdout
	}

	return &StdLogger{w, f, LevelToOrder(lvl)}
}

// Level returns the minimum level.
func (l StdLogger) Level() Level {
	return OrderToLevel(l.lvlOrder)
}

// Fields returns the fields associated with this logger.
func (l StdLogger) Fields() map[string]interface{} {
	return l.fields
}

// Panic logging
func (l StdLogger) Panic(args ...interface{}) {
	l.log(PanicLevelOrder, panic, args...)
}

// Panicf logging
func (l StdLogger) Panicf(msg string, args ...interface{}) {
	l.logf(PanicLevelOrder, panic, msg, args...)
}

// Fatal logging
func (l StdLogger) Fatal(args ...interface{}) {
	l.log(FatalLevelOrder, fatal, args...)
}

// Fatalf logging
func (l StdLogger) Fatalf(msg string, args ...interface{}) {
	l.logf(FatalLevelOrder, fatal, msg, args...)
}

// Error logging
func (l StdLogger) Error(args ...interface{}) {
	l.log(ErrorLevelOrder, err, args...)
}

// Errorf logging
func (l StdLogger) Errorf(msg string, args ...interface{}) {
	l.logf(ErrorLevelOrder, err, msg, args...)
}

// Warn logging
func (l StdLogger) Warn(args ...interface{}) {
	l.log(WarnLevelOrder, warning, args...)
}

// Warnf logging
func (l StdLogger) Warnf(msg string, args ...interface{}) {
	l.logf(WarnLevelOrder, warning, msg, args...)
}

// Info logging
func (l StdLogger) Info(args ...interface{}) {
	l.log(InfoLevelOrder, info, args...)
}

// Infof logging
func (l StdLogger) Infof(msg string, args ...interface{}) {
	l.logf(InfoLevelOrder, info, msg, args...)
}

// Debug logging
func (l StdLogger) Debug(args ...interface{}) {
	l.log(DebugLevelOrder, debug, args...)
}

// Debugf logging
func (l StdLogger) Debugf(msg string, args ...interface{}) {
	l.logf(DebugLevelOrder, debug, msg, args...)
}

func (l StdLogger) log(ord LevelOrder, lvl string, args ...interface{}) {
	if !ShouldLog(l.lvlOrder, ord) {
		return
	}
	l.logMessage(lvl, fmt.Sprint(args...))
}

func (l StdLogger) logf(ord LevelOrder, lvl string, msg string, args ...interface{}) {
	if !ShouldLog(l.lvlOrder, ord) {
		return
	}
	l.logMessage(lvl, fmt.Sprintf(msg, args...))
}

func (l StdLogger) logMessage(lvl string, msg string) {
	ts := time.Now().UTC().Format(time.RFC3339Nano)
	fmt.Fprintf(l.w, "[%s] %s %s%s\n", ts, lvl, l.getFieldsMessage(), msg)
}

func (l StdLogger) getFieldsMessage() string {

	if len(l.fields) == 0 {
		return ""
	}

	var m string
	for k, v := range l.fields {
		m += fmt.Sprintf("%s=%s ", k, v)
	}
	return m
}
