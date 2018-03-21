package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStdLogger_WithoutFields(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	l := NewStdLogger(w, nil)
	l.Info("test")
	assert.Len(l.Fields(), 0)
}

func TestStdLogger_Fields(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	assert.EqualValues(f, NewStdLogger(w, f).Fields())
}

func TestStdLogger_Panic(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Panic("panic")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] PANIC test=testing panic", w.String())
}

func TestStdLogger_Panicf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Panicf("panic %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] PANIC test=testing panic 1", w.String())
}

func TestStdLogger_Fatal(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Fatal("fatal")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] FATAL test=testing fatal", w.String())
}

func TestStdLogger_Fatalf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Fatalf("fatal %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] FATAL test=testing fatal 1", w.String())
}

func TestStdLogger_Error(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Error("error")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] ERROR test=testing error", w.String())
}

func TestStdLogger_Errorf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Errorf("error %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] ERROR test=testing error 1", w.String())
}

func TestStdLogger_Warn(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Warn("warn")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] WARN test=testing warn", w.String())
}

func TestStdLogger_Warnf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Warnf("warn %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] WARN test=testing warn 1", w.String())
}

func TestStdLogger_Info(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Info("info")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] INFO test=testing info", w.String())
}

func TestStdLogger_Infof(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Infof("info %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] INFO test=testing info 1", w.String())
}

func TestStdLogger_Debug(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Debug("debug")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] DEBUG test=testing debug", w.String())
}

func TestStdLogger_Debugf(t *testing.T) {
	assert := assert.New(t)
	w := &bytes.Buffer{}
	f := make(map[string]interface{})
	f["test"] = "testing"
	NewStdLogger(w, f).Debugf("debug %s", "1")
	assert.Regexp("\\[\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d+Z\\] DEBUG test=testing debug 1", w.String())
}
