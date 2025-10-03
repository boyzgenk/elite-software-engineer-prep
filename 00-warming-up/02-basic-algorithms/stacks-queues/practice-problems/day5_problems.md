# 📝 Day 5 Practice Problems: Stack Fundamentals
## LIFO Principle and Stack Applications

### 🎯 Today's Goals
- Master stack operations and LIFO principle
- Solve parentheses matching problems
- Understand monotonic stack patterns
- Practice stack-based problem solving
- Build confidence with stack applications

---

## Problem 1: Valid Parentheses (Easy) ⭐⭐
**LeetCode #20**

Determine if the input string of brackets is valid. Valid means:
1. Open brackets must be closed by the same type of brackets
2. Open brackets must be closed in the correct order

**Example:**
```
Input: s = "()"
Output: true

Input: s = "()[]{}"
Output: true

Input: s = "(]"
Output: false

Input: s = "([)]"
Output: false
```

**Your Task:**
1. Use stack to track opening brackets
2. When encountering closing bracket, check if it matches top of stack
3. Stack should be empty at the end

**Template:**
```go
func isValid(s string) bool {
    stack := []rune{}
    
    // Create mapping of closing to opening brackets
    pairs := map[rune]rune{
        ')': '(',
        '}': '{',
        ']': '[',
    }
    
    // Your implementation here
    
    return len(stack) == 0
}

// Test cases
// isValid("()") should return true
// isValid("()[]{}" should return true
// isValid("(]") should return false
// isValid("([)]") should return false
```

**Learning Focus:**
- Stack for matching pairs
- Hash map for quick lookups
- String processing in Go

---

## Problem 2: Min Stack (Easy) ⭐⭐⭐
**LeetCode #155**

Design a stack that supports push, pop, top, and retrieving the minimum element in constant time.

**Example:**
```
Input:
["MinStack","push","push","push","getMin","pop","top","getMin"]
[[],[-2],[0],[-3],[],[],[],[]]

Output:
[null,null,null,null,-3,null,0,-2]
```

**Your Task:**
1. Maintain a separate stack for minimums
2. Ensure all operations are O(1)
3. Handle edge cases properly

**Template:**
```go
type MinStack struct {
    stack    []int
    minStack []int
}

func Constructor() MinStack {
    return MinStack{
        stack:    []int{},
        minStack: []int{},
    }
}

func (s *MinStack) Push(val int) {
    // Your implementation here
}

func (s *MinStack) Pop() {
    // Your implementation here
}

func (s *MinStack) Top() int {
    // Your implementation here
    return 0
}

func (s *MinStack) GetMin() int {
    // Your implementation here
    return 0
}

// Test with the example above
```

**Learning Focus:**
- Auxiliary data structures
- Maintaining invariants
- Constant time operations

---

## Problem 3: Evaluate Reverse Polish Notation (Medium) ⭐⭐⭐
**LeetCode #150**

Evaluate the value of an arithmetic expression in Reverse Polish Notation (RPN).

**Example:**
```
Input: tokens = ["2","1","+","3","*"]
Output: 9
Explanation: ((2 + 1) * 3) = 9

Input: tokens = ["4","13","5","/","+"]
Output: 6
Explanation: (4 + (13 / 5)) = 6
```

**Your Task:**
1. Use stack to store operands
2. When encountering operator, pop two operands and apply operation
3. Push result back to stack

**Template:**
```go
func evalRPN(tokens []string) int {
    stack := []int{}
    
    for _, token := range tokens {
        switch token {
        case "+":
            // Your implementation for addition
        case "-":
            // Your implementation for subtraction
        case "*":
            // Your implementation for multiplication
        case "/":
            // Your implementation for division
        default:
            // Convert string to int and push to stack
        }
    }
    
    return stack[0]
}

// Test cases
// evalRPN(["2","1","+","3","*"]) should return 9
// evalRPN(["4","13","5","/","+"]) should return 6
```

**Learning Focus:**
- Postfix expression evaluation
- String to integer conversion in Go
- Operator precedence handling

---

## Problem 4: Daily Temperatures (Medium) ⭐⭐⭐⭐
**LeetCode #739**

Find how many days you have to wait until a warmer temperature. Return array where answer[i] is the number of days after day i until a warmer temperature.

**Example:**
```
Input: temperatures = [73,74,75,71,69,72,76,73]
Output: [1,1,4,2,1,1,0,0]
```

**Your Task:**
1. Use monotonic decreasing stack
2. Stack stores indices, not values
3. When current temperature is higher than stack top, we found the answer

**Template:**
```go
func dailyTemperatures(temperatures []int) []int {
    result := make([]int, len(temperatures))
    stack := []int{} // Stack of indices
    
    for i, temp := range temperatures {
        // While stack not empty and current temp > temp at stack top
        for len(stack) > 0 && temp > temperatures[stack[len(stack)-1]] {
            // Your implementation here
            // Pop from stack and calculate days difference
        }
        
        // Push current index to stack
        stack = append(stack, i)
    }
    
    return result
}

// Test case
// dailyTemperatures([73,74,75,71,69,72,76,73]) should return [1,1,4,2,1,1,0,0]
```

**Learning Focus:**
- Monotonic stack pattern
- Stack of indices technique
- Next greater element problems

---

## Problem 5: Baseball Game (Easy) ⭐⭐
**LeetCode #682**

Calculate the sum of scores after all operations in a baseball game.

**Operations:**
- Integer: Record a new score
- "+": Record sum of last two scores
- "D": Record double of last score  
- "C": Cancel last score

**Example:**
```
Input: ops = ["5","2","C","D","+"]
Output: 30
Explanation: 
"5" -> [5]
"2" -> [5, 2]
"C" -> [5] (cancel 2)
"D" -> [5, 10] (double of 5)
"+" -> [5, 10, 15] (5 + 10)
Sum = 30
```

**Template:**
```go
func calPoints(operations []string) int {
    stack := []int{}
    
    for _, op := range operations {
        switch op {
        case "+":
            // Add sum of last two scores
        case "D":
            // Add double of last score
        case "C":
            // Cancel last score
        default:
            // Convert string to int and add to stack
        }
    }
    
    // Calculate sum of all scores
    sum := 0
    for _, score := range stack {
        sum += score
    }
    
    return sum
}

// Test case
// calPoints(["5","2","C","D","+"]) should return 30
```

**Learning Focus:**
- Stack for undo operations
- String operations in Go
- Multiple conditional handling

---

## Bonus Problem: Largest Rectangle in Histogram (Hard) ⭐⭐⭐⭐⭐
**LeetCode #84**

Find the area of the largest rectangle that can be formed in a histogram.

**Example:**
```
Input: heights = [2,1,5,6,2,3]
Output: 10
Explanation: Rectangle with height 5 and width 2 has area 10
```

**Your Task:**
1. Use monotonic increasing stack
2. When we find a bar shorter than stack top, calculate area
3. This is an advanced problem - focus on understanding the pattern

**Template:**
```go
func largestRectangleArea(heights []int) int {
    stack := []int{} // Stack of indices
    maxArea := 0
    
    for i, height := range heights {
        // While stack not empty and current height < height at stack top
        for len(stack) > 0 && height < heights[stack[len(stack)-1]] {
            // Calculate area with stack top as smallest bar
            // This is complex - focus on understanding the approach
        }
        stack = append(stack, i)
    }
    
    // Handle remaining bars in stack
    // ...
    
    return maxArea
}
```

**Learning Focus:**
- Advanced monotonic stack usage
- Area calculation with constraints
- Complex index manipulation

---

## 🧠 Stack Patterns Summary

### Pattern 1: Matching Pairs
```go
// Use stack to match opening/closing elements
func matchingPattern(s string) bool {
    stack := []rune{}
    for _, char := range s {
        if isOpening(char) {
            stack = append(stack, char)
        } else if len(stack) > 0 && matches(stack[len(stack)-1], char) {
            stack = stack[:len(stack)-1] // Pop
        } else {
            return false
        }
    }
    return len(stack) == 0
}
```

### Pattern 2: Monotonic Stack
```go
// Maintain stack in increasing/decreasing order
func monotonicStack(arr []int) []int {
    stack := []int{} // Usually stores indices
    result := make([]int, len(arr))
    
    for i, val := range arr {
        // Pop elements that violate monotonic property
        for len(stack) > 0 && shouldPop(val, arr[stack[len(stack)-1]]) {
            idx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            // Process the popped element
            result[idx] = calculateResult(i, idx)
        }
        stack = append(stack, i)
    }
    
    return result
}
```

### Pattern 3: Expression Evaluation
```go
// Use stack to evaluate expressions
func evaluateExpression(tokens []string) int {
    stack := []int{}
    
    for _, token := range tokens {
        if isOperator(token) {
            b := stack[len(stack)-1]; stack = stack[:len(stack)-1]
            a := stack[len(stack)-1]; stack = stack[:len(stack)-1]
            result := applyOperator(a, b, token)
            stack = append(stack, result)
        } else {
            num := stringToInt(token)
            stack = append(stack, num)
        }
    }
    
    return stack[0]
}
```

---

## 🎯 Problem-Solving Framework for Stacks

### 1. Identify Stack Use Cases
- **Need to reverse order?** → Stack
- **Matching pairs/brackets?** → Stack  
- **Undo operations?** → Stack
- **Next greater/smaller element?** → Monotonic Stack
- **Expression evaluation?** → Stack

### 2. Choose What to Store
- **Values:** For simple operations
- **Indices:** For position-dependent problems
- **Pairs:** For complex state tracking

### 3. Determine Stack Property
- **Monotonic Increasing:** For next smaller element
- **Monotonic Decreasing:** For next greater element
- **Regular:** For matching and evaluation

---

## 📊 Self-Assessment Checklist

After solving each problem, verify:

**Understanding:**
- [ ] Do I understand when and why to use a stack?
- [ ] Can I identify the stack pattern used?
- [ ] Do I know what to store in the stack?

**Implementation:**
- [ ] Are my push/pop operations correct?
- [ ] Do I handle empty stack cases?
- [ ] Is my time complexity optimal?

**Testing:**
- [ ] Did I test with empty input?
- [ ] Did I test edge cases (single element, all same elements)?
- [ ] Does my solution handle all given examples?

---

## 🚀 Tomorrow Preview: Queues

**Day 6 Topics:**
- Queue implementations (circular array, linked list)
- BFS traversal concepts
- Queue-based problem solving
- Implementing stack using queues

**Preparation:**
- Review FIFO (First In, First Out) principle
- Think about circular arrays and modular arithmetic
- Consider when order of processing matters

**Key Insight:**
Stacks are about **reversing order** and **backtracking**. Once you master these patterns, many complex problems become straightforward! 

Keep practicing and focus on recognizing when to use each pattern. Tomorrow we'll explore queues and see how order of processing creates different algorithmic opportunities! 📚
