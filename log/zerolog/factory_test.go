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

func TestFactory_CreateSub(t *testing.T) {
	assert := assert.New(t)
	var b bytes.Buffer
	f := createFactory(&b)
	tests := []struct {
		name           string
		fields         map[string]interface{}
		expectedFields int
	}{
		{"without fields", nil, 0},
		{"with fields", map[string]interface{}{"key": "value"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			l := f.CreateSub(f.Create(nil), tt.fields)

			assert.NotNil(l)
			assert.Len(l.Fields(), tt.expectedFields)
		})
	}
}

func createFactory(wr io.Writer) log.Factory {
	l := zerolog.New(wr)
	return NewFactory(&l, log.DebugLevel)
}
