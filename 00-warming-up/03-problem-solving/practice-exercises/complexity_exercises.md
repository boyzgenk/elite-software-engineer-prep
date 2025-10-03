# Big-O Analysis Practice Exercises

## Exercise Set 1: Basic Time Complexity Analysis

### Exercise 1: Analyze the following code snippets

**A. Simple Loop**
```go
func simpleLoop(n int) int {
    sum := 0
    for i := 0; i < n; i++ {
        sum += i
    }
    return sum
}
```
**Question**: What is the time complexity?
**Answer**: O(n) - single loop that runs n times

---

**B. Nested Loops**
```go
func nestedLoops(n int) int {
    sum := 0
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            sum += i * j
        }
    }
    return sum
}
```
**Question**: What is the time complexity?
**Answer**: O(n²) - nested loops each running n times

---

**C. Sequential Operations**
```go
func sequential(arr []int) int {
    // First loop
    sum1 := 0
    for i := 0; i < len(arr); i++ {
        sum1 += arr[i]
    }
    
    // Second loop
    sum2 := 0
    for i := 0; i < len(arr); i++ {
        sum2 += arr[i] * 2
    }
    
    return sum1 + sum2
}
```
**Question**: What is the time complexity?
**Answer**: O(n) - two sequential O(n) operations = O(n) + O(n) = O(n)

---

## Exercise Set 2: More Complex Analysis

### Exercise 2: Binary Search Variations

**A. Standard Binary Search**
```go
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
```
**Question**: What is the time complexity?
**Answer**: O(log n) - halves search space each iteration

---

**B. Linear Search in Each Row of 2D Array**
```go
func searchMatrix(matrix [][]int, target int) bool {
    for i := 0; i < len(matrix); i++ {
        for j := 0; j < len(matrix[i]); j++ {
            if matrix[i][j] == target {
                return true
            }
        }
    }
    return false
}
```
**Question**: If matrix is m×n, what is the time complexity?
**Answer**: O(m×n) - visiting each element once

---

### Exercise 3: Recursive Complexity

**A. Factorial**
```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}
```
**Question**: What is the time and space complexity?
**Answer**: Time: O(n), Space: O(n) - n recursive calls on call stack

---

**B. Fibonacci (naive)**
```go
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}
```
**Question**: What is the time complexity?
**Answer**: O(2^n) - each call spawns two more calls

---

**C. Fibonacci (memoized)**
```go
func fibonacciMemo(n int, memo map[int]int) int {
    if n <= 1 {
        return n
    }
    
    if val, exists := memo[n]; exists {
        return val
    }
    
    memo[n] = fibonacciMemo(n-1, memo) + fibonacciMemo(n-2, memo)
    return memo[n]
}
```
**Question**: What is the time complexity with memoization?
**Answer**: O(n) - each value calculated once, stored in memo

---

## Exercise Set 3: Advanced Analysis

### Exercise 4: Tree Algorithms

**A. Tree Traversal**
```go
func inorderTraversal(root *TreeNode) {
    if root == nil {
        return
    }
    
    inorderTraversal(root.Left)
    fmt.Print(root.Val)
    inorderTraversal(root.Right)
}
```
**Question**: For a tree with n nodes, what is the complexity?
**Answer**: Time: O(n), Space: O(h) where h is height

---

**B. Find Height**
```go
func maxDepth(root *TreeNode) int {
    if root == nil {
        return 0
    }
    
    leftDepth := maxDepth(root.Left)
    rightDepth := maxDepth(root.Right)
    
    return 1 + max(leftDepth, rightDepth)
}
```
**Question**: What is the time and space complexity?
**Answer**: Time: O(n), Space: O(h) - visit all nodes, recursion depth = height

---

### Exercise 5: Sorting Algorithm Analysis

**A. Merge Sort**
```go
func mergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])
    right := mergeSort(arr[mid:])
    
    return merge(left, right) // O(n) operation
}
```
**Question**: What is the time complexity?
**Answer**: O(n log n) - log n levels, each level does O(n) work

---

**B. Quick Sort (average case)**
```go
func quickSort(arr []int, low, high int) {
    if low < high {
        pivotIndex := partition(arr, low, high) // O(n)
        quickSort(arr, low, pivotIndex-1)
        quickSort(arr, pivotIndex+1, high)
    }
}
```
**Question**: What is the average case time complexity?
**Answer**: O(n log n) - good pivot splits array roughly in half

**Question**: What is the worst case time complexity?
**Answer**: O(n²) - poor pivot choice leads to unbalanced splits

---

## Exercise Set 4: Space Complexity Analysis

### Exercise 6: Memory Usage

**A. Creating New Array**
```go
func doubleArray(arr []int) []int {
    doubled := make([]int, len(arr))
    for i, val := range arr {
        doubled[i] = val * 2
    }
    return doubled
}
```
**Question**: What is the space complexity?
**Answer**: O(n) - creates new array of size n

---

**B. In-place Modification**
```go
func doubleInPlace(arr []int) {
    for i := range arr {
        arr[i] *= 2
    }
}
```
**Question**: What is the space complexity?
**Answer**: O(1) - modifies existing array, no extra space

---

**C. Recursive Factorial Space**
```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}
```
**Question**: What is the space complexity?
**Answer**: O(n) - n function calls on the call stack

---

## Exercise Set 5: Best, Average, Worst Case Analysis

### Exercise 7: Search Algorithms

**A. Linear Search**
```go
func linearSearch(arr []int, target int) int {
    for i, val := range arr {
        if val == target {
            return i
        }
    }
    return -1
}
```
**Questions**:
- Best case: O(1) - target is first element
- Average case: O(n) - target is in middle on average
- Worst case: O(n) - target is last element or not present

---

**B. Binary Search**
```go
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
```
**Questions**:
- Best case: O(1) - target is middle element
- Average case: O(log n) - typically need log n comparisons
- Worst case: O(log n) - target not found, exhaust all log n levels

---

## Exercise Set 6: Challenge Problems

### Exercise 8: Complex Algorithm Analysis

**A. Find All Pairs with Sum**
```go
func findPairsWithSum(arr []int, target int) [][]int {
    pairs := [][]int{}
    
    for i := 0; i < len(arr); i++ {
        for j := i + 1; j < len(arr); j++ {
            if arr[i] + arr[j] == target {
                pairs = append(pairs, []int{arr[i], arr[j]})
            }
        }
    }
    
    return pairs
}
```
**Question**: What is the time complexity?
**Answer**: O(n²) - nested loops check all pairs

---

**B. Optimized Version with Hash Map**
```go
func findPairsOptimized(arr []int, target int) [][]int {
    seen := make(map[int]bool)
    pairs := [][]int{}
    
    for _, num := range arr {
        complement := target - num
        if seen[complement] {
            pairs = append(pairs, []int{complement, num})
        }
        seen[num] = true
    }
    
    return pairs
}
```
**Question**: What is the time and space complexity?
**Answer**: Time: O(n), Space: O(n) - single pass with hash map

---

### Exercise 9: Dynamic Programming

**A. Fibonacci Bottom-Up**
```go
func fibonacciDP(n int) int {
    if n <= 1 {
        return n
    }
    
    dp := make([]int, n+1)
    dp[0], dp[1] = 0, 1
    
    for i := 2; i <= n; i++ {
        dp[i] = dp[i-1] + dp[i-2]
    }
    
    return dp[n]
}
```
**Question**: What is the time and space complexity?
**Answer**: Time: O(n), Space: O(n) - array of size n+1

---

**B. Space-Optimized Fibonacci**
```go
func fibonacciOptimized(n int) int {
    if n <= 1 {
        return n
    }
    
    prev, curr := 0, 1
    
    for i := 2; i <= n; i++ {
        prev, curr = curr, prev+curr
    }
    
    return curr
}
```
**Question**: What is the space complexity improvement?
**Answer**: O(1) - only store previous two values instead of entire array

---

## Quick Reference: Common Complexities

### Time Complexities (from best to worst):
1. **O(1)** - Constant: hash table lookup, array access
2. **O(log n)** - Logarithmic: binary search, balanced tree operations
3. **O(n)** - Linear: single loop, linear search
4. **O(n log n)** - Linearithmic: efficient sorting (merge sort, heap sort)
5. **O(n²)** - Quadratic: nested loops, bubble sort
6. **O(2^n)** - Exponential: naive recursive algorithms
7. **O(n!)** - Factorial: generating all permutations

### Space Complexities:
- **O(1)** - Constant extra space
- **O(log n)** - Recursion depth in balanced trees
- **O(n)** - Linear extra space (new array, hash table)
- **O(n²)** - 2D arrays or matrices

### Tips for Analysis:
1. **Drop constants**: O(2n) → O(n)
2. **Drop lower terms**: O(n² + n) → O(n²)
3. **Consider worst case** unless specified otherwise
4. **Account for space used by recursion** (call stack)
5. **Distinguish between time and space** complexity

---

## Practice Strategy:
1. Start with simple loops and build complexity
2. Practice identifying patterns (nested loops = O(n²), halving = O(log n))
3. Always consider both time and space complexity
4. Think about best/average/worst cases
5. Practice explaining your analysis clearly