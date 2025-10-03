# Searching Algorithms Practice Problems
## Module 1: CS Fundamentals - Basic Algorithms

### 🎯 Learning Objectives
By completing these problems, you will:
- Master linear and binary search implementations
- Understand when to apply different searching strategies
- Solve advanced binary search variations (common in interviews)
- Develop intuition for divide-and-conquer problem solving

---

## 📚 Problem Categories

### **EASY LEVEL (Foundation Building)**

#### Problem 1: Linear Search with Multiple Occurrences
**Difficulty:** Easy | **Time:** 15-20 minutes

**Problem Statement:**
Implement linear search that finds ALL occurrences of a target value and returns their indices.

```go
func LinearSearchAll(arr []int, target int) []int {
    // Return slice of all indices where target is found
    // Return empty slice if target not found
}
```

**Test Cases:**
```go
arr1 := []int{3, 5, 2, 5, 8, 5, 1}
target1 := 5  // Expected: [1, 3, 5]

arr2 := []int{1, 2, 3, 4, 5}
target2 := 6  // Expected: []

arr3 := []int{7, 7, 7, 7}
target3 := 7  // Expected: [0, 1, 2, 3]
```

---

#### Problem 2: Binary Search Implementation
**Difficulty:** Easy | **Time:** 20-25 minutes

**Problem Statement:**
Implement standard binary search with detailed step tracking.

```go
func BinarySearchWithSteps(arr []int, target int) (int, int) {
    // Return index (or -1 if not found) and number of comparisons made
    // Print each step of the search process
}
```

**Test Cases:**
```go
arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

// Test existing elements
BinarySearchWithSteps(arr, 7)   // Should find at index 3
BinarySearchWithSteps(arr, 1)   // Should find at index 0  
BinarySearchWithSteps(arr, 19)  // Should find at index 9

// Test non-existing elements
BinarySearchWithSteps(arr, 4)   // Should return -1
BinarySearchWithSteps(arr, 20)  // Should return -1
```

**Expected Output Format:**
```
Step 1: left=0, right=9, mid=4, arr[4]=9, target=7, search left
Step 2: left=0, right=3, mid=1, arr[1]=3, target=7, search right
Step 3: left=2, right=3, mid=2, arr[2]=5, target=7, search right
Step 4: left=3, right=3, mid=3, arr[3]=7, FOUND!
Found 7 at index 3 using 4 comparisons
```

---

#### Problem 3: Search in Array with Unknown Size
**Difficulty:** Easy-Medium | **Time:** 25-30 minutes

**Problem Statement:**
You have a sorted array but don't know its size. The array has a special `get(index)` function that returns the element at index, or `math.MaxInt` if index is out of bounds.

```go
// Given interface - don't modify
type UnknownSizeArray interface {
    Get(index int) int  // Returns math.MaxInt if out of bounds
}

func SearchUnknownSize(arr UnknownSizeArray, target int) int {
    // Find target in sorted array of unknown size
    // Return index or -1 if not found
}
```

**Approach:**
1. First find the bounds of the array
2. Then perform binary search within those bounds

---

### **MEDIUM LEVEL (Core Algorithm Mastery)**

#### Problem 4: Find First and Last Position
**Difficulty:** Medium | **Time:** 30-35 minutes

**Problem Statement:**
Given a sorted array with duplicates, find the first and last position of a target value.

```go
func SearchRange(nums []int, target int) []int {
    // Return [first_index, last_index] or [-1, -1] if not found
}

// Helper functions
func findFirst(nums []int, target int) int {
    // Find leftmost occurrence
}

func findLast(nums []int, target int) int {
    // Find rightmost occurrence  
}
```

**Test Cases:**
```go
nums1 := []int{5, 7, 7, 8, 8, 10}
target1 := 8  // Expected: [3, 4]

nums2 := []int{5, 7, 7, 8, 8, 10}
target2 := 6  // Expected: [-1, -1]

nums3 := []int{1}
target3 := 1  // Expected: [0, 0]

nums4 := []int{2, 2, 2, 2, 2}
target4 := 2  // Expected: [0, 4]
```

---

#### Problem 5: Search Insert Position
**Difficulty:** Medium | **Time:** 20-25 minutes

**Problem Statement:**
Given a sorted array and a target value, return the index where target should be inserted to maintain sorted order.

```go
func SearchInsertPosition(nums []int, target int) int {
    // Return index where target should be inserted
    // If target exists, return its index
}
```

**Test Cases:**
```go
nums1 := []int{1, 3, 5, 6}
target1 := 5  // Expected: 2 (found at index 2)

nums2 := []int{1, 3, 5, 6}  
target2 := 2  // Expected: 1 (insert between 1 and 3)

nums3 := []int{1, 3, 5, 6}
target3 := 7  // Expected: 4 (insert at end)

nums4 := []int{1, 3, 5, 6}
target4 := 0  // Expected: 0 (insert at beginning)
```

---

#### Problem 6: Search in Rotated Sorted Array
**Difficulty:** Medium | **Time:** 35-40 minutes

**Problem Statement:**
Search for a target value in a rotated sorted array. The array was sorted then rotated at some pivot.

```go
func SearchInRotatedArray(nums []int, target int) int {
    // Return index of target, or -1 if not found
    // Original: [0,1,2,4,5,6,7] might be rotated to [4,5,6,7,0,1,2]
}
```

**Test Cases:**
```go
nums1 := []int{4, 5, 6, 7, 0, 1, 2}
target1 := 0  // Expected: 4

nums2 := []int{4, 5, 6, 7, 0, 1, 2}
target2 := 3  // Expected: -1

nums3 := []int{1}
target3 := 0  // Expected: -1

nums4 := []int{1, 3}
target4 := 3  // Expected: 1
```

**Key Insight:** At any point, at least one half of the array is properly sorted. Use this to determine which half to search.

---

### **ADVANCED LEVEL (Interview Challenges)**

#### Problem 7: Find Peak Element
**Difficulty:** Medium-Hard | **Time:** 30-40 minutes

**Problem Statement:**
Find a peak element in an array. A peak element is greater than its neighbors. You can assume `nums[-1] = nums[n] = -∞`.

```go
func FindPeakElement(nums []int) int {
    // Return index of any peak element
    // Multiple peaks may exist, return any one
}
```

**Test Cases:**
```go
nums1 := []int{1, 2, 3, 1}        // Expected: 2 (index of element 3)
nums2 := []int{1, 2, 1, 3, 5, 6, 4}  // Expected: 1 or 5 (multiple peaks exist)
nums3 := []int{1}                 // Expected: 0
nums4 := []int{1, 2}              // Expected: 1
```

**Challenge:** Solve in O(log n) time using binary search approach.

---

#### Problem 8: Find Minimum in Rotated Sorted Array
**Difficulty:** Medium-Hard | **Time:** 25-35 minutes

**Problem Statement:**
Find the minimum element in a rotated sorted array (assume all elements are unique).

```go
func FindMinInRotatedArray(nums []int) int {
    // Return the minimum element value
}
```

**Test Cases:**
```go
nums1 := []int{3, 4, 5, 1, 2}     // Expected: 1
nums2 := []int{4, 5, 6, 7, 0, 1, 2}  // Expected: 0  
nums3 := []int{11, 13, 15, 17}    // Expected: 11 (no rotation)
nums4 := []int{2, 1}              // Expected: 1
```

**Follow-up:** What if duplicates are allowed? How does this change the solution?

---

#### Problem 9: Search in 2D Matrix
**Difficulty:** Medium | **Time:** 30-35 minutes

**Problem Statement:**
Search for a value in an m×n matrix where:
- Each row is sorted left to right
- The first integer of each row is greater than the last integer of the previous row

```go
func SearchMatrix(matrix [][]int, target int) bool {
    // Return true if target is found, false otherwise
}
```

**Test Cases:**
```go
matrix1 := [][]int{
    {1,  4,  7,  11},
    {2,  5,  8,  12},
    {3,  6,  9,  16},
    {10, 13, 14, 17},
}
target1 := 5  // Expected: true

matrix2 := [][]int{
    {1,  4,  7,  11},
    {2,  5,  8,  12},
    {3,  6,  9,  16},
    {10, 13, 14, 17},
}
target2 := 20  // Expected: false
```

**Approach Options:**
1. **Two-step binary search:** Find row, then search within row
2. **Treat as 1D array:** Use coordinate transformation
3. **Start from corner:** Eliminate row or column at each step

---

#### Problem 10: Kth Smallest Element in Sorted Matrix
**Difficulty:** Hard | **Time:** 45-60 minutes

**Problem Statement:**
Given an n×n matrix where each row and column is sorted in ascending order, find the kth smallest element.

```go
func KthSmallestInMatrix(matrix [][]int, k int) int {
    // Return the kth smallest element (1-indexed)
}
```

**Test Cases:**
```go
matrix1 := [][]int{
    {1,  5,  9},
    {10, 11, 13},
    {12, 13, 15},
}
k1 := 8  // Expected: 13

matrix2 := [][]int{
    {-5},
}
k2 := 1  // Expected: -5
```

**Approaches to consider:**
1. **Min-heap approach:** O(k log n)
2. **Binary search on value range:** O(n log(max-min))

---

## 🧪 Testing Framework

### Performance Testing Suite

```go
package main

import (
    "fmt"
    "math/rand"
    "sort"
    "time"
)

// Benchmark different search algorithms
func BenchmarkSearchAlgorithms() {
    sizes := []int{100, 1000, 10000, 100000}
    
    for _, size := range sizes {
        fmt.Printf("\n📊 Testing with array size: %d\n", size)
        fmt.Printf("%s\n", strings.Repeat("-", 40))
        
        // Generate sorted array
        arr := make([]int, size)
        for i := 0; i < size; i++ {
            arr[i] = i * 2  // Even numbers for testing
        }
        
        target := arr[size/2]  // Middle element
        
        // Test Linear Search
        start := time.Now()
        LinearSearchAll(arr, target)
        linearTime := time.Since(start)
        
        // Test Binary Search  
        start = time.Now()
        BinarySearchWithSteps(arr, target)
        binaryTime := time.Since(start)
        
        fmt.Printf("Linear Search:  %v\n", linearTime)
        fmt.Printf("Binary Search:  %v\n", binaryTime)
        
        if linearTime > binaryTime {
            speedup := float64(linearTime) / float64(binaryTime)
            fmt.Printf("Binary Search is %.1fx faster\n", speedup)
        }
    }
}

// Test correctness of search algorithms
func TestSearchCorrectness() {
    testCases := []struct {
        name     string
        arr      []int
        target   int
        expected int
    }{
        {"Found at beginning", []int{1, 2, 3, 4, 5}, 1, 0},
        {"Found at end", []int{1, 2, 3, 4, 5}, 5, 4},
        {"Found in middle", []int{1, 2, 3, 4, 5}, 3, 2},
        {"Not found - too small", []int{1, 2, 3, 4, 5}, 0, -1},
        {"Not found - too large", []int{1, 2, 3, 4, 5}, 6, -1},
        {"Not found - in range", []int{1, 3, 5, 7, 9}, 4, -1},
        {"Single element - found", []int{42}, 42, 0},
        {"Single element - not found", []int{42}, 41, -1},
    }
    
    for _, tc := range testCases {
        result, _ := BinarySearchWithSteps(tc.arr, tc.target)
        if result == tc.expected {
            fmt.Printf("✅ %s: PASS\n", tc.name)
        } else {
            fmt.Printf("❌ %s: FAIL (expected %d, got %d)\n", 
                tc.name, tc.expected, result)
        }
    }
}
```

---

## 💡 Interview Strategy Guide

### Common Binary Search Patterns

#### Pattern 1: Standard Binary Search Template
```go
func binarySearchTemplate(arr []int, target int) int {
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
```

#### Pattern 2: Find Boundary (First/Last Occurrence)
```go
func findLeftBoundary(arr []int, target int) int {
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

#### Pattern 3: Search for Condition
```go
func searchCondition(arr []int, condition func(int) bool) int {
    left, right := 0, len(arr)
    
    for left < right {
        mid := left + (right-left)/2
        
        if condition(arr[mid]) {
            right = mid
        } else {
            left = mid + 1
        }
    }
    
    return left
}
```

### Key Interview Points

1. **Always ask about array properties:**
   - Is it sorted?
   - Are there duplicates?
   - What about edge cases (empty, single element)?

2. **Common mistakes to avoid:**
   - Integer overflow: Use `mid = left + (right-left)/2`
   - Infinite loops: Check boundary update logic
   - Off-by-one errors: Be careful with `<=` vs `<`

3. **Time complexity analysis:**
   - Linear search: O(n)
   - Binary search: O(log n)
   - Always mention space complexity too

4. **When to use each approach:**
   - Linear: Unsorted data, small arrays, simple implementation
   - Binary: Large sorted data, repeated searches, optimization needed

---

## 🎯 Mastery Checklist

### Basic Skills
- [ ] Implement linear search with multiple variations
- [ ] Implement binary search from memory without bugs
- [ ] Handle edge cases (empty array, single element, duplicates)
- [ ] Explain why binary search requires sorted data

### Intermediate Skills  
- [ ] Find first and last occurrence of target in sorted array
- [ ] Search in rotated sorted array
- [ ] Find insertion position for maintaining sorted order
- [ ] Search in arrays with unknown size

### Advanced Skills
- [ ] Find peak element using binary search
- [ ] Search in 2D sorted matrices
- [ ] Find minimum in rotated sorted array
- [ ] Solve problems that require binary search on answer space

### Problem-Solving Skills
- [ ] Recognize when binary search can be applied to non-obvious problems
- [ ] Choose appropriate search algorithm based on data characteristics
- [ ] Optimize search algorithms for specific constraints
- [ ] Debug binary search boundary issues quickly

**Target:** Complete 8+ problems with bug-free implementations and optimal time complexity.