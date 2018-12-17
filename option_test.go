package patron

import (
	"testing"

	"github.com/mantzas/patron/sync/http"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
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
		{"success", args{rr: []http.Route{http.NewRoute("/", "GET", nil, true, nil)}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New("test", "1.0.0")
			assert.NoError(t, err)
			err = Routes(tt.args.rr)(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
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
			s, err := New("test", "1.0.0")
			assert.NoError(t, err)
			err = HealthCheck(tt.args.hcf)(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestComponents(t *testing.T) {
	type args struct {
		c Component
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"failure due to empty components", args{}, true},
		{"failure due to nil components", args{c: nil}, true},
		{"success", args{c: &testComponent{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New("test", "1.0.0")
			assert.NoError(t, err)
			err = Components(tt.args.c)(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDocs(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{file: "testdata/test.md"}, wantErr: false},
		{name: "doc file missing", args: args{file: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New("test", "1.0.0")
			assert.NoError(t, err)
			err = Docs(tt.args.file)(s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
