// Context Patterns - Advanced Request Lifecycle Management
// Used by: Google (gRPC), Kubernetes, Docker, every major Go HTTP service
// Problem: Handle request timeouts, cancellation, and metadata propagation
// Solution: Context-driven patterns for distributed systems

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ==============================================================================
// 1. TIMEOUT PATTERNS - Time-based Cancellation
// ==============================================================================

// TimeoutService demonstrates various timeout strategies
// Used in: HTTP clients, database operations, external API calls
type TimeoutService struct {
	name string
}

// NewTimeoutService creates a new timeout service
func NewTimeoutService(name string) *TimeoutService {
	return &TimeoutService{name: name}
}

// SlowOperation simulates a potentially long-running operation
func (ts *TimeoutService) SlowOperation(ctx context.Context, duration time.Duration) (string, error) {
	// Create a channel to signal completion
	resultChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	// Start the operation in a goroutine
	go func() {
		// Simulate work
		select {
		case <-time.After(duration):
			resultChan <- fmt.Sprintf("%s: Operation completed after %v", ts.name, duration)
		case <-ctx.Done():
			errorChan <- fmt.Errorf("%s: Operation cancelled: %v", ts.name, ctx.Err())
			return
		}
	}()

	// Wait for either completion or context cancellation
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return "", err
	case <-ctx.Done():
		return "", fmt.Errorf("%s: Context timeout: %v", ts.name, ctx.Err())
	}
}

// OperationWithDeadline demonstrates deadline-based operations
func (ts *TimeoutService) OperationWithDeadline(ctx context.Context) (string, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return "", errors.New("no deadline set in context")
	}

	timeRemaining := time.Until(deadline)
	fmt.Printf("⏰ %s: %v remaining until deadline\n", ts.name, timeRemaining)

	if timeRemaining < 100*time.Millisecond {
		return "", errors.New("insufficient time remaining")
	}

	// Perform operation with awareness of remaining time
	operationDuration := timeRemaining / 2 // Use half the remaining time
	return ts.SlowOperation(ctx, operationDuration)
}

// ChainedOperations demonstrates context propagation through multiple operations
func (ts *TimeoutService) ChainedOperations(ctx context.Context) ([]string, error) {
	var results []string

	// Operation 1: Quick operation
	ctx1, cancel1 := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel1()

	result1, err := ts.SlowOperation(ctx1, 100*time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("operation 1 failed: %v", err)
	}
	results = append(results, result1)

	// Operation 2: Medium operation
	ctx2, cancel2 := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel2()

	result2, err := ts.SlowOperation(ctx2, 300*time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("operation 2 failed: %v", err)
	}
	results = append(results, result2)

	return results, nil
}

// ==============================================================================
// 2. CANCELLATION PATTERNS - Graceful Shutdown
// ==============================================================================

// CancellationService demonstrates cancellation propagation
// Used in: HTTP servers, background workers, cleanup operations
type CancellationService struct {
	name     string
	isActive int32 // atomic boolean
}

// NewCancellationService creates a new cancellation service
func NewCancellationService(name string) *CancellationService {
	return &CancellationService{
		name:     name,
		isActive: 1,
	}
}

// LongRunningTask demonstrates cancellation-aware long-running task
func (cs *CancellationService) LongRunningTask(ctx context.Context) error {
	fmt.Printf("🚀 %s: Starting long-running task\n", cs.name)
	atomic.StoreInt32(&cs.isActive, 1)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	stepCount := 0
	maxSteps := 50 // Total duration: 5 seconds

	for {
		select {
		case <-ctx.Done():
			atomic.StoreInt32(&cs.isActive, 0)
			fmt.Printf("🛑 %s: Task cancelled at step %d/%d: %v\n", cs.name, stepCount, maxSteps, ctx.Err())
			return ctx.Err()
		case <-ticker.C:
			stepCount++
			fmt.Printf("   %s: Progress %d/%d\n", cs.name, stepCount, maxSteps)

			if stepCount >= maxSteps {
				atomic.StoreInt32(&cs.isActive, 0)
				fmt.Printf("✅ %s: Task completed successfully\n", cs.name)
				return nil
			}
		}
	}
}

// BatchProcessor demonstrates cancellation in batch processing
func (cs *CancellationService) BatchProcessor(ctx context.Context, items []int) ([]int, error) {
	var processed []int
	var mu sync.Mutex

	// Process items in parallel with cancellation support
	semaphore := make(chan struct{}, 3) // Limit to 3 concurrent workers
	var wg sync.WaitGroup
	errChan := make(chan error, len(items))

	for i, item := range items {
		select {
		case <-ctx.Done():
			return processed, ctx.Err()
		case semaphore <- struct{}{}:
			wg.Add(1)
			go func(index, value int) {
				defer wg.Done()
				defer func() { <-semaphore }()

				// Simulate processing with cancellation check
				select {
				case <-ctx.Done():
					errChan <- fmt.Errorf("item %d cancelled", index)
					return
				case <-time.After(time.Duration(50+value%100) * time.Millisecond):
					// Processing completed
					result := value * 2

					mu.Lock()
					processed = append(processed, result)
					mu.Unlock()

					fmt.Printf("   Processed item %d: %d -> %d\n", index, value, result)
				}
			}(i, item)
		}
	}

	// Wait for completion or cancellation
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		close(errChan)
		// Check for any errors
		for err := range errChan {
			if err != nil {
				return processed, err
			}
		}
		return processed, nil
	case <-ctx.Done():
		return processed, ctx.Err()
	}
}

// ==============================================================================
// 3. CONTEXT VALUES - Request Metadata Propagation
// ==============================================================================

// ContextKey represents a key for context values
type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
	TraceIDKey   ContextKey = "trace_id"
	SessionIDKey ContextKey = "session_id"
)

// MetadataService demonstrates context value propagation
// Used in: HTTP middleware, logging, tracing, authentication
type MetadataService struct {
	name string
}

// NewMetadataService creates a new metadata service
func NewMetadataService(name string) *MetadataService {
	return &MetadataService{name: name}
}

// ProcessRequest demonstrates context value usage
func (ms *MetadataService) ProcessRequest(ctx context.Context, data string) (string, error) {
	// Extract metadata from context
	requestID := ms.getRequestID(ctx)
	userID := ms.getUserID(ctx)
	traceID := ms.getTraceID(ctx)

	fmt.Printf("🔍 %s: Processing request\n", ms.name)
	fmt.Printf("   Request ID: %s\n", requestID)
	fmt.Printf("   User ID: %s\n", userID)
	fmt.Printf("   Trace ID: %s\n", traceID)

	// Simulate processing
	time.Sleep(50 * time.Millisecond)

	result := fmt.Sprintf("%s processed '%s' for user %s (req: %s)", ms.name, data, userID, requestID)
	return result, nil
}

// CallDownstream demonstrates context propagation to other services
func (ms *MetadataService) CallDownstream(ctx context.Context, service *MetadataService, data string) (string, error) {
	// Add service-specific metadata
	enrichedCtx := context.WithValue(ctx, ContextKey("caller"), ms.name)

	fmt.Printf("📞 %s: Calling downstream service %s\n", ms.name, service.name)
	return service.ProcessRequest(enrichedCtx, data)
}

// Helper methods for context value extraction
func (ms *MetadataService) getRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return "unknown"
}

func (ms *MetadataService) getUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return "anonymous"
}

func (ms *MetadataService) getTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		return id
	}
	return "no-trace"
}

// ==============================================================================
// 4. DISTRIBUTED TRACING PATTERN - Request Tracking
// ==============================================================================

// TraceSpan represents a distributed tracing span
type TraceSpan struct {
	TraceID   string
	SpanID    string
	ParentID  string
	Service   string
	Operation string
	StartTime time.Time
	Duration  time.Duration
	Tags      map[string]string
}

// TracingService demonstrates distributed tracing with context
// Used in: Jaeger, Zipkin, OpenTelemetry implementations
type TracingService struct {
	serviceName string
	spans       []TraceSpan
	mu          sync.Mutex
}

// NewTracingService creates a new tracing service
func NewTracingService(serviceName string) *TracingService {
	return &TracingService{
		serviceName: serviceName,
		spans:       make([]TraceSpan, 0),
	}
}

// StartSpan creates a new tracing span
func (ts *TracingService) StartSpan(ctx context.Context, operation string) (context.Context, func()) {
	traceID := ts.getTraceID(ctx)
	if traceID == "" {
		traceID = fmt.Sprintf("trace-%d", time.Now().UnixNano())
	}

	parentSpanID := ts.getSpanID(ctx)
	spanID := fmt.Sprintf("span-%d", time.Now().UnixNano())

	span := TraceSpan{
		TraceID:   traceID,
		SpanID:    spanID,
		ParentID:  parentSpanID,
		Service:   ts.serviceName,
		Operation: operation,
		StartTime: time.Now(),
		Tags:      make(map[string]string),
	}

	// Add span to context
	spanCtx := context.WithValue(ctx, ContextKey("span_id"), spanID)
	spanCtx = context.WithValue(spanCtx, TraceIDKey, traceID)

	fmt.Printf("📊 Started span: %s.%s (trace: %s, parent: %s)\n",
		ts.serviceName, operation, traceID, parentSpanID)

	// Return context and finish function
	return spanCtx, func() {
		span.Duration = time.Since(span.StartTime)

		ts.mu.Lock()
		ts.spans = append(ts.spans, span)
		ts.mu.Unlock()

		fmt.Printf("📊 Finished span: %s.%s (duration: %v)\n",
			ts.serviceName, operation, span.Duration)
	}
}

// TracedOperation demonstrates operation with tracing
func (ts *TracingService) TracedOperation(ctx context.Context, operationName string, duration time.Duration) error {
	spanCtx, finish := ts.StartSpan(ctx, operationName)
	defer finish()

	// Simulate operation with context checking
	select {
	case <-time.After(duration):
		fmt.Printf("✅ %s.%s completed\n", ts.serviceName, operationName)
		return nil
	case <-spanCtx.Done():
		fmt.Printf("❌ %s.%s cancelled: %v\n", ts.serviceName, operationName, spanCtx.Err())
		return spanCtx.Err()
	}
}

// GetSpans returns collected spans for analysis
func (ts *TracingService) GetSpans() []TraceSpan {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	spans := make([]TraceSpan, len(ts.spans))
	copy(spans, ts.spans)
	return spans
}

// Helper methods for span extraction
func (ts *TracingService) getTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		return id
	}
	return ""
}

func (ts *TracingService) getSpanID(ctx context.Context) string {
	if id, ok := ctx.Value(ContextKey("span_id")).(string); ok {
		return id
	}
	return ""
}

// ==============================================================================
// 5. HTTP SERVER PATTERNS - Real-world Context Usage
// ==============================================================================

// HTTPHandler demonstrates context patterns in HTTP servers
// Used in: Gin, Echo, standard library HTTP servers
type HTTPHandler struct {
	services map[string]*MetadataService
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{
		services: map[string]*MetadataService{
			"auth":    NewMetadataService("auth-service"),
			"user":    NewMetadataService("user-service"),
			"payment": NewMetadataService("payment-service"),
		},
	}
}

// SimulateHTTPRequest demonstrates a complete HTTP request flow with context
func (h *HTTPHandler) SimulateHTTPRequest(requestTimeout time.Duration) error {
	// Create request context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Add request metadata
	requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
	userID := "user-12345"
	sessionID := fmt.Sprintf("session-%d", time.Now().UnixNano())

	ctx = context.WithValue(ctx, RequestIDKey, requestID)
	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctx = context.WithValue(ctx, SessionIDKey, sessionID)

	fmt.Printf("\n🌐 Simulating HTTP request (timeout: %v)\n", requestTimeout)
	fmt.Printf("   Request ID: %s\n", requestID)
	fmt.Printf("   User ID: %s\n", userID)

	// Step 1: Authentication
	authResult, err := h.services["auth"].ProcessRequest(ctx, "login-token")
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	fmt.Printf("   Auth: %s\n", authResult)

	// Step 2: User service call
	userResult, err := h.services["user"].ProcessRequest(ctx, "get-profile")
	if err != nil {
		return fmt.Errorf("user service failed: %v", err)
	}
	fmt.Printf("   User: %s\n", userResult)

	// Step 3: Payment service call (with potential timeout)
	paymentResult, err := h.services["payment"].ProcessRequest(ctx, "process-payment")
	if err != nil {
		return fmt.Errorf("payment service failed: %v", err)
	}
	fmt.Printf("   Payment: %s\n", paymentResult)

	fmt.Printf("✅ HTTP request completed successfully\n")
	return nil
}

func main() {
	fmt.Println("🔗 Context Patterns Demo")
	fmt.Println("========================")

	// Demo 1: Timeout Patterns
	fmt.Println("\n1️⃣ Timeout Patterns")
	timeoutService := NewTimeoutService("timeout-service")

	// Test successful operation
	ctx1, cancel1 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel1()

	result, err := timeoutService.SlowOperation(ctx1, 200*time.Millisecond)
	if err != nil {
		fmt.Printf("❌ Operation failed: %v\n", err)
	} else {
		fmt.Printf("✅ %s\n", result)
	}

	// Test timeout
	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()

	result, err = timeoutService.SlowOperation(ctx2, 300*time.Millisecond)
	if err != nil {
		fmt.Printf("⏰ Expected timeout: %v\n", err)
	}

	// Demo 2: Cancellation Patterns
	fmt.Println("\n2️⃣ Cancellation Patterns")
	cancelService := NewCancellationService("cancel-service")

	ctx3, cancel3 := context.WithCancel(context.Background())

	// Start long-running task
	go func() {
		err := cancelService.LongRunningTask(ctx3)
		if err != nil {
			fmt.Printf("Task ended with: %v\n", err)
		}
	}()

	// Cancel after 1 second
	time.Sleep(1 * time.Second)
	cancel3()
	time.Sleep(200 * time.Millisecond) // Wait for cancellation

	// Demo 3: Context Values and Metadata
	fmt.Println("\n3️⃣ Context Values and Metadata")

	// Create context with metadata
	baseCtx := context.Background()
	baseCtx = context.WithValue(baseCtx, RequestIDKey, "req-001")
	baseCtx = context.WithValue(baseCtx, UserIDKey, "user-alice")
	baseCtx = context.WithValue(baseCtx, TraceIDKey, "trace-xyz")

	service1 := NewMetadataService("api-gateway")
	service2 := NewMetadataService("business-logic")

	// Process request with metadata propagation
	result, err = service1.ProcessRequest(baseCtx, "user-data")
	if err != nil {
		fmt.Printf("❌ Service 1 failed: %v\n", err)
	} else {
		fmt.Printf("Service 1: %s\n", result)
	}

	// Call downstream service
	result, err = service1.CallDownstream(baseCtx, service2, "processed-data")
	if err != nil {
		fmt.Printf("❌ Downstream call failed: %v\n", err)
	} else {
		fmt.Printf("Downstream: %s\n", result)
	}

	// Demo 4: Distributed Tracing
	fmt.Println("\n4️⃣ Distributed Tracing")

	tracer1 := NewTracingService("web-service")
	tracer2 := NewTracingService("database-service")

	// Create trace context
	traceCtx := context.WithValue(context.Background(), TraceIDKey, "trace-distributed-001")

	// Start root span
	spanCtx1, finish1 := tracer1.StartSpan(traceCtx, "handle-request")

	// Perform operation
	err = tracer1.TracedOperation(spanCtx1, "validate-input", 100*time.Millisecond)
	if err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	}

	// Call another service
	spanCtx2, finish2 := tracer2.StartSpan(spanCtx1, "database-query")
	err = tracer2.TracedOperation(spanCtx2, "fetch-data", 150*time.Millisecond)
	if err != nil {
		fmt.Printf("❌ Database query failed: %v\n", err)
	}
	finish2()

	finish1()

	// Print collected spans
	fmt.Println("\n📊 Collected Spans:")
	for _, span := range tracer1.GetSpans() {
		fmt.Printf("   %s.%s: %v (trace: %s)\n", span.Service, span.Operation, span.Duration, span.TraceID)
	}
	for _, span := range tracer2.GetSpans() {
		fmt.Printf("   %s.%s: %v (trace: %s)\n", span.Service, span.Operation, span.Duration, span.TraceID)
	}

	// Demo 5: HTTP Server Simulation
	fmt.Println("\n5️⃣ HTTP Server Context Patterns")

	handler := NewHTTPHandler()

	// Successful request
	err = handler.SimulateHTTPRequest(2 * time.Second)
	if err != nil {
		fmt.Printf("❌ HTTP request failed: %v\n", err)
	}

	// Request with timeout
	err = handler.SimulateHTTPRequest(50 * time.Millisecond)
	if err != nil {
		fmt.Printf("⏰ HTTP request timed out: %v\n", err)
	}

	fmt.Println("\n🎯 Context Patterns Complete!")
	fmt.Println("Key Takeaways:")
	fmt.Println("• Use WithTimeout() for operation deadlines")
	fmt.Println("• Use WithCancel() for graceful shutdown")
	fmt.Println("• Propagate request metadata through context values")
	fmt.Println("• Implement distributed tracing for observability")
	fmt.Println("• Always check context.Done() in long-running operations")
	fmt.Println("• Context should be first parameter in function signatures")
}
