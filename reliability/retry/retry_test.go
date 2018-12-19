package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		retries int
		delay   time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{retries: 3, delay: 3 * time.Second}, wantErr: false},
		{name: "invalid retries", args: args{retries: -1, delay: 3 * time.Second}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.retries, tt.args.delay)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestRetry_Execute_Success(t *testing.T) {
	r, err := New(3, 10*time.Millisecond)
	assert.NoError(t, err)
	act := mockAction{errors: 1}
	res, err := r.Execute(func() (interface{}, error) {
		return act.Execute()
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, act.executions)
}

func TestRetry_Execute_Failed(t *testing.T) {
	r, err := New(3, 10*time.Millisecond)
	assert.NoError(t, err)
	act := mockAction{errors: 3}
	res, err := r.Execute(func() (interface{}, error) {
		return act.Execute()
	})
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, 3, act.executions)
}

type mockAction struct {
	errors     int
	executions int
}

func (ma *mockAction) Execute() (string, error) {
	defer func() {
		ma.errors--
		ma.executions++
	}()
	if ma.errors > 0 {
		return "", errors.New("TEST")
	}
	return "TEST", nil
}

var err error

func BenchmarkCircuitBreaker_Execute(b *testing.B) {
	var r *Retry
	r, err = New(3, 0)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = r.Execute(testSuccessAction)
	}
}

func testSuccessAction() (interface{}, error) {
	return "test", nil
}
