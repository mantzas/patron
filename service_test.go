package patron

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

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
		{"success", args{"test", cps, options}, false},
		{"failed missing name", args{"", cps, options}, true},
		{"failed missing components", args{"test", []Component{}, options}, true},
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
		{"failed to run", []Component{&testComponent{true, false}}, true, false},
		{"failed to shutdown", []Component{&testComponent{false, true}}, false, true},
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
	errorRunning     bool
	errorShutingDown bool
}

func (ts testComponent) Run(ctx context.Context) error {
	if ts.errorRunning {
		return errors.New("failed to run component")
	}
	return nil
}

func (ts testComponent) Shutdown(ctx context.Context) error {
	if ts.errorShutingDown {
		return errors.New("failed to shut down")
	}
	return nil
}
