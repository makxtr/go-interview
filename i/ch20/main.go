// Package challenge20 contains the implementation for Challenge 20: Circuit Breaker Pattern
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// String returns the string representation of the state
func (s State) String() string {
	switch s {
	case StateClosed:
		return "Closed"
	case StateOpen:
		return "Open"
	case StateHalfOpen:
		return "Half-Open"
	default:
		return "Unknown"
	}
}

// Metrics represents the circuit breaker metrics
type Metrics struct {
	Requests            int64
	Successes           int64
	Failures            int64
	ConsecutiveFailures int64
	LastFailureTime     time.Time
}

// Config represents the configuration for the circuit breaker
type Config struct {
	MaxRequests   uint32                                  // Max requests allowed in half-open state
	Interval      time.Duration                           // Statistical window for closed state
	Timeout       time.Duration                           // Time to wait before half-open
	ReadyToTrip   func(Metrics) bool                      // Function to determine when to trip
	OnStateChange func(name string, from State, to State) // State change callback
}

// CircuitBreaker interface defines the operations for a circuit breaker
type CircuitBreaker interface {
	Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error)
	GetState() State
	GetMetrics() Metrics
}

// circuitBreakerImpl is the concrete implementation of CircuitBreaker
type circuitBreakerImpl struct {
	name             string
	config           Config
	state            State
	metrics          Metrics
	lastStateChange  time.Time
	halfOpenRequests uint32
	mutex            sync.Mutex
}

// Error definitions
var (
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
	ErrTooManyRequests    = errors.New("too many requests in half-open state")
)

// NewCircuitBreaker creates a new circuit breaker with the given configuration
func NewCircuitBreaker(config Config) CircuitBreaker {
	if config.MaxRequests == 0 {
		config.MaxRequests = 1
	}
	if config.Interval == 0 {
		config.Interval = time.Minute
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.ReadyToTrip == nil {
		config.ReadyToTrip = func(m Metrics) bool {
			return m.ConsecutiveFailures >= 5
		}
	}

	return &circuitBreakerImpl{
		name:            "circuit-breaker",
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

// Call executes the given operation through the circuit breaker
func (cb *circuitBreakerImpl) Call(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := cb.beforeCall(); err != nil {
		return nil, err
	}

	result, err := operation()
	cb.afterCall(err)

	if err != nil {
		return nil, err
	}
	return result, nil
}

// beforeCall validates the breaker permits execution and performs any eager
// state transition (e.g. Open -> HalfOpen once Timeout has elapsed).
func (cb *circuitBreakerImpl) beforeCall() error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == StateOpen {
		if time.Since(cb.lastStateChange) < cb.config.Timeout {
			return ErrCircuitBreakerOpen
		}
		cb.setStateLocked(StateHalfOpen)
	}

	if cb.state == StateHalfOpen {
		if cb.halfOpenRequests >= cb.config.MaxRequests {
			return ErrTooManyRequests
		}
		cb.halfOpenRequests++
	}

	return nil
}

// afterCall records the operation outcome and drives state transitions.
func (cb *circuitBreakerImpl) afterCall(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.metrics.Requests++
	if err != nil {
		cb.metrics.Failures++
		cb.metrics.ConsecutiveFailures++
		cb.metrics.LastFailureTime = time.Now()

		switch cb.state {
		case StateClosed:
			if cb.config.ReadyToTrip(cb.metrics) {
				cb.setStateLocked(StateOpen)
			}
		case StateHalfOpen:
			cb.setStateLocked(StateOpen)
		}
		return
	}

	cb.metrics.Successes++
	cb.metrics.ConsecutiveFailures = 0
	if cb.state == StateHalfOpen {
		cb.setStateLocked(StateClosed)
	}
}

// GetState returns the current state of the circuit breaker
func (cb *circuitBreakerImpl) GetState() State {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return cb.state
}

// GetMetrics returns the current metrics of the circuit breaker
func (cb *circuitBreakerImpl) GetMetrics() Metrics {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	return cb.metrics
}

// setStateLocked transitions state and resets window-scoped counters.
// Must be called with cb.mutex held.
func (cb *circuitBreakerImpl) setStateLocked(newState State) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = time.Now()
	cb.halfOpenRequests = 0

	if newState == StateClosed {
		cb.metrics = Metrics{}
	}

	if cb.config.OnStateChange != nil {
		cb.config.OnStateChange(cb.name, oldState, newState)
	}
}

// Example usage and testing helper functions
func main() {
	fmt.Println("Circuit Breaker Pattern Example")

	config := Config{
		MaxRequests: 3,
		Interval:    time.Minute,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(m Metrics) bool {
			return m.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from State, to State) {
			fmt.Printf("Circuit breaker %s: %s -> %s\n", name, from, to)
		},
	}

	cb := NewCircuitBreaker(config)
	ctx := context.Background()

	result, err := cb.Call(ctx, func() (interface{}, error) {
		return "success", nil
	})
	fmt.Printf("Result: %v, Error: %v\n", result, err)

	result, err = cb.Call(ctx, func() (interface{}, error) {
		return nil, errors.New("simulated failure")
	})
	fmt.Printf("Result: %v, Error: %v\n", result, err)

	fmt.Printf("Current state: %v\n", cb.GetState())
	fmt.Printf("Current metrics: %+v\n", cb.GetMetrics())
}