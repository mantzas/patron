package circuitbreaker

import (
	"math"
	"sync"
	"time"
)

// OpenError defines a open error.
type OpenError struct {
}

func (oe OpenError) Error() string {
	return "circuit is open"
}

type status int

const (
	close status = iota
	open
)

var (
	tsFuture  = int64(math.MaxInt64)
	openError = new(OpenError)
)

// Setting definition.
type Setting struct {
	// The threshold for the circuit to open.
	FailureThreshold int
	// The timeout after which we set the state to half-open and allow a retry.
	RetryTimeout time.Duration
	// The threshold of the retry successes which returns the state to open.
	RetrySuccessThreshold int
	// The threshold of how many retry executions are allowed when the status is half-open.
	MaxRetryExecutionThreshold int
}

// Action function to execute in circuit breaker.
type Action func() (interface{}, error)

// Executor interface.
type Executor interface {
	Execute(act Action) (interface{}, error)
}

// CircuitBreaker implementation.
type CircuitBreaker struct {
	set Setting
	sync.RWMutex
	status     status
	executions int
	failures   int
	retries    int
	nextRetry  int64
}

// New constructor.
func New(s Setting) *CircuitBreaker {
	return &CircuitBreaker{
		set:        s,
		executions: 0,
		failures:   0,
		retries:    0,
		nextRetry:  tsFuture,
	}
}

func (cb *CircuitBreaker) isHalfOpen() bool {
	cb.RLock()
	defer cb.RUnlock()
	if cb.status == open && cb.nextRetry <= time.Now().UnixNano() {
		return true
	}
	return false
}

func (cb *CircuitBreaker) isOpen() bool {
	cb.RLock()
	defer cb.RUnlock()
	if cb.status == open && cb.nextRetry > time.Now().UnixNano() {
		return true
	}
	return false
}

func (cb *CircuitBreaker) isClose() bool {
	cb.RLock()
	defer cb.RUnlock()
	if cb.status == close {
		return true
	}
	return false
}

// Execute the function enclosed.
func (cb *CircuitBreaker) Execute(act Action) (interface{}, error) {
	if cb.isOpen() {
		return nil, openError
	}

	resp, err := act()
	if err != nil {
		cb.incFailure()
		return nil, err
	}
	cb.incSuccess()

	return resp, err
}

func (cb *CircuitBreaker) incFailure() {
	// allow closed and half open to transition to open
	if cb.isOpen() {
		return
	}
	cb.Lock()
	defer cb.Unlock()

	cb.failures++

	if cb.status == close && cb.failures >= cb.set.FailureThreshold {
		cb.transitionToOpen()
		return
	}

	cb.executions++

	if cb.executions < cb.set.MaxRetryExecutionThreshold {
		return
	}

	cb.transitionToOpen()
}

func (cb *CircuitBreaker) incSuccess() {
	// allow only half open in order to transition to closed
	if !cb.isHalfOpen() {
		return
	}
	cb.Lock()
	defer cb.Unlock()

	cb.retries++
	cb.executions++

	if cb.retries < cb.set.RetrySuccessThreshold {
		return
	}
	cb.transitionToClose()
}

func (cb *CircuitBreaker) transitionToOpen() {
	cb.status = open
	cb.failures = 0
	cb.executions = 0
	cb.retries = 0
	cb.nextRetry = time.Now().Add(cb.set.RetryTimeout).UnixNano()
}

func (cb *CircuitBreaker) transitionToClose() {
	cb.status = close
	cb.failures = 0
	cb.executions = 0
	cb.retries = 0
	cb.nextRetry = tsFuture
}
