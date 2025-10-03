// Rate Limiting Pattern Implementation
// Used by: Netflix (API throttling), Twitter (tweet rate limits), AWS (request throttling)
// Problem: Control request/operation rates to prevent system overload
// Solution: Token bucket and sliding window algorithms with Go channels

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TokenBucketLimiter implements token bucket algorithm
// Used by Netflix for API rate limiting and AWS for request throttling
type TokenBucketLimiter struct {
	tokens   chan struct{}
	capacity int
	refill   time.Duration
	ticker   *time.Ticker
	done     chan struct{}
	once     sync.Once
}

// NewTokenBucketLimiter creates a new token bucket rate limiter
// capacity: maximum tokens in bucket
// refillRate: how often to add tokens (e.g., time.Second/10 = 10 tokens per second)
func NewTokenBucketLimiter(capacity int, refillRate time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		tokens:   make(chan struct{}, capacity),
		capacity: capacity,
		refill:   refillRate,
		ticker:   time.NewTicker(refillRate),
		done:     make(chan struct{}),
	}

	// Fill bucket initially
	for i := 0; i < capacity; i++ {
		limiter.tokens <- struct{}{}
	}

	// Start refill goroutine
	go limiter.refillTokens()

	return limiter
}

// refillTokens adds tokens to bucket at regular intervals
func (t *TokenBucketLimiter) refillTokens() {
	for {
		select {
		case <-t.ticker.C:
			// Non-blocking send - don't block if bucket is full
			select {
			case t.tokens <- struct{}{}:
				// Token added successfully
			default:
				// Bucket is full, skip this refill
			}
		case <-t.done:
			t.ticker.Stop()
			return
		}
	}
}

// Allow attempts to acquire a token
// Returns true if token acquired, false if rate limited
func (t *TokenBucketLimiter) Allow() bool {
	select {
	case <-t.tokens:
		return true
	default:
		return false
	}
}

// Wait waits until a token is available
func (t *TokenBucketLimiter) Wait(ctx context.Context) error {
	select {
	case <-t.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close stops the rate limiter
func (t *TokenBucketLimiter) Close() {
	t.once.Do(func() {
		close(t.done)
	})
}

// SlidingWindowLimiter implements sliding window algorithm
// Used by Twitter for tweet rate limiting and GitHub for API requests
type SlidingWindowLimiter struct {
	mu          sync.Mutex
	requests    []time.Time
	maxRequests int
	window      time.Duration
}

// NewSlidingWindowLimiter creates a new sliding window rate limiter
// maxRequests: maximum requests allowed in the time window
// window: time window duration (e.g., time.Minute for per-minute limiting)
func NewSlidingWindowLimiter(maxRequests int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		requests:    make([]time.Time, 0),
		maxRequests: maxRequests,
		window:      window,
	}
}

// Allow checks if request is allowed under sliding window
func (s *SlidingWindowLimiter) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-s.window)

	// Remove expired requests
	var validRequests []time.Time
	for _, reqTime := range s.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	s.requests = validRequests

	// Check if we can allow this request
	if len(s.requests) < s.maxRequests {
		s.requests = append(s.requests, now)
		return true
	}

	return false
}

// RateLimitedWorker demonstrates rate-limited request processing
// Like Netflix's API gateway handling user requests
func RateLimitedWorker(limiter *TokenBucketLimiter, requests <-chan string, results chan<- string) {
	for req := range requests {
		// Try to acquire token
		if limiter.Allow() {
			// Process request
			time.Sleep(50 * time.Millisecond) // Simulate processing
			results <- fmt.Sprintf("✅ Processed: %s", req)
		} else {
			// Rate limited
			results <- fmt.Sprintf("🚫 Rate limited: %s", req)
		}
	}
}

// AdaptiveRateLimiter adjusts rate based on system load
// Similar to AWS Auto Scaling and Netflix's adaptive throttling
type AdaptiveRateLimiter struct {
	limiter     *TokenBucketLimiter
	mu          sync.RWMutex
	successRate float64
	adjustment  time.Duration
}

// NewAdaptiveRateLimiter creates adaptive rate limiter
func NewAdaptiveRateLimiter(initialCapacity int, initialRate time.Duration) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		limiter:    NewTokenBucketLimiter(initialCapacity, initialRate),
		adjustment: initialRate,
	}
}

// AdjustRate modifies rate based on success/failure feedback
func (a *AdaptiveRateLimiter) AdjustRate(success bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if success {
		a.successRate = a.successRate*0.9 + 0.1 // Exponential moving average
		if a.successRate > 0.8 {
			// Increase rate if success rate is high
			a.adjustment = time.Duration(float64(a.adjustment) * 0.9)
		}
	} else {
		a.successRate = a.successRate * 0.9 // No success bonus
		if a.successRate < 0.5 {
			// Decrease rate if success rate is low
			a.adjustment = time.Duration(float64(a.adjustment) * 1.1)
		}
	}

	// Recreate limiter with new rate
	a.limiter.Close()
	a.limiter = NewTokenBucketLimiter(10, a.adjustment)
}

// Allow checks if request is allowed
func (a *AdaptiveRateLimiter) Allow() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.limiter.Allow()
}

func main() {
	fmt.Println("🚀 Rate Limiting Patterns Demo")
	fmt.Println("=====================================")

	// Demo 1: Token Bucket Rate Limiter
	fmt.Println("\n1️⃣ Token Bucket Rate Limiter (Netflix-style)")
	tokenLimiter := NewTokenBucketLimiter(5, time.Second/2) // 2 tokens per second, burst of 5
	defer tokenLimiter.Close()

	// Create request channels
	requests := make(chan string, 20)
	results := make(chan string, 20)

	// Start rate-limited worker
	go RateLimitedWorker(tokenLimiter, requests, results)

	// Send burst of requests
	for i := 1; i <= 10; i++ {
		requests <- fmt.Sprintf("Request-%d", i)
	}
	close(requests)

	// Collect results
	for i := 0; i < 10; i++ {
		fmt.Println(<-results)
	}

	// Demo 2: Sliding Window Rate Limiter
	fmt.Println("\n2️⃣ Sliding Window Rate Limiter (Twitter-style)")
	windowLimiter := NewSlidingWindowLimiter(3, 2*time.Second) // 3 requests per 2 seconds

	for i := 1; i <= 8; i++ {
		if windowLimiter.Allow() {
			fmt.Printf("✅ Request %d: Allowed\n", i)
		} else {
			fmt.Printf("🚫 Request %d: Rate limited\n", i)
		}
		time.Sleep(300 * time.Millisecond)
	}

	// Demo 3: Rate Limiting with Context and Graceful Handling
	fmt.Println("\n3️⃣ Context-Aware Rate Limiting")
	contextLimiter := NewTokenBucketLimiter(2, time.Second) // 1 token per second, burst of 2
	defer contextLimiter.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if err := contextLimiter.Wait(ctx); err != nil {
				fmt.Printf("⏰ Request %d: Timeout (%v)\n", id, err)
				return
			}

			fmt.Printf("✅ Request %d: Processed successfully\n", id)
		}(i)
	}

	wg.Wait()

	// Demo 4: Adaptive Rate Limiting
	fmt.Println("\n4️⃣ Adaptive Rate Limiting (AWS-style)")
	adaptiveLimiter := NewAdaptiveRateLimiter(5, time.Second/3)

	// Simulate varying success rates
	successPatterns := []bool{true, true, false, false, false, true, true, true, true, false}

	for i, success := range successPatterns {
		if adaptiveLimiter.Allow() {
			fmt.Printf("✅ Request %d: Allowed", i+1)
			if success {
				fmt.Printf(" (Success)\n")
			} else {
				fmt.Printf(" (Failed)\n")
			}
			adaptiveLimiter.AdjustRate(success)
		} else {
			fmt.Printf("🚫 Request %d: Rate limited\n", i+1)
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("\n🎯 Rate Limiting Patterns Complete!")
	fmt.Println("Key Takeaways:")
	fmt.Println("• Token bucket: Good for bursty traffic with sustained rate")
	fmt.Println("• Sliding window: Precise time-based limiting")
	fmt.Println("• Context-aware: Graceful handling with timeouts")
	fmt.Println("• Adaptive: Responds to system feedback")
}
