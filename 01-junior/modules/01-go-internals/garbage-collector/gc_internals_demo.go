// Package gc demonstrates Go's garbage collector internals
// Essential knowledge for FAANG backend engineer interviews

package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func main() {
	fmt.Println("🔥 Go Garbage Collector Deep Dive")
	fmt.Println("==================================")

	// Show initial GC state
	showGCStats("Initial State")

	// Demonstrate tricolor marking
	demonstrateTricolorMarking()

	// Show concurrent GC behavior
	demonstrateConcurrentGC()

	// GC tuning examples
	demonstrateGCTuning()

	// Memory pressure scenarios
	demonstrateMemoryPressure()
}

// Demonstrate tricolor marking algorithm
func demonstrateTricolorMarking() {
	fmt.Println("\n🎯 Tricolor Marking Algorithm")
	fmt.Println("-----------------------------")

	// Create objects with different reference patterns
	type Node struct {
		data     [1024]byte // Make it substantial enough to trigger GC
		children []*Node
	}

	// Create a tree structure to demonstrate marking
	root := &Node{}

	// Create referenced objects (will be marked as reachable)
	for i := 0; i < 100; i++ {
		child := &Node{}
		root.children = append(root.children, child)

		// Create grandchildren (deeper references)
		for j := 0; j < 10; j++ {
			grandchild := &Node{}
			child.children = append(child.children, grandchild)
		}
	}

	fmt.Printf("Created tree with root and %d children\n", len(root.children))
	showGCStats("After creating referenced objects")

	// Create unreferenced objects (will be swept)
	for i := 0; i < 1000; i++ {
		_ = &Node{} // These become immediately unreachable
	}

	showGCStats("After creating unreferenced objects")

	// Force GC to see tricolor marking in action
	runtime.GC()
	showGCStats("After forced GC")

	// Keep root alive to prevent optimization
	_ = root
}

// Demonstrate concurrent GC behavior
func demonstrateConcurrentGC() {
	fmt.Println("\n⚡ Concurrent GC Behavior")
	fmt.Println("------------------------")

	// Start memory allocation in background
	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			// Allocate large slices to trigger GC
			data := make([][]byte, 1000)
			for j := range data {
				data[j] = make([]byte, 1024)
			}

			time.Sleep(50 * time.Millisecond)
			// data goes out of scope and becomes eligible for GC
		}
		done <- true
	}()

	// Monitor GC activity while allocation happens
	initialGCCount := getGCCount()

	<-done

	finalGCCount := getGCCount()
	fmt.Printf("GC cycles during concurrent allocation: %d\n",
		finalGCCount-initialGCCount)

	showGCStats("After concurrent allocation")
}

// Demonstrate GC tuning parameters
func demonstrateGCTuning() {
	fmt.Println("\n🔧 GC Tuning Parameters")
	fmt.Println("-----------------------")

	// Show current GC target percentage
	gcPercent := debug.SetGCPercent(-1) // Get current value
	debug.SetGCPercent(gcPercent)       // Restore it

	fmt.Printf("Current GOGC (GC target percentage): %d%%\n", gcPercent)

	// Demonstrate effect of different GC percentages
	fmt.Println("\nTesting different GOGC values:")

	testGCPercent(50)  // Aggressive GC
	testGCPercent(100) // Default GC
	testGCPercent(200) // Conservative GC

	// Restore default
	debug.SetGCPercent(100)
}

func testGCPercent(percent int) {
	debug.SetGCPercent(percent)

	start := time.Now()
	initialGC := getGCCount()

	// Allocate memory
	var data [][]byte
	for i := 0; i < 1000; i++ {
		data = append(data, make([]byte, 1024))
	}

	finalGC := getGCCount()
	duration := time.Since(start)

	fmt.Printf("GOGC=%d%% -> GC cycles: %d, Duration: %v\n",
		percent, finalGC-initialGC, duration)

	// Keep data alive during test
	_ = data
}

// Demonstrate memory pressure scenarios
func demonstrateMemoryPressure() {
	fmt.Println("\n💾 Memory Pressure Scenarios")
	fmt.Println("----------------------------")

	showGCStats("Before memory pressure")

	// Scenario 1: Rapid allocation and deallocation
	fmt.Println("Scenario 1: Rapid allocation/deallocation")
	for i := 0; i < 100; i++ {
		data := make([]byte, 1024*1024) // 1MB allocation
		_ = data                        // Use it briefly
		// Goes out of scope immediately
	}

	runtime.GC() // Force collection
	showGCStats("After rapid alloc/dealloc")

	// Scenario 2: Long-lived large objects
	fmt.Println("\nScenario 2: Long-lived large objects")
	longLived := make([][]byte, 100)
	for i := range longLived {
		longLived[i] = make([]byte, 1024*1024) // 1MB each
	}

	runtime.GC()
	showGCStats("After creating long-lived objects")

	// Scenario 3: Mixed allocation patterns
	fmt.Println("\nScenario 3: Mixed allocation patterns")
	go func() {
		for i := 0; i < 1000; i++ {
			_ = make([]byte, 1024) // Short-lived
			time.Sleep(time.Millisecond)
		}
	}()

	time.Sleep(500 * time.Millisecond)
	runtime.GC()
	showGCStats("After mixed allocation patterns")

	// Keep long-lived data alive
	_ = longLived
}

// Helper function to get GC count
func getGCCount() uint32 {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return stats.NumGC
}

// Helper function to show GC statistics
func showGCStats(label string) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	fmt.Printf("\n📊 %s:\n", label)
	fmt.Printf("  GC Cycles: %d\n", stats.NumGC)
	fmt.Printf("  Heap Size: %.2f MB\n", float64(stats.HeapAlloc)/1024/1024)
	fmt.Printf("  Heap Objects: %d\n", stats.HeapObjects)
	fmt.Printf("  GC Pause Total: %v\n", time.Duration(stats.PauseTotalNs))

	if stats.NumGC > 0 {
		recentPause := stats.PauseNs[(stats.NumGC+255)%256]
		fmt.Printf("  Recent GC Pause: %v\n", time.Duration(recentPause))
	}
}

/*
🔥 FAANG Interview Key Points:

1. Tricolor Marking Algorithm:
   - White: Objects not yet scanned
   - Gray: Objects scanned but children not scanned
   - Black: Objects and all children scanned
   - Only white objects are collected

2. Concurrent Collection:
   - GC runs concurrently with application
   - Write barriers ensure consistency during concurrent marking
   - STW (Stop-The-World) phases are minimized

3. GC Triggers:
   - Heap growth threshold (GOGC percentage)
   - Manual runtime.GC() calls
   - Memory pressure from OS

4. Performance Tuning:
   - GOGC environment variable (default 100%)
   - Lower values = more frequent GC, less memory
   - Higher values = less frequent GC, more memory

5. Memory Management Best Practices:
   - Reuse objects when possible (sync.Pool)
   - Avoid frequent large allocations
   - Use appropriate data structures for access patterns
   - Monitor GC metrics in production

6. Production Monitoring:
   - Track GC pause times
   - Monitor heap growth patterns
   - Alert on excessive GC frequency
   - Use pprof for memory profiling

This deep understanding of Go's GC is essential for:
- Optimizing high-performance Go applications
- Debugging memory-related issues
- System design interviews involving performance
- Senior backend engineer roles at scale
*/
