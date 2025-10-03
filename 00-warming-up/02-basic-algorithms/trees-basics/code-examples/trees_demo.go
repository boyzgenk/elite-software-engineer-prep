package main

import (
	"fmt"
	"math"
)

// TreeNode represents a binary tree node
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ============= BINARY TREE CREATION =============

// CreateSampleTree creates a sample binary tree for testing
func CreateSampleTree() *TreeNode {
	/*
	     3
	   /   \
	  9     20
	       /  \
	      15   7
	*/
	root := &TreeNode{Val: 3}
	root.Left = &TreeNode{Val: 9}
	root.Right = &TreeNode{Val: 20}
	root.Right.Left = &TreeNode{Val: 15}
	root.Right.Right = &TreeNode{Val: 7}
	return root
}

// CreateBST creates a sample Binary Search Tree
func CreateBST() *TreeNode {
	/*
		     5
		   /   \
		  3     8
		 / \   / \
		2   4 7   9
	*/
	root := &TreeNode{Val: 5}
	root.Left = &TreeNode{Val: 3}
	root.Right = &TreeNode{Val: 8}
	root.Left.Left = &TreeNode{Val: 2}
	root.Left.Right = &TreeNode{Val: 4}
	root.Right.Left = &TreeNode{Val: 7}
	root.Right.Right = &TreeNode{Val: 9}
	return root
}

// ============= TREE TRAVERSALS =============

// InorderTraversal (Left, Root, Right) - for BST gives sorted order
func InorderTraversal(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	result := []int{}
	result = append(result, InorderTraversal(root.Left)...)
	result = append(result, root.Val)
	result = append(result, InorderTraversal(root.Right)...)

	return result
}

// PreorderTraversal (Root, Left, Right) - good for tree reconstruction
func PreorderTraversal(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	result := []int{}
	result = append(result, root.Val)
	result = append(result, PreorderTraversal(root.Left)...)
	result = append(result, PreorderTraversal(root.Right)...)

	return result
}

// PostorderTraversal (Left, Right, Root) - good for deletion
func PostorderTraversal(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	result := []int{}
	result = append(result, PostorderTraversal(root.Left)...)
	result = append(result, PostorderTraversal(root.Right)...)
	result = append(result, root.Val)

	return result
}

// LevelOrderTraversal (BFS) - level by level
func LevelOrderTraversal(root *TreeNode) [][]int {
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

// ============= BASIC TREE OPERATIONS =============

// MaxDepth calculates the maximum depth of a binary tree
func MaxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	leftDepth := MaxDepth(root.Left)
	rightDepth := MaxDepth(root.Right)

	return 1 + max(leftDepth, rightDepth)
}

// MinDepth calculates the minimum depth to a leaf node
func MinDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	if root.Left == nil && root.Right == nil {
		return 1
	}

	if root.Left == nil {
		return 1 + MinDepth(root.Right)
	}

	if root.Right == nil {
		return 1 + MinDepth(root.Left)
	}

	return 1 + min(MinDepth(root.Left), MinDepth(root.Right))
}

// IsSymmetric checks if a tree is symmetric (mirror of itself)
func IsSymmetric(root *TreeNode) bool {
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

// ============= BINARY SEARCH TREE OPERATIONS =============

// SearchBST searches for a value in BST
func SearchBST(root *TreeNode, val int) *TreeNode {
	if root == nil || root.Val == val {
		return root
	}

	if val < root.Val {
		return SearchBST(root.Left, val)
	}
	return SearchBST(root.Right, val)
}

// InsertIntoBST inserts a value into BST
func InsertIntoBST(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return &TreeNode{Val: val}
	}

	if val < root.Val {
		root.Left = InsertIntoBST(root.Left, val)
	} else {
		root.Right = InsertIntoBST(root.Right, val)
	}

	return root
}

// IsValidBST checks if a tree is a valid BST
func IsValidBST(root *TreeNode) bool {
	return validateBST(root, math.MinInt64, math.MaxInt64)
}

func validateBST(node *TreeNode, min, max int) bool {
	if node == nil {
		return true
	}

	if node.Val <= min || node.Val >= max {
		return false
	}

	return validateBST(node.Left, min, node.Val) &&
		validateBST(node.Right, node.Val, max)
}

// ============= PATH AND SUM PROBLEMS =============

// HasPathSum checks if there's a path from root to leaf with given sum
func HasPathSum(root *TreeNode, targetSum int) bool {
	if root == nil {
		return false
	}

	if root.Left == nil && root.Right == nil {
		return root.Val == targetSum
	}

	remainingSum := targetSum - root.Val
	return HasPathSum(root.Left, remainingSum) || HasPathSum(root.Right, remainingSum)
}

// PathSum returns all paths from root to leaf with given sum
func PathSum(root *TreeNode, targetSum int) [][]int {
	result := [][]int{}
	if root == nil {
		return result
	}

	var dfs func(*TreeNode, int, []int)
	dfs = func(node *TreeNode, remaining int, path []int) {
		if node == nil {
			return
		}

		path = append(path, node.Val)
		remaining -= node.Val

		if node.Left == nil && node.Right == nil && remaining == 0 {
			// Make a copy of the path
			pathCopy := make([]int, len(path))
			copy(pathCopy, path)
			result = append(result, pathCopy)
		}

		dfs(node.Left, remaining, path)
		dfs(node.Right, remaining, path)
	}

	dfs(root, targetSum, []int{})
	return result
}

// ============= UTILITY FUNCTIONS =============

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PrintTree prints the tree in a readable format (level order)
func PrintTree(root *TreeNode) {
	if root == nil {
		fmt.Println("Empty tree")
		return
	}

	levels := LevelOrderTraversal(root)
	for i, level := range levels {
		fmt.Printf("Level %d: %v\n", i, level)
	}
}

// ============= DEMONSTRATION =============

func main() {
	fmt.Println("=== Binary Tree Examples ===")

	// Create sample trees
	tree := CreateSampleTree()
	bst := CreateBST()

	fmt.Println("\nSample Tree:")
	PrintTree(tree)

	fmt.Println("\nSample BST:")
	PrintTree(bst)

	// Demonstrate traversals
	fmt.Println("\n=== Tree Traversals ===")
	fmt.Printf("Inorder (BST): %v\n", InorderTraversal(bst))
	fmt.Printf("Preorder: %v\n", PreorderTraversal(bst))
	fmt.Printf("Postorder: %v\n", PostorderTraversal(bst))
	fmt.Printf("Level Order: %v\n", LevelOrderTraversal(bst))

	// Demonstrate basic operations
	fmt.Println("\n=== Basic Operations ===")
	fmt.Printf("Max Depth: %d\n", MaxDepth(bst))
	fmt.Printf("Min Depth: %d\n", MinDepth(bst))
	fmt.Printf("Is Symmetric: %t\n", IsSymmetric(bst))

	// Demonstrate BST operations
	fmt.Println("\n=== BST Operations ===")
	fmt.Printf("Search for 7: %v\n", SearchBST(bst, 7) != nil)
	fmt.Printf("Search for 10: %v\n", SearchBST(bst, 10) != nil)
	fmt.Printf("Is Valid BST: %t\n", IsValidBST(bst))

	// Insert new value
	bst = InsertIntoBST(bst, 6)
	fmt.Printf("After inserting 6: %v\n", InorderTraversal(bst))

	// Demonstrate path problems
	fmt.Println("\n=== Path Problems ===")
	fmt.Printf("Has path sum 12: %t\n", HasPathSum(bst, 12))
	fmt.Printf("Has path sum 20: %t\n", HasPathSum(bst, 20))

	pathSums := PathSum(bst, 20)
	fmt.Printf("All paths with sum 20: %v\n", pathSums)
}
