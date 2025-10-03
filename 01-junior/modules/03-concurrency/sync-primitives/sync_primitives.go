// Synchronization Primitives - Advanced Concurrency Control
// Used by: Every major Go application (Docker, Kubernetes, Prometheus)
// Problem: Coordinate access to shared resources safely and efficiently
// Solution: Mutex, RWMutex, atomic operations, and advanced sync patterns

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
// 1. MUTEX PATTERNS - Mutual Exclusion
// ==============================================================================

// SafeCounter demonstrates basic mutex usage for protecting shared state
// Used in: Prometheus metrics, HTTP server request counters
type SafeCounter struct {
	mu    sync.Mutex
	value int64
}

// Increment safely increments the counter
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Value safely reads the counter value
func (c *SafeCounter) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// CacheWithMutex demonstrates a thread-safe cache implementation
// Used in: Redis clients, HTTP caches, session stores
type CacheWithMutex struct {
	mu   sync.Mutex
	data map[string]interface{}
}

// NewCache creates a new thread-safe cache
func NewCache() *CacheWithMutex {
	return &CacheWithMutex{
		data: make(map[string]interface{}),
	}
}

// Set stores a value in the cache
func (c *CacheWithMutex) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Get retrieves a value from the cache
func (c *CacheWithMutex) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, exists := c.data[key]
	return value, exists
}

// Delete removes a value from the cache
func (c *CacheWithMutex) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Size returns the number of items in the cache
func (c *CacheWithMutex) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}

// ==============================================================================
// 2. RWMUTEX PATTERNS - Reader-Writer Lock
// ==============================================================================

// ReadWriteCache demonstrates RWMutex for read-heavy workloads
// Used in: Configuration stores, DNS caches, routing tables
type ReadWriteCache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// NewReadWriteCache creates a new read-optimized cache
func NewReadWriteCache() *ReadWriteCache {
	return &ReadWriteCache{
		data: make(map[string]interface{}),
	}
}

// Set stores a value (write operation)
func (c *ReadWriteCache) Set(key string, value interface{}) {
	c.mu.Lock() // Exclusive lock for writes
	defer c.mu.Unlock()
	c.data[key] = value
}

// Get retrieves a value (read operation)
func (c *ReadWriteCache) Get(key string) (interface{}, bool) {
	c.mu.RLock() // Shared lock for reads
	defer c.mu.RUnlock()
	value, exists := c.data[key]
	return value, exists
}

// GetMultiple retrieves multiple values efficiently
func (c *ReadWriteCache) GetMultiple(keys []string) map[string]interface{} {
	c.mu.RLock() // Single read lock for multiple operations
	defer c.mu.RUnlock()

	result := make(map[string]interface{})
	for _, key := range keys {
		if value, exists := c.data[key]; exists {
			result[key] = value
		}
	}
	return result
}

// ConfigManager demonstrates advanced RWMutex usage for configuration
// Used in: Microservice configuration, feature flags, A/B testing
type ConfigManager struct {
	mu      sync.RWMutex
	config  map[string]string
	version int64
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		config:  make(map[string]string),
		version: 0,
	}
}

// UpdateConfig updates multiple configuration values atomically
func (cm *ConfigManager) UpdateConfig(newConfig map[string]string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Replace entire config atomically
	cm.config = make(map[string]string)
	for k, v := range newConfig {
		cm.config[k] = v
	}
	cm.version++
	fmt.Printf("📝 Config updated to version %d\n", cm.version)
}

// GetConfig returns a snapshot of current configuration
func (cm *ConfigManager) GetConfig() (map[string]string, int64) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Create a copy to avoid race conditions
	snapshot := make(map[string]string)
	for k, v := range cm.config {
		snapshot[k] = v
	}
	return snapshot, cm.version
}

// GetValue retrieves a specific configuration value
func (cm *ConfigManager) GetValue(key string) (string, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	value, exists := cm.config[key]
	return value, exists
}

// ==============================================================================
// 3. ATOMIC OPERATIONS - Lock-free Programming
// ==============================================================================

// AtomicCounter demonstrates atomic operations for high-performance counters
// Used in: Metrics collection, rate limiting, load balancing
type AtomicCounter struct {
	value int64
}

// Increment atomically increments the counter
func (c *AtomicCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Decrement atomically decrements the counter
func (c *AtomicCounter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Add atomically adds a value to the counter
func (c *AtomicCounter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

// Value atomically reads the counter value
func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// CompareAndSwap atomically compares and swaps the value
func (c *AtomicCounter) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.value, old, new)
}

// Reset atomically resets the counter to zero
func (c *AtomicCounter) Reset() int64 {
	return atomic.SwapInt64(&c.value, 0)
}

// LoadBalancer demonstrates atomic operations for request distribution
// Used in: Nginx, HAProxy-style load balancing, service discovery
type LoadBalancer struct {
	servers []string
	current int64
}

// NewLoadBalancer creates a new round-robin load balancer
func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
		current: 0,
	}
}

// NextServer returns the next server using atomic round-robin
func (lb *LoadBalancer) NextServer() string {
	if len(lb.servers) == 0 {
		return ""
	}

	// Atomic increment and wrap around
	next := atomic.AddInt64(&lb.current, 1)
	index := (next - 1) % int64(len(lb.servers))
	return lb.servers[index]
}

// GetStats returns current load balancer statistics
func (lb *LoadBalancer) GetStats() (int64, int) {
	current := atomic.LoadInt64(&lb.current)
	return current, len(lb.servers)
}

// ==============================================================================
// 4. ADVANCED SYNC PATTERNS - Complex Coordination
// ==============================================================================

// WorkerPool demonstrates advanced synchronization with multiple primitives
// Used in: HTTP servers, background job processors, data pipelines
type WorkerPool struct {
	mu         sync.RWMutex
	workers    int
	activeJobs int64 // atomic
	totalJobs  int64 // atomic
	isShutdown int32 // atomic boolean
	jobQueue   chan func()
	wg         sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers, queueSize int) *WorkerPool {
	wp := &WorkerPool{
		workers:  workers,
		jobQueue: make(chan func(), queueSize),
	}

	// Start workers
	for i := 0; i < workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	return wp
}

// worker processes jobs from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for job := range wp.jobQueue {
		if atomic.LoadInt32(&wp.isShutdown) == 1 {
			break
		}

		atomic.AddInt64(&wp.activeJobs, 1)
		atomic.AddInt64(&wp.totalJobs, 1)

		// Execute job
		job()

		atomic.AddInt64(&wp.activeJobs, -1)
	}

	fmt.Printf("🏁 Worker %d shutting down\n", id)
}

// Submit submits a job to the worker pool
func (wp *WorkerPool) Submit(job func()) error {
	if atomic.LoadInt32(&wp.isShutdown) == 1 {
		return fmt.Errorf("worker pool is shutdown")
	}

	select {
	case wp.jobQueue <- job:
		return nil
	default:
		return fmt.Errorf("job queue is full")
	}
}

// Stats returns current worker pool statistics
func (wp *WorkerPool) Stats() (active, total int64, workers int) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	return atomic.LoadInt64(&wp.activeJobs),
		atomic.LoadInt64(&wp.totalJobs),
		wp.workers
}

// Shutdown gracefully shuts down the worker pool
func (wp *WorkerPool) Shutdown(ctx context.Context) error {
	atomic.StoreInt32(&wp.isShutdown, 1)
	close(wp.jobQueue)

	// Wait for workers to finish with timeout
	done := make(chan struct{})
	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("✅ Worker pool shutdown completed")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout: %v", ctx.Err())
	}
}

// ==============================================================================
// 5. PERFORMANCE BENCHMARKS - Mutex vs RWMutex vs Atomic
// ==============================================================================

// BenchmarkManager compares performance of different synchronization methods
type BenchmarkManager struct {
	mutexCounter  *SafeCounter
	atomicCounter *AtomicCounter
	mutexCache    *CacheWithMutex
	rwCache       *ReadWriteCache
}

// NewBenchmarkManager creates a new benchmark manager
func NewBenchmarkManager() *BenchmarkManager {
	return &BenchmarkManager{
		mutexCounter:  &SafeCounter{},
		atomicCounter: &AtomicCounter{},
		mutexCache:    NewCache(),
		rwCache:       NewReadWriteCache(),
	}
}

// BenchmarkCounters compares mutex vs atomic counters
func (bm *BenchmarkManager) BenchmarkCounters(operations int) {
	fmt.Printf("\n🏃‍♂️ Benchmarking Counters (%d operations)\n", operations)

	// Benchmark mutex counter
	start := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bm.mutexCounter.Increment()
		}()
	}
	wg.Wait()
	mutexDuration := time.Since(start)

	// Benchmark atomic counter
	start = time.Now()
	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bm.atomicCounter.Increment()
		}()
	}
	wg.Wait()
	atomicDuration := time.Since(start)

	fmt.Printf("   Mutex Counter:  %v (Final: %d)\n", mutexDuration, bm.mutexCounter.Value())
	fmt.Printf("   Atomic Counter: %v (Final: %d)\n", atomicDuration, bm.atomicCounter.Value())
	fmt.Printf("   Atomic is %.2fx faster!\n", float64(mutexDuration.Nanoseconds())/float64(atomicDuration.Nanoseconds()))
}

// BenchmarkCaches compares mutex vs rwmutex for read-heavy workloads
func (bm *BenchmarkManager) BenchmarkCaches(reads, writes int) {
	fmt.Printf("\n📚 Benchmarking Caches (%d reads, %d writes)\n", reads, writes)

	// Populate caches
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		bm.mutexCache.Set(key, value)
		bm.rwCache.Set(key, value)
	}

	// Benchmark mutex cache
	start := time.Now()
	var wg sync.WaitGroup

	// Read operations
	for i := 0; i < reads; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i%100)
			bm.mutexCache.Get(key)
		}(i)
	}

	// Write operations
	for i := 0; i < writes; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("write-key-%d", i)
			bm.mutexCache.Set(key, fmt.Sprintf("write-value-%d", i))
		}(i)
	}

	wg.Wait()
	mutexDuration := time.Since(start)

	// Benchmark RWMutex cache
	start = time.Now()

	// Read operations
	for i := 0; i < reads; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i%100)
			bm.rwCache.Get(key)
		}(i)
	}

	// Write operations
	for i := 0; i < writes; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("write-key-%d", i)
			bm.rwCache.Set(key, fmt.Sprintf("write-value-%d", i))
		}(i)
	}

	wg.Wait()
	rwDuration := time.Since(start)

	fmt.Printf("   Mutex Cache:    %v\n", mutexDuration)
	fmt.Printf("   RWMutex Cache:  %v\n", rwDuration)
	fmt.Printf("   RWMutex is %.2fx faster for read-heavy workload!\n",
		float64(mutexDuration.Nanoseconds())/float64(rwDuration.Nanoseconds()))
}

func main() {
	fmt.Println("🔒 Synchronization Primitives Demo")
	fmt.Println("=====================================")

	// Demo 1: Basic Mutex Usage
	fmt.Println("\n1️⃣ Basic Mutex - Thread-Safe Counter")
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// Start multiple goroutines incrementing counter
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			counter.Increment()
			if id%20 == 0 {
				fmt.Printf("   Goroutine %d: Counter = %d\n", id, counter.Value())
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("✅ Final counter value: %d\n", counter.Value())

	// Demo 2: RWMutex for Read-Heavy Workloads
	fmt.Println("\n2️⃣ RWMutex - Configuration Manager")
	configMgr := NewConfigManager()

	// Initial configuration
	initialConfig := map[string]string{
		"database_url": "postgres://localhost:5432/app",
		"redis_url":    "redis://localhost:6379",
		"log_level":    "info",
		"debug_mode":   "false",
	}
	configMgr.UpdateConfig(initialConfig)

	// Multiple readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if value, exists := configMgr.GetValue("database_url"); exists {
				if id%10 == 0 {
					fmt.Printf("   Reader %d: database_url = %s\n", id, value)
				}
			}
		}(i)
	}

	// Occasional writer
	go func() {
		time.Sleep(100 * time.Millisecond)
		updatedConfig := map[string]string{
			"database_url": "postgres://prod-db:5432/app",
			"redis_url":    "redis://prod-redis:6379",
			"log_level":    "warn",
			"debug_mode":   "false",
		}
		configMgr.UpdateConfig(updatedConfig)
	}()

	wg.Wait()

	// Demo 3: Atomic Operations for High Performance
	fmt.Println("\n3️⃣ Atomic Operations - Load Balancer")
	servers := []string{"server-1", "server-2", "server-3", "server-4"}
	lb := NewLoadBalancer(servers)

	// Simulate concurrent requests
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			server := lb.NextServer()
			fmt.Printf("   Request %d routed to: %s\n", id+1, server)
		}(i)
	}
	wg.Wait()

	requests, serverCount := lb.GetStats()
	fmt.Printf("✅ Total requests: %d, Servers: %d\n", requests, serverCount)

	// Demo 4: Advanced Worker Pool
	fmt.Println("\n4️⃣ Advanced Worker Pool")
	pool := NewWorkerPool(3, 10)

	// Submit jobs
	for i := 0; i < 15; i++ {
		jobID := i + 1
		err := pool.Submit(func() {
			// Simulate work
			duration := time.Duration(rand.Intn(100)+50) * time.Millisecond
			time.Sleep(duration)
			fmt.Printf("   ✅ Job %d completed (took %v)\n", jobID, duration)
		})

		if err != nil {
			fmt.Printf("   ❌ Failed to submit job %d: %v\n", jobID, err)
		}
	}

	// Monitor progress
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(200 * time.Millisecond)
			active, total, workers := pool.Stats()
			fmt.Printf("   📊 Stats: %d active, %d total, %d workers\n", active, total, workers)
		}
	}()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("   ❌ Shutdown error: %v\n", err)
	}

	// Demo 5: Performance Benchmarks
	fmt.Println("\n5️⃣ Performance Benchmarks")
	benchMgr := NewBenchmarkManager()

	// Benchmark counters
	benchMgr.BenchmarkCounters(10000)

	// Benchmark caches (read-heavy workload)
	benchMgr.BenchmarkCaches(50000, 1000)

	fmt.Println("\n🎯 Synchronization Primitives Complete!")
	fmt.Printf("Key Insights:\n")
	fmt.Printf("• Mutex: Use for simple exclusive access to shared resources\n")
	fmt.Printf("• RWMutex: Use for read-heavy workloads (10:1 read/write ratio or higher)\n")
	fmt.Printf("• Atomic: Use for simple counters and flags (10-100x faster than mutex)\n")
	fmt.Printf("• Worker Pool: Combine all primitives for complex coordination\n")
	fmt.Printf("• Always measure: Profile your specific use case for optimal choice\n")

	fmt.Printf("\nRuntime Stats:\n")
	fmt.Printf("• Goroutines: %d\n", runtime.NumGoroutine())
	fmt.Printf("• CPUs: %d\n", runtime.NumCPU())
}
