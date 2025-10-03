package main

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
)

// ===== HASH TABLE WITH SEPARATE CHAINING =====

// Entry represents a key-value pair in hash table
type Entry struct {
	Key   string
	Value interface{}
	Next  *Entry
}

// HashTable with separate chaining collision resolution
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

// hash computes hash value for string key using FNV-1a
func (ht *HashTable) hash(key string) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	return int(hasher.Sum32()) % ht.size
}

// Put inserts or updates key-value pair - O(1) average
func (ht *HashTable) Put(key string, value interface{}) {
	index := ht.hash(key)

	// Check if key already exists (update case)
	current := ht.buckets[index]
	for current != nil {
		if current.Key == key {
			current.Value = value
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

	// Auto-resize if load factor exceeds 0.75
	if ht.LoadFactor() > 0.75 {
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

// LoadFactor returns current load factor
func (ht *HashTable) LoadFactor() float64 {
	return float64(ht.count) / float64(ht.size)
}

// resize doubles the hash table size and rehashes all elements
func (ht *HashTable) resize() {
	fmt.Printf("Resizing hash table from %d to %d\n", ht.size, ht.size*2)

	oldBuckets := ht.buckets
	oldCount := ht.count

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

	fmt.Printf("Rehashed %d entries\n", oldCount)
}

// Display prints the hash table structure
func (ht *HashTable) Display() {
	fmt.Printf("\nHash Table (size: %d, count: %d, load factor: %.2f)\n",
		ht.size, ht.count, ht.LoadFactor())
	fmt.Println(strings.Repeat("=", 50))

	for i, head := range ht.buckets {
		if head != nil {
			fmt.Printf("Bucket %2d: ", i)
			current := head
			for current != nil {
				fmt.Printf("[%s: %v]", current.Key, current.Value)
				if current.Next != nil {
					fmt.Print(" -> ")
				}
				current = current.Next
			}
			fmt.Println()
		}
	}
	fmt.Println()
}

// ===== HASH SET IMPLEMENTATION =====

type HashSet struct {
	table map[string]bool
}

func NewHashSet() *HashSet {
	return &HashSet{
		table: make(map[string]bool),
	}
}

func (hs *HashSet) Add(item string) {
	hs.table[item] = true
}

func (hs *HashSet) Contains(item string) bool {
	return hs.table[item]
}

func (hs *HashSet) Remove(item string) {
	delete(hs.table, item)
}

func (hs *HashSet) Size() int {
	return len(hs.table)
}

func (hs *HashSet) Items() []string {
	items := make([]string, 0, len(hs.table))
	for item := range hs.table {
		items = append(items, item)
	}
	sort.Strings(items) // For consistent output
	return items
}

// ===== COMMON HASH TABLE PATTERNS =====

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

	return []int{} // No solution found
}

// Check if two strings are anagrams using frequency counting
func areAnagrams(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	charCount := make(map[rune]int)

	// Count characters in first string
	for _, char := range s1 {
		charCount[char]++
	}

	// Subtract characters from second string
	for _, char := range s2 {
		charCount[char]--
		if charCount[char] == 0 {
			delete(charCount, char)
		}
	}

	return len(charCount) == 0
}

// Group anagrams together
func groupAnagrams(strs []string) [][]string {
	groups := make(map[string][]string)

	for _, str := range strs {
		// Sort characters to create key
		chars := []rune(str)
		sort.Slice(chars, func(i, j int) bool {
			return chars[i] < chars[j]
		})
		key := string(chars)

		groups[key] = append(groups[key], str)
	}

	result := make([][]string, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}

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
		currentLength := i - start + 1
		if currentLength > maxLength {
			maxLength = currentLength
		}
	}

	return maxLength
}

// ===== DEMONSTRATION FUNCTION =====

func main() {
	fmt.Println("🗂️ Hash Tables Demonstration in Golang")
	fmt.Println("=====================================")

	// Create hash table
	ht := NewHashTable(4)

	// Insert some data
	fmt.Println("Inserting key-value pairs...")
	ht.Put("apple", 5)
	ht.Put("banana", 3)
	ht.Put("orange", 8)
	ht.Put("grape", 12)
	ht.Put("kiwi", 7)
	ht.Put("mango", 15) // This should trigger resize

	ht.Display()

	// Test retrieval
	fmt.Println("Testing retrieval:")
	if value, exists := ht.Get("banana"); exists {
		fmt.Printf("banana: %v\n", value)
	}

	if value, exists := ht.Get("nonexistent"); exists {
		fmt.Printf("nonexistent: %v\n", value)
	} else {
		fmt.Println("nonexistent: key not found")
	}

	// Test deletion
	fmt.Println("\nDeleting 'orange'...")
	ht.Delete("orange")
	ht.Display()

	// Hash Set demonstration
	fmt.Println("Hash Set demonstration:")
	set := NewHashSet()
	set.Add("red")
	set.Add("green")
	set.Add("blue")
	set.Add("red") // Duplicate - should not increase size

	fmt.Printf("Set size: %d\n", set.Size())
	fmt.Printf("Contains 'red': %t\n", set.Contains("red"))
	fmt.Printf("Contains 'yellow': %t\n", set.Contains("yellow"))
	fmt.Printf("All items: %v\n", set.Items())

	// Pattern demonstrations
	fmt.Println("\nPattern demonstrations:")

	// Two Sum
	nums := []int{2, 7, 11, 15}
	target := 9
	result := twoSum(nums, target)
	fmt.Printf("Two Sum: nums=%v, target=%d, result=%v\n", nums, target, result)

	// Anagrams
	s1, s2 := "listen", "silent"
	fmt.Printf("Are '%s' and '%s' anagrams? %t\n", s1, s2, areAnagrams(s1, s2))

	// Longest substring
	testStr := "abcabcbb"
	fmt.Printf("Longest substring without repeating chars in '%s': %d\n",
		testStr, lengthOfLongestSubstring(testStr))

	// Group anagrams
	words := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
	groups := groupAnagrams(words)
	fmt.Printf("Grouped anagrams: %v\n", groups)
}
