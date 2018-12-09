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
	assert.Equal(t, utcFuture, cb.nextRetry)
}

// func TestCircuitBreaker_Closed(t *testing.T) {
// 	set := Setting{FailureThreshold: 1, RetryTimeout: 1 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
// 	cb := New(set)
// 	_, err := cb.Execute(testSuccessAction)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 0, cb.failures)
// 	assert.Equal(t, 0, cb.executions)
// 	assert.Equal(t, 0, cb.retries)
// 	assert.Equal(t, utcFuture, cb.nextRetry)
// }

// func TestCircuitBreaker_Open(t *testing.T) {
// 	set := Setting{FailureThreshold: 1, RetryTimeout: 1 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
// 	cb := New(set)
// 	_, err := cb.Execute(testFailureAction)
// 	assert.Error(t, err)
// 	assert.Equal(t, 1, cb.failures)
// 	assert.Equal(t, 0, cb.executions)
// 	assert.Equal(t, 0, cb.retries)
// 	assert.True(t, time.Now().UTC().Before(cb.nextRetry))
// 	_, err = cb.Execute(testFailureAction)
// 	assert.Error(t, err)
// 	assert.Equal(t, 1, cb.failures)
// 	assert.Equal(t, 0, cb.executions)
// 	assert.Equal(t, 0, cb.retries)
// 	assert.True(t, time.Now().UTC().Before(cb.nextRetry))
// }

// func TestCircuitBreaker_HalfOpen_Closed(t *testing.T) {
// 	set := Setting{FailureThreshold: 1, RetryTimeout: 5 * time.Millisecond, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
// 	cb := New(set)
// 	_, err := cb.Execute(testFailureAction)
// 	assert.Error(t, err)
// 	assert.Equal(t, 1, cb.failures)
// 	assert.Equal(t, 0, cb.executions)
// 	assert.Equal(t, 0, cb.retries)
// 	assert.True(t, time.Now().UTC().Before(cb.nextRetry))
// 	time.Sleep(10 * time.Millisecond)
// 	_, err = cb.Execute(testSuccessAction)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, cb.failures)
// 	assert.Equal(t, 1, cb.executions)
// 	assert.Equal(t, 1, cb.retries)
// 	assert.True(t, time.Now().UTC().After(cb.nextRetry))
// 	_, err = cb.Execute(testSuccessAction)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 0, cb.failures)
// 	assert.Equal(t, 0, cb.executions)
// 	assert.Equal(t, 0, cb.retries)
// 	assert.Equal(t, utcFuture, cb.nextRetry)
// }

func TestCircuitBreaker_HalfOpen_Open(t *testing.T) {
	set := Setting{
		FailureThreshold:           1,
		RetryTimeout:               5 * time.Millisecond,
		RetrySuccessThreshold:      3,
		MaxRetryExecutionThreshold: 2,
	}
	cb := New(set)
	_, err := cb.Execute(testFailureAction)
	assert.Error(t, err)
	assert.Equal(t, 1, cb.failures)
	assert.Equal(t, 0, cb.executions)
	assert.Equal(t, 0, cb.retries)
	assert.True(t, time.Now().UTC().Before(cb.nextRetry))
	assert.Equal(t, closed, cb.status)
	time.Sleep(10 * time.Millisecond)
	_, err = cb.Execute(testSuccessAction)
	assert.NoError(t, err)
	assert.Equal(t, 1, cb.failures)
	assert.Equal(t, 1, cb.executions)
	assert.Equal(t, 1, cb.retries)
	assert.True(t, time.Now().UTC().After(cb.nextRetry))
	assert.Equal(t, halfOpen, cb.status)
	_, err = cb.Execute(testSuccessAction)
	assert.NoError(t, err)
	assert.Equal(t, 1, cb.failures)
	assert.Equal(t, 2, cb.executions)
	assert.Equal(t, 2, cb.retries)
	assert.True(t, time.Now().UTC().After(cb.nextRetry))
	assert.Equal(t, halfOpen, cb.status)
	_, err = cb.Execute(testSuccessAction)
	assert.NoError(t, err)
	assert.Equal(t, 1, cb.failures)
	assert.Equal(t, 3, cb.executions)
	assert.Equal(t, 3, cb.retries)
	assert.True(t, time.Now().UTC().After(cb.nextRetry))
	assert.Equal(t, halfOpen, cb.status)
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
