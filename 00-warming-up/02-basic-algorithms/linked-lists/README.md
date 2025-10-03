# 🔗 Linked Lists Fundamentals
## Week 1, Days 3-4 | Pointer Mastery for Junior Engineers

### 🎯 Learning Objectives
By the end of this section, you will:
- [ ] Understand pointer concepts and memory references
- [ ] Implement singly and doubly linked lists from scratch
- [ ] Master basic operations: insert, delete, search, reverse
- [ ] Solve 8-10 linked list problems with confidence
- [ ] Compare arrays vs linked lists trade-offs

---

## 📚 Concepts Overview

### 🔸 What is a Linked List?

**Definition:**
- Collection of nodes where each node contains data and a reference/pointer to the next node
- Nodes can be stored anywhere in memory (unlike arrays)
- Dynamic size - can grow and shrink during runtime

**Key Characteristics:**
- **Sequential Access:** Must traverse from head to reach any element - O(n)
- **Dynamic Size:** Can add/remove elements without declaring size upfront
- **Memory Efficient:** Only allocates memory when needed
- **No Random Access:** Cannot jump directly to index like arrays

**Node Structure in Golang:**
```go
type ListNode struct {
    Val  int
    Next *ListNode
}
```

### 🔸 Singly Linked List Operations

**Basic Operations Time Complexity:**
- **Search:** O(n) - must traverse from head
- **Insert at beginning:** O(1) - just update head
- **Insert at end:** O(n) - must find tail first
- **Delete from beginning:** O(1) - update head
- **Delete from middle:** O(n) - must find node first

**Memory Layout:**
```
Head -> [Data|Next] -> [Data|Next] -> [Data|NULL]
         Node1          Node2         Node3
```

### 🔸 Arrays vs Linked Lists Comparison

| Operation | Array | Linked List |
|-----------|-------|-------------|
| Access by index | O(1) | O(n) |
| Search | O(n) | O(n) |
| Insert at beginning | O(n) | O(1) |
| Insert at end | O(1)* | O(n) |
| Delete from beginning | O(n) | O(1) |
| Memory usage | Contiguous | Scattered + pointers |
| Cache performance | Better | Worse |

*Amortized for dynamic arrays

---

## 🛠️ Implementation in Golang

### Basic Singly Linked List

```go
package main

import "fmt"

// ListNode represents a node in the linked list
type ListNode struct {
    Val  int
    Next *ListNode
}

// LinkedList represents the linked list structure
type LinkedList struct {
    Head *ListNode
    Size int
}

// NewLinkedList creates a new empty linked list
func NewLinkedList() *LinkedList {
    return &LinkedList{
        Head: nil,
        Size: 0,
    }
}

// Insert adds a new node at the beginning - O(1)
func (ll *LinkedList) Insert(val int) {
    newNode := &ListNode{
        Val:  val,
        Next: ll.Head,
    }
    ll.Head = newNode
    ll.Size++
}

// Append adds a new node at the end - O(n)
func (ll *LinkedList) Append(val int) {
    newNode := &ListNode{
        Val:  val,
        Next: nil,
    }
    
    if ll.Head == nil {
        ll.Head = newNode
        ll.Size++
        return
    }
    
    current := ll.Head
    for current.Next != nil {
        current = current.Next
    }
    current.Next = newNode
    ll.Size++
}

// Delete removes the first occurrence of val - O(n)
func (ll *LinkedList) Delete(val int) bool {
    if ll.Head == nil {
        return false
    }
    
    // If head node is to be deleted
    if ll.Head.Val == val {
        ll.Head = ll.Head.Next
        ll.Size--
        return true
    }
    
    current := ll.Head
    for current.Next != nil {
        if current.Next.Val == val {
            current.Next = current.Next.Next
            ll.Size--
            return true
        }
        current = current.Next
    }
    
    return false // Value not found
}

// Search finds if a value exists in the list - O(n)
func (ll *LinkedList) Search(val int) bool {
    current := ll.Head
    for current != nil {
        if current.Val == val {
            return true
        }
        current = current.Next
    }
    return false
}

// Display prints all elements in the list - O(n)
func (ll *LinkedList) Display() {
    if ll.Head == nil {
        fmt.Println("List is empty")
        return
    }
    
    current := ll.Head
    fmt.Print("List: ")
    for current != nil {
        fmt.Print(current.Val)
        if current.Next != nil {
            fmt.Print(" -> ")
        }
        current = current.Next
    }
    fmt.Println()
}

// Reverse reverses the linked list - O(n)
func (ll *LinkedList) Reverse() {
    var prev *ListNode = nil
    current := ll.Head
    
    for current != nil {
        nextTemp := current.Next  // Store next node
        current.Next = prev       // Reverse the link
        prev = current           // Move prev forward
        current = nextTemp       // Move current forward
    }
    
    ll.Head = prev
}

// GetSize returns the number of elements - O(1)
func (ll *LinkedList) GetSize() int {
    return ll.Size
}
```

### Doubly Linked List Implementation

```go
// DoublyListNode represents a node in doubly linked list
type DoublyListNode struct {
    Val  int
    Next *DoublyListNode
    Prev *DoublyListNode
}

// DoublyLinkedList represents doubly linked list
type DoublyLinkedList struct {
    Head *DoublyListNode
    Tail *DoublyListNode
    Size int
}

// NewDoublyLinkedList creates new doubly linked list
func NewDoublyLinkedList() *DoublyLinkedList {
    return &DoublyLinkedList{
        Head: nil,
        Tail: nil,
        Size: 0,
    }
}

// InsertFront adds node at the beginning - O(1)
func (dll *DoublyLinkedList) InsertFront(val int) {
    newNode := &DoublyListNode{
        Val:  val,
        Next: dll.Head,
        Prev: nil,
    }
    
    if dll.Head != nil {
        dll.Head.Prev = newNode
    } else {
        dll.Tail = newNode // First node
    }
    
    dll.Head = newNode
    dll.Size++
}

// InsertBack adds node at the end - O(1)
func (dll *DoublyLinkedList) InsertBack(val int) {
    newNode := &DoublyListNode{
        Val:  val,
        Next: nil,
        Prev: dll.Tail,
    }
    
    if dll.Tail != nil {
        dll.Tail.Next = newNode
    } else {
        dll.Head = newNode // First node
    }
    
    dll.Tail = newNode
    dll.Size++
}

// DeleteFront removes first node - O(1)
func (dll *DoublyLinkedList) DeleteFront() bool {
    if dll.Head == nil {
        return false
    }
    
    if dll.Head == dll.Tail {
        // Only one node
        dll.Head = nil
        dll.Tail = nil
    } else {
        dll.Head = dll.Head.Next
        dll.Head.Prev = nil
    }
    
    dll.Size--
    return true
}

// DeleteBack removes last node - O(1)
func (dll *DoublyLinkedList) DeleteBack() bool {
    if dll.Tail == nil {
        return false
    }
    
    if dll.Head == dll.Tail {
        // Only one node
        dll.Head = nil
        dll.Tail = nil
    } else {
        dll.Tail = dll.Tail.Prev
        dll.Tail.Next = nil
    }
    
    dll.Size--
    return true
}
```

---

## 🎯 Essential Problem-Solving Patterns

### Pattern 1: Two Pointers (Fast & Slow)

**When to Use:**
- Finding middle of linked list
- Detecting cycles
- Finding nth node from end

**Example: Find Middle Node**
```go
func findMiddle(head *ListNode) *ListNode {
    if head == nil {
        return nil
    }
    
    slow := head
    fast := head
    
    // Fast moves 2 steps, slow moves 1 step
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    
    return slow // Slow will be at middle
}
```

### Pattern 2: Dummy Head Technique

**When to Use:**
- When head might change (insertions/deletions)
- Simplifies edge case handling
- Makes code cleaner and less error-prone

**Example: Remove Elements**
```go
func removeElements(head *ListNode, val int) *ListNode {
    // Create dummy head to handle edge cases
    dummy := &ListNode{Next: head}
    current := dummy
    
    for current.Next != nil {
        if current.Next.Val == val {
            current.Next = current.Next.Next
        } else {
            current = current.Next
        }
    }
    
    return dummy.Next
}
```

### Pattern 3: Reverse Linked List

**When to Use:**
- Palindrome checking
- Reversing sublists
- Various manipulation problems

**Example: Reverse List Iteratively**
```go
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode = nil
    current := head
    
    for current != nil {
        nextTemp := current.Next
        current.Next = prev
        prev = current
        current = nextTemp
    }
    
    return prev
}
```

**Example: Reverse List Recursively**
```go
func reverseListRecursive(head *ListNode) *ListNode {
    // Base case
    if head == nil || head.Next == nil {
        return head
    }
    
    // Recursively reverse the rest
    newHead := reverseListRecursive(head.Next)
    
    // Reverse current connection
    head.Next.Next = head
    head.Next = nil
    
    return newHead
}
```

---

## 💻 Practice Problems

### Easy Level (Day 3)

#### Problem 1: Remove Linked List Elements (LeetCode #203)
Remove all elements from linked list that have value `val`.

```go
func removeElements(head *ListNode, val int) *ListNode {
    // Your implementation here
    return nil
}

// Test cases
// Input: head = [1,2,6,3,4,5,6], val = 6
// Output: [1,2,3,4,5]
```

#### Problem 2: Reverse Linked List (LeetCode #206)
```go
func reverseList(head *ListNode) *ListNode {
    // Implement both iterative and recursive solutions
    return nil
}
```

#### Problem 3: Merge Two Sorted Lists (LeetCode #21)
```go
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
    // Merge two sorted linked lists
    return nil
}
```

### Medium Level (Day 4)

#### Problem 4: Linked List Cycle (LeetCode #141)
```go
func hasCycle(head *ListNode) bool {
    // Use Floyd's cycle detection algorithm
    return false
}
```

#### Problem 5: Remove Nth Node From End (LeetCode #19)
```go
func removeNthFromEnd(head *ListNode, n int) *ListNode {
    // Use two pointers with n gap
    return nil
}
```

#### Problem 6: Palindrome Linked List (LeetCode #234)
```go
func isPalindrome(head *ListNode) bool {
    // Find middle, reverse second half, compare
    return false
}
```

---

## 🔍 Common Patterns & Techniques

### 1. Finding Length
```go
func getLength(head *ListNode) int {
    length := 0
    current := head
    for current != nil {
        length++
        current = current.Next
    }
    return length
}
```

### 2. Finding Nth Node
```go
func getNthNode(head *ListNode, n int) *ListNode {
    current := head
    for i := 0; i < n && current != nil; i++ {
        current = current.Next
    }
    return current
}
```

### 3. Cycle Detection (Floyd's Algorithm)
```go
func detectCycle(head *ListNode) *ListNode {
    if head == nil || head.Next == nil {
        return nil
    }
    
    slow := head
    fast := head
    
    // Phase 1: Detect if cycle exists
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        
        if slow == fast {
            break // Cycle detected
        }
    }
    
    // No cycle found
    if fast == nil || fast.Next == nil {
        return nil
    }
    
    // Phase 2: Find cycle start
    slow = head
    for slow != fast {
        slow = slow.Next
        fast = fast.Next
    }
    
    return slow // Start of cycle
}
```

---

## 📊 Time & Space Complexity Analysis

### Singly Linked List Operations
| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Insert at head | O(1) | O(1) |
| Insert at tail | O(n) | O(1) |
| Insert at position | O(n) | O(1) |
| Delete from head | O(1) | O(1) |
| Delete from tail | O(n) | O(1) |
| Search | O(n) | O(1) |
| Reverse | O(n) | O(1) |

### Doubly Linked List Operations
| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Insert at head | O(1) | O(1) |
| Insert at tail | O(1) | O(1) |
| Delete from head | O(1) | O(1) |
| Delete from tail | O(1) | O(1) |
| Search | O(n) | O(1) |

---

## 🎓 Study Tips & Best Practices

### 1. Drawing is Essential
```
Before coding, always draw the list:

Original: A -> B -> C -> D -> NULL
After deletion of B: A -> C -> D -> NULL

Step by step:
1. Find node before B (A)
2. Set A.Next = B.Next (which is C)
3. B is now disconnected
```

### 2. Handle Edge Cases
```go
// Always check for these cases:
// 1. Empty list (head == nil)
// 2. Single node list
// 3. Operation on head node
// 4. Operation on tail node
// 5. Node not found
```

### 3. Pointer Safety in Golang
```go
// Always check nil before dereferencing
if current != nil && current.Next != nil {
    // Safe to access current.Next.Val
}

// Use pointers to pointers for head changes
func insertAtHead(head **ListNode, val int) {
    newNode := &ListNode{Val: val, Next: *head}
    *head = newNode
}
```

### 4. Memory Management
```go
// In Golang, garbage collector handles memory
// But be aware of circular references in cyclic lists
// Break cycles before losing references

func breakCycle(head *ListNode) {
    // ... detect cycle first
    // Then break it by setting one Next pointer to nil
}
```

---

## 🚀 Next Steps

### Day 3 Goals
- [ ] Implement complete singly linked list
- [ ] Solve remove elements and reverse problems
- [ ] Master dummy head technique
- [ ] Practice drawing solutions on paper

### Day 4 Goals
- [ ] Implement doubly linked list
- [ ] Master two-pointer technique
- [ ] Solve cycle detection problems
- [ ] Complete palindrome linked list

### Preparation for Day 5 (Stacks)
- Understand LIFO (Last In, First Out) principle
- Think about call stack and function calls
- Review recursion concepts
- Consider when you need to "undo" operations

Remember: Linked lists are the foundation for many advanced data structures. Master pointer manipulation here, and trees, graphs, and other structures will be much easier! 🔗

---

## 🔧 Debug Tips

### Common Mistakes:
1. **Null pointer dereference** - Always check `current != nil`
2. **Lost references** - Store `next` before modifying pointers
3. **Infinite loops** - Make sure you're advancing pointers
4. **Memory leaks** - In languages with manual memory management

### Debugging Techniques:
1. **Print statements** - Add debug prints to track pointer movement
2. **Step through manually** - Trace execution on paper
3. **Draw the list** - Visual representation helps catch errors
4. **Test edge cases first** - Empty list, single node, etc.

**Master these fundamentals and you'll be ready for any linked list interview question!** 💪
