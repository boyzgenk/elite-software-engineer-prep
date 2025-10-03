package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// 🏥 FAANG System Design - Health Check & Circuit Breaker Patterns
// Essential for production reliability and system design interviews

// HealthStatus represents the health state of a service
type HealthStatus int32

const (
	Healthy HealthStatus = iota
	Degraded
	Unhealthy
)

func (h HealthStatus) String() string {
	switch h {
	case Healthy:
		return "healthy"
	case Degraded:
		return "degraded"
	case Unhealthy:
		return "unhealthy"
	default:
		return "unknown"
	}
}

// HealthCheckResult contains the result of a health check
type HealthCheckResult struct {
	ServiceName string                 `json:"service_name"`
	Status      HealthStatus           `json:"status"`
	Timestamp   time.Time              `json:"timestamp"`
	Duration    time.Duration          `json:"duration_ms"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// ServiceHealthChecker defines the interface for health checking
type ServiceHealthChecker interface {
	Check(ctx context.Context) HealthCheckResult
	Name() string
}

// 🔍 Database Health Checker
type DatabaseHealthChecker struct {
	name     string
	dbURL    string
	timeout  time.Duration
	pingFunc func() error // Mock ping function for demo
}

func NewDatabaseHealthChecker(name, dbURL string, timeout time.Duration) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{
		name:    name,
		dbURL:   dbURL,
		timeout: timeout,
		// Mock ping function - in real world, this would ping the database
		pingFunc: func() error {
			// Simulate database latency
			time.Sleep(10 * time.Millisecond)
			return nil // Always healthy for demo
		},
	}
}

func (d *DatabaseHealthChecker) Name() string {
	return d.name
}

func (d *DatabaseHealthChecker) Check(ctx context.Context) HealthCheckResult {
	start := time.Now()
	result := HealthCheckResult{
		ServiceName: d.name,
		Timestamp:   start,
		Details:     make(map[string]interface{}),
	}

	// Create timeout context
	checkCtx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	// Channel to receive ping result
	done := make(chan error, 1)

	go func() {
		done <- d.pingFunc()
	}()

	select {
	case err := <-done:
		result.Duration = time.Since(start)
		if err != nil {
			result.Status = Unhealthy
			result.Error = err.Error()
		} else {
			// Check response time for degraded state
			if result.Duration > d.timeout/2 {
				result.Status = Degraded
				result.Details["warning"] = "slow response time"
			} else {
				result.Status = Healthy
			}
		}
		result.Details["database_url"] = d.dbURL
		result.Details["response_time_ms"] = result.Duration.Milliseconds()

	case <-checkCtx.Done():
		result.Duration = time.Since(start)
		result.Status = Unhealthy
		result.Error = "health check timeout"
		result.Details["timeout_ms"] = d.timeout.Milliseconds()
	}

	return result
}

// 🌐 HTTP Service Health Checker
type HTTPServiceHealthChecker struct {
	name     string
	endpoint string
	timeout  time.Duration
	client   *http.Client
	expected int // Expected HTTP status code
}

func NewHTTPServiceHealthChecker(name, endpoint string, timeout time.Duration) *HTTPServiceHealthChecker {
	return &HTTPServiceHealthChecker{
		name:     name,
		endpoint: endpoint,
		timeout:  timeout,
		expected: 200,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (h *HTTPServiceHealthChecker) Name() string {
	return h.name
}

func (h *HTTPServiceHealthChecker) Check(ctx context.Context) HealthCheckResult {
	start := time.Now()
	result := HealthCheckResult{
		ServiceName: h.name,
		Timestamp:   start,
		Details:     make(map[string]interface{}),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", h.endpoint, nil)
	if err != nil {
		result.Duration = time.Since(start)
		result.Status = Unhealthy
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		return result
	}

	resp, err := h.client.Do(req)
	result.Duration = time.Since(start)

	if err != nil {
		result.Status = Unhealthy
		result.Error = fmt.Sprintf("request failed: %v", err)
	} else {
		defer resp.Body.Close()
		result.Details["status_code"] = resp.StatusCode
		result.Details["response_time_ms"] = result.Duration.Milliseconds()

		if resp.StatusCode == h.expected {
			if result.Duration > h.timeout/2 {
				result.Status = Degraded
				result.Details["warning"] = "slow response time"
			} else {
				result.Status = Healthy
			}
		} else {
			result.Status = Unhealthy
			result.Error = fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
		}
	}

	return result
}

// 📊 Memory Health Checker
type MemoryHealthChecker struct {
	name          string
	maxMemoryMB   int64
	warningThresh float64 // Warning threshold as percentage (0.0-1.0)
}

func NewMemoryHealthChecker(name string, maxMemoryMB int64, warningThresh float64) *MemoryHealthChecker {
	return &MemoryHealthChecker{
		name:          name,
		maxMemoryMB:   maxMemoryMB,
		warningThresh: warningThresh,
	}
}

func (m *MemoryHealthChecker) Name() string {
	return m.name
}

func (m *MemoryHealthChecker) Check(ctx context.Context) HealthCheckResult {
	start := time.Now()
	result := HealthCheckResult{
		ServiceName: m.name,
		Timestamp:   start,
		Details:     make(map[string]interface{}),
	}

	// Simulate memory usage check (in real world, use runtime.MemStats)
	currentMemoryMB := int64(128) // Mock current memory usage
	memoryUsagePercent := float64(currentMemoryMB) / float64(m.maxMemoryMB)

	result.Duration = time.Since(start)
	result.Details["current_memory_mb"] = currentMemoryMB
	result.Details["max_memory_mb"] = m.maxMemoryMB
	result.Details["usage_percent"] = fmt.Sprintf("%.1f%%", memoryUsagePercent*100)

	if memoryUsagePercent >= 0.9 { // 90% threshold for unhealthy
		result.Status = Unhealthy
		result.Error = "memory usage too high"
	} else if memoryUsagePercent >= m.warningThresh {
		result.Status = Degraded
		result.Details["warning"] = "memory usage above warning threshold"
	} else {
		result.Status = Healthy
	}

	return result
}

// 🏗️ Composite Health Manager
// Aggregates multiple health checkers for overall service health
type HealthManager struct {
	checkers    []ServiceHealthChecker
	lastResults map[string]HealthCheckResult
	mutex       sync.RWMutex
	interval    time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewHealthManager(interval time.Duration) *HealthManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &HealthManager{
		checkers:    make([]ServiceHealthChecker, 0),
		lastResults: make(map[string]HealthCheckResult),
		interval:    interval,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (hm *HealthManager) AddChecker(checker ServiceHealthChecker) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	hm.checkers = append(hm.checkers, checker)
}

func (hm *HealthManager) StartChecking() {
	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	log.Printf("Health manager started with %v interval", hm.interval)

	// Initial health check
	hm.runHealthChecks()

	for {
		select {
		case <-hm.ctx.Done():
			log.Println("Health manager stopped")
			return
		case <-ticker.C:
			hm.runHealthChecks()
		}
	}
}

func (hm *HealthManager) Stop() {
	hm.cancel()
}

func (hm *HealthManager) runHealthChecks() {
	hm.mutex.Lock()
	checkers := make([]ServiceHealthChecker, len(hm.checkers))
	copy(checkers, hm.checkers)
	hm.mutex.Unlock()

	// Run all health checks concurrently
	results := make(chan HealthCheckResult, len(checkers))
	ctx, cancel := context.WithTimeout(hm.ctx, 30*time.Second)
	defer cancel()

	for _, checker := range checkers {
		go func(c ServiceHealthChecker) {
			result := c.Check(ctx)
			results <- result
		}(checker)
	}

	// Collect results
	newResults := make(map[string]HealthCheckResult)
	for i := 0; i < len(checkers); i++ {
		result := <-results
		newResults[result.ServiceName] = result

		// Log status changes
		hm.mutex.RLock()
		if lastResult, exists := hm.lastResults[result.ServiceName]; exists {
			if lastResult.Status != result.Status {
				log.Printf("Health status changed for %s: %s -> %s",
					result.ServiceName, lastResult.Status, result.Status)
			}
		}
		hm.mutex.RUnlock()
	}

	// Update stored results
	hm.mutex.Lock()
	hm.lastResults = newResults
	hm.mutex.Unlock()
}

func (hm *HealthManager) GetOverallHealth() HealthCheckResult {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	overall := HealthCheckResult{
		ServiceName: "overall",
		Timestamp:   time.Now(),
		Status:      Healthy,
		Details:     make(map[string]interface{}),
	}

	if len(hm.lastResults) == 0 {
		overall.Status = Unhealthy
		overall.Error = "no health checks available"
		return overall
	}

	healthyCount := 0
	degradedCount := 0
	unhealthyCount := 0
	services := make(map[string]string)

	for name, result := range hm.lastResults {
		services[name] = result.Status.String()
		switch result.Status {
		case Healthy:
			healthyCount++
		case Degraded:
			degradedCount++
		case Unhealthy:
			unhealthyCount++
		}
	}

	total := len(hm.lastResults)
	overall.Details["services"] = services
	overall.Details["healthy_count"] = healthyCount
	overall.Details["degraded_count"] = degradedCount
	overall.Details["unhealthy_count"] = unhealthyCount
	overall.Details["total_services"] = total

	// Determine overall status
	if unhealthyCount > 0 {
		overall.Status = Unhealthy
		overall.Error = fmt.Sprintf("%d/%d services unhealthy", unhealthyCount, total)
	} else if degradedCount > 0 {
		overall.Status = Degraded
		overall.Details["warning"] = fmt.Sprintf("%d/%d services degraded", degradedCount, total)
	} else {
		overall.Status = Healthy
	}

	return overall
}

func (hm *HealthManager) GetServiceHealth(serviceName string) (HealthCheckResult, bool) {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	result, exists := hm.lastResults[serviceName]
	return result, exists
}

// 🔄 Circuit Breaker Pattern
// Critical for preventing cascade failures in microservices
type CircuitBreakerState int32

const (
	CBClosed CircuitBreakerState = iota
	CBOpen
	CBHalfOpen
)

func (cb CircuitBreakerState) String() string {
	switch cb {
	case CBClosed:
		return "closed"
	case CBOpen:
		return "open"
	case CBHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

type CircuitBreakerStats struct {
	State               CircuitBreakerState `json:"state"`
	FailureCount        int64               `json:"failure_count"`
	SuccessCount        int64               `json:"success_count"`
	LastFailureTime     time.Time           `json:"last_failure_time"`
	LastStateChangeTime time.Time           `json:"last_state_change_time"`
	TotalRequests       int64               `json:"total_requests"`
}

type CircuitBreaker struct {
	name          string
	failureThresh int64         // Number of failures to open circuit
	timeout       time.Duration // Time to wait before trying half-open
	maxRetries    int           // Max retries in half-open state

	state           int32 // CircuitBreakerState (atomic)
	failureCount    int64 // Consecutive failure count (atomic)
	successCount    int64 // Success count in half-open (atomic)
	totalRequests   int64 // Total requests processed (atomic)
	lastFailure     int64 // Last failure timestamp (atomic)
	lastStateChange int64 // Last state change timestamp (atomic)
}

func NewCircuitBreaker(name string, failureThresh int64, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:            name,
		failureThresh:   failureThresh,
		timeout:         timeout,
		maxRetries:      3,
		state:           int32(CBClosed),
		lastStateChange: time.Now().UnixNano(),
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	// Check if circuit breaker allows the call
	if !cb.allowRequest() {
		return fmt.Errorf("circuit breaker %s is open", cb.name)
	}

	atomic.AddInt64(&cb.totalRequests, 1)

	// Execute the function
	err := fn()

	if err != nil {
		cb.onFailure()
		return err
	}

	cb.onSuccess()
	return nil
}

func (cb *CircuitBreaker) allowRequest() bool {
	state := CircuitBreakerState(atomic.LoadInt32(&cb.state))

	switch state {
	case CBClosed:
		return true
	case CBOpen:
		// Check if timeout has elapsed
		lastFailure := atomic.LoadInt64(&cb.lastFailure)
		if time.Now().UnixNano()-lastFailure > cb.timeout.Nanoseconds() {
			// Try to transition to half-open
			if atomic.CompareAndSwapInt32(&cb.state, int32(CBOpen), int32(CBHalfOpen)) {
				atomic.StoreInt64(&cb.lastStateChange, time.Now().UnixNano())
				atomic.StoreInt64(&cb.successCount, 0)
				log.Printf("Circuit breaker %s transitioning to half-open", cb.name)
				return true
			}
		}
		return false
	case CBHalfOpen:
		// Allow limited requests in half-open state
		successCount := atomic.LoadInt64(&cb.successCount)
		return successCount < int64(cb.maxRetries)
	default:
		return false
	}
}

func (cb *CircuitBreaker) onSuccess() {
	state := CircuitBreakerState(atomic.LoadInt32(&cb.state))

	if state == CBHalfOpen {
		successCount := atomic.AddInt64(&cb.successCount, 1)
		if successCount >= int64(cb.maxRetries) {
			// Enough successes in half-open, transition to closed
			if atomic.CompareAndSwapInt32(&cb.state, int32(CBHalfOpen), int32(CBClosed)) {
				atomic.StoreInt64(&cb.failureCount, 0)
				atomic.StoreInt64(&cb.lastStateChange, time.Now().UnixNano())
				log.Printf("Circuit breaker %s transitioning to closed", cb.name)
			}
		}
	} else if state == CBClosed {
		// Reset failure count on success
		atomic.StoreInt64(&cb.failureCount, 0)
	}
}

func (cb *CircuitBreaker) onFailure() {
	atomic.StoreInt64(&cb.lastFailure, time.Now().UnixNano())

	state := CircuitBreakerState(atomic.LoadInt32(&cb.state))
	failureCount := atomic.AddInt64(&cb.failureCount, 1)

	if state == CBClosed && failureCount >= cb.failureThresh {
		// Transition to open
		if atomic.CompareAndSwapInt32(&cb.state, int32(CBClosed), int32(CBOpen)) {
			atomic.StoreInt64(&cb.lastStateChange, time.Now().UnixNano())
			log.Printf("Circuit breaker %s transitioning to open after %d failures", cb.name, failureCount)
		}
	} else if state == CBHalfOpen {
		// Failure in half-open, go back to open
		if atomic.CompareAndSwapInt32(&cb.state, int32(CBHalfOpen), int32(CBOpen)) {
			atomic.StoreInt64(&cb.lastStateChange, time.Now().UnixNano())
			log.Printf("Circuit breaker %s back to open after failure in half-open", cb.name)
		}
	}
}

func (cb *CircuitBreaker) GetStats() CircuitBreakerStats {
	return CircuitBreakerStats{
		State:               CircuitBreakerState(atomic.LoadInt32(&cb.state)),
		FailureCount:        atomic.LoadInt64(&cb.failureCount),
		SuccessCount:        atomic.LoadInt64(&cb.successCount),
		LastFailureTime:     time.Unix(0, atomic.LoadInt64(&cb.lastFailure)),
		LastStateChangeTime: time.Unix(0, atomic.LoadInt64(&cb.lastStateChange)),
		TotalRequests:       atomic.LoadInt64(&cb.totalRequests),
	}
}

// 🧪 Demo Functions for FAANG Interview Practice
func demoServiceHealthChecking() {
	fmt.Println("🏥 FAANG System Design - Health Checking & Circuit Breaker Demo")
	fmt.Println("=" + fmt.Sprintf("%60s", "="))

	// Create health manager
	healthManager := NewHealthManager(3 * time.Second)

	// Add various health checkers
	dbChecker := NewDatabaseHealthChecker("postgres", "postgres://localhost:5432/app", 2*time.Second)
	redisChecker := NewHTTPServiceHealthChecker("redis", "http://localhost:6379/ping", 1*time.Second)
	memoryChecker := NewMemoryHealthChecker("memory", 1024, 0.8) // 1GB max, 80% warning

	healthManager.AddChecker(dbChecker)
	healthManager.AddChecker(redisChecker)
	healthManager.AddChecker(memoryChecker)

	// Start health checking in background
	go healthManager.StartChecking()

	// Wait for initial health checks
	time.Sleep(1 * time.Second)

	fmt.Println("\n📊 Health Check Results:")
	overall := healthManager.GetOverallHealth()
	prettyPrint("Overall Health", overall)

	// Check individual services
	if result, exists := healthManager.GetServiceHealth("postgres"); exists {
		prettyPrint("Database Health", result)
	}

	healthManager.Stop()
}

func demoCircuitBreaker() {
	fmt.Println("\n🔄 Circuit Breaker Pattern Demo:")
	fmt.Println("=" + fmt.Sprintf("%40s", "="))

	cb := NewCircuitBreaker("external-api", 3, 5*time.Second)

	// Simulate API calls with failures
	unstableAPI := func() error {
		// Simulate 70% failure rate
		if time.Now().UnixNano()%10 < 7 {
			return fmt.Errorf("external API error")
		}
		return nil
	}

	fmt.Println("Simulating API calls with 70% failure rate:")

	for i := 0; i < 10; i++ {
		err := cb.Call(unstableAPI)
		stats := cb.GetStats()

		if err != nil {
			fmt.Printf("Call %d: FAILED - %v (State: %s, Failures: %d)\n",
				i+1, err, stats.State, stats.FailureCount)
		} else {
			fmt.Printf("Call %d: SUCCESS (State: %s, Failures: %d)\n",
				i+1, stats.State, stats.FailureCount)
		}

		time.Sleep(100 * time.Millisecond)
	}

	// Show final stats
	finalStats := cb.GetStats()
	fmt.Printf("\nFinal Circuit Breaker Stats:\n")
	prettyPrint("Circuit Breaker Stats", finalStats)

	// Simulate recovery after timeout
	fmt.Println("\nWaiting for circuit breaker timeout...")
	time.Sleep(6 * time.Second)

	// Try calls again (should attempt half-open)
	stableAPI := func() error {
		return nil // Always succeeds
	}

	fmt.Println("Testing recovery with stable API:")
	for i := 0; i < 5; i++ {
		err := cb.Call(stableAPI)
		stats := cb.GetStats()

		if err != nil {
			fmt.Printf("Recovery call %d: FAILED - %v (State: %s)\n", i+1, err, stats.State)
		} else {
			fmt.Printf("Recovery call %d: SUCCESS (State: %s, Successes: %d)\n",
				i+1, stats.State, stats.SuccessCount)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func prettyPrint(title string, data interface{}) {
	fmt.Printf("\n%s:\n", title)
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}
	fmt.Println(string(jsonBytes))
}

func runHealthCheckDemo() {
	fmt.Println("🎯 FAANG System Design Interview - Health Checking Mastery")
	fmt.Println("Essential patterns for production reliability and L3/L4 interviews")
	fmt.Println()

	// Run demos
	demoServiceHealthChecking()
	demoCircuitBreaker()

	fmt.Println("\n🎪 Key Takeaways for FAANG Interviews:")
	fmt.Println("=" + fmt.Sprintf("%50s", "="))
	fmt.Println("1. Health Checks: Proactive monitoring of service dependencies")
	fmt.Println("2. Circuit Breaker: Prevents cascade failures in microservices")
	fmt.Println("3. Graceful Degradation: Service can operate with reduced functionality")
	fmt.Println("4. Observability: Comprehensive monitoring and alerting")
	fmt.Println("5. Atomic Operations: Thread-safe state management")
	fmt.Println("\n💪 Practice designing resilient systems like Netflix/Amazon!")
}

// Uncomment to run this demo independently
// func main() {
//     runHealthCheckDemo()
// }
