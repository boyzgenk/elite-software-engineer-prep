# Recursion Practice Problems - Day 7

## Learning Objectives
By the end of this session, you should be able to:
- Identify base cases and recursive cases in problems
- Apply recursion to solve array, string, and mathematical problems
- Implement backtracking algorithms for generating combinations
- Understand and apply divide-and-conquer strategies
- Debug recursive functions effectively

## Problem Templates

### Problem 1: Climbing Stairs (Easy)
**LeetCode #70**

**Problem Statement:**
You are climbing a staircase. It takes n steps to reach the top. Each time you can either climb 1 or 2 steps. In how many distinct ways can you climb to the top?

**Examples:**
```
Input: n = 2
Output: 2
Explanation: 1. 1 step + 1 step, 2. 2 steps

Input: n = 3
Output: 3
Explanation: 1. 1+1+1, 2. 1+2, 3. 2+1
```

**Golang Template:**
```go
func climbStairs(n int) int {
    // TODO: Implement using recursion with memoization
    // Think: How many ways to reach step n?
    // Base cases: n=1 (1 way), n=2 (2 ways)
    // Recursive: ways(n) = ways(n-1) + ways(n-2)
    
}
```

**Approach:**
1. Base cases: n=1 returns 1, n=2 returns 2
2. Recursive case: ways to reach n = ways to reach (n-1) + ways to reach (n-2)
3. Use memoization to avoid recalculating same values

---

### Problem 2: Generate Parentheses (Medium)
**LeetCode #22**

**Problem Statement:**
Given n pairs of parentheses, write a function to generate all combinations of well-formed parentheses.

**Examples:**
```
Input: n = 3
Output: ["((()))","(()())","(())()","()(())","()()()"]

Input: n = 1
Output: ["()"]
```

**Golang Template:**
```go
func generateParenthesis(n int) []string {
    // TODO: Implement using backtracking
    // Track: current string, open count, close count
    // Rules: open < n, close < open
    
}

func backtrack(current string, open, close, max int, result *[]string) {
    // TODO: Base case and recursive cases
    
}
```

**Approach:**
1. Use backtracking to build valid combinations
2. Track current string, open parentheses count, close parentheses count
3. Add '(' if open count < n
4. Add ')' if close count < open count

---

### Problem 3: Power of Two (Easy)
**LeetCode #231**

**Problem Statement:**
Given an integer n, return true if it is a power of two. Otherwise, return false.

**Examples:**
```
Input: n = 1
Output: true
Explanation: 2^0 = 1

Input: n = 16
Output: true
Explanation: 2^4 = 16

Input: n = 3
Output: false
```

**Golang Template:**
```go
func isPowerOfTwo(n int) bool {
    // TODO: Implement using recursion
    // Base cases: n=1 (true), n<1 or n odd (false)
    // Recursive: isPowerOfTwo(n/2)
    
}
```

**Approach:**
1. Base case: if n = 1, return true (2^0 = 1)
2. Base case: if n < 1 or n is odd, return false
3. Recursive case: check if n/2 is power of two

---

### Problem 4: Reverse String (Easy)
**LeetCode #344 - Modified for recursion**

**Problem Statement:**
Write a function that reverses a string using recursion. The input string is given as an array of characters.

**Examples:**
```
Input: s = ["h","e","l","l","o"]
Output: ["o","l","l","e","h"]

Input: s = ["H","a","n","n","a","h"]
Output: ["h","a","n","n","a","H"]
```

**Golang Template:**
```go
func reverseString(s []byte) {
    // TODO: Implement using recursion with two pointers
    // Helper function with left and right indices
    
}

func reverseHelper(s []byte, left, right int) {
    // TODO: Base case and swap logic
    
}
```

**Approach:**
1. Use helper function with left and right pointers
2. Base case: left >= right (nothing to swap)
3. Recursive case: swap characters at left and right, recurse with left+1, right-1

---

### Problem 5: Merge Two Sorted Lists (Easy)
**LeetCode #21**

**Problem Statement:**
You are given the heads of two sorted linked lists list1 and list2. Merge the two lists in a sorted manner and return the head of the merged linked list.

**Examples:**
```
Input: list1 = [1,2,4], list2 = [1,3,4]
Output: [1,1,2,3,4,4]

Input: list1 = [], list2 = []
Output: []
```

**Golang Template:**
```go
type ListNode struct {
    Val  int
    Next *ListNode
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
    // TODO: Implement using recursion
    // Base cases: one list is nil
    // Recursive: choose smaller head, merge rest
    
}
```

**Approach:**
1. Base cases: if either list is nil, return the other
2. Compare heads: choose smaller one
3. Recursively merge rest of chosen list with other list

---

## Advanced Problems

### Problem 6: Letter Combinations of Phone Number (Medium)
**LeetCode #17**

**Problem Statement:**
Given a string containing digits from 2-9 inclusive, return all possible letter combinations that the number could represent.

**Examples:**
```
Input: digits = "23"
Output: ["ad","ae","af","bd","be","bf","cd","ce","cf"]

Input: digits = ""
Output: []
```

**Golang Template:**
```go
func letterCombinations(digits string) []string {
    if len(digits) == 0 {
        return []string{}
    }
    
    // TODO: Define digit to letters mapping
    // TODO: Implement backtracking
    
}

func backtrackLetters(digits string, index int, current string, mapping map[byte]string, result *[]string) {
    // TODO: Base case and recursive exploration
    
}
```

**Approach:**
1. Create mapping from digits to letters
2. Use backtracking to build combinations
3. Base case: processed all digits, add current to result
4. For each letter of current digit, recurse with next digit

---

### Problem 7: Subsets (Medium)
**LeetCode #78**

**Problem Statement:**
Given an integer array nums of unique elements, return all possible subsets (the power set).

**Examples:**
```
Input: nums = [1,2,3]
Output: [[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]

Input: nums = [0]
Output: [[],[0]]
```

**Golang Template:**
```go
func subsets(nums []int) [][]int {
    result := [][]int{}
    // TODO: Implement backtracking
    // For each element: include it or don't include it
    
    return result
}

func backtrackSubsets(nums []int, index int, current []int, result *[][]int) {
    // TODO: Base case and recursive choices
    
}
```

**Approach:**
1. Use backtracking with inclusion/exclusion choices
2. Base case: processed all elements, add current subset
3. Two recursive calls: include current element, exclude current element

---

### Problem 8: Binary Tree Maximum Path Sum (Hard)
**LeetCode #124**

**Problem Statement:**
A path in a binary tree is a sequence of nodes where each pair of adjacent nodes has an edge connecting them. The path sum is the sum of the node's values in the path. Find the maximum path sum of any non-empty path.

**Examples:**
```
Input: root = [1,2,3]
Output: 6
Explanation: The optimal path is 2 -> 1 -> 3 with a path sum of 2 + 1 + 3 = 6.

Input: root = [-10,9,20,null,null,15,7]
Output: 42
Explanation: The optimal path is 15 -> 20 -> 7 with a path sum of 15 + 20 + 7 = 42.
```

**Golang Template:**
```go
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

func maxPathSum(root *TreeNode) int {
    maxSum := math.MinInt32
    // TODO: Implement helper function
    // Helper returns max path sum starting from node going down
    
    return maxSum
}

func maxPathHelper(node *TreeNode, maxSum *int) int {
    // TODO: Calculate max path through this node
    // Update global maximum
    // Return max path going down from this node
    
}
```

**Approach:**
1. For each node, consider path passing through it
2. Path through node = left_max + node.val + right_max
3. Update global maximum
4. Return max path going down (node + max(left, right))

---

## Daily Study Plan

### Morning Session (45 minutes)
1. **Review Theory (10 min)**: Recursion patterns, base cases, stack frames
2. **Solve Problems 1-3 (30 min)**: Focus on basic recursive thinking
3. **Debug Practice (5 min)**: Trace through recursive calls step by step

### Evening Session (45 minutes)
1. **Solve Problems 4-5 (25 min)**: Practice with different data structures
2. **Attempt Advanced Problems 6-8 (15 min)**: Challenge backtracking problems
3. **Reflection (5 min)**: Identify common patterns and mistakes

## Success Checklist
- [ ] Can identify base cases quickly
- [ ] Understand recursive case construction
- [ ] Comfortable with backtracking template
- [ ] Can trace recursive calls manually
- [ ] Know when to use memoization
- [ ] Can handle multiple recursive calls per function
- [ ] Understand divide-and-conquer approach

## Key Recursion Patterns

### 1. **Linear Recursion** (single recursive call)
```go
func linearRecursion(n int) int {
    // Base case
    if n <= baseCondition {
        return baseValue
    }
    
    // Recursive case
    return operation(n, linearRecursion(n-1))
}
```

### 2. **Binary Recursion** (two recursive calls)
```go
func binaryRecursion(n int) int {
    // Base case
    if n <= baseCondition {
        return baseValue
    }
    
    // Recursive cases
    return combine(
        binaryRecursion(n-1),
        binaryRecursion(n-2)
    )
}
```

### 3. **Backtracking Template**
```go
func backtrack(state, choices, result) {
    // Base case: solution found
    if isComplete(state) {
        result = append(result, copy(state))
        return
    }
    
    // Try each choice
    for choice in choices {
        // Choose
        makeChoice(state, choice)
        
        // Explore
        backtrack(newState, newChoices, result)
        
        // Unchoose (backtrack)
        undoChoice(state, choice)
    }
}
```

### 4. **Tree Recursion**
```go
func treeRecursion(node *TreeNode) returnType {
    // Base case
    if node == nil {
        return baseValue
    }
    
    // Process current node
    current := processNode(node)
    
    // Recurse on children
    left := treeRecursion(node.Left)
    right := treeRecursion(node.Right)
    
    // Combine results
    return combine(current, left, right)
}
```

### 5. **Memoization Pattern**
```go
func memoizedRecursion(n int) int {
    memo := make(map[int]int)
    return helper(n, memo)
}

func helper(n int, memo map[int]int) int {
    // Check memo first
    if val, exists := memo[n]; exists {
        return val
    }
    
    // Base case
    if baseCondition {
        return baseValue
    }
    
    // Calculate and store
    result := recursiveCalculation(n)
    memo[n] = result
    return result
}
```

## Common Mistakes to Avoid
1. **Missing base case** - Always define when recursion stops
2. **Infinite recursion** - Ensure progress toward base case
3. **Wrong base case** - Test edge cases (n=0, n=1, empty input)
4. **Stack overflow** - Consider iterative solution for deep recursion
5. **Inefficient recursion** - Use memoization when subproblems repeat
6. **Incorrect combination** - Verify how recursive results are combined

## Debugging Recursion Tips
1. **Trace small examples** - Follow execution step by step
2. **Print statements** - Add debug prints to see call stack
3. **Trust the recursion** - Assume recursive calls work correctly
4. **Check base cases first** - Ensure they handle all edge cases
5. **Verify progress** - Each call should move closer to base case

## Time/Space Complexity Quick Reference
- **Linear recursion**: O(n) time, O(n) space (call stack)
- **Binary recursion**: O(2^n) time, O(n) space (without memoization)
- **With memoization**: Often O(n) time, O(n) space
- **Tree recursion**: O(nodes) time, O(height) space
- **Divide and conquer**: Often O(n log n) time, O(log n) space