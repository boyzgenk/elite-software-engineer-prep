# 📚 Module 2: Basic Algorithms & Data Structures  
## Warming-Up Track - Programming Fundamentals

### 🎯 Learning Objectives
By the end of this section, you will:
- [ ] Understand and implement fundamental sorting algorithms
- [ ] Master binary search and its variations
- [ ] Recognize when to apply different sorting/searching strategies
- [ ] Analyze time and space complexity of algorithms
- [ ] Solve 8-10 algorithm-based problems confidently

---

## 📚 Sorting Algorithms

### 🔸 Why Learn Sorting?

**Importance:**
- Foundation for many advanced algorithms
- Demonstrates divide-and-conquer and other paradigms
- Optimizes search operations (sorted data enables binary search)
- Common interview topic across all levels

**When Sorting is Useful:**
- Preprocessing data for faster searches
- Finding duplicates or ranges
- Implementing efficient algorithms (merge operations)
- Data analysis and statistics

### 🔸 Elementary Sorting Algorithms

#### Bubble Sort - O(n²)
**Concept:** Repeatedly compare adjacent elements and swap if wrong order
**Use Case:** Educational purposes, nearly sorted data

```go
package main

import "fmt"

// BubbleSort implements bubble sort algorithm - O(n²) time, O(1) space
func BubbleSort(arr []int) {
    n := len(arr)
    
    for i := 0; i < n-1; i++ {
        swapped := false
        
        // Last i elements are already in place
        for j := 0; j < n-i-1; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
                swapped = true
            }
        }
        
        // If no swapping occurred, array is sorted
        if !swapped {
            break
        }
    }
}

// Optimized bubble sort with early termination
func BubbleSortOptimized(arr []int) {
    n := len(arr)
    
    for i := 0; i < n-1; i++ {
        swapped := false
        
        for j := 0; j < n-i-1; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
                swapped = true
            }
        }
        
        if !swapped {
            fmt.Printf("Array sorted after %d passes\n", i+1)
            break
        }
    }
}
```

#### Selection Sort - O(n²)
**Concept:** Find minimum element and place it at beginning, repeat for rest
**Use Case:** Small datasets, when memory writes are expensive

```go
// SelectionSort implements selection sort - O(n²) time, O(1) space
func SelectionSort(arr []int) {
    n := len(arr)
    
    for i := 0; i < n-1; i++ {
        minIndex := i
        
        // Find minimum element in remaining array
        for j := i + 1; j < n; j++ {
            if arr[j] < arr[minIndex] {
                minIndex = j
            }
        }
        
        // Swap minimum element with first element
        if minIndex != i {
            arr[i], arr[minIndex] = arr[minIndex], arr[i]
        }
    }
}
```

#### Insertion Sort - O(n²)
**Concept:** Build sorted array one element at a time, like sorting playing cards
**Use Case:** Small arrays, nearly sorted data, online algorithms

```go
// InsertionSort implements insertion sort - O(n²) time, O(1) space
func InsertionSort(arr []int) {
    n := len(arr)
    
    for i := 1; i < n; i++ {
        key := arr[i]
        j := i - 1
        
        // Move elements greater than key one position ahead
        for j >= 0 && arr[j] > key {
            arr[j+1] = arr[j]
            j--
        }
        
        arr[j+1] = key
    }
}

// InsertionSortBinary uses binary search to find insertion position
func InsertionSortBinary(arr []int) {
    for i := 1; i < len(arr); i++ {
        key := arr[i]
        
        // Find location to insert using binary search
        pos := binarySearch(arr, 0, i-1, key)
        
        // Shift elements to make space
        for j := i; j > pos; j-- {
            arr[j] = arr[j-1]
        }
        
        arr[pos] = key
    }
}

func binarySearch(arr []int, left, right, key int) int {
    if left >= right {
        if arr[left] > key {
            return left
        }
        return left + 1
    }
    
    mid := left + (right-left)/2
    
    if arr[mid] == key {
        return mid + 1
    } else if arr[mid] > key {
        return binarySearch(arr, left, mid-1, key)
    } else {
        return binarySearch(arr, mid+1, right, key)
    }
}
```

### 🔸 Efficient Sorting Algorithms

#### Merge Sort - O(n log n)
**Concept:** Divide array into halves, sort recursively, then merge
**Use Case:** Large datasets, stable sorting required, linked lists

```go
// MergeSort implements merge sort - O(n log n) time, O(n) space
func MergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    mid := len(arr) / 2
    left := MergeSort(arr[:mid])
    right := MergeSort(arr[mid:])
    
    return merge(left, right)
}

// merge combines two sorted arrays into one sorted array
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
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    
    return result
}

// MergeSortInPlace sorts array in-place to save space
func MergeSortInPlace(arr []int) {
    if len(arr) <= 1 {
        return
    }
    
    mergeSortHelper(arr, 0, len(arr)-1)
}

func mergeSortHelper(arr []int, left, right int) {
    if left >= right {
        return
    }
    
    mid := left + (right-left)/2
    mergeSortHelper(arr, left, mid)
    mergeSortHelper(arr, mid+1, right)
    mergeInPlace(arr, left, mid, right)
}

func mergeInPlace(arr []int, left, mid, right int) {
    // Create temporary arrays for left and right subarrays
    leftArr := make([]int, mid-left+1)
    rightArr := make([]int, right-mid)
    
    copy(leftArr, arr[left:mid+1])
    copy(rightArr, arr[mid+1:right+1])
    
    i, j, k := 0, 0, left
    
    // Merge back into original array
    for i < len(leftArr) && j < len(rightArr) {
        if leftArr[i] <= rightArr[j] {
            arr[k] = leftArr[i]
            i++
        } else {
            arr[k] = rightArr[j]
            j++
        }
        k++
    }
    
    // Copy remaining elements
    for i < len(leftArr) {
        arr[k] = leftArr[i]
        i++
        k++
    }
    
    for j < len(rightArr) {
        arr[k] = rightArr[j]
        j++
        k++
    }
}
```

#### Quick Sort - O(n log n) average, O(n²) worst
**Concept:** Choose pivot, partition around it, recursively sort partitions
**Use Case:** General purpose, in-place sorting, when average case is acceptable

```go
// QuickSort implements quicksort algorithm
func QuickSort(arr []int) {
    if len(arr) <= 1 {
        return
    }
    quickSortHelper(arr, 0, len(arr)-1)
}

func quickSortHelper(arr []int, low, high int) {
    if low < high {
        // Partition the array and get pivot index
        pivotIndex := partition(arr, low, high)
        
        // Recursively sort elements before and after partition
        quickSortHelper(arr, low, pivotIndex-1)
        quickSortHelper(arr, pivotIndex+1, high)
    }
}

// partition rearranges array so elements <= pivot are on left
func partition(arr []int, low, high int) int {
    // Choose rightmost element as pivot
    pivot := arr[high]
    i := low - 1 // Index of smaller element
    
    for j := low; j < high; j++ {
        if arr[j] <= pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    
    // Place pivot in correct position
    arr[i+1], arr[high] = arr[high], arr[i+1]
    return i + 1
}

// QuickSortRandomized uses random pivot for better average case
func QuickSortRandomized(arr []int) {
    if len(arr) <= 1 {
        return
    }
    quickSortRandomizedHelper(arr, 0, len(arr)-1)
}

func quickSortRandomizedHelper(arr []int, low, high int) {
    if low < high {
        // Randomly choose pivot
        randomIndex := low + rand.Intn(high-low+1)
        arr[randomIndex], arr[high] = arr[high], arr[randomIndex]
        
        pivotIndex := partition(arr, low, high)
        quickSortRandomizedHelper(arr, low, pivotIndex-1)
        quickSortRandomizedHelper(arr, pivotIndex+1, high)
    }
}
```

---

## 🔍 Searching Algorithms

### 🔸 Linear Search - O(n)
**Use Case:** Unsorted data, small datasets, finding all occurrences

```go
// LinearSearch finds first occurrence of target - O(n) time, O(1) space
func LinearSearch(arr []int, target int) int {
    for i, value := range arr {
        if value == target {
            return i
        }
    }
    return -1 // Not found
}

// LinearSearchAll finds all occurrences of target
func LinearSearchAll(arr []int, target int) []int {
    indices := []int{}
    for i, value := range arr {
        if value == target {
            indices = append(indices, i)
        }
    }
    return indices
}
```

### 🔸 Binary Search - O(log n)
**Prerequisite:** Array must be sorted
**Use Case:** Large sorted datasets, range queries, optimization problems

```go
// BinarySearch finds target in sorted array - O(log n) time, O(1) space
func BinarySearch(arr []int, target int) int {
    left, right := 0, len(arr)-1
    
    for left <= right {
        mid := left + (right-left)/2 // Avoid overflow
        
        if arr[mid] == target {
            return mid
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    
    return -1 // Not found
}

// BinarySearchRecursive implements recursive binary search
func BinarySearchRecursive(arr []int, target int) int {
    return binarySearchHelper(arr, target, 0, len(arr)-1)
}

func binarySearchHelper(arr []int, target, left, right int) int {
    if left > right {
        return -1
    }
    
    mid := left + (right-left)/2
    
    if arr[mid] == target {
        return mid
    } else if arr[mid] < target {
        return binarySearchHelper(arr, target, mid+1, right)
    } else {
        return binarySearchHelper(arr, target, left, mid-1)
    }
}

// BinarySearchLeftmost finds leftmost occurrence of target
func BinarySearchLeftmost(arr []int, target int) int {
    left, right := 0, len(arr)-1
    result := -1
    
    for left <= right {
        mid := left + (right-left)/2
        
        if arr[mid] == target {
            result = mid
            right = mid - 1 // Continue searching left
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    
    return result
}

// BinarySearchRightmost finds rightmost occurrence of target
func BinarySearchRightmost(arr []int, target int) int {
    left, right := 0, len(arr)-1
    result := -1
    
    for left <= right {
        mid := left + (right-left)/2
        
        if arr[mid] == target {
            result = mid
            left = mid + 1 // Continue searching right
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    
    return result
}

// BinarySearchInsertPosition finds position where target should be inserted
func BinarySearchInsertPosition(arr []int, target int) int {
    left, right := 0, len(arr)
    
    for left < right {
        mid := left + (right-left)/2
        
        if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid
        }
    }
    
    return left
}
```

### 🔸 Binary Search Variations

```go
// SearchInRotatedArray searches in rotated sorted array
func SearchInRotatedArray(nums []int, target int) int {
    left, right := 0, len(nums)-1
    
    for left <= right {
        mid := left + (right-left)/2
        
        if nums[mid] == target {
            return mid
        }
        
        // Determine which half is sorted
        if nums[left] <= nums[mid] {
            // Left half is sorted
            if nums[left] <= target && target < nums[mid] {
                right = mid - 1
            } else {
                left = mid + 1
            }
        } else {
            // Right half is sorted
            if nums[mid] < target && target <= nums[right] {
                left = mid + 1
            } else {
                right = mid - 1
            }
        }
    }
    
    return -1
}

// FindPeakElement finds any peak element (element greater than neighbors)
func FindPeakElement(nums []int) int {
    left, right := 0, len(nums)-1
    
    for left < right {
        mid := left + (right-left)/2
        
        if nums[mid] > nums[mid+1] {
            // Peak is in left half (including mid)
            right = mid
        } else {
            // Peak is in right half
            left = mid + 1
        }
    }
    
    return left
}

// FindMinInRotatedArray finds minimum element in rotated sorted array
func FindMinInRotatedArray(nums []int) int {
    left, right := 0, len(nums)-1
    
    for left < right {
        mid := left + (right-left)/2
        
        if nums[mid] > nums[right] {
            // Minimum is in right half
            left = mid + 1
        } else {
            // Minimum is in left half (including mid)
            right = mid
        }
    }
    
    return nums[left]
}
```

---

## 📊 Algorithm Comparison

### Sorting Algorithms Comparison

| Algorithm | Best Case | Average Case | Worst Case | Space | Stable | Notes |
|-----------|-----------|--------------|------------|-------|--------|-------|
| Bubble Sort | O(n) | O(n²) | O(n²) | O(1) | Yes | Simple, good for small/nearly sorted |
| Selection Sort | O(n²) | O(n²) | O(n²) | O(1) | No | Consistent performance |
| Insertion Sort | O(n) | O(n²) | O(n²) | O(1) | Yes | Efficient for small/nearly sorted |
| Merge Sort | O(n log n) | O(n log n) | O(n log n) | O(n) | Yes | Consistent, good for large data |
| Quick Sort | O(n log n) | O(n log n) | O(n²) | O(log n) | No | Fast average case, in-place |

### When to Use Which Sorting Algorithm

**Use Insertion Sort when:**
- Array size < 50
- Array is nearly sorted
- Need stable sort with minimal extra space

**Use Merge Sort when:**
- Need guaranteed O(n log n) performance
- Stability is required
- Have extra memory available
- Sorting linked lists

**Use Quick Sort when:**
- Average case performance is most important
- Memory is limited (in-place sorting)
- Array is randomly ordered

**Use Built-in Sort when:**
- Production code (optimized implementations)
- No specific requirements mandate custom implementation

---

## 🎯 Common Algorithmic Patterns

### Pattern 1: Divide and Conquer

```go
// FindMaxSubarraySum finds maximum sum of contiguous subarray (Kadane's algorithm)
func FindMaxSubarraySum(arr []int) int {
    if len(arr) == 0 {
        return 0
    }
    
    maxSoFar := arr[0]
    maxEndingHere := arr[0]
    
    for i := 1; i < len(arr); i++ {
        maxEndingHere = max(arr[i], maxEndingHere+arr[i])
        maxSoFar = max(maxSoFar, maxEndingHere)
    }
    
    return maxSoFar
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// MergeKSortedArrays merges k sorted arrays into one sorted array
func MergeKSortedArrays(arrays [][]int) []int {
    if len(arrays) == 0 {
        return []int{}
    }
    if len(arrays) == 1 {
        return arrays[0]
    }
    
    // Divide arrays into two halves
    mid := len(arrays) / 2
    left := MergeKSortedArrays(arrays[:mid])
    right := MergeKSortedArrays(arrays[mid:])
    
    // Merge the two halves
    return merge(left, right)
}
```

### Pattern 2: Two Pointers with Sorted Data

```go
// TwoSumSorted finds two numbers in sorted array that sum to target
func TwoSumSorted(arr []int, target int) []int {
    left, right := 0, len(arr)-1
    
    for left < right {
        sum := arr[left] + arr[right]
        if sum == target {
            return []int{left, right}
        } else if sum < target {
            left++
        } else {
            right--
        }
    }
    
    return []int{} // No solution found
}

// ThreeSum finds all unique triplets that sum to zero
func ThreeSum(nums []int) [][]int {
    result := [][]int{}
    if len(nums) < 3 {
        return result
    }
    
    // Sort array first
    sort.Ints(nums)
    
    for i := 0; i < len(nums)-2; i++ {
        // Skip duplicates for first number
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }
        
        left, right := i+1, len(nums)-1
        target := -nums[i]
        
        for left < right {
            sum := nums[left] + nums[right]
            if sum == target {
                result = append(result, []int{nums[i], nums[left], nums[right]})
                
                // Skip duplicates
                for left < right && nums[left] == nums[left+1] {
                    left++
                }
                for left < right && nums[right] == nums[right-1] {
                    right--
                }
                
                left++
                right--
            } else if sum < target {
                left++
            } else {
                right--
            }
        }
    }
    
    return result
}
```

---

## 💻 Practice Problems

### Easy Level (Day 1)

#### Problem 1: Sort Colors (LeetCode #75)
```go
func sortColors(nums []int) {
    // Dutch National Flag algorithm
    // Sort array with values 0, 1, 2
}
```

#### Problem 2: Merge Sorted Array (LeetCode #88)
```go
func merge(nums1 []int, m int, nums2 []int, n int) {
    // Merge nums2 into nums1 in-place
}
```

### Medium Level (Day 2-3)

#### Problem 3: Search in Rotated Sorted Array (LeetCode #33)
```go
func search(nums []int, target int) int {
    // Binary search in rotated array
    return -1
}
```

#### Problem 4: Find First and Last Position (LeetCode #34)
```go
func searchRange(nums []int, target int) []int {
    // Find first and last occurrence using binary search
    return []int{-1, -1}
}
```

#### Problem 5: Kth Largest Element (LeetCode #215)
```go
func findKthLargest(nums []int, k int) int {
    // Find kth largest using quickselect or heap
    return 0
}
```

---

## 🚀 Next Steps

### Day 1 Goals (Basic Sorting)
- [ ] Implement bubble, selection, and insertion sort
- [ ] Understand when to use each elementary sort
- [ ] Practice sorting small arrays by hand
- [ ] Analyze time complexity of each algorithm

### Day 2 Goals (Advanced Sorting)
- [ ] Implement merge sort (recursive and iterative)
- [ ] Master quick sort with different pivot strategies
- [ ] Understand stability and in-place sorting concepts
- [ ] Practice divide-and-conquer thinking

### Day 3 Goals (Searching)
- [ ] Master binary search and its variations
- [ ] Solve rotated array search problems
- [ ] Practice range finding and insertion position
- [ ] Understand when binary search applies

### Preparation for Recursion (Days 4-5)
- Review function call stack concepts
- Practice breaking problems into smaller subproblems
- Understand base cases and recursive relations
- Think about tree-like problem decomposition

**Algorithms are the tools that turn brute force into elegant solutions! Master these fundamentals and you'll recognize patterns everywhere. This is strategic thinking in code!** ⚡

---

## 🔧 Algorithm Implementation Tips

### Sorting Tips:
```go
// Always test with edge cases
testCases := [][]int{
    {},                    // Empty array
    {1},                   // Single element
    {2, 1},               // Two elements
    {3, 1, 2},            // Small array
    {5, 2, 8, 1, 9},      // Random order
    {1, 2, 3, 4, 5},      // Already sorted
    {5, 4, 3, 2, 1},      // Reverse sorted
    {3, 3, 3, 3},         // All equal
}

// Verify sorting correctness
func isSorted(arr []int) bool {
    for i := 1; i < len(arr); i++ {
        if arr[i] < arr[i-1] {
            return false
        }
    }
    return true
}
```

### Binary Search Tips:
```go
// Common binary search template
func binarySearchTemplate(arr []int, target int) int {
    left, right := 0, len(arr)-1
    
    for left <= right {
        mid := left + (right-left)/2 // Prevent overflow
        
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

// For finding insertion position, use left < right
func insertPosition(arr []int, target int) int {
    left, right := 0, len(arr)
    
    for left < right {
        mid := left + (right-left)/2
        
        if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid
        }
    }
    
    return left
}
```

**Remember: Good algorithms are like good strategies - they work efficiently even when the problem size grows. Think big picture, code small steps!** 💪