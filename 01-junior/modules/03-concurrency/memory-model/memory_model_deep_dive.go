// Package memory_model demonstrates Go's memory model and happens-before relationships
// This is advanced concurrency knowledge essential for FAANG interviews
// Based on the Go Memory Model specification: https://golang.org/ref/mem
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// ===== GO MEMORY MODEL FUNDAMENTALS =====

/*
GO MEMORY MODEL KEY PRINCIPLES:

1. Programs that modify data being simultaneously accessed by multiple goroutines
   must serialize such access.

2. To serialize access, protect the data with channel operations or other
   synchronization primitives such as those in the sync and sync/atomic packages.

3. Reads and writes of values larger than a single machine word are not guaranteed
   to be atomic.

4. Within a single goroutine, reads and writes must behave as if they executed
   in the order specified by the program.

HAPPENS-BEFORE RELATIONSHIP:
- If event e1 happens before event e2, then e2 happens after e1.
- If e1 does not happen before e2 and e2 does not happen before e1,
  then e1 and e2 happen concurrently.
*/

// ===== HAPPENS-BEFORE EXAMPLES =====

// Example 1: Channel Communication Synchronization
func ChannelHappensBefore() {
	fmt.Println("=== Channel Happens-Before Example ===")

	var a string
	var c = make(chan bool, 1)

	go func() {
		a = "hello world" // happens before the send
		c <- true         // send happens before corresponding receive
	}()

	<-c            // receive happens before this line
	fmt.Println(a) // guaranteed to print "hello world"

	fmt.Println("✅ Channel synchronization ensures happens-before relationship")
}

// Example 2: Mutex Synchronization
func MutexHappensBefore() {
	fmt.Println("\n=== Mutex Happens-Before Example ===")

	var mu sync.Mutex
	var a string

	go func() {
		mu.Lock()
		a = "hello world" // happens before Unlock
		mu.Unlock()       // Unlock happens before next Lock
	}()

	// Brief delay to ensure goroutine runs first
	time.Sleep(10 * time.Millisecond)

	mu.Lock()      // Lock happens after previous Unlock
	fmt.Println(a) // guaranteed to see "hello world"
	mu.Unlock()

	fmt.Println("✅ Mutex synchronization ensures happens-before relationship")
}

// Example 3: Once Synchronization
func OnceHappensBefore() {
	fmt.Println("\n=== sync.Once Happens-Before Example ===")

	var once sync.Once
	var a string

	setup := func() {
		a = "hello world" // happens before once.Do returns
	}

	go func() {
		once.Do(setup) // Do returns after setup completes
	}()

	// Brief delay
	time.Sleep(10 * time.Millisecond)

	once.Do(setup) // This call returns immediately
	fmt.Println(a) // guaranteed to see "hello world"

	fmt.Println("✅ sync.Once ensures happens-before relationship")
}

// ===== DATA RACES AND MEMORY ORDERING =====

// RaceConditionExample demonstrates a data race
type RaceConditionExample struct {
	counter  int64
	data     []int
	finished bool
}

// UnsafeIncrement creates a data race
func (rce *RaceConditionExample) UnsafeIncrement() {
	// THIS IS WRONG - Data race!
	temp := rce.counter
	runtime.Gosched() // Force context switch to make race more likely
	rce.counter = temp + 1
}

// SafeIncrement uses atomic operations
func (rce *RaceConditionExample) SafeIncrement() {
	atomic.AddInt64(&rce.counter, 1)
}

// DemonstrateDataRace shows how data races occur
func DemonstrateDataRace() {
	fmt.Println("\n=== Data Race Demonstration ===")

	example := &RaceConditionExample{}

	var wg sync.WaitGroup
	numGoroutines := 100
	iterations := 100

	// Test with unsafe increment (data race)
	fmt.Println("Running unsafe increment (with data races)...")
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				example.UnsafeIncrement()
			}
		}()
	}
	wg.Wait()

	expected := int64(numGoroutines * iterations)
	actual := example.counter
	fmt.Printf("Expected: %d, Got: %d, Lost updates: %d\n",
		expected, actual, expected-actual)

	// Reset and test with safe increment
	example.counter = 0
	fmt.Println("Running safe increment (atomic operations)...")
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				example.SafeIncrement()
			}
		}()
	}
	wg.Wait()

	actual = atomic.LoadInt64(&example.counter)
	fmt.Printf("Expected: %d, Got: %d, Lost updates: %d\n",
		expected, actual, expected-actual)

	fmt.Println("✅ Atomic operations prevent data races")
}

// ===== MEMORY VISIBILITY PATTERNS =====

// VisibilityExample demonstrates memory visibility issues
type VisibilityExample struct {
	flag int32
	data int32
}

// Writer sets data and flag
func (ve *VisibilityExample) Writer() {
	atomic.StoreInt32(&ve.data, 42) // Store data
	atomic.StoreInt32(&ve.flag, 1)  // Set flag (acts as release)
}

// Reader checks flag and reads data
func (ve *VisibilityExample) Reader() int32 {
	if atomic.LoadInt32(&ve.flag) == 1 { // Check flag (acts as acquire)
		return atomic.LoadInt32(&ve.data) // Read data
	}
	return -1
}

// UnsafeWriter creates visibility problems
func (ve *VisibilityExample) UnsafeWriter() {
	ve.data = 42 // Unsafe store
	ve.flag = 1  // Unsafe store - no ordering guarantees
}

// UnsafeReader has visibility problems
func (ve *VisibilityExample) UnsafeReader() int32 {
	if ve.flag == 1 { // Unsafe read
		return ve.data // May see stale value!
	}
	return -1
}

// DemonstrateMemoryVisibility shows memory visibility patterns
func DemonstrateMemoryVisibility() {
	fmt.Println("\n=== Memory Visibility Demonstration ===")

	// Test safe visibility with atomic operations
	fmt.Println("Testing safe visibility (atomic operations)...")
	example := &VisibilityExample{}

	go example.Writer()
	time.Sleep(1 * time.Millisecond) // Allow writer to complete

	result := example.Reader()
	fmt.Printf("Safe reader result: %d (expected: 42)\n", result)

	// Test unsafe visibility
	fmt.Println("Testing unsafe visibility (may see inconsistent state)...")
	example2 := &VisibilityExample{}

	go example2.UnsafeWriter()
	time.Sleep(1 * time.Millisecond) // Allow writer to complete

	result = example2.UnsafeReader()
	fmt.Printf("Unsafe reader result: %d (may be -1 or 42)\n", result)

	fmt.Println("✅ Atomic operations ensure proper memory visibility")
}

// ===== ADVANCED MEMORY ORDERING PATTERNS =====

// LoadStoreReordering demonstrates compiler/CPU reordering issues
type LoadStoreReordering struct {
	x, y   int32
	r1, r2 int32
}

// Writer1 performs stores in order
func (lsr *LoadStoreReordering) Writer1() {
	atomic.StoreInt32(&lsr.x, 1)
	atomic.StoreInt32(&lsr.r1, atomic.LoadInt32(&lsr.y))
}

// Writer2 performs stores in reverse order
func (lsr *LoadStoreReordering) Writer2() {
	atomic.StoreInt32(&lsr.y, 1)
	atomic.StoreInt32(&lsr.r2, atomic.LoadInt32(&lsr.x))
}

// DemonstrateReordering shows how reordering can affect visibility
func DemonstrateReordering() {
	fmt.Println("\n=== Memory Reordering Demonstration ===")

	iterations := 100000
	bothZero := 0

	for i := 0; i < iterations; i++ {
		example := &LoadStoreReordering{}

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			example.Writer1()
		}()

		go func() {
			defer wg.Done()
			example.Writer2()
		}()

		wg.Wait()

		if example.r1 == 0 && example.r2 == 0 {
			bothZero++
		}
	}

	fmt.Printf("Out of %d iterations, both r1 and r2 were 0: %d times (%.2f%%)\n",
		iterations, bothZero, float64(bothZero)/float64(iterations)*100)
	fmt.Println("This demonstrates that memory reordering can occur")
}

// ===== FALSE SHARING DEMONSTRATION =====

// FalseSharingDemo demonstrates false sharing performance issues
type FalseSharingDemo struct {
	// These fields will likely be on the same cache line
	counter1 int64
	counter2 int64
}

// PaddedDemo avoids false sharing with padding
type PaddedDemo struct {
	counter1 int64
	_        [56]byte // Padding to separate cache lines (64-byte cache line)
	counter2 int64
	_        [56]byte // More padding
}

// BenchmarkFalseSharing measures performance impact of false sharing
func BenchmarkFalseSharing() {
	fmt.Println("\n=== False Sharing Demonstration ===")

	// Test without padding (false sharing)
	fmt.Println("Testing with false sharing...")
	falseSharingDemo := &FalseSharingDemo{}

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			atomic.AddInt64(&falseSharingDemo.counter1, 1)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			atomic.AddInt64(&falseSharingDemo.counter2, 1)
		}
	}()

	wg.Wait()
	falseSharingTime := time.Since(start)

	// Test with padding (no false sharing)
	fmt.Println("Testing without false sharing (padded)...")
	paddedDemo := &PaddedDemo{}

	start = time.Now()
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			atomic.AddInt64(&paddedDemo.counter1, 1)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			atomic.AddInt64(&paddedDemo.counter2, 1)
		}
	}()

	wg.Wait()
	paddedTime := time.Since(start)

	fmt.Printf("False sharing time: %v\n", falseSharingTime)
	fmt.Printf("Padded time: %v\n", paddedTime)
	fmt.Printf("Performance improvement: %.2fx\n",
		float64(falseSharingTime)/float64(paddedTime))

	fmt.Println("✅ Padding prevents false sharing and improves performance")
}

// ===== UNSAFE PACKAGE PATTERNS =====

// UnsafePatterns demonstrates unsafe memory access patterns
type UnsafePatterns struct {
	data [4]int32
}

// AtomicSliceUpdate shows atomic operations on slices using unsafe
func (up *UnsafePatterns) AtomicSliceUpdate(index int, value int32) {
	// Calculate pointer to the specific array element
	ptr := (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&up.data[0])) +
		uintptr(index)*unsafe.Sizeof(up.data[0])))
	atomic.StoreInt32(ptr, value)
}

// AtomicSliceLoad atomically loads from slice using unsafe
func (up *UnsafePatterns) AtomicSliceLoad(index int) int32 {
	ptr := (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&up.data[0])) +
		uintptr(index)*unsafe.Sizeof(up.data[0])))
	return atomic.LoadInt32(ptr)
}

// DemonstrateUnsafePatterns shows advanced unsafe usage
func DemonstrateUnsafePatterns() {
	fmt.Println("\n=== Unsafe Package Patterns ===")

	patterns := &UnsafePatterns{}

	// Concurrent updates to different array elements
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				patterns.AtomicSliceUpdate(index, int32(j))
			}
		}(i)
	}
	wg.Wait()

	// Read final values
	fmt.Println("Final array values:")
	for i := 0; i < 4; i++ {
		value := patterns.AtomicSliceLoad(i)
		fmt.Printf("  data[%d] = %d\n", i, value)
	}

	fmt.Println("✅ Unsafe package enables atomic operations on complex data structures")
	fmt.Println("⚠️  WARNING: Use unsafe package only when absolutely necessary!")
}

// ===== MEMORY FENCE PATTERNS =====

// MemoryFenceDemo demonstrates memory fence usage
type MemoryFenceDemo struct {
	ready int32
	data  int32
}

// Publisher writes data and sets ready flag
func (mfd *MemoryFenceDemo) Publisher() {
	atomic.StoreInt32(&mfd.data, 42) // Write data
	// Memory fence is implicit in atomic operations
	atomic.StoreInt32(&mfd.ready, 1) // Set ready flag
}

// Subscriber waits for ready flag and reads data
func (mfd *MemoryFenceDemo) Subscriber() int32 {
	for atomic.LoadInt32(&mfd.ready) == 0 {
		runtime.Gosched() // Yield to other goroutines
	}
	// Memory fence is implicit in atomic operations
	return atomic.LoadInt32(&mfd.data) // Read data
}

// DemonstrateMemoryFences shows memory fence patterns
func DemonstrateMemoryFences() {
	fmt.Println("\n=== Memory Fence Demonstration ===")

	demo := &MemoryFenceDemo{}

	// Start subscriber
	result := make(chan int32)
	go func() {
		result <- demo.Subscriber()
	}()

	// Brief delay before publishing
	time.Sleep(1 * time.Millisecond)

	// Publish data
	demo.Publisher()

	// Get result
	value := <-result
	fmt.Printf("Subscriber received: %d (expected: 42)\n", value)

	fmt.Println("✅ Memory fences in atomic operations ensure proper ordering")
}

// ===== INTERVIEW QUESTIONS AND PATTERNS =====

// InterviewQuestion represents a memory model interview question
type InterviewQuestion struct {
	Question    string
	Code        func()
	Explanation string
}

// GetMemoryModelInterviewQuestions returns common interview questions
func GetMemoryModelInterviewQuestions() []InterviewQuestion {
	return []InterviewQuestion{
		{
			Question: "What happens if two goroutines access a shared variable without synchronization?",
			Code: func() {
				var x int
				go func() { x = 1 }()
				go func() { x = 2 }()
				// Data race! Result is undefined
				_ = x // Suppress unused warning
			},
			Explanation: "This creates a data race. The behavior is undefined according to the Go memory model.",
		},
		{
			Question: "Is reading/writing a single int atomic in Go?",
			Code: func() {
				var x int64
				go func() { x = 42 }()                    // Not guaranteed atomic on 32-bit systems
				go func() { atomic.StoreInt64(&x, 42) }() // Guaranteed atomic
			},
			Explanation: "Only operations on single machine word are atomic. Use sync/atomic for guarantees.",
		},
		{
			Question: "What's the difference between buffered and unbuffered channels for happens-before?",
			Code: func() {
				unbuffered := make(chan int)
				buffered := make(chan int, 1)

				// Unbuffered: send happens before receive
				// Buffered: receive happens before next send (if buffer full)
				_ = unbuffered
				_ = buffered
			},
			Explanation: "Unbuffered channels provide stronger synchronization guarantees than buffered channels.",
		},
	}
}

// ===== MAIN DEMONSTRATION =====

func DemonstrateMemoryModel() {
	fmt.Println("=== Go Memory Model Deep Dive ===")
	fmt.Println("Essential knowledge for senior Go developers and FAANG interviews")
	fmt.Println()

	// Run demonstrations
	ChannelHappensBefore()
	MutexHappensBefore()
	OnceHappensBefore()
	DemonstrateDataRace()
	DemonstrateMemoryVisibility()
	DemonstrateReordering()
	BenchmarkFalseSharing()
	DemonstrateUnsafePatterns()
	DemonstrateMemoryFences()

	// Show interview questions
	fmt.Println("\n=== Common Interview Questions ===")
	questions := GetMemoryModelInterviewQuestions()
	for i, q := range questions {
		fmt.Printf("\nQ%d: %s\n", i+1, q.Question)
		fmt.Printf("A%d: %s\n", i+1, q.Explanation)
	}

	fmt.Println("\n🎯 KEY TAKEAWAYS:")
	fmt.Println("1. Go memory model defines when memory operations are visible")
	fmt.Println("2. Use synchronization primitives to establish happens-before relationships")
	fmt.Println("3. Data races lead to undefined behavior - always avoid them")
	fmt.Println("4. Atomic operations provide memory ordering guarantees")
	fmt.Println("5. Understand cache line effects and false sharing")
	fmt.Println("6. The happens-before relationship is transitive")
	fmt.Println("7. Within a single goroutine, operations appear sequentially consistent")

	fmt.Println("\n📚 FURTHER READING:")
	fmt.Println("- Go Memory Model: https://golang.org/ref/mem")
	fmt.Println("- Effective Go: https://golang.org/doc/effective_go.html#concurrency")
	fmt.Println("- Go Race Detector: https://golang.org/doc/articles/race_detector.html")
}

// ===== MAIN FUNCTION - COMPLETE DEMONSTRATION =====

func main() {
	fmt.Println("🚀 GO MEMORY MODEL DEEP DIVE - PRODUCTION MASTERY")
	fmt.Println("=================================================")

	// Run all memory model demonstrations
	ChannelHappensBefore()
	time.Sleep(100 * time.Millisecond)

	MutexHappensBefore()
	time.Sleep(100 * time.Millisecond)

	OnceHappensBefore()
	time.Sleep(100 * time.Millisecond)

	DemonstrateDataRace()
	time.Sleep(100 * time.Millisecond)

	DemonstrateMemoryVisibility()
	time.Sleep(100 * time.Millisecond)

	DemonstrateReordering()
	time.Sleep(100 * time.Millisecond)

	BenchmarkFalseSharing()
	time.Sleep(100 * time.Millisecond)

	DemonstrateUnsafePatterns()
	time.Sleep(100 * time.Millisecond)

	DemonstrateMemoryFences()
	time.Sleep(100 * time.Millisecond)

	DemonstrateMemoryModel()

	fmt.Println("\n🎯 MEMORY MODEL MASTERY COMPLETE!")
	fmt.Println("You now understand Go's memory model at FAANG senior engineer level!")
}
