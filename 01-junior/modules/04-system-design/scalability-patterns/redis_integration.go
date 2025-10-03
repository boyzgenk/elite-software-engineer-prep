package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 🔴 FAANG System Design - Redis Integration Patterns
// Essential caching patterns for L3/L4 interviews with Redis

// RedisClient interface defines Redis operations for caching
type RedisClient interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
	Del(key string) error
	Exists(key string) (bool, error)
	TTL(key string) (time.Duration, error)
	IncrBy(key string, increment int64) (int64, error)
	HSet(key, field, value string) error
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	SAdd(key string, members ...string) error
	SMembers(key string) ([]string, error)
	ZAdd(key string, score float64, member string) error
	ZRange(key string, start, stop int64) ([]string, error)
	ZRangeByScore(key string, min, max float64) ([]string, error)
	Ping() error
	FlushAll() error
}

// 💾 Mock Redis Client for Demo Purposes
// In production, use github.com/go-redis/redis or similar
type MockRedisClient struct {
	data       map[string]redisValue
	hashes     map[string]map[string]string   // Hash data structures
	sets       map[string]map[string]struct{} // Set data structures
	sortedSets map[string][]zsetMember        // Sorted set data structures
	mutex      sync.RWMutex
	stats      RedisStats
}

type redisValue struct {
	value     string
	expiresAt time.Time
	dataType  string
}

type zsetMember struct {
	member string
	score  float64
}

type RedisStats struct {
	Gets       int64   `json:"gets"`
	Sets       int64   `json:"sets"`
	Deletes    int64   `json:"deletes"`
	HitRatio   float64 `json:"hit_ratio"`
	KeyCount   int     `json:"key_count"`
	MemoryUsed int64   `json:"memory_used_bytes"`
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data:       make(map[string]redisValue),
		hashes:     make(map[string]map[string]string),
		sets:       make(map[string]map[string]struct{}),
		sortedSets: make(map[string][]zsetMember),
		stats:      RedisStats{},
	}
}

func (r *MockRedisClient) Get(key string) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	atomic.AddInt64(&r.stats.Gets, 1)

	if val, exists := r.data[key]; exists {
		if !val.expiresAt.IsZero() && time.Now().After(val.expiresAt) {
			// Key expired
			delete(r.data, key)
			return "", fmt.Errorf("key not found: %s", key)
		}
		return val.value, nil
	}

	return "", fmt.Errorf("key not found: %s", key)
}

func (r *MockRedisClient) Set(key string, value string, expiration time.Duration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	atomic.AddInt64(&r.stats.Sets, 1)

	var expiresAt time.Time
	if expiration > 0 {
		expiresAt = time.Now().Add(expiration)
	}

	r.data[key] = redisValue{
		value:     value,
		expiresAt: expiresAt,
		dataType:  "string",
	}

	return nil
}

func (r *MockRedisClient) Del(key string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	atomic.AddInt64(&r.stats.Deletes, 1)

	delete(r.data, key)
	delete(r.hashes, key)
	delete(r.sets, key)
	delete(r.sortedSets, key)

	return nil
}

func (r *MockRedisClient) Exists(key string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.data[key]
	return exists, nil
}

func (r *MockRedisClient) TTL(key string) (time.Duration, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, exists := r.data[key]; exists {
		if val.expiresAt.IsZero() {
			return -1, nil // No expiration
		}
		remaining := time.Until(val.expiresAt)
		if remaining < 0 {
			return 0, nil // Expired
		}
		return remaining, nil
	}

	return -2, fmt.Errorf("key not found") // Key doesn't exist
}

func (r *MockRedisClient) IncrBy(key string, increment int64) (int64, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var currentValue int64 = 0
	if val, exists := r.data[key]; exists {
		parsed, err := strconv.ParseInt(val.value, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("value is not an integer: %s", val.value)
		}
		currentValue = parsed
	}

	newValue := currentValue + increment
	r.data[key] = redisValue{
		value:    strconv.FormatInt(newValue, 10),
		dataType: "string",
	}

	return newValue, nil
}

func (r *MockRedisClient) HSet(key, field, value string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.hashes[key] == nil {
		r.hashes[key] = make(map[string]string)
	}
	r.hashes[key][field] = value

	return nil
}

func (r *MockRedisClient) HGet(key, field string) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if hash, exists := r.hashes[key]; exists {
		if value, fieldExists := hash[field]; fieldExists {
			return value, nil
		}
	}

	return "", fmt.Errorf("field not found: %s.%s", key, field)
}

func (r *MockRedisClient) HGetAll(key string) (map[string]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if hash, exists := r.hashes[key]; exists {
		result := make(map[string]string)
		for k, v := range hash {
			result[k] = v
		}
		return result, nil
	}

	return make(map[string]string), nil
}

func (r *MockRedisClient) SAdd(key string, members ...string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.sets[key] == nil {
		r.sets[key] = make(map[string]struct{})
	}

	for _, member := range members {
		r.sets[key][member] = struct{}{}
	}

	return nil
}

func (r *MockRedisClient) SMembers(key string) ([]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if set, exists := r.sets[key]; exists {
		members := make([]string, 0, len(set))
		for member := range set {
			members = append(members, member)
		}
		return members, nil
	}

	return []string{}, nil
}

func (r *MockRedisClient) ZAdd(key string, score float64, member string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.sortedSets[key] == nil {
		r.sortedSets[key] = make([]zsetMember, 0)
	}

	// Remove existing member if it exists
	for i, existing := range r.sortedSets[key] {
		if existing.member == member {
			r.sortedSets[key] = append(r.sortedSets[key][:i], r.sortedSets[key][i+1:]...)
			break
		}
	}

	// Add new member
	r.sortedSets[key] = append(r.sortedSets[key], zsetMember{member: member, score: score})

	// Sort by score
	members := r.sortedSets[key]
	for i := 0; i < len(members)-1; i++ {
		for j := i + 1; j < len(members); j++ {
			if members[i].score > members[j].score {
				members[i], members[j] = members[j], members[i]
			}
		}
	}

	return nil
}

func (r *MockRedisClient) ZRange(key string, start, stop int64) ([]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if members, exists := r.sortedSets[key]; exists {
		length := int64(len(members))
		if start < 0 {
			start = length + start
		}
		if stop < 0 {
			stop = length + stop
		}

		if start < 0 {
			start = 0
		}
		if stop >= length {
			stop = length - 1
		}

		if start > stop {
			return []string{}, nil
		}

		result := make([]string, 0, stop-start+1)
		for i := start; i <= stop; i++ {
			result = append(result, members[i].member)
		}
		return result, nil
	}

	return []string{}, nil
}

func (r *MockRedisClient) ZRangeByScore(key string, min, max float64) ([]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if members, exists := r.sortedSets[key]; exists {
		result := make([]string, 0)
		for _, member := range members {
			if member.score >= min && member.score <= max {
				result = append(result, member.member)
			}
		}
		return result, nil
	}

	return []string{}, nil
}

func (r *MockRedisClient) Ping() error {
	return nil // Always successful for mock
}

func (r *MockRedisClient) FlushAll() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.data = make(map[string]redisValue)
	r.hashes = make(map[string]map[string]string)
	r.sets = make(map[string]map[string]struct{})
	r.sortedSets = make(map[string][]zsetMember)

	return nil
}

func (r *MockRedisClient) GetStats() RedisStats {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	stats := r.stats
	stats.KeyCount = len(r.data) + len(r.hashes) + len(r.sets) + len(r.sortedSets)

	// Calculate hit ratio
	total := stats.Gets
	if total > 0 {
		hits := total - stats.Gets/10 // Mock calculation
		stats.HitRatio = float64(hits) / float64(total) * 100
	}

	return stats
}

// 🏪 Redis-Based Cache Service
// Production-ready caching service using Redis patterns
type RedisCacheService struct {
	client     RedisClient
	prefix     string
	defaultTTL time.Duration
}

func NewRedisCacheService(client RedisClient, prefix string, defaultTTL time.Duration) *RedisCacheService {
	return &RedisCacheService{
		client:     client,
		prefix:     prefix,
		defaultTTL: defaultTTL,
	}
}

func (rcs *RedisCacheService) buildKey(key string) string {
	if rcs.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", rcs.prefix, key)
}

// 🎯 Cache-Aside Pattern with Redis
func (rcs *RedisCacheService) GetOrSet(key string, fetchFunc func() (interface{}, error)) (interface{}, error) {
	redisKey := rcs.buildKey(key)

	// Try to get from cache first
	cached, err := rcs.client.Get(redisKey)
	if err == nil {
		var result interface{}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			log.Printf("Cache HIT for key: %s", key)
			return result, nil
		}
	}

	log.Printf("Cache MISS for key: %s", key)

	// Cache miss - fetch data
	data, err := fetchFunc()
	if err != nil {
		return nil, err
	}

	// Store in cache
	jsonData, err := json.Marshal(data)
	if err == nil {
		rcs.client.Set(redisKey, string(jsonData), rcs.defaultTTL)
	}

	return data, nil
}

func (rcs *RedisCacheService) Set(key string, value interface{}, ttl time.Duration) error {
	redisKey := rcs.buildKey(key)
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if ttl == 0 {
		ttl = rcs.defaultTTL
	}

	return rcs.client.Set(redisKey, string(jsonData), ttl)
}

func (rcs *RedisCacheService) Get(key string) (interface{}, error) {
	redisKey := rcs.buildKey(key)
	cached, err := rcs.client.Get(redisKey)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal([]byte(cached), &result)
	return result, err
}

func (rcs *RedisCacheService) Delete(key string) error {
	redisKey := rcs.buildKey(key)
	return rcs.client.Del(redisKey)
}

// 📊 Redis Counter Pattern
type RedisCounter struct {
	client RedisClient
	prefix string
}

func NewRedisCounter(client RedisClient, prefix string) *RedisCounter {
	return &RedisCounter{
		client: client,
		prefix: prefix,
	}
}

func (rc *RedisCounter) Increment(key string, delta int64) (int64, error) {
	redisKey := fmt.Sprintf("%s:counter:%s", rc.prefix, key)
	return rc.client.IncrBy(redisKey, delta)
}

func (rc *RedisCounter) Get(key string) (int64, error) {
	redisKey := fmt.Sprintf("%s:counter:%s", rc.prefix, key)
	value, err := rc.client.Get(redisKey)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(value, 10, 64)
}

// 📈 Redis Leaderboard Pattern
type RedisLeaderboard struct {
	client  RedisClient
	key     string
	maxSize int
}

func NewRedisLeaderboard(client RedisClient, key string, maxSize int) *RedisLeaderboard {
	return &RedisLeaderboard{
		client:  client,
		key:     key,
		maxSize: maxSize,
	}
}

func (rl *RedisLeaderboard) AddScore(player string, score float64) error {
	return rl.client.ZAdd(rl.key, score, player)
}

func (rl *RedisLeaderboard) GetTopN(n int) ([]string, error) {
	// Get top N players (highest scores first)
	return rl.client.ZRange(rl.key, -int64(n), -1)
}

func (rl *RedisLeaderboard) GetRankRange(start, end int) ([]string, error) {
	return rl.client.ZRange(rl.key, int64(start), int64(end))
}

func (rl *RedisLeaderboard) GetPlayersByScoreRange(minScore, maxScore float64) ([]string, error) {
	return rl.client.ZRangeByScore(rl.key, minScore, maxScore)
}

// 🔄 Redis Session Store Pattern
type RedisSessionStore struct {
	client     RedisClient
	keyPrefix  string
	expiration time.Duration
}

func NewRedisSessionStore(client RedisClient, keyPrefix string, expiration time.Duration) *RedisSessionStore {
	return &RedisSessionStore{
		client:     client,
		keyPrefix:  keyPrefix,
		expiration: expiration,
	}
}

func (rss *RedisSessionStore) CreateSession(sessionID string, data map[string]string) error {
	sessionKey := fmt.Sprintf("%s:session:%s", rss.keyPrefix, sessionID)

	for field, value := range data {
		if err := rss.client.HSet(sessionKey, field, value); err != nil {
			return err
		}
	}

	// Set expiration on the session
	return rss.client.Set(sessionKey+":ttl", "1", rss.expiration)
}

func (rss *RedisSessionStore) GetSession(sessionID string) (map[string]string, error) {
	sessionKey := fmt.Sprintf("%s:session:%s", rss.keyPrefix, sessionID)
	return rss.client.HGetAll(sessionKey)
}

func (rss *RedisSessionStore) UpdateSession(sessionID, field, value string) error {
	sessionKey := fmt.Sprintf("%s:session:%s", rss.keyPrefix, sessionID)
	return rss.client.HSet(sessionKey, field, value)
}

func (rss *RedisSessionStore) DeleteSession(sessionID string) error {
	sessionKey := fmt.Sprintf("%s:session:%s", rss.keyPrefix, sessionID)
	rss.client.Del(sessionKey)
	return rss.client.Del(sessionKey + ":ttl")
}

// 🎮 Rate Limiter using Redis
type RedisRateLimiter struct {
	client RedisClient
	prefix string
}

func NewRedisRateLimiter(client RedisClient, prefix string) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		prefix: prefix,
	}
}

func (rrl *RedisRateLimiter) IsAllowed(key string, limit int64, window time.Duration) (bool, error) {
	redisKey := fmt.Sprintf("%s:ratelimit:%s", rrl.prefix, key)

	// Get current count
	current, err := rrl.client.IncrBy(redisKey, 1)
	if err != nil {
		return false, err
	}

	// Set expiration on first request
	if current == 1 {
		rrl.client.Set(redisKey+":ttl", "1", window)
	}

	return current <= limit, nil
}

func (rrl *RedisRateLimiter) GetCurrentCount(key string) (int64, error) {
	redisKey := fmt.Sprintf("%s:ratelimit:%s", rrl.prefix, key)
	value, err := rrl.client.Get(redisKey)
	if err != nil {
		return 0, nil // Key doesn't exist, count is 0
	}

	return strconv.ParseInt(value, 10, 64)
}

// 🧪 Demo Functions for FAANG Interview Practice
func demoRedisBasicOperations() {
	fmt.Println("🔴 FAANG System Design - Redis Integration Patterns Demo")
	fmt.Println(strings.Repeat("=", 60))

	client := NewMockRedisClient()

	fmt.Println("\n📝 Basic Redis Operations:")

	// String operations
	client.Set("user:123:name", "Alice Johnson", 5*time.Minute)
	client.Set("user:123:email", "alice@example.com", 5*time.Minute)

	name, _ := client.Get("user:123:name")
	email, _ := client.Get("user:123:email")
	fmt.Printf("User Name: %s, Email: %s\n", name, email)

	// Counter operations
	client.IncrBy("page:views", 1)
	client.IncrBy("page:views", 5)
	client.IncrBy("page:views", 3)

	views, _ := client.Get("page:views")
	fmt.Printf("Total page views: %s\n", views)

	// Hash operations
	client.HSet("user:456", "name", "Bob Smith")
	client.HSet("user:456", "email", "bob@example.com")
	client.HSet("user:456", "role", "admin")

	userHash, _ := client.HGetAll("user:456")
	fmt.Printf("User Hash: %+v\n", userHash)

	// Set operations
	client.SAdd("user:123:skills", "Go", "Redis", "Docker", "Kubernetes")
	skills, _ := client.SMembers("user:123:skills")
	fmt.Printf("User Skills: %v\n", skills)

	// Sorted set operations (leaderboard)
	client.ZAdd("game:leaderboard", 1500, "Alice")
	client.ZAdd("game:leaderboard", 2100, "Bob")
	client.ZAdd("game:leaderboard", 1800, "Charlie")
	client.ZAdd("game:leaderboard", 2500, "Diana")

	topPlayers, _ := client.ZRange("game:leaderboard", -3, -1)
	fmt.Printf("Top 3 Players: %v\n", topPlayers)

	stats := client.GetStats()
	fmt.Printf("\nRedis Stats: %+v\n", stats)
}

func demoCacheAsidePattern() {
	fmt.Println("\n🎯 Cache-Aside Pattern with Redis:")
	fmt.Println(strings.Repeat("=", 40))

	client := NewMockRedisClient()
	cacheService := NewRedisCacheService(client, "app", 10*time.Minute)

	// Simulate database fetch function
	fetchUserFromDB := func() (interface{}, error) {
		fmt.Println("  -> Fetching from database (expensive operation)")
		time.Sleep(100 * time.Millisecond) // Simulate DB latency
		return map[string]interface{}{
			"id":    123,
			"name":  "Alice Johnson",
			"email": "alice@example.com",
			"role":  "admin",
		}, nil
	}

	// First call - cache miss
	start := time.Now()
	user1, _ := cacheService.GetOrSet("user:123", fetchUserFromDB)
	fmt.Printf("First call (cache miss): %v - took %v\n", user1, time.Since(start))

	// Second call - cache hit
	start = time.Now()
	user2, _ := cacheService.GetOrSet("user:123", fetchUserFromDB)
	fmt.Printf("Second call (cache hit): %v - took %v\n", user2, time.Since(start))

	// Third call - cache hit
	start = time.Now()
	user3, _ := cacheService.GetOrSet("user:123", fetchUserFromDB)
	fmt.Printf("Third call (cache hit): %v - took %v\n", user3, time.Since(start))
}

func demoRedisCounters() {
	fmt.Println("\n📊 Redis Counter Pattern:")
	fmt.Println(strings.Repeat("=", 30))

	client := NewMockRedisClient()
	counter := NewRedisCounter(client, "analytics")

	// Simulate API calls tracking
	endpoints := []string{"api/users", "api/posts", "api/comments"}

	fmt.Println("Simulating API endpoint tracking:")
	for i := 0; i < 10; i++ {
		endpoint := endpoints[rand.Intn(len(endpoints))]
		count, _ := counter.Increment(endpoint, 1)
		fmt.Printf("  %s called - total: %d\n", endpoint, count)
	}

	fmt.Println("\nFinal API call counts:")
	for _, endpoint := range endpoints {
		count, _ := counter.Get(endpoint)
		fmt.Printf("  %s: %d calls\n", endpoint, count)
	}
}

func demoRedisLeaderboard() {
	fmt.Println("\n📈 Redis Leaderboard Pattern:")
	fmt.Println(strings.Repeat("=", 35))

	client := NewMockRedisClient()
	leaderboard := NewRedisLeaderboard(client, "game:scores", 10)

	// Add player scores
	players := map[string]float64{
		"Alice":   2500,
		"Bob":     1800,
		"Charlie": 3200,
		"Diana":   2100,
		"Eve":     1500,
		"Frank":   2800,
	}

	fmt.Println("Adding player scores:")
	for player, score := range players {
		leaderboard.AddScore(player, score)
		fmt.Printf("  %s: %.0f points\n", player, score)
	}

	// Get top 3 players
	top3, _ := leaderboard.GetTopN(3)
	fmt.Printf("\nTop 3 Players: %v\n", top3)

	// Get players with scores between 2000-3000
	midRange, _ := leaderboard.GetPlayersByScoreRange(2000, 3000)
	fmt.Printf("Players with 2000-3000 points: %v\n", midRange)
}

func demoRedisSession() {
	fmt.Println("\n🔄 Redis Session Store Pattern:")
	fmt.Println(strings.Repeat("=", 35))

	client := NewMockRedisClient()
	sessionStore := NewRedisSessionStore(client, "webapp", 30*time.Minute)

	// Create session
	sessionID := "sess_abc123"
	sessionData := map[string]string{
		"user_id":  "123",
		"username": "alice",
		"role":     "admin",
		"login_at": time.Now().Format(time.RFC3339),
	}

	sessionStore.CreateSession(sessionID, sessionData)
	fmt.Printf("Created session: %s\n", sessionID)

	// Retrieve session
	retrievedSession, _ := sessionStore.GetSession(sessionID)
	fmt.Printf("Retrieved session: %+v\n", retrievedSession)

	// Update session
	sessionStore.UpdateSession(sessionID, "last_activity", time.Now().Format(time.RFC3339))
	fmt.Println("Updated session with last activity")

	// Get updated session
	updatedSession, _ := sessionStore.GetSession(sessionID)
	fmt.Printf("Updated session: %+v\n", updatedSession)
}

func demoRedisRateLimiter() {
	fmt.Println("\n🎮 Redis Rate Limiter Pattern:")
	fmt.Println(strings.Repeat("=", 35))

	client := NewMockRedisClient()
	rateLimiter := NewRedisRateLimiter(client, "api")

	userID := "user123"
	limit := int64(5)
	window := 1 * time.Minute

	fmt.Printf("Rate limiting user %s to %d requests per minute:\n", userID, limit)

	// Simulate requests
	for i := 1; i <= 8; i++ {
		allowed, _ := rateLimiter.IsAllowed(userID, limit, window)
		currentCount, _ := rateLimiter.GetCurrentCount(userID)

		if allowed {
			fmt.Printf("  Request %d: ALLOWED (count: %d/%d)\n", i, currentCount, limit)
		} else {
			fmt.Printf("  Request %d: BLOCKED (count: %d/%d)\n", i, currentCount, limit)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func runRedisDemo() {
	fmt.Println("🎯 FAANG System Design Interview - Redis Integration Mastery")
	fmt.Println("Essential caching patterns for L3/L4 technical interviews")
	fmt.Println()

	// Run all demos
	demoRedisBasicOperations()
	demoCacheAsidePattern()
	demoRedisCounters()
	demoRedisLeaderboard()
	demoRedisSession()
	demoRedisRateLimiter()

	fmt.Println("\n🎪 Key Takeaways for FAANG Interviews:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. Cache-Aside: Most common Redis caching pattern")
	fmt.Println("2. Counters: Track metrics and analytics efficiently")
	fmt.Println("3. Leaderboards: Sorted sets for ranking systems")
	fmt.Println("4. Sessions: Hash-based user session management")
	fmt.Println("5. Rate Limiting: Prevent abuse with counters + TTL")
	fmt.Println("6. Data Structures: Choose right Redis type for use case")
	fmt.Println("\n💪 Practice designing Instagram/Twitter with Redis caching!")
}

// Uncomment to run this demo independently
// func main() {
//     rand.Seed(time.Now().UnixNano())
//     runRedisDemo()
// }
