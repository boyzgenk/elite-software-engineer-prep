# Sorting Algorithms Practice Problems
## Module 1: CS Fundamentals - Basic Algorithms

### 🎯 Learning Objectives
By completing these problems, you will:
- Master implementation of fundamental sorting algorithms
- Understand time and space complexity trade-offs
- Apply appropriate sorting algorithms to different scenarios
- Develop strong foundation for more advanced algorithms

---

## 📚 Problem Categories

### **EASY LEVEL (Foundation Building)**

#### Problem 1: Implement Bubble Sort with Optimization
**Difficulty:** Easy | **Time:** 15-20 minutes

**Problem Statement:**
Implement bubble sort algorithm that stops early if the array becomes sorted before all passes are complete.

```go
func BubbleSortOptimized(arr []int) []int {
    // Your implementation here
    // Should detect when no swaps occur and terminate early
}
```

**Test Cases:**
```go
// Test 1: Random array
input1 := []int{64, 34, 25, 12, 22, 11, 90}
expected1 := []int{11, 12, 22, 25, 34, 64, 90}

// Test 2: Already sorted (should terminate after 1 pass)
input2 := []int{1, 2, 3, 4, 5}
expected2 := []int{1, 2, 3, 4, 5}

// Test 3: Reverse sorted (worst case)
input3 := []int{5, 4, 3, 2, 1}
expected3 := []int{1, 2, 3, 4, 5}
```

**Expected Output:**
- Return sorted array
- Print number of passes needed
- Demonstrate early termination for sorted arrays

---

#### Problem 2: Selection Sort with Minimum Tracking
**Difficulty:** Easy | **Time:** 15-20 minutes

**Problem Statement:**
Implement selection sort and track how many times you find a new minimum element.

```go
func SelectionSortWithTracking(arr []int) ([]int, int) {
    // Return sorted array and count of minimum element updates
}
```

**Test Cases:**
```go
input := []int{29, 10, 14, 37, 13}
expectedArray := []int{10, 13, 14, 29, 37}
// Should return array and count of minimum updates
```

---

#### Problem 3: Insertion Sort for Nearly Sorted Array
**Difficulty:** Easy | **Time:** 20 minutes

**Problem Statement:**
Implement insertion sort and measure its performance on nearly sorted arrays vs random arrays.

```go
func InsertionSortWithMetrics(arr []int) ([]int, int) {
    // Return sorted array and number of shifts performed
}
```

**Test Cases:**
```go
// Nearly sorted (should be very efficient)
nearlySorted := []int{1, 2, 4, 3, 5, 6, 8, 7, 9}

// Random array (more shifts needed)
random := []int{9, 3, 7, 1, 5, 2, 8, 4, 6}
```

---

### **MEDIUM LEVEL (Core Algorithm Implementation)**

#### Problem 4: Merge Sort Implementation
**Difficulty:** Medium | **Time:** 30-40 minutes

**Problem Statement:**
Implement merge sort algorithm with detailed step tracking for educational purposes.

```go
func MergeSortWithSteps(arr []int) []int {
    // Implement merge sort with step-by-step logging
    // Show divide and conquer process clearly
}

func merge(left, right []int) []int {
    // Implement merge function
}
```

**Requirements:**
- Print division steps: "Dividing [4,2,7,1] into [4,2] and [7,1]"
- Print merge steps: "Merging [2,4] and [1,7] into [1,2,4,7]"
- Handle arrays of any size (including odd lengths)

**Test Cases:**
```go
input1 := []int{38, 27, 43, 3, 9, 82, 10}
input2 := []int{5, 4, 3, 2, 1}
input3 := []int{1} // Single element
input4 := []int{}  // Empty array
```

---

#### Problem 5: Quick Sort with Random Pivot
**Difficulty:** Medium | **Time:** 35-45 minutes

**Problem Statement:**
Implement quicksort with random pivot selection to avoid worst-case performance on sorted arrays.

```go
func QuickSortRandom(arr []int) []int {
    // Implement quicksort with random pivot selection
}

func partitionRandom(arr []int, low, high int) int {
    // Implement partition with random pivot
}
```

**Requirements:**
- Use random pivot to avoid O(n²) on sorted arrays
- Show partitioning process with detailed logging
- Handle duplicate elements correctly

**Test Cases:**
```go
sorted := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}     // Worst case for basic quicksort
reverse := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}     // Another worst case
duplicates := []int{5, 5, 5, 1, 5, 5, 3, 5, 5, 5}   // Many duplicates
```

---

#### Problem 6: Hybrid Sorting Algorithm
**Difficulty:** Medium | **Time:** 40-50 minutes

**Problem Statement:**
Create a hybrid sorting algorithm that uses insertion sort for small arrays (≤ 10 elements) and merge sort for larger arrays.

```go
func HybridSort(arr []int) []int {
    // Use insertion sort for small arrays, merge sort for large ones
    // Threshold: 10 elements
}
```

**Requirements:**
- Switch to insertion sort when subarray size ≤ 10
- Measure and compare performance with pure merge sort
- Handle edge cases properly

---

### **ADVANCED LEVEL (Interview Challenge Problems)**

#### Problem 7: Sort Colors (Dutch National Flag)
**Difficulty:** Medium-Hard | **Time:** 25-35 minutes

**Problem Statement:**
Given an array with 0s, 1s, and 2s, sort it in-place in a single pass.

```go
func SortColors(nums []int) {
    // Sort array of 0s, 1s, and 2s in-place
    // Must be done in single pass with O(1) space
}
```

**Test Cases:**
```go
input1 := []int{2, 0, 2, 1, 1, 0}     // Expected: [0, 0, 1, 1, 2, 2]
input2 := []int{2, 0, 1}              // Expected: [0, 1, 2]
input3 := []int{0}                    // Expected: [0]
input4 := []int{1, 1, 1, 1}           // Expected: [1, 1, 1, 1]
```

**Follow-up:** Can you solve it in O(n) time and O(1) space using only one pass?

---

#### Problem 8: Kth Largest Element Using QuickSelect
**Difficulty:** Hard | **Time:** 35-45 minutes

**Problem Statement:**
Find the kth largest element in an unsorted array using the QuickSelect algorithm (based on quicksort partitioning).

```go
func FindKthLargest(nums []int, k int) int {
    // Use QuickSelect algorithm
    // Average O(n) time complexity, O(1) space
}
```

**Test Cases:**
```go
nums1 := []int{3, 2, 1, 5, 6, 4}
k1 := 2  // Expected: 5 (2nd largest)

nums2 := []int{3, 2, 3, 1, 2, 4, 5, 5, 6}
k2 := 4  // Expected: 4 (4th largest)
```

---

#### Problem 9: Merge K Sorted Arrays
**Difficulty:** Hard | **Time:** 45-60 minutes

**Problem Statement:**
Given k sorted arrays, merge them into one sorted array efficiently.

```go
func MergeKSortedArrays(arrays [][]int) []int {
    // Merge k sorted arrays efficiently
    // Try both naive and optimized approaches
}
```

**Test Cases:**
```go
arrays1 := [][]int{
    {1, 4, 5},
    {1, 3, 4},
    {2, 6},
}
// Expected: [1, 1, 2, 3, 4, 4, 5, 6]

arrays2 := [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}
// Expected: [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

**Approaches to implement:**
1. **Naive:** Merge arrays one by one - O(kn log(kn))
2. **Optimized:** Use divide-and-conquer - O(nk log k)
3. **Heap-based:** Use min-heap - O(nk log k)

---

## 🧪 Testing Framework

### Running Your Solutions

```go
package main

import (
    "fmt"
    "reflect"
    "time"
)

// Test framework for sorting problems
func TestSortingFunction(name string, sortFunc func([]int) []int, testCases []TestCase) {
    fmt.Printf("\n🧪 Testing %s\n", name)
    fmt.Printf("%s\n", strings.Repeat("-", 40))
    
    allPassed := true
    
    for i, tc := range testCases {
        // Make a copy of input to avoid modifying original
        input := make([]int, len(tc.Input))
        copy(input, tc.Input)
        
        start := time.Now()
        result := sortFunc(input)
        duration := time.Since(start)
        
        passed := reflect.DeepEqual(result, tc.Expected)
        if passed {
            fmt.Printf("  ✅ Test %d: PASS (Time: %v)\n", i+1, duration)
        } else {
            fmt.Printf("  ❌ Test %d: FAIL\n", i+1)
            fmt.Printf("     Input:    %v\n", tc.Input)
            fmt.Printf("     Expected: %v\n", tc.Expected)
            fmt.Printf("     Got:      %v\n", result)
            allPassed = false
        }
    }
    
    if allPassed {
        fmt.Printf("🎉 All tests passed for %s!\n", name)
    } else {
        fmt.Printf("❌ Some tests failed for %s\n", name)
    }
}

type TestCase struct {
    Input    []int
    Expected []int
    Name     string
}

// Example usage:
func main() {
    testCases := []TestCase{
        {[]int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}, "Random array"},
        {[]int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}, "Reverse sorted"},
        {[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}, "Already sorted"},
        {[]int{1}, []int{1}, "Single element"},
        {[]int{}, []int{}, "Empty array"},
    }
    
    // Test your implementations
    TestSortingFunction("Bubble Sort", BubbleSortOptimized, testCases)
    TestSortingFunction("Merge Sort", MergeSortWithSteps, testCases)
    // Add more tests...
}
```

---

## 📊 Performance Analysis Exercises

### Exercise A: Compare Algorithm Performance

Create a program that:
1. Generates arrays of different sizes (100, 1000, 10000 elements)
2. Tests each sorting algorithm with the same data
3. Measures and compares execution times
4. Creates performance charts/reports

### Exercise B: Best/Average/Worst Case Analysis

For each algorithm, create test cases that demonstrate:
- **Best case performance** (e.g., already sorted for insertion sort)
- **Average case performance** (random data)
- **Worst case performance** (e.g., reverse sorted for bubble sort)

### Exercise C: Memory Usage Analysis

Implement versions that track:
- Number of comparisons made
- Number of swaps/moves performed  
- Maximum additional memory used
- Recursion depth (for recursive algorithms)

---

## 💡 Interview Tips

### Common Questions to Expect:
1. **"Implement merge sort"** - Focus on the merge function
2. **"When would you use insertion sort over quicksort?"** - Small arrays, nearly sorted data
3. **"What's the difference between stable and unstable sorting?"** - Preserve relative order of equal elements
4. **"How would you sort a million integers with limited memory?"** - External sorting, merge sort
5. **"Sort an array of 0s, 1s, and 2s"** - Dutch national flag problem

### Implementation Tips:
- Always test with edge cases: empty array, single element, all duplicates
- Watch out for integer overflow in pivot calculations
- Practice implementing without IDE autocomplete
- Explain time/space complexity while coding
- Discuss trade-offs between different algorithms

---

## 🎯 Mastery Checklist

- [ ] Can implement bubble sort with early termination
- [ ] Can implement insertion sort and understand its best-case O(n) behavior
- [ ] Can implement selection sort and explain why it's always O(n²)
- [ ] Can implement merge sort from scratch, including the merge function
- [ ] Can implement quicksort with proper partitioning
- [ ] Understand when to use each algorithm based on data characteristics
- [ ] Can analyze and compare time/space complexity of all algorithms
- [ ] Can solve advanced problems like Dutch National Flag
- [ ] Can implement hybrid approaches for optimal performance
- [ ] Can explain trade-offs and choose appropriate algorithm for given constraints

**Target:** Complete 7+ problems with working solutions and pass all test cases.