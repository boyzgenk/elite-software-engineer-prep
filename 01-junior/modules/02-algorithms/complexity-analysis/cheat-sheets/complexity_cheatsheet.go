// Complexity Analysis Cheat Sheet for Go Developers
// Quick reference guide for Big O analysis and optimization patterns

package main

import (
	"fmt"
	"sort"
	"time"
)

// ========================================================================
// TIME COMPLEXITY CHEAT SHEET
// ========================================================================

// O(1) - Constant Time Examples
func constant_examples() {
	// Array/slice index access
	arr := []int{1, 2, 3, 4, 5}
	fmt.Println(arr[2]) // Always takes same time regardless of array size

	// Map access (average case)
	m := make(map[string]int)
	m["key"] = 42
	value := m["key"] // Hash table lookup

	// Stack operations
	stack := []int{1, 2, 3}
	stack = append(stack, 4)     // Push - O(1) amortized
	stack = stack[:len(stack)-1] // Pop - O(1)

	fmt.Printf("Constant time operations completed: %d\n", value)
}

// O(log n) - Logarithmic Time Examples
func logarithmic_examples() {
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

	// Binary search
	target := 7
	result := binarySearch(arr, target)
	fmt.Printf("Binary search result: %d\n", result)

	// Balanced tree operations (conceptual)
	// Height of balanced BST = log n
	// Each level eliminates half the possibilities
}

func binarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1

	for left <= right {
		mid := left + (right-left)/2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

// O(n) - Linear Time Examples
func linear_examples() {
	arr := []int{5, 2, 8, 1, 9, 3}

	// Single pass through array
	sum := 0
	for _, val := range arr {
		sum += val // Visit each element once
	}

	// Linear search
	target := 8
	index := linearSearch(arr, target)

	// String/slice operations
	s := "hello world"
	count := 0
	for _, char := range s {
		if char == 'l' {
			count++
		}
	}

	fmt.Printf("Linear operations: sum=%d, index=%d, count=%d\n", sum, index, count)
}

func linearSearch(arr []int, target int) int {
	for i, val := range arr {
		if val == target {
			return i
		}
	}
	return -1
}

// O(n log n) - Linearithmic Time Examples
func linearithmic_examples() {
	arr := []int{64, 34, 25, 12, 22, 11, 90}

	// Efficient sorting algorithms
	arrCopy := make([]int, len(arr))
	copy(arrCopy, arr)

	// Go's built-in sort (Introsort - hybrid of quicksort, heapsort, insertion sort)
	sort.Ints(arrCopy)
	fmt.Printf("Sorted array: %v\n", arrCopy)

	// Custom merge sort implementation
	mergeSort(arr, 0, len(arr)-1)
	fmt.Printf("Merge sorted: %v\n", arr)
}

func mergeSort(arr []int, left, right int) {
	if left < right {
		mid := left + (right-left)/2
		mergeSort(arr, left, mid)
		mergeSort(arr, mid+1, right)
		merge(arr, left, mid, right)
	}
}

func merge(arr []int, left, mid, right int) {
	n1 := mid - left + 1
	n2 := right - mid

	leftArr := make([]int, n1)
	rightArr := make([]int, n2)

	copy(leftArr, arr[left:mid+1])
	copy(rightArr, arr[mid+1:right+1])

	i, j, k := 0, 0, left

	for i < n1 && j < n2 {
		if leftArr[i] <= rightArr[j] {
			arr[k] = leftArr[i]
			i++
		} else {
			arr[k] = rightArr[j]
			j++
		}
		k++
	}

	for i < n1 {
		arr[k] = leftArr[i]
		i++
		k++
	}

	for j < n2 {
		arr[k] = rightArr[j]
		j++
		k++
	}
}

// O(n²) - Quadratic Time Examples
func quadratic_examples() {
	arr := []int{64, 34, 25, 12, 22, 11, 90}

	// Bubble sort - nested loops
	bubbleSort(arr)
	fmt.Printf("Bubble sorted: %v\n", arr)

	// Finding all pairs
	pairs := findAllPairs([]int{1, 2, 3, 4})
	fmt.Printf("All pairs: %v\n", pairs)

	// Matrix operations
	matrix := [][]int{{1, 2}, {3, 4}}
	result := matrixMultiply(matrix, matrix)
	fmt.Printf("Matrix multiplication: %v\n", result)
}

func bubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func findAllPairs(arr []int) [][]int {
	var pairs [][]int
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			pairs = append(pairs, []int{arr[i], arr[j]})
		}
	}
	return pairs
}

func matrixMultiply(a, b [][]int) [][]int {
	n := len(a)
	result := make([][]int, n)
	for i := range result {
		result[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

// ========================================================================
// SPACE COMPLEXITY CHEAT SHEET
// ========================================================================

// O(1) - Constant Space
func constantSpace(arr []int) int {
	// Only using fixed amount of extra variables
	max := arr[0]
	for _, val := range arr {
		if val > max {
			max = val
		}
	}
	return max // No extra space proportional to input size
}

// O(n) - Linear Space
func linearSpace(arr []int) []int {
	// Creating new slice proportional to input
	reversed := make([]int, len(arr))
	for i, val := range arr {
		reversed[len(arr)-1-i] = val
	}
	return reversed
}

// O(log n) - Logarithmic Space (typically recursion stack)
func logSpace(n int) int {
	if n <= 1 {
		return n
	}
	// Recursion depth = log n for binary operations
	return logSpace(n/2) + logSpace(n/2)
}

// ========================================================================
// GO-SPECIFIC OPTIMIZATION PATTERNS
// ========================================================================

type OptimizationDemo struct {
	data []int
}

// Memory pool pattern for reducing allocations
var intSlicePool = make(chan []int, 100)

func getSlice(size int) []int {
	select {
	case slice := <-intSlicePool:
		if cap(slice) >= size {
			return slice[:size]
		}
	default:
	}
	return make([]int, size)
}

func putSlice(slice []int) {
	if cap(slice) < 1000 { // Avoid keeping huge slices
		select {
		case intSlicePool <- slice[:0]:
		default:
		}
	}
}

// Benchmark-driven optimization example
func (o *OptimizationDemo) ProcessData() {
	// Use object pooling for frequent allocations
	tempSlice := getSlice(len(o.data))
	defer putSlice(tempSlice)

	// Process with pooled memory
	copy(tempSlice, o.data)
	sort.Ints(tempSlice)
}

// ========================================================================
// COMPLEXITY ANALYSIS QUICK REFERENCE
// ========================================================================

func printComplexityReference() {
	fmt.Println("=== TIME COMPLEXITY QUICK REFERENCE ===")
	fmt.Println("O(1)     - Constant     - Hash table access, array indexing")
	fmt.Println("O(log n) - Logarithmic  - Binary search, balanced tree operations")
	fmt.Println("O(n)     - Linear       - Single loop, linear search")
	fmt.Println("O(n log n) - Linearithmic - Efficient sorting (merge, heap, quick)")
	fmt.Println("O(n²)    - Quadratic    - Nested loops, bubble sort")
	fmt.Println("O(n³)    - Cubic        - Triple nested loops")
	fmt.Println("O(2ⁿ)    - Exponential  - Recursive fibonacci, subset generation")
	fmt.Println("O(n!)    - Factorial    - Permutation generation")

	fmt.Println("\n=== SPACE COMPLEXITY QUICK REFERENCE ===")
	fmt.Println("O(1)     - Constant     - Fixed variables, in-place algorithms")
	fmt.Println("O(log n) - Logarithmic  - Recursion stack for divide & conquer")
	fmt.Println("O(n)     - Linear       - Creating copy of input, recursion stack")
	fmt.Println("O(n²)    - Quadratic    - 2D matrix, memoization table")

	fmt.Println("\n=== GO OPTIMIZATION TIPS ===")
	fmt.Println("• Use make() with capacity to avoid slice reallocations")
	fmt.Println("• Prefer maps over slices for lookups (O(1) vs O(n))")
	fmt.Println("• Use sync.Pool for object reuse in hot paths")
	fmt.Println("• Profile with go tool pprof to identify bottlenecks")
	fmt.Println("• Consider memory layout - structs vs pointers")
	fmt.Println("• Use buffered channels to reduce goroutine blocking")
}

// ========================================================================
// PERFORMANCE MEASUREMENT UTILITIES
// ========================================================================

func measureTime(name string, fn func()) {
	start := time.Now()
	fn()
	duration := time.Since(start)
	fmt.Printf("%s took: %v\n", name, duration)
}

func demonstrateComplexityDifferences() {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := range arr {
			arr[i] = size - i // Reverse sorted for worst case
		}

		fmt.Printf("\n=== Array size: %d ===\n", size)

		// O(n) - Linear search
		measureTime("Linear Search O(n)", func() {
			linearSearch(arr, 1)
		})

		// O(log n) - Binary search (need sorted array)
		sorted := make([]int, len(arr))
		copy(sorted, arr)
		sort.Ints(sorted)
		measureTime("Binary Search O(log n)", func() {
			binarySearch(sorted, 1)
		})

		// O(n²) - Bubble sort (on copy to avoid modifying original)
		if size <= 1000 { // Skip for large arrays as it's too slow
			bubbleCopy := make([]int, len(arr))
			copy(bubbleCopy, arr)
			measureTime("Bubble Sort O(n²)", func() {
				bubbleSort(bubbleCopy)
			})
		}

		// O(n log n) - Merge sort
		mergeCopy := make([]int, len(arr))
		copy(mergeCopy, arr)
		measureTime("Merge Sort O(n log n)", func() {
			mergeSort(mergeCopy, 0, len(mergeCopy)-1)
		})
	}
}

func main() {
	fmt.Println("🚀 COMPLEXITY ANALYSIS CHEAT SHEET DEMO")
	fmt.Println("========================================")

	printComplexityReference()

	fmt.Println("\n=== EXAMPLES BY COMPLEXITY ===")
	constant_examples()
	logarithmic_examples()
	linear_examples()
	linearithmic_examples()
	quadratic_examples()

	fmt.Println("\n=== PERFORMANCE COMPARISON ===")
	demonstrateComplexityDifferences()

	fmt.Println("\n🎯 Use this cheat sheet for quick complexity analysis reference!")
}
