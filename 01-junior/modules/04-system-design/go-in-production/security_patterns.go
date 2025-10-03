package production

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ======================== AUTHENTICATION AND AUTHORIZATION ========================

// AuthProvider handles authentication operations
type AuthProvider struct {
	mu           sync.RWMutex
	users        map[string]*User
	sessions     map[string]*Session
	jwtSecretKey []byte
	sessionTTL   time.Duration
}

// User represents a system user
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never serialize passwords
	Roles        []string  `json:"roles"`
	Permissions  []string  `json:"permissions"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	LastLogin    time.Time `json:"last_login"`
	MFAEnabled   bool      `json:"mfa_enabled"`
	MFASecret    string    `json:"-"` // TOTP secret
}

// Session represents an authenticated session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	IsActive  bool      `json:"is_active"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	ExpiresAt   int64    `json:"exp"`
	IssuedAt    int64    `json:"iat"`
	Issuer      string   `json:"iss"`
}

// NewAuthProvider creates a new authentication provider
func NewAuthProvider(jwtSecret string, sessionTTL time.Duration) *AuthProvider {
	return &AuthProvider{
		users:        make(map[string]*User),
		sessions:     make(map[string]*Session),
		jwtSecretKey: []byte(jwtSecret),
		sessionTTL:   sessionTTL,
	}
}

// CreateUser creates a new user with hashed password
func (ap *AuthProvider) CreateUser(username, email, password string, roles []string) (*User, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	// Check if user already exists
	for _, user := range ap.users {
		if user.Username == username || user.Email == email {
			return nil, fmt.Errorf("user already exists")
		}
	}

	// Hash password (simple implementation)
	hash := sha256.Sum256([]byte(password + "salt"))
	passwordHash := fmt.Sprintf("%x", hash)

	// Generate user ID
	userID := generateSecureID()

	user := &User{
		ID:           userID,
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
		Roles:        roles,
		Permissions:  resolvePermissions(roles),
		IsActive:     true,
		CreatedAt:    time.Now(),
		MFAEnabled:   false,
	}

	ap.users[userID] = user
	return user, nil
}

// Authenticate verifies user credentials and returns a session
func (ap *AuthProvider) Authenticate(username, password string) (*Session, error) {
	ap.mu.RLock()
	var user *User
	for _, u := range ap.users {
		if u.Username == username && u.IsActive {
			user = u
			break
		}
	}
	ap.mu.RUnlock()

	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password (simple implementation)
	hash := sha256.Sum256([]byte(password + "salt"))
	expectedHash := fmt.Sprintf("%x", hash)
	if user.PasswordHash != expectedHash {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Create session
	session := &Session{
		ID:        generateSecureID(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ap.sessionTTL),
		IsActive:  true,
	}

	ap.mu.Lock()
	ap.sessions[session.ID] = session
	user.LastLogin = time.Now()
	ap.mu.Unlock()

	return session, nil
}

// ValidateSession validates a session token
func (ap *AuthProvider) ValidateSession(sessionID string) (*User, error) {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	session, exists := ap.sessions[sessionID]
	if !exists || !session.IsActive || time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("invalid or expired session")
	}

	user, exists := ap.users[session.UserID]
	if !exists || !user.IsActive {
		return nil, fmt.Errorf("user not found or inactive")
	}

	return user, nil
}

// GenerateJWT generates a JWT token for a user
func (ap *AuthProvider) GenerateJWT(user *User) (string, error) {
	claims := JWTClaims{
		UserID:      user.ID,
		Username:    user.Username,
		Roles:       user.Roles,
		Permissions: user.Permissions,
		ExpiresAt:   time.Now().Add(24 * time.Hour).Unix(),
		IssuedAt:    time.Now().Unix(),
		Issuer:      "production-system",
	}

	// Simple JWT implementation (in production, use a proper library)
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	message := header + "." + payload
	signature := ap.signHMAC(message)

	return message + "." + signature, nil
}

// ValidateJWT validates and parses a JWT token
func (ap *AuthProvider) ValidateJWT(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	message := parts[0] + "." + parts[1]
	expectedSignature := ap.signHMAC(message)

	if parts[2] != expectedSignature {
		return nil, fmt.Errorf("invalid token signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

func (ap *AuthProvider) signHMAC(message string) string {
	hash := sha256.New()
	hash.Write([]byte(message))
	hash.Write(ap.jwtSecretKey)
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

// ======================== ROLE-BASED ACCESS CONTROL (RBAC) ========================

// RBACManager manages role-based access control
type RBACManager struct {
	mu          sync.RWMutex
	roles       map[string]*Role
	permissions map[string]*Permission
}

// Role represents a role with associated permissions
type Role struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

// Permission represents a permission
type Permission struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
}

// NewRBACManager creates a new RBAC manager
func NewRBACManager() *RBACManager {
	rbac := &RBACManager{
		roles:       make(map[string]*Role),
		permissions: make(map[string]*Permission),
	}

	// Initialize default permissions and roles
	rbac.initializeDefaults()
	return rbac
}

// initializeDefaults sets up default roles and permissions
func (rbac *RBACManager) initializeDefaults() {
	// Default permissions
	permissions := []*Permission{
		{Name: "user:read", Description: "Read user data", Resource: "user", Action: "read"},
		{Name: "user:write", Description: "Write user data", Resource: "user", Action: "write"},
		{Name: "user:delete", Description: "Delete user data", Resource: "user", Action: "delete"},
		{Name: "admin:*", Description: "Full admin access", Resource: "*", Action: "*"},
		{Name: "config:read", Description: "Read configuration", Resource: "config", Action: "read"},
		{Name: "config:write", Description: "Write configuration", Resource: "config", Action: "write"},
		{Name: "metrics:read", Description: "Read metrics", Resource: "metrics", Action: "read"},
	}

	for _, perm := range permissions {
		rbac.permissions[perm.Name] = perm
	}

	// Default roles
	roles := []*Role{
		{
			Name:        "admin",
			Description: "System administrator with full access",
			Permissions: []string{"admin:*"},
			CreatedAt:   time.Now(),
		},
		{
			Name:        "user",
			Description: "Regular user with limited access",
			Permissions: []string{"user:read", "metrics:read"},
			CreatedAt:   time.Now(),
		},
		{
			Name:        "operator",
			Description: "Operations user with config access",
			Permissions: []string{"user:read", "config:read", "config:write", "metrics:read"},
			CreatedAt:   time.Now(),
		},
	}

	for _, role := range roles {
		rbac.roles[role.Name] = role
	}
}

// HasPermission checks if roles have a specific permission
func (rbac *RBACManager) HasPermission(roles []string, permission string) bool {
	rbac.mu.RLock()
	defer rbac.mu.RUnlock()

	for _, roleName := range roles {
		role, exists := rbac.roles[roleName]
		if !exists {
			continue
		}

		for _, perm := range role.Permissions {
			if perm == permission || perm == "admin:*" {
				return true
			}

			// Check wildcard permissions
			if strings.HasSuffix(perm, ":*") {
				resource := strings.TrimSuffix(perm, ":*")
				if strings.HasPrefix(permission, resource+":") {
					return true
				}
			}
		}
	}

	return false
}

// resolvePermissions converts roles to permissions list
func resolvePermissions(roles []string) []string {
	rbac := NewRBACManager()
	permissions := make(map[string]bool)

	for _, roleName := range roles {
		if role, exists := rbac.roles[roleName]; exists {
			for _, perm := range role.Permissions {
				permissions[perm] = true
			}
		}
	}

	result := make([]string, 0, len(permissions))
	for perm := range permissions {
		result = append(result, perm)
	}

	return result
}

// ======================== SSL/TLS CERTIFICATE MANAGEMENT ========================

// CertificateManager manages SSL/TLS certificates
type CertificateManager struct {
	mu           sync.RWMutex
	certificates map[string]*Certificate
	autoRenew    bool
	renewBefore  time.Duration
}

// Certificate represents an SSL/TLS certificate
type Certificate struct {
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	CertPEM   string    `json:"cert_pem"`
	KeyPEM    string    `json:"key_pem"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
	Issuer    string    `json:"issuer"`
	IsValid   bool      `json:"is_valid"`
}

// NewCertificateManager creates a new certificate manager
func NewCertificateManager(autoRenew bool) *CertificateManager {
	cm := &CertificateManager{
		certificates: make(map[string]*Certificate),
		autoRenew:    autoRenew,
		renewBefore:  30 * 24 * time.Hour, // Renew 30 days before expiry
	}

	if autoRenew {
		go cm.renewalWatcher()
	}

	return cm
}

// LoadCertificate loads a certificate from PEM data
func (cm *CertificateManager) LoadCertificate(name, domain, certPEM, keyPEM string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Parse certificate to get expiration
	certBlock, _ := pem.Decode([]byte(certPEM))
	if certBlock == nil {
		return fmt.Errorf("failed to decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	certificate := &Certificate{
		Name:      name,
		Domain:    domain,
		CertPEM:   certPEM,
		KeyPEM:    keyPEM,
		ExpiresAt: cert.NotAfter,
		IssuedAt:  cert.NotBefore,
		Issuer:    cert.Issuer.String(),
		IsValid:   time.Now().Before(cert.NotAfter),
	}

	cm.certificates[name] = certificate
	return nil
}

// GetTLSConfig returns TLS configuration for a domain
func (cm *CertificateManager) GetTLSConfig(domain string) (*tls.Config, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var cert *Certificate
	for _, c := range cm.certificates {
		if c.Domain == domain && c.IsValid {
			cert = c
			break
		}
	}

	if cert == nil {
		return nil, fmt.Errorf("no valid certificate found for domain %s", domain)
	}

	tlsCert, err := tls.X509KeyPair([]byte(cert.CertPEM), []byte(cert.KeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to create X509 key pair: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ServerName:   domain,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// renewalWatcher watches for certificates that need renewal
func (cm *CertificateManager) renewalWatcher() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		cm.checkRenewals()
	}
}

// checkRenewals checks and initiates certificate renewals
func (cm *CertificateManager) checkRenewals() {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for name, cert := range cm.certificates {
		if time.Until(cert.ExpiresAt) < cm.renewBefore {
			// In production, this would integrate with Let's Encrypt or other CA
			fmt.Printf("Certificate %s for domain %s needs renewal\n", name, cert.Domain)
			// cm.renewCertificate(cert)
		}
	}
}

// ======================== RATE LIMITING AND DDOS PROTECTION ========================

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*SimpleLimiter
	config   RateLimitConfig
}

// SimpleLimiter provides simple rate limiting
type SimpleLimiter struct {
	tokens     int
	lastRefill time.Time
	rate       int
	capacity   int
	mu         sync.Mutex
}

// RateLimitConfig defines rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int           `json:"requests_per_second"`
	BurstSize         int           `json:"burst_size"`
	WindowSize        time.Duration `json:"window_size"`
	BlockDuration     time.Duration `json:"block_duration"`
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*SimpleLimiter),
		config:   config,
	}
}

// Allow checks if a request from the given identifier should be allowed
func (rl *RateLimiter) Allow(identifier string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[identifier]
	if !exists {
		limiter = &SimpleLimiter{
			tokens:     rl.config.BurstSize,
			lastRefill: time.Now(),
			rate:       rl.config.RequestsPerSecond,
			capacity:   rl.config.BurstSize,
		}
		rl.limiters[identifier] = limiter
	}

	return limiter.Allow()
}

// AllowN checks if N requests from the given identifier should be allowed
func (rl *RateLimiter) AllowN(identifier string, n int) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[identifier]
	if !exists {
		limiter = &SimpleLimiter{
			tokens:     rl.config.BurstSize,
			lastRefill: time.Now(),
			rate:       rl.config.RequestsPerSecond,
			capacity:   rl.config.BurstSize,
		}
		rl.limiters[identifier] = limiter
	}

	return limiter.AllowN(n)
}

// Allow checks if a single request is allowed
func (sl *SimpleLimiter) Allow() bool {
	return sl.AllowN(1)
}

// AllowN checks if N requests are allowed
func (sl *SimpleLimiter) AllowN(n int) bool {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(sl.lastRefill)

	// Refill tokens based on elapsed time
	tokensToAdd := int(elapsed.Seconds()) * sl.rate
	sl.tokens += tokensToAdd
	if sl.tokens > sl.capacity {
		sl.tokens = sl.capacity
	}
	sl.lastRefill = now

	// Check if we have enough tokens
	if sl.tokens >= n {
		sl.tokens -= n
		return true
	}

	return false
}

// DDoSProtection provides DDoS protection mechanisms
type DDoSProtection struct {
	mu             sync.RWMutex
	rateLimiter    *RateLimiter
	blockedIPs     map[string]time.Time
	suspiciousIPs  map[string]int
	alertThreshold int
	blockDuration  time.Duration
	analysisWindow time.Duration
}

// NewDDoSProtection creates a new DDoS protection system
func NewDDoSProtection(rateLimiter *RateLimiter) *DDoSProtection {
	ddos := &DDoSProtection{
		rateLimiter:    rateLimiter,
		blockedIPs:     make(map[string]time.Time),
		suspiciousIPs:  make(map[string]int),
		alertThreshold: 100, // requests per minute
		blockDuration:  1 * time.Hour,
		analysisWindow: 1 * time.Minute,
	}

	go ddos.cleanupWorker()
	return ddos
}

// IsBlocked checks if an IP address is currently blocked
func (ddos *DDoSProtection) IsBlocked(ipAddress string) bool {
	ddos.mu.RLock()
	defer ddos.mu.RUnlock()

	if blockTime, exists := ddos.blockedIPs[ipAddress]; exists {
		if time.Since(blockTime) < ddos.blockDuration {
			return true
		}
		// Remove expired blocks
		delete(ddos.blockedIPs, ipAddress)
	}

	return false
}

// RecordRequest records a request and analyzes for DDoS patterns
func (ddos *DDoSProtection) RecordRequest(ipAddress string) {
	ddos.mu.Lock()
	defer ddos.mu.Unlock()

	ddos.suspiciousIPs[ipAddress]++

	// Check if IP should be blocked
	if ddos.suspiciousIPs[ipAddress] > ddos.alertThreshold {
		ddos.blockedIPs[ipAddress] = time.Now()
		fmt.Printf("IP %s blocked due to suspicious activity\n", ipAddress)
	}
}

// cleanupWorker periodically cleans up tracking data
func (ddos *DDoSProtection) cleanupWorker() {
	ticker := time.NewTicker(ddos.analysisWindow)
	defer ticker.Stop()

	for range ticker.C {
		ddos.mu.Lock()
		// Reset suspicious IP counters
		ddos.suspiciousIPs = make(map[string]int)

		// Clean expired blocked IPs
		for ip, blockTime := range ddos.blockedIPs {
			if time.Since(blockTime) >= ddos.blockDuration {
				delete(ddos.blockedIPs, ip)
			}
		}
		ddos.mu.Unlock()
	}
}

// ======================== INPUT VALIDATION AND SANITIZATION ========================

// Validator provides input validation and sanitization
type Validator struct {
	patterns map[string]*regexp.Regexp
}

// NewValidator creates a new input validator
func NewValidator() *Validator {
	v := &Validator{
		patterns: make(map[string]*regexp.Regexp),
	}

	// Initialize common validation patterns
	v.patterns["email"] = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	v.patterns["username"] = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	v.patterns["password"] = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`)
	v.patterns["ipv4"] = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	v.patterns["uuid"] = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

	return v
}

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(email string) error {
	if !v.patterns["email"].MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidateUsername validates username format
func (v *Validator) ValidateUsername(username string) error {
	if !v.patterns["username"].MatchString(username) {
		return fmt.Errorf("username must be 3-20 characters, alphanumeric and underscores only")
	}
	return nil
}

// ValidatePassword validates password strength
func (v *Validator) ValidatePassword(password string) error {
	if !v.patterns["password"].MatchString(password) {
		return fmt.Errorf("password must be at least 8 characters with uppercase, lowercase, digit, and special character")
	}
	return nil
}

// SanitizeHTML removes potentially dangerous HTML tags and attributes
func (v *Validator) SanitizeHTML(input string) string {
	// Basic HTML sanitization - in production, use a proper library like bluemonday
	dangerous := []string{"<script", "</script", "<iframe", "</iframe", "javascript:", "on"}

	sanitized := input
	for _, tag := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, tag, "")
	}

	return sanitized
}

// ValidateSQL checks for SQL injection patterns
func (v *Validator) ValidateSQL(input string) error {
	sqlPatterns := []string{
		`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`,
		`(?i)(or|and)\s+\d+\s*=\s*\d+`,
		`(?i)(or|and)\s+['"].*['"]`,
		`(?i)(\-\-)|(\#)|(/\*)|(\*/)`,
	}

	for _, pattern := range sqlPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			return fmt.Errorf("potential SQL injection detected")
		}
	}

	return nil
}

// ======================== AUDIT LOGGING AND COMPLIANCE ========================

// AuditLogger provides comprehensive audit logging
type AuditLogger struct {
	mu     sync.Mutex
	events []AuditEvent
	config AuditConfig
}

// AuditEvent represents an audit log event
type AuditEvent struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	UserID     string                 `json:"user_id,omitempty"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id,omitempty"`
	Result     AuditResult            `json:"result"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Risk       RiskLevel              `json:"risk_level"`
}

// AuditResult represents the result of an audited action
type AuditResult string

const (
	AuditResultSuccess AuditResult = "success"
	AuditResultFailure AuditResult = "failure"
	AuditResultBlocked AuditResult = "blocked"
)

// RiskLevel represents the risk level of an action
type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

// AuditConfig defines audit logging configuration
type AuditConfig struct {
	EnabledActions []string      `json:"enabled_actions"`
	RetentionDays  int           `json:"retention_days"`
	AlertOnHigh    bool          `json:"alert_on_high"`
	ExportInterval time.Duration `json:"export_interval"`
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(config AuditConfig) *AuditLogger {
	al := &AuditLogger{
		events: make([]AuditEvent, 0),
		config: config,
	}

	go al.retentionWorker()
	return al
}

// LogEvent logs an audit event
func (al *AuditLogger) LogEvent(event AuditEvent) {
	al.mu.Lock()
	defer al.mu.Unlock()

	event.ID = generateSecureID()
	event.Timestamp = time.Now()

	// Determine risk level if not set
	if event.Risk == "" {
		event.Risk = al.determineRiskLevel(event)
	}

	al.events = append(al.events, event)

	// Alert on high-risk events
	if al.config.AlertOnHigh && event.Risk == RiskLevelHigh {
		al.sendAlert(event)
	}
}

// LogAuthentication logs authentication events
func (al *AuditLogger) LogAuthentication(userID, ipAddress string, success bool) {
	result := AuditResultSuccess
	if !success {
		result = AuditResultFailure
	}

	event := AuditEvent{
		UserID:    userID,
		Action:    "authentication",
		Resource:  "user_session",
		Result:    result,
		IPAddress: ipAddress,
		Risk:      RiskLevelMedium,
	}

	al.LogEvent(event)
}

// LogDataAccess logs data access events
func (al *AuditLogger) LogDataAccess(userID, resource, resourceID, action string, success bool) {
	result := AuditResultSuccess
	if !success {
		result = AuditResultFailure
	}

	event := AuditEvent{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Result:     result,
		Risk:       al.determineDataAccessRisk(action, resource),
	}

	al.LogEvent(event)
}

// determineRiskLevel determines the risk level of an event
func (al *AuditLogger) determineRiskLevel(event AuditEvent) RiskLevel {
	// High-risk actions
	highRiskActions := []string{"delete", "admin_access", "privilege_escalation", "config_change"}
	for _, action := range highRiskActions {
		if event.Action == action {
			return RiskLevelHigh
		}
	}

	// Medium-risk actions
	mediumRiskActions := []string{"authentication", "authorization", "data_export"}
	for _, action := range mediumRiskActions {
		if event.Action == action {
			return RiskLevelMedium
		}
	}

	return RiskLevelLow
}

// determineDataAccessRisk determines risk level for data access
func (al *AuditLogger) determineDataAccessRisk(action, resource string) RiskLevel {
	if action == "delete" || strings.Contains(resource, "sensitive") {
		return RiskLevelHigh
	}
	if action == "write" || action == "update" {
		return RiskLevelMedium
	}
	return RiskLevelLow
}

// sendAlert sends alerts for high-risk events
func (al *AuditLogger) sendAlert(event AuditEvent) {
	// In production, this would integrate with alerting systems
	fmt.Printf("HIGH-RISK AUDIT EVENT: %s by user %s at %s\n",
		event.Action, event.UserID, event.Timestamp.Format(time.RFC3339))
}

// retentionWorker manages audit log retention
func (al *AuditLogger) retentionWorker() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		al.cleanupOldEvents()
	}
}

// cleanupOldEvents removes events older than retention period
func (al *AuditLogger) cleanupOldEvents() {
	al.mu.Lock()
	defer al.mu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -al.config.RetentionDays)
	var newEvents []AuditEvent

	for _, event := range al.events {
		if event.Timestamp.After(cutoff) {
			newEvents = append(newEvents, event)
		}
	}

	al.events = newEvents
}

// GetEvents returns audit events within a time range
func (al *AuditLogger) GetEvents(start, end time.Time) []AuditEvent {
	al.mu.Lock()
	defer al.mu.Unlock()

	var result []AuditEvent
	for _, event := range al.events {
		if event.Timestamp.After(start) && event.Timestamp.Before(end) {
			result = append(result, event)
		}
	}

	return result
}

// ======================== SECURE CODING PRACTICES ========================

// SecureStorage provides secure data storage patterns
type SecureStorage struct {
	encryptionKey []byte
}

// NewSecureStorage creates a new secure storage manager
func NewSecureStorage(key []byte) *SecureStorage {
	return &SecureStorage{
		encryptionKey: key,
	}
}

// StoreSensitiveData securely stores sensitive data
func (ss *SecureStorage) StoreSensitiveData(data []byte) ([]byte, error) {
	// Add timestamp and integrity check
	wrapper := struct {
		Data      []byte    `json:"data"`
		Timestamp time.Time `json:"timestamp"`
		Hash      string    `json:"hash"`
	}{
		Data:      data,
		Timestamp: time.Now(),
		Hash:      ss.calculateHash(data),
	}

	wrapperBytes, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}

	// Encrypt the wrapped data
	return ss.encrypt(wrapperBytes)
}

// RetrieveSensitiveData securely retrieves and validates sensitive data
func (ss *SecureStorage) RetrieveSensitiveData(encryptedData []byte) ([]byte, error) {
	// Decrypt the data
	decrypted, err := ss.decrypt(encryptedData)
	if err != nil {
		return nil, err
	}

	// Unmarshal wrapper
	var wrapper struct {
		Data      []byte    `json:"data"`
		Timestamp time.Time `json:"timestamp"`
		Hash      string    `json:"hash"`
	}

	if err := json.Unmarshal(decrypted, &wrapper); err != nil {
		return nil, err
	}

	// Verify integrity
	if wrapper.Hash != ss.calculateHash(wrapper.Data) {
		return nil, fmt.Errorf("data integrity check failed")
	}

	// Check if data is too old (optional)
	if time.Since(wrapper.Timestamp) > 24*time.Hour {
		return nil, fmt.Errorf("data expired")
	}

	return wrapper.Data, nil
}

func (ss *SecureStorage) encrypt(data []byte) ([]byte, error) {
	// Simple encryption implementation - in production, use proper AES-GCM
	hash := sha256.Sum256(append(data, ss.encryptionKey...))
	return append(data, hash[:]...), nil
}

func (ss *SecureStorage) decrypt(data []byte) ([]byte, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("invalid encrypted data")
	}

	payload := data[:len(data)-32]
	providedHash := data[len(data)-32:]

	expectedHash := sha256.Sum256(append(payload, ss.encryptionKey...))

	// Constant time comparison
	var equal byte = 0
	for i := 0; i < 32; i++ {
		equal |= providedHash[i] ^ expectedHash[i]
	}

	if equal != 0 {
		return nil, fmt.Errorf("decryption failed")
	}

	return payload, nil
}

func (ss *SecureStorage) calculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// ======================== UTILITY FUNCTIONS ========================

// generateSecureID generates a cryptographically secure ID
func generateSecureID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// ======================== HTTP SECURITY MIDDLEWARE ========================

// SecurityMiddleware provides HTTP security middleware
type SecurityMiddleware struct {
	authProvider   *AuthProvider
	rateLimiter    *RateLimiter
	ddosProtection *DDoSProtection
	validator      *Validator
	auditLogger    *AuditLogger
}

// NewSecurityMiddleware creates security middleware
func NewSecurityMiddleware() *SecurityMiddleware {
	rateLimitConfig := RateLimitConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		WindowSize:        time.Minute,
		BlockDuration:     time.Hour,
	}

	rateLimiter := NewRateLimiter(rateLimitConfig)
	auditConfig := AuditConfig{
		EnabledActions: []string{"authentication", "data_access", "config_change"},
		RetentionDays:  90,
		AlertOnHigh:    true,
		ExportInterval: 24 * time.Hour,
	}

	return &SecurityMiddleware{
		authProvider:   NewAuthProvider("secret-key", 24*time.Hour),
		rateLimiter:    rateLimiter,
		ddosProtection: NewDDoSProtection(rateLimiter),
		validator:      NewValidator(),
		auditLogger:    NewAuditLogger(auditConfig),
	}
}

// AuthenticateMiddleware validates JWT tokens
func (sm *SecurityMiddleware) AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for health checks
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		// Extract JWT token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Validate JWT
		claims, err := sm.authProvider.ValidateJWT(tokenParts[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			sm.auditLogger.LogAuthentication(claims.UserID, r.RemoteAddr, false)
			return
		}

		// Add user context to request
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "roles", claims.Roles)

		sm.auditLogger.LogAuthentication(claims.UserID, r.RemoteAddr, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RateLimitMiddleware applies rate limiting
func (sm *SecurityMiddleware) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		// Check if IP is blocked by DDoS protection
		if sm.ddosProtection.IsBlocked(clientIP) {
			http.Error(w, "IP temporarily blocked", http.StatusTooManyRequests)
			return
		}

		// Apply rate limiting
		if !sm.rateLimiter.Allow(clientIP) {
			sm.ddosProtection.RecordRequest(clientIP)
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware adds security headers
func (sm *SecurityMiddleware) SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		next.ServeHTTP(w, r)
	})
}

// InputValidationMiddleware validates request inputs
func (sm *SecurityMiddleware) InputValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate common injection patterns
		for _, value := range r.URL.Query() {
			for _, v := range value {
				if err := sm.validator.ValidateSQL(v); err != nil {
					http.Error(w, "Invalid input detected", http.StatusBadRequest)
					sm.auditLogger.LogEvent(AuditEvent{
						Action:    "input_validation_failure",
						Resource:  "request",
						Result:    AuditResultBlocked,
						IPAddress: r.RemoteAddr,
						Risk:      RiskLevelHigh,
						Metadata:  map[string]interface{}{"reason": err.Error()},
					})
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
