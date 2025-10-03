// Searching Algorithms Implementation - Junior Software Engineer Interview Prep
// Module 1: CS Fundamentals - Basic Algorithms
// Complete Golang implementations with educational focus

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// =============================================================================
// LINEAR SEARCH - O(n) Time Complexity, O(1) Space Complexity
// =============================================================================

// LinearSearch searches for a target value in an unsorted array
// Best for: Small arrays, unsorted data, simple implementation
// Time Complexity: O(n) - must check every element in worst case
// Space Complexity: O(1) - no extra space needed
func LinearSearch(arr []int, target int) (int, bool) {
	fmt.Printf("🔍 Linear Search for %d in %v\n", target, arr)

	for i, value := range arr {
		fmt.Printf("  Step %d: Checking arr[%d] = %d", i+1, i, value)
		if value == target {
			fmt.Printf(" ✅ FOUND!\n")
			fmt.Printf("📍 Linear Search Result: Found %d at index %d\n\n", target, i)
			return i, true
		}
		fmt.Printf(" ❌ Continue...\n")
	}

	fmt.Printf("❌ Linear Search Result: %d not found in array\n\n", target)
	return -1, false
}

// LinearSearchWithCount returns index, found status, and comparison count
func LinearSearchWithCount(arr []int, target int) (int, bool, int) {
	comparisons := 0

	for i, value := range arr {
		comparisons++
		if value == target {
			return i, true, comparisons
		}
	}

	return -1, false, comparisons
}

// =============================================================================
// BINARY SEARCH - O(log n) Time Complexity, O(1) Space Complexity
// =============================================================================

// BinarySearch searches for a target value in a SORTED array
// PREREQUISITE: Array must be sorted!
// Best for: Large sorted arrays, frequent searches
// Time Complexity: O(log n) - eliminates half the search space each step
// Space Complexity: O(1) - iterative version uses constant space
func BinarySearch(arr []int, target int) (int, bool) {
	fmt.Printf("🎯 Binary Search for %d in %v\n", target, arr)

	left, right := 0, len(arr)-1
	step := 1

	for left <= right {
		mid := left + (right-left)/2 // Prevents integer overflow
		midValue := arr[mid]

		fmt.Printf("  Step %d: left=%d, right=%d, mid=%d, arr[mid]=%d\n",
			step, left, right, mid, midValue)

		if midValue == target {
			fmt.Printf("✅ Binary Search Result: Found %d at index %d\n\n", target, mid)
			return mid, true
		} else if midValue < target {
			fmt.Printf("    %d < %d, search right half\n", midValue, target)
			left = mid + 1
		} else {
			fmt.Printf("    %d > %d, search left half\n", midValue, target)
			right = mid - 1
		}
		step++
	}

	fmt.Printf("❌ Binary Search Result: %d not found in array\n\n", target)
	return -1, false
}

// BinarySearchRecursive implements binary search using recursion
// Time Complexity: O(log n)
// Space Complexity: O(log n) - due to recursive call stack
func BinarySearchRecursive(arr []int, target int) (int, bool) {
	fmt.Printf("🔄 Recursive Binary Search for %d in %v\n", target, arr)
	index := binarySearchRecursiveHelper(arr, target, 0, len(arr)-1, 1)

	if index != -1 {
		fmt.Printf("✅ Recursive Binary Search Result: Found %d at index %d\n\n", target, index)
		return index, true
	}

	fmt.Printf("❌ Recursive Binary Search Result: %d not found\n\n", target)
	return -1, false
}

func binarySearchRecursiveHelper(arr []int, target, left, right, step int) int {
	if left > right {
		return -1
	}

	mid := left + (right-left)/2
	midValue := arr[mid]

	fmt.Printf("  Step %d: left=%d, right=%d, mid=%d, arr[mid]=%d\n",
		step, left, right, mid, midValue)

	if midValue == target {
		return mid
	} else if midValue < target {
		fmt.Printf("    %d < %d, search right half\n", midValue, target)
		return binarySearchRecursiveHelper(arr, target, mid+1, right, step+1)
	} else {
		fmt.Printf("    %d > %d, search left half\n", midValue, target)
		return binarySearchRecursiveHelper(arr, target, left, mid-1, step+1)
	}
}

// =============================================================================
// BINARY SEARCH VARIATIONS - Common Interview Patterns
// =============================================================================

// FindFirstOccurrence finds the first occurrence of target in sorted array with duplicates
// Example: [1, 2, 2, 2, 3, 4] target=2 → returns index 1 (first 2)
func FindFirstOccurrence(arr []int, target int) int {
	fmt.Printf("🎯 Finding FIRST occurrence of %d in %v\n", target, arr)

	left, right := 0, len(arr)-1
	result := -1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			result = mid    // Found target, but continue searching left
			right = mid - 1 // Look for earlier occurrence
			fmt.Printf("  Found at %d, continue searching left for first occurrence\n", mid)
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	if result != -1 {
		fmt.Printf("✅ First occurrence of %d found at index %d\n\n", target, result)
	} else {
		fmt.Printf("❌ %d not found in array\n\n", target)
	}

	return result
}

// FindLastOccurrence finds the last occurrence of target in sorted array with duplicates
// Example: [1, 2, 2, 2, 3, 4] target=2 → returns index 3 (last 2)
func FindLastOccurrence(arr []int, target int) int {
	fmt.Printf("🎯 Finding LAST occurrence of %d in %v\n", target, arr)

	left, right := 0, len(arr)-1
	result := -1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			result = mid   // Found target, but continue searching right
			left = mid + 1 // Look for later occurrence
			fmt.Printf("  Found at %d, continue searching right for last occurrence\n", mid)
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	if result != -1 {
		fmt.Printf("✅ Last occurrence of %d found at index %d\n\n", target, result)
	} else {
		fmt.Printf("❌ %d not found in array\n\n", target)
	}

	return result
}

// FindInsertionPosition finds where target should be inserted to maintain sorted order
// This is the "lower bound" - first position where we could insert target
func FindInsertionPosition(arr []int, target int) int {
	fmt.Printf("📍 Finding insertion position for %d in %v\n", target, arr)

	left, right := 0, len(arr)

	for left < right {
		mid := left + (right-left)/2

		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}

	fmt.Printf("✅ Insertion position for %d: index %d\n", target, left)
	fmt.Printf("   Result array would be: %v\n\n", insertAtPosition(arr, target, left))

	return left
}

// Helper function to show what array would look like after insertion
func insertAtPosition(arr []int, value, position int) []int {
	result := make([]int, len(arr)+1)
	copy(result[:position], arr[:position])
	result[position] = value
	copy(result[position+1:], arr[position:])
	return result
}

// =============================================================================
// SEARCH IN ROTATED SORTED ARRAY - Advanced Binary Search
// =============================================================================

// SearchInRotatedSortedArray searches in a rotated sorted array
// Example: [4, 5, 6, 7, 0, 1, 2] is a rotated version of [0, 1, 2, 4, 5, 6, 7]
// This is a common interview question!
func SearchInRotatedSortedArray(arr []int, target int) int {
	fmt.Printf("🌀 Searching for %d in rotated sorted array: %v\n", target, arr)

	left, right := 0, len(arr)-1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			fmt.Printf("✅ Found %d at index %d\n\n", target, mid)
			return mid
		}

		// Determine which half is properly sorted
		if arr[left] <= arr[mid] {
			// Left half is sorted
			fmt.Printf("  Left half [%d:%d] is sorted\n", left, mid)
			if arr[left] <= target && target < arr[mid] {
				fmt.Printf("  Target %d is in left half\n", target)
				right = mid - 1
			} else {
				fmt.Printf("  Target %d is in right half\n", target)
				left = mid + 1
			}
		} else {
			// Right half is sorted
			fmt.Printf("  Right half [%d:%d] is sorted\n", mid, right)
			if arr[mid] < target && target <= arr[right] {
				fmt.Printf("  Target %d is in right half\n", target)
				left = mid + 1
			} else {
				fmt.Printf("  Target %d is in left half\n", target)
				right = mid - 1
			}
		}
	}

	fmt.Printf("❌ %d not found in rotated sorted array\n\n", target)
	return -1
}

// =============================================================================
// PERFORMANCE TESTING AND COMPARISON
// =============================================================================

// TimeSearchAlgorithm measures execution time of a search function
func TimeSearchAlgorithm(name string, searchFunc func([]int, int) (int, bool), arr []int, target int) time.Duration {
	start := time.Now()
	_, _ = searchFunc(arr, target)
	duration := time.Since(start)
	fmt.Printf("⏱️  %s took: %v\n", name, duration)
	return duration
}

// GenerateRandomSortedArray creates a sorted array for binary search testing
func GenerateRandomSortedArray(size int, maxValue int) []int {
	rand.Seed(time.Now().UnixNano())
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(maxValue)
	}
	sort.Ints(arr) // Sort the array for binary search
	return arr
}

// CompareSearchAlgorithms tests linear vs binary search performance
func CompareSearchAlgorithms(arr []int, target int) {
	fmt.Printf("\n🏁 SEARCH ALGORITHM COMPARISON\n")
	fmt.Printf("Array size: %d, Target: %d\n", len(arr), target)
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	// Linear search (works on any array)
	fmt.Println("\n--- Linear Search ---")
	_, found1, comparisons1 := LinearSearchWithCount(arr, target)

	// Binary search (requires sorted array)
	fmt.Println("\n--- Binary Search ---")
	_, found2 := BinarySearch(arr, target)

	fmt.Printf("\n📊 COMPARISON RESULTS:\n")
	fmt.Printf("Linear Search: %d comparisons, Found: %v\n", comparisons1, found1)
	fmt.Printf("Binary Search: ~%d comparisons (log₂(%d)), Found: %v\n",
		int(logBase2(float64(len(arr))))+1, len(arr), found2)

	if len(arr) > 10 {
		improvement := float64(comparisons1) / (logBase2(float64(len(arr))) + 1)
		fmt.Printf("Binary Search is ~%.1fx faster for this size!\n", improvement)
	}
}

// Simple log base 2 calculation
func logBase2(n float64) float64 {
	if n <= 1 {
		return 0
	}
	return 1 + logBase2(n/2)
}

// =============================================================================
// PRACTICAL APPLICATIONS AND EXAMPLES
// =============================================================================

// DemonstrateBinarySearchVariations shows different binary search use cases
func DemonstrateBinarySearchVariations() {
	fmt.Println("🎓 BINARY SEARCH VARIATIONS - Interview Patterns")
	fmt.Printf("%s\n", strings.Repeat("=", 60))

	// Array with duplicates for testing first/last occurrence
	duplicateArray := []int{1, 2, 2, 2, 2, 3, 4, 4, 5}
	fmt.Printf("Test Array (with duplicates): %v\n\n", duplicateArray)

	target := 2
	FindFirstOccurrence(duplicateArray, target)
	FindLastOccurrence(duplicateArray, target)

	// Test insertion position
	sortedArray := []int{1, 3, 5, 7, 9}
	fmt.Printf("Sorted Array: %v\n", sortedArray)
	FindInsertionPosition(sortedArray, 4)
	FindInsertionPosition(sortedArray, 0)
	FindInsertionPosition(sortedArray, 10)

	// Test rotated sorted array
	rotatedArray := []int{4, 5, 6, 7, 0, 1, 2}
	SearchInRotatedSortedArray(rotatedArray, 0)
	SearchInRotatedSortedArray(rotatedArray, 3)
}

// =============================================================================
// INTERVIEW TIPS AND PATTERNS
// =============================================================================

func PrintInterviewTips() {
	fmt.Println("\n💡 INTERVIEW TIPS FOR SEARCHING ALGORITHMS")
	fmt.Printf("%s\n", strings.Repeat("=", 60))

	tips := []string{
		"Always ask if the array is sorted - this determines which search to use",
		"Binary search requires sorted data - don't forget to mention this!",
		"Watch out for integer overflow: use mid = left + (right-left)/2",
		"For duplicates, clarify if you need first, last, or any occurrence",
		"Rotated sorted array is a common variation - practice the pattern",
		"Binary search template: while left <= right (inclusive bounds)",
		"Time complexity: Linear O(n), Binary O(log n)",
		"Space complexity: Both can be O(1) with iterative implementation",
	}

	for i, tip := range tips {
		fmt.Printf("  %d. %s\n", i+1, tip)
	}

	fmt.Println("\n🔥 COMMON BINARY SEARCH MISTAKES TO AVOID:")
	mistakes := []string{
		"Using binary search on unsorted array",
		"Integer overflow with (left + right) / 2",
		"Infinite loops due to incorrect boundary updates",
		"Off-by-one errors in boundary conditions",
		"Not handling empty array edge case",
	}

	for i, mistake := range mistakes {
		fmt.Printf("  ❌ %d. %s\n", i+1, mistake)
	}
}

// =============================================================================
// MAIN DEMONSTRATION FUNCTION
// =============================================================================

func main() {
	fmt.Println("🔍 SEARCHING ALGORITHMS DEMONSTRATION")
	fmt.Println("Module 1: CS Fundamentals - Junior Software Engineer Interview Prep")
	fmt.Printf("%s\n\n", strings.Repeat("=", 70))

	// Basic demonstration with small arrays
	fmt.Println("📚 BASIC SEARCHING DEMO")

	// Unsorted array for linear search
	unsortedArray := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("Unsorted Array: %v\n\n", unsortedArray)

	target := 22
	LinearSearch(unsortedArray, target)
	LinearSearch(unsortedArray, 99) // Not found case

	// Sorted array for binary search
	sortedArray := []int{11, 12, 22, 25, 34, 64, 90}
	fmt.Printf("Sorted Array: %v\n\n", sortedArray)

	BinarySearch(sortedArray, target)
	BinarySearch(sortedArray, 99) // Not found case

	BinarySearchRecursive(sortedArray, 25)

	// Advanced binary search variations
	DemonstrateBinarySearchVariations()

	// Performance comparison
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("📊 PERFORMANCE ANALYSIS")

	// Small array comparison
	fmt.Println("\n🔬 Small Array Test (Size: 10)")
	smallSorted := GenerateRandomSortedArray(10, 100)
	testTarget := smallSorted[5] // Guaranteed to exist
	CompareSearchAlgorithms(smallSorted, testTarget)

	// Large array comparison (demonstrate the power of binary search)
	fmt.Println("\n🏋️ Large Array Test (Size: 1000)")
	largeSorted := GenerateRandomSortedArray(1000, 10000)
	testTarget2 := largeSorted[500] // Guaranteed to exist

	fmt.Printf("Testing with array of size %d...\n", len(largeSorted))

	// Time both algorithms
	start := time.Now()
	LinearSearchWithCount(largeSorted, testTarget2)
	linearTime := time.Since(start)

	start = time.Now()
	BinarySearch(largeSorted, testTarget2)
	binaryTime := time.Since(start)

	fmt.Printf("\n⚡ PERFORMANCE RESULTS:\n")
	fmt.Printf("Linear Search: %v\n", linearTime)
	fmt.Printf("Binary Search: %v\n", binaryTime)

	if linearTime > binaryTime {
		improvement := float64(linearTime) / float64(binaryTime)
		fmt.Printf("Binary Search was %.1fx faster!\n", improvement)
	}

	// Interview tips
	PrintInterviewTips()

	fmt.Println("\n🎯 KEY TAKEAWAYS:")
	fmt.Println("• Linear Search: O(n) - works on any array, simple to implement")
	fmt.Println("• Binary Search: O(log n) - requires sorted array, much faster for large data")
	fmt.Println("• Always clarify array properties (sorted? duplicates?) in interviews")
	fmt.Println("• Practice binary search variations - they're common interview questions")
	fmt.Println("• Watch out for integer overflow and boundary conditions")
	fmt.Println("\n💪 Master these patterns and you'll crush searching algorithm interviews!")
}
