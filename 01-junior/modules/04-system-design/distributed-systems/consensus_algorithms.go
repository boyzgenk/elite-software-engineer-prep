package distributed_systems

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Raft Consensus Algorithm Implementation
// Used by etcd, Consul, and many distributed systems

// RaftState represents the state of a Raft node
type RaftState int

const (
	Follower RaftState = iota
	Candidate
	Leader
)

// LogEntry represents an entry in the Raft log
type LogEntry struct {
	Term    uint64
	Index   uint64
	Command interface{}
	Type    string
}

// RaftNode represents a single node in a Raft cluster
type RaftNode struct {
	id          string
	state       RaftState
	currentTerm uint64
	votedFor    string
	log         []LogEntry
	commitIndex uint64
	lastApplied uint64

	// Leader state
	nextIndex  map[string]uint64
	matchIndex map[string]uint64

	// Configuration
	peers    []string
	majority int

	// Channels and control
	voteCh      chan VoteRequest
	appendCh    chan AppendEntriesRequest
	heartbeatCh chan bool
	leaderCh    chan bool
	shutdown    chan bool

	mu           sync.RWMutex
	lastActivity time.Time
}

// VoteRequest represents a RequestVote RPC
type VoteRequest struct {
	Term         uint64
	CandidateID  string
	LastLogIndex uint64
	LastLogTerm  uint64
	ResponseCh   chan VoteResponse
}

// VoteResponse represents a RequestVote response
type VoteResponse struct {
	Term        uint64
	VoteGranted bool
	NodeID      string
}

// AppendEntriesRequest represents an AppendEntries RPC
type AppendEntriesRequest struct {
	Term         uint64
	LeaderID     string
	PrevLogIndex uint64
	PrevLogTerm  uint64
	Entries      []LogEntry
	LeaderCommit uint64
	ResponseCh   chan AppendEntriesResponse
}

// AppendEntriesResponse represents an AppendEntries response
type AppendEntriesResponse struct {
	Term      uint64
	Success   bool
	NodeID    string
	LastIndex uint64
}

// NewRaftNode creates a new Raft node
func NewRaftNode(id string, peers []string) *RaftNode {
	node := &RaftNode{
		id:           id,
		state:        Follower,
		currentTerm:  0,
		votedFor:     "",
		log:          make([]LogEntry, 0),
		commitIndex:  0,
		lastApplied:  0,
		nextIndex:    make(map[string]uint64),
		matchIndex:   make(map[string]uint64),
		peers:        peers,
		majority:     len(peers)/2 + 1,
		voteCh:       make(chan VoteRequest, 10),
		appendCh:     make(chan AppendEntriesRequest, 10),
		heartbeatCh:  make(chan bool, 10),
		leaderCh:     make(chan bool, 10),
		shutdown:     make(chan bool),
		lastActivity: time.Now(),
	}

	// Initialize leader state
	for _, peer := range peers {
		if peer != id {
			node.nextIndex[peer] = 1
			node.matchIndex[peer] = 0
		}
	}

	return node
}

// Start begins the Raft consensus protocol
func (rn *RaftNode) Start() {
	go rn.runStateMachine()
}

// Stop gracefully stops the node
func (rn *RaftNode) Stop() {
	close(rn.shutdown)
}

// runStateMachine runs the main Raft state machine
func (rn *RaftNode) runStateMachine() {
	for {
		select {
		case <-rn.shutdown:
			return
		default:
			rn.mu.RLock()
			state := rn.state
			rn.mu.RUnlock()

			switch state {
			case Follower:
				rn.runFollower()
			case Candidate:
				rn.runCandidate()
			case Leader:
				rn.runLeader()
			}
		}
	}
}

// runFollower implements follower behavior
func (rn *RaftNode) runFollower() {
	electionTimeout := time.Duration(150+rand.Intn(150)) * time.Millisecond
	timer := time.NewTimer(electionTimeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// Election timeout - become candidate
			rn.mu.Lock()
			rn.state = Candidate
			rn.mu.Unlock()
			return

		case vote := <-rn.voteCh:
			rn.handleVoteRequest(vote)
			timer.Reset(electionTimeout)

		case append := <-rn.appendCh:
			rn.handleAppendEntries(append)
			timer.Reset(electionTimeout)

		case <-rn.shutdown:
			return
		}
	}
}

// runCandidate implements candidate behavior (leader election)
func (rn *RaftNode) runCandidate() {
	rn.mu.Lock()
	rn.currentTerm++
	rn.votedFor = rn.id
	currentTerm := rn.currentTerm
	rn.mu.Unlock()

	votes := 1 // Vote for self
	votesNeeded := rn.majority

	// Send vote requests to all peers
	for _, peer := range rn.peers {
		if peer != rn.id {
			go rn.sendVoteRequest(peer, currentTerm)
		}
	}

	electionTimeout := time.Duration(150+rand.Intn(150)) * time.Millisecond
	timer := time.NewTimer(electionTimeout)
	defer timer.Stop()

	for {
		select {
		case vote := <-rn.voteCh:
			response := rn.handleVoteRequest(vote)
			if response.Term > currentTerm {
				// Higher term discovered, become follower
				rn.mu.Lock()
				rn.currentTerm = response.Term
				rn.state = Follower
				rn.votedFor = ""
				rn.mu.Unlock()
				return
			}

		case append := <-rn.appendCh:
			response := rn.handleAppendEntries(append)
			if response.Success {
				// Valid leader found, become follower
				rn.mu.Lock()
				rn.state = Follower
				rn.mu.Unlock()
				return
			}

		case <-timer.C:
			// Election timeout - start new election
			return

		case <-rn.shutdown:
			return
		}

		// Check if we have majority votes
		if votes >= votesNeeded {
			rn.mu.Lock()
			rn.state = Leader
			rn.mu.Unlock()
			rn.initializeLeaderState()
			return
		}
	}
}

// runLeader implements leader behavior
func (rn *RaftNode) runLeader() {
	// Send initial heartbeats
	rn.sendHeartbeats()

	heartbeatInterval := 50 * time.Millisecond
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rn.sendHeartbeats()

		case vote := <-rn.voteCh:
			response := rn.handleVoteRequest(vote)
			if response.Term > rn.currentTerm {
				// Higher term discovered, step down
				rn.mu.Lock()
				rn.currentTerm = response.Term
				rn.state = Follower
				rn.votedFor = ""
				rn.mu.Unlock()
				return
			}

		case append := <-rn.appendCh:
			rn.handleAppendEntries(append)

		case <-rn.shutdown:
			return
		}
	}
}

// handleVoteRequest processes vote requests
func (rn *RaftNode) handleVoteRequest(req VoteRequest) VoteResponse {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	response := VoteResponse{
		Term:        rn.currentTerm,
		VoteGranted: false,
		NodeID:      rn.id,
	}

	// Update term if higher
	if req.Term > rn.currentTerm {
		rn.currentTerm = req.Term
		rn.votedFor = ""
		rn.state = Follower
	}

	// Grant vote if conditions are met
	if req.Term >= rn.currentTerm &&
		(rn.votedFor == "" || rn.votedFor == req.CandidateID) &&
		rn.isLogUpToDate(req.LastLogIndex, req.LastLogTerm) {

		response.VoteGranted = true
		rn.votedFor = req.CandidateID
		rn.lastActivity = time.Now()
	}

	response.Term = rn.currentTerm

	// Send response
	go func() {
		select {
		case req.ResponseCh <- response:
		case <-time.After(100 * time.Millisecond):
		}
	}()

	return response
}

// handleAppendEntries processes append entries requests
func (rn *RaftNode) handleAppendEntries(req AppendEntriesRequest) AppendEntriesResponse {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	response := AppendEntriesResponse{
		Term:      rn.currentTerm,
		Success:   false,
		NodeID:    rn.id,
		LastIndex: rn.getLastLogIndex(),
	}

	// Update term if higher
	if req.Term > rn.currentTerm {
		rn.currentTerm = req.Term
		rn.votedFor = ""
		rn.state = Follower
	}

	// Reject if term is lower
	if req.Term < rn.currentTerm {
		response.Term = rn.currentTerm
		go func() {
			select {
			case req.ResponseCh <- response:
			case <-time.After(100 * time.Millisecond):
			}
		}()
		return response
	}

	// Reset election timeout
	rn.lastActivity = time.Now()

	// Check log consistency
	if req.PrevLogIndex > 0 {
		if req.PrevLogIndex > rn.getLastLogIndex() {
			// Log too short
			go func() {
				select {
				case req.ResponseCh <- response:
				case <-time.After(100 * time.Millisecond):
				}
			}()
			return response
		}

		if rn.log[req.PrevLogIndex-1].Term != req.PrevLogTerm {
			// Log inconsistency
			go func() {
				select {
				case req.ResponseCh <- response:
				case <-time.After(100 * time.Millisecond):
				}
			}()
			return response
		}
	}

	// Append new entries
	if len(req.Entries) > 0 {
		// Remove conflicting entries
		rn.log = rn.log[:req.PrevLogIndex]
		rn.log = append(rn.log, req.Entries...)
	}

	// Update commit index
	if req.LeaderCommit > rn.commitIndex {
		rn.commitIndex = min(req.LeaderCommit, rn.getLastLogIndex())
	}

	response.Success = true
	response.Term = rn.currentTerm
	response.LastIndex = rn.getLastLogIndex()

	go func() {
		select {
		case req.ResponseCh <- response:
		case <-time.After(100 * time.Millisecond):
		}
	}()

	return response
}

// sendVoteRequest sends a vote request to a peer
func (rn *RaftNode) sendVoteRequest(peer string, term uint64) {
	rn.mu.RLock()
	lastLogIndex := rn.getLastLogIndex()
	lastLogTerm := uint64(0)
	if lastLogIndex > 0 {
		lastLogTerm = rn.log[lastLogIndex-1].Term
	}
	rn.mu.RUnlock()

	_ = VoteRequest{
		Term:         term,
		CandidateID:  rn.id,
		LastLogIndex: lastLogIndex,
		LastLogTerm:  lastLogTerm,
		ResponseCh:   make(chan VoteResponse, 1),
	}

	// Simulate network call
	go func() {
		// In real implementation, this would be an RPC call
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
		// Simulate response - in real implementation this comes from the peer
	}()
}

// sendHeartbeats sends heartbeat messages to all followers
func (rn *RaftNode) sendHeartbeats() {
	rn.mu.RLock()
	currentTerm := rn.currentTerm
	rn.mu.RUnlock()

	for _, peer := range rn.peers {
		if peer != rn.id {
			go rn.sendAppendEntries(peer, currentTerm, nil)
		}
	}
}

// sendAppendEntries sends append entries to a peer
func (rn *RaftNode) sendAppendEntries(peer string, term uint64, entries []LogEntry) {
	rn.mu.RLock()
	prevLogIndex := rn.nextIndex[peer] - 1
	prevLogTerm := uint64(0)
	if prevLogIndex > 0 && prevLogIndex <= uint64(len(rn.log)) {
		prevLogTerm = rn.log[prevLogIndex-1].Term
	}
	commitIndex := rn.commitIndex
	rn.mu.RUnlock()

	_ = AppendEntriesRequest{
		Term:         term,
		LeaderID:     rn.id,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
		LeaderCommit: commitIndex,
		ResponseCh:   make(chan AppendEntriesResponse, 1),
	}

	// Simulate network call
	go func() {
		// In real implementation, this would be an RPC call
		time.Sleep(time.Duration(rand.Intn(25)) * time.Millisecond)
		// Handle response would update nextIndex and matchIndex
	}()
}

// Append submits a command to the Raft cluster
func (rn *RaftNode) Append(command interface{}) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	if rn.state != Leader {
		return fmt.Errorf("not leader")
	}

	// Add entry to log
	entry := LogEntry{
		Term:    rn.currentTerm,
		Index:   rn.getLastLogIndex() + 1,
		Command: command,
		Type:    "command",
	}

	rn.log = append(rn.log, entry)

	// Replicate to followers
	go rn.replicateEntry(entry)

	return nil
}

// replicateEntry replicates an entry to all followers
func (rn *RaftNode) replicateEntry(entry LogEntry) {
	for _, peer := range rn.peers {
		if peer != rn.id {
			go rn.sendAppendEntries(peer, entry.Term, []LogEntry{entry})
		}
	}
}

// Helper methods
func (rn *RaftNode) getLastLogIndex() uint64 {
	if len(rn.log) == 0 {
		return 0
	}
	return rn.log[len(rn.log)-1].Index
}

func (rn *RaftNode) getLastLogTerm() uint64 {
	if len(rn.log) == 0 {
		return 0
	}
	return rn.log[len(rn.log)-1].Term
}

func (rn *RaftNode) isLogUpToDate(lastLogIndex, lastLogTerm uint64) bool {
	myLastTerm := rn.getLastLogTerm()
	myLastIndex := rn.getLastLogIndex()

	if lastLogTerm != myLastTerm {
		return lastLogTerm > myLastTerm
	}
	return lastLogIndex >= myLastIndex
}

func (rn *RaftNode) initializeLeaderState() {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	lastIndex := rn.getLastLogIndex()
	for _, peer := range rn.peers {
		if peer != rn.id {
			rn.nextIndex[peer] = lastIndex + 1
			rn.matchIndex[peer] = 0
		}
	}
}

// Byzantine Fault Tolerance (BFT) Implementation
// Handles malicious nodes in addition to crash failures

// BFTNode represents a Byzantine fault tolerant node
type BFTNode struct {
	id          string
	view        uint64
	sequenceNum uint64

	// Message stores
	prepares map[string]map[string]*PrepareMsg
	commits  map[string]map[string]*CommitMsg

	// Request tracking
	requests map[string]*Request
	executed map[string]bool

	// Configuration
	nodes      []string
	faultCount int // f nodes can be faulty
	threshold  int // 2f + 1 threshold

	mu sync.RWMutex
}

// BFT message types
type Request struct {
	ID        string
	Operation string
	Timestamp time.Time
	ClientID  string
}

type PrepareMsg struct {
	View        uint64
	SequenceNum uint64
	RequestID   string
	NodeID      string
	Signature   string
}

type CommitMsg struct {
	View        uint64
	SequenceNum uint64
	RequestID   string
	NodeID      string
	Signature   string
}

// NewBFTNode creates a new Byzantine fault tolerant node
func NewBFTNode(id string, nodes []string, faultCount int) *BFTNode {
	return &BFTNode{
		id:          id,
		view:        0,
		sequenceNum: 0,
		prepares:    make(map[string]map[string]*PrepareMsg),
		commits:     make(map[string]map[string]*CommitMsg),
		requests:    make(map[string]*Request),
		executed:    make(map[string]bool),
		nodes:       nodes,
		faultCount:  faultCount,
		threshold:   2*faultCount + 1,
	}
}

// SubmitRequest submits a request to the BFT system
func (bn *BFTNode) SubmitRequest(req *Request) error {
	bn.mu.Lock()
	defer bn.mu.Unlock()

	// Store request
	bn.requests[req.ID] = req

	// Start consensus process
	bn.sequenceNum++

	// Send prepare messages
	prepare := &PrepareMsg{
		View:        bn.view,
		SequenceNum: bn.sequenceNum,
		RequestID:   req.ID,
		NodeID:      bn.id,
		Signature:   bn.signMessage(req.ID),
	}

	bn.broadcastPrepare(prepare)
	return nil
}

// HandlePrepare processes prepare messages
func (bn *BFTNode) HandlePrepare(prepare *PrepareMsg) {
	bn.mu.Lock()
	defer bn.mu.Unlock()

	// Validate message
	if !bn.validatePrepare(prepare) {
		return
	}

	// Store prepare message
	if bn.prepares[prepare.RequestID] == nil {
		bn.prepares[prepare.RequestID] = make(map[string]*PrepareMsg)
	}
	bn.prepares[prepare.RequestID][prepare.NodeID] = prepare

	// Check if we have enough prepares
	if len(bn.prepares[prepare.RequestID]) >= bn.threshold {
		bn.sendCommit(prepare)
	}
}

// HandleCommit processes commit messages
func (bn *BFTNode) HandleCommit(commit *CommitMsg) {
	bn.mu.Lock()
	defer bn.mu.Unlock()

	// Validate message
	if !bn.validateCommit(commit) {
		return
	}

	// Store commit message
	if bn.commits[commit.RequestID] == nil {
		bn.commits[commit.RequestID] = make(map[string]*CommitMsg)
	}
	bn.commits[commit.RequestID][commit.NodeID] = commit

	// Check if we have enough commits
	if len(bn.commits[commit.RequestID]) >= bn.threshold {
		bn.executeRequest(commit.RequestID)
	}
}

// Quorum-based Consensus System
// Used in systems like Cassandra, DynamoDB

// QuorumSystem implements quorum-based consensus
type QuorumSystem struct {
	nodes       []*QuorumNode
	replication int
	readQuorum  int
	writeQuorum int
}

// QuorumNode represents a node in quorum system
type QuorumNode struct {
	id        string
	data      map[string]*VersionedValue
	available bool
	mu        sync.RWMutex
}

// VersionedValue represents a value with version information
type VersionedValue struct {
	Value     interface{}
	Version   uint64
	Timestamp time.Time
	NodeID    string
}

// NewQuorumSystem creates a new quorum-based system
func NewQuorumSystem(nodeCount, replication int) *QuorumSystem {
	nodes := make([]*QuorumNode, nodeCount)
	for i := 0; i < nodeCount; i++ {
		nodes[i] = &QuorumNode{
			id:        fmt.Sprintf("node-%d", i),
			data:      make(map[string]*VersionedValue),
			available: true,
		}
	}

	return &QuorumSystem{
		nodes:       nodes,
		replication: replication,
		readQuorum:  replication/2 + 1,
		writeQuorum: replication/2 + 1,
	}
}

// QuorumWrite performs a quorum write
func (qs *QuorumSystem) QuorumWrite(key string, value interface{}) error {
	// Select nodes for replication
	nodes := qs.selectNodes(qs.replication)

	version := uint64(time.Now().UnixNano())
	versionedValue := &VersionedValue{
		Value:     value,
		Version:   version,
		Timestamp: time.Now(),
		NodeID:    nodes[0].id,
	}

	successCount := int64(0)
	var wg sync.WaitGroup

	// Write to selected nodes
	for _, node := range nodes {
		wg.Add(1)
		go func(n *QuorumNode) {
			defer wg.Done()
			if qs.writeToNode(n, key, versionedValue) {
				atomic.AddInt64(&successCount, 1)
			}
		}(node)
	}

	wg.Wait()

	if int(successCount) >= qs.writeQuorum {
		return nil
	}
	return fmt.Errorf("failed to achieve write quorum: %d/%d", successCount, qs.writeQuorum)
}

// QuorumRead performs a quorum read
func (qs *QuorumSystem) QuorumRead(key string) (interface{}, error) {
	nodes := qs.selectNodes(qs.readQuorum)

	values := make([]*VersionedValue, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Read from quorum nodes
	for _, node := range nodes {
		wg.Add(1)
		go func(n *QuorumNode) {
			defer wg.Done()
			if value := qs.readFromNode(n, key); value != nil {
				mu.Lock()
				values = append(values, value)
				mu.Unlock()
			}
		}(node)
	}

	wg.Wait()

	if len(values) >= qs.readQuorum {
		// Return most recent version
		var latest *VersionedValue
		for _, v := range values {
			if latest == nil || v.Version > latest.Version {
				latest = v
			}
		}

		// Trigger read repair if needed
		go qs.readRepair(key, values)

		return latest.Value, nil
	}

	return nil, fmt.Errorf("failed to achieve read quorum")
}

// Helper methods for BFT
func (bn *BFTNode) validatePrepare(prepare *PrepareMsg) bool {
	// Validate signature, view, sequence number, etc.
	return true
}

func (bn *BFTNode) validateCommit(commit *CommitMsg) bool {
	// Validate signature, view, sequence number, etc.
	return true
}

func (bn *BFTNode) signMessage(message string) string {
	// In real implementation, use cryptographic signature
	return fmt.Sprintf("sig_%s_%s", bn.id, message)
}

func (bn *BFTNode) broadcastPrepare(prepare *PrepareMsg) {
	// Broadcast to all nodes
	for _, nodeID := range bn.nodes {
		if nodeID != bn.id {
			// Send prepare message
		}
	}
}

func (bn *BFTNode) sendCommit(prepare *PrepareMsg) {
	_ = &CommitMsg{
		View:        prepare.View,
		SequenceNum: prepare.SequenceNum,
		RequestID:   prepare.RequestID,
		NodeID:      bn.id,
		Signature:   bn.signMessage(prepare.RequestID + "_commit"),
	}

	// Broadcast commit
	for _, nodeID := range bn.nodes {
		if nodeID != bn.id {
			// Send commit message
		}
	}
}

func (bn *BFTNode) executeRequest(requestID string) {
	if bn.executed[requestID] {
		return
	}

	req := bn.requests[requestID]
	if req != nil {
		// Execute the request
		bn.executed[requestID] = true
	}
}

// Helper methods for Quorum
func (qs *QuorumSystem) selectNodes(count int) []*QuorumNode {
	available := make([]*QuorumNode, 0)
	for _, node := range qs.nodes {
		if node.available {
			available = append(available, node)
		}
	}

	if len(available) < count {
		return available
	}

	// Shuffle and select
	rand.Shuffle(len(available), func(i, j int) {
		available[i], available[j] = available[j], available[i]
	})

	return available[:count]
}

func (qs *QuorumSystem) writeToNode(node *QuorumNode, key string, value *VersionedValue) bool {
	node.mu.Lock()
	defer node.mu.Unlock()

	if !node.available {
		return false
	}

	node.data[key] = value
	return true
}

func (qs *QuorumSystem) readFromNode(node *QuorumNode, key string) *VersionedValue {
	node.mu.RLock()
	defer node.mu.RUnlock()

	if !node.available {
		return nil
	}

	return node.data[key]
}

func (qs *QuorumSystem) readRepair(key string, values []*VersionedValue) {
	// Find the most recent version
	var latest *VersionedValue
	for _, v := range values {
		if latest == nil || v.Version > latest.Version {
			latest = v
		}
	}

	// Update nodes with stale data
	for _, node := range qs.nodes {
		current := qs.readFromNode(node, key)
		if current != nil && current.Version < latest.Version {
			qs.writeToNode(node, key, latest)
		}
	}
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// Interview Discussion Points:
//
// 1. Raft vs Paxos:
//    - Raft: Simpler to understand, leader-based, strong leader
//    - Paxos: More flexible, multiple proposers, complex to implement
//    - Multi-Paxos: More practical variant used in production
//
// 2. Leader Election:
//    - Randomized timeouts prevent split votes
//    - Majority voting ensures single leader
//    - Term numbers prevent stale leaders
//
// 3. Log Replication:
//    - Sequential consistency through log ordering
//    - Commit only after majority acknowledgment
//    - Read operations can go to leader or followers (depending on consistency needs)
//
// 4. Byzantine Fault Tolerance:
//    - Handles malicious/corrupted nodes
//    - Requires 3f+1 nodes to tolerate f Byzantine faults
//    - Practical BFT (pBFT) used in blockchain systems
//
// 5. Split-brain Prevention:
//    - Majority quorums prevent multiple leaders
//    - Network partition handling
//    - Witness nodes for even cluster sizes
//
// 6. Real-world Applications:
//    - etcd: Raft for Kubernetes cluster state
//    - Consul: Raft for service discovery
//    - MongoDB: Raft for replica set leader election
//    - CockroachDB: Raft for range consensus
//    - Blockchain: Various consensus (PoW, PoS, pBFT)
//
// 7. Performance Considerations:
//    - Latency: Round-trips for consensus
//    - Throughput: Batch operations, pipeline
//    - Network partitions: Availability vs consistency trade-offs
//
// 8. Operational Aspects:
//    - Cluster membership changes
//    - Node failure detection and recovery
//    - Log compaction and snapshots
//    - Monitoring consensus health
