# 🌳 Trees Fundamentals  
## Week 2, Days 4-7 | Hierarchical Data Structures

### 🎯 Learning Objectives
By the end of this section, you will:
- [ ] Understand tree terminology and hierarchical relationships
- [ ] Implement binary trees and binary search trees from scratch
- [ ] Master tree traversal algorithms (in-order, pre-order, post-order, level-order)
- [ ] Solve 10-12 tree problems using recursive and iterative approaches
- [ ] Analyze tree operations and their time complexities

---

## 📚 Tree Fundamentals

### 🔸 What is a Tree?

**Definition:**
- Hierarchical data structure consisting of nodes connected by edges
- Each node has at most one parent (except root) and zero or more children
- No cycles exist in a tree structure
- Represents hierarchical relationships (family trees, file systems, decision trees)

**Key Terminology:**
- **Root:** Top node with no parent
- **Leaf:** Node with no children
- **Parent:** Node with child nodes
- **Child:** Node with a parent
- **Sibling:** Nodes with same parent
- **Depth:** Distance from root to node
- **Height:** Maximum depth in tree
- **Level:** All nodes at same depth

**Real-world Applications:**
- File system hierarchy
- HTML DOM structure  
- Database indexing (B-trees)
- Expression parsing
- Decision making algorithms
- Organizational charts

### 🔸 Binary Tree Properties

**Binary Tree:**
- Each node has at most 2 children (left and right)
- Children are ordered (left vs right matters)
- Can be empty (null tree)

**Types of Binary Trees:**
1. **Full Binary Tree:** Every node has 0 or 2 children
2. **Complete Binary Tree:** All levels filled except possibly last, filled left to right
3. **Perfect Binary Tree:** All internal nodes have 2 children, all leaves at same level
4. **Balanced Binary Tree:** Height difference between left and right subtrees ≤ 1

**Binary Search Tree (BST) Properties:**
- Left subtree contains values less than node
- Right subtree contains values greater than node
- No duplicate values (typically)
- In-order traversal gives sorted sequence

---

## 🛠️ Implementation in Golang

### Basic Binary Tree Node

```go
package main

import (
    "fmt"
    "strconv"
    "strings"
)

// TreeNode represents a node in binary tree
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

// NewTreeNode creates a new tree node
func NewTreeNode(val int) *TreeNode {
    return &TreeNode{
        Val:   val,
        Left:  nil,
        Right: nil,
    }
}

// BinaryTree represents a binary tree structure
type BinaryTree struct {
    Root *TreeNode
}

// NewBinaryTree creates a new binary tree
func NewBinaryTree() *BinaryTree {
    return &BinaryTree{Root: nil}
}

// Insert adds a value to the tree (level-order insertion)
func (bt *BinaryTree) Insert(val int) {
    newNode := NewTreeNode(val)
    
    if bt.Root == nil {
        bt.Root = newNode
        return
    }
    
    // Level-order insertion using queue
    queue := []*TreeNode{bt.Root}
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        if current.Left == nil {
            current.Left = newNode
            return
        } else if current.Right == nil {
            current.Right = newNode
            return
        } else {
            queue = append(queue, current.Left)
            queue = append(queue, current.Right)
        }
    }
}

// Search finds a value in the tree - O(n) for general binary tree
func (bt *BinaryTree) Search(val int) bool {
    return bt.searchHelper(bt.Root, val)
}

func (bt *BinaryTree) searchHelper(node *TreeNode, val int) bool {
    if node == nil {
        return false
    }
    
    if node.Val == val {
        return true
    }
    
    return bt.searchHelper(node.Left, val) || bt.searchHelper(node.Right, val)
}

// Height returns the height of the tree
func (bt *BinaryTree) Height() int {
    return bt.heightHelper(bt.Root)
}

func (bt *BinaryTree) heightHelper(node *TreeNode) int {
    if node == nil {
        return -1 // Height of empty tree is -1
    }
    
    leftHeight := bt.heightHelper(node.Left)
    rightHeight := bt.heightHelper(node.Right)
    
    return max(leftHeight, rightHeight) + 1
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// Size returns number of nodes in the tree
func (bt *BinaryTree) Size() int {
    return bt.sizeHelper(bt.Root)
}

func (bt *BinaryTree) sizeHelper(node *TreeNode) int {
    if node == nil {
        return 0
    }
    
    return 1 + bt.sizeHelper(node.Left) + bt.sizeHelper(node.Right)
}
```

### Binary Search Tree Implementation

```go
// BST represents a Binary Search Tree
type BST struct {
    Root *TreeNode
}

// NewBST creates a new binary search tree
func NewBST() *BST {
    return &BST{Root: nil}
}

// Insert adds value maintaining BST property - O(log n) average, O(n) worst
func (bst *BST) Insert(val int) {
    bst.Root = bst.insertHelper(bst.Root, val)
}

func (bst *BST) insertHelper(node *TreeNode, val int) *TreeNode {
    // Base case: empty spot found
    if node == nil {
        return NewTreeNode(val)
    }
    
    // Recursive case: navigate to correct position
    if val < node.Val {
        node.Left = bst.insertHelper(node.Left, val)
    } else if val > node.Val {
        node.Right = bst.insertHelper(node.Right, val)
    }
    // Equal values are ignored (no duplicates)
    
    return node
}

// Search finds value in BST - O(log n) average, O(n) worst
func (bst *BST) Search(val int) bool {
    return bst.searchHelper(bst.Root, val)
}

func (bst *BST) searchHelper(node *TreeNode, val int) bool {
    if node == nil {
        return false
    }
    
    if val == node.Val {
        return true
    } else if val < node.Val {
        return bst.searchHelper(node.Left, val)
    } else {
        return bst.searchHelper(node.Right, val)
    }
}

// Delete removes value from BST - O(log n) average
func (bst *BST) Delete(val int) {
    bst.Root = bst.deleteHelper(bst.Root, val)
}

func (bst *BST) deleteHelper(node *TreeNode, val int) *TreeNode {
    if node == nil {
        return nil
    }
    
    if val < node.Val {
        node.Left = bst.deleteHelper(node.Left, val)
    } else if val > node.Val {
        node.Right = bst.deleteHelper(node.Right, val)
    } else {
        // Node to delete found
        
        // Case 1: No children (leaf node)
        if node.Left == nil && node.Right == nil {
            return nil
        }
        
        // Case 2: One child
        if node.Left == nil {
            return node.Right
        }
        if node.Right == nil {
            return node.Left
        }
        
        // Case 3: Two children
        // Find inorder successor (smallest in right subtree)
        successor := bst.findMin(node.Right)
        node.Val = successor.Val
        node.Right = bst.deleteHelper(node.Right, successor.Val)
    }
    
    return node
}

// findMin finds minimum value in subtree
func (bst *BST) findMin(node *TreeNode) *TreeNode {
    for node.Left != nil {
        node = node.Left
    }
    return node
}

// FindMin returns minimum value in BST
func (bst *BST) FindMin() (int, bool) {
    if bst.Root == nil {
        return 0, false
    }
    node := bst.findMin(bst.Root)
    return node.Val, true
}

// FindMax returns maximum value in BST
func (bst *BST) FindMax() (int, bool) {
    if bst.Root == nil {
        return 0, false
    }
    
    current := bst.Root
    for current.Right != nil {
        current = current.Right
    }
    
    return current.Val, true
}

// IsValidBST checks if tree satisfies BST property
func (bst *BST) IsValidBST() bool {
    return bst.isValidBSTHelper(bst.Root, nil, nil)
}

func (bst *BST) isValidBSTHelper(node *TreeNode, min, max *int) bool {
    if node == nil {
        return true
    }
    
    // Check current node's value against bounds
    if (min != nil && node.Val <= *min) || (max != nil && node.Val >= *max) {
        return false
    }
    
    // Recursively validate left and right subtrees with updated bounds
    return bst.isValidBSTHelper(node.Left, min, &node.Val) &&
           bst.isValidBSTHelper(node.Right, &node.Val, max)
}
```

---

## 🎯 Tree Traversal Algorithms

### Depth-First Traversals

```go
// ===== RECURSIVE TRAVERSALS =====

// InorderTraversal: Left -> Root -> Right (gives sorted order for BST)
func (bt *BinaryTree) InorderTraversal() []int {
    result := []int{}
    bt.inorderHelper(bt.Root, &result)
    return result
}

func (bt *BinaryTree) inorderHelper(node *TreeNode, result *[]int) {
    if node == nil {
        return
    }
    bt.inorderHelper(node.Left, result)   // Left
    *result = append(*result, node.Val)   // Root
    bt.inorderHelper(node.Right, result)  // Right
}

// PreorderTraversal: Root -> Left -> Right (useful for copying tree)
func (bt *BinaryTree) PreorderTraversal() []int {
    result := []int{}
    bt.preorderHelper(bt.Root, &result)
    return result
}

func (bt *BinaryTree) preorderHelper(node *TreeNode, result *[]int) {
    if node == nil {
        return
    }
    *result = append(*result, node.Val)   // Root
    bt.preorderHelper(node.Left, result)  // Left
    bt.preorderHelper(node.Right, result) // Right
}

// PostorderTraversal: Left -> Right -> Root (useful for deleting tree)
func (bt *BinaryTree) PostorderTraversal() []int {
    result := []int{}
    bt.postorderHelper(bt.Root, &result)
    return result
}

func (bt *BinaryTree) postorderHelper(node *TreeNode, result *[]int) {
    if node == nil {
        return
    }
    bt.postorderHelper(node.Left, result)  // Left
    bt.postorderHelper(node.Right, result) // Right
    *result = append(*result, node.Val)    // Root
}

// ===== ITERATIVE TRAVERSALS =====

// InorderIterative using stack
func (bt *BinaryTree) InorderIterative() []int {
    result := []int{}
    stack := []*TreeNode{}
    current := bt.Root
    
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

// PreorderIterative using stack
func (bt *BinaryTree) PreorderIterative() []int {
    if bt.Root == nil {
        return []int{}
    }
    
    result := []int{}
    stack := []*TreeNode{bt.Root}
    
    for len(stack) > 0 {
        // Pop node from stack
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        
        // Process current node
        result = append(result, current.Val)
        
        // Push right first, then left (stack is LIFO)
        if current.Right != nil {
            stack = append(stack, current.Right)
        }
        if current.Left != nil {
            stack = append(stack, current.Left)
        }
    }
    
    return result
}

// PostorderIterative using two stacks
func (bt *BinaryTree) PostorderIterative() []int {
    if bt.Root == nil {
        return []int{}
    }
    
    result := []int{}
    stack1 := []*TreeNode{bt.Root}
    stack2 := []*TreeNode{}
    
    // First stack for traversal, second stack for result order
    for len(stack1) > 0 {
        current := stack1[len(stack1)-1]
        stack1 = stack1[:len(stack1)-1]
        stack2 = append(stack2, current)
        
        if current.Left != nil {
            stack1 = append(stack1, current.Left)
        }
        if current.Right != nil {
            stack1 = append(stack1, current.Right)
        }
    }
    
    // Pop from second stack to get postorder
    for len(stack2) > 0 {
        current := stack2[len(stack2)-1]
        stack2 = stack2[:len(stack2)-1]
        result = append(result, current.Val)
    }
    
    return result
}
```

### Breadth-First Traversal (Level Order)

```go
// LevelOrderTraversal: Process nodes level by level
func (bt *BinaryTree) LevelOrderTraversal() []int {
    if bt.Root == nil {
        return []int{}
    }
    
    result := []int{}
    queue := []*TreeNode{bt.Root}
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        result = append(result, current.Val)
        
        if current.Left != nil {
            queue = append(queue, current.Left)
        }
        if current.Right != nil {
            queue = append(queue, current.Right)
        }
    }
    
    return result
}

// LevelOrderByLevels: Return each level as separate slice
func (bt *BinaryTree) LevelOrderByLevels() [][]int {
    if bt.Root == nil {
        return [][]int{}
    }
    
    result := [][]int{}
    queue := []*TreeNode{bt.Root}
    
    for len(queue) > 0 {
        levelSize := len(queue)
        level := []int{}
        
        for i := 0; i < levelSize; i++ {
            current := queue[0]
            queue = queue[1:]
            
            level = append(level, current.Val)
            
            if current.Left != nil {
                queue = append(queue, current.Left)
            }
            if current.Right != nil {
                queue = append(queue, current.Right)
            }
        }
        
        result = append(result, level)
    }
    
    return result
}

// ZigzagLevelOrder: Alternate left-to-right and right-to-left
func (bt *BinaryTree) ZigzagLevelOrder() [][]int {
    if bt.Root == nil {
        return [][]int{}
    }
    
    result := [][]int{}
    queue := []*TreeNode{bt.Root}
    leftToRight := true
    
    for len(queue) > 0 {
        levelSize := len(queue)
        level := make([]int, levelSize)
        
        for i := 0; i < levelSize; i++ {
            current := queue[0]
            queue = queue[1:]
            
            // Fill level array differently based on direction
            if leftToRight {
                level[i] = current.Val
            } else {
                level[levelSize-1-i] = current.Val
            }
            
            if current.Left != nil {
                queue = append(queue, current.Left)
            }
            if current.Right != nil {
                queue = append(queue, current.Right)
            }
        }
        
        result = append(result, level)
        leftToRight = !leftToRight
    }
    
    return result
}
```

---

## 🎯 Essential Problem-Solving Patterns

### Pattern 1: Tree Properties and Validation

```go
// Check if tree is symmetric (mirror of itself)
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

// Check if tree is balanced (height difference ≤ 1)
func isBalanced(root *TreeNode) bool {
    _, balanced := checkBalance(root)
    return balanced
}

func checkBalance(node *TreeNode) (int, bool) {
    if node == nil {
        return 0, true
    }
    
    leftHeight, leftBalanced := checkBalance(node.Left)
    if !leftBalanced {
        return 0, false
    }
    
    rightHeight, rightBalanced := checkBalance(node.Right)
    if !rightBalanced {
        return 0, false
    }
    
    heightDiff := abs(leftHeight - rightHeight)
    balanced := heightDiff <= 1
    height := max(leftHeight, rightHeight) + 1
    
    return height, balanced
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

// Find diameter of tree (longest path between any two nodes)
func diameterOfBinaryTree(root *TreeNode) int {
    diameter := 0
    calculateHeight(root, &diameter)
    return diameter
}

func calculateHeight(node *TreeNode, diameter *int) int {
    if node == nil {
        return 0
    }
    
    leftHeight := calculateHeight(node.Left, diameter)
    rightHeight := calculateHeight(node.Right, diameter)
    
    // Update diameter if path through current node is longer
    currentDiameter := leftHeight + rightHeight
    if currentDiameter > *diameter {
        *diameter = currentDiameter
    }
    
    return max(leftHeight, rightHeight) + 1
}
```

### Pattern 2: Path Problems

```go
// Find all root-to-leaf paths
func binaryTreePaths(root *TreeNode) []string {
    if root == nil {
        return []string{}
    }
    
    paths := []string{}
    findPaths(root, "", &paths)
    return paths
}

func findPaths(node *TreeNode, currentPath string, paths *[]string) {
    if currentPath == "" {
        currentPath = strconv.Itoa(node.Val)
    } else {
        currentPath += "->" + strconv.Itoa(node.Val)
    }
    
    // If leaf node, add path to result
    if node.Left == nil && node.Right == nil {
        *paths = append(*paths, currentPath)
        return
    }
    
    // Continue path exploration
    if node.Left != nil {
        findPaths(node.Left, currentPath, paths)
    }
    if node.Right != nil {
        findPaths(node.Right, currentPath, paths)
    }
}

// Check if path exists with given sum
func hasPathSum(root *TreeNode, targetSum int) bool {
    if root == nil {
        return false
    }
    
    // If leaf node, check if remaining sum matches node value
    if root.Left == nil && root.Right == nil {
        return root.Val == targetSum
    }
    
    // Continue with remaining sum
    remainingSum := targetSum - root.Val
    return hasPathSum(root.Left, remainingSum) || hasPathSum(root.Right, remainingSum)
}

// Find maximum path sum in tree (path can start/end at any node)
func maxPathSum(root *TreeNode) int {
    maxSum := root.Val // Initialize with root value
    calculateMaxPath(root, &maxSum)
    return maxSum
}

func calculateMaxPath(node *TreeNode, maxSum *int) int {
    if node == nil {
        return 0
    }
    
    // Calculate max path from left and right (ignore negative sums)
    leftMax := max(0, calculateMaxPath(node.Left, maxSum))
    rightMax := max(0, calculateMaxPath(node.Right, maxSum))
    
    // Check if path through current node gives better sum
    currentMax := node.Val + leftMax + rightMax
    if currentMax > *maxSum {
        *maxSum = currentMax
    }
    
    // Return max path ending at current node
    return node.Val + max(leftMax, rightMax)
}
```

### Pattern 3: Tree Construction

```go
// Build tree from inorder and preorder traversals
func buildTree(preorder []int, inorder []int) *TreeNode {
    if len(preorder) == 0 || len(inorder) == 0 {
        return nil
    }
    
    // First element in preorder is always root
    root := NewTreeNode(preorder[0])
    
    // Find root position in inorder
    var rootIndex int
    for i, val := range inorder {
        if val == preorder[0] {
            rootIndex = i
            break
        }
    }
    
    // Recursively build left and right subtrees
    root.Left = buildTree(preorder[1:rootIndex+1], inorder[:rootIndex])
    root.Right = buildTree(preorder[rootIndex+1:], inorder[rootIndex+1:])
    
    return root
}

// Convert sorted array to balanced BST
func sortedArrayToBST(nums []int) *TreeNode {
    if len(nums) == 0 {
        return nil
    }
    
    return buildBalancedBST(nums, 0, len(nums)-1)
}

func buildBalancedBST(nums []int, left, right int) *TreeNode {
    if left > right {
        return nil
    }
    
    mid := left + (right-left)/2
    root := NewTreeNode(nums[mid])
    
    root.Left = buildBalancedBST(nums, left, mid-1)
    root.Right = buildBalancedBST(nums, mid+1, right)
    
    return root
}
```

---

## 💻 Practice Problems

### Easy Level (Day 4-5)

#### Problem 1: Maximum Depth of Binary Tree (LeetCode #104)
```go
func maxDepth(root *TreeNode) int {
    // Calculate height of tree
    return 0
}
```

#### Problem 2: Same Tree (LeetCode #100)
```go
func isSameTree(p *TreeNode, q *TreeNode) bool {
    // Check if two trees are identical
    return false
}
```

#### Problem 3: Invert Binary Tree (LeetCode #226)
```go
func invertTree(root *TreeNode) *TreeNode {
    // Swap left and right children recursively
    return nil
}
```

### Medium Level (Day 6-7)

#### Problem 4: Validate Binary Search Tree (LeetCode #98)
```go
func isValidBST(root *TreeNode) bool {
    // Check BST property with bounds
    return false
}
```

#### Problem 5: Binary Tree Level Order Traversal (LeetCode #102)
```go
func levelOrder(root *TreeNode) [][]int {
    // Return level-by-level traversal
    return nil
}
```

#### Problem 6: Construct Binary Tree from Preorder and Inorder (LeetCode #105)
```go
func buildTree(preorder []int, inorder []int) *TreeNode {
    // Build tree from traversal arrays
    return nil
}
```

---

## 📊 Time & Space Complexity Analysis

### Binary Tree Operations
| Operation | Average | Worst Case | Space |
|-----------|---------|------------|-------|
| Search | O(n) | O(n) | O(h) |
| Insert | O(n) | O(n) | O(h) |
| Delete | O(n) | O(n) | O(h) |
| Traversal | O(n) | O(n) | O(h) |

### Binary Search Tree Operations
| Operation | Average | Worst Case | Space |
|-----------|---------|------------|-------|
| Search | O(log n) | O(n) | O(h) |
| Insert | O(log n) | O(n) | O(h) |
| Delete | O(log n) | O(n) | O(h) |
| Find Min/Max | O(log n) | O(n) | O(h) |

Where h = height of tree (log n for balanced, n for skewed)

### Traversal Complexities
| Traversal | Time | Space (Recursive) | Space (Iterative) |
|-----------|------|-------------------|-------------------|
| Inorder | O(n) | O(h) | O(h) |
| Preorder | O(n) | O(h) | O(h) |
| Postorder | O(n) | O(h) | O(h) |
| Level-order | O(n) | O(w) | O(w) |

Where w = maximum width of tree

---

## 🚀 Next Steps

### Day 4 Goals (Basic Trees)
- [ ] Implement binary tree with basic operations
- [ ] Master recursive traversals (in/pre/post-order)
- [ ] Solve tree property problems (height, size, symmetric)
- [ ] Practice tree visualization and drawing

### Day 5 Goals (Tree Traversals)
- [ ] Implement iterative traversals using stacks
- [ ] Master level-order traversal with queues
- [ ] Solve path-finding problems
- [ ] Practice zigzag and vertical traversals

### Day 6 Goals (Binary Search Trees)
- [ ] Implement BST with insert/delete/search
- [ ] Understand BST validation and properties
- [ ] Solve BST-specific problems
- [ ] Practice tree construction from traversals

### Day 7 Goals (Advanced Tree Problems)
- [ ] Solve diameter and path sum problems
- [ ] Practice tree transformation problems
- [ ] Master tree construction patterns
- [ ] Review and consolidate all tree concepts

### Preparation for Week 3 (Basic Algorithms)
- Review recursion and divide-and-conquer concepts
- Understand sorting and searching principles
- Practice analyzing algorithm efficiency
- Think about problem decomposition strategies

**Trees are the foundation for many advanced algorithms! Master these recursive patterns and graph problems will become much more manageable. Strategic thinking starts with hierarchical data! 🌳**

---

## 🔧 Tree Debugging Tips

### Common Tree Mistakes:
1. **Null pointer access** - Always check `node != nil`
2. **Infinite recursion** - Ensure base cases are correct
3. **Wrong traversal order** - Draw the tree and trace through
4. **BST property violation** - Validate bounds correctly

### Debugging Techniques:
```go
// Add debug prints to trace recursion
func debugTraversal(node *TreeNode, depth int) {
    if node == nil {
        fmt.Printf("%snil\n", strings.Repeat("  ", depth))
        return
    }
    
    fmt.Printf("%s%d\n", strings.Repeat("  ", depth), node.Val)
    debugTraversal(node.Left, depth+1)
    debugTraversal(node.Right, depth+1)
}

// Visualize tree structure
func (bt *BinaryTree) PrintTree() {
    bt.printHelper(bt.Root, "", true)
}

func (bt *BinaryTree) printHelper(node *TreeNode, prefix string, isLast bool) {
    if node == nil {
        return
    }
    
    fmt.Print(prefix)
    if isLast {
        fmt.Print("└── ")
        prefix += "    "
    } else {
        fmt.Print("├── ")
        prefix += "│   "
    }
    fmt.Println(node.Val)
    
    if node.Left != nil || node.Right != nil {
        if node.Right != nil {
            bt.printHelper(node.Right, prefix, node.Left == nil)
        }
        if node.Left != nil {
            bt.printHelper(node.Left, prefix, true)
        }
    }
}
```

**Remember: Trees are recursive by nature - embrace the recursive thinking! Each subtree is a smaller version of the same problem. That's the key to tree mastery!** 💪