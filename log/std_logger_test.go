package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var f = map[string]interface{}{"test": "testing"}

func TestNewStdLogger_WithoutFields(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	l := NewStdLogger(w, InfoLevel, nil)
	l.Info("test")
	assert.Len(l.Fields(), 0)
}

func TestStdLogger_Fields(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	assert.EqualValues(f, NewStdLogger(w, InfoLevel, f).Fields())
}

func TestStdLogger_Panic(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, InfoLevel, f).Panic("panic")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] PANIC test=testing panic", w.String())
}

func TestStdLogger_Panicf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, InfoLevel, f).Panicf("panic %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] PANIC test=testing panic 1", w.String())
}

func TestStdLogger_Fatal(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, InfoLevel, f).Fatal("fatal")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] FATAL test=testing fatal", w.String())
}

func TestStdLogger_Fatalf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, InfoLevel, f).Fatalf("fatal %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] FATAL test=testing fatal 1", w.String())
}

func TestStdLogger_Error(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, InfoLevel, f).Error("error")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] ERROR test=testing error", w.String())
}

func TestStdLogger_Errorf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, InfoLevel, f).Errorf("error %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] ERROR test=testing error 1", w.String())
}

func TestStdLogger_Warn(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, InfoLevel, f).Warn("warn")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] WARN test=testing warn", w.String())
}

func TestStdLogger_Warnf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, InfoLevel, f).Warnf("warn %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] WARN test=testing warn 1", w.String())
}

func TestStdLogger_Info(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, InfoLevel, f).Info("info")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] INFO test=testing info", w.String())
}

func TestStdLogger_Infof(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, InfoLevel, f).Infof("info %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] INFO test=testing info 1", w.String())
}

func TestStdLogger_Debug(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, DebugLevel, f).Debug("debug")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] DEBUG test=testing debug", w.String())
}

func TestStdLogger_Debugf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, DebugLevel, f).Debugf("debug %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] DEBUG test=testing debug 1", w.String())
}

func TestStdLogger_Debug_NoLog(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, NoLevel, f).Debug("debug")
	assert.Equal("", w.String())
}

func TestStdLogger_Debugf_NoLog(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	NewStdLogger(w, NoLevel, f).Debugf("debug %s", "1")
	assert.Equal("", w.String())
}
