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
	assert.Len(l.Fields(), 0)
}

func createFactory(wr io.Writer) log.Factory {
	l := zerolog.New(wr)
	return NewFactory(&l, log.DebugLevel)
}
