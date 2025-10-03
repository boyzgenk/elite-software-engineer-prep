package production

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// ======================== BLUE-GREEN DEPLOYMENT ========================

// Environment represents a deployment environment
type Environment struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	IsActive  bool      `json:"is_active"`
	Health    string    `json:"health"` // healthy, degraded, unhealthy
	Instances []string  `json:"instances"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BlueGreenDeployment manages zero-downtime deployments
type BlueGreenDeployment struct {
	mu            sync.RWMutex
	blue          *Environment
	green         *Environment
	activeColor   string // "blue" or "green"
	healthChecker HealthChecker
	loadBalancer  LoadBalancer
}

type HealthChecker interface {
	CheckHealth(ctx context.Context, instances []string) error
}

type LoadBalancer interface {
	SwitchTraffic(ctx context.Context, targetEnv *Environment) error
	GetCurrentTraffic() string
}

// NewBlueGreenDeployment creates a new blue-green deployment manager
func NewBlueGreenDeployment(healthChecker HealthChecker, lb LoadBalancer) *BlueGreenDeployment {
	return &BlueGreenDeployment{
		blue: &Environment{
			Name:      "blue",
			Health:    "healthy",
			Instances: make([]string, 0),
			UpdatedAt: time.Now(),
		},
		green: &Environment{
			Name:      "green",
			Health:    "healthy",
			Instances: make([]string, 0),
			UpdatedAt: time.Now(),
		},
		activeColor:   "blue",
		healthChecker: healthChecker,
		loadBalancer:  lb,
	}
}

// Deploy performs a blue-green deployment
func (bg *BlueGreenDeployment) Deploy(ctx context.Context, newVersion string, instances []string) error {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	// Determine target environment (inactive one)
	var targetEnv, currentEnv *Environment
	if bg.activeColor == "blue" {
		targetEnv = bg.green
		currentEnv = bg.blue
	} else {
		targetEnv = bg.blue
		currentEnv = bg.green
	}

	// Update target environment
	targetEnv.Version = newVersion
	targetEnv.Instances = instances
	targetEnv.UpdatedAt = time.Now()
	targetEnv.Health = "deploying"

	// Wait for deployment to complete (simulate)
	time.Sleep(10 * time.Second)

	// Health check new environment
	if err := bg.healthChecker.CheckHealth(ctx, targetEnv.Instances); err != nil {
		targetEnv.Health = "unhealthy"
		return fmt.Errorf("health check failed: %w", err)
	}

	targetEnv.Health = "healthy"

	// Switch traffic to new environment
	if err := bg.loadBalancer.SwitchTraffic(ctx, targetEnv); err != nil {
		return fmt.Errorf("failed to switch traffic: %w", err)
	}

	// Update active environment
	targetEnv.IsActive = true
	currentEnv.IsActive = false
	bg.activeColor = targetEnv.Name

	return nil
}

// Rollback quickly switches back to the previous environment
func (bg *BlueGreenDeployment) Rollback(ctx context.Context) error {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	var rollbackEnv *Environment
	if bg.activeColor == "blue" {
		rollbackEnv = bg.green
	} else {
		rollbackEnv = bg.blue
	}

	// Ensure rollback environment is healthy
	if rollbackEnv.Health != "healthy" {
		return fmt.Errorf("rollback environment is not healthy: %s", rollbackEnv.Health)
	}

	// Switch traffic back
	if err := bg.loadBalancer.SwitchTraffic(ctx, rollbackEnv); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	// Update active status
	rollbackEnv.IsActive = true
	if bg.activeColor == "blue" {
		bg.blue.IsActive = false
		bg.activeColor = "green"
	} else {
		bg.green.IsActive = false
		bg.activeColor = "blue"
	}

	return nil
}

// GetStatus returns current deployment status
func (bg *BlueGreenDeployment) GetStatus() map[string]*Environment {
	bg.mu.RLock()
	defer bg.mu.RUnlock()

	return map[string]*Environment{
		"blue":   bg.blue,
		"green":  bg.green,
		"active": bg.getCurrentEnvironment(),
	}
}

func (bg *BlueGreenDeployment) getCurrentEnvironment() *Environment {
	if bg.activeColor == "blue" {
		return bg.blue
	}
	return bg.green
}

// ======================== CANARY DEPLOYMENT ========================

// CanaryDeployment manages gradual traffic shifting with automatic rollback
type CanaryDeployment struct {
	mu                sync.RWMutex
	production        *Environment
	canary            *Environment
	trafficSplit      int32 // percentage of traffic to canary (0-100)
	maxTrafficSplit   int32
	stepSize          int32
	stepDuration      time.Duration
	rollbackThreshold float64 // error rate threshold for automatic rollback
	metrics           *DeploymentMetrics
	isDeploying       bool
	stopCh            chan struct{}
}

type DeploymentMetrics struct {
	mu                   sync.RWMutex
	canaryErrorRate      float64
	productionErrorRate  float64
	canaryLatencyP99     time.Duration
	productionLatencyP99 time.Duration
	lastUpdated          time.Time
}

// NewCanaryDeployment creates a new canary deployment manager
func NewCanaryDeployment(stepSize int32, stepDuration time.Duration, rollbackThreshold float64) *CanaryDeployment {
	return &CanaryDeployment{
		production: &Environment{
			Name:      "production",
			Health:    "healthy",
			IsActive:  true,
			Instances: make([]string, 0),
			UpdatedAt: time.Now(),
		},
		canary: &Environment{
			Name:      "canary",
			Health:    "healthy",
			IsActive:  false,
			Instances: make([]string, 0),
			UpdatedAt: time.Now(),
		},
		maxTrafficSplit:   100,
		stepSize:          stepSize,
		stepDuration:      stepDuration,
		rollbackThreshold: rollbackThreshold,
		metrics:           &DeploymentMetrics{},
		stopCh:            make(chan struct{}),
	}
}

// Deploy starts a canary deployment with gradual traffic increase
func (cd *CanaryDeployment) Deploy(ctx context.Context, newVersion string, instances []string) error {
	cd.mu.Lock()
	if cd.isDeploying {
		cd.mu.Unlock()
		return fmt.Errorf("deployment already in progress")
	}
	cd.isDeploying = true
	cd.mu.Unlock()

	defer func() {
		cd.mu.Lock()
		cd.isDeploying = false
		cd.mu.Unlock()
	}()

	// Update canary environment
	cd.canary.Version = newVersion
	cd.canary.Instances = instances
	cd.canary.UpdatedAt = time.Now()
	cd.canary.Health = "healthy"
	cd.canary.IsActive = true

	// Start with 0% traffic to canary
	atomic.StoreInt32(&cd.trafficSplit, 0)

	// Gradual traffic increase
	for trafficPct := cd.stepSize; trafficPct <= cd.maxTrafficSplit; trafficPct += cd.stepSize {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-cd.stopCh:
			return fmt.Errorf("deployment stopped")
		default:
		}

		// Update traffic split
		atomic.StoreInt32(&cd.trafficSplit, trafficPct)

		// Wait for metrics to stabilize
		time.Sleep(cd.stepDuration)

		// Check metrics and auto-rollback if needed
		if err := cd.checkMetricsAndRollback(ctx); err != nil {
			return err
		}
	}

	// Deployment successful - promote canary to production
	return cd.promoteCanary(ctx)
}

// checkMetricsAndRollback monitors metrics and triggers rollback if thresholds are exceeded
func (cd *CanaryDeployment) checkMetricsAndRollback(ctx context.Context) error {
	cd.metrics.mu.RLock()
	canaryErrorRate := cd.metrics.canaryErrorRate
	productionErrorRate := cd.metrics.productionErrorRate
	cd.metrics.mu.RUnlock()

	// Auto-rollback if canary error rate exceeds threshold
	if canaryErrorRate > cd.rollbackThreshold {
		return cd.Rollback(ctx, fmt.Sprintf("canary error rate %.2f%% exceeds threshold %.2f%%",
			canaryErrorRate*100, cd.rollbackThreshold*100))
	}

	// Auto-rollback if canary performs significantly worse than production
	if canaryErrorRate > productionErrorRate*2 && canaryErrorRate > 0.01 {
		return cd.Rollback(ctx, fmt.Sprintf("canary error rate %.2f%% significantly higher than production %.2f%%",
			canaryErrorRate*100, productionErrorRate*100))
	}

	return nil
}

// Rollback immediately stops canary deployment and reverts to production
func (cd *CanaryDeployment) Rollback(ctx context.Context, reason string) error {
	atomic.StoreInt32(&cd.trafficSplit, 0)
	cd.canary.IsActive = false
	cd.canary.Health = "rolled_back"

	// Signal deployment to stop
	select {
	case cd.stopCh <- struct{}{}:
	default:
	}

	return fmt.Errorf("canary deployment rolled back: %s", reason)
}

// promoteCanary promotes canary to production
func (cd *CanaryDeployment) promoteCanary(ctx context.Context) error {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	// Swap environments
	cd.production, cd.canary = cd.canary, cd.production
	cd.production.Name = "production"
	cd.canary.Name = "canary"
	cd.canary.IsActive = false
	atomic.StoreInt32(&cd.trafficSplit, 0)

	return nil
}

// GetTrafficSplit returns current traffic distribution
func (cd *CanaryDeployment) GetTrafficSplit() (production, canary int32) {
	canaryPct := atomic.LoadInt32(&cd.trafficSplit)
	return 100 - canaryPct, canaryPct
}

// UpdateMetrics updates deployment metrics for decision making
func (cd *CanaryDeployment) UpdateMetrics(canaryErrorRate, productionErrorRate float64,
	canaryLatency, productionLatency time.Duration) {
	cd.metrics.mu.Lock()
	defer cd.metrics.mu.Unlock()

	cd.metrics.canaryErrorRate = canaryErrorRate
	cd.metrics.productionErrorRate = productionErrorRate
	cd.metrics.canaryLatencyP99 = canaryLatency
	cd.metrics.productionLatencyP99 = productionLatency
	cd.metrics.lastUpdated = time.Now()
}

// ======================== ROLLING DEPLOYMENT ========================

// RollingDeployment manages sequential instance updates with health checks
type RollingDeployment struct {
	mu              sync.RWMutex
	instances       []Instance
	maxUnavailable  int
	maxSurge        int
	healthChecker   HealthChecker
	deploymentState DeploymentState
}

type Instance struct {
	ID        string
	Version   string
	Status    InstanceStatus
	Health    string
	UpdatedAt time.Time
}

type InstanceStatus int

const (
	StatusRunning InstanceStatus = iota
	StatusUpdating
	StatusStopped
	StatusFailed
)

type DeploymentState struct {
	TotalInstances    int
	UpdatedInstances  int
	FailedInstances   int
	CurrentVersion    string
	TargetVersion     string
	StartTime         time.Time
	EstimatedComplete time.Time
}

// NewRollingDeployment creates a new rolling deployment manager
func NewRollingDeployment(instances []string, maxUnavailable, maxSurge int, healthChecker HealthChecker) *RollingDeployment {
	instanceList := make([]Instance, len(instances))
	for i, id := range instances {
		instanceList[i] = Instance{
			ID:        id,
			Version:   "v1.0.0",
			Status:    StatusRunning,
			Health:    "healthy",
			UpdatedAt: time.Now(),
		}
	}

	return &RollingDeployment{
		instances:      instanceList,
		maxUnavailable: maxUnavailable,
		maxSurge:       maxSurge,
		healthChecker:  healthChecker,
		deploymentState: DeploymentState{
			TotalInstances: len(instances),
			CurrentVersion: "v1.0.0",
		},
	}
}

// Deploy performs rolling deployment with configurable surge and unavailability
func (rd *RollingDeployment) Deploy(ctx context.Context, targetVersion string) error {
	rd.mu.Lock()
	defer rd.mu.Unlock()

	rd.deploymentState.TargetVersion = targetVersion
	rd.deploymentState.StartTime = time.Now()
	rd.deploymentState.UpdatedInstances = 0
	rd.deploymentState.FailedInstances = 0

	// Calculate batches based on maxUnavailable and maxSurge
	batchSize := rd.calculateBatchSize()

	for i := 0; i < len(rd.instances); i += batchSize {
		end := i + batchSize
		if end > len(rd.instances) {
			end = len(rd.instances)
		}

		batch := rd.instances[i:end]
		if err := rd.deployBatch(ctx, batch, targetVersion); err != nil {
			return fmt.Errorf("batch deployment failed: %w", err)
		}
	}

	rd.deploymentState.CurrentVersion = targetVersion
	rd.deploymentState.EstimatedComplete = time.Now()
	return nil
}

// calculateBatchSize determines optimal batch size based on constraints
func (rd *RollingDeployment) calculateBatchSize() int {
	totalInstances := len(rd.instances)

	// Ensure we don't exceed maxUnavailable constraint
	maxBatch := rd.maxUnavailable
	if maxBatch <= 0 {
		maxBatch = 1
	}

	// Consider maxSurge for additional instances
	if rd.maxSurge > 0 {
		maxBatch += rd.maxSurge
	}

	if maxBatch > totalInstances {
		maxBatch = totalInstances
	}

	return maxBatch
}

// deployBatch updates a batch of instances
func (rd *RollingDeployment) deployBatch(ctx context.Context, batch []Instance, targetVersion string) error {
	// Mark instances as updating
	for i := range batch {
		batch[i].Status = StatusUpdating
		batch[i].UpdatedAt = time.Now()
	}

	// Simulate deployment time
	time.Sleep(5 * time.Second)

	// Update instances to new version
	instanceIDs := make([]string, len(batch))
	for i := range batch {
		batch[i].Version = targetVersion
		batch[i].Status = StatusRunning
		batch[i].UpdatedAt = time.Now()
		instanceIDs[i] = batch[i].ID
		rd.deploymentState.UpdatedInstances++
	}

	// Health check batch
	if err := rd.healthChecker.CheckHealth(ctx, instanceIDs); err != nil {
		// Mark failed instances
		for i := range batch {
			batch[i].Status = StatusFailed
			batch[i].Health = "unhealthy"
			rd.deploymentState.FailedInstances++
		}
		return fmt.Errorf("health check failed for batch: %w", err)
	}

	// Mark instances as healthy
	for i := range batch {
		batch[i].Health = "healthy"
	}

	return nil
}

// GetDeploymentState returns current deployment progress
func (rd *RollingDeployment) GetDeploymentState() DeploymentState {
	rd.mu.RLock()
	defer rd.mu.RUnlock()
	return rd.deploymentState
}

// ======================== FEATURE FLAGS & CIRCUIT BREAKERS ========================

// FeatureFlag manages runtime feature toggles
type FeatureFlag struct {
	Name       string            `json:"name"`
	Enabled    bool              `json:"enabled"`
	Percentage int               `json:"percentage"` // 0-100
	Conditions map[string]string `json:"conditions"`
	UserGroups []string          `json:"user_groups"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// FeatureFlagManager manages feature flags and A/B testing
type FeatureFlagManager struct {
	mu    sync.RWMutex
	flags map[string]*FeatureFlag
	hash  func(string) uint32
}

// NewFeatureFlagManager creates a new feature flag manager
func NewFeatureFlagManager() *FeatureFlagManager {
	return &FeatureFlagManager{
		flags: make(map[string]*FeatureFlag),
		hash:  hashString,
	}
}

// CreateFlag creates a new feature flag
func (ffm *FeatureFlagManager) CreateFlag(name string, enabled bool, percentage int) {
	ffm.mu.Lock()
	defer ffm.mu.Unlock()

	ffm.flags[name] = &FeatureFlag{
		Name:       name,
		Enabled:    enabled,
		Percentage: percentage,
		Conditions: make(map[string]string),
		UserGroups: make([]string, 0),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// IsEnabled checks if a feature is enabled for a user
func (ffm *FeatureFlagManager) IsEnabled(flagName, userID string) bool {
	ffm.mu.RLock()
	flag, exists := ffm.flags[flagName]
	ffm.mu.RUnlock()

	if !exists || !flag.Enabled {
		return false
	}

	// Check percentage rollout
	if flag.Percentage < 100 {
		hash := ffm.hash(userID + flagName)
		if int(hash%100) >= flag.Percentage {
			return false
		}
	}

	return true
}

// UpdateFlag updates feature flag configuration
func (ffm *FeatureFlagManager) UpdateFlag(name string, enabled bool, percentage int) error {
	ffm.mu.Lock()
	defer ffm.mu.Unlock()

	flag, exists := ffm.flags[name]
	if !exists {
		return fmt.Errorf("feature flag %s not found", name)
	}

	flag.Enabled = enabled
	flag.Percentage = percentage
	flag.UpdatedAt = time.Now()

	return nil
}

// GetFlags returns all feature flags
func (ffm *FeatureFlagManager) GetFlags() map[string]*FeatureFlag {
	ffm.mu.RLock()
	defer ffm.mu.RUnlock()

	result := make(map[string]*FeatureFlag, len(ffm.flags))
	for k, v := range ffm.flags {
		flagCopy := *v
		result[k] = &flagCopy
	}
	return result
}

// CircuitBreakerState represents circuit breaker states
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker prevents cascading failures
type CircuitBreaker struct {
	mu              sync.Mutex
	name            string
	maxFailures     int64
	resetTimeout    time.Duration
	state           CircuitBreakerState
	failures        int64
	lastFailureTime time.Time
	nextRetryTime   time.Time
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, maxFailures int64, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:         name,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}
}

// Execute runs a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.canExecute() {
		return fmt.Errorf("circuit breaker %s is open", cb.name)
	}

	err := fn()
	cb.recordResult(err)
	return err
}

// canExecute checks if the circuit breaker allows execution
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if now.After(cb.nextRetryTime) {
			cb.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}

	return false
}

// recordResult updates circuit breaker state based on execution result
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailureTime = time.Now()

		if cb.state == StateHalfOpen || cb.failures >= cb.maxFailures {
			cb.state = StateOpen
			cb.nextRetryTime = time.Now().Add(cb.resetTimeout)
		}
	} else {
		cb.failures = 0
		cb.state = StateClosed
	}
}

// GetState returns current circuit breaker state
func (cb *CircuitBreaker) GetState() (CircuitBreakerState, int64, time.Time) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state, cb.failures, cb.lastFailureTime
}

// ======================== A/B TESTING INFRASTRUCTURE ========================

// Experiment represents an A/B test experiment
type Experiment struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Status       ExperimentStatus  `json:"status"`
	Variants     []Variant         `json:"variants"`
	TrafficSplit map[string]int    `json:"traffic_split"` // variant -> percentage
	Conditions   map[string]string `json:"conditions"`
	Metrics      []string          `json:"metrics"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      time.Time         `json:"end_time"`
	CreatedAt    time.Time         `json:"created_at"`
}

type ExperimentStatus int

const (
	ExperimentDraft ExperimentStatus = iota
	ExperimentRunning
	ExperimentPaused
	ExperimentCompleted
)

type Variant struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Config    map[string]interface{} `json:"config"`
	IsControl bool                   `json:"is_control"`
}

// ABTestingManager manages A/B testing experiments
type ABTestingManager struct {
	mu          sync.RWMutex
	experiments map[string]*Experiment
	assignments map[string]string // userID -> variantID
}

// NewABTestingManager creates a new A/B testing manager
func NewABTestingManager() *ABTestingManager {
	return &ABTestingManager{
		experiments: make(map[string]*Experiment),
		assignments: make(map[string]string),
	}
}

// CreateExperiment creates a new A/B test experiment
func (ab *ABTestingManager) CreateExperiment(id, name, description string, variants []Variant, trafficSplit map[string]int) error {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	// Validate traffic split totals 100%
	total := 0
	for _, pct := range trafficSplit {
		total += pct
	}
	if total != 100 {
		return fmt.Errorf("traffic split must total 100%%, got %d%%", total)
	}

	ab.experiments[id] = &Experiment{
		ID:           id,
		Name:         name,
		Description:  description,
		Status:       ExperimentDraft,
		Variants:     variants,
		TrafficSplit: trafficSplit,
		Conditions:   make(map[string]string),
		Metrics:      make([]string, 0),
		CreatedAt:    time.Now(),
	}

	return nil
}

// StartExperiment starts an A/B test experiment
func (ab *ABTestingManager) StartExperiment(experimentID string) error {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	experiment, exists := ab.experiments[experimentID]
	if !exists {
		return fmt.Errorf("experiment %s not found", experimentID)
	}

	if experiment.Status != ExperimentDraft {
		return fmt.Errorf("experiment %s is not in draft status", experimentID)
	}

	experiment.Status = ExperimentRunning
	experiment.StartTime = time.Now()

	return nil
}

// GetVariant returns the variant assignment for a user
func (ab *ABTestingManager) GetVariant(experimentID, userID string) (Variant, error) {
	ab.mu.RLock()
	experiment, exists := ab.experiments[experimentID]
	ab.mu.RUnlock()

	if !exists {
		return Variant{}, fmt.Errorf("experiment %s not found", experimentID)
	}

	if experiment.Status != ExperimentRunning {
		// Return control variant for non-running experiments
		for _, variant := range experiment.Variants {
			if variant.IsControl {
				return variant, nil
			}
		}
		return experiment.Variants[0], nil
	}

	// Check existing assignment
	assignmentKey := fmt.Sprintf("%s:%s", experimentID, userID)
	ab.mu.RLock()
	variantID, hasAssignment := ab.assignments[assignmentKey]
	ab.mu.RUnlock()

	if hasAssignment {
		// Return existing assignment
		for _, variant := range experiment.Variants {
			if variant.ID == variantID {
				return variant, nil
			}
		}
	}

	// Create new assignment based on traffic split
	hash := hashString(userID + experimentID)
	bucket := int(hash % 100)

	cumulative := 0
	for variantID, percentage := range experiment.TrafficSplit {
		cumulative += percentage
		if bucket < cumulative {
			// Store assignment
			ab.mu.Lock()
			ab.assignments[assignmentKey] = variantID
			ab.mu.Unlock()

			// Return variant
			for _, variant := range experiment.Variants {
				if variant.ID == variantID {
					return variant, nil
				}
			}
			break
		}
	}

	// Fallback to control variant
	for _, variant := range experiment.Variants {
		if variant.IsControl {
			return variant, nil
		}
	}

	return experiment.Variants[0], nil
}

// GetExperimentStatus returns experiment status and statistics
func (ab *ABTestingManager) GetExperimentStatus(experimentID string) (map[string]interface{}, error) {
	ab.mu.RLock()
	defer ab.mu.RUnlock()

	experiment, exists := ab.experiments[experimentID]
	if !exists {
		return nil, fmt.Errorf("experiment %s not found", experimentID)
	}

	// Count assignments per variant
	variantCounts := make(map[string]int)
	for assignmentKey, variantID := range ab.assignments {
		if assignmentKey[:len(experimentID)] == experimentID {
			variantCounts[variantID]++
		}
	}

	return map[string]interface{}{
		"experiment":     experiment,
		"variant_counts": variantCounts,
		"total_users":    len(variantCounts),
		"running_days":   time.Since(experiment.StartTime).Hours() / 24,
	}, nil
}

// ======================== UTILITY FUNCTIONS ========================

// hashString provides consistent hashing for user assignments
func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// SimpleHealthChecker implements basic health checking
type SimpleHealthChecker struct{}

func (shc *SimpleHealthChecker) CheckHealth(ctx context.Context, instances []string) error {
	// Simulate health check logic
	for _, instance := range instances {
		if rand.Float32() < 0.05 { // 5% chance of failure
			return fmt.Errorf("instance %s failed health check", instance)
		}
	}
	return nil
}

// SimpleLoadBalancer implements basic load balancer functionality
type SimpleLoadBalancer struct {
	currentTarget string
}

func (slb *SimpleLoadBalancer) SwitchTraffic(ctx context.Context, targetEnv *Environment) error {
	// Simulate traffic switching
	slb.currentTarget = targetEnv.Name
	time.Sleep(1 * time.Second)
	return nil
}

func (slb *SimpleLoadBalancer) GetCurrentTraffic() string {
	return slb.currentTarget
}

// ======================== HTTP HANDLERS FOR DEPLOYMENT MANAGEMENT ========================

// DeploymentHandler provides HTTP endpoints for deployment management
type DeploymentHandler struct {
	blueGreen *BlueGreenDeployment
	canary    *CanaryDeployment
	rolling   *RollingDeployment
	features  *FeatureFlagManager
	abtesting *ABTestingManager
}

// NewDeploymentHandler creates HTTP handlers for deployment operations
func NewDeploymentHandler() *DeploymentHandler {
	healthChecker := &SimpleHealthChecker{}
	loadBalancer := &SimpleLoadBalancer{}

	return &DeploymentHandler{
		blueGreen: NewBlueGreenDeployment(healthChecker, loadBalancer),
		canary:    NewCanaryDeployment(10, 2*time.Minute, 0.05),
		rolling:   NewRollingDeployment([]string{"i1", "i2", "i3", "i4"}, 2, 1, healthChecker),
		features:  NewFeatureFlagManager(),
		abtesting: NewABTestingManager(),
	}
}

// HandleBlueGreenDeploy handles blue-green deployment requests
func (dh *DeploymentHandler) HandleBlueGreenDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Version   string   `json:"version"`
		Instances []string `json:"instances"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := dh.blueGreen.Deploy(r.Context(), req.Version, req.Instances); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deployment started"})
}

// HandleFeatureFlag handles feature flag operations
func (dh *DeploymentHandler) HandleFeatureFlag(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		flags := dh.features.GetFlags()
		json.NewEncoder(w).Encode(flags)
	case http.MethodPost:
		var req struct {
			Name       string `json:"name"`
			Enabled    bool   `json:"enabled"`
			Percentage int    `json:"percentage"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dh.features.CreateFlag(req.Name, req.Enabled, req.Percentage)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "feature flag created"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
