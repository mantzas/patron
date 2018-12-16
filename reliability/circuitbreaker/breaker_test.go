package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCircuitBreaker(t *testing.T) {
	set := Setting{FailureThreshold: 1, RetryTimeout: 1 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	cb := New(set)
	assert.Equal(t, 0, cb.failures)
	assert.Equal(t, 0, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.Equal(t, tsFuture, cb.nextRetry)
}

func TestCircuitBreaker_isHalfOpen(t *testing.T) {
	type fields struct {
		status    status
		nextRetry int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "closed", fields: fields{status: close, nextRetry: tsFuture}, want: false},
		{name: "open", fields: fields{status: open, nextRetry: time.Now().Add(1 * time.Hour).UnixNano()}, want: false},
		{name: "half open", fields: fields{status: open, nextRetry: time.Now().Add(-1 * time.Minute).UnixNano()}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := &CircuitBreaker{
				status:    tt.fields.status,
				nextRetry: tt.fields.nextRetry,
			}
			assert.Equal(t, tt.want, cb.isHalfOpen())
		})
	}
}

func TestCircuitBreaker_isOpen(t *testing.T) {
	type fields struct {
		status    status
		nextRetry int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "closed", fields: fields{status: close, nextRetry: tsFuture}, want: false},
		{name: "half open", fields: fields{status: open, nextRetry: time.Now().Add(-1 * time.Minute).UnixNano()}, want: false},
		{name: "open", fields: fields{status: open, nextRetry: time.Now().Add(1 * time.Hour).UnixNano()}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := &CircuitBreaker{
				status:    tt.fields.status,
				nextRetry: tt.fields.nextRetry,
			}
			assert.Equal(t, tt.want, cb.isOpen())
		})
	}
}

func TestCircuitBreaker_isClose(t *testing.T) {
	type fields struct {
		status    status
		nextRetry int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "closed", fields: fields{status: close, nextRetry: tsFuture}, want: true},
		{name: "half open", fields: fields{status: open, nextRetry: time.Now().Add(-1 * time.Minute).UnixNano()}, want: false},
		{name: "open", fields: fields{status: open, nextRetry: time.Now().Add(1 * time.Hour).UnixNano()}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := &CircuitBreaker{
				status:    tt.fields.status,
				nextRetry: tt.fields.nextRetry,
			}
			assert.Equal(t, tt.want, cb.isClose())
		})
	}
}

func TestCircuitBreaker_Close_Open_HalfOpen_Open_HalfOpen_Close(t *testing.T) {
	retryTimeout := 5 * time.Millisecond
	waitRetryTimeout := 7 * time.Millisecond

	set := Setting{FailureThreshold: 1, RetryTimeout: retryTimeout, RetrySuccessThreshold: 2, MaxRetryExecutionThreshold: 2}
	cb := New(set)
	_, err := cb.Execute(testSuccessAction)
	assert.NoError(t, err)
	assert.Equal(t, 0, cb.failures)
	assert.Equal(t, 0, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.True(t, cb.isClose())
	assert.Equal(t, tsFuture, cb.nextRetry)
	// will transition to open
	_, err = cb.Execute(testFailureAction)
	assert.EqualError(t, err, "Test error")
	assert.Equal(t, 0, cb.failures)
	assert.Equal(t, 0, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.True(t, cb.isOpen())
	assert.True(t, cb.nextRetry < tsFuture)
	// open, returns err immediately
	_, err = cb.Execute(testSuccessAction)
	assert.EqualError(t, err, "circuit is open")
	assert.Equal(t, 0, cb.failures)
	assert.Equal(t, 0, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.True(t, cb.isOpen())
	assert.True(t, cb.nextRetry < tsFuture)
	// should be half open now and will stay in there
	time.Sleep(waitRetryTimeout)
	_, err = cb.Execute(testFailureAction)
	assert.EqualError(t, err, "Test error")
	assert.Equal(t, 1, cb.failures)
	assert.Equal(t, 1, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.True(t, cb.isHalfOpen())
	assert.True(t, cb.nextRetry < tsFuture)
	// should be half open now and will transition to open
	_, err = cb.Execute(testFailureAction)
	assert.EqualError(t, err, "Test error")
	assert.Equal(t, 0, cb.failures)
	assert.Equal(t, 0, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.True(t, cb.isOpen())
	assert.True(t, cb.nextRetry < tsFuture)
	// should be half open now and will transition to close
	time.Sleep(waitRetryTimeout)
	_, err = cb.Execute(testSuccessAction)
	assert.NoError(t, err)
	assert.Equal(t, 0, cb.failures)
	assert.Equal(t, 1, cb.executions)
	assert.Equal(t, 1, cb.retries)
	assert.True(t, cb.isHalfOpen())
	assert.True(t, cb.nextRetry < tsFuture)
	_, err = cb.Execute(testSuccessAction)
	assert.NoError(t, err)
	assert.Equal(t, 0, cb.failures)
	assert.Equal(t, 0, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.True(t, cb.isClose())
	assert.Equal(t, tsFuture, cb.nextRetry)
}

var err error

func BenchmarkCircuitBreaker_Execute(b *testing.B) {
	set := Setting{FailureThreshold: 1, RetryTimeout: 1 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	c := New(set)

	for i := 0; i < b.N; i++ {
		_, err = c.Execute(testFailureAction)
	}
}

func testSuccessAction() (interface{}, error) {
	return "test", nil
}

func testFailureAction() (interface{}, error) {
	return "", errors.New("Test error")
}
