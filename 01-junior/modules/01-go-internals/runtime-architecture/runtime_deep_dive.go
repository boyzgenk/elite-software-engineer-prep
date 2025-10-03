// Package runtime demonstrates Go's runtime architecture internals
// Essential production-level knowledge for FAANG L3/L4+ interviews

package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Go Runtime Architecture Components:
// - Bootstrap process and initialization
// - Stack management and growth
// - Memory allocator integration
// - System call handling
// - Signal handling and preemption

func main() {
	fmt.Println("🔥 Go Runtime Architecture Deep Dive")
	fmt.Println("===================================")

	// Runtime initialization and bootstrap
	demonstrateRuntimeBootstrap()

	// Stack management and growth
	demonstrateStackManagement()

	// Memory allocator integration
	demonstrateMemoryAllocator()

	// System call handling
	demonstrateSystemCalls()

	// Runtime coordination patterns
	demonstrateRuntimeCoordination()

	// Advanced runtime inspection
	demonstrateRuntimeInspection()
}

// Runtime bootstrap and initialization sequence
func demonstrateRuntimeBootstrap() {
	fmt.Println("\n🚀 Runtime Bootstrap & Initialization")
	fmt.Println("-----------------------------------")

	// Show runtime version and build info
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Compiler: %s\n", runtime.Compiler)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("Operating System: %s\n", runtime.GOOS)

	// Show initial runtime state
	fmt.Printf("Initial GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("Initial NumGoroutine: %d\n", runtime.NumGoroutine())
	fmt.Printf("Initial NumCPU: %d\n", runtime.NumCPU())

	// Demonstrate runtime parameter adjustment
	oldMaxProcs := runtime.GOMAXPROCS(2)
	fmt.Printf("Changed GOMAXPROCS from %d to 2\n", oldMaxProcs)

	// Restore original setting
	runtime.GOMAXPROCS(oldMaxProcs)
	fmt.Printf("Restored GOMAXPROCS to %d\n", oldMaxProcs)
}

// Stack management: growth, shrinking, and optimization
func demonstrateStackManagement() {
	fmt.Println("\n📚 Stack Management & Growth")
	fmt.Println("---------------------------")

	// Show initial stack usage
	showStackUsage("Initial")

	// Demonstrate stack growth through deep recursion
	fmt.Println("Testing stack growth with recursion...")
	result := fibonacci(25) // Moderate recursion to trigger stack growth
	fmt.Printf("Fibonacci(25) = %d\n", result)

	showStackUsage("After recursion")

	// Demonstrate goroutine stack size differences
	demonstrateGoroutineStacks()
}

// Recursive function to demonstrate stack growth
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}

	// Create some stack variables to increase stack usage
	var localVars [10]int64
	for i := range localVars {
		localVars[i] = int64(i * n)
	}

	_ = localVars // Use the variables
	return fibonacci(n-1) + fibonacci(n-2)
}

// Demonstrate goroutine stack characteristics
func demonstrateGoroutineStacks() {
	fmt.Println("\nGoroutine Stack Characteristics:")

	var wg sync.WaitGroup
	stackSizes := make([]uintptr, 5)

	// Create multiple goroutines to show stack behavior
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Get stack bounds for this goroutine
			stack := make([]byte, 64<<10) // 64KB buffer
			stackSize := runtime.Stack(stack, false)
			stackSizes[id] = uintptr(stackSize)

			fmt.Printf("Goroutine %d stack size: %d bytes\n", id, stackSize)

			// Do some work to potentially grow stack
			deepWork(10, id)
		}(i)
	}

	wg.Wait()

	// Analyze stack size patterns
	var totalSize uintptr
	for _, size := range stackSizes {
		totalSize += size
	}
	avgSize := totalSize / uintptr(len(stackSizes))
	fmt.Printf("Average goroutine stack usage: %d bytes\n", avgSize)
}

// Helper function for stack growth demonstration
func deepWork(depth int, id int) {
	if depth <= 0 {
		return
	}

	// Create local variables to use stack space
	var buffer [1024]byte
	for i := range buffer {
		buffer[i] = byte(id + depth)
	}

	// Recurse to grow stack
	deepWork(depth-1, id)

	// Use buffer to prevent optimization
	_ = buffer
}

// Memory allocator integration with runtime
func demonstrateMemoryAllocator() {
	fmt.Println("\n💾 Memory Allocator Integration")
	fmt.Println("------------------------------")

	// Show memory stats before allocation
	showMemoryStats("Before allocation")

	// Demonstrate different allocation patterns
	demonstrateSmallAllocations()
	demonstrateLargeAllocations()
	demonstratePooledAllocations()

	// Force GC and show final stats
	runtime.GC()
	showMemoryStats("After GC")
}

// Small object allocation patterns
func demonstrateSmallAllocations() {
	fmt.Println("Small allocations (< 32KB):")

	// Allocate many small objects
	objects := make([]*smallObject, 1000)
	for i := range objects {
		objects[i] = &smallObject{
			id:   i,
			data: make([]byte, 64), // Small allocation
		}
	}

	showMemoryStats("After small allocations")

	// Keep objects alive
	_ = objects
}

// Large object allocation patterns
func demonstrateLargeAllocations() {
	fmt.Println("Large allocations (> 32KB):")

	// Allocate large objects (go directly to heap)
	largeObjects := make([]*largeObject, 10)
	for i := range largeObjects {
		largeObjects[i] = &largeObject{
			id:   i,
			data: make([]byte, 64*1024), // 64KB - large allocation
		}
	}

	showMemoryStats("After large allocations")

	// Keep objects alive
	_ = largeObjects
}

// Object pooling to reduce allocator pressure
func demonstratePooledAllocations() {
	fmt.Println("Pooled allocations:")

	// Create object pool
	pool := &sync.Pool{
		New: func() interface{} {
			return &pooledObject{
				data: make([]byte, 256),
			}
		},
	}

	// Use pooled objects
	for i := 0; i < 100; i++ {
		obj := pool.Get().(*pooledObject)
		obj.id = i

		// Simulate work
		for j := range obj.data {
			obj.data[j] = byte(i + j)
		}

		pool.Put(obj)
	}

	showMemoryStats("After pooled allocations")
}

// System call handling and runtime coordination
func demonstrateSystemCalls() {
	fmt.Println("\n🔧 System Call Handling")
	fmt.Println("----------------------")

	// Show goroutines before I/O operations
	fmt.Printf("Goroutines before I/O: %d\n", runtime.NumGoroutine())

	var wg sync.WaitGroup

	// Demonstrate blocking system calls
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Simulate blocking I/O (sleep is a blocking syscall)
			fmt.Printf("Goroutine %d: Starting blocking operation\n", id)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Goroutine %d: Completed blocking operation\n", id)
		}(i)
	}

	// Show goroutines during I/O operations
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Goroutines during I/O: %d\n", runtime.NumGoroutine())

	wg.Wait()
	fmt.Printf("Goroutines after I/O: %d\n", runtime.NumGoroutine())
}

// Runtime coordination between components
func demonstrateRuntimeCoordination() {
	fmt.Println("\n⚡ Runtime Coordination Patterns")
	fmt.Println("------------------------------")

	// Demonstrate scheduler coordination
	demonstrateSchedulerCoordination()

	// Demonstrate GC coordination
	demonstrateGCCoordination()

	// Demonstrate signal handling
	demonstrateSignalHandling()
}

// Scheduler coordination during high concurrency
func demonstrateSchedulerCoordination() {
	fmt.Println("Scheduler coordination:")

	const numGoroutines = 100
	var counter int64
	var wg sync.WaitGroup

	start := time.Now()

	// Create many goroutines to test scheduler
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Atomic increment to test coordination
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)

				// Yield to scheduler occasionally
				if j%100 == 0 {
					runtime.Gosched()
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("Coordinated %d goroutines, final counter: %d, duration: %v\n",
		numGoroutines, counter, duration)
}

// GC coordination with runtime
func demonstrateGCCoordination() {
	fmt.Println("GC coordination:")

	// Show GC stats before
	var stats1 runtime.MemStats
	runtime.ReadMemStats(&stats1)

	// Allocate memory
	data := make([][]byte, 1000)
	for i := range data {
		data[i] = make([]byte, 1024)
	}

	// Force GC and measure coordination
	start := time.Now()
	runtime.GC()
	gcDuration := time.Since(start)

	var stats2 runtime.MemStats
	runtime.ReadMemStats(&stats2)

	fmt.Printf("GC coordination took: %v\n", gcDuration)
	fmt.Printf("GC cycles: %d -> %d\n", stats1.NumGC, stats2.NumGC)

	// Keep data alive
	_ = data
}

// Signal handling demonstration
func demonstrateSignalHandling() {
	fmt.Println("Signal handling patterns:")

	// The runtime handles signals for preemption
	// This demonstrates the coordination between signal handling and scheduling

	done := make(chan bool)

	// CPU-intensive goroutine
	go func() {
		count := 0
		start := time.Now()

		// Tight loop that should be preempted
		for time.Since(start) < 200*time.Millisecond {
			count++
		}

		fmt.Printf("CPU-intensive work completed: %d iterations\n", count)
		done <- true
	}()

	// This should still execute due to preemptive scheduling
	go func() {
		for i := 0; i < 4; i++ {
			fmt.Printf("Preempted goroutine tick: %d\n", i+1)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	<-done
	time.Sleep(100 * time.Millisecond)
}

// Advanced runtime inspection techniques
func demonstrateRuntimeInspection() {
	fmt.Println("\n🔍 Advanced Runtime Inspection")
	fmt.Println("-----------------------------")

	// Goroutine tracing
	demonstrateGoroutineTracing()

	// Memory layout inspection
	demonstrateMemoryLayout()

	// Performance characteristics
	demonstratePerformanceCharacteristics()
}

// Goroutine state and stack tracing
func demonstrateGoroutineTracing() {
	fmt.Println("Goroutine tracing:")

	// Get stack trace of all goroutines
	buf := make([]byte, 1<<16)            // 64KB buffer
	stackSize := runtime.Stack(buf, true) // all goroutines

	fmt.Printf("Stack trace size: %d bytes\n", stackSize)

	// Count goroutines in different states
	numGoroutines := runtime.NumGoroutine()
	fmt.Printf("Total goroutines: %d\n", numGoroutines)

	// Create goroutines in different states for demonstration
	ch := make(chan int)

	// Blocked goroutine
	go func() {
		<-ch // Will block
	}()

	// Running goroutine
	go func() {
		for i := 0; i < 1000000; i++ {
			_ = i * i
		}
	}()

	time.Sleep(10 * time.Millisecond)
	fmt.Printf("Goroutines after creating blocked/running: %d\n", runtime.NumGoroutine())

	close(ch) // Unblock the first goroutine
}

// Memory layout and alignment inspection
func demonstrateMemoryLayout() {
	fmt.Println("Memory layout inspection:")

	// Demonstrate struct alignment and padding
	type AlignedStruct struct {
		a int8  // 1 byte
		b int64 // 8 bytes (will be aligned)
		c int8  // 1 byte
		d int32 // 4 bytes (will be aligned)
	}

	var s AlignedStruct
	fmt.Printf("Struct size: %d bytes\n", unsafe.Sizeof(s))
	fmt.Printf("Field 'a' offset: %d\n", unsafe.Offsetof(s.a))
	fmt.Printf("Field 'b' offset: %d\n", unsafe.Offsetof(s.b))
	fmt.Printf("Field 'c' offset: %d\n", unsafe.Offsetof(s.c))
	fmt.Printf("Field 'd' offset: %d\n", unsafe.Offsetof(s.d))

	// Show pointer size and alignment
	var ptr *int
	fmt.Printf("Pointer size: %d bytes\n", unsafe.Sizeof(ptr))
	fmt.Printf("Pointer alignment: %d bytes\n", unsafe.Alignof(ptr))
}

// Performance characteristics analysis
func demonstratePerformanceCharacteristics() {
	fmt.Println("Performance characteristics:")

	// Goroutine creation overhead
	start := time.Now()
	const numGoroutines = 10000

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Minimal work
		}()
	}
	wg.Wait()

	creationTime := time.Since(start)
	fmt.Printf("Created %d goroutines in %v (%.2f ns/goroutine)\n",
		numGoroutines, creationTime, float64(creationTime.Nanoseconds())/float64(numGoroutines))

	// Channel operation overhead
	ch := make(chan int, 1)
	start = time.Now()
	const numOps = 100000

	for i := 0; i < numOps; i++ {
		ch <- i
		<-ch
	}

	channelTime := time.Since(start)
	fmt.Printf("Performed %d channel operations in %v (%.2f ns/operation)\n",
		numOps*2, channelTime, float64(channelTime.Nanoseconds())/float64(numOps*2))
}

// Helper types for demonstrations
type smallObject struct {
	id   int
	data []byte
}

type largeObject struct {
	id   int
	data []byte
}

type pooledObject struct {
	id   int
	data []byte
}

// Helper functions
func showStackUsage(label string) {
	buf := make([]byte, 1<<12) // 4KB buffer
	stackSize := runtime.Stack(buf, false)
	fmt.Printf("%s stack usage: %d bytes\n", label, stackSize)
}

func showMemoryStats(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("\n📊 %s:\n", label)
	fmt.Printf("  Heap Alloc: %.2f MB\n", float64(m.HeapAlloc)/1024/1024)
	fmt.Printf("  Heap Objects: %d\n", m.HeapObjects)
	fmt.Printf("  GC Cycles: %d\n", m.NumGC)
	fmt.Printf("  Goroutines: %d\n", runtime.NumGoroutine())
}

/*
🔥 FAANG Interview Key Points:

1. Runtime Bootstrap:
   - Go runtime initializes before main()
   - Sets up scheduler, GC, and signal handlers
   - GOMAXPROCS determines parallelism level

2. Stack Management:
   - Goroutines start with 2KB stack
   - Stacks grow and shrink dynamically
   - Stack copying during growth is expensive

3. Memory Allocator Integration:
   - Small objects (<32KB) use thread-local caches
   - Large objects go directly to heap
   - sync.Pool reduces allocation pressure

4. System Call Integration:
   - Blocking syscalls don't block OS threads
   - Runtime manages M:N threading model
   - Netpoller handles network I/O efficiently

5. Runtime Coordination:
   - Scheduler coordinates with GC
   - Signal-based preemption (Go 1.14+)
   - Work-stealing maintains balance

6. Production Implications:
   - Understand goroutine lifecycle costs
   - Memory allocation patterns affect performance
   - Runtime tuning affects application behavior
   - Profiling reveals runtime bottlenecks

This deep runtime knowledge is essential for:
- Optimizing high-performance Go applications
- Debugging complex runtime issues
- System design discussions involving Go
- Senior/Staff engineer positions at FAANG companies
*/
