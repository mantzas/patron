package circuitbreaker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewState(t *testing.T) {
	setting := Setting{FailureThreshold: 1, RetryTimeout: time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	state := NewState(setting)
	assert.Equal(t, 0, state.currentExecutions)
	assert.Equal(t, 0, state.currentFailureCount)
	assert.Equal(t, 0, state.retrySuccessCount)
	assert.Equal(t, time.Date(9999, 12, 31, 23, 59, 59, 999999, time.UTC), state.lastFailureTimestamp)
}

func TestState_Reset(t *testing.T) {
	setting := Setting{FailureThreshold: 1, RetryTimeout: time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	state := NewState(setting)
	state.IncreaseFailure()
	state.IncrementRetrySuccessCount()

	state.Reset()

	assert.Equal(t, 0, state.currentFailureCount)
	assert.Equal(t, 0, state.retrySuccessCount)
}

func TestState_IncreaseFailure(t *testing.T) {
	setting := Setting{FailureThreshold: 1, RetryTimeout: time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	state := NewState(setting)
	state.IncreaseFailure()
	assert.Equal(t, 1, state.currentFailureCount)
}

func TestState_IncrementRetrySuccessCount(t *testing.T) {
	setting := Setting{FailureThreshold: 1, RetryTimeout: time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	state := NewState(setting)
	state.IncrementRetrySuccessCount()
	assert.Equal(t, 1, state.retrySuccessCount)
}

func TestState_IncreaseExecutions(t *testing.T) {
	setting := Setting{FailureThreshold: 1, RetryTimeout: time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	state := NewState(setting)
	state.IncreaseExecutions()
	assert.Equal(t, 1, state.currentExecutions)
}

func TestState_DecreaseExecutions(t *testing.T) {
	setting := Setting{FailureThreshold: 1, RetryTimeout: time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	state := NewState(setting)
	state.DecreaseExecutions()
	assert.Equal(t, -1, state.currentExecutions)
}

func TestState_GetStatus(t *testing.T) {
	setting := Setting{FailureThreshold: 1, RetryTimeout: time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	stateClosed := NewState(setting)
	stateClosed.IncreaseFailure()
	stateHalf := NewState(setting)
	stateHalf.IncreaseFailure()
	stateHalf.lastFailureTimestamp = stateHalf.lastFailureTimestamp.Add(-2 * time.Second)

	stateOpenMaxRetry := NewState(setting)
	stateOpenMaxRetry.IncreaseFailure()
	stateOpenMaxRetry.lastFailureTimestamp = stateHalf.lastFailureTimestamp.Add(-2 * time.Second)
	stateOpenMaxRetry.IncreaseExecutions()
	stateOpenMaxRetry.IncreaseExecutions()

	stateClosedRetrySuccess := NewState(setting)
	stateClosedRetrySuccess.IncreaseFailure()
	stateClosedRetrySuccess.lastFailureTimestamp = stateHalf.lastFailureTimestamp.Add(-2 * time.Second)
	stateClosedRetrySuccess.IncrementRetrySuccessCount()

	tests := []struct {
		name string
		s    *State
		want Status
	}{
		{"Closed", NewState(setting), Closed},
		{"Open", stateClosed, Open},
		{"HalfOpen", stateHalf, HalfOpen},
		{"Open Max Retry", stateOpenMaxRetry, Open},
		{"Closes after retry success", stateClosedRetrySuccess, Closed},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.s.GetStatus())
	}
}

var s Status

func BenchmarkState_GetStatus(b *testing.B) {
	setting := Setting{FailureThreshold: 1, RetryTimeout: time.Second, RetrySuccessThreshold: 1, MaxRetryExecutionThreshold: 1}
	state := NewState(setting)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s = state.GetStatus()
	}
}
