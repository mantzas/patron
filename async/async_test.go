package async

import (
	"testing"

	"github.com/mantzas/patron/encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestDetermineDecoder(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		contentType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{json.ContentType}, false},
		{"failure", args{"XXX"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetermineDecoder(tt.args.contentType)
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
