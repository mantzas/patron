package patron

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		name    string
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		want    *Service
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, testCreateHandler, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ListenAndServer_Shutdown(t *testing.T) {
	assert := assert.New(t)
	s, err := New("test", testCreateHandler)
	assert.NoError(err)
	go func() {
		s.Run()
	}()
	err = s.shutdown()
	assert.NoError(err)
}

func Test_createHTTPServer(t *testing.T) {
	assert := assert.New(t)
	s := createHTTPServer(10000, http.DefaultServeMux)
	assert.Equal(":10000", s.Addr)
	assert.Equal(5*time.Second, s.ReadTimeout)
	assert.Equal(60*time.Second, s.WriteTimeout)
	assert.Equal(120*time.Second, s.IdleTimeout)
	assert.Equal(s.Handler, http.DefaultServeMux)
}

func testCreateHandler(routes []Route) http.Handler {
	return http.NewServeMux()
}
