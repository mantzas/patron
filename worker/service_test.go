package worker

import (
	"testing"

	"github.com/mantzas/patron/worker/work"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type testAcquirer struct {
	returnError bool
}

func (ta testAcquirer) Acquire() ([]work.Work, error) {

	if ta.returnError {
		return nil, errors.New("error returned")
	}

	return []work.Work{"One", "Two"}, nil
}

type testProcessor struct {
	count int
}

func (tp *testProcessor) Process(items []work.Work) error {

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
		acq  work.Acquirer
		prc  work.Processor
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"test", &testAcquirer{}, &testProcessor{}}, false},
		{"failed with missing name", args{"", &testAcquirer{}, &testProcessor{}}, true},
		{"failed with missing acquirer", args{"test", nil, &testProcessor{}}, true},
		{"failed with missing processor", args{"test", &testAcquirer{}, nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.name, tt.args.acq, tt.args.prc)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestService_Run_ProcessorError(t *testing.T) {
	assert := assert.New(t)
	s, err := New("test", &testAcquirer{}, &testProcessor{})
	assert.NoError(err)
	assert.Error(s.Run())
}

func TestService_Run_AcquirerError(t *testing.T) {
	assert := assert.New(t)
	s, err := New("test", &testAcquirer{true}, &testProcessor{})
	assert.NoError(err)
	assert.Error(s.Run())
}
