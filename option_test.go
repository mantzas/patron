package patron

import (
	"testing"

	"github.com/mantzas/patron/sync/http"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		rr []http.Route
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"failure due to empty routes", args{rr: []http.Route{}}, true},
		{"failure due to nil routes", args{rr: nil}, true},
		{"success", args{rr: []http.Route{http.NewRoute("/", "GET", nil, true)}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{name: "test"}
			err := Routes(tt.args.rr)(&s)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		hcf http.HealthCheckFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"failure due to nil health check", args{hcf: nil}, true},
		{"success", args{hcf: http.DefaultHealthCheck}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{name: "test"}
			err := HealthCheck(tt.args.hcf)(&s)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestComponents(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		c Component
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"failure due to empty components", args{}, true},
		{"failure due to nil components", args{nil}, true},
		{"success", args{&testComponent{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{name: "test"}
			err := Components(tt.args.c)(&s)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
