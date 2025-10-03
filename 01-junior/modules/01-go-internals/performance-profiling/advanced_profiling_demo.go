// Package profiling demonstrates Go's performance profiling and optimization
// Production-level profiling knowledge for FAANG performance engineering roles

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof" // Enable pprof endpoints
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"sync"
	"time"
)

// Types for benchmarking
type Processor interface {
	Process(int) int
}

type ConcreteProcessor struct{}

func (cp ConcreteProcessor) Process(x int) int {
	return x * 2
}

// Performance Profiling Core Areas:
// - CPU profiling with pprof
// - Memory profiling and leak detection
// - Goroutine profiling and analysis
// - Execution tracing
// - Mutex contention profiling
// - HTTP profiling endpoints
// - Production monitoring patterns

func main() {
	fmt.Println("🔥 Go Performance Profiling & Optimization")
	fmt.Println("==========================================")

	// Start HTTP server for pprof endpoints in background
	go startPprofServer()

	// CPU profiling demonstration
	demonstrateCPUProfiling()

	// Memory profiling demonstration
	demonstrateMemoryProfiling()

	// Goroutine profiling demonstration
	demonstrateGoroutineProfiling()

	// Execution tracing demonstration
	demonstrateExecutionTracing()

	// Mutex profiling demonstration
	demonstrateMutexProfiling()

	// Production monitoring patterns
	demonstrateProductionMonitoring()

	// Performance optimization techniques
	demonstrateOptimizationTechniques()

	fmt.Println("\n🌐 pprof web interface available at:")
	fmt.Println("  http://localhost:6060/debug/pprof/")
	fmt.Println("  Run: go tool pprof http://localhost:6060/debug/pprof/profile")
}

// Start HTTP server with pprof endpoints
func startPprofServer() {
	fmt.Println("Starting pprof HTTP server on :6060")
	log.Println(http.ListenAndServe("localhost:6060", nil))
}

// CPU profiling techniques and analysis
func demonstrateCPUProfiling() {
	fmt.Println("\n🚀 CPU Profiling Demonstration")
	fmt.Println("------------------------------")

	// Create CPU profile file
	cpuFile, err := os.Create("cpu_profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer cpuFile.Close()

	// Start CPU profiling
	fmt.Println("Starting CPU profiling...")
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Fatal(err)
	}

	// Run CPU-intensive workload
	runCPUIntensiveWorkload()

	// Stop CPU profiling
	pprof.StopCPUProfile()
	fmt.Println("CPU profile saved to cpu_profile.prof")
	fmt.Println("Analyze with: go tool pprof cpu_profile.prof")
}

// CPU-intensive workload for profiling
func runCPUIntensiveWorkload() {
	fmt.Println("Running CPU-intensive workload...")

	// Mathematical computation
	result := computeExpensiveOperation(1000000)
	fmt.Printf("Expensive operation result: %d\n", result)

	// String processing
	processStrings(50000)

	// Slice operations
	processSlices(100000)
}

// Expensive mathematical computation
func computeExpensiveOperation(n int) int64 {
	var result int64
	for i := 0; i < n; i++ {
		// Simulate complex calculation
		for j := 0; j < 100; j++ {
			result += int64(i * j)
			result = result ^ int64(i+j) // XOR operation
		}
	}
	return result
}

// String processing workload
func processStrings(count int) {
	var result []string
	for i := 0; i < count; i++ {
		// String concatenation and manipulation
		str := fmt.Sprintf("processing_string_%d", i)
		str += "_suffix"
		result = append(result, str)
	}
	fmt.Printf("Processed %d strings\n", len(result))
}

// Slice operations workload
func processSlices(count int) {
	data := make([]int, 0, count)
	for i := 0; i < count; i++ {
		data = append(data, i*i)
	}

	// Sort and search operations
	quickSort(data, 0, len(data)-1)
	_ = binarySearch(data, count/2)
	fmt.Printf("Processed slice with %d elements\n", len(data))
}

// Quick sort implementation for profiling
func quickSort(arr []int, low, high int) {
	if low < high {
		pi := partition(arr, low, high)
		quickSort(arr, low, pi-1)
		quickSort(arr, pi+1, high)
	}
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// Binary search implementation
func binarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := left + (right-left)/2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

// Memory profiling and leak detection
func demonstrateMemoryProfiling() {
	fmt.Println("\n💾 Memory Profiling Demonstration")
	fmt.Println("--------------------------------")

	// Run memory-intensive workload
	runMemoryIntensiveWorkload()

	// Create heap profile
	heapFile, err := os.Create("heap_profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer heapFile.Close()

	// Force GC and write heap profile
	runtime.GC()
	if err := pprof.WriteHeapProfile(heapFile); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Heap profile saved to heap_profile.prof")
	fmt.Println("Analyze with: go tool pprof heap_profile.prof")

	// Memory statistics analysis
	analyzeMemoryStats()
}

// Memory-intensive workload
func runMemoryIntensiveWorkload() {
	fmt.Println("Running memory-intensive workload...")

	// Large slice allocations
	var largeSlices [][]byte
	for i := 0; i < 100; i++ {
		slice := make([]byte, 1024*1024) // 1MB each
		for j := range slice {
			slice[j] = byte(j % 256)
		}
		largeSlices = append(largeSlices, slice)
	}

	// Map with many entries
	largeMap := make(map[string][]int)
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := make([]int, 100)
		for j := range value {
			value[j] = i + j
		}
		largeMap[key] = value
	}

	// Struct allocations
	type ComplexStruct struct {
		ID       int
		Name     string
		Data     []byte
		Children []*ComplexStruct
	}

	var structures []*ComplexStruct
	for i := 0; i < 1000; i++ {
		s := &ComplexStruct{
			ID:   i,
			Name: fmt.Sprintf("struct_%d", i),
			Data: make([]byte, 512),
		}
		// Create some child structures
		for j := 0; j < 5; j++ {
			child := &ComplexStruct{
				ID:   i*10 + j,
				Name: fmt.Sprintf("child_%d_%d", i, j),
				Data: make([]byte, 256),
			}
			s.Children = append(s.Children, child)
		}
		structures = append(structures, s)
	}

	fmt.Printf("Allocated: %d large slices, %d map entries, %d structures\n",
		len(largeSlices), len(largeMap), len(structures))

	// Keep references to prevent immediate GC
	_ = largeSlices
	_ = largeMap
	_ = structures
}

// Analyze detailed memory statistics
func analyzeMemoryStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Println("Detailed Memory Analysis:")
	fmt.Printf("  Heap Allocated: %.2f MB\n", float64(m.HeapAlloc)/1024/1024)
	fmt.Printf("  Heap System: %.2f MB\n", float64(m.HeapSys)/1024/1024)
	fmt.Printf("  Heap Idle: %.2f MB\n", float64(m.HeapIdle)/1024/1024)
	fmt.Printf("  Heap In Use: %.2f MB\n", float64(m.HeapInuse)/1024/1024)
	fmt.Printf("  Heap Objects: %d\n", m.HeapObjects)
	fmt.Printf("  Stack In Use: %.2f MB\n", float64(m.StackInuse)/1024/1024)
	fmt.Printf("  Stack System: %.2f MB\n", float64(m.StackSys)/1024/1024)
	fmt.Printf("  Total Allocated: %.2f MB\n", float64(m.TotalAlloc)/1024/1024)
	fmt.Printf("  Mallocs: %d\n", m.Mallocs)
	fmt.Printf("  Frees: %d\n", m.Frees)
	fmt.Printf("  GC Cycles: %d\n", m.NumGC)
	fmt.Printf("  GC CPU Fraction: %.4f\n", m.GCCPUFraction)
}

// Goroutine profiling and analysis
func demonstrateGoroutineProfiling() {
	fmt.Println("\n🔀 Goroutine Profiling Demonstration")
	fmt.Println("-----------------------------------")

	// Create various goroutine patterns
	createGoroutineWorkload()

	// Create goroutine profile
	goroutineFile, err := os.Create("goroutine_profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer goroutineFile.Close()

	// Write goroutine profile
	if err := pprof.Lookup("goroutine").WriteTo(goroutineFile, 0); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Goroutine profile saved to goroutine_profile.prof")
	fmt.Println("Analyze with: go tool pprof goroutine_profile.prof")

	// Analyze goroutine states
	analyzeGoroutineStates()
}

// Create diverse goroutine workload
func createGoroutineWorkload() {
	fmt.Println("Creating goroutine workload...")

	var wg sync.WaitGroup

	// CPU-intensive goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// CPU-intensive work
			result := 0
			for j := 0; j < 1000000; j++ {
				result += j * id
			}
			_ = result
		}(i)
	}

	// I/O blocking goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Simulate I/O blocking
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}(i)
	}

	// Channel communication goroutines
	ch := make(chan int, 5)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				ch <- id*10 + j
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Channel consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for count < 30 {
			<-ch
			count++
		}
	}()

	wg.Wait()
	fmt.Printf("Goroutine workload completed. Current goroutines: %d\n", runtime.NumGoroutine())
}

// Analyze goroutine states
func analyzeGoroutineStates() {
	// Get stack traces of all goroutines
	buf := make([]byte, 1<<20) // 1MB buffer
	stackSize := runtime.Stack(buf, true)

	fmt.Printf("Stack trace size: %d bytes\n", stackSize)
	fmt.Printf("Active goroutines: %d\n", runtime.NumGoroutine())

	// Write stack trace to file for analysis
	stackFile, err := os.Create("goroutine_stacks.txt")
	if err != nil {
		log.Printf("Error creating stack file: %v", err)
		return
	}
	defer stackFile.Close()

	stackFile.Write(buf[:stackSize])
	fmt.Println("Goroutine stack traces saved to goroutine_stacks.txt")
}

// Execution tracing demonstration
func demonstrateExecutionTracing() {
	fmt.Println("\n📊 Execution Tracing Demonstration")
	fmt.Println("---------------------------------")

	// Create trace file
	traceFile, err := os.Create("execution_trace.trace")
	if err != nil {
		log.Fatal(err)
	}
	defer traceFile.Close()

	// Start execution tracing
	fmt.Println("Starting execution tracing...")
	if err := trace.Start(traceFile); err != nil {
		log.Fatal(err)
	}

	// Run workload with tracing
	runTracedWorkload()

	// Stop tracing
	trace.Stop()
	fmt.Println("Execution trace saved to execution_trace.trace")
	fmt.Println("Analyze with: go tool trace execution_trace.trace")
}

// Workload designed for trace analysis
func runTracedWorkload() {
	fmt.Println("Running traced workload...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create multiple concurrent workers
	var wg sync.WaitGroup
	numWorkers := 4

	// Work queue
	workCh := make(chan int, 100)
	resultCh := make(chan int, 100)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			tracedWorker(ctx, workerID, workCh, resultCh)
		}(i)
	}

	// Send work
	go func() {
		for i := 0; i < 100; i++ {
			select {
			case workCh <- i:
			case <-ctx.Done():
				return
			}
		}
		close(workCh)
	}()

	// Collect results
	go func() {
		resultsCollected := 0
		for range resultCh {
			resultsCollected++
			if resultsCollected >= 100 {
				break
			}
		}
	}()

	wg.Wait()
	close(resultCh)
	fmt.Println("Traced workload completed")
}

// Worker function for tracing
func tracedWorker(ctx context.Context, workerID int, workCh <-chan int, resultCh chan<- int) {
	for {
		select {
		case work, ok := <-workCh:
			if !ok {
				return
			}
			// Simulate work with trace regions
			taskCtx, task := trace.NewTask(ctx, "process_work")
			_ = taskCtx // Use task context if needed
			result := processWork(work)
			resultCh <- result
			task.End()

		case <-ctx.Done():
			return
		}
	}
}

// Process work with trace regions
func processWork(work int) int {
	// Simulate different types of work
	switch work % 3 {
	case 0:
		// CPU-intensive work
		trace.WithRegion(context.Background(), "cpu_work", func() {
			result := 0
			for i := 0; i < 10000; i++ {
				result += i * work
			}
		})
	case 1:
		// Memory allocation work
		trace.WithRegion(context.Background(), "memory_work", func() {
			data := make([]int, 1000)
			for i := range data {
				data[i] = i + work
			}
		})
	case 2:
		// Sleep work (I/O simulation)
		trace.WithRegion(context.Background(), "io_work", func() {
			time.Sleep(time.Millisecond)
		})
	}
	return work * 2
}

// Mutex contention profiling
func demonstrateMutexProfiling() {
	fmt.Println("\n🔒 Mutex Contention Profiling")
	fmt.Println("-----------------------------")

	// Enable mutex profiling
	runtime.SetMutexProfileFraction(1)

	// Run workload with mutex contention
	runMutexContentionWorkload()

	// Create mutex profile
	mutexFile, err := os.Create("mutex_profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer mutexFile.Close()

	// Write mutex profile
	if err := pprof.Lookup("mutex").WriteTo(mutexFile, 0); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mutex profile saved to mutex_profile.prof")
	fmt.Println("Analyze with: go tool pprof mutex_profile.prof")
}

// Workload with intentional mutex contention
func runMutexContentionWorkload() {
	fmt.Println("Running mutex contention workload...")

	// Shared resource with contention
	type Counter struct {
		mu    sync.Mutex
		value int
	}

	counter := &Counter{}
	var wg sync.WaitGroup

	// High contention scenario
	numGoroutines := 20
	operationsPerGoroutine := 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				// Intentionally hold lock for longer to create contention
				counter.mu.Lock()
				counter.value++
				// Simulate some work while holding lock
				time.Sleep(time.Microsecond * 10)
				counter.mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("Mutex contention workload completed. Final counter value: %d\n", counter.value)
}

// Production monitoring patterns
func demonstrateProductionMonitoring() {
	fmt.Println("\n📈 Production Monitoring Patterns")
	fmt.Println("--------------------------------")

	// Continuous monitoring goroutine
	go continuousMonitoring()

	// Performance metrics collection
	collectPerformanceMetrics()

	// Health check endpoint simulation
	simulateHealthChecks()
}

// Continuous monitoring in production
func continuousMonitoring() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 3; i++ { // Run for demo purposes
		select {
		case <-ticker.C:
			monitorSystemHealth()
		}
	}
}

// Monitor system health metrics
func monitorSystemHealth() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("📊 System Health Check:\n")
	fmt.Printf("  Goroutines: %d\n", runtime.NumGoroutine())
	fmt.Printf("  Heap: %.2f MB\n", float64(m.HeapAlloc)/1024/1024)
	fmt.Printf("  GC Cycles: %d\n", m.NumGC)
	fmt.Printf("  GC Pause: %v\n", time.Duration(m.PauseTotalNs))

	// Alert conditions
	if runtime.NumGoroutine() > 1000 {
		fmt.Printf("⚠️  WARNING: High goroutine count: %d\n", runtime.NumGoroutine())
	}

	if float64(m.HeapAlloc)/1024/1024 > 500 {
		fmt.Printf("⚠️  WARNING: High memory usage: %.2f MB\n", float64(m.HeapAlloc)/1024/1024)
	}
}

// Collect performance metrics
func collectPerformanceMetrics() {
	fmt.Println("Collecting performance metrics...")

	// Simulate metric collection
	metrics := map[string]interface{}{
		"heap_size_mb":      float64(getHeapSize()) / 1024 / 1024,
		"goroutine_count":   runtime.NumGoroutine(),
		"gc_cycles":         getGCCycles(),
		"cpu_usage_percent": getCPUUsage(),
		"memory_usage_mb":   getMemoryUsage() / 1024 / 1024,
	}

	fmt.Println("Performance Metrics:")
	for key, value := range metrics {
		fmt.Printf("  %s: %v\n", key, value)
	}
}

// Health check simulation
func simulateHealthChecks() {
	fmt.Println("Simulating health check endpoints...")

	// Simulate different health check types
	checks := []struct {
		name     string
		status   string
		duration time.Duration
	}{
		{"database", "healthy", 15 * time.Millisecond},
		{"cache", "healthy", 5 * time.Millisecond},
		{"external_api", "healthy", 25 * time.Millisecond},
		{"disk_space", "healthy", 2 * time.Millisecond},
	}

	for _, check := range checks {
		start := time.Now()
		// Simulate health check
		time.Sleep(check.duration)
		duration := time.Since(start)

		fmt.Printf("  Health Check - %s: %s (took %v)\n",
			check.name, check.status, duration)
	}
}

// Performance optimization techniques
func demonstrateOptimizationTechniques() {
	fmt.Println("\n⚡ Performance Optimization Techniques")
	fmt.Println("------------------------------------")

	// Benchmark different approaches
	benchmarkStringConcatenation()
	benchmarkSliceOperations()
	benchmarkMapOperations()
	benchmarkInterfaceVsDirectCall()
}

// Benchmark string concatenation methods
func benchmarkStringConcatenation() {
	fmt.Println("String concatenation benchmarks:")

	const iterations = 10000

	// Method 1: String concatenation with +
	start := time.Now()
	result1 := ""
	for i := 0; i < iterations; i++ {
		result1 += "test"
	}
	time1 := time.Since(start)

	// Method 2: Using fmt.Sprintf
	start = time.Now()
	result2 := ""
	for i := 0; i < iterations; i++ {
		result2 = fmt.Sprintf("%s%s", result2, "test")
	}
	time2 := time.Since(start)

	// Method 3: Using strings.Builder
	start = time.Now()
	var builder strings.Builder
	for i := 0; i < iterations; i++ {
		builder.WriteString("test")
	}
	result3 := builder.String()
	time3 := time.Since(start)

	fmt.Printf("  String + operator: %v (length: %d)\n", time1, len(result1))
	fmt.Printf("  fmt.Sprintf: %v (length: %d)\n", time2, len(result2))
	fmt.Printf("  strings.Builder: %v (length: %d)\n", time3, len(result3))
}

// Benchmark slice operations
func benchmarkSliceOperations() {
	fmt.Println("Slice operations benchmarks:")

	const size = 100000

	// Pre-allocated slice
	start := time.Now()
	prealloc := make([]int, 0, size)
	for i := 0; i < size; i++ {
		prealloc = append(prealloc, i)
	}
	time1 := time.Since(start)

	// Growing slice
	start = time.Now()
	growing := make([]int, 0)
	for i := 0; i < size; i++ {
		growing = append(growing, i)
	}
	time2 := time.Since(start)

	fmt.Printf("  Pre-allocated slice: %v\n", time1)
	fmt.Printf("  Growing slice: %v\n", time2)
	fmt.Printf("  Performance improvement: %.2fx\n", float64(time2)/float64(time1))
}

// Benchmark map operations
func benchmarkMapOperations() {
	fmt.Println("Map operations benchmarks:")

	const size = 50000

	// Pre-sized map
	start := time.Now()
	presized := make(map[int]string, size)
	for i := 0; i < size; i++ {
		presized[i] = fmt.Sprintf("value_%d", i)
	}
	time1 := time.Since(start)

	// Regular map
	start = time.Now()
	regular := make(map[int]string)
	for i := 0; i < size; i++ {
		regular[i] = fmt.Sprintf("value_%d", i)
	}
	time2 := time.Since(start)

	fmt.Printf("  Pre-sized map: %v\n", time1)
	fmt.Printf("  Regular map: %v\n", time2)
	fmt.Printf("  Performance improvement: %.2fx\n", float64(time2)/float64(time1))
}

// Benchmark interface vs direct calls
func benchmarkInterfaceVsDirectCall() {
	fmt.Println("Interface vs direct call benchmarks:")

	concrete := ConcreteProcessor{}
	var iface Processor = concrete

	const iterations = 1000000

	// Direct call
	start := time.Now()
	var result1 int
	for i := 0; i < iterations; i++ {
		result1 += concrete.Process(i)
	}
	time1 := time.Since(start)

	// Interface call
	start = time.Now()
	var result2 int
	for i := 0; i < iterations; i++ {
		result2 += iface.Process(i)
	}
	time2 := time.Since(start)

	fmt.Printf("  Direct call: %v (result: %d)\n", time1, result1)
	fmt.Printf("  Interface call: %v (result: %d)\n", time2, result2)
	fmt.Printf("  Overhead: %.2fx\n", float64(time2)/float64(time1))
}

// Helper functions for metrics
func getHeapSize() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

func getGCCycles() uint32 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.NumGC
}

func getCPUUsage() float64 {
	// Simplified CPU usage calculation
	return float64(runtime.NumGoroutine()) * 0.1
}

func getMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Sys
}

/*
🔥 FAANG Interview Key Points:

1. CPU Profiling:
   - Use pprof.StartCPUProfile() for CPU-intensive analysis
   - Identify hot paths and optimization opportunities
   - Profile before and after optimization changes
   - Focus on functions with highest CPU time

2. Memory Profiling:
   - Heap profiles show allocation patterns
   - Track memory leaks and excessive allocations
   - Monitor GC frequency and pause times
   - Use runtime.MemStats for real-time monitoring

3. Goroutine Profiling:
   - Detect goroutine leaks and blocking issues
   - Analyze goroutine states and stack traces
   - Monitor goroutine growth patterns
   - Identify synchronization bottlenecks

4. Execution Tracing:
   - Visualize goroutine scheduling and execution
   - Identify GC impact on application performance
   - Analyze processor utilization patterns
   - Debug complex concurrency issues

5. Production Monitoring:
   - Continuous health monitoring
   - Performance metrics collection
   - Alerting on threshold breaches
   - HTTP endpoints for runtime profiling

6. Optimization Techniques:
   - Benchmark different implementation approaches
   - Pre-allocate slices and maps when size is known
   - Use strings.Builder for string concatenation
   - Minimize interface overhead in hot paths
   - Profile-guided optimization decisions

This profiling expertise is essential for:
- Performance engineering roles at FAANG companies
- Production system optimization and troubleshooting
- Capacity planning and resource optimization
- Technical leadership in performance-critical systems
*/
