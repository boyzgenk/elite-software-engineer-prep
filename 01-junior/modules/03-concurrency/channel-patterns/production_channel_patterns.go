// Package channel_patterns demonstrates production-grade channel patterns
// used at Netflix, Uber, and other high-scale systems
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// FAN-OUT / FAN-IN PATTERN - Netflix Style Data Processing
// =============================================================================

// Job represents work to be processed
type Job struct {
	ID   int
	Data string
}

// Result represents processed job result
type Result struct {
	JobID       int
	Output      string
	ProcessedAt time.Time
	WorkerID    int
}

// FanOutFanIn demonstrates distributing work to multiple workers
// and collecting results - essential pattern for microservices
func FanOutFanIn(jobs []Job, numWorkers int) []Result {
	// Fan-out: Distribute jobs to multiple workers
	jobChan := make(chan Job, len(jobs))
	resultChan := make(chan Result, len(jobs))

	// Start workers (fan-out)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i+1, jobChan, resultChan, &wg)
	}

	// Send jobs to workers
	go func() {
		defer close(jobChan)
		for _, job := range jobs {
			jobChan <- job
		}
	}()

	// Close result channel when all workers finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Fan-in: Collect all results
	var results []Result
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// worker processes jobs - simulates CPU-intensive work
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Simulate processing time
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		result := Result{
			JobID:       job.ID,
			Output:      fmt.Sprintf("Processed: %s", job.Data),
			ProcessedAt: time.Now(),
			WorkerID:    id,
		}

		results <- result
	}
}

// =============================================================================
// PIPELINE PATTERN - Uber Style Data Transformation
// =============================================================================

// Pipeline represents a data processing pipeline
type Pipeline struct {
	stages []func(<-chan interface{}) <-chan interface{}
}

// NewPipeline creates a new processing pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// AddStage adds a processing stage to the pipeline
func (p *Pipeline) AddStage(stage func(<-chan interface{}) <-chan interface{}) {
	p.stages = append(p.stages, stage)
}

// Process runs data through the entire pipeline
func (p *Pipeline) Process(input <-chan interface{}) <-chan interface{} {
	current := input

	// Chain all stages together
	for _, stage := range p.stages {
		current = stage(current)
	}

	return current
}

// Example pipeline stages for processing user events
func validateStage(input <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)
		for data := range input {
			// Simulate validation logic
			if str, ok := data.(string); ok && len(str) > 0 {
				output <- fmt.Sprintf("validated_%s", str)
			}
		}
	}()

	return output
}

func enrichStage(input <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)
		for data := range input {
			// Simulate data enrichment
			enriched := fmt.Sprintf("%s_enriched_at_%d", data, time.Now().Unix())
			output <- enriched
		}
	}()

	return output
}

func persistStage(input <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)
		for data := range input {
			// Simulate database persistence
			time.Sleep(10 * time.Millisecond)
			output <- fmt.Sprintf("persisted_%s", data)
		}
	}()

	return output
}

// =============================================================================
// RATE LIMITING PATTERN - Production API Gateway Style
// =============================================================================

// RateLimiter implements token bucket algorithm
type RateLimiter struct {
	tokens chan struct{}
	ticker *time.Ticker
	done   chan bool
}

// NewRateLimiter creates rate limiter with specified rate and burst
func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, burst),
		ticker: time.NewTicker(rate),
		done:   make(chan bool),
	}

	// Fill initial tokens
	for i := 0; i < burst; i++ {
		rl.tokens <- struct{}{}
	}

	// Start token replenishment
	go rl.replenish()

	return rl
}

// Allow checks if request is allowed (non-blocking)
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Wait blocks until request is allowed
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// replenish adds tokens at specified rate
func (rl *RateLimiter) replenish() {
	for {
		select {
		case <-rl.ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Bucket is full, discard token
			}
		case <-rl.done:
			rl.ticker.Stop()
			return
		}
	}
}

// Close stops the rate limiter
func (rl *RateLimiter) Close() {
	close(rl.done)
}

// =============================================================================
// CIRCUIT BREAKER PATTERN - Netflix Hystrix Style
// =============================================================================

// CircuitState represents circuit breaker states
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements circuit breaker pattern
type CircuitBreaker struct {
	mu               sync.RWMutex
	state            CircuitState
	failureCount     int
	successCount     int
	failureThreshold int
	successThreshold int
	timeout          time.Duration
	nextAttempt      time.Time
}

// NewCircuitBreaker creates a circuit breaker
func NewCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
	}
}

// Execute runs function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.canExecute() {
		return fmt.Errorf("circuit breaker is open")
	}

	err := fn()
	cb.recordResult(err == nil)

	return err
}

// canExecute checks if circuit allows execution
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		return time.Now().After(cb.nextAttempt)
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// recordResult updates circuit breaker state based on result
func (cb *CircuitBreaker) recordResult(success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if success {
		cb.successCount++
		cb.failureCount = 0

		if cb.state == StateHalfOpen && cb.successCount >= cb.successThreshold {
			cb.state = StateClosed
			cb.successCount = 0
		}
	} else {
		cb.failureCount++
		cb.successCount = 0

		if cb.failureCount >= cb.failureThreshold {
			cb.state = StateOpen
			cb.nextAttempt = time.Now().Add(cb.timeout)
		}
	}

	// Transition from Open to Half-Open
	if cb.state == StateOpen && time.Now().After(cb.nextAttempt) {
		cb.state = StateHalfOpen
	}
}

// GetState returns current circuit breaker state
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// =============================================================================
// DEMONSTRATION MAIN FUNCTION
// =============================================================================

func main() {
	fmt.Println("🔥 Production Channel Patterns Demo")
	fmt.Println("=====================================")

	// Demo 1: Fan-Out/Fan-In Pattern
	fmt.Println("\n1. Fan-Out/Fan-In Pattern (Netflix Style)")
	jobs := []Job{
		{1, "user_event_1"}, {2, "user_event_2"}, {3, "user_event_3"},
		{4, "user_event_4"}, {5, "user_event_5"}, {6, "user_event_6"},
	}

	start := time.Now()
	results := FanOutFanIn(jobs, 3)
	duration := time.Since(start)

	fmt.Printf("Processed %d jobs in %v using 3 workers:\n", len(results), duration)
	for _, result := range results {
		fmt.Printf("  Job %d -> Worker %d: %s\n", result.JobID, result.WorkerID, result.Output)
	}

	// Demo 2: Pipeline Pattern
	fmt.Println("\n2. Pipeline Pattern (Uber Style)")
	input := make(chan interface{}, 5)

	// Create pipeline
	pipeline := NewPipeline()
	pipeline.AddStage(validateStage)
	pipeline.AddStage(enrichStage)
	pipeline.AddStage(persistStage)

	// Send data through pipeline
	go func() {
		defer close(input)
		for i := 1; i <= 3; i++ {
			input <- fmt.Sprintf("event_%d", i)
		}
	}()

	// Process through pipeline
	output := pipeline.Process(input)
	for result := range output {
		fmt.Printf("  Pipeline result: %s\n", result)
	}

	// Demo 3: Rate Limiting
	fmt.Println("\n3. Rate Limiting Pattern (API Gateway Style)")
	rateLimiter := NewRateLimiter(100*time.Millisecond, 2)
	defer rateLimiter.Close()

	// Test rate limiting
	for i := 1; i <= 5; i++ {
		if rateLimiter.Allow() {
			fmt.Printf("  Request %d: ALLOWED\n", i)
		} else {
			fmt.Printf("  Request %d: RATE LIMITED\n", i)
		}
	}

	// Demo 4: Circuit Breaker
	fmt.Println("\n4. Circuit Breaker Pattern (Netflix Hystrix Style)")
	cb := NewCircuitBreaker(3, 2, 1*time.Second)

	// Simulate failing service
	failingService := func() error {
		if rand.Float32() < 0.7 { // 70% failure rate
			return fmt.Errorf("service failure")
		}
		return nil
	}

	// Test circuit breaker
	for i := 1; i <= 10; i++ {
		err := cb.Execute(failingService)
		state := cb.GetState()

		if err != nil {
			fmt.Printf("  Call %d: FAILED (State: %v) - %v\n", i, state, err)
		} else {
			fmt.Printf("  Call %d: SUCCESS (State: %v)\n", i, state)
		}

		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("\n🎯 These patterns are essential for FAANG L3/L4 interviews!")
	fmt.Println("Master these and you'll handle production concurrency with confidence.")
}
