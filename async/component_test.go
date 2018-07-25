package async

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	proc := mockProcessor{}
	type args struct {
		name string
		p    ProcessorFunc
		cns  Consumer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{name: "test", p: proc.Process, cns: &mockConsumer{}}, wantErr: false},
		{name: "failed, missing name", args: args{name: "", p: proc.Process, cns: &mockConsumer{}}, wantErr: true},
		{name: "failed, missing processor func", args: args{name: "test", p: nil, cns: &mockConsumer{}}, wantErr: true},
		{name: "failed, missing consumer", args: args{name: "test", p: proc.Process, cns: nil}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.p, tt.args.cns)
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

func TestRun_ReturnsError(t *testing.T) {
	assert := assert.New(t)
	cnr := mockConsumer{consumeError: true}
	proc := mockProcessor{}
	cmp, err := New("test", proc.Process, &cnr)
	assert.NoError(err)
	ctx := context.Background()
	err = cmp.Run(ctx)
	assert.Error(err)
}

func TestRun_Process_Error(t *testing.T) {
	assert := assert.New(t)
	cnr := mockConsumer{
		chMsg: make(chan Message, 10),
		chErr: make(chan error, 10),
	}
	proc := mockProcessor{retError: true}
	cmp, err := New("test", proc.Process, &cnr)
	assert.NoError(err)
	ctx := context.Background()
	cnr.chMsg <- &mockMessage{ctx}
	err = cmp.Run(ctx)
	assert.Error(err)
}

func TestRun_ConsumeError(t *testing.T) {
	assert := assert.New(t)
	cnr := mockConsumer{
		chMsg: make(chan Message, 10),
		chErr: make(chan error, 10),
	}
	proc := mockProcessor{retError: true}
	cmp, err := New("test", proc.Process, &cnr)
	assert.NoError(err)
	ctx := context.Background()
	cnr.chErr <- errors.New("CONSUMER ERROR")
	err = cmp.Run(ctx)
	assert.Error(err)
}

func TestRun_Process_Shutdown(t *testing.T) {
	assert := assert.New(t)
	cnr := mockConsumer{
		chMsg: make(chan Message, 10),
		chErr: make(chan error, 10),
	}
	proc := mockProcessor{retError: false}
	cmp, err := New("test", proc.Process, &cnr)
	assert.NoError(err)
	cnr.chMsg <- &mockMessage{context.Background()}
	ch := make(chan bool)
	ctx := context.Background()
	go func() {
		err := cmp.Run(ctx)
		assert.NoError(err)
		ch <- true
	}()
	time.Sleep(10 * time.Millisecond)
	err = cmp.Shutdown(ctx)
	assert.NoError(err)
	<-ch
}

type mockMessage struct {
	ctx context.Context
}

func (mm *mockMessage) Context() context.Context {
	return mm.ctx
}

func (mm *mockMessage) Decode(v interface{}) error {
	return nil
}

func (mm *mockMessage) Ack() error {
	return nil
}

func (mm *mockMessage) Nack() error {
	return nil
}

type mockProcessor struct {
	retError bool
}

func (mp *mockProcessor) Process(msg Message) error {
	if mp.retError {
		return errors.New("PROC ERROR")
	}
	return nil
}

type mockConsumer struct {
	consumeError bool
	chMsg        chan Message
	chErr        chan error
}

func (mc *mockConsumer) Consume(context.Context) (<-chan Message, <-chan error, error) {
	if mc.consumeError {
		return nil, nil, errors.New("CONSUMER ERROR")
	}
	return mc.chMsg, mc.chErr, nil
}

func (mc *mockConsumer) Close() error {
	return nil
}
