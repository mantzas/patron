package httprouter

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mantzas/patron/http/route"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		routes []route.Route
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"failed with nil routes", args{nil}, true},
		{"failed with nil routes", args{[]route.Route{}}, true},
		{"success", args{[]route.Route{route.New("/", "GET", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Welcome!\n")
		})}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateHandler(tt.args.routes)
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
