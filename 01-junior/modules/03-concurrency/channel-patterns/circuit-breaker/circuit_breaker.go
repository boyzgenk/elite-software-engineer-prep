// Circuit Breaker Pattern Implementation
// Used by: Netflix (Hystrix), Uber (fault tolerance), AWS Lambda (error handling)
// Problem: Prevent cascading failures when downstream services are failing
// Solution: Monitor failures and temporarily block requests to failing services

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// CircuitState represents the circuit breaker state
type CircuitState int32

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	MaxFailures      int           // Failure threshold to trip circuit
	Timeout          time.Duration // Time to wait before attempting recovery
	ResetTimeout     time.Duration // Time to wait in half-open state
	SuccessThreshold int           // Successes needed to close circuit from half-open
}

// DefaultConfig provides sensible defaults for circuit breaker
func DefaultConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxFailures:      5,
		Timeout:          10 * time.Second,
		ResetTimeout:     30 * time.Second,
		SuccessThreshold: 3,
	}
}

// CircuitBreaker implements the circuit breaker pattern with Go channels
// Similar to Netflix Hystrix but optimized for Go's concurrency model
type CircuitBreaker struct {
	config          CircuitBreakerConfig
	state           int32 // atomic access for CircuitState
	failures        int32 // atomic counter
	successes       int32 // atomic counter for half-open state
	lastFailureTime int64 // atomic timestamp
	mu              sync.RWMutex
	requests        chan Request
	responses       chan Response
	done            chan struct{}
	metrics         *CircuitMetrics
}

// Request represents a service request
type Request struct {
	ID       string
	Payload  interface{}
	Response chan Response
	Context  context.Context
}

// Response represents a service response
type Response struct {
	ID     string
	Result interface{}
	Error  error
}

// CircuitMetrics tracks circuit breaker statistics
type CircuitMetrics struct {
	TotalRequests int64
	SuccessCount  int64
	FailureCount  int64
	TimeoutCount  int64
	CircuitOpens  int64
	CircuitCloses int64
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	cb := &CircuitBreaker{
		config:    config,
		state:     int32(StateClosed),
		requests:  make(chan Request, 100), // Buffered for high throughput
		responses: make(chan Response, 100),
		done:      make(chan struct{}),
		metrics:   &CircuitMetrics{},
	}

	// Start the circuit breaker processor
	go cb.processRequests()

	return cb
}

// processRequests handles incoming requests based on circuit state
func (cb *CircuitBreaker) processRequests() {
	for {
		select {
		case req := <-cb.requests:
			cb.handleRequest(req)
		case <-cb.done:
			return
		}
	}
}

// handleRequest processes individual requests based on circuit state
func (cb *CircuitBreaker) handleRequest(req Request) {
	atomic.AddInt64(&cb.metrics.TotalRequests, 1)
	currentState := CircuitState(atomic.LoadInt32(&cb.state))

	switch currentState {
	case StateClosed:
		cb.executeRequest(req)

	case StateOpen:
		// Check if enough time has passed to try half-open
		lastFailure := atomic.LoadInt64(&cb.lastFailureTime)
		if time.Since(time.Unix(0, lastFailure)) > cb.config.Timeout {
			// Transition to half-open
			if atomic.CompareAndSwapInt32(&cb.state, int32(StateOpen), int32(StateHalfOpen)) {
				atomic.StoreInt32(&cb.successes, 0)
				fmt.Printf("🔄 Circuit breaker transitioning to HALF-OPEN\n")
			}
			cb.executeRequest(req)
		} else {
			// Circuit is open, reject request
			cb.rejectRequest(req, errors.New("circuit breaker is OPEN"))
		}

	case StateHalfOpen:
		cb.executeRequest(req)
	}
}

// executeRequest actually executes the request
func (cb *CircuitBreaker) executeRequest(req Request) {
	// Create a timeout context for the request
	ctx, cancel := context.WithTimeout(req.Context, 5*time.Second)
	defer cancel()

	// Execute in separate goroutine to handle timeouts
	resultChan := make(chan Response, 1)
	go func() {
		// Simulate service call - replace with actual service logic
		result, err := cb.simulateServiceCall(req.Payload)
		resultChan <- Response{
			ID:     req.ID,
			Result: result,
			Error:  err,
		}
	}()

	select {
	case response := <-resultChan:
		cb.handleResponse(req, response)
	case <-ctx.Done():
		// Timeout occurred
		atomic.AddInt64(&cb.metrics.TimeoutCount, 1)
		cb.handleFailure()
		cb.rejectRequest(req, errors.New("request timeout"))
	}
}

// simulateServiceCall simulates a service call that might fail
func (cb *CircuitBreaker) simulateServiceCall(payload interface{}) (interface{}, error) {
	// Simulate processing time
	time.Sleep(time.Duration(50+cb.failures*10) * time.Millisecond)

	// Simulate failure based on current failure count (more failures = higher chance)
	failureRate := float64(atomic.LoadInt32(&cb.failures)) / float64(cb.config.MaxFailures)
	if time.Now().UnixNano()%100 < int64(failureRate*50) {
		return nil, errors.New("service unavailable")
	}

	return fmt.Sprintf("Processed: %v", payload), nil
}

// handleResponse processes successful responses
func (cb *CircuitBreaker) handleResponse(req Request, response Response) {
	if response.Error != nil {
		cb.handleFailure()
		cb.sendResponse(req, response)
		return
	}

	// Success case
	atomic.AddInt64(&cb.metrics.SuccessCount, 1)
	atomic.StoreInt32(&cb.failures, 0) // Reset failure counter

	currentState := CircuitState(atomic.LoadInt32(&cb.state))
	if currentState == StateHalfOpen {
		successes := atomic.AddInt32(&cb.successes, 1)
		if successes >= int32(cb.config.SuccessThreshold) {
			// Close the circuit
			if atomic.CompareAndSwapInt32(&cb.state, int32(StateHalfOpen), int32(StateClosed)) {
				atomic.AddInt64(&cb.metrics.CircuitCloses, 1)
				fmt.Printf("✅ Circuit breaker CLOSED (enough successes)\n")
			}
		}
	}

	cb.sendResponse(req, response)
}

// handleFailure handles failed requests
func (cb *CircuitBreaker) handleFailure() {
	atomic.AddInt64(&cb.metrics.FailureCount, 1)
	failures := atomic.AddInt32(&cb.failures, 1)
	atomic.StoreInt64(&cb.lastFailureTime, time.Now().UnixNano())

	// Check if we should open the circuit
	if failures >= int32(cb.config.MaxFailures) {
		currentState := atomic.LoadInt32(&cb.state)
		if currentState == int32(StateClosed) || currentState == int32(StateHalfOpen) {
			if atomic.CompareAndSwapInt32(&cb.state, currentState, int32(StateOpen)) {
				atomic.AddInt64(&cb.metrics.CircuitOpens, 1)
				fmt.Printf("🚨 Circuit breaker OPENED (too many failures: %d)\n", failures)
			}
		}
	}
}

// rejectRequest rejects a request due to circuit being open
func (cb *CircuitBreaker) rejectRequest(req Request, err error) {
	response := Response{
		ID:    req.ID,
		Error: err,
	}
	cb.sendResponse(req, response)
}

// sendResponse sends response back to requester
func (cb *CircuitBreaker) sendResponse(req Request, response Response) {
	select {
	case req.Response <- response:
		// Response sent successfully
	case <-req.Context.Done():
		// Request context cancelled, don't send response
	}
}

// Execute makes a request through the circuit breaker
func (cb *CircuitBreaker) Execute(ctx context.Context, id string, payload interface{}) (interface{}, error) {
	responseChan := make(chan Response, 1)

	request := Request{
		ID:       id,
		Payload:  payload,
		Response: responseChan,
		Context:  ctx,
	}

	select {
	case cb.requests <- request:
		// Request queued successfully
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case response := <-responseChan:
		return response.Result, response.Error
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetState returns current circuit breaker state
func (cb *CircuitBreaker) GetState() CircuitState {
	return CircuitState(atomic.LoadInt32(&cb.state))
}

// GetMetrics returns current metrics
func (cb *CircuitBreaker) GetMetrics() CircuitMetrics {
	return CircuitMetrics{
		TotalRequests: atomic.LoadInt64(&cb.metrics.TotalRequests),
		SuccessCount:  atomic.LoadInt64(&cb.metrics.SuccessCount),
		FailureCount:  atomic.LoadInt64(&cb.metrics.FailureCount),
		TimeoutCount:  atomic.LoadInt64(&cb.metrics.TimeoutCount),
		CircuitOpens:  atomic.LoadInt64(&cb.metrics.CircuitOpens),
		CircuitCloses: atomic.LoadInt64(&cb.metrics.CircuitCloses),
	}
}

// Close shuts down the circuit breaker
func (cb *CircuitBreaker) Close() {
	close(cb.done)
}

// CircuitBreakerPool manages multiple circuit breakers for different services
// Like Netflix's approach of having separate circuit breakers per service
type CircuitBreakerPool struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
	config   CircuitBreakerConfig
}

// NewCircuitBreakerPool creates a new pool of circuit breakers
func NewCircuitBreakerPool(config CircuitBreakerConfig) *CircuitBreakerPool {
	return &CircuitBreakerPool{
		breakers: make(map[string]*CircuitBreaker),
		config:   config,
	}
}

// GetBreaker returns circuit breaker for a service (creates if not exists)
func (pool *CircuitBreakerPool) GetBreaker(serviceName string) *CircuitBreaker {
	pool.mu.RLock()
	breaker, exists := pool.breakers[serviceName]
	pool.mu.RUnlock()

	if exists {
		return breaker
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Double-check pattern
	if breaker, exists := pool.breakers[serviceName]; exists {
		return breaker
	}

	breaker = NewCircuitBreaker(pool.config)
	pool.breakers[serviceName] = breaker
	return breaker
}

// Execute executes request through service-specific circuit breaker
func (pool *CircuitBreakerPool) Execute(ctx context.Context, serviceName, requestID string, payload interface{}) (interface{}, error) {
	breaker := pool.GetBreaker(serviceName)
	return breaker.Execute(ctx, requestID, payload)
}

func main() {
	fmt.Println("🔌 Circuit Breaker Pattern Demo")
	fmt.Println("=====================================")

	// Demo 1: Basic Circuit Breaker
	fmt.Println("\n1️⃣ Basic Circuit Breaker (Netflix Hystrix-style)")

	config := DefaultConfig()
	config.MaxFailures = 3 // Lower threshold for demo
	cb := NewCircuitBreaker(config)
	defer cb.Close()

	ctx := context.Background()

	// Test circuit breaker behavior
	for i := 1; i <= 15; i++ {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

		result, err := cb.Execute(ctx, fmt.Sprintf("req-%d", i), fmt.Sprintf("data-%d", i))

		state := cb.GetState()
		if err != nil {
			fmt.Printf("❌ Request %d: %v (State: %s)\n", i, err, state)
		} else {
			fmt.Printf("✅ Request %d: %v (State: %s)\n", i, result, state)
		}

		cancel()
		time.Sleep(500 * time.Millisecond)
	}

	// Demo 2: Circuit Breaker Pool for Multiple Services
	fmt.Println("\n2️⃣ Circuit Breaker Pool (Multi-service)")

	pool := NewCircuitBreakerPool(config)
	services := []string{"user-service", "payment-service", "inventory-service"}

	var wg sync.WaitGroup
	for _, service := range services {
		wg.Add(1)
		go func(serviceName string) {
			defer wg.Done()

			for i := 1; i <= 5; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

				result, err := pool.Execute(ctx, serviceName, fmt.Sprintf("%s-req-%d", serviceName, i), fmt.Sprintf("data-%d", i))

				breaker := pool.GetBreaker(serviceName)
				state := breaker.GetState()

				if err != nil {
					fmt.Printf("❌ %s Request %d: %v (State: %s)\n", serviceName, i, err, state)
				} else {
					fmt.Printf("✅ %s Request %d: %v (State: %s)\n", serviceName, i, result, state)
				}

				cancel()
				time.Sleep(200 * time.Millisecond)
			}
		}(service)
	}

	wg.Wait()

	// Demo 3: Metrics and Monitoring
	fmt.Println("\n3️⃣ Circuit Breaker Metrics")

	metrics := cb.GetMetrics()
	fmt.Printf("📊 Circuit Breaker Metrics:\n")
	fmt.Printf("   Total Requests: %d\n", metrics.TotalRequests)
	fmt.Printf("   Successes: %d\n", metrics.SuccessCount)
	fmt.Printf("   Failures: %d\n", metrics.FailureCount)
	fmt.Printf("   Timeouts: %d\n", metrics.TimeoutCount)
	fmt.Printf("   Circuit Opens: %d\n", metrics.CircuitOpens)
	fmt.Printf("   Circuit Closes: %d\n", metrics.CircuitCloses)

	if metrics.TotalRequests > 0 {
		successRate := float64(metrics.SuccessCount) / float64(metrics.TotalRequests) * 100
		fmt.Printf("   Success Rate: %.2f%%\n", successRate)
	}

	fmt.Println("\n🎯 Circuit Breaker Pattern Complete!")
	fmt.Println("Key Benefits:")
	fmt.Println("• Prevents cascading failures in distributed systems")
	fmt.Println("• Provides fail-fast behavior when services are down")
	fmt.Println("• Automatic recovery testing with half-open state")
	fmt.Println("• Detailed metrics for monitoring and alerting")
	fmt.Println("• Per-service isolation with circuit breaker pools")
}
