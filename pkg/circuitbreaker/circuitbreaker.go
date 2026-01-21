package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	mu                   sync.RWMutex
	state                State
	failureThreshold     int
	successThreshold     int
	timeout              time.Duration
	consecutiveFailures  int
	consecutiveSuccesses int
	lastErrorTime        time.Time
}

func New(failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
	}
}

func (cb *CircuitBreaker) Execute(op func() (interface{}, error)) (interface{}, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastErrorTime) > cb.timeout {
			cb.state = StateHalfOpen
			cb.consecutiveSuccesses = 0
		} else {
			return nil, ErrCircuitOpen
		}
	case StateHalfOpen:
	case StateClosed:
	}

	res, err := op()
	if err != nil {
		cb.handleFailure()
		return nil, err
	}

	cb.handleSuccess()
	return res, nil
}

func (cb *CircuitBreaker) handleFailure() {
	cb.consecutiveFailures++
	if cb.state == StateHalfOpen || cb.consecutiveFailures >= cb.failureThreshold {
		cb.state = StateOpen
		cb.lastErrorTime = time.Now()
	}
}

func (cb *CircuitBreaker) handleSuccess() {
	if cb.state == StateHalfOpen {
		cb.consecutiveSuccesses++
		if cb.consecutiveSuccesses >= cb.successThreshold {
			cb.state = StateClosed
			cb.consecutiveFailures = 0
		}
	} else {
		cb.consecutiveFailures = 0
	}
}
