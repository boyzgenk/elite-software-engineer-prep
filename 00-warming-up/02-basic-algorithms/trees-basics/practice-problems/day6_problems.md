# Trees Practice Problems - Day 6

## Learning Objectives
By the end of this session, you should be able to:
- Implement and understand binary tree traversals (DFS and BFS)
- Solve basic Binary Search Tree (BST) problems
- Handle tree depth and path-related questions
- Apply recursion effectively for tree problems

## Problem Templates

### Problem 1: Maximum Depth of Binary Tree (Easy)
**LeetCode #104**

**Problem Statement:**
Given the root of a binary tree, return its maximum depth. Maximum depth is the number of nodes along the longest path from the root node down to the farthest leaf node.

**Examples:**
```
Input: root = [3,9,20,null,null,15,7]
Output: 3

Input: root = [1,null,2]
Output: 2
```

**Golang Template:**
```go
func maxDepth(root *TreeNode) int {
    // TODO: Implement using recursion
    // Time: O(n), Space: O(h) where h is height
    
}
```

**Approach:**
1. Base case: if root is nil, return 0
2. Recursively find depth of left and right subtrees
3. Return 1 + max(leftDepth, rightDepth)

---

### Problem 2: Binary Tree Inorder Traversal (Easy)
**LeetCode #94**

**Problem Statement:**
Given the root of a binary tree, return the inorder traversal of its nodes' values.

**Examples:**
```
Input: root = [1,null,2,3]
Output: [1,3,2]

Input: root = []
Output: []

Input: root = [1]
Output: [1]
```

**Golang Template:**
```go
func inorderTraversal(root *TreeNode) []int {
    // TODO: Implement using recursion (Left, Root, Right)
    // Time: O(n), Space: O(h)
    
}
```

**Approach:**
1. Base case: if root is nil, return empty slice
2. Recursively traverse left subtree
3. Process current node (add to result)
4. Recursively traverse right subtree

---

### Problem 3: Symmetric Tree (Easy)
**LeetCode #101**

**Problem Statement:**
Given the root of a binary tree, check whether it is a mirror of itself (i.e., symmetric around its center).

**Examples:**
```
Input: root = [1,2,2,3,4,4,3]
Output: true

Input: root = [1,2,2,null,3,null,3]
Output: false
```

**Golang Template:**
```go
func isSymmetric(root *TreeNode) bool {
    // TODO: Implement using helper function
    // Time: O(n), Space: O(h)
    
}

func isMirror(left, right *TreeNode) bool {
    // TODO: Check if left and right subtrees are mirrors
    
}
```

**Approach:**
1. Use helper function to check if two trees are mirrors
2. Two trees are mirrors if:
   - Both are nil
   - Both have same root value
   - Left subtree of first mirrors right subtree of second
   - Right subtree of first mirrors left subtree of second

---

### Problem 4: Binary Tree Level Order Traversal (Medium)
**LeetCode #102**

**Problem Statement:**
Given the root of a binary tree, return the level order traversal of its nodes' values (i.e., from left to right, level by level).

**Examples:**
```
Input: root = [3,9,20,null,null,15,7]
Output: [[3],[9,20],[15,7]]

Input: root = [1]
Output: [[1]]

Input: root = []
Output: []
```

**Golang Template:**
```go
func levelOrder(root *TreeNode) [][]int {
    // TODO: Implement using BFS with queue
    // Time: O(n), Space: O(w) where w is max width
    
}
```

**Approach:**
1. Use queue to store nodes level by level
2. Process all nodes at current level before moving to next
3. Track level size to group nodes correctly

---

### Problem 5: Validate Binary Search Tree (Medium)
**LeetCode #98**

**Problem Statement:**
Given the root of a binary tree, determine if it is a valid binary search tree (BST).

**Examples:**
```
Input: root = [2,1,3]
Output: true

Input: root = [5,1,4,null,null,3,6]
Output: false
Explanation: The root node's value is 5 but its right child's value is 4.
```

**Golang Template:**
```go
func isValidBST(root *TreeNode) bool {
    // TODO: Implement using bounds checking
    // Time: O(n), Space: O(h)
    
}

func validate(node *TreeNode, min, max int) bool {
    // TODO: Check if node value is within bounds
    
}
```

**Approach:**
1. Use helper function with min/max bounds
2. Each node must be within its allowed range
3. Left subtree: upper bound becomes current node value
4. Right subtree: lower bound becomes current node value

---

## Advanced Problems

### Problem 6: Binary Tree Paths (Medium)
**LeetCode #257**

**Problem Statement:**
Given the root of a binary tree, return all root-to-leaf paths in any order.

**Examples:**
```
Input: root = [1,2,3,null,5]
Output: ["1->2->5","1->3"]

Input: root = [1]
Output: ["1"]
```

**Golang Template:**
```go
func binaryTreePaths(root *TreeNode) []string {
    // TODO: Implement using DFS with path tracking
    // Time: O(n), Space: O(n)
    
}
```

**Approach:**
1. Use DFS to explore all paths
2. Keep track of current path as string
3. When reaching leaf, add path to result
4. Backtrack by removing current node when returning

---

### Problem 7: Path Sum (Easy)
**LeetCode #112**

**Problem Statement:**
Given the root of a binary tree and an integer targetSum, return true if the tree has a root-to-leaf path such that adding up all the values along the path equals targetSum.

**Examples:**
```
Input: root = [5,4,8,11,null,13,4,7,2,null,null,null,1], targetSum = 22
Output: true

Input: root = [1,2,3], targetSum = 5
Output: false
```

**Golang Template:**
```go
func hasPathSum(root *TreeNode, targetSum int) bool {
    // TODO: Implement using recursion
    // Time: O(n), Space: O(h)
    
}
```

**Approach:**
1. Base case: if root is nil, return false
2. If leaf node, check if value equals remaining sum
3. Recursively check left and right subtrees with updated sum

---

### Problem 8: Lowest Common Ancestor of BST (Easy)
**LeetCode #235**

**Problem Statement:**
Given a binary search tree (BST), find the lowest common ancestor (LCA) of two given nodes in the BST.

**Examples:**
```
Input: root = [6,2,8,0,4,7,9,null,null,3,5], p = 2, q = 8
Output: 6
Explanation: The LCA of nodes 2 and 8 is 6.

Input: root = [6,2,8,0,4,7,9,null,null,3,5], p = 2, q = 4
Output: 2
Explanation: The LCA of nodes 2 and 4 is 2.
```

**Golang Template:**
```go
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    // TODO: Implement using BST property
    // Time: O(h), Space: O(1) iterative / O(h) recursive
    
}
```

**Approach:**
1. Use BST property: LCA is where paths to p and q diverge
2. If both p and q are smaller than root, LCA is in left subtree
3. If both p and q are larger than root, LCA is in right subtree
4. Otherwise, current root is the LCA

---

## Daily Study Plan

### Morning Session (45 minutes)
1. **Review Theory (10 min)**: Tree terminology, traversals, BST properties
2. **Solve Problems 1-3 (30 min)**: Focus on basic recursion patterns
3. **Code Review (5 min)**: Analyze time/space complexity

### Evening Session (45 minutes)
1. **Solve Problems 4-5 (25 min)**: Practice BFS and BST validation
2. **Attempt Advanced Problems 6-8 (15 min)**: Challenge problems
3. **Reflection (5 min)**: Note common patterns and techniques

## Success Checklist
- [ ] Can implement all three DFS traversals (inorder, preorder, postorder)
- [ ] Understand BFS level-order traversal with queue
- [ ] Comfortable with tree recursion patterns
- [ ] Can validate BST correctly using bounds
- [ ] Understand path-tracking techniques
- [ ] Can handle edge cases (empty tree, single node)
- [ ] Know when to use DFS vs BFS

## Key Patterns to Remember

### 1. **Basic Tree Recursion**
```go
func treeFunction(root *TreeNode) returnType {
    // Base case
    if root == nil {
        return baseValue
    }
    
    // Process current node
    // Recurse on children
    left := treeFunction(root.Left)
    right := treeFunction(root.Right)
    
    // Combine results
    return combine(root.Val, left, right)
}
```

### 2. **Path Tracking DFS**
```go
func pathDFS(node *TreeNode, path []int, target condition) {
    if node == nil {
        return
    }
    
    path = append(path, node.Val)
    
    if isLeaf(node) && condition {
        // Process complete path
    }
    
    pathDFS(node.Left, path, updatedCondition)
    pathDFS(node.Right, path, updatedCondition)
    // Note: path automatically backtracks due to slice semantics
}
```

### 3. **Level Order BFS**
```go
func levelOrder(root *TreeNode) [][]int {
    if root == nil {
        return [][]int{}
    }
    
    result := [][]int{}
    queue := []*TreeNode{root}
    
    for len(queue) > 0 {
        levelSize := len(queue)
        currentLevel := []int{}
        
        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:]
            
            currentLevel = append(currentLevel, node.Val)
            
            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        
        result = append(result, currentLevel)
    }
    
    return result
}
```

### 4. **BST Bounds Validation**
```go
func isValidBST(root *TreeNode) bool {
    return validate(root, math.MinInt64, math.MaxInt64)
}

func validate(node *TreeNode, min, max int) bool {
    if node == nil {
        return true
    }
    
    if node.Val <= min || node.Val >= max {
        return false
    }
    
    return validate(node.Left, min, node.Val) && 
           validate(node.Right, node.Val, max)
}
```

## Common Mistakes to Avoid
1. **Forgetting base cases** - Always handle nil nodes
2. **Incorrect BST validation** - Use bounds, not just parent-child comparison
3. **Path tracking errors** - Remember backtracking in DFS
4. **BFS level confusion** - Process entire level before moving to next
5. **Leaf node detection** - Check both left and right children are nil

## Time/Space Complexity Quick Reference
- **Tree traversal**: O(n) time, O(h) space (recursion stack)
- **BFS level order**: O(n) time, O(w) space (queue width)
- **BST search**: O(h) time, O(1) space iterative
- **Path problems**: O(n) time, O(h) space + path length
- Where n = number of nodes, h = height, w = max width