// Package gospecific implements Go-specific complexity analysis
// Understanding how Go runtime features affect theoretical complexity
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 🔧 GO-SPECIFIC COMPLEXITY CONSIDERATIONS
// How Go runtime affects algorithmic complexity analysis

// 📊 SLICE COMPLEXITY ANALYSIS

// SliceGrowthPattern analyzes Go slice append complexity
type SliceGrowthPattern struct {
	InitialCap    int
	FinalSize     int
	Appends       int
	Reallocations int
	TotalCopies   int64
}

// AnalyzeSliceGrowth demonstrates amortized O(1) append complexity
func AnalyzeSliceGrowth(targetSize int) SliceGrowthPattern {
	pattern := SliceGrowthPattern{InitialCap: 0}

	var slice []int
	totalCopies := int64(0)
	reallocations := 0

	for i := 0; i < targetSize; i++ {
		oldCap := cap(slice)
		slice = append(slice, i)
		newCap := cap(slice)

		// Track reallocation
		if newCap > oldCap {
			reallocations++
			totalCopies += int64(len(slice) - 1) // All previous elements copied
		}
	}

	pattern.FinalSize = len(slice)
	pattern.Appends = targetSize
	pattern.Reallocations = reallocations
	pattern.TotalCopies = totalCopies

	return pattern
}

// SliceOperationComplexity demonstrates slice operation complexities
func SliceOperationComplexity() {
	fmt.Println("📊 GO SLICE OPERATION COMPLEXITIES")
	fmt.Println(strings.Repeat("=", 60))

	operations := []struct {
		operation   string
		complexity  string
		explanation string
		example     string
	}{
		{
			operation:   "append() - amortized",
			complexity:  "O(1)",
			explanation: "Geometric growth ensures amortized constant time",
			example:     "slice = append(slice, value)",
		},
		{
			operation:   "append() - worst case",
			complexity:  "O(n)",
			explanation: "When reallocation occurs, all elements are copied",
			example:     "slice = append(slice, value) // when cap exceeded",
		},
		{
			operation:   "copy()",
			complexity:  "O(min(len(dst), len(src)))",
			explanation: "Copies min(len(dst), len(src)) elements",
			example:     "copy(dst, src)",
		},
		{
			operation:   "Access by index",
			complexity:  "O(1)",
			explanation: "Direct memory access with bounds checking",
			example:     "value := slice[index]",
		},
		{
			operation:   "Slice creation [i:j]",
			complexity:  "O(1)",
			explanation: "Creates new slice header, shares underlying array",
			example:     "subSlice := slice[i:j]",
		},
	}

	for _, op := range operations {
		fmt.Printf("\n%s: %s\n", op.operation, op.complexity)
		fmt.Printf("  %s\n", op.explanation)
		fmt.Printf("  Example: %s\n", op.example)
	}
}

// 🗺️ MAP COMPLEXITY ANALYSIS

// MapComplexityAnalysis demonstrates Go map operation complexities
type MapComplexityAnalysis struct {
	LoadFactor    float64
	Collisions    int
	Rehashes      int
	AverageProbes float64
}

// AnalyzeMapComplexity demonstrates hash map complexity in Go
func AnalyzeMapComplexity(size int) MapComplexityAnalysis {
	hashMap := make(map[string]int, size)

	// Fill map with predictable collision patterns
	collisions := 0
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("key_%d", i)
		hashMap[key] = i

		// Simulate collision detection (simplified)
		if i > 0 && (i%100) == 0 {
			collisions++
		}
	}

	loadFactor := float64(len(hashMap)) / float64(size)

	return MapComplexityAnalysis{
		LoadFactor:    loadFactor,
		Collisions:    collisions,
		Rehashes:      int(loadFactor * 2), // Simplified estimation
		AverageProbes: 1.0 + loadFactor/2,  // Expected probes in open addressing
	}
}

// MapOperationComplexity demonstrates map operation complexities
func MapOperationComplexity() {
	fmt.Println("\n🗺️ GO MAP OPERATION COMPLEXITIES")
	fmt.Println(strings.Repeat("=", 60))

	operations := []struct {
		operation   string
		avgCase     string
		worstCase   string
		explanation string
	}{
		{
			operation:   "Insert/Update m[key] = value",
			avgCase:     "O(1)",
			worstCase:   "O(n)",
			explanation: "Average O(1), worst O(n) during rehashing",
		},
		{
			operation:   "Lookup value := m[key]",
			avgCase:     "O(1)",
			worstCase:   "O(n)",
			explanation: "Hash collision chains can degrade to O(n)",
		},
		{
			operation:   "Delete delete(m, key)",
			avgCase:     "O(1)",
			worstCase:   "O(n)",
			explanation: "Similar to lookup, depends on collision resolution",
		},
		{
			operation:   "Range iteration",
			avgCase:     "O(n)",
			worstCase:   "O(n)",
			explanation: "Must visit all key-value pairs",
		},
		{
			operation:   "len(m)",
			avgCase:     "O(1)",
			worstCase:   "O(1)",
			explanation: "Constant time, maintained as metadata",
		},
	}

	fmt.Printf("%-25s | %-10s | %-10s | %s\n", "Operation", "Avg Case", "Worst Case", "Notes")
	fmt.Println(strings.Repeat("-", 80))

	for _, op := range operations {
		fmt.Printf("%-25s | %-10s | %-10s | %s\n",
			op.operation, op.avgCase, op.worstCase, op.explanation)
	}
}

// 🚀 GOROUTINE COMPLEXITY ANALYSIS

// GoroutineOverhead measures goroutine creation and switching costs
type GoroutineOverhead struct {
	CreationTime       time.Duration
	MemoryPerGoroutine int64
	StackSize          int64
	ContextSwitchTime  time.Duration
}

// MeasureGoroutineOverhead analyzes goroutine performance characteristics
func MeasureGoroutineOverhead(numGoroutines int) GoroutineOverhead {
	var wg sync.WaitGroup

	// Measure creation time
	start := time.Now()
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Minimal work
			_ = 1 + 1
		}()
	}
	wg.Wait()
	creationTime := time.Since(start)

	// Estimate memory usage (simplified)
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Create goroutines that stay alive briefly
	done := make(chan bool)
	for i := 0; i < 1000; i++ {
		go func() {
			<-done
		}()
	}

	runtime.ReadMemStats(&m2)
	close(done)

	memoryPerGoroutine := int64(m2.Alloc-m1.Alloc) / 1000

	return GoroutineOverhead{
		CreationTime:       creationTime / time.Duration(numGoroutines),
		MemoryPerGoroutine: memoryPerGoroutine,
		StackSize:          2048,             // Default Go stack size
		ContextSwitchTime:  time.Microsecond, // Approximate
	}
}

// ChannelComplexity demonstrates channel operation complexities
func ChannelComplexity() {
	fmt.Println("\n📡 GO CHANNEL OPERATION COMPLEXITIES")
	fmt.Println(strings.Repeat("=", 60))

	operations := []struct {
		operation   string
		complexity  string
		explanation string
	}{
		{
			operation:   "Buffered channel send (space available)",
			complexity:  "O(1)",
			explanation: "Direct write to buffer slot",
		},
		{
			operation:   "Buffered channel send (buffer full)",
			complexity:  "O(1) + block",
			explanation: "Constant time + potential goroutine blocking",
		},
		{
			operation:   "Unbuffered channel send/receive",
			complexity:  "O(1) + sync",
			explanation: "Constant time + synchronization overhead",
		},
		{
			operation:   "Channel close",
			complexity:  "O(g)",
			explanation: "Notify all g waiting goroutines",
		},
		{
			operation:   "Select statement",
			complexity:  "O(n)",
			explanation: "Check n channel cases (pseudo-random order)",
		},
	}

	for _, op := range operations {
		fmt.Printf("\n%s: %s\n", op.operation, op.complexity)
		fmt.Printf("  %s\n", op.explanation)
	}
}

// 🗑️ GARBAGE COLLECTOR IMPACT ANALYSIS

// GCImpactAnalysis measures garbage collection effects on performance
type GCImpactAnalysis struct {
	AllocationRate   float64 // MB/s
	GCFrequency      time.Duration
	PauseTime        time.Duration
	HeapSize         int64
	ThroughputImpact float64 // Percentage
}

// AnalyzeGCImpact demonstrates how GC affects algorithmic performance
func AnalyzeGCImpact(allocSize int) GCImpactAnalysis {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	start := time.Now()

	// Perform allocations that will trigger GC
	data := make([][]byte, 0, allocSize)
	for i := 0; i < allocSize; i++ {
		// Allocate 1KB chunks
		chunk := make([]byte, 1024)
		rand.Read(chunk) // Fill with random data
		data = append(data, chunk)

		// Force some garbage creation
		if i%100 == 0 {
			_ = make([]int, 1000) // This will be garbage collected
		}
	}

	elapsed := time.Since(start)
	runtime.ReadMemStats(&m2)

	allocationRate := float64(m2.TotalAlloc-m1.TotalAlloc) / 1024 / 1024 / elapsed.Seconds()
	avgPause := time.Duration(m2.PauseTotalNs / uint64(m2.NumGC))

	return GCImpactAnalysis{
		AllocationRate:   allocationRate,
		GCFrequency:      time.Duration(elapsed.Nanoseconds() / int64(m2.NumGC-m1.NumGC)),
		PauseTime:        avgPause,
		HeapSize:         int64(m2.HeapInuse),
		ThroughputImpact: float64(m2.PauseTotalNs) / float64(elapsed.Nanoseconds()) * 100,
	}
}

// 📏 MEMORY LAYOUT IMPACT ON COMPLEXITY

// CacheEffectAnalysis demonstrates cache-friendly vs cache-unfriendly access patterns
type CacheEffectAnalysis struct {
	SequentialTime time.Duration
	RandomTime     time.Duration
	Speedup        float64
}

// AnalyzeCacheEffects compares sequential vs random memory access
func AnalyzeCacheEffects(size int) CacheEffectAnalysis {
	data := make([]int64, size)
	for i := range data {
		data[i] = int64(i)
	}

	// Sequential access
	start := time.Now()
	sum1 := int64(0)
	for i := 0; i < len(data); i++ {
		sum1 += data[i]
	}
	sequentialTime := time.Since(start)

	// Random access
	indices := make([]int, size)
	for i := range indices {
		indices[i] = rand.Intn(size)
	}

	start = time.Now()
	sum2 := int64(0)
	for _, idx := range indices {
		sum2 += data[idx]
	}
	randomTime := time.Since(start)

	// Prevent optimization
	_ = sum1 + sum2

	speedup := float64(randomTime) / float64(sequentialTime)

	return CacheEffectAnalysis{
		SequentialTime: sequentialTime,
		RandomTime:     randomTime,
		Speedup:        speedup,
	}
}

// 🔍 PRACTICAL GO OPTIMIZATION EXAMPLES

// OptimizationExample demonstrates complexity improvements through Go-specific optimizations
func OptimizationExample() {
	fmt.Println("\n🔍 GO-SPECIFIC OPTIMIZATION EXAMPLES")
	fmt.Println(strings.Repeat("=", 60))

	examples := []struct {
		scenario    string
		before      string
		after       string
		improvement string
		explanation string
	}{
		{
			scenario:    "String concatenation in loop",
			before:      "O(n²)",
			after:       "O(n)",
			improvement: "strings.Builder",
			explanation: "strings.Builder avoids O(n) copies on each append",
		},
		{
			scenario:    "Slice preallocation",
			before:      "O(n log n)",
			after:       "O(n)",
			improvement: "make([]T, 0, capacity)",
			explanation: "Preallocating prevents multiple reallocations",
		},
		{
			scenario:    "Map preallocation",
			before:      "O(n log n)",
			after:       "O(n)",
			improvement: "make(map[K]V, size)",
			explanation: "Prevents rehashing during population",
		},
		{
			scenario:    "Pointer vs value in large structs",
			before:      "O(n * structSize)",
			after:       "O(n * pointerSize)",
			improvement: "*BigStruct instead of BigStruct",
			explanation: "Reduces copy overhead for large data structures",
		},
	}

	for _, ex := range examples {
		fmt.Printf("\n%s:\n", ex.scenario)
		fmt.Printf("  Before: %s → After: %s\n", ex.before, ex.after)
		fmt.Printf("  Technique: %s\n", ex.improvement)
		fmt.Printf("  Why: %s\n", ex.explanation)
	}
}

// Main demonstration function
func main() {
	fmt.Println("🔧 GO-SPECIFIC COMPLEXITY ANALYSIS")
	fmt.Println("Understanding how Go runtime affects algorithmic complexity")
	fmt.Println(strings.Repeat("=", 80))

	// Slice analysis
	fmt.Println("\n📊 SLICE GROWTH ANALYSIS")
	pattern := AnalyzeSliceGrowth(10000)
	fmt.Printf("Appending %d elements:\n", pattern.Appends)
	fmt.Printf("  Reallocations: %d\n", pattern.Reallocations)
	fmt.Printf("  Total copies: %d\n", pattern.TotalCopies)
	fmt.Printf("  Amortized copies per append: %.2f\n",
		float64(pattern.TotalCopies)/float64(pattern.Appends))

	// Slice operations
	SliceOperationComplexity()

	// Map analysis
	fmt.Println("\n🗺️ MAP ANALYSIS")
	mapAnalysis := AnalyzeMapComplexity(10000)
	fmt.Printf("Map with 10,000 elements:\n")
	fmt.Printf("  Load factor: %.2f\n", mapAnalysis.LoadFactor)
	fmt.Printf("  Estimated collisions: %d\n", mapAnalysis.Collisions)
	fmt.Printf("  Average probes: %.2f\n", mapAnalysis.AverageProbes)

	MapOperationComplexity()

	// Goroutine analysis
	fmt.Println("\n🚀 GOROUTINE OVERHEAD ANALYSIS")
	overhead := MeasureGoroutineOverhead(1000)
	fmt.Printf("Per-goroutine metrics:\n")
	fmt.Printf("  Creation time: %v\n", overhead.CreationTime)
	fmt.Printf("  Memory overhead: %d bytes\n", overhead.MemoryPerGoroutine)
	fmt.Printf("  Stack size: %d bytes\n", overhead.StackSize)

	ChannelComplexity()

	// GC impact analysis
	fmt.Println("\n🗑️ GARBAGE COLLECTION IMPACT")
	gcAnalysis := AnalyzeGCImpact(10000)
	fmt.Printf("GC performance with heavy allocation:\n")
	fmt.Printf("  Allocation rate: %.2f MB/s\n", gcAnalysis.AllocationRate)
	fmt.Printf("  Average pause: %v\n", gcAnalysis.PauseTime)
	fmt.Printf("  Throughput impact: %.2f%%\n", gcAnalysis.ThroughputImpact)

	// Cache effects
	fmt.Println("\n📏 MEMORY ACCESS PATTERNS")
	cacheAnalysis := AnalyzeCacheEffects(1000000)
	fmt.Printf("Sequential vs Random access (1M elements):\n")
	fmt.Printf("  Sequential: %v\n", cacheAnalysis.SequentialTime)
	fmt.Printf("  Random: %v\n", cacheAnalysis.RandomTime)
	fmt.Printf("  Random slowdown: %.2fx\n", cacheAnalysis.Speedup)

	// Optimization examples
	OptimizationExample()

	fmt.Println("\n✅ Go-specific complexity analysis complete!")
	fmt.Println("💡 Key insight: Go runtime characteristics can significantly impact")
	fmt.Println("   theoretical complexity bounds in practice!")
}
