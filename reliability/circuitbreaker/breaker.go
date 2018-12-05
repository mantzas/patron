package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

type status int

const (
	closed status = iota
	halfOpen
	open
)

var (
	utcFuture = time.Date(9999, 12, 31, 23, 59, 59, 999999, time.UTC)
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
	sync.Mutex
	executions  int
	failures    int
	retries     int
	lastFailure time.Time
}

// New constructor.
func New(s Setting) *CircuitBreaker {
	return &CircuitBreaker{
		set:         s,
		executions:  0,
		failures:    0,
		retries:     0,
		lastFailure: utcFuture,
	}
}

// Execute the function enclosed.
func (cb *CircuitBreaker) Execute(act Action) (interface{}, error) {
	status := cb.status()
	if status == open {
		return nil, errors.New("circuit is open")
	}

	cb.incrExecutions()

	resp, err := act()
	if err != nil {
		cb.incFailure()
		return nil, err
	}

	if status == halfOpen {
		cb.incRetrySuccess()
	} else {
		cb.reset()
	}

	return resp, nil
}

func (cb *CircuitBreaker) status() status {
	cb.Lock()
	defer cb.Unlock()

	if cb.set.FailureThreshold > cb.failures {
		return closed
	}

	retry := cb.lastFailure.Add(cb.set.RetryTimeout)
	now := time.Now().UTC()

	if retry.Before(now) || retry.Equal(now) {

		if cb.retries >= cb.set.RetrySuccessThreshold {
			return closed
		}

		if cb.executions > cb.set.MaxRetryExecutionThreshold {
			return open
		}

		return halfOpen
	}

	return open
}

func (cb *CircuitBreaker) incrExecutions() {
	cb.Lock()
	defer cb.Unlock()

	cb.executions++
}

func (cb *CircuitBreaker) incFailure() {
	cb.Lock()
	defer cb.Unlock()

	cb.failures++
	cb.executions = 0
	cb.lastFailure = time.Now().UTC()
}

func (cb *CircuitBreaker) incRetrySuccess() {
	cb.Lock()
	defer cb.Unlock()

	cb.retries++
}

func (cb *CircuitBreaker) reset() {
	cb.Lock()
	defer cb.Unlock()

	cb.failures = 0
	cb.executions = 0
	cb.retries = 0
	cb.lastFailure = utcFuture
}
