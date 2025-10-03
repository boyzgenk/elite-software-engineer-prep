// Package patterns demonstrates production-grade Go concurrency patterns
// Used in Netflix, Uber, Google scale backend systems

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("🔥 Production Go Concurrency Patterns")
	fmt.Println("======================================")

	// Pattern 1: Fan-In/Fan-Out for parallel processing
	demonstrateFanInFanOut()

	// Pattern 2: Worker Pool with graceful shutdown
	demonstrateWorkerPool()

	// Pattern 3: Circuit Breaker pattern
	demonstrateCircuitBreaker()

	// Pattern 4: Rate Limiting with Token Bucket
	demonstrateRateLimiting()

	// Pattern 5: Pipeline Processing
	demonstratePipeline()
}

// Pattern 1: Fan-In/Fan-Out - Distribute work across multiple goroutines
func demonstrateFanInFanOut() {
	fmt.Println("\n🚀 Fan-In/Fan-Out Pattern")
	fmt.Println("-------------------------")

	// Input data
	numbers := make(chan int, 100)
	for i := 1; i <= 100; i++ {
		numbers <- i
	}
	close(numbers)

	// Fan-Out: Distribute work to multiple workers
	const numWorkers = 5
	results := make([]<-chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		output := make(chan int, 20)
		results[i] = output

		// Each worker processes different numbers
		go func(input <-chan int, output chan<- int, workerID int) {
			defer close(output)
			processed := 0

			for num := range input {
				// Simulate expensive computation
				result := num * num
				output <- result
				processed++
			}

			fmt.Printf("Worker %d processed %d numbers\n", workerID, processed)
		}(numbers, output, i)
	}

	// Fan-In: Merge results from multiple workers
	merged := fanIn(results...)

	// Collect all results
	var total int
	count := 0
	for result := range merged {
		total += result
		count++
	}

	fmt.Printf("Processed %d results, sum: %d\n", count, total)
}

// fanIn merges multiple channels into one
func fanIn(inputs ...<-chan int) <-chan int {
	output := make(chan int)
	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, input := range inputs {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for value := range ch {
				output <- value
			}
		}(input)
	}

	// Close output when all inputs are exhausted
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// Pattern 2: Worker Pool with Graceful Shutdown
func demonstrateWorkerPool() {
	fmt.Println("\n⚡ Worker Pool Pattern")
	fmt.Println("---------------------")

	const numWorkers = 3
	const numJobs = 20

	// Create worker pool
	pool := NewWorkerPool(numWorkers)
	pool.Start()

	// Submit jobs
	for i := 1; i <= numJobs; i++ {
		job := Job{
			ID:   i,
			Data: fmt.Sprintf("Task %d", i),
		}
		pool.Submit(job)
	}

	// Graceful shutdown
	fmt.Println("Initiating graceful shutdown...")
	pool.Stop()
	fmt.Println("Worker pool shutdown complete")
}

type Job struct {
	ID   int
	Data string
}

type WorkerPool struct {
	jobs     chan Job
	shutdown chan struct{}
	wg       sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		jobs:     make(chan Job, 100),
		shutdown: make(chan struct{}),
	}

	// Start workers
	for i := 0; i < numWorkers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case job := <-p.jobs:
			// Process job
			fmt.Printf("Worker %d processing job %d: %s\n", id, job.ID, job.Data)
			time.Sleep(time.Duration(rand.IntN(500)) * time.Millisecond)
			fmt.Printf("Worker %d completed job %d\n", id, job.ID)

		case <-p.shutdown:
			fmt.Printf("Worker %d shutting down\n", id)
			return
		}
	}
}

func (p *WorkerPool) Start() {
	// Workers already started in NewWorkerPool
}

func (p *WorkerPool) Submit(job Job) {
	select {
	case p.jobs <- job:
	case <-p.shutdown:
		fmt.Printf("Cannot submit job %d: pool is shutting down\n", job.ID)
	}
}

func (p *WorkerPool) Stop() {
	close(p.shutdown)
	p.wg.Wait()
}

// Pattern 3: Circuit Breaker
func demonstrateCircuitBreaker() {
	fmt.Println("\n🔧 Circuit Breaker Pattern")
	fmt.Println("--------------------------")

	// Simulate flaky service
	flakyService := func() error {
		if rand.Float64() < 0.7 { // 70% failure rate
			return fmt.Errorf("service unavailable")
		}
		return nil
	}

	cb := NewCircuitBreaker(3, 2*time.Second) // 3 failures, 2 second timeout

	// Test circuit breaker behavior
	for i := 0; i < 15; i++ {
		err := cb.Call(flakyService)
		status := "SUCCESS"
		if err != nil {
			status = fmt.Sprintf("FAILED: %v", err)
		}

		fmt.Printf("Call %d: %s (State: %s)\n", i+1, status, cb.State())
		time.Sleep(300 * time.Millisecond)
	}
}

type CircuitState int

const (
	Closed CircuitState = iota
	Open
	HalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case Closed:
		return "CLOSED"
	case Open:
		return "OPEN"
	case HalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}

type CircuitBreaker struct {
	maxFailures  int
	timeout      time.Duration
	failures     int64
	lastFailTime time.Time
	state        CircuitState
	mutex        sync.RWMutex
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures: maxFailures,
		timeout:     timeout,
		state:       Closed,
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Check if we should transition from Open to HalfOpen
	if cb.state == Open && time.Since(cb.lastFailTime) > cb.timeout {
		cb.state = HalfOpen
		atomic.StoreInt64(&cb.failures, 0)
	}

	// Reject calls if circuit is open
	if cb.state == Open {
		return fmt.Errorf("circuit breaker is open")
	}

	// Execute the function
	err := fn()

	if err != nil {
		// Increment failure count
		failures := atomic.AddInt64(&cb.failures, 1)
		cb.lastFailTime = time.Now()

		// Trip circuit if max failures reached
		if int(failures) >= cb.maxFailures {
			cb.state = Open
		}

		return err
	}

	// Success - reset circuit
	if cb.state == HalfOpen {
		cb.state = Closed
	}
	atomic.StoreInt64(&cb.failures, 0)

	return nil
}

func (cb *CircuitBreaker) State() CircuitState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// Pattern 4: Rate Limiting with Token Bucket
func demonstrateRateLimiting() {
	fmt.Println("\n⏱️  Rate Limiting Pattern")
	fmt.Println("-----------------------")

	// Create rate limiter: 5 requests per second, burst of 10
	limiter := NewTokenBucket(5, 10)

	// Simulate rapid requests
	for i := 0; i < 20; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: ALLOWED\n", i+1)
		} else {
			fmt.Printf("Request %d: RATE LIMITED\n", i+1)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

type TokenBucket struct {
	tokens     int64
	capacity   int64
	rate       int64
	lastRefill time.Time
	mutex      sync.Mutex
}

func NewTokenBucket(rate, capacity int64) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		rate:       rate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	// Add tokens based on elapsed time
	tokensToAdd := int64(elapsed.Seconds()) * tb.rate
	if tokensToAdd > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}

	// Check if request can be allowed
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// Pattern 5: Pipeline Processing
func demonstratePipeline() {
	fmt.Println("\n🔄 Pipeline Processing Pattern")
	fmt.Println("------------------------------")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stage 1: Generate numbers
	numbers := generateNumbers(ctx, 1, 20)

	// Stage 2: Square the numbers
	squared := squareNumbers(ctx, numbers)

	// Stage 3: Filter even numbers
	evenOnly := filterEven(ctx, squared)

	// Stage 4: Sum the results
	sum := sumNumbers(ctx, evenOnly)

	fmt.Printf("Pipeline result: %d\n", sum)
}

func generateNumbers(ctx context.Context, start, end int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for i := start; i <= end; i++ {
			select {
			case output <- i:
				time.Sleep(50 * time.Millisecond) // Simulate work
			case <-ctx.Done():
				return
			}
		}
	}()

	return output
}

func squareNumbers(ctx context.Context, input <-chan int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for num := range input {
			select {
			case output <- num * num:
			case <-ctx.Done():
				return
			}
		}
	}()

	return output
}

func filterEven(ctx context.Context, input <-chan int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for num := range input {
			if num%2 == 0 {
				select {
				case output <- num:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return output
}

func sumNumbers(ctx context.Context, input <-chan int) int {
	sum := 0
	for num := range input {
		sum += num
		select {
		case <-ctx.Done():
			log.Printf("Sum calculation cancelled")
			return sum
		default:
		}
	}
	return sum
}

/*
🔥 FAANG Interview Key Points:

1. Fan-In/Fan-Out:
   - Distribute work across multiple goroutines (fan-out)
   - Merge results from multiple sources (fan-in)
   - Essential for parallel processing at scale

2. Worker Pool:
   - Fixed number of workers processing jobs from queue
   - Graceful shutdown with sync.WaitGroup
   - Prevents goroutine explosion under load

3. Circuit Breaker:
   - Prevents cascading failures in distributed systems
   - States: Closed (normal), Open (failing), Half-Open (testing)
   - Used in microservices for fault tolerance

4. Rate Limiting:
   - Token bucket algorithm for controlling request rate
   - Protects services from being overwhelmed
   - Critical for API design and SLA compliance

5. Pipeline Processing:
   - Chain of processing stages connected by channels
   - Context cancellation for proper cleanup
   - Efficient memory usage with streaming data

These patterns are essential for:
- Building resilient microservices
- Handling high-concurrency scenarios
- System design interviews
- Production Go backend systems

Real-world usage:
- Netflix: Circuit breakers in microservice architecture
- Uber: Rate limiting for API protection
- Google: Pipeline processing for data streams
*/
