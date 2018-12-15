package circuitbreaker

import (
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
	halfOpen
	open
)

var (
	tsFuture  = time.Date(9999, 12, 31, 23, 59, 59, 999999, time.UTC).UnixNano()
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
	sync.Mutex
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

// Execute the function enclosed.
func (cb *CircuitBreaker) Execute(act Action) (interface{}, error) {

	// calculate status
	cb.calcStatus()

	if cb.status == open {
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

func (cb *CircuitBreaker) calcStatus() {
	cb.Lock()
	defer cb.Unlock()

	switch cb.status {
	case close:
		if cb.failures >= cb.set.FailureThreshold {
			cb.transitionToOpen()
			return
		}

	case open:
		if cb.nextRetry >= time.Now().UnixNano() {
			cb.transitionToHalfOpen()
			return
		}

	case halfOpen:
		if cb.retries >= cb.set.RetrySuccessThreshold {
			cb.transitionToClose()
			return
		}

		if cb.executions >= cb.set.MaxRetryExecutionThreshold {
			cb.transitionToOpen()
		}
	}

	return
}

func (cb *CircuitBreaker) transitionToOpen() {
	cb.status = open
	cb.failures = 0
	cb.executions = 0
	cb.retries = 0
	cb.nextRetry = time.Now().Add(cb.set.RetryTimeout).UnixNano()
}

func (cb *CircuitBreaker) transitionToHalfOpen() {
	cb.status = halfOpen
	cb.failures = 0
	cb.executions = 0
	cb.retries = 0
}

func (cb *CircuitBreaker) transitionToClose() {
	cb.status = close
	cb.failures = 0
	cb.executions = 0
	cb.retries = 0
	cb.nextRetry = tsFuture
}

func (cb *CircuitBreaker) incFailure() {
	cb.Lock()
	defer cb.Unlock()

	cb.failures++
	cb.executions++
}

func (cb *CircuitBreaker) incSuccess() {
	cb.Lock()
	defer cb.Unlock()

	cb.retries++
	cb.executions++
}
