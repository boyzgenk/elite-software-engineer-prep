// Package concurrency_testing demonstrates advanced testing patterns for concurrent Go code
// This is essential for FAANG interviews where race conditions are commonly tested
package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ===== RACE DETECTION PATTERNS =====

// RaceProneCounter demonstrates a classic race condition
// This pattern is frequently tested in FAANG interviews
type RaceProneCounter struct {
	value int64
	mu    sync.Mutex
}

// UnsafeIncrement creates a race condition (for demonstration)
func (c *RaceProneCounter) UnsafeIncrement() {
	// This will cause data races when run with -race flag
	temp := c.value
	runtime.Gosched() // Force context switch to increase race probability
	c.value = temp + 1
}

// SafeIncrement uses proper synchronization
func (c *RaceProneCounter) SafeIncrement() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// AtomicIncrement uses atomic operations for lock-free increment
func (c *RaceProneCounter) AtomicIncrement() {
	atomic.AddInt64(&c.value, 1)
}

// GetValue safely reads the counter value
func (c *RaceProneCounter) GetValue() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// GetValueAtomic reads using atomic operation
func (c *RaceProneCounter) GetValueAtomic() int64 {
	return atomic.LoadInt64(&c.value)
}

// ===== RACE DETECTION UTILITIES =====

// RaceDetector helps identify race conditions in concurrent code
type RaceDetector struct {
	iterations int
	goroutines int
}

// NewRaceDetector creates a new race detector with specified parameters
func NewRaceDetector(iterations, goroutines int) *RaceDetector {
	return &RaceDetector{
		iterations: iterations,
		goroutines: goroutines,
	}
}

// TestForRaces runs a function concurrently to detect race conditions
func (rd *RaceDetector) TestForRaces(fn func()) bool {
	var wg sync.WaitGroup
	raceDetected := int32(0)

	// Run the test multiple times to increase race detection probability
	for i := 0; i < 100; i++ {
		wg.Add(rd.goroutines)

		for j := 0; j < rd.goroutines; j++ {
			go func() {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						atomic.StoreInt32(&raceDetected, 1)
					}
				}()

				for k := 0; k < rd.iterations; k++ {
					fn()
					runtime.Gosched() // Encourage context switching
				}
			}()
		}

		wg.Wait()

		if atomic.LoadInt32(&raceDetected) == 1 {
			return true
		}
	}

	return false
}

// ===== SHARED RESOURCE RACE PATTERNS =====

// SharedMap demonstrates race conditions with maps
type SharedMap struct {
	data map[string]int
	mu   sync.RWMutex
}

// NewSharedMap creates a new shared map
func NewSharedMap() *SharedMap {
	return &SharedMap{
		data: make(map[string]int),
	}
}

// UnsafeSet creates a race condition with map writes
func (sm *SharedMap) UnsafeSet(key string, value int) {
	// This will panic under race conditions
	sm.data[key] = value
}

// SafeSet uses proper locking for map writes
func (sm *SharedMap) SafeSet(key string, value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

// SafeGet uses read lock for map reads
func (sm *SharedMap) SafeGet(key string) (int, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, exists := sm.data[key]
	return value, exists
}

// ===== CHANNEL RACE PATTERNS =====

// ChannelRaceDemo demonstrates race conditions with channels
type ChannelRaceDemo struct {
	ch chan int
}

// NewChannelRaceDemo creates a new channel race demo
func NewChannelRaceDemo() *ChannelRaceDemo {
	return &ChannelRaceDemo{
		ch: make(chan int, 10),
	}
}

// RaceProneProducer creates potential race conditions
func (crd *ChannelRaceDemo) RaceProneProducer(ctx context.Context, data []int) {
	for _, value := range data {
		select {
		case crd.ch <- value:
			// Race condition: channel might be closed by another goroutine
		case <-ctx.Done():
			return
		}
	}
}

// SafeProducer handles channel operations safely
func (crd *ChannelRaceDemo) SafeProducer(ctx context.Context, data []int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	for _, value := range data {
		select {
		case crd.ch <- value:
		case <-ctx.Done():
			return
		}
	}
}

// ===== BENCHMARK RACE DETECTION =====

// BenchmarkRaceDetection compares performance of different synchronization methods
func BenchmarkRaceDetection(b *testing.B) {
	counter := &RaceProneCounter{}

	b.Run("Unsafe", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				counter.UnsafeIncrement()
			}
		})
	})

	b.Run("Mutex", func(b *testing.B) {
		counter.value = 0
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				counter.SafeIncrement()
			}
		})
	})

	b.Run("Atomic", func(b *testing.B) {
		counter.value = 0
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				counter.AtomicIncrement()
			}
		})
	})
}

// ===== ADVANCED RACE DETECTION PATTERNS =====

// SliceRaceDemo demonstrates race conditions with slices
type SliceRaceDemo struct {
	data []int
	mu   sync.Mutex
}

// NewSliceRaceDemo creates a new slice race demo
func NewSliceRaceDemo() *SliceRaceDemo {
	return &SliceRaceDemo{
		data: make([]int, 0),
	}
}

// UnsafeAppend creates race condition with slice growth
func (srd *SliceRaceDemo) UnsafeAppend(value int) {
	// Race condition: slice header modification is not atomic
	srd.data = append(srd.data, value)
}

// SafeAppend uses proper synchronization
func (srd *SliceRaceDemo) SafeAppend(value int) {
	srd.mu.Lock()
	defer srd.mu.Unlock()
	srd.data = append(srd.data, value)
}

// ===== TESTING UTILITIES =====

// TestRaceConditions demonstrates how to test for race conditions
func TestRaceConditions(t *testing.T) {
	// Test counter races
	t.Run("CounterRaces", func(t *testing.T) {
		counter := &RaceProneCounter{}
		detector := NewRaceDetector(1000, 10)

		// This should detect races in unsafe increment
		hasRace := detector.TestForRaces(func() {
			counter.UnsafeIncrement()
		})

		if !hasRace {
			t.Log("Race condition not detected (may need more iterations)")
		}
	})

	// Test map races
	t.Run("MapRaces", func(t *testing.T) {
		sharedMap := NewSharedMap()

		// This will panic with concurrent map writes
		go func() {
			for i := 0; i < 1000; i++ {
				sharedMap.UnsafeSet(fmt.Sprintf("key%d", i), i)
			}
		}()

		go func() {
			for i := 0; i < 1000; i++ {
				sharedMap.UnsafeSet(fmt.Sprintf("key%d", i+1000), i)
			}
		}()

		time.Sleep(100 * time.Millisecond)
		// The test might panic before reaching here
	})
}

// ===== INTERVIEW SIMULATION PATTERNS =====

// InterviewRaceQuestion simulates common FAANG interview race condition questions
type InterviewRaceQuestion struct {
	name string
	code func()
}

// GetCommonRaceQuestions returns typical race condition interview questions
func GetCommonRaceQuestions() []InterviewRaceQuestion {
	questions := []InterviewRaceQuestion{
		{
			name: "Bank Account Transfer Race",
			code: func() {
				// Simulate bank account transfer with race conditions
				balance := int64(1000)

				transfer := func(amount int64) {
					// Race condition: read-modify-write is not atomic
					current := atomic.LoadInt64(&balance)
					runtime.Gosched()
					atomic.StoreInt64(&balance, current-amount)
				}

				var wg sync.WaitGroup
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						transfer(50)
					}()
				}
				wg.Wait()

				fmt.Printf("Final balance: %d (should be 500)\n", balance)
			},
		},
		{
			name: "Producer-Consumer Race",
			code: func() {
				buffer := make([]int, 0)
				var mu sync.Mutex

				// Producer
				go func() {
					for i := 0; i < 100; i++ {
						mu.Lock()
						buffer = append(buffer, i)
						mu.Unlock()
						runtime.Gosched()
					}
				}()

				// Consumer
				go func() {
					for len(buffer) < 100 {
						mu.Lock()
						if len(buffer) > 0 {
							buffer = buffer[1:]
						}
						mu.Unlock()
						runtime.Gosched()
					}
				}()

				time.Sleep(100 * time.Millisecond)
				fmt.Printf("Buffer length: %d\n", len(buffer))
			},
		},
	}

	return questions
}

// ===== MAIN DEMONSTRATION =====

func main() {
	fmt.Println("=== Go Concurrency Race Detection Patterns ===")
	fmt.Println("Run with: go run -race race_detection.go")
	fmt.Println()

	// Demonstrate race detection
	fmt.Println("1. Testing Race-Prone Counter:")
	counter := &RaceProneCounter{}

	var wg sync.WaitGroup
	numGoroutines := 10
	iterations := 1000

	// Test unsafe increment (will show races with -race flag)
	fmt.Println("   Running unsafe increment...")
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.UnsafeIncrement()
			}
		}()
	}
	wg.Wait()

	fmt.Printf("   Unsafe counter result: %d (expected: %d)\n",
		counter.GetValue(), numGoroutines*iterations)

	// Reset and test safe increment
	counter.value = 0
	fmt.Println("   Running safe increment...")
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.SafeIncrement()
			}
		}()
	}
	wg.Wait()

	fmt.Printf("   Safe counter result: %d (expected: %d)\n",
		counter.GetValue(), numGoroutines*iterations)

	// Test shared map races
	fmt.Println("\n2. Testing Shared Map Races:")
	sharedMap := NewSharedMap()

	// Safe operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d_%d", id, j)
				sharedMap.SafeSet(key, j)

				if value, exists := sharedMap.SafeGet(key); exists {
					_ = value // Use the value to prevent optimization
				}
			}
		}(i)
	}
	wg.Wait()

	fmt.Println("   Shared map operations completed safely")

	// Demonstrate interview questions
	fmt.Println("\n3. Running Interview Simulation Questions:")
	questions := GetCommonRaceQuestions()

	for i, question := range questions {
		fmt.Printf("   Question %d: %s\n", i+1, question.name)
		question.code()
	}

	fmt.Println("\n=== Race Detection Complete ===")
	fmt.Println("Remember to run with -race flag to detect race conditions!")
	fmt.Println("Example: go run -race race_detection.go")
}
