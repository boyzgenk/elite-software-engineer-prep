// Package redblacktree implements a Red-Black Tree data structure.
// Red-Black Trees maintain balance through color properties rather than strict height balance,
// providing O(log n) operations with better insertion/deletion performance than AVL trees.
//
// Red-Black Tree Properties:
// 1. Every node is either red or black
// 2. Root is always black
// 3. All leaves (NIL nodes) are black
// 4. Red nodes cannot have red children (no two red nodes adjacent)
// 5. Every path from root to leaf contains same number of black nodes
//
// Time Complexity: O(log n) for all operations
// Space Complexity: O(n)
//
// Use Cases:
// - Language standard libraries (C++ std::map, Java TreeMap)
// - Database indexing systems
// - Memory allocation algorithms
// - High-frequency trading systems requiring consistent performance
package redblacktree

import (
	"fmt"
	"strings"
)

// Color represents the color of a Red-Black Tree node
type Color bool

const (
	Red   Color = false
	Black Color = true
)

// String returns string representation of color
func (c Color) String() string {
	if c == Red {
		return "Red"
	}
	return "Black"
}

// Node represents a node in the Red-Black Tree
type Node struct {
	Key    int
	Value  interface{}
	Color  Color
	Left   *Node
	Right  *Node
	Parent *Node
}

// NewNode creates a new Red-Black Tree node with given key and value
func NewNode(key int, value interface{}) *Node {
	return &Node{
		Key:   key,
		Value: value,
		Color: Red, // New nodes are always inserted as red
	}
}

// IsRed checks if node is red (handles nil nodes as black)
func (n *Node) IsRed() bool {
	return n != nil && n.Color == Red
}

// IsBlack checks if node is black (handles nil nodes as black)
func (n *Node) IsBlack() bool {
	return n == nil || n.Color == Black
}

// Grandparent returns the grandparent node
func (n *Node) Grandparent() *Node {
	if n != nil && n.Parent != nil {
		return n.Parent.Parent
	}
	return nil
}

// Uncle returns the uncle node (parent's sibling)
func (n *Node) Uncle() *Node {
	gp := n.Grandparent()
	if gp == nil {
		return nil
	}
	if n.Parent == gp.Left {
		return gp.Right
	}
	return gp.Left
}

// Sibling returns the sibling node
func (n *Node) Sibling() *Node {
	if n == nil || n.Parent == nil {
		return nil
	}
	if n == n.Parent.Left {
		return n.Parent.Right
	}
	return n.Parent.Left
}

// RedBlackTree represents a Red-Black Tree
type RedBlackTree struct {
	Root *Node
	Size int
}

// NewRedBlackTree creates a new empty Red-Black Tree
func NewRedBlackTree() *RedBlackTree {
	return &RedBlackTree{}
}

// rotateLeft performs left rotation around node x
//
//	  x               y
//	 / \             / \
//	a   y    -->    x   c
//	   / \         / \
//	  b   c       a   b
func (rbt *RedBlackTree) rotateLeft(x *Node) {
	if x == nil || x.Right == nil {
		return
	}

	y := x.Right
	x.Right = y.Left

	if y.Left != nil {
		y.Left.Parent = x
	}

	y.Parent = x.Parent

	if x.Parent == nil {
		rbt.Root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Left = x
	x.Parent = y
}

// rotateRight performs right rotation around node y
//
//	    y             x
//	   / \           / \
//	  x   c   -->   a   y
//	 / \               / \
//	a   b             b   c
func (rbt *RedBlackTree) rotateRight(y *Node) {
	if y == nil || y.Left == nil {
		return
	}

	x := y.Left
	y.Left = x.Right

	if x.Right != nil {
		x.Right.Parent = y
	}

	x.Parent = y.Parent

	if y.Parent == nil {
		rbt.Root = x
	} else if y == y.Parent.Left {
		y.Parent.Left = x
	} else {
		y.Parent.Right = x
	}

	x.Right = y
	y.Parent = x
}

// Insert adds a new key-value pair to the tree
func (rbt *RedBlackTree) Insert(key int, value interface{}) {
	newNode := NewNode(key, value)

	if rbt.Root == nil {
		newNode.Color = Black
		rbt.Root = newNode
		rbt.Size++
		return
	}

	// Standard BST insertion
	current := rbt.Root
	for {
		if key < current.Key {
			if current.Left == nil {
				current.Left = newNode
				newNode.Parent = current
				break
			}
			current = current.Left
		} else if key > current.Key {
			if current.Right == nil {
				current.Right = newNode
				newNode.Parent = current
				break
			}
			current = current.Right
		} else {
			// Key already exists, update value
			current.Value = value
			return
		}
	}

	rbt.Size++
	rbt.insertFixup(newNode)
}

// insertFixup maintains Red-Black Tree properties after insertion
func (rbt *RedBlackTree) insertFixup(node *Node) {
	for node.Parent != nil && node.Parent.IsRed() {
		if node.Parent == node.Grandparent().Left {
			uncle := node.Uncle()

			if uncle.IsRed() {
				// Case 1: Uncle is red
				node.Parent.Color = Black
				uncle.Color = Black
				node.Grandparent().Color = Red
				node = node.Grandparent()
			} else {
				if node == node.Parent.Right {
					// Case 2: Uncle is black, node is right child
					node = node.Parent
					rbt.rotateLeft(node)
				}
				// Case 3: Uncle is black, node is left child
				node.Parent.Color = Black
				node.Grandparent().Color = Red
				rbt.rotateRight(node.Grandparent())
			}
		} else {
			uncle := node.Uncle()

			if uncle.IsRed() {
				// Case 1: Uncle is red (symmetric)
				node.Parent.Color = Black
				uncle.Color = Black
				node.Grandparent().Color = Red
				node = node.Grandparent()
			} else {
				if node == node.Parent.Left {
					// Case 2: Uncle is black, node is left child (symmetric)
					node = node.Parent
					rbt.rotateRight(node)
				}
				// Case 3: Uncle is black, node is right child (symmetric)
				node.Parent.Color = Black
				node.Grandparent().Color = Red
				rbt.rotateLeft(node.Grandparent())
			}
		}
	}

	rbt.Root.Color = Black // Root is always black
}

// Delete removes a node with the given key from the tree
func (rbt *RedBlackTree) Delete(key int) bool {
	node, found := rbt.Search(key)
	if !found {
		return false
	}

	rbt.deleteNode(node)
	rbt.Size--
	return true
}

// deleteNode removes the specified node from the tree
func (rbt *RedBlackTree) deleteNode(node *Node) {
	var nodeToDelete *Node
	var child *Node

	// Find the node to actually delete (node or its successor)
	if node.Left == nil || node.Right == nil {
		nodeToDelete = node
	} else {
		// Node has two children, find successor
		successor := rbt.minNode(node.Right)
		node.Key = successor.Key
		node.Value = successor.Value
		nodeToDelete = successor
	}

	// Get the child that will replace the deleted node
	if nodeToDelete.Left != nil {
		child = nodeToDelete.Left
	} else {
		child = nodeToDelete.Right
	}

	// Link child to parent
	if child != nil {
		child.Parent = nodeToDelete.Parent
	}

	if nodeToDelete.Parent == nil {
		rbt.Root = child
	} else if nodeToDelete == nodeToDelete.Parent.Left {
		nodeToDelete.Parent.Left = child
	} else {
		nodeToDelete.Parent.Right = child
	}

	// Fix Red-Black properties if deleted node was black
	if nodeToDelete.IsBlack() {
		rbt.deleteFixup(child, nodeToDelete.Parent)
	}
}

// deleteFixup maintains Red-Black Tree properties after deletion
func (rbt *RedBlackTree) deleteFixup(node *Node, parent *Node) {
	for node != rbt.Root && (node == nil || node.IsBlack()) {
		if node == parent.Left {
			sibling := parent.Right

			// Case 1: Sibling is red
			if sibling != nil && sibling.IsRed() {
				sibling.Color = Black
				parent.Color = Red
				rbt.rotateLeft(parent)
				sibling = parent.Right
			}

			// Case 2: Sibling is black with two black children
			if (sibling == nil || sibling.IsBlack()) &&
				(sibling == nil || (sibling.Left == nil || sibling.Left.IsBlack())) &&
				(sibling == nil || (sibling.Right == nil || sibling.Right.IsBlack())) {
				if sibling != nil {
					sibling.Color = Red
				}
				node = parent
				parent = node.Parent
			} else {
				// Case 3: Sibling is black with red left child and black right child
				if sibling != nil &&
					(sibling.Right == nil || sibling.Right.IsBlack()) &&
					sibling.Left != nil && sibling.Left.IsRed() {
					sibling.Left.Color = Black
					sibling.Color = Red
					rbt.rotateRight(sibling)
					sibling = parent.Right
				}

				// Case 4: Sibling is black with red right child
				if sibling != nil {
					sibling.Color = parent.Color
					parent.Color = Black
					if sibling.Right != nil {
						sibling.Right.Color = Black
					}
					rbt.rotateLeft(parent)
				}
				node = rbt.Root // This will terminate the loop
			}
		} else {
			// Symmetric cases when node is right child
			sibling := parent.Left

			// Case 1: Sibling is red (symmetric)
			if sibling != nil && sibling.IsRed() {
				sibling.Color = Black
				parent.Color = Red
				rbt.rotateRight(parent)
				sibling = parent.Left
			}

			// Case 2: Sibling is black with two black children (symmetric)
			if (sibling == nil || sibling.IsBlack()) &&
				(sibling == nil || (sibling.Left == nil || sibling.Left.IsBlack())) &&
				(sibling == nil || (sibling.Right == nil || sibling.Right.IsBlack())) {
				if sibling != nil {
					sibling.Color = Red
				}
				node = parent
				parent = node.Parent
			} else {
				// Case 3: Sibling is black with red right child and black left child (symmetric)
				if sibling != nil &&
					(sibling.Left == nil || sibling.Left.IsBlack()) &&
					sibling.Right != nil && sibling.Right.IsRed() {
					sibling.Right.Color = Black
					sibling.Color = Red
					rbt.rotateLeft(sibling)
					sibling = parent.Left
				}

				// Case 4: Sibling is black with red left child (symmetric)
				if sibling != nil {
					sibling.Color = parent.Color
					parent.Color = Black
					if sibling.Left != nil {
						sibling.Left.Color = Black
					}
					rbt.rotateRight(parent)
				}
				node = rbt.Root // This will terminate the loop
			}
		}
	}

	if node != nil {
		node.Color = Black
	}
}

// Search finds a node with the given key
func (rbt *RedBlackTree) Search(key int) (*Node, bool) {
	current := rbt.Root
	for current != nil {
		if key == current.Key {
			return current, true
		} else if key < current.Key {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	return nil, false
}

// Get retrieves the value associated with the given key
func (rbt *RedBlackTree) Get(key int) (interface{}, bool) {
	node, found := rbt.Search(key)
	if found {
		return node.Value, true
	}
	return nil, false
}

// Contains checks if the tree contains the given key
func (rbt *RedBlackTree) Contains(key int) bool {
	_, found := rbt.Search(key)
	return found
}

// Min returns the node with minimum key
func (rbt *RedBlackTree) Min() *Node {
	if rbt.Root == nil {
		return nil
	}
	return rbt.minNode(rbt.Root)
}

// minNode finds the minimum node in subtree rooted at node
func (rbt *RedBlackTree) minNode(node *Node) *Node {
	for node.Left != nil {
		node = node.Left
	}
	return node
}

// Max returns the node with maximum key
func (rbt *RedBlackTree) Max() *Node {
	if rbt.Root == nil {
		return nil
	}
	return rbt.maxNode(rbt.Root)
}

// maxNode finds the maximum node in subtree rooted at node
func (rbt *RedBlackTree) maxNode(node *Node) *Node {
	for node.Right != nil {
		node = node.Right
	}
	return node
}

// Height returns the height of the tree
func (rbt *RedBlackTree) Height() int {
	return rbt.height(rbt.Root)
}

// height calculates height of subtree rooted at node
func (rbt *RedBlackTree) height(node *Node) int {
	if node == nil {
		return 0
	}
	leftHeight := rbt.height(node.Left)
	rightHeight := rbt.height(node.Right)
	if leftHeight > rightHeight {
		return leftHeight + 1
	}
	return rightHeight + 1
}

// BlackHeight returns the black height of the tree
func (rbt *RedBlackTree) BlackHeight() int {
	return rbt.blackHeight(rbt.Root)
}

// blackHeight calculates black height from node to any leaf
func (rbt *RedBlackTree) blackHeight(node *Node) int {
	if node == nil {
		return 1 // NIL nodes are black
	}

	leftBH := rbt.blackHeight(node.Left)
	if node.IsBlack() {
		return leftBH + 1
	}
	return leftBH
}

// IsValid validates all Red-Black Tree properties
func (rbt *RedBlackTree) IsValid() bool {
	if rbt.Root == nil {
		return true
	}

	// Property 2: Root is black
	if rbt.Root.IsRed() {
		return false
	}

	// Check all properties recursively
	_, valid := rbt.validateProperties(rbt.Root)
	return valid
}

// validateProperties checks Red-Black Tree properties recursively
func (rbt *RedBlackTree) validateProperties(node *Node) (int, bool) {
	if node == nil {
		return 1, true // NIL nodes are black
	}

	// Property 4: Red nodes cannot have red children
	if node.IsRed() {
		if (node.Left != nil && node.Left.IsRed()) ||
			(node.Right != nil && node.Right.IsRed()) {
			return 0, false
		}
	}

	// Check subtrees
	leftBH, leftValid := rbt.validateProperties(node.Left)
	rightBH, rightValid := rbt.validateProperties(node.Right)

	if !leftValid || !rightValid {
		return 0, false
	}

	// Property 5: Same black height for all paths
	if leftBH != rightBH {
		return 0, false
	}

	if node.IsBlack() {
		return leftBH + 1, true
	}
	return leftBH, true
}

// InorderTraversal returns keys in sorted order
func (rbt *RedBlackTree) InorderTraversal() []int {
	var result []int
	rbt.inorder(rbt.Root, &result)
	return result
}

// inorder performs inorder traversal
func (rbt *RedBlackTree) inorder(node *Node, result *[]int) {
	if node != nil {
		rbt.inorder(node.Left, result)
		*result = append(*result, node.Key)
		rbt.inorder(node.Right, result)
	}
}

// String returns string representation of the tree
func (rbt *RedBlackTree) String() string {
	if rbt.Root == nil {
		return "Empty Red-Black Tree"
	}

	var sb strings.Builder
	sb.WriteString("Red-Black Tree:\n")
	rbt.printTree(rbt.Root, "", true, &sb)
	return sb.String()
}

// printTree creates visual representation of the tree
func (rbt *RedBlackTree) printTree(node *Node, prefix string, isLast bool, sb *strings.Builder) {
	if node != nil {
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		colorStr := "R"
		if node.IsBlack() {
			colorStr = "B"
		}

		sb.WriteString(fmt.Sprintf("%s%s%d(%s)\n", prefix, connector, node.Key, colorStr))

		childPrefix := prefix
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}

		children := []*Node{node.Left, node.Right}
		for i, child := range children {
			if child != nil {
				rbt.printTree(child, childPrefix, i == 1 && children[1] != nil, sb)
			} else if node.Left != nil || node.Right != nil {
				sb.WriteString(fmt.Sprintf("%s%s(nil)\n", childPrefix,
					func() string {
						if i == 1 {
							return "└── "
						} else {
							return "├── "
						}
					}()))
			}
		}
	}
}
