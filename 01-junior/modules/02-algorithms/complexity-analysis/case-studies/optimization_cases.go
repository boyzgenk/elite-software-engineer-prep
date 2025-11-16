// Package casestudies demonstrates real-world complexity optimization examples
// Production optimization cases from FAANG-level engineering
package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// 🎯 CASE STUDY 1: SOCIAL MEDIA FEED OPTIMIZATION
// From O(n²) to O(n log k) - Real Instagram/Twitter scale optimization

// SocialMediaPost represents a post in social media feed
type SocialMediaPost struct {
	ID        int64
	UserID    int64
	Timestamp time.Time
	Score     float64 // Engagement score
	Content   string
}

// NaiveFeedGeneration - O(n²) implementation (BEFORE optimization)
func NaiveFeedGeneration(allPosts []SocialMediaPost, followedUsers []int64, limit int) []SocialMediaPost {
	// O(n²) approach: For each post, check if user is followed
	var relevantPosts []SocialMediaPost

	for _, post := range allPosts { // O(n)
		for _, userID := range followedUsers { // O(m) where m = followed users
			if post.UserID == userID {
				relevantPosts = append(relevantPosts, post)
				break
			}
		}
	}

	// Sort by engagement score - O(k log k) where k = relevant posts
	sort.Slice(relevantPosts, func(i, j int) bool {
		return relevantPosts[i].Score > relevantPosts[j].Score
	})

	// Return top posts
	if len(relevantPosts) > limit {
		return relevantPosts[:limit]
	}
	return relevantPosts
}

// OptimizedFeedGeneration - O(n + k log k) implementation (AFTER optimization)
func OptimizedFeedGeneration(allPosts []SocialMediaPost, followedUsers []int64, limit int) []SocialMediaPost {
	// O(m) preprocessing: Create hash set of followed users
	followedSet := make(map[int64]bool)
	for _, userID := range followedUsers {
		followedSet[userID] = true
	}

	// O(n) filtering: Single pass through all posts
	relevantPosts := make([]SocialMediaPost, 0, len(allPosts)/10) // Estimate capacity
	for _, post := range allPosts {
		if followedSet[post.UserID] { // O(1) lookup
			relevantPosts = append(relevantPosts, post)
		}
	}

	// O(k log k) sorting where k << n
	sort.Slice(relevantPosts, func(i, j int) bool {
		return relevantPosts[i].Score > relevantPosts[j].Score
	})

	if len(relevantPosts) > limit {
		return relevantPosts[:limit]
	}
	return relevantPosts
}

// BenchmarkFeedGeneration compares naive vs optimized approaches
func BenchmarkFeedGeneration() {
	fmt.Println("📱 CASE STUDY 1: SOCIAL MEDIA FEED OPTIMIZATION")
	fmt.Println(strings.Repeat("=", 70))

	// Generate test data
	allPosts := make([]SocialMediaPost, 100000)
	followedUsers := make([]int64, 1000)

	for i := range allPosts {
		allPosts[i] = SocialMediaPost{
			ID:        int64(i),
			UserID:    int64(i % 10000), // 10k unique users
			Timestamp: time.Now().Add(-time.Duration(i) * time.Minute),
			Score:     float64(i%100) + 0.1*float64(i%10),
		}
	}

	for i := range followedUsers {
		followedUsers[i] = int64(i * 10) // Every 10th user is followed
	}

	// Benchmark naive approach
	start := time.Now()
	naiveResult := NaiveFeedGeneration(allPosts, followedUsers, 50)
	naiveDuration := time.Since(start)

	// Benchmark optimized approach
	start = time.Now()
	optimizedResult := OptimizedFeedGeneration(allPosts, followedUsers, 50)
	optimizedDuration := time.Since(start)

	fmt.Printf("Naive O(n²) approach:      %v (%d posts)\n", naiveDuration, len(naiveResult))
	fmt.Printf("Optimized O(n log k):      %v (%d posts)\n", optimizedDuration, len(optimizedResult))
	fmt.Printf("Speedup:                   %.2fx faster\n", float64(naiveDuration)/float64(optimizedDuration))
	fmt.Printf("Complexity improvement:    O(n²) → O(n + k log k)\n")
}

// 🎯 CASE STUDY 2: AUTOCOMPLETE OPTIMIZATION
// From O(n*m) to O(m) - Search optimization using Trie

// AutocompleteNaive - O(n*m) implementation scanning all strings
func AutocompleteNaive(dictionary []string, prefix string) []string {
	var matches []string

	// O(n*m): For each string, check if it starts with prefix
	for _, word := range dictionary { // O(n)
		if len(word) >= len(prefix) && word[:len(prefix)] == prefix { // O(m)
			matches = append(matches, word)
		}
	}

	return matches
}

// TrieNode for optimized autocomplete
type TrieNode struct {
	Children map[rune]*TrieNode
	IsEnd    bool
	Words    []string // Store words ending at this node
}

// Trie structure for O(m) autocomplete
type Trie struct {
	Root *TrieNode
}

// NewTrie creates a new trie structure
func NewTrie() *Trie {
	return &Trie{
		Root: &TrieNode{
			Children: make(map[rune]*TrieNode),
			Words:    make([]string, 0),
		},
	}
}

// Insert adds a word to the trie - O(m) where m = word length
func (t *Trie) Insert(word string) {
	current := t.Root

	for _, char := range word {
		if _, exists := current.Children[char]; !exists {
			current.Children[char] = &TrieNode{
				Children: make(map[rune]*TrieNode),
				Words:    make([]string, 0),
			}
		}
		current = current.Children[char]
		current.Words = append(current.Words, word) // Store at each level
	}
	current.IsEnd = true
}

// Search finds all words with given prefix - O(m) where m = prefix length
func (t *Trie) Search(prefix string) []string {
	current := t.Root

	// Navigate to prefix node - O(m)
	for _, char := range prefix {
		if _, exists := current.Children[char]; !exists {
			return []string{} // No matches
		}
		current = current.Children[char]
	}

	// Return all words at this node - O(1) retrieval
	return current.Words
}

// AutocompleteOptimized - O(m) implementation using Trie
func AutocompleteOptimized(trie *Trie, prefix string) []string {
	return trie.Search(prefix) // O(m) where m = prefix length
}

// BenchmarkAutocomplete compares naive vs optimized autocomplete
func BenchmarkAutocomplete() {
	fmt.Println("\n🔍 CASE STUDY 2: AUTOCOMPLETE OPTIMIZATION")
	fmt.Println(strings.Repeat("=", 70))

	// Generate dictionary
	dictionary := []string{
		"algorithm", "algorithmic", "algorithms", "algebra", "alphabetical",
		"autocomplete", "automatic", "automobile", "autonomy", "autumn",
		"application", "apple", "apply", "approach", "appropriate",
		"architecture", "archive", "argument", "arithmetic", "arrangement",
		"basketball", "baseball", "badminton", "bicycle", "biology",
		"computer", "computing", "complexity", "complete", "complicated",
		"data", "database", "datastore", "development", "developer",
	}

	// Build trie (preprocessing cost)
	trie := NewTrie()
	for _, word := range dictionary {
		trie.Insert(word)
	}

	testPrefix := "algo"

	// Benchmark naive approach
	start := time.Now()
	var naiveResults []string
	for i := 0; i < 10000; i++ { // Simulate multiple searches
		naiveResults = AutocompleteNaive(dictionary, testPrefix)
	}
	naiveDuration := time.Since(start)

	// Benchmark optimized approach
	start = time.Now()
	var optimizedResults []string
	for i := 0; i < 10000; i++ { // Simulate multiple searches
		optimizedResults = AutocompleteOptimized(trie, testPrefix)
	}
	optimizedDuration := time.Since(start)

	fmt.Printf("Dictionary size:           %d words\n", len(dictionary))
	fmt.Printf("Search prefix:             '%s'\n", testPrefix)
	fmt.Printf("Naive O(n*m):             %v (%d matches)\n", naiveDuration, len(naiveResults))
	fmt.Printf("Optimized O(m):           %v (%d matches)\n", optimizedDuration, len(optimizedResults))
	fmt.Printf("Speedup:                  %.2fx faster\n", float64(naiveDuration)/float64(optimizedDuration))
	fmt.Printf("Complexity improvement:   O(n*m) → O(m)\n")
}

// 🎯 CASE STUDY 3: DATABASE QUERY OPTIMIZATION
// From O(n²) to O(n log n) - Join optimization

// UserRecord represents a user in the database
type UserRecord struct {
	ID   int64
	Name string
	Age  int
}

// OrderRecord represents an order in the database
type OrderRecord struct {
	ID     int64
	UserID int64
	Amount float64
}

// NaiveJoin - O(n*m) nested loop join
func NaiveJoin(users []UserRecord, orders []OrderRecord) []struct {
	User  UserRecord
	Order OrderRecord
} {
	var results []struct {
		User  UserRecord
		Order OrderRecord
	}

	// Nested loop join - O(n*m)
	for _, user := range users { // O(n)
		for _, order := range orders { // O(m)
			if user.ID == order.UserID {
				results = append(results, struct {
					User  UserRecord
					Order OrderRecord
				}{User: user, Order: order})
			}
		}
	}

	return results
}

// OptimizedJoin - O(n + m) hash join
func OptimizedJoin(users []UserRecord, orders []OrderRecord) []struct {
	User  UserRecord
	Order OrderRecord
} {
	// Build hash table of users - O(n)
	userMap := make(map[int64]UserRecord)
	for _, user := range users {
		userMap[user.ID] = user
	}

	var results []struct {
		User  UserRecord
		Order OrderRecord
	}

	// Hash join - O(m)
	for _, order := range orders {
		if user, exists := userMap[order.UserID]; exists { // O(1) lookup
			results = append(results, struct {
				User  UserRecord
				Order OrderRecord
			}{User: user, Order: order})
		}
	}

	return results
}

// BenchmarkJoin compares naive vs optimized database joins
func BenchmarkJoin() {
	fmt.Println("\n💾 CASE STUDY 3: DATABASE JOIN OPTIMIZATION")
	fmt.Println(strings.Repeat("=", 70))

	// Generate test data
	users := make([]UserRecord, 10000)
	orders := make([]OrderRecord, 50000)

	for i := range users {
		users[i] = UserRecord{
			ID:   int64(i + 1),
			Name: fmt.Sprintf("User_%d", i+1),
			Age:  20 + (i % 60),
		}
	}

	for i := range orders {
		orders[i] = OrderRecord{
			ID:     int64(i + 1),
			UserID: int64((i % len(users)) + 1), // Random user assignment
			Amount: float64(10 + (i % 1000)),
		}
	}

	// Benchmark naive join
	start := time.Now()
	naiveResults := NaiveJoin(users, orders)
	naiveDuration := time.Since(start)

	// Benchmark optimized join
	start = time.Now()
	optimizedResults := OptimizedJoin(users, orders)
	optimizedDuration := time.Since(start)

	fmt.Printf("Users:                     %d records\n", len(users))
	fmt.Printf("Orders:                    %d records\n", len(orders))
	fmt.Printf("Naive O(n*m):             %v (%d results)\n", naiveDuration, len(naiveResults))
	fmt.Printf("Optimized O(n+m):         %v (%d results)\n", optimizedDuration, len(optimizedResults))
	fmt.Printf("Speedup:                  %.2fx faster\n", float64(naiveDuration)/float64(optimizedDuration))
	fmt.Printf("Complexity improvement:   O(n*m) → O(n+m)\n")
}

// 🎯 CASE STUDY 4: SPACE-TIME TRADEOFF ANALYSIS
// Fibonacci: Recursive vs Memoized vs Tabulated

// FibonacciRecursive - O(2^n) time, O(n) space (call stack)
func FibonacciRecursive(n int) int64 {
	if n <= 1 {
		return int64(n)
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

// FibonacciMemoized - O(n) time, O(n) space
var memo = make(map[int]int64)

func FibonacciMemoized(n int) int64 {
	if n <= 1 {
		return int64(n)
	}

	if val, exists := memo[n]; exists {
		return val
	}

	memo[n] = FibonacciMemoized(n-1) + FibonacciMemoized(n-2)
	return memo[n]
}

// FibonacciTabulated - O(n) time, O(n) space
func FibonacciTabulated(n int) int64 {
	if n <= 1 {
		return int64(n)
	}

	dp := make([]int64, n+1)
	dp[0], dp[1] = 0, 1

	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n]
}

// FibonacciOptimal - O(n) time, O(1) space
func FibonacciOptimal(n int) int64 {
	if n <= 1 {
		return int64(n)
	}

	prev2, prev1 := int64(0), int64(1)

	for i := 2; i <= n; i++ {
		current := prev1 + prev2
		prev2, prev1 = prev1, current
	}

	return prev1
}

// BenchmarkFibonacci compares different Fibonacci implementations
func BenchmarkFibonacci() {
	fmt.Println("\n🔢 CASE STUDY 4: SPACE-TIME TRADEOFF - FIBONACCI")
	fmt.Println(strings.Repeat("=", 70))

	testValues := []int{10, 20, 30, 40}

	for _, n := range testValues {
		fmt.Printf("\nFibonacci(%d):\n", n)

		// Clear memoization
		memo = make(map[int]int64)

		// Recursive (only for small values)
		if n <= 35 {
			start := time.Now()
			resultRecursive := FibonacciRecursive(n)
			recursiveDuration := time.Since(start)
			fmt.Printf("  Recursive O(2^n):       %v (result: %d)\n", recursiveDuration, resultRecursive)
		} else {
			fmt.Printf("  Recursive O(2^n):       [SKIPPED - too slow]\n")
		}

		// Memoized
		start := time.Now()
		resultMemoized := FibonacciMemoized(n)
		memoizedDuration := time.Since(start)
		fmt.Printf("  Memoized O(n):          %v (result: %d)\n", memoizedDuration, resultMemoized)

		// Tabulated
		start = time.Now()
		resultTabulated := FibonacciTabulated(n)
		tabulatedDuration := time.Since(start)
		fmt.Printf("  Tabulated O(n):         %v (result: %d)\n", tabulatedDuration, resultTabulated)

		// Optimal
		start = time.Now()
		resultOptimal := FibonacciOptimal(n)
		optimalDuration := time.Since(start)
		fmt.Printf("  Optimal O(1) space:     %v (result: %d)\n", optimalDuration, resultOptimal)
	}

	fmt.Println("\n💡 Key Insight: Space-time tradeoffs allow us to optimize for")
	fmt.Println("   different constraints (memory vs computation vs development time)")
}

// 🎯 MAIN CASE STUDIES DEMONSTRATION

func main() {
	fmt.Println("🎯 REAL-WORLD COMPLEXITY OPTIMIZATION CASE STUDIES")
	fmt.Println("Production examples from FAANG-scale engineering")
	fmt.Println(strings.Repeat("=", 80))

	// Run all case studies
	BenchmarkFeedGeneration()
	BenchmarkAutocomplete()
	BenchmarkJoin()
	BenchmarkFibonacci()

	// Summary
	fmt.Println("\n📊 OPTIMIZATION SUMMARY")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("✅ Social Media Feed:     O(n²) → O(n log k)  [Hash set preprocessing]")
	fmt.Println("✅ Autocomplete:          O(n*m) → O(m)      [Trie data structure]")
	fmt.Println("✅ Database Join:         O(n*m) → O(n+m)    [Hash join algorithm]")
	fmt.Println("✅ Fibonacci:             O(2^n) → O(n)      [Dynamic programming]")
	fmt.Println("\n💡 These optimizations are used daily at Google, Meta, Netflix, Amazon!")
}
