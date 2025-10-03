package microservices

import (
	"context"
	"crypto/rand"
	"crypto/tls"
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

// ServiceMesh represents a comprehensive service mesh implementation
// covering service discovery, load balancing, circuit breakers, and distributed tracing
// Key FAANG Interview Topics:
// - Service discovery patterns and health checking
// - Load balancing algorithms (round-robin, weighted, least-connections)
// - Circuit breaker implementation for fault tolerance
// - Distributed tracing and request correlation
// - mTLS simulation for secure service-to-service communication
// - Sidecar proxy pattern for transparent service mesh

// ServiceInstance represents a service instance in the mesh
type ServiceInstance struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	Metadata map[string]string `json:"metadata"`
	Health   HealthStatus      `json:"health"`
	Weight   int               `json:"weight"`
	LastSeen time.Time         `json:"last_seen"`
}

// HealthStatus represents service health states
type HealthStatus string

const (
	HealthHealthy   HealthStatus = "healthy"
	HealthUnhealthy HealthStatus = "unhealthy"
	HealthCritical  HealthStatus = "critical"
)

// ServiceRegistry handles service discovery and registration
type ServiceRegistry interface {
	Register(ctx context.Context, service *ServiceInstance) error
	Deregister(ctx context.Context, serviceID string) error
	Discover(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	HealthCheck(ctx context.Context, serviceID string) (HealthStatus, error)
	Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInstance, error)
}

// LoadBalancer defines load balancing strategies
type LoadBalancer interface {
	Select(instances []*ServiceInstance) (*ServiceInstance, error)
	UpdateWeights(weights map[string]int)
	GetAlgorithm() string
}

// CircuitBreaker provides fault tolerance for service calls
type CircuitBreaker interface {
	Call(ctx context.Context, fn func() (any, error)) (any, error)
	State() CircuitState
	Metrics() CircuitMetrics
}

// CircuitState represents circuit breaker states
type CircuitState string

const (
	CircuitClosed   CircuitState = "closed"
	CircuitOpen     CircuitState = "open"
	CircuitHalfOpen CircuitState = "half_open"
)

// CircuitMetrics provides circuit breaker observability
type CircuitMetrics struct {
	Requests         int64        `json:"requests"`
	Successes        int64        `json:"successes"`
	Failures         int64        `json:"failures"`
	ConsecutiveFails int64        `json:"consecutive_fails"`
	LastFailure      time.Time    `json:"last_failure"`
	State            CircuitState `json:"state"`
}

// TraceContext represents distributed tracing context
type TraceContext struct {
	TraceID  string            `json:"trace_id"`
	SpanID   string            `json:"span_id"`
	ParentID string            `json:"parent_id,omitempty"`
	Baggage  map[string]string `json:"baggage"`
}

// Tracer handles distributed tracing
type Tracer interface {
	StartSpan(ctx context.Context, operationName string) (context.Context, Span)
	InjectHeaders(ctx context.Context, headers map[string]string)
	ExtractHeaders(headers map[string]string) context.Context
}

// Span represents a tracing span
type Span interface {
	SetTag(key string, value any)
	SetBaggage(key, value string)
	LogEvent(event string)
	Finish()
	Context() TraceContext
}

// ServiceMeshProxy represents the sidecar proxy pattern
type ServiceMeshProxy interface {
	InterceptRequest(req *http.Request) (*http.Response, error)
	ApplyPolicy(serviceName string, policy *MeshPolicy) error
	GetMetrics() ProxyMetrics
}

// MeshPolicy defines service mesh policies
type MeshPolicy struct {
	ServiceName     string        `json:"service_name"`
	TimeoutDuration time.Duration `json:"timeout_duration"`
	RetryAttempts   int           `json:"retry_attempts"`
	CircuitBreaker  bool          `json:"circuit_breaker"`
	RateLimitRPS    int           `json:"rate_limit_rps"`
	RequireMTLS     bool          `json:"require_mtls"`
}

// ProxyMetrics provides proxy observability
type ProxyMetrics struct {
	RequestsTotal   int64   `json:"requests_total"`
	RequestsSuccess int64   `json:"requests_success"`
	RequestsFailure int64   `json:"requests_failure"`
	AvgLatencyMS    float64 `json:"avg_latency_ms"`
	P99LatencyMS    float64 `json:"p99_latency_ms"`
}

// DefaultServiceRegistry implements in-memory service discovery
type DefaultServiceRegistry struct {
	mu       sync.RWMutex
	services map[string]map[string]*ServiceInstance // serviceName -> serviceID -> instance
	watchers map[string][]chan []*ServiceInstance   // serviceName -> channels
	ttl      time.Duration
}

// NewServiceRegistry creates a new service registry with TTL-based cleanup
func NewServiceRegistry(ttl time.Duration) *DefaultServiceRegistry {
	registry := &DefaultServiceRegistry{
		services: make(map[string]map[string]*ServiceInstance),
		watchers: make(map[string][]chan []*ServiceInstance),
		ttl:      ttl,
	}

	// Start cleanup goroutine for expired services
	go registry.cleanupExpired()

	return registry
}

func (r *DefaultServiceRegistry) Register(ctx context.Context, service *ServiceInstance) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.services[service.Name] == nil {
		r.services[service.Name] = make(map[string]*ServiceInstance)
	}

	service.LastSeen = time.Now()
	r.services[service.Name][service.ID] = service

	// Notify watchers of service changes
	r.notifyWatchers(service.Name)

	return nil
}

func (r *DefaultServiceRegistry) Deregister(ctx context.Context, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for serviceName, instances := range r.services {
		if _, exists := instances[serviceID]; exists {
			delete(instances, serviceID)
			r.notifyWatchers(serviceName)
			return nil
		}
	}

	return errors.New("service not found")
}

func (r *DefaultServiceRegistry) Discover(ctx context.Context, serviceName string) ([]*ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	instances, exists := r.services[serviceName]
	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	var healthy []*ServiceInstance
	for _, instance := range instances {
		if instance.Health == HealthHealthy {
			healthy = append(healthy, instance)
		}
	}

	return healthy, nil
}

func (r *DefaultServiceRegistry) HealthCheck(ctx context.Context, serviceID string) (HealthStatus, error) {
	// Simulate health check via HTTP call
	// In production, this would make actual HTTP requests to service health endpoints
	return HealthHealthy, nil
}

func (r *DefaultServiceRegistry) Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInstance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ch := make(chan []*ServiceInstance, 1)
	r.watchers[serviceName] = append(r.watchers[serviceName], ch)

	// Send initial state
	if instances, exists := r.services[serviceName]; exists {
		var healthy []*ServiceInstance
		for _, instance := range instances {
			if instance.Health == HealthHealthy {
				healthy = append(healthy, instance)
			}
		}
		select {
		case ch <- healthy:
		default:
		}
	}

	return ch, nil
}

func (r *DefaultServiceRegistry) notifyWatchers(serviceName string) {
	if watchers, exists := r.watchers[serviceName]; exists {
		instances := r.services[serviceName]
		var healthy []*ServiceInstance
		for _, instance := range instances {
			if instance.Health == HealthHealthy {
				healthy = append(healthy, instance)
			}
		}

		for _, ch := range watchers {
			select {
			case ch <- healthy:
			default:
			}
		}
	}
}

func (r *DefaultServiceRegistry) cleanupExpired() {
	ticker := time.NewTicker(r.ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		now := time.Now()

		for serviceName, instances := range r.services {
			changed := false
			for serviceID, instance := range instances {
				if now.Sub(instance.LastSeen) > r.ttl {
					delete(instances, serviceID)
					changed = true
				}
			}
			if changed {
				r.notifyWatchers(serviceName)
			}
		}
		r.mu.Unlock()
	}
}

// RoundRobinLoadBalancer implements round-robin load balancing
type RoundRobinLoadBalancer struct {
	counter int64
}

func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{}
}

func (lb *RoundRobinLoadBalancer) Select(instances []*ServiceInstance) (*ServiceInstance, error) {
	if len(instances) == 0 {
		return nil, errors.New("no healthy instances available")
	}

	index := atomic.AddInt64(&lb.counter, 1) % int64(len(instances))
	return instances[index], nil
}

func (lb *RoundRobinLoadBalancer) UpdateWeights(weights map[string]int) {
	// Round-robin doesn't use weights
}

func (lb *RoundRobinLoadBalancer) GetAlgorithm() string {
	return "round_robin"
}

// WeightedLoadBalancer implements weighted round-robin load balancing
type WeightedLoadBalancer struct {
	mu      sync.RWMutex
	weights map[string]int
	current map[string]int
}

func NewWeightedLoadBalancer() *WeightedLoadBalancer {
	return &WeightedLoadBalancer{
		weights: make(map[string]int),
		current: make(map[string]int),
	}
}

func (lb *WeightedLoadBalancer) Select(instances []*ServiceInstance) (*ServiceInstance, error) {
	if len(instances) == 0 {
		return nil, errors.New("no healthy instances available")
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Use weighted round-robin algorithm
	var selected *ServiceInstance
	maxCurrentWeight := -1

	for _, instance := range instances {
		weight := lb.weights[instance.ID]
		if weight == 0 {
			weight = instance.Weight
			if weight == 0 {
				weight = 1 // Default weight
			}
		}

		lb.current[instance.ID] += weight

		if lb.current[instance.ID] > maxCurrentWeight {
			maxCurrentWeight = lb.current[instance.ID]
			selected = instance
		}
	}

	if selected != nil {
		lb.current[selected.ID] -= lb.getTotalWeight(instances)
	}

	return selected, nil
}

func (lb *WeightedLoadBalancer) UpdateWeights(weights map[string]int) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for id, weight := range weights {
		lb.weights[id] = weight
	}
}

func (lb *WeightedLoadBalancer) GetAlgorithm() string {
	return "weighted_round_robin"
}

func (lb *WeightedLoadBalancer) getTotalWeight(instances []*ServiceInstance) int {
	total := 0
	for _, instance := range instances {
		weight := lb.weights[instance.ID]
		if weight == 0 {
			weight = instance.Weight
			if weight == 0 {
				weight = 1
			}
		}
		total += weight
	}
	return total
}

// DefaultCircuitBreaker implements circuit breaker pattern
type DefaultCircuitBreaker struct {
	mu               sync.RWMutex
	maxFailures      int64
	resetTimeout     time.Duration
	state            CircuitState
	consecutiveFails int64
	lastFailure      time.Time
	requests         int64
	successes        int64
	failures         int64
}

func NewCircuitBreaker(maxFailures int64, resetTimeout time.Duration) *DefaultCircuitBreaker {
	return &DefaultCircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        CircuitClosed,
	}
}

func (cb *DefaultCircuitBreaker) Call(ctx context.Context, fn func() (any, error)) (any, error) {
	cb.mu.Lock()

	// Check if circuit should transition to half-open
	if cb.state == CircuitOpen && time.Since(cb.lastFailure) > cb.resetTimeout {
		cb.state = CircuitHalfOpen
		cb.consecutiveFails = 0
	}

	// Fail fast if circuit is open
	if cb.state == CircuitOpen {
		cb.mu.Unlock()
		return nil, errors.New("circuit breaker is open")
	}

	cb.requests++
	cb.mu.Unlock()

	// Execute the function
	result, err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.consecutiveFails++
		cb.lastFailure = time.Now()

		// Trip circuit if max failures reached
		if cb.consecutiveFails >= cb.maxFailures {
			cb.state = CircuitOpen
		}

		return nil, err
	}

	cb.successes++
	cb.consecutiveFails = 0

	// Close circuit if in half-open state and call succeeded
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
	}

	return result, nil
}

func (cb *DefaultCircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *DefaultCircuitBreaker) Metrics() CircuitMetrics {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return CircuitMetrics{
		Requests:         cb.requests,
		Successes:        cb.successes,
		Failures:         cb.failures,
		ConsecutiveFails: cb.consecutiveFails,
		LastFailure:      cb.lastFailure,
		State:            cb.state,
	}
}

// DefaultTracer implements distributed tracing
type DefaultTracer struct {
	mu    sync.RWMutex
	spans map[string]*DefaultSpan
}

func NewTracer() *DefaultTracer {
	return &DefaultTracer{
		spans: make(map[string]*DefaultSpan),
	}
}

type DefaultSpan struct {
	mu            sync.RWMutex
	traceContext  TraceContext
	operationName string
	tags          map[string]any
	logs          []SpanLog
	startTime     time.Time
	endTime       time.Time
	finished      bool
}

type SpanLog struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
}

func (t *DefaultTracer) StartSpan(ctx context.Context, operationName string) (context.Context, Span) {
	// Extract parent context if exists
	var parentTraceID, parentSpanID string
	if traceCtx := extractTraceContext(ctx); traceCtx != nil {
		parentTraceID = traceCtx.TraceID
		parentSpanID = traceCtx.SpanID
	}

	// Generate new trace ID if this is root span
	traceID := parentTraceID
	if traceID == "" {
		traceID = generateID()
	}

	spanID := generateID()

	span := &DefaultSpan{
		traceContext: TraceContext{
			TraceID:  traceID,
			SpanID:   spanID,
			ParentID: parentSpanID,
			Baggage:  make(map[string]string),
		},
		operationName: operationName,
		tags:          make(map[string]any),
		startTime:     time.Now(),
	}

	t.mu.Lock()
	t.spans[spanID] = span
	t.mu.Unlock()

	// Inject trace context into new context
	newCtx := injectTraceContext(ctx, &span.traceContext)

	return newCtx, span
}

func (t *DefaultTracer) InjectHeaders(ctx context.Context, headers map[string]string) {
	if traceCtx := extractTraceContext(ctx); traceCtx != nil {
		headers["X-Trace-ID"] = traceCtx.TraceID
		headers["X-Span-ID"] = traceCtx.SpanID
		if traceCtx.ParentID != "" {
			headers["X-Parent-ID"] = traceCtx.ParentID
		}

		// Inject baggage
		if baggage, err := json.Marshal(traceCtx.Baggage); err == nil {
			headers["X-Baggage"] = string(baggage)
		}
	}
}

func (t *DefaultTracer) ExtractHeaders(headers map[string]string) context.Context {
	traceID := headers["X-Trace-ID"]
	spanID := headers["X-Span-ID"]
	parentID := headers["X-Parent-ID"]

	if traceID == "" || spanID == "" {
		return context.Background()
	}

	baggage := make(map[string]string)
	if baggageStr := headers["X-Baggage"]; baggageStr != "" {
		json.Unmarshal([]byte(baggageStr), &baggage)
	}

	traceCtx := &TraceContext{
		TraceID:  traceID,
		SpanID:   spanID,
		ParentID: parentID,
		Baggage:  baggage,
	}

	return injectTraceContext(context.Background(), traceCtx)
}

func (s *DefaultSpan) SetTag(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tags[key] = value
}

func (s *DefaultSpan) SetBaggage(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.traceContext.Baggage[key] = value
}

func (s *DefaultSpan) LogEvent(event string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logs = append(s.logs, SpanLog{
		Timestamp: time.Now(),
		Event:     event,
	})
}

func (s *DefaultSpan) Finish() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.endTime = time.Now()
	s.finished = true
}

func (s *DefaultSpan) Context() TraceContext {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.traceContext
}

// DefaultServiceMeshProxy implements sidecar proxy pattern
type DefaultServiceMeshProxy struct {
	mu             sync.RWMutex
	policies       map[string]*MeshPolicy
	loadBalancer   LoadBalancer
	circuitBreaker CircuitBreaker
	tracer         Tracer
	registry       ServiceRegistry
	client         *http.Client
	metrics        ProxyMetrics
	tlsConfig      *tls.Config
}

func NewServiceMeshProxy(
	registry ServiceRegistry,
	loadBalancer LoadBalancer,
	tracer Tracer,
) *DefaultServiceMeshProxy {
	return &DefaultServiceMeshProxy{
		policies:       make(map[string]*MeshPolicy),
		loadBalancer:   loadBalancer,
		circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
		tracer:         tracer,
		registry:       registry,
		client:         &http.Client{Timeout: 30 * time.Second},
		tlsConfig:      &tls.Config{InsecureSkipVerify: false},
	}
}

func (p *DefaultServiceMeshProxy) InterceptRequest(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	// Start distributed tracing span
	ctx, span := p.tracer.StartSpan(ctx, fmt.Sprintf("%s %s", req.Method, req.URL.Path))
	defer span.Finish()

	span.SetTag("http.method", req.Method)
	span.SetTag("http.url", req.URL.String())

	// Extract service name from request (simplified)
	serviceName := extractServiceName(req.URL.Path)

	// Get policy for service
	policy := p.getPolicy(serviceName)
	if policy == nil {
		policy = &MeshPolicy{
			ServiceName:     serviceName,
			TimeoutDuration: 30 * time.Second,
			RetryAttempts:   3,
			CircuitBreaker:  true,
		}
	}

	// Apply mTLS if required
	if policy.RequireMTLS {
		req.Header.Set("X-Client-Cert", "simulated-mtls-cert")
	}

	// Inject tracing headers
	headers := make(map[string]string)
	p.tracer.InjectHeaders(ctx, headers)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Discover healthy service instances
	instances, err := p.registry.Discover(ctx, serviceName)
	if err != nil {
		atomic.AddInt64(&p.metrics.RequestsFailure, 1)
		return nil, fmt.Errorf("service discovery failed: %w", err)
	}

	// Select instance using load balancer
	instance, err := p.loadBalancer.Select(instances)
	if err != nil {
		atomic.AddInt64(&p.metrics.RequestsFailure, 1)
		return nil, fmt.Errorf("load balancing failed: %w", err)
	}

	// Update request URL with selected instance
	req.URL.Host = fmt.Sprintf("%s:%d", instance.Host, instance.Port)
	req.URL.Scheme = "http"
	if policy.RequireMTLS {
		req.URL.Scheme = "https"
	}

	// Execute request with circuit breaker and retries
	startTime := time.Now()

	result, err := p.circuitBreaker.Call(ctx, func() (any, error) {
		return p.executeWithRetries(req, policy.RetryAttempts)
	})

	latency := time.Since(startTime)

	// Update metrics
	atomic.AddInt64(&p.metrics.RequestsTotal, 1)
	if err != nil {
		atomic.AddInt64(&p.metrics.RequestsFailure, 1)
		span.SetTag("error", true)
		span.LogEvent(fmt.Sprintf("request failed: %v", err))
		return nil, err
	}

	atomic.AddInt64(&p.metrics.RequestsSuccess, 1)
	span.SetTag("http.status_code", result.(*http.Response).StatusCode)
	span.LogEvent("request completed successfully")

	// Update latency metrics (simplified)
	p.updateLatencyMetrics(latency)

	return result.(*http.Response), nil
}

func (p *DefaultServiceMeshProxy) ApplyPolicy(serviceName string, policy *MeshPolicy) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.policies[serviceName] = policy
	return nil
}

func (p *DefaultServiceMeshProxy) GetMetrics() ProxyMetrics {
	return ProxyMetrics{
		RequestsTotal:   atomic.LoadInt64(&p.metrics.RequestsTotal),
		RequestsSuccess: atomic.LoadInt64(&p.metrics.RequestsSuccess),
		RequestsFailure: atomic.LoadInt64(&p.metrics.RequestsFailure),
		AvgLatencyMS:    p.metrics.AvgLatencyMS,
		P99LatencyMS:    p.metrics.P99LatencyMS,
	}
}

func (p *DefaultServiceMeshProxy) getPolicy(serviceName string) *MeshPolicy {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.policies[serviceName]
}

func (p *DefaultServiceMeshProxy) executeWithRetries(req *http.Request, maxRetries int) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := p.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		lastErr = err
		if resp != nil {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
		}

		if attempt < maxRetries {
			// Exponential backoff
			backoff := time.Duration(1<<uint(attempt)) * 100 * time.Millisecond
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func (p *DefaultServiceMeshProxy) updateLatencyMetrics(latency time.Duration) {
	// Simplified latency tracking - in production, use proper percentile calculation
	latencyMS := float64(latency.Nanoseconds()) / 1e6

	// Update average (simplified)
	currentAvg := p.metrics.AvgLatencyMS
	if currentAvg == 0 {
		p.metrics.AvgLatencyMS = latencyMS
	} else {
		p.metrics.AvgLatencyMS = (currentAvg + latencyMS) / 2
	}

	// Update P99 (simplified)
	if latencyMS > p.metrics.P99LatencyMS {
		p.metrics.P99LatencyMS = latencyMS
	}
}

// Helper functions

func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func extractServiceName(path string) string {
	// Extract service name from path - simplified implementation
	// In production, this would use more sophisticated routing rules
	if len(path) > 1 {
		parts := []rune(path[1:])
		for i, r := range parts {
			if r == '/' {
				return string(parts[:i])
			}
		}
		return string(parts)
	}
	return "unknown"
}

// Context helpers for tracing
type traceContextKey struct{}

func injectTraceContext(ctx context.Context, traceCtx *TraceContext) context.Context {
	return context.WithValue(ctx, traceContextKey{}, traceCtx)
}

func extractTraceContext(ctx context.Context) *TraceContext {
	if traceCtx, ok := ctx.Value(traceContextKey{}).(*TraceContext); ok {
		return traceCtx
	}
	return nil
}

// ServiceMeshExample demonstrates comprehensive service mesh usage
type ServiceMeshExample struct {
	registry       ServiceRegistry
	loadBalancer   LoadBalancer
	circuitBreaker CircuitBreaker
	tracer         Tracer
	proxy          ServiceMeshProxy
}

func NewServiceMeshExample() *ServiceMeshExample {
	registry := NewServiceRegistry(30 * time.Second)
	loadBalancer := NewWeightedLoadBalancer()
	tracer := NewTracer()
	proxy := NewServiceMeshProxy(registry, loadBalancer, tracer)

	return &ServiceMeshExample{
		registry:       registry,
		loadBalancer:   loadBalancer,
		circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
		tracer:         tracer,
		proxy:          proxy,
	}
}

// DemonstrateServiceMesh shows comprehensive service mesh patterns
func (sm *ServiceMeshExample) DemonstrateServiceMesh() {
	ctx := context.Background()

	// Register service instances
	userService1 := &ServiceInstance{
		ID:       "user-service-1",
		Name:     "user-service",
		Host:     "10.0.1.10",
		Port:     8080,
		Health:   HealthHealthy,
		Weight:   100,
		Metadata: map[string]string{"version": "v1.2.0", "datacenter": "us-west-1"},
	}

	userService2 := &ServiceInstance{
		ID:       "user-service-2",
		Name:     "user-service",
		Host:     "10.0.1.11",
		Port:     8080,
		Health:   HealthHealthy,
		Weight:   150,
		Metadata: map[string]string{"version": "v1.2.0", "datacenter": "us-west-2"},
	}

	sm.registry.Register(ctx, userService1)
	sm.registry.Register(ctx, userService2)

	// Apply service mesh policy
	policy := &MeshPolicy{
		ServiceName:     "user-service",
		TimeoutDuration: 5 * time.Second,
		RetryAttempts:   3,
		CircuitBreaker:  true,
		RateLimitRPS:    1000,
		RequireMTLS:     true,
	}

	sm.proxy.ApplyPolicy("user-service", policy)

	// Demonstrate service discovery and load balancing
	instances, _ := sm.registry.Discover(ctx, "user-service")
	log.Printf("Discovered %d healthy instances", len(instances))

	// Demonstrate load balancing
	for i := 0; i < 5; i++ {
		selected, _ := sm.loadBalancer.Select(instances)
		log.Printf("Selected instance: %s (weight: %d)", selected.ID, selected.Weight)
	}

	// Demonstrate circuit breaker
	for i := 0; i < 10; i++ {
		_, err := sm.circuitBreaker.Call(ctx, func() (any, error) {
			// Simulate service call with 30% failure rate
			if i%3 == 0 {
				return nil, errors.New("simulated failure")
			}
			return "success", nil
		})

		log.Printf("Call %d: %v (circuit state: %s)", i+1, err, sm.circuitBreaker.State())
	}

	// Demonstrate distributed tracing
	ctx, span := sm.tracer.StartSpan(ctx, "user.GetProfile")
	span.SetTag("user.id", "12345")
	span.SetBaggage("request.id", "req-abc-123")

	// Simulate child span
	ctx, childSpan := sm.tracer.StartSpan(ctx, "database.query")
	childSpan.SetTag("db.table", "users")
	childSpan.LogEvent("query executed")
	childSpan.Finish()

	span.LogEvent("profile retrieved")
	span.Finish()

	log.Printf("Trace completed: %+v", span.Context())

	// Demonstrate proxy metrics
	metrics := sm.proxy.GetMetrics()
	log.Printf("Proxy Metrics: %+v", metrics)
}

// FAANG Interview Discussion Points:
//
// 1. Service Discovery Patterns:
//    - Health checking mechanisms (active vs passive)
//    - Consensus algorithms for distributed registries (Raft, Gossip)
//    - Service mesh vs client-side discovery trade-offs
//    - DNS-based vs API-based service discovery
//
// 2. Load Balancing Strategies:
//    - Consistent hashing for stateful services
//    - Locality-aware load balancing (zone/region preference)
//    - Adaptive load balancing based on response times
//    - Connection pooling and multiplexing considerations
//
// 3. Circuit Breaker Design:
//    - Fail-fast vs fail-silent strategies
//    - Bulkhead pattern for resource isolation
//    - Adaptive circuit breaking based on success rates
//    - Integration with monitoring and alerting systems
//
// 4. Distributed Tracing:
//    - Sampling strategies for high-traffic systems
//    - Trace correlation across async boundaries
//    - Performance impact of tracing instrumentation
//    - Privacy and security considerations for trace data
//
// 5. Service Mesh Architecture:
//    - Sidecar vs library-based implementations
//    - Control plane vs data plane separation
//    - mTLS certificate rotation and management
//    - Policy enforcement and service authorization
//
// 6. Scalability Considerations:
//    - Registry partitioning for large deployments
//    - Caching strategies for service metadata
//    - Graceful handling of registry unavailability
//    - Cross-region service discovery latency
