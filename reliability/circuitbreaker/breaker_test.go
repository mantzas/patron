package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCircuitBreaker(t *testing.T) {
	c := NewCircuitBreaker(Setting{})
	assert.NotNil(t, c)
}

func TestExecute_Closed(t *testing.T) {
	set := Setting{FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	res, err := NewCircuitBreaker(set).Execute(testSuccessAction)

	assert.Nil(t, err)
	assert.Equal(t, "test", res)
}

func TestExecute_Open(t *testing.T) {
	set := Setting{FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}

	cb := NewCircuitBreaker(set)
	cb.state.IncreaseFailure()

	_, err := cb.Execute(testSuccessAction)

	assert.NotNil(t, err)
}

func TestExecute_Failed(t *testing.T) {
	set := Setting{FailureThreshold: 1, RetryTimeout: 10 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}

	_, err := NewCircuitBreaker(set).Execute(testFailureAction)

	assert.NotNil(t, err)
}

func TestExecute_SuccessAfterFailed(t *testing.T) {
	set := Setting{FailureThreshold: 1, RetryTimeout: 10 * time.Millisecond, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}

	cb := NewCircuitBreaker(set)
	_, err := cb.Execute(testFailureAction)
	assert.NotNil(t, err)
	time.Sleep(20 * time.Millisecond)
	_, err = cb.Execute(testSuccessAction)

	assert.Nil(t, err)
}

var err error

func BenchmarkCircuitBreaker_Execute(b *testing.B) {
	set := Setting{FailureThreshold: 1, RetryTimeout: 1 * time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	c := NewCircuitBreaker(set)

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
