package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// Test basic insertion and search functionality
func TestTrieBasicOperations(t *testing.T) {
	trie := NewTrie()

	// Test empty trie
	if !trie.IsEmpty() {
		t.Error("New trie should be empty")
	}

	if trie.Size() != 0 {
		t.Error("New trie size should be 0")
	}

	// Test insertion
	words := []string{"hello", "world", "hi", "help", "hell"}
	for _, word := range words {
		trie.Insert(word)
	}

	if trie.Size() != 5 {
		t.Errorf("Expected size 5, got %d", trie.Size())
	}

	// Test search
	for _, word := range words {
		if !trie.Search(word) {
			t.Errorf("Word '%s' should exist in trie", word)
		}
	}

	// Test non-existent words
	nonExistent := []string{"hel", "helper", "wo", "worlds"}
	for _, word := range nonExistent {
		if trie.Search(word) {
			t.Errorf("Word '%s' should not exist in trie", word)
		}
	}
}

// Test prefix functionality
func TestTriePrefixOperations(t *testing.T) {
	trie := NewTrie()
	words := []string{"apple", "app", "application", "apply", "banana", "band", "bandana"}

	for _, word := range words {
		trie.Insert(word)
	}

	// Test StartsWith
	testCases := []struct {
		prefix   string
		expected bool
	}{
		{"app", true},
		{"ban", true},
		{"band", true},
		{"orange", false},
		{"xyz", false},
		{"", true}, // Empty prefix should match
	}

	for _, tc := range testCases {
		result := trie.StartsWith(tc.prefix)
		if result != tc.expected {
			t.Errorf("StartsWith('%s') = %v, expected %v", tc.prefix, result, tc.expected)
		}
	}

	// Test GetWordsWithPrefix
	appWords := trie.GetWordsWithPrefix("app")
	expectedAppWords := []string{"app", "apple", "application", "apply"}
	if len(appWords) != len(expectedAppWords) {
		t.Errorf("Expected %d words with prefix 'app', got %d", len(expectedAppWords), len(appWords))
	}

	banWords := trie.GetWordsWithPrefix("ban")
	expectedBanWords := []string{"banana", "band", "bandana"}
	if len(banWords) != len(expectedBanWords) {
		t.Errorf("Expected %d words with prefix 'ban', got %d", len(expectedBanWords), len(banWords))
	}
}

// Test deletion functionality
func TestTrieDeletion(t *testing.T) {
	trie := NewTrie()
	words := []string{"cat", "car", "card", "care", "careful"}

	for _, word := range words {
		trie.Insert(word)
	}

	initialSize := trie.Size()

	// Test deleting existing word
	if !trie.Delete("care") {
		t.Error("Should be able to delete existing word 'care'")
	}

	if trie.Search("care") {
		t.Error("Word 'care' should not exist after deletion")
	}

	if trie.Size() != initialSize-1 {
		t.Errorf("Size should decrease by 1 after deletion, got %d", trie.Size())
	}

	// Test that prefix still works for remaining words
	if !trie.StartsWith("car") {
		t.Error("Prefix 'car' should still exist")
	}

	// Test deleting non-existent word
	if trie.Delete("nonexistent") {
		t.Error("Should not be able to delete non-existent word")
	}

	// Test deleting word that is prefix of another
	if !trie.Delete("car") {
		t.Error("Should be able to delete word 'car'")
	}

	if trie.Search("car") {
		t.Error("Word 'car' should not exist after deletion")
	}

	if !trie.Search("card") {
		t.Error("Word 'card' should still exist after deleting 'car'")
	}
}

// Test frequency counting and autocomplete
func TestTrieFrequencyAndAutocomplete(t *testing.T) {
	trie := NewTrie()

	// Insert words with different frequencies
	words := map[string]int{
		"apple":   5,
		"app":     3,
		"apply":   2,
		"apricot": 1,
	}

	for word, freq := range words {
		for i := 0; i < freq; i++ {
			trie.Insert(word)
		}
	}

	// Test word count
	for word, expectedFreq := range words {
		actualFreq := trie.GetWordCount(word)
		if actualFreq != expectedFreq {
			t.Errorf("Word '%s' frequency = %d, expected %d", word, actualFreq, expectedFreq)
		}
	}

	// Test autocomplete
	suggestions := trie.AutoComplete("app", 3)
	if len(suggestions) != 3 {
		t.Errorf("Expected 3 suggestions, got %d", len(suggestions))
	}

	// Should be sorted by frequency (descending)
	expectedOrder := []string{"apple", "app", "apply"}
	for i, expected := range expectedOrder {
		if i < len(suggestions) && suggestions[i] != expected {
			t.Errorf("Suggestion[%d] = '%s', expected '%s'", i, suggestions[i], expected)
		}
	}
}

// Test Unicode support
func TestTrieUnicodeSupport(t *testing.T) {
	trie := NewTrie()

	unicodeWords := []string{"café", "naïve", "résumé", "日本語", "🚀rocket"}

	for _, word := range unicodeWords {
		trie.Insert(word)
	}

	for _, word := range unicodeWords {
		if !trie.Search(word) {
			t.Errorf("Unicode word '%s' should exist in trie", word)
		}
	}

	// Test prefix with Unicode
	if !trie.StartsWith("caf") {
		t.Error("Should find prefix 'caf' for word 'café'")
	}
}

// Test edge cases
func TestTrieEdgeCases(t *testing.T) {
	trie := NewTrie()

	// Test empty string
	trie.Insert("")
	if trie.Search("") {
		t.Error("Empty string should not be searchable")
	}

	// Test single character
	trie.Insert("a")
	if !trie.Search("a") {
		t.Error("Single character 'a' should be searchable")
	}

	// Test very long string
	longString := strings.Repeat("abcdefghijk", 100)
	trie.Insert(longString)
	if !trie.Search(longString) {
		t.Error("Long string should be searchable")
	}

	// Test duplicate insertions
	initialSize := trie.Size()
	trie.Insert("duplicate")
	trie.Insert("duplicate")
	trie.Insert("duplicate")

	if trie.Size() != initialSize+1 {
		t.Error("Duplicate insertions should not increase size")
	}

	if trie.GetWordCount("duplicate") != 3 {
		t.Error("Duplicate insertions should increase frequency count")
	}
}

// Test longest common prefix
func TestTrieLongestCommonPrefix(t *testing.T) {
	trie := NewTrie()

	// Test with common prefix
	words := []string{"flower", "flow", "flight"}
	for _, word := range words {
		trie.Insert(word)
	}

	lcp := trie.LongestCommonPrefix()
	if lcp != "fl" {
		t.Errorf("Expected longest common prefix 'fl', got '%s'", lcp)
	}

	// Test with no common prefix
	trie.Clear()
	words = []string{"dog", "cat", "bird"}
	for _, word := range words {
		trie.Insert(word)
	}

	lcp = trie.LongestCommonPrefix()
	if lcp != "" {
		t.Errorf("Expected empty longest common prefix, got '%s'", lcp)
	}

	// Test empty trie
	trie.Clear()
	lcp = trie.LongestCommonPrefix()
	if lcp != "" {
		t.Errorf("Expected empty longest common prefix for empty trie, got '%s'", lcp)
	}
}

// Test counting words with prefix
func TestTrieCountWordsWithPrefix(t *testing.T) {
	trie := NewTrie()
	words := []string{"apple", "app", "application", "apply", "banana", "band"}

	for _, word := range words {
		trie.Insert(word)
	}

	testCases := []struct {
		prefix   string
		expected int
	}{
		{"app", 4},
		{"ban", 2},
		{"ap", 4},
		{"xyz", 0},
		{"", 6}, // All words
	}

	for _, tc := range testCases {
		count := trie.CountWordsWithPrefix(tc.prefix)
		if count != tc.expected {
			t.Errorf("CountWordsWithPrefix('%s') = %d, expected %d", tc.prefix, count, tc.expected)
		}
	}
}

// Test clear functionality
func TestTrieClear(t *testing.T) {
	trie := NewTrie()
	words := []string{"hello", "world", "test"}

	for _, word := range words {
		trie.Insert(word)
	}

	if trie.IsEmpty() {
		t.Error("Trie should not be empty after insertions")
	}

	trie.Clear()

	if !trie.IsEmpty() {
		t.Error("Trie should be empty after clear")
	}

	if trie.Size() != 0 {
		t.Error("Trie size should be 0 after clear")
	}

	for _, word := range words {
		if trie.Search(word) {
			t.Errorf("Word '%s' should not exist after clear", word)
		}
	}
}

// Benchmark basic operations
func BenchmarkTrieInsert(b *testing.B) {
	trie := NewTrie()
	words := generateWords(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := words[i%len(words)]
		trie.Insert(word)
	}
}

func BenchmarkTrieSearch(b *testing.B) {
	trie := NewTrie()
	words := generateWords(1000)

	// Pre-populate trie
	for _, word := range words {
		trie.Insert(word)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := words[i%len(words)]
		trie.Search(word)
	}
}

func BenchmarkTrieStartsWith(b *testing.B) {
	trie := NewTrie()
	words := generateWords(1000)

	// Pre-populate trie
	for _, word := range words {
		trie.Insert(word)
	}

	prefixes := make([]string, 100)
	for i := 0; i < 100; i++ {
		word := words[i%len(words)]
		if len(word) > 2 {
			prefixes[i] = word[:len(word)/2]
		} else {
			prefixes[i] = word
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prefix := prefixes[i%len(prefixes)]
		trie.StartsWith(prefix)
	}
}

func BenchmarkTrieAutocomplete(b *testing.B) {
	trie := NewTrie()
	words := generateWords(1000)

	// Pre-populate trie with frequencies
	for _, word := range words {
		// Insert each word 1-5 times randomly
		freq := (len(word) % 5) + 1
		for j := 0; j < freq; j++ {
			trie.Insert(word)
		}
	}

	prefixes := make([]string, 100)
	for i := 0; i < 100; i++ {
		word := words[i%len(words)]
		if len(word) > 2 {
			prefixes[i] = word[:2]
		} else {
			prefixes[i] = word
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prefix := prefixes[i%len(prefixes)]
		trie.AutoComplete(prefix, 5)
	}
}

// Benchmark vs naive string operations
func BenchmarkNaiveStringSearch(b *testing.B) {
	words := generateWords(1000)
	searchWords := words[:100] // Search in first 100 words

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchWord := searchWords[i%len(searchWords)]
		found := false
		for _, word := range words {
			if word == searchWord {
				found = true
				break
			}
		}
		_ = found
	}
}

func BenchmarkNaivePrefixSearch(b *testing.B) {
	words := generateWords(1000)

	prefixes := make([]string, 100)
	for i := 0; i < 100; i++ {
		word := words[i%len(words)]
		if len(word) > 2 {
			prefixes[i] = word[:2]
		} else {
			prefixes[i] = word
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prefix := prefixes[i%len(prefixes)]
		var matches []string
		for _, word := range words {
			if strings.HasPrefix(word, prefix) {
				matches = append(matches, word)
			}
		}
		_ = matches
	}
}

// Helper function to generate test words
func generateWords(count int) []string {
	words := make([]string, count)
	prefixes := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew"}
	suffixes := []string{"", "s", "ing", "ed", "er", "est", "ly", "tion", "ness", "ment"}

	for i := 0; i < count; i++ {
		prefix := prefixes[i%len(prefixes)]
		suffix := suffixes[i%len(suffixes)]
		words[i] = prefix + suffix
	}

	return words
}

// Test comprehensive example usage
func TestTrieExample(t *testing.T) {
	fmt.Println("\n=== Trie Data Structure Demo ===")

	trie := NewTrie()

	// Build a dictionary
	dictionary := []string{
		"apple", "app", "application", "apply", "apricot",
		"banana", "band", "bandana", "bandage",
		"cat", "car", "card", "care", "careful", "careless",
	}

	fmt.Println("Building dictionary...")
	for _, word := range dictionary {
		trie.Insert(word)
		// Insert some words multiple times for frequency
		if word == "apple" || word == "app" {
			trie.Insert(word)
			trie.Insert(word)
		}
	}

	fmt.Printf("Dictionary size: %d words\n", trie.Size())

	// Search examples
	searchWords := []string{"apple", "app", "application", "xyz", "car"}
	fmt.Println("\nSearch results:")
	for _, word := range searchWords {
		exists := trie.Search(word)
		count := trie.GetWordCount(word)
		fmt.Printf("  '%s': exists=%v, count=%d\n", word, exists, count)
	}

	// Prefix examples
	prefixes := []string{"app", "ban", "car", "xyz"}
	fmt.Println("\nPrefix search results:")
	for _, prefix := range prefixes {
		hasPrefix := trie.StartsWith(prefix)
		words := trie.GetWordsWithPrefix(prefix)
		count := trie.CountWordsWithPrefix(prefix)
		fmt.Printf("  '%s': exists=%v, count=%d, words=%v\n", prefix, hasPrefix, count, words)
	}

	// Autocomplete examples
	fmt.Println("\nAutocomplete suggestions (top 3):")
	autocompletePrefixes := []string{"app", "ban", "car"}
	for _, prefix := range autocompletePrefixes {
		suggestions := trie.AutoComplete(prefix, 3)
		fmt.Printf("  '%s': %v\n", prefix, suggestions)
	}

	// Longest common prefix
	lcp := trie.LongestCommonPrefix()
	fmt.Printf("\nLongest common prefix: '%s'\n", lcp)

	// Performance comparison
	fmt.Println("\n=== Performance Comparison ===")

	// Measure trie operations
	start := time.Now()
	for i := 0; i < 10000; i++ {
		word := dictionary[i%len(dictionary)]
		trie.Search(word)
	}
	trieTime := time.Since(start)

	// Measure naive search
	start = time.Now()
	for i := 0; i < 10000; i++ {
		searchWord := dictionary[i%len(dictionary)]
		found := false
		for _, word := range dictionary {
			if word == searchWord {
				found = true
				break
			}
		}
		_ = found
	}
	naiveTime := time.Since(start)

	fmt.Printf("Trie search (10k ops): %v\n", trieTime)
	fmt.Printf("Naive search (10k ops): %v\n", naiveTime)
	fmt.Printf("Speedup: %.2fx\n", float64(naiveTime)/float64(trieTime))

	fmt.Println("\n=== Demo Complete ===")
}
