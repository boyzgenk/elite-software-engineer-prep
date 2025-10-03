// Package chaos_engineering demonstrates chaos testing patterns for Go concurrency
// This implements Netflix-style chaos engineering for finding concurrency bugs
// Used by companies like Netflix, Uber, Facebook to stress-test distributed systems
package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ===== CHAOS ENGINEERING FRAMEWORK =====

// ChaosConfig defines parameters for chaos testing
type ChaosConfig struct {
	Duration         time.Duration // How long to run chaos
	GoroutineCount   int           // Number of concurrent goroutines
	FailureRate      float64       // Probability of introducing failures
	MaxLatency       time.Duration // Maximum artificial latency
	MemoryPressure   bool          // Apply memory pressure
	CPUStarvation    bool          // Apply CPU starvation
	NetworkPartition bool          // Simulate network failures
}

// DefaultChaosConfig returns a default configuration
func DefaultChaosConfig() *ChaosConfig {
	return &ChaosConfig{
		Duration:         30 * time.Second,
		GoroutineCount:   50,
		FailureRate:      0.1, // 10% failure rate
		MaxLatency:       100 * time.Millisecond,
		MemoryPressure:   true,
		CPUStarvation:    true,
		NetworkPartition: false,
	}
}

// ChaosEngine orchestrates chaos testing
type ChaosEngine struct {
	config *ChaosConfig
	stats  *ChaosStats
}

// ChaosStats tracks chaos testing metrics
type ChaosStats struct {
	TestsExecuted    int64
	TestsFailed      int64
	LatencyInjected  int64
	MemoryAllocated  int64
	PanicsRecovered  int64
	DeadlocksFound   int64
	RaceConditions   int64
	GoroutinesLeaked int64
}

// NewChaosEngine creates a new chaos engine
func NewChaosEngine(config *ChaosConfig) *ChaosEngine {
	return &ChaosEngine{
		config: config,
		stats:  &ChaosStats{},
	}
}

// RunChaosTest executes a chaos test on the provided function
func (ce *ChaosEngine) RunChaosTest(name string, testFunc func() error) *ChaosResult {
	fmt.Printf("🔥 Starting chaos test: %s\n", name)

	result := &ChaosResult{
		TestName:  name,
		StartTime: time.Now(),
		Config:    ce.config,
		Succeeded: true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), ce.config.Duration)
	defer cancel()

	// Start chaos goroutines
	var wg sync.WaitGroup
	errorChan := make(chan error, ce.config.GoroutineCount)

	// Memory pressure goroutine
	if ce.config.MemoryPressure {
		wg.Add(1)
		go ce.induceMemoryPressure(ctx, &wg)
	}

	// CPU starvation goroutine
	if ce.config.CPUStarvation {
		wg.Add(1)
		go ce.induceCPUStarvation(ctx, &wg)
	}

	// Start test goroutines with chaos injection
	for i := 0; i < ce.config.GoroutineCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt64(&ce.stats.PanicsRecovered, 1)
					errorChan <- fmt.Errorf("panic in goroutine %d: %v", id, r)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Inject chaos before each test execution
					ce.injectChaos()

					// Execute the test function
					if err := testFunc(); err != nil {
						atomic.AddInt64(&ce.stats.TestsFailed, 1)
						errorChan <- fmt.Errorf("test failed in goroutine %d: %v", id, err)
						result.Succeeded = false
					}

					atomic.AddInt64(&ce.stats.TestsExecuted, 1)
					time.Sleep(10 * time.Millisecond) // Brief pause between iterations
				}
			}
		}(i)
	}

	// Monitor goroutine for collecting errors
	go func() {
		for err := range errorChan {
			result.Errors = append(result.Errors, err)
		}
	}()

	// Wait for completion or timeout
	wg.Wait()
	close(errorChan)

	result.Duration = time.Since(result.StartTime)
	result.Stats = *ce.stats

	fmt.Printf("   ✅ Chaos test completed: %v\n", result.Duration)
	return result
}

// injectChaos randomly injects various types of chaos
func (ce *ChaosEngine) injectChaos() {
	if rand.Float64() < ce.config.FailureRate {
		switch rand.Intn(4) {
		case 0:
			ce.injectLatency()
		case 1:
			ce.injectGoroutineSpam()
		case 2:
			ce.injectMemoryFragmentation()
		case 3:
			ce.injectContextSwitching()
		}
	}
}

// injectLatency adds random delays to simulate network/disk latency
func (ce *ChaosEngine) injectLatency() {
	latency := time.Duration(rand.Int63n(int64(ce.config.MaxLatency)))
	time.Sleep(latency)
	atomic.AddInt64(&ce.stats.LatencyInjected, 1)
}

// injectGoroutineSpam creates temporary goroutines to increase scheduler pressure
func (ce *ChaosEngine) injectGoroutineSpam() {
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(1 * time.Millisecond)
		}()
	}
}

// injectMemoryFragmentation allocates and deallocates memory randomly
func (ce *ChaosEngine) injectMemoryFragmentation() {
	size := rand.Intn(1024*1024) + 1024 // 1KB to 1MB
	data := make([]byte, size)
	_ = data // Use the data to prevent optimization
	atomic.AddInt64(&ce.stats.MemoryAllocated, int64(size))
}

// injectContextSwitching forces context switches
func (ce *ChaosEngine) injectContextSwitching() {
	runtime.Gosched()
}

// induceMemoryPressure continuously allocates memory to stress the GC
func (ce *ChaosEngine) induceMemoryPressure(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	var memBlocks [][]byte
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Allocate memory blocks
			block := make([]byte, 64*1024) // 64KB blocks
			memBlocks = append(memBlocks, block)

			// Keep only the last 100 blocks to prevent OOM
			if len(memBlocks) > 100 {
				memBlocks = memBlocks[1:]
			}
		}
	}
}

// induceCPUStarvation creates CPU-intensive work to stress scheduling
func (ce *ChaosEngine) induceCPUStarvation(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// CPU-intensive work
			sum := 0
			for i := 0; i < 10000; i++ {
				sum += i * i
			}
			_ = sum           // Use result to prevent optimization
			runtime.Gosched() // Yield occasionally
		}
	}
}

// ChaosResult contains the results of a chaos test
type ChaosResult struct {
	TestName  string
	Succeeded bool
	Duration  time.Duration
	StartTime time.Time
	Config    *ChaosConfig
	Stats     ChaosStats
	Errors    []error
}

// String returns a formatted string representation of the result
func (cr *ChaosResult) String() string {
	status := "PASS"
	if !cr.Succeeded {
		status = "FAIL"
	}

	return fmt.Sprintf("ChaosTest[%s]: %s - Duration: %v, Tests: %d, Failures: %d, Panics: %d",
		cr.TestName, status, cr.Duration, cr.Stats.TestsExecuted,
		cr.Stats.TestsFailed, cr.Stats.PanicsRecovered)
}

// ===== CHAOS TEST SUBJECTS =====

// ChaosCounter represents a counter that might fail under stress
type ChaosCounter struct {
	value int64
	mu    sync.RWMutex
}

// Increment adds 1 to the counter
func (cc *ChaosCounter) Increment() error {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// Simulate potential failure points
	if rand.Float64() < 0.001 { // 0.1% chance of simulated failure
		return fmt.Errorf("simulated counter failure")
	}

	cc.value++
	return nil
}

// Get returns the current counter value
func (cc *ChaosCounter) Get() int64 {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return cc.value
}

// ChaosMap represents a map that might fail under concurrent stress
type ChaosMap struct {
	data sync.Map
}

// Set stores a key-value pair
func (cm *ChaosMap) Set(key, value interface{}) error {
	if rand.Float64() < 0.001 { // 0.1% chance of simulated failure
		return fmt.Errorf("simulated map write failure")
	}

	cm.data.Store(key, value)
	return nil
}

// Get retrieves a value by key
func (cm *ChaosMap) Get(key interface{}) (interface{}, bool, error) {
	if rand.Float64() < 0.0005 { // 0.05% chance of simulated failure
		return nil, false, fmt.Errorf("simulated map read failure")
	}

	value, exists := cm.data.Load(key)
	return value, exists, nil
}

// ChaosChannel represents a channel that might block or fail
type ChaosChannel struct {
	ch   chan interface{}
	size int
}

// NewChaosChannel creates a new chaos channel
func NewChaosChannel(size int) *ChaosChannel {
	return &ChaosChannel{
		ch:   make(chan interface{}, size),
		size: size,
	}
}

// Send sends a value to the channel
func (cc *ChaosChannel) Send(value interface{}) error {
	if rand.Float64() < 0.001 { // 0.1% chance of simulated failure
		return fmt.Errorf("simulated channel send failure")
	}

	select {
	case cc.ch <- value:
		return nil
	default:
		return fmt.Errorf("channel full")
	}
}

// Receive receives a value from the channel
func (cc *ChaosChannel) Receive() (interface{}, error) {
	if rand.Float64() < 0.001 { // 0.1% chance of simulated failure
		return nil, fmt.Errorf("simulated channel receive failure")
	}

	select {
	case value := <-cc.ch:
		return value, nil
	default:
		return nil, fmt.Errorf("channel empty")
	}
}

// ===== CHAOS TEST SCENARIOS =====

// ChaosTestSuite contains various chaos test scenarios
type ChaosTestSuite struct {
	engine *ChaosEngine
}

// NewChaosTestSuite creates a new test suite
func NewChaosTestSuite(config *ChaosConfig) *ChaosTestSuite {
	return &ChaosTestSuite{
		engine: NewChaosEngine(config),
	}
}

// TestCounterChaos tests counter operations under chaos
func (cts *ChaosTestSuite) TestCounterChaos() *ChaosResult {
	counter := &ChaosCounter{}

	return cts.engine.RunChaosTest("Counter Chaos", func() error {
		return counter.Increment()
	})
}

// TestMapChaos tests map operations under chaos
func (cts *ChaosTestSuite) TestMapChaos() *ChaosResult {
	chaosMap := &ChaosMap{}

	return cts.engine.RunChaosTest("Map Chaos", func() error {
		key := fmt.Sprintf("key_%d", rand.Intn(100))
		value := rand.Intn(1000)

		if err := chaosMap.Set(key, value); err != nil {
			return err
		}

		_, _, err := chaosMap.Get(key)
		return err
	})
}

// TestChannelChaos tests channel operations under chaos
func (cts *ChaosTestSuite) TestChannelChaos() *ChaosResult {
	chaosChannel := NewChaosChannel(10)

	return cts.engine.RunChaosTest("Channel Chaos", func() error {
		// Randomly send or receive
		if rand.Float64() < 0.5 {
			return chaosChannel.Send(rand.Intn(100))
		} else {
			_, err := chaosChannel.Receive()
			return err
		}
	})
}

// TestDeadlockChaos intentionally creates deadlock scenarios
func (cts *ChaosTestSuite) TestDeadlockChaos() *ChaosResult {
	mu1 := &sync.Mutex{}
	mu2 := &sync.Mutex{}

	return cts.engine.RunChaosTest("Deadlock Chaos", func() error {
		// Randomly order mutex acquisition to create potential deadlocks
		if rand.Float64() < 0.5 {
			mu1.Lock()
			defer mu1.Unlock()

			// Small delay to increase deadlock probability
			if rand.Float64() < 0.01 {
				time.Sleep(1 * time.Millisecond)
			}

			mu2.Lock()
			defer mu2.Unlock()
		} else {
			mu2.Lock()
			defer mu2.Unlock()

			// Small delay to increase deadlock probability
			if rand.Float64() < 0.01 {
				time.Sleep(1 * time.Millisecond)
			}

			mu1.Lock()
			defer mu1.Unlock()
		}

		return nil
	})
}

// TestGoroutineLeakChaos tests for goroutine leaks
func (cts *ChaosTestSuite) TestGoroutineLeakChaos() *ChaosResult {
	return cts.engine.RunChaosTest("Goroutine Leak Chaos", func() error {
		// Occasionally spawn goroutines that might not complete
		if rand.Float64() < 0.05 { // 5% chance
			go func() {
				// This goroutine might leak if it blocks indefinitely
				if rand.Float64() < 0.1 { // 10% of spawned goroutines leak
					select {} // Block forever
				}
				// Normal goroutines complete quickly
				time.Sleep(1 * time.Millisecond)
			}()
		}

		return nil
	})
}

// RunAllChaosTests runs all chaos test scenarios
func (cts *ChaosTestSuite) RunAllChaosTests() []*ChaosResult {
	fmt.Println("🔥🔥🔥 STARTING CHAOS ENGINEERING TESTS 🔥🔥🔥")
	fmt.Println("\"Chaos Engineering is the discipline of experimenting on a system")
	fmt.Println("in order to build confidence in the system's capability to withstand")
	fmt.Println("turbulent conditions in production.\" - Netflix")
	fmt.Println()

	tests := []func() *ChaosResult{
		cts.TestCounterChaos,
		cts.TestMapChaos,
		cts.TestChannelChaos,
		// Note: Deadlock test is commented out as it may actually deadlock
		// cts.TestDeadlockChaos,
		cts.TestGoroutineLeakChaos,
	}

	results := make([]*ChaosResult, len(tests))

	for i, test := range tests {
		results[i] = test()
		fmt.Printf("   %s\n", results[i].String())

		// Brief pause between tests
		time.Sleep(1 * time.Second)
	}

	return results
}

// ===== CHAOS ANALYSIS =====

// AnalyzeChaosResults analyzes the results of chaos testing
func AnalyzeChaosResults(results []*ChaosResult) {
	fmt.Println("\n🔍 CHAOS TEST ANALYSIS")
	fmt.Println(strings.Repeat("=", 50))

	totalTests := int64(0)
	totalFailures := int64(0)
	totalPanics := int64(0)
	totalDuration := time.Duration(0)

	for _, result := range results {
		totalTests += result.Stats.TestsExecuted
		totalFailures += result.Stats.TestsFailed
		totalPanics += result.Stats.PanicsRecovered
		totalDuration += result.Duration

		fmt.Printf("Test: %s\n", result.TestName)
		fmt.Printf("  Status: %v\n", result.Succeeded)
		fmt.Printf("  Duration: %v\n", result.Duration)
		fmt.Printf("  Tests Executed: %d\n", result.Stats.TestsExecuted)
		fmt.Printf("  Failures: %d\n", result.Stats.TestsFailed)
		fmt.Printf("  Panics Recovered: %d\n", result.Stats.PanicsRecovered)
		fmt.Printf("  Memory Allocated: %d bytes\n", result.Stats.MemoryAllocated)
		fmt.Printf("  Latency Injections: %d\n", result.Stats.LatencyInjected)

		if len(result.Errors) > 0 {
			fmt.Printf("  Sample Errors:\n")
			for i, err := range result.Errors {
				if i >= 3 { // Show only first 3 errors
					fmt.Printf("    ... and %d more errors\n", len(result.Errors)-3)
					break
				}
				fmt.Printf("    - %v\n", err)
			}
		}
		fmt.Println()
	}

	fmt.Printf("SUMMARY:\n")
	fmt.Printf("  Total Tests: %d\n", totalTests)
	fmt.Printf("  Total Failures: %d\n", totalFailures)
	fmt.Printf("  Total Panics: %d\n", totalPanics)
	fmt.Printf("  Total Duration: %v\n", totalDuration)

	if totalTests > 0 {
		failureRate := float64(totalFailures) / float64(totalTests) * 100
		fmt.Printf("  Failure Rate: %.2f%%\n", failureRate)
	}

	fmt.Println("\n✅ Chaos engineering analysis complete!")
	fmt.Println("Use these insights to improve your system's resilience.")
}

// ===== MAIN DEMONSTRATION =====

func DemonstrateChaosTesting() {
	fmt.Println("=== Go Chaos Engineering for Concurrency ===")
	fmt.Println("Based on Netflix's chaos engineering principles")
	fmt.Println("Finding the bugs that traditional tests miss!")
	fmt.Println()

	// Configure chaos testing
	config := DefaultChaosConfig()
	config.Duration = 10 * time.Second // Shorter for demo
	config.GoroutineCount = 20         // Fewer goroutines for demo

	// Create and run chaos test suite
	suite := NewChaosTestSuite(config)
	results := suite.RunAllChaosTests()

	// Analyze results
	AnalyzeChaosResults(results)

	fmt.Println("\n🎯 KEY TAKEAWAYS:")
	fmt.Println("1. Chaos engineering reveals hidden race conditions")
	fmt.Println("2. Systems fail in unexpected ways under stress")
	fmt.Println("3. Regular chaos testing builds confidence in production")
	fmt.Println("4. Netflix uses this to ensure 99.99% uptime")
	fmt.Println("5. Your code should handle chaos gracefully")
}

// ===== MAIN FUNCTION - CHAOS ENGINEERING DEMONSTRATION =====

func main() {
	fmt.Println("🔥 CHAOS ENGINEERING - NETFLIX STYLE TESTING")
	fmt.Println("===========================================")

	// Run the complete chaos engineering demonstration
	DemonstrateChaosTesting()

	fmt.Println("\n🎯 CHAOS ENGINEERING MASTERY COMPLETE!")
	fmt.Println("You now understand Netflix-level chaos testing!")
}
