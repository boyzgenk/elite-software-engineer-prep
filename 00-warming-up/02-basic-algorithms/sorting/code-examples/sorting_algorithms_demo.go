// Sorting Algorithms Implementation - Junior Software Engineer Interview Prep
// Module 1: CS Fundamentals - Basic Algorithms
// Complete Golang implementations with educational focus

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// =============================================================================
// BUBBLE SORT - O(n²) Time Complexity, O(1) Space Complexity
// =============================================================================

// BubbleSort implements the bubble sort algorithm
// Best for: Understanding basic sorting concepts, small datasets
// Time Complexity: O(n²) average and worst case, O(n) best case (optimized version)
// Space Complexity: O(1)
func BubbleSort(arr []int) []int {
	fmt.Println("🫧 Starting Bubble Sort...")
	result := make([]int, len(arr))
	copy(result, arr)
	n := len(result)

	// Track number of swaps to optimize for already sorted arrays
	for i := 0; i < n-1; i++ {
		swapped := false
		fmt.Printf("Pass %d: ", i+1)

		// Last i elements are already in place
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				// Swap elements
				result[j], result[j+1] = result[j+1], result[j]
				swapped = true
				fmt.Printf("Swap(%d,%d) ", result[j+1], result[j])
			}
		}

		if !swapped {
			fmt.Printf("No swaps needed - array is sorted!")
			break
		}
		fmt.Printf("→ %v\n", result)
	}

	fmt.Printf("✅ Bubble Sort Complete: %v\n\n", result)
	return result
}

// OptimizedBubbleSort with early termination
func OptimizedBubbleSort(arr []int) []int {
	fmt.Println("🚀 Starting Optimized Bubble Sort...")
	result := make([]int, len(arr))
	copy(result, arr)
	n := len(result)

	for i := 0; i < n-1; i++ {
		swapped := false

		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
				swapped = true
			}
		}

		// If no swapping occurred, array is already sorted
		if !swapped {
			fmt.Printf("✅ Early termination at pass %d - array sorted!\n", i+1)
			break
		}
	}

	fmt.Printf("✅ Optimized Bubble Sort Complete: %v\n\n", result)
	return result
}

// =============================================================================
// MERGE SORT - O(n log n) Time Complexity, O(n) Space Complexity
// =============================================================================

// MergeSort implements the divide-and-conquer merge sort algorithm
// Best for: Large datasets, stable sorting, guaranteed O(n log n) performance
// Time Complexity: O(n log n) in all cases
// Space Complexity: O(n)
func MergeSort(arr []int) []int {
	fmt.Println("🔀 Starting Merge Sort...")
	result := make([]int, len(arr))
	copy(result, arr)

	if len(result) <= 1 {
		return result
	}

	return mergeSortRecursive(result, 0, len(result)-1)
}

func mergeSortRecursive(arr []int, left, right int) []int {
	if left >= right {
		return []int{arr[left]}
	}

	mid := left + (right-left)/2
	fmt.Printf("Dividing: arr[%d:%d] → left[%d:%d] right[%d:%d]\n",
		left, right, left, mid, mid+1, right)

	leftHalf := mergeSortRecursive(arr, left, mid)
	rightHalf := mergeSortRecursive(arr, mid+1, right)

	merged := merge(leftHalf, rightHalf)
	fmt.Printf("Merging: %v + %v → %v\n", leftHalf, rightHalf, merged)

	return merged
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	// Merge elements in sorted order
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// Add remaining elements
	for i < len(left) {
		result = append(result, left[i])
		i++
	}

	for j < len(right) {
		result = append(result, right[j])
		j++
	}

	return result
}

// =============================================================================
// QUICK SORT - O(n log n) Average, O(n²) Worst Case, O(log n) Space
// =============================================================================

// QuickSort implements the quicksort algorithm with random pivot
// Best for: Average case performance, in-place sorting, general purpose
// Time Complexity: O(n log n) average, O(n²) worst case
// Space Complexity: O(log n) due to recursion stack
func QuickSort(arr []int) []int {
	fmt.Println("⚡ Starting Quick Sort...")
	result := make([]int, len(arr))
	copy(result, arr)

	if len(result) <= 1 {
		return result
	}

	quickSortRecursive(result, 0, len(result)-1)
	fmt.Printf("✅ Quick Sort Complete: %v\n\n", result)
	return result
}

func quickSortRecursive(arr []int, low, high int) {
	if low < high {
		// Partition the array and get pivot index
		pivotIndex := partition(arr, low, high)
		fmt.Printf("Partitioned around pivot %d at index %d: %v\n",
			arr[pivotIndex], pivotIndex, arr[low:high+1])

		// Recursively sort elements before and after partition
		quickSortRecursive(arr, low, pivotIndex-1)
		quickSortRecursive(arr, pivotIndex+1, high)
	}
}

func partition(arr []int, low, high int) int {
	// Choose random pivot to avoid worst-case on sorted arrays
	randomIndex := low + rand.Intn(high-low+1)
	arr[randomIndex], arr[high] = arr[high], arr[randomIndex]

	pivot := arr[high] // Last element as pivot
	fmt.Printf("Chosen pivot: %d\n", pivot)

	i := low - 1 // Index of smaller element

	for j := low; j < high; j++ {
		// If current element is smaller than or equal to pivot
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
			fmt.Printf("  Swap: %d ↔ %d → %v\n", arr[j], arr[i], arr[low:high+1])
		}
	}

	// Place pivot in correct position
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// =============================================================================
// INSERTION SORT - O(n²) Time, O(1) Space, Efficient for Small Arrays
// =============================================================================

// InsertionSort implements insertion sort algorithm
// Best for: Small datasets, nearly sorted arrays, online algorithm
// Time Complexity: O(n²) worst case, O(n) best case
// Space Complexity: O(1)
func InsertionSort(arr []int) []int {
	fmt.Println("📝 Starting Insertion Sort...")
	result := make([]int, len(arr))
	copy(result, arr)

	for i := 1; i < len(result); i++ {
		key := result[i]
		j := i - 1

		fmt.Printf("Inserting %d: ", key)

		// Move elements greater than key one position ahead
		for j >= 0 && result[j] > key {
			result[j+1] = result[j]
			j--
		}

		result[j+1] = key
		fmt.Printf("%v\n", result)
	}

	fmt.Printf("✅ Insertion Sort Complete: %v\n\n", result)
	return result
}

// =============================================================================
// SELECTION SORT - O(n²) Time, O(1) Space, Minimal Swaps
// =============================================================================

// SelectionSort implements selection sort algorithm
// Best for: Situations where swap cost is high, small datasets
// Time Complexity: O(n²) in all cases
// Space Complexity: O(1)
func SelectionSort(arr []int) []int {
	fmt.Println("🎯 Starting Selection Sort...")
	result := make([]int, len(arr))
	copy(result, arr)
	n := len(result)

	for i := 0; i < n-1; i++ {
		minIndex := i

		// Find the minimum element in remaining unsorted array
		for j := i + 1; j < n; j++ {
			if result[j] < result[minIndex] {
				minIndex = j
			}
		}

		// Swap the found minimum element with the first element
		if minIndex != i {
			result[i], result[minIndex] = result[minIndex], result[i]
			fmt.Printf("Step %d: Swap %d ↔ %d → %v\n",
				i+1, result[minIndex], result[i], result)
		}
	}

	fmt.Printf("✅ Selection Sort Complete: %v\n\n", result)
	return result
}

// =============================================================================
// PERFORMANCE TESTING AND COMPARISON
// =============================================================================

// TimeSortingAlgorithm measures execution time of a sorting function
func TimeSortingAlgorithm(name string, sortFunc func([]int) []int, arr []int) time.Duration {
	start := time.Now()
	_ = sortFunc(arr)
	duration := time.Since(start)
	fmt.Printf("⏱️  %s took: %v\n", name, duration)
	return duration
}

// GenerateRandomArray creates a random array for testing
func GenerateRandomArray(size int, maxValue int) []int {
	rand.Seed(time.Now().UnixNano())
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(maxValue)
	}
	return arr
}

// GenerateWorstCaseArray creates reverse sorted array (worst case for some algorithms)
func GenerateWorstCaseArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = size - i
	}
	return arr
}

// GenerateBestCaseArray creates already sorted array (best case for some algorithms)
func GenerateBestCaseArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = i + 1
	}
	return arr
}

// IsSorted checks if an array is sorted
func IsSorted(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i-1] > arr[i] {
			return false
		}
	}
	return true
}

// CompareAllSortingAlgorithms tests all algorithms with the same input
func CompareAllSortingAlgorithms(arr []int) {
	fmt.Printf("\n🏁 SORTING ALGORITHM COMPARISON\n")
	fmt.Printf("Input Array: %v (length: %d)\n", arr, len(arr))
	fmt.Printf("%s\n", strings.Repeat("=", 60))

	algorithms := []struct {
		name string
		fn   func([]int) []int
	}{
		{"Bubble Sort", BubbleSort},
		{"Optimized Bubble Sort", OptimizedBubbleSort},
		{"Insertion Sort", InsertionSort},
		{"Selection Sort", SelectionSort},
		{"Merge Sort", MergeSort},
		{"Quick Sort", QuickSort},
	}

	results := make(map[string]time.Duration)

	for _, algo := range algorithms {
		fmt.Printf("\n--- Testing %s ---\n", algo.name)
		// Make a copy of the array for each algorithm
		testArray := make([]int, len(arr))
		copy(testArray, arr)

		start := time.Now()
		sorted := algo.fn(testArray)
		duration := time.Since(start)
		results[algo.name] = duration

		// Verify the result is sorted
		if IsSorted(sorted) {
			fmt.Printf("✅ %s: CORRECT (Time: %v)\n", algo.name, duration)
		} else {
			fmt.Printf("❌ %s: INCORRECT SORTING!\n", algo.name)
		}
	}

	fmt.Printf("\n📊 PERFORMANCE SUMMARY:\n")
	for name, duration := range results {
		fmt.Printf("%-25s: %v\n", name, duration)
	}
}

// =============================================================================
// MAIN DEMONSTRATION FUNCTION
// =============================================================================

func main() {
	fmt.Println("🚀 SORTING ALGORITHMS DEMONSTRATION")
	fmt.Println("Module 1: CS Fundamentals - Junior Software Engineer Interview Prep")
	fmt.Printf("%s\n\n", strings.Repeat("=", 70))

	// Test with a small array for educational purposes
	fmt.Println("📚 EDUCATIONAL DEMO - Small Array")
	smallArray := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("Original Array: %v\n\n", smallArray)

	// Demonstrate each algorithm step by step
	BubbleSort(append([]int{}, smallArray...))
	InsertionSort(append([]int{}, smallArray...))
	SelectionSort(append([]int{}, smallArray...))
	MergeSort(append([]int{}, smallArray...))
	QuickSort(append([]int{}, smallArray...))

	// Performance comparison with different scenarios
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("📊 PERFORMANCE TESTING")

	// Test 1: Random small array
	fmt.Println("\n🎲 Test 1: Random Array (Size: 10)")
	randomArray := GenerateRandomArray(10, 100)
	CompareAllSortingAlgorithms(randomArray)

	// Test 2: Already sorted array (best case for some algorithms)
	fmt.Println("\n✅ Test 2: Already Sorted Array (Size: 10)")
	sortedArray := GenerateBestCaseArray(10)
	CompareAllSortingAlgorithms(sortedArray)

	// Test 3: Reverse sorted array (worst case for some algorithms)
	fmt.Println("\n❌ Test 3: Reverse Sorted Array (Size: 10)")
	reverseArray := GenerateWorstCaseArray(10)
	CompareAllSortingAlgorithms(reverseArray)

	// Test 4: Performance with larger array (be careful with O(n²) algorithms)
	fmt.Println("\n🏋️ Test 4: Larger Array Performance (Size: 100)")
	largeArray := GenerateRandomArray(100, 1000)

	// Only test efficient algorithms with large arrays
	fmt.Printf("Testing efficient algorithms only...\n")
	TimeSortingAlgorithm("Merge Sort", MergeSort, append([]int{}, largeArray...))
	TimeSortingAlgorithm("Quick Sort", QuickSort, append([]int{}, largeArray...))

	fmt.Println("\n🎯 KEY TAKEAWAYS:")
	fmt.Println("• Bubble Sort: Simple but inefficient O(n²) - good for learning")
	fmt.Println("• Insertion Sort: Efficient for small/nearly sorted arrays")
	fmt.Println("• Selection Sort: Minimal swaps, but still O(n²)")
	fmt.Println("• Merge Sort: Guaranteed O(n log n), stable, uses extra space")
	fmt.Println("• Quick Sort: Average O(n log n), in-place, can be O(n²) worst case")
	fmt.Println("\n💡 For interviews: Focus on Merge Sort and Quick Sort implementations!")
}
