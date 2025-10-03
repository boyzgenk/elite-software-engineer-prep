package microservices

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// InterServiceCommunication represents comprehensive inter-service communication patterns
// covering synchronous/asynchronous communication, event-driven architecture, and distributed transactions
// Key FAANG Interview Topics:
// - Synchronous HTTP communication with retries and circuit breakers
// - Asynchronous messaging patterns (pub/sub, message queues)
// - Event-driven architecture and event sourcing
// - Saga pattern for distributed transactions and compensation
// - gRPC-style communication simulation
// - Message delivery guarantees and ordering

// Message represents a communication message between services
type Message struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Source      string            `json:"source"`
	Destination string            `json:"destination"`
	Payload     json.RawMessage   `json:"payload"`
	Headers     map[string]string `json:"headers"`
	Timestamp   time.Time         `json:"timestamp"`
	TTL         time.Duration     `json:"ttl"`
	Retry       int               `json:"retry"`
	MaxRetries  int               `json:"max_retries"`
}

// Event represents an event in event-driven architecture
type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Source    string          `json:"source"`
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
	Version   string          `json:"version"`
	Metadata  map[string]any  `json:"metadata"`
}

// SagaTransaction represents a distributed transaction using Saga pattern
type SagaTransaction struct {
	ID               string         `json:"id"`
	Steps            []*SagaStep    `json:"steps"`
	CurrentStep      int            `json:"current_step"`
	State            SagaState      `json:"state"`
	CompensationData map[string]any `json:"compensation_data"`
	StartTime        time.Time      `json:"start_time"`
	EndTime          time.Time      `json:"end_time"`
	Error            string         `json:"error,omitempty"`
}

type SagaState string

const (
	SagaPending      SagaState = "pending"
	SagaExecuting    SagaState = "executing"
	SagaCompleted    SagaState = "completed"
	SagaCompensating SagaState = "compensating"
	SagaFailed       SagaState = "failed"
)

// SagaStep represents a single step in a saga transaction
type SagaStep struct {
	ID           string          `json:"id"`
	ServiceName  string          `json:"service_name"`
	Action       string          `json:"action"`
	Payload      json.RawMessage `json:"payload"`
	Compensation string          `json:"compensation"`
	Completed    bool            `json:"completed"`
	Compensated  bool            `json:"compensated"`
	Error        string          `json:"error,omitempty"`
}

// SyncHTTPClient provides synchronous HTTP communication with retries
type SyncHTTPClient interface {
	CallService(ctx context.Context, request *ServiceRequest) (*ServiceResponse, error)
	CallWithRetry(ctx context.Context, request *ServiceRequest, retryPolicy *RetryPolicy) (*ServiceResponse, error)
	GetMetrics() HTTPClientMetrics
}

// ServiceRequest represents a service call request
type ServiceRequest struct {
	ServiceName string            `json:"service_name"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Headers     map[string]string `json:"headers"`
	Body        []byte            `json:"body"`
	Timeout     time.Duration     `json:"timeout"`
}

// ServiceResponse represents a service call response
type ServiceResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
	Latency    time.Duration     `json:"latency"`
}

// RetryPolicy defines retry behavior for service calls
type RetryPolicy struct {
	MaxAttempts       int           `json:"max_attempts"`
	InitialDelay      time.Duration `json:"initial_delay"`
	MaxDelay          time.Duration `json:"max_delay"`
	BackoffFactor     float64       `json:"backoff_factor"`
	RetryableStatuses []int         `json:"retryable_statuses"`
}

// HTTPClientMetrics provides HTTP client observability
type HTTPClientMetrics struct {
	RequestsTotal   int64   `json:"requests_total"`
	RequestsSuccess int64   `json:"requests_success"`
	RequestsFailure int64   `json:"requests_failure"`
	RequestsTimeout int64   `json:"requests_timeout"`
	AvgLatencyMS    float64 `json:"avg_latency_ms"`
	RetriesTotal    int64   `json:"retries_total"`
}

// MessageQueue provides asynchronous messaging capabilities
type MessageQueue interface {
	Publish(ctx context.Context, topic string, message *Message) error
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Unsubscribe(ctx context.Context, topic string, subscriberID string) error
	CreateTopic(ctx context.Context, topic string, config *TopicConfig) error
	GetMetrics() QueueMetrics
}

// MessageHandler processes incoming messages
type MessageHandler interface {
	HandleMessage(ctx context.Context, message *Message) error
	GetSubscriberID() string
}

// TopicConfig defines topic configuration
type TopicConfig struct {
	Partitions        int           `json:"partitions"`
	ReplicationFactor int           `json:"replication_factor"`
	RetentionTime     time.Duration `json:"retention_time"`
	MessageTTL        time.Duration `json:"message_ttl"`
	DeadLetterQueue   bool          `json:"dead_letter_queue"`
}

// QueueMetrics provides message queue observability
type QueueMetrics struct {
	MessagesPublished int64 `json:"messages_published"`
	MessagesConsumed  int64 `json:"messages_consumed"`
	MessagesFailed    int64 `json:"messages_failed"`
	TopicsTotal       int64 `json:"topics_total"`
	SubscribersTotal  int64 `json:"subscribers_total"`
}

// EventBus provides event-driven communication
type EventBus interface {
	PublishEvent(ctx context.Context, event *Event) error
	SubscribeToEvent(ctx context.Context, eventType string, handler EventHandler) error
	UnsubscribeFromEvent(ctx context.Context, eventType string, subscriberID string) error
	GetEventHistory(ctx context.Context, filters *EventFilters) ([]*Event, error)
}

// EventHandler processes events
type EventHandler interface {
	HandleEvent(ctx context.Context, event *Event) error
	GetSubscriberID() string
	GetEventTypes() []string
}

// EventFilters defines event filtering criteria
type EventFilters struct {
	EventTypes []string  `json:"event_types"`
	Sources    []string  `json:"sources"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Limit      int       `json:"limit"`
}

// SagaOrchestrator manages distributed transactions using Saga pattern
type SagaOrchestrator interface {
	StartSaga(ctx context.Context, sagaID string, steps []*SagaStep) error
	ExecuteStep(ctx context.Context, sagaID string, stepID string) error
	CompensateStep(ctx context.Context, sagaID string, stepID string) error
	GetSagaStatus(ctx context.Context, sagaID string) (*SagaTransaction, error)
	ListActiveSagas(ctx context.Context) ([]*SagaTransaction, error)
}

// gRPCStyleClient simulates gRPC-style communication
type gRPCStyleClient interface {
	Call(ctx context.Context, service, method string, request any) (any, error)
	CallStream(ctx context.Context, service, method string, requests <-chan any) (<-chan any, error)
	RegisterService(serviceName string, handler ServiceHandler) error
}

// ServiceHandler handles gRPC-style service calls
type ServiceHandler interface {
	HandleCall(ctx context.Context, method string, request any) (any, error)
	HandleStream(ctx context.Context, method string, requests <-chan any) (<-chan any, error)
}

// DefaultSyncHTTPClient implements synchronous HTTP communication
type DefaultSyncHTTPClient struct {
	client  *http.Client
	metrics HTTPClientMetrics
	mu      sync.RWMutex
}

func NewSyncHTTPClient(timeout time.Duration) *DefaultSyncHTTPClient {
	return &DefaultSyncHTTPClient{
		client: &http.Client{Timeout: timeout},
	}
}

func (c *DefaultSyncHTTPClient) CallService(ctx context.Context, request *ServiceRequest) (*ServiceResponse, error) {
	startTime := time.Now()
	atomic.AddInt64(&c.metrics.RequestsTotal, 1)

	// Build HTTP request
	url := fmt.Sprintf("http://%s%s", request.ServiceName, request.Path)

	httpReq, err := http.NewRequestWithContext(ctx, request.Method, url, nil)
	if err != nil {
		atomic.AddInt64(&c.metrics.RequestsFailure, 1)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for k, v := range request.Headers {
		httpReq.Header.Set(k, v)
	}

	// Apply timeout if specified
	if request.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, request.Timeout)
		defer cancel()
		httpReq = httpReq.WithContext(ctx)
	}

	// Execute request
	resp, err := c.client.Do(httpReq)
	latency := time.Since(startTime)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			atomic.AddInt64(&c.metrics.RequestsTimeout, 1)
		} else {
			atomic.AddInt64(&c.metrics.RequestsFailure, 1)
		}
		return nil, fmt.Errorf("service call failed: %w", err)
	}
	defer resp.Body.Close()

	// Build response
	response := &ServiceResponse{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
		Latency:    latency,
	}

	// Copy response headers
	for k, v := range resp.Header {
		if len(v) > 0 {
			response.Headers[k] = v[0]
		}
	}

	atomic.AddInt64(&c.metrics.RequestsSuccess, 1)
	c.updateLatencyMetrics(latency)

	return response, nil
}

func (c *DefaultSyncHTTPClient) CallWithRetry(ctx context.Context, request *ServiceRequest, retryPolicy *RetryPolicy) (*ServiceResponse, error) {
	var lastErr error
	delay := retryPolicy.InitialDelay

	for attempt := 0; attempt < retryPolicy.MaxAttempts; attempt++ {
		response, err := c.CallService(ctx, request)

		if err == nil && !c.isRetryableStatus(response.StatusCode, retryPolicy.RetryableStatuses) {
			return response, nil
		}

		lastErr = err
		if err == nil {
			lastErr = fmt.Errorf("retryable status code: %d", response.StatusCode)
		}

		// Don't retry on last attempt
		if attempt == retryPolicy.MaxAttempts-1 {
			break
		}

		atomic.AddInt64(&c.metrics.RetriesTotal, 1)

		// Wait with exponential backoff
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		// Calculate next delay with exponential backoff
		delay = time.Duration(float64(delay) * retryPolicy.BackoffFactor)
		if delay > retryPolicy.MaxDelay {
			delay = retryPolicy.MaxDelay
		}
	}

	return nil, fmt.Errorf("all retry attempts failed: %w", lastErr)
}

func (c *DefaultSyncHTTPClient) GetMetrics() HTTPClientMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return HTTPClientMetrics{
		RequestsTotal:   atomic.LoadInt64(&c.metrics.RequestsTotal),
		RequestsSuccess: atomic.LoadInt64(&c.metrics.RequestsSuccess),
		RequestsFailure: atomic.LoadInt64(&c.metrics.RequestsFailure),
		RequestsTimeout: atomic.LoadInt64(&c.metrics.RequestsTimeout),
		AvgLatencyMS:    c.metrics.AvgLatencyMS,
		RetriesTotal:    atomic.LoadInt64(&c.metrics.RetriesTotal),
	}
}

func (c *DefaultSyncHTTPClient) isRetryableStatus(statusCode int, retryableStatuses []int) bool {
	for _, code := range retryableStatuses {
		if statusCode == code {
			return true
		}
	}
	return false
}

func (c *DefaultSyncHTTPClient) updateLatencyMetrics(latency time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	latencyMS := float64(latency.Nanoseconds()) / 1e6
	if c.metrics.AvgLatencyMS == 0 {
		c.metrics.AvgLatencyMS = latencyMS
	} else {
		c.metrics.AvgLatencyMS = (c.metrics.AvgLatencyMS + latencyMS) / 2
	}
}

// InMemoryMessageQueue implements message queue using in-memory channels
type InMemoryMessageQueue struct {
	mu          sync.RWMutex
	topics      map[string]*Topic
	subscribers map[string]map[string]MessageHandler
	metrics     QueueMetrics
}

type Topic struct {
	Name    string
	Config  *TopicConfig
	Channel chan *Message
	DLQ     chan *Message // Dead Letter Queue
}

func NewInMemoryMessageQueue() *InMemoryMessageQueue {
	return &InMemoryMessageQueue{
		topics:      make(map[string]*Topic),
		subscribers: make(map[string]map[string]MessageHandler),
	}
}

func (mq *InMemoryMessageQueue) CreateTopic(ctx context.Context, topicName string, config *TopicConfig) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if _, exists := mq.topics[topicName]; exists {
		return fmt.Errorf("topic %s already exists", topicName)
	}

	topic := &Topic{
		Name:    topicName,
		Config:  config,
		Channel: make(chan *Message, 1000), // Buffered channel
	}

	if config.DeadLetterQueue {
		topic.DLQ = make(chan *Message, 100)
	}

	mq.topics[topicName] = topic
	mq.subscribers[topicName] = make(map[string]MessageHandler)

	// Start message processing goroutine
	go mq.processMessages(topicName)

	atomic.AddInt64(&mq.metrics.TopicsTotal, 1)

	return nil
}

func (mq *InMemoryMessageQueue) Publish(ctx context.Context, topicName string, message *Message) error {
	mq.mu.RLock()
	topic, exists := mq.topics[topicName]
	mq.mu.RUnlock()

	if !exists {
		return fmt.Errorf("topic %s does not exist", topicName)
	}

	// Set message metadata
	if message.ID == "" {
		message.ID = generateMessageID()
	}
	message.Timestamp = time.Now()

	// Check TTL
	if message.TTL > 0 && time.Since(message.Timestamp) > message.TTL {
		return fmt.Errorf("message TTL exceeded")
	}

	select {
	case topic.Channel <- message:
		atomic.AddInt64(&mq.metrics.MessagesPublished, 1)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("topic %s is full", topicName)
	}
}

func (mq *InMemoryMessageQueue) Subscribe(ctx context.Context, topicName string, handler MessageHandler) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if _, exists := mq.topics[topicName]; !exists {
		return fmt.Errorf("topic %s does not exist", topicName)
	}

	subscribers := mq.subscribers[topicName]
	subscribers[handler.GetSubscriberID()] = handler

	atomic.AddInt64(&mq.metrics.SubscribersTotal, 1)

	return nil
}

func (mq *InMemoryMessageQueue) Unsubscribe(ctx context.Context, topicName string, subscriberID string) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if subscribers, exists := mq.subscribers[topicName]; exists {
		if _, subscriberExists := subscribers[subscriberID]; subscriberExists {
			delete(subscribers, subscriberID)
			atomic.AddInt64(&mq.metrics.SubscribersTotal, -1)
			return nil
		}
	}

	return fmt.Errorf("subscriber %s not found for topic %s", subscriberID, topicName)
}

func (mq *InMemoryMessageQueue) GetMetrics() QueueMetrics {
	return QueueMetrics{
		MessagesPublished: atomic.LoadInt64(&mq.metrics.MessagesPublished),
		MessagesConsumed:  atomic.LoadInt64(&mq.metrics.MessagesConsumed),
		MessagesFailed:    atomic.LoadInt64(&mq.metrics.MessagesFailed),
		TopicsTotal:       atomic.LoadInt64(&mq.metrics.TopicsTotal),
		SubscribersTotal:  atomic.LoadInt64(&mq.metrics.SubscribersTotal),
	}
}

func (mq *InMemoryMessageQueue) processMessages(topicName string) {
	mq.mu.RLock()
	topic := mq.topics[topicName]
	mq.mu.RUnlock()

	for message := range topic.Channel {
		mq.mu.RLock()
		subscribers := mq.subscribers[topicName]
		mq.mu.RUnlock()

		// Deliver message to all subscribers
		for _, handler := range subscribers {
			go func(h MessageHandler, msg *Message) {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				err := h.HandleMessage(ctx, msg)
				if err != nil {
					atomic.AddInt64(&mq.metrics.MessagesFailed, 1)
					log.Printf("Message handling failed for subscriber %s: %v", h.GetSubscriberID(), err)

					// Send to DLQ if configured
					if topic.DLQ != nil {
						select {
						case topic.DLQ <- msg:
						default:
							log.Printf("DLQ is full for topic %s", topicName)
						}
					}
				} else {
					atomic.AddInt64(&mq.metrics.MessagesConsumed, 1)
				}
			}(handler, message)
		}
	}
}

// DefaultEventBus implements event-driven communication
type DefaultEventBus struct {
	mu         sync.RWMutex
	handlers   map[string]map[string]EventHandler // eventType -> subscriberID -> handler
	eventStore []*Event
	maxEvents  int
}

func NewEventBus(maxEvents int) *DefaultEventBus {
	return &DefaultEventBus{
		handlers:   make(map[string]map[string]EventHandler),
		eventStore: make([]*Event, 0),
		maxEvents:  maxEvents,
	}
}

func (eb *DefaultEventBus) PublishEvent(ctx context.Context, event *Event) error {
	// Set event metadata
	if event.ID == "" {
		event.ID = generateEventID()
	}
	event.Timestamp = time.Now()

	// Store event
	eb.mu.Lock()
	eb.eventStore = append(eb.eventStore, event)
	if len(eb.eventStore) > eb.maxEvents {
		eb.eventStore = eb.eventStore[1:]
	}

	// Get handlers for this event type
	handlers := eb.handlers[event.Type]
	eb.mu.Unlock()

	// Deliver event to handlers
	for _, handler := range handlers {
		go func(h EventHandler, e *Event) {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := h.HandleEvent(ctx, e); err != nil {
				log.Printf("Event handling failed for subscriber %s: %v", h.GetSubscriberID(), err)
			}
		}(handler, event)
	}

	return nil
}

func (eb *DefaultEventBus) SubscribeToEvent(ctx context.Context, eventType string, handler EventHandler) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.handlers[eventType] == nil {
		eb.handlers[eventType] = make(map[string]EventHandler)
	}

	eb.handlers[eventType][handler.GetSubscriberID()] = handler

	return nil
}

func (eb *DefaultEventBus) UnsubscribeFromEvent(ctx context.Context, eventType string, subscriberID string) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if handlers, exists := eb.handlers[eventType]; exists {
		delete(handlers, subscriberID)
		if len(handlers) == 0 {
			delete(eb.handlers, eventType)
		}
	}

	return nil
}

func (eb *DefaultEventBus) GetEventHistory(ctx context.Context, filters *EventFilters) ([]*Event, error) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	var result []*Event
	count := 0

	for i := len(eb.eventStore) - 1; i >= 0 && (filters.Limit == 0 || count < filters.Limit); i-- {
		event := eb.eventStore[i]

		// Apply filters
		if len(filters.EventTypes) > 0 {
			found := false
			for _, eventType := range filters.EventTypes {
				if event.Type == eventType {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if len(filters.Sources) > 0 {
			found := false
			for _, source := range filters.Sources {
				if event.Source == source {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if !filters.StartTime.IsZero() && event.Timestamp.Before(filters.StartTime) {
			continue
		}

		if !filters.EndTime.IsZero() && event.Timestamp.After(filters.EndTime) {
			continue
		}

		result = append(result, event)
		count++
	}

	return result, nil
}

// DefaultSagaOrchestrator implements the Saga pattern for distributed transactions
type DefaultSagaOrchestrator struct {
	mu         sync.RWMutex
	sagas      map[string]*SagaTransaction
	httpClient SyncHTTPClient
}

func NewSagaOrchestrator(httpClient SyncHTTPClient) *DefaultSagaOrchestrator {
	return &DefaultSagaOrchestrator{
		sagas:      make(map[string]*SagaTransaction),
		httpClient: httpClient,
	}
}

func (so *DefaultSagaOrchestrator) StartSaga(ctx context.Context, sagaID string, steps []*SagaStep) error {
	so.mu.Lock()
	defer so.mu.Unlock()

	if _, exists := so.sagas[sagaID]; exists {
		return fmt.Errorf("saga %s already exists", sagaID)
	}

	saga := &SagaTransaction{
		ID:               sagaID,
		Steps:            steps,
		CurrentStep:      0,
		State:            SagaPending,
		CompensationData: make(map[string]any),
		StartTime:        time.Now(),
	}

	so.sagas[sagaID] = saga

	// Start saga execution
	go so.executeSaga(ctx, sagaID)

	return nil
}

func (so *DefaultSagaOrchestrator) ExecuteStep(ctx context.Context, sagaID string, stepID string) error {
	so.mu.RLock()
	saga, exists := so.sagas[sagaID]
	so.mu.RUnlock()

	if !exists {
		return fmt.Errorf("saga %s not found", sagaID)
	}

	// Find step
	var step *SagaStep
	for _, s := range saga.Steps {
		if s.ID == stepID {
			step = s
			break
		}
	}

	if step == nil {
		return fmt.Errorf("step %s not found in saga %s", stepID, sagaID)
	}

	// Execute step action
	request := &ServiceRequest{
		ServiceName: step.ServiceName,
		Method:      "POST",
		Path:        "/" + step.Action,
		Body:        step.Payload,
		Timeout:     30 * time.Second,
	}

	response, err := so.httpClient.CallService(ctx, request)
	if err != nil {
		step.Error = err.Error()
		so.compensateSaga(ctx, sagaID)
		return fmt.Errorf("step execution failed: %w", err)
	}

	if response.StatusCode >= 400 {
		step.Error = fmt.Sprintf("HTTP %d", response.StatusCode)
		so.compensateSaga(ctx, sagaID)
		return fmt.Errorf("step failed with status %d", response.StatusCode)
	}

	step.Completed = true

	return nil
}

func (so *DefaultSagaOrchestrator) CompensateStep(ctx context.Context, sagaID string, stepID string) error {
	so.mu.RLock()
	saga, exists := so.sagas[sagaID]
	so.mu.RUnlock()

	if !exists {
		return fmt.Errorf("saga %s not found", sagaID)
	}

	// Find step
	var step *SagaStep
	for _, s := range saga.Steps {
		if s.ID == stepID {
			step = s
			break
		}
	}

	if step == nil {
		return fmt.Errorf("step %s not found in saga %s", stepID, sagaID)
	}

	if !step.Completed {
		return nil // Nothing to compensate
	}

	// Execute compensation action
	request := &ServiceRequest{
		ServiceName: step.ServiceName,
		Method:      "POST",
		Path:        "/" + step.Compensation,
		Body:        step.Payload,
		Timeout:     30 * time.Second,
	}

	_, err := so.httpClient.CallService(ctx, request)
	if err != nil {
		return fmt.Errorf("compensation failed: %w", err)
	}

	step.Compensated = true

	return nil
}

func (so *DefaultSagaOrchestrator) GetSagaStatus(ctx context.Context, sagaID string) (*SagaTransaction, error) {
	so.mu.RLock()
	defer so.mu.RUnlock()

	saga, exists := so.sagas[sagaID]
	if !exists {
		return nil, fmt.Errorf("saga %s not found", sagaID)
	}

	return saga, nil
}

func (so *DefaultSagaOrchestrator) ListActiveSagas(ctx context.Context) ([]*SagaTransaction, error) {
	so.mu.RLock()
	defer so.mu.RUnlock()

	var active []*SagaTransaction
	for _, saga := range so.sagas {
		if saga.State == SagaExecuting || saga.State == SagaCompensating {
			active = append(active, saga)
		}
	}

	return active, nil
}

func (so *DefaultSagaOrchestrator) executeSaga(ctx context.Context, sagaID string) {
	so.mu.Lock()
	saga := so.sagas[sagaID]
	saga.State = SagaExecuting
	so.mu.Unlock()

	// Execute each step sequentially
	for i, step := range saga.Steps {
		saga.CurrentStep = i

		err := so.ExecuteStep(ctx, sagaID, step.ID)
		if err != nil {
			log.Printf("Saga %s step %s failed: %v", sagaID, step.ID, err)
			return // compensateSaga is called from ExecuteStep
		}
	}

	// All steps completed successfully
	so.mu.Lock()
	saga.State = SagaCompleted
	saga.EndTime = time.Now()
	so.mu.Unlock()

	log.Printf("Saga %s completed successfully", sagaID)
}

func (so *DefaultSagaOrchestrator) compensateSaga(ctx context.Context, sagaID string) {
	so.mu.Lock()
	saga := so.sagas[sagaID]
	saga.State = SagaCompensating
	so.mu.Unlock()

	// Compensate completed steps in reverse order
	for i := saga.CurrentStep; i >= 0; i-- {
		step := saga.Steps[i]
		if step.Completed && !step.Compensated {
			err := so.CompensateStep(ctx, sagaID, step.ID)
			if err != nil {
				log.Printf("Compensation failed for saga %s step %s: %v", sagaID, step.ID, err)
			}
		}
	}

	so.mu.Lock()
	saga.State = SagaFailed
	saga.EndTime = time.Now()
	so.mu.Unlock()

	log.Printf("Saga %s compensation completed", sagaID)
}

// Helper functions

func generateMessageID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "msg-" + hex.EncodeToString(bytes)
}

func generateEventID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "evt-" + hex.EncodeToString(bytes)
}

// Sample message handler
type OrderMessageHandler struct {
	subscriberID string
}

func NewOrderMessageHandler(subscriberID string) *OrderMessageHandler {
	return &OrderMessageHandler{subscriberID: subscriberID}
}

func (h *OrderMessageHandler) HandleMessage(ctx context.Context, message *Message) error {
	log.Printf("OrderHandler [%s] processing message: %s", h.subscriberID, message.Type)

	// Simulate message processing
	time.Sleep(100 * time.Millisecond)

	// Simulate occasional failures
	if message.Type == "order.cancel" && message.Retry > 2 {
		return errors.New("order cancellation failed")
	}

	return nil
}

func (h *OrderMessageHandler) GetSubscriberID() string {
	return h.subscriberID
}

// Sample event handler
type PaymentEventHandler struct {
	subscriberID string
}

func NewPaymentEventHandler(subscriberID string) *PaymentEventHandler {
	return &PaymentEventHandler{subscriberID: subscriberID}
}

func (h *PaymentEventHandler) HandleEvent(ctx context.Context, event *Event) error {
	log.Printf("PaymentHandler [%s] processing event: %s", h.subscriberID, event.Type)

	// Simulate event processing based on type
	switch event.Type {
	case "order.created":
		log.Printf("Processing payment for order creation")
	case "order.cancelled":
		log.Printf("Refunding payment for order cancellation")
	default:
		log.Printf("Unknown event type: %s", event.Type)
	}

	return nil
}

func (h *PaymentEventHandler) GetSubscriberID() string {
	return h.subscriberID
}

func (h *PaymentEventHandler) GetEventTypes() []string {
	return []string{"order.created", "order.cancelled", "payment.processed"}
}

// InterServiceExample demonstrates comprehensive inter-service communication patterns
type InterServiceExample struct {
	httpClient       SyncHTTPClient
	messageQueue     MessageQueue
	eventBus         EventBus
	sagaOrchestrator SagaOrchestrator
}

func NewInterServiceExample() *InterServiceExample {
	httpClient := NewSyncHTTPClient(30 * time.Second)
	messageQueue := NewInMemoryMessageQueue()
	eventBus := NewEventBus(1000)
	sagaOrchestrator := NewSagaOrchestrator(httpClient)

	return &InterServiceExample{
		httpClient:       httpClient,
		messageQueue:     messageQueue,
		eventBus:         eventBus,
		sagaOrchestrator: sagaOrchestrator,
	}
}

func (example *InterServiceExample) DemonstrateInterServiceCommunication() {
	ctx := context.Background()

	// 1. Demonstrate synchronous HTTP communication with retries
	log.Println("=== Synchronous HTTP Communication ===")

	retryPolicy := &RetryPolicy{
		MaxAttempts:       3,
		InitialDelay:      100 * time.Millisecond,
		MaxDelay:          1 * time.Second,
		BackoffFactor:     2.0,
		RetryableStatuses: []int{500, 502, 503, 504},
	}

	request := &ServiceRequest{
		ServiceName: "user-service",
		Method:      "GET",
		Path:        "/users/123",
		Headers:     map[string]string{"Authorization": "Bearer token"},
		Timeout:     5 * time.Second,
	}

	// This would normally fail in a real environment, but demonstrates the pattern
	_, err := example.httpClient.CallWithRetry(ctx, request, retryPolicy)
	if err != nil {
		log.Printf("HTTP call failed (expected): %v", err)
	}

	metrics := example.httpClient.GetMetrics()
	log.Printf("HTTP Client Metrics: %+v", metrics)

	// 2. Demonstrate asynchronous messaging
	log.Println("\n=== Asynchronous Messaging ===")

	// Create topic
	topicConfig := &TopicConfig{
		Partitions:        3,
		ReplicationFactor: 2,
		RetentionTime:     24 * time.Hour,
		MessageTTL:        1 * time.Hour,
		DeadLetterQueue:   true,
	}

	example.messageQueue.CreateTopic(ctx, "orders", topicConfig)

	// Subscribe handlers
	orderHandler := NewOrderMessageHandler("order-service-1")
	example.messageQueue.Subscribe(ctx, "orders", orderHandler)

	// Publish messages
	messages := []*Message{
		{
			Type:        "order.created",
			Source:      "order-service",
			Destination: "inventory-service",
			Payload:     json.RawMessage(`{"orderId": "12345", "items": [{"sku": "ABC123", "qty": 2}]}`),
			MaxRetries:  3,
		},
		{
			Type:        "order.updated",
			Source:      "order-service",
			Destination: "notification-service",
			Payload:     json.RawMessage(`{"orderId": "12345", "status": "confirmed"}`),
			MaxRetries:  3,
		},
	}

	for _, msg := range messages {
		err := example.messageQueue.Publish(ctx, "orders", msg)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		} else {
			log.Printf("Published message: %s", msg.Type)
		}
	}

	time.Sleep(500 * time.Millisecond) // Allow message processing

	queueMetrics := example.messageQueue.GetMetrics()
	log.Printf("Queue Metrics: %+v", queueMetrics)

	// 3. Demonstrate event-driven architecture
	log.Println("\n=== Event-Driven Architecture ===")

	// Subscribe event handlers
	paymentHandler := NewPaymentEventHandler("payment-service")
	example.eventBus.SubscribeToEvent(ctx, "order.created", paymentHandler)
	example.eventBus.SubscribeToEvent(ctx, "order.cancelled", paymentHandler)

	// Publish events
	events := []*Event{
		{
			Type:    "order.created",
			Source:  "order-service",
			Data:    json.RawMessage(`{"orderId": "67890", "customerId": "cust-123", "amount": 99.99}`),
			Version: "1.0",
			Metadata: map[string]any{
				"correlation_id": "corr-abc-123",
				"source_system":  "web",
			},
		},
		{
			Type:    "order.cancelled",
			Source:  "order-service",
			Data:    json.RawMessage(`{"orderId": "67890", "reason": "customer_request"}`),
			Version: "1.0",
		},
	}

	for _, event := range events {
		err := example.eventBus.PublishEvent(ctx, event)
		if err != nil {
			log.Printf("Failed to publish event: %v", err)
		} else {
			log.Printf("Published event: %s", event.Type)
		}
	}

	time.Sleep(500 * time.Millisecond) // Allow event processing

	// Query event history
	filters := &EventFilters{
		EventTypes: []string{"order.created", "order.cancelled"},
		Limit:      10,
	}

	eventHistory, _ := example.eventBus.GetEventHistory(ctx, filters)
	log.Printf("Event History: %d events found", len(eventHistory))

	// 4. Demonstrate Saga pattern for distributed transactions
	log.Println("\n=== Saga Pattern for Distributed Transactions ===")

	sagaSteps := []*SagaStep{
		{
			ID:           "reserve-inventory",
			ServiceName:  "inventory-service",
			Action:       "reserve",
			Payload:      json.RawMessage(`{"orderId": "saga-123", "items": [{"sku": "ABC123", "qty": 1}]}`),
			Compensation: "release",
		},
		{
			ID:           "process-payment",
			ServiceName:  "payment-service",
			Action:       "charge",
			Payload:      json.RawMessage(`{"orderId": "saga-123", "amount": 49.99}`),
			Compensation: "refund",
		},
		{
			ID:           "create-shipment",
			ServiceName:  "shipping-service",
			Action:       "create_shipment",
			Payload:      json.RawMessage(`{"orderId": "saga-123", "address": {"street": "123 Main St"}}`),
			Compensation: "cancel_shipment",
		},
	}

	err = example.sagaOrchestrator.StartSaga(ctx, "order-saga-123", sagaSteps)
	if err != nil {
		log.Printf("Failed to start saga: %v", err)
	} else {
		log.Printf("Started saga: order-saga-123")
	}

	time.Sleep(2 * time.Second) // Allow saga processing

	// Check saga status
	sagaStatus, _ := example.sagaOrchestrator.GetSagaStatus(ctx, "order-saga-123")
	if sagaStatus != nil {
		log.Printf("Saga Status: %s (Step: %d/%d)", sagaStatus.State, sagaStatus.CurrentStep+1, len(sagaStatus.Steps))
	}

	activeSagas, _ := example.sagaOrchestrator.ListActiveSagas(ctx)
	log.Printf("Active Sagas: %d", len(activeSagas))
}

// FAANG Interview Discussion Points:
//
// 1. Synchronous Communication Patterns:
//    - Circuit breaker integration for fault tolerance
//    - Connection pooling and HTTP/2 multiplexing
//    - Request/response correlation and tracing
//    - Timeout strategies and cascading failure prevention
//
// 2. Asynchronous Messaging Architecture:
//    - Message ordering guarantees and partitioning strategies
//    - At-least-once vs exactly-once delivery semantics
//    - Dead letter queues and poison message handling
//    - Consumer group patterns and load balancing
//
// 3. Event-Driven Architecture Considerations:
//    - Event schema evolution and versioning
//    - Event sourcing vs traditional database patterns
//    - CQRS (Command Query Responsibility Segregation)
//    - Event replay and time-travel debugging
//
// 4. Distributed Transaction Patterns:
//    - Two-Phase Commit (2PC) vs Saga pattern trade-offs
//    - Choreography vs Orchestration in Saga patterns
//    - Compensation action design and idempotency
//    - Long-running transaction timeout and recovery
//
// 5. Message Delivery Guarantees:
//    - Network partition handling and split-brain scenarios
//    - Message deduplication strategies
//    - Consensus algorithms (Raft, PBFT) for distributed queues
//    - CAP theorem implications for messaging systems
//
// 6. Performance and Scalability:
//    - Message batching and compression strategies
//    - Horizontal scaling of message brokers
//    - Hot partition problems and load balancing
//    - Cross-region replication and latency optimization
//
// 7. Monitoring and Observability:
//    - End-to-end request tracing across async boundaries
//    - Message flow visualization and dependency mapping
//    - Saga execution monitoring and alerting
//    - Performance metrics and SLA monitoring
