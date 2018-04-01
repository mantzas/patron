package worker

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type testProcessor struct {
	count int
}

func (tp *testProcessor) Process() error {

	if tp.count == 1 {
		return errors.New("failed to process")
	}

	tp.count++
	return nil
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		name string
		p    Processor
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"test", &testProcessor{}}, false},
		{"failed with missing name", args{"", &testProcessor{}}, true},
		{"failed with nil processor", args{"test", nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.name, tt.args.p)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestService_Run(t *testing.T) {
	assert := assert.New(t)
	s, err := New("test", &testProcessor{})
	assert.NoError(err)
	assert.NoError(s.Run())
}
