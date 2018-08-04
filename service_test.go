package patron

import (
	"context"
	"testing"

	"github.com/mantzas/patron/sync/http"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	assert := assert.New(t)
	route := http.NewRoute("/", "GET", nil, true)
	goodCfg := &Config{Name: "test"}
	badCfg := &Config{}
	type args struct {
		cfg *Config
		opt OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{cfg: goodCfg, opt: Routes([]http.Route{route})}, false},
		{"failed missing config", args{cfg: nil, opt: Routes([]http.Route{route})}, true},
		{"failed missing name", args{cfg: badCfg, opt: Routes([]http.Route{route})}, true},
		{"failed missing routes", args{cfg: goodCfg, opt: Routes([]http.Route{})}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg, tt.args.opt)
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
		cp              Component
		wantRunErr      bool
		wantShutdownErr bool
	}{
		{"success", &testComponent{}, false, false},
		{"failed to run", &testComponent{errorRunning: true, errorShuttingDown: false}, true, false},
		{"failed to shutdown", &testComponent{errorRunning: false, errorShuttingDown: true}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(&Config{Name: "test", Version: ""}, Components(tt.cp))
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
