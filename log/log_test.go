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

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	Setup(testFactory{})
	l := Create()
	assert.NotNil(l)
}

func TestCreateWithFields(t *testing.T) {
	assert := assert.New(t)
	Setup(testFactory{})
	l := CreateWithFields(make(map[string]interface{}))
	assert.NotNil(l)
}

func TestCreateSub(t *testing.T) {
	assert := assert.New(t)
	Setup(testFactory{})
	l := CreateSub(&testLogger{}, make(map[string]interface{}))
	assert.NotNil(l)
}

type testFactory struct {
}

func (f testFactory) Create() Logger {
	return &testLogger{}
}

func (f testFactory) CreateWithFields(fields map[string]interface{}) Logger {
	return &testLogger{}
}

func (f testFactory) CreateSub(logger Logger, fields map[string]interface{}) Logger {
	return &testLogger{}
}

type testLogger struct {
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
