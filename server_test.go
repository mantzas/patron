package patron

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		name     string
		services []ServiceInt
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"test", []ServiceInt{&testService{}}}, false},
		{"failed missing name", args{"", []ServiceInt{&testService{}}}, true},
		{"failed missing services", args{"test", []ServiceInt{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := NewServer(tt.args.name, tt.args.services...)
			if tt.wantErr {
				assert.Error(err)
				assert.Nil(got)
			} else {
				assert.NoError(err)
				assert.NotNil(got)
			}
		})
	}
}

func TestServer_Run_Shutdown(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name            string
		service         ServiceInt
		wantRunErr      bool
		wantShutdownErr bool
	}{
		{"success", &testService{}, false, false},
		{"failed to run", &testService{true, false}, true, false},
		{"failed to shutdown", &testService{false, true}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s, err := NewServer("test", tt.service)
			assert.NoError(err)
			err = s.Run()
			if tt.wantRunErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			err = s.Shutdown()
			if tt.wantShutdownErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

type testService struct {
	errorRunning     bool
	errorShutingDown bool
}

func (ts testService) Run(ctx context.Context) error {
	if ts.errorRunning {
		return errors.New("failed to run service")
	}
	return nil
}

func (ts testService) Shutdown(ctx context.Context) error {
	if ts.errorShutingDown {
		return errors.New("failed to shut down")
	}
	return nil
}
