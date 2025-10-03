package database_design

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"log"
	"sort"
	"sync"
	"time"
)

// Database Sharding Strategies for System Design Interviews
// Interview Focus: Horizontal scaling, data distribution, and query routing

// ShardingStrategy defines interface for different sharding approaches
type ShardingStrategy interface {
	GetShardID(key string) (string, error)
	AddShard(shardID string) error
	RemoveShard(shardID string) error
	GetAllShards() []string
	RebalanceShards() error
}

// HashBasedSharding implements hash-based data distribution
// FAANG Interview Point: Simple, even distribution, no range queries
type HashBasedSharding struct {
	shards   []string
	mu       sync.RWMutex
	hashFunc func(string) uint32
	replicas int // Virtual nodes for consistent hashing
}

func NewHashBasedSharding(shards []string) *HashBasedSharding {
	return &HashBasedSharding{
		shards:   shards,
		hashFunc: func(key string) uint32 { return crc32.ChecksumIEEE([]byte(key)) },
		replicas: 150, // Common choice for consistent hashing
	}
}

// GetShardID returns shard for given key using consistent hashing
// FAANG Interview Point: Handles shard additions/removals gracefully
func (h *HashBasedSharding) GetShardID(key string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.shards) == 0 {
		return "", fmt.Errorf("no shards available")
	}

	// Simple hash for demonstration (in production, use consistent hashing)
	hash := h.hashFunc(key)
	shardIndex := int(hash) % len(h.shards)
	return h.shards[shardIndex], nil
}

func (h *HashBasedSharding) AddShard(shardID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.shards = append(h.shards, shardID)
	log.Printf("Added shard: %s. Total shards: %d", shardID, len(h.shards))
	return nil
}

func (h *HashBasedSharding) RemoveShard(shardID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i, shard := range h.shards {
		if shard == shardID {
			h.shards = append(h.shards[:i], h.shards[i+1:]...)
			log.Printf("Removed shard: %s. Total shards: %d", shardID, len(h.shards))
			return nil
		}
	}
	return fmt.Errorf("shard %s not found", shardID)
}

func (h *HashBasedSharding) GetAllShards() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]string, len(h.shards))
	copy(result, h.shards)
	return result
}

func (h *HashBasedSharding) RebalanceShards() error {
	// Hash-based sharding doesn't require explicit rebalancing
	// New hash function distributes data automatically
	log.Println("Hash-based sharding: no explicit rebalancing needed")
	return nil
}

// RangeBasedSharding implements range-based data partitioning
// FAANG Interview Point: Supports range queries but can create hotspots
type RangeBasedSharding struct {
	ranges map[string]ShardRange // shardID -> range
	mu     sync.RWMutex
}

type ShardRange struct {
	Start string
	End   string
	Shard string
}

func NewRangeBasedSharding() *RangeBasedSharding {
	return &RangeBasedSharding{
		ranges: make(map[string]ShardRange),
	}
}

// AddShardRange adds a new shard with specific range
func (r *RangeBasedSharding) AddShardRange(shardID, start, end string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ranges[shardID] = ShardRange{
		Start: start,
		End:   end,
		Shard: shardID,
	}

	log.Printf("Added shard %s with range [%s, %s)", shardID, start, end)
	return nil
}

func (r *RangeBasedSharding) GetShardID(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, shardRange := range r.ranges {
		if key >= shardRange.Start && key < shardRange.End {
			return shardRange.Shard, nil
		}
	}

	return "", fmt.Errorf("no shard found for key: %s", key)
}

func (r *RangeBasedSharding) AddShard(shardID string) error {
	// Range-based sharding requires explicit range assignment
	return fmt.Errorf("use AddShardRange for range-based sharding")
}

func (r *RangeBasedSharding) RemoveShard(shardID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.ranges, shardID)
	log.Printf("Removed shard: %s", shardID)
	return nil
}

func (r *RangeBasedSharding) GetAllShards() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var shards []string
	for _, shardRange := range r.ranges {
		shards = append(shards, shardRange.Shard)
	}
	return shards
}

func (r *RangeBasedSharding) RebalanceShards() error {
	// Range rebalancing requires splitting hot ranges
	log.Println("Range-based rebalancing: split hot ranges, redistribute load")
	return nil
}

// ConsistentHashing implements consistent hashing with virtual nodes
// FAANG Interview Point: Industry standard for distributed systems
type ConsistentHashing struct {
	ring         map[uint32]string // hash -> shard
	sortedHashes []uint32
	shards       map[string]bool
	replicas     int
	mu           sync.RWMutex
}

func NewConsistentHashing(replicas int) *ConsistentHashing {
	return &ConsistentHashing{
		ring:     make(map[uint32]string),
		shards:   make(map[string]bool),
		replicas: replicas,
	}
}

func (c *ConsistentHashing) hash(key string) uint32 {
	h := md5.Sum([]byte(key))
	return uint32(h[0])<<24 | uint32(h[1])<<16 | uint32(h[2])<<8 | uint32(h[3])
}

func (c *ConsistentHashing) AddShard(shardID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.shards[shardID] {
		return fmt.Errorf("shard %s already exists", shardID)
	}

	// Add virtual nodes for better distribution
	for i := 0; i < c.replicas; i++ {
		virtualKey := fmt.Sprintf("%s:%d", shardID, i)
		hash := c.hash(virtualKey)
		c.ring[hash] = shardID
		c.sortedHashes = append(c.sortedHashes, hash)
	}

	sort.Slice(c.sortedHashes, func(i, j int) bool {
		return c.sortedHashes[i] < c.sortedHashes[j]
	})

	c.shards[shardID] = true
	log.Printf("Added shard %s with %d virtual nodes", shardID, c.replicas)
	return nil
}

func (c *ConsistentHashing) RemoveShard(shardID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.shards[shardID] {
		return fmt.Errorf("shard %s not found", shardID)
	}

	// Remove all virtual nodes for this shard
	for i := 0; i < c.replicas; i++ {
		virtualKey := fmt.Sprintf("%s:%d", shardID, i)
		hash := c.hash(virtualKey)
		delete(c.ring, hash)

		// Remove from sorted hashes
		for j, h := range c.sortedHashes {
			if h == hash {
				c.sortedHashes = append(c.sortedHashes[:j], c.sortedHashes[j+1:]...)
				break
			}
		}
	}

	delete(c.shards, shardID)
	log.Printf("Removed shard %s", shardID)
	return nil
}

func (c *ConsistentHashing) GetShardID(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.ring) == 0 {
		return "", fmt.Errorf("no shards available")
	}

	hash := c.hash(key)

	// Find the first shard >= hash (clockwise on ring)
	for _, h := range c.sortedHashes {
		if h >= hash {
			return c.ring[h], nil
		}
	}

	// Wrap around to beginning of ring
	return c.ring[c.sortedHashes[0]], nil
}

func (c *ConsistentHashing) GetAllShards() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var shards []string
	for shard := range c.shards {
		shards = append(shards, shard)
	}
	return shards
}

func (c *ConsistentHashing) RebalanceShards() error {
	// Consistent hashing automatically rebalances
	log.Println("Consistent hashing: automatic rebalancing via virtual nodes")
	return nil
}

// Shard Key Selection Strategies
// Interview Focus: Critical for even data distribution and query performance

// ShardKeyAnalyzer helps select optimal shard keys
type ShardKeyAnalyzer struct {
	metrics map[string]ShardKeyMetrics
}

type ShardKeyMetrics struct {
	Cardinality  int     // Number of distinct values
	Distribution float64 // How evenly distributed (0-1, 1 is perfect)
	QueryPattern QueryPatternType
	GrowthRate   float64 // How fast key space grows
	Hotspots     int     // Number of hot keys
}

type QueryPatternType int

const (
	PointLookup QueryPatternType = iota
	RangeQuery
	PrefixQuery
	Analytics
)

// AnalyzeShardKey evaluates shard key effectiveness
// FAANG Interview Point: Data-driven shard key selection
func (ska *ShardKeyAnalyzer) AnalyzeShardKey(keyName string, sampleData []string) ShardKeyMetrics {
	// Calculate cardinality
	uniqueKeys := make(map[string]bool)
	for _, key := range sampleData {
		uniqueKeys[key] = true
	}
	cardinality := len(uniqueKeys)

	// Simulate distribution analysis
	distribution := float64(cardinality) / float64(len(sampleData))
	if distribution > 0.8 {
		distribution = 0.9 // High cardinality = good distribution
	}

	// Simulate hotspot detection
	keyFreq := make(map[string]int)
	for _, key := range sampleData {
		keyFreq[key]++
	}

	hotspots := 0
	avgFreq := len(sampleData) / cardinality
	for _, freq := range keyFreq {
		if freq > avgFreq*3 { // 3x average = hotspot
			hotspots++
		}
	}

	metrics := ShardKeyMetrics{
		Cardinality:  cardinality,
		Distribution: distribution,
		QueryPattern: PointLookup, // Default
		GrowthRate:   0.1,         // 10% growth
		Hotspots:     hotspots,
	}

	log.Printf("Shard key '%s' analysis: cardinality=%d, distribution=%.2f, hotspots=%d",
		keyName, cardinality, distribution, hotspots)

	return metrics
}

// RecommendShardKey provides shard key recommendations
// FAANG Interview Point: System-specific shard key selection
func (ska *ShardKeyAnalyzer) RecommendShardKey(systemType string) ShardKeyRecommendation {
	recommendations := map[string]ShardKeyRecommendation{
		"social_media": {
			Primary:        "user_id",
			Reasoning:      "High cardinality, even distribution, supports user-centric queries",
			Alternatives:   []string{"hash(user_id)", "user_id + timestamp"},
			Considerations: "Avoid celebrity user hotspots, consider hash prefix",
		},
		"e_commerce": {
			Primary:        "customer_id",
			Reasoning:      "Natural partition boundary, supports customer queries",
			Alternatives:   []string{"order_id", "hash(customer_id + product_category)"},
			Considerations: "Watch for high-volume customers, geographic distribution",
		},
		"messaging": {
			Primary:        "conversation_id",
			Reasoning:      "Groups related messages, supports conversation queries",
			Alternatives:   []string{"hash(sender_id + receiver_id)", "message_id"},
			Considerations: "Group chats can create hotspots, consider hash-based approach",
		},
		"analytics": {
			Primary:        "date + hash(user_id)",
			Reasoning:      "Time-based partitioning with user distribution",
			Alternatives:   []string{"hash(event_type + user_id)", "timestamp"},
			Considerations: "Balance between time-based queries and even distribution",
		},
		"gaming": {
			Primary:        "game_session_id",
			Reasoning:      "Natural boundary, supports session-based queries",
			Alternatives:   []string{"player_id", "hash(player_id + game_type)"},
			Considerations: "Popular games create hotspots, consider game-type distribution",
		},
	}

	if rec, exists := recommendations[systemType]; exists {
		return rec
	}

	return ShardKeyRecommendation{
		Primary:        "id",
		Reasoning:      "Default primary key sharding",
		Alternatives:   []string{"hash(id)"},
		Considerations: "Analyze query patterns and data distribution",
	}
}

type ShardKeyRecommendation struct {
	Primary        string
	Reasoning      string
	Alternatives   []string
	Considerations string
}

// Cross-Shard Query Handling
// Interview Focus: Distributed query execution and aggregation

// CrossShardQueryExecutor handles queries spanning multiple shards
type CrossShardQueryExecutor struct {
	sharding ShardingStrategy
	timeout  time.Duration
}

func NewCrossShardQueryExecutor(sharding ShardingStrategy) *CrossShardQueryExecutor {
	return &CrossShardQueryExecutor{
		sharding: sharding,
		timeout:  5 * time.Second,
	}
}

// QueryResult represents result from a single shard
type QueryResult struct {
	ShardID string
	Data    interface{}
	Error   error
	Latency time.Duration
}

// ExecuteScatterGatherQuery executes query on all shards
// FAANG Interview Point: Distributed aggregation patterns
func (cq *CrossShardQueryExecutor) ExecuteScatterGatherQuery(query string) ([]QueryResult, error) {
	shards := cq.sharding.GetAllShards()
	resultChan := make(chan QueryResult, len(shards))

	// Execute query on all shards in parallel
	for _, shardID := range shards {
		go func(shard string) {
			start := time.Now()

			// Simulate query execution
			result := cq.executeQueryOnShard(shard, query)
			result.ShardID = shard
			result.Latency = time.Since(start)

			resultChan <- result
		}(shardID)
	}

	// Collect results with timeout
	var results []QueryResult
	timeout := time.After(cq.timeout)

	for i := 0; i < len(shards); i++ {
		select {
		case result := <-resultChan:
			results = append(results, result)
		case <-timeout:
			return results, fmt.Errorf("query timeout after %v", cq.timeout)
		}
	}

	return results, nil
}

// ExecuteRoutedQuery executes query on specific shards based on routing key
func (cq *CrossShardQueryExecutor) ExecuteRoutedQuery(routingKey, query string) (QueryResult, error) {
	shardID, err := cq.sharding.GetShardID(routingKey)
	if err != nil {
		return QueryResult{}, err
	}

	start := time.Now()
	result := cq.executeQueryOnShard(shardID, query)
	result.ShardID = shardID
	result.Latency = time.Since(start)

	return result, result.Error
}

func (cq *CrossShardQueryExecutor) executeQueryOnShard(shardID, query string) QueryResult {
	// Simulate query execution with varying latency
	latency := time.Duration(10+len(query)%50) * time.Millisecond
	time.Sleep(latency)

	// Simulate success/failure
	if len(query)%7 == 0 {
		return QueryResult{
			Error: fmt.Errorf("query failed on shard %s", shardID),
		}
	}

	return QueryResult{
		Data: fmt.Sprintf("Results from shard %s for query: %s", shardID, query),
	}
}

// Shard Rebalancing Algorithms
// Interview Focus: Dynamic load balancing and data migration

// ShardRebalancer handles automatic shard rebalancing
type ShardRebalancer struct {
	sharding    ShardingStrategy
	metrics     map[string]ShardMetrics
	mu          sync.RWMutex
	rebalancing bool
}

type ShardMetrics struct {
	DataSize        int64   // Bytes
	QueryRate       float64 // QPS
	CPUUsage        float64 // 0-1
	MemoryUsage     float64 // 0-1
	ConnectionCount int
	LastUpdated     time.Time
}

func NewShardRebalancer(sharding ShardingStrategy) *ShardRebalancer {
	return &ShardRebalancer{
		sharding: sharding,
		metrics:  make(map[string]ShardMetrics),
	}
}

// UpdateShardMetrics records current shard performance
func (sr *ShardRebalancer) UpdateShardMetrics(shardID string, metrics ShardMetrics) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	metrics.LastUpdated = time.Now()
	sr.metrics[shardID] = metrics
}

// AnalyzeImbalance detects load imbalances across shards
// FAANG Interview Point: Automated load balancing triggers
func (sr *ShardRebalancer) AnalyzeImbalance() ImbalanceAnalysis {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	if len(sr.metrics) < 2 {
		return ImbalanceAnalysis{Balanced: true}
	}

	var totalQPS, totalCPU, totalMemory float64
	var maxQPS, minQPS, maxCPU, minCPU float64

	first := true
	for _, metrics := range sr.metrics {
		totalQPS += metrics.QueryRate
		totalCPU += metrics.CPUUsage
		totalMemory += metrics.MemoryUsage

		if first {
			maxQPS, minQPS = metrics.QueryRate, metrics.QueryRate
			maxCPU, minCPU = metrics.CPUUsage, metrics.CPUUsage
			first = false
		} else {
			if metrics.QueryRate > maxQPS {
				maxQPS = metrics.QueryRate
			}
			if metrics.QueryRate < minQPS {
				minQPS = metrics.QueryRate
			}
			if metrics.CPUUsage > maxCPU {
				maxCPU = metrics.CPUUsage
			}
			if metrics.CPUUsage < minCPU {
				minCPU = metrics.CPUUsage
			}
		}
	}

	shardCount := float64(len(sr.metrics))
	avgQPS := totalQPS / shardCount
	avgCPU := totalCPU / shardCount

	// Calculate imbalance ratios
	qpsImbalance := maxQPS / avgQPS
	cpuImbalance := maxCPU / avgCPU

	balanced := qpsImbalance < 2.0 && cpuImbalance < 1.5 // Thresholds

	analysis := ImbalanceAnalysis{
		Balanced:          balanced,
		QPSImbalance:      qpsImbalance,
		CPUImbalance:      cpuImbalance,
		AvgQPS:            avgQPS,
		MaxQPS:            maxQPS,
		MinQPS:            minQPS,
		RecommendedAction: "",
	}

	if !balanced {
		if qpsImbalance > 3.0 {
			analysis.RecommendedAction = "Split hot shards or add read replicas"
		} else if cpuImbalance > 2.0 {
			analysis.RecommendedAction = "Migrate data from overloaded shards"
		} else {
			analysis.RecommendedAction = "Monitor and consider gradual rebalancing"
		}
	}

	return analysis
}

type ImbalanceAnalysis struct {
	Balanced          bool
	QPSImbalance      float64
	CPUImbalance      float64
	AvgQPS            float64
	MaxQPS            float64
	MinQPS            float64
	RecommendedAction string
}

// ExecuteRebalancing performs automatic shard rebalancing
func (sr *ShardRebalancer) ExecuteRebalancing() error {
	sr.mu.Lock()
	if sr.rebalancing {
		sr.mu.Unlock()
		return fmt.Errorf("rebalancing already in progress")
	}
	sr.rebalancing = true
	sr.mu.Unlock()

	defer func() {
		sr.mu.Lock()
		sr.rebalancing = false
		sr.mu.Unlock()
	}()

	analysis := sr.AnalyzeImbalance()
	if analysis.Balanced {
		log.Println("Shards are balanced, no rebalancing needed")
		return nil
	}

	log.Printf("Starting rebalancing: %s", analysis.RecommendedAction)

	// Simulate rebalancing steps
	steps := []string{
		"Identifying hot spots and cold shards",
		"Calculating optimal data migration plan",
		"Creating new shards if needed",
		"Migrating data in background",
		"Updating routing tables",
		"Verifying data consistency",
	}

	for i, step := range steps {
		log.Printf("Rebalancing step %d/%d: %s", i+1, len(steps), step)
		time.Sleep(100 * time.Millisecond) // Simulate work
	}

	log.Println("Rebalancing completed successfully")
	return nil
}

// Hotspot Management
// Interview Focus: Handling uneven data access patterns

// HotspotDetector identifies and manages database hotspots
type HotspotDetector struct {
	accessPatterns map[string]AccessPattern
	mu             sync.RWMutex
	threshold      float64 // Hotspot threshold (QPS ratio)
}

type AccessPattern struct {
	Key          string
	AccessCount  int64
	LastAccessed time.Time
	AvgLatency   time.Duration
	ShardID      string
}

func NewHotspotDetector(threshold float64) *HotspotDetector {
	return &HotspotDetector{
		accessPatterns: make(map[string]AccessPattern),
		threshold:      threshold,
	}
}

// RecordAccess tracks key access patterns
func (hd *HotspotDetector) RecordAccess(key, shardID string, latency time.Duration) {
	hd.mu.Lock()
	defer hd.mu.Unlock()

	pattern, exists := hd.accessPatterns[key]
	if !exists {
		pattern = AccessPattern{
			Key:        key,
			ShardID:    shardID,
			AvgLatency: latency,
		}
	}

	pattern.AccessCount++
	pattern.LastAccessed = time.Now()

	// Update running average latency
	if exists {
		pattern.AvgLatency = (pattern.AvgLatency + latency) / 2
	}

	hd.accessPatterns[key] = pattern
}

// DetectHotspots identifies keys with unusually high access rates
// FAANG Interview Point: Proactive hotspot detection and mitigation
func (hd *HotspotDetector) DetectHotspots() []HotspotInfo {
	hd.mu.RLock()
	defer hd.mu.RUnlock()

	// Calculate average access count
	var totalAccess int64
	for _, pattern := range hd.accessPatterns {
		totalAccess += pattern.AccessCount
	}

	if len(hd.accessPatterns) == 0 {
		return nil
	}

	avgAccess := float64(totalAccess) / float64(len(hd.accessPatterns))
	hotspotThreshold := avgAccess * hd.threshold

	var hotspots []HotspotInfo
	for key, pattern := range hd.accessPatterns {
		if float64(pattern.AccessCount) > hotspotThreshold {
			hotspots = append(hotspots, HotspotInfo{
				Key:         key,
				AccessCount: pattern.AccessCount,
				ShardID:     pattern.ShardID,
				Severity:    float64(pattern.AccessCount) / avgAccess,
				AvgLatency:  pattern.AvgLatency,
			})
		}
	}

	// Sort by severity (highest first)
	sort.Slice(hotspots, func(i, j int) bool {
		return hotspots[i].Severity > hotspots[j].Severity
	})

	return hotspots
}

type HotspotInfo struct {
	Key         string
	AccessCount int64
	ShardID     string
	Severity    float64 // Access ratio vs average
	AvgLatency  time.Duration
}

// MitigateHotspots implements hotspot mitigation strategies
func (hd *HotspotDetector) MitigateHotspots(hotspots []HotspotInfo) []MitigationAction {
	var actions []MitigationAction

	for _, hotspot := range hotspots {
		var action MitigationAction

		if hotspot.Severity > 10.0 {
			// Critical hotspot - immediate action needed
			action = MitigationAction{
				Type:        "CACHE_HOT_DATA",
				Target:      hotspot.Key,
				Description: fmt.Sprintf("Cache key %s with %dx average traffic", hotspot.Key, int(hotspot.Severity)),
				Priority:    "HIGH",
			}
		} else if hotspot.Severity > 5.0 {
			// Significant hotspot - create read replicas
			action = MitigationAction{
				Type:        "CREATE_READ_REPLICA",
				Target:      hotspot.ShardID,
				Description: fmt.Sprintf("Add read replica for shard %s", hotspot.ShardID),
				Priority:    "MEDIUM",
			}
		} else {
			// Moderate hotspot - monitor and consider load balancing
			action = MitigationAction{
				Type:        "MONITOR_AND_BALANCE",
				Target:      hotspot.Key,
				Description: fmt.Sprintf("Monitor key %s and balance load", hotspot.Key),
				Priority:    "LOW",
			}
		}

		actions = append(actions, action)
	}

	return actions
}

type MitigationAction struct {
	Type        string
	Target      string
	Description string
	Priority    string
}

// ShardingSummaryReporter provides comprehensive sharding analysis
// FAANG Interview Point: Operational visibility and monitoring
type ShardingSummaryReporter struct {
	strategy        ShardingStrategy
	rebalancer      *ShardRebalancer
	hotspotDetector *HotspotDetector
}

func NewShardingSummaryReporter(strategy ShardingStrategy, rebalancer *ShardRebalancer, detector *HotspotDetector) *ShardingSummaryReporter {
	return &ShardingSummaryReporter{
		strategy:        strategy,
		rebalancer:      rebalancer,
		hotspotDetector: detector,
	}
}

// GenerateShardingReport creates comprehensive sharding status report
func (ssr *ShardingSummaryReporter) GenerateShardingReport() ShardingReport {
	shards := ssr.strategy.GetAllShards()
	imbalanceAnalysis := ssr.rebalancer.AnalyzeImbalance()
	hotspots := ssr.hotspotDetector.DetectHotspots()

	report := ShardingReport{
		Timestamp:       time.Now(),
		TotalShards:     len(shards),
		ShardIDs:        shards,
		IsBalanced:      imbalanceAnalysis.Balanced,
		Imbalance:       imbalanceAnalysis,
		HotspotCount:    len(hotspots),
		Hotspots:        hotspots,
		Recommendations: []string{},
	}

	// Generate recommendations
	if !report.IsBalanced {
		report.Recommendations = append(report.Recommendations,
			fmt.Sprintf("Load imbalance detected: %s", imbalanceAnalysis.RecommendedAction))
	}

	if len(hotspots) > 0 {
		report.Recommendations = append(report.Recommendations,
			fmt.Sprintf("Found %d hotspots requiring attention", len(hotspots)))
	}

	if len(shards) > 20 {
		report.Recommendations = append(report.Recommendations,
			"Consider shard consolidation for operational simplicity")
	}

	if len(report.Recommendations) == 0 {
		report.Recommendations = append(report.Recommendations, "Sharding configuration is healthy")
	}

	return report
}

type ShardingReport struct {
	Timestamp       time.Time
	TotalShards     int
	ShardIDs        []string
	IsBalanced      bool
	Imbalance       ImbalanceAnalysis
	HotspotCount    int
	Hotspots        []HotspotInfo
	Recommendations []string
}
