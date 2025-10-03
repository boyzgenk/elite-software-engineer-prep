package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// FAANG Interview Focus: Stateless Service Design Patterns
// Key Topics: Horizontal scaling, session management, external state storage

// StatelessSessionData represents user session information
type StatelessSessionData struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ExternalSessionStore manages sessions in external storage (Redis simulation)
type ExternalSessionStore struct {
	sessions map[string]StatelessSessionData
	mu       sync.RWMutex
	prefix   string
}

// NewExternalSessionStore creates a session store (simulating Redis)
func NewExternalSessionStore() *ExternalSessionStore {
	return &ExternalSessionStore{
		sessions: make(map[string]StatelessSessionData),
		prefix:   "session:",
	}
}

// CreateSession stores session data externally
func (s *ExternalSessionStore) CreateSession(ctx context.Context, sessionID string, data StatelessSessionData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.prefix + sessionID
	s.sessions[key] = data

	// Simulate TTL cleanup in real Redis
	go s.cleanupExpiredSession(key, data.ExpiresAt)

	return nil
}

// GetSession retrieves session data from external storage
func (s *ExternalSessionStore) GetSession(ctx context.Context, sessionID string) (*StatelessSessionData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.prefix + sessionID

	data, exists := s.sessions[key]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	if time.Now().After(data.ExpiresAt) {
		go s.DeleteSession(ctx, sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return &data, nil
}

// DeleteSession removes session from external storage
func (s *ExternalSessionStore) DeleteSession(ctx context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.prefix + sessionID
	delete(s.sessions, key)

	return nil
}

// cleanupExpiredSession removes expired sessions
func (s *ExternalSessionStore) cleanupExpiredSession(key string, expiry time.Time) {
	time.Sleep(time.Until(expiry))
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, key)
}

// StatelessService demonstrates stateless microservice patterns
type StatelessService struct {
	config       *StatelessServiceConfig
	sessionStore *ExternalSessionStore
	healthInfo   StatelessHealthStatus
	mu           sync.RWMutex
}

// StatelessServiceConfig holds external configuration
type StatelessServiceConfig struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	Environment string
	ServiceName string
	Version     string
}

// LoadStatelessConfigFromEnv demonstrates 12-factor app configuration
func LoadStatelessConfigFromEnv() *StatelessServiceConfig {
	return &StatelessServiceConfig{
		Port:        getStatelessEnvOrDefault("PORT", "8080"),
		DatabaseURL: getStatelessEnvOrDefault("DATABASE_URL", "postgres://localhost/mydb"),
		RedisURL:    getStatelessEnvOrDefault("REDIS_URL", "localhost:6379"),
		JWTSecret:   getStatelessEnvOrDefault("JWT_SECRET", "default-secret"),
		Environment: getStatelessEnvOrDefault("ENVIRONMENT", "development"),
		ServiceName: getStatelessEnvOrDefault("SERVICE_NAME", "stateless-service"),
		Version:     getStatelessEnvOrDefault("VERSION", "1.0.0"),
	}
}

func getStatelessEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// StatelessHealthStatus represents service health
type StatelessHealthStatus struct {
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
	Checks      map[string]string `json:"checks"`
}

// NewStatelessService creates a new stateless service instance
func NewStatelessService(config *StatelessServiceConfig) *StatelessService {
	sessionStore := NewExternalSessionStore()

	return &StatelessService{
		config:       config,
		sessionStore: sessionStore,
		healthInfo: StatelessHealthStatus{
			Status:      "healthy",
			Timestamp:   time.Now(),
			Version:     config.Version,
			Environment: config.Environment,
			Checks:      make(map[string]string),
		},
	}
}

// UserProcessRequest represents incoming user data
type UserProcessRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Action   string `json:"action"`
}

// UserProcessResponse represents processed response
type UserProcessResponse struct {
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	ProcessedBy string    `json:"processed_by"`
}

// ProcessUserRequest handles requests without server state
func (s *StatelessService) ProcessUserRequest(ctx context.Context, req UserProcessRequest) (*UserProcessResponse, error) {
	// Stateless processing - no server state stored
	userID := generateStatelessUserID(req.Username, req.Email)

	// All state comes from external sources (DB, cache, etc.)
	response := &UserProcessResponse{
		ID:          userID,
		Message:     fmt.Sprintf("Processed %s for user %s", req.Action, req.Username),
		Timestamp:   time.Now(),
		ProcessedBy: s.config.ServiceName,
	}

	// Could store result in external database here
	log.Printf("Processed request for user %s, action: %s", req.Username, req.Action)

	return response, nil
}

func generateStatelessUserID(username, email string) string {
	return fmt.Sprintf("user-%s-%d", username, time.Now().Unix())
}

// HTTP Handlers for stateless operations

// StatelessHealthCheckHandler provides health endpoint for load balancers
func (s *StatelessService) StatelessHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check external dependencies
	s.healthInfo.Checks = s.checkStatelessDependencies(r.Context())
	s.healthInfo.Timestamp = time.Now()

	// Determine overall status
	healthy := true
	for _, status := range s.healthInfo.Checks {
		if status != "healthy" {
			healthy = false
			break
		}
	}

	if healthy {
		s.healthInfo.Status = "healthy"
		w.WriteHeader(http.StatusOK)
	} else {
		s.healthInfo.Status = "unhealthy"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.healthInfo)
}

// checkStatelessDependencies verifies external service health
func (s *StatelessService) checkStatelessDependencies(ctx context.Context) map[string]string {
	checks := make(map[string]string)

	// Check session store connection (simulated)
	if s.sessionStore != nil {
		checks["session_store"] = "healthy"
	} else {
		checks["session_store"] = "unhealthy"
	}

	// Could add database, external API checks here
	checks["service"] = "healthy"

	return checks
}

// StatelessReadinessHandler checks if service is ready to accept traffic
func (s *StatelessService) StatelessReadinessHandler(w http.ResponseWriter, r *http.Request) {
	// Check if service can handle requests
	if s.isStatelessReady(r.Context()) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
	}
}

func (s *StatelessService) isStatelessReady(ctx context.Context) bool {
	// Check critical dependencies
	return s.sessionStore != nil
}

// StatelessProcessRequestHandler handles stateless HTTP requests
func (s *StatelessService) StatelessProcessRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req UserProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response, err := s.ProcessUserRequest(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StatelessSessionHandler manages sessions without server state
func (s *StatelessService) StatelessSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		s.createStatelessSessionHandler(w, r, sessionID)
	case http.MethodGet:
		s.getStatelessSessionHandler(w, r, sessionID)
	case http.MethodDelete:
		s.deleteStatelessSessionHandler(w, r, sessionID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *StatelessService) createStatelessSessionHandler(w http.ResponseWriter, r *http.Request, sessionID string) {
	var sessionData StatelessSessionData
	if err := json.NewDecoder(r.Body).Decode(&sessionData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	sessionData.CreatedAt = time.Now()
	sessionData.ExpiresAt = time.Now().Add(24 * time.Hour)

	if err := s.sessionStore.CreateSession(r.Context(), sessionID, sessionData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sessionData)
}

func (s *StatelessService) getStatelessSessionHandler(w http.ResponseWriter, r *http.Request, sessionID string) {
	sessionData, err := s.sessionStore.GetSession(r.Context(), sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessionData)
}

func (s *StatelessService) deleteStatelessSessionHandler(w http.ResponseWriter, r *http.Request, sessionID string) {
	if err := s.sessionStore.DeleteSession(r.Context(), sessionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SetupStatelessRoutes configures HTTP routes for the stateless service
func (s *StatelessService) SetupStatelessRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health endpoints for load balancer integration
	mux.HandleFunc("/health", s.StatelessHealthCheckHandler)
	mux.HandleFunc("/ready", s.StatelessReadinessHandler)

	// API endpoints
	mux.HandleFunc("/api/v1/process", s.StatelessProcessRequestHandler)
	mux.HandleFunc("/api/v1/session", s.StatelessSessionHandler)

	return mux
}

// StatelessLoadBalancerConfig demonstrates load balancer integration
type StatelessLoadBalancerConfig struct {
	Instances   []StatelessServiceInstance `json:"instances"`
	Algorithm   string                     `json:"algorithm"` // round-robin, least-conn, ip-hash
	HealthCheck StatelessHealthCheckConfig `json:"health_check"`
}

type StatelessServiceInstance struct {
	Address string `json:"address"`
	Weight  int    `json:"weight"`
	Status  string `json:"status"` // active, inactive, draining
}

type StatelessHealthCheckConfig struct {
	Endpoint string        `json:"endpoint"`
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
	Retries  int           `json:"retries"`
}

// HorizontalStatelessScalingDemo demonstrates scaling patterns
func HorizontalStatelessScalingDemo() {
	// Example configuration for multiple instances
	scalingConfig := StatelessLoadBalancerConfig{
		Instances: []StatelessServiceInstance{
			{Address: "10.0.1.1:8080", Weight: 100, Status: "active"},
			{Address: "10.0.1.2:8080", Weight: 100, Status: "active"},
			{Address: "10.0.1.3:8080", Weight: 50, Status: "draining"},
		},
		Algorithm: "round-robin",
		HealthCheck: StatelessHealthCheckConfig{
			Endpoint: "/health",
			Interval: 30 * time.Second,
			Timeout:  5 * time.Second,
			Retries:  3,
		},
	}

	log.Printf("Stateless scaling configuration: %+v", scalingConfig)
}

// StatelessConfigurationManager demonstrates configuration patterns
type StatelessConfigurationManager struct {
	configs map[string]interface{}
	mu      sync.RWMutex
}

// NewStatelessConfigurationManager creates a configuration manager
func NewStatelessConfigurationManager() *StatelessConfigurationManager {
	return &StatelessConfigurationManager{
		configs: make(map[string]interface{}),
	}
}

// GetConfig retrieves configuration values
func (cm *StatelessConfigurationManager) GetConfig(key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	value, exists := cm.configs[key]
	return value, exists
}

// SetConfig sets configuration values
func (cm *StatelessConfigurationManager) SetConfig(key string, value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.configs[key] = value
}

// LoadFromEnvironment loads config from environment variables
func (cm *StatelessConfigurationManager) LoadFromEnvironment() {
	envVars := []string{"PORT", "DATABASE_URL", "REDIS_URL", "JWT_SECRET", "ENVIRONMENT"}

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			cm.SetConfig(envVar, value)
		}
	}
}

// StatelessMetrics demonstrates metrics collection for stateless services
type StatelessMetrics struct {
	RequestCount   int64            `json:"request_count"`
	ResponseTimes  []time.Duration  `json:"response_times"`
	ErrorCount     int64            `json:"error_count"`
	ActiveSessions int64            `json:"active_sessions"`
	LastUpdated    time.Time        `json:"last_updated"`
	CustomMetrics  map[string]int64 `json:"custom_metrics"`
	mu             sync.RWMutex
}

// NewStatelessMetrics creates a metrics collector
func NewStatelessMetrics() *StatelessMetrics {
	return &StatelessMetrics{
		ResponseTimes: make([]time.Duration, 0),
		CustomMetrics: make(map[string]int64),
		LastUpdated:   time.Now(),
	}
}

// IncrementRequests increments request counter
func (m *StatelessMetrics) IncrementRequests() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RequestCount++
	m.LastUpdated = time.Now()
}

// RecordResponseTime records response time
func (m *StatelessMetrics) RecordResponseTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResponseTimes = append(m.ResponseTimes, duration)

	// Keep only last 1000 response times
	if len(m.ResponseTimes) > 1000 {
		m.ResponseTimes = m.ResponseTimes[1:]
	}
	m.LastUpdated = time.Now()
}

// IncrementErrors increments error counter
func (m *StatelessMetrics) IncrementErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorCount++
	m.LastUpdated = time.Now()
}

// GetMetrics returns current metrics
func (m *StatelessMetrics) GetMetrics() StatelessMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions
	metrics := StatelessMetrics{
		RequestCount:   m.RequestCount,
		ErrorCount:     m.ErrorCount,
		ActiveSessions: m.ActiveSessions,
		LastUpdated:    m.LastUpdated,
		CustomMetrics:  make(map[string]int64),
	}

	// Copy custom metrics
	for k, v := range m.CustomMetrics {
		metrics.CustomMetrics[k] = v
	}

	// Copy response times
	metrics.ResponseTimes = make([]time.Duration, len(m.ResponseTimes))
	copy(metrics.ResponseTimes, m.ResponseTimes)

	return metrics
}

// StatelessPerformanceBenchmark demonstrates performance testing
type StatelessPerformanceBenchmark struct {
	service  *StatelessService
	metrics  *StatelessMetrics
	workers  int
	requests int
}

// NewStatelessPerformanceBenchmark creates a benchmark
func NewStatelessPerformanceBenchmark(service *StatelessService, workers, requests int) *StatelessPerformanceBenchmark {
	return &StatelessPerformanceBenchmark{
		service:  service,
		metrics:  NewStatelessMetrics(),
		workers:  workers,
		requests: requests,
	}
}

// RunBenchmark executes performance test
func (b *StatelessPerformanceBenchmark) RunBenchmark() {
	requestsPerWorker := b.requests / b.workers
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < b.workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < requestsPerWorker; j++ {
				requestStart := time.Now()

				req := UserProcessRequest{
					Username: fmt.Sprintf("user%d_%d", workerID, j),
					Email:    fmt.Sprintf("user%d_%d@example.com", workerID, j),
					Action:   "benchmark_test",
				}

				_, err := b.service.ProcessUserRequest(context.Background(), req)

				duration := time.Since(requestStart)
				b.metrics.RecordResponseTime(duration)
				b.metrics.IncrementRequests()

				if err != nil {
					b.metrics.IncrementErrors()
				}
			}
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(start)

	log.Printf("Benchmark completed: %d requests in %v", b.requests, totalDuration)
	log.Printf("Requests per second: %.2f", float64(b.requests)/totalDuration.Seconds())
}

// StatelessServiceDemo demonstrates comprehensive stateless service usage
func StatelessServiceDemo() {
	// Load configuration
	config := LoadStatelessConfigFromEnv()

	// Create service
	service := NewStatelessService(config)

	// Setup routes
	mux := service.SetupStatelessRoutes()

	// Add middleware for metrics
	metrics := NewStatelessMetrics()
	handler := addStatelessMetricsMiddleware(mux, metrics)

	// Start server (would be done in actual deployment)
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler,
	}

	log.Printf("Starting stateless service on port %s", config.Port)

	// Demo scaling configuration
	HorizontalStatelessScalingDemo()

	// Demo configuration management
	configManager := NewStatelessConfigurationManager()
	configManager.LoadFromEnvironment()

	// Demo performance testing
	benchmark := NewStatelessPerformanceBenchmark(service, 10, 1000)
	benchmark.RunBenchmark()

	// Print final metrics
	finalMetrics := metrics.GetMetrics()
	log.Printf("Final metrics: %+v", finalMetrics)

	// In production, you would call server.ListenAndServe()
	log.Printf("Server configuration ready: %+v", server.Addr)
}

// addStatelessMetricsMiddleware adds metrics collection to HTTP handlers
func addStatelessMetricsMiddleware(next http.Handler, metrics *StatelessMetrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		metrics.IncrementRequests()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Record response time
		metrics.RecordResponseTime(time.Since(start))
	})
}

// FAANG Interview Discussion Points:

// 1. Stateless Design Benefits:
//    - Horizontal scaling without session affinity
//    - Rolling deployments without downtime
//    - Fault tolerance through instance replacement
//    - Load balancer flexibility

// 2. External State Management:
//    - Session storage in Redis/database
//    - Configuration from environment/config service
//    - No local file system dependencies
//    - Shared cache across instances

// 3. Health Check Strategies:
//    - Liveness vs Readiness probes
//    - Dependency health verification
//    - Graceful shutdown handling
//    - Circuit breaker integration

// 4. Session Handling Patterns:
//    - JWT tokens for stateless auth
//    - Redis session storage
//    - Database session management
//    - Session replication strategies

// 5. Configuration Management:
//    - 12-factor app principles
//    - Environment-based config
//    - Secret management
//    - Feature flags integration

// Performance Considerations:
// - External state access latency
// - Connection pooling for Redis/DB
// - Cache warming strategies
// - Session cleanup policies

// Real-world Use Cases:
// - Microservices architecture
// - API gateways
// - Web application backends
// - Container orchestration (Kubernetes)

// Scaling Challenges:
// - Session consistency across instances
// - Database connection limits
// - External service dependencies
// - Configuration synchronization

// Interview Questions to Discuss:
// 1. How do you handle session state in a stateless service?
// 2. What are the trade-offs between stateful and stateless services?
// 3. How do you implement health checks for load balancers?
// 4. What configuration patterns work best for microservices?
// 5. How do you measure and monitor stateless service performance?
