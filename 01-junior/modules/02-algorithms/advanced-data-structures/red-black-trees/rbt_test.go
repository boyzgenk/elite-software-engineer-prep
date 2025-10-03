package redblacktree

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

// TestNewRedBlackTree tests tree creation
func TestNewRedBlackTree(t *testing.T) {
	rbt := NewRedBlackTree()

	if rbt.Root != nil {
		t.Error("New tree should have nil root")
	}
	if rbt.Size != 0 {
		t.Error("New tree should have size 0")
	}
	if !rbt.IsValid() {
		t.Error("New tree should be valid")
	}
}

// TestSingleInsertion tests inserting a single element
func TestSingleInsertion(t *testing.T) {
	rbt := NewRedBlackTree()
	rbt.Insert(10, "ten")

	if rbt.Root == nil {
		t.Fatal("Root should not be nil after insertion")
	}
	if rbt.Root.Key != 10 {
		t.Error("Root key should be 10")
	}
	if rbt.Root.Color != Black {
		t.Error("Root should be black")
	}
	if rbt.Size != 1 {
		t.Error("Size should be 1")
	}
	if !rbt.IsValid() {
		t.Error("Tree should be valid after single insertion")
	}
}

// TestMultipleInsertions tests inserting multiple elements
func TestMultipleInsertions(t *testing.T) {
	rbt := NewRedBlackTree()
	keys := []int{10, 5, 15, 3, 7, 12, 18, 1, 4, 6, 8, 11, 13, 16, 20}

	for i, key := range keys {
		rbt.Insert(key, fmt.Sprintf("value_%d", key))

		if rbt.Size != i+1 {
			t.Errorf("Size should be %d after inserting %d elements", i+1, i+1)
		}
		if !rbt.IsValid() {
			t.Errorf("Tree should be valid after inserting key %d", key)
			t.Logf("Tree state:\n%s", rbt.String())
		}
	}

	// Test all keys are present
	for _, key := range keys {
		if !rbt.Contains(key) {
			t.Errorf("Tree should contain key %d", key)
		}
		value, found := rbt.Get(key)
		if !found {
			t.Errorf("Should find key %d", key)
		}
		expectedValue := fmt.Sprintf("value_%d", key)
		if value != expectedValue {
			t.Errorf("Expected value %s, got %s", expectedValue, value)
		}
	}
}

// TestInsertionBalance tests that insertions maintain balance
func TestInsertionBalance(t *testing.T) {
	rbt := NewRedBlackTree()

	// Insert in ascending order (worst case for unbalanced BST)
	for i := 1; i <= 15; i++ {
		rbt.Insert(i, i)
		if !rbt.IsValid() {
			t.Errorf("Tree invalid after inserting %d", i)
		}
	}

	// Height should be logarithmic
	height := rbt.Height()
	maxHeight := 2 * rbt.BlackHeight() // Red-Black Tree height property
	if height > maxHeight {
		t.Errorf("Height %d exceeds maximum expected %d", height, maxHeight)
	}

	// Test inorder traversal gives sorted order
	keys := rbt.InorderTraversal()
	if !sort.IntsAreSorted(keys) {
		t.Error("Inorder traversal should give sorted keys")
	}
	if len(keys) != 15 {
		t.Errorf("Expected 15 keys, got %d", len(keys))
	}
}

// TestDuplicateInsertion tests inserting duplicate keys
func TestDuplicateInsertion(t *testing.T) {
	rbt := NewRedBlackTree()

	rbt.Insert(10, "first")
	originalSize := rbt.Size

	rbt.Insert(10, "second")

	if rbt.Size != originalSize {
		t.Error("Size should not change when inserting duplicate key")
	}

	value, found := rbt.Get(10)
	if !found {
		t.Error("Should find key 10")
	}
	if value != "second" {
		t.Error("Value should be updated to 'second'")
	}

	if !rbt.IsValid() {
		t.Error("Tree should remain valid after duplicate insertion")
	}
}

// TestDeletion tests basic deletion operations
func TestDeletion(t *testing.T) {
	rbt := NewRedBlackTree()
	keys := []int{10, 5, 15, 3, 7, 12, 18, 1, 4, 6, 8, 11, 13, 16, 20}

	// Insert all keys
	for _, key := range keys {
		rbt.Insert(key, key)
	}

	// Delete some keys
	toDelete := []int{1, 4, 6, 11, 16, 20}
	for _, key := range toDelete {
		initialSize := rbt.Size
		deleted := rbt.Delete(key)

		if !deleted {
			t.Errorf("Should successfully delete key %d", key)
		}
		if rbt.Size != initialSize-1 {
			t.Errorf("Size should decrease by 1 after deleting key %d", key)
		}
		if rbt.Contains(key) {
			t.Errorf("Tree should not contain deleted key %d", key)
		}
		if !rbt.IsValid() {
			t.Errorf("Tree should be valid after deleting key %d", key)
			t.Logf("Tree state:\n%s", rbt.String())
		}
	}

	// Test remaining keys are still present
	remaining := []int{10, 5, 15, 3, 7, 12, 18, 8, 13}
	for _, key := range remaining {
		if !rbt.Contains(key) {
			t.Errorf("Tree should still contain key %d", key)
		}
	}
}

// TestDeleteNonExistent tests deleting non-existent keys
func TestDeleteNonExistent(t *testing.T) {
	rbt := NewRedBlackTree()
	rbt.Insert(10, "ten")

	originalSize := rbt.Size
	deleted := rbt.Delete(20)

	if deleted {
		t.Error("Should not delete non-existent key")
	}
	if rbt.Size != originalSize {
		t.Error("Size should not change when deleting non-existent key")
	}
	if !rbt.IsValid() {
		t.Error("Tree should remain valid after failed deletion")
	}
}

// TestDeleteRoot tests deleting root node in various scenarios
func TestDeleteRoot(t *testing.T) {
	// Test 1: Delete root with no children
	rbt := NewRedBlackTree()
	rbt.Insert(10, "ten")
	rbt.Delete(10)

	if rbt.Root != nil {
		t.Error("Root should be nil after deleting only node")
	}
	if rbt.Size != 0 {
		t.Error("Size should be 0 after deleting only node")
	}
	if !rbt.IsValid() {
		t.Error("Empty tree should be valid")
	}

	// Test 2: Delete root with children
	rbt = NewRedBlackTree()
	keys := []int{10, 5, 15, 3, 7, 12, 18}
	for _, key := range keys {
		rbt.Insert(key, key)
	}

	rbt.Delete(10) // Delete root

	if rbt.Root == nil {
		t.Error("Root should not be nil after deleting root with children")
	}
	if rbt.Contains(10) {
		t.Error("Tree should not contain deleted root")
	}
	if !rbt.IsValid() {
		t.Error("Tree should be valid after deleting root")
	}
}

// TestComplexDeletionScenarios tests various deletion edge cases
func TestComplexDeletionScenarios(t *testing.T) {
	rbt := NewRedBlackTree()

	// Create a complex tree structure
	keys := []int{50, 25, 75, 10, 30, 60, 80, 5, 15, 27, 35, 55, 65, 77, 85}
	for _, key := range keys {
		rbt.Insert(key, key)
	}

	// Delete nodes that will trigger different rebalancing cases
	deletionOrder := []int{5, 15, 85, 77, 65, 35, 27, 30, 10, 25}

	for _, key := range deletionOrder {
		rbt.Delete(key)
		if !rbt.IsValid() {
			t.Errorf("Tree invalid after deleting %d", key)
			t.Logf("Tree state:\n%s", rbt.String())
			break
		}
	}
}

// TestMinMax tests finding minimum and maximum elements
func TestMinMax(t *testing.T) {
	rbt := NewRedBlackTree()

	// Test empty tree
	if rbt.Min() != nil {
		t.Error("Min of empty tree should be nil")
	}
	if rbt.Max() != nil {
		t.Error("Max of empty tree should be nil")
	}

	// Insert elements
	keys := []int{10, 5, 15, 3, 7, 12, 18, 1, 20}
	for _, key := range keys {
		rbt.Insert(key, key)
	}

	min := rbt.Min()
	if min == nil || min.Key != 1 {
		t.Error("Min should be 1")
	}

	max := rbt.Max()
	if max == nil || max.Key != 20 {
		t.Error("Max should be 20")
	}

	// Delete min and max, test again
	rbt.Delete(1)
	min = rbt.Min()
	if min == nil || min.Key != 3 {
		t.Error("Min should be 3 after deleting 1")
	}

	rbt.Delete(20)
	max = rbt.Max()
	if max == nil || max.Key != 18 {
		t.Error("Max should be 18 after deleting 20")
	}
}

// TestHeightAndBlackHeight tests height calculations
func TestHeightAndBlackHeight(t *testing.T) {
	rbt := NewRedBlackTree()

	// Empty tree
	if rbt.Height() != 0 {
		t.Error("Height of empty tree should be 0")
	}
	if rbt.BlackHeight() != 1 {
		t.Error("Black height of empty tree should be 1")
	}

	// Single node
	rbt.Insert(10, "ten")
	if rbt.Height() != 1 {
		t.Error("Height of single-node tree should be 1")
	}
	if rbt.BlackHeight() != 2 {
		t.Error("Black height of single-node tree should be 2")
	}

	// Multiple nodes
	keys := []int{5, 15, 3, 7, 12, 18}
	for _, key := range keys {
		rbt.Insert(key, key)
	}

	height := rbt.Height()
	blackHeight := rbt.BlackHeight()

	// Red-Black Tree height property: height ≤ 2 * blackHeight - 1
	if height > 2*blackHeight-1 {
		t.Errorf("Height %d violates Red-Black property with black height %d", height, blackHeight)
	}
}

// TestRedBlackProperties tests all Red-Black Tree properties
func TestRedBlackProperties(t *testing.T) {
	rbt := NewRedBlackTree()

	// Test with various insertion patterns
	patterns := [][]int{
		{1, 2, 3, 4, 5, 6, 7}, // Ascending
		{7, 6, 5, 4, 3, 2, 1}, // Descending
		{4, 2, 6, 1, 3, 5, 7}, // Balanced
		{1, 3, 2, 5, 4, 7, 6}, // Random-ish
	}

	for i, pattern := range patterns {
		rbt = NewRedBlackTree()

		for _, key := range pattern {
			rbt.Insert(key, key)
			if !rbt.IsValid() {
				t.Errorf("Pattern %d: Tree invalid after inserting %d", i, key)
				break
			}
		}

		// Test after full insertion
		if !rbt.IsValid() {
			t.Errorf("Pattern %d: Tree invalid after complete insertion", i)
		}

		// Test after some deletions
		for j, key := range pattern[:len(pattern)/2] {
			rbt.Delete(key)
			if !rbt.IsValid() {
				t.Errorf("Pattern %d: Tree invalid after deleting %d (step %d)", i, key, j)
				break
			}
		}
	}
}

// TestLargeDataset tests performance with larger datasets
func TestLargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	rbt := NewRedBlackTree()
	n := 10000

	// Generate random keys
	rand.Seed(time.Now().UnixNano())
	keys := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = rand.Intn(n * 2)
	}

	// Insert all keys
	start := time.Now()
	for _, key := range keys {
		rbt.Insert(key, key)
	}
	insertTime := time.Since(start)

	t.Logf("Inserted %d elements in %v", n, insertTime)
	t.Logf("Tree size: %d, Height: %d, Black Height: %d",
		rbt.Size, rbt.Height(), rbt.BlackHeight())

	if !rbt.IsValid() {
		t.Error("Tree should be valid after large insertions")
	}

	// Verify height is logarithmic
	maxHeight := 2 * rbt.BlackHeight()
	if rbt.Height() > maxHeight {
		t.Errorf("Height %d exceeds maximum expected %d", rbt.Height(), maxHeight)
	}

	// Delete half the elements
	start = time.Now()
	for i := 0; i < len(keys)/2; i++ {
		rbt.Delete(keys[i])
	}
	deleteTime := time.Since(start)

	t.Logf("Deleted %d elements in %v", len(keys)/2, deleteTime)

	if !rbt.IsValid() {
		t.Error("Tree should be valid after large deletions")
	}
}

// Benchmark tests
func BenchmarkInsert(b *testing.B) {
	rbt := NewRedBlackTree()
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = rand.Intn(b.N * 2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbt.Insert(keys[i], keys[i])
	}
}

func BenchmarkSearch(b *testing.B) {
	rbt := NewRedBlackTree()
	n := 10000

	// Pre-populate tree
	for i := 0; i < n; i++ {
		rbt.Insert(rand.Intn(n*2), i)
	}

	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = rand.Intn(n * 2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbt.Search(keys[i])
	}
}

func BenchmarkDelete(b *testing.B) {
	// Pre-populate tree
	rbt := NewRedBlackTree()
	keys := make([]int, b.N*2)
	for i := 0; i < b.N*2; i++ {
		keys[i] = rand.Intn(b.N * 4)
		rbt.Insert(keys[i], keys[i])
	}

	deleteKeys := keys[:b.N]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbt.Delete(deleteKeys[i])
	}
}

// BenchmarkRBTvsAVL compares Red-Black Tree with AVL Tree performance
func BenchmarkRBTvsAVL(b *testing.B) {
	// This would require importing the AVL implementation
	// For now, we'll just benchmark RBT operations

	b.Run("RBT-Insert", func(b *testing.B) {
		rbt := NewRedBlackTree()
		for i := 0; i < b.N; i++ {
			rbt.Insert(rand.Intn(b.N*2), i)
		}
	})

	b.Run("RBT-Search", func(b *testing.B) {
		rbt := NewRedBlackTree()
		// Pre-populate
		for i := 0; i < 10000; i++ {
			rbt.Insert(rand.Intn(20000), i)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rbt.Search(rand.Intn(20000))
		}
	})

	b.Run("RBT-Delete", func(b *testing.B) {
		rbt := NewRedBlackTree()
		keys := make([]int, b.N*2)

		// Pre-populate
		for i := 0; i < b.N*2; i++ {
			keys[i] = rand.Intn(b.N * 4)
			rbt.Insert(keys[i], i)
		}

		b.ResetTimer()
		for i := 0; i < b.N && i < len(keys); i++ {
			rbt.Delete(keys[i])
		}
	})
}
