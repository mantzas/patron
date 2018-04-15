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

func TestServer_Run_ReturnsError(t *testing.T) {
	assert := assert.New(t)
	s, err := NewServer("test", &testService{true, false})
	assert.NoError(err)
	assert.Error(s.Run())
}

func TestServer_Shutdown_ReturnsError(t *testing.T) {
	assert := assert.New(t)
	s, err := NewServer("test", &testService{false, true})
	assert.NoError(err)
	go func() {
		s.Run()
	}()
	assert.Error(s.Shutdown())
}

func TestServer_Shutdown(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		service ServiceInt
		wantErr bool
	}{
		{"success", &testService{}, false},
		{"failed to shutdown", &testService{false, true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s, err := NewServer("test", tt.service)
			assert.NoError(err)
			go func() {
				s.Run()
			}()
			err = s.Shutdown()
			if tt.wantErr {
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
