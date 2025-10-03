package main

import (
	"fmt"
	"testing"
)

// Test basic RadixTree operations
func TestRadixTreeBasicOperations(t *testing.T) {
	rt := NewRadixTree()

	// Test empty tree
	if !rt.IsEmpty() {
		t.Error("New radix tree should be empty")
	}

	if rt.Size() != 0 {
		t.Error("New radix tree size should be 0")
	}

	// Test insertion
	words := []string{"hello", "hell", "help", "hero", "heroine"}
	for _, word := range words {
		rt.Insert(word)
	}

	if rt.Size() != 5 {
		t.Errorf("Expected size 5, got %d", rt.Size())
	}

	// Test search
	for _, word := range words {
		if !rt.Search(word) {
			t.Errorf("Word '%s' should exist in radix tree", word)
		}
	}

	// Test non-existent words
	nonExistent := []string{"he", "helper", "heroic"}
	for _, word := range nonExistent {
		if rt.Search(word) {
			t.Errorf("Word '%s' should not exist in radix tree", word)
		}
	}
}

// Test RadixTree prefix operations
func TestRadixTreePrefixOperations(t *testing.T) {
	rt := NewRadixTree()
	words := []string{"apple", "app", "application", "apply", "banana", "band", "bandana"}

	for _, word := range words {
		rt.Insert(word)
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
	}

	for _, tc := range testCases {
		result := rt.StartsWith(tc.prefix)
		if result != tc.expected {
			t.Errorf("StartsWith('%s') = %v, expected %v", tc.prefix, result, tc.expected)
		}
	}

	// Test GetWordsWithPrefix
	appWords := rt.GetWordsWithPrefix("app")
	expectedAppWords := 4
	if len(appWords) != expectedAppWords {
		t.Errorf("Expected %d words with prefix 'app', got %d", expectedAppWords, len(appWords))
	}

	banWords := rt.GetWordsWithPrefix("ban")
	expectedBanWords := 3
	if len(banWords) != expectedBanWords {
		t.Errorf("Expected %d words with prefix 'ban', got %d", expectedBanWords, len(banWords))
	}
}

// Test RadixTree compression efficiency
func TestRadixTreeCompression(t *testing.T) {
	rt := NewRadixTree()

	// Insert words with common prefixes
	words := []string{
		"test", "testing", "tester", "tests",
		"technology", "technical", "technique",
		"car", "card", "cards", "care", "careful",
	}

	for _, word := range words {
		rt.Insert(word)
	}

	// Get memory efficiency comparison
	radixNodes, trieNodes := rt.GetMemoryEfficiency()

	fmt.Printf("Radix Tree nodes: %d\n", radixNodes)
	fmt.Printf("Equivalent Trie nodes: %d\n", trieNodes)
	fmt.Printf("Space savings: %.2f%%\n", float64(trieNodes-radixNodes)/float64(trieNodes)*100)

	// Radix tree should use significantly fewer nodes
	if radixNodes >= trieNodes {
		t.Error("Radix tree should use fewer nodes than equivalent trie")
	}
}

// Test RadixTree autocomplete
func TestRadixTreeAutocomplete(t *testing.T) {
	rt := NewRadixTree()

	// Insert words with different frequencies
	words := map[string]int{
		"apple":       5,
		"app":         3,
		"application": 2,
		"apply":       1,
	}

	for word, freq := range words {
		for i := 0; i < freq; i++ {
			rt.Insert(word)
		}
	}

	// Test word count
	for word, expectedFreq := range words {
		actualFreq := rt.GetWordCount(word)
		if actualFreq != expectedFreq {
			t.Errorf("Word '%s' frequency = %d, expected %d", word, actualFreq, expectedFreq)
		}
	}

	// Test autocomplete
	suggestions := rt.AutoComplete("app", 3)
	if len(suggestions) != 3 {
		t.Errorf("Expected 3 suggestions, got %d", len(suggestions))
	}

	// Should be sorted by frequency (descending)
	expectedOrder := []string{"apple", "app", "application"}
	for i, expected := range expectedOrder {
		if i < len(suggestions) && suggestions[i] != expected {
			t.Errorf("Suggestion[%d] = '%s', expected '%s'", i, suggestions[i], expected)
		}
	}
}

// Test RadixTree edge cases
func TestRadixTreeEdgeCases(t *testing.T) {
	rt := NewRadixTree()

	// Test single character
	rt.Insert("a")
	if !rt.Search("a") {
		t.Error("Single character 'a' should be searchable")
	}

	// Test overlapping words
	rt.Insert("car")
	rt.Insert("card")
	rt.Insert("care")
	rt.Insert("careful")

	if !rt.Search("car") {
		t.Error("Word 'car' should exist")
	}
	if !rt.Search("card") {
		t.Error("Word 'card' should exist")
	}
	if !rt.Search("care") {
		t.Error("Word 'care' should exist")
	}
	if !rt.Search("careful") {
		t.Error("Word 'careful' should exist")
	}

	// Test prefix that is also a word
	if !rt.StartsWith("car") {
		t.Error("Should find prefix 'car'")
	}

	if rt.CountWordsWithPrefix("car") != 4 {
		t.Errorf("Expected 4 words with prefix 'car', got %d", rt.CountWordsWithPrefix("car"))
	}
}

// Test RadixTree deletion
func TestRadixTreeDeletion(t *testing.T) {
	rt := NewRadixTree()
	words := []string{"car", "card", "care", "careful"}

	for _, word := range words {
		rt.Insert(word)
	}

	initialSize := rt.Size()

	// Test deleting existing word
	if !rt.Delete("care") {
		t.Error("Should be able to delete existing word 'care'")
	}

	if rt.Search("care") {
		t.Error("Word 'care' should not exist after deletion")
	}

	if rt.Size() != initialSize-1 {
		t.Errorf("Size should decrease by 1 after deletion, got %d", rt.Size())
	}

	// Test that other words still exist
	if !rt.Search("car") {
		t.Error("Word 'car' should still exist")
	}
	if !rt.Search("card") {
		t.Error("Word 'card' should still exist")
	}
	if !rt.Search("careful") {
		t.Error("Word 'careful' should still exist")
	}
}

// Benchmark RadixTree vs Trie comparison
func BenchmarkRadixTreeInsert(b *testing.B) {
	rt := NewRadixTree()
	words := generateWords(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := words[i%len(words)]
		rt.Insert(word)
	}
}

func BenchmarkRadixTreeSearch(b *testing.B) {
	rt := NewRadixTree()
	words := generateWords(1000)

	// Pre-populate radix tree
	for _, word := range words {
		rt.Insert(word)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := words[i%len(words)]
		rt.Search(word)
	}
}

// Test comprehensive RadixTree example
func TestRadixTreeExample(t *testing.T) {
	fmt.Println("\n=== Radix Tree (Compressed Trie) Demo ===")

	rt := NewRadixTree()

	// Build a dictionary with many common prefixes
	dictionary := []string{
		"test", "testing", "tester", "tests", "testify",
		"car", "card", "cards", "care", "careful", "careless", "caring",
		"app", "apple", "application", "apply", "applying", "applied",
		"book", "booking", "bookish", "bookstore", "bookmark",
	}

	fmt.Println("Building dictionary with common prefixes...")
	for _, word := range dictionary {
		rt.Insert(word)
		// Insert some words multiple times for frequency
		if word == "test" || word == "car" || word == "app" {
			rt.Insert(word)
			rt.Insert(word)
		}
	}

	fmt.Printf("Dictionary size: %d words\n", rt.Size())

	// Show memory efficiency
	radixNodes, trieNodes := rt.GetMemoryEfficiency()
	fmt.Printf("Memory efficiency: %d radix nodes vs %d trie nodes (%.1f%% savings)\n",
		radixNodes, trieNodes, float64(trieNodes-radixNodes)/float64(trieNodes)*100)

	// Search examples
	searchWords := []string{"test", "testing", "car", "xyz", "book"}
	fmt.Println("\nSearch results:")
	for _, word := range searchWords {
		exists := rt.Search(word)
		count := rt.GetWordCount(word)
		fmt.Printf("  '%s': exists=%v, count=%d\n", word, exists, count)
	}

	// Prefix examples
	prefixes := []string{"test", "car", "app", "book"}
	fmt.Println("\nPrefix search results:")
	for _, prefix := range prefixes {
		hasPrefix := rt.StartsWith(prefix)
		words := rt.GetWordsWithPrefix(prefix)
		count := rt.CountWordsWithPrefix(prefix)
		fmt.Printf("  '%s': exists=%v, count=%d, words=%v\n", prefix, hasPrefix, count, words[:min(3, len(words))])
	}

	// Autocomplete examples
	fmt.Println("\nAutocomplete suggestions (top 3):")
	autocompletePrefixes := []string{"test", "car", "app"}
	for _, prefix := range autocompletePrefixes {
		suggestions := rt.AutoComplete(prefix, 3)
		fmt.Printf("  '%s': %v\n", prefix, suggestions)
	}

	fmt.Println("\n=== Radix Tree Demo Complete ===")
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
