package patron

import (
	"context"
	"testing"

	"github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/sync/http"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	assert := assert.New(t)
	route := http.NewRoute("/", "GET", nil, true)
	type args struct {
		name string
		opt  OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{name: "test", opt: Routes([]http.Route{route})}, false},
		{"failed missing name", args{name: "", opt: Routes([]http.Route{route})}, true},
		{"failed missing routes", args{name: "test", opt: Routes([]http.Route{})}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, "", tt.args.opt)
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
		name       string
		cp         Component
		wantRunErr bool
	}{
		{"success", &testComponent{}, false},
		{"failed to run", &testComponent{errorRunning: true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New("test", "", Components(tt.cp, tt.cp, tt.cp))
			assert.NoError(err)
			err = s.Run()
			if tt.wantRunErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

type testComponent struct {
	errorRunning bool
}

func (ts testComponent) Run(ctx context.Context) error {
	if ts.errorRunning {
		return errors.New("failed to run component")
	}
	return nil
}
