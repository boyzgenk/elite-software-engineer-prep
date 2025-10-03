package distributed_systems

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Two-Phase Commit (2PC) Protocol Implementation

// TransactionState represents the state of a transaction
type TransactionState int

const (
	TransactionActive TransactionState = iota
	TransactionPrepared
	TransactionCommitted
	TransactionAborted
)

// TwoPhaseCommitCoordinator manages 2PC protocol
type TwoPhaseCommitCoordinator struct {
	transactionID string
	participants  []Participant
	state         TransactionState
	timeout       time.Duration
	mu            sync.RWMutex
}

// Participant represents a node participating in 2PC
type Participant interface {
	Prepare(txID string) error
	Commit(txID string) error
	Abort(txID string) error
	GetID() string
}

// ParticipantNode implements a simple participant
type ParticipantNode struct {
	id           string
	transactions map[string]*TransactionLog
	mu           sync.RWMutex
}

// TransactionLog represents transaction state at a participant
type TransactionLog struct {
	ID         string
	State      TransactionState
	Operations []Operation
	Timestamp  time.Time
}

// Operation represents a transactional operation
type Operation struct {
	Type     string // "INSERT", "UPDATE", "DELETE"
	Table    string
	Key      string
	Value    interface{}
	OldValue interface{}
}

// NewTwoPhaseCommitCoordinator creates a new 2PC coordinator
func NewTwoPhaseCommitCoordinator(txID string, participants []Participant, timeout time.Duration) *TwoPhaseCommitCoordinator {
	return &TwoPhaseCommitCoordinator{
		transactionID: txID,
		participants:  participants,
		state:         TransactionActive,
		timeout:       timeout,
	}
}

// ExecuteTransaction executes a distributed transaction using 2PC
func (coord *TwoPhaseCommitCoordinator) ExecuteTransaction() error {
	// Phase 1: Prepare
	if err := coord.preparePhase(); err != nil {
		coord.abortTransaction()
		return fmt.Errorf("prepare phase failed: %w", err)
	}

	// Phase 2: Commit
	return coord.commitPhase()
}

// preparePhase executes the prepare phase of 2PC
func (coord *TwoPhaseCommitCoordinator) preparePhase() error {
	coord.mu.Lock()
	defer coord.mu.Unlock()

	log.Printf("Starting prepare phase for transaction %s", coord.transactionID)

	// Send prepare messages to all participants
	prepareChan := make(chan error, len(coord.participants))

	for _, participant := range coord.participants {
		go func(p Participant) {
			ctx, cancel := context.WithTimeout(context.Background(), coord.timeout)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				done <- p.Prepare(coord.transactionID)
			}()

			select {
			case err := <-done:
				prepareChan <- err
			case <-ctx.Done():
				prepareChan <- fmt.Errorf("participant %s prepare timeout", p.GetID())
			}
		}(participant)
	}

	// Wait for all prepare responses
	for i := 0; i < len(coord.participants); i++ {
		if err := <-prepareChan; err != nil {
			return err
		}
	}

	coord.state = TransactionPrepared
	log.Printf("All participants prepared for transaction %s", coord.transactionID)
	return nil
}

// commitPhase executes the commit phase of 2PC
func (coord *TwoPhaseCommitCoordinator) commitPhase() error {
	coord.mu.Lock()
	defer coord.mu.Unlock()

	log.Printf("Starting commit phase for transaction %s", coord.transactionID)

	// Send commit messages to all participants
	commitChan := make(chan error, len(coord.participants))

	for _, participant := range coord.participants {
		go func(p Participant) {
			ctx, cancel := context.WithTimeout(context.Background(), coord.timeout)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				done <- p.Commit(coord.transactionID)
			}()

			select {
			case err := <-done:
				commitChan <- err
			case <-ctx.Done():
				// If commit times out, we still consider it successful
				// as the participant should eventually commit
				commitChan <- nil
			}
		}(participant)
	}

	// Wait for all commit responses
	commitErrors := 0
	for i := 0; i < len(coord.participants); i++ {
		if err := <-commitChan; err != nil {
			log.Printf("Commit error from participant: %v", err)
			commitErrors++
		}
	}

	coord.state = TransactionCommitted
	log.Printf("Transaction %s committed (%d errors)", coord.transactionID, commitErrors)
	return nil
}

// abortTransaction aborts the transaction
func (coord *TwoPhaseCommitCoordinator) abortTransaction() {
	coord.mu.Lock()
	defer coord.mu.Unlock()

	log.Printf("Aborting transaction %s", coord.transactionID)

	// Send abort messages to all participants
	for _, participant := range coord.participants {
		go func(p Participant) {
			if err := p.Abort(coord.transactionID); err != nil {
				log.Printf("Abort error from participant %s: %v", p.GetID(), err)
			}
		}(participant)
	}

	coord.state = TransactionAborted
}

// NewParticipantNode creates a new participant node
func NewParticipantNode(id string) *ParticipantNode {
	return &ParticipantNode{
		id:           id,
		transactions: make(map[string]*TransactionLog),
	}
}

// Prepare implements the prepare phase for a participant
func (p *ParticipantNode) Prepare(txID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf("Participant %s preparing transaction %s", p.id, txID)

	// Check if transaction exists
	txLog, exists := p.transactions[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found", txID)
	}

	// Validate operations (simulate business logic validation)
	if err := p.validateOperations(txLog.Operations); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Mark as prepared
	txLog.State = TransactionPrepared
	txLog.Timestamp = time.Now()

	log.Printf("Participant %s prepared transaction %s", p.id, txID)
	return nil
}

// Commit implements the commit phase for a participant
func (p *ParticipantNode) Commit(txID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf("Participant %s committing transaction %s", p.id, txID)

	txLog, exists := p.transactions[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found", txID)
	}

	if txLog.State != TransactionPrepared {
		return fmt.Errorf("transaction %s not in prepared state", txID)
	}

	// Apply operations (simulate actual data changes)
	if err := p.applyOperations(txLog.Operations); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	txLog.State = TransactionCommitted
	txLog.Timestamp = time.Now()

	log.Printf("Participant %s committed transaction %s", p.id, txID)
	return nil
}

// Abort implements the abort phase for a participant
func (p *ParticipantNode) Abort(txID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf("Participant %s aborting transaction %s", p.id, txID)

	txLog, exists := p.transactions[txID]
	if !exists {
		// Already aborted or never existed
		return nil
	}

	txLog.State = TransactionAborted
	txLog.Timestamp = time.Now()

	// Clean up any prepared resources
	p.cleanupTransaction(txID)

	log.Printf("Participant %s aborted transaction %s", p.id, txID)
	return nil
}

// GetID returns the participant ID
func (p *ParticipantNode) GetID() string {
	return p.id
}

// AddOperation adds an operation to a transaction
func (p *ParticipantNode) AddOperation(txID string, op Operation) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	txLog, exists := p.transactions[txID]
	if !exists {
		txLog = &TransactionLog{
			ID:         txID,
			State:      TransactionActive,
			Operations: make([]Operation, 0),
			Timestamp:  time.Now(),
		}
		p.transactions[txID] = txLog
	}

	txLog.Operations = append(txLog.Operations, op)
	return nil
}

// Helper methods for ParticipantNode
func (p *ParticipantNode) validateOperations(operations []Operation) error {
	// Simulate validation logic
	for _, op := range operations {
		if op.Key == "" {
			return fmt.Errorf("empty key in operation")
		}
	}
	return nil
}

func (p *ParticipantNode) applyOperations(operations []Operation) error {
	// Simulate applying operations to actual data store
	for _, op := range operations {
		log.Printf("Applying operation: %s on %s.%s", op.Type, op.Table, op.Key)
	}
	return nil
}

func (p *ParticipantNode) cleanupTransaction(txID string) {
	delete(p.transactions, txID)
}

// Saga Pattern Implementation for Long-Running Transactions

// SagaStep represents a step in a saga
type SagaStep struct {
	ID                string
	Execute           func() error
	Compensate        func() error
	ExecuteTimeout    time.Duration
	CompensateTimeout time.Duration
}

// SagaExecution represents the execution state of a saga
type SagaExecution struct {
	ID               string
	Steps            []SagaStep
	CompletedSteps   []int
	CompensatedSteps []int
	CurrentStep      int
	State            SagaState
	mu               sync.RWMutex
}

// SagaState represents the state of saga execution
type SagaState int

const (
	SagaExecuting SagaState = iota
	SagaCompleted
	SagaCompensating
	SagaCompensated
	SagaFailed
)

// SagaOrchestrator manages saga execution
type SagaOrchestrator struct {
	executions map[string]*SagaExecution
	mu         sync.RWMutex
}

// NewSagaOrchestrator creates a new saga orchestrator
func NewSagaOrchestrator() *SagaOrchestrator {
	return &SagaOrchestrator{
		executions: make(map[string]*SagaExecution),
	}
}

// ExecuteSaga executes a saga with automatic compensation on failure
func (so *SagaOrchestrator) ExecuteSaga(sagaID string, steps []SagaStep) error {
	execution := &SagaExecution{
		ID:               sagaID,
		Steps:            steps,
		CompletedSteps:   make([]int, 0),
		CompensatedSteps: make([]int, 0),
		CurrentStep:      0,
		State:            SagaExecuting,
	}

	so.mu.Lock()
	so.executions[sagaID] = execution
	so.mu.Unlock()

	return so.executeSteps(execution)
}

// executeSteps executes saga steps sequentially
func (so *SagaOrchestrator) executeSteps(execution *SagaExecution) error {
	execution.mu.Lock()
	defer execution.mu.Unlock()

	log.Printf("Starting saga execution: %s", execution.ID)

	// Execute steps sequentially
	for i, step := range execution.Steps {
		execution.CurrentStep = i

		log.Printf("Executing step %d (%s) in saga %s", i, step.ID, execution.ID)

		if err := so.executeStepWithTimeout(step); err != nil {
			log.Printf("Step %d failed in saga %s: %v", i, execution.ID, err)
			execution.State = SagaCompensating
			return so.compensateSteps(execution)
		}

		execution.CompletedSteps = append(execution.CompletedSteps, i)
		log.Printf("Step %d completed in saga %s", i, execution.ID)
	}

	execution.State = SagaCompleted
	log.Printf("Saga %s completed successfully", execution.ID)
	return nil
}

// executeStepWithTimeout executes a step with timeout
func (so *SagaOrchestrator) executeStepWithTimeout(step SagaStep) error {
	done := make(chan error, 1)

	go func() {
		done <- step.Execute()
	}()

	timeout := step.ExecuteTimeout
	if timeout == 0 {
		timeout = 30 * time.Second // Default timeout
	}

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("step %s execution timeout", step.ID)
	}
}

// compensateSteps compensates completed steps in reverse order
func (so *SagaOrchestrator) compensateSteps(execution *SagaExecution) error {
	log.Printf("Starting compensation for saga %s", execution.ID)

	// Compensate in reverse order
	for i := len(execution.CompletedSteps) - 1; i >= 0; i-- {
		stepIndex := execution.CompletedSteps[i]
		step := execution.Steps[stepIndex]

		log.Printf("Compensating step %d (%s) in saga %s", stepIndex, step.ID, execution.ID)

		if err := so.compensateStepWithTimeout(step); err != nil {
			log.Printf("Compensation failed for step %d in saga %s: %v", stepIndex, execution.ID, err)
			execution.State = SagaFailed
			return err
		}

		execution.CompensatedSteps = append(execution.CompensatedSteps, stepIndex)
		log.Printf("Step %d compensated in saga %s", stepIndex, execution.ID)
	}

	execution.State = SagaCompensated
	log.Printf("Saga %s compensated successfully", execution.ID)
	return nil
}

// compensateStepWithTimeout compensates a step with timeout
func (so *SagaOrchestrator) compensateStepWithTimeout(step SagaStep) error {
	if step.Compensate == nil {
		return fmt.Errorf("no compensation function for step %s", step.ID)
	}

	done := make(chan error, 1)

	go func() {
		done <- step.Compensate()
	}()

	timeout := step.CompensateTimeout
	if timeout == 0 {
		timeout = 30 * time.Second // Default timeout
	}

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("step %s compensation timeout", step.ID)
	}
}

// GetSagaStatus returns the status of a saga execution
func (so *SagaOrchestrator) GetSagaStatus(sagaID string) (*SagaExecution, error) {
	so.mu.RLock()
	defer so.mu.RUnlock()

	execution, exists := so.executions[sagaID]
	if !exists {
		return nil, fmt.Errorf("saga %s not found", sagaID)
	}

	// Return a copy to avoid race conditions
	return &SagaExecution{
		ID:               execution.ID,
		Steps:            execution.Steps,
		CompletedSteps:   execution.CompletedSteps,
		CompensatedSteps: execution.CompensatedSteps,
		CurrentStep:      execution.CurrentStep,
		State:            execution.State,
	}, nil
}

// Outbox Pattern Implementation

// OutboxEvent represents an event in the outbox
type OutboxEvent struct {
	ID          string
	AggregateID string
	EventType   string
	Payload     interface{}
	CreatedAt   time.Time
	PublishedAt *time.Time
	Retries     int
	MaxRetries  int
}

// OutboxPublisher manages reliable event publishing
type OutboxPublisher struct {
	events     map[string]*OutboxEvent
	publishers map[string]EventPublisher
	maxRetries int
	retryDelay time.Duration
	mu         sync.RWMutex
	running    bool
	stopCh     chan struct{}
}

// EventPublisher interface for publishing events
type EventPublisher interface {
	Publish(event *OutboxEvent) error
	GetName() string
}

// SimpleEventPublisher implements a simple event publisher
type SimpleEventPublisher struct {
	name        string
	failureRate float64 // Simulate failures
}

// NewOutboxPublisher creates a new outbox publisher
func NewOutboxPublisher(maxRetries int, retryDelay time.Duration) *OutboxPublisher {
	return &OutboxPublisher{
		events:     make(map[string]*OutboxEvent),
		publishers: make(map[string]EventPublisher),
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		stopCh:     make(chan struct{}),
	}
}

// AddPublisher adds an event publisher
func (op *OutboxPublisher) AddPublisher(publisher EventPublisher) {
	op.mu.Lock()
	defer op.mu.Unlock()
	op.publishers[publisher.GetName()] = publisher
}

// AddEvent adds an event to the outbox
func (op *OutboxPublisher) AddEvent(event *OutboxEvent) {
	op.mu.Lock()
	defer op.mu.Unlock()

	if event.MaxRetries == 0 {
		event.MaxRetries = op.maxRetries
	}

	op.events[event.ID] = event
	log.Printf("Added event %s to outbox", event.ID)
}

// Start starts the outbox publisher
func (op *OutboxPublisher) Start() {
	op.mu.Lock()
	if op.running {
		op.mu.Unlock()
		return
	}
	op.running = true
	op.mu.Unlock()

	go op.publishLoop()
}

// Stop stops the outbox publisher
func (op *OutboxPublisher) Stop() {
	op.mu.Lock()
	if !op.running {
		op.mu.Unlock()
		return
	}
	op.running = false
	op.mu.Unlock()

	close(op.stopCh)
}

// publishLoop continuously publishes events from the outbox
func (op *OutboxPublisher) publishLoop() {
	ticker := time.NewTicker(op.retryDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			op.publishPendingEvents()
		case <-op.stopCh:
			return
		}
	}
}

// publishPendingEvents publishes all pending events
func (op *OutboxPublisher) publishPendingEvents() {
	op.mu.Lock()
	pendingEvents := make([]*OutboxEvent, 0)

	for _, event := range op.events {
		if event.PublishedAt == nil && event.Retries < event.MaxRetries {
			pendingEvents = append(pendingEvents, event)
		}
	}
	op.mu.Unlock()

	for _, event := range pendingEvents {
		op.publishEvent(event)
	}
}

// publishEvent publishes a single event to all publishers
func (op *OutboxPublisher) publishEvent(event *OutboxEvent) {
	op.mu.Lock()
	publishers := make([]EventPublisher, 0, len(op.publishers))
	for _, publisher := range op.publishers {
		publishers = append(publishers, publisher)
	}
	op.mu.Unlock()

	allSuccess := true

	for _, publisher := range publishers {
		if err := publisher.Publish(event); err != nil {
			log.Printf("Failed to publish event %s to %s: %v", event.ID, publisher.GetName(), err)
			allSuccess = false
		}
	}

	op.mu.Lock()
	if allSuccess {
		now := time.Now()
		event.PublishedAt = &now
		log.Printf("Successfully published event %s", event.ID)
	} else {
		event.Retries++
		log.Printf("Event %s publication failed, retries: %d/%d", event.ID, event.Retries, event.MaxRetries)
	}
	op.mu.Unlock()
}

// GetUnpublishedEvents returns events that haven't been published
func (op *OutboxPublisher) GetUnpublishedEvents() []*OutboxEvent {
	op.mu.RLock()
	defer op.mu.RUnlock()

	unpublished := make([]*OutboxEvent, 0)
	for _, event := range op.events {
		if event.PublishedAt == nil {
			unpublished = append(unpublished, event)
		}
	}

	return unpublished
}

// NewSimpleEventPublisher creates a simple event publisher
func NewSimpleEventPublisher(name string, failureRate float64) *SimpleEventPublisher {
	return &SimpleEventPublisher{
		name:        name,
		failureRate: failureRate,
	}
}

// Publish publishes an event (with simulated failures)
func (sep *SimpleEventPublisher) Publish(event *OutboxEvent) error {
	// Simulate random failures
	if rand.Float64() < sep.failureRate {
		return fmt.Errorf("simulated publish failure")
	}

	log.Printf("Publisher %s published event %s", sep.name, event.ID)
	return nil
}

// GetName returns the publisher name
func (sep *SimpleEventPublisher) GetName() string {
	return sep.name
}

// Event Sourcing with Snapshots Implementation

// EventStore manages events and snapshots
type EventStore struct {
	events    map[string][]DomainEvent
	snapshots map[string]*Snapshot
	mu        sync.RWMutex
}

// DomainEvent represents a domain event in event sourcing
type DomainEvent struct {
	ID          string
	AggregateID string
	Type        string
	Data        interface{}
	Version     int64
	Timestamp   time.Time
}

// Snapshot represents an aggregate snapshot
type Snapshot struct {
	AggregateID string
	Version     int64
	Data        interface{}
	Timestamp   time.Time
}

// EventSourcedAggregate interface for event-sourced aggregates
type EventSourcedAggregate interface {
	GetID() string
	GetVersion() int64
	ApplyEvent(event DomainEvent)
	GetUncommittedEvents() []DomainEvent
	ClearUncommittedEvents()
	CreateSnapshot() interface{}
	LoadFromSnapshot(data interface{})
}

// NewEventStore creates a new event store
func NewEventStore() *EventStore {
	return &EventStore{
		events:    make(map[string][]DomainEvent),
		snapshots: make(map[string]*Snapshot),
	}
}

// SaveEvents saves events for an aggregate
func (es *EventStore) SaveEvents(aggregateID string, expectedVersion int64, events []DomainEvent) error {
	es.mu.Lock()
	defer es.mu.Unlock()

	existingEvents := es.events[aggregateID]
	currentVersion := int64(0)
	if len(existingEvents) > 0 {
		currentVersion = existingEvents[len(existingEvents)-1].Version
	}

	if currentVersion != expectedVersion {
		return fmt.Errorf("concurrency conflict: expected version %d, actual version %d",
			expectedVersion, currentVersion)
	}

	// Append new events with incremented versions
	for i, event := range events {
		event.Version = expectedVersion + int64(i) + 1
		event.Timestamp = time.Now()
		existingEvents = append(existingEvents, event)
	}

	es.events[aggregateID] = existingEvents
	log.Printf("Saved %d events for aggregate %s", len(events), aggregateID)

	return nil
}

// LoadEvents loads events for an aggregate from a specific version
func (es *EventStore) LoadEvents(aggregateID string, fromVersion int64) ([]DomainEvent, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	allEvents := es.events[aggregateID]
	if len(allEvents) == 0 {
		return []DomainEvent{}, nil
	}

	// Find events after fromVersion
	var events []DomainEvent
	for _, event := range allEvents {
		if event.Version > fromVersion {
			events = append(events, event)
		}
	}

	return events, nil
}

// SaveSnapshot saves a snapshot for an aggregate
func (es *EventStore) SaveSnapshot(snapshot *Snapshot) error {
	es.mu.Lock()
	defer es.mu.Unlock()

	snapshot.Timestamp = time.Now()
	es.snapshots[snapshot.AggregateID] = snapshot

	log.Printf("Saved snapshot for aggregate %s at version %d",
		snapshot.AggregateID, snapshot.Version)

	return nil
}

// LoadSnapshot loads the latest snapshot for an aggregate
func (es *EventStore) LoadSnapshot(aggregateID string) (*Snapshot, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	snapshot, exists := es.snapshots[aggregateID]
	if !exists {
		return nil, fmt.Errorf("no snapshot found for aggregate %s", aggregateID)
	}

	return snapshot, nil
}

// LoadAggregate loads an aggregate from events and snapshots
func (es *EventStore) LoadAggregate(aggregateID string, aggregate EventSourcedAggregate) error {
	// Try to load from snapshot first
	snapshot, err := es.LoadSnapshot(aggregateID)
	fromVersion := int64(0)

	if err == nil {
		aggregate.LoadFromSnapshot(snapshot.Data)
		fromVersion = snapshot.Version
		log.Printf("Loaded aggregate %s from snapshot at version %d", aggregateID, fromVersion)
	}

	// Load events after snapshot
	events, err := es.LoadEvents(aggregateID, fromVersion)
	if err != nil {
		return err
	}

	// Apply events to rebuild current state
	for _, event := range events {
		aggregate.ApplyEvent(event)
	}

	log.Printf("Loaded aggregate %s with %d events applied", aggregateID, len(events))
	return nil
}

// SaveAggregate saves an aggregate's uncommitted events
func (es *EventStore) SaveAggregate(aggregate EventSourcedAggregate) error {
	uncommittedEvents := aggregate.GetUncommittedEvents()
	if len(uncommittedEvents) == 0 {
		return nil // Nothing to save
	}

	expectedVersion := aggregate.GetVersion() - int64(len(uncommittedEvents))

	if err := es.SaveEvents(aggregate.GetID(), expectedVersion, uncommittedEvents); err != nil {
		return err
	}

	aggregate.ClearUncommittedEvents()

	// Create snapshot if enough events have been accumulated
	if es.shouldCreateSnapshot(aggregate.GetID(), aggregate.GetVersion()) {
		snapshot := &Snapshot{
			AggregateID: aggregate.GetID(),
			Version:     aggregate.GetVersion(),
			Data:        aggregate.CreateSnapshot(),
		}
		es.SaveSnapshot(snapshot)
	}

	return nil
}

// shouldCreateSnapshot determines if a snapshot should be created
func (es *EventStore) shouldCreateSnapshot(aggregateID string, currentVersion int64) bool {
	// Create snapshot every 100 events
	return currentVersion%100 == 0
}

// GetEventStatistics returns statistics about stored events
func (es *EventStore) GetEventStatistics() map[string]interface{} {
	es.mu.RLock()
	defer es.mu.RUnlock()

	totalEvents := 0
	aggregateCount := len(es.events)
	snapshotCount := len(es.snapshots)

	for _, events := range es.events {
		totalEvents += len(events)
	}

	return map[string]interface{}{
		"total_events":             totalEvents,
		"aggregate_count":          aggregateCount,
		"snapshot_count":           snapshotCount,
		"avg_events_per_aggregate": float64(totalEvents) / float64(aggregateCount),
	}
}

// Example Event-Sourced Aggregate
type OrderAggregate struct {
	id                string
	version           int64
	status            string
	items             map[string]int
	totalAmount       float64
	uncommittedEvents []DomainEvent
	mu                sync.RWMutex
}

// NewOrderAggregate creates a new order aggregate
func NewOrderAggregate(id string) *OrderAggregate {
	return &OrderAggregate{
		id:                id,
		version:           0,
		status:            "created",
		items:             make(map[string]int),
		totalAmount:       0,
		uncommittedEvents: make([]DomainEvent, 0),
	}
}

// Implement EventSourcedAggregate interface
func (o *OrderAggregate) GetID() string {
	return o.id
}

func (o *OrderAggregate) GetVersion() int64 {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.version
}

func (o *OrderAggregate) ApplyEvent(event DomainEvent) {
	o.mu.Lock()
	defer o.mu.Unlock()

	switch event.Type {
	case "OrderCreated":
		o.status = "created"
	case "ItemAdded":
		data := event.Data.(map[string]interface{})
		item := data["item"].(string)
		quantity := data["quantity"].(int)
		o.items[item] = quantity
	case "OrderConfirmed":
		o.status = "confirmed"
	}

	o.version = event.Version
}

func (o *OrderAggregate) GetUncommittedEvents() []DomainEvent {
	o.mu.RLock()
	defer o.mu.RUnlock()

	// Return a copy to prevent race conditions
	events := make([]DomainEvent, len(o.uncommittedEvents))
	copy(events, o.uncommittedEvents)
	return events
}

func (o *OrderAggregate) ClearUncommittedEvents() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.uncommittedEvents = make([]DomainEvent, 0)
}

func (o *OrderAggregate) CreateSnapshot() interface{} {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return map[string]interface{}{
		"status":      o.status,
		"items":       o.items,
		"totalAmount": o.totalAmount,
	}
}

func (o *OrderAggregate) LoadFromSnapshot(data interface{}) {
	o.mu.Lock()
	defer o.mu.Unlock()

	snapshot := data.(map[string]interface{})
	o.status = snapshot["status"].(string)
	o.items = snapshot["items"].(map[string]int)
	o.totalAmount = snapshot["totalAmount"].(float64)
}

// Business methods that generate events
func (o *OrderAggregate) AddItem(item string, quantity int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	event := DomainEvent{
		AggregateID: o.id,
		Type:        "ItemAdded",
		Data: map[string]interface{}{
			"item":     item,
			"quantity": quantity,
		},
	}

	o.uncommittedEvents = append(o.uncommittedEvents, event)
	o.ApplyEvent(event)
}

func (o *OrderAggregate) ConfirmOrder() {
	o.mu.Lock()
	defer o.mu.Unlock()

	event := DomainEvent{
		AggregateID: o.id,
		Type:        "OrderConfirmed",
		Data:        map[string]interface{}{},
	}

	o.uncommittedEvents = append(o.uncommittedEvents, event)
	o.ApplyEvent(event)
}

// Benchmark and Performance Testing
type TransactionBenchmark struct {
	coordinators   []*TwoPhaseCommitCoordinator
	participants   []Participant
	completedTx    int64
	failedTx       int64
	avgLatency     int64
	maxConcurrency int
}

// NewTransactionBenchmark creates a benchmark for distributed transactions
func NewTransactionBenchmark(participantCount int, maxConcurrency int) *TransactionBenchmark {
	participants := make([]Participant, participantCount)
	for i := 0; i < participantCount; i++ {
		participants[i] = NewParticipantNode(fmt.Sprintf("participant-%d", i))
	}

	return &TransactionBenchmark{
		participants:   participants,
		maxConcurrency: maxConcurrency,
	}
}

// RunBenchmark runs a benchmark of distributed transactions
func (tb *TransactionBenchmark) RunBenchmark(transactionCount int, duration time.Duration) map[string]interface{} {
	start := time.Now()
	sem := make(chan struct{}, tb.maxConcurrency)
	var wg sync.WaitGroup

	for i := 0; i < transactionCount; i++ {
		wg.Add(1)
		go func(txID int) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			tb.runSingleTransaction(fmt.Sprintf("tx-%d", txID))
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	return map[string]interface{}{
		"total_transactions": transactionCount,
		"completed":          atomic.LoadInt64(&tb.completedTx),
		"failed":             atomic.LoadInt64(&tb.failedTx),
		"duration":           elapsed,
		"tps":                float64(transactionCount) / elapsed.Seconds(),
		"avg_latency":        time.Duration(atomic.LoadInt64(&tb.avgLatency)),
		"max_concurrency":    tb.maxConcurrency,
	}
}

// runSingleTransaction runs a single 2PC transaction
func (tb *TransactionBenchmark) runSingleTransaction(txID string) {
	start := time.Now()

	// Add some operations to participants
	for i, participant := range tb.participants {
		if pNode, ok := participant.(*ParticipantNode); ok {
			pNode.AddOperation(txID, Operation{
				Type:  "INSERT",
				Table: "orders",
				Key:   fmt.Sprintf("order-%s-%d", txID, i),
				Value: fmt.Sprintf("value-%d", i),
			})
		}
	}

	// Execute 2PC
	coordinator := NewTwoPhaseCommitCoordinator(txID, tb.participants, 5*time.Second)

	if err := coordinator.ExecuteTransaction(); err != nil {
		atomic.AddInt64(&tb.failedTx, 1)
	} else {
		atomic.AddInt64(&tb.completedTx, 1)
	}

	// Update average latency
	latency := time.Since(start).Nanoseconds()
	atomic.StoreInt64(&tb.avgLatency, latency)
}

// Interview Discussion Points:
//
// 1. Two-Phase Commit (2PC):
//    - Coordinator and participants pattern
//    - Blocking protocol - can hang if coordinator fails
//    - Used in XA transactions, distributed databases
//    - Problems: coordinator single point of failure, blocking
//
// 2. Three-Phase Commit (3PC):
//    - Adds "pre-commit" phase to reduce blocking
//    - More network overhead but better availability
//    - Rarely used in practice due to complexity
//
// 3. Saga Pattern:
//    - Choreography vs Orchestration approaches
//    - Forward recovery vs backward compensation
//    - Used by Netflix, Uber for microservices transactions
//    - Better availability than 2PC but eventual consistency
//
// 4. Outbox Pattern:
//    - Reliable event publishing using local transactions
//    - Ensures at-least-once delivery semantics
//    - Used with event sourcing and CQRS
//    - Prevents dual-write problem
//
// 5. Event Sourcing:
//    - Store events instead of current state
//    - Complete audit trail and time travel
//    - Snapshots for performance optimization
//    - Used by banking, trading systems
//
// 6. Performance Considerations:
//    - 2PC: High latency due to coordination overhead
//    - Saga: Better throughput but complex error handling
//    - Event Sourcing: Fast writes, potentially slow reads
//    - Network partitions affect all distributed transactions
//
// 7. Real-world Applications:
//    - Banking: 2PC for ACID transactions across accounts
//    - E-commerce: Saga for order processing workflows
//    - Microservices: Outbox pattern for service integration
//    - Audit systems: Event sourcing for compliance
//
// 8. Trade-offs:
//    - Consistency vs Availability (CAP theorem applies)
//    - Complexity vs Reliability
//    - Performance vs Correctness
//    - Operational overhead vs Business requirements
