// Package memory demonstrates Go's memory management internals
// Critical knowledge for FAANG performance optimization interviews

package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

// Memory Management Core Concepts:
// - Stack vs Heap allocation decisions
// - Escape analysis and compiler optimizations
// - Memory alignment and padding
// - Memory pool patterns and optimizations
// - Memory profiling and leak detection

func main() {
	fmt.Println("🔥 Go Memory Management Deep Dive")
	fmt.Println("=================================")

	// Stack vs Heap allocation patterns
	demonstrateStackVsHeap()

	// Escape analysis in action
	demonstrateEscapeAnalysis()

	// Memory alignment and padding
	demonstrateMemoryAlignment()

	// Memory pools and optimization
	demonstrateMemoryPools()

	// Memory profiling techniques
	demonstrateMemoryProfiling()

	// Advanced memory patterns
	demonstrateAdvancedPatterns()
}

// Stack vs Heap allocation decision factors
func demonstrateStackVsHeap() {
	fmt.Println("\n📚 Stack vs Heap Allocation")
	fmt.Println("---------------------------")

	// Show memory stats before allocations
	showMemoryStats("Initial state")

	// Stack allocation examples
	demonstrateStackAllocations()

	// Heap allocation examples
	demonstrateHeapAllocations()

	// Mixed allocation patterns
	demonstrateMixedAllocations()

	showMemoryStats("After allocations")
}

// Examples of stack allocations (escape analysis keeps them on stack)
func demonstrateStackAllocations() {
	fmt.Println("Stack allocation examples:")

	// Local variables that don't escape
	stackExample1()
	stackExample2()
	stackExample3()
}

// Simple local variable - stays on stack
func stackExample1() {
	var x int = 42
	var arr [10]int

	// These operations keep data on stack
	for i := range arr {
		arr[i] = x + i
	}

	// Calculate sum (stays on stack)
	sum := 0
	for _, v := range arr {
		sum += v
	}

	fmt.Printf("Stack example 1 - sum: %d\n", sum)
	// x, arr, sum all allocated on stack and cleaned up automatically
}

// Local struct that doesn't escape
func stackExample2() {
	type Point struct {
		X, Y float64
	}

	// Struct allocated on stack
	p := Point{X: 3.0, Y: 4.0}

	// Operations on stack-allocated struct
	distance := func(p Point) float64 {
		return p.X*p.X + p.Y*p.Y // sqrt omitted for simplicity
	}(p)

	fmt.Printf("Stack example 2 - distance: %.2f\n", distance)
	// Point struct stays on stack
}

// Slice with known size that doesn't escape
func stackExample3() {
	// Small slice that might stay on stack (compiler dependent)
	nums := make([]int, 5)
	for i := range nums {
		nums[i] = i * i
	}

	// Process without escaping
	var total int
	for _, num := range nums {
		total += num
	}

	fmt.Printf("Stack example 3 - total: %d\n", total)
}

// Examples of heap allocations (escape analysis forces heap)
func demonstrateHeapAllocations() {
	fmt.Println("Heap allocation examples:")

	// Returning pointer causes escape to heap
	ptr1 := heapExample1()
	fmt.Printf("Heap example 1 - value: %d\n", *ptr1)

	// Large objects go to heap
	slice2 := heapExample2()
	fmt.Printf("Heap example 2 - length: %d\n", len(slice2))

	// Interface assignment causes escape
	intf3 := heapExample3()
	fmt.Printf("Heap example 3 - value: %v\n", intf3)

	// Closure capture causes escape
	fn4 := heapExample4()
	fmt.Printf("Heap example 4 - result: %d\n", fn4())
}

// Returning pointer to local variable - escapes to heap
func heapExample1() *int {
	x := 42 // This will be allocated on heap due to escape
	return &x
}

// Large slice allocation - goes to heap
func heapExample2() []int {
	// Large allocation always goes to heap
	return make([]int, 10000)
}

// Interface assignment - causes escape to heap
func heapExample3() interface{} {
	x := 123 // Escapes to heap when assigned to interface{}
	return x
}

// Closure capturing variable - causes escape to heap
func heapExample4() func() int {
	counter := 0 // Escapes to heap due to closure capture
	return func() int {
		counter++
		return counter
	}
}

// Mixed allocation patterns showing decision factors
func demonstrateMixedAllocations() {
	fmt.Println("Mixed allocation patterns:")

	// Pattern 1: Local processing (stack) vs returning (heap)
	localResult := processLocally([]int{1, 2, 3, 4, 5})
	heapResult := processAndReturn([]int{1, 2, 3, 4, 5})

	fmt.Printf("Local processing result: %d\n", localResult)
	fmt.Printf("Heap processing result: %d\n", len(heapResult))
}

// Processes data locally without escaping
func processLocally(data []int) int {
	// Local variables stay on stack
	sum := 0
	for _, v := range data {
		sum += v
	}
	return sum // Only the result escapes, not intermediate data
}

// Returns slice - causes escape to heap
func processAndReturn(data []int) []int {
	// This slice will be allocated on heap
	result := make([]int, len(data))
	for i, v := range data {
		result[i] = v * 2
	}
	return result // Entire slice escapes to heap
}

// Escape analysis demonstration with compiler flags
func demonstrateEscapeAnalysis() {
	fmt.Println("\n🔍 Escape Analysis Demonstration")
	fmt.Println("-------------------------------")

	fmt.Println("Run with: go run -gcflags='-m' memory_management.go")
	fmt.Println("To see escape analysis decisions.")
	fmt.Println()

	// Various escape analysis scenarios
	demonstrateEscapeScenarios()
	demonstrateEscapePrevention()
}

// Different scenarios that trigger escape analysis
func demonstrateEscapeScenarios() {
	fmt.Println("Escape analysis scenarios:")

	// Scenario 1: Taking address of local variable
	scenario1()

	// Scenario 2: Slice growth beyond initial capacity
	scenario2()

	// Scenario 3: Interface conversion
	scenario3()

	// Scenario 4: Sending to channel
	scenario4()
}

func scenario1() {
	type Data struct {
		value int
	}

	d := Data{value: 42}
	ptr := &d // Taking address - may cause escape

	fmt.Printf("Scenario 1 - Address: %p, Value: %d\n", ptr, ptr.value)
}

func scenario2() {
	// Start with small slice
	slice := make([]int, 0, 2)

	// Add elements beyond capacity - may cause escape
	for i := 0; i < 5; i++ {
		slice = append(slice, i)
	}

	fmt.Printf("Scenario 2 - Slice length: %d, capacity: %d\n", len(slice), cap(slice))
}

func scenario3() {
	value := 123
	var intf interface{} = value // Interface conversion - causes escape

	fmt.Printf("Scenario 3 - Interface value: %v\n", intf)
}

func scenario4() {
	ch := make(chan *int, 1)

	x := 42
	ch <- &x // Sending pointer to channel - causes escape

	received := <-ch
	fmt.Printf("Scenario 4 - Received: %d\n", *received)
}

// Techniques to prevent unnecessary escapes
func demonstrateEscapePrevention() {
	fmt.Println("Escape prevention techniques:")

	// Technique 1: Use values instead of pointers when possible
	preventEscape1()

	// Technique 2: Avoid interface{} when type is known
	preventEscape2()

	// Technique 3: Use arrays instead of slices for fixed size
	preventEscape3()
}

func preventEscape1() {
	type Point struct {
		X, Y float64
	}

	// Pass by value to avoid escape
	calculateDistance := func(p Point) float64 {
		return p.X*p.X + p.Y*p.Y
	}

	p := Point{X: 3, Y: 4}
	dist := calculateDistance(p) // Passes copy, original stays on stack

	fmt.Printf("Prevention 1 - Distance: %.2f\n", dist)
}

func preventEscape2() {
	// Use specific types instead of interface{}
	processInt := func(x int) int {
		return x * 2
	}

	value := 42
	result := processInt(value) // No interface conversion, stays on stack

	fmt.Printf("Prevention 2 - Result: %d\n", result)
}

func preventEscape3() {
	// Use array instead of slice for fixed-size data
	var arr [5]int
	for i := range arr {
		arr[i] = i * i
	}

	// Process array without escaping
	sum := 0
	for _, v := range arr {
		sum += v
	}

	fmt.Printf("Prevention 3 - Array sum: %d\n", sum)
}

// Memory alignment and padding demonstration
func demonstrateMemoryAlignment() {
	fmt.Println("\n📐 Memory Alignment & Padding")
	fmt.Println("----------------------------")

	// Demonstrate struct alignment
	demonstrateStructAlignment()

	// Demonstrate optimal struct ordering
	demonstrateOptimalOrdering()

	// Demonstrate padding effects
	demonstratePaddingEffects()
}

// Show how struct fields are aligned in memory
func demonstrateStructAlignment() {
	fmt.Println("Struct alignment examples:")

	// Poorly aligned struct
	type PoorlyAligned struct {
		a int8  // 1 byte
		b int64 // 8 bytes, will be aligned to 8-byte boundary
		c int8  // 1 byte
		d int32 // 4 bytes, will be aligned to 4-byte boundary
	}

	// Well aligned struct
	type WellAligned struct {
		b int64 // 8 bytes
		d int32 // 4 bytes
		a int8  // 1 byte
		c int8  // 1 byte
		// Compiler will add 2 bytes padding here for 8-byte alignment
	}

	var poor PoorlyAligned
	var well WellAligned

	fmt.Printf("Poorly aligned struct size: %d bytes\n", unsafe.Sizeof(poor))
	fmt.Printf("Well aligned struct size: %d bytes\n", unsafe.Sizeof(well))

	// Show field offsets
	fmt.Printf("Poor alignment offsets: a=%d, b=%d, c=%d, d=%d\n",
		unsafe.Offsetof(poor.a), unsafe.Offsetof(poor.b),
		unsafe.Offsetof(poor.c), unsafe.Offsetof(poor.d))

	fmt.Printf("Good alignment offsets: b=%d, d=%d, a=%d, c=%d\n",
		unsafe.Offsetof(well.b), unsafe.Offsetof(well.d),
		unsafe.Offsetof(well.a), unsafe.Offsetof(well.c))
}

// Demonstrate optimal struct field ordering
func demonstrateOptimalOrdering() {
	fmt.Println("Optimal struct ordering:")

	// Before optimization (24 bytes)
	type Before struct {
		flag1 bool  // 1 byte + 7 bytes padding
		value int64 // 8 bytes
		flag2 bool  // 1 byte + 7 bytes padding
	}

	// After optimization (16 bytes)
	type After struct {
		value int64 // 8 bytes
		flag1 bool  // 1 byte
		flag2 bool  // 1 byte + 6 bytes padding
	}

	var before Before
	var after After

	fmt.Printf("Before optimization: %d bytes\n", unsafe.Sizeof(before))
	fmt.Printf("After optimization: %d bytes\n", unsafe.Sizeof(after))
	fmt.Printf("Memory savings: %d bytes (%.1f%%)\n",
		int(unsafe.Sizeof(before)-unsafe.Sizeof(after)),
		100.0*float64(unsafe.Sizeof(before)-unsafe.Sizeof(after))/float64(unsafe.Sizeof(before)))
}

// Show effects of padding on performance
func demonstratePaddingEffects() {
	fmt.Println("Padding effects on performance:")

	// Test struct with padding
	type WithPadding struct {
		a int8
		b int64
		c int8
	}

	// Test struct without padding
	type WithoutPadding struct {
		a int64
		b int8
		c int8
	}

	const iterations = 1000000

	// Benchmark with padding
	start := time.Now()
	var withPadding [iterations]WithPadding
	for i := 0; i < iterations; i++ {
		withPadding[i] = WithPadding{a: 1, b: 2, c: 3}
	}
	paddingTime := time.Since(start)

	// Benchmark without padding
	start = time.Now()
	var withoutPadding [iterations]WithoutPadding
	for i := 0; i < iterations; i++ {
		withoutPadding[i] = WithoutPadding{a: 2, b: 1, c: 3}
	}
	noPaddingTime := time.Since(start)

	fmt.Printf("With padding (%d bytes): %v\n", unsafe.Sizeof(withPadding[0]), paddingTime)
	fmt.Printf("Without padding (%d bytes): %v\n", unsafe.Sizeof(withoutPadding[0]), noPaddingTime)

	// Keep arrays alive
	_ = withPadding
	_ = withoutPadding
}

// Memory pools and optimization patterns
func demonstrateMemoryPools() {
	fmt.Println("\n🏊 Memory Pools & Optimization")
	fmt.Println("-----------------------------")

	// Object pooling to reduce GC pressure
	demonstrateObjectPooling()

	// Arena allocation patterns
	demonstrateArenaAllocation()

	// Buffer reuse patterns
	demonstrateBufferReuse()
}

// Object pooling to reduce allocations
func demonstrateObjectPooling() {
	fmt.Println("Object pooling demonstration:")

	type ExpensiveObject struct {
		data [1024]byte
		id   int
	}

	showMemoryStats("Before pooling test")

	// Without pooling - many allocations
	start := time.Now()
	for i := 0; i < 1000; i++ {
		obj := &ExpensiveObject{id: i}
		// Simulate work
		obj.data[0] = byte(i)
		// Object becomes garbage immediately
	}
	withoutPoolingTime := time.Since(start)

	showMemoryStats("After non-pooled allocations")

	// Force GC to see the difference
	runtime.GC()
	showMemoryStats("After GC")

	// With pooling - reuse objects
	pool := make([]*ExpensiveObject, 0, 100)

	start = time.Now()
	for i := 0; i < 1000; i++ {
		var obj *ExpensiveObject

		// Try to reuse from pool
		if len(pool) > 0 {
			obj = pool[len(pool)-1]
			pool = pool[:len(pool)-1]
		} else {
			obj = &ExpensiveObject{}
		}

		obj.id = i
		obj.data[0] = byte(i)

		// Return to pool
		if len(pool) < cap(pool) {
			pool = append(pool, obj)
		}
	}
	withPoolingTime := time.Since(start)

	fmt.Printf("Without pooling: %v\n", withoutPoolingTime)
	fmt.Printf("With pooling: %v\n", withPoolingTime)
	fmt.Printf("Performance improvement: %.2fx\n",
		float64(withoutPoolingTime)/float64(withPoolingTime))

	showMemoryStats("After pooled allocations")
}

// Arena allocation pattern for batch processing
func demonstrateArenaAllocation() {
	fmt.Println("Arena allocation pattern:")

	// Arena holds all allocations for a batch
	type Arena struct {
		data   []byte
		offset int
	}

	// Allocate from arena
	allocateFromArena := func(arena *Arena, size int) []byte {
		if arena.offset+size > len(arena.data) {
			// Arena full - in production, might allocate new arena
			return nil
		}

		result := arena.data[arena.offset : arena.offset+size]
		arena.offset += size
		return result
	}

	// Create arena for batch processing
	arena := &Arena{
		data: make([]byte, 1024*1024), // 1MB arena
	}

	showMemoryStats("Before arena allocation")

	// Allocate many small objects from arena
	objects := make([][]byte, 1000)
	for i := range objects {
		objects[i] = allocateFromArena(arena, 512) // 512 byte objects
	}

	fmt.Printf("Allocated %d objects from arena\n", len(objects))
	fmt.Printf("Arena utilization: %d/%d bytes (%.1f%%)\n",
		arena.offset, len(arena.data),
		100.0*float64(arena.offset)/float64(len(arena.data)))

	showMemoryStats("After arena allocation")

	// In arena pattern, all objects are freed together when arena is released
	arena = nil // Free entire arena at once
	objects = nil

	runtime.GC()
	showMemoryStats("After arena release")
}

// Buffer reuse patterns
func demonstrateBufferReuse() {
	fmt.Println("Buffer reuse patterns:")

	// Buffer pool for reusing byte slices
	type BufferPool struct {
		buffers [][]byte
		size    int
	}

	newBufferPool := func(bufferSize int) *BufferPool {
		return &BufferPool{
			buffers: make([][]byte, 0, 10),
			size:    bufferSize,
		}
	}

	getBuffer := func(pool *BufferPool) []byte {
		if len(pool.buffers) > 0 {
			// Reuse existing buffer
			buffer := pool.buffers[len(pool.buffers)-1]
			pool.buffers = pool.buffers[:len(pool.buffers)-1]
			return buffer[:0] // Reset length but keep capacity
		}
		// Allocate new buffer
		return make([]byte, 0, pool.size)
	}

	putBuffer := func(pool *BufferPool, buffer []byte) {
		if cap(buffer) == pool.size && len(pool.buffers) < cap(pool.buffers) {
			pool.buffers = append(pool.buffers, buffer)
		}
	}

	pool := newBufferPool(1024)

	showMemoryStats("Before buffer reuse test")

	// Simulate many operations requiring buffers
	for i := 0; i < 1000; i++ {
		buffer := getBuffer(pool)

		// Simulate work - append some data
		for j := 0; j < 100; j++ {
			buffer = append(buffer, byte(j))
		}

		// Return buffer to pool
		putBuffer(pool, buffer)
	}

	fmt.Printf("Buffer pool utilization: %d buffers cached\n", len(pool.buffers))
	showMemoryStats("After buffer reuse test")
}

// Memory profiling techniques
func demonstrateMemoryProfiling() {
	fmt.Println("\n📊 Memory Profiling Techniques")
	fmt.Println("-----------------------------")

	// Runtime memory statistics
	demonstrateMemoryStats()

	// Allocation tracking
	demonstrateAllocationTracking()

	// Memory leak detection patterns
	demonstrateMemoryLeakDetection()
}

// Detailed memory statistics analysis
func demonstrateMemoryStats() {
	fmt.Println("Detailed memory statistics:")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("General statistics:\n")
	fmt.Printf("  Alloc: %d KB (bytes allocated and not yet freed)\n", m.Alloc/1024)
	fmt.Printf("  TotalAlloc: %d KB (cumulative bytes allocated)\n", m.TotalAlloc/1024)
	fmt.Printf("  Sys: %d KB (bytes obtained from system)\n", m.Sys/1024)
	fmt.Printf("  Lookups: %d (number of pointer lookups)\n", m.Lookups)
	fmt.Printf("  Mallocs: %d (number of mallocs)\n", m.Mallocs)
	fmt.Printf("  Frees: %d (number of frees)\n", m.Frees)

	fmt.Printf("\nHeap statistics:\n")
	fmt.Printf("  HeapAlloc: %d KB (bytes allocated and not yet freed)\n", m.HeapAlloc/1024)
	fmt.Printf("  HeapSys: %d KB (bytes obtained from system)\n", m.HeapSys/1024)
	fmt.Printf("  HeapIdle: %d KB (bytes in idle spans)\n", m.HeapIdle/1024)
	fmt.Printf("  HeapInuse: %d KB (bytes in non-idle spans)\n", m.HeapInuse/1024)
	fmt.Printf("  HeapReleased: %d KB (bytes released to OS)\n", m.HeapReleased/1024)
	fmt.Printf("  HeapObjects: %d (number of allocated objects)\n", m.HeapObjects)

	fmt.Printf("\nGC statistics:\n")
	fmt.Printf("  NextGC: %d KB (next GC target)\n", m.NextGC/1024)
	fmt.Printf("  LastGC: %v (last GC time)\n", time.Unix(0, int64(m.LastGC)))
	fmt.Printf("  PauseTotalNs: %v (total pause time)\n", time.Duration(m.PauseTotalNs))
	fmt.Printf("  NumGC: %d (number of GC cycles)\n", m.NumGC)
	fmt.Printf("  GCCPUFraction: %.4f (fraction of CPU time used by GC)\n", m.GCCPUFraction)
}

// Track allocations for specific operations
func demonstrateAllocationTracking() {
	fmt.Println("Allocation tracking:")

	// Get baseline
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// Perform allocations
	var data [][]byte
	for i := 0; i < 100; i++ {
		data = append(data, make([]byte, 1024))
	}

	// Measure allocations
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	fmt.Printf("Allocations made:\n")
	fmt.Printf("  Objects: %d\n", m2.Mallocs-m1.Mallocs)
	fmt.Printf("  Bytes: %d KB\n", (m2.TotalAlloc-m1.TotalAlloc)/1024)
	fmt.Printf("  Heap growth: %d KB\n", int64(m2.HeapAlloc-m1.HeapAlloc)/1024)

	// Keep data alive
	_ = data
}

// Memory leak detection patterns
func demonstrateMemoryLeakDetection() {
	fmt.Println("Memory leak detection patterns:")

	// Simulate potential leak scenario
	type LeakyResource struct {
		data []byte
		id   int
	}

	// Track resources
	activeResources := make(map[int]*LeakyResource)

	// Allocate resources
	for i := 0; i < 100; i++ {
		resource := &LeakyResource{
			data: make([]byte, 1024),
			id:   i,
		}
		activeResources[i] = resource
	}

	fmt.Printf("Active resources before cleanup: %d\n", len(activeResources))

	// Simulate partial cleanup (potential leak)
	for i := 0; i < 50; i++ {
		delete(activeResources, i)
	}

	fmt.Printf("Active resources after cleanup: %d\n", len(activeResources))

	// Check for potential leaks
	if len(activeResources) > 0 {
		fmt.Printf("WARNING: Potential memory leak - %d resources not cleaned up\n",
			len(activeResources))
	}

	// Proper cleanup
	for k := range activeResources {
		delete(activeResources, k)
	}

	fmt.Printf("Active resources after proper cleanup: %d\n", len(activeResources))
}

// Advanced memory management patterns
func demonstrateAdvancedPatterns() {
	fmt.Println("\n🎯 Advanced Memory Patterns")
	fmt.Println("--------------------------")

	// Memory-mapped I/O simulation
	demonstrateMemoryMapping()

	// Zero-copy patterns
	demonstrateZeroCopy()

	// Memory-efficient data structures
	demonstrateMemoryEfficientStructures()
}

// Memory mapping simulation
func demonstrateMemoryMapping() {
	fmt.Println("Memory mapping patterns:")

	// Simulate memory-mapped file behavior
	type MemoryMappedBuffer struct {
		data   []byte
		offset int
		size   int
	}

	// Create mapped buffer
	mappedSize := 1024 * 1024 // 1MB
	buffer := &MemoryMappedBuffer{
		data: make([]byte, mappedSize),
		size: mappedSize,
	}

	// Write operations (simulate file I/O)
	writeToMapped := func(buf *MemoryMappedBuffer, offset int, data []byte) bool {
		if offset+len(data) > buf.size {
			return false
		}
		copy(buf.data[offset:], data)
		return true
	}

	// Read operations
	readFromMapped := func(buf *MemoryMappedBuffer, offset, length int) []byte {
		if offset+length > buf.size {
			return nil
		}
		return buf.data[offset : offset+length]
	}

	// Simulate I/O operations
	testData := []byte("Hello, memory-mapped world!")
	writeToMapped(buffer, 0, testData)

	readData := readFromMapped(buffer, 0, len(testData))
	fmt.Printf("Memory-mapped I/O: %s\n", string(readData))

	fmt.Printf("Mapped buffer size: %d KB\n", buffer.size/1024)
}

// Zero-copy pattern demonstration
func demonstrateZeroCopy() {
	fmt.Println("Zero-copy patterns:")

	// Traditional copy approach
	traditionalCopy := func(src []byte) []byte {
		dst := make([]byte, len(src)) // Allocation + copy
		copy(dst, src)
		return dst
	}

	// Zero-copy approach using slicing
	zeroCopy := func(src []byte, start, end int) []byte {
		return src[start:end] // No allocation, shares underlying array
	}

	testData := make([]byte, 1024)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	// Compare approaches
	start := time.Now()
	for i := 0; i < 10000; i++ {
		_ = traditionalCopy(testData)
	}
	copyTime := time.Since(start)

	start = time.Now()
	for i := 0; i < 10000; i++ {
		_ = zeroCopy(testData, 0, len(testData))
	}
	zeroCopyTime := time.Since(start)

	fmt.Printf("Traditional copy: %v\n", copyTime)
	fmt.Printf("Zero-copy: %v\n", zeroCopyTime)
	fmt.Printf("Performance improvement: %.2fx\n",
		float64(copyTime)/float64(zeroCopyTime))
}

// Memory-efficient data structures
func demonstrateMemoryEfficientStructures() {
	fmt.Println("Memory-efficient data structures:")

	// Compact bit array
	type BitArray struct {
		bits []uint64
		size int
	}

	newBitArray := func(size int) *BitArray {
		return &BitArray{
			bits: make([]uint64, (size+63)/64), // Round up to 64-bit words
			size: size,
		}
	}

	setBit := func(ba *BitArray, index int) {
		if index >= ba.size {
			return
		}
		wordIndex := index / 64
		bitIndex := index % 64
		ba.bits[wordIndex] |= 1 << bitIndex
	}

	getBit := func(ba *BitArray, index int) bool {
		if index >= ba.size {
			return false
		}
		wordIndex := index / 64
		bitIndex := index % 64
		return (ba.bits[wordIndex] & (1 << bitIndex)) != 0
	}

	// Compare memory usage
	const arraySize = 10000

	// Traditional bool array
	boolArray := make([]bool, arraySize)
	for i := 0; i < arraySize; i += 3 {
		boolArray[i] = true
	}

	// Compact bit array
	bitArray := newBitArray(arraySize)
	for i := 0; i < arraySize; i += 3 {
		setBit(bitArray, i)
	}

	// Test correctness
	matches := 0
	for i := 0; i < arraySize; i++ {
		if boolArray[i] == getBit(bitArray, i) {
			matches++
		}
	}

	fmt.Printf("Bool array size: %d bytes\n", len(boolArray))
	fmt.Printf("Bit array size: %d bytes\n", len(bitArray.bits)*8)
	fmt.Printf("Memory savings: %.1fx\n",
		float64(len(boolArray))/float64(len(bitArray.bits)*8))
	fmt.Printf("Correctness: %d/%d matches\n", matches, arraySize)
}

// Helper function to show memory statistics
func showMemoryStats(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("📊 %s:\n", label)
	fmt.Printf("  Heap: %.2f MB, Objects: %d, GC: %d\n",
		float64(m.HeapAlloc)/1024/1024, m.HeapObjects, m.NumGC)
}

/*
🔥 FAANG Interview Key Points:

1. Stack vs Heap Allocation:
   - Stack: Fast, automatic cleanup, limited size
   - Heap: Flexible size, GC overhead, slower allocation
   - Escape analysis determines allocation location

2. Escape Analysis Triggers:
   - Returning pointers to local variables
   - Interface{} assignment
   - Closure variable capture
   - Channel communication
   - Large object allocation

3. Memory Alignment:
   - CPU requires aligned memory access for performance
   - Struct field ordering affects memory usage
   - Padding adds overhead but improves performance

4. Memory Optimization Patterns:
   - Object pooling reduces GC pressure
   - Arena allocation for batch processing
   - Buffer reuse minimizes allocations
   - Zero-copy avoids unnecessary data copying

5. Memory Profiling:
   - runtime.MemStats provides detailed metrics
   - Track allocation patterns over time
   - Monitor for memory leaks in long-running applications
   - Use pprof for production profiling

6. Production Implications:
   - Memory allocation patterns affect GC frequency
   - Struct design impacts cache performance
   - Pool patterns reduce allocation overhead
   - Escape analysis understanding enables optimization

This deep memory management knowledge is essential for:
- High-performance Go application development
- System optimization and troubleshooting
- Technical architecture discussions
- Senior/Staff engineer roles at FAANG companies
*/
