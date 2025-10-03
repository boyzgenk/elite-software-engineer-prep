package segmenttree

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// TestNewSegmentTree tests basic tree creation
func TestNewSegmentTree(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	// Test sum segment tree
	st := NewSegmentTree(arr, Sum)
	if st == nil {
		t.Fatal("Failed to create segment tree")
	}
	if st.Size() != 5 {
		t.Errorf("Expected size 5, got %d", st.Size())
	}
	if st.Operation() != Sum {
		t.Errorf("Expected Sum operation, got %s", st.Operation())
	}

	// Test empty array
	emptySt := NewSegmentTree([]int{}, Sum)
	if emptySt != nil {
		t.Error("Empty array should return nil segment tree")
	}
}

// TestSumSegmentTree tests range sum queries
func TestSumSegmentTree(t *testing.T) {
	arr := []int{1, 3, 5, 7, 9, 11}
	st := NewSumSegmentTree(arr)

	// Test individual elements
	for i, expected := range arr {
		result := st.RangeSum(i, i)
		if result != expected {
			t.Errorf("Query(%d,%d): expected %d, got %d", i, i, expected, result)
		}
	}

	// Test range queries
	testCases := []struct {
		l, r, expected int
	}{
		{0, 0, 1},  // Single element
		{0, 1, 4},  // First two: 1+3
		{1, 3, 15}, // Middle range: 3+5+7
		{0, 5, 36}, // Entire array: 1+3+5+7+9+11
		{2, 4, 21}, // 5+7+9
		{3, 5, 27}, // 7+9+11
	}

	for _, tc := range testCases {
		result := st.RangeSum(tc.l, tc.r)
		if result != tc.expected {
			t.Errorf("RangeSum(%d,%d): expected %d, got %d",
				tc.l, tc.r, tc.expected, result)
		}
	}
}

// TestMinSegmentTree tests range minimum queries
func TestMinSegmentTree(t *testing.T) {
	arr := []int{4, 2, 8, 1, 9, 3, 6}
	st := NewMinSegmentTree(arr)

	testCases := []struct {
		l, r, expected int
	}{
		{0, 0, 4}, // Single element
		{0, 1, 2}, // min(4,2)
		{1, 3, 1}, // min(2,8,1)
		{0, 6, 1}, // Entire array minimum
		{2, 4, 1}, // min(8,1,9)
		{4, 6, 3}, // min(9,3,6)
		{5, 6, 3}, // min(3,6)
	}

	for _, tc := range testCases {
		result := st.RangeMin(tc.l, tc.r)
		if result != tc.expected {
			t.Errorf("RangeMin(%d,%d): expected %d, got %d",
				tc.l, tc.r, tc.expected, result)
		}
	}
}

// TestMaxSegmentTree tests range maximum queries
func TestMaxSegmentTree(t *testing.T) {
	arr := []int{4, 2, 8, 1, 9, 3, 6}
	st := NewMaxSegmentTree(arr)

	testCases := []struct {
		l, r, expected int
	}{
		{0, 0, 4}, // Single element
		{0, 1, 4}, // max(4,2)
		{1, 3, 8}, // max(2,8,1)
		{0, 6, 9}, // Entire array maximum
		{2, 4, 9}, // max(8,1,9)
		{4, 6, 9}, // max(9,3,6)
		{5, 6, 6}, // max(3,6)
	}

	for _, tc := range testCases {
		result := st.RangeMax(tc.l, tc.r)
		if result != tc.expected {
			t.Errorf("RangeMax(%d,%d): expected %d, got %d",
				tc.l, tc.r, tc.expected, result)
		}
	}
}

// TestPointUpdate tests single element updates
func TestPointUpdate(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	st := NewSumSegmentTree(arr)

	// Update middle element
	st.Update(2, 10)
	expected := []int{1, 2, 10, 4, 5}
	result := st.GetArray()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After update: expected %v, got %v", expected, result)
	}

	// Test range sum after update
	totalSum := st.RangeSum(0, 4)
	expectedSum := 22 // 1+2+10+4+5
	if totalSum != expectedSum {
		t.Errorf("Sum after update: expected %d, got %d", expectedSum, totalSum)
	}

	// Update first element
	st.Update(0, 100)
	firstThreeSum := st.RangeSum(0, 2)
	expectedFirstThree := 112 // 100+2+10
	if firstThreeSum != expectedFirstThree {
		t.Errorf("First three sum: expected %d, got %d", expectedFirstThree, firstThreeSum)
	}
}

// TestRangeUpdate tests range updates with lazy propagation
func TestRangeUpdate(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	st := NewSumSegmentTree(arr)

	// Add 10 to range [1, 3]
	st.RangeUpdate(1, 3, 10)

	// Check individual elements
	expected := []int{1, 12, 13, 14, 5}
	result := st.GetArray()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After range update: expected %v, got %v", expected, result)
	}

	// Test range sum after update
	updatedRangeSum := st.RangeSum(1, 3)
	expectedRangeSum := 39 // 12+13+14
	if updatedRangeSum != expectedRangeSum {
		t.Errorf("Updated range sum: expected %d, got %d", expectedRangeSum, updatedRangeSum)
	}

	// Multiple range updates
	st.RangeUpdate(0, 4, 5) // Add 5 to entire array
	finalSum := st.RangeSum(0, 4)
	expectedFinalSum := 70 // (1+12+13+14+5) + 5*5 = 45 + 25
	if finalSum != expectedFinalSum {
		t.Errorf("Final sum: expected %d, got %d", expectedFinalSum, finalSum)
	}
}

// TestCustomSegmentTree tests custom combine functions
func TestCustomSegmentTree(t *testing.T) {
	// Test GCD segment tree
	arr := []int{12, 18, 24, 36}
	gcdSt := NewGCDSegmentTree(arr)

	// GCD of entire array should be 6
	result := gcdSt.Query(0, 3)
	expected := 6
	if result != expected {
		t.Errorf("GCD of entire array: expected %d, got %d", expected, result)
	}

	// GCD of [12, 18] should be 6
	result = gcdSt.Query(0, 1)
	expected = 6
	if result != expected {
		t.Errorf("GCD of [12,18]: expected %d, got %d", expected, result)
	}

	// Test LCM segment tree
	arr2 := []int{4, 6, 8, 12}
	lcmSt := NewLCMSegmentTree(arr2)

	// LCM of [4, 6] should be 12
	result = lcmSt.Query(0, 1)
	expected = 12
	if result != expected {
		t.Errorf("LCM of [4,6]: expected %d, got %d", expected, result)
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	// Single element array
	arr := []int{42}
	st := NewSumSegmentTree(arr)

	if st.RangeSum(0, 0) != 42 {
		t.Error("Single element query failed")
	}

	st.Update(0, 100)
	if st.RangeSum(0, 0) != 100 {
		t.Error("Single element update failed")
	}

	// Invalid range queries
	arr2 := []int{1, 2, 3, 4, 5}
	st2 := NewSumSegmentTree(arr2)

	// Out of bounds queries should be clamped
	if st2.Query(-1, 0) != 1 {
		t.Error("Negative index should be clamped to valid range")
	}
	if st2.Query(0, 10) != st2.Query(0, 4) {
		t.Error("Out of bounds end should be clamped")
	}
	if st2.Query(3, 2) != 0 {
		t.Error("Invalid range (l > r) should return identity")
	}

	// Edge case updates
	st2.Update(-1, 100) // Should be ignored
	st2.Update(10, 100) // Should be ignored

	originalSum := st2.RangeSum(0, 4)
	expectedSum := 15 // 1+2+3+4+5
	if originalSum != expectedSum {
		t.Error("Invalid updates should not affect tree")
	}
}

// TestLazyPropagation tests complex lazy propagation scenarios
func TestLazyPropagation(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8}
	st := NewSumSegmentTree(arr)

	// Multiple overlapping range updates
	st.RangeUpdate(1, 5, 10) // Add 10 to indices 1-5
	st.RangeUpdate(3, 7, 5)  // Add 5 to indices 3-7
	st.RangeUpdate(0, 2, 3)  // Add 3 to indices 0-2

	// Test specific positions
	testCases := []struct {
		pos, expected int
	}{
		{0, 4},  // 1 + 3 = 4
		{1, 15}, // 2 + 10 + 3 = 15
		{2, 16}, // 3 + 10 + 3 = 16
		{3, 19}, // 4 + 10 + 5 = 19
		{4, 20}, // 5 + 10 + 5 = 20
		{5, 21}, // 6 + 10 + 5 = 21
		{6, 12}, // 7 + 5 = 12
		{7, 13}, // 8 + 5 = 13
	}

	for _, tc := range testCases {
		result := st.Query(tc.pos, tc.pos)
		if result != tc.expected {
			t.Errorf("Position %d: expected %d, got %d", tc.pos, tc.expected, result)
		}
	}
}

// TestValidation tests tree structure validation
func TestValidation(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8}

	// Test all operation types
	operations := []Operation{Sum, Min, Max}
	for _, op := range operations {
		st := NewSegmentTree(arr, op)
		if !st.Validate() {
			t.Errorf("Tree validation failed for operation %s", op)
		}

		// Test after updates
		st.Update(3, 100)
		if !st.Validate() {
			t.Errorf("Tree validation failed after update for operation %s", op)
		}

		// Test after range updates (only for sum)
		if op == Sum {
			st.RangeUpdate(1, 5, 10)
			if !st.Validate() {
				t.Errorf("Tree validation failed after range update for operation %s", op)
			}
		}
	}
}

// TestLargeArray tests performance with larger arrays
func TestLargeArray(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large array test in short mode")
	}

	n := 100000
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(1000)
	}

	st := NewSumSegmentTree(arr)

	// Test multiple operations
	start := time.Now()

	// Perform 1000 range queries
	for i := 0; i < 1000; i++ {
		l := rand.Intn(n)
		r := rand.Intn(n)
		if l > r {
			l, r = r, l
		}
		st.RangeSum(l, r)
	}

	queryTime := time.Since(start)
	t.Logf("1000 range queries on %d elements: %v", n, queryTime)

	// Perform 1000 updates
	start = time.Now()
	for i := 0; i < 1000; i++ {
		pos := rand.Intn(n)
		val := rand.Intn(1000)
		st.Update(pos, val)
	}

	updateTime := time.Since(start)
	t.Logf("1000 updates on %d elements: %v", n, updateTime)

	// Validate tree is still correct
	if !st.Validate() {
		t.Error("Tree validation failed after large operations")
	}
}

// Benchmark tests comparing segment tree with naive approaches

func BenchmarkSegmentTreeQuery(b *testing.B) {
	arr := make([]int, 10000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000)
	}

	st := NewSumSegmentTree(arr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := rand.Intn(len(arr))
		r := rand.Intn(len(arr))
		if l > r {
			l, r = r, l
		}
		st.RangeSum(l, r)
	}
}

func BenchmarkNaiveRangeSum(b *testing.B) {
	arr := make([]int, 10000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := rand.Intn(len(arr))
		r := rand.Intn(len(arr))
		if l > r {
			l, r = r, l
		}

		// Naive O(n) sum
		sum := 0
		for j := l; j <= r; j++ {
			sum += arr[j]
		}
	}
}

func BenchmarkSegmentTreeUpdate(b *testing.B) {
	arr := make([]int, 10000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000)
	}

	st := NewSumSegmentTree(arr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos := rand.Intn(len(arr))
		val := rand.Intn(1000)
		st.Update(pos, val)
	}
}

func BenchmarkSegmentTreeRangeUpdate(b *testing.B) {
	arr := make([]int, 10000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000)
	}

	st := NewSumSegmentTree(arr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := rand.Intn(len(arr))
		r := rand.Intn(len(arr))
		if l > r {
			l, r = r, l
		}
		val := rand.Intn(100)
		st.RangeUpdate(l, r, val)
	}
}

func BenchmarkMinSegmentTree(b *testing.B) {
	arr := make([]int, 10000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000)
	}

	st := NewMinSegmentTree(arr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := rand.Intn(len(arr))
		r := rand.Intn(len(arr))
		if l > r {
			l, r = r, l
		}
		st.RangeMin(l, r)
	}
}

func BenchmarkNaiveRangeMin(b *testing.B) {
	arr := make([]int, 10000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := rand.Intn(len(arr))
		r := rand.Intn(len(arr))
		if l > r {
			l, r = r, l
		}

		// Naive O(n) minimum
		minVal := math.MaxInt32
		for j := l; j <= r; j++ {
			if arr[j] < minVal {
				minVal = arr[j]
			}
		}
	}
}
