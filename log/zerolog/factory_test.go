package zerolog

import (
	"bytes"
	"io"
	"testing"

	"github.com/mantzas/patron/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewFactory(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	f := createFactory(&b)
	assert.NotNil(f)
}

func TestDefaultFactory(t *testing.T) {
	assert := assert.New(t)
	f := DefaultFactory(log.InfoLevel)
	assert.NotNil(f)
}

func TestFactory_Create(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	f := createFactory(&b)
	l := f.Create(nil)
	assert.NotNil(l)
}

func Test_getSource(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		file    string
		wantSrc string
	}{
		{name: "empty", file: "", wantSrc: ""},
		{name: "no parent folder", file: "main.go", wantSrc: "main.go:10"},
		{name: "with parent folder", file: "/home/patron/main.go", wantSrc: "patron/main.go:10"},
	}
	for _, tt := range tests {
		assert.Equal(tt.wantSrc, getSource(tt.file, 10))
	}
}

func Test_sourceFields(t *testing.T) {
	assert := assert.New(t)
	key, src, ok := sourceFields(1)
	assert.True(ok)
	assert.Equal("src", key)
	assert.Equal("zerolog/factory_test.go:52", src)
}

func createFactory(wr io.Writer) log.Factory {
	l := zerolog.New(wr)
	return NewFactory(&l, log.DebugLevel)
}
