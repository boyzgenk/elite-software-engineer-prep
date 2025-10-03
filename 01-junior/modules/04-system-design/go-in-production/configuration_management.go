package production

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ======================== ENVIRONMENT-BASED CONFIGURATION ========================

// ConfigEnvironment represents different deployment environments
type ConfigEnvironment string

const (
	EnvDevelopment ConfigEnvironment = "development"
	EnvStaging     ConfigEnvironment = "staging"
	EnvProduction  ConfigEnvironment = "production"
	EnvTest        ConfigEnvironment = "test"
)

// ConfigManager manages application configuration across environments
type ConfigManager struct {
	mu            sync.RWMutex
	environment   Environment
	configs       map[string]*ConfigValue
	watchers      []ConfigWatcher
	secretManager *SecretManager
	validator     *ConfigValidator
	hotReloadCh   chan ConfigChange
	stopCh        chan struct{}
}

// ConfigValue represents a configuration value with metadata
type ConfigValue struct {
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Type        ConfigType  `json:"type"`
	Environment Environment `json:"environment"`
	Secret      bool        `json:"secret"`
	Required    bool        `json:"required"`
	DefaultVal  interface{} `json:"default"`
	Description string      `json:"description"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Version     int         `json:"version"`
}

// ConfigType represents the type of configuration value
type ConfigType string

const (
	ConfigTypeString ConfigType = "string"
	ConfigTypeInt    ConfigType = "int"
	ConfigTypeFloat  ConfigType = "float"
	ConfigTypeBool   ConfigType = "bool"
	ConfigTypeJSON   ConfigType = "json"
	ConfigTypeArray  ConfigType = "array"
)

// ConfigChange represents a configuration change event
type ConfigChange struct {
	Key      string      `json:"key"`
	OldValue interface{} `json:"old_value"`
	NewValue interface{} `json:"new_value"`
	Type     ChangeType  `json:"type"`
}

// ChangeType represents the type of configuration change
type ChangeType int

const (
	ChangeTypeAdd ChangeType = iota
	ChangeTypeUpdate
	ChangeTypeDelete
)

// ConfigWatcher is called when configuration changes
type ConfigWatcher interface {
	OnConfigChange(change ConfigChange) error
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(env Environment) *ConfigManager {
	cm := &ConfigManager{
		environment:   env,
		configs:       make(map[string]*ConfigValue),
		watchers:      make([]ConfigWatcher, 0),
		secretManager: NewSecretManager(),
		validator:     NewConfigValidator(),
		hotReloadCh:   make(chan ConfigChange, 100),
		stopCh:        make(chan struct{}),
	}

	// Load initial configuration
	cm.loadFromEnvironment()
	cm.loadFromFiles()

	// Start hot reload processor
	go cm.processHotReload()

	return cm
}

// loadFromEnvironment loads configuration from environment variables
func (cm *ConfigManager) loadFromEnvironment() {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// Skip system environment variables
		if isSystemEnvVar(key) {
			continue
		}

		configType := inferConfigType(value)
		parsedValue := parseConfigValue(value, configType)

		cm.configs[key] = &ConfigValue{
			Key:         key,
			Value:       parsedValue,
			Type:        configType,
			Environment: cm.environment,
			Secret:      isSecretKey(key),
			Required:    false,
			UpdatedAt:   time.Now(),
			Version:     1,
		}
	}
}

// loadFromFiles loads configuration from JSON/YAML files
func (cm *ConfigManager) loadFromFiles() {
	configFiles := []string{
		fmt.Sprintf("config/%s.json", cm.environment),
		"config/default.json",
		"config/app.json",
	}

	for _, filename := range configFiles {
		if data, err := os.ReadFile(filename); err == nil {
			var config map[string]interface{}
			if err := json.Unmarshal(data, &config); err == nil {
				cm.loadConfigFromMap(config, "")
			}
		}
	}
}

// loadConfigFromMap loads configuration from a map with optional prefix
func (cm *ConfigManager) loadConfigFromMap(config map[string]interface{}, prefix string) {
	for key, value := range config {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		// Handle nested objects
		if nested, ok := value.(map[string]interface{}); ok {
			cm.loadConfigFromMap(nested, fullKey)
			continue
		}

		configType := inferConfigTypeFromValue(value)

		cm.configs[fullKey] = &ConfigValue{
			Key:         fullKey,
			Value:       value,
			Type:        configType,
			Environment: cm.environment,
			Secret:      isSecretKey(fullKey),
			Required:    false,
			UpdatedAt:   time.Now(),
			Version:     1,
		}
	}
}

// Set sets a configuration value
func (cm *ConfigManager) Set(key string, value interface{}) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	configType := inferConfigTypeFromValue(value)

	// Validate the new value
	if err := cm.validator.Validate(key, value, configType); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	oldConfig := cm.configs[key]
	var oldValue interface{}
	if oldConfig != nil {
		oldValue = oldConfig.Value
	}

	// Handle secrets
	if isSecretKey(key) {
		encryptedValue, err := cm.secretManager.Encrypt(fmt.Sprintf("%v", value))
		if err != nil {
			return fmt.Errorf("failed to encrypt secret: %w", err)
		}
		value = encryptedValue
	}

	newConfig := &ConfigValue{
		Key:         key,
		Value:       value,
		Type:        configType,
		Environment: cm.environment,
		Secret:      isSecretKey(key),
		Required:    false,
		UpdatedAt:   time.Now(),
		Version:     1,
	}

	if oldConfig != nil {
		newConfig.Version = oldConfig.Version + 1
		newConfig.Required = oldConfig.Required
		newConfig.DefaultVal = oldConfig.DefaultVal
		newConfig.Description = oldConfig.Description
	}

	cm.configs[key] = newConfig

	// Notify hot reload
	change := ConfigChange{
		Key:      key,
		OldValue: oldValue,
		NewValue: value,
		Type:     ChangeTypeUpdate,
	}

	if oldConfig == nil {
		change.Type = ChangeTypeAdd
	}

	select {
	case cm.hotReloadCh <- change:
	default:
		// Channel full, skip notification
	}

	return nil
}

// Get retrieves a configuration value
func (cm *ConfigManager) Get(key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	config, exists := cm.configs[key]
	if !exists {
		return nil, false
	}

	value := config.Value

	// Decrypt secrets
	if config.Secret {
		decrypted, err := cm.secretManager.Decrypt(fmt.Sprintf("%v", value))
		if err != nil {
			return nil, false
		}
		value = decrypted
	}

	return value, true
}

// GetString retrieves a string configuration value
func (cm *ConfigManager) GetString(key string) string {
	if value, exists := cm.Get(key); exists {
		if str, ok := value.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", value)
	}
	return ""
}

// GetInt retrieves an integer configuration value
func (cm *ConfigManager) GetInt(key string) int {
	if value, exists := cm.Get(key); exists {
		if i, ok := value.(int); ok {
			return i
		}
		if f, ok := value.(float64); ok {
			return int(f)
		}
		if str, ok := value.(string); ok {
			if i, err := strconv.Atoi(str); err == nil {
				return i
			}
		}
	}
	return 0
}

// GetBool retrieves a boolean configuration value
func (cm *ConfigManager) GetBool(key string) bool {
	if value, exists := cm.Get(key); exists {
		if b, ok := value.(bool); ok {
			return b
		}
		if str, ok := value.(string); ok {
			return strings.ToLower(str) == "true" || str == "1"
		}
	}
	return false
}

// GetFloat retrieves a float configuration value
func (cm *ConfigManager) GetFloat(key string) float64 {
	if value, exists := cm.Get(key); exists {
		if f, ok := value.(float64); ok {
			return f
		}
		if i, ok := value.(int); ok {
			return float64(i)
		}
		if str, ok := value.(string); ok {
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				return f
			}
		}
	}
	return 0.0
}

// GetStringSlice retrieves a string slice configuration value
func (cm *ConfigManager) GetStringSlice(key string) []string {
	if value, exists := cm.Get(key); exists {
		if slice, ok := value.([]string); ok {
			return slice
		}
		if arr, ok := value.([]interface{}); ok {
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = fmt.Sprintf("%v", v)
			}
			return result
		}
		if str, ok := value.(string); ok {
			return strings.Split(str, ",")
		}
	}
	return []string{}
}

// Delete removes a configuration value
func (cm *ConfigManager) Delete(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if oldConfig, exists := cm.configs[key]; exists {
		delete(cm.configs, key)

		// Notify hot reload
		change := ConfigChange{
			Key:      key,
			OldValue: oldConfig.Value,
			NewValue: nil,
			Type:     ChangeTypeDelete,
		}

		select {
		case cm.hotReloadCh <- change:
		default:
		}
	}
}

// AddWatcher adds a configuration change watcher
func (cm *ConfigManager) AddWatcher(watcher ConfigWatcher) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.watchers = append(cm.watchers, watcher)
}

// processHotReload processes hot reload events
func (cm *ConfigManager) processHotReload() {
	for {
		select {
		case change := <-cm.hotReloadCh:
			cm.mu.RLock()
			watchers := make([]ConfigWatcher, len(cm.watchers))
			copy(watchers, cm.watchers)
			cm.mu.RUnlock()

			for _, watcher := range watchers {
				if err := watcher.OnConfigChange(change); err != nil {
					// Log error but continue
					fmt.Printf("Config watcher error: %v\n", err)
				}
			}
		case <-cm.stopCh:
			return
		}
	}
}

// Stop stops the configuration manager
func (cm *ConfigManager) Stop() {
	close(cm.stopCh)
}

// Export exports all configuration (excluding secrets)
func (cm *ConfigManager) Export() map[string]*ConfigValue {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	result := make(map[string]*ConfigValue)
	for key, config := range cm.configs {
		if !config.Secret {
			configCopy := *config
			result[key] = &configCopy
		} else {
			// Export metadata only for secrets
			result[key] = &ConfigValue{
				Key:         config.Key,
				Type:        config.Type,
				Environment: config.Environment,
				Secret:      true,
				Required:    config.Required,
				Description: config.Description,
				UpdatedAt:   config.UpdatedAt,
				Version:     config.Version,
			}
		}
	}
	return result
}

// ======================== SECRET MANAGEMENT AND ENCRYPTION ========================

// SecretManager handles encryption and decryption of sensitive configuration
type SecretManager struct {
	mu         sync.RWMutex
	key        []byte
	gcm        cipher.AEAD
	keyRotator *KeyRotator
}

// KeyRotator manages encryption key rotation
type KeyRotator struct {
	currentKey       []byte
	oldKeys          [][]byte
	rotationInterval time.Duration
	lastRotation     time.Time
}

// NewSecretManager creates a new secret manager
func NewSecretManager() *SecretManager {
	key := generateOrLoadKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("failed to create cipher: %v", err))
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(fmt.Sprintf("failed to create GCM: %v", err))
	}

	sm := &SecretManager{
		key: key,
		gcm: gcm,
		keyRotator: &KeyRotator{
			currentKey:       key,
			oldKeys:          make([][]byte, 0),
			rotationInterval: 24 * time.Hour,
			lastRotation:     time.Now(),
		},
	}

	go sm.rotateKeys()
	return sm
}

// Encrypt encrypts a plaintext value
func (sm *SecretManager) Encrypt(plaintext string) (string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	nonce := make([]byte, sm.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := sm.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts an encrypted value
func (sm *SecretManager) Decrypt(encrypted string) (string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	nonceSize := sm.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Try current key first
	if plaintext, err := sm.gcm.Open(nil, nonce, ciphertext, nil); err == nil {
		return string(plaintext), nil
	}

	// Try old keys for backward compatibility
	for _, oldKey := range sm.keyRotator.oldKeys {
		if block, err := aes.NewCipher(oldKey); err == nil {
			if gcm, err := cipher.NewGCM(block); err == nil {
				if plaintext, err := gcm.Open(nil, nonce, ciphertext, nil); err == nil {
					return string(plaintext), nil
				}
			}
		}
	}

	return "", fmt.Errorf("failed to decrypt")
}

// rotateKeys periodically rotates encryption keys
func (sm *SecretManager) rotateKeys() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.Lock()

		if time.Since(sm.keyRotator.lastRotation) >= sm.keyRotator.rotationInterval {
			// Generate new key
			newKey := generateKey()

			// Store old key
			sm.keyRotator.oldKeys = append(sm.keyRotator.oldKeys, sm.keyRotator.currentKey)

			// Keep only last 5 keys
			if len(sm.keyRotator.oldKeys) > 5 {
				sm.keyRotator.oldKeys = sm.keyRotator.oldKeys[1:]
			}

			// Update current key
			sm.keyRotator.currentKey = newKey
			sm.key = newKey

			// Update cipher
			if block, err := aes.NewCipher(newKey); err == nil {
				if gcm, err := cipher.NewGCM(block); err == nil {
					sm.gcm = gcm
					sm.keyRotator.lastRotation = time.Now()
				}
			}
		}

		sm.mu.Unlock()
	}
}

// generateOrLoadKey generates a new key or loads existing one
func generateOrLoadKey() []byte {
	// In production, this would load from secure key management service
	// For demo, generate a key based on environment
	envKey := os.Getenv("CONFIG_ENCRYPTION_KEY")
	if envKey != "" {
		hash := sha256.Sum256([]byte(envKey))
		return hash[:]
	}

	return generateKey()
}

// generateKey generates a random 256-bit key
func generateKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)
	return key
}

// ======================== CONFIGURATION HOT-RELOADING ========================

// HotReloader monitors configuration files and triggers reloads
type HotReloader struct {
	configManager *ConfigManager
	watchedFiles  map[string]time.Time
	pollInterval  time.Duration
	stopCh        chan struct{}
}

// NewHotReloader creates a new hot reloader
func NewHotReloader(cm *ConfigManager, pollInterval time.Duration) *HotReloader {
	hr := &HotReloader{
		configManager: cm,
		watchedFiles:  make(map[string]time.Time),
		pollInterval:  pollInterval,
		stopCh:        make(chan struct{}),
	}

	// Add default files to watch
	files := []string{
		fmt.Sprintf("config/%s.json", cm.environment),
		"config/default.json",
		"config/app.json",
	}

	for _, file := range files {
		if info, err := os.Stat(file); err == nil {
			hr.watchedFiles[file] = info.ModTime()
		}
	}

	go hr.watch()
	return hr
}

// watch monitors files for changes
func (hr *HotReloader) watch() {
	ticker := time.NewTicker(hr.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hr.checkFiles()
		case <-hr.stopCh:
			return
		}
	}
}

// checkFiles checks for file modifications
func (hr *HotReloader) checkFiles() {
	for filename, lastModTime := range hr.watchedFiles {
		if info, err := os.Stat(filename); err == nil {
			if info.ModTime().After(lastModTime) {
				hr.watchedFiles[filename] = info.ModTime()
				hr.reloadFile(filename)
			}
		}
	}
}

// reloadFile reloads a specific configuration file
func (hr *HotReloader) reloadFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read config file %s: %v\n", filename, err)
		return
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Failed to parse config file %s: %v\n", filename, err)
		return
	}

	// Apply changes
	hr.applyConfigChanges(config, "")
}

// applyConfigChanges applies configuration changes from a map
func (hr *HotReloader) applyConfigChanges(config map[string]interface{}, prefix string) {
	for key, value := range config {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		// Handle nested objects
		if nested, ok := value.(map[string]interface{}); ok {
			hr.applyConfigChanges(nested, fullKey)
			continue
		}

		// Update configuration
		if err := hr.configManager.Set(fullKey, value); err != nil {
			fmt.Printf("Failed to update config %s: %v\n", fullKey, err)
		}
	}
}

// Stop stops the hot reloader
func (hr *HotReloader) Stop() {
	close(hr.stopCh)
}

// ======================== FEATURE FLAGS AND TOGGLES ========================

// FeatureToggle represents a feature toggle with advanced targeting
type FeatureToggle struct {
	Name      string                 `json:"name"`
	Enabled   bool                   `json:"enabled"`
	Rollout   *RolloutConfig         `json:"rollout,omitempty"`
	Targeting *TargetingConfig       `json:"targeting,omitempty"`
	Variants  map[string]*Variant    `json:"variants,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// RolloutConfig defines rollout percentage and strategy
type RolloutConfig struct {
	Percentage int             `json:"percentage"`
	Strategy   RolloutStrategy `json:"strategy"`
	Seed       string          `json:"seed"`
}

// RolloutStrategy defines how rollout is calculated
type RolloutStrategy string

const (
	RolloutStrategyRandom  RolloutStrategy = "random"
	RolloutStrategyUserID  RolloutStrategy = "user_id"
	RolloutStrategySession RolloutStrategy = "session"
)

// TargetingConfig defines targeting rules
type TargetingConfig struct {
	Rules []TargetingRule `json:"rules"`
}

// TargetingRule defines a single targeting rule
type TargetingRule struct {
	Attribute string   `json:"attribute"`
	Operator  Operator `json:"operator"`
	Values    []string `json:"values"`
}

// Operator defines comparison operators for targeting
type Operator string

const (
	OperatorEquals    Operator = "equals"
	OperatorNotEquals Operator = "not_equals"
	OperatorIn        Operator = "in"
	OperatorNotIn     Operator = "not_in"
	OperatorContains  Operator = "contains"
	OperatorMatches   Operator = "matches"
)

// Variant represents a feature variant for A/B testing
type Variant struct {
	Name   string                 `json:"name"`
	Weight int                    `json:"weight"`
	Config map[string]interface{} `json:"config"`
}

// EvaluationContext provides context for feature evaluation
type EvaluationContext struct {
	UserID     string            `json:"user_id"`
	SessionID  string            `json:"session_id"`
	Attributes map[string]string `json:"attributes"`
}

// FeatureToggleManager manages feature toggles and evaluation
type FeatureToggleManager struct {
	mu      sync.RWMutex
	toggles map[string]*FeatureToggle
	hasher  func(string) uint32
}

// NewFeatureToggleManager creates a new feature toggle manager
func NewFeatureToggleManager() *FeatureToggleManager {
	return &FeatureToggleManager{
		toggles: make(map[string]*FeatureToggle),
		hasher:  hashString32,
	}
}

// CreateToggle creates a new feature toggle
func (ftm *FeatureToggleManager) CreateToggle(toggle *FeatureToggle) {
	ftm.mu.Lock()
	defer ftm.mu.Unlock()

	toggle.CreatedAt = time.Now()
	toggle.UpdatedAt = time.Now()
	ftm.toggles[toggle.Name] = toggle
}

// UpdateToggle updates an existing feature toggle
func (ftm *FeatureToggleManager) UpdateToggle(name string, updates map[string]interface{}) error {
	ftm.mu.Lock()
	defer ftm.mu.Unlock()

	toggle, exists := ftm.toggles[name]
	if !exists {
		return fmt.Errorf("toggle %s not found", name)
	}

	// Apply updates using reflection for simplicity
	for field, value := range updates {
		if err := ftm.updateToggleField(toggle, field, value); err != nil {
			return fmt.Errorf("failed to update field %s: %w", field, err)
		}
	}

	toggle.UpdatedAt = time.Now()
	return nil
}

// updateToggleField updates a specific field of a toggle
func (ftm *FeatureToggleManager) updateToggleField(toggle *FeatureToggle, field string, value interface{}) error {
	switch field {
	case "enabled":
		if b, ok := value.(bool); ok {
			toggle.Enabled = b
		}
	case "rollout.percentage":
		if i, ok := value.(int); ok {
			if toggle.Rollout == nil {
				toggle.Rollout = &RolloutConfig{}
			}
			toggle.Rollout.Percentage = i
		}
	default:
		return fmt.Errorf("unknown field: %s", field)
	}
	return nil
}

// IsEnabled evaluates if a feature is enabled for a given context
func (ftm *FeatureToggleManager) IsEnabled(featureName string, ctx EvaluationContext) bool {
	ftm.mu.RLock()
	toggle, exists := ftm.toggles[featureName]
	ftm.mu.RUnlock()

	if !exists || !toggle.Enabled {
		return false
	}

	// Check targeting rules
	if toggle.Targeting != nil {
		if !ftm.evaluateTargeting(toggle.Targeting, ctx) {
			return false
		}
	}

	// Check rollout percentage
	if toggle.Rollout != nil {
		return ftm.evaluateRollout(toggle.Rollout, ctx, featureName)
	}

	return true
}

// GetVariant returns the variant for a feature (for A/B testing)
func (ftm *FeatureToggleManager) GetVariant(featureName string, ctx EvaluationContext) *Variant {
	ftm.mu.RLock()
	toggle, exists := ftm.toggles[featureName]
	ftm.mu.RUnlock()

	if !exists || !toggle.Enabled || len(toggle.Variants) == 0 {
		return nil
	}

	// Check if feature is enabled first
	if !ftm.IsEnabled(featureName, ctx) {
		return nil
	}

	// Calculate variant based on weights
	totalWeight := 0
	for _, variant := range toggle.Variants {
		totalWeight += variant.Weight
	}

	if totalWeight == 0 {
		return nil
	}

	// Generate hash based on user ID and feature name
	hash := ftm.hasher(ctx.UserID + featureName)
	bucket := int(hash % uint32(totalWeight))

	currentWeight := 0
	for _, variant := range toggle.Variants {
		currentWeight += variant.Weight
		if bucket < currentWeight {
			return variant
		}
	}

	return nil
}

// evaluateTargeting evaluates targeting rules
func (ftm *FeatureToggleManager) evaluateTargeting(targeting *TargetingConfig, ctx EvaluationContext) bool {
	for _, rule := range targeting.Rules {
		if !ftm.evaluateRule(rule, ctx) {
			return false
		}
	}
	return true
}

// evaluateRule evaluates a single targeting rule
func (ftm *FeatureToggleManager) evaluateRule(rule TargetingRule, ctx EvaluationContext) bool {
	var value string

	switch rule.Attribute {
	case "user_id":
		value = ctx.UserID
	case "session_id":
		value = ctx.SessionID
	default:
		value = ctx.Attributes[rule.Attribute]
	}

	switch rule.Operator {
	case OperatorEquals:
		return len(rule.Values) > 0 && value == rule.Values[0]
	case OperatorNotEquals:
		return len(rule.Values) > 0 && value != rule.Values[0]
	case OperatorIn:
		for _, v := range rule.Values {
			if value == v {
				return true
			}
		}
		return false
	case OperatorNotIn:
		for _, v := range rule.Values {
			if value == v {
				return false
			}
		}
		return true
	case OperatorContains:
		return len(rule.Values) > 0 && strings.Contains(value, rule.Values[0])
	default:
		return false
	}
}

// evaluateRollout evaluates rollout percentage
func (ftm *FeatureToggleManager) evaluateRollout(rollout *RolloutConfig, ctx EvaluationContext, featureName string) bool {
	var key string

	switch rollout.Strategy {
	case RolloutStrategyUserID:
		key = ctx.UserID
	case RolloutStrategySession:
		key = ctx.SessionID
	default:
		key = ctx.UserID
	}

	if rollout.Seed != "" {
		key = rollout.Seed + key
	}

	key += featureName
	hash := ftm.hasher(key)
	bucket := int(hash % 100)

	return bucket < rollout.Percentage
}

// ======================== CONFIGURATION VALIDATION ========================

// ConfigValidator validates configuration values
type ConfigValidator struct {
	rules map[string]ValidationRule
}

// ValidationRule defines validation constraints for a config key
type ValidationRule struct {
	Required bool                    `json:"required"`
	Type     ConfigType              `json:"type"`
	Min      interface{}             `json:"min,omitempty"`
	Max      interface{}             `json:"max,omitempty"`
	Pattern  string                  `json:"pattern,omitempty"`
	Enum     []interface{}           `json:"enum,omitempty"`
	Custom   func(interface{}) error `json:"-"`
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	cv := &ConfigValidator{
		rules: make(map[string]ValidationRule),
	}

	// Add default validation rules
	cv.addDefaultRules()
	return cv
}

// addDefaultRules adds common validation rules
func (cv *ConfigValidator) addDefaultRules() {
	cv.rules["port"] = ValidationRule{
		Required: true,
		Type:     ConfigTypeInt,
		Min:      1,
		Max:      65535,
	}

	cv.rules["database.host"] = ValidationRule{
		Required: true,
		Type:     ConfigTypeString,
	}

	cv.rules["timeout"] = ValidationRule{
		Type:    ConfigTypeString,
		Pattern: `^\d+[smh]$`, // e.g., "30s", "5m", "1h"
	}
}

// AddRule adds a validation rule for a configuration key
func (cv *ConfigValidator) AddRule(key string, rule ValidationRule) {
	cv.rules[key] = rule
}

// Validate validates a configuration value against its rules
func (cv *ConfigValidator) Validate(key string, value interface{}, configType ConfigType) error {
	rule, hasRule := cv.rules[key]
	if !hasRule {
		return nil // No validation rule defined
	}

	// Check if required
	if rule.Required && value == nil {
		return fmt.Errorf("config key %s is required", key)
	}

	if value == nil {
		return nil // Optional value not provided
	}

	// Type validation
	if rule.Type != "" && rule.Type != configType {
		return fmt.Errorf("config key %s must be of type %s", key, rule.Type)
	}

	// Range validation
	if err := cv.validateRange(key, value, rule); err != nil {
		return err
	}

	// Enum validation
	if len(rule.Enum) > 0 {
		if err := cv.validateEnum(key, value, rule.Enum); err != nil {
			return err
		}
	}

	// Pattern validation
	if rule.Pattern != "" {
		if err := cv.validatePattern(key, value, rule.Pattern); err != nil {
			return err
		}
	}

	// Custom validation
	if rule.Custom != nil {
		if err := rule.Custom(value); err != nil {
			return fmt.Errorf("custom validation failed for %s: %w", key, err)
		}
	}

	return nil
}

// validateRange validates value is within specified range
func (cv *ConfigValidator) validateRange(key string, value interface{}, rule ValidationRule) error {
	if rule.Min == nil && rule.Max == nil {
		return nil
	}

	switch v := value.(type) {
	case int:
		if rule.Min != nil {
			if min, ok := rule.Min.(int); ok && v < min {
				return fmt.Errorf("config key %s value %d is less than minimum %d", key, v, min)
			}
		}
		if rule.Max != nil {
			if max, ok := rule.Max.(int); ok && v > max {
				return fmt.Errorf("config key %s value %d is greater than maximum %d", key, v, max)
			}
		}
	case float64:
		if rule.Min != nil {
			if min, ok := rule.Min.(float64); ok && v < min {
				return fmt.Errorf("config key %s value %f is less than minimum %f", key, v, min)
			}
		}
		if rule.Max != nil {
			if max, ok := rule.Max.(float64); ok && v > max {
				return fmt.Errorf("config key %s value %f is greater than maximum %f", key, v, max)
			}
		}
	case string:
		if rule.Min != nil {
			if min, ok := rule.Min.(int); ok && len(v) < min {
				return fmt.Errorf("config key %s string length %d is less than minimum %d", key, len(v), min)
			}
		}
		if rule.Max != nil {
			if max, ok := rule.Max.(int); ok && len(v) > max {
				return fmt.Errorf("config key %s string length %d is greater than maximum %d", key, len(v), max)
			}
		}
	}

	return nil
}

// validateEnum validates value is in allowed enum values
func (cv *ConfigValidator) validateEnum(key string, value interface{}, enum []interface{}) error {
	for _, allowed := range enum {
		if reflect.DeepEqual(value, allowed) {
			return nil
		}
	}
	return fmt.Errorf("config key %s value %v is not in allowed values %v", key, value, enum)
}

// validatePattern validates string value matches pattern
func (cv *ConfigValidator) validatePattern(key string, value interface{}, pattern string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("config key %s pattern validation requires string value", key)
	}

	// Simple pattern matching - in production, use regexp
	// For demo, just check if pattern is contained in string
	if !strings.Contains(str, strings.Trim(pattern, "^$")) {
		return fmt.Errorf("config key %s value %s does not match pattern %s", key, str, pattern)
	}

	return nil
}

// ======================== MULTI-ENVIRONMENT DEPLOYMENT ========================

// EnvironmentManager manages configuration across multiple environments
type EnvironmentManager struct {
	environments map[Environment]*ConfigManager
	currentEnv   Environment
	promoter     *ConfigPromoter
}

// ConfigPromoter handles configuration promotion between environments
type ConfigPromoter struct {
	approvalRequired  bool
	approvers         []string
	pendingPromotions map[string]*PromotionRequest
	mu                sync.RWMutex
}

// PromotionRequest represents a request to promote config between environments
type PromotionRequest struct {
	ID        string          `json:"id"`
	FromEnv   Environment     `json:"from_env"`
	ToEnv     Environment     `json:"to_env"`
	ConfigKey string          `json:"config_key"`
	Value     interface{}     `json:"value"`
	Requester string          `json:"requester"`
	Approvals []string        `json:"approvals"`
	Status    PromotionStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}

// PromotionStatus represents the status of a promotion request
type PromotionStatus string

const (
	PromotionStatusPending  PromotionStatus = "pending"
	PromotionStatusApproved PromotionStatus = "approved"
	PromotionStatusRejected PromotionStatus = "rejected"
	PromotionStatusApplied  PromotionStatus = "applied"
)

// NewEnvironmentManager creates a new environment manager
func NewEnvironmentManager() *EnvironmentManager {
	return &EnvironmentManager{
		environments: map[Environment]*ConfigManager{
			EnvDevelopment: NewConfigManager(EnvDevelopment),
			EnvStaging:     NewConfigManager(EnvStaging),
			EnvProduction:  NewConfigManager(EnvProduction),
		},
		currentEnv: EnvDevelopment,
		promoter: &ConfigPromoter{
			approvalRequired:  true,
			approvers:         []string{"admin", "lead"},
			pendingPromotions: make(map[string]*PromotionRequest),
		},
	}
}

// GetConfigManager returns the config manager for an environment
func (em *EnvironmentManager) GetConfigManager(env Environment) *ConfigManager {
	return em.environments[env]
}

// PromoteConfig promotes configuration from one environment to another
func (em *EnvironmentManager) PromoteConfig(fromEnv, toEnv Environment, configKey, requester string) (string, error) {
	fromCM := em.environments[fromEnv]
	if fromCM == nil {
		return "", fmt.Errorf("source environment %s not found", fromEnv)
	}

	toCM := em.environments[toEnv]
	if toCM == nil {
		return "", fmt.Errorf("target environment %s not found", toEnv)
	}

	// Get value from source environment
	value, exists := fromCM.Get(configKey)
	if !exists {
		return "", fmt.Errorf("config key %s not found in environment %s", configKey, fromEnv)
	}

	// Create promotion request
	requestID := fmt.Sprintf("promo_%d", time.Now().Unix())

	request := &PromotionRequest{
		ID:        requestID,
		FromEnv:   fromEnv,
		ToEnv:     toEnv,
		ConfigKey: configKey,
		Value:     value,
		Requester: requester,
		Approvals: make([]string, 0),
		Status:    PromotionStatusPending,
		CreatedAt: time.Now(),
	}

	em.promoter.mu.Lock()
	em.promoter.pendingPromotions[requestID] = request
	em.promoter.mu.Unlock()

	// Auto-approve if promotion doesn't require approval
	if !em.promoter.approvalRequired {
		return requestID, em.ApprovePromotion(requestID, "system")
	}

	return requestID, nil
}

// ApprovePromotion approves a promotion request
func (em *EnvironmentManager) ApprovePromotion(requestID, approver string) error {
	em.promoter.mu.Lock()
	defer em.promoter.mu.Unlock()

	request, exists := em.promoter.pendingPromotions[requestID]
	if !exists {
		return fmt.Errorf("promotion request %s not found", requestID)
	}

	if request.Status != PromotionStatusPending {
		return fmt.Errorf("promotion request %s is not pending", requestID)
	}

	// Check if approver is authorized
	authorized := false
	for _, validApprover := range em.promoter.approvers {
		if approver == validApprover {
			authorized = true
			break
		}
	}

	if !authorized && em.promoter.approvalRequired {
		return fmt.Errorf("approver %s is not authorized", approver)
	}

	// Add approval
	request.Approvals = append(request.Approvals, approver)

	// Check if we have enough approvals (for simplicity, require 1)
	if len(request.Approvals) >= 1 {
		request.Status = PromotionStatusApproved

		// Apply the promotion
		toCM := em.environments[request.ToEnv]
		if err := toCM.Set(request.ConfigKey, request.Value); err != nil {
			return fmt.Errorf("failed to apply promotion: %w", err)
		}

		request.Status = PromotionStatusApplied
	}

	return nil
}

// ======================== UTILITY FUNCTIONS ========================

// isSystemEnvVar checks if an environment variable is a system variable
func isSystemEnvVar(key string) bool {
	systemVars := []string{"PATH", "HOME", "USER", "SHELL", "TERM", "PWD"}
	for _, sysVar := range systemVars {
		if key == sysVar || strings.HasPrefix(key, sysVar+"_") {
			return true
		}
	}
	return false
}

// isSecretKey checks if a configuration key represents a secret
func isSecretKey(key string) bool {
	secretKeywords := []string{"password", "secret", "key", "token", "credential", "private"}
	lowerKey := strings.ToLower(key)

	for _, keyword := range secretKeywords {
		if strings.Contains(lowerKey, keyword) {
			return true
		}
	}
	return false
}

// inferConfigType infers the configuration type from a string value
func inferConfigType(value string) ConfigType {
	// Try to parse as different types
	if _, err := strconv.Atoi(value); err == nil {
		return ConfigTypeInt
	}

	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return ConfigTypeFloat
	}

	if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
		return ConfigTypeBool
	}

	if strings.HasPrefix(value, "[") || strings.HasPrefix(value, "{") {
		return ConfigTypeJSON
	}

	return ConfigTypeString
}

// inferConfigTypeFromValue infers the configuration type from a value
func inferConfigTypeFromValue(value interface{}) ConfigType {
	switch value.(type) {
	case int, int32, int64:
		return ConfigTypeInt
	case float32, float64:
		return ConfigTypeFloat
	case bool:
		return ConfigTypeBool
	case []interface{}, []string:
		return ConfigTypeArray
	case map[string]interface{}:
		return ConfigTypeJSON
	default:
		return ConfigTypeString
	}
}

// parseConfigValue parses a string value according to the specified type
func parseConfigValue(value string, configType ConfigType) interface{} {
	switch configType {
	case ConfigTypeInt:
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	case ConfigTypeFloat:
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	case ConfigTypeBool:
		return strings.ToLower(value) == "true" || value == "1"
	case ConfigTypeJSON:
		var result interface{}
		if err := json.Unmarshal([]byte(value), &result); err == nil {
			return result
		}
	case ConfigTypeArray:
		return strings.Split(value, ",")
	}

	return value
}

// hashString32 provides 32-bit hash for consistent bucketing
func hashString32(s string) uint32 {
	hash := uint32(2166136261)
	for _, c := range []byte(s) {
		hash = (hash ^ uint32(c)) * 16777619
	}
	return hash
}

// ======================== HTTP HANDLERS FOR CONFIGURATION MANAGEMENT ========================

// ConfigHandler provides HTTP endpoints for configuration management
type ConfigHandler struct {
	configManager *ConfigManager
	envManager    *EnvironmentManager
	toggleManager *FeatureToggleManager
}

// NewConfigHandler creates HTTP handlers for configuration management
func NewConfigHandler() *ConfigHandler {
	envManager := NewEnvironmentManager()

	return &ConfigHandler{
		configManager: envManager.GetConfigManager(EnvDevelopment),
		envManager:    envManager,
		toggleManager: NewFeatureToggleManager(),
	}
}

// HandleConfig handles configuration CRUD operations
func (ch *ConfigHandler) HandleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		configs := ch.configManager.Export()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(configs)

	case http.MethodPost:
		var req struct {
			Key   string      `json:"key"`
			Value interface{} `json:"value"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := ch.configManager.Set(req.Key, req.Value); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "configuration updated"})

	case http.MethodDelete:
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key parameter required", http.StatusBadRequest)
			return
		}

		ch.configManager.Delete(key)
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleFeatureToggle handles feature toggle operations
func (ch *ConfigHandler) HandleFeatureToggle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var toggle FeatureToggle
		if err := json.NewDecoder(r.Body).Decode(&toggle); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ch.toggleManager.CreateToggle(&toggle)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(toggle)

	case http.MethodGet:
		featureName := r.URL.Query().Get("feature")
		userID := r.URL.Query().Get("user_id")

		if featureName == "" {
			http.Error(w, "feature parameter required", http.StatusBadRequest)
			return
		}

		ctx := EvaluationContext{UserID: userID}
		enabled := ch.toggleManager.IsEnabled(featureName, ctx)
		variant := ch.toggleManager.GetVariant(featureName, ctx)

		response := map[string]interface{}{
			"enabled": enabled,
			"variant": variant,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandlePromotion handles configuration promotion between environments
func (ch *ConfigHandler) HandlePromotion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FromEnv   string `json:"from_env"`
		ToEnv     string `json:"to_env"`
		ConfigKey string `json:"config_key"`
		Requester string `json:"requester"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestID, err := ch.envManager.PromoteConfig(
		Environment(req.FromEnv),
		Environment(req.ToEnv),
		req.ConfigKey,
		req.Requester,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"request_id": requestID,
		"status":     "promotion requested",
	})
}
