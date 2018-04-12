package zerolog

import (
	"bytes"
	"os"
	"testing"

	"github.com/bouk/monkey"
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
			assert.Equal(l.Level(), tt.lvl)
			assert.Len(l.Fields(), tt.expectedFields)
		})
	}
}

func TestLogger_Fields(t *testing.T) {
	assert := assert.New(t)
	l := NewLogger(&zerolog.Logger{}, log.DebugLevel, f)
	assert.NotNil(l)
	assert.Equal(f, l.Fields())
}

func TestLogger_Panic(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	assert.Panics(func() { l.Panic("testing") })
	assert.Equal("{\"level\":\"panic\",\"key\":\"value\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Panicf(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	assert.Panics(func() { l.Panicf("testing %d", 1) })
	assert.Equal("{\"level\":\"panic\",\"key\":\"value\",\"message\":\"testing 1\"}\n", b.String())
}

func TestLogger_Fatal(t *testing.T) {
	fakeExit := func(int) {
		panic("os.Exit called")
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	assert.PanicsWithValue("os.Exit called", func() { l.Fatal("testing") })
	assert.Equal("{\"level\":\"fatal\",\"key\":\"value\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Fatalf(t *testing.T) {
	fakeExit := func(int) {
		panic("os.Exit called")
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	assert.PanicsWithValue("os.Exit called", func() { l.Fatalf("testing %d", 1) })
	assert.Equal("{\"level\":\"fatal\",\"key\":\"value\",\"message\":\"testing 1\"}\n", b.String())
}

func TestLogger_Error(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Error("testing")
	assert.Equal("{\"level\":\"error\",\"key\":\"value\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Errorf(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Errorf("testing %d", 1)
	assert.Equal("{\"level\":\"error\",\"key\":\"value\",\"message\":\"testing 1\"}\n", b.String())
}

func TestLogger_Warn(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Warn("testing")
	assert.Equal("{\"level\":\"warn\",\"key\":\"value\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Warnf(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Warnf("testing %d", 1)
	assert.Equal("{\"level\":\"warn\",\"key\":\"value\",\"message\":\"testing 1\"}\n", b.String())
}

func TestLogger_Info(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Info("testing")
	assert.Equal("{\"level\":\"info\",\"key\":\"value\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Infof(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Infof("testing %d", 1)
	assert.Equal("{\"level\":\"info\",\"key\":\"value\",\"message\":\"testing 1\"}\n", b.String())
}

func TestLogger_Debug(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Debug("testing")
	assert.Equal("{\"level\":\"debug\",\"key\":\"value\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Debugf(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, log.DebugLevel, f)
	l.Debugf("testing %d", 1)
	assert.Equal("{\"level\":\"debug\",\"key\":\"value\",\"message\":\"testing 1\"}\n", b.String())
}
