package circuitbreaker

import (
	"time"

	"github.com/pkg/errors"
)

// Status of the circuit breaker.
type Status int

const (
	// Closed allow execution.
	Closed Status = iota
	// HalfOpen allowing execution to check if resource works again.
	HalfOpen
	// Open disallowing execution.
	Open
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

// CircuitBreaker implementation.
type CircuitBreaker struct {
	state *State
}

// NewCircuitBreaker constructor.
func NewCircuitBreaker(s Setting) *CircuitBreaker {
	return &CircuitBreaker{state: NewState(s)}
}

// Execute the function enclosed
func (cb *CircuitBreaker) Execute(act Action) (interface{}, error) {
	status := cb.state.GetStatus()
	if status == Open {
		return nil, errors.New("circuit is open")
	}

	cb.state.IncreaseExecutions()
	defer cb.state.DecreaseExecutions()

	resp, err := act()
	if err != nil {
		cb.state.IncreaseFailure()
		return nil, errors.Wrap(err, "Execution return error")
	}

	if status == HalfOpen {
		cb.state.IncrementRetrySuccessCount()
	}

	return resp, nil
}
