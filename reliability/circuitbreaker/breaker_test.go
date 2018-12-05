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
	assert.Equal(t, utcFuture, cb.lastFailure)
}

func TestCircuitBreaker_Closed(t *testing.T) {
	set := Setting{FailureThreshold: 1, RetryTimeout: 1 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	cb := New(set)
	_, err := cb.Execute(testSuccessAction)
	assert.NoError(t, err)
	assert.Equal(t, 0, cb.failures)
	assert.Equal(t, 0, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.Equal(t, utcFuture, cb.lastFailure)
}

func TestCircuitBreaker_Open(t *testing.T) {
	set := Setting{FailureThreshold: 1, RetryTimeout: 1 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	cb := New(set)
	_, err := cb.Execute(testFailureAction)
	assert.Error(t, err)
	assert.Equal(t, 1, cb.failures)
	assert.Equal(t, 0, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.True(t, time.Now().UTC().After(cb.lastFailure))
	_, err = cb.Execute(testFailureAction)
	assert.Error(t, err)
	assert.Equal(t, 1, cb.failures)
	assert.Equal(t, 0, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.True(t, time.Now().UTC().After(cb.lastFailure))
}

var err error

func BenchmarkCircuitBreaker_Execute(b *testing.B) {
	set := Setting{FailureThreshold: 1, RetryTimeout: 1 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	c := New(set)

	for i := 0; i < b.N; i++ {
		_, err = c.Execute(testSuccessAction)
	}
}

func testSuccessAction() (interface{}, error) {
	return "test", nil
}

func testFailureAction() (interface{}, error) {
	return "", errors.New("Test error")
}
