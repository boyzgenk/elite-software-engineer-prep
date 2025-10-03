package main

import (
	"fmt"
	"sort"
	"strings"
)

// RadixNode represents a single node in the Radix Tree (Compressed Trie)
type RadixNode struct {
	// Edge label - the string that this edge represents
	label string

	// Children nodes mapped by the first character of their edge labels
	children map[rune]*RadixNode

	// Marks if this node represents the end of a word
	isEnd bool

	// Count of words ending at this node (for frequency tracking)
	count int
}

// NewRadixNode creates a new RadixNode
func NewRadixNode(label string) *RadixNode {
	return &RadixNode{
		label:    label,
		children: make(map[rune]*RadixNode),
		isEnd:    false,
		count:    0,
	}
}

// RadixTree represents the Radix Tree (Compressed Trie) data structure
type RadixTree struct {
	root *RadixNode
	size int // Total number of unique words in the tree
}

// NewRadixTree creates a new RadixTree
func NewRadixTree() *RadixTree {
	return &RadixTree{
		root: NewRadixNode(""), // Root has empty label
		size: 0,
	}
}

// Insert adds a word to the radix tree
// Time Complexity: O(m) where m is the length of the word
// Space Complexity: O(m) in worst case
func (rt *RadixTree) Insert(word string) {
	if word == "" {
		return
	}

	current := rt.root
	remainingWord := word

	for remainingWord != "" {
		found := false

		// Look for a child that starts with the first character
		firstChar := rune(remainingWord[0])
		if child, exists := current.children[firstChar]; exists {
			// Find the common prefix
			commonPrefixLen := rt.longestCommonPrefix(child.label, remainingWord)

			if commonPrefixLen == len(child.label) {
				// The child's label is a prefix of remaining word
				current = child
				remainingWord = remainingWord[commonPrefixLen:]
				found = true
			} else if commonPrefixLen > 0 {
				// Need to split the child node
				// Create new internal node with common prefix
				commonPrefix := child.label[:commonPrefixLen]
				oldSuffix := child.label[commonPrefixLen:]
				newSuffix := remainingWord[commonPrefixLen:]

				// Create new internal node
				internalNode := NewRadixNode(commonPrefix)

				// Update child's label to the remaining part
				child.label = oldSuffix

				// Add child as child of internal node
				if oldSuffix != "" {
					internalNode.children[rune(oldSuffix[0])] = child
				} else {
					// The split point is at the end of child's label
					internalNode.isEnd = child.isEnd
					internalNode.count = child.count
					child.isEnd = false
					child.count = 0
				}

				// Replace child with internal node
				current.children[firstChar] = internalNode

				// Continue with internal node
				current = internalNode
				remainingWord = newSuffix
				found = true
			}
		}

		if !found {
			// No matching child found, create new leaf
			newNode := NewRadixNode(remainingWord)
			newNode.isEnd = true
			newNode.count = 1
			current.children[rune(remainingWord[0])] = newNode
			rt.size++
			return
		}

		// If we've consumed all the word
		if remainingWord == "" {
			if !current.isEnd {
				current.isEnd = true
				rt.size++
			}
			current.count++
			return
		}
	}
}

// longestCommonPrefix finds the longest common prefix between two strings
func (rt *RadixTree) longestCommonPrefix(s1, s2 string) int {
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}

	for i := 0; i < minLen; i++ {
		if s1[i] != s2[i] {
			return i
		}
	}
	return minLen
}

// Search checks if a word exists in the radix tree
// Time Complexity: O(m) where m is the length of the word
// Space Complexity: O(1)
func (rt *RadixTree) Search(word string) bool {
	node, remainingWord := rt.searchNode(word)
	return node != nil && remainingWord == "" && node.isEnd
}

// StartsWith checks if there are any words starting with the given prefix
// Time Complexity: O(p) where p is the length of the prefix
// Space Complexity: O(1)
func (rt *RadixTree) StartsWith(prefix string) bool {
	node, remainingWord := rt.searchNode(prefix)
	return node != nil && remainingWord == ""
}

// searchNode finds the node corresponding to a word/prefix
// Returns the node and any remaining unmatched part of the word
func (rt *RadixTree) searchNode(word string) (*RadixNode, string) {
	current := rt.root
	remainingWord := word

	for remainingWord != "" {
		firstChar := rune(remainingWord[0])
		child, exists := current.children[firstChar]

		if !exists {
			return nil, remainingWord
		}

		// Check if the child's label matches the beginning of remaining word
		if len(child.label) <= len(remainingWord) &&
			child.label == remainingWord[:len(child.label)] {
			// Perfect match, continue traversal
			current = child
			remainingWord = remainingWord[len(child.label):]
		} else if strings.HasPrefix(child.label, remainingWord) {
			// The remaining word is a prefix of child's label
			return child, ""
		} else {
			// No match
			return nil, remainingWord
		}
	}

	return current, remainingWord
}

// Delete removes a word from the radix tree
// Time Complexity: O(m) where m is the length of the word
// Space Complexity: O(1)
func (rt *RadixTree) Delete(word string) bool {
	if word == "" {
		return false
	}

	// First check if word exists
	if !rt.Search(word) {
		return false
	}

	// Simple approach: mark as not end and decrease count
	node, remainingWord := rt.searchNode(word)
	if node != nil && remainingWord == "" && node.isEnd {
		node.isEnd = false
		node.count = 0
		rt.size--
		return true
	}

	return false
}

// GetAllWords returns all words in the radix tree
// Time Complexity: O(n*m) where n is number of words, m is average word length
// Space Complexity: O(n*m) for storing all words
func (rt *RadixTree) GetAllWords() []string {
	var words []string
	rt.dfs(rt.root, "", &words)
	return words
}

// GetWordsWithPrefix returns all words that start with the given prefix
// Time Complexity: O(p + n*m) where p is prefix length, n is number of matching words
// Space Complexity: O(n*m) for storing matching words
func (rt *RadixTree) GetWordsWithPrefix(prefix string) []string {
	var words []string
	prefixNode, remainingPrefix := rt.searchNode(prefix)

	if prefixNode == nil || remainingPrefix != "" {
		return words
	}

	rt.dfs(prefixNode, prefix, &words)
	return words
}

// dfs performs depth-first search to collect all words from a given node
func (rt *RadixTree) dfs(node *RadixNode, currentWord string, words *[]string) {
	if node.isEnd {
		*words = append(*words, currentWord)
	}

	for _, child := range node.children {
		rt.dfs(child, currentWord+child.label, words)
	}
}

// wordFreq holds word and its frequency for autocomplete functionality
type radixWordFreq struct {
	word string
	freq int
}

// AutoComplete returns top k words that start with the given prefix, sorted by frequency
// Time Complexity: O(p + n*m + n*log(n)) where p is prefix length, n is matching words
// Space Complexity: O(n*m) for storing matching words
func (rt *RadixTree) AutoComplete(prefix string, k int) []string {
	var results []radixWordFreq
	prefixNode, remainingPrefix := rt.searchNode(prefix)

	if prefixNode == nil || remainingPrefix != "" {
		return []string{}
	}

	rt.dfsWithFreq(prefixNode, prefix, &results)

	// Sort by frequency (descending) then by lexicographic order
	sort.Slice(results, func(i, j int) bool {
		if results[i].freq == results[j].freq {
			return results[i].word < results[j].word
		}
		return results[i].freq > results[j].freq
	})

	// Return top k results
	if k > len(results) {
		k = len(results)
	}

	words := make([]string, k)
	for i := 0; i < k; i++ {
		words[i] = results[i].word
	}

	return words
}

// dfsWithFreq performs DFS and collects words with their frequencies
func (rt *RadixTree) dfsWithFreq(node *RadixNode, currentWord string, results *[]radixWordFreq) {
	if node.isEnd {
		*results = append(*results, radixWordFreq{
			word: currentWord,
			freq: node.count,
		})
	}

	for _, child := range node.children {
		rt.dfsWithFreq(child, currentWord+child.label, results)
	}
}

// Size returns the number of unique words in the radix tree
func (rt *RadixTree) Size() int {
	return rt.size
}

// IsEmpty checks if the radix tree is empty
func (rt *RadixTree) IsEmpty() bool {
	return rt.size == 0
}

// GetWordCount returns the frequency count of a specific word
func (rt *RadixTree) GetWordCount(word string) int {
	node, remainingWord := rt.searchNode(word)
	if node != nil && remainingWord == "" && node.isEnd {
		return node.count
	}
	return 0
}

// CountWordsWithPrefix counts how many words start with the given prefix
// Time Complexity: O(p + n) where p is prefix length, n is number of nodes in subtree
// Space Complexity: O(1)
func (rt *RadixTree) CountWordsWithPrefix(prefix string) int {
	prefixNode, remainingPrefix := rt.searchNode(prefix)
	if prefixNode == nil || remainingPrefix != "" {
		return 0
	}

	return rt.countWords(prefixNode)
}

// countWords recursively counts all words in a subtree
func (rt *RadixTree) countWords(node *RadixNode) int {
	count := 0
	if node.isEnd {
		count = 1
	}

	for _, child := range node.children {
		count += rt.countWords(child)
	}

	return count
}

// Clear removes all words from the radix tree
func (rt *RadixTree) Clear() {
	rt.root = NewRadixNode("")
	rt.size = 0
}

// String returns a string representation of the radix tree
func (rt *RadixTree) String() string {
	words := rt.GetAllWords()
	return fmt.Sprintf("RadixTree{size: %d, words: %v}", rt.size, words)
}

// PrintTree prints the radix tree structure for debugging
func (rt *RadixTree) PrintTree() {
	fmt.Println("Radix Tree Structure:")
	rt.printNode(rt.root, "", "")
}

// printNode recursively prints the radix tree structure
func (rt *RadixTree) printNode(node *RadixNode, prefix, indent string) {
	if node.isEnd {
		fmt.Printf("%s%s%s (end, count: %d)\n", indent, prefix, node.label, node.count)
	} else if node.label != "" {
		fmt.Printf("%s%s%s\n", indent, prefix, node.label)
	}

	for _, child := range node.children {
		rt.printNode(child, "", indent+"  ")
	}
}

// GetMemoryEfficiency returns a comparison of memory usage vs regular trie
func (rt *RadixTree) GetMemoryEfficiency() (int, int) {
	radixNodes := rt.countNodes(rt.root)

	// Estimate equivalent trie nodes (each character would be a separate node)
	totalChars := 0
	words := rt.GetAllWords()
	for _, word := range words {
		totalChars += len(word)
	}

	// In a regular trie, we'd have roughly totalChars nodes
	return radixNodes, totalChars
}

// countNodes recursively counts the number of nodes in the tree
func (rt *RadixTree) countNodes(node *RadixNode) int {
	count := 1
	for _, child := range node.children {
		count += rt.countNodes(child)
	}
	return count
}
