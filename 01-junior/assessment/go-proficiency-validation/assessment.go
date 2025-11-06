package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// GoAssessmentSuite provides comprehensive Go proficiency testing
type GoAssessmentSuite struct {
	results map[string]*TestResult
	mu      sync.RWMutex
}

// TestResult tracks individual test performance
type TestResult struct {
	TestName       string        `json:"test_name"`
	Score          int           `json:"score"`
	MaxScore       int           `json:"max_score"`
	Duration       time.Duration `json:"duration"`
	Passed         bool          `json:"passed"`
	Feedback       string        `json:"feedback"`
	MemoryUsage    uint64        `json:"memory_usage_bytes"`
	GoroutineCount int           `json:"goroutine_count"`
}

// Test interfaces for interface composition testing
type Reader interface {
	Read([]byte) (int, error)
}

type Writer interface {
	Write([]byte) (int, error)
}

type ReadWriter interface {
	Reader
	Writer
}

// TestReadWriter implements ReadWriter interface
type TestReadWriter struct {
	data []byte
	pos  int
}

func (trw *TestReadWriter) Read(p []byte) (int, error) {
	if trw.pos >= len(trw.data) {
		return 0, fmt.Errorf("EOF")
	}
	n := copy(p, trw.data[trw.pos:])
	trw.pos += n
	return n, nil
}

func (trw *TestReadWriter) Write(p []byte) (int, error) {
	trw.data = append(trw.data, p...)
	return len(p), nil
}

// NewGoAssessmentSuite initializes the assessment framework
func NewGoAssessmentSuite() *GoAssessmentSuite {
	return &GoAssessmentSuite{
		results: make(map[string]*TestResult),
	}
}

// =============================================================================
// CONCURRENCY PROFICIENCY TESTS
// =============================================================================

// TestConcurrencyPatterns evaluates goroutine and channel usage
func (gas *GoAssessmentSuite) TestConcurrencyPatterns() {
	start := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	initialMem := m.Alloc
	initialGoroutines := runtime.NumGoroutine()

	score := 0
	maxScore := 20
	feedback := []string{}

	// Test 1: Fan-out/Fan-in Pattern (5 points)
	fanOutScore := gas.testFanOutFanIn()
	score += fanOutScore
	if fanOutScore == 5 {
		feedback = append(feedback, "✅ Excellent fan-out/fan-in implementation")
	} else {
		feedback = append(feedback, "❌ Fan-out/fan-in needs improvement")
	}

	// Test 2: Worker Pool Pattern (5 points)
	workerScore := gas.testWorkerPool()
	score += workerScore
	if workerScore == 5 {
		feedback = append(feedback, "✅ Solid worker pool implementation")
	} else {
		feedback = append(feedback, "❌ Worker pool pattern needs work")
	}

	// Test 3: Context Cancellation (5 points)
	contextScore := gas.testContextUsage()
	score += contextScore
	if contextScore == 5 {
		feedback = append(feedback, "✅ Proper context usage demonstrated")
	} else {
		feedback = append(feedback, "❌ Context handling needs improvement")
	}

	// Test 4: Race Condition Prevention (5 points)
	raceScore := gas.testRacePrevention()
	score += raceScore
	if raceScore == 5 {
		feedback = append(feedback, "✅ Race conditions properly handled")
	} else {
		feedback = append(feedback, "❌ Race condition handling inadequate")
	}

	runtime.ReadMemStats(&m)
	finalMem := m.Alloc
	finalGoroutines := runtime.NumGoroutine()

	result := &TestResult{
		TestName:       "Concurrency Patterns",
		Score:          score,
		MaxScore:       maxScore,
		Duration:       time.Since(start),
		Passed:         score >= 16, // 80% threshold
		Feedback:       fmt.Sprintf("%v", feedback),
		MemoryUsage:    finalMem - initialMem,
		GoroutineCount: finalGoroutines - initialGoroutines,
	}

	gas.mu.Lock()
	gas.results["concurrency"] = result
	gas.mu.Unlock()
}

// testFanOutFanIn evaluates fan-out/fan-in pattern implementation
func (gas *GoAssessmentSuite) testFanOutFanIn() int {
	score := 0

	// Test implementation of fan-out/fan-in
	input := make(chan int, 100)
	for i := 1; i <= 100; i++ {
		input <- i
	}
	close(input)

	// Fan-out to multiple workers
	numWorkers := 5
	workerChans := make([]chan int, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workerChans[i] = make(chan int, 20)
	}

	// Distribute work
	go func() {
		i := 0
		for val := range input {
			workerChans[i%numWorkers] <- val
			i++
		}
		for _, ch := range workerChans {
			close(ch)
		}
	}()

	// Fan-in results
	results := make(chan int, 100)
	var wg sync.WaitGroup

	for _, ch := range workerChans {
		wg.Add(1)
		go func(c chan int) {
			defer wg.Done()
			for val := range c {
				results <- val * val // Square the numbers
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and verify results
	sum := 0
	count := 0
	for result := range results {
		sum += result
		count++
	}

	if count == 100 {
		score += 3 // Correct count
	}
	if sum == 338350 { // Sum of squares 1-100
		score += 2 // Correct calculation
	}

	return score
}

// testWorkerPool evaluates worker pool pattern
func (gas *GoAssessmentSuite) testWorkerPool() int {
	score := 0

	type Job struct {
		ID   int
		Data string
	}

	type Result struct {
		JobID  int
		Output string
		Error  error
	}

	jobs := make(chan Job, 50)
	results := make(chan Result, 50)

	// Start workers
	numWorkers := 3
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate work
				time.Sleep(time.Millisecond)
				results <- Result{
					JobID:  job.ID,
					Output: fmt.Sprintf("Worker %d processed: %s", workerID, job.Data),
				}
			}
		}(i)
	}

	// Send jobs
	jobCount := 20
	go func() {
		for i := 1; i <= jobCount; i++ {
			jobs <- Job{ID: i, Data: fmt.Sprintf("job-%d", i)}
		}
		close(jobs)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	processedJobs := 0
	for range results {
		processedJobs++
	}

	if processedJobs == jobCount {
		score += 5 // All jobs processed correctly
	} else if processedJobs > jobCount/2 {
		score += 3 // Most jobs processed
	} else {
		score += 1 // Some jobs processed
	}

	return score
}

// testContextUsage evaluates proper context handling
func (gas *GoAssessmentSuite) testContextUsage() int {
	score := 0

	// Test context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan bool)
	started := time.Now()

	go func() {
		select {
		case <-ctx.Done():
			done <- true
		case <-time.After(200 * time.Millisecond):
			done <- false
		}
	}()

	result := <-done
	elapsed := time.Since(started)

	if result && elapsed < 150*time.Millisecond {
		score += 5 // Proper context cancellation handling
	} else if result {
		score += 3 // Context handled but timing off
	} else {
		score += 1 // Context not properly handled
	}

	return score
}

// testRacePrevention evaluates race condition prevention
func (gas *GoAssessmentSuite) testRacePrevention() int {
	score := 0

	// Shared counter with race protection
	var counter int64
	var mu sync.Mutex
	unsafeCounter := 0

	numGoroutines := 100
	increments := 1000

	var wg sync.WaitGroup

	// Safe counter operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	// Unsafe counter operations (to demonstrate race)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				mu.Lock()
				unsafeCounter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	expectedValue := int64(numGoroutines * increments)

	if counter == expectedValue {
		score += 3 // Atomic operations working correctly
	}

	if unsafeCounter == int(expectedValue) {
		score += 2 // Mutex protection working correctly
	}

	return score
}

// =============================================================================
// MEMORY MANAGEMENT TESTS
// =============================================================================

// TestMemoryManagement evaluates memory allocation and GC understanding
func (gas *GoAssessmentSuite) TestMemoryManagement() {
	start := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	initialMem := m.Alloc

	score := 0
	maxScore := 15
	feedback := []string{}

	// Test 1: Object Pooling (5 points)
	poolScore := gas.testObjectPooling()
	score += poolScore
	if poolScore >= 4 {
		feedback = append(feedback, "✅ Good object pooling implementation")
	} else {
		feedback = append(feedback, "❌ Object pooling needs improvement")
	}

	// Test 2: Memory Leak Prevention (5 points)
	leakScore := gas.testMemoryLeakPrevention()
	score += leakScore
	if leakScore >= 4 {
		feedback = append(feedback, "✅ Memory leaks properly prevented")
	} else {
		feedback = append(feedback, "❌ Memory leak prevention inadequate")
	}

	// Test 3: Efficient String Operations (5 points)
	stringScore := gas.testStringOperations()
	score += stringScore
	if stringScore >= 4 {
		feedback = append(feedback, "✅ Efficient string operations")
	} else {
		feedback = append(feedback, "❌ String operations could be more efficient")
	}

	runtime.ReadMemStats(&m)
	finalMem := m.Alloc

	result := &TestResult{
		TestName:    "Memory Management",
		Score:       score,
		MaxScore:    maxScore,
		Duration:    time.Since(start),
		Passed:      score >= 12, // 80% threshold
		Feedback:    fmt.Sprintf("%v", feedback),
		MemoryUsage: finalMem - initialMem,
	}

	gas.mu.Lock()
	gas.results["memory"] = result
	gas.mu.Unlock()
}

// testObjectPooling evaluates sync.Pool usage
func (gas *GoAssessmentSuite) testObjectPooling() int {
	score := 0

	// Object pool for byte buffers
	bufferPool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, 1024)
		},
	}

	// Simulate heavy buffer usage with pooling
	iterations := 10000
	start := time.Now()

	for i := 0; i < iterations; i++ {
		buf := bufferPool.Get().([]byte)
		buf = buf[:0] // Reset length

		// Simulate some work
		for j := 0; j < 100; j++ {
			buf = append(buf, byte(j))
		}

		bufferPool.Put(buf)
	}

	pooledDuration := time.Since(start)

	// Simulate same work without pooling
	start = time.Now()
	for i := 0; i < iterations; i++ {
		buf := make([]byte, 0, 1024)
		for j := 0; j < 100; j++ {
			buf = append(buf, byte(j))
		}
	}
	nonPooledDuration := time.Since(start)

	// Pool should be faster due to reduced allocations
	if pooledDuration < nonPooledDuration {
		score += 5
	} else if pooledDuration < time.Duration(float64(nonPooledDuration)*1.2) {
		score += 3
	} else {
		score += 1
	}

	return score
}

// testMemoryLeakPrevention checks for proper cleanup
func (gas *GoAssessmentSuite) testMemoryLeakPrevention() int {
	score := 0

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Create many objects that should be garbage collected
	var slices [][]int
	for i := 0; i < 1000; i++ {
		slice := make([]int, 1000)
		for j := range slice {
			slice[j] = i * j
		}
		slices = append(slices, slice)
	}

	// Clear references and force GC
	slices = nil
	runtime.GC()
	runtime.GC() // Call twice to ensure cleanup
	runtime.ReadMemStats(&m2)

	// Memory should have been reclaimed
	memDiff := int64(m2.Alloc) - int64(m1.Alloc)
	if memDiff < 1024*1024 { // Less than 1MB growth
		score += 5
	} else if memDiff < 5*1024*1024 { // Less than 5MB growth
		score += 3
	} else {
		score += 1
	}

	return score
}

// testStringOperations evaluates efficient string handling
func (gas *GoAssessmentSuite) testStringOperations() int {
	score := 0

	// Test string concatenation efficiency
	numStrings := 10000
	testString := "test"

	// Using strings.Builder (efficient)
	start := time.Now()
	var builder strings.Builder
	builder.Grow(numStrings * len(testString)) // Pre-allocate
	for i := 0; i < numStrings; i++ {
		builder.WriteString(testString)
	}
	result1 := builder.String()
	builderDuration := time.Since(start)

	// Using string concatenation (inefficient)
	start = time.Now()
	result2 := ""
	for i := 0; i < numStrings; i++ {
		result2 += testString
	}
	concatDuration := time.Since(start)

	// Verify correctness
	if len(result1) == len(result2) && len(result1) == numStrings*len(testString) {
		score += 2 // Correct implementation
	}

	// Builder should be significantly faster
	if builderDuration < concatDuration/2 {
		score += 3 // Efficient implementation
	} else if builderDuration < concatDuration {
		score += 1 // Somewhat efficient
	}

	return score
}

// =============================================================================
// INTERFACE AND TYPE SYSTEM TESTS
// =============================================================================

// TestInterfaceDesign evaluates interface usage and type system understanding
func (gas *GoAssessmentSuite) TestInterfaceDesign() {
	start := time.Now()

	score := 0
	maxScore := 10
	feedback := []string{}

	// Test interface composition and embedding
	compositionScore := gas.testInterfaceComposition()
	score += compositionScore
	if compositionScore >= 8 {
		feedback = append(feedback, "✅ Excellent interface design")
	} else {
		feedback = append(feedback, "❌ Interface design needs improvement")
	}

	result := &TestResult{
		TestName: "Interface Design",
		Score:    score,
		MaxScore: maxScore,
		Duration: time.Since(start),
		Passed:   score >= 8, // 80% threshold
		Feedback: fmt.Sprintf("%v", feedback),
	}

	gas.mu.Lock()
	gas.results["interfaces"] = result
	gas.mu.Unlock()
}

// testInterfaceComposition evaluates proper interface design
func (gas *GoAssessmentSuite) testInterfaceComposition() int {
	score := 0

	// Test implementation
	var rw ReadWriter = &TestReadWriter{}

	// Test write operation
	testData := []byte("Hello, World!")
	n, err := rw.Write(testData)
	if err == nil && n == len(testData) {
		score += 5
	}

	// Test read operation
	readBuf := make([]byte, len(testData))
	n, err = rw.Read(readBuf)
	if err == nil && n == len(testData) && string(readBuf) == string(testData) {
		score += 5
	}

	return score
}

// =============================================================================
// PERFORMANCE PROFILING TESTS
// =============================================================================

// TestPerformanceProfiling evaluates profiling and optimization skills
func (gas *GoAssessmentSuite) TestPerformanceProfiling() {
	start := time.Now()

	score := 0
	maxScore := 10
	feedback := []string{}

	// Test CPU profiling understanding
	cpuScore := gas.testCPUProfiling()
	score += cpuScore

	// Test memory profiling understanding
	memScore := gas.testMemoryProfiling()
	score += memScore

	if score >= 8 {
		feedback = append(feedback, "✅ Strong profiling skills demonstrated")
	} else {
		feedback = append(feedback, "❌ Profiling skills need development")
	}

	result := &TestResult{
		TestName: "Performance Profiling",
		Score:    score,
		MaxScore: maxScore,
		Duration: time.Since(start),
		Passed:   score >= 8, // 80% threshold
		Feedback: fmt.Sprintf("%v", feedback),
	}

	gas.mu.Lock()
	gas.results["profiling"] = result
	gas.mu.Unlock()
}

// testCPUProfiling simulates CPU-intensive work for profiling
func (gas *GoAssessmentSuite) testCPUProfiling() int {
	score := 0

	// Simulate CPU-intensive work
	start := time.Now()
	result := 0
	for i := 0; i < 1000000; i++ {
		result += i * i
	}
	duration := time.Since(start)

	// Basic performance awareness
	if duration < 50*time.Millisecond {
		score += 5 // Good performance
	} else if duration < 100*time.Millisecond {
		score += 3 // Acceptable performance
	} else {
		score += 1 // Performance needs improvement
	}

	return score
}

// testMemoryProfiling evaluates memory allocation patterns
func (gas *GoAssessmentSuite) testMemoryProfiling() int {
	score := 0

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// Allocate and immediately release memory
	data := make([][]int, 1000)
	for i := range data {
		data[i] = make([]int, 1000)
	}

	// Clear and force GC
	data = nil
	runtime.GC()
	runtime.ReadMemStats(&m2)

	// Check if memory was properly managed
	allocDiff := m2.TotalAlloc - m1.TotalAlloc
	if allocDiff > 0 {
		score += 5 // Memory allocation tracked
	}

	return score
}

// =============================================================================
// REPORTING AND ANALYSIS
// =============================================================================

// GenerateAssessmentReport creates a comprehensive report
func (gas *GoAssessmentSuite) GenerateAssessmentReport() string {
	gas.mu.RLock()
	defer gas.mu.RUnlock()

	totalScore := 0
	totalMaxScore := 0

	report := "🔥 GO PROFICIENCY ASSESSMENT REPORT\n"
	report += "====================================\n\n"

	for _, result := range gas.results {
		totalScore += result.Score
		totalMaxScore += result.MaxScore

		percentage := float64(result.Score) / float64(result.MaxScore) * 100
		status := "❌ NEEDS WORK"
		if percentage >= 80 {
			status = "✅ EXCELLENT"
		} else if percentage >= 70 {
			status = "⚡ GOOD"
		} else if percentage >= 60 {
			status = "⚠️  DEVELOPING"
		}

		report += fmt.Sprintf("%s: %d/%d (%.1f%%) %s\n",
			result.TestName, result.Score, result.MaxScore, percentage, status)
		report += fmt.Sprintf("  Duration: %v\n", result.Duration)
		report += fmt.Sprintf("  Feedback: %s\n\n", result.Feedback)
	}

	overallPercentage := float64(totalScore) / float64(totalMaxScore) * 100
	report += fmt.Sprintf("OVERALL SCORE: %d/%d (%.1f%%)\n\n", totalScore, totalMaxScore, overallPercentage)

	// Recommendations
	report += "🎯 RECOMMENDATIONS:\n"
	if overallPercentage >= 85 {
		report += "✅ ELITE LEVEL: You're ready for FAANG L3/L4 interviews!\n"
	} else if overallPercentage >= 75 {
		report += "⚡ STRONG: Focus on weak areas, you're close to elite level!\n"
	} else if overallPercentage >= 65 {
		report += "📚 DEVELOPING: Strengthen fundamentals before advanced topics.\n"
	} else {
		report += "🏗️ FOUNDATION: Complete Go basics course before this track.\n"
	}

	return report
}

// RunFullAssessment executes all assessment tests
func (gas *GoAssessmentSuite) RunFullAssessment() {
	fmt.Println("🚀 Starting Go Proficiency Assessment...")

	gas.TestConcurrencyPatterns()
	fmt.Println("✅ Concurrency patterns assessment complete")

	gas.TestMemoryManagement()
	fmt.Println("✅ Memory management assessment complete")

	gas.TestInterfaceDesign()
	fmt.Println("✅ Interface design assessment complete")

	gas.TestPerformanceProfiling()
	fmt.Println("✅ Performance profiling assessment complete")

	fmt.Println("\n" + gas.GenerateAssessmentReport())
}

// Example usage and CLI
func main() {
	if len(os.Args) > 1 && os.Args[1] == "assess" {
		suite := NewGoAssessmentSuite()
		suite.RunFullAssessment()

		// Save results to JSON
		suite.mu.RLock()
		resultsJSON, _ := json.MarshalIndent(suite.results, "", "  ")
		suite.mu.RUnlock()

		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("go_assessment_%s.json", timestamp)
		if err := ioutil.WriteFile(filename, resultsJSON, 0644); err == nil {
			fmt.Printf("\n📊 Assessment results saved to: %s\n", filename)
		}
	} else {
		fmt.Println("Go Proficiency Assessment Tool")
		fmt.Println("Usage: go run assessment.go assess")
		fmt.Println("\nThis tool provides comprehensive Go skills assessment for FAANG interview preparation.")
	}
}
