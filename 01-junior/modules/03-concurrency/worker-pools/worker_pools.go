// Worker Pool Patterns - Advanced Concurrency Management
// Used by: HTTP servers, background job processors, data pipelines
// Problem: Efficiently manage concurrent workers for different workload patterns
// Solution: Multiple worker pool strategies optimized for different scenarios

package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ==============================================================================
// 1. FIXED WORKER POOL - Classic Pattern
// ==============================================================================

// FixedWorkerPool implements traditional fixed-size worker pool
// Used in: HTTP servers (like Nginx), database connection pools
type FixedWorkerPool struct {
	workers    int
	jobQueue   chan Job
	resultChan chan JobResult
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc

	// Metrics
	totalJobs     int64
	completedJobs int64
	failedJobs    int64
	activeWorkers int64
}

// Job represents a unit of work
type Job struct {
	ID       string
	Data     interface{}
	Priority int
}

// JobResult represents the result of job execution
type JobResult struct {
	JobID    string
	Result   interface{}
	Error    error
	Duration time.Duration
}

// NewFixedWorkerPool creates a new fixed worker pool
func NewFixedWorkerPool(workers, bufferSize int) *FixedWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &FixedWorkerPool{
		workers:    workers,
		jobQueue:   make(chan Job, bufferSize),
		resultChan: make(chan JobResult, bufferSize),
		ctx:        ctx,
		cancel:     cancel,
	}

	// Start workers
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

// worker executes jobs from the queue
func (p *FixedWorkerPool) worker(id int) {
	defer p.wg.Done()
	atomic.AddInt64(&p.activeWorkers, 1)

	fmt.Printf("🏃‍♂️ Worker %d started\n", id)

	for {
		select {
		case job, ok := <-p.jobQueue:
			if !ok {
				fmt.Printf("🏁 Worker %d shutting down (queue closed)\n", id)
				atomic.AddInt64(&p.activeWorkers, -1)
				return
			}

			// Process job
			result := p.processJob(job, id)

			// Send result
			select {
			case p.resultChan <- result:
			case <-p.ctx.Done():
				fmt.Printf("🏁 Worker %d shutting down (context cancelled)\n", id)
				atomic.AddInt64(&p.activeWorkers, -1)
				return
			}

		case <-p.ctx.Done():
			fmt.Printf("🏁 Worker %d shutting down (context cancelled)\n", id)
			atomic.AddInt64(&p.activeWorkers, -1)
			return
		}
	}
}

// processJob executes a single job
func (p *FixedWorkerPool) processJob(job Job, workerID int) JobResult {
	start := time.Now()

	fmt.Printf("   Worker %d processing job %s\n", workerID, job.ID)

	// Simulate work based on job data
	var result interface{}
	var err error

	if data, ok := job.Data.(int); ok {
		// Simulate processing time
		processingTime := time.Duration(50+data%100) * time.Millisecond
		time.Sleep(processingTime)

		// Simulate occasional failure
		if data%10 == 0 {
			err = fmt.Errorf("simulated failure for job %s", job.ID)
			atomic.AddInt64(&p.failedJobs, 1)
		} else {
			result = data * 2 // Simple processing
			atomic.AddInt64(&p.completedJobs, 1)
		}
	}

	duration := time.Since(start)

	return JobResult{
		JobID:    job.ID,
		Result:   result,
		Error:    err,
		Duration: duration,
	}
}

// Submit adds a job to the pool
func (p *FixedWorkerPool) Submit(job Job) error {
	atomic.AddInt64(&p.totalJobs, 1)

	select {
	case p.jobQueue <- job:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("pool is shutting down")
	default:
		return fmt.Errorf("job queue is full")
	}
}

// Results returns the result channel
func (p *FixedWorkerPool) Results() <-chan JobResult {
	return p.resultChan
}

// Stats returns pool statistics
func (p *FixedWorkerPool) Stats() (total, completed, failed, active int64) {
	return atomic.LoadInt64(&p.totalJobs),
		atomic.LoadInt64(&p.completedJobs),
		atomic.LoadInt64(&p.failedJobs),
		atomic.LoadInt64(&p.activeWorkers)
}

// Shutdown gracefully shuts down the pool
func (p *FixedWorkerPool) Shutdown() {
	close(p.jobQueue)
	p.wg.Wait()
	close(p.resultChan)
	p.cancel()
}

// ==============================================================================
// 2. DYNAMIC WORKER POOL - Auto-scaling Pattern
// ==============================================================================

// DynamicWorkerPool implements auto-scaling worker pool
// Used in: Cloud services, adaptive load handling
type DynamicWorkerPool struct {
	minWorkers     int
	maxWorkers     int
	currentWorkers int64
	jobQueue       chan Job
	resultChan     chan JobResult
	workerWG       sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc

	// Scaling metrics
	queueLength        int64
	avgProcessTime     int64 // nanoseconds
	scaleUpThreshold   int
	scaleDownThreshold int
	lastScaleTime      time.Time
	scaleCooldown      time.Duration
	mu                 sync.RWMutex

	// Statistics
	totalJobs     int64
	completedJobs int64
	scalingEvents int64
}

// NewDynamicWorkerPool creates a new dynamic worker pool
func NewDynamicWorkerPool(minWorkers, maxWorkers, bufferSize int) *DynamicWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &DynamicWorkerPool{
		minWorkers:         minWorkers,
		maxWorkers:         maxWorkers,
		currentWorkers:     0,
		jobQueue:           make(chan Job, bufferSize),
		resultChan:         make(chan JobResult, bufferSize),
		ctx:                ctx,
		cancel:             cancel,
		scaleUpThreshold:   bufferSize / 2,  // Scale up when queue is 50% full
		scaleDownThreshold: bufferSize / 10, // Scale down when queue is 10% full
		scaleCooldown:      5 * time.Second,
	}

	// Start minimum number of workers
	for i := 0; i < minWorkers; i++ {
		pool.startWorker()
	}

	// Start scaling monitor
	go pool.scalingMonitor()

	return pool
}

// startWorker starts a new worker
func (p *DynamicWorkerPool) startWorker() {
	workerID := atomic.AddInt64(&p.currentWorkers, 1)

	p.workerWG.Add(1)
	go func(id int64) {
		defer p.workerWG.Done()
		defer atomic.AddInt64(&p.currentWorkers, -1)

		fmt.Printf("🚀 Dynamic worker %d started (total: %d)\n", id, atomic.LoadInt64(&p.currentWorkers))

		for {
			select {
			case job, ok := <-p.jobQueue:
				if !ok {
					fmt.Printf("🏁 Dynamic worker %d shutting down\n", id)
					return
				}

				start := time.Now()
				result := p.processJob(job, int(id))
				duration := time.Since(start)

				// Update average processing time
				atomic.StoreInt64(&p.avgProcessTime, duration.Nanoseconds())
				atomic.AddInt64(&p.completedJobs, 1)

				// Send result
				select {
				case p.resultChan <- result:
				case <-p.ctx.Done():
					fmt.Printf("🏁 Dynamic worker %d shutting down (context cancelled)\n", id)
					return
				}

			case <-p.ctx.Done():
				fmt.Printf("🏁 Dynamic worker %d shutting down (context cancelled)\n", id)
				return
			}
		}
	}(workerID)
}

// scalingMonitor monitors queue and scales workers
func (p *DynamicWorkerPool) scalingMonitor() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.checkAndScale()
		case <-p.ctx.Done():
			return
		}
	}
}

// checkAndScale evaluates if scaling is needed
func (p *DynamicWorkerPool) checkAndScale() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Cooldown period to prevent thrashing
	if time.Since(p.lastScaleTime) < p.scaleCooldown {
		return
	}

	currentWorkers := int(atomic.LoadInt64(&p.currentWorkers))
	queueLen := len(p.jobQueue)
	atomic.StoreInt64(&p.queueLength, int64(queueLen))

	// Scale up condition
	if queueLen >= p.scaleUpThreshold && currentWorkers < p.maxWorkers {
		p.startWorker()
		p.lastScaleTime = time.Now()
		atomic.AddInt64(&p.scalingEvents, 1)
		fmt.Printf("📈 Scaled UP: %d -> %d workers (queue: %d)\n",
			currentWorkers, currentWorkers+1, queueLen)
	}

	// Scale down condition (only if we have more than minimum workers)
	if queueLen <= p.scaleDownThreshold && currentWorkers > p.minWorkers {
		// Scale down by not replacing a worker that finishes
		// This is a simplified approach; production systems use more sophisticated methods
		fmt.Printf("📉 Scale DOWN opportunity: %d workers (queue: %d)\n",
			currentWorkers, queueLen)
	}
}

// processJob executes a job (similar to FixedWorkerPool but with metrics)
func (p *DynamicWorkerPool) processJob(job Job, workerID int) JobResult {
	start := time.Now()

	var result interface{}
	var err error

	if data, ok := job.Data.(int); ok {
		// Variable processing time
		processingTime := time.Duration(30+data%150) * time.Millisecond
		time.Sleep(processingTime)
		result = data * 3
	}

	return JobResult{
		JobID:    job.ID,
		Result:   result,
		Error:    err,
		Duration: time.Since(start),
	}
}

// Submit adds a job to the dynamic pool
func (p *DynamicWorkerPool) Submit(job Job) error {
	atomic.AddInt64(&p.totalJobs, 1)

	select {
	case p.jobQueue <- job:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("pool is shutting down")
	default:
		return fmt.Errorf("job queue is full")
	}
}

// Results returns the result channel
func (p *DynamicWorkerPool) Results() <-chan JobResult {
	return p.resultChan
}

// Stats returns comprehensive pool statistics
func (p *DynamicWorkerPool) Stats() (total, completed, workers, scalingEvents int64, avgTime time.Duration) {
	avgNanos := atomic.LoadInt64(&p.avgProcessTime)
	return atomic.LoadInt64(&p.totalJobs),
		atomic.LoadInt64(&p.completedJobs),
		atomic.LoadInt64(&p.currentWorkers),
		atomic.LoadInt64(&p.scalingEvents),
		time.Duration(avgNanos)
}

// Shutdown gracefully shuts down the dynamic pool
func (p *DynamicWorkerPool) Shutdown() {
	p.cancel()
	close(p.jobQueue)
	p.workerWG.Wait()
	close(p.resultChan)
}

// ==============================================================================
// 3. PRIORITY WORKER POOL - Priority-based Processing
// ==============================================================================

// PriorityWorkerPool implements priority-based job processing
// Used in: Task schedulers, request prioritization systems
type PriorityWorkerPool struct {
	workers      int
	highPriority chan Job
	medPriority  chan Job
	lowPriority  chan Job
	resultChan   chan JobResult
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc

	// Priority statistics
	highProcessed int64
	medProcessed  int64
	lowProcessed  int64
}

// NewPriorityWorkerPool creates a new priority-based worker pool
func NewPriorityWorkerPool(workers, bufferSize int) *PriorityWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &PriorityWorkerPool{
		workers:      workers,
		highPriority: make(chan Job, bufferSize),
		medPriority:  make(chan Job, bufferSize),
		lowPriority:  make(chan Job, bufferSize),
		resultChan:   make(chan JobResult, bufferSize*3),
		ctx:          ctx,
		cancel:       cancel,
	}

	// Start workers
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.priorityWorker(i)
	}

	return pool
}

// priorityWorker processes jobs with priority consideration
func (p *PriorityWorkerPool) priorityWorker(id int) {
	defer p.wg.Done()

	fmt.Printf("🎯 Priority worker %d started\n", id)

	for {
		select {
		// Highest priority first
		case job := <-p.highPriority:
			result := p.processJobWithPriority(job, "HIGH", id)
			atomic.AddInt64(&p.highProcessed, 1)
			p.resultChan <- result

		// Check medium priority if no high priority jobs
		case job := <-p.medPriority:
			result := p.processJobWithPriority(job, "MEDIUM", id)
			atomic.AddInt64(&p.medProcessed, 1)
			p.resultChan <- result

		// Check low priority if no high or medium priority jobs
		case job := <-p.lowPriority:
			result := p.processJobWithPriority(job, "LOW", id)
			atomic.AddInt64(&p.lowProcessed, 1)
			p.resultChan <- result

		case <-p.ctx.Done():
			fmt.Printf("🏁 Priority worker %d shutting down\n", id)
			return
		}
	}
}

// processJobWithPriority processes a job with priority awareness
func (p *PriorityWorkerPool) processJobWithPriority(job Job, priority string, workerID int) JobResult {
	start := time.Now()

	fmt.Printf("   Worker %d processing %s priority job %s\n", workerID, priority, job.ID)

	var result interface{}
	var err error

	if data, ok := job.Data.(int); ok {
		// Processing time varies by priority
		var processingTime time.Duration
		switch priority {
		case "HIGH":
			processingTime = time.Duration(20+data%30) * time.Millisecond // Fast processing
		case "MEDIUM":
			processingTime = time.Duration(50+data%50) * time.Millisecond // Medium processing
		case "LOW":
			processingTime = time.Duration(100+data%100) * time.Millisecond // Slower processing
		}

		time.Sleep(processingTime)
		result = fmt.Sprintf("%s priority result: %d", priority, data*job.Priority)
	}

	return JobResult{
		JobID:    job.ID,
		Result:   result,
		Error:    err,
		Duration: time.Since(start),
	}
}

// Submit submits a job based on its priority
func (p *PriorityWorkerPool) Submit(job Job) error {
	var targetChan chan Job

	switch job.Priority {
	case 3: // High priority
		targetChan = p.highPriority
	case 2: // Medium priority
		targetChan = p.medPriority
	case 1: // Low priority
		targetChan = p.lowPriority
	default:
		targetChan = p.lowPriority // Default to low priority
	}

	select {
	case targetChan <- job:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("pool is shutting down")
	default:
		return fmt.Errorf("priority queue is full")
	}
}

// Results returns the result channel
func (p *PriorityWorkerPool) Results() <-chan JobResult {
	return p.resultChan
}

// Stats returns priority-based statistics
func (p *PriorityWorkerPool) Stats() (high, medium, low int64) {
	return atomic.LoadInt64(&p.highProcessed),
		atomic.LoadInt64(&p.medProcessed),
		atomic.LoadInt64(&p.lowProcessed)
}

// Shutdown gracefully shuts down the priority pool
func (p *PriorityWorkerPool) Shutdown() {
	p.cancel()
	close(p.highPriority)
	close(p.medPriority)
	close(p.lowPriority)
	p.wg.Wait()
	close(p.resultChan)
}

// ==============================================================================
// 4. BENCHMARK AND COMPARISON
// ==============================================================================

// PoolBenchmark compares different worker pool implementations
type PoolBenchmark struct {
	jobCount int
	jobs     []Job
}

// NewPoolBenchmark creates a new benchmark suite
func NewPoolBenchmark(jobCount int) *PoolBenchmark {
	jobs := make([]Job, jobCount)
	for i := 0; i < jobCount; i++ {
		jobs[i] = Job{
			ID:       fmt.Sprintf("job-%d", i+1),
			Data:     rand.Intn(1000),
			Priority: (i % 3) + 1, // Priority 1-3
		}
	}

	return &PoolBenchmark{
		jobCount: jobCount,
		jobs:     jobs,
	}
}

// BenchmarkFixedPool benchmarks the fixed worker pool
func (pb *PoolBenchmark) BenchmarkFixedPool(workers int) time.Duration {
	fmt.Printf("\n📊 Benchmarking Fixed Pool (%d workers, %d jobs)\n", workers, pb.jobCount)

	pool := NewFixedWorkerPool(workers, pb.jobCount)
	defer pool.Shutdown()

	start := time.Now()

	// Submit all jobs
	go func() {
		for _, job := range pb.jobs {
			pool.Submit(job)
		}
	}()

	// Collect results
	completed := 0
	for result := range pool.Results() {
		if result.Error != nil {
			fmt.Printf("   Job %s failed: %v\n", result.JobID, result.Error)
		}
		completed++
		if completed >= pb.jobCount {
			break
		}
	}

	duration := time.Since(start)
	total, comp, failed, active := pool.Stats()

	fmt.Printf("   ✅ Completed: %d/%d, Failed: %d, Active Workers: %d\n", comp, total, failed, active)
	fmt.Printf("   ⏱️  Total Time: %v\n", duration)

	return duration
}

// BenchmarkDynamicPool benchmarks the dynamic worker pool
func (pb *PoolBenchmark) BenchmarkDynamicPool(minWorkers, maxWorkers int) time.Duration {
	fmt.Printf("\n📊 Benchmarking Dynamic Pool (%d-%d workers, %d jobs)\n", minWorkers, maxWorkers, pb.jobCount)

	pool := NewDynamicWorkerPool(minWorkers, maxWorkers, pb.jobCount)
	defer pool.Shutdown()

	start := time.Now()

	// Submit jobs with some delay to trigger scaling
	go func() {
		for i, job := range pb.jobs {
			pool.Submit(job)
			if i%50 == 0 { // Add small delay every 50 jobs to trigger scaling
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Collect results
	completed := 0
	for range pool.Results() {
		completed++
		if completed >= pb.jobCount {
			break
		}
	}

	duration := time.Since(start)
	total, comp, workers, scalingEvents, avgTime := pool.Stats()

	fmt.Printf("   ✅ Completed: %d/%d, Final Workers: %d, Scaling Events: %d\n", comp, total, workers, scalingEvents)
	fmt.Printf("   ⏱️  Total Time: %v, Avg Job Time: %v\n", duration, avgTime)

	return duration
}

// BenchmarkPriorityPool benchmarks the priority worker pool
func (pb *PoolBenchmark) BenchmarkPriorityPool(workers int) time.Duration {
	fmt.Printf("\n📊 Benchmarking Priority Pool (%d workers, %d jobs)\n", workers, pb.jobCount)

	pool := NewPriorityWorkerPool(workers, pb.jobCount/3)
	defer pool.Shutdown()

	start := time.Now()

	// Submit all jobs
	go func() {
		for _, job := range pb.jobs {
			pool.Submit(job)
		}
	}()

	// Collect results
	completed := 0
	for range pool.Results() {
		completed++
		if completed >= pb.jobCount {
			break
		}
	}

	duration := time.Since(start)
	high, medium, low := pool.Stats()

	fmt.Printf("   ✅ High: %d, Medium: %d, Low: %d\n", high, medium, low)
	fmt.Printf("   ⏱️  Total Time: %v\n", duration)

	return duration
}

func main() {
	fmt.Println("👷‍♂️ Worker Pool Patterns Demo")
	fmt.Println("===============================")

	// Demo 1: Fixed Worker Pool
	fmt.Println("\n1️⃣ Fixed Worker Pool Demo")
	fixedPool := NewFixedWorkerPool(3, 10)

	// Submit some jobs
	jobs := []Job{
		{ID: "job-1", Data: 10, Priority: 1},
		{ID: "job-2", Data: 20, Priority: 2},
		{ID: "job-3", Data: 30, Priority: 3},
		{ID: "job-4", Data: 40, Priority: 1},
		{ID: "job-5", Data: 50, Priority: 2},
	}

	for _, job := range jobs {
		fixedPool.Submit(job)
	}

	// Collect results
	go func() {
		for i := 0; i < len(jobs); i++ {
			result := <-fixedPool.Results()
			if result.Error != nil {
				fmt.Printf("❌ %s failed: %v\n", result.JobID, result.Error)
			} else {
				fmt.Printf("✅ %s: %v (took %v)\n", result.JobID, result.Result, result.Duration)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	fixedPool.Shutdown()

	// Demo 2: Dynamic Worker Pool
	fmt.Println("\n2️⃣ Dynamic Worker Pool Demo")
	dynamicPool := NewDynamicWorkerPool(2, 6, 20)

	// Submit jobs in bursts to trigger scaling
	go func() {
		for i := 0; i < 15; i++ {
			job := Job{
				ID:   fmt.Sprintf("dynamic-job-%d", i+1),
				Data: rand.Intn(100),
			}
			dynamicPool.Submit(job)

			if i == 7 { // Create a burst
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Monitor results
	go func() {
		for i := 0; i < 15; i++ {
			result := <-dynamicPool.Results()
			fmt.Printf("⚡ %s: %v (took %v)\n", result.JobID, result.Result, result.Duration)
		}
	}()

	time.Sleep(5 * time.Second)
	dynamicPool.Shutdown()

	// Demo 3: Priority Worker Pool
	fmt.Println("\n3️⃣ Priority Worker Pool Demo")
	priorityPool := NewPriorityWorkerPool(2, 5)

	// Submit jobs with different priorities
	priorityJobs := []Job{
		{ID: "low-1", Data: 100, Priority: 1},  // Low priority
		{ID: "high-1", Data: 200, Priority: 3}, // High priority
		{ID: "med-1", Data: 150, Priority: 2},  // Medium priority
		{ID: "low-2", Data: 110, Priority: 1},  // Low priority
		{ID: "high-2", Data: 220, Priority: 3}, // High priority
		{ID: "med-2", Data: 160, Priority: 2},  // Medium priority
	}

	for _, job := range priorityJobs {
		priorityPool.Submit(job)
	}

	// Collect results (high priority should be processed first)
	go func() {
		for i := 0; i < len(priorityJobs); i++ {
			result := <-priorityPool.Results()
			fmt.Printf("🎯 %s: %v (took %v)\n", result.JobID, result.Result, result.Duration)
		}
	}()

	time.Sleep(3 * time.Second)
	priorityPool.Shutdown()

	// Demo 4: Performance Benchmark
	fmt.Println("\n4️⃣ Performance Benchmark")
	benchmark := NewPoolBenchmark(200) // 200 jobs

	// Benchmark different configurations
	fixedTime := benchmark.BenchmarkFixedPool(4)
	dynamicTime := benchmark.BenchmarkDynamicPool(2, 8)
	priorityTime := benchmark.BenchmarkPriorityPool(4)

	fmt.Printf("\n🏆 Performance Summary:\n")
	fmt.Printf("   Fixed Pool (4 workers):    %v\n", fixedTime)
	fmt.Printf("   Dynamic Pool (2-8 workers): %v\n", dynamicTime)
	fmt.Printf("   Priority Pool (4 workers):  %v\n", priorityTime)

	// Find the fastest
	fastest := "Fixed"
	fastestTime := fixedTime

	if dynamicTime < fastestTime {
		fastest = "Dynamic"
		fastestTime = dynamicTime
	}

	if priorityTime < fastestTime {
		fastest = "Priority"
		fastestTime = priorityTime
	}

	fmt.Printf("   🥇 Fastest: %s Pool (%v)\n", fastest, fastestTime)

	fmt.Println("\n🎯 Worker Pool Patterns Complete!")
	fmt.Printf("Key Insights:\n")
	fmt.Printf("• Fixed Pool: Predictable resource usage, consistent performance\n")
	fmt.Printf("• Dynamic Pool: Adapts to load, efficient resource utilization\n")
	fmt.Printf("• Priority Pool: Ensures critical tasks are processed first\n")
	fmt.Printf("• Choose pattern based on: workload predictability, resource constraints, SLA requirements\n")

	fmt.Printf("\nRuntime Stats:\n")
	fmt.Printf("• Goroutines: %d\n", runtime.NumGoroutine())
	fmt.Printf("• CPUs: %d\n", runtime.NumCPU())
}
