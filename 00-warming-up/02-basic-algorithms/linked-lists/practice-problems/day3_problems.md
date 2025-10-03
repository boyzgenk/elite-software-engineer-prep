# 📝 Day 3 Practice Problems: Linked Lists Fundamentals
## Pointer Manipulation and Basic Operations

### 🎯 Today's Goals
- Master basic linked list operations
- Understand pointer manipulation in Golang
- Practice drawing linked list operations
- Solve 4-5 fundamental linked list problems
- Build confidence with nil pointer handling

---

## Problem 1: Remove Linked List Elements (Easy) ⭐⭐
**LeetCode #203**

Remove all elements from a linked list of integers that have value `val`.

**Example:**
```
Input: head = [1,2,6,3,4,5,6], val = 6
Output: [1,2,3,4,5]

Input: head = [], val = 1
Output: []

Input: head = [7,7,7,7], val = 7
Output: []
```

**Your Task:**
1. Handle the case where head node needs to be removed
2. Use dummy head technique to simplify code
3. Make sure to handle empty list edge case

**Template:**
```go
func removeElements(head *ListNode, val int) *ListNode {
    // Create dummy head to handle edge cases
    dummy := &ListNode{Next: head}
    current := dummy
    
    // Your code here
    
    return dummy.Next
}

// Test cases
// removeElements([1,2,6,3,4,5,6], 6) should return [1,2,3,4,5]
// removeElements([], 1) should return []
// removeElements([7,7,7,7], 7) should return []
```

**Learning Focus:**
- Dummy head technique
- Pointer manipulation
- Edge case handling

---

## Problem 2: Reverse Linked List (Easy) ⭐⭐⭐
**LeetCode #206**

Reverse a singly linked list.

**Example:**
```
Input: head = [1,2,3,4,5]
Output: [5,4,3,2,1]

Input: head = [1,2]
Output: [2,1]

Input: head = []
Output: []
```

**Your Task:**
1. Implement iterative solution first
2. Then try recursive solution
3. Understand the pointer manipulation pattern

**Template:**
```go
// Iterative approach
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode = nil
    current := head
    
    // Your iterative solution here
    
    return prev
}

// Recursive approach (bonus)
func reverseListRecursive(head *ListNode) *ListNode {
    // Base case
    if head == nil || head.Next == nil {
        return head
    }
    
    // Your recursive solution here
    
    return nil // Replace with actual return
}

// Test cases
// reverseList([1,2,3,4,5]) should return [5,4,3,2,1]
// reverseList([1,2]) should return [2,1]
// reverseList([]) should return []
```

**Learning Focus:**
- Three-pointer technique (prev, current, next)
- Recursive thinking
- Understanding reversal pattern

---

## Problem 3: Merge Two Sorted Lists (Easy) ⭐⭐
**LeetCode #21**

Merge two sorted linked lists and return it as a sorted list.

**Example:**
```
Input: list1 = [1,2,4], list2 = [1,3,4]
Output: [1,1,2,3,4,4]

Input: list1 = [], list2 = []
Output: []

Input: list1 = [], list2 = [0]
Output: [0]
```

**Your Task:**
1. Use dummy head to build result list
2. Compare values and choose smaller one
3. Handle remaining elements after one list is exhausted

**Template:**
```go
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
    // Create dummy head for result
    dummy := &ListNode{}
    current := dummy
    
    // Your merge logic here
    
    return dummy.Next
}

// Test cases
// mergeTwoLists([1,2,4], [1,3,4]) should return [1,1,2,3,4,4]
// mergeTwoLists([], []) should return []
// mergeTwoLists([], [0]) should return [0]
```

**Learning Focus:**
- Two-pointer traversal
- Building result list
- Handling different list lengths

---

## Problem 4: Middle of the Linked List (Easy) ⭐⭐
**LeetCode #876**

Find the middle node of a linked list. If there are two middle nodes, return the second middle node.

**Example:**
```
Input: head = [1,2,3,4,5]
Output: [3,4,5] (return node with value 3)

Input: head = [1,2,3,4,5,6]
Output: [4,5,6] (return node with value 4)
```

**Your Task:**
1. Use fast and slow pointer technique
2. When fast reaches end, slow will be at middle
3. Handle edge cases (empty list, single node)

**Template:**
```go
func middleNode(head *ListNode) *ListNode {
    if head == nil {
        return nil
    }
    
    slow := head
    fast := head
    
    // Your two-pointer logic here
    
    return slow
}

// Test cases
// middleNode([1,2,3,4,5]) should return node with value 3
// middleNode([1,2,3,4,5,6]) should return node with value 4
// middleNode([1]) should return node with value 1
```

**Learning Focus:**
- Fast and slow pointer technique
- Mathematical insight (why this works)
- Handling even vs odd length lists

---

## Problem 5: Delete Node in a Linked List (Easy) ⭐⭐⭐
**LeetCode #237**

Delete a node (except the tail) in a singly linked list, given only access to that node.

**Example:**
```
Input: head = [4,5,1,9], node = 5
Output: [4,1,9]
Explanation: You are given the second node with value 5, 
the linked list should become 4 -> 1 -> 9 after calling your function.
```

**Your Task:**
1. You don't have access to head, only the node to delete
2. Think creatively - you can't actually delete the node
3. What can you do instead?

**Template:**
```go
func deleteNode(node *ListNode) {
    // You only have access to the node to be "deleted"
    // Think: what if you copy the next node's value?
    
    // Your creative solution here
}

// Test cases
// Given node with value 5 in list [4,5,1,9]
// After deleteNode(node), list should be [4,1,9]
```

**Learning Focus:**
- Creative problem solving
- Understanding the difference between "deleting" and "removing"
- Pointer manipulation tricks

---

## Bonus Problem: Intersection of Two Linked Lists (Easy) ⭐⭐⭐⭐
**LeetCode #160**

Find the node at which the intersection of two singly linked lists begins.

**Example:**
```
Input: intersectVal = 8, listA = [4,1,8,4,5], listB = [5,6,1,8,4,5], skipA = 2, skipB = 3
Output: Reference to the node with value = 8
```

**Your Task:**
1. Find where two lists intersect (same node reference, not just same value)
2. Use two-pointer technique with length difference handling
3. Handle case where lists don't intersect

**Template:**
```go
func getIntersectionNode(headA, headB *ListNode) *ListNode {
    if headA == nil || headB == nil {
        return nil
    }
    
    // Calculate lengths or use two-pointer trick
    // Your intersection logic here
    
    return nil
}

// This is more advanced - focus on understanding the approach
```

**Learning Focus:**
- Advanced two-pointer technique
- Understanding node references vs values
- Mathematical approach to synchronize traversal

---

## 🧠 Problem-Solving Patterns Review

### Pattern 1: Dummy Head Technique
```go
// Use when head might change or for edge case handling
func someOperation(head *ListNode) *ListNode {
    dummy := &ListNode{Next: head}
    current := dummy
    
    // Process nodes
    for current.Next != nil {
        // Your logic here
        current = current.Next
    }
    
    return dummy.Next
}
```

### Pattern 2: Two Pointers (Fast & Slow)
```go
// Use for finding middle, cycle detection, etc.
func twoPointerOperation(head *ListNode) *ListNode {
    slow := head
    fast := head
    
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    
    return slow // Or use the relationship between pointers
}
```

### Pattern 3: Three Pointers (Reversal)
```go
// Use for reversing linked list
func reversePattern(head *ListNode) *ListNode {
    var prev *ListNode = nil
    current := head
    
    for current != nil {
        next := current.Next    // Store next
        current.Next = prev     // Reverse link
        prev = current          // Move prev
        current = next          // Move current
    }
    
    return prev
}
```

---

## 🎯 Daily Practice Structure

### Step-by-Step Approach for Each Problem:

#### 1. Understand (5 minutes)
- Draw the linked list on paper
- Identify what changes need to happen
- List edge cases (empty list, single node, etc.)

#### 2. Plan (10 minutes)
- Choose appropriate pattern (dummy head, two pointers, etc.)
- Sketch the algorithm steps
- Identify loop conditions and termination

#### 3. Code (20 minutes)
- Start with basic structure
- Handle main logic first
- Add edge case handling
- Use meaningful variable names

#### 4. Test (10 minutes)
- Trace through your algorithm with given examples
- Test edge cases manually
- Check for nil pointer access

#### 5. Draw & Verify (5 minutes)
- Draw before and after states
- Verify pointer connections are correct
- Make sure no nodes are lost

---

## 📊 Self-Assessment Checklist

After solving each problem, verify:

**Understanding:**
- [ ] Can I draw the operation step by step?
- [ ] Do I understand why this approach works?
- [ ] Have I identified all edge cases?

**Implementation:**
- [ ] Does my code handle nil pointers safely?
- [ ] Are all nodes properly connected?
- [ ] Did I test with empty lists and single nodes?

**Complexity:**
- [ ] What's the time complexity? Can I explain it?
- [ ] What's the space complexity?
- [ ] Is this the most efficient approach?

**Communication:**
- [ ] Can I explain my approach clearly?
- [ ] Can I walk through the pointer movements?
- [ ] Can I justify my choice of technique?

---

## 🔧 Common Mistakes & Debug Tips

### Typical Mistakes:
1. **Nil pointer dereference** - Always check `node != nil` before accessing `node.Val` or `node.Next`
2. **Lost references** - Store `next` pointer before modifying `current.Next`
3. **Infinite loops** - Make sure pointers are advancing correctly
4. **Edge cases** - Test with empty lists, single nodes, and boundary conditions

### Debug Strategy:
```go
// Add debug prints to trace pointer movement
func debugList(head *ListNode, label string) {
    fmt.Printf("%s: ", label)
    current := head
    for current != nil {
        fmt.Printf("%d -> ", current.Val)
        current = current.Next
    }
    fmt.Println("nil")
}

// Use this in your solutions for debugging
debugList(head, "Original")
// ... your operations ...
debugList(result, "Result")
```

---

## 🚀 Tomorrow Preview: Advanced Linked Lists

**Day 4 Topics:**
- Cycle detection with Floyd's algorithm
- Finding nth node from end
- Palindrome linked list checking
- More complex two-pointer problems

**Preparation:**
- Make sure you're comfortable with basic pointer manipulation
- Practice drawing linked list operations
- Review recursion concepts (we'll use recursive solutions)

**Key Insight for Today:**
Linked lists are all about **pointer manipulation**. Master the patterns today (dummy head, two pointers, three pointers for reversal), and tomorrow's advanced problems will be much easier!

Remember: **Draw first, code second!** Visual understanding makes pointer manipulation much clearer. 🎯
