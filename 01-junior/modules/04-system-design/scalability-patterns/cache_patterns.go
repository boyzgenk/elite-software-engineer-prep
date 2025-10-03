package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// 🎯 FAANG System Design - Caching Strategies Implementation
// Essential patterns for L3/L4 interviews: Cache-aside, Write-through, Write-behind

// CacheItem represents a cached data item with metadata
type CacheItem struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	ExpiresAt  time.Time   `json:"expires_at"`
	AccessedAt time.Time   `json:"accessed_at"`
	Size       int         `json:"size"` // bytes
}

// IsExpired checks if cache item has expired
func (c *CacheItem) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// UpdateAccess updates the last accessed time (for LRU eviction)
func (c *CacheItem) UpdateAccess() {
	c.AccessedAt = time.Now()
}

// CacheStats holds cache performance metrics
type CacheStats struct {
	Hits      int64 `json:"hits"`
	Misses    int64 `json:"misses"`
	Sets      int64 `json:"sets"`
	Deletes   int64 `json:"deletes"`
	Evictions int64 `json:"evictions"`
	TotalSize int   `json:"total_size"`
	ItemCount int   `json:"item_count"`
}

// HitRatio calculates cache hit ratio percentage
func (s *CacheStats) HitRatio() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0.0
	}
	return float64(s.Hits) / float64(total) * 100.0
}

// Cache interface defines caching operations
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	Stats() CacheStats
	Size() int
}

// 💾 In-Memory Cache with LRU Eviction
// Essential for FAANG interviews - implements common caching patterns
type InMemoryCache struct {
	data       map[string]*CacheItem
	maxSize    int // Maximum number of items
	maxMemory  int // Maximum memory in bytes
	defaultTTL time.Duration
	stats      CacheStats
	mutex      sync.RWMutex
}

func NewInMemoryCache(maxSize int, maxMemory int, defaultTTL time.Duration) *InMemoryCache {
	return &InMemoryCache{
		data:       make(map[string]*CacheItem),
		maxSize:    maxSize,
		maxMemory:  maxMemory,
		defaultTTL: defaultTTL,
		stats:      CacheStats{},
	}
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		c.stats.Misses++
		return nil, false
	}

	// Check expiration
	if item.IsExpired() {
		delete(c.data, key)
		c.stats.TotalSize -= item.Size
		c.stats.ItemCount--
		c.stats.Misses++
		return nil, false
	}

	// Update access time for LRU
	item.UpdateAccess()
	c.stats.Hits++
	return item.Value, true
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ttl == 0 {
		ttl = c.defaultTTL
	}

	// Calculate item size (simplified)
	valueBytes, _ := json.Marshal(value)
	itemSize := len(key) + len(valueBytes) + 100 // overhead estimation

	// Check if we need to evict items
	if err := c.evictIfNeeded(itemSize); err != nil {
		return err
	}

	// Create new cache item
	item := &CacheItem{
		Key:        key,
		Value:      value,
		ExpiresAt:  time.Now().Add(ttl),
		AccessedAt: time.Now(),
		Size:       itemSize,
	}

	// Update existing item or add new one
	if existingItem, exists := c.data[key]; exists {
		c.stats.TotalSize -= existingItem.Size
	} else {
		c.stats.ItemCount++
	}

	c.data[key] = item
	c.stats.TotalSize += itemSize
	c.stats.Sets++

	return nil
}

func (c *InMemoryCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, exists := c.data[key]; exists {
		delete(c.data, key)
		c.stats.TotalSize -= item.Size
		c.stats.ItemCount--
		c.stats.Deletes++
	}

	return nil
}

func (c *InMemoryCache) Clear() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*CacheItem)
	c.stats.TotalSize = 0
	c.stats.ItemCount = 0

	return nil
}

func (c *InMemoryCache) Stats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.stats
}

func (c *InMemoryCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.stats.ItemCount
}

// evictIfNeeded performs LRU eviction if cache is full
func (c *InMemoryCache) evictIfNeeded(newItemSize int) error {
	// Check size limit
	for c.stats.ItemCount >= c.maxSize {
		if err := c.evictLRU(); err != nil {
			return err
		}
	}

	// Check memory limit
	for c.stats.TotalSize+newItemSize > c.maxMemory {
		if err := c.evictLRU(); err != nil {
			return err
		}
	}

	return nil
}

// evictLRU removes the least recently used item
func (c *InMemoryCache) evictLRU() error {
	if len(c.data) == 0 {
		return fmt.Errorf("cache is empty, cannot evict")
	}

	var oldestKey string
	var oldestTime time.Time = time.Now()

	// Find least recently used item
	for key, item := range c.data {
		if item.AccessedAt.Before(oldestTime) {
			oldestTime = item.AccessedAt
			oldestKey = key
		}
	}

	// Remove oldest item
	if oldestItem, exists := c.data[oldestKey]; exists {
		delete(c.data, oldestKey)
		c.stats.TotalSize -= oldestItem.Size
		c.stats.ItemCount--
		c.stats.Evictions++
	}

	return nil
}

// 🏪 Database simulation for caching patterns
type Database struct {
	data  map[string]interface{}
	mutex sync.RWMutex
	stats struct {
		reads  int64
		writes int64
	}
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]interface{}),
	}
}

func (db *Database) Get(key string) (interface{}, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	// Simulate database latency
	time.Sleep(10 * time.Millisecond)

	db.stats.reads++
	value, exists := db.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return value, nil
}

func (db *Database) Set(key string, value interface{}) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Simulate database write latency
	time.Sleep(20 * time.Millisecond)

	db.data[key] = value
	db.stats.writes++
	return nil
}

func (db *Database) Stats() (int64, int64) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return db.stats.reads, db.stats.writes
}

// 🎯 Cache-Aside Pattern (Lazy Loading)
// Most common pattern in FAANG interviews
type CacheAsideService struct {
	cache Cache
	db    *Database
}

func NewCacheAsideService(cache Cache, db *Database) *CacheAsideService {
	return &CacheAsideService{
		cache: cache,
		db:    db,
	}
}

func (s *CacheAsideService) Get(key string) (interface{}, error) {
	// 1. Try cache first
	if value, found := s.cache.Get(key); found {
		log.Printf("Cache HIT for key: %s", key)
		return value, nil
	}

	log.Printf("Cache MISS for key: %s", key)

	// 2. Cache miss - fetch from database
	value, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}

	// 3. Store in cache for next time
	s.cache.Set(key, value, 5*time.Minute)

	return value, nil
}

func (s *CacheAsideService) Set(key string, value interface{}) error {
	// 1. Write to database first
	if err := s.db.Set(key, value); err != nil {
		return err
	}

	// 2. Invalidate cache (or update it)
	s.cache.Delete(key)

	return nil
}

// 📝 Write-Through Cache Pattern
// Critical for data consistency discussions in FAANG interviews
type WriteThroughService struct {
	cache Cache
	db    *Database
}

func NewWriteThroughService(cache Cache, db *Database) *WriteThroughService {
	return &WriteThroughService{
		cache: cache,
		db:    db,
	}
}

func (s *WriteThroughService) Get(key string) (interface{}, error) {
	// Try cache first
	if value, found := s.cache.Get(key); found {
		return value, nil
	}

	// Cache miss - fetch from database
	value, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}

	// Store in cache
	s.cache.Set(key, value, 10*time.Minute)
	return value, nil
}

func (s *WriteThroughService) Set(key string, value interface{}) error {
	// 1. Write to database first (synchronous)
	if err := s.db.Set(key, value); err != nil {
		return err
	}

	// 2. Write to cache (data is always consistent)
	if err := s.cache.Set(key, value, 10*time.Minute); err != nil {
		log.Printf("Warning: Cache write failed for key %s: %v", key, err)
		// Continue - database write succeeded
	}

	return nil
}

// 📤 Write-Behind (Write-Back) Cache Pattern
// Advanced pattern for high-write scenarios
type WriteBehindService struct {
	cache      Cache
	db         *Database
	writeQueue chan writeOperation
	ctx        context.Context
	cancel     context.CancelFunc
}

type writeOperation struct {
	key   string
	value interface{}
	retry int
}

func NewWriteBehindService(cache Cache, db *Database) *WriteBehindService {
	ctx, cancel := context.WithCancel(context.Background())
	service := &WriteBehindService{
		cache:      cache,
		db:         db,
		writeQueue: make(chan writeOperation, 1000),
		ctx:        ctx,
		cancel:     cancel,
	}

	// Start background writer
	go service.backgroundWriter()

	return service
}

func (s *WriteBehindService) Get(key string) (interface{}, error) {
	// Try cache first
	if value, found := s.cache.Get(key); found {
		return value, nil
	}

	// Cache miss - fetch from database
	value, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}

	s.cache.Set(key, value, 15*time.Minute)
	return value, nil
}

func (s *WriteBehindService) Set(key string, value interface{}) error {
	// 1. Write to cache immediately (fast response)
	if err := s.cache.Set(key, value, 15*time.Minute); err != nil {
		return err
	}

	// 2. Queue database write for background processing
	select {
	case s.writeQueue <- writeOperation{key: key, value: value, retry: 0}:
		// Successfully queued
	default:
		// Queue is full - could implement backpressure logic here
		log.Printf("Warning: Write queue full for key %s", key)
	}

	return nil
}

func (s *WriteBehindService) backgroundWriter() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case op := <-s.writeQueue:
			if err := s.db.Set(op.key, op.value); err != nil {
				log.Printf("Background write failed for key %s: %v", op.key, err)

				// Retry logic
				if op.retry < 3 {
					op.retry++
					time.Sleep(time.Duration(op.retry) * time.Second)

					select {
					case s.writeQueue <- op:
						// Retry queued
					default:
						log.Printf("Failed to retry write for key %s", op.key)
					}
				}
			}
		}
	}
}

func (s *WriteBehindService) Stop() {
	s.cancel()
	close(s.writeQueue)
}

// 🌐 CDN Simulation for Content Delivery
type CDNCache struct {
	regions map[string]Cache // Different cache instances for different regions
	mutex   sync.RWMutex
}

func NewCDNCache(regions []string) *CDNCache {
	cdn := &CDNCache{
		regions: make(map[string]Cache),
	}

	for _, region := range regions {
		cdn.regions[region] = NewInMemoryCache(1000, 10*1024*1024, 1*time.Hour)
	}

	return cdn
}

func (cdn *CDNCache) Get(region, key string) (interface{}, bool) {
	cdn.mutex.RLock()
	defer cdn.mutex.RUnlock()

	if cache, exists := cdn.regions[region]; exists {
		return cache.Get(key)
	}

	return nil, false
}

func (cdn *CDNCache) Set(region, key string, value interface{}, ttl time.Duration) error {
	cdn.mutex.RLock()
	defer cdn.mutex.RUnlock()

	if cache, exists := cdn.regions[region]; exists {
		return cache.Set(key, value, ttl)
	}

	return fmt.Errorf("region not found: %s", region)
}

func (cdn *CDNCache) GlobalSet(key string, value interface{}, ttl time.Duration) {
	cdn.mutex.RLock()
	defer cdn.mutex.RUnlock()

	// Distribute content to all regions
	for region, cache := range cdn.regions {
		if err := cache.Set(key, value, ttl); err != nil {
			log.Printf("Failed to cache in region %s: %v", region, err)
		}
	}
}

func (cdn *CDNCache) GetRegionStats(region string) CacheStats {
	cdn.mutex.RLock()
	defer cdn.mutex.RUnlock()

	if cache, exists := cdn.regions[region]; exists {
		return cache.Stats()
	}

	return CacheStats{}
}

// 🧪 Demo Functions for FAANG Interview Practice
func demoCachePatterns() {
	fmt.Println("🚀 FAANG System Design - Caching Strategies Demo")
	fmt.Println(strings.Repeat("=", 60))

	// Initialize components
	db := NewDatabase()
	cache := NewInMemoryCache(100, 1024*1024, 5*time.Minute)

	// Pre-populate database
	testData := map[string]interface{}{
		"user:123":    map[string]string{"name": "Alice", "email": "alice@example.com"},
		"user:456":    map[string]string{"name": "Bob", "email": "bob@example.com"},
		"product:789": map[string]interface{}{"name": "Laptop", "price": 999.99},
	}

	for key, value := range testData {
		db.Set(key, value)
	}

	fmt.Println("\n🎯 Cache-Aside Pattern Demo:")
	cacheAside := NewCacheAsideService(cache, db)

	// First access - cache miss
	start := time.Now()
	value, _ := cacheAside.Get("user:123")
	fmt.Printf("First access (cache miss): %v - took %v\n", value, time.Since(start))

	// Second access - cache hit
	start = time.Now()
	value, _ = cacheAside.Get("user:123")
	fmt.Printf("Second access (cache hit): %v - took %v\n", value, time.Since(start))

	// Update data - cache invalidated
	cacheAside.Set("user:123", map[string]string{"name": "Alice Updated", "email": "alice.new@example.com"})
	value, _ = cacheAside.Get("user:123")
	fmt.Printf("After update: %v\n", value)

	fmt.Println("\n📝 Write-Through Pattern Demo:")
	writeThrough := NewWriteThroughService(cache, db)

	start = time.Now()
	writeThrough.Set("user:999", map[string]string{"name": "Charlie", "email": "charlie@example.com"})
	fmt.Printf("Write-through set took: %v\n", time.Since(start))

	start = time.Now()
	value, _ = writeThrough.Get("user:999")
	fmt.Printf("Immediate read (cache hit): %v - took %v\n", value, time.Since(start))

	fmt.Println("\n📤 Write-Behind Pattern Demo:")
	writeBehind := NewWriteBehindService(cache, db)

	start = time.Now()
	writeBehind.Set("user:888", map[string]string{"name": "David", "email": "david@example.com"})
	fmt.Printf("Write-behind set took: %v (async DB write)\n", time.Since(start))

	// Wait for background write to complete
	time.Sleep(100 * time.Millisecond)

	// Verify data is in database
	dbValue, _ := db.Get("user:888")
	fmt.Printf("Background write completed: %v\n", dbValue)

	writeBehind.Stop()
}

func demoCDNCache() {
	fmt.Println("\n🌐 CDN Cache Pattern Demo:")
	fmt.Println(strings.Repeat("=", 40))

	regions := []string{"us-east", "us-west", "eu-west", "asia-pacific"}
	cdn := NewCDNCache(regions)

	// Simulate content distribution
	content := map[string]interface{}{
		"image.jpg": "binary_image_data_here",
		"script.js": "console.log('Hello CDN');",
		"style.css": "body { margin: 0; }",
	}

	fmt.Println("Distributing content globally...")
	for key, value := range content {
		cdn.GlobalSet(key, value, 2*time.Hour)
	}

	// Simulate regional access
	for _, region := range regions {
		fmt.Printf("\nRegion %s access patterns:\n", region)

		for key := range content {
			start := time.Now()
			_, found := cdn.Get(region, key)
			if found {
				fmt.Printf("  %s: Cache HIT - took %v\n", key, time.Since(start))
			}
		}

		stats := cdn.GetRegionStats(region)
		fmt.Printf("  Stats: Hit ratio %.1f%%, Items: %d\n", stats.HitRatio(), stats.ItemCount)
	}
}

func demoPerformanceComparison() {
	fmt.Println("\n📊 Performance Comparison for Interview Discussion:")
	fmt.Println(strings.Repeat("=", 60))

	db := NewDatabase()
	cache := NewInMemoryCache(1000, 10*1024*1024, 10*time.Minute)

	// Pre-populate data
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key:%d", i)
		value := map[string]interface{}{
			"id":   i,
			"data": generateRandomString(100),
		}
		db.Set(key, value)
	}

	cacheAside := NewCacheAsideService(cache, db)
	_ = NewWriteThroughService(cache, db)

	// Benchmark reads
	fmt.Println("\nRead Performance (100 operations):")

	// Cache-aside reads (with warm cache)
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("key:%d", i)
		cacheAside.Get(key)
	}

	start := time.Now()
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key:%d", i%50)
		cacheAside.Get(key)
	}
	cacheAsideTime := time.Since(start)

	// Direct database reads
	start = time.Now()
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key:%d", i%50)
		db.Get(key)
	}
	dbTime := time.Since(start)

	fmt.Printf("Cache-aside reads: %v\n", cacheAsideTime)
	fmt.Printf("Direct DB reads: %v\n", dbTime)
	fmt.Printf("Performance improvement: %.1fx faster\n", float64(dbTime)/float64(cacheAsideTime))

	// Display cache statistics
	stats := cache.Stats()
	fmt.Printf("\nCache Statistics:\n")
	fmt.Printf("Hit Ratio: %.1f%%\n", stats.HitRatio())
	fmt.Printf("Total Items: %d\n", stats.ItemCount)
	fmt.Printf("Memory Usage: %d bytes\n", stats.TotalSize)
	fmt.Printf("Evictions: %d\n", stats.Evictions)

	// Display database statistics
	dbReads, dbWrites := db.Stats()
	fmt.Printf("\nDatabase Statistics:\n")
	fmt.Printf("Reads: %d\n", dbReads)
	fmt.Printf("Writes: %d\n", dbWrites)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// 🎯 Cache Key Generation Strategies
func demonstrateKeyStrategies() {
	fmt.Println("\n🔑 Cache Key Generation Strategies:")
	fmt.Println(strings.Repeat("=", 45))

	examples := []struct {
		description string
		key         string
	}{
		{"User profile", "user:profile:12345"},
		{"User posts", "user:posts:12345:page:1"},
		{"Product details", "product:details:67890"},
		{"Search results", "search:results:" + hashString("laptop computers")},
		{"API response", "api:v1:users:12345:friends"},
		{"Session data", "session:" + hashString("session_token_abc123")},
	}

	for _, example := range examples {
		fmt.Printf("%-20s -> %s\n", example.description, example.key)
	}

	fmt.Println("\n💡 Key Design Principles:")
	fmt.Println("- Use consistent naming conventions (namespace:type:id)")
	fmt.Println("- Hash long or sensitive values")
	fmt.Println("- Include version info for API responses")
	fmt.Println("- Use hierarchical structure for easy invalidation")
}

func hashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:8]) // Use first 8 bytes for shorter keys
}

func runCacheDemo() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("🎯 FAANG System Design Interview - Caching Mastery")
	fmt.Println("Essential patterns for L3/L4 technical interviews")
	fmt.Println()

	// Run all demos
	demoCachePatterns()
	demoCDNCache()
	demoPerformanceComparison()
	demonstrateKeyStrategies()

	fmt.Println("\n🎪 Key Takeaways for FAANG Interviews:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. Cache-Aside: Best for read-heavy workloads")
	fmt.Println("2. Write-Through: Strong consistency, slower writes")
	fmt.Println("3. Write-Behind: High write performance, eventual consistency")
	fmt.Println("4. CDN: Geographic distribution for global scale")
	fmt.Println("5. LRU Eviction: Memory management under pressure")
	fmt.Println("6. TTL: Automatic cache invalidation strategy")
	fmt.Println("\n💪 Practice with Twitter/Netflix/Instagram cache design!")
}

func main() {
	runCacheDemo()
}
