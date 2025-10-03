package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

// FAANG Interview Focus: Data Partitioning and Sharding Strategies
// Key Topics: Horizontal/vertical partitioning, consistent hashing, rebalancing

// PartitionStrategy defines different partitioning strategies
type PartitionStrategy int

const (
	HashBasedPartitioning PartitionStrategy = iota
	RangeBasedPartitioning
	DirectoryBasedPartitioning
	ConsistentHashPartitioning
)

// ShardedRecord represents a record with partitioning metadata
type ShardedRecord struct {
	ID          string      `json:"id"`
	Data        interface{} `json:"data"`
	ShardKey    string      `json:"shard_key"`
	Timestamp   time.Time   `json:"timestamp"`
	PartitionID int         `json:"partition_id"`
}

// PartitionNode represents a database shard/partition
type PartitionNode struct {
	ID       int                      `json:"id"`
	Name     string                   `json:"name"`
	Address  string                   `json:"address"`
	Capacity int64                    `json:"capacity"`
	Used     int64                    `json:"used"`
	Status   string                   `json:"status"` // active, readonly, draining
	Records  map[string]ShardedRecord `json:"records"`
	mu       sync.RWMutex
}

// NewPartitionNode creates a new partition node
func NewPartitionNode(id int, name, address string, capacity int64) *PartitionNode {
	return &PartitionNode{
		ID:       id,
		Name:     name,
		Address:  address,
		Capacity: capacity,
		Used:     0,
		Status:   "active",
		Records:  make(map[string]ShardedRecord),
	}
}

// AddRecord adds a record to the partition
func (pn *PartitionNode) AddRecord(record ShardedRecord) error {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	if pn.Status == "readonly" {
		return fmt.Errorf("partition %d is readonly", pn.ID)
	}

	if pn.Used >= pn.Capacity {
		return fmt.Errorf("partition %d is at capacity", pn.ID)
	}

	record.PartitionID = pn.ID
	pn.Records[record.ID] = record
	pn.Used++

	return nil
}

// GetRecord retrieves a record from the partition
func (pn *PartitionNode) GetRecord(recordID string) (*ShardedRecord, bool) {
	pn.mu.RLock()
	defer pn.mu.RUnlock()

	record, exists := pn.Records[recordID]
	return &record, exists
}

// RemoveRecord removes a record from the partition
func (pn *PartitionNode) RemoveRecord(recordID string) error {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	if _, exists := pn.Records[recordID]; !exists {
		return fmt.Errorf("record %s not found in partition %d", recordID, pn.ID)
	}

	delete(pn.Records, recordID)
	pn.Used--

	return nil
}

// GetStats returns partition statistics
func (pn *PartitionNode) GetStats() map[string]interface{} {
	pn.mu.RLock()
	defer pn.mu.RUnlock()

	return map[string]interface{}{
		"id":           pn.ID,
		"name":         pn.Name,
		"status":       pn.Status,
		"capacity":     pn.Capacity,
		"used":         pn.Used,
		"utilization":  float64(pn.Used) / float64(pn.Capacity),
		"record_count": len(pn.Records),
	}
}

// HashBasedPartitioner implements hash-based sharding
type HashBasedPartitioner struct {
	partitions []*PartitionNode
	numShards  int
	hashFunc   func(string) uint32
}

// NewHashBasedPartitioner creates a hash-based partitioner
func NewHashBasedPartitioner(partitions []*PartitionNode) *HashBasedPartitioner {
	return &HashBasedPartitioner{
		partitions: partitions,
		numShards:  len(partitions),
		hashFunc:   defaultHashFunction,
	}
}

// defaultHashFunction provides consistent hashing
func defaultHashFunction(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// GetPartition determines which partition a key should go to
func (hbp *HashBasedPartitioner) GetPartition(shardKey string) int {
	hash := hbp.hashFunc(shardKey)
	return int(hash) % hbp.numShards
}

// AddRecord adds a record using hash-based partitioning
func (hbp *HashBasedPartitioner) AddRecord(record ShardedRecord) error {
	partitionID := hbp.GetPartition(record.ShardKey)

	if partitionID >= len(hbp.partitions) {
		return fmt.Errorf("invalid partition ID: %d", partitionID)
	}

	return hbp.partitions[partitionID].AddRecord(record)
}

// GetRecord retrieves a record using hash-based partitioning
func (hbp *HashBasedPartitioner) GetRecord(recordID, shardKey string) (*ShardedRecord, error) {
	partitionID := hbp.GetPartition(shardKey)

	if partitionID >= len(hbp.partitions) {
		return nil, fmt.Errorf("invalid partition ID: %d", partitionID)
	}

	record, exists := hbp.partitions[partitionID].GetRecord(recordID)
	if !exists {
		return nil, fmt.Errorf("record %s not found", recordID)
	}

	return record, nil
}

// RangeBasedPartitioner implements range-based sharding
type RangeBasedPartitioner struct {
	partitions []*PartitionNode
	ranges     []PartitionRange
}

// PartitionRange defines a range for range-based partitioning
type PartitionRange struct {
	StartKey    string `json:"start_key"`
	EndKey      string `json:"end_key"`
	PartitionID int    `json:"partition_id"`
}

// NewRangeBasedPartitioner creates a range-based partitioner
func NewRangeBasedPartitioner(partitions []*PartitionNode, ranges []PartitionRange) *RangeBasedPartitioner {
	// Sort ranges by start key
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].StartKey < ranges[j].StartKey
	})

	return &RangeBasedPartitioner{
		partitions: partitions,
		ranges:     ranges,
	}
}

// GetPartition determines which partition a key should go to based on ranges
func (rbp *RangeBasedPartitioner) GetPartition(shardKey string) int {
	for _, r := range rbp.ranges {
		if shardKey >= r.StartKey && (r.EndKey == "" || shardKey < r.EndKey) {
			return r.PartitionID
		}
	}

	// Default to last partition if no range matches
	if len(rbp.ranges) > 0 {
		return rbp.ranges[len(rbp.ranges)-1].PartitionID
	}

	return 0
}

// AddRecord adds a record using range-based partitioning
func (rbp *RangeBasedPartitioner) AddRecord(record ShardedRecord) error {
	partitionID := rbp.GetPartition(record.ShardKey)

	if partitionID >= len(rbp.partitions) {
		return fmt.Errorf("invalid partition ID: %d", partitionID)
	}

	return rbp.partitions[partitionID].AddRecord(record)
}

// ConsistentHashRing implements consistent hashing for data partitioning
type ConsistentHashRing struct {
	ring         map[uint32]*PartitionNode
	sortedHashes []uint32
	virtualNodes int
	mu           sync.RWMutex
}

// NewConsistentHashRing creates a consistent hash ring
func NewConsistentHashRing(virtualNodes int) *ConsistentHashRing {
	return &ConsistentHashRing{
		ring:         make(map[uint32]*PartitionNode),
		virtualNodes: virtualNodes,
	}
}

// AddNode adds a node to the consistent hash ring
func (chr *ConsistentHashRing) AddNode(node *PartitionNode) {
	chr.mu.Lock()
	defer chr.mu.Unlock()

	for i := 0; i < chr.virtualNodes; i++ {
		hash := chr.hashNode(fmt.Sprintf("%s:%d:%d", node.Address, node.ID, i))
		chr.ring[hash] = node
		chr.sortedHashes = append(chr.sortedHashes, hash)
	}

	sort.Slice(chr.sortedHashes, func(i, j int) bool {
		return chr.sortedHashes[i] < chr.sortedHashes[j]
	})
}

// RemoveNode removes a node from the consistent hash ring
func (chr *ConsistentHashRing) RemoveNode(node *PartitionNode) {
	chr.mu.Lock()
	defer chr.mu.Unlock()

	for i := 0; i < chr.virtualNodes; i++ {
		hash := chr.hashNode(fmt.Sprintf("%s:%d:%d", node.Address, node.ID, i))
		delete(chr.ring, hash)

		// Remove from sorted hashes
		for j, h := range chr.sortedHashes {
			if h == hash {
				chr.sortedHashes = append(chr.sortedHashes[:j], chr.sortedHashes[j+1:]...)
				break
			}
		}
	}
}

// hashNode creates a hash for a node
func (chr *ConsistentHashRing) hashNode(nodeKey string) uint32 {
	hash := md5.Sum([]byte(nodeKey))
	return binary.LittleEndian.Uint32(hash[:4])
}

// GetNode gets the node responsible for a key
func (chr *ConsistentHashRing) GetNode(key string) *PartitionNode {
	chr.mu.RLock()
	defer chr.mu.RUnlock()

	if len(chr.sortedHashes) == 0 {
		return nil
	}

	hash := chr.hashNode(key)

	// Find the first node hash >= key hash
	idx := sort.Search(len(chr.sortedHashes), func(i int) bool {
		return chr.sortedHashes[i] >= hash
	})

	// Wrap around to the first node if we've gone past the end
	if idx == len(chr.sortedHashes) {
		idx = 0
	}

	return chr.ring[chr.sortedHashes[idx]]
}

// GetNodesForReplication gets N nodes for replication
func (chr *ConsistentHashRing) GetNodesForReplication(key string, replicationFactor int) []*PartitionNode {
	chr.mu.RLock()
	defer chr.mu.RUnlock()

	if len(chr.sortedHashes) == 0 {
		return nil
	}

	hash := chr.hashNode(key)
	var nodes []*PartitionNode
	seenNodes := make(map[*PartitionNode]bool)

	// Find starting position
	idx := sort.Search(len(chr.sortedHashes), func(i int) bool {
		return chr.sortedHashes[i] >= hash
	})

	// Collect unique nodes
	for len(nodes) < replicationFactor && len(nodes) < len(chr.ring) {
		if idx >= len(chr.sortedHashes) {
			idx = 0
		}

		node := chr.ring[chr.sortedHashes[idx]]
		if !seenNodes[node] {
			nodes = append(nodes, node)
			seenNodes[node] = true
		}

		idx++

		// Prevent infinite loop
		if idx == 0 && len(nodes) == 0 {
			break
		}
	}

	return nodes
}

// VerticalPartitioner handles vertical partitioning (column-based)
type VerticalPartitioner struct {
	schemas    map[string]TableSchema
	partitions map[string]*PartitionNode
}

// TableSchema defines which columns belong to which partition
type TableSchema struct {
	TableName string                     `json:"table_name"`
	Columns   map[string]ColumnPartition `json:"columns"`
}

// ColumnPartition defines column partitioning information
type ColumnPartition struct {
	PartitionName string `json:"partition_name"`
	DataType      string `json:"data_type"`
	Indexed       bool   `json:"indexed"`
}

// NewVerticalPartitioner creates a vertical partitioner
func NewVerticalPartitioner() *VerticalPartitioner {
	return &VerticalPartitioner{
		schemas:    make(map[string]TableSchema),
		partitions: make(map[string]*PartitionNode),
	}
}

// AddTableSchema adds a table schema with column partitioning
func (vp *VerticalPartitioner) AddTableSchema(schema TableSchema) {
	vp.schemas[schema.TableName] = schema
}

// AddPartition adds a named partition
func (vp *VerticalPartitioner) AddPartition(name string, node *PartitionNode) {
	vp.partitions[name] = node
}

// ShardRebalancer handles shard rebalancing operations
type ShardRebalancer struct {
	partitioner  interface{}
	rebalanceLog []RebalanceOperation
	mu           sync.Mutex
}

// RebalanceOperation records rebalancing operations
type RebalanceOperation struct {
	Timestamp       time.Time `json:"timestamp"`
	SourcePartition int       `json:"source_partition"`
	TargetPartition int       `json:"target_partition"`
	RecordsMoved    int       `json:"records_moved"`
	Operation       string    `json:"operation"` // split, merge, move
	Status          string    `json:"status"`    // pending, completed, failed
}

// NewShardRebalancer creates a shard rebalancer
func NewShardRebalancer(partitioner interface{}) *ShardRebalancer {
	return &ShardRebalancer{
		partitioner:  partitioner,
		rebalanceLog: make([]RebalanceOperation, 0),
	}
}

// RebalancePartitions rebalances data across partitions
func (sr *ShardRebalancer) RebalancePartitions(sourcePartitionID, targetPartitionID int, recordCount int) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	operation := RebalanceOperation{
		Timestamp:       time.Now(),
		SourcePartition: sourcePartitionID,
		TargetPartition: targetPartitionID,
		RecordsMoved:    recordCount,
		Operation:       "move",
		Status:          "pending",
	}

	// Simulate rebalancing operation
	log.Printf("Starting rebalance: moving %d records from partition %d to partition %d",
		recordCount, sourcePartitionID, targetPartitionID)

	// In a real implementation, this would:
	// 1. Lock the source partition
	// 2. Copy records to target partition
	// 3. Verify data integrity
	// 4. Update routing tables
	// 5. Delete records from source partition
	// 6. Update partition metadata

	operation.Status = "completed"
	sr.rebalanceLog = append(sr.rebalanceLog, operation)

	log.Printf("Rebalance completed: %+v", operation)
	return nil
}

// GetRebalanceHistory returns rebalancing history
func (sr *ShardRebalancer) GetRebalanceHistory() []RebalanceOperation {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	// Return a copy to avoid race conditions
	history := make([]RebalanceOperation, len(sr.rebalanceLog))
	copy(history, sr.rebalanceLog)

	return history
}

// CrossShardQuery handles queries across multiple shards
type CrossShardQuery struct {
	QueryID    string              `json:"query_id"`
	Query      string              `json:"query"`
	Partitions []int               `json:"partitions"`
	Results    map[int]interface{} `json:"results"`
	Status     string              `json:"status"`
	StartTime  time.Time           `json:"start_time"`
	EndTime    time.Time           `json:"end_time"`
	mu         sync.RWMutex
}

// CrossShardQueryExecutor executes queries across multiple shards
type CrossShardQueryExecutor struct {
	partitions []*PartitionNode
	queries    map[string]*CrossShardQuery
	mu         sync.RWMutex
}

// NewCrossShardQueryExecutor creates a cross-shard query executor
func NewCrossShardQueryExecutor(partitions []*PartitionNode) *CrossShardQueryExecutor {
	return &CrossShardQueryExecutor{
		partitions: partitions,
		queries:    make(map[string]*CrossShardQuery),
	}
}

// ExecuteQuery executes a query across multiple shards
func (csqe *CrossShardQueryExecutor) ExecuteQuery(queryID, query string, partitionIDs []int) *CrossShardQuery {
	csqe.mu.Lock()
	defer csqe.mu.Unlock()

	crossQuery := &CrossShardQuery{
		QueryID:    queryID,
		Query:      query,
		Partitions: partitionIDs,
		Results:    make(map[int]interface{}),
		Status:     "running",
		StartTime:  time.Now(),
	}

	csqe.queries[queryID] = crossQuery

	// Execute query on each partition concurrently
	go csqe.executeQueryConcurrently(crossQuery)

	return crossQuery
}

// executeQueryConcurrently executes query on multiple partitions
func (csqe *CrossShardQueryExecutor) executeQueryConcurrently(crossQuery *CrossShardQuery) {
	var wg sync.WaitGroup
	var resultMu sync.Mutex

	for _, partitionID := range crossQuery.Partitions {
		if partitionID >= len(csqe.partitions) {
			continue
		}

		wg.Add(1)
		go func(pID int) {
			defer wg.Done()

			// Simulate query execution on partition
			partition := csqe.partitions[pID]
			result := csqe.executeQueryOnPartition(partition, crossQuery.Query)

			resultMu.Lock()
			crossQuery.Results[pID] = result
			resultMu.Unlock()
		}(partitionID)
	}

	wg.Wait()

	crossQuery.mu.Lock()
	crossQuery.Status = "completed"
	crossQuery.EndTime = time.Now()
	crossQuery.mu.Unlock()

	log.Printf("Cross-shard query %s completed in %v",
		crossQuery.QueryID, crossQuery.EndTime.Sub(crossQuery.StartTime))
}

// executeQueryOnPartition executes query on a single partition
func (csqe *CrossShardQueryExecutor) executeQueryOnPartition(partition *PartitionNode, query string) interface{} {
	// Simulate query execution
	time.Sleep(time.Millisecond * 50) // Simulate processing time

	stats := partition.GetStats()
	return map[string]interface{}{
		"partition_id":   partition.ID,
		"record_count":   stats["record_count"],
		"query":          query,
		"execution_time": "50ms",
	}
}

// PartitioningDemo demonstrates various partitioning strategies
func PartitioningDemo() {
	log.Println("=== Data Partitioning and Sharding Demo ===")

	// Create partition nodes
	partitions := []*PartitionNode{
		NewPartitionNode(0, "shard-0", "10.0.1.1:5432", 1000),
		NewPartitionNode(1, "shard-1", "10.0.1.2:5432", 1000),
		NewPartitionNode(2, "shard-2", "10.0.1.3:5432", 1000),
	}

	// Demo 1: Hash-based partitioning
	log.Println("\n--- Hash-based Partitioning ---")
	hashPartitioner := NewHashBasedPartitioner(partitions)

	// Add records with different shard keys
	testRecords := []ShardedRecord{
		{ID: "user1", Data: "User 1 Data", ShardKey: "user1@example.com"},
		{ID: "user2", Data: "User 2 Data", ShardKey: "user2@example.com"},
		{ID: "user3", Data: "User 3 Data", ShardKey: "user3@example.com"},
	}

	for _, record := range testRecords {
		partitionID := hashPartitioner.GetPartition(record.ShardKey)
		log.Printf("Record %s with shard key %s goes to partition %d",
			record.ID, record.ShardKey, partitionID)

		if err := hashPartitioner.AddRecord(record); err != nil {
			log.Printf("Error adding record: %v", err)
		}
	}

	// Demo 2: Range-based partitioning
	log.Println("\n--- Range-based Partitioning ---")
	ranges := []PartitionRange{
		{StartKey: "A", EndKey: "H", PartitionID: 0},
		{StartKey: "H", EndKey: "P", PartitionID: 1},
		{StartKey: "P", EndKey: "", PartitionID: 2}, // End key empty means "to end"
	}

	rangePartitioner := NewRangeBasedPartitioner(partitions, ranges)

	rangeTestRecords := []ShardedRecord{
		{ID: "alice", Data: "Alice Data", ShardKey: "Alice"},
		{ID: "bob", Data: "Bob Data", ShardKey: "Bob"},
		{ID: "charlie", Data: "Charlie Data", ShardKey: "Charlie"},
		{ID: "john", Data: "John Data", ShardKey: "John"},
		{ID: "susan", Data: "Susan Data", ShardKey: "Susan"},
		{ID: "zoe", Data: "Zoe Data", ShardKey: "Zoe"},
	}

	for _, record := range rangeTestRecords {
		partitionID := rangePartitioner.GetPartition(record.ShardKey)
		log.Printf("Record %s with shard key %s goes to partition %d",
			record.ID, record.ShardKey, partitionID)

		if err := rangePartitioner.AddRecord(record); err != nil {
			log.Printf("Error adding record: %v", err)
		}
	}

	// Demo 3: Consistent hashing
	log.Println("\n--- Consistent Hashing ---")
	consistentRing := NewConsistentHashRing(150) // 150 virtual nodes per physical node

	for _, partition := range partitions {
		consistentRing.AddNode(partition)
	}

	consistentTestKeys := []string{"user1", "user2", "user3", "user4", "user5"}
	for _, key := range consistentTestKeys {
		node := consistentRing.GetNode(key)
		if node != nil {
			log.Printf("Key %s maps to node %s (partition %d)", key, node.Name, node.ID)
		}

		// Demo replication
		replicaNodes := consistentRing.GetNodesForReplication(key, 2)
		log.Printf("Key %s replicas: ", key)
		for _, replica := range replicaNodes {
			fmt.Printf("  %s (partition %d)", replica.Name, replica.ID)
		}
		fmt.Println()
	}

	// Demo 4: Cross-shard queries
	log.Println("\n--- Cross-shard Queries ---")
	queryExecutor := NewCrossShardQueryExecutor(partitions)

	crossQuery := queryExecutor.ExecuteQuery(
		"query1",
		"SELECT COUNT(*) FROM users WHERE age > 25",
		[]int{0, 1, 2},
	)

	// Wait for query completion
	for crossQuery.Status == "running" {
		time.Sleep(time.Millisecond * 10)
	}

	log.Printf("Cross-shard query results: %+v", crossQuery.Results)

	// Demo 5: Rebalancing
	log.Println("\n--- Shard Rebalancing ---")
	rebalancer := NewShardRebalancer(hashPartitioner)

	// Simulate rebalancing from partition 0 to partition 1
	if err := rebalancer.RebalancePartitions(0, 1, 100); err != nil {
		log.Printf("Rebalancing error: %v", err)
	}

	history := rebalancer.GetRebalanceHistory()
	log.Printf("Rebalance history: %+v", history)

	// Print final partition stats
	log.Println("\n--- Final Partition Statistics ---")
	for _, partition := range partitions {
		stats := partition.GetStats()
		log.Printf("Partition %d stats: %+v", partition.ID, stats)
	}
}

// PerformanceComparison compares different partitioning strategies
func PerformanceComparison() {
	log.Println("\n=== Partitioning Performance Comparison ===")

	partitions := []*PartitionNode{
		NewPartitionNode(0, "shard-0", "10.0.1.1:5432", 10000),
		NewPartitionNode(1, "shard-1", "10.0.1.2:5432", 10000),
		NewPartitionNode(2, "shard-2", "10.0.1.3:5432", 10000),
	}

	numRecords := 10000

	// Test hash-based partitioning performance
	start := time.Now()
	hashPartitioner := NewHashBasedPartitioner(partitions)
	for i := 0; i < numRecords; i++ {
		record := ShardedRecord{
			ID:       fmt.Sprintf("record%d", i),
			Data:     fmt.Sprintf("Data %d", i),
			ShardKey: fmt.Sprintf("key%d", i),
		}
		hashPartitioner.AddRecord(record)
	}
	hashDuration := time.Since(start)

	log.Printf("Hash-based partitioning: %d records in %v (%.2f records/sec)",
		numRecords, hashDuration, float64(numRecords)/hashDuration.Seconds())

	// Test consistent hashing performance
	start = time.Now()
	consistentRing := NewConsistentHashRing(100)
	for _, partition := range partitions {
		consistentRing.AddNode(partition)
	}

	for i := 0; i < numRecords; i++ {
		key := fmt.Sprintf("key%d", i)
		consistentRing.GetNode(key)
	}
	consistentDuration := time.Since(start)

	log.Printf("Consistent hashing: %d lookups in %v (%.2f lookups/sec)",
		numRecords, consistentDuration, float64(numRecords)/consistentDuration.Seconds())

	// Distribution analysis
	log.Println("\n--- Distribution Analysis ---")
	distribution := make(map[int]int)
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key%d", i)
		partitionID := hashPartitioner.GetPartition(key)
		distribution[partitionID]++
	}

	log.Printf("Hash-based distribution: %+v", distribution)

	// Calculate standard deviation to measure uniformity
	mean := float64(numRecords) / float64(len(partitions))
	variance := 0.0
	for _, count := range distribution {
		variance += math.Pow(float64(count)-mean, 2)
	}
	variance /= float64(len(partitions))
	stddev := math.Sqrt(variance)

	log.Printf("Distribution std deviation: %.2f (lower is better)", stddev)
}

// FAANG Interview Discussion Points:

// 1. Horizontal vs Vertical Partitioning:
//    - Horizontal: Split rows across shards (sharding)
//    - Vertical: Split columns across different stores
//    - Trade-offs in query complexity and scalability

// 2. Partitioning Strategies:
//    - Hash-based: Even distribution, simple routing
//    - Range-based: Good for range queries, hotspot risk
//    - Directory-based: Flexible but adds complexity
//    - Consistent hashing: Good for dynamic scaling

// 3. Consistent Hashing Benefits:
//    - Minimal data movement on node addition/removal
//    - Virtual nodes for better load distribution
//    - Replication factor control

// 4. Cross-shard Query Challenges:
//    - Joins across shards are expensive
//    - Distributed transactions complexity
//    - Eventual consistency considerations
//    - Query optimization strategies

// 5. Rebalancing Strategies:
//    - Live migration vs offline migration
//    - Data consistency during rebalancing
//    - Monitoring partition hotspots
//    - Automated vs manual rebalancing

// Real-world Applications:
// - MongoDB: Range and hash-based sharding
// - Cassandra: Consistent hashing with virtual nodes
// - PostgreSQL: Table partitioning (range, hash, list)
// - Redis Cluster: Hash slots for consistent sharding

// Interview Questions to Discuss:
// 1. How do you choose between different partitioning strategies?
// 2. How do you handle hotspots in hash-based partitioning?
// 3. What are the challenges of cross-shard transactions?
// 4. How do you rebalance shards without downtime?
// 5. How do you ensure data consistency during partition splits?
