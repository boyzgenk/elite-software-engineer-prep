package distributed_systems

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// CAP Theorem: Consistency, Availability, Partition tolerance
// You can only guarantee two out of three properties in a distributed system

// ConsistencyLevel defines different consistency guarantees
type ConsistencyLevel int

const (
	StrongConsistency ConsistencyLevel = iota
	EventualConsistency
	CausalConsistency
	SessionConsistency
	MonotonicReadConsistency
	MonotonicWriteConsistency
)

// SystemMode represents CAP theorem trade-offs
type SystemMode int

const (
	CPMode SystemMode = iota // Consistency + Partition tolerance (sacrifice Availability)
	APMode                   // Availability + Partition tolerance (sacrifice Consistency)
	CAMode                   // Consistency + Availability (sacrifice Partition tolerance)
)

// CAPSystem demonstrates CAP theorem trade-offs
type CAPSystem struct {
	nodes            []*Node
	mode             SystemMode
	consistencyLevel ConsistencyLevel
	partitioned      bool
	quorumSize       int
	mu               sync.RWMutex
}

// Node represents a distributed system node
type Node struct {
	id           string
	available    bool
	data         map[string]interface{}
	version      uint64
	vectorClock  map[string]uint64
	partitioned  bool
	mu           sync.RWMutex
	lastActivity time.Time
}

// NewCAPSystem creates a distributed system with specified characteristics
func NewCAPSystem(nodeCount int, mode SystemMode, consistencyLevel ConsistencyLevel) *CAPSystem {
	nodes := make([]*Node, nodeCount)
	for i := 0; i < nodeCount; i++ {
		nodes[i] = &Node{
			id:           fmt.Sprintf("node-%d", i),
			available:    true,
			data:         make(map[string]interface{}),
			version:      0,
			vectorClock:  make(map[string]uint64),
			partitioned:  false,
			lastActivity: time.Now(),
		}
	}

	return &CAPSystem{
		nodes:            nodes,
		mode:             mode,
		consistencyLevel: consistencyLevel,
		partitioned:      false,
		quorumSize:       (nodeCount / 2) + 1,
	}
}

// Write demonstrates different write behaviors based on CAP trade-offs
func (cs *CAPSystem) Write(key string, value interface{}) error {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	switch cs.mode {
	case CPMode:
		return cs.cpWrite(key, value)
	case APMode:
		return cs.apWrite(key, value)
	case CAMode:
		return cs.caWrite(key, value)
	default:
		return fmt.Errorf("unknown system mode")
	}
}

// cpWrite: Consistency + Partition tolerance (CP system like etcd, Consul)
func (cs *CAPSystem) cpWrite(key string, value interface{}) error {
	if cs.partitioned {
		availableNodes := cs.getAvailableNodes()
		if len(availableNodes) < cs.quorumSize {
			return fmt.Errorf("insufficient nodes for quorum: need %d, have %d",
				cs.quorumSize, len(availableNodes))
		}
	}

	// Require majority consensus for writes
	return cs.writeWithQuorum(key, value)
}

// apWrite: Availability + Partition tolerance (AP system like Cassandra, DynamoDB)
func (cs *CAPSystem) apWrite(key string, value interface{}) error {
	// Write to any available node, handle conflicts later
	availableNodes := cs.getAvailableNodes()
	if len(availableNodes) == 0 {
		return fmt.Errorf("no available nodes")
	}

	// Write to at least one node (best effort)
	node := availableNodes[0]
	return cs.writeToNode(node, key, value)
}

// caWrite: Consistency + Availability (CA system like traditional RDBMS)
func (cs *CAPSystem) caWrite(key string, value interface{}) error {
	if cs.partitioned {
		return fmt.Errorf("network partition detected, cannot maintain CA guarantees")
	}

	// All nodes must be available and consistent
	return cs.writeToAllNodes(key, value)
}

// Read demonstrates different read behaviors based on consistency levels
func (cs *CAPSystem) Read(key string) (interface{}, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	switch cs.consistencyLevel {
	case StrongConsistency:
		return cs.strongConsistentRead(key)
	case EventualConsistency:
		return cs.eventuallyConsistentRead(key)
	case CausalConsistency:
		return cs.causalConsistentRead(key)
	case SessionConsistency:
		return cs.sessionConsistentRead(key)
	default:
		return cs.eventuallyConsistentRead(key)
	}
}

// strongConsistentRead requires quorum read
func (cs *CAPSystem) strongConsistentRead(key string) (interface{}, error) {
	availableNodes := cs.getAvailableNodes()
	if len(availableNodes) < cs.quorumSize {
		return nil, fmt.Errorf("insufficient nodes for quorum read")
	}

	// Read from majority of nodes
	values := make(map[interface{}]int)
	versions := make(map[interface{}]uint64)

	for i := 0; i < cs.quorumSize; i++ {
		node := availableNodes[i]
		node.mu.RLock()
		value, exists := node.data[key]
		if exists {
			values[value]++
			versions[value] = node.version
		}
		node.mu.RUnlock()
	}

	// Return the value with highest version
	var latestValue interface{}
	var latestVersion uint64
	for value, version := range versions {
		if version > latestVersion {
			latestVersion = version
			latestValue = value
		}
	}

	return latestValue, nil
}

// eventuallyConsistentRead reads from any available node
func (cs *CAPSystem) eventuallyConsistentRead(key string) (interface{}, error) {
	availableNodes := cs.getAvailableNodes()
	if len(availableNodes) == 0 {
		return nil, fmt.Errorf("no available nodes")
	}

	// Read from first available node
	node := availableNodes[rand.Intn(len(availableNodes))]
	node.mu.RLock()
	defer node.mu.RUnlock()

	value, exists := node.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found")
	}
	return value, nil
}

// causalConsistentRead ensures causally related reads are consistent
func (cs *CAPSystem) causalConsistentRead(key string) (interface{}, error) {
	// Implementation would track causal dependencies using vector clocks
	return cs.eventuallyConsistentRead(key)
}

// sessionConsistentRead ensures read-your-writes consistency
func (cs *CAPSystem) sessionConsistentRead(key string) (interface{}, error) {
	// Implementation would track session state
	return cs.eventuallyConsistentRead(key)
}

// SimulateNetworkPartition creates a network partition
func (cs *CAPSystem) SimulateNetworkPartition(nodeIndices []int) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.partitioned = true
	for _, idx := range nodeIndices {
		if idx < len(cs.nodes) {
			cs.nodes[idx].partitioned = true
		}
	}
}

// HealNetworkPartition removes network partition
func (cs *CAPSystem) HealNetworkPartition() {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.partitioned = false
	for _, node := range cs.nodes {
		node.partitioned = false
	}

	// Trigger read repair and anti-entropy
	go cs.readRepair()
}

// writeWithQuorum writes to majority of nodes
func (cs *CAPSystem) writeWithQuorum(key string, value interface{}) error {
	availableNodes := cs.getAvailableNodes()
	if len(availableNodes) < cs.quorumSize {
		return fmt.Errorf("insufficient nodes for quorum")
	}

	successCount := 0
	for i := 0; i < cs.quorumSize; i++ {
		if err := cs.writeToNode(availableNodes[i], key, value); err == nil {
			successCount++
		}
	}

	if successCount >= cs.quorumSize {
		return nil
	}
	return fmt.Errorf("failed to achieve quorum write")
}

// writeToNode writes to a single node
func (cs *CAPSystem) writeToNode(node *Node, key string, value interface{}) error {
	node.mu.Lock()
	defer node.mu.Unlock()

	if !node.available || node.partitioned {
		return fmt.Errorf("node %s unavailable", node.id)
	}

	node.data[key] = value
	node.version++
	node.lastActivity = time.Now()

	// Update vector clock
	node.vectorClock[node.id]++

	return nil
}

// writeToAllNodes writes to all available nodes
func (cs *CAPSystem) writeToAllNodes(key string, value interface{}) error {
	availableNodes := cs.getAvailableNodes()
	if len(availableNodes) != len(cs.nodes) {
		return fmt.Errorf("not all nodes available")
	}

	for _, node := range availableNodes {
		if err := cs.writeToNode(node, key, value); err != nil {
			return err
		}
	}
	return nil
}

// getAvailableNodes returns non-partitioned, available nodes
func (cs *CAPSystem) getAvailableNodes() []*Node {
	var available []*Node
	for _, node := range cs.nodes {
		if node.available && !node.partitioned {
			available = append(available, node)
		}
	}
	return available
}

// readRepair performs background consistency repair
func (cs *CAPSystem) readRepair() {
	// Anti-entropy process to sync nodes after partition healing
	for key := range cs.getAllKeys() {
		cs.repairKey(key)
	}
}

// repairKey repairs inconsistencies for a specific key
func (cs *CAPSystem) repairKey(key string) {
	nodeValues := make(map[string]interface{})
	nodeVersions := make(map[string]uint64)

	// Collect values from all nodes
	for _, node := range cs.nodes {
		node.mu.RLock()
		if value, exists := node.data[key]; exists {
			nodeValues[node.id] = value
			nodeVersions[node.id] = node.version
		}
		node.mu.RUnlock()
	}

	// Find the most recent version
	var latestVersion uint64
	var latestValue interface{}
	for nodeID, version := range nodeVersions {
		if version > latestVersion {
			latestVersion = version
			latestValue = nodeValues[nodeID]
		}
	}

	// Update nodes with stale data
	for _, node := range cs.nodes {
		node.mu.Lock()
		if currentVersion, exists := nodeVersions[node.id]; exists {
			if currentVersion < latestVersion {
				node.data[key] = latestValue
				node.version = latestVersion
			}
		}
		node.mu.Unlock()
	}
}

// getAllKeys returns all keys across all nodes
func (cs *CAPSystem) getAllKeys() map[string]bool {
	keys := make(map[string]bool)
	for _, node := range cs.nodes {
		node.mu.RLock()
		for key := range node.data {
			keys[key] = true
		}
		node.mu.RUnlock()
	}
	return keys
}

// GetSystemStatus returns current system status and metrics
func (cs *CAPSystem) GetSystemStatus() map[string]interface{} {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	availableNodes := len(cs.getAvailableNodes())
	totalNodes := len(cs.nodes)

	return map[string]interface{}{
		"mode":              cs.mode,
		"consistency_level": cs.consistencyLevel,
		"partitioned":       cs.partitioned,
		"available_nodes":   availableNodes,
		"total_nodes":       totalNodes,
		"quorum_size":       cs.quorumSize,
		"availability":      float64(availableNodes) / float64(totalNodes),
	}
}

// CAPTheoryDemo demonstrates different CAP theorem scenarios
type CAPTheoryDemo struct {
	systems map[string]*CAPSystem
}

// NewCAPTheoryDemo creates demo systems for different CAP trade-offs
func NewCAPTheoryDemo() *CAPTheoryDemo {
	return &CAPTheoryDemo{
		systems: map[string]*CAPSystem{
			"etcd":      NewCAPSystem(5, CPMode, StrongConsistency),
			"cassandra": NewCAPSystem(5, APMode, EventualConsistency),
			"mysql":     NewCAPSystem(3, CAMode, StrongConsistency),
		},
	}
}

// DemonstrateScenario shows how different systems handle specific scenarios
func (demo *CAPTheoryDemo) DemonstrateScenario(scenario string) {
	switch scenario {
	case "network_partition":
		demo.networkPartitionScenario()
	case "node_failure":
		demo.nodeFailureScenario()
	case "consistency_levels":
		demo.consistencyLevelDemo()
	}
}

// networkPartitionScenario demonstrates behavior during network partitions
func (demo *CAPTheoryDemo) networkPartitionScenario() {
	fmt.Println("=== Network Partition Scenario ===")

	for name, system := range demo.systems {
		fmt.Printf("\n%s System (Mode: %v):\n", name, system.mode)

		// Normal operation
		err := system.Write("user:123", "John Doe")
		fmt.Printf("  Normal write: %v\n", err)

		// Simulate partition
		system.SimulateNetworkPartition([]int{3, 4}) // Partition 2 nodes
		err = system.Write("user:124", "Jane Doe")
		fmt.Printf("  Partitioned write: %v\n", err)

		value, err := system.Read("user:123")
		fmt.Printf("  Partitioned read: %v, error: %v\n", value, err)

		// Heal partition
		system.HealNetworkPartition()
		fmt.Printf("  Partition healed\n")
	}
}

// nodeFailureScenario demonstrates behavior during node failures
func (demo *CAPTheoryDemo) nodeFailureScenario() {
	fmt.Println("=== Node Failure Scenario ===")

	for name, system := range demo.systems {
		fmt.Printf("\n%s System:\n", name)

		// Simulate node failure
		system.nodes[0].available = false
		system.nodes[1].available = false

		err := system.Write("order:456", "pending")
		fmt.Printf("  Write with 2 nodes down: %v\n", err)

		status := system.GetSystemStatus()
		fmt.Printf("  Availability: %.2f%%\n", status["availability"].(float64)*100)
	}
}

// consistencyLevelDemo shows different consistency guarantees
func (demo *CAPTheoryDemo) consistencyLevelDemo() {
	fmt.Println("=== Consistency Levels Demo ===")

	system := NewCAPSystem(3, APMode, EventualConsistency)

	// Write to subset of nodes (simulating replication lag)
	system.writeToNode(system.nodes[0], "account:balance", 1000)
	system.writeToNode(system.nodes[1], "account:balance", 950) // Different value

	// Demonstrate different read consistency
	levels := []ConsistencyLevel{
		EventualConsistency,
		StrongConsistency,
	}

	for _, level := range levels {
		system.consistencyLevel = level
		value, err := system.Read("account:balance")
		fmt.Printf("  %v read: %v, error: %v\n", level, value, err)
	}
}

// BenchmarkCAPOperations measures performance implications of CAP choices
func BenchmarkCAPOperations() map[string]time.Duration {
	systems := map[string]*CAPSystem{
		"CP_strong":   NewCAPSystem(5, CPMode, StrongConsistency),
		"AP_eventual": NewCAPSystem(5, APMode, EventualConsistency),
	}

	results := make(map[string]time.Duration)

	for name, system := range systems {
		start := time.Now()

		// Perform batch operations
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key:%d", i)
			system.Write(key, fmt.Sprintf("value:%d", i))
		}

		results[name+"_write"] = time.Since(start)

		start = time.Now()
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key:%d", i)
			system.Read(key)
		}

		results[name+"_read"] = time.Since(start)
	}

	return results
}

// Interview Discussion Points:
//
// 1. CAP Theorem Trade-offs:
//    - CP Systems: etcd, Consul, HBase (sacrifice availability during partitions)
//    - AP Systems: Cassandra, DynamoDB, DNS (sacrifice consistency for availability)
//    - CA Systems: RDBMS in single datacenter (cannot handle partitions)
//
// 2. Real-world Examples:
//    - Google Spanner: Uses synchronized clocks to achieve strong consistency globally
//    - Amazon DynamoDB: Eventually consistent by default, optional strong consistency
//    - MongoDB: CP system that becomes unavailable if can't reach majority
//
// 3. Consistency Models:
//    - Strong: All reads receive the most recent write
//    - Eventual: System will become consistent over time
//    - Causal: Causally related operations are seen in same order
//    - Session: Read-your-writes consistency within a session
//
// 4. Practical Considerations:
//    - Network partitions are rare but real
//    - Most systems choose AP and handle consistency at application level
//    - Quorum-based systems (R + W > N) provide tunable consistency
//    - Multi-region deployments must consider CAP trade-offs carefully
//
// 5. Performance Implications:
//    - Strong consistency requires coordination (slower writes)
//    - Eventual consistency allows local writes (faster, more available)
//    - Read-heavy vs write-heavy workloads affect consistency choices
//    - Network latency significantly impacts distributed consistency protocols
