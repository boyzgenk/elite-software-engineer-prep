package distributed_systems

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Failure Detection Mechanisms

// FailureDetector interface for different failure detection strategies
type FailureDetector interface {
	IsAlive(nodeID string) bool
	MarkSuspected(nodeID string)
	MarkAlive(nodeID string)
	GetSuspectedNodes() []string
}

// PhiAccrualFailureDetector implements the Phi Accrual Failure Detector
// Used by Cassandra and Akka for adaptive failure detection
type PhiAccrualFailureDetector struct {
	threshold       float64
	maxSampleSize   int
	minStdDeviation float64
	windowSize      time.Duration

	// Per-node state
	heartbeats     map[string]*HeartbeatHistory
	suspectedNodes map[string]bool
	mu             sync.RWMutex
}

// HeartbeatHistory tracks heartbeat intervals for a node
type HeartbeatHistory struct {
	intervals   []float64
	lastArrival time.Time
	lastPhi     float64
}

// NewPhiAccrualFailureDetector creates a new Phi Accrual detector
func NewPhiAccrualFailureDetector(threshold float64) *PhiAccrualFailureDetector {
	return &PhiAccrualFailureDetector{
		threshold:       threshold,
		maxSampleSize:   200,
		minStdDeviation: 0.5,
		windowSize:      time.Second * 5,
		heartbeats:      make(map[string]*HeartbeatHistory),
		suspectedNodes:  make(map[string]bool),
	}
}

// Heartbeat records a heartbeat from a node
func (pfd *PhiAccrualFailureDetector) Heartbeat(nodeID string) {
	pfd.mu.Lock()
	defer pfd.mu.Unlock()

	now := time.Now()

	if history, exists := pfd.heartbeats[nodeID]; exists {
		// Calculate interval since last heartbeat
		interval := float64(now.Sub(history.lastArrival).Nanoseconds()) / 1e6 // milliseconds

		// Update interval history
		history.intervals = append(history.intervals, interval)
		if len(history.intervals) > pfd.maxSampleSize {
			history.intervals = history.intervals[1:]
		}

		history.lastArrival = now

		// Node is alive
		delete(pfd.suspectedNodes, nodeID)
	} else {
		// First heartbeat
		pfd.heartbeats[nodeID] = &HeartbeatHistory{
			intervals:   make([]float64, 0),
			lastArrival: now,
		}
	}
}

// IsAlive checks if a node is considered alive based on Phi value
func (pfd *PhiAccrualFailureDetector) IsAlive(nodeID string) bool {
	phi := pfd.calculatePhi(nodeID)
	return phi < pfd.threshold
}

// calculatePhi computes the Phi value for a node
func (pfd *PhiAccrualFailureDetector) calculatePhi(nodeID string) float64 {
	pfd.mu.RLock()
	defer pfd.mu.RUnlock()

	history, exists := pfd.heartbeats[nodeID]
	if !exists || len(history.intervals) < 2 {
		return 0.0 // Not enough data
	}

	now := time.Now()
	timeSinceLastHeartbeat := float64(now.Sub(history.lastArrival).Nanoseconds()) / 1e6

	// Calculate mean and standard deviation
	mean := pfd.calculateMean(history.intervals)
	stdDev := pfd.calculateStdDev(history.intervals, mean)
	if stdDev < pfd.minStdDeviation {
		stdDev = pfd.minStdDeviation
	}

	// Calculate Phi using normal distribution
	y := (timeSinceLastHeartbeat - mean) / stdDev
	phi := -math.Log10(pfd.cumulativeNormalDistribution(y))

	history.lastPhi = phi
	return phi
}

// MarkSuspected manually marks a node as suspected
func (pfd *PhiAccrualFailureDetector) MarkSuspected(nodeID string) {
	pfd.mu.Lock()
	defer pfd.mu.Unlock()
	pfd.suspectedNodes[nodeID] = true
}

// MarkAlive manually marks a node as alive
func (pfd *PhiAccrualFailureDetector) MarkAlive(nodeID string) {
	pfd.mu.Lock()
	defer pfd.mu.Unlock()
	delete(pfd.suspectedNodes, nodeID)
}

// GetSuspectedNodes returns list of suspected nodes
func (pfd *PhiAccrualFailureDetector) GetSuspectedNodes() []string {
	pfd.mu.RLock()
	defer pfd.mu.RUnlock()

	var suspected []string
	for nodeID := range pfd.suspectedNodes {
		suspected = append(suspected, nodeID)
	}

	// Also check nodes with high Phi values
	for nodeID := range pfd.heartbeats {
		if pfd.calculatePhi(nodeID) >= pfd.threshold {
			suspected = append(suspected, nodeID)
		}
	}

	return suspected
}

// Helper methods for Phi calculation
func (pfd *PhiAccrualFailureDetector) calculateMean(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (pfd *PhiAccrualFailureDetector) calculateStdDev(values []float64, mean float64) float64 {
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(values)))
}

func (pfd *PhiAccrualFailureDetector) cumulativeNormalDistribution(x float64) float64 {
	// Approximation of CDF for standard normal distribution
	return 0.5 * (1.0 + math.Erf(x/math.Sqrt(2)))
}

// Circuit Breaker Pattern for Failure Handling

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name         string
	maxFailures  int
	timeout      time.Duration
	resetTimeout time.Duration

	// State
	state           CircuitState
	failures        int
	lastFailureTime time.Time
	successCount    int

	// Metrics
	totalRequests  int64
	totalFailures  int64
	totalSuccesses int64

	mu            sync.RWMutex
	onStateChange func(from, to CircuitState)
}

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	Name         string
	MaxFailures  int
	Timeout      time.Duration
	ResetTimeout time.Duration
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		name:         config.Name,
		maxFailures:  config.MaxFailures,
		timeout:      config.Timeout,
		resetTimeout: config.ResetTimeout,
		state:        CircuitClosed,
	}
}

// Execute runs a function through the circuit breaker
func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	cb.mu.Lock()

	atomic.AddInt64(&cb.totalRequests, 1)

	// Check if circuit is open and should remain open
	if cb.state == CircuitOpen {
		if time.Since(cb.lastFailureTime) < cb.resetTimeout {
			cb.mu.Unlock()
			return nil, fmt.Errorf("circuit breaker %s is open", cb.name)
		}
		// Try to transition to half-open
		cb.setState(CircuitHalfOpen)
	}

	cb.mu.Unlock()

	// Execute the function
	result, err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure()
		return nil, err
	}

	cb.onSuccess()
	return result, nil
}

// onFailure handles failure cases
func (cb *CircuitBreaker) onFailure() {
	cb.failures++
	cb.lastFailureTime = time.Now()
	atomic.AddInt64(&cb.totalFailures, 1)

	switch cb.state {
	case CircuitClosed:
		if cb.failures >= cb.maxFailures {
			cb.setState(CircuitOpen)
		}
	case CircuitHalfOpen:
		cb.setState(CircuitOpen)
	}
}

// onSuccess handles success cases
func (cb *CircuitBreaker) onSuccess() {
	cb.successCount++
	atomic.AddInt64(&cb.totalSuccesses, 1)

	switch cb.state {
	case CircuitHalfOpen:
		if cb.successCount >= 3 { // Require multiple successes
			cb.setState(CircuitClosed)
			cb.failures = 0
			cb.successCount = 0
		}
	case CircuitClosed:
		if cb.failures > 0 {
			cb.failures = 0 // Reset failure count on success
		}
	}
}

// setState transitions the circuit breaker state
func (cb *CircuitBreaker) setState(newState CircuitState) {
	if cb.state != newState {
		oldState := cb.state
		cb.state = newState

		if cb.onStateChange != nil {
			cb.onStateChange(oldState, newState)
		}

		log.Printf("Circuit breaker %s state changed: %v -> %v", cb.name, oldState, newState)
	}
}

// GetState returns current state
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetMetrics returns circuit breaker metrics
func (cb *CircuitBreaker) GetMetrics() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"name":           cb.name,
		"state":          cb.state,
		"failures":       cb.failures,
		"success_count":  cb.successCount,
		"total_requests": atomic.LoadInt64(&cb.totalRequests),
		"total_failures": atomic.LoadInt64(&cb.totalFailures),
		"total_success":  atomic.LoadInt64(&cb.totalSuccesses),
	}
}

// Bulkhead Pattern Implementation

// Bulkhead isolates resources to prevent cascade failures
type Bulkhead struct {
	name        string
	pools       map[string]*ResourcePool
	defaultPool *ResourcePool
	mu          sync.RWMutex
}

// ResourcePool represents a pool of resources for a specific service
type ResourcePool struct {
	name          string
	maxSize       int
	currentSize   int
	activeCount   int64
	totalCount    int64
	rejectedCount int64
	semaphore     chan struct{}
	mu            sync.RWMutex
}

// NewBulkhead creates a new bulkhead with resource pools
func NewBulkhead(name string, defaultPoolSize int) *Bulkhead {
	defaultPool := &ResourcePool{
		name:      "default",
		maxSize:   defaultPoolSize,
		semaphore: make(chan struct{}, defaultPoolSize),
	}

	return &Bulkhead{
		name:        name,
		pools:       make(map[string]*ResourcePool),
		defaultPool: defaultPool,
	}
}

// AddPool adds a dedicated resource pool for a service
func (b *Bulkhead) AddPool(serviceName string, maxSize int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	pool := &ResourcePool{
		name:      serviceName,
		maxSize:   maxSize,
		semaphore: make(chan struct{}, maxSize),
	}

	b.pools[serviceName] = pool
}

// Execute runs a function with bulkhead isolation
func (b *Bulkhead) Execute(serviceName string, timeout time.Duration, fn func() error) error {
	pool := b.getPool(serviceName)

	// Try to acquire resource from pool
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case pool.semaphore <- struct{}{}:
		// Resource acquired
		atomic.AddInt64(&pool.totalCount, 1)
		atomic.AddInt64(&pool.activeCount, 1)

		defer func() {
			<-pool.semaphore // Release resource
			atomic.AddInt64(&pool.activeCount, -1)
		}()

		return fn()

	case <-ctx.Done():
		// Timeout or rejection
		atomic.AddInt64(&pool.rejectedCount, 1)
		return fmt.Errorf("bulkhead %s: resource pool %s exhausted", b.name, pool.name)
	}
}

// getPool returns the appropriate resource pool
func (b *Bulkhead) getPool(serviceName string) *ResourcePool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if pool, exists := b.pools[serviceName]; exists {
		return pool
	}
	return b.defaultPool
}

// GetPoolMetrics returns metrics for all pools
func (b *Bulkhead) GetPoolMetrics() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	metrics := make(map[string]interface{})

	// Default pool metrics
	metrics["default"] = map[string]interface{}{
		"max_size":       b.defaultPool.maxSize,
		"active_count":   atomic.LoadInt64(&b.defaultPool.activeCount),
		"total_count":    atomic.LoadInt64(&b.defaultPool.totalCount),
		"rejected_count": atomic.LoadInt64(&b.defaultPool.rejectedCount),
	}

	// Service-specific pool metrics
	for name, pool := range b.pools {
		metrics[name] = map[string]interface{}{
			"max_size":       pool.maxSize,
			"active_count":   atomic.LoadInt64(&pool.activeCount),
			"total_count":    atomic.LoadInt64(&pool.totalCount),
			"rejected_count": atomic.LoadInt64(&pool.rejectedCount),
		}
	}

	return metrics
}

// Retry Mechanisms with Exponential Backoff

// RetryPolicy defines retry behavior
type RetryPolicy struct {
	MaxAttempts    int
	BaseDelay      time.Duration
	MaxDelay       time.Duration
	BackoffFactor  float64
	Jitter         bool
	RetryableError func(error) bool
}

// ExponentialBackoffRetry implements exponential backoff retry
func ExponentialBackoffRetry(policy RetryPolicy, operation func() error) error {
	var lastErr error

	for attempt := 0; attempt < policy.MaxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if policy.RetryableError != nil && !policy.RetryableError(err) {
			return err // Non-retryable error
		}

		// Don't sleep after last attempt
		if attempt == policy.MaxAttempts-1 {
			break
		}

		// Calculate delay with exponential backoff
		delay := time.Duration(float64(policy.BaseDelay) * math.Pow(policy.BackoffFactor, float64(attempt)))
		if delay > policy.MaxDelay {
			delay = policy.MaxDelay
		}

		// Add jitter to prevent thundering herd
		if policy.Jitter {
			jitter := time.Duration(rand.Float64() * float64(delay) * 0.1)
			delay += jitter
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("operation failed after %d attempts: %w", policy.MaxAttempts, lastErr)
}

// Graceful Degradation Implementation

// DegradationLevel represents different levels of service degradation
type DegradationLevel int

const (
	FullService DegradationLevel = iota
	PartialService
	MinimalService
	EmergencyMode
)

// ServiceDegradationManager handles graceful degradation
type ServiceDegradationManager struct {
	currentLevel    DegradationLevel
	healthChecks    map[string]HealthCheck
	levelThresholds map[DegradationLevel]int

	// Feature flags for degradation
	features map[string]bool

	mu sync.RWMutex
}

// HealthCheck represents a health check function
type HealthCheck struct {
	Name     string
	Check    func() error
	Weight   int
	Critical bool
}

// NewServiceDegradationManager creates a new degradation manager
func NewServiceDegradationManager() *ServiceDegradationManager {
	return &ServiceDegradationManager{
		currentLevel: FullService,
		healthChecks: make(map[string]HealthCheck),
		levelThresholds: map[DegradationLevel]int{
			FullService:    90,
			PartialService: 70,
			MinimalService: 50,
			EmergencyMode:  0,
		},
		features: make(map[string]bool),
	}
}

// AddHealthCheck adds a health check
func (sdm *ServiceDegradationManager) AddHealthCheck(check HealthCheck) {
	sdm.mu.Lock()
	defer sdm.mu.Unlock()
	sdm.healthChecks[check.Name] = check
}

// RunHealthChecks executes all health checks and updates degradation level
func (sdm *ServiceDegradationManager) RunHealthChecks() {
	sdm.mu.Lock()
	defer sdm.mu.Unlock()

	totalWeight := 0
	healthyWeight := 0
	criticalFailures := 0

	for _, check := range sdm.healthChecks {
		totalWeight += check.Weight

		if err := check.Check(); err != nil {
			log.Printf("Health check %s failed: %v", check.Name, err)
			if check.Critical {
				criticalFailures++
			}
		} else {
			healthyWeight += check.Weight
		}
	}

	// Calculate health percentage
	healthPercentage := 0
	if totalWeight > 0 {
		healthPercentage = (healthyWeight * 100) / totalWeight
	}

	// Determine degradation level
	newLevel := sdm.calculateDegradationLevel(healthPercentage, criticalFailures > 0)

	if newLevel != sdm.currentLevel {
		log.Printf("Service degradation level changed: %v -> %v (health: %d%%)",
			sdm.currentLevel, newLevel, healthPercentage)
		sdm.currentLevel = newLevel
		sdm.updateFeatureFlags()
	}
}

// calculateDegradationLevel determines degradation level based on health
func (sdm *ServiceDegradationManager) calculateDegradationLevel(healthPercentage int, criticalFailure bool) DegradationLevel {
	if criticalFailure {
		return EmergencyMode
	}

	for level := FullService; level <= EmergencyMode; level++ {
		if healthPercentage >= sdm.levelThresholds[level] {
			return level
		}
	}

	return EmergencyMode
}

// updateFeatureFlags updates feature availability based on degradation level
func (sdm *ServiceDegradationManager) updateFeatureFlags() {
	switch sdm.currentLevel {
	case FullService:
		sdm.features["analytics"] = true
		sdm.features["recommendations"] = true
		sdm.features["search"] = true
		sdm.features["write"] = true

	case PartialService:
		sdm.features["analytics"] = false
		sdm.features["recommendations"] = true
		sdm.features["search"] = true
		sdm.features["write"] = true

	case MinimalService:
		sdm.features["analytics"] = false
		sdm.features["recommendations"] = false
		sdm.features["search"] = true
		sdm.features["write"] = true

	case EmergencyMode:
		sdm.features["analytics"] = false
		sdm.features["recommendations"] = false
		sdm.features["search"] = false
		sdm.features["write"] = false
	}
}

// IsFeatureEnabled checks if a feature is enabled at current degradation level
func (sdm *ServiceDegradationManager) IsFeatureEnabled(feature string) bool {
	sdm.mu.RLock()
	defer sdm.mu.RUnlock()
	return sdm.features[feature]
}

// GetCurrentLevel returns the current degradation level
func (sdm *ServiceDegradationManager) GetCurrentLevel() DegradationLevel {
	sdm.mu.RLock()
	defer sdm.mu.RUnlock()
	return sdm.currentLevel
}

// Chaos Engineering Implementation

// ChaosMonkey implements chaos engineering principles
type ChaosMonkey struct {
	enabled     bool
	experiments map[string]*ChaosExperiment
	mu          sync.RWMutex
}

// ChaosExperiment represents a chaos experiment
type ChaosExperiment struct {
	Name        string
	Probability float64
	Duration    time.Duration
	Action      func() error
	Enabled     bool
	LastRun     time.Time
	RunCount    int
}

// NewChaosMonkey creates a new chaos monkey
func NewChaosMonkey() *ChaosMonkey {
	return &ChaosMonkey{
		enabled:     false,
		experiments: make(map[string]*ChaosExperiment),
	}
}

// AddExperiment adds a chaos experiment
func (cm *ChaosMonkey) AddExperiment(experiment *ChaosExperiment) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.experiments[experiment.Name] = experiment
}

// Enable enables chaos experiments
func (cm *ChaosMonkey) Enable() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.enabled = true
}

// Disable disables chaos experiments
func (cm *ChaosMonkey) Disable() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.enabled = false
}

// RunExperiments runs enabled chaos experiments
func (cm *ChaosMonkey) RunExperiments() {
	cm.mu.RLock()
	enabled := cm.enabled
	experiments := make([]*ChaosExperiment, 0, len(cm.experiments))
	for _, exp := range cm.experiments {
		if exp.Enabled {
			experiments = append(experiments, exp)
		}
	}
	cm.mu.RUnlock()

	if !enabled {
		return
	}

	for _, exp := range experiments {
		// Check if experiment should run based on probability
		if rand.Float64() < exp.Probability {
			log.Printf("Running chaos experiment: %s", exp.Name)

			go func(e *ChaosExperiment) {
				if err := e.Action(); err != nil {
					log.Printf("Chaos experiment %s failed: %v", e.Name, err)
				}

				cm.mu.Lock()
				e.LastRun = time.Now()
				e.RunCount++
				cm.mu.Unlock()
			}(exp)
		}
	}
}

// GetExperimentStatus returns status of all experiments
func (cm *ChaosMonkey) GetExperimentStatus() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	status := make(map[string]interface{})
	status["enabled"] = cm.enabled

	experiments := make(map[string]interface{})
	for name, exp := range cm.experiments {
		experiments[name] = map[string]interface{}{
			"enabled":     exp.Enabled,
			"probability": exp.Probability,
			"last_run":    exp.LastRun,
			"run_count":   exp.RunCount,
		}
	}
	status["experiments"] = experiments

	return status
}

// Disaster Recovery Patterns

// DisasterRecoveryManager handles disaster recovery scenarios
type DisasterRecoveryManager struct {
	primaryRegion  string
	backupRegions  []string
	currentRegion  string
	failoverPolicy FailoverPolicy

	// State tracking
	regionHealth  map[string]bool
	lastFailover  time.Time
	failoverCount int

	mu sync.RWMutex
}

// FailoverPolicy defines failover behavior
type FailoverPolicy struct {
	AutoFailover          bool
	HealthCheckInterval   time.Duration
	FailoverCooldown      time.Duration
	MaxFailoversPerHour   int
	RequireManualApproval bool
}

// NewDisasterRecoveryManager creates a new DR manager
func NewDisasterRecoveryManager(primary string, backups []string, policy FailoverPolicy) *DisasterRecoveryManager {
	regionHealth := make(map[string]bool)
	regionHealth[primary] = true
	for _, region := range backups {
		regionHealth[region] = true
	}

	return &DisasterRecoveryManager{
		primaryRegion:  primary,
		backupRegions:  backups,
		currentRegion:  primary,
		failoverPolicy: policy,
		regionHealth:   regionHealth,
	}
}

// CheckRegionHealth performs health checks on all regions
func (drm *DisasterRecoveryManager) CheckRegionHealth() {
	drm.mu.Lock()
	defer drm.mu.Unlock()

	// Check current region health
	currentHealthy := drm.performHealthCheck(drm.currentRegion)
	drm.regionHealth[drm.currentRegion] = currentHealthy

	// If current region is unhealthy and auto-failover is enabled
	if !currentHealthy && drm.failoverPolicy.AutoFailover {
		if drm.canFailover() {
			bestRegion := drm.findBestRegion()
			if bestRegion != "" && bestRegion != drm.currentRegion {
				drm.performFailover(bestRegion)
			}
		}
	}

	// Check backup regions
	for _, region := range drm.backupRegions {
		healthy := drm.performHealthCheck(region)
		drm.regionHealth[region] = healthy
	}
}

// performHealthCheck checks if a region is healthy
func (drm *DisasterRecoveryManager) performHealthCheck(region string) bool {
	// Simulate health check - in real implementation, this would check:
	// - Network connectivity
	// - Service availability
	// - Database accessibility
	// - Resource utilization

	// For demo, randomly simulate failures
	return rand.Float64() > 0.1 // 90% healthy
}

// canFailover checks if failover is allowed based on policy
func (drm *DisasterRecoveryManager) canFailover() bool {
	now := time.Now()

	// Check cooldown period
	if now.Sub(drm.lastFailover) < drm.failoverPolicy.FailoverCooldown {
		return false
	}

	// Check failover rate limit
	hourAgo := now.Add(-time.Hour)
	if drm.lastFailover.After(hourAgo) && drm.failoverCount >= drm.failoverPolicy.MaxFailoversPerHour {
		return false
	}

	return true
}

// findBestRegion finds the best available region for failover
func (drm *DisasterRecoveryManager) findBestRegion() string {
	// Try primary region first (if not current)
	if drm.currentRegion != drm.primaryRegion && drm.regionHealth[drm.primaryRegion] {
		return drm.primaryRegion
	}

	// Try backup regions
	for _, region := range drm.backupRegions {
		if region != drm.currentRegion && drm.regionHealth[region] {
			return region
		}
	}

	return "" // No healthy region found
}

// performFailover executes failover to target region
func (drm *DisasterRecoveryManager) performFailover(targetRegion string) {
	log.Printf("Performing failover from %s to %s", drm.currentRegion, targetRegion)

	// Update state
	oldRegion := drm.currentRegion
	drm.currentRegion = targetRegion
	drm.lastFailover = time.Now()
	drm.failoverCount++

	// In real implementation, this would:
	// - Update DNS records
	// - Redirect traffic
	// - Update load balancer configuration
	// - Notify monitoring systems
	// - Update service discovery

	log.Printf("Failover completed: %s -> %s", oldRegion, targetRegion)
}

// ManualFailover performs manual failover to specified region
func (drm *DisasterRecoveryManager) ManualFailover(targetRegion string) error {
	drm.mu.Lock()
	defer drm.mu.Unlock()

	// Validate target region
	validRegion := targetRegion == drm.primaryRegion
	for _, region := range drm.backupRegions {
		if region == targetRegion {
			validRegion = true
			break
		}
	}

	if !validRegion {
		return fmt.Errorf("invalid target region: %s", targetRegion)
	}

	if !drm.regionHealth[targetRegion] {
		return fmt.Errorf("target region %s is unhealthy", targetRegion)
	}

	drm.performFailover(targetRegion)
	return nil
}

// GetDRStatus returns disaster recovery status
func (drm *DisasterRecoveryManager) GetDRStatus() map[string]interface{} {
	drm.mu.RLock()
	defer drm.mu.RUnlock()

	return map[string]interface{}{
		"primary_region": drm.primaryRegion,
		"current_region": drm.currentRegion,
		"backup_regions": drm.backupRegions,
		"region_health":  drm.regionHealth,
		"last_failover":  drm.lastFailover,
		"failover_count": drm.failoverCount,
		"auto_failover":  drm.failoverPolicy.AutoFailover,
	}
}

// Interview Discussion Points:
//
// 1. Failure Detection:
//    - Phi Accrual: Adaptive, accounts for network variability
//    - Heartbeat-based: Simple but can have false positives
//    - Gossip protocols: Distributed failure detection
//
// 2. Circuit Breaker Pattern:
//    - Prevents cascade failures
//    - States: Closed, Open, Half-Open
//    - Used by Netflix Hystrix, AWS Lambda
//
// 3. Bulkhead Pattern:
//    - Isolates resources to prevent total system failure
//    - Thread pools, connection pools, rate limiting
//    - Ship compartments analogy
//
// 4. Retry Strategies:
//    - Exponential backoff prevents thundering herd
//    - Jitter adds randomness
//    - Circuit breakers complement retries
//
// 5. Graceful Degradation:
//    - Feature flags for progressive functionality reduction
//    - Health-based service levels
//    - User experience preservation during failures
//
// 6. Chaos Engineering:
//    - Netflix Chaos Monkey, Chaos Gorilla
//    - Proactive failure testing
//    - Build confidence in system resilience
//
// 7. Disaster Recovery:
//    - RTO (Recovery Time Objective)
//    - RPO (Recovery Point Objective)
//    - Multi-region architectures
//    - Automated vs manual failover
//
// 8. Real-world Applications:
//    - AWS: Multi-AZ, Cross-region replication
//    - Google: Global load balancing, Spanner
//    - Netflix: Regional failover, Hystrix
//    - Uber: Multiple regions, graceful degradation
