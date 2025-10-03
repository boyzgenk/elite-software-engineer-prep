// Package goroutine_leak_detection demonstrates advanced techniques for detecting goroutine leaks
// This is critical for production Go applications and frequently tested in FAANG interviews
package main

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ===== GOROUTINE LEAK DETECTOR =====

// GoroutineLeakDetector monitors and detects goroutine leaks
type GoroutineLeakDetector struct {
	initialCount   int
	threshold      int
	checkInterval  time.Duration
	timeout        time.Duration
	ignorePatterns []string
	mu             sync.RWMutex
	stackDumps     []string
	leaksDetected  int64
}

// NewGoroutineLeakDetector creates a new leak detector
func NewGoroutineLeakDetector() *GoroutineLeakDetector {
	return &GoroutineLeakDetector{
		initialCount:  runtime.NumGoroutine(),
		threshold:     10, // Alert if more than 10 extra goroutines
		checkInterval: 100 * time.Millisecond,
		timeout:       30 * time.Second,
		ignorePatterns: []string{
			"testing.(*T).Run",   // Test framework goroutines
			"testing.tRunner",    // Test runner goroutines
			"runtime.goexit",     // Runtime goroutines
			"go.uber.org/goleak", // Goleak test goroutines
		},
	}
}

// GoroutineSnapshot represents a snapshot of running goroutines
type GoroutineSnapshot struct {
	Count     int
	Timestamp time.Time
	Stacks    []GoroutineStack
}

// GoroutineStack represents a single goroutine's stack trace
type GoroutineStack struct {
	ID       int
	State    string
	Function string
	File     string
	Line     int
	Stack    string
}

// TakeSnapshot captures current goroutine state
func (gld *GoroutineLeakDetector) TakeSnapshot() *GoroutineSnapshot {
	gld.mu.Lock()
	defer gld.mu.Unlock()

	count := runtime.NumGoroutine()

	// Get stack traces for all goroutines
	buf := make([]byte, 1<<20) // 1MB buffer for stack traces
	stackSize := runtime.Stack(buf, true)
	stackData := string(buf[:stackSize])

	stacks := gld.parseStackTrace(stackData)

	return &GoroutineSnapshot{
		Count:     count,
		Timestamp: time.Now(),
		Stacks:    stacks,
	}
}

// parseStackTrace parses the runtime.Stack output into structured data
func (gld *GoroutineLeakDetector) parseStackTrace(stackData string) []GoroutineStack {
	lines := strings.Split(stackData, "\n")
	var stacks []GoroutineStack
	var currentStack *GoroutineStack

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// New goroutine header: "goroutine 123 [running]:"
		if strings.HasPrefix(line, "goroutine ") && strings.HasSuffix(line, ":") {
			if currentStack != nil {
				stacks = append(stacks, *currentStack)
			}

			currentStack = &GoroutineStack{}
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				currentStack.State = strings.Trim(parts[2], "[]:")
			}
			continue
		}

		// Function call: "main.main()"
		if currentStack != nil && !strings.Contains(line, "/") && !strings.Contains(line, "\\") {
			if currentStack.Function == "" {
				currentStack.Function = line
			}
			currentStack.Stack += line + "\n"
			continue
		}

		// File location: "/path/to/file.go:123 +0x456"
		if currentStack != nil && (strings.Contains(line, ".go:") || strings.Contains(line, "+0x")) {
			if currentStack.File == "" && strings.Contains(line, ".go:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					currentStack.File = parts[0]
				}
			}
			currentStack.Stack += line + "\n"
			continue
		}

		// Other stack trace lines
		if currentStack != nil {
			currentStack.Stack += line + "\n"
		}
	}

	// Add the last stack if any
	if currentStack != nil {
		stacks = append(stacks, *currentStack)
	}

	return stacks
}

// DetectLeaks monitors for goroutine leaks during test execution
func (gld *GoroutineLeakDetector) DetectLeaks(ctx context.Context, testFunc func() error) (*LeakReport, error) {
	report := &LeakReport{
		StartTime:       time.Now(),
		InitialCount:    gld.initialCount,
		InitialSnapshot: gld.TakeSnapshot(),
	}

	// Run the test function
	testErr := testFunc()

	// Wait for potential cleanup
	time.Sleep(100 * time.Millisecond)

	// Take final snapshot
	report.FinalSnapshot = gld.TakeSnapshot()
	report.FinalCount = report.FinalSnapshot.Count
	report.LeakCount = report.FinalCount - report.InitialCount
	report.Duration = time.Since(report.StartTime)

	// Check for leaks
	if report.LeakCount > gld.threshold {
		atomic.AddInt64(&gld.leaksDetected, 1)
		report.LeakDetected = true
		report.LeakedGoroutines = gld.identifyLeakedGoroutines(report.InitialSnapshot, report.FinalSnapshot)
	}

	return report, testErr
}

// identifyLeakedGoroutines compares snapshots to identify leaked goroutines
func (gld *GoroutineLeakDetector) identifyLeakedGoroutines(initial, final *GoroutineSnapshot) []GoroutineStack {
	var leaked []GoroutineStack

	// Create a map of initial goroutines by their function signature
	initialFuncs := make(map[string]int)
	for _, stack := range initial.Stacks {
		key := stack.Function + ":" + stack.State
		initialFuncs[key]++
	}

	// Find final goroutines that weren't in initial state
	finalFuncs := make(map[string]int)
	for _, stack := range final.Stacks {
		key := stack.Function + ":" + stack.State
		finalFuncs[key]++
	}

	// Identify potentially leaked goroutines
	for key, finalCount := range finalFuncs {
		initialCount := initialFuncs[key]
		if finalCount > initialCount {
			// Find a representative stack for this leak
			for _, stack := range final.Stacks {
				stackKey := stack.Function + ":" + stack.State
				if stackKey == key && !gld.shouldIgnoreStack(stack) {
					leaked = append(leaked, stack)
					break // Only add one representative per type
				}
			}
		}
	}

	return leaked
}

// shouldIgnoreStack checks if a stack should be ignored based on patterns
func (gld *GoroutineLeakDetector) shouldIgnoreStack(stack GoroutineStack) bool {
	for _, pattern := range gld.ignorePatterns {
		if strings.Contains(stack.Stack, pattern) || strings.Contains(stack.Function, pattern) {
			return true
		}
	}
	return false
}

// LeakReport contains the results of leak detection
type LeakReport struct {
	StartTime        time.Time
	Duration         time.Duration
	InitialCount     int
	FinalCount       int
	LeakCount        int
	LeakDetected     bool
	InitialSnapshot  *GoroutineSnapshot
	FinalSnapshot    *GoroutineSnapshot
	LeakedGoroutines []GoroutineStack
}

// String returns a formatted string representation of the leak report
func (lr *LeakReport) String() string {
	status := "NO LEAKS"
	if lr.LeakDetected {
		status = "LEAKS DETECTED"
	}

	return fmt.Sprintf("LeakReport[%s]: %d -> %d goroutines (+%d) in %v",
		status, lr.InitialCount, lr.FinalCount, lr.LeakCount, lr.Duration)
}

// DetailedReport returns a detailed report with stack traces
func (lr *LeakReport) DetailedReport() string {
	var report strings.Builder

	report.WriteString(lr.String())
	report.WriteString("\n")

	if lr.LeakDetected && len(lr.LeakedGoroutines) > 0 {
		report.WriteString("\nLeaked Goroutines:\n")
		report.WriteString(strings.Repeat("=", 50))
		report.WriteString("\n")

		for i, stack := range lr.LeakedGoroutines {
			report.WriteString(fmt.Sprintf("\nLeak #%d:\n", i+1))
			report.WriteString(fmt.Sprintf("  Function: %s\n", stack.Function))
			report.WriteString(fmt.Sprintf("  State: %s\n", stack.State))
			report.WriteString(fmt.Sprintf("  File: %s:%d\n", stack.File, stack.Line))
			report.WriteString("  Stack Trace:\n")

			stackLines := strings.Split(strings.TrimSpace(stack.Stack), "\n")
			for _, line := range stackLines {
				if strings.TrimSpace(line) != "" {
					report.WriteString(fmt.Sprintf("    %s\n", line))
				}
			}
		}
	}

	return report.String()
}

// ===== GOROUTINE LEAK TEST CASES =====

// LeakyFunction demonstrates different types of goroutine leaks
type LeakyFunction struct {
	name string
	fn   func() error
}

// GetLeakyFunctions returns common goroutine leak patterns
func GetLeakyFunctions() []LeakyFunction {
	return []LeakyFunction{
		{
			name: "Channel Deadlock",
			fn: func() error {
				ch := make(chan int)

				// This goroutine will leak because channel is never read
				go func() {
					ch <- 42 // Blocks forever
				}()

				time.Sleep(10 * time.Millisecond)
				return nil
			},
		},
		{
			name: "Infinite Loop",
			fn: func() error {
				// This goroutine will leak because loop never exits
				go func() {
					for {
						time.Sleep(1 * time.Millisecond)
						// No exit condition
					}
				}()

				time.Sleep(10 * time.Millisecond)
				return nil
			},
		},
		{
			name: "Context Not Used",
			fn: func() error {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel() // Cancel is called, but goroutine doesn't check context

				go func() {
					for {
						time.Sleep(1 * time.Millisecond)
						// Should check: select { case <-ctx.Done(): return }
						// But this goroutine ignores ctx completely, causing leak
						_ = ctx // Suppress unused variable warning
					}
				}()

				time.Sleep(10 * time.Millisecond)
				return nil
			},
		},
		{
			name: "Ticker Not Stopped",
			fn: func() error {
				ticker := time.NewTicker(1 * time.Millisecond)
				// Missing: defer ticker.Stop()

				go func() {
					for range ticker.C {
						// Process ticks, but ticker never stops
					}
				}()

				time.Sleep(10 * time.Millisecond)
				return nil
			},
		},
		{
			name: "WaitGroup Never Done",
			fn: func() error {
				var wg sync.WaitGroup

				wg.Add(1)
				go func() {
					// Missing: defer wg.Done()
					time.Sleep(5 * time.Millisecond)
				}()

				// This would block forever: wg.Wait()
				time.Sleep(10 * time.Millisecond)
				return nil
			},
		},
		{
			name: "No Leak (Control)",
			fn: func() error {
				// This function should not leak goroutines
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				var wg sync.WaitGroup
				wg.Add(1)

				go func() {
					defer wg.Done()
					ticker := time.NewTicker(1 * time.Millisecond)
					defer ticker.Stop()

					for {
						select {
						case <-ctx.Done():
							return
						case <-ticker.C:
							// Do work
						}
					}
				}()

				time.Sleep(10 * time.Millisecond)
				cancel()
				wg.Wait()
				return nil
			},
		},
	}
}

// ===== GOLEAK INTEGRATION =====

// GoleakDetector provides integration with uber-go/goleak for production use
type GoleakDetector struct {
	options []interface{} // Would be goleak.Option in real implementation
}

// NewGoleakDetector creates a detector that mimics uber-go/goleak
func NewGoleakDetector() *GoleakDetector {
	return &GoleakDetector{}
}

// VerifyNoLeaks mimics goleak.VerifyNoLeaks functionality
func (gd *GoleakDetector) VerifyNoLeaks(testFunc func()) error {
	detector := NewGoroutineLeakDetector()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	report, err := detector.DetectLeaks(ctx, func() error {
		testFunc()
		return nil
	})

	if err != nil {
		return err
	}

	if report.LeakDetected {
		return fmt.Errorf("goroutine leak detected: %s", report.DetailedReport())
	}

	return nil
}

// ===== GOROUTINE MONITORING =====

// GoroutineMonitor continuously monitors goroutine count
type GoroutineMonitor struct {
	running       int32
	maxGoroutines int
	samples       []GoroutineSample
	mu            sync.RWMutex
}

// GoroutineSample represents a monitoring sample
type GoroutineSample struct {
	Timestamp      time.Time
	GoroutineCount int
	MemoryUsage    uint64
}

// NewGoroutineMonitor creates a new goroutine monitor
func NewGoroutineMonitor(maxGoroutines int) *GoroutineMonitor {
	return &GoroutineMonitor{
		maxGoroutines: maxGoroutines,
		samples:       make([]GoroutineSample, 0, 1000),
	}
}

// Start begins monitoring goroutines
func (gm *GoroutineMonitor) Start(ctx context.Context) {
	if !atomic.CompareAndSwapInt32(&gm.running, 0, 1) {
		return // Already running
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			atomic.StoreInt32(&gm.running, 0)
			return
		case <-ticker.C:
			gm.takeSample()
		}
	}
}

// takeSample records current goroutine count and memory usage
func (gm *GoroutineMonitor) takeSample() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	sample := GoroutineSample{
		Timestamp:      time.Now(),
		GoroutineCount: runtime.NumGoroutine(),
		MemoryUsage:    m.Alloc,
	}

	gm.mu.Lock()
	defer gm.mu.Unlock()

	gm.samples = append(gm.samples, sample)

	// Keep only the last 1000 samples
	if len(gm.samples) > 1000 {
		gm.samples = gm.samples[1:]
	}

	// Alert if goroutine count exceeds threshold
	if sample.GoroutineCount > gm.maxGoroutines {
		fmt.Printf("⚠️  Goroutine count alert: %d (threshold: %d)\n",
			sample.GoroutineCount, gm.maxGoroutines)
	}
}

// GetStats returns monitoring statistics
func (gm *GoroutineMonitor) GetStats() MonitoringStats {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	if len(gm.samples) == 0 {
		return MonitoringStats{}
	}

	counts := make([]int, len(gm.samples))
	for i, sample := range gm.samples {
		counts[i] = sample.GoroutineCount
	}

	sort.Ints(counts)

	stats := MonitoringStats{
		SampleCount:   len(gm.samples),
		MinGoroutines: counts[0],
		MaxGoroutines: counts[len(counts)-1],
		CurrentCount:  runtime.NumGoroutine(),
		StartTime:     gm.samples[0].Timestamp,
		EndTime:       gm.samples[len(gm.samples)-1].Timestamp,
	}

	// Calculate average
	sum := 0
	for _, count := range counts {
		sum += count
	}
	stats.AvgGoroutines = float64(sum) / float64(len(counts))

	// Calculate median
	if len(counts)%2 == 0 {
		mid := len(counts) / 2
		stats.MedianGoroutines = float64(counts[mid-1]+counts[mid]) / 2
	} else {
		stats.MedianGoroutines = float64(counts[len(counts)/2])
	}

	return stats
}

// MonitoringStats contains goroutine monitoring statistics
type MonitoringStats struct {
	SampleCount      int
	MinGoroutines    int
	MaxGoroutines    int
	AvgGoroutines    float64
	MedianGoroutines float64
	CurrentCount     int
	StartTime        time.Time
	EndTime          time.Time
}

// String returns a formatted string representation of the stats
func (ms *MonitoringStats) String() string {
	duration := ms.EndTime.Sub(ms.StartTime)
	return fmt.Sprintf("GoroutineStats: %d samples over %v | Min: %d, Max: %d, Avg: %.1f, Median: %.1f, Current: %d",
		ms.SampleCount, duration, ms.MinGoroutines, ms.MaxGoroutines,
		ms.AvgGoroutines, ms.MedianGoroutines, ms.CurrentCount)
}

// ===== MAIN DEMONSTRATION =====

func DemonstrateGoroutineLeakDetection() {
	fmt.Println("=== Go Goroutine Leak Detection ===")
	fmt.Println("Essential for production Go applications")
	fmt.Println("Used by companies like Uber, Google, and Netflix")
	fmt.Println()

	detector := NewGoroutineLeakDetector()
	leakyFunctions := GetLeakyFunctions()

	fmt.Printf("Initial goroutine count: %d\n", runtime.NumGoroutine())
	fmt.Println()

	for _, leakyFunc := range leakyFunctions {
		fmt.Printf("🔍 Testing: %s\n", leakyFunc.name)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		report, err := detector.DetectLeaks(ctx, leakyFunc.fn)
		cancel()

		if err != nil {
			fmt.Printf("   Error: %v\n", err)
		} else {
			fmt.Printf("   %s\n", report.String())
			if report.LeakDetected {
				fmt.Printf("   ⚠️  Leaked %d goroutines\n", len(report.LeakedGoroutines))
			} else {
				fmt.Printf("   ✅ No leaks detected\n")
			}
		}
		fmt.Println()

		// Brief pause between tests
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("Final goroutine count: %d\n", runtime.NumGoroutine())

	fmt.Println("\n🎯 KEY TAKEAWAYS:")
	fmt.Println("1. Always close channels, cancel contexts, and stop tickers")
	fmt.Println("2. Use uber-go/goleak in your test suite")
	fmt.Println("3. Monitor goroutine count in production")
	fmt.Println("4. Implement proper cleanup in your goroutines")
	fmt.Println("5. Test for leaks as part of your CI/CD pipeline")
}

// ===== MAIN FUNCTION - GOROUTINE LEAK DETECTION DEMO =====

func main() {
	fmt.Println("🔍 GOROUTINE LEAK DETECTION - UBER-GO/GOLEAK MASTERY")
	fmt.Println("=================================================")

	// Run the complete goroutine leak detection demonstration
	DemonstrateGoroutineLeakDetection()

	fmt.Println("\n🎯 GOROUTINE LEAK DETECTION MASTERY COMPLETE!")
	fmt.Println("You now have production-grade leak detection skills!")
}
