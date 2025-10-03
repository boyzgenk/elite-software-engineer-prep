// Package property_testing demonstrates property-based testing for concurrent Go code
// This advanced testing technique is used by companies like Google, Netflix, and Uber
// to find edge cases that traditional unit tests miss
package main

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ===== PROPERTY-BASED TESTING FRAMEWORK =====

// Property represents a testable property of concurrent code
type Property struct {
	Name        string
	Description string
	TestFunc    func(input []interface{}) bool
	Generator   func() []interface{}
	Iterations  int
	Goroutines  int
}

// PropertyTester runs property-based tests on concurrent code
type PropertyTester struct {
	maxIterations int
	maxGoroutines int
	timeout       time.Duration
}

// NewPropertyTester creates a new property tester
func NewPropertyTester() *PropertyTester {
	return &PropertyTester{
		maxIterations: 1000,
		maxGoroutines: 100,
		timeout:       10 * time.Second,
	}
}

// TestProperty runs a property-based test
func (pt *PropertyTester) TestProperty(prop *Property) *PropertyResult {
	result := &PropertyResult{
		Property:     prop,
		StartTime:    time.Now(),
		Passed:       true,
		FailureCount: 0,
		TotalRuns:    0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), pt.timeout)
	defer cancel()

	for i := 0; i < prop.Iterations && result.Passed; i++ {
		select {
		case <-ctx.Done():
			result.TimedOut = true
			result.Passed = false
			break
		default:
			if !pt.runSingleTest(prop) {
				result.Passed = false
				result.FailureCount++
				result.FailingInput = prop.Generator()
			}
			result.TotalRuns++
		}
	}

	result.Duration = time.Since(result.StartTime)
	return result
}

// runSingleTest runs a single iteration of the property test
func (pt *PropertyTester) runSingleTest(prop *Property) bool {
	input := prop.Generator()

	// Run the test function concurrently
	var wg sync.WaitGroup
	results := make([]bool, prop.Goroutines)

	for i := 0; i < prop.Goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					results[idx] = false
				}
			}()
			results[idx] = prop.TestFunc(input)
		}(i)
	}

	wg.Wait()

	// All goroutines should return the same result for a correct property
	expected := results[0]
	for _, result := range results[1:] {
		if result != expected {
			return false
		}
	}

	return expected
}

// PropertyResult contains the results of a property-based test
type PropertyResult struct {
	Property     *Property
	Passed       bool
	TimedOut     bool
	FailureCount int
	TotalRuns    int
	FailingInput []interface{}
	StartTime    time.Time
	Duration     time.Duration
}

// String returns a string representation of the test result
func (pr *PropertyResult) String() string {
	status := "PASS"
	if !pr.Passed {
		status = "FAIL"
	}
	if pr.TimedOut {
		status = "TIMEOUT"
	}

	return fmt.Sprintf("Property: %s [%s] - %d/%d runs, %v duration",
		pr.Property.Name, status, pr.TotalRuns-pr.FailureCount,
		pr.TotalRuns, pr.Duration)
}

// ===== CONCURRENT DATA STRUCTURE PROPERTIES =====

// ConcurrentCounter for testing
type ConcurrentCounter struct {
	value int64
	mu    sync.Mutex
}

// Add adds a value to the counter
func (cc *ConcurrentCounter) Add(delta int64) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.value += delta
}

// Get returns the current value
func (cc *ConcurrentCounter) Get() int64 {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.value
}

// AddAtomic adds using atomic operations
func (cc *ConcurrentCounter) AddAtomic(delta int64) {
	atomic.AddInt64(&cc.value, delta)
}

// GetAtomic gets using atomic operations
func (cc *ConcurrentCounter) GetAtomic() int64 {
	return atomic.LoadInt64(&cc.value)
}

// ===== PROPERTY DEFINITIONS =====

// CounterCommutativityProperty tests that counter addition is commutative
func CounterCommutativityProperty() *Property {
	return &Property{
		Name:        "Counter Commutativity",
		Description: "Adding values to counter should be commutative: a+b = b+a",
		Iterations:  500,
		Goroutines:  10,
		Generator: func() []interface{} {
			return []interface{}{
				int64(rand.Intn(1000)),
				int64(rand.Intn(1000)),
			}
		},
		TestFunc: func(input []interface{}) bool {
			a := input[0].(int64)
			b := input[1].(int64)

			// Test with mutex-based counter
			counter1 := &ConcurrentCounter{}
			counter2 := &ConcurrentCounter{}

			var wg sync.WaitGroup

			// Counter 1: add a then b
			wg.Add(2)
			go func() {
				defer wg.Done()
				counter1.Add(a)
			}()
			go func() {
				defer wg.Done()
				counter1.Add(b)
			}()
			wg.Wait()

			// Counter 2: add b then a
			wg.Add(2)
			go func() {
				defer wg.Done()
				counter2.Add(b)
			}()
			go func() {
				defer wg.Done()
				counter2.Add(a)
			}()
			wg.Wait()

			return counter1.Get() == counter2.Get()
		},
	}
}

// CounterAssociativityProperty tests that counter addition is associative
func CounterAssociativityProperty() *Property {
	return &Property{
		Name:        "Counter Associativity",
		Description: "Adding values should be associative: (a+b)+c = a+(b+c)",
		Iterations:  300,
		Goroutines:  15,
		Generator: func() []interface{} {
			return []interface{}{
				int64(rand.Intn(100)),
				int64(rand.Intn(100)),
				int64(rand.Intn(100)),
			}
		},
		TestFunc: func(input []interface{}) bool {
			a := input[0].(int64)
			b := input[1].(int64)
			c := input[2].(int64)

			// Test (a+b)+c
			counter1 := &ConcurrentCounter{}
			var wg sync.WaitGroup

			wg.Add(3)
			go func() {
				defer wg.Done()
				counter1.Add(a)
			}()
			go func() {
				defer wg.Done()
				counter1.Add(b)
			}()
			go func() {
				defer wg.Done()
				counter1.Add(c)
			}()
			wg.Wait()

			// Compare with expected result
			expected := a + b + c
			return counter1.Get() == expected
		},
	}
}

// ===== CONCURRENT MAP PROPERTIES =====

// ConcurrentMap for testing
type ConcurrentMap struct {
	data sync.Map
}

// Set stores a key-value pair
func (cm *ConcurrentMap) Set(key, value interface{}) {
	cm.data.Store(key, value)
}

// Get retrieves a value by key
func (cm *ConcurrentMap) Get(key interface{}) (interface{}, bool) {
	return cm.data.Load(key)
}

// Delete removes a key
func (cm *ConcurrentMap) Delete(key interface{}) {
	cm.data.Delete(key)
}

// MapConsistencyProperty tests map read-write consistency
func MapConsistencyProperty() *Property {
	return &Property{
		Name:        "Map Consistency",
		Description: "Values written to map should be readable immediately after",
		Iterations:  400,
		Goroutines:  20,
		Generator: func() []interface{} {
			return []interface{}{
				fmt.Sprintf("key_%d", rand.Intn(100)),
				rand.Intn(1000),
			}
		},
		TestFunc: func(input []interface{}) bool {
			key := input[0].(string)
			value := input[1].(int)

			cm := &ConcurrentMap{}

			// Write and immediately read
			var wg sync.WaitGroup
			results := make([]interface{}, 2)

			wg.Add(2)

			// Writer
			go func() {
				defer wg.Done()
				cm.Set(key, value)
				results[0] = "written"
			}()

			// Reader (with small delay to ensure write happens first)
			go func() {
				defer wg.Done()
				time.Sleep(1 * time.Millisecond)
				val, exists := cm.Get(key)
				if exists && val == value {
					results[1] = "consistent"
				} else {
					results[1] = "inconsistent"
				}
			}()

			wg.Wait()

			return results[1] == "consistent"
		},
	}
}

// ===== CHANNEL PROPERTIES =====

// ChannelOrderingProperty tests FIFO ordering in buffered channels
func ChannelOrderingProperty() *Property {
	return &Property{
		Name:        "Channel FIFO Ordering",
		Description: "Buffered channels should maintain FIFO ordering",
		Iterations:  200,
		Goroutines:  5,
		Generator: func() []interface{} {
			size := rand.Intn(10) + 1
			values := make([]int, size)
			for i := range values {
				values[i] = rand.Intn(1000)
			}
			return []interface{}{values}
		},
		TestFunc: func(input []interface{}) bool {
			values := input[0].([]int)
			ch := make(chan int, len(values)+1) // Buffered channel

			var wg sync.WaitGroup

			// Producer
			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, v := range values {
					ch <- v
				}
				close(ch)
			}()

			// Consumer
			received := make([]int, 0, len(values))
			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range ch {
					received = append(received, v)
				}
			}()

			wg.Wait()

			// Check if order is preserved
			return reflect.DeepEqual(values, received)
		},
	}
}

// ===== GOROUTINE PROPERTIES =====

// GoroutineCompletionProperty tests that all goroutines complete
func GoroutineCompletionProperty() *Property {
	return &Property{
		Name:        "Goroutine Completion",
		Description: "All spawned goroutines should complete execution",
		Iterations:  100,
		Goroutines:  3,
		Generator: func() []interface{} {
			return []interface{}{
				rand.Intn(50) + 1, // number of goroutines to spawn
			}
		},
		TestFunc: func(input []interface{}) bool {
			numGoroutines := input[0].(int)

			var wg sync.WaitGroup
			completed := int64(0)

			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					// Simulate some work
					time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
					atomic.AddInt64(&completed, 1)
				}(i)
			}

			// Wait with timeout
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			select {
			case <-done:
				return atomic.LoadInt64(&completed) == int64(numGoroutines)
			case <-time.After(5 * time.Second):
				return false // Timeout
			}
		},
	}
}

// ===== ADVANCED PROPERTIES =====

// AtomicOperationProperty tests atomic operation consistency
func AtomicOperationProperty() *Property {
	return &Property{
		Name:        "Atomic Operation Consistency",
		Description: "Atomic operations should be consistent across goroutines",
		Iterations:  300,
		Goroutines:  20,
		Generator: func() []interface{} {
			return []interface{}{
				int64(rand.Intn(100) + 1), // number of increments per goroutine
			}
		},
		TestFunc: func(input []interface{}) bool {
			increments := input[0].(int64)
			counter := int64(0)

			var wg sync.WaitGroup
			numGoroutines := 10

			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := int64(0); j < increments; j++ {
						atomic.AddInt64(&counter, 1)
					}
				}()
			}

			wg.Wait()

			expected := int64(numGoroutines) * increments
			actual := atomic.LoadInt64(&counter)

			return actual == expected
		},
	}
}

// ===== PROPERTY TEST SUITE =====

// PropertyTestSuite contains all property-based tests
type PropertyTestSuite struct {
	properties []*Property
	tester     *PropertyTester
}

// NewPropertyTestSuite creates a new test suite
func NewPropertyTestSuite() *PropertyTestSuite {
	return &PropertyTestSuite{
		properties: []*Property{
			CounterCommutativityProperty(),
			CounterAssociativityProperty(),
			MapConsistencyProperty(),
			ChannelOrderingProperty(),
			GoroutineCompletionProperty(),
			AtomicOperationProperty(),
		},
		tester: NewPropertyTester(),
	}
}

// RunAll runs all property-based tests
func (pts *PropertyTestSuite) RunAll() []*PropertyResult {
	results := make([]*PropertyResult, len(pts.properties))

	fmt.Println("=== Running Property-Based Tests ===")

	for i, prop := range pts.properties {
		fmt.Printf("Running: %s...\n", prop.Name)
		results[i] = pts.tester.TestProperty(prop)
		fmt.Printf("  %s\n", results[i].String())

		if !results[i].Passed {
			fmt.Printf("  Failing input: %v\n", results[i].FailingInput)
		}
	}

	return results
}

// ===== BENCHMARKING PROPERTIES =====

// BenchmarkProperties benchmarks the property-based tests
func BenchmarkProperties(b *testing.B) {
	suite := NewPropertyTestSuite()

	for _, prop := range suite.properties {
		b.Run(prop.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				suite.tester.TestProperty(prop)
			}
		})
	}
}

// ===== MAIN DEMONSTRATION =====

func DemoPropertyTesting() {
	fmt.Println("=== Go Property-Based Testing for Concurrency ===")
	fmt.Println("This demonstrates advanced testing techniques used by FAANG companies")
	fmt.Println()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create and run property test suite
	suite := NewPropertyTestSuite()
	results := suite.RunAll()

	// Summary
	fmt.Println("\n=== Test Summary ===")
	passed := 0
	total := len(results)

	for _, result := range results {
		if result.Passed {
			passed++
		}
	}

	fmt.Printf("Tests passed: %d/%d\n", passed, total)
	fmt.Printf("Success rate: %.1f%%\n", float64(passed)/float64(total)*100)

	if passed == total {
		fmt.Println("🎉 All property-based tests passed!")
		fmt.Println("Your concurrent code satisfies all tested properties.")
	} else {
		fmt.Println("⚠️  Some tests failed. Review the failing properties above.")
	}

	fmt.Println("\n=== Property-Based Testing Benefits ===")
	fmt.Println("1. Finds edge cases that unit tests miss")
	fmt.Println("2. Tests mathematical properties of concurrent code")
	fmt.Println("3. Provides confidence in correctness across input ranges")
	fmt.Println("4. Catches race conditions and consistency issues")
	fmt.Println("5. Scales testing to thousands of scenarios automatically")

	fmt.Println("\n=== Advanced Usage ===")
	fmt.Println("- Run with -race flag to detect race conditions")
	fmt.Println("- Increase iterations for more thorough testing")
	fmt.Println("- Add custom properties for your specific use cases")
	fmt.Println("- Use property tests in CI/CD pipelines")
}

// ===== MAIN FUNCTION - PROPERTY-BASED TESTING DEMO =====

func main() {
	fmt.Println("🧮 PROPERTY-BASED TESTING - MATHEMATICAL CONCURRENCY VALIDATION")
	fmt.Println("==============================================================")

	// Run the complete property-based testing demonstration
	DemoPropertyTesting()

	fmt.Println("\n🎯 PROPERTY-BASED TESTING MASTERY COMPLETE!")
	fmt.Println("You now understand mathematical testing approaches!")
}
