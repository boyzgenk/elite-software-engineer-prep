// Complexity Analysis Practice Problems for Go Developers
// Hands-on exercises to master Big O analysis and optimization

package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// ========================================================================
// PROBLEM 1: TIME COMPLEXITY ANALYSIS
// ========================================================================

// Analyze the time complexity of these functions
// Solutions provided at the bottom of each function

func problem1_a(arr []int) int {
	// What's the time complexity?
	for i := 0; i < len(arr); i++ {
		if arr[i] == 42 {
			return i
		}
	}
	return -1
	// SOLUTION: O(n) - worst case we check every element
}

func problem1_b(matrix [][]int) int {
	// What's the time complexity?
	sum := 0
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			sum += matrix[i][j]
		}
	}
	return sum
	// SOLUTION: O(n²) or O(rows × cols) - nested loops over 2D structure
}

func problem1_c(n int) int {
	// What's the time complexity?
	if n <= 1 {
		return n
	}
	return problem1_c(n-1) + problem1_c(n-2)
	// SOLUTION: O(2ⁿ) - exponential due to overlapping subproblems
}

func problem1_d(arr []int, target int) int {
	// What's the time complexity? (assume arr is sorted)
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
	// SOLUTION: O(log n) - binary search eliminates half each iteration
}

func problem1_e(arr []int) {
	// What's the time complexity?
	for i := 0; i < len(arr)-1; i++ {
		for j := 0; j < len(arr)-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	// SOLUTION: O(n²) - bubble sort with nested loops
}

// ========================================================================
// PROBLEM 2: OPTIMIZATION CHALLENGES
// ========================================================================

// CHALLENGE 2A: Optimize this O(n²) solution to O(n)
func twoSum_slow(nums []int, target int) []int {
	// SLOW: O(n²) nested loop solution
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

func twoSum_fast(nums []int, target int) []int {
	// OPTIMIZED: O(n) using hash map
	seen := make(map[int]int)
	for i, num := range nums {
		complement := target - num
		if j, found := seen[complement]; found {
			return []int{j, i}
		}
		seen[num] = i
	}
	return nil
	// SOLUTION: Use hash map to store seen numbers and their indices
	// Trade space (O(n)) for time improvement (O(n²) → O(n))
}

// CHALLENGE 2B: Optimize this string concatenation
func buildString_slow(words []string) string {
	// SLOW: O(n²) due to string immutability
	result := ""
	for _, word := range words {
		result += word + " " // Each concatenation creates new string
	}
	return result
}

func buildString_fast(words []string) string {
	// OPTIMIZED: O(n) using strings.Builder
	var builder strings.Builder
	for i, word := range words {
		builder.WriteString(word)
		if i < len(words)-1 {
			builder.WriteString(" ")
		}
	}
	return builder.String()
	// SOLUTION: strings.Builder amortizes allocation cost
}

// CHALLENGE 2C: Optimize this duplicate finder
func findDuplicates_slow(arr []int) []int {
	// SLOW: O(n²) nested loops
	var duplicates []int
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				// Check if already added
				found := false
				for _, dup := range duplicates {
					if dup == arr[i] {
						found = true
						break
					}
				}
				if !found {
					duplicates = append(duplicates, arr[i])
				}
			}
		}
	}
	return duplicates
}

func findDuplicates_fast(arr []int) []int {
	// OPTIMIZED: O(n) using hash map
	seen := make(map[int]bool)
	duplicates := make(map[int]bool)
	var result []int

	for _, num := range arr {
		if seen[num] && !duplicates[num] {
			result = append(result, num)
			duplicates[num] = true
		} else {
			seen[num] = true
		}
	}
	return result
	// SOLUTION: Single pass with hash maps for tracking
}

// ========================================================================
// PROBLEM 3: SPACE COMPLEXITY OPTIMIZATION
// ========================================================================

// CHALLENGE 3A: Reduce space complexity of this solution
func reverseArray_extraSpace(arr []int) []int {
	// O(n) extra space
	reversed := make([]int, len(arr))
	for i, val := range arr {
		reversed[len(arr)-1-i] = val
	}
	return reversed
}

func reverseArray_inPlace(arr []int) {
	// O(1) extra space - modify input array
	for i := 0; i < len(arr)/2; i++ {
		j := len(arr) - 1 - i
		arr[i], arr[j] = arr[j], arr[i]
	}
	// SOLUTION: Two-pointer technique, swap elements in place
}

// CHALLENGE 3B: Check if string is palindrome with minimal space
func isPalindrome_extraSpace(s string) bool {
	// O(n) extra space - create reversed string
	cleaned := strings.ToLower(s)
	reversed := ""
	for i := len(cleaned) - 1; i >= 0; i-- {
		reversed += string(cleaned[i])
	}
	return cleaned == reversed
}

func isPalindrome_constantSpace(s string) bool {
	// O(1) extra space - two pointers
	s = strings.ToLower(s)
	left, right := 0, len(s)-1

	for left < right {
		// Skip non-alphanumeric characters
		for left < right && !isAlphanumeric(s[left]) {
			left++
		}
		for left < right && !isAlphanumeric(s[right]) {
			right--
		}

		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}
	return true
}

func isAlphanumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}

// ========================================================================
// PROBLEM 4: ALGORITHM SELECTION CHALLENGES
// ========================================================================

// Choose the best algorithm based on constraints

type SearchProblem struct {
	data   []int
	sorted bool
	size   int
}

func (sp *SearchProblem) OptimalSearch(target int) int {
	// CHALLENGE: Choose optimal search strategy based on data characteristics

	if sp.size < 10 {
		// Small arrays: linear search is fine due to constant overhead
		return sp.linearSearch(target)
	}

	if sp.sorted {
		// Sorted data: binary search is optimal O(log n)
		return sp.binarySearch(target)
	} else {
		// Unsorted data: linear search is only option O(n)
		// Unless we sort first, but that's O(n log n) + O(log n)
		return sp.linearSearch(target)
	}
}

func (sp *SearchProblem) linearSearch(target int) int {
	for i, val := range sp.data {
		if val == target {
			return i
		}
	}
	return -1
}

func (sp *SearchProblem) binarySearch(target int) int {
	left, right := 0, len(sp.data)-1
	for left <= right {
		mid := left + (right-left)/2
		if sp.data[mid] == target {
			return mid
		} else if sp.data[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

// ========================================================================
// PROBLEM 5: PRACTICAL OPTIMIZATION SCENARIOS
// ========================================================================

// SCENARIO A: Large dataset processing
type DataProcessor struct {
	cache map[string]int
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		cache: make(map[string]int),
	}
}

func (dp *DataProcessor) ProcessDataSlow(data []string) int {
	// SLOW: Recomputing expensive operations
	total := 0
	for _, item := range data {
		total += dp.expensiveComputation(item)
	}
	return total
}

func (dp *DataProcessor) ProcessDataFast(data []string) int {
	// OPTIMIZED: Using memoization
	total := 0
	for _, item := range data {
		if cached, found := dp.cache[item]; found {
			total += cached
		} else {
			result := dp.expensiveComputation(item)
			dp.cache[item] = result
			total += result
		}
	}
	return total
}

func (dp *DataProcessor) expensiveComputation(item string) int {
	// Simulate expensive operation
	time.Sleep(1 * time.Millisecond)
	return len(item) * len(item)
}

// SCENARIO B: Memory-efficient data structure selection
type FrequencyCounter struct {
	// Choose appropriate data structure based on use case
}

func (fc *FrequencyCounter) CountFrequency_Map(items []string) map[string]int {
	// Good for: Unknown number of unique items, need all frequencies
	// O(n) time, O(k) space where k = unique items
	freq := make(map[string]int)
	for _, item := range items {
		freq[item]++
	}
	return freq
}

func (fc *FrequencyCounter) CountFrequency_Array(items []int, maxValue int) []int {
	// Good for: Known range of values, memory is critical
	// O(n) time, O(maxValue) space - more memory efficient if maxValue < unique items
	freq := make([]int, maxValue+1)
	for _, item := range items {
		if item >= 0 && item <= maxValue {
			freq[item]++
		}
	}
	return freq
}

// ========================================================================
// PROBLEM 6: COMPLEXITY TRADE-OFF ANALYSIS
// ========================================================================

// Different approaches to find K largest elements

func findKLargest_Sort(arr []int, k int) []int {
	// TIME: O(n log n), SPACE: O(n) if we copy
	sorted := make([]int, len(arr))
	copy(sorted, arr)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	return sorted[:k]
	// ANALYSIS: Simple but overkill - we sort entire array just for k elements
}

func findKLargest_MinHeap(arr []int, k int) []int {
	// TIME: O(n log k), SPACE: O(k)
	// Maintain min-heap of size k
	heap := make([]int, 0, k)

	for _, num := range arr {
		if len(heap) < k {
			heap = append(heap, num)
			heapifyUp(heap, len(heap)-1)
		} else if num > heap[0] {
			heap[0] = num
			heapifyDown(heap, 0)
		}
	}

	return heap
	// ANALYSIS: Better for large n, small k
}

func findKLargest_QuickSelect(arr []int, k int) []int {
	// TIME: O(n) average, O(n²) worst, SPACE: O(1) if in-place
	arrCopy := make([]int, len(arr))
	copy(arrCopy, arr)

	quickSelectKth(arrCopy, 0, len(arrCopy)-1, k)
	result := arrCopy[len(arrCopy)-k:]
	sort.Sort(sort.Reverse(sort.IntSlice(result)))
	return result
	// ANALYSIS: Best average case, but can degrade
}

// Helper functions for heap operations
func heapifyUp(heap []int, index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if heap[index] >= heap[parent] {
			break
		}
		heap[index], heap[parent] = heap[parent], heap[index]
		index = parent
	}
}

func heapifyDown(heap []int, index int) {
	for {
		smallest := index
		left := 2*index + 1
		right := 2*index + 2

		if left < len(heap) && heap[left] < heap[smallest] {
			smallest = left
		}
		if right < len(heap) && heap[right] < heap[smallest] {
			smallest = right
		}

		if smallest == index {
			break
		}

		heap[index], heap[smallest] = heap[smallest], heap[index]
		index = smallest
	}
}

func quickSelectKth(arr []int, left, right, k int) {
	if left < right {
		pivot := partition(arr, left, right)
		if pivot == len(arr)-k {
			return
		} else if pivot > len(arr)-k {
			quickSelectKth(arr, left, pivot-1, k)
		} else {
			quickSelectKth(arr, pivot+1, right, k)
		}
	}
}

func partition(arr []int, left, right int) int {
	pivot := arr[right]
	i := left - 1

	for j := left; j < right; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[right] = arr[right], arr[i+1]
	return i + 1
}

// ========================================================================
// PROBLEM 7: BENCHMARKING AND PROFILING EXERCISE
// ========================================================================

func benchmarkComplexityDemo() {
	fmt.Println("🔬 COMPLEXITY BENCHMARKING DEMO")
	fmt.Println("===============================")

	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		fmt.Printf("\n--- Dataset size: %d ---\n", size)

		// Generate test data
		data := make([]int, size)
		for i := range data {
			data[i] = size - i // Reverse sorted for worst case
		}

		target := 1 // Worst case - element at end

		// Benchmark different search algorithms
		benchmarkFunction("Linear Search", func() {
			linearSearchBench(data, target)
		})

		// Sort data for binary search
		sortedData := make([]int, len(data))
		copy(sortedData, data)
		sort.Ints(sortedData)

		benchmarkFunction("Binary Search", func() {
			binarySearchBench(sortedData, target)
		})

		// Benchmark different sorting algorithms (only for smaller sizes)
		if size <= 10000 {
			unsorted := make([]int, len(data))
			copy(unsorted, data)
			benchmarkFunction("Bubble Sort", func() {
				bubbleSortBench(unsorted)
			})
		}

		unsorted2 := make([]int, len(data))
		copy(unsorted2, data)
		benchmarkFunction("Go Built-in Sort", func() {
			sort.Ints(unsorted2)
		})
	}
}

func benchmarkFunction(name string, fn func()) {
	start := time.Now()
	fn()
	duration := time.Since(start)
	fmt.Printf("%-20s: %v\n", name, duration)
}

func linearSearchBench(arr []int, target int) int {
	for i, val := range arr {
		if val == target {
			return i
		}
	}
	return -1
}

func binarySearchBench(arr []int, target int) int {
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

func bubbleSortBench(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

// ========================================================================
// PROBLEM 8: ADVANCED ANALYSIS CHALLENGES
// ========================================================================

// Analyze and optimize these complex scenarios

func matrixChainMultiplication(dimensions []int) int {
	// CHALLENGE: Find optimal parenthesization for matrix multiplication
	// Classic DP problem: O(n³) time, O(n²) space
	n := len(dimensions) - 1
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	// l is chain length
	for l := 2; l <= n; l++ {
		for i := 0; i <= n-l; i++ {
			j := i + l - 1
			dp[i][j] = math.MaxInt32

			for k := i; k < j; k++ {
				cost := dp[i][k] + dp[k+1][j] + dimensions[i]*dimensions[k+1]*dimensions[j+1]
				if cost < dp[i][j] {
					dp[i][j] = cost
				}
			}
		}
	}

	return dp[0][n-1]
	// ANALYSIS: Classic O(n³) DP solution, optimal for this problem
}

func longestCommonSubsequence(text1, text2 string) int {
	// CHALLENGE: Find LCS length
	// Can we optimize space complexity?

	m, n := len(text1), len(text2)

	// Space-optimized version: O(min(m,n)) instead of O(m*n)
	if m < n {
		text1, text2 = text2, text1
		m, n = n, m
	}

	prev := make([]int, n+1)
	curr := make([]int, n+1)

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				curr[j] = prev[j-1] + 1
			} else {
				curr[j] = max(prev[j], curr[j-1])
			}
		}
		prev, curr = curr, prev
	}

	return prev[n]
	// ANALYSIS: Reduced space from O(m*n) to O(min(m,n))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ========================================================================
// MAIN FUNCTION - PRACTICE PROBLEMS RUNNER
// ========================================================================

func main() {
	fmt.Println("🎯 COMPLEXITY ANALYSIS PRACTICE PROBLEMS")
	fmt.Println("========================================")

	fmt.Println("\n=== PROBLEM 1: Time Complexity Analysis ===")
	fmt.Println("Analyze the functions problem1_a through problem1_e")
	fmt.Println("Solutions are provided as comments in the code")

	fmt.Println("\n=== PROBLEM 2: Optimization Challenges ===")

	// Demonstrate two sum optimization
	nums := []int{2, 7, 11, 15}
	target := 9

	fmt.Printf("Two Sum Example - nums: %v, target: %d\n", nums, target)
	fmt.Printf("Slow O(n²): %v\n", twoSum_slow(nums, target))
	fmt.Printf("Fast O(n): %v\n", twoSum_fast(nums, target))

	// Demonstrate string building optimization
	words := []string{"Hello", "World", "from", "Go"}
	fmt.Printf("\nString Building - words: %v\n", words)
	fmt.Printf("Slow O(n²): '%s'\n", buildString_slow(words))
	fmt.Printf("Fast O(n): '%s'\n", buildString_fast(words))

	fmt.Println("\n=== PROBLEM 3: Space Complexity Optimization ===")

	arr := []int{1, 2, 3, 4, 5}
	fmt.Printf("Original array: %v\n", arr)

	reversed := reverseArray_extraSpace(arr)
	fmt.Printf("Reversed (extra space): %v\n", reversed)

	arrCopy := make([]int, len(arr))
	copy(arrCopy, arr)
	reverseArray_inPlace(arrCopy)
	fmt.Printf("Reversed (in-place): %v\n", arrCopy)

	fmt.Println("\n=== PROBLEM 4: Algorithm Selection ===")

	searchProb := &SearchProblem{
		data:   []int{1, 3, 5, 7, 9, 11, 13, 15},
		sorted: true,
		size:   8,
	}

	result := searchProb.OptimalSearch(7)
	fmt.Printf("Optimal search for 7 in sorted array: index %d\n", result)

	fmt.Println("\n=== PROBLEM 5: K Largest Elements Comparison ===")

	testArr := []int{3, 2, 1, 5, 6, 4}
	k := 2

	fmt.Printf("Find %d largest in %v:\n", k, testArr)
	fmt.Printf("Sort approach: %v\n", findKLargest_Sort(testArr, k))
	fmt.Printf("Min-heap approach: %v\n", findKLargest_MinHeap(testArr, k))
	fmt.Printf("QuickSelect approach: %v\n", findKLargest_QuickSelect(testArr, k))

	fmt.Println("\n=== PROBLEM 6: Benchmarking Demo ===")
	benchmarkComplexityDemo()

	fmt.Println("\n=== PROBLEM 7: Advanced DP Examples ===")

	// Matrix chain multiplication example
	matrices := []int{1, 2, 3, 4} // 3 matrices: 1x2, 2x3, 3x4
	minOps := matrixChainMultiplication(matrices)
	fmt.Printf("Min scalar multiplications for matrices %v: %d\n", matrices, minOps)

	// LCS example
	text1, text2 := "abcde", "ace"
	lcsLen := longestCommonSubsequence(text1, text2)
	fmt.Printf("LCS length of '%s' and '%s': %d\n", text1, text2, lcsLen)

	fmt.Println("\n🏆 Practice Problems Complete!")
	fmt.Println("Next steps:")
	fmt.Println("1. Implement your own solutions without looking at the code")
	fmt.Println("2. Analyze complexity of real-world algorithms in your projects")
	fmt.Println("3. Use benchmarking to validate your optimization decisions")
	fmt.Println("4. Practice on LeetCode/competitive programming platforms")
}
