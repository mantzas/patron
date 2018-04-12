package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		f       Factory
		wantErr bool
	}{
		{"failure with nil factory", nil, true},
		{"success", testFactory{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := Setup(tt.f)

			if tt.wantErr {
				assert.Error(err, "expected error")
			} else {
				assert.NoError(err, "error not expected")
			}
		})
	}
}

func TestSub(t *testing.T) {
	assert := assert.New(t)
	Setup(testFactory{})
	l := Sub(make(map[string]interface{}))
	assert.NotNil(l)
}

func TestLog_Panic(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, InfoLevel))
	AppendField("test", "testing")
	Panic("panic")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] PANIC test=testing panic", w.String())
}

func TestLog_Panicf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, InfoLevel))
	AppendField("test", "testing")
	Panicf("panic %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] PANIC test=testing panic 1", w.String())
}

func TestLog_Fatal(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, InfoLevel))
	AppendField("test", "testing")
	Fatal("fatal")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] FATAL test=testing fatal", w.String())
}

func TestLog_Fatalf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, InfoLevel))
	AppendField("test", "testing")
	Fatalf("fatal %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] FATAL test=testing fatal 1", w.String())
}

func TestLog_Error(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, InfoLevel))
	AppendField("test", "testing")
	Error("error")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] ERROR test=testing error", w.String())
}

func TestLog_Errorf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, InfoLevel))
	AppendField("test", "testing")
	Errorf("error %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] ERROR test=testing error 1", w.String())
}

func TestLog_Warn(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, InfoLevel))
	AppendField("test", "testing")
	Warn("warn")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] WARN test=testing warn", w.String())
}

func TestLog_Warnf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, InfoLevel))
	AppendField("test", "testing")
	Warnf("warn %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] WARN test=testing warn 1", w.String())
}

func TestLog_Info(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, InfoLevel))
	AppendField("test", "testing")
	Info("info")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] INFO test=testing info", w.String())
}

func TestLog_Infof(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, InfoLevel))
	AppendField("test", "testing")
	Infof("info %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] INFO test=testing info 1", w.String())
}

func TestLog_Debug(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, DebugLevel))
	AppendField("test", "testing")
	Debug("debug")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] DEBUG test=testing debug", w.String())
}

func TestLog_Debugf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	Setup(NewStdFactory(w, DebugLevel))
	AppendField("test", "testing")
	Debugf("debug %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] DEBUG test=testing debug 1", w.String())
}

type testFactory struct {
}

func (f testFactory) Create(fields map[string]interface{}) Logger {
	return &testLogger{}
}

func (f testFactory) CreateSub(logger Logger, fields map[string]interface{}) Logger {
	return &testLogger{}
}

type testLogger struct {
}

func (t testLogger) Level() Level {
	return DebugLevel
}

func (t testLogger) Fields() map[string]interface{} {
	return make(map[string]interface{})
}

func (t testLogger) Panic(args ...interface{}) {
}

func (t testLogger) Panicf(msg string, args ...interface{}) {
}

func (t testLogger) Fatal(args ...interface{}) {
}

func (t testLogger) Fatalf(msg string, args ...interface{}) {
}

func (t testLogger) Error(args ...interface{}) {
}

func (t testLogger) Errorf(msg string, args ...interface{}) {
}

func (t testLogger) Warn(args ...interface{}) {
}

func (t testLogger) Warnf(msg string, args ...interface{}) {
}

func (t testLogger) Info(args ...interface{}) {
}

func (t testLogger) Infof(msg string, args ...interface{}) {
}

func (t testLogger) Debug(args ...interface{}) {
}

func (t testLogger) Debugf(msg string, args ...interface{}) {
}
