package kafka

import (
	"testing"

	"github.com/mantzas/patron/errors"
	"github.com/stretchr/testify/assert"
)

func ErrorOption() OptionFunc {
	return func(p *Producer) error {
		return errors.New("TEST")
	}
}

func TestNewProducer(t *testing.T) {
	brokers := []string{"xxx"}
	type args struct {
		brokers []string
		oo      []OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		async   bool
		wantErr bool
	}{
		{name: "sync failed, no brokers", args: args{}, wantErr: true},
		{name: "sync failed, invalid option", args: args{brokers: brokers, oo: []OptionFunc{ErrorOption()}}, wantErr: true},
		{name: "sync success", args: args{brokers: brokers}, wantErr: false},
		{name: "async failed, no brokers", async: true, args: args{}, wantErr: true},
		{name: "async failed, invalid option", async: true, args: args{brokers: brokers, oo: []OptionFunc{ErrorOption()}}, wantErr: true},
		{name: "async success", async: true, args: args{brokers: brokers}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var got *Producer
			if tt.async {
				got, err = NewAsyncProducer(tt.args.brokers, tt.args.oo...)
			} else {
				got, err = NewProducer(tt.args.brokers, tt.args.oo...)
			}
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				if tt.async {
					assert.NotNil(t, got.Results())
				}
			}
		})
	}
}
