package fenwicktree

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// TestLSBFunction tests the core LSB bit manipulation
func TestLSBFunction(t *testing.T) {
	testCases := []struct {
		input, expected int
	}{
		{1, 1},   // 001 -> 001
		{2, 2},   // 010 -> 010
		{3, 1},   // 011 -> 001
		{4, 4},   // 100 -> 100
		{5, 1},   // 101 -> 001
		{6, 2},   // 110 -> 010
		{7, 1},   // 111 -> 001
		{8, 8},   // 1000 -> 1000
		{12, 4},  // 1100 -> 0100
		{16, 16}, // 10000 -> 10000
	}

	for _, tc := range testCases {
		result := lsb(tc.input)
		if result != tc.expected {
			t.Errorf("lsb(%d): expected %d, got %d", tc.input, tc.expected, result)
		}
	}
}

// TestNewFenwickTree tests basic tree creation
func TestNewFenwickTree(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	ft := NewFenwickTree(arr)

	if ft.Size() != 5 {
		t.Errorf("Expected size 5, got %d", ft.Size())
	}

	// Test empty array
	emptyFt := NewFenwickTree([]int{})
	if emptyFt.Size() != 0 {
		t.Error("Empty array should have size 0")
	}

	// Test single element
	singleFt := NewFenwickTree([]int{42})
	if singleFt.Size() != 1 || singleFt.Get(0) != 42 {
		t.Error("Single element tree creation failed")
	}
}

// TestBasicOperations tests fundamental Fenwick Tree operations
func TestBasicOperations(t *testing.T) {
	arr := []int{1, 3, 5, 7, 9, 11}
	ft := NewFenwickTree(arr)

	// Test individual gets
	for i, expected := range arr {
		result := ft.Get(i)
		if result != expected {
			t.Errorf("Get(%d): expected %d, got %d", i, expected, result)
		}
	}

	// Test prefix sums
	prefixTests := []struct {
		index    int
		expected int64
	}{
		{0, 1},  // sum[0..0] = 1
		{1, 4},  // sum[0..1] = 1+3
		{2, 9},  // sum[0..2] = 1+3+5
		{3, 16}, // sum[0..3] = 1+3+5+7
		{4, 25}, // sum[0..4] = 1+3+5+7+9
		{5, 36}, // sum[0..5] = 1+3+5+7+9+11
	}

	for _, tc := range prefixTests {
		result := ft.PrefixSum(tc.index)
		if result != tc.expected {
			t.Errorf("PrefixSum(%d): expected %d, got %d", tc.index, tc.expected, result)
		}
	}

	// Test range sums
	rangeTests := []struct {
		l, r, expected int
	}{
		{0, 0, 1},  // Single element
		{0, 1, 4},  // First two
		{1, 3, 15}, // Middle range: 3+5+7
		{0, 5, 36}, // Entire array
		{2, 4, 21}, // 5+7+9
		{3, 5, 27}, // 7+9+11
	}

	for _, tc := range rangeTests {
		result := ft.RangeSum(tc.l, tc.r)
		if result != tc.expected {
			t.Errorf("RangeSum(%d,%d): expected %d, got %d",
				tc.l, tc.r, tc.expected, result)
		}
	}
}

// TestUpdates tests point updates and their effects
func TestUpdates(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	ft := NewFenwickTree(arr)

	// Test point update
	ft.Update(2, 10) // Add 10 to index 2
	expected := []int{1, 2, 13, 4, 5}
	result := ft.ToArray()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After update: expected %v, got %v", expected, result)
	}

	// Test Set operation
	ft.Set(1, 100) // Set index 1 to 100
	expected = []int{1, 100, 13, 4, 5}
	result = ft.ToArray()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After set: expected %v, got %v", expected, result)
	}

	// Test total sum
	expectedSum := int64(1 + 100 + 13 + 4 + 5)
	totalSum := ft.TotalSum()
	if totalSum != expectedSum {
		t.Errorf("Total sum: expected %d, got %d", expectedSum, totalSum)
	}
}

// TestEdgeCases tests boundary conditions and edge cases
func TestEdgeCases(t *testing.T) {
	ft := NewFenwickTree([]int{10, 20, 30})

	// Test out of bounds operations
	ft.Update(-1, 5) // Should be ignored
	ft.Update(10, 5) // Should be ignored
	ft.Set(-1, 100)  // Should be ignored
	ft.Set(10, 100)  // Should be ignored

	// Array should remain unchanged
	expected := []int{10, 20, 30}
	result := ft.ToArray()
	if !reflect.DeepEqual(result, expected) {
		t.Error("Out of bounds operations should be ignored")
	}

	// Test out of bounds gets
	if ft.Get(-1) != 0 {
		t.Error("Get with negative index should return 0")
	}
	if ft.Get(10) != 0 {
		t.Error("Get with out of bounds index should return 0")
	}

	// Test prefix sum edge cases
	if ft.PrefixSum(-1) != 0 {
		t.Error("PrefixSum with negative index should return 0")
	}

	// Test range sum edge cases
	if ft.RangeSum(2, 1) != 0 {
		t.Error("Invalid range (l > r) should return 0")
	}
	if ft.RangeSum(-5, -1) != 0 {
		t.Error("Negative range should return 0")
	}
}

// TestRangeUpdateFenwickTree tests range updates with difference array
func TestRangeUpdateFenwickTree(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	ruft := NewRangeUpdateFenwickTree(arr)

	// Test initial state
	expected := []int{1, 2, 3, 4, 5}
	result := ruft.ToArray()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Initial state: expected %v, got %v", expected, result)
	}

	// Test range update
	ruft.RangeUpdate(1, 3, 10) // Add 10 to indices 1-3
	expected = []int{1, 12, 13, 14, 5}
	result = ruft.ToArray()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After range update: expected %v, got %v", expected, result)
	}

	// Test point update
	ruft.PointUpdate(0, 100) // Add 100 to index 0
	expected = []int{101, 12, 13, 14, 5}
	result = ruft.ToArray()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After point update: expected %v, got %v", expected, result)
	}

	// Test overlapping range updates
	ruft.RangeUpdate(2, 4, 5) // Add 5 to indices 2-4
	expected = []int{101, 12, 18, 19, 10}
	result = ruft.ToArray()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After overlapping update: expected %v, got %v", expected, result)
	}
}

// TestFenwickTree2D tests 2D matrix operations
func TestFenwickTree2D(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	ft2d := NewFenwickTree2D(matrix)

	// Test prefix sum queries
	testCases := []struct {
		row, col int
		expected int64
	}{
		{0, 0, 1},  // Just (0,0)
		{0, 1, 3},  // Sum of first row up to col 1: 1+2
		{1, 1, 12}, // Sum of 2x2 top-left: 1+2+4+5
		{2, 2, 45}, // Sum of entire matrix
		{1, 2, 21}, // Sum of first 2 rows: 1+2+3+4+5+6
	}

	for _, tc := range testCases {
		result := ft2d.PrefixSum2D(tc.row, tc.col)
		if result != tc.expected {
			t.Errorf("PrefixSum2D(%d,%d): expected %d, got %d",
				tc.row, tc.col, tc.expected, result)
		}
	}

	// Test range sum queries
	rangeTests := []struct {
		r1, c1, r2, c2 int
		expected       int64
	}{
		{0, 0, 0, 0, 1},  // Single cell
		{0, 0, 1, 1, 12}, // 2x2 top-left
		{1, 1, 2, 2, 28}, // 2x2 bottom-right: 5+6+8+9
		{0, 1, 1, 2, 16}, // Rectangle: 2+3+5+6
	}

	for _, tc := range rangeTests {
		result := ft2d.RangeSum2D(tc.r1, tc.c1, tc.r2, tc.c2)
		if result != tc.expected {
			t.Errorf("RangeSum2D(%d,%d,%d,%d): expected %d, got %d",
				tc.r1, tc.c1, tc.r2, tc.c2, tc.expected, result)
		}
	}

	// Test updates
	ft2d.Update(1, 1, 10) // Add 10 to position (1,1)

	// Check that affected queries change correctly
	if ft2d.PrefixSum2D(1, 1) != 22 { // Was 12, now 12+10
		t.Error("Update should affect prefix sum")
	}
	if ft2d.RangeSum2D(1, 1, 1, 1) != 15 { // Was 5, now 5+10
		t.Error("Update should affect range sum")
	}
}

// TestFrequencyFenwickTree tests frequency tracking and order statistics
func TestFrequencyFenwickTree(t *testing.T) {
	values := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	fft := NewFrequencyFenwickTree(values)

	// Test frequency counting
	freqTests := []struct {
		value, expected int
	}{
		{1, 2}, // appears twice
		{3, 2}, // appears twice
		{5, 2}, // appears twice
		{2, 1}, // appears once
		{4, 1}, // appears once
		{7, 0}, // doesn't appear
	}

	for _, tc := range freqTests {
		result := fft.Count(tc.value)
		if result != tc.expected {
			t.Errorf("Count(%d): expected %d, got %d", tc.value, tc.expected, result)
		}
	}

	// Test total count
	if fft.TotalCount() != len(values) {
		t.Errorf("Total count: expected %d, got %d", len(values), fft.TotalCount())
	}

	// Test count less than
	lessTests := []struct {
		value, expected int
	}{
		{1, 0}, // No values less than 1
		{2, 2}, // Two 1's are less than 2
		{4, 5}, // 1,1,2,3,3 are less than 4
	}

	for _, tc := range lessTests {
		result := fft.CountLess(tc.value)
		if result != tc.expected {
			t.Errorf("CountLess(%d): expected %d, got %d", tc.value, tc.expected, result)
		}
	}

	// Test count range
	rangeCount := fft.CountRange(2, 5)
	expected := 6 // 2,3,3,4,5,5
	if rangeCount != expected {
		t.Errorf("CountRange(2,5): expected %d, got %d", expected, rangeCount)
	}

	// Test k-th smallest
	kthTests := []struct {
		k        int
		expected int
		exists   bool
	}{
		{1, 1, true},   // 1st smallest is 1
		{3, 2, true},   // 3rd smallest is 2
		{5, 3, true},   // 5th smallest is 3
		{0, 0, false},  // Invalid k
		{20, 0, false}, // k too large
	}

	for _, tc := range kthTests {
		result, exists := fft.KthSmallest(tc.k)
		if exists != tc.exists {
			t.Errorf("KthSmallest(%d) existence: expected %v, got %v", tc.k, tc.exists, exists)
		}
		if exists && result != tc.expected {
			t.Errorf("KthSmallest(%d): expected %d, got %d", tc.k, tc.expected, result)
		}
	}

	// Test insert and delete
	fft.Insert(3)          // Insert existing value
	if fft.Count(3) != 3 { // Was 2, now should be 3
		t.Error("Insert of existing value should increase count")
	}

	fft.Delete(3)
	if fft.Count(3) != 2 { // Back to 2
		t.Error("Delete should decrease count")
	}

	// Test delete non-existent
	fft.Delete(100) // Should not crash
	if fft.TotalCount() != len(values) {
		t.Error("Delete non-existent should not change total count")
	}
}

// TestPowerOfTwoBehavior tests Fenwick Tree with power-of-2 sizes
func TestPowerOfTwoBehavior(t *testing.T) {
	// Test with various power-of-2 sizes
	sizes := []int{1, 2, 4, 8, 16, 32}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = i + 1
		}

		ft := NewFenwickTree(arr)

		// Test that all operations work correctly
		expectedSum := size * (size + 1) / 2
		if int(ft.TotalSum()) != expectedSum {
			t.Errorf("Size %d: expected sum %d, got %d", size, expectedSum, int(ft.TotalSum()))
		}

		// Test range sum over entire array
		rangeSum := ft.RangeSum(0, size-1)
		if rangeSum != expectedSum {
			t.Errorf("Size %d: range sum expected %d, got %d", size, expectedSum, rangeSum)
		}
	}
}

// TestLargeDataset tests performance with larger datasets
func TestLargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	n := 100000
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(1000)
	}

	// Build Fenwick Tree
	start := time.Now()
	ft := NewFenwickTree(arr)
	buildTime := time.Since(start)
	t.Logf("Built Fenwick Tree with %d elements in %v", n, buildTime)

	// Test multiple operations
	start = time.Now()
	for i := 0; i < 1000; i++ {
		l := rand.Intn(n)
		r := rand.Intn(n)
		if l > r {
			l, r = r, l
		}
		ft.RangeSum(l, r)
	}
	queryTime := time.Since(start)
	t.Logf("1000 range queries: %v", queryTime)

	start = time.Now()
	for i := 0; i < 1000; i++ {
		pos := rand.Intn(n)
		delta := rand.Intn(100) - 50
		ft.Update(pos, delta)
	}
	updateTime := time.Since(start)
	t.Logf("1000 updates: %v", updateTime)
}

// Benchmark tests
func BenchmarkFenwickTreeQuery(b *testing.B) {
	arr := make([]int, 10000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000)
	}

	ft := NewFenwickTree(arr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := rand.Intn(len(arr))
		r := rand.Intn(len(arr))
		if l > r {
			l, r = r, l
		}
		ft.RangeSum(l, r)
	}
}

func BenchmarkFenwickTreeUpdate(b *testing.B) {
	arr := make([]int, 10000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000)
	}

	ft := NewFenwickTree(arr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos := rand.Intn(len(arr))
		delta := rand.Intn(100)
		ft.Update(pos, delta)
	}
}

func BenchmarkFenwickTreePrefixSum(b *testing.B) {
	arr := make([]int, 10000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000)
	}

	ft := NewFenwickTree(arr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos := rand.Intn(len(arr))
		ft.PrefixSum(pos)
	}
}

func BenchmarkRangeUpdateFenwick(b *testing.B) {
	arr := make([]int, 10000)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(1000)
	}

	ruft := NewRangeUpdateFenwickTree(arr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := rand.Intn(len(arr))
		r := rand.Intn(len(arr))
		if l > r {
			l, r = r, l
		}
		delta := rand.Intn(100)
		ruft.RangeUpdate(l, r, delta)
	}
}

func BenchmarkFenwick2D(b *testing.B) {
	size := 100
	matrix := make([][]int, size)
	for i := 0; i < size; i++ {
		matrix[i] = make([]int, size)
		for j := 0; j < size; j++ {
			matrix[i][j] = rand.Intn(100)
		}
	}

	ft2d := NewFenwickTree2D(matrix)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r1 := rand.Intn(size)
		c1 := rand.Intn(size)
		r2 := rand.Intn(size)
		c2 := rand.Intn(size)
		if r1 > r2 {
			r1, r2 = r2, r1
		}
		if c1 > c2 {
			c1, c2 = c2, c1
		}
		ft2d.RangeSum2D(r1, c1, r2, c2)
	}
}

func BenchmarkFrequencyFenwick(b *testing.B) {
	values := make([]int, 10000)
	for i := 0; i < len(values); i++ {
		values[i] = rand.Intn(1000)
	}

	fft := NewFrequencyFenwickTree(values)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch rand.Intn(4) {
		case 0:
			fft.Count(rand.Intn(1000))
		case 1:
			fft.CountLess(rand.Intn(1000))
		case 2:
			fft.KthSmallest(rand.Intn(len(values)) + 1)
		case 3:
			fft.Insert(rand.Intn(1000))
		}
	}
}
