package avl

import (
	"testing"
)

// TestNewAVLTree tests AVL tree creation
func TestNewAVLTree(t *testing.T) {
	avl := NewAVLTree()

	if avl.Root != nil {
		t.Error("New AVL tree should have nil root")
	}

	if avl.Size != 0 {
		t.Error("New AVL tree should have size 0")
	}

	if !avl.IsEmpty() {
		t.Error("New AVL tree should be empty")
	}
}

// TestInsertSingle tests inserting a single element
func TestInsertSingle(t *testing.T) {
	avl := NewAVLTree()
	avl.Insert(10)

	if avl.Root == nil {
		t.Error("Root should not be nil after insertion")
	}

	if avl.Root.Key != 10 {
		t.Errorf("Root key should be 10, got %d", avl.Root.Key)
	}

	if avl.Size != 1 {
		t.Errorf("Size should be 1, got %d", avl.Size)
	}

	if avl.GetHeight() != 1 {
		t.Errorf("Height should be 1, got %d", avl.GetHeight())
	}
}

// TestInsertMultiple tests inserting multiple elements
func TestInsertMultiple(t *testing.T) {
	avl := NewAVLTree()
	values := []int{10, 20, 30, 40, 50, 25}

	for _, val := range values {
		avl.Insert(val)
	}

	if avl.Size != len(values) {
		t.Errorf("Size should be %d, got %d", len(values), avl.Size)
	}

	// Check if tree is balanced
	if !avl.IsBalanced() {
		t.Error("Tree should be balanced after insertions")
	}

	// Check if all values are present
	for _, val := range values {
		if !avl.Search(val) {
			t.Errorf("Value %d should be found in tree", val)
		}
	}
}

// TestInsertDuplicates tests inserting duplicate values
func TestInsertDuplicates(t *testing.T) {
	avl := NewAVLTree()
	avl.Insert(10)
	avl.Insert(10) // Duplicate

	if avl.Size != 1 {
		t.Errorf("Size should be 1 after inserting duplicate, got %d", avl.Size)
	}
}

// TestLeftRotation tests left rotation scenario
func TestLeftRotation(t *testing.T) {
	avl := NewAVLTree()

	// Insert values that will cause left rotation
	avl.Insert(10)
	avl.Insert(20)
	avl.Insert(30) // This should trigger left rotation

	// After rotation, root should be 20
	if avl.Root.Key != 20 {
		t.Errorf("Root should be 20 after left rotation, got %d", avl.Root.Key)
	}

	if avl.Root.Left.Key != 10 {
		t.Errorf("Left child of root should be 10, got %d", avl.Root.Left.Key)
	}

	if avl.Root.Right.Key != 30 {
		t.Errorf("Right child of root should be 30, got %d", avl.Root.Right.Key)
	}

	if !avl.IsBalanced() {
		t.Error("Tree should be balanced after rotation")
	}
}

// TestRightRotation tests right rotation scenario
func TestRightRotation(t *testing.T) {
	avl := NewAVLTree()

	// Insert values that will cause right rotation
	avl.Insert(30)
	avl.Insert(20)
	avl.Insert(10) // This should trigger right rotation

	// After rotation, root should be 20
	if avl.Root.Key != 20 {
		t.Errorf("Root should be 20 after right rotation, got %d", avl.Root.Key)
	}

	if avl.Root.Left.Key != 10 {
		t.Errorf("Left child of root should be 10, got %d", avl.Root.Left.Key)
	}

	if avl.Root.Right.Key != 30 {
		t.Errorf("Right child of root should be 30, got %d", avl.Root.Right.Key)
	}

	if !avl.IsBalanced() {
		t.Error("Tree should be balanced after rotation")
	}
}

// TestLeftRightRotation tests left-right rotation scenario
func TestLeftRightRotation(t *testing.T) {
	avl := NewAVLTree()

	// Insert values that will cause left-right rotation
	avl.Insert(30)
	avl.Insert(10)
	avl.Insert(20) // This should trigger left-right rotation

	// After rotation, root should be 20
	if avl.Root.Key != 20 {
		t.Errorf("Root should be 20 after left-right rotation, got %d", avl.Root.Key)
	}

	if !avl.IsBalanced() {
		t.Error("Tree should be balanced after rotation")
	}
}

// TestRightLeftRotation tests right-left rotation scenario
func TestRightLeftRotation(t *testing.T) {
	avl := NewAVLTree()

	// Insert values that will cause right-left rotation
	avl.Insert(10)
	avl.Insert(30)
	avl.Insert(20) // This should trigger right-left rotation

	// After rotation, root should be 20
	if avl.Root.Key != 20 {
		t.Errorf("Root should be 20 after right-left rotation, got %d", avl.Root.Key)
	}

	if !avl.IsBalanced() {
		t.Error("Tree should be balanced after rotation")
	}
}

// TestSearch tests search functionality
func TestSearch(t *testing.T) {
	avl := NewAVLTree()
	values := []int{50, 30, 70, 20, 40, 60, 80}

	for _, val := range values {
		avl.Insert(val)
	}

	// Test existing values
	for _, val := range values {
		if !avl.Search(val) {
			t.Errorf("Value %d should be found", val)
		}
	}

	// Test non-existing values
	nonExisting := []int{10, 25, 35, 45, 90}
	for _, val := range nonExisting {
		if avl.Search(val) {
			t.Errorf("Value %d should not be found", val)
		}
	}
}

// TestDelete tests deletion functionality
func TestDelete(t *testing.T) {
	avl := NewAVLTree()
	values := []int{50, 30, 70, 20, 40, 60, 80}

	for _, val := range values {
		avl.Insert(val)
	}

	initialSize := avl.Size

	// Delete a leaf node
	avl.Delete(20)
	if avl.Search(20) {
		t.Error("Value 20 should be deleted")
	}
	if avl.Size != initialSize-1 {
		t.Errorf("Size should be %d after deletion, got %d", initialSize-1, avl.Size)
	}
	if !avl.IsBalanced() {
		t.Error("Tree should remain balanced after deletion")
	}

	// Delete a node with one child
	avl.Delete(30)
	if avl.Search(30) {
		t.Error("Value 30 should be deleted")
	}
	if !avl.IsBalanced() {
		t.Error("Tree should remain balanced after deletion")
	}

	// Delete a node with two children
	avl.Delete(50)
	if avl.Search(50) {
		t.Error("Value 50 should be deleted")
	}
	if !avl.IsBalanced() {
		t.Error("Tree should remain balanced after deletion")
	}
}

// TestDeleteNonExisting tests deleting non-existing values
func TestDeleteNonExisting(t *testing.T) {
	avl := NewAVLTree()
	avl.Insert(10)

	initialSize := avl.Size
	avl.Delete(20) // Non-existing value

	if avl.Size != initialSize {
		t.Error("Size should not change when deleting non-existing value")
	}
}

// TestInorderTraversal tests inorder traversal
func TestInorderTraversal(t *testing.T) {
	avl := NewAVLTree()
	values := []int{50, 30, 70, 20, 40, 60, 80}
	expected := []int{20, 30, 40, 50, 60, 70, 80}

	for _, val := range values {
		avl.Insert(val)
	}

	result := avl.InorderTraversal()

	if len(result) != len(expected) {
		t.Errorf("Inorder result length should be %d, got %d", len(expected), len(result))
	}

	for i, val := range result {
		if val != expected[i] {
			t.Errorf("Inorder result[%d] should be %d, got %d", i, expected[i], val)
		}
	}
}

// TestComplexOperations tests a complex sequence of operations
func TestComplexOperations(t *testing.T) {
	avl := NewAVLTree()

	// Insert a large number of values
	values := []int{50, 30, 70, 20, 40, 60, 80, 10, 25, 35, 45, 55, 65, 75, 85}
	for _, val := range values {
		avl.Insert(val)
	}

	// Verify tree is balanced
	if !avl.IsBalanced() {
		t.Error("Tree should be balanced after complex insertions")
	}

	// Delete some values
	toDelete := []int{10, 25, 75, 85}
	for _, val := range toDelete {
		avl.Delete(val)
	}

	// Verify tree is still balanced
	if !avl.IsBalanced() {
		t.Error("Tree should remain balanced after deletions")
	}

	// Verify remaining values are still searchable
	remaining := []int{50, 30, 70, 20, 40, 60, 80, 35, 45, 55, 65}
	for _, val := range remaining {
		if !avl.Search(val) {
			t.Errorf("Remaining value %d should be found", val)
		}
	}

	// Verify deleted values are not found
	for _, val := range toDelete {
		if avl.Search(val) {
			t.Errorf("Deleted value %d should not be found", val)
		}
	}
}

// TestEmptyTreeOperations tests operations on empty tree
func TestEmptyTreeOperations(t *testing.T) {
	avl := NewAVLTree()

	if avl.Search(10) {
		t.Error("Search on empty tree should return false")
	}

	avl.Delete(10) // Should not crash

	if avl.GetHeight() != 0 {
		t.Error("Empty tree height should be 0")
	}

	result := avl.InorderTraversal()
	if len(result) != 0 {
		t.Error("Inorder traversal of empty tree should be empty")
	}
}

// BenchmarkInsert benchmarks insertion performance
func BenchmarkInsert(b *testing.B) {
	avl := NewAVLTree()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		avl.Insert(i)
	}
}

// BenchmarkSearch benchmarks search performance
func BenchmarkSearch(b *testing.B) {
	avl := NewAVLTree()

	// Prepare tree with 1000 elements
	for i := 0; i < 1000; i++ {
		avl.Insert(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		avl.Search(i % 1000)
	}
}

// BenchmarkDelete benchmarks deletion performance
func BenchmarkDelete(b *testing.B) {
	// Prepare multiple trees for deletion
	trees := make([]*AVLTree, b.N)
	for i := 0; i < b.N; i++ {
		trees[i] = NewAVLTree()
		for j := 0; j < 100; j++ {
			trees[i].Insert(j)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trees[i].Delete(i % 100)
	}
}
