package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 🎯 FAANG System Design - Load Balancer Implementations
// Essential patterns for L3/L4 interviews: Round-Robin, Weighted, Consistent Hashing

// Server represents a backend server with health checking capabilities
type Server struct {
	URL       string
	Weight    int   // For weighted load balancing
	IsHealthy int32 // Atomic flag for health status
	Mutex     sync.RWMutex
	Requests  int64 // Request counter for monitoring
}

// IsServerHealthy checks if server is healthy (atomic read)
func (s *Server) IsServerHealthy() bool {
	return atomic.LoadInt32(&s.IsHealthy) == 1
}

// SetHealth sets server health status (atomic write)
func (s *Server) SetHealth(healthy bool) {
	if healthy {
		atomic.StoreInt32(&s.IsHealthy, 1)
	} else {
		atomic.StoreInt32(&s.IsHealthy, 0)
	}
}

// IncrementRequests atomically increments request counter
func (s *Server) IncrementRequests() {
	atomic.AddInt64(&s.Requests, 1)
}

// GetRequestCount returns current request count
func (s *Server) GetRequestCount() int64 {
	return atomic.LoadInt64(&s.Requests)
}

// LoadBalancer interface defines load balancing strategies
type LoadBalancer interface {
	NextServer() *Server
	AddServer(server *Server)
	RemoveServer(url string)
	GetServers() []*Server
	GetStats() map[string]interface{}
}

// 🔄 Round-Robin Load Balancer
// Most common pattern in FAANG interviews - simple but effective
type RoundRobinLB struct {
	servers []*Server
	current int32
	mutex   sync.RWMutex
}

func NewRoundRobinLB() *RoundRobinLB {
	return &RoundRobinLB{
		servers: make([]*Server, 0),
		current: -1,
	}
}

func (rr *RoundRobinLB) NextServer() *Server {
	rr.mutex.RLock()
	defer rr.mutex.RUnlock()

	if len(rr.servers) == 0 {
		return nil
	}

	// Find next healthy server using round-robin
	healthyServers := rr.getHealthyServers()
	if len(healthyServers) == 0 {
		return nil
	}

	// Atomic increment with wraparound
	next := atomic.AddInt32(&rr.current, 1)
	server := healthyServers[int(next)%len(healthyServers)]
	server.IncrementRequests()

	return server
}

func (rr *RoundRobinLB) getHealthyServers() []*Server {
	healthy := make([]*Server, 0)
	for _, server := range rr.servers {
		if server.IsServerHealthy() {
			healthy = append(healthy, server)
		}
	}
	return healthy
}

func (rr *RoundRobinLB) AddServer(server *Server) {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()
	server.SetHealth(true) // New servers start healthy
	rr.servers = append(rr.servers, server)
}

func (rr *RoundRobinLB) RemoveServer(url string) {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	for i, server := range rr.servers {
		if server.URL == url {
			rr.servers = append(rr.servers[:i], rr.servers[i+1:]...)
			break
		}
	}
}

func (rr *RoundRobinLB) GetServers() []*Server {
	rr.mutex.RLock()
	defer rr.mutex.RUnlock()
	return rr.servers
}

func (rr *RoundRobinLB) GetStats() map[string]interface{} {
	rr.mutex.RLock()
	defer rr.mutex.RUnlock()

	totalRequests := int64(0)
	healthyCount := 0

	for _, server := range rr.servers {
		totalRequests += server.GetRequestCount()
		if server.IsServerHealthy() {
			healthyCount++
		}
	}

	return map[string]interface{}{
		"type":            "round-robin",
		"total_servers":   len(rr.servers),
		"healthy_servers": healthyCount,
		"total_requests":  totalRequests,
	}
}

// ⚖️ Weighted Load Balancer
// Critical for FAANG interviews - handles heterogeneous server capacities
type WeightedLB struct {
	servers     []*Server
	totalWeight int
	mutex       sync.RWMutex
}

func NewWeightedLB() *WeightedLB {
	return &WeightedLB{
		servers: make([]*Server, 0),
	}
}

func (w *WeightedLB) NextServer() *Server {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	if len(w.servers) == 0 {
		return nil
	}

	// Calculate total weight of healthy servers
	healthyServers := make([]*Server, 0)
	totalWeight := 0

	for _, server := range w.servers {
		if server.IsServerHealthy() {
			healthyServers = append(healthyServers, server)
			totalWeight += server.Weight
		}
	}

	if len(healthyServers) == 0 || totalWeight == 0 {
		return nil
	}

	// Weighted random selection
	randWeight := rand.Intn(totalWeight)
	currentWeight := 0

	for _, server := range healthyServers {
		currentWeight += server.Weight
		if randWeight < currentWeight {
			server.IncrementRequests()
			return server
		}
	}

	// Fallback (should never reach here)
	return healthyServers[0]
}

func (w *WeightedLB) AddServer(server *Server) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if server.Weight <= 0 {
		server.Weight = 1 // Default weight
	}

	server.SetHealth(true)
	w.servers = append(w.servers, server)
	w.totalWeight += server.Weight
}

func (w *WeightedLB) RemoveServer(url string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for i, server := range w.servers {
		if server.URL == url {
			w.totalWeight -= server.Weight
			w.servers = append(w.servers[:i], w.servers[i+1:]...)
			break
		}
	}
}

func (w *WeightedLB) GetServers() []*Server {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.servers
}

func (w *WeightedLB) GetStats() map[string]interface{} {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	totalRequests := int64(0)
	healthyCount := 0
	totalWeight := 0

	for _, server := range w.servers {
		totalRequests += server.GetRequestCount()
		if server.IsServerHealthy() {
			healthyCount++
			totalWeight += server.Weight
		}
	}

	return map[string]interface{}{
		"type":            "weighted",
		"total_servers":   len(w.servers),
		"healthy_servers": healthyCount,
		"total_requests":  totalRequests,
		"total_weight":    totalWeight,
	}
}

// 🔄 Consistent Hashing Load Balancer
// Advanced pattern for FAANG L4+ interviews - handles server addition/removal gracefully
type ConsistentHashLB struct {
	servers      map[string]*Server
	ring         map[uint32]string // Hash ring: hash -> server URL
	sortedHashes []uint32          // Sorted hash keys for binary search
	replicas     int               // Virtual nodes per server
	mutex        sync.RWMutex
}

func NewConsistentHashLB(replicas int) *ConsistentHashLB {
	return &ConsistentHashLB{
		servers:      make(map[string]*Server),
		ring:         make(map[uint32]string),
		sortedHashes: make([]uint32, 0),
		replicas:     replicas,
	}
}

func (c *ConsistentHashLB) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (c *ConsistentHashLB) NextServer() *Server {
	// For demonstration, use current timestamp as key
	// In real scenarios, this would be client IP or request ID
	key := fmt.Sprintf("request_%d", time.Now().UnixNano())
	return c.GetServerForKey(key)
}

func (c *ConsistentHashLB) GetServerForKey(key string) *Server {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.sortedHashes) == 0 {
		return nil
	}

	hash := c.hash(key)

	// Find the first server with hash >= target hash
	idx := c.binarySearch(hash)
	serverURL := c.ring[c.sortedHashes[idx]]

	server := c.servers[serverURL]
	if server != nil && server.IsServerHealthy() {
		server.IncrementRequests()
		return server
	}

	// If selected server is unhealthy, try next healthy server
	for i := 1; i < len(c.sortedHashes); i++ {
		nextIdx := (idx + i) % len(c.sortedHashes)
		nextServerURL := c.ring[c.sortedHashes[nextIdx]]
		nextServer := c.servers[nextServerURL]

		if nextServer != nil && nextServer.IsServerHealthy() {
			nextServer.IncrementRequests()
			return nextServer
		}
	}

	return nil // No healthy servers available
}

func (c *ConsistentHashLB) binarySearch(hash uint32) int {
	// Binary search for first hash >= target
	left, right := 0, len(c.sortedHashes)

	for left < right {
		mid := (left + right) / 2
		if c.sortedHashes[mid] >= hash {
			right = mid
		} else {
			left = mid + 1
		}
	}

	// Wrap around if we've gone past the end
	if left >= len(c.sortedHashes) {
		left = 0
	}

	return left
}

func (c *ConsistentHashLB) AddServer(server *Server) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	server.SetHealth(true)
	c.servers[server.URL] = server

	// Add virtual nodes to ring
	for i := 0; i < c.replicas; i++ {
		virtualKey := fmt.Sprintf("%s:%d", server.URL, i)
		hash := c.hash(virtualKey)
		c.ring[hash] = server.URL
		c.sortedHashes = append(c.sortedHashes, hash)
	}

	// Re-sort the hash ring
	c.sortHashes()
}

func (c *ConsistentHashLB) RemoveServer(url string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.servers, url)

	// Remove virtual nodes from ring
	for i := 0; i < c.replicas; i++ {
		virtualKey := fmt.Sprintf("%s:%d", url, i)
		hash := c.hash(virtualKey)
		delete(c.ring, hash)

		// Remove from sorted hashes
		for j, h := range c.sortedHashes {
			if h == hash {
				c.sortedHashes = append(c.sortedHashes[:j], c.sortedHashes[j+1:]...)
				break
			}
		}
	}
}

func (c *ConsistentHashLB) sortHashes() {
	// Simple selection sort (good enough for demo)
	n := len(c.sortedHashes)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if c.sortedHashes[j] < c.sortedHashes[minIdx] {
				minIdx = j
			}
		}
		c.sortedHashes[i], c.sortedHashes[minIdx] = c.sortedHashes[minIdx], c.sortedHashes[i]
	}
}

func (c *ConsistentHashLB) GetServers() []*Server {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	servers := make([]*Server, 0, len(c.servers))
	for _, server := range c.servers {
		servers = append(servers, server)
	}
	return servers
}

func (c *ConsistentHashLB) GetStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	totalRequests := int64(0)
	healthyCount := 0

	for _, server := range c.servers {
		totalRequests += server.GetRequestCount()
		if server.IsServerHealthy() {
			healthyCount++
		}
	}

	return map[string]interface{}{
		"type":            "consistent-hash",
		"total_servers":   len(c.servers),
		"healthy_servers": healthyCount,
		"total_requests":  totalRequests,
		"replicas":        c.replicas,
		"ring_size":       len(c.ring),
	}
}

// 🏥 Health Checker - Essential for production load balancers
type HealthChecker struct {
	loadBalancer LoadBalancer
	interval     time.Duration
	timeout      time.Duration
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewHealthChecker(lb LoadBalancer, interval, timeout time.Duration) *HealthChecker {
	ctx, cancel := context.WithCancel(context.Background())
	return &HealthChecker{
		loadBalancer: lb,
		interval:     interval,
		timeout:      timeout,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	log.Printf("Health checker started with %v interval", hc.interval)

	for {
		select {
		case <-hc.ctx.Done():
			log.Println("Health checker stopped")
			return
		case <-ticker.C:
			hc.checkHealth()
		}
	}
}

func (hc *HealthChecker) Stop() {
	hc.cancel()
}

func (hc *HealthChecker) checkHealth() {
	servers := hc.loadBalancer.GetServers()

	for _, server := range servers {
		go hc.checkServerHealth(server)
	}
}

func (hc *HealthChecker) checkServerHealth(server *Server) {
	client := &http.Client{
		Timeout: hc.timeout,
	}

	// Simple health check - GET request to server
	resp, err := client.Get(server.URL + "/health")

	if err != nil || resp.StatusCode != 200 {
		if server.IsServerHealthy() {
			log.Printf("Server %s marked as unhealthy", server.URL)
			server.SetHealth(false)
		}
	} else {
		if !server.IsServerHealthy() {
			log.Printf("Server %s recovered and marked as healthy", server.URL)
			server.SetHealth(true)
		}
		resp.Body.Close()
	}
}

// 🧪 Demo Functions for FAANG Interview Practice
func demoLoadBalancers() {
	fmt.Println("🚀 FAANG System Design - Load Balancer Implementations Demo")
	fmt.Println(strings.Repeat("=", 60))

	// Create test servers
	servers := []*Server{
		{URL: "http://server1:8001", Weight: 3},
		{URL: "http://server2:8002", Weight: 2},
		{URL: "http://server3:8003", Weight: 1},
	}

	// Demo Round-Robin Load Balancer
	fmt.Println("\n🔄 Round-Robin Load Balancer Demo:")
	rrLB := NewRoundRobinLB()
	for _, server := range servers {
		rrLB.AddServer(server)
	}

	for i := 0; i < 6; i++ {
		server := rrLB.NextServer()
		fmt.Printf("Request %d -> %s\n", i+1, server.URL)
	}
	fmt.Printf("Stats: %+v\n", rrLB.GetStats())

	// Demo Weighted Load Balancer
	fmt.Println("\n⚖️ Weighted Load Balancer Demo:")
	wLB := NewWeightedLB()
	for _, server := range servers {
		wLB.AddServer(server)
	}

	requestCounts := make(map[string]int)
	for i := 0; i < 100; i++ {
		server := wLB.NextServer()
		requestCounts[server.URL]++
	}

	fmt.Println("Request distribution (100 requests):")
	for url, count := range requestCounts {
		fmt.Printf("%s: %d requests\n", url, count)
	}
	fmt.Printf("Stats: %+v\n", wLB.GetStats())

	// Demo Consistent Hashing Load Balancer
	fmt.Println("\n🔄 Consistent Hashing Load Balancer Demo:")
	chLB := NewConsistentHashLB(3) // 3 virtual nodes per server
	for _, server := range servers {
		chLB.AddServer(server)
	}

	// Test key distribution
	keys := []string{"user123", "user456", "user789", "user000", "user111"}
	fmt.Println("Key-to-server mapping:")
	for _, key := range keys {
		server := chLB.GetServerForKey(key)
		fmt.Printf("Key '%s' -> %s\n", key, server.URL)
	}
	fmt.Printf("Stats: %+v\n", chLB.GetStats())

	// Demonstrate server removal (key feature of consistent hashing)
	fmt.Println("\nRemoving server2 - demonstrating minimal key remapping:")
	chLB.RemoveServer("http://server2:8002")
	for _, key := range keys {
		server := chLB.GetServerForKey(key)
		fmt.Printf("Key '%s' -> %s (after removal)\n", key, server.URL)
	}
}

// 📊 Performance Comparison for FAANG Interview Discussions
func performanceComparison() {
	fmt.Println("\n📊 Performance Comparison for Interview Discussion:")
	fmt.Println(strings.Repeat("=", 60))

	servers := []*Server{
		{URL: "http://server1:8001", Weight: 1},
		{URL: "http://server2:8002", Weight: 1},
		{URL: "http://server3:8003", Weight: 1},
	}

	// Benchmark Round-Robin
	rrLB := NewRoundRobinLB()
	for _, server := range servers {
		rrLB.AddServer(server)
	}

	start := time.Now()
	for i := 0; i < 100000; i++ {
		rrLB.NextServer()
	}
	rrTime := time.Since(start)

	// Benchmark Weighted
	wLB := NewWeightedLB()
	for _, server := range servers {
		wLB.AddServer(server)
	}

	start = time.Now()
	for i := 0; i < 100000; i++ {
		wLB.NextServer()
	}
	wTime := time.Since(start)

	// Benchmark Consistent Hashing
	chLB := NewConsistentHashLB(3)
	for _, server := range servers {
		chLB.AddServer(server)
	}

	start = time.Now()
	for i := 0; i < 100000; i++ {
		chLB.NextServer()
	}
	chTime := time.Since(start)

	fmt.Printf("Round-Robin: %v (100k requests)\n", rrTime)
	fmt.Printf("Weighted: %v (100k requests)\n", wTime)
	fmt.Printf("Consistent Hash: %v (100k requests)\n", chTime)

	fmt.Println("\n💡 Interview Discussion Points:")
	fmt.Println("- Round-Robin: O(1) time, simple but no stickiness")
	fmt.Println("- Weighted: O(1) time, handles heterogeneous servers")
	fmt.Println("- Consistent Hash: O(log n) time, minimal redistribution on changes")
}

// 🎯 Health Check Demo with Mock Servers
func demoHealthChecking() {
	fmt.Println("\n🏥 Health Checking Demo:")
	fmt.Println(strings.Repeat("=", 40))

	// Create mock HTTP servers for health checking
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		}
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(503) // Simulate unhealthy server
			w.Write([]byte("Service Unavailable"))
		}
	}))
	defer server2.Close()

	// Create load balancer with health checking
	lb := NewRoundRobinLB()
	lb.AddServer(&Server{URL: server1.URL})
	lb.AddServer(&Server{URL: server2.URL})

	// Start health checker
	healthChecker := NewHealthChecker(lb, 2*time.Second, 1*time.Second)
	go healthChecker.Start()

	// Wait for health checks to run
	time.Sleep(3 * time.Second)

	fmt.Println("Testing load balancer with health checking:")
	for i := 0; i < 5; i++ {
		server := lb.NextServer()
		if server != nil {
			fmt.Printf("Request %d -> %s (healthy)\n", i+1, server.URL)
		} else {
			fmt.Printf("Request %d -> no healthy servers available\n", i+1)
		}
	}

	healthChecker.Stop()
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("🎯 FAANG System Design Interview - Load Balancer Mastery")
	fmt.Println("Essential patterns for L3/L4 technical interviews")
	fmt.Println()

	// Run all demos
	demoLoadBalancers()
	performanceComparison()
	demoHealthChecking()

	fmt.Println("\n🎪 Key Takeaways for FAANG Interviews:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. Round-Robin: Simple, fair distribution, O(1)")
	fmt.Println("2. Weighted: Handles different server capacities")
	fmt.Println("3. Consistent Hashing: Minimal redistribution, cache-friendly")
	fmt.Println("4. Health Checking: Essential for production reliability")
	fmt.Println("5. Atomic Operations: Thread-safe implementations")
	fmt.Println("\n💪 Practice designing Twitter/Netflix with these patterns!")
}
