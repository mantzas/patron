package zero

import (
	"bytes"
	"os"
	"testing"

	"github.com/bouk/monkey"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name           string
		f              map[string]interface{}
		expectedFields int
	}{
		{"without fields", nil, 0},
		{"with fields", map[string]interface{}{"key": "value"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			l := NewLogger(&zerolog.Logger{}, tt.f)
			assert.NotNil(l)
			assert.Len(l.Fields(), tt.expectedFields)
		})
	}
}

func TestLogger_Fields(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	l := NewLogger(&zerolog.Logger{}, f)
	assert.NotNil(l)
	assert.Equal(f, l.Fields())
}

func TestLogger_Panic(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	assert.Panics(func() { l.Panic("testing") })
	assert.Equal("{\"level\":\"panic\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Panicf(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)

	assert.Panics(func() { l.Panicf("testing %d", 1) })
	assert.Equal("{\"level\":\"panic\",\"message\":\"testing 1\"}\n", b.String())
}

func TestLogger_Fatal(t *testing.T) {
	fakeExit := func(int) {
		panic("os.Exit called")
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	assert.PanicsWithValue("os.Exit called", func() { l.Fatal("testing") })
	assert.Equal("{\"level\":\"fatal\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Fatalf(t *testing.T) {
	fakeExit := func(int) {
		panic("os.Exit called")
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	assert.PanicsWithValue("os.Exit called", func() { l.Fatalf("testing %d", 1) })
	assert.Equal("{\"level\":\"fatal\",\"message\":\"testing 1\"}\n", b.String())
}

func TestLogger_Error(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	l.Error("testing")
	assert.Equal("{\"level\":\"error\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Errorf(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	l.Errorf("testing %d", 1)
	assert.Equal("{\"level\":\"error\",\"message\":\"testing 1\"}\n", b.String())
}

func TestLogger_Warn(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	l.Warn("testing")
	assert.Equal("{\"level\":\"warn\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Warnf(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	l.Warnf("testing %d", 1)
	assert.Equal("{\"level\":\"warn\",\"message\":\"testing 1\"}\n", b.String())
}

func TestLogger_Info(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	l.Info("testing")
	assert.Equal("{\"level\":\"info\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Infof(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	l.Infof("testing %d", 1)
	assert.Equal("{\"level\":\"info\",\"message\":\"testing 1\"}\n", b.String())
}

func TestLogger_Debug(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	l.Debug("testing")
	assert.Equal("{\"level\":\"debug\",\"message\":\"testing\"}\n", b.String())
}

func TestLogger_Debugf(t *testing.T) {
	assert := assert.New(t)
	f := map[string]interface{}{"key": "value"}
	var b bytes.Buffer
	zl := zerolog.New(&b)
	l := NewLogger(&zl, f)
	l.Debugf("testing %d", 1)
	assert.Equal("{\"level\":\"debug\",\"message\":\"testing 1\"}\n", b.String())
}
