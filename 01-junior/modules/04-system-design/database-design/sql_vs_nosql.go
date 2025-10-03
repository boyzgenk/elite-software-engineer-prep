package database_design

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// ACID Properties Implementation Examples
// Interview Focus: Understanding transaction guarantees and trade-offs

// ACIDTransaction demonstrates ACID properties in a SQL database
type ACIDTransaction struct {
	db *sql.DB
	mu sync.RWMutex
}

// TransferFunds demonstrates Atomicity and Consistency
// FAANG Interview Point: Banking system, payment processing
func (a *ACIDTransaction) TransferFunds(fromAccount, toAccount string, amount float64) error {
	// Begin transaction to ensure Atomicity
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Check sender balance (Consistency)
	var senderBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ? FOR UPDATE", fromAccount).Scan(&senderBalance)
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %w", err)
	}

	if senderBalance < amount {
		return fmt.Errorf("insufficient funds: balance %.2f, requested %.2f", senderBalance, amount)
	}

	// Debit sender account
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromAccount)
	if err != nil {
		return fmt.Errorf("failed to debit sender: %w", err)
	}

	// Credit receiver account
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toAccount)
	if err != nil {
		return fmt.Errorf("failed to credit receiver: %w", err)
	}

	// Log transaction (Durability - survives system failure)
	_, err = tx.Exec("INSERT INTO transactions (from_account, to_account, amount, timestamp) VALUES (?, ?, ?, ?)",
		fromAccount, toAccount, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction: %w", err)
	}

	// Commit ensures Atomicity and Durability
	return tx.Commit()
}

// IsolationLevels demonstrates different isolation levels
// FAANG Interview Point: Concurrency control, read phenomena
type IsolationLevels struct {
	db *sql.DB
}

func (il *IsolationLevels) DemonstrateReadCommitted() {
	// Prevents dirty reads but allows non-repeatable reads
	tx, _ := il.db.Begin()
	tx.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")

	// This won't see uncommitted changes from other transactions
	var balance float64
	tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", "user123").Scan(&balance)

	// But subsequent reads might see different values if other transactions commit
	time.Sleep(100 * time.Millisecond) // Simulate concurrent transaction
	tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", "user123").Scan(&balance)

	tx.Commit()
}

func (il *IsolationLevels) DemonstrateSerializable() {
	// Highest isolation level - prevents all read phenomena
	tx, _ := il.db.Begin()
	tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")

	// This provides snapshot isolation - consistent view throughout transaction
	var balance float64
	tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", "user123").Scan(&balance)

	// Subsequent reads will see the same value (snapshot consistency)
	time.Sleep(100 * time.Millisecond)
	tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", "user123").Scan(&balance)

	tx.Commit()
}

// CAP Theorem Trade-offs
// Interview Focus: Distributed systems theory and practical implications

// CAPTradeoffs demonstrates CAP theorem in different database scenarios
type CAPTradeoffs struct {
	sqlDB   *sql.DB
	noSQLDB interface{} // Generic NoSQL database interface
}

// ConsistencyAvailability (CA) - Traditional RDBMS approach
// FAANG Interview Point: When network partitions are rare (single datacenter)
func (cap *CAPTradeoffs) DemonstrateCASQLSystem() error {
	// Strong consistency with ACID transactions
	// High availability within single datacenter
	// Cannot handle network partitions between datacenters gracefully

	tx, err := cap.sqlDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Strong consistency - all reads see latest write
	_, err = tx.Exec("UPDATE user_profiles SET last_login = ? WHERE user_id = ?", time.Now(), "user123")
	if err != nil {
		return err
	}

	// Immediate consistency - other transactions will see this update after commit
	return tx.Commit()
}

// ConsistencyPartition (CP) - NoSQL with majority write concern
// FAANG Interview Point: Distributed databases that prioritize consistency
func (cap *CAPTradeoffs) DemonstrateCPNoSQLSystem(ctx context.Context) error {
	// NoSQL database with majority write concern (e.g., MongoDB)
	// Ensures consistency across replica set
	// May become unavailable during network partitions

	// Simulate NoSQL update with majority write concern
	log.Printf("Updating user profile with majority write concern")
	log.Printf("Ensuring update is replicated to majority of replicas before acknowledging")

	// In a real MongoDB implementation:
	// collection := cap.noSQLDB.Collection("user_profiles")
	// update := bson.M{"$set": bson.M{"last_login": time.Now()}}
	// _, err := collection.UpdateOne(ctx, bson.M{"user_id": "user123"}, update,
	//     options.Update().SetWriteConcern(&writeconcern.WriteConcern{W: "majority"}))

	// Simulate consistency check
	time.Sleep(50 * time.Millisecond) // Simulate replication delay
	log.Printf("Update confirmed on majority of replicas")

	return nil
}

// AvailabilityPartition (AP) - Eventually consistent systems
// FAANG Interview Point: Social media feeds, recommendation systems
func (cap *CAPTradeoffs) DemonstrateAPEventualConsistency() {
	// Cassandra-style eventually consistent system
	// High availability even during network partitions
	// May serve stale data temporarily

	// Simulate AP system behavior
	type EventuallyConsistentStore struct {
		data map[string]interface{}
		mu   sync.RWMutex
	}

	store := &EventuallyConsistentStore{
		data: make(map[string]interface{}),
	}

	// Write always succeeds (high availability)
	store.mu.Lock()
	store.data["user123"] = map[string]interface{}{
		"status":    "online",
		"timestamp": time.Now(),
	}
	store.mu.Unlock()

	// Reads might return stale data but system remains available
	store.mu.RLock()
	userData := store.data["user123"]
	store.mu.RUnlock()

	log.Printf("User data (may be stale): %+v", userData)
}

// Use Case Decision Matrix
// Interview Focus: Choosing the right database for specific requirements

// DatabaseSelector provides decision framework for database selection
type DatabaseSelector struct {
	requirements DatabaseRequirements
}

type DatabaseRequirements struct {
	DataStructure    DataStructureType
	ScalePattern     ScalePatternType
	ConsistencyLevel ConsistencyLevelType
	QueryComplexity  QueryComplexityType
	OperationalNeeds OperationalNeedsType
}

type DataStructureType int

const (
	StructuredRelational DataStructureType = iota
	SemiStructuredDocument
	KeyValueSimple
	GraphRelationships
	TimeSeriesEvents
)

type ScalePatternType int

const (
	ReadHeavy ScalePatternType = iota
	WriteHeavy
	BalancedReadWrite
	AnalyticalWorkload
)

type ConsistencyLevelType int

const (
	StrongConsistency ConsistencyLevelType = iota
	EventualConsistency
	SessionConsistency
	TunableConsistency
)

type QueryComplexityType int

const (
	SimpleKeyValue QueryComplexityType = iota
	SimpleFiltering
	ComplexJoins
	FullTextSearch
	GraphTraversal
)

type OperationalNeedsType int

const (
	HighOperationalComplexity OperationalNeedsType = iota
	MediumOperationalComplexity
	LowOperationalComplexity
)

// SelectDatabase provides decision matrix for database selection
// FAANG Interview Point: Systematic approach to technology decisions
func (ds *DatabaseSelector) SelectDatabase() DatabaseRecommendation {
	req := ds.requirements

	// Social Media Platform Example
	if req.DataStructure == SemiStructuredDocument &&
		req.ScalePattern == ReadHeavy &&
		req.ConsistencyLevel == EventualConsistency {
		return DatabaseRecommendation{
			Primary:      "MongoDB",
			Reasoning:    "Document structure fits user profiles, read-heavy for feeds, eventual consistency acceptable for social features",
			Alternatives: []string{"DynamoDB", "Couchbase"},
			Caching:      "Redis for hot user data and timeline cache",
			Analytics:    "Elasticsearch for user search, Kafka + ClickHouse for analytics",
		}
	}

	// Financial Trading System Example
	if req.DataStructure == StructuredRelational &&
		req.ConsistencyLevel == StrongConsistency {
		return DatabaseRecommendation{
			Primary:      "PostgreSQL",
			Reasoning:    "ACID compliance critical for financial data, strong consistency required for trading",
			Alternatives: []string{"MySQL", "Oracle"},
			Caching:      "Redis with write-through pattern",
			Analytics:    "Separate OLAP system (Snowflake/BigQuery) via CDC",
		}
	}

	// Real-time Analytics Example
	if req.DataStructure == TimeSeriesEvents &&
		req.ScalePattern == WriteHeavy {
		return DatabaseRecommendation{
			Primary:      "ClickHouse",
			Reasoning:    "Columnar storage optimized for time-series, handles high write throughput",
			Alternatives: []string{"TimescaleDB", "InfluxDB"},
			Caching:      "Redis for real-time dashboards",
			Analytics:    "Built-in analytical capabilities",
		}
	}

	// Session Store Example
	if req.DataStructure == KeyValueSimple &&
		req.ScalePattern == ReadHeavy {
		return DatabaseRecommendation{
			Primary:      "Redis",
			Reasoning:    "In-memory performance for session data, simple key-value operations",
			Alternatives: []string{"DynamoDB", "Memcached"},
			Caching:      "Primary storage is cache",
			Analytics:    "Export to warehouse for user behavior analysis",
		}
	}

	// Default recommendation
	return DatabaseRecommendation{
		Primary:      "PostgreSQL",
		Reasoning:    "Safe default choice with good performance characteristics and feature completeness",
		Alternatives: []string{"MySQL", "MongoDB"},
		Caching:      "Redis",
		Analytics:    "Separate analytics database via ETL",
	}
}

type DatabaseRecommendation struct {
	Primary      string
	Reasoning    string
	Alternatives []string
	Caching      string
	Analytics    string
}

// Performance Characteristics Analysis
// Interview Focus: Understanding performance implications and bottlenecks

// PerformanceAnalyzer provides performance comparison framework
type PerformanceAnalyzer struct {
	metrics map[string]PerformanceMetrics
}

type PerformanceMetrics struct {
	ReadLatencyP50   time.Duration
	ReadLatencyP99   time.Duration
	WriteLatencyP50  time.Duration
	WriteLatencyP99  time.Duration
	ThroughputReads  int // operations per second
	ThroughputWrites int
	StorageOverhead  float64 // percentage
	MemoryUsage      int64   // bytes
}

// ComparePerformanceProfiles shows typical performance characteristics
// FAANG Interview Point: Performance-driven architecture decisions
func (pa *PerformanceAnalyzer) ComparePerformanceProfiles() map[string]PerformanceMetrics {
	return map[string]PerformanceMetrics{
		"Redis": {
			ReadLatencyP50:   100 * time.Microsecond,
			ReadLatencyP99:   500 * time.Microsecond,
			WriteLatencyP50:  150 * time.Microsecond,
			WriteLatencyP99:  800 * time.Microsecond,
			ThroughputReads:  100000,
			ThroughputWrites: 80000,
			StorageOverhead:  10.0,               // Memory overhead
			MemoryUsage:      1024 * 1024 * 1024, // 1GB
		},
		"PostgreSQL": {
			ReadLatencyP50:   2 * time.Millisecond,
			ReadLatencyP99:   10 * time.Millisecond,
			WriteLatencyP50:  3 * time.Millisecond,
			WriteLatencyP99:  15 * time.Millisecond,
			ThroughputReads:  20000,
			ThroughputWrites: 15000,
			StorageOverhead:  20.0,              // Index overhead
			MemoryUsage:      512 * 1024 * 1024, // 512MB shared buffers
		},
		"MongoDB": {
			ReadLatencyP50:   1 * time.Millisecond,
			ReadLatencyP99:   8 * time.Millisecond,
			WriteLatencyP50:  2 * time.Millisecond,
			WriteLatencyP99:  12 * time.Millisecond,
			ThroughputReads:  25000,
			ThroughputWrites: 18000,
			StorageOverhead:  30.0,               // Document structure overhead
			MemoryUsage:      1024 * 1024 * 1024, // 1GB WiredTiger cache
		},
		"Cassandra": {
			ReadLatencyP50:   3 * time.Millisecond,
			ReadLatencyP99:   20 * time.Millisecond,
			WriteLatencyP50:  1 * time.Millisecond,
			WriteLatencyP99:  5 * time.Millisecond,
			ThroughputReads:  15000,
			ThroughputWrites: 50000, // Write-optimized
			StorageOverhead:  25.0,
			MemoryUsage:      2048 * 1024 * 1024, // 2GB
		},
	}
}

// Data Consistency Models
// Interview Focus: Understanding consistency guarantees in distributed systems

// ConsistencyModel represents different consistency models
type ConsistencyModel interface {
	Read(key string) (interface{}, error)
	Write(key string, value interface{}) error
	GetConsistencyLevel() string
}

// StrongConsistencyModel - Linearizability
// FAANG Interview Point: Bank account balances, inventory systems
type StrongConsistencyModel struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func (sc *StrongConsistencyModel) Read(key string) (interface{}, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// All reads see the most recent write
	// Provides linearizability guarantees
	value, exists := sc.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found")
	}
	return value, nil
}

func (sc *StrongConsistencyModel) Write(key string, value interface{}) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Writes are immediately visible to all subsequent reads
	sc.data[key] = value
	return nil
}

func (sc *StrongConsistencyModel) GetConsistencyLevel() string {
	return "Strong Consistency (Linearizability)"
}

// EventualConsistencyModel - Convergence without ordering guarantees
// FAANG Interview Point: Social media likes, recommendation systems
type EventualConsistencyModel struct {
	replicas []map[string]interface{}
	mu       sync.RWMutex
}

func (ec *EventualConsistencyModel) Read(key string) (interface{}, error) {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	// May return stale data from any replica
	// Eventually all replicas will converge
	for _, replica := range ec.replicas {
		if value, exists := replica[key]; exists {
			return value, nil
		}
	}
	return nil, fmt.Errorf("key not found")
}

func (ec *EventualConsistencyModel) Write(key string, value interface{}) error {
	// Asynchronously propagate to all replicas
	go func() {
		ec.mu.Lock()
		defer ec.mu.Unlock()

		// Eventually propagate to all replicas (simulated delay)
		for i := range ec.replicas {
			time.Sleep(10 * time.Millisecond) // Simulate network delay
			ec.replicas[i][key] = value
		}
	}()
	return nil
}

func (ec *EventualConsistencyModel) GetConsistencyLevel() string {
	return "Eventual Consistency"
}

// Transaction Patterns for Distributed Systems
// Interview Focus: Handling transactions across multiple databases/services

// SagaPattern implements distributed transaction pattern
// FAANG Interview Point: E-commerce order processing, microservices transactions
type SagaPattern struct {
	steps         []SagaStep
	compensations []CompensationStep
}

type SagaStep interface {
	Execute() error
	GetStepName() string
}

type CompensationStep interface {
	Compensate() error
	GetStepName() string
}

// E-commerce order saga example
type ReserveInventoryStep struct {
	ProductID string
	Quantity  int
}

func (r *ReserveInventoryStep) Execute() error {
	// Reserve inventory in inventory service
	log.Printf("Reserving %d units of product %s", r.Quantity, r.ProductID)
	return nil // Simulate success
}

func (r *ReserveInventoryStep) GetStepName() string {
	return "ReserveInventory"
}

type ReserveInventoryCompensation struct {
	ProductID string
	Quantity  int
}

func (r *ReserveInventoryCompensation) Compensate() error {
	// Release reserved inventory
	log.Printf("Releasing %d units of product %s", r.Quantity, r.ProductID)
	return nil
}

func (r *ReserveInventoryCompensation) GetStepName() string {
	return "CompensateInventoryReservation"
}

// ExecuteSaga runs distributed transaction with compensation
func (s *SagaPattern) ExecuteSaga() error {
	completedSteps := 0

	// Execute all steps
	for i, step := range s.steps {
		if err := step.Execute(); err != nil {
			log.Printf("Step %s failed: %v", step.GetStepName(), err)

			// Compensate completed steps in reverse order
			for j := i - 1; j >= 0; j-- {
				if err := s.compensations[j].Compensate(); err != nil {
					log.Printf("Compensation failed for step %s: %v",
						s.compensations[j].GetStepName(), err)
				}
			}
			return fmt.Errorf("saga failed at step %s: %w", step.GetStepName(), err)
		}
		completedSteps++
	}

	log.Printf("Saga completed successfully with %d steps", completedSteps)
	return nil
}

// OutboxPattern ensures reliable event publishing
// FAANG Interview Point: Event-driven architectures, microservices communication
type OutboxPattern struct {
	db *sql.DB
}

type OutboxEvent struct {
	ID        string
	EventType string
	Payload   json.RawMessage
	CreatedAt time.Time
	Published bool
}

func (o *OutboxPattern) ProcessBusinessTransaction(userID, productID string, quantity int) error {
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Update business data
	_, err = tx.Exec("UPDATE inventory SET quantity = quantity - ? WHERE product_id = ?",
		quantity, productID)
	if err != nil {
		return err
	}

	// 2. Insert event into outbox table (same transaction)
	eventPayload, _ := json.Marshal(map[string]interface{}{
		"user_id":    userID,
		"product_id": productID,
		"quantity":   quantity,
		"timestamp":  time.Now(),
	})

	_, err = tx.Exec(`
		INSERT INTO outbox_events (id, event_type, payload, created_at, published)
		VALUES (?, ?, ?, ?, ?)`,
		fmt.Sprintf("order-%d", time.Now().UnixNano()),
		"order.created",
		eventPayload,
		time.Now(),
		false)
	if err != nil {
		return err
	}

	// 3. Commit both business data and event atomically
	return tx.Commit()
}

// PublishOutboxEvents processes unpublished events
func (o *OutboxPattern) PublishOutboxEvents() error {
	rows, err := o.db.Query(`
		SELECT id, event_type, payload FROM outbox_events 
		WHERE published = false 
		ORDER BY created_at LIMIT 100`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var event OutboxEvent
		if err := rows.Scan(&event.ID, &event.EventType, &event.Payload); err != nil {
			continue
		}

		// Publish to message queue/event bus
		if err := o.publishEvent(event); err == nil {
			// Mark as published
			o.db.Exec("UPDATE outbox_events SET published = true WHERE id = ?", event.ID)
		}
	}

	return nil
}

func (o *OutboxPattern) publishEvent(event OutboxEvent) error {
	// Simulate event publishing to Kafka/SQS/etc
	log.Printf("Publishing event %s: %s", event.EventType, string(event.Payload))
	return nil
}
