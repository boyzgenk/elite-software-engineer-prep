# 📊 Big-O Analysis Fundamentals
## Week 4, Days 1-2 | Algorithm Complexity Analysis

### 🎯 Learning Objectives
By the end of this section, you will:
- [ ] Master Big-O, Big-Θ, and Big-Ω notation concepts
- [ ] Analyze time and space complexity of algorithms
- [ ] Recognize common complexity patterns instantly
- [ ] Compare algorithms based on efficiency
- [ ] Make informed decisions about algorithm trade-offs

---

## 📚 What is Big-O Notation?

### 🔸 Definition and Purpose

**Big-O Notation:**
- Mathematical notation describing limiting behavior of functions
- Represents upper bound of algorithm's growth rate
- Focuses on worst-case scenario performance
- Ignores constants and lower-order terms
- Enables comparison of algorithms independent of hardware/implementation

**Why Big-O Matters:**
- Predicts performance as input size grows
- Helps choose right algorithm for the job
- Essential for system design and optimization  
- Standard language for discussing algorithm efficiency
- Critical for technical interviews

**Real-world Analogy:**
Think of Big-O like speed limits on roads:
- O(1): Teleportation - instant regardless of distance
- O(log n): Highway - efficient even for long distances
- O(n): City streets - time proportional to distance
- O(n²): Walking through a maze - time grows quadratically
- O(2^n): Visiting every house in every neighborhood - exponential explosion

### 🔸 Big-O Families (from best to worst)

```go
package main

import (
    "fmt"
    "math"
    "time"
)

// O(1) - Constant Time
func constantTime(arr []int) int {
    // Always takes the same time regardless of input size
    if len(arr) > 0 {
        return arr[0] // First element access
    }
    return -1
}

// O(log n) - Logarithmic Time  
func logarithmicTime(arr []int, target int) int {
    // Binary search - eliminates half the data each step
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

// O(n) - Linear Time
func linearTime(arr []int, target int) int {
    // Must potentially check every element
    for i, value := range arr {
        if value == target {
            return i
        }
    }
    return -1
}

// O(n log n) - Linearithmic Time
func mergeSort(arr []int) []int {
    // Divide (log n levels) and Conquer (n work per level)
    if len(arr) <= 1 {
        return arr
    }
    
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])   // log n recursive levels
    right := mergeSort(arr[mid:])
    
    return merge(left, right)      // O(n) work to merge
}

func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result
}

// O(n²) - Quadratic Time
func quadraticTime(arr []int) [][]int {
    // Nested loops - for each element, check all other elements
    pairs := [][]int{}
    n := len(arr)
    
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            pairs = append(pairs, []int{arr[i], arr[j]})
        }
    }
    return pairs
}

// O(2^n) - Exponential Time
func exponentialTime(n int) int {
    // Fibonacci without memoization - branches exponentially
    if n <= 1 {
        return n
    }
    return exponentialTime(n-1) + exponentialTime(n-2)
}

// O(n!) - Factorial Time  
func factorialTime(arr []int) [][]int {
    // Generate all permutations
    if len(arr) <= 1 {
        return [][]int{arr}
    }
    
    result := [][]int{}
    for i, val := range arr {
        // Remove element at index i
        remaining := make([]int, 0, len(arr)-1)
        remaining = append(remaining, arr[:i]...)
        remaining = append(remaining, arr[i+1:]...)
        
        // Get permutations of remaining elements
        subPerms := factorialTime(remaining)
        
        // Add current element to front of each permutation
        for _, perm := range subPerms {
            newPerm := make([]int, 0, len(perm)+1)
            newPerm = append(newPerm, val)
            newPerm = append(newPerm, perm...)
            result = append(result, newPerm)
        }
    }
    return result
}
```

---

## 📈 Growth Rate Comparison

### Visual Understanding of Growth Rates

```go
// Demonstrate growth rates with actual timing
func demonstrateGrowthRates() {
    sizes := []int{10, 100, 1000, 10000}
    
    fmt.Println("Algorithm Performance Comparison")
    fmt.Println("================================")
    fmt.Printf("%-10s %-12s %-12s %-12s %-12s\n", 
               "Size", "O(1)", "O(log n)", "O(n)", "O(n²)")
    
    for _, n := range sizes {
        // O(1) - constant operations
        const1 := 1
        
        // O(log n) - binary search steps
        logN := int(math.Log2(float64(n)))
        
        // O(n) - linear operations
        linear := n
        
        // O(n²) - quadratic operations
        quadratic := n * n
        
        fmt.Printf("%-10d %-12d %-12d %-12d %-12d\n", 
                   n, const1, logN, linear, quadratic)
    }
    
    fmt.Println("\nNotice how O(n²) grows much faster than others!")
}

// Practical timing comparison
func timeComplexityDemo() {
    sizes := []int{1000, 2000, 4000, 8000}
    
    fmt.Println("\nActual Runtime Comparison (milliseconds)")
    fmt.Println("=======================================")
    
    for _, size := range sizes {
        arr := make([]int, size)
        for i := range arr {
            arr[i] = i
        }
        
        // Time O(n) operation
        start := time.Now()
        sum := 0
        for _, val := range arr {
            sum += val
        }
        linearTime := time.Since(start).Nanoseconds() / 1000000
        
        // Time O(n²) operation  
        start = time.Now()
        count := 0
        for i := 0; i < len(arr); i++ {
            for j := i + 1; j < len(arr); j++ {
                if arr[i] > arr[j] {
                    count++
                }
            }
        }
        quadraticTime := time.Since(start).Nanoseconds() / 1000000
        
        fmt.Printf("Size: %d, O(n): %dms, O(n²): %dms, Ratio: %.1fx\n", 
                   size, linearTime, quadraticTime, 
                   float64(quadraticTime)/float64(linearTime))
    }
}
```

---

## 🔍 Analyzing Common Algorithms

### Data Structure Operations

```go
// Array/Slice operations analysis
func analyzeArrayOperations() {
    fmt.Println("Array Operations Complexity")
    fmt.Println("===========================")
    
    arr := []int{1, 2, 3, 4, 5}
    
    // O(1) - Access by index
    value := arr[2]
    fmt.Printf("Access arr[2]: O(1) - %d\n", value)
    
    // O(1) - Append to end (amortized)
    arr = append(arr, 6)
    fmt.Printf("Append to end: O(1) amortized\n")
    
    // O(n) - Insert at beginning
    arr = append([]int{0}, arr...)
    fmt.Printf("Insert at beginning: O(n)\n")
    
    // O(n) - Search for value
    target := 3
    for i, val := range arr {
        if val == target {
            fmt.Printf("Linear search found %d at index %d: O(n)\n", target, i)
            break
        }
    }
    
    // O(n) - Delete from middle
    index := 2
    arr = append(arr[:index], arr[index+1:]...)
    fmt.Printf("Delete from middle: O(n)\n")
}

// Hash table operations analysis
func analyzeHashTableOperations() {
    fmt.Println("\nHash Table Operations Complexity")
    fmt.Println("================================")
    
    hashMap := make(map[string]int)
    
    // O(1) average - Insert
    hashMap["apple"] = 5
    fmt.Printf("Insert key-value: O(1) average\n")
    
    // O(1) average - Access
    value, exists := hashMap["apple"]
    if exists {
        fmt.Printf("Access by key: O(1) average - %d\n", value)
    }
    
    // O(1) average - Delete
    delete(hashMap, "apple")
    fmt.Printf("Delete by key: O(1) average\n")
}

// Tree operations analysis (BST)
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

func analyzeBSTOperations() {
    fmt.Println("\nBinary Search Tree Operations")
    fmt.Println("=============================")
    
    // Build a balanced BST for demonstration
    root := &TreeNode{Val: 4}
    root.Left = &TreeNode{Val: 2}
    root.Right = &TreeNode{Val: 6}
    root.Left.Left = &TreeNode{Val: 1}
    root.Left.Right = &TreeNode{Val: 3}
    root.Right.Left = &TreeNode{Val: 5}
    root.Right.Right = &TreeNode{Val: 7}
    
    fmt.Printf("Search in BST: O(log n) average, O(n) worst\n")
    fmt.Printf("Insert in BST: O(log n) average, O(n) worst\n")
    fmt.Printf("Delete in BST: O(log n) average, O(n) worst\n")
    fmt.Printf("Tree traversal: O(n) - must visit all nodes\n")
}
```

### Sorting Algorithm Analysis

```go
func analyzeSortingComplexity() {
    fmt.Println("\nSorting Algorithms Complexity")
    fmt.Println("=============================")
    
    algorithms := []struct {
        name      string
        best      string
        average   string
        worst     string
        space     string
        stable    bool
    }{
        {"Bubble Sort", "O(n)", "O(n²)", "O(n²)", "O(1)", true},
        {"Selection Sort", "O(n²)", "O(n²)", "O(n²)", "O(1)", false},
        {"Insertion Sort", "O(n)", "O(n²)", "O(n²)", "O(1)", true},
        {"Merge Sort", "O(n log n)", "O(n log n)", "O(n log n)", "O(n)", true},
        {"Quick Sort", "O(n log n)", "O(n log n)", "O(n²)", "O(log n)", false},
        {"Heap Sort", "O(n log n)", "O(n log n)", "O(n log n)", "O(1)", false},
    }
    
    fmt.Printf("%-15s %-10s %-10s %-10s %-8s %-6s\n",
               "Algorithm", "Best", "Average", "Worst", "Space", "Stable")
    fmt.Println(strings.Repeat("-", 70))
    
    for _, algo := range algorithms {
        fmt.Printf("%-15s %-10s %-10s %-10s %-8s %-6t\n",
                   algo.name, algo.best, algo.average, algo.worst, 
                   algo.space, algo.stable)
    }
}
```

---

## 🎯 Complexity Analysis Techniques

### Analyzing Loops

```go
// Single loop - O(n)
func singleLoop(n int) {
    for i := 0; i < n; i++ {
        // O(1) operation
        fmt.Printf("%d ", i)
    }
    // Total: O(n)
}

// Nested loops - O(n²)
func nestedLoops(n int) {
    for i := 0; i < n; i++ {           // Outer loop: n iterations
        for j := 0; j < n; j++ {       // Inner loop: n iterations each
            // O(1) operation
            fmt.Printf("(%d,%d) ", i, j)
        }
    }
    // Total: O(n × n) = O(n²)
}

// Nested loops with dependency - O(n²) 
func triangularLoops(n int) {
    for i := 0; i < n; i++ {           // Outer: n iterations
        for j := 0; j < i; j++ {       // Inner: 0,1,2,...,n-1 iterations
            // O(1) operation
            fmt.Printf("(%d,%d) ", i, j)
        }
    }
    // Total: 0+1+2+...+(n-1) = n(n-1)/2 = O(n²)
}

// Loop with halving - O(log n)
func halvingLoop(n int) {
    i := n
    for i > 1 {
        fmt.Printf("%d ", i)
        i = i / 2                      // Halving each iteration
    }
    // Number of iterations: log₂(n), so O(log n)
}

// Loop with multiplication - O(log n)
func multiplicationLoop(n int) {
    i := 1
    for i < n {
        fmt.Printf("%d ", i)
        i = i * 2                      // Doubling each iteration
    }
    // Number of iterations: log₂(n), so O(log n)
}
```

### Analyzing Recursive Algorithms

```go
// Linear recursion - O(n)
func factorialRecursive(n int) int {
    if n <= 1 {
        return 1                       // Base case: O(1)
    }
    return n * factorialRecursive(n-1) // Recursive case: T(n-1) + O(1)
}
// Recurrence relation: T(n) = T(n-1) + O(1) = O(n)
// Space complexity: O(n) due to call stack

// Binary recursion - O(2^n)
func fibonacciNaive(n int) int {
    if n <= 1 {
        return n                       // Base case: O(1)
    }
    return fibonacciNaive(n-1) + fibonacciNaive(n-2) // Two recursive calls
}
// Recurrence relation: T(n) = T(n-1) + T(n-2) + O(1) = O(2^n)
// Space complexity: O(n) due to maximum call stack depth

// Divide and conquer - O(n log n)
func mergeSortAnalysis(arr []int) []int {
    if len(arr) <= 1 {
        return arr                     // Base case: O(1)
    }
    
    mid := len(arr) / 2
    left := mergeSortAnalysis(arr[:mid])   // T(n/2)
    right := mergeSortAnalysis(arr[mid:])  // T(n/2)
    
    return mergeAnalysis(left, right)      // O(n)
}
// Recurrence relation: T(n) = 2T(n/2) + O(n) = O(n log n)
// Space complexity: O(n) for temporary arrays + O(log n) call stack

func mergeAnalysis(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    
    // This loop runs at most len(left) + len(right) times = O(n)
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result
}
```

### Master Theorem for Divide and Conquer

```go
// Master Theorem: T(n) = aT(n/b) + f(n)
// where a ≥ 1, b > 1, and f(n) is asymptotically positive

func masterTheoremExamples() {
    fmt.Println("Master Theorem Examples")
    fmt.Println("======================")
    
    examples := []struct {
        algorithm   string
        recurrence  string
        a, b        int
        f           string
        complexity  string
    }{
        {"Binary Search", "T(n) = T(n/2) + O(1)", 1, 2, "O(1)", "O(log n)"},
        {"Merge Sort", "T(n) = 2T(n/2) + O(n)", 2, 2, "O(n)", "O(n log n)"},
        {"Quick Sort (avg)", "T(n) = 2T(n/2) + O(n)", 2, 2, "O(n)", "O(n log n)"},
        {"Tree Traversal", "T(n) = 2T(n/2) + O(1)", 2, 2, "O(1)", "O(n)"},
        {"Matrix Multiply", "T(n) = 8T(n/2) + O(n²)", 8, 2, "O(n²)", "O(n³)"},
    }
    
    fmt.Printf("%-15s %-20s %-8s %-8s %-8s %-12s\n",
               "Algorithm", "Recurrence", "a", "b", "f(n)", "Complexity")
    fmt.Println(strings.Repeat("-", 75))
    
    for _, ex := range examples {
        fmt.Printf("%-15s %-20s %-8d %-8d %-8s %-12s\n",
                   ex.algorithm, ex.recurrence, ex.a, ex.b, ex.f, ex.complexity)
    }
}
```

---

## 📊 Space Complexity Analysis

### Understanding Space Complexity

```go
func spaceComplexityExamples() {
    fmt.Println("\nSpace Complexity Examples")
    fmt.Println("=========================")
    
    n := 1000
    
    // O(1) - Constant space
    fmt.Println("O(1) Space - Few variables regardless of input size:")
    sum := 0
    max := 0
    for i := 0; i < n; i++ {
        sum += i
        if i > max {
            max = i
        }
    }
    
    // O(n) - Linear space
    fmt.Println("O(n) Space - Array/slice proportional to input:")
    arr := make([]int, n)
    for i := range arr {
        arr[i] = i * i
    }
    
    // O(n²) - Quadratic space
    fmt.Println("O(n²) Space - 2D matrix:")
    matrix := make([][]int, n)
    for i := range matrix {
        matrix[i] = make([]int, n)
        for j := range matrix[i] {
            matrix[i][j] = i * j
        }
    }
    
    // O(log n) - Logarithmic space (recursive call stack)
    fmt.Printf("O(log n) Space - Binary search call stack depth: %d\n", 
               int(math.Log2(float64(n))))
}

// Auxiliary space vs Total space
func spaceAnalysisExample(arr []int) []int {
    // Input space: O(n) - the input array
    // Auxiliary space: O(n) - additional space used by algorithm
    // Total space: O(n) - input + auxiliary
    
    result := make([]int, len(arr)) // O(n) auxiliary space
    
    for i, val := range arr {       // O(1) auxiliary space for loop variables
        result[i] = val * 2
    }
    
    return result                   // Total space: O(n)
}

// In-place vs Out-of-place algorithms
func inPlaceVsOutOfPlace() {
    fmt.Println("\nIn-place vs Out-of-place Algorithms")
    fmt.Println("===================================")
    
    arr := []int{5, 2, 8, 1, 9}
    
    // In-place: O(1) auxiliary space
    fmt.Println("In-place bubble sort: O(1) auxiliary space")
    bubbleSortInPlace(arr)
    
    arr2 := []int{5, 2, 8, 1, 9}
    
    // Out-of-place: O(n) auxiliary space
    fmt.Println("Out-of-place merge sort: O(n) auxiliary space")
    sorted := mergeSort(arr2)
    fmt.Printf("Result: %v\n", sorted)
}

func bubbleSortInPlace(arr []int) {
    n := len(arr)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j] // O(1) space for swap
            }
        }
    }
}
```

---

## 🎯 Practical Analysis Examples

### Real Interview Problems

```go
// Problem: Find two numbers that sum to target
func twoSumAnalysis(nums []int, target int) []int {
    // Approach 1: Brute Force - O(n²) time, O(1) space
    for i := 0; i < len(nums); i++ {
        for j := i + 1; j < len(nums); j++ {
            if nums[i]+nums[j] == target {
                return []int{i, j}
            }
        }
    }
    
    // Approach 2: Hash Table - O(n) time, O(n) space
    numMap := make(map[int]int)
    for i, num := range nums {
        complement := target - num
        if index, exists := numMap[complement]; exists {
            return []int{index, i}
        }
        numMap[num] = i
    }
    
    return nil
}

// Problem: Reverse array in-place
func reverseArrayAnalysis(arr []int) {
    // Two pointers approach: O(n) time, O(1) space
    left, right := 0, len(arr)-1
    
    for left < right {
        arr[left], arr[right] = arr[right], arr[left] // O(1) swap
        left++
        right--
    }
    // Time: O(n/2) = O(n)
    // Space: O(1) - only using two pointer variables
}

// Problem: Check if array is palindrome
func isPalindromeAnalysis(arr []int) bool {
    // Approach 1: Two pointers - O(n) time, O(1) space
    left, right := 0, len(arr)-1
    
    for left < right {
        if arr[left] != arr[right] {
            return false
        }
        left++
        right--
    }
    return true
    
    // Approach 2: Reverse and compare - O(n) time, O(n) space
    // reversed := make([]int, len(arr))
    // for i, val := range arr {
    //     reversed[len(arr)-1-i] = val
    // }
    // return reflect.DeepEqual(arr, reversed)
}

// Problem: Find maximum subarray sum
func maxSubarrayAnalysis(nums []int) int {
    // Kadane's Algorithm: O(n) time, O(1) space
    maxSoFar := nums[0]
    maxEndingHere := nums[0]
    
    for i := 1; i < len(nums); i++ {
        maxEndingHere = max(nums[i], maxEndingHere+nums[i])
        maxSoFar = max(maxSoFar, maxEndingHere)
    }
    
    return maxSoFar
    
    // Brute force would be O(n³):
    // for i := 0; i < len(nums); i++ {
    //     for j := i; j < len(nums); j++ {
    //         sum := 0
    //         for k := i; k <= j; k++ {
    //             sum += nums[k]
    //         }
    //         if sum > maxSum {
    //             maxSum = sum
    //         }
    //     }
    // }
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

---

## 🚀 Optimization Strategies

### Common Optimization Techniques

```go
// 1. Use better data structures
func optimizeWithDataStructures() {
    fmt.Println("Optimization with Better Data Structures")
    fmt.Println("=======================================")
    
    // Problem: Check if element exists
    arr := []int{1, 5, 3, 9, 2, 8, 4, 7, 6}
    target := 5
    
    // Inefficient: Linear search in array - O(n)
    found := false
    for _, val := range arr {
        if val == target {
            found = true
            break
        }
    }
    fmt.Printf("Linear search found %d: %t (O(n))\n", target, found)
    
    // Efficient: Use hash set - O(1)
    set := make(map[int]bool)
    for _, val := range arr {
        set[val] = true
    }
    found = set[target]
    fmt.Printf("Hash set lookup found %d: %t (O(1))\n", target, found)
}

// 2. Avoid nested loops when possible
func optimizeNestedLoops() {
    fmt.Println("\nOptimizing Nested Loops")
    fmt.Println("======================")
    
    arr1 := []int{1, 3, 5, 7}
    arr2 := []int{2, 4, 6, 8}
    
    // Inefficient: Nested loops to find common elements - O(n×m)
    common := []int{}
    for _, val1 := range arr1 {
        for _, val2 := range arr2 {
            if val1 == val2 {
                common = append(common, val1)
            }
        }
    }
    
    // Efficient: Use hash set - O(n+m)
    set := make(map[int]bool)
    for _, val := range arr1 {
        set[val] = true
    }
    
    commonOptimized := []int{}
    for _, val := range arr2 {
        if set[val] {
            commonOptimized = append(commonOptimized, val)
        }
    }
    
    fmt.Printf("Common elements: %v\n", commonOptimized)
}

// 3. Use memoization for recursive algorithms
var fibMemo = make(map[int]int)

func fibonacciOptimized(n int) int {
    if n <= 1 {
        return n
    }
    
    if val, exists := fibMemo[n]; exists {
        return val
    }
    
    fibMemo[n] = fibonacciOptimized(n-1) + fibonacciOptimized(n-2)
    return fibMemo[n]
}

// 4. Early termination and pruning
func optimizeWithEarlyTermination(arr []int, target int) bool {
    // If array is sorted, can terminate early
    for _, val := range arr {
        if val == target {
            return true
        }
        if val > target { // Early termination if sorted
            return false
        }
    }
    return false
}
```

---

## 💻 Practice Problems

### Complexity Analysis Exercises

```go
// Exercise 1: Analyze this function
func mysteryFunction1(n int) int {
    count := 0
    for i := 1; i < n; i *= 2 {
        for j := 0; j < i; j++ {
            count++
        }
    }
    return count
}
// Answer: O(n) - geometric series sum

// Exercise 2: Analyze this function  
func mysteryFunction2(arr []int) bool {
    n := len(arr)
    for i := 0; i < n-1; i++ {
        for j := i + 1; j < n; j++ {
            if arr[i] == arr[j] {
                return true
            }
        }
    }
    return false
}
// Answer: O(n²) time, O(1) space - checking for duplicates

// Exercise 3: Analyze this recursive function
func mysteryFunction3(n int) int {
    if n <= 1 {
        return 1
    }
    return mysteryFunction3(n/2) + mysteryFunction3(n/2)
}
// Answer: O(n) time, O(log n) space

// Exercise 4: Compare these two approaches
func approach1(matrix [][]int) int {
    sum := 0
    for i := 0; i < len(matrix); i++ {
        for j := 0; j < len(matrix[i]); j++ {
            sum += matrix[i][j]
        }
    }
    return sum
}

func approach2(matrix [][]int) int {
    sum := 0
    for _, row := range matrix {
        for _, val := range row {
            sum += val
        }
    }
    return sum
}
// Both are O(n×m) where n=rows, m=cols
```

---

## 🚀 Next Steps

### Final Integration Goals
- [ ] Analyze all previously learned algorithms
- [ ] Practice complexity analysis on interview problems
- [ ] Understand when to optimize vs when "good enough" is sufficient
- [ ] Master the art of explaining complexity to interviewers

### Preparation for Module 2 (Problem-Solving)
- Review all data structures and their complexities
- Practice explaining algorithmic choices
- Focus on clean, readable code that demonstrates understanding
- Prepare to discuss trade-offs between different approaches

**Big-O analysis is your strategic compass - it guides you toward efficient solutions and helps you avoid algorithmic dead ends. Master this and you'll think like a true computer scientist! Strategic optimization is the ultimate power move!** 📊

---

## 🔧 Common Big-O Mistakes to Avoid

### Mistake 1: Confusing Best, Average, and Worst Case
```go
// Quick Sort analysis
func quickSortAnalysis() {
    fmt.Println("Quick Sort Complexity:")
    fmt.Println("Best case (balanced partitions): O(n log n)")
    fmt.Println("Average case (random data): O(n log n)")  
    fmt.Println("Worst case (sorted/reverse sorted): O(n²)")
    fmt.Println("Space: O(log n) average, O(n) worst (call stack)")
}
```

### Mistake 2: Ignoring Hidden Complexities
```go
// Be careful with built-in functions
func hiddenComplexities(arr []int) {
    // This looks like O(n), but append might be O(n) if resize needed
    result := []int{}
    for _, val := range arr {
        result = append(result, val*2) // Potentially O(n) per operation
    }
    
    // Better: pre-allocate
    result2 := make([]int, len(arr))
    for i, val := range arr {
        result2[i] = val * 2 // Guaranteed O(1) per operation
    }
}
```

### Mistake 3: Over-optimizing Prematurely
```go
func prematureOptimization() {
    fmt.Println("Optimization Guidelines:")
    fmt.Println("1. Make it work first")
    fmt.Println("2. Make it readable") 
    fmt.Println("3. Make it fast (only if needed)")
    fmt.Println("4. Profile before optimizing")
    fmt.Println("5. Optimize the bottleneck, not everything")
}
```

**Remember: The best algorithm is the one that solves the problem correctly, efficiently, and can be understood by your team. Strategic thinking includes maintainability!** 💪