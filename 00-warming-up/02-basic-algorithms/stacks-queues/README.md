# 📚 Stacks & Queues Fundamentals
## Week 1, Days 5-6 | LIFO and FIFO Data Structures

### 🎯 Learning Objectives
By the end of this section, you will:
- [ ] Understand LIFO (Stack) and FIFO (Queue) principles
- [ ] Implement stacks and queues using arrays and linked lists
- [ ] Master stack applications: parentheses matching, expression evaluation
- [ ] Understand queue applications: BFS traversal, task scheduling
- [ ] Solve 8-10 stack and queue problems confidently

---

## 📚 Stack Fundamentals

### 🔸 What is a Stack?

**Definition:**
- Linear data structure that follows LIFO (Last In, First Out) principle
- Elements are added and removed from the same end called "top"
- Think of it like a stack of plates - you can only add/remove from the top

**Key Operations:**
- **Push:** Add element to top - O(1)
- **Pop:** Remove and return top element - O(1)
- **Top/Peek:** View top element without removing - O(1)
- **IsEmpty:** Check if stack is empty - O(1)
- **Size:** Get number of elements - O(1)

**Real-world Applications:**
- Function call stack (recursion)
- Undo operations in editors
- Expression evaluation and syntax parsing
- Browser back button functionality
- Balancing parentheses in code

### 🔸 Stack Implementation Approaches

**1. Array-based Stack:**
```go
type ArrayStack struct {
    items []int
    top   int
}
```

**Pros:** Simple, cache-friendly, no extra memory for pointers
**Cons:** Fixed size (if using fixed array), resizing overhead

**2. Linked List-based Stack:**
```go
type Node struct {
    data int
    next *Node
}

type LinkedStack struct {
    top *Node
}
```

**Pros:** Dynamic size, no wasted memory
**Cons:** Extra memory for pointers, less cache-friendly

---

## 📚 Queue Fundamentals

### 🔸 What is a Queue?

**Definition:**
- Linear data structure that follows FIFO (First In, First Out) principle
- Elements are added at one end (rear) and removed from other end (front)
- Think of it like a line of people - first person in line is served first

**Key Operations:**
- **Enqueue:** Add element to rear - O(1)
- **Dequeue:** Remove and return front element - O(1)
- **Front:** View front element without removing - O(1)
- **IsEmpty:** Check if queue is empty - O(1)
- **Size:** Get number of elements - O(1)

**Real-world Applications:**
- Task scheduling in operating systems
- BFS (Breadth-First Search) in graphs
- Print job management
- Handling requests in web servers
- Buffer for data streams

### 🔸 Queue Implementation Approaches

**1. Array-based Queue (Circular):**
```go
type CircularQueue struct {
    items []int
    front int
    rear  int
    size  int
    capacity int
}
```

**2. Linked List-based Queue:**
```go
type QueueNode struct {
    data int
    next *QueueNode
}

type LinkedQueue struct {
    front *QueueNode
    rear  *QueueNode
}
```

---

## 🛠️ Implementation in Golang

### Stack Implementation

```go
package main

import (
    "errors"
    "fmt"
)

// ===== ARRAY-BASED STACK =====

type ArrayStack struct {
    items    []int
    capacity int
}

// NewArrayStack creates a new array-based stack
func NewArrayStack(capacity int) *ArrayStack {
    return &ArrayStack{
        items:    make([]int, 0, capacity),
        capacity: capacity,
    }
}

// Push adds element to top of stack - O(1)
func (s *ArrayStack) Push(item int) error {
    if len(s.items) >= s.capacity {
        return errors.New("stack overflow")
    }
    s.items = append(s.items, item)
    return nil
}

// Pop removes and returns top element - O(1)
func (s *ArrayStack) Pop() (int, error) {
    if s.IsEmpty() {
        return 0, errors.New("stack underflow")
    }
    index := len(s.items) - 1
    item := s.items[index]
    s.items = s.items[:index]
    return item, nil
}

// Top returns top element without removing - O(1)
func (s *ArrayStack) Top() (int, error) {
    if s.IsEmpty() {
        return 0, errors.New("stack is empty")
    }
    return s.items[len(s.items)-1], nil
}

// IsEmpty checks if stack is empty - O(1)
func (s *ArrayStack) IsEmpty() bool {
    return len(s.items) == 0
}

// Size returns number of elements - O(1)
func (s *ArrayStack) Size() int {
    return len(s.items)
}

// Display prints all elements
func (s *ArrayStack) Display() {
    if s.IsEmpty() {
        fmt.Println("Stack is empty")
        return
    }
    fmt.Print("Stack (top to bottom): ")
    for i := len(s.items) - 1; i >= 0; i-- {
        fmt.Printf("%d ", s.items[i])
    }
    fmt.Println()
}

// ===== LINKED LIST-BASED STACK =====

type StackNode struct {
    data int
    next *StackNode
}

type LinkedStack struct {
    top *StackNode
}

// NewLinkedStack creates a new linked list-based stack
func NewLinkedStack() *LinkedStack {
    return &LinkedStack{top: nil}
}

// Push adds element to top - O(1)
func (s *LinkedStack) Push(item int) {
    newNode := &StackNode{
        data: item,
        next: s.top,
    }
    s.top = newNode
}

// Pop removes and returns top element - O(1)
func (s *LinkedStack) Pop() (int, error) {
    if s.IsEmpty() {
        return 0, errors.New("stack underflow")
    }
    item := s.top.data
    s.top = s.top.next
    return item, nil
}

// Top returns top element without removing - O(1)
func (s *LinkedStack) Top() (int, error) {
    if s.IsEmpty() {
        return 0, errors.New("stack is empty")
    }
    return s.top.data, nil
}

// IsEmpty checks if stack is empty - O(1)
func (s *LinkedStack) IsEmpty() bool {
    return s.top == nil
}

// Size returns number of elements - O(n)
func (s *LinkedStack) Size() int {
    count := 0
    current := s.top
    for current != nil {
        count++
        current = current.next
    }
    return count
}

// Display prints all elements
func (s *LinkedStack) Display() {
    if s.IsEmpty() {
        fmt.Println("Stack is empty")
        return
    }
    fmt.Print("Stack (top to bottom): ")
    current := s.top
    for current != nil {
        fmt.Printf("%d ", current.data)
        current = current.next
    }
    fmt.Println()
}
```

### Queue Implementation

```go
// ===== CIRCULAR ARRAY-BASED QUEUE =====

type CircularQueue struct {
    items    []int
    front    int
    rear     int
    size     int
    capacity int
}

// NewCircularQueue creates a new circular queue
func NewCircularQueue(capacity int) *CircularQueue {
    return &CircularQueue{
        items:    make([]int, capacity),
        front:    0,
        rear:     -1,
        size:     0,
        capacity: capacity,
    }
}

// Enqueue adds element to rear - O(1)
func (q *CircularQueue) Enqueue(item int) error {
    if q.IsFull() {
        return errors.New("queue is full")
    }
    q.rear = (q.rear + 1) % q.capacity
    q.items[q.rear] = item
    q.size++
    return nil
}

// Dequeue removes and returns front element - O(1)
func (q *CircularQueue) Dequeue() (int, error) {
    if q.IsEmpty() {
        return 0, errors.New("queue is empty")
    }
    item := q.items[q.front]
    q.front = (q.front + 1) % q.capacity
    q.size--
    return item, nil
}

// Front returns front element without removing - O(1)
func (q *CircularQueue) Front() (int, error) {
    if q.IsEmpty() {
        return 0, errors.New("queue is empty")
    }
    return q.items[q.front], nil
}

// IsEmpty checks if queue is empty - O(1)
func (q *CircularQueue) IsEmpty() bool {
    return q.size == 0
}

// IsFull checks if queue is full - O(1)
func (q *CircularQueue) IsFull() bool {
    return q.size == q.capacity
}

// Size returns number of elements - O(1)
func (q *CircularQueue) Size() int {
    return q.size
}

// Display prints all elements
func (q *CircularQueue) Display() {
    if q.IsEmpty() {
        fmt.Println("Queue is empty")
        return
    }
    fmt.Print("Queue (front to rear): ")
    for i := 0; i < q.size; i++ {
        index := (q.front + i) % q.capacity
        fmt.Printf("%d ", q.items[index])
    }
    fmt.Println()
}

// ===== LINKED LIST-BASED QUEUE =====

type QueueNode struct {
    data int
    next *QueueNode
}

type LinkedQueue struct {
    front *QueueNode
    rear  *QueueNode
}

// NewLinkedQueue creates a new linked queue
func NewLinkedQueue() *LinkedQueue {
    return &LinkedQueue{
        front: nil,
        rear:  nil,
    }
}

// Enqueue adds element to rear - O(1)
func (q *LinkedQueue) Enqueue(item int) {
    newNode := &QueueNode{
        data: item,
        next: nil,
    }
    
    if q.rear == nil {
        // First element
        q.front = newNode
        q.rear = newNode
    } else {
        q.rear.next = newNode
        q.rear = newNode
    }
}

// Dequeue removes and returns front element - O(1)
func (q *LinkedQueue) Dequeue() (int, error) {
    if q.IsEmpty() {
        return 0, errors.New("queue is empty")
    }
    
    item := q.front.data
    q.front = q.front.next
    
    if q.front == nil {
        // Queue became empty
        q.rear = nil
    }
    
    return item, nil
}

// Front returns front element without removing - O(1)
func (q *LinkedQueue) Front() (int, error) {
    if q.IsEmpty() {
        return 0, errors.New("queue is empty")
    }
    return q.front.data, nil
}

// IsEmpty checks if queue is empty - O(1)
func (q *LinkedQueue) IsEmpty() bool {
    return q.front == nil
}

// Size returns number of elements - O(n)
func (q *LinkedQueue) Size() int {
    count := 0
    current := q.front
    for current != nil {
        count++
        current = current.next
    }
    return count
}

// Display prints all elements
func (q *LinkedQueue) Display() {
    if q.IsEmpty() {
        fmt.Println("Queue is empty")
        return
    }
    fmt.Print("Queue (front to rear): ")
    current := q.front
    for current != nil {
        fmt.Printf("%d ", current.data)
        current = current.next
    }
    fmt.Println()
}
```

---

## 🎯 Essential Problem-Solving Patterns

### Pattern 1: Parentheses Matching

**When to Use:** Validating balanced brackets, expressions
**Key Insight:** Use stack to match opening and closing brackets

```go
func isValid(s string) bool {
    stack := []rune{}
    pairs := map[rune]rune{
        ')': '(',
        '}': '{',
        ']': '[',
    }
    
    for _, char := range s {
        if char == '(' || char == '{' || char == '[' {
            // Opening bracket - push to stack
            stack = append(stack, char)
        } else {
            // Closing bracket - check if matches top
            if len(stack) == 0 {
                return false
            }
            if stack[len(stack)-1] != pairs[char] {
                return false
            }
            stack = stack[:len(stack)-1] // Pop
        }
    }
    
    return len(stack) == 0
}
```

### Pattern 2: Monotonic Stack

**When to Use:** Finding next greater/smaller element, histogram problems
**Key Insight:** Maintain stack in increasing/decreasing order

```go
func nextGreaterElement(nums []int) []int {
    result := make([]int, len(nums))
    stack := []int{} // Stack of indices
    
    for i := len(nums) - 1; i >= 0; i-- {
        // Pop smaller elements
        for len(stack) > 0 && nums[stack[len(stack)-1]] <= nums[i] {
            stack = stack[:len(stack)-1]
        }
        
        if len(stack) == 0 {
            result[i] = -1
        } else {
            result[i] = nums[stack[len(stack)-1]]
        }
        
        stack = append(stack, i)
    }
    
    return result
}
```

### Pattern 3: Stack for Recursion Simulation

**When to Use:** Converting recursive algorithms to iterative
**Key Insight:** Stack simulates function call stack

```go
func iterativeInorder(root *TreeNode) []int {
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

### Pattern 4: Queue for BFS

**When to Use:** Level-order traversal, shortest path problems
**Key Insight:** Process nodes level by level

```go
func levelOrder(root *TreeNode) [][]int {
    if root == nil {
        return [][]int{}
    }
    
    result := [][]int{}
    queue := []*TreeNode{root}
    
    for len(queue) > 0 {
        levelSize := len(queue)
        level := []int{}
        
        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:] // Dequeue
            level = append(level, node.Val)
            
            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        
        result = append(result, level)
    }
    
    return result
}
```

---

## 📊 Time & Space Complexity Comparison

### Stack Operations
| Implementation | Push | Pop | Top | IsEmpty | Size |
|----------------|------|-----|-----|---------|------|
| Array-based | O(1)* | O(1) | O(1) | O(1) | O(1) |
| Linked List | O(1) | O(1) | O(1) | O(1) | O(n) |

*Amortized for dynamic arrays

### Queue Operations
| Implementation | Enqueue | Dequeue | Front | IsEmpty | Size |
|----------------|---------|---------|-------|---------|------|
| Circular Array | O(1) | O(1) | O(1) | O(1) | O(1) |
| Linked List | O(1) | O(1) | O(1) | O(1) | O(n) |

### Space Complexity
- **Array-based:** O(n) where n is capacity
- **Linked List-based:** O(n) where n is number of elements + pointer overhead

---

## 🎓 When to Use Which?

### Use Stack When:
- Need to reverse order of processing
- Implementing recursion iteratively  
- Matching pairs (parentheses, tags)
- Undo operations
- Expression evaluation
- DFS traversal

### Use Queue When:
- First-come, first-served processing
- BFS traversal
- Task scheduling
- Buffer for streaming data
- Level-order processing

### Array vs Linked List Implementation:

**Choose Array-based when:**
- Memory is a concern (no pointer overhead)
- Cache performance matters
- Maximum size is known in advance
- Need O(1) size operation

**Choose Linked List-based when:**
- Dynamic size requirements
- Memory allocation flexibility
- Unknown maximum size
- Frequent insertions/deletions

---

## 🚀 Next Steps

### Day 5 Goals (Stacks)
- [ ] Implement both array and linked list stacks
- [ ] Solve parentheses matching problems
- [ ] Master monotonic stack patterns
- [ ] Practice expression evaluation

### Day 6 Goals (Queues)  
- [ ] Implement circular and linked queues
- [ ] Understand BFS traversal concept
- [ ] Solve queue-based problems
- [ ] Practice implementing stack using queues

### Preparation for Week 2 (Hash Tables)
- Review array indexing and hash functions
- Understand collision resolution concepts
- Think about key-value pair applications

**Stack and Queue are fundamental building blocks for more complex algorithms. Master these, and graph traversals, tree operations, and recursion will become much clearer!** 📚

---

## 🔧 Implementation Tips

### Stack Tips:
```go
// Always check for underflow
if stack.IsEmpty() {
    return errors.New("stack underflow")
}

// Use slices efficiently in Go
stack = append(stack, item)    // Push
item = stack[len(stack)-1]     // Top
stack = stack[:len(stack)-1]   // Pop
```

### Queue Tips:
```go
// Circular queue index calculation
nextIndex = (currentIndex + 1) % capacity

// Linked queue - update both front and rear
if q.front == nil {
    q.front = newNode
    q.rear = newNode
}
```

**Remember: Stacks and queues are about controlling the order of processing. Master this concept and you'll solve many algorithmic problems with ease!** 💪
