// Package performance_profiling demonstrates pprof integration for concurrency debugging
// This is essential for production Go applications and senior-level interview questions
package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Import for side effects
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"sync/atomic"
	"time"
)

// ===== PPROF INTEGRATION PATTERNS =====

// ProfiledApplication represents an application with built-in profiling
type ProfiledApplication struct {
	name            string
	profilingServer *http.Server
	workers         []*Worker
	stats           *ApplicationStats
}

// ApplicationStats tracks application metrics
type ApplicationStats struct {
	RequestsProcessed int64
	ErrorsOccurred    int64
	ActiveGoroutines  int64
	MemoryAllocated   uint64
	GCPauses          []time.Duration
}

// Worker represents a concurrent worker
type Worker struct {
	id      int
	tasks   chan Task
	stats   *WorkerStats
	running int32
	stopCh  chan struct{}
}

// Task represents a unit of work
type Task struct {
	ID       int
	Type     string
	Data     []byte
	Started  time.Time
	Duration time.Duration
}

// WorkerStats tracks per-worker metrics
type WorkerStats struct {
	TasksCompleted int64
	TotalDuration  time.Duration
	LastError      error
}

// NewProfiledApplication creates a new application with profiling enabled
func NewProfiledApplication(name string, numWorkers int) *ProfiledApplication {
	app := &ProfiledApplication{
		name:    name,
		workers: make([]*Worker, numWorkers),
		stats:   &ApplicationStats{},
	}

	// Create workers
	for i := 0; i < numWorkers; i++ {
		app.workers[i] = &Worker{
			id:     i,
			tasks:  make(chan Task, 100),
			stats:  &WorkerStats{},
			stopCh: make(chan struct{}),
		}
	}

	return app
}

// StartProfiling enables HTTP profiling endpoint
func (app *ProfiledApplication) StartProfiling(port string) error {
	mux := http.NewServeMux()

	// Default pprof handlers
	mux.HandleFunc("/debug/pprof/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))

	// Custom metrics endpoint
	mux.HandleFunc("/metrics", app.metricsHandler)
	mux.HandleFunc("/health", app.healthHandler)

	app.profilingServer = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("🔍 Profiling server started on :%s\n", port)
	fmt.Printf("   - CPU Profile: http://localhost:%s/debug/pprof/profile\n", port)
	fmt.Printf("   - Heap Profile: http://localhost:%s/debug/pprof/heap\n", port)
	fmt.Printf("   - Goroutine Profile: http://localhost:%s/debug/pprof/goroutine\n", port)
	fmt.Printf("   - Metrics: http://localhost:%s/metrics\n", port)

	go func() {
		if err := app.profilingServer.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Profiling server error: %v\n", err)
		}
	}()

	return nil
}

// metricsHandler provides custom application metrics
func (app *ProfiledApplication) metricsHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# Application Metrics for %s\n", app.name)
	fmt.Fprintf(w, "requests_processed_total %d\n", atomic.LoadInt64(&app.stats.RequestsProcessed))
	fmt.Fprintf(w, "errors_occurred_total %d\n", atomic.LoadInt64(&app.stats.ErrorsOccurred))
	fmt.Fprintf(w, "active_goroutines %d\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "memory_allocated_bytes %d\n", m.Alloc)
	fmt.Fprintf(w, "memory_sys_bytes %d\n", m.Sys)
	fmt.Fprintf(w, "gc_count %d\n", m.NumGC)
	fmt.Fprintf(w, "gc_pause_ns %d\n", m.PauseTotalNs)

	// Worker-specific metrics
	for i, worker := range app.workers {
		fmt.Fprintf(w, "worker_%d_tasks_completed %d\n", i, atomic.LoadInt64(&worker.stats.TasksCompleted))
		fmt.Fprintf(w, "worker_%d_running %d\n", i, atomic.LoadInt32(&worker.running))
	}
}

// healthHandler provides health check endpoint
func (app *ProfiledApplication) healthHandler(w http.ResponseWriter, r *http.Request) {
	status := "healthy"
	activeWorkers := 0

	for _, worker := range app.workers {
		if atomic.LoadInt32(&worker.running) == 1 {
			activeWorkers++
		}
	}

	if activeWorkers < len(app.workers)/2 {
		status = "degraded"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "%s", "active_workers": %d, "total_workers": %d}`,
		status, activeWorkers, len(app.workers))
}

// StartWorkers starts all worker goroutines
func (app *ProfiledApplication) StartWorkers(ctx context.Context) {
	fmt.Printf("🚀 Starting %d workers\n", len(app.workers))

	for _, worker := range app.workers {
		go worker.Run(ctx)
	}
}

// SubmitTask submits a task to be processed
func (app *ProfiledApplication) SubmitTask(task Task) {
	// Simple round-robin task distribution
	workerIndex := int(atomic.AddInt64(&app.stats.RequestsProcessed, 1)-1) % len(app.workers)

	select {
	case app.workers[workerIndex].tasks <- task:
		// Task submitted successfully
	default:
		// Worker queue full, increment error count
		atomic.AddInt64(&app.stats.ErrorsOccurred, 1)
	}
}

// Run executes a worker
func (w *Worker) Run(ctx context.Context) {
	atomic.StoreInt32(&w.running, 1)
	defer atomic.StoreInt32(&w.running, 0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		case task := <-w.tasks:
			w.processTask(task)
		}
	}
}

// processTask processes a single task
func (w *Worker) processTask(task Task) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		atomic.AddInt64(&w.stats.TasksCompleted, 1)
		w.stats.TotalDuration += duration
	}()

	// Simulate different types of work
	switch task.Type {
	case "cpu_intensive":
		w.cpuIntensiveWork()
	case "memory_intensive":
		w.memoryIntensiveWork()
	case "io_intensive":
		w.ioIntensiveWork()
	default:
		w.defaultWork()
	}
}

// cpuIntensiveWork simulates CPU-bound work
func (w *Worker) cpuIntensiveWork() {
	// Prime number calculation
	for i := 2; i < 10000; i++ {
		isPrime := true
		for j := 2; j*j <= i; j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}
		_ = isPrime
	}
}

// memoryIntensiveWork simulates memory-bound work
func (w *Worker) memoryIntensiveWork() {
	// Allocate and manipulate large slices
	data := make([][]byte, 1000)
	for i := range data {
		data[i] = make([]byte, 1024) // 1KB each
		// Fill with pattern
		for j := range data[i] {
			data[i][j] = byte(i + j)
		}
	}

	// Process data to prevent optimization
	sum := 0
	for _, slice := range data {
		for _, b := range slice {
			sum += int(b)
		}
	}
	_ = sum
}

// ioIntensiveWork simulates I/O-bound work
func (w *Worker) ioIntensiveWork() {
	// Simulate network/disk I/O with sleep
	time.Sleep(time.Duration(10+w.id) * time.Millisecond)
}

// defaultWork simulates balanced work
func (w *Worker) defaultWork() {
	// Mix of CPU and I/O
	for i := 0; i < 1000; i++ {
		_ = i * i
	}
	time.Sleep(1 * time.Millisecond)
}

// ===== PROFILING UTILITIES =====

// ProfileManager handles different types of profiling
type ProfileManager struct {
	outputDir string
}

// NewProfileManager creates a new profile manager
func NewProfileManager(outputDir string) *ProfileManager {
	// Create output directory if it doesn't exist
	os.MkdirAll(outputDir, 0755)

	return &ProfileManager{
		outputDir: outputDir,
	}
}

// CPUProfile runs CPU profiling for the specified duration
func (pm *ProfileManager) CPUProfile(duration time.Duration, workload func()) error {
	filename := fmt.Sprintf("%s/cpu_profile_%d.prof", pm.outputDir, time.Now().Unix())

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create CPU profile: %w", err)
	}
	defer f.Close()

	fmt.Printf("🔍 Starting CPU profile (%v) -> %s\n", duration, filename)

	if err := pprof.StartCPUProfile(f); err != nil {
		return fmt.Errorf("could not start CPU profile: %w", err)
	}
	defer pprof.StopCPUProfile()

	// Run workload
	done := make(chan struct{})
	go func() {
		defer close(done)
		workload()
	}()

	select {
	case <-done:
		fmt.Println("✅ CPU profile completed (workload finished)")
	case <-time.After(duration):
		fmt.Println("✅ CPU profile completed (timeout)")
	}

	return nil
}

// MemoryProfile captures a heap profile
func (pm *ProfileManager) MemoryProfile(workload func()) error {
	filename := fmt.Sprintf("%s/heap_profile_%d.prof", pm.outputDir, time.Now().Unix())

	// Run workload first
	fmt.Println("🔍 Running workload for memory profile...")
	workload()

	// Force GC to get accurate heap snapshot
	runtime.GC()

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create memory profile: %w", err)
	}
	defer f.Close()

	fmt.Printf("📊 Writing memory profile -> %s\n", filename)

	if err := pprof.WriteHeapProfile(f); err != nil {
		return fmt.Errorf("could not write memory profile: %w", err)
	}

	fmt.Println("✅ Memory profile completed")
	return nil
}

// GoroutineProfile captures a goroutine profile
func (pm *ProfileManager) GoroutineProfile(workload func()) error {
	filename := fmt.Sprintf("%s/goroutine_profile_%d.prof", pm.outputDir, time.Now().Unix())

	// Run workload
	fmt.Println("🔍 Running workload for goroutine profile...")
	workload()

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create goroutine profile: %w", err)
	}
	defer f.Close()

	fmt.Printf("🦍 Writing goroutine profile -> %s\n", filename)

	profile := pprof.Lookup("goroutine")
	if err := profile.WriteTo(f, 0); err != nil {
		return fmt.Errorf("could not write goroutine profile: %w", err)
	}

	fmt.Println("✅ Goroutine profile completed")
	return nil
}

// TraceProfile captures an execution trace
func (pm *ProfileManager) TraceProfile(duration time.Duration, workload func()) error {
	filename := fmt.Sprintf("%s/trace_%d.trace", pm.outputDir, time.Now().Unix())

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create trace file: %w", err)
	}
	defer f.Close()

	fmt.Printf("📋 Starting execution trace (%v) -> %s\n", duration, filename)

	if err := trace.Start(f); err != nil {
		return fmt.Errorf("could not start trace: %w", err)
	}
	defer trace.Stop()

	// Run workload
	done := make(chan struct{})
	go func() {
		defer close(done)
		workload()
	}()

	select {
	case <-done:
		fmt.Println("✅ Trace completed (workload finished)")
	case <-time.After(duration):
		fmt.Println("✅ Trace completed (timeout)")
	}

	return nil
}

// ===== WORKLOAD GENERATORS =====

// ConcurrencyWorkload generates concurrent workload for profiling
func ConcurrencyWorkload(numGoroutines int, duration time.Duration) func() {
	return func() {
		var wg sync.WaitGroup
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		// CPU-intensive workers
		for i := 0; i < numGoroutines/3; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						// CPU work
						sum := 0
						for j := 0; j < 10000; j++ {
							sum += j * j
						}
						_ = sum
					}
				}
			}(i)
		}

		// Memory-intensive workers
		for i := 0; i < numGoroutines/3; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						// Memory allocation
						data := make([]byte, 64*1024) // 64KB
						for j := range data {
							data[j] = byte(j % 256)
						}
						time.Sleep(10 * time.Millisecond)
					}
				}
			}(i)
		}

		// I/O-intensive workers
		remaining := numGoroutines - 2*(numGoroutines/3)
		for i := 0; i < remaining; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						// Simulate I/O
						time.Sleep(50 * time.Millisecond)
					}
				}
			}(i)
		}

		wg.Wait()
	}
}

// ===== MAIN DEMONSTRATION =====

func DemonstrateProfiling() {
	fmt.Println("=== Go Performance Profiling for Concurrency ===")
	fmt.Println("Production-grade profiling techniques used by FAANG companies")
	fmt.Println()

	// Create profile manager
	profileManager := NewProfileManager("./profiles")

	// Create workload
	workload := ConcurrencyWorkload(50, 3*time.Second)

	fmt.Println("🎯 Running different profiling scenarios...")

	// CPU Profile
	if err := profileManager.CPUProfile(2*time.Second, workload); err != nil {
		fmt.Printf("CPU profiling error: %v\n", err)
	}

	time.Sleep(500 * time.Millisecond)

	// Memory Profile
	if err := profileManager.MemoryProfile(workload); err != nil {
		fmt.Printf("Memory profiling error: %v\n", err)
	}

	time.Sleep(500 * time.Millisecond)

	// Goroutine Profile
	if err := profileManager.GoroutineProfile(workload); err != nil {
		fmt.Printf("Goroutine profiling error: %v\n", err)
	}

	time.Sleep(500 * time.Millisecond)

	// Execution Trace
	if err := profileManager.TraceProfile(2*time.Second, workload); err != nil {
		fmt.Printf("Trace profiling error: %v\n", err)
	}

	// Demonstration of profiled application
	fmt.Println("\n🏭 Demonstrating Profiled Application...")
	app := NewProfiledApplication("demo-app", 4)

	// Start profiling server (but don't block on it)
	app.StartProfiling("8080")

	// Start workers
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	app.StartWorkers(ctx)

	// Submit tasks
	go func() {
		taskID := 0
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				taskTypes := []string{"cpu_intensive", "memory_intensive", "io_intensive", "default"}
				taskType := taskTypes[taskID%len(taskTypes)]

				task := Task{
					ID:      taskID,
					Type:    taskType,
					Data:    make([]byte, 256),
					Started: time.Now(),
				}

				app.SubmitTask(task)
				taskID++
			}
		}
	}()

	// Wait for completion
	<-ctx.Done()

	// Final stats
	fmt.Printf("\n📊 Final Statistics:\n")
	fmt.Printf("   Requests Processed: %d\n", atomic.LoadInt64(&app.stats.RequestsProcessed))
	fmt.Printf("   Errors Occurred: %d\n", atomic.LoadInt64(&app.stats.ErrorsOccurred))
	fmt.Printf("   Active Goroutines: %d\n", runtime.NumGoroutine())

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("   Memory Allocated: %d bytes\n", m.Alloc)
	fmt.Printf("   GC Count: %d\n", m.NumGC)

	fmt.Println("\n🎯 KEY TAKEAWAYS:")
	fmt.Println("1. Use pprof for production performance analysis")
	fmt.Println("2. CPU profiles identify hot code paths")
	fmt.Println("3. Memory profiles find allocation bottlenecks")
	fmt.Println("4. Goroutine profiles detect concurrency issues")
	fmt.Println("5. Execution traces show scheduling behavior")
	fmt.Println("6. Integrate profiling endpoints in your applications")
	fmt.Println("7. Monitor continuously in production")

	fmt.Println("\n📋 ANALYSIS COMMANDS:")
	fmt.Println("   go tool pprof cpu_profile.prof")
	fmt.Println("   go tool pprof -http=:8081 heap_profile.prof")
	fmt.Println("   go tool trace trace.trace")
	fmt.Println("   curl http://localhost:8080/debug/pprof/goroutine?debug=1")
}

// ===== MAIN FUNCTION - COMPLETE PROFILING DEMONSTRATION =====

func main() {
	fmt.Println("🚀 GO PERFORMANCE PROFILING MASTERY")
	fmt.Println("===================================")

	// Run the complete profiling demonstration
	DemonstrateProfiling()

	fmt.Println("\n🎯 PROFILING MASTERY COMPLETE!")
	fmt.Println("You now have production-grade Go profiling skills!")
}
