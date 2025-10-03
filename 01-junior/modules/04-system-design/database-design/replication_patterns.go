package database_design

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Database Replication Patterns for System Design Interviews
// Interview Focus: High availability, read scaling, and data consistency

// ReplicationStrategy defines interface for different replication approaches
type ReplicationStrategy interface {
	WriteData(key string, data interface{}) error
	ReadData(key string) (interface{}, error)
	GetReplicationFactor() int
	HandleNodeFailure(nodeID string) error
	GetHealthStatus() ReplicationHealth
}

// MasterSlaveReplication implements single-writer, multiple-reader pattern
// FAANG Interview Point: Read scaling, simple consistency model
type MasterSlaveReplication struct {
	master   *DatabaseNode
	slaves   []*DatabaseNode
	mu       sync.RWMutex
	syncMode SyncMode
}

type SyncMode int

const (
	AsyncReplication SyncMode = iota
	SyncReplication
	SemiSyncReplication
)

type DatabaseNode struct {
	ID       string
	Address  string
	Data     map[string]interface{}
	IsAlive  bool
	LastSync time.Time
	mu       sync.RWMutex
}

func NewMasterSlaveReplication(masterID string, slaveIDs []string) *MasterSlaveReplication {
	master := &DatabaseNode{
		ID:      masterID,
		Address: fmt.Sprintf("master-%s:5432", masterID),
		Data:    make(map[string]interface{}),
		IsAlive: true,
	}

	var slaves []*DatabaseNode
	for _, slaveID := range slaveIDs {
		slave := &DatabaseNode{
			ID:      slaveID,
			Address: fmt.Sprintf("slave-%s:5432", slaveID),
			Data:    make(map[string]interface{}),
			IsAlive: true,
		}
		slaves = append(slaves, slave)
	}

	return &MasterSlaveReplication{
		master:   master,
		slaves:   slaves,
		syncMode: AsyncReplication,
	}
}

// WriteData writes to master and replicates to slaves
// FAANG Interview Point: Write path consistency and replication lag
func (msr *MasterSlaveReplication) WriteData(key string, data interface{}) error {
	msr.mu.Lock()
	defer msr.mu.Unlock()

	if !msr.master.IsAlive {
		return fmt.Errorf("master node is unavailable")
	}

	// Write to master first
	msr.master.mu.Lock()
	msr.master.Data[key] = data
	msr.master.mu.Unlock()

	// Replicate to slaves based on sync mode
	switch msr.syncMode {
	case SyncReplication:
		return msr.syncReplicateToSlaves(key, data)
	case SemiSyncReplication:
		return msr.semiSyncReplicateToSlaves(key, data)
	default:
		go msr.asyncReplicateToSlaves(key, data)
		return nil
	}
}

func (msr *MasterSlaveReplication) syncReplicateToSlaves(key string, data interface{}) error {
	// Synchronous replication - wait for all slaves
	var wg sync.WaitGroup
	errChan := make(chan error, len(msr.slaves))

	for _, slave := range msr.slaves {
		if !slave.IsAlive {
			continue
		}

		wg.Add(1)
		go func(s *DatabaseNode) {
			defer wg.Done()

			s.mu.Lock()
			defer s.mu.Unlock()

			// Simulate network delay
			time.Sleep(10 * time.Millisecond)
			s.Data[key] = data
			s.LastSync = time.Now()

			errChan <- nil
		}(slave)
	}

	wg.Wait()
	close(errChan)

	// Check for replication errors
	for err := range errChan {
		if err != nil {
			return fmt.Errorf("synchronous replication failed: %w", err)
		}
	}

	log.Printf("Synchronous replication completed for key: %s", key)
	return nil
}

func (msr *MasterSlaveReplication) semiSyncReplicateToSlaves(key string, data interface{}) error {
	// Semi-synchronous replication - wait for at least one slave
	if len(msr.slaves) == 0 {
		return nil
	}

	// Find first available slave
	for _, slave := range msr.slaves {
		if slave.IsAlive {
			slave.mu.Lock()
			slave.Data[key] = data
			slave.LastSync = time.Now()
			slave.mu.Unlock()

			// Replicate to remaining slaves asynchronously
			go msr.asyncReplicateToRemainingSlaves(key, data, slave.ID)
			return nil
		}
	}

	return fmt.Errorf("no available slaves for semi-sync replication")
}

func (msr *MasterSlaveReplication) asyncReplicateToSlaves(key string, data interface{}) {
	for _, slave := range msr.slaves {
		if !slave.IsAlive {
			continue
		}

		go func(s *DatabaseNode) {
			// Simulate replication lag
			time.Sleep(time.Duration(50+len(key)%100) * time.Millisecond)

			s.mu.Lock()
			defer s.mu.Unlock()

			s.Data[key] = data
			s.LastSync = time.Now()
		}(slave)
	}
}

func (msr *MasterSlaveReplication) asyncReplicateToRemainingSlaves(key string, data interface{}, skipSlaveID string) {
	for _, slave := range msr.slaves {
		if !slave.IsAlive || slave.ID == skipSlaveID {
			continue
		}

		go func(s *DatabaseNode) {
			time.Sleep(20 * time.Millisecond)

			s.mu.Lock()
			defer s.mu.Unlock()

			s.Data[key] = data
			s.LastSync = time.Now()
		}(slave)
	}
}

// ReadData reads from slaves for load distribution
// FAANG Interview Point: Read scaling and eventual consistency
func (msr *MasterSlaveReplication) ReadData(key string) (interface{}, error) {
	msr.mu.RLock()
	defer msr.mu.RUnlock()

	// Try reading from available slaves first (read scaling)
	for _, slave := range msr.slaves {
		if slave.IsAlive {
			slave.mu.RLock()
			if data, exists := slave.Data[key]; exists {
				slave.mu.RUnlock()
				return data, nil
			}
			slave.mu.RUnlock()
		}
	}

	// Fallback to master if slaves don't have the data
	if msr.master.IsAlive {
		msr.master.mu.RLock()
		defer msr.master.mu.RUnlock()

		if data, exists := msr.master.Data[key]; exists {
			return data, nil
		}
	}

	return nil, fmt.Errorf("key not found: %s", key)
}

func (msr *MasterSlaveReplication) GetReplicationFactor() int {
	return len(msr.slaves) + 1 // slaves + master
}

func (msr *MasterSlaveReplication) HandleNodeFailure(nodeID string) error {
	msr.mu.Lock()
	defer msr.mu.Unlock()

	if msr.master.ID == nodeID {
		msr.master.IsAlive = false
		log.Printf("Master node %s failed - manual failover required", nodeID)
		return fmt.Errorf("master node failure requires manual intervention")
	}

	for _, slave := range msr.slaves {
		if slave.ID == nodeID {
			slave.IsAlive = false
			log.Printf("Slave node %s marked as failed", nodeID)
			return nil
		}
	}

	return fmt.Errorf("node not found: %s", nodeID)
}

func (msr *MasterSlaveReplication) GetHealthStatus() ReplicationHealth {
	msr.mu.RLock()
	defer msr.mu.RUnlock()

	aliveSlaves := 0
	var maxLag time.Duration

	for _, slave := range msr.slaves {
		if slave.IsAlive {
			aliveSlaves++
			lag := time.Since(slave.LastSync)
			if lag > maxLag {
				maxLag = lag
			}
		}
	}

	return ReplicationHealth{
		MasterAlive:       msr.master.IsAlive,
		AliveSlavesCount:  aliveSlaves,
		TotalSlavesCount:  len(msr.slaves),
		MaxReplicationLag: maxLag,
		HealthScore:       float64(aliveSlaves) / float64(len(msr.slaves)),
	}
}

// MasterMasterReplication implements multi-writer pattern
// FAANG Interview Point: High availability, conflict resolution complexity
type MasterMasterReplication struct {
	masters          []*DatabaseNode
	mu               sync.RWMutex
	conflictResolver ConflictResolver
}

type ConflictResolver interface {
	ResolveConflict(key string, values []ConflictingValue) (interface{}, error)
}

type ConflictingValue struct {
	Value     interface{}
	Timestamp time.Time
	NodeID    string
}

func NewMasterMasterReplication(masterIDs []string) *MasterMasterReplication {
	var masters []*DatabaseNode
	for _, masterID := range masterIDs {
		master := &DatabaseNode{
			ID:      masterID,
			Address: fmt.Sprintf("master-%s:5432", masterID),
			Data:    make(map[string]interface{}),
			IsAlive: true,
		}
		masters = append(masters, master)
	}

	return &MasterMasterReplication{
		masters:          masters,
		conflictResolver: &LastWriteWinsResolver{},
	}
}

// WriteData writes to any available master
// FAANG Interview Point: Write availability vs consistency trade-offs
func (mmr *MasterMasterReplication) WriteData(key string, data interface{}) error {
	mmr.mu.RLock()
	defer mmr.mu.RUnlock()

	// Find first available master
	for _, master := range mmr.masters {
		if master.IsAlive {
			master.mu.Lock()
			master.Data[key] = TimestampedValue{
				Value:     data,
				Timestamp: time.Now(),
				NodeID:    master.ID,
			}
			master.mu.Unlock()

			// Asynchronously replicate to other masters
			go mmr.replicateToOtherMasters(key, data, master.ID)
			return nil
		}
	}

	return fmt.Errorf("no available masters")
}

type TimestampedValue struct {
	Value     interface{}
	Timestamp time.Time
	NodeID    string
}

func (mmr *MasterMasterReplication) replicateToOtherMasters(key string, data interface{}, sourceID string) {
	for _, master := range mmr.masters {
		if master.ID == sourceID || !master.IsAlive {
			continue
		}

		go func(m *DatabaseNode) {
			// Simulate network delay
			time.Sleep(30 * time.Millisecond)

			m.mu.Lock()
			defer m.mu.Unlock()

			// Check for conflicts
			if existingData, exists := m.Data[key]; exists {
				if tv, ok := existingData.(TimestampedValue); ok {
					// Conflict detected - use resolver
					conflictingValues := []ConflictingValue{
						{Value: tv.Value, Timestamp: tv.Timestamp, NodeID: tv.NodeID},
						{Value: data, Timestamp: time.Now(), NodeID: sourceID},
					}

					resolvedValue, err := mmr.conflictResolver.ResolveConflict(key, conflictingValues)
					if err == nil {
						m.Data[key] = TimestampedValue{
							Value:     resolvedValue,
							Timestamp: time.Now(),
							NodeID:    m.ID,
						}
					}
					return
				}
			}

			// No conflict - just replicate
			m.Data[key] = TimestampedValue{
				Value:     data,
				Timestamp: time.Now(),
				NodeID:    sourceID,
			}
		}(master)
	}
}

// ReadData reads from any available master
func (mmr *MasterMasterReplication) ReadData(key string) (interface{}, error) {
	mmr.mu.RLock()
	defer mmr.mu.RUnlock()

	for _, master := range mmr.masters {
		if master.IsAlive {
			master.mu.RLock()
			if data, exists := master.Data[key]; exists {
				if tv, ok := data.(TimestampedValue); ok {
					master.mu.RUnlock()
					return tv.Value, nil
				}
			}
			master.mu.RUnlock()
		}
	}

	return nil, fmt.Errorf("key not found: %s", key)
}

func (mmr *MasterMasterReplication) GetReplicationFactor() int {
	return len(mmr.masters)
}

func (mmr *MasterMasterReplication) HandleNodeFailure(nodeID string) error {
	mmr.mu.Lock()
	defer mmr.mu.Unlock()

	for _, master := range mmr.masters {
		if master.ID == nodeID {
			master.IsAlive = false
			log.Printf("Master node %s failed - other masters continue operating", nodeID)
			return nil
		}
	}

	return fmt.Errorf("node not found: %s", nodeID)
}

func (mmr *MasterMasterReplication) GetHealthStatus() ReplicationHealth {
	mmr.mu.RLock()
	defer mmr.mu.RUnlock()

	aliveMasters := 0
	for _, master := range mmr.masters {
		if master.IsAlive {
			aliveMasters++
		}
	}

	return ReplicationHealth{
		MasterAlive:       aliveMasters > 0,
		AliveSlavesCount:  aliveMasters - 1, // Other masters act as "slaves"
		TotalSlavesCount:  len(mmr.masters) - 1,
		MaxReplicationLag: 50 * time.Millisecond, // Async replication
		HealthScore:       float64(aliveMasters) / float64(len(mmr.masters)),
	}
}

// Conflict Resolution Strategies
// Interview Focus: Handling concurrent writes in distributed systems

// LastWriteWinsResolver resolves conflicts using timestamp
type LastWriteWinsResolver struct{}

func (lww *LastWriteWinsResolver) ResolveConflict(key string, values []ConflictingValue) (interface{}, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("no values to resolve")
	}

	latest := values[0]
	for _, value := range values[1:] {
		if value.Timestamp.After(latest.Timestamp) {
			latest = value
		}
	}

	log.Printf("Conflict resolved for key %s using last-write-wins: node %s wins", key, latest.NodeID)
	return latest.Value, nil
}

// ApplicationSpecificResolver allows custom conflict resolution
type ApplicationSpecificResolver struct {
	ResolverFunc func(key string, values []ConflictingValue) (interface{}, error)
}

func (asr *ApplicationSpecificResolver) ResolveConflict(key string, values []ConflictingValue) (interface{}, error) {
	return asr.ResolverFunc(key, values)
}

// Read Replicas for Geographic Distribution
// Interview Focus: Global read scaling and latency optimization

type ReadReplicaManager struct {
	primary  *DatabaseNode
	replicas map[string]*GeographicReplica
	mu       sync.RWMutex
}

type GeographicReplica struct {
	*DatabaseNode
	Region       string
	Latency      time.Duration
	ReadCount    int64
	LastReadTime time.Time
}

func NewReadReplicaManager(primaryID string) *ReadReplicaManager {
	primary := &DatabaseNode{
		ID:      primaryID,
		Address: fmt.Sprintf("primary-%s:5432", primaryID),
		Data:    make(map[string]interface{}),
		IsAlive: true,
	}

	return &ReadReplicaManager{
		primary:  primary,
		replicas: make(map[string]*GeographicReplica),
	}
}

// AddReadReplica adds a geographic read replica
func (rrm *ReadReplicaManager) AddReadReplica(replicaID, region string, latency time.Duration) error {
	rrm.mu.Lock()
	defer rrm.mu.Unlock()

	replica := &GeographicReplica{
		DatabaseNode: &DatabaseNode{
			ID:      replicaID,
			Address: fmt.Sprintf("replica-%s:5432", replicaID),
			Data:    make(map[string]interface{}),
			IsAlive: true,
		},
		Region:  region,
		Latency: latency,
	}

	rrm.replicas[replicaID] = replica
	log.Printf("Added read replica %s in region %s with %v latency", replicaID, region, latency)
	return nil
}

// WriteData writes to primary and replicates asynchronously
func (rrm *ReadReplicaManager) WriteData(key string, data interface{}) error {
	if !rrm.primary.IsAlive {
		return fmt.Errorf("primary node unavailable")
	}

	// Write to primary
	rrm.primary.mu.Lock()
	rrm.primary.Data[key] = data
	rrm.primary.mu.Unlock()

	// Replicate to all read replicas asynchronously
	go rrm.replicateToReadReplicas(key, data)
	return nil
}

func (rrm *ReadReplicaManager) replicateToReadReplicas(key string, data interface{}) {
	rrm.mu.RLock()
	defer rrm.mu.RUnlock()

	for _, replica := range rrm.replicas {
		if !replica.IsAlive {
			continue
		}

		go func(r *GeographicReplica) {
			// Simulate replication delay based on geographic distance
			time.Sleep(r.Latency)

			r.mu.Lock()
			defer r.mu.Unlock()

			r.Data[key] = data
			r.LastSync = time.Now()
		}(replica)
	}
}

// ReadDataFromClosestReplica reads from the replica with lowest latency
// FAANG Interview Point: Geographic load balancing and latency optimization
func (rrm *ReadReplicaManager) ReadDataFromClosestReplica(key string, userRegion string) (interface{}, error) {
	rrm.mu.RLock()
	defer rrm.mu.RUnlock()

	// Find replicas in the same region first
	for _, replica := range rrm.replicas {
		if replica.Region == userRegion && replica.IsAlive {
			replica.mu.RLock()
			if data, exists := replica.Data[key]; exists {
				replica.ReadCount++
				replica.LastReadTime = time.Now()
				replica.mu.RUnlock()
				return data, nil
			}
			replica.mu.RUnlock()
		}
	}

	// Find closest replica by latency
	var closestReplica *GeographicReplica
	minLatency := time.Hour

	for _, replica := range rrm.replicas {
		if replica.IsAlive && replica.Latency < minLatency {
			replica.mu.RLock()
			if _, exists := replica.Data[key]; exists {
				closestReplica = replica
				minLatency = replica.Latency
			}
			replica.mu.RUnlock()
		}
	}

	if closestReplica != nil {
		closestReplica.mu.RLock()
		defer closestReplica.mu.RUnlock()

		if data, exists := closestReplica.Data[key]; exists {
			closestReplica.ReadCount++
			closestReplica.LastReadTime = time.Now()
			return data, nil
		}
	}

	// Fallback to primary
	rrm.primary.mu.RLock()
	defer rrm.primary.mu.RUnlock()

	if data, exists := rrm.primary.Data[key]; exists {
		return data, nil
	}

	return nil, fmt.Errorf("key not found: %s", key)
}

// Lag Monitoring and Management
// Interview Focus: Operational visibility and automated remediation

type ReplicationLagMonitor struct {
	nodes          map[string]*DatabaseNode
	lagThresholds  map[string]time.Duration
	alertCallbacks []AlertCallback
	mu             sync.RWMutex
}

type AlertCallback func(alert LagAlert)

type LagAlert struct {
	NodeID    string
	LagAmount time.Duration
	Severity  AlertSeverity
	Timestamp time.Time
	Message   string
}

type AlertSeverity int

const (
	AlertInfo AlertSeverity = iota
	AlertWarning
	AlertCritical
)

func NewReplicationLagMonitor() *ReplicationLagMonitor {
	return &ReplicationLagMonitor{
		nodes: make(map[string]*DatabaseNode),
		lagThresholds: map[string]time.Duration{
			"warning":  5 * time.Second,
			"critical": 30 * time.Second,
		},
	}
}

func (rlm *ReplicationLagMonitor) AddNode(node *DatabaseNode) {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()

	rlm.nodes[node.ID] = node
}

func (rlm *ReplicationLagMonitor) AddAlertCallback(callback AlertCallback) {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()

	rlm.alertCallbacks = append(rlm.alertCallbacks, callback)
}

// MonitorLag continuously monitors replication lag
func (rlm *ReplicationLagMonitor) MonitorLag(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rlm.checkReplicationLag()
		}
	}
}

func (rlm *ReplicationLagMonitor) checkReplicationLag() {
	rlm.mu.RLock()
	defer rlm.mu.RUnlock()

	for nodeID, node := range rlm.nodes {
		if !node.IsAlive {
			continue
		}

		lag := time.Since(node.LastSync)

		var severity AlertSeverity
		var message string

		if lag > rlm.lagThresholds["critical"] {
			severity = AlertCritical
			message = fmt.Sprintf("Critical replication lag detected: %v", lag)
		} else if lag > rlm.lagThresholds["warning"] {
			severity = AlertWarning
			message = fmt.Sprintf("Warning: replication lag increasing: %v", lag)
		} else {
			continue // No alert needed
		}

		alert := LagAlert{
			NodeID:    nodeID,
			LagAmount: lag,
			Severity:  severity,
			Timestamp: time.Now(),
			Message:   message,
		}

		// Send alerts to all registered callbacks
		for _, callback := range rlm.alertCallbacks {
			go callback(alert)
		}
	}
}

// Failover Mechanisms
// Interview Focus: Automated recovery and split-brain prevention

type FailoverManager struct {
	strategy           ReplicationStrategy
	healthChecker      *HealthChecker
	failoverInProgress bool
	mu                 sync.RWMutex
}

type HealthChecker struct {
	nodes           map[string]*DatabaseNode
	checkInterval   time.Duration
	timeoutDuration time.Duration
}

func NewFailoverManager(strategy ReplicationStrategy) *FailoverManager {
	return &FailoverManager{
		strategy: strategy,
		healthChecker: &HealthChecker{
			nodes:           make(map[string]*DatabaseNode),
			checkInterval:   1 * time.Second,
			timeoutDuration: 3 * time.Second,
		},
	}
}

// AutomaticFailover performs automatic failover when master fails
// FAANG Interview Point: Zero-downtime failover strategies
func (fm *FailoverManager) AutomaticFailover(ctx context.Context) error {
	fm.mu.Lock()
	if fm.failoverInProgress {
		fm.mu.Unlock()
		return fmt.Errorf("failover already in progress")
	}
	fm.failoverInProgress = true
	fm.mu.Unlock()

	defer func() {
		fm.mu.Lock()
		fm.failoverInProgress = false
		fm.mu.Unlock()
	}()

	log.Println("Starting automatic failover process...")

	// Step 1: Detect failure
	health := fm.strategy.GetHealthStatus()
	if health.MasterAlive {
		return fmt.Errorf("master is alive, no failover needed")
	}

	// Step 2: Choose new master (implementation depends on strategy)
	if msr, ok := fm.strategy.(*MasterSlaveReplication); ok {
		return fm.promoteSlaveToMaster(msr)
	}

	log.Println("Failover completed successfully")
	return nil
}

func (fm *FailoverManager) promoteSlaveToMaster(msr *MasterSlaveReplication) error {
	// Find the most up-to-date slave
	var bestSlave *DatabaseNode
	var latestSync time.Time

	for _, slave := range msr.slaves {
		if slave.IsAlive && slave.LastSync.After(latestSync) {
			bestSlave = slave
			latestSync = slave.LastSync
		}
	}

	if bestSlave == nil {
		return fmt.Errorf("no suitable slave found for promotion")
	}

	log.Printf("Promoting slave %s to master", bestSlave.ID)

	// Promote slave to master
	msr.master = bestSlave

	// Remove promoted slave from slaves list
	for i, slave := range msr.slaves {
		if slave.ID == bestSlave.ID {
			msr.slaves = append(msr.slaves[:i], msr.slaves[i+1:]...)
			break
		}
	}

	log.Printf("Slave %s successfully promoted to master", bestSlave.ID)
	return nil
}

// Health status tracking
type ReplicationHealth struct {
	MasterAlive       bool
	AliveSlavesCount  int
	TotalSlavesCount  int
	MaxReplicationLag time.Duration
	HealthScore       float64 // 0.0 to 1.0
}

// ReplicationMetricsCollector provides operational insights
// FAANG Interview Point: Monitoring and observability
type ReplicationMetricsCollector struct {
	strategy   ReplicationStrategy
	lagMonitor *ReplicationLagMonitor
	metrics    ReplicationMetrics
	mu         sync.RWMutex
}

type ReplicationMetrics struct {
	TotalWrites           int64
	TotalReads            int64
	AverageReplicationLag time.Duration
	FailoverCount         int64
	LastFailoverTime      time.Time
	HealthScore           float64
}

func NewReplicationMetricsCollector(strategy ReplicationStrategy, lagMonitor *ReplicationLagMonitor) *ReplicationMetricsCollector {
	return &ReplicationMetricsCollector{
		strategy:   strategy,
		lagMonitor: lagMonitor,
	}
}

func (rmc *ReplicationMetricsCollector) RecordWrite() {
	rmc.mu.Lock()
	defer rmc.mu.Unlock()

	rmc.metrics.TotalWrites++
}

func (rmc *ReplicationMetricsCollector) RecordRead() {
	rmc.mu.Lock()
	defer rmc.mu.Unlock()

	rmc.metrics.TotalReads++
}

func (rmc *ReplicationMetricsCollector) RecordFailover() {
	rmc.mu.Lock()
	defer rmc.mu.Unlock()

	rmc.metrics.FailoverCount++
	rmc.metrics.LastFailoverTime = time.Now()
}

func (rmc *ReplicationMetricsCollector) GetMetrics() ReplicationMetrics {
	rmc.mu.RLock()
	defer rmc.mu.RUnlock()

	health := rmc.strategy.GetHealthStatus()
	rmc.metrics.HealthScore = health.HealthScore
	rmc.metrics.AverageReplicationLag = health.MaxReplicationLag

	return rmc.metrics
}

// ReplicationRecommendationEngine provides optimization suggestions
// FAANG Interview Point: Automated optimization and scaling decisions
type ReplicationRecommendationEngine struct {
	metricsCollector *ReplicationMetricsCollector
	thresholds       ReplicationThresholds
}

type ReplicationThresholds struct {
	MaxAcceptableLag    time.Duration
	MinHealthScore      float64
	MaxReadWriteRatio   float64
	MaxSlaveUtilization float64
}

func NewReplicationRecommendationEngine(collector *ReplicationMetricsCollector) *ReplicationRecommendationEngine {
	return &ReplicationRecommendationEngine{
		metricsCollector: collector,
		thresholds: ReplicationThresholds{
			MaxAcceptableLag:    5 * time.Second,
			MinHealthScore:      0.8,
			MaxReadWriteRatio:   10.0,
			MaxSlaveUtilization: 0.8,
		},
	}
}

func (rre *ReplicationRecommendationEngine) GenerateRecommendations() []ReplicationRecommendation {
	metrics := rre.metricsCollector.GetMetrics()
	var recommendations []ReplicationRecommendation

	// Check replication lag
	if metrics.AverageReplicationLag > rre.thresholds.MaxAcceptableLag {
		recommendations = append(recommendations, ReplicationRecommendation{
			Type:        "REDUCE_LAG",
			Description: fmt.Sprintf("Replication lag is %v, consider upgrading network or reducing write load", metrics.AverageReplicationLag),
			Priority:    "HIGH",
			Action:      "Increase replication bandwidth or add read replicas",
		})
	}

	// Check health score
	if metrics.HealthScore < rre.thresholds.MinHealthScore {
		recommendations = append(recommendations, ReplicationRecommendation{
			Type:        "IMPROVE_HEALTH",
			Description: fmt.Sprintf("Health score is %.2f, below threshold of %.2f", metrics.HealthScore, rre.thresholds.MinHealthScore),
			Priority:    "MEDIUM",
			Action:      "Add more replicas or fix failed nodes",
		})
	}

	// Check read/write ratio
	if metrics.TotalReads > 0 {
		readWriteRatio := float64(metrics.TotalReads) / float64(metrics.TotalWrites)
		if readWriteRatio > rre.thresholds.MaxReadWriteRatio {
			recommendations = append(recommendations, ReplicationRecommendation{
				Type:        "SCALE_READS",
				Description: fmt.Sprintf("Read/write ratio is %.2f, consider adding read replicas", readWriteRatio),
				Priority:    "LOW",
				Action:      "Add geographic read replicas for better read scaling",
			})
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, ReplicationRecommendation{
			Type:        "HEALTHY",
			Description: "Replication system is operating within acceptable parameters",
			Priority:    "INFO",
			Action:      "Continue monitoring",
		})
	}

	return recommendations
}

type ReplicationRecommendation struct {
	Type        string
	Description string
	Priority    string
	Action      string
}
