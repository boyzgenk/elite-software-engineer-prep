package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

// 🌐 FAANG System Design - CDN Patterns & Edge Caching
// Essential content delivery patterns for L3/L4 interviews

// CDNNode represents a CDN edge server/node
type CDNNode struct {
	ID       string
	Region   string
	Location string
	Cache    Cache
	Stats    CDNNodeStats
	mutex    sync.RWMutex
}

type CDNNodeStats struct {
	Requests    int64   `json:"requests"`
	CacheHits   int64   `json:"cache_hits"`
	CacheMisses int64   `json:"cache_misses"`
	Bandwidth   int64   `json:"bandwidth_bytes"`
	HitRatio    float64 `json:"hit_ratio"`
	Latency     float64 `json:"avg_latency_ms"`
}

func NewCDNNode(id, region, location string, cacheSize int) *CDNNode {
	return &CDNNode{
		ID:       id,
		Region:   region,
		Location: location,
		Cache:    NewInMemoryCache(cacheSize, 100*1024*1024, 24*time.Hour), // 100MB, 24h TTL
		Stats:    CDNNodeStats{},
	}
}

func (n *CDNNode) Get(key string) (interface{}, bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.Stats.Requests++

	value, found := n.Cache.Get(key)
	if found {
		n.Stats.CacheHits++
		log.Printf("CDN HIT [%s]: %s", n.ID, key)
	} else {
		n.Stats.CacheMisses++
		log.Printf("CDN MISS [%s]: %s", n.ID, key)
	}

	n.updateHitRatio()
	return value, found
}

func (n *CDNNode) Set(key string, value interface{}, ttl time.Duration) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	return n.Cache.Set(key, value, ttl)
}

func (n *CDNNode) updateHitRatio() {
	total := n.Stats.CacheHits + n.Stats.CacheMisses
	if total > 0 {
		n.Stats.HitRatio = float64(n.Stats.CacheHits) / float64(total) * 100.0
	}
}

func (n *CDNNode) GetStats() CDNNodeStats {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.Stats
}

// 🗺️ CDN Network - Global Distribution System
type CDNNetwork struct {
	nodes   map[string]*CDNNode
	origins map[string]string // key -> origin URL mapping
	routing CDNRoutingStrategy
	mutex   sync.RWMutex
}

type CDNRoutingStrategy interface {
	SelectNode(userLocation string, nodes []*CDNNode) *CDNNode
}

func NewCDNNetwork(routing CDNRoutingStrategy) *CDNNetwork {
	return &CDNNetwork{
		nodes:   make(map[string]*CDNNode),
		origins: make(map[string]string),
		routing: routing,
	}
}

func (cdn *CDNNetwork) AddNode(node *CDNNode) {
	cdn.mutex.Lock()
	defer cdn.mutex.Unlock()
	cdn.nodes[node.ID] = node
}

func (cdn *CDNNetwork) AddOrigin(key, originURL string) {
	cdn.mutex.Lock()
	defer cdn.mutex.Unlock()
	cdn.origins[key] = originURL
}

func (cdn *CDNNetwork) Get(key, userLocation string) (interface{}, error) {
	cdn.mutex.RLock()
	nodes := make([]*CDNNode, 0, len(cdn.nodes))
	for _, node := range cdn.nodes {
		nodes = append(nodes, node)
	}
	cdn.mutex.RUnlock()

	// Select best CDN node based on routing strategy
	selectedNode := cdn.routing.SelectNode(userLocation, nodes)
	if selectedNode == nil {
		return nil, fmt.Errorf("no CDN nodes available")
	}

	// Try to get from selected node
	value, found := selectedNode.Get(key)
	if found {
		return value, nil
	}

	// Cache miss - fetch from origin
	originValue, err := cdn.fetchFromOrigin(key)
	if err != nil {
		return nil, err
	}

	// Cache in selected node and propagate to nearby nodes
	selectedNode.Set(key, originValue, 24*time.Hour)
	go cdn.propagateToNearbyNodes(key, originValue, selectedNode)

	return originValue, nil
}

func (cdn *CDNNetwork) fetchFromOrigin(key string) (interface{}, error) {
	// Simulate origin fetch
	time.Sleep(200 * time.Millisecond) // Origin latency

	// Mock content based on key
	content := map[string]interface{}{
		"image.jpg":     fmt.Sprintf("binary_image_data_%s", key),
		"style.css":     fmt.Sprintf("body { background: #%s; }", generateHash(key)[:6]),
		"script.js":     fmt.Sprintf("console.log('Script for %s');", key),
		"video.mp4":     fmt.Sprintf("video_stream_data_%s", key),
		"api/data.json": fmt.Sprintf(`{"data": "%s", "timestamp": "%s"}`, key, time.Now().Format(time.RFC3339)),
	}

	if data, exists := content[key]; exists {
		return data, nil
	}

	return fmt.Sprintf("content_for_%s", key), nil
}

func (cdn *CDNNetwork) propagateToNearbyNodes(key string, value interface{}, sourceNode *CDNNode) {
	// Propagate to nodes in the same region
	cdn.mutex.RLock()
	defer cdn.mutex.RUnlock()

	for _, node := range cdn.nodes {
		if node.ID != sourceNode.ID && node.Region == sourceNode.Region {
			node.Set(key, value, 12*time.Hour) // Shorter TTL for propagated content
			log.Printf("Propagated %s to node %s", key, node.ID)
		}
	}
}

func (cdn *CDNNetwork) GetGlobalStats() map[string]CDNNodeStats {
	cdn.mutex.RLock()
	defer cdn.mutex.RUnlock()

	stats := make(map[string]CDNNodeStats)
	for id, node := range cdn.nodes {
		stats[id] = node.GetStats()
	}
	return stats
}

// 📍 Geographic Routing Strategy
type GeographicRouting struct{}

func (gr *GeographicRouting) SelectNode(userLocation string, nodes []*CDNNode) *CDNNode {
	// Simple geographic routing based on region matching
	for _, node := range nodes {
		if strings.Contains(userLocation, node.Region) {
			return node
		}
	}

	// Fallback to any available node
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

// ⚡ Latency-Based Routing Strategy
type LatencyBasedRouting struct {
	latencyMap map[string]map[string]float64 // userLocation -> nodeID -> latency
	mutex      sync.RWMutex
}

func NewLatencyBasedRouting() *LatencyBasedRouting {
	return &LatencyBasedRouting{
		latencyMap: make(map[string]map[string]float64),
	}
}

func (lbr *LatencyBasedRouting) SetLatency(userLocation, nodeID string, latency float64) {
	lbr.mutex.Lock()
	defer lbr.mutex.Unlock()

	if lbr.latencyMap[userLocation] == nil {
		lbr.latencyMap[userLocation] = make(map[string]float64)
	}
	lbr.latencyMap[userLocation][nodeID] = latency
}

func (lbr *LatencyBasedRouting) SelectNode(userLocation string, nodes []*CDNNode) *CDNNode {
	lbr.mutex.RLock()
	defer lbr.mutex.RUnlock()

	if nodeLatencies, exists := lbr.latencyMap[userLocation]; exists {
		var bestNode *CDNNode
		bestLatency := float64(1000) // 1 second max

		for _, node := range nodes {
			if latency, hasLatency := nodeLatencies[node.ID]; hasLatency {
				if latency < bestLatency {
					bestLatency = latency
					bestNode = node
				}
			}
		}

		if bestNode != nil {
			return bestNode
		}
	}

	// Fallback to geographic routing
	geographic := &GeographicRouting{}
	return geographic.SelectNode(userLocation, nodes)
}

// 🔄 Load-Based Routing Strategy
type LoadBasedRouting struct{}

func (lbr *LoadBasedRouting) SelectNode(userLocation string, nodes []*CDNNode) *CDNNode {
	if len(nodes) == 0 {
		return nil
	}

	// Select node with lowest request count
	var bestNode *CDNNode
	lowestLoad := int64(1000000)

	for _, node := range nodes {
		stats := node.GetStats()
		if stats.Requests < lowestLoad {
			lowestLoad = stats.Requests
			bestNode = node
		}
	}

	return bestNode
}

// 🎯 Content Invalidation System
type CDNInvalidation struct {
	cdn       *CDNNetwork
	patterns  []string
	timestamp time.Time
}

func NewCDNInvalidation(cdn *CDNNetwork) *CDNInvalidation {
	return &CDNInvalidation{
		cdn:       cdn,
		patterns:  make([]string, 0),
		timestamp: time.Now(),
	}
}

func (ci *CDNInvalidation) InvalidatePattern(pattern string) error {
	ci.patterns = append(ci.patterns, pattern)
	ci.timestamp = time.Now()

	log.Printf("Invalidating pattern: %s", pattern)

	// Invalidate across all nodes
	ci.cdn.mutex.RLock()
	defer ci.cdn.mutex.RUnlock()

	for _, node := range ci.cdn.nodes {
		go ci.invalidateOnNode(node, pattern)
	}

	return nil
}

func (ci *CDNInvalidation) invalidateOnNode(node *CDNNode, pattern string) {
	// Simple pattern matching - in production, use more sophisticated matching
	// For demo, we'll clear all cache if pattern contains "*"
	if strings.Contains(pattern, "*") {
		node.Cache.Clear()
		log.Printf("Cleared all cache on node %s", node.ID)
	} else {
		node.Cache.Delete(pattern)
		log.Printf("Deleted %s from node %s", pattern, node.ID)
	}
}

func (ci *CDNInvalidation) GetInvalidationHistory() []string {
	return ci.patterns
}

// 📊 CDN Analytics & Monitoring
type CDNAnalytics struct {
	cdn     *CDNNetwork
	metrics map[string]interface{}
	mutex   sync.RWMutex
}

func NewCDNAnalytics(cdn *CDNNetwork) *CDNAnalytics {
	return &CDNAnalytics{
		cdn:     cdn,
		metrics: make(map[string]interface{}),
	}
}

func (ca *CDNAnalytics) CollectMetrics() map[string]interface{} {
	ca.mutex.Lock()
	defer ca.mutex.Unlock()

	globalStats := ca.cdn.GetGlobalStats()

	totalRequests := int64(0)
	totalHits := int64(0)
	totalBandwidth := int64(0)
	nodeCount := len(globalStats)

	for _, stats := range globalStats {
		totalRequests += stats.Requests
		totalHits += stats.CacheHits
		totalBandwidth += stats.Bandwidth
	}

	globalHitRatio := float64(0)
	if totalRequests > 0 {
		globalHitRatio = float64(totalHits) / float64(totalRequests) * 100.0
	}

	ca.metrics = map[string]interface{}{
		"total_requests":   totalRequests,
		"total_hits":       totalHits,
		"global_hit_ratio": globalHitRatio,
		"total_bandwidth":  totalBandwidth,
		"active_nodes":     nodeCount,
		"nodes_stats":      globalStats,
		"collection_time":  time.Now(),
	}

	return ca.metrics
}

func (ca *CDNAnalytics) GetTopPerformingNodes(limit int) []string {
	ca.mutex.RLock()
	defer ca.mutex.RUnlock()

	globalStats := ca.cdn.GetGlobalStats()

	type nodePerformance struct {
		id       string
		hitRatio float64
	}

	performances := make([]nodePerformance, 0, len(globalStats))
	for nodeID, stats := range globalStats {
		performances = append(performances, nodePerformance{
			id:       nodeID,
			hitRatio: stats.HitRatio,
		})
	}

	// Sort by hit ratio
	sort.Slice(performances, func(i, j int) bool {
		return performances[i].hitRatio > performances[j].hitRatio
	})

	result := make([]string, 0, limit)
	for i, perf := range performances {
		if i >= limit {
			break
		}
		result = append(result, perf.id)
	}

	return result
}

// 🧪 Demo Functions for FAANG Interview Practice
func demoCDNBasicOperations() {
	fmt.Println("🌐 FAANG System Design - CDN Patterns Demo")
	fmt.Println(strings.Repeat("=", 50))

	// Create CDN network with geographic routing
	routing := &GeographicRouting{}
	cdn := NewCDNNetwork(routing)

	// Add CDN nodes across different regions
	cdn.AddNode(NewCDNNode("us-east-1", "us-east", "New York", 1000))
	cdn.AddNode(NewCDNNode("us-west-1", "us-west", "San Francisco", 1000))
	cdn.AddNode(NewCDNNode("eu-west-1", "eu-west", "London", 1000))
	cdn.AddNode(NewCDNNode("ap-south-1", "ap-south", "Mumbai", 1000))

	// Add origin mappings
	cdn.AddOrigin("image.jpg", "https://origin.example.com/image.jpg")
	cdn.AddOrigin("style.css", "https://origin.example.com/style.css")
	cdn.AddOrigin("script.js", "https://origin.example.com/script.js")

	fmt.Println("\n📍 Testing Geographic Routing:")

	// Simulate requests from different locations
	locations := []string{"us-east", "us-west", "eu-west", "ap-south"}
	content := []string{"image.jpg", "style.css", "script.js"}

	for _, location := range locations {
		fmt.Printf("\nRequests from %s:\n", location)
		for _, item := range content {
			start := time.Now()
			data, err := cdn.Get(item, location)
			latency := time.Since(start)

			if err != nil {
				fmt.Printf("  %s: ERROR - %v\n", item, err)
			} else {
				fmt.Printf("  %s: SUCCESS (took %v) - %s\n",
					item, latency, truncateString(fmt.Sprintf("%v", data), 50))
			}
		}
	}
}

func demoCDNLatencyRouting() {
	fmt.Println("\n⚡ Latency-Based Routing Demo:")
	fmt.Println(strings.Repeat("=", 35))

	// Create CDN with latency-based routing
	latencyRouting := NewLatencyBasedRouting()
	cdn := NewCDNNetwork(latencyRouting)

	// Add nodes
	cdn.AddNode(NewCDNNode("edge-1", "global", "Close", 500))
	cdn.AddNode(NewCDNNode("edge-2", "global", "Medium", 500))
	cdn.AddNode(NewCDNNode("edge-3", "global", "Far", 500))

	// Set up latency matrix (user -> node latencies in ms)
	latencyRouting.SetLatency("user-location", "edge-1", 50)  // Closest
	latencyRouting.SetLatency("user-location", "edge-2", 150) // Medium
	latencyRouting.SetLatency("user-location", "edge-3", 300) // Farthest

	fmt.Println("Latency routing will prefer edge-1 (50ms) over others")

	// Test routing - should prefer edge-1
	for i := 0; i < 5; i++ {
		data, err := cdn.Get("test-content", "user-location")
		if err != nil {
			fmt.Printf("Request %d: ERROR - %v\n", i+1, err)
		} else {
			fmt.Printf("Request %d: SUCCESS - %s\n", i+1,
				truncateString(fmt.Sprintf("%v", data), 30))
		}
	}
}

func demoCDNInvalidation() {
	fmt.Println("\n🎯 CDN Invalidation Demo:")
	fmt.Println(strings.Repeat("=", 30))

	routing := &GeographicRouting{}
	cdn := NewCDNNetwork(routing)

	// Add nodes
	cdn.AddNode(NewCDNNode("cache-1", "region-1", "Location 1", 100))
	cdn.AddNode(NewCDNNode("cache-2", "region-1", "Location 2", 100))

	// Cache some content
	cdn.Get("style.css", "region-1")
	cdn.Get("script.js", "region-1")

	fmt.Println("Cached content in CDN nodes")

	// Set up invalidation system
	invalidation := NewCDNInvalidation(cdn)

	// Invalidate specific file
	invalidation.InvalidatePattern("style.css")
	fmt.Println("Invalidated style.css across all nodes")

	// Invalidate all files matching pattern
	invalidation.InvalidatePattern("*.js")
	fmt.Println("Invalidated all JS files across all nodes")

	history := invalidation.GetInvalidationHistory()
	fmt.Printf("Invalidation history: %v\n", history)
}

func demoCDNAnalytics() {
	fmt.Println("\n📊 CDN Analytics Demo:")
	fmt.Println(strings.Repeat("=", 25))

	routing := &LoadBasedRouting{}
	cdn := NewCDNNetwork(routing)

	// Add multiple nodes
	nodes := []string{"analytics-1", "analytics-2", "analytics-3"}
	for _, nodeID := range nodes {
		cdn.AddNode(NewCDNNode(nodeID, "global", "Test Location", 200))
	}

	// Simulate traffic
	fmt.Println("Simulating CDN traffic...")
	for i := 0; i < 50; i++ {
		userLocation := fmt.Sprintf("user-%d", i%3)
		content := fmt.Sprintf("content-%d.jpg", i%10)
		cdn.Get(content, userLocation)
	}

	// Collect and display analytics
	analytics := NewCDNAnalytics(cdn)
	metrics := analytics.CollectMetrics()

	fmt.Println("\nCDN Analytics Results:")
	for key, value := range metrics {
		if key != "nodes_stats" && key != "collection_time" {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	// Show top performing nodes
	topNodes := analytics.GetTopPerformingNodes(2)
	fmt.Printf("\nTop Performing Nodes: %v\n", topNodes)
}

func demoCDNContentTypes() {
	fmt.Println("\n🎵 CDN Content Type Optimization:")
	fmt.Println(strings.Repeat("=", 40))

	routing := &GeographicRouting{}
	cdn := NewCDNNetwork(routing)
	cdn.AddNode(NewCDNNode("media-cdn", "global", "Media Server", 1000))

	contentTypes := map[string]time.Duration{
		"image.jpg":     24 * time.Hour,     // Images - long cache
		"video.mp4":     7 * 24 * time.Hour, // Videos - very long cache
		"style.css":     12 * time.Hour,     // CSS - medium cache
		"script.js":     6 * time.Hour,      // JS - shorter cache
		"api/data.json": 1 * time.Hour,      // API data - short cache
	}

	fmt.Println("Content-specific caching strategies:")
	for contentType, ttl := range contentTypes {
		data, _ := cdn.Get(contentType, "global")
		fmt.Printf("  %s: TTL %v - %s\n",
			contentType, ttl, truncateString(fmt.Sprintf("%v", data), 40))
	}

	fmt.Println("\nCDN optimizes cache duration based on content type:")
	fmt.Println("- Static assets (images, videos): Long cache duration")
	fmt.Println("- Dynamic content (API data): Short cache duration")
	fmt.Println("- CSS/JS: Medium cache duration with versioning")
}

func generateHash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func runCDNDemo() {
	fmt.Println("🎯 FAANG System Design Interview - CDN Mastery")
	fmt.Println("Essential content delivery patterns for L3/L4 interviews")
	fmt.Println()

	// Run all demos
	demoCDNBasicOperations()
	demoCDNLatencyRouting()
	demoCDNInvalidation()
	demoCDNAnalytics()
	demoCDNContentTypes()

	fmt.Println("\n🎪 Key Takeaways for FAANG Interviews:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. Geographic Routing: Route users to nearest edge server")
	fmt.Println("2. Latency-Based: Dynamic routing based on real-time latency")
	fmt.Println("3. Cache Invalidation: Propagate updates across CDN network")
	fmt.Println("4. Content Optimization: Different TTL for different content types")
	fmt.Println("5. Analytics: Monitor hit ratios and performance metrics")
	fmt.Println("6. Load Balancing: Distribute traffic across edge servers")
	fmt.Println("\n💪 Practice designing YouTube/Netflix CDN architecture!")
}

// Uncomment to run this demo independently
// func main() {
//     rand.Seed(time.Now().UnixNano())
//     runCDNDemo()
// }
