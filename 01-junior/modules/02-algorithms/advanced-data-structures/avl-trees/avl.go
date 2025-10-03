// Package avl implements an AVL (Adelson-Velsky and Landis) tree
// AVL trees are self-balancing binary search trees where the height difference
// between left and right subtrees is at most 1 for every node.
//
// Time Complexity:
//   - Search: O(log n)
//   - Insert: O(log n)
//   - Delete: O(log n)
//   - Height: O(log n) guaranteed
//
// Space Complexity: O(n) for storage, O(log n) for recursion stack
package avl

import "fmt"

// Node represents a single node in the AVL tree
type Node struct {
	Key    int   // The value stored in the node
	Height int   // Height of this node (leaf nodes have height 1)
	Left   *Node // Left child
	Right  *Node // Right child
}

// AVLTree represents the AVL tree structure
type AVLTree struct {
	Root *Node // Root node of the tree
	Size int   // Number of nodes in the tree
}

// NewAVLTree creates a new empty AVL tree
func NewAVLTree() *AVLTree {
	return &AVLTree{
		Root: nil,
		Size: 0,
	}
}

// height returns the height of a node (0 for nil nodes)
func height(node *Node) int {
	if node == nil {
		return 0
	}
	return node.Height
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// updateHeight updates the height of a node based on its children
func updateHeight(node *Node) {
	if node != nil {
		node.Height = 1 + max(height(node.Left), height(node.Right))
	}
}

// getBalance returns the balance factor of a node
// Balance factor = height(left) - height(right)
// AVL property: balance factor must be -1, 0, or 1
func getBalance(node *Node) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

// rightRotate performs a right rotation on the given node
//
//	   y               x
//	  / \             / \
//	 x   T3   -->    T1   y
//	/ \                  / \
//
// T1  T2               T2  T3
func rightRotate(y *Node) *Node {
	x := y.Left
	T2 := x.Right

	// Perform rotation
	x.Right = y
	y.Left = T2

	// Update heights
	updateHeight(y)
	updateHeight(x)

	// Return new root
	return x
}

// leftRotate performs a left rotation on the given node
//
//	 x                   y
//	/ \                 / \
//
// T1   y      -->     x   T3
//
//	 / \            / \
//	T2  T3         T1  T2
func leftRotate(x *Node) *Node {
	y := x.Right
	T2 := y.Left

	// Perform rotation
	y.Left = x
	x.Right = T2

	// Update heights
	updateHeight(x)
	updateHeight(y)

	// Return new root
	return y
}

// Insert adds a new key to the AVL tree
// Time Complexity: O(log n)
// Space Complexity: O(log n) due to recursion stack
func (avl *AVLTree) Insert(key int) {
	avl.Root = avl.insertNode(avl.Root, key)
}

// insertNode is a helper function to recursively insert a node
func (avl *AVLTree) insertNode(node *Node, key int) *Node {
	// Step 1: Perform normal BST insertion
	if node == nil {
		avl.Size++
		return &Node{
			Key:    key,
			Height: 1,
			Left:   nil,
			Right:  nil,
		}
	}

	if key < node.Key {
		node.Left = avl.insertNode(node.Left, key)
	} else if key > node.Key {
		node.Right = avl.insertNode(node.Right, key)
	} else {
		// Duplicate keys not allowed, return unchanged node
		return node
	}

	// Step 2: Update height of current node
	updateHeight(node)

	// Step 3: Get the balance factor
	balance := getBalance(node)

	// Step 4: Perform rotations if needed to maintain AVL property

	// Left Left Case
	if balance > 1 && key < node.Left.Key {
		return rightRotate(node)
	}

	// Right Right Case
	if balance < -1 && key > node.Right.Key {
		return leftRotate(node)
	}

	// Left Right Case
	if balance > 1 && key > node.Left.Key {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// Right Left Case
	if balance < -1 && key < node.Right.Key {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	// Return unchanged node if no rotation needed
	return node
}

// Search looks for a key in the AVL tree
// Time Complexity: O(log n)
// Space Complexity: O(log n) due to recursion stack
func (avl *AVLTree) Search(key int) bool {
	return avl.searchNode(avl.Root, key)
}

// searchNode is a helper function to recursively search for a node
func (avl *AVLTree) searchNode(node *Node, key int) bool {
	if node == nil {
		return false
	}

	if key == node.Key {
		return true
	} else if key < node.Key {
		return avl.searchNode(node.Left, key)
	} else {
		return avl.searchNode(node.Right, key)
	}
}

// findMin finds the node with minimum key in a subtree
func findMin(node *Node) *Node {
	for node.Left != nil {
		node = node.Left
	}
	return node
}

// Delete removes a key from the AVL tree
// Time Complexity: O(log n)
// Space Complexity: O(log n) due to recursion stack
func (avl *AVLTree) Delete(key int) {
	avl.Root = avl.deleteNode(avl.Root, key)
}

// deleteNode is a helper function to recursively delete a node
func (avl *AVLTree) deleteNode(node *Node, key int) *Node {
	// Step 1: Perform normal BST deletion
	if node == nil {
		return node
	}

	if key < node.Key {
		node.Left = avl.deleteNode(node.Left, key)
	} else if key > node.Key {
		node.Right = avl.deleteNode(node.Right, key)
	} else {
		// Node to be deleted found
		avl.Size--

		// Node with only one child or no child
		if node.Left == nil || node.Right == nil {
			var temp *Node
			if node.Left == nil {
				temp = node.Right
			} else {
				temp = node.Left
			}

			// No child case
			if temp == nil {
				temp = node
				node = nil
			} else {
				// One child case
				*node = *temp
			}
		} else {
			// Node with two children
			// Get the inorder successor (smallest in the right subtree)
			temp := findMin(node.Right)

			// Copy the inorder successor's key to this node
			node.Key = temp.Key

			// Delete the inorder successor
			node.Right = avl.deleteNode(node.Right, temp.Key)
		}
	}

	// If the tree had only one node, return
	if node == nil {
		return node
	}

	// Step 2: Update height of current node
	updateHeight(node)

	// Step 3: Get the balance factor
	balance := getBalance(node)

	// Step 4: Perform rotations if needed

	// Left Left Case
	if balance > 1 && getBalance(node.Left) >= 0 {
		return rightRotate(node)
	}

	// Left Right Case
	if balance > 1 && getBalance(node.Left) < 0 {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// Right Right Case
	if balance < -1 && getBalance(node.Right) <= 0 {
		return leftRotate(node)
	}

	// Right Left Case
	if balance < -1 && getBalance(node.Right) > 0 {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	// Return unchanged node
	return node
}

// InorderTraversal returns the keys in sorted order
// Time Complexity: O(n)
// Space Complexity: O(n) for result slice + O(log n) for recursion
func (avl *AVLTree) InorderTraversal() []int {
	var result []int
	avl.inorderHelper(avl.Root, &result)
	return result
}

// inorderHelper is a helper function for inorder traversal
func (avl *AVLTree) inorderHelper(node *Node, result *[]int) {
	if node != nil {
		avl.inorderHelper(node.Left, result)
		*result = append(*result, node.Key)
		avl.inorderHelper(node.Right, result)
	}
}

// GetSize returns the number of nodes in the tree
func (avl *AVLTree) GetSize() int {
	return avl.Size
}

// GetHeight returns the height of the tree
func (avl *AVLTree) GetHeight() int {
	return height(avl.Root)
}

// IsEmpty checks if the tree is empty
func (avl *AVLTree) IsEmpty() bool {
	return avl.Root == nil
}

// IsBalanced checks if the tree maintains AVL property
// This is mainly for testing purposes
func (avl *AVLTree) IsBalanced() bool {
	return avl.isBalancedHelper(avl.Root)
}

// isBalancedHelper recursively checks if subtree is balanced
func (avl *AVLTree) isBalancedHelper(node *Node) bool {
	if node == nil {
		return true
	}

	balance := getBalance(node)
	if balance < -1 || balance > 1 {
		return false
	}

	return avl.isBalancedHelper(node.Left) && avl.isBalancedHelper(node.Right)
}

// PrintTree prints the tree structure (for debugging)
func (avl *AVLTree) PrintTree() {
	fmt.Println("AVL Tree Structure:")
	avl.printHelper(avl.Root, "", true)
}

// printHelper is a helper function to print tree structure
func (avl *AVLTree) printHelper(node *Node, prefix string, isLast bool) {
	if node != nil {
		fmt.Printf("%s", prefix)
		if isLast {
			fmt.Printf("└── ")
		} else {
			fmt.Printf("├── ")
		}
		fmt.Printf("%d (h:%d, b:%d)\n", node.Key, node.Height, getBalance(node))

		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		if node.Left != nil || node.Right != nil {
			if node.Right != nil {
				avl.printHelper(node.Right, newPrefix, node.Left == nil)
			}
			if node.Left != nil {
				avl.printHelper(node.Left, newPrefix, true)
			}
		}
	}
}
