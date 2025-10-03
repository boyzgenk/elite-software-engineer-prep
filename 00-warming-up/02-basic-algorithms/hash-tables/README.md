# 🗂️ Hash Tables (Maps) Fundamentals
## Week 2, Days 1-3 | Key-Value Data Structure Mastery

### 🎯 Learning Objectives
By the end of this section, you will:
- [ ] Understand hash function concepts and collision resolution
- [ ] Implement hash table using separate chaining and open addressing
- [ ] Master common hash table patterns: frequency counting, caching
- [ ] Solve 10-12 hash table problems efficiently
- [ ] Analyze time complexity trade-offs in different scenarios

---

## 📚 Hash Table Fundamentals

### 🔸 What is a Hash Table?

**Definition:**
- Data structure that implements key-value pairs using hash functions
- Provides average O(1) time complexity for insertion, deletion, and lookup
- Uses hash function to compute index where element should be stored
- Also known as HashMap, Dictionary, or Associative Array

**Key Components:**
- **Hash Function:** Converts key to array index
- **Bucket Array:** Stores the actual key-value pairs
- **Collision Resolution:** Handles when multiple keys hash to same index

**Real-world Applications:**
- Database indexing and caching systems
- Implementing sets and maps in programming languages
- Symbol tables in compilers
- Routers for network packet forwarding
- Browser history and bookmarks

### 🔸 Hash Functions

**Properties of Good Hash Function:**
1. **Deterministic:** Same key always produces same hash
2. **Uniform Distribution:** Keys spread evenly across array
3. **Fast Computation:** O(1) time to compute hash
4. **Avalanche Effect:** Small key changes → big hash changes

**Common Hash Functions:**
```go
// Simple hash for integers (Division Method)
func simpleHash(key, tableSize int) int {
    return key % tableSize
}

// Hash for strings (djb2 algorithm)
func stringHash(key string, tableSize int) int {
    hash := 5381
    for _, char := range key {
        hash = ((hash << 5) + hash) + int(char) // hash * 33 + c
    }
    return abs(hash) % tableSize
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
```

### 🔸 Collision Resolution Strategies

**1. Separate Chaining (Closed Addressing):**
- Each bucket contains a linked list of elements
- Multiple elements can exist at same index
- Simple to implement, handles high load factors well

**2. Open Addressing:**
- All elements stored in the hash table array itself
- When collision occurs, probe for next available slot
- Types: Linear Probing, Quadratic Probing, Double Hashing

---

## 🛠️ Implementation in Golang

### Hash Table with Separate Chaining

```go
package main

import (
    "fmt"
    "hash/fnv"
)

// Entry represents a key-value pair
type Entry struct {
    Key   string
    Value interface{}
    Next  *Entry
}

// HashTable represents hash table with separate chaining
type HashTable struct {
    buckets []*Entry
    size    int
    count   int
}

// NewHashTable creates a new hash table
func NewHashTable(size int) *HashTable {
    return &HashTable{
        buckets: make([]*Entry, size),
        size:    size,
        count:   0,
    }
}

// hash computes hash value for string key
func (ht *HashTable) hash(key string) int {
    hasher := fnv.New32a()
    hasher.Write([]byte(key))
    return int(hasher.Sum32()) % ht.size
}

// Put inserts or updates key-value pair - O(1) average
func (ht *HashTable) Put(key string, value interface{}) {
    index := ht.hash(key)
    
    // Check if key already exists
    current := ht.buckets[index]
    for current != nil {
        if current.Key == key {
            current.Value = value // Update existing
            return
        }
        current = current.Next
    }
    
    // Insert new entry at beginning of chain
    newEntry := &Entry{
        Key:   key,
        Value: value,
        Next:  ht.buckets[index],
    }
    ht.buckets[index] = newEntry
    ht.count++
    
    // Resize if load factor exceeds 0.75
    if float64(ht.count)/float64(ht.size) > 0.75 {
        ht.resize()
    }
}

// Get retrieves value by key - O(1) average
func (ht *HashTable) Get(key string) (interface{}, bool) {
    index := ht.hash(key)
    current := ht.buckets[index]
    
    for current != nil {
        if current.Key == key {
            return current.Value, true
        }
        current = current.Next
    }
    
    return nil, false
}

// Delete removes key-value pair - O(1) average
func (ht *HashTable) Delete(key string) bool {
    index := ht.hash(key)
    
    if ht.buckets[index] == nil {
        return false
    }
    
    // If first entry matches
    if ht.buckets[index].Key == key {
        ht.buckets[index] = ht.buckets[index].Next
        ht.count--
        return true
    }
    
    // Search in chain
    current := ht.buckets[index]
    for current.Next != nil {
        if current.Next.Key == key {
            current.Next = current.Next.Next
            ht.count--
            return true
        }
        current = current.Next
    }
    
    return false
}

// Contains checks if key exists - O(1) average
func (ht *HashTable) Contains(key string) bool {
    _, exists := ht.Get(key)
    return exists
}

// Size returns number of key-value pairs
func (ht *HashTable) Size() int {
    return ht.count
}

// IsEmpty checks if hash table is empty
func (ht *HashTable) IsEmpty() bool {
    return ht.count == 0
}

// LoadFactor returns current load factor
func (ht *HashTable) LoadFactor() float64 {
    return float64(ht.count) / float64(ht.size)
}

// resize doubles the hash table size and rehashes all elements
func (ht *HashTable) resize() {
    oldBuckets := ht.buckets
    ht.size *= 2
    ht.buckets = make([]*Entry, ht.size)
    ht.count = 0
    
    // Rehash all existing entries
    for _, head := range oldBuckets {
        current := head
        for current != nil {
            next := current.Next
            ht.Put(current.Key, current.Value)
            current = next
        }
    }
}

// Keys returns all keys in the hash table
func (ht *HashTable) Keys() []string {
    keys := make([]string, 0, ht.count)
    for _, head := range ht.buckets {
        current := head
        for current != nil {
            keys = append(keys, current.Key)
            current = current.Next
        }
    }
    return keys
}

// Display prints all key-value pairs
func (ht *HashTable) Display() {
    fmt.Printf("Hash Table (size: %d, count: %d, load factor: %.2f)\n", 
               ht.size, ht.count, ht.LoadFactor())
    
    for i, head := range ht.buckets {
        if head != nil {
            fmt.Printf("Bucket %d: ", i)
            current := head
            for current != nil {
                fmt.Printf("(%s: %v)", current.Key, current.Value)
                if current.Next != nil {
                    fmt.Print(" -> ")
                }
                current = current.Next
            }
            fmt.Println()
        }
    }
}
```

### Hash Set Implementation

```go
// HashSet implements a set using hash table
type HashSet struct {
    table map[string]bool
}

// NewHashSet creates a new hash set
func NewHashSet() *HashSet {
    return &HashSet{
        table: make(map[string]bool),
    }
}

// Add inserts element into set - O(1) average
func (hs *HashSet) Add(item string) {
    hs.table[item] = true
}

// Remove deletes element from set - O(1) average
func (hs *HashSet) Remove(item string) {
    delete(hs.table, item)
}

// Contains checks if element exists - O(1) average
func (hs *HashSet) Contains(item string) bool {
    return hs.table[item]
}

// Size returns number of elements
func (hs *HashSet) Size() int {
    return len(hs.table)
}

// IsEmpty checks if set is empty
func (hs *HashSet) IsEmpty() bool {
    return len(hs.table) == 0
}

// Clear removes all elements
func (hs *HashSet) Clear() {
    hs.table = make(map[string]bool)
}

// Items returns slice of all elements
func (hs *HashSet) Items() []string {
    items := make([]string, 0, len(hs.table))
    for item := range hs.table {
        items = append(items, item)
    }
    return items
}

// Union returns new set with elements from both sets
func (hs *HashSet) Union(other *HashSet) *HashSet {
    result := NewHashSet()
    
    // Add all items from current set
    for item := range hs.table {
        result.Add(item)
    }
    
    // Add all items from other set
    for item := range other.table {
        result.Add(item)
    }
    
    return result
}

// Intersection returns new set with common elements
func (hs *HashSet) Intersection(other *HashSet) *HashSet {
    result := NewHashSet()
    
    // Add items that exist in both sets
    for item := range hs.table {
        if other.Contains(item) {
            result.Add(item)
        }
    }
    
    return result
}

// Difference returns new set with elements in current but not in other
func (hs *HashSet) Difference(other *HashSet) *HashSet {
    result := NewHashSet()
    
    for item := range hs.table {
        if !other.Contains(item) {
            result.Add(item)
        }
    }
    
    return result
}
```

---

## 🎯 Essential Problem-Solving Patterns

### Pattern 1: Frequency Counting

**When to Use:** Count occurrences of elements, find duplicates, anagrams
**Key Insight:** Use map to count frequency of each element

```go
// Count character frequency in string
func charFrequency(s string) map[rune]int {
    freq := make(map[rune]int)
    for _, char := range s {
        freq[char]++
    }
    return freq
}

// Check if two strings are anagrams
func areAnagrams(s1, s2 string) bool {
    if len(s1) != len(s2) {
        return false
    }
    
    freq1 := charFrequency(s1)
    freq2 := charFrequency(s2)
    
    // Compare frequency maps
    if len(freq1) != len(freq2) {
        return false
    }
    
    for char, count := range freq1 {
        if freq2[char] != count {
            return false
        }
    }
    
    return true
}
```

### Pattern 2: Two Sum Pattern

**When to Use:** Finding pairs that sum to target, complement problems
**Key Insight:** Store complement in hash map for O(1) lookup

```go
// Two Sum - Find indices of two numbers that add up to target
func twoSum(nums []int, target int) []int {
    numMap := make(map[int]int) // value -> index
    
    for i, num := range nums {
        complement := target - num
        if index, exists := numMap[complement]; exists {
            return []int{index, i}
        }
        numMap[num] = i
    }
    
    return nil // No solution found
}

// Three Sum - Find all unique triplets that sum to zero
func threeSum(nums []int) [][]int {
    result := [][]int{}
    if len(nums) < 3 {
        return result
    }
    
    // Sort array first
    sort.Ints(nums)
    
    for i := 0; i < len(nums)-2; i++ {
        // Skip duplicates for first element
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }
        
        left, right := i+1, len(nums)-1
        target := -nums[i]
        
        for left < right {
            sum := nums[left] + nums[right]
            if sum == target {
                result = append(result, []int{nums[i], nums[left], nums[right]})
                
                // Skip duplicates
                for left < right && nums[left] == nums[left+1] {
                    left++
                }
                for left < right && nums[right] == nums[right-1] {
                    right--
                }
                
                left++
                right--
            } else if sum < target {
                left++
            } else {
                right--
            }
        }
    }
    
    return result
}
```

### Pattern 3: Sliding Window with Hash Map

**When to Use:** Substring problems, window-based counting
**Key Insight:** Use hash map to track window contents

```go
// Find longest substring without repeating characters
func lengthOfLongestSubstring(s string) int {
    charMap := make(map[rune]int) // character -> last seen index
    maxLength := 0
    start := 0
    
    for i, char := range s {
        if lastIndex, exists := charMap[char]; exists && lastIndex >= start {
            start = lastIndex + 1
        }
        
        charMap[char] = i
        maxLength = max(maxLength, i-start+1)
    }
    
    return maxLength
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// Find all anagrams of pattern in string
func findAnagrams(s string, p string) []int {
    result := []int{}
    if len(s) < len(p) {
        return result
    }
    
    // Count frequency of pattern
    patternCount := make(map[rune]int)
    for _, char := range p {
        patternCount[char]++
    }
    
    windowCount := make(map[rune]int)
    windowSize := len(p)
    
    // Process first window
    for i := 0; i < windowSize; i++ {
        char := rune(s[i])
        windowCount[char]++
    }
    
    // Check if first window is anagram
    if mapsEqual(patternCount, windowCount) {
        result = append(result, 0)
    }
    
    // Slide window
    for i := windowSize; i < len(s); i++ {
        // Add new character
        newChar := rune(s[i])
        windowCount[newChar]++
        
        // Remove old character
        oldChar := rune(s[i-windowSize])
        windowCount[oldChar]--
        if windowCount[oldChar] == 0 {
            delete(windowCount, oldChar)
        }
        
        // Check if current window is anagram
        if mapsEqual(patternCount, windowCount) {
            result = append(result, i-windowSize+1)
        }
    }
    
    return result
}

func mapsEqual(map1, map2 map[rune]int) bool {
    if len(map1) != len(map2) {
        return false
    }
    
    for key, value := range map1 {
        if map2[key] != value {
            return false
        }
    }
    
    return true
}
```

### Pattern 4: Caching/Memoization

**When to Use:** Expensive computations, dynamic programming
**Key Insight:** Store computed results for reuse

```go
// LRU Cache implementation
type LRUCache struct {
    capacity int
    cache    map[int]*Node
    head     *Node
    tail     *Node
}

type Node struct {
    key   int
    value int
    prev  *Node
    next  *Node
}

func NewLRUCache(capacity int) *LRUCache {
    cache := &LRUCache{
        capacity: capacity,
        cache:    make(map[int]*Node),
        head:     &Node{}, // Dummy head
        tail:     &Node{}, // Dummy tail
    }
    cache.head.next = cache.tail
    cache.tail.prev = cache.head
    return cache
}

func (lru *LRUCache) Get(key int) int {
    if node, exists := lru.cache[key]; exists {
        // Move to front (most recently used)
        lru.moveToFront(node)
        return node.value
    }
    return -1
}

func (lru *LRUCache) Put(key int, value int) {
    if node, exists := lru.cache[key]; exists {
        // Update existing node
        node.value = value
        lru.moveToFront(node)
    } else {
        // Add new node
        if len(lru.cache) >= lru.capacity {
            // Remove least recently used
            lru.removeLRU()
        }
        
        newNode := &Node{key: key, value: value}
        lru.cache[key] = newNode
        lru.addToFront(newNode)
    }
}

func (lru *LRUCache) moveToFront(node *Node) {
    lru.removeNode(node)
    lru.addToFront(node)
}

func (lru *LRUCache) addToFront(node *Node) {
    node.prev = lru.head
    node.next = lru.head.next
    lru.head.next.prev = node
    lru.head.next = node
}

func (lru *LRUCache) removeNode(node *Node) {
    node.prev.next = node.next
    node.next.prev = node.prev
}

func (lru *LRUCache) removeLRU() {
    lru_node := lru.tail.prev
    lru.removeNode(lru_node)
    delete(lru.cache, lru_node.key)
}
```

---

## 💻 Practice Problems

### Easy Level (Day 1-2)

#### Problem 1: Two Sum (LeetCode #1)
```go
func twoSum(nums []int, target int) []int {
    // Implement using hash map
    return nil
}

// Test cases:
// Input: nums = [2,7,11,15], target = 9
// Output: [0,1]
```

#### Problem 2: Valid Anagram (LeetCode #242)
```go
func isAnagram(s string, t string) bool {
    // Use frequency counting
    return false
}
```

#### Problem 3: Contains Duplicate (LeetCode #217)
```go
func containsDuplicate(nums []int) bool {
    // Use hash set
    return false
}
```

### Medium Level (Day 3)

#### Problem 4: Group Anagrams (LeetCode #49)
```go
func groupAnagrams(strs []string) [][]string {
    // Group strings by their sorted form
    return nil
}
```

#### Problem 5: Longest Substring Without Repeating Characters (LeetCode #3)
```go
func lengthOfLongestSubstring(s string) int {
    // Use sliding window with hash map
    return 0
}
```

#### Problem 6: Top K Frequent Elements (LeetCode #347)
```go
func topKFrequent(nums []int, k int) []int {
    // Count frequency then find top k
    return nil
}
```

---

## 📊 Time & Space Complexity Analysis

### Hash Table Operations
| Operation | Average Case | Worst Case | Space |
|-----------|--------------|-------------|-------|
| Insert | O(1) | O(n) | O(1) |
| Delete | O(1) | O(n) | O(1) |
| Search | O(1) | O(n) | O(1) |
| Space | O(n) | O(n) | O(n) |

### Load Factor Impact
- **Low Load Factor (< 0.5):** Better performance, more memory waste
- **High Load Factor (> 0.8):** More collisions, worse performance
- **Optimal Load Factor (0.6-0.75):** Good balance of time and space

### Collision Resolution Comparison
| Method | Average Insert | Average Search | Space Overhead |
|--------|----------------|----------------|----------------|
| Separate Chaining | O(1) | O(1 + α) | High (pointers) |
| Linear Probing | O(1) | O(1) | Low |
| Quadratic Probing | O(1) | O(1) | Low |

Where α = load factor (n/m)

---

## 🎓 Study Tips & Best Practices

### 1. Choosing Hash Function
```go
// Good for integers
func intHash(key, size int) int {
    return ((key * 2654435761) >> 22) % size
}

// Good for strings (avoid collisions)
func stringHashDJB2(key string, size int) int {
    hash := 5381
    for _, char := range key {
        hash = ((hash << 5) + hash) + int(char)
    }
    return abs(hash) % size
}
```

### 2. Handling Edge Cases
```go
// Always check for nil/empty inputs
if key == "" || value == nil {
    return errors.New("invalid input")
}

// Handle hash table resizing
if loadFactor() > threshold {
    resize()
}

// Handle collisions gracefully
if bucket[index] != nil {
    // Collision detected - use chaining or probing
}
```

### 3. Memory Management
```go
// Be aware of memory usage
type HashTable struct {
    buckets [][]Entry  // Slice of slices
    size    int
    count   int
    maxLoadFactor float64
}

// Clean up when deleting
func (ht *HashTable) Delete(key string) {
    // Remove entry and decrease count
    // Consider shrinking if load factor too low
}
```

---

## 🚀 Next Steps

### Day 1 Goals
- [ ] Implement basic hash table with separate chaining
- [ ] Understand hash function properties
- [ ] Solve two sum and anagram problems
- [ ] Practice frequency counting patterns

### Day 2 Goals
- [ ] Implement hash set operations
- [ ] Master sliding window with hash map
- [ ] Solve substring problems
- [ ] Understand collision resolution trade-offs

### Day 3 Goals
- [ ] Implement LRU cache
- [ ] Practice advanced hash table patterns
- [ ] Solve grouping and counting problems
- [ ] Master time complexity analysis

### Preparation for Trees (Days 4-7)
- Review recursive thinking
- Understand parent-child relationships
- Practice tree drawing and visualization
- Think about hierarchical data organization

**Hash tables are the secret weapon for many interview problems! Master these patterns and you'll solve problems that seem impossible in linear time. Checkmate!** 🗂️

---

## 🔧 Common Interview Tips

### Hash Table Red Flags:
1. **"Find pair/complement"** → Think Two Sum pattern
2. **"Count frequency"** → Think frequency map
3. **"Group similar items"** → Think hash map grouping
4. **"Fast lookup required"** → Think hash table
5. **"Seen before/cache"** → Think hash set/map

### Performance Optimization:
```go
// Use appropriate initial capacity
hashMap := make(map[string]int, expectedSize)

// Choose right data structure
// map[string]int vs map[int]int vs []int (if keys are small integers)

// Consider memory vs time trade-offs
// Sometimes O(n) scan is better than hash table overhead
```

**Remember: Hash tables turn many O(n²) problems into O(n) solutions. That's pure strategic advantage!** 💪