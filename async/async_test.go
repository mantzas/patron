package async

import (
	"testing"

	"github.com/mantzas/patron/encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(NewMessage(nil, nil))
}

func TestMessage_Decode(t *testing.T) {
	assert := assert.New(t)

	j, err := json.Encode("string")
	assert.NoError(err)

	req := NewMessage(j, json.DecodeRaw)
	assert.NotNil(req)

	var data string

	err = req.Decode(&data)
	assert.NoError(err)

	assert.Equal("string", data)
}

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
		{"success", args{contentType: json.ContentType}, false},
		{"failure", args{contentType: "XXX"}, true},
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
