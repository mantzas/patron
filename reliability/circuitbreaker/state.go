package circuitbreaker

import (
	"sync"
	"time"
)

// State definition
type State struct {
	set Setting
	sync.Mutex
	currentFailureCount  int
	retrySuccessCount    int
	currentExecutions    int
	lastFailureTimestamp time.Time
}

var (
	utcFuture = time.Date(9999, 12, 31, 23, 59, 59, 999999, time.UTC)
)

// NewState creates a new state. If no metric counter is provided
// the default null counter is used.
func NewState(s Setting) *State {
	return &State{set: s, currentFailureCount: 0, retrySuccessCount: 0, currentExecutions: 0, lastFailureTimestamp: utcFuture}
}

// Reset the state
func (s *State) Reset() {
	s.Lock()
	defer s.Unlock()
	s.innerReset()
}

func (s *State) innerReset() {
	s.currentFailureCount = 0
	s.retrySuccessCount = 0
	s.lastFailureTimestamp = utcFuture
}

// IncreaseFailure increases the failure count
func (s *State) IncreaseFailure() {
	s.Lock()
	defer s.Unlock()

	s.currentFailureCount++
	s.lastFailureTimestamp = time.Now().UTC()
}

// IncrementRetrySuccessCount increments the retry success count
func (s *State) IncrementRetrySuccessCount() {
	s.Lock()
	defer s.Unlock()

	s.retrySuccessCount++
}

// IncreaseExecutions increases the current execution count
func (s *State) IncreaseExecutions() {
	s.Lock()
	defer s.Unlock()

	s.currentExecutions++
}

// DecreaseExecutions decreases the current execution count
func (s *State) DecreaseExecutions() {
	s.Lock()
	defer s.Unlock()

	s.currentExecutions--
}

// GetStatus returns the status of the circuit
func (s *State) GetStatus() Status {
	s.Lock()
	defer s.Unlock()

	if s.set.FailureThreshold > s.currentFailureCount {
		return Closed
	}

	retry := s.lastFailureTimestamp.Add(s.set.RetryTimeout)
	now := time.Now().UTC()

	if retry.Before(now) || retry.Equal(now) {

		if s.retrySuccessCount >= s.set.RetrySuccessThreshold {
			s.innerReset()
			return Closed
		}

		if s.currentExecutions > s.set.MaxRetryExecutionThreshold {
			return Open
		}

		return HalfOpen
	}

	return Open
}
