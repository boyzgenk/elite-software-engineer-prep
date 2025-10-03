package main

import (
	"fmt"
	"sort"
	"strings"
)

// TrieNode represents a single node in the Trie
type TrieNode struct {
	children map[rune]*TrieNode // Using rune for Unicode support
	isEnd    bool               // Marks the end of a word
	count    int                // Number of words ending at this node (for frequency tracking)
}

// NewTrieNode creates a new TrieNode
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEnd:    false,
		count:    0,
	}
}

// Trie represents the Trie data structure
type Trie struct {
	root *TrieNode
	size int // Total number of unique words in the trie
}

// NewTrie creates a new Trie
func NewTrie() *Trie {
	return &Trie{
		root: NewTrieNode(),
		size: 0,
	}
}

// Insert adds a word to the trie
// Time Complexity: O(m) where m is the length of the word
// Space Complexity: O(m) in worst case (no common prefixes)
func (t *Trie) Insert(word string) {
	if word == "" {
		return
	}

	current := t.root
	for _, char := range word {
		if current.children[char] == nil {
			current.children[char] = NewTrieNode()
		}
		current = current.children[char]
	}

	if !current.isEnd {
		current.isEnd = true
		t.size++
	}
	current.count++
}

// Search checks if a word exists in the trie
// Time Complexity: O(m) where m is the length of the word
// Space Complexity: O(1)
func (t *Trie) Search(word string) bool {
	node := t.searchNode(word)
	return node != nil && node.isEnd
}

// StartsWith checks if there are any words starting with the given prefix
// Time Complexity: O(m) where m is the length of the prefix
// Space Complexity: O(1)
func (t *Trie) StartsWith(prefix string) bool {
	return t.searchNode(prefix) != nil
}

// searchNode is a helper function that finds the node corresponding to a word/prefix
func (t *Trie) searchNode(word string) *TrieNode {
	current := t.root
	for _, char := range word {
		if current.children[char] == nil {
			return nil
		}
		current = current.children[char]
	}
	return current
}

// Delete removes a word from the trie
// Time Complexity: O(m) where m is the length of the word
// Space Complexity: O(m) for recursion stack
func (t *Trie) Delete(word string) bool {
	if word == "" {
		return false
	}

	deleted := t.deleteHelper(t.root, word, 0)
	return deleted
}

// deleteHelper recursively deletes a word from the trie
func (t *Trie) deleteHelper(node *TrieNode, word string, index int) bool {
	if index == len([]rune(word)) {
		// Reached end of word
		if !node.isEnd {
			return false // Word doesn't exist
		}

		node.isEnd = false
		node.count = 0
		t.size-- // Decrease size when we actually delete a word

		// Return true to indicate word was successfully deleted
		// The parent will decide whether to delete this node based on its state
		return true
	}

	chars := []rune(word)
	char := chars[index]
	childNode := node.children[char]

	if childNode == nil {
		return false // Word doesn't exist
	}

	wordDeleted := t.deleteHelper(childNode, word, index+1)

	if wordDeleted {
		// If child node is now empty (not end of word and no children), remove it
		if !childNode.isEnd && len(childNode.children) == 0 {
			delete(node.children, char)
		}
		return true // Propagate that word was deleted
	}

	return false
}

// GetAllWords returns all words in the trie
// Time Complexity: O(n*m) where n is number of words, m is average word length
// Space Complexity: O(n*m) for storing all words
func (t *Trie) GetAllWords() []string {
	var words []string
	t.dfs(t.root, "", &words)
	return words
}

// GetWordsWithPrefix returns all words that start with the given prefix
// Time Complexity: O(p + n*m) where p is prefix length, n is number of matching words
// Space Complexity: O(n*m) for storing matching words
func (t *Trie) GetWordsWithPrefix(prefix string) []string {
	var words []string
	prefixNode := t.searchNode(prefix)

	if prefixNode == nil {
		return words
	}

	t.dfs(prefixNode, prefix, &words)
	return words
}

// dfs performs depth-first search to collect all words from a given node
func (t *Trie) dfs(node *TrieNode, currentWord string, words *[]string) {
	if node.isEnd {
		*words = append(*words, currentWord)
	}

	for char, childNode := range node.children {
		t.dfs(childNode, currentWord+string(char), words)
	}
}

// wordFreq holds word and its frequency for autocomplete functionality
type wordFreq struct {
	word string
	freq int
}

// AutoComplete returns top k words that start with the given prefix, sorted by frequency
// Time Complexity: O(p + n*m + n*log(n)) where p is prefix length, n is matching words
// Space Complexity: O(n*m) for storing matching words
func (t *Trie) AutoComplete(prefix string, k int) []string {
	var results []wordFreq
	prefixNode := t.searchNode(prefix)

	if prefixNode == nil {
		return []string{}
	}

	t.dfsWithFreq(prefixNode, prefix, &results)

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
func (t *Trie) dfsWithFreq(node *TrieNode, currentWord string, results *[]wordFreq) {
	if node.isEnd {
		*results = append(*results, wordFreq{
			word: currentWord,
			freq: node.count,
		})
	}

	for char, childNode := range node.children {
		t.dfsWithFreq(childNode, currentWord+string(char), results)
	}
}

// Size returns the number of unique words in the trie
func (t *Trie) Size() int {
	return t.size
}

// IsEmpty checks if the trie is empty
func (t *Trie) IsEmpty() bool {
	return t.size == 0
}

// GetWordCount returns the frequency count of a specific word
func (t *Trie) GetWordCount(word string) int {
	node := t.searchNode(word)
	if node != nil && node.isEnd {
		return node.count
	}
	return 0
}

// LongestCommonPrefix finds the longest common prefix among all words in the trie
// Time Complexity: O(m) where m is the length of the longest common prefix
// Space Complexity: O(m) for the result string
func (t *Trie) LongestCommonPrefix() string {
	if t.size == 0 {
		return ""
	}

	var prefix strings.Builder
	current := t.root

	for len(current.children) == 1 && !current.isEnd {
		for char, child := range current.children {
			prefix.WriteRune(char)
			current = child
			break
		}
	}

	return prefix.String()
}

// CountWordsWithPrefix counts how many words start with the given prefix
// Time Complexity: O(p + n) where p is prefix length, n is number of nodes in subtree
// Space Complexity: O(1)
func (t *Trie) CountWordsWithPrefix(prefix string) int {
	prefixNode := t.searchNode(prefix)
	if prefixNode == nil {
		return 0
	}

	return t.countWords(prefixNode)
}

// countWords recursively counts all words in a subtree
func (t *Trie) countWords(node *TrieNode) int {
	count := 0
	if node.isEnd {
		count = 1
	}

	for _, child := range node.children {
		count += t.countWords(child)
	}

	return count
}

// Clear removes all words from the trie
func (t *Trie) Clear() {
	t.root = NewTrieNode()
	t.size = 0
}

// String returns a string representation of the trie
func (t *Trie) String() string {
	words := t.GetAllWords()
	return fmt.Sprintf("Trie{size: %d, words: %v}", t.size, words)
}

// PrintTrie prints the trie structure for debugging
func (t *Trie) PrintTrie() {
	fmt.Println("Trie Structure:")
	t.printNode(t.root, "", "")
}

// printNode recursively prints the trie structure
func (t *Trie) printNode(node *TrieNode, prefix, indent string) {
	if node.isEnd {
		fmt.Printf("%s%s (end, count: %d)\n", indent, prefix, node.count)
	}

	for char, child := range node.children {
		t.printNode(child, string(char), indent+"  ")
	}
}
