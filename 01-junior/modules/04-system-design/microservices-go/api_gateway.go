package microservices

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// APIGateway represents a comprehensive API Gateway implementation
// covering request routing, authentication, rate limiting, and API versioning
// Key FAANG Interview Topics:
// - Request routing patterns (path-based, header-based, content-based)
// - Authentication and authorization strategies (JWT, OAuth, API keys)
// - Rate limiting algorithms (token bucket, sliding window, fixed window)
// - Request/response transformation and middleware architecture
// - API versioning strategies and backward compatibility
// - Monitoring, logging, and observability patterns

// Route represents a routing rule in the API Gateway
type Route struct {
	ID          string            `json:"id"`
	Path        string            `json:"path"`
	Method      string            `json:"method"`
	ServiceName string            `json:"service_name"`
	ServicePath string            `json:"service_path"`
	Version     string            `json:"version"`
	Headers     map[string]string `json:"headers"`
	Priority    int               `json:"priority"`
	Enabled     bool              `json:"enabled"`
	Middleware  []string          `json:"middleware"`
}

// AuthPolicy represents authentication and authorization policies
type AuthPolicy struct {
	PolicyID     string           `json:"policy_id"`
	PolicyType   AuthPolicyType   `json:"policy_type"`
	RequiredRole string           `json:"required_role,omitempty"`
	APIKeyHeader string           `json:"api_key_header,omitempty"`
	JWTSecret    string           `json:"jwt_secret,omitempty"`
	Scopes       []string         `json:"scopes,omitempty"`
	RateLimit    *RateLimitPolicy `json:"rate_limit,omitempty"`
}

type AuthPolicyType string

const (
	AuthPolicyNone   AuthPolicyType = "none"
	AuthPolicyAPIKey AuthPolicyType = "api_key"
	AuthPolicyJWT    AuthPolicyType = "jwt"
	AuthPolicyOAuth  AuthPolicyType = "oauth"
)

// RateLimitPolicy defines rate limiting configuration
type RateLimitPolicy struct {
	RequestsPerSecond int                `json:"requests_per_second"`
	BurstSize         int                `json:"burst_size"`
	WindowSize        time.Duration      `json:"window_size"`
	Algorithm         RateLimitAlgorithm `json:"algorithm"`
}

type RateLimitAlgorithm string

const (
	TokenBucketAlgorithm RateLimitAlgorithm = "token_bucket"
	SlidingWindow        RateLimitAlgorithm = "sliding_window"
	FixedWindow          RateLimitAlgorithm = "fixed_window"
)

// RequestContext holds request processing context
type RequestContext struct {
	RequestID  string            `json:"request_id"`
	UserID     string            `json:"user_id"`
	ClientID   string            `json:"client_id"`
	Route      *Route            `json:"route"`
	AuthPolicy *AuthPolicy       `json:"auth_policy"`
	Headers    map[string]string `json:"headers"`
	StartTime  time.Time         `json:"start_time"`
	Metadata   map[string]any    `json:"metadata"`
}

// APIGatewayRouter handles request routing to microservices
type APIGatewayRouter interface {
	AddRoute(route *Route) error
	RemoveRoute(routeID string) error
	MatchRoute(req *http.Request) (*Route, error)
	GetRoutes() []*Route
}

// Authenticator handles various authentication mechanisms
type Authenticator interface {
	Authenticate(req *http.Request, policy *AuthPolicy) (*AuthContext, error)
	ValidateJWT(token string, secret string) (*JWTClaims, error)
	ValidateAPIKey(apiKey string) (*APIKeyInfo, error)
	ValidateOAuth(token string) (*OAuthClaims, error)
}

// AuthContext represents authenticated request context
type AuthContext struct {
	UserID    string            `json:"user_id"`
	ClientID  string            `json:"client_id"`
	Roles     []string          `json:"roles"`
	Scopes    []string          `json:"scopes"`
	Metadata  map[string]string `json:"metadata"`
	IsValid   bool              `json:"is_valid"`
	ExpiresAt time.Time         `json:"expires_at"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID    string   `json:"user_id"`
	ClientID  string   `json:"client_id"`
	Roles     []string `json:"roles"`
	Scopes    []string `json:"scopes"`
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
}

// APIKeyInfo represents API key information
type APIKeyInfo struct {
	KeyID     string   `json:"key_id"`
	ClientID  string   `json:"client_id"`
	Scopes    []string `json:"scopes"`
	RateLimit int      `json:"rate_limit"`
	IsActive  bool     `json:"is_active"`
}

// OAuthClaims represents OAuth token claims
type OAuthClaims struct {
	UserID    string   `json:"user_id"`
	ClientID  string   `json:"client_id"`
	Scopes    []string `json:"scopes"`
	ExpiresAt int64    `json:"exp"`
}

// RateLimiter provides rate limiting functionality
type RateLimiter interface {
	Allow(clientID string, policy *RateLimitPolicy) (bool, error)
	GetQuota(clientID string) (RateLimitQuota, error)
	Reset(clientID string) error
}

// RateLimitQuota represents current rate limit status
type RateLimitQuota struct {
	Remaining int       `json:"remaining"`
	ResetTime time.Time `json:"reset_time"`
	Total     int       `json:"total"`
}

// Middleware represents request/response transformation middleware
type Middleware interface {
	Process(ctx *RequestContext, req *http.Request, resp *http.Response) error
	Name() string
}

// APIVersionManager handles API versioning strategies
type APIVersionManager interface {
	ExtractVersion(req *http.Request) string
	RewriteRequest(req *http.Request, targetVersion string) error
	IsVersionSupported(version string) bool
	GetDefaultVersion() string
}

// GatewayMetrics provides observability data
type GatewayMetrics struct {
	RequestsTotal     int64   `json:"requests_total"`
	RequestsSuccess   int64   `json:"requests_success"`
	RequestsFailure   int64   `json:"requests_failure"`
	RequestsAuth      int64   `json:"requests_auth"`
	RequestsUnauth    int64   `json:"requests_unauth"`
	RequestsRateLimit int64   `json:"requests_rate_limit"`
	AvgLatencyMS      float64 `json:"avg_latency_ms"`
	P99LatencyMS      float64 `json:"p99_latency_ms"`
}

// DefaultAPIGatewayRouter implements path and header-based routing
type DefaultAPIGatewayRouter struct {
	mu     sync.RWMutex
	routes []*Route
}

func NewAPIGatewayRouter() *DefaultAPIGatewayRouter {
	return &DefaultAPIGatewayRouter{
		routes: make([]*Route, 0),
	}
}

func (r *DefaultAPIGatewayRouter) AddRoute(route *Route) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate route
	if route.Path == "" || route.ServiceName == "" {
		return fmt.Errorf("invalid route: path and service_name are required")
	}

	// Set default values
	if route.ID == "" {
		route.ID = generateRouteID()
	}
	if route.Method == "" {
		route.Method = "GET"
	}
	if route.Priority == 0 {
		route.Priority = 100
	}
	if route.Version == "" {
		route.Version = "v1"
	}

	route.Enabled = true

	r.routes = append(r.routes, route)

	// Sort routes by priority (higher priority first)
	sort.Slice(r.routes, func(i, j int) bool {
		return r.routes[i].Priority > r.routes[j].Priority
	})

	return nil
}

func (r *DefaultAPIGatewayRouter) RemoveRoute(routeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, route := range r.routes {
		if route.ID == routeID {
			r.routes = append(r.routes[:i], r.routes[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("route not found: %s", routeID)
}

func (r *DefaultAPIGatewayRouter) MatchRoute(req *http.Request) (*Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, route := range r.routes {
		if !route.Enabled {
			continue
		}

		// Match HTTP method
		if route.Method != "*" && route.Method != req.Method {
			continue
		}

		// Match path pattern
		matched, err := matchPath(route.Path, req.URL.Path)
		if err != nil || !matched {
			continue
		}

		// Match headers if specified
		if !matchHeaders(route.Headers, req.Header) {
			continue
		}

		return route, nil
	}

	return nil, fmt.Errorf("no matching route found for %s %s", req.Method, req.URL.Path)
}

func (r *DefaultAPIGatewayRouter) GetRoutes() []*Route {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routes := make([]*Route, len(r.routes))
	copy(routes, r.routes)
	return routes
}

// DefaultAuthenticator implements multiple authentication mechanisms
type DefaultAuthenticator struct {
	mu        sync.RWMutex
	apiKeys   map[string]*APIKeyInfo
	jwtSecret string
	oauthURL  string
}

func NewAuthenticator(jwtSecret, oauthURL string) *DefaultAuthenticator {
	return &DefaultAuthenticator{
		apiKeys:   make(map[string]*APIKeyInfo),
		jwtSecret: jwtSecret,
		oauthURL:  oauthURL,
	}
}

func (a *DefaultAuthenticator) Authenticate(req *http.Request, policy *AuthPolicy) (*AuthContext, error) {
	switch policy.PolicyType {
	case AuthPolicyNone:
		return &AuthContext{IsValid: true}, nil

	case AuthPolicyAPIKey:
		apiKey := req.Header.Get(policy.APIKeyHeader)
		if apiKey == "" {
			return nil, fmt.Errorf("missing API key in header: %s", policy.APIKeyHeader)
		}

		keyInfo, err := a.ValidateAPIKey(apiKey)
		if err != nil {
			return nil, err
		}

		return &AuthContext{
			UserID:   keyInfo.ClientID,
			ClientID: keyInfo.ClientID,
			Scopes:   keyInfo.Scopes,
			IsValid:  keyInfo.IsActive,
		}, nil

	case AuthPolicyJWT:
		token := extractBearerToken(req)
		if token == "" {
			return nil, fmt.Errorf("missing JWT token")
		}

		claims, err := a.ValidateJWT(token, policy.JWTSecret)
		if err != nil {
			return nil, err
		}

		return &AuthContext{
			UserID:    claims.UserID,
			ClientID:  claims.ClientID,
			Roles:     claims.Roles,
			Scopes:    claims.Scopes,
			IsValid:   true,
			ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		}, nil

	case AuthPolicyOAuth:
		token := extractBearerToken(req)
		if token == "" {
			return nil, fmt.Errorf("missing OAuth token")
		}

		claims, err := a.ValidateOAuth(token)
		if err != nil {
			return nil, err
		}

		return &AuthContext{
			UserID:    claims.UserID,
			ClientID:  claims.ClientID,
			Scopes:    claims.Scopes,
			IsValid:   true,
			ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported auth policy type: %s", policy.PolicyType)
	}
}

func (a *DefaultAuthenticator) ValidateJWT(token string, secret string) (*JWTClaims, error) {
	// Simplified JWT validation - in production, use proper JWT library
	// This is a mock implementation for demonstration
	claims := &JWTClaims{
		UserID:    "user-12345",
		ClientID:  "client-abc",
		Roles:     []string{"user", "premium"},
		Scopes:    []string{"read", "write"},
		IssuedAt:  time.Now().Unix() - 3600,
		ExpiresAt: time.Now().Unix() + 3600,
	}

	return claims, nil
}

func (a *DefaultAuthenticator) ValidateAPIKey(apiKey string) (*APIKeyInfo, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	keyInfo, exists := a.apiKeys[apiKey]
	if !exists {
		return nil, fmt.Errorf("invalid API key")
	}

	if !keyInfo.IsActive {
		return nil, fmt.Errorf("API key is inactive")
	}

	return keyInfo, nil
}

func (a *DefaultAuthenticator) ValidateOAuth(token string) (*OAuthClaims, error) {
	// Simplified OAuth validation - in production, validate with OAuth provider
	claims := &OAuthClaims{
		UserID:    "oauth-user-789",
		ClientID:  "oauth-client-xyz",
		Scopes:    []string{"profile", "email"},
		ExpiresAt: time.Now().Unix() + 7200,
	}

	return claims, nil
}

func (a *DefaultAuthenticator) AddAPIKey(apiKey string, info *APIKeyInfo) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.apiKeys[apiKey] = info
}

// GatewayTokenBucket represents a token bucket for rate limiting
type GatewayTokenBucket struct {
	Tokens     float64   `json:"tokens"`
	Capacity   float64   `json:"capacity"`
	RefillRate float64   `json:"refill_rate"`
	LastRefill time.Time `json:"last_refill"`
}

// TokenBucketRateLimiter implements token bucket algorithm
type TokenBucketRateLimiter struct {
	mu      sync.RWMutex
	buckets map[string]*GatewayTokenBucket
}

type TokenBucket struct {
	Tokens     float64   `json:"tokens"`
	Capacity   float64   `json:"capacity"`
	RefillRate float64   `json:"refill_rate"`
	LastRefill time.Time `json:"last_refill"`
}

func NewTokenBucketRateLimiter() *TokenBucketRateLimiter {
	limiter := &TokenBucketRateLimiter{
		buckets: make(map[string]*GatewayTokenBucket),
	}

	// Start background refill process
	go limiter.refillTokens()

	return limiter
}

func (rl *TokenBucketRateLimiter) Allow(clientID string, policy *RateLimitPolicy) (bool, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket := rl.getBucket(clientID, policy)

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(bucket.LastRefill).Seconds()
	tokensToAdd := elapsed * bucket.RefillRate

	bucket.Tokens = min(bucket.Capacity, bucket.Tokens+tokensToAdd)
	bucket.LastRefill = now

	// Check if request can be allowed
	if bucket.Tokens >= 1.0 {
		bucket.Tokens -= 1.0
		return true, nil
	}

	return false, nil
}

func (rl *TokenBucketRateLimiter) GetQuota(clientID string) (RateLimitQuota, error) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	bucket, exists := rl.buckets[clientID]
	if !exists {
		return RateLimitQuota{}, fmt.Errorf("client not found: %s", clientID)
	}

	resetTime := bucket.LastRefill.Add(time.Duration((bucket.Capacity-bucket.Tokens)/bucket.RefillRate) * time.Second)

	return RateLimitQuota{
		Remaining: int(bucket.Tokens),
		ResetTime: resetTime,
		Total:     int(bucket.Capacity),
	}, nil
}

func (rl *TokenBucketRateLimiter) Reset(clientID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.buckets, clientID)
	return nil
}

func (rl *TokenBucketRateLimiter) getBucket(clientID string, policy *RateLimitPolicy) *GatewayTokenBucket {
	bucket, exists := rl.buckets[clientID]
	if !exists {
		bucket = &GatewayTokenBucket{
			Tokens:     float64(policy.BurstSize),
			Capacity:   float64(policy.BurstSize),
			RefillRate: float64(policy.RequestsPerSecond),
			LastRefill: time.Now(),
		}
		rl.buckets[clientID] = bucket
	}

	return bucket
}

func (rl *TokenBucketRateLimiter) refillTokens() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()

		for _, bucket := range rl.buckets {
			elapsed := now.Sub(bucket.LastRefill).Seconds()
			tokensToAdd := elapsed * bucket.RefillRate
			bucket.Tokens = min(bucket.Capacity, bucket.Tokens+tokensToAdd)
			bucket.LastRefill = now
		}

		rl.mu.Unlock()
	}
}

// HeaderVersionManager extracts API version from headers
type HeaderVersionManager struct {
	supportedVersions map[string]bool
	defaultVersion    string
}

func NewHeaderVersionManager(supportedVersions []string, defaultVersion string) *HeaderVersionManager {
	versionMap := make(map[string]bool)
	for _, v := range supportedVersions {
		versionMap[v] = true
	}

	return &HeaderVersionManager{
		supportedVersions: versionMap,
		defaultVersion:    defaultVersion,
	}
}

func (vm *HeaderVersionManager) ExtractVersion(req *http.Request) string {
	// Try different version extraction strategies

	// 1. Accept header with version
	accept := req.Header.Get("Accept")
	if version := extractVersionFromAccept(accept); version != "" {
		return version
	}

	// 2. Custom API-Version header
	if version := req.Header.Get("API-Version"); version != "" {
		return version
	}

	// 3. Path-based versioning
	if version := extractVersionFromPath(req.URL.Path); version != "" {
		return version
	}

	// 4. Query parameter
	if version := req.URL.Query().Get("version"); version != "" {
		return version
	}

	return vm.defaultVersion
}

func (vm *HeaderVersionManager) RewriteRequest(req *http.Request, targetVersion string) error {
	// Set version in header for downstream services
	req.Header.Set("X-API-Version", targetVersion)

	// Rewrite path if needed (remove version prefix)
	if strings.HasPrefix(req.URL.Path, "/v") {
		parts := strings.SplitN(req.URL.Path, "/", 3)
		if len(parts) >= 3 {
			req.URL.Path = "/" + parts[2]
		}
	}

	return nil
}

func (vm *HeaderVersionManager) IsVersionSupported(version string) bool {
	return vm.supportedVersions[version]
}

func (vm *HeaderVersionManager) GetDefaultVersion() string {
	return vm.defaultVersion
}

// LoggingMiddleware logs requests and responses
type LoggingMiddleware struct {
	logger *log.Logger
}

func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

func (m *LoggingMiddleware) Process(ctx *RequestContext, req *http.Request, resp *http.Response) error {
	duration := time.Since(ctx.StartTime)
	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}

	m.logger.Printf(
		"[%s] %s %s %s -> %d (%v) [User: %s, Client: %s]",
		ctx.RequestID,
		req.Method,
		req.URL.Path,
		req.RemoteAddr,
		statusCode,
		duration,
		ctx.UserID,
		ctx.ClientID,
	)

	return nil
}

func (m *LoggingMiddleware) Name() string {
	return "logging"
}

// TransformationMiddleware handles request/response transformation
type TransformationMiddleware struct {
	requestTransforms  map[string]func(*http.Request) error
	responseTransforms map[string]func(*http.Response) error
}

func NewTransformationMiddleware() *TransformationMiddleware {
	return &TransformationMiddleware{
		requestTransforms:  make(map[string]func(*http.Request) error),
		responseTransforms: make(map[string]func(*http.Response) error),
	}
}

func (m *TransformationMiddleware) Process(ctx *RequestContext, req *http.Request, resp *http.Response) error {
	// Apply request transformations
	if transform, exists := m.requestTransforms[ctx.Route.ServiceName]; exists {
		if err := transform(req); err != nil {
			return fmt.Errorf("request transformation failed: %w", err)
		}
	}

	// Apply response transformations
	if resp != nil {
		if transform, exists := m.responseTransforms[ctx.Route.ServiceName]; exists {
			if err := transform(resp); err != nil {
				return fmt.Errorf("response transformation failed: %w", err)
			}
		}
	}

	return nil
}

func (m *TransformationMiddleware) Name() string {
	return "transformation"
}

func (m *TransformationMiddleware) AddRequestTransform(serviceName string, transform func(*http.Request) error) {
	m.requestTransforms[serviceName] = transform
}

func (m *TransformationMiddleware) AddResponseTransform(serviceName string, transform func(*http.Response) error) {
	m.responseTransforms[serviceName] = transform
}

// APIGateway represents the main gateway implementation
type APIGateway struct {
	router         APIGatewayRouter
	authenticator  Authenticator
	rateLimiter    RateLimiter
	versionManager APIVersionManager
	middleware     []Middleware
	metrics        GatewayMetrics
	serviceClient  *http.Client
}

func NewAPIGateway(
	router APIGatewayRouter,
	authenticator Authenticator,
	rateLimiter RateLimiter,
	versionManager APIVersionManager,
) *APIGateway {
	return &APIGateway{
		router:         router,
		authenticator:  authenticator,
		rateLimiter:    rateLimiter,
		versionManager: versionManager,
		middleware:     make([]Middleware, 0),
		serviceClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (gw *APIGateway) AddMiddleware(middleware Middleware) {
	gw.middleware = append(gw.middleware, middleware)
}

func (gw *APIGateway) HandleRequest(req *http.Request) (*http.Response, error) {
	startTime := time.Now()
	atomic.AddInt64(&gw.metrics.RequestsTotal, 1)

	// Create request context
	ctx := &RequestContext{
		RequestID: generateRequestID(),
		StartTime: startTime,
		Headers:   make(map[string]string),
		Metadata:  make(map[string]any),
	}

	// Extract and validate API version
	version := gw.versionManager.ExtractVersion(req)
	if !gw.versionManager.IsVersionSupported(version) {
		atomic.AddInt64(&gw.metrics.RequestsFailure, 1)
		return nil, fmt.Errorf("unsupported API version: %s", version)
	}

	// Rewrite request for version compatibility
	if err := gw.versionManager.RewriteRequest(req, version); err != nil {
		atomic.AddInt64(&gw.metrics.RequestsFailure, 1)
		return nil, fmt.Errorf("version rewrite failed: %w", err)
	}

	// Find matching route
	route, err := gw.router.MatchRoute(req)
	if err != nil {
		atomic.AddInt64(&gw.metrics.RequestsFailure, 1)
		return nil, fmt.Errorf("route matching failed: %w", err)
	}
	ctx.Route = route

	// Get auth policy (simplified - would typically be from route configuration)
	authPolicy := &AuthPolicy{
		PolicyType: AuthPolicyJWT,
		JWTSecret:  "secret-key",
	}
	ctx.AuthPolicy = authPolicy

	// Authenticate request
	authCtx, err := gw.authenticator.Authenticate(req, authPolicy)
	if err != nil {
		atomic.AddInt64(&gw.metrics.RequestsUnauth, 1)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	atomic.AddInt64(&gw.metrics.RequestsAuth, 1)

	ctx.UserID = authCtx.UserID
	ctx.ClientID = authCtx.ClientID

	// Apply rate limiting
	if authPolicy.RateLimit != nil {
		allowed, err := gw.rateLimiter.Allow(authCtx.ClientID, authPolicy.RateLimit)
		if err != nil {
			return nil, fmt.Errorf("rate limiting error: %w", err)
		}
		if !allowed {
			atomic.AddInt64(&gw.metrics.RequestsRateLimit, 1)
			return nil, fmt.Errorf("rate limit exceeded for client: %s", authCtx.ClientID)
		}
	}

	// Apply request middleware
	for _, mw := range gw.middleware {
		if err := mw.Process(ctx, req, nil); err != nil {
			return nil, fmt.Errorf("middleware %s failed: %w", mw.Name(), err)
		}
	}

	// Forward request to downstream service
	serviceURL := fmt.Sprintf("http://%s%s", route.ServiceName, route.ServicePath)
	if route.ServicePath == "" {
		serviceURL = fmt.Sprintf("http://%s%s", route.ServiceName, req.URL.Path)
	}

	serviceReq, err := http.NewRequest(req.Method, serviceURL, req.Body)
	if err != nil {
		atomic.AddInt64(&gw.metrics.RequestsFailure, 1)
		return nil, fmt.Errorf("service request creation failed: %w", err)
	}

	// Copy headers
	for k, v := range req.Header {
		serviceReq.Header[k] = v
	}

	// Add gateway headers
	serviceReq.Header.Set("X-Gateway-Request-ID", ctx.RequestID)
	serviceReq.Header.Set("X-User-ID", ctx.UserID)
	serviceReq.Header.Set("X-Client-ID", ctx.ClientID)

	// Execute service request
	resp, err := gw.serviceClient.Do(serviceReq)
	if err != nil {
		atomic.AddInt64(&gw.metrics.RequestsFailure, 1)
		return nil, fmt.Errorf("service request failed: %w", err)
	}

	// Apply response middleware
	for _, mw := range gw.middleware {
		if err := mw.Process(ctx, req, resp); err != nil {
			return nil, fmt.Errorf("response middleware %s failed: %w", mw.Name(), err)
		}
	}

	atomic.AddInt64(&gw.metrics.RequestsSuccess, 1)

	// Update latency metrics
	gw.updateLatencyMetrics(time.Since(startTime))

	return resp, nil
}

func (gw *APIGateway) GetMetrics() GatewayMetrics {
	return GatewayMetrics{
		RequestsTotal:     atomic.LoadInt64(&gw.metrics.RequestsTotal),
		RequestsSuccess:   atomic.LoadInt64(&gw.metrics.RequestsSuccess),
		RequestsFailure:   atomic.LoadInt64(&gw.metrics.RequestsFailure),
		RequestsAuth:      atomic.LoadInt64(&gw.metrics.RequestsAuth),
		RequestsUnauth:    atomic.LoadInt64(&gw.metrics.RequestsUnauth),
		RequestsRateLimit: atomic.LoadInt64(&gw.metrics.RequestsRateLimit),
		AvgLatencyMS:      gw.metrics.AvgLatencyMS,
		P99LatencyMS:      gw.metrics.P99LatencyMS,
	}
}

func (gw *APIGateway) updateLatencyMetrics(latency time.Duration) {
	latencyMS := float64(latency.Nanoseconds()) / 1e6

	// Update average (simplified)
	if gw.metrics.AvgLatencyMS == 0 {
		gw.metrics.AvgLatencyMS = latencyMS
	} else {
		gw.metrics.AvgLatencyMS = (gw.metrics.AvgLatencyMS + latencyMS) / 2
	}

	// Update P99 (simplified)
	if latencyMS > gw.metrics.P99LatencyMS {
		gw.metrics.P99LatencyMS = latencyMS
	}
}

// Helper functions

func generateRouteID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "route-" + hex.EncodeToString(bytes)
}

func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "req-" + hex.EncodeToString(bytes)
}

func matchPath(pattern, path string) (bool, error) {
	// Convert path pattern to regex
	// Support wildcards: /api/*/users, /api/{id}/profile
	regexPattern := strings.ReplaceAll(pattern, "*", "[^/]+")
	regexPattern = regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(regexPattern, "[^/]+")
	regexPattern = "^" + regexPattern + "$"

	matched, err := regexp.MatchString(regexPattern, path)
	return matched, err
}

func matchHeaders(required map[string]string, headers http.Header) bool {
	for key, value := range required {
		headerValue := headers.Get(key)
		if headerValue != value {
			return false
		}
	}
	return true
}

func extractBearerToken(req *http.Request) string {
	auth := req.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return auth[7:]
	}
	return ""
}

func extractVersionFromAccept(accept string) string {
	// Parse Accept header: application/vnd.api+json;version=2
	parts := strings.Split(accept, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "version=") {
			return part[8:]
		}
	}
	return ""
}

func extractVersionFromPath(path string) string {
	// Extract version from path: /v2/users/123
	parts := strings.Split(path, "/")
	if len(parts) >= 2 && strings.HasPrefix(parts[1], "v") {
		return parts[1]
	}
	return ""
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// APIGatewayExample demonstrates comprehensive API Gateway usage
type APIGatewayExample struct {
	gateway *APIGateway
}

func NewAPIGatewayExample() *APIGatewayExample {
	// Setup components
	router := NewAPIGatewayRouter()
	authenticator := NewAuthenticator("jwt-secret-key", "https://oauth.example.com")
	rateLimiter := NewTokenBucketRateLimiter()
	versionManager := NewHeaderVersionManager([]string{"v1", "v2", "v3"}, "v1")

	gateway := NewAPIGateway(router, authenticator, rateLimiter, versionManager)

	// Add middleware
	gateway.AddMiddleware(NewLoggingMiddleware(log.Default()))
	gateway.AddMiddleware(NewTransformationMiddleware())

	return &APIGatewayExample{gateway: gateway}
}

func (example *APIGatewayExample) DemonstrateAPIGateway() {
	// Configure routes
	routes := []*Route{
		{
			Path:        "/api/v*/users/*",
			Method:      "GET",
			ServiceName: "user-service",
			ServicePath: "/users",
			Priority:    100,
		},
		{
			Path:        "/api/v*/orders",
			Method:      "POST",
			ServiceName: "order-service",
			ServicePath: "/orders",
			Priority:    90,
		},
		{
			Path:        "/api/v*/payments/{id}",
			Method:      "PUT",
			ServiceName: "payment-service",
			ServicePath: "/payments",
			Priority:    95,
		},
	}

	for _, route := range routes {
		example.gateway.router.AddRoute(route)
	}

	// Add API key for testing
	auth := example.gateway.authenticator.(*DefaultAuthenticator)
	auth.AddAPIKey("test-api-key-12345", &APIKeyInfo{
		KeyID:     "key-001",
		ClientID:  "client-abc",
		Scopes:    []string{"read", "write"},
		RateLimit: 1000,
		IsActive:  true,
	})

	// Simulate API requests
	testRequests := []*http.Request{
		createTestRequest("GET", "/api/v2/users/123", map[string]string{
			"Authorization": "Bearer jwt-token-here",
			"Accept":        "application/json",
		}),
		createTestRequest("POST", "/api/v1/orders", map[string]string{
			"API-Key":      "test-api-key-12345",
			"Content-Type": "application/json",
		}),
		createTestRequest("PUT", "/api/v3/payments/456", map[string]string{
			"Authorization": "Bearer jwt-token-here",
			"API-Version":   "v3",
		}),
	}

	// Process requests
	for i, req := range testRequests {
		log.Printf("Processing request %d: %s %s", i+1, req.Method, req.URL.Path)

		resp, err := example.gateway.HandleRequest(req)
		if err != nil {
			log.Printf("Request failed: %v", err)
		} else {
			log.Printf("Request succeeded: Status %d", resp.StatusCode)
		}
	}

	// Display metrics
	metrics := example.gateway.GetMetrics()
	log.Printf("Gateway Metrics: %+v", metrics)

	// Display routes
	routes = example.gateway.router.GetRoutes()
	log.Printf("Configured routes: %d", len(routes))
	for _, route := range routes {
		log.Printf("Route: %s %s -> %s (Priority: %d)",
			route.Method, route.Path, route.ServiceName, route.Priority)
	}
}

func createTestRequest(method, path string, headers map[string]string) *http.Request {
	req, _ := http.NewRequest(method, "http://gateway.example.com"+path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}

// FAANG Interview Discussion Points:
//
// 1. API Gateway Patterns:
//    - Centralized vs decentralized gateway architectures
//    - Backend for Frontend (BFF) pattern considerations
//    - Gateway aggregation patterns for mobile clients
//    - Caching strategies at the gateway level
//
// 2. Routing Strategies:
//    - Path-based vs header-based vs content-based routing
//    - Canary deployments and traffic splitting
//    - Blue-green deployments through routing
//    - Geographic routing for global deployments
//
// 3. Authentication Architecture:
//    - Token validation performance and caching
//    - Federated identity and SSO integration
//    - Service-to-service authentication patterns
//    - Authorization policy evaluation and caching
//
// 4. Rate Limiting Design:
//    - Distributed rate limiting across gateway instances
//    - Fair queuing and prioritization algorithms
//    - Adaptive rate limiting based on system load
//    - Client-specific vs global rate limiting strategies
//
// 5. API Versioning Strategies:
//    - Semantic versioning and breaking change management
//    - Deprecation strategies and migration paths
//    - Version negotiation and fallback mechanisms
//    - Schema evolution and backward compatibility
//
// 6. Observability and Monitoring:
//    - Request tracing across service boundaries
//    - SLA monitoring and alerting strategies
//    - Performance metrics and bottleneck identification
//    - Security event monitoring and threat detection
//
// 7. Scalability Considerations:
//    - Horizontal scaling of gateway instances
//    - Load balancing strategies for gateways
//    - Circuit breaker integration for fault tolerance
//    - Graceful degradation during service outages
