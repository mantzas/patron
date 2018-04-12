package log

import (
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
	Setup(&testFactory{})
	l := Sub(make(map[string]interface{}))
	assert.NotNil(l)
}

func TestLog_AppendField(t *testing.T) {
	assert := assert.New(t)
	Setup(&testFactory{})
	AppendField("test", "testing")
	assert.Equal("testing", fields["test"])
}

func TestLog_Panic(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Panic("panic")
	assert.Equal(1, l.panicCount)
}

func TestLog_Panicf(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Panicf("panic %s", "1")
	assert.Equal(1, l.panicCount)
}

func TestLog_Fatal(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Fatal("fatal")
	assert.Equal(1, l.fatalCount)
}

func TestLog_Fatalf(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Fatalf("fatal %s", "1")
	assert.Equal(1, l.fatalCount)
}

func TestLog_Error(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Error("error")
	assert.Equal(1, l.errorCount)
}

func TestLog_Errorf(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Errorf("error %s", "1")
	assert.Equal(1, l.errorCount)
}

func TestLog_Warn(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Warn("warn")
	assert.Equal(1, l.warnCount)
}

func TestLog_Warnf(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Warnf("warn %s", "1")
	assert.Equal(1, l.warnCount)
}

func TestLog_Info(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Info("info")
	assert.Equal(1, l.infoCount)
}

func TestLog_Infof(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Infof("info %s", "1")
	assert.Equal(1, l.infoCount)
}

func TestLog_Debug(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Debug("debug")
	assert.Equal(1, l.debugCount)
}

func TestLog_Debugf(t *testing.T) {
	assert := assert.New(t)
	l := testLogger{}
	logger = &l
	Debugf("debug %s", "1")
	assert.Equal(1, l.debugCount)
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
	debugCount int
	infoCount  int
	warnCount  int
	errorCount int
	fatalCount int
	panicCount int
}

func (t *testLogger) Level() Level {
	return DebugLevel
}

func (t *testLogger) Fields() map[string]interface{} {
	return make(map[string]interface{})
}

func (t *testLogger) Panic(args ...interface{}) {
	t.panicCount++
}

func (t *testLogger) Panicf(msg string, args ...interface{}) {
	t.panicCount++
}

func (t *testLogger) Fatal(args ...interface{}) {
	t.fatalCount++
}

func (t *testLogger) Fatalf(msg string, args ...interface{}) {
	t.fatalCount++
}

func (t *testLogger) Error(args ...interface{}) {
	t.errorCount++
}

func (t *testLogger) Errorf(msg string, args ...interface{}) {
	t.errorCount++
}

func (t *testLogger) Warn(args ...interface{}) {
	t.warnCount++
}

func (t *testLogger) Warnf(msg string, args ...interface{}) {
	t.warnCount++
}

func (t *testLogger) Info(args ...interface{}) {
	t.infoCount++
}

func (t *testLogger) Infof(msg string, args ...interface{}) {
	t.infoCount++
}

func (t *testLogger) Debug(args ...interface{}) {
	t.debugCount++
}

func (t *testLogger) Debugf(msg string, args ...interface{}) {
	t.debugCount++
}
