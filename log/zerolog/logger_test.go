package zerolog

import (
	"bytes"
	"testing"

	"github.com/mantzas/patron/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var f = map[string]interface{}{"key": "value"}

func TestNewLogger(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name           string
		f              map[string]interface{}
		lvl            log.Level
		expectedFields int
	}{
		{"without fields", nil, log.DebugLevel, 0},
		{"with fields", f, log.DebugLevel, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			l := NewLogger(&zerolog.Logger{}, tt.lvl, tt.f)
			assert.NotNil(l)
		})
	}
}

func TestLogger_Panic(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	assert.Panics(func() { l.Panic("testing") })
	assert.Equal("{\"lvl\":\"panic\",\"key\":\"value\",\"msg\":\"testing\"}\n", b.String())
}

func TestLogger_Panicf(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	assert.Panics(func() { l.Panicf("testing %d", 1) })
	assert.Equal("{\"lvl\":\"panic\",\"key\":\"value\",\"msg\":\"testing 1\"}\n", b.String())
}

func TestLogger_Error(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Error("testing")
	assert.Equal("{\"lvl\":\"error\",\"key\":\"value\",\"msg\":\"testing\"}\n", b.String())
}

func TestLogger_Errorf(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Errorf("testing %d", 1)
	assert.Equal("{\"lvl\":\"error\",\"key\":\"value\",\"msg\":\"testing 1\"}\n", b.String())
}

func TestLogger_Warn(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Warn("testing")
	assert.Equal("{\"lvl\":\"warn\",\"key\":\"value\",\"msg\":\"testing\"}\n", b.String())
}

func TestLogger_Warnf(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Warnf("testing %d", 1)
	assert.Equal("{\"lvl\":\"warn\",\"key\":\"value\",\"msg\":\"testing 1\"}\n", b.String())
}

func TestLogger_Info(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Info("testing")
	assert.Equal("{\"lvl\":\"info\",\"key\":\"value\",\"msg\":\"testing\"}\n", b.String())
}

func TestLogger_Infof(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Infof("testing %d", 1)
	assert.Equal("{\"lvl\":\"info\",\"key\":\"value\",\"msg\":\"testing 1\"}\n", b.String())
}

func TestLogger_Debug(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Debug("testing")
	assert.Equal("{\"lvl\":\"debug\",\"key\":\"value\",\"msg\":\"testing\"}\n", b.String())
}

func TestLogger_Debugf(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Debugf("testing %d", 1)
	assert.Equal("{\"lvl\":\"debug\",\"key\":\"value\",\"msg\":\"testing 1\"}\n", b.String())
}

var t int

func Benchmark_LoggingEnabled(b *testing.B) {

	var bf bytes.Buffer
	zl := zerolog.New(&bf)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Debugf("testing %d", 1)

	for n := 0; n < b.N; n++ {
		l.Debugf("testing %d", 1)
		t = n
	}
}

func Benchmark_LoggingDisabled(b *testing.B) {

	var bf bytes.Buffer
	zl := zerolog.New(&bf)
	l := NewLogger(&zl, log.NoLevel, f)
	l.Debugf("testing %d", 1)

	for n := 0; n < b.N; n++ {
		l.Debugf("testing %d", 1)
		t = n
	}
}
