package patron

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func ErrorOption() Option {
	return func(s *Service) error {
		return errors.New("TEST")
	}
}

func TestNewServer(t *testing.T) {
	assert := assert.New(t)

	cps := []Component{&testComponent{}}
	options := []Option{}

	type args struct {
		name    string
		cps     []Component
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{name: "test", cps: cps, options: options}, false},
		{"failed missing name", args{name: "", cps: cps, options: options}, true},
		{"failed missing components", args{name: "test", cps: []Component{}, options: options}, true},
		{"failed error option", args{name: "test", cps: cps, options: []Option{ErrorOption()}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.cps, tt.args.options...)
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
		cps             []Component
		wantRunErr      bool
		wantShutdownErr bool
	}{
		{"success", []Component{&testComponent{}}, false, false},
		{"failed to run", []Component{&testComponent{errorRunning: true, errorShuttingDown: false}}, true, false},
		{"failed to shutdown", []Component{&testComponent{errorRunning: false, errorShuttingDown: true}}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s, err := New("test", tt.cps)
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

type testComponent struct {
	errorRunning      bool
	errorShuttingDown bool
}

func (ts testComponent) Run(ctx context.Context) error {
	if ts.errorRunning {
		return errors.New("failed to run component")
	}
	return nil
}

func (ts testComponent) Shutdown(ctx context.Context) error {
	if ts.errorShuttingDown {
		return errors.New("failed to shut down")
	}
	return nil
}
