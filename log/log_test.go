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
		{"success", &nilFactory{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Setup(tt.f, nil)
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
	t.Run("success", func(t *testing.T) {
		fields = map[string]interface{}{"key": "val"}
		expected := map[string]interface{}{"key": "val", "src": "log/log_test.go"}
		l := Create()
		assert.NotNil(l)
		assert.Equal(DebugLevel, l.Level())
		assert.Equal(expected, l.Fields())
	})
	t.Run("factory nil", func(t *testing.T) {
		factory = nil
		l := Create()
		assert.Nil(l)
	})
}

func Test_getSource(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		file    string
		wantSrc string
	}{
		{name: "empty", file: "", wantSrc: ""},
		{name: "no parent folder", file: "main.go", wantSrc: "main.go"},
		{name: "with parent folder", file: "/home/patron/main.go", wantSrc: "patron/main.go"},
	}
	for _, tt := range tests {
		assert.Equal(tt.wantSrc, getSource(tt.file))
	}
}
