# 🔄 Recursion Fundamentals
## Week 3, Days 4-5 | Mastering Recursive Problem Solving

### 🎯 Learning Objectives
By the end of this section, you will:
- [ ] Understand recursion concepts and call stack mechanics
- [ ] Master the recursive problem-solving framework
- [ ] Implement classic recursive algorithms confidently  
- [ ] Convert between recursive and iterative solutions
- [ ] Solve 8-10 recursive problems using systematic approach

---

## 📚 Recursion Fundamentals

### 🔸 What is Recursion?

**Definition:**
- Function that calls itself to solve smaller instances of the same problem
- Must have a base case to prevent infinite recursion
- Each recursive call should move closer to the base case
- Natural fit for problems with self-similar substructure

**Key Components:**
1. **Base Case:** Condition where recursion stops
2. **Recursive Case:** Function calls itself with modified parameters
3. **Progress:** Each call should reduce problem size

**Real-world Analogies:**
- Russian nesting dolls (Matryoshka)
- Fractals in nature (tree branches, snowflakes)
- Mathematical definitions (factorial, Fibonacci)
- File system directory traversal

### 🔸 How Recursion Works - Call Stack

```go
package main

import "fmt"

// Simple example to demonstrate call stack
func countdown(n int) {
    fmt.Printf("Entering countdown(%d)\n", n)
    
    // Base case
    if n <= 0 {
        fmt.Println("Blastoff!")
        return
    }
    
    fmt.Printf("About to call countdown(%d)\n", n-1)
    countdown(n - 1) // Recursive call
    
    fmt.Printf("Returning from countdown(%d)\n", n)
}

/*
Call Stack Visualization for countdown(3):

countdown(3) calls countdown(2)
│
├── countdown(2) calls countdown(1) 
│   │
│   ├── countdown(1) calls countdown(0)
│   │   │
│   │   └── countdown(0) → "Blastoff!" (base case)
│   │   
│   └── countdown(1) returns
│   
└── countdown(2) returns
countdown(3) returns
*/
```

### 🔸 Recursive Problem-Solving Framework

**Step-by-Step Approach:**
1. **Identify the problem structure:** Can it be broken into smaller similar problems?
2. **Define the base case:** When should recursion stop?
3. **Define recursive relation:** How does problem relate to smaller instances?
4. **Ensure progress:** Each call must move toward base case
5. **Combine results:** How to use subproblem solutions for main problem

---

## 🛠️ Classic Recursive Algorithms

### Mathematical Recursion

```go
import "fmt"

// Factorial: n! = n × (n-1)!
func factorial(n int) int {
    // Base case
    if n <= 1 {
        return 1
    }
    
    // Recursive case: n! = n × (n-1)!
    return n * factorial(n-1)
}

// Fibonacci: F(n) = F(n-1) + F(n-2)
func fibonacci(n int) int {
    // Base cases
    if n <= 1 {
        return n
    }
    
    // Recursive case
    return fibonacci(n-1) + fibonacci(n-2)
}

// Optimized Fibonacci with memoization
func fibonacciMemo(n int) int {
    memo := make(map[int]int)
    return fibHelper(n, memo)
}

func fibHelper(n int, memo map[int]int) int {
    if n <= 1 {
        return n
    }
    
    if val, exists := memo[n]; exists {
        return val
    }
    
    memo[n] = fibHelper(n-1, memo) + fibHelper(n-2, memo)
    return memo[n]
}

// Power function: a^n = a × a^(n-1)
func power(base, exponent int) int {
    // Base case
    if exponent == 0 {
        return 1
    }
    
    // Negative exponent handling (for completeness)
    if exponent < 0 {
        return 0 // Simplified - normally would return 1/power(base, -exponent)
    }
    
    // Recursive case
    return base * power(base, exponent-1)
}

// Optimized power using fast exponentiation: O(log n)
func powerFast(base, exponent int) int {
    if exponent == 0 {
        return 1
    }
    
    if exponent < 0 {
        return 0 // Simplified
    }
    
    // If exponent is even: a^n = (a^(n/2))^2
    if exponent%2 == 0 {
        half := powerFast(base, exponent/2)
        return half * half
    }
    
    // If exponent is odd: a^n = a × a^(n-1)
    return base * powerFast(base, exponent-1)
}

// Greatest Common Divisor using Euclidean algorithm
func gcd(a, b int) int {
    // Base case
    if b == 0 {
        return a
    }
    
    // Recursive case: gcd(a,b) = gcd(b, a%b)
    return gcd(b, a%b)
}
```

### Array and String Recursion

```go
// Sum of array elements
func arraySum(arr []int) int {
    // Base case: empty array
    if len(arr) == 0 {
        return 0
    }
    
    // Recursive case: first element + sum of rest
    return arr[0] + arraySum(arr[1:])
}

// Find maximum element in array
func arrayMax(arr []int) int {
    // Base case: single element
    if len(arr) == 1 {
        return arr[0]
    }
    
    // Recursive case: max of first and max of rest
    restMax := arrayMax(arr[1:])
    if arr[0] > restMax {
        return arr[0]
    }
    return restMax
}

// Binary search recursively
func binarySearchRec(arr []int, target, left, right int) int {
    // Base case: element not found
    if left > right {
        return -1
    }
    
    mid := left + (right-left)/2
    
    // Base case: element found
    if arr[mid] == target {
        return mid
    }
    
    // Recursive cases
    if arr[mid] > target {
        return binarySearchRec(arr, target, left, mid-1)
    } else {
        return binarySearchRec(arr, target, mid+1, right)
    }
}

// Check if string is palindrome
func isPalindromeRec(s string) bool {
    // Base cases
    if len(s) <= 1 {
        return true
    }
    
    // Check first and last characters
    if s[0] != s[len(s)-1] {
        return false
    }
    
    // Recursive case: check middle substring
    return isPalindromeRec(s[1 : len(s)-1])
}

// Reverse string recursively
func reverseString(s string) string {
    // Base case
    if len(s) <= 1 {
        return s
    }
    
    // Recursive case: last char + reverse of rest
    return string(s[len(s)-1]) + reverseString(s[:len(s)-1])
}

// Generate all permutations of string
func permutations(s string) []string {
    if len(s) <= 1 {
        return []string{s}
    }
    
    result := []string{}
    
    // For each character, make it first and permute the rest
    for i, char := range s {
        // Remove character at position i
        remaining := s[:i] + s[i+1:]
        
        // Get all permutations of remaining characters
        subPerms := permutations(remaining)
        
        // Add current character to front of each permutation
        for _, perm := range subPerms {
            result = append(result, string(char)+perm)
        }
    }
    
    return result
}
```

### Tree Recursion

```go
// TreeNode definition (from previous section)
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

// Tree traversals (recursive)
func inorderTraversal(root *TreeNode, result *[]int) {
    if root == nil {
        return
    }
    
    inorderTraversal(root.Left, result)   // Left
    *result = append(*result, root.Val)   // Root
    inorderTraversal(root.Right, result)  // Right
}

func preorderTraversal(root *TreeNode, result *[]int) {
    if root == nil {
        return
    }
    
    *result = append(*result, root.Val)   // Root
    preorderTraversal(root.Left, result)  // Left
    preorderTraversal(root.Right, result) // Right
}

func postorderTraversal(root *TreeNode, result *[]int) {
    if root == nil {
        return
    }
    
    postorderTraversal(root.Left, result)  // Left
    postorderTraversal(root.Right, result) // Right
    *result = append(*result, root.Val)    // Root
}

// Calculate tree height
func treeHeight(root *TreeNode) int {
    if root == nil {
        return -1 // Height of empty tree
    }
    
    leftHeight := treeHeight(root.Left)
    rightHeight := treeHeight(root.Right)
    
    return max(leftHeight, rightHeight) + 1
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// Count total nodes in tree
func countNodes(root *TreeNode) int {
    if root == nil {
        return 0
    }
    
    return 1 + countNodes(root.Left) + countNodes(root.Right)
}

// Check if tree is symmetric
func isSymmetric(root *TreeNode) bool {
    if root == nil {
        return true
    }
    return isMirror(root.Left, root.Right)
}

func isMirror(left, right *TreeNode) bool {
    if left == nil && right == nil {
        return true
    }
    if left == nil || right == nil {
        return false
    }
    
    return left.Val == right.Val &&
           isMirror(left.Left, right.Right) &&
           isMirror(left.Right, right.Left)
}
```

---

## 🎯 Advanced Recursive Patterns

### Backtracking

```go
// Generate all subsets of array (power set)
func generateSubsets(nums []int) [][]int {
    result := [][]int{}
    current := []int{}
    backtrackSubsets(nums, 0, current, &result)
    return result
}

func backtrackSubsets(nums []int, start int, current []int, result *[][]int) {
    // Add current subset to result
    subset := make([]int, len(current))
    copy(subset, current)
    *result = append(*result, subset)
    
    // Try including each remaining element
    for i := start; i < len(nums); i++ {
        current = append(current, nums[i])    // Choose
        backtrackSubsets(nums, i+1, current, result) // Explore
        current = current[:len(current)-1]    // Unchoose (backtrack)
    }
}

// Generate all combinations of size k
func combinations(n, k int) [][]int {
    result := [][]int{}
    current := []int{}
    backtrackCombinations(n, k, 1, current, &result)
    return result
}

func backtrackCombinations(n, k, start int, current []int, result *[][]int) {
    // Base case: we have k elements
    if len(current) == k {
        combination := make([]int, len(current))
        copy(combination, current)
        *result = append(*result, combination)
        return
    }
    
    // Try each number from start to n
    for i := start; i <= n; i++ {
        current = append(current, i)          // Choose
        backtrackCombinations(n, k, i+1, current, result) // Explore
        current = current[:len(current)-1]    // Unchoose
    }
}

// Solve N-Queens problem
func solveNQueens(n int) [][]string {
    solutions := [][]string{}
    board := make([][]bool, n)
    for i := range board {
        board[i] = make([]bool, n)
    }
    
    backtrackQueens(board, 0, &solutions)
    return solutions
}

func backtrackQueens(board [][]bool, row int, solutions *[][]string) {
    n := len(board)
    
    // Base case: all queens placed
    if row == n {
        *solutions = append(*solutions, boardToStrings(board))
        return
    }
    
    // Try placing queen in each column of current row
    for col := 0; col < n; col++ {
        if isSafeQueens(board, row, col) {
            board[row][col] = true              // Place queen
            backtrackQueens(board, row+1, solutions) // Recurse
            board[row][col] = false             // Remove queen
        }
    }
}

func isSafeQueens(board [][]bool, row, col int) bool {
    n := len(board)
    
    // Check column
    for i := 0; i < row; i++ {
        if board[i][col] {
            return false
        }
    }
    
    // Check diagonal (top-left to bottom-right)
    for i, j := row-1, col-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
        if board[i][j] {
            return false
        }
    }
    
    // Check diagonal (top-right to bottom-left)
    for i, j := row-1, col+1; i >= 0 && j < n; i, j = i-1, j+1 {
        if board[i][j] {
            return false
        }
    }
    
    return true
}

func boardToStrings(board [][]bool) []string {
    result := []string{}
    for _, row := range board {
        s := ""
        for _, hasQueen := range row {
            if hasQueen {
                s += "Q"
            } else {
                s += "."
            }
        }
        result = append(result, s)
    }
    return result
}
```

### Divide and Conquer

```go
// Merge sort implementation
func mergeSort(arr []int) []int {
    // Base case: arrays of size 1 are already sorted
    if len(arr) <= 1 {
        return arr
    }
    
    // Divide: split array in half
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])   // Conquer left half
    right := mergeSort(arr[mid:])  // Conquer right half
    
    // Combine: merge sorted halves
    return merge(left, right)
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
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    
    return result
}

// Quick sort implementation
func quickSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    pivot := arr[len(arr)/2]
    var less, equal, greater []int
    
    // Partition around pivot
    for _, value := range arr {
        if value < pivot {
            less = append(less, value)
        } else if value == pivot {
            equal = append(equal, value)
        } else {
            greater = append(greater, value)
        }
    }
    
    // Recursively sort and combine
    result := []int{}
    result = append(result, quickSort(less)...)
    result = append(result, equal...)
    result = append(result, quickSort(greater)...)
    
    return result
}

// Find maximum subarray sum using divide and conquer
func maxSubarraySum(arr []int) int {
    if len(arr) == 0 {
        return 0
    }
    return maxSubarrayHelper(arr, 0, len(arr)-1)
}

func maxSubarrayHelper(arr []int, left, right int) int {
    // Base case: single element
    if left == right {
        return arr[left]
    }
    
    mid := left + (right-left)/2
    
    // Find max sum in left and right halves
    leftSum := maxSubarrayHelper(arr, left, mid)
    rightSum := maxSubarrayHelper(arr, mid+1, right)
    
    // Find max sum crossing the middle
    crossSum := maxCrossingSum(arr, left, mid, right)
    
    // Return maximum of the three
    return max(max(leftSum, rightSum), crossSum)
}

func maxCrossingSum(arr []int, left, mid, right int) int {
    // Find max sum for left side (including mid)
    leftSum := arr[mid]
    sum := arr[mid]
    for i := mid - 1; i >= left; i-- {
        sum += arr[i]
        if sum > leftSum {
            leftSum = sum
        }
    }
    
    // Find max sum for right side (excluding mid)
    rightSum := arr[mid+1]
    sum = arr[mid+1]
    for i := mid + 2; i <= right; i++ {
        sum += arr[i]
        if sum > rightSum {
            rightSum = sum
        }
    }
    
    return leftSum + rightSum
}
```

---

## 🔄 Recursion vs Iteration

### Converting Recursive to Iterative

```go
// Factorial - Recursive vs Iterative
func factorialRecursive(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorialRecursive(n-1)
}

func factorialIterative(n int) int {
    result := 1
    for i := 2; i <= n; i++ {
        result *= i
    }
    return result
}

// Fibonacci - Recursive vs Iterative
func fibonacciRecursive(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacciRecursive(n-1) + fibonacciRecursive(n-2)
}

func fibonacciIterative(n int) int {
    if n <= 1 {
        return n
    }
    
    prev, curr := 0, 1
    for i := 2; i <= n; i++ {
        prev, curr = curr, prev+curr
    }
    return curr
}

// Tree traversal - Recursive vs Iterative
func inorderRecursive(root *TreeNode) []int {
    result := []int{}
    if root == nil {
        return result
    }
    
    result = append(result, inorderRecursive(root.Left)...)
    result = append(result, root.Val)
    result = append(result, inorderRecursive(root.Right)...)
    
    return result
}

func inorderIterative(root *TreeNode) []int {
    result := []int{}
    stack := []*TreeNode{}
    current := root
    
    for current != nil || len(stack) > 0 {
        // Go to leftmost node
        for current != nil {
            stack = append(stack, current)
            current = current.Left
        }
        
        // Process current node
        current = stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        result = append(result, current.Val)
        
        // Move to right subtree
        current = current.Right
    }
    
    return result
}
```

---

## 💻 Practice Problems

### Easy Level (Day 4)

#### Problem 1: Power of Two (LeetCode #231)
```go
func isPowerOfTwo(n int) bool {
    // Check if n is power of 2 using recursion
    return false
}
```

#### Problem 2: Reverse String (LeetCode #344)
```go
func reverseStringRec(s []byte) {
    // Reverse string in-place using recursion
}
```

#### Problem 3: Merge Two Sorted Lists (LeetCode #21)
```go
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
    // Merge using recursion
    return nil
}
```

### Medium Level (Day 5)

#### Problem 4: Generate Parentheses (LeetCode #22)
```go
func generateParenthesis(n int) []string {
    // Generate all valid combinations using backtracking
    return []string{}
}
```

#### Problem 5: Subsets (LeetCode #78)
```go
func subsets(nums []int) [][]int {
    // Generate all subsets using recursion
    return [][]int{}
}
```

#### Problem 6: Tree Path Sum (LeetCode #112)
```go
func hasPathSum(root *TreeNode, targetSum int) bool {
    // Check if path exists with given sum
    return false
}
```

---

## 📊 Recursion Analysis

### Time and Space Complexity

| Pattern | Time Complexity | Space Complexity | Example |
|---------|----------------|------------------|---------|
| Linear Recursion | O(n) | O(n) | Factorial, Array Sum |
| Binary Recursion | O(2^n) | O(n) | Fibonacci (naive) |
| Divide & Conquer | O(n log n) | O(log n) | Merge Sort |
| Tree Recursion | O(nodes) | O(height) | Tree Traversal |
| Backtracking | O(b^d) | O(d) | N-Queens, Subsets |

### When to Use Recursion vs Iteration

**Use Recursion when:**
- Problem has recursive structure (trees, fractals)
- Divide and conquer applies naturally
- Backtracking is needed
- Code clarity is more important than efficiency

**Use Iteration when:**
- Memory usage is critical (no stack overhead)
- Performance is crucial
- Simple repetitive operations
- Tail recursion can be optimized

---

## 🚀 Next Steps

### Day 4 Goals (Basic Recursion)
- [ ] Master factorial, Fibonacci, power functions
- [ ] Understand call stack mechanics
- [ ] Practice array and string recursion
- [ ] Convert simple recursive to iterative

### Day 5 Goals (Advanced Recursion)
- [ ] Implement backtracking algorithms
- [ ] Master tree recursive patterns
- [ ] Practice divide and conquer
- [ ] Solve complex recursive problems

### Preparation for Big-O Analysis (Days 6-7)
- Review algorithm complexities learned so far
- Practice analyzing recursive time complexity
- Understand space-time tradeoffs
- Think about optimization opportunities

**Recursion is the art of solving big problems by solving smaller versions of themselves. Master this mindset and complex algorithms become elegant solutions! Strategic decomposition at its finest!** 🔄

---

## 🔧 Recursion Debugging Tips

### Common Recursion Mistakes:
1. **Missing base case** - Leads to infinite recursion
2. **Wrong base case** - Incorrect termination condition  
3. **No progress toward base case** - Parameters don't change appropriately
4. **Stack overflow** - Too many recursive calls for large inputs

### Debugging Techniques:
```go
// Add debugging to trace recursive calls
func factorialDebug(n int, depth int) int {
    indent := strings.Repeat("  ", depth)
    fmt.Printf("%sEntering factorial(%d)\n", indent, n)
    
    if n <= 1 {
        fmt.Printf("%sBase case reached: factorial(1) = 1\n", indent)
        return 1
    }
    
    result := n * factorialDebug(n-1, depth+1)
    fmt.Printf("%sReturning factorial(%d) = %d\n", indent, n, result)
    return result
}

// Iterative version for comparison/verification
func factorialIterativeCheck(n int) int {
    result := 1
    for i := 2; i <= n; i++ {
        result *= i
    }
    return result
}
```

### Optimization Strategies:
```go
// Memoization for expensive recursive calls
func fibonacciMemoized(n int) int {
    memo := make([]int, n+1)
    for i := range memo {
        memo[i] = -1 // Mark as uncomputed
    }
    return fibMemo(n, memo)
}

func fibMemo(n int, memo []int) int {
    if n <= 1 {
        return n
    }
    
    if memo[n] != -1 {
        return memo[n] // Return cached result
    }
    
    memo[n] = fibMemo(n-1, memo) + fibMemo(n-2, memo)
    return memo[n]
}

// Tail recursion optimization (Go doesn't optimize, but good to know)
func factorialTailRec(n int) int {
    return factorialTailHelper(n, 1)
}

func factorialTailHelper(n, accumulator int) int {
    if n <= 1 {
        return accumulator
    }
    return factorialTailHelper(n-1, n*accumulator)
}
```

**Remember: Recursion is thinking in terms of smaller subproblems. Each recursive call should be a step toward the solution, not just a repetition. Think strategically about the problem structure!** 💪