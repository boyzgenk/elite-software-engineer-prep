// Package segmenttree implements various Segment Tree data structures.
// Segment Trees provide efficient range queries and updates in O(log n) time.
// They're essential for competitive programming and advanced system algorithms.
//
// Key Operations:
// - Range Query: O(log n) - sum, min, max, or custom operations over ranges
// - Point Update: O(log n) - update single elements
// - Range Update: O(log n) - update entire ranges (with lazy propagation)
// - Build: O(n) - construct tree from initial array
//
// Space Complexity: O(4n) ≈ O(n)
//
// Use Cases:
// - Range sum/min/max queries with updates
// - Competitive programming problems
// - Database indexing for range operations
// - Real-time analytics systems
// - Game development (collision detection, scoring)
package segmenttree

import (
	"fmt"
	"math"
	"strings"
)

// Operation defines the type of operation supported by the segment tree
type Operation int

const (
	Sum Operation = iota
	Min
	Max
	Custom
)

// String returns string representation of operation
func (op Operation) String() string {
	switch op {
	case Sum:
		return "Sum"
	case Min:
		return "Min"
	case Max:
		return "Max"
	case Custom:
		return "Custom"
	default:
		return "Unknown"
	}
}

// CombineFunc defines a function that combines two values
type CombineFunc func(a, b int) int

// SegmentTree represents a segment tree data structure
type SegmentTree struct {
	tree      []int       // Tree stored in array form
	lazy      []int       // Lazy propagation array
	size      int         // Size of original array
	operation Operation   // Type of operation (sum, min, max, custom)
	combine   CombineFunc // Function to combine values
	identity  int         // Identity element for the operation
}

// NewSegmentTree creates a new segment tree with the specified operation
func NewSegmentTree(arr []int, op Operation) *SegmentTree {
	n := len(arr)
	if n == 0 {
		return nil
	}

	st := &SegmentTree{
		tree:      make([]int, 4*n), // 4*n is sufficient for any array size
		lazy:      make([]int, 4*n),
		size:      n,
		operation: op,
	}

	// Set combine function and identity based on operation
	switch op {
	case Sum:
		st.combine = func(a, b int) int { return a + b }
		st.identity = 0
	case Min:
		st.combine = func(a, b int) int {
			if a < b {
				return a
			}
			return b
		}
		st.identity = math.MaxInt32
	case Max:
		st.combine = func(a, b int) int {
			if a > b {
				return a
			}
			return b
		}
		st.identity = math.MinInt32
	}

	st.build(arr, 1, 0, n-1)
	return st
}

// NewCustomSegmentTree creates a segment tree with custom combine function
func NewCustomSegmentTree(arr []int, combine CombineFunc, identity int) *SegmentTree {
	n := len(arr)
	if n == 0 {
		return nil
	}

	st := &SegmentTree{
		tree:      make([]int, 4*n),
		lazy:      make([]int, 4*n),
		size:      n,
		operation: Custom,
		combine:   combine,
		identity:  identity,
	}

	st.build(arr, 1, 0, n-1)
	return st
}

// build constructs the segment tree from the input array
func (st *SegmentTree) build(arr []int, vertex, tl, tr int) {
	if tl == tr {
		// Leaf node
		st.tree[vertex] = arr[tl]
	} else {
		// Internal node
		tm := (tl + tr) / 2
		st.build(arr, 2*vertex, tl, tm)
		st.build(arr, 2*vertex+1, tm+1, tr)
		st.tree[vertex] = st.combine(st.tree[2*vertex], st.tree[2*vertex+1])
	}
}

// pushLazy pushes lazy propagation down to children
func (st *SegmentTree) pushLazy(vertex, tl, tr int) {
	if st.lazy[vertex] != 0 {
		// Apply lazy value to current node
		if st.operation == Sum {
			st.tree[vertex] += st.lazy[vertex] * (tr - tl + 1)
		} else {
			st.tree[vertex] += st.lazy[vertex]
		}

		// Push to children if not leaf
		if tl != tr {
			st.lazy[2*vertex] += st.lazy[vertex]
			st.lazy[2*vertex+1] += st.lazy[vertex]
		}

		st.lazy[vertex] = 0
	}
}

// Query performs a range query from l to r (inclusive)
func (st *SegmentTree) Query(l, r int) int {
	if l > r {
		return st.identity
	}

	// Clamp bounds to valid range
	if l < 0 {
		l = 0
	}
	if r >= st.size {
		r = st.size - 1
	}

	return st.query(1, 0, st.size-1, l, r)
}

// query is the internal recursive query function
func (st *SegmentTree) query(vertex, tl, tr, l, r int) int {
	// Push lazy propagation
	st.pushLazy(vertex, tl, tr)

	if l > r {
		return st.identity
	}
	if l == tl && r == tr {
		return st.tree[vertex]
	}

	tm := (tl + tr) / 2
	leftResult := st.query(2*vertex, tl, tm, l, min(r, tm))
	rightResult := st.query(2*vertex+1, tm+1, tr, max(l, tm+1), r)

	return st.combine(leftResult, rightResult)
}

// Update performs a point update at position pos with value val
func (st *SegmentTree) Update(pos, val int) {
	if pos < 0 || pos >= st.size {
		return
	}
	st.update(1, 0, st.size-1, pos, val)
}

// update is the internal recursive update function
func (st *SegmentTree) update(vertex, tl, tr, pos, val int) {
	// Push lazy propagation
	st.pushLazy(vertex, tl, tr)

	if tl == tr {
		// Leaf node - update value
		st.tree[vertex] = val
	} else {
		tm := (tl + tr) / 2
		if pos <= tm {
			st.update(2*vertex, tl, tm, pos, val)
		} else {
			st.update(2*vertex+1, tm+1, tr, pos, val)
		}

		// Push lazy to children before combining
		st.pushLazy(2*vertex, tl, tm)
		st.pushLazy(2*vertex+1, tm+1, tr)

		st.tree[vertex] = st.combine(st.tree[2*vertex], st.tree[2*vertex+1])
	}
}

// RangeUpdate performs a range update from l to r with value val
func (st *SegmentTree) RangeUpdate(l, r, val int) {
	if l < 0 || r >= st.size || l > r {
		return
	}
	st.rangeUpdate(1, 0, st.size-1, l, r, val)
}

// rangeUpdate is the internal recursive range update function
func (st *SegmentTree) rangeUpdate(vertex, tl, tr, l, r, val int) {
	if l > r {
		return
	}

	if l == tl && r == tr {
		// Complete overlap - apply lazy propagation
		st.lazy[vertex] += val
		st.pushLazy(vertex, tl, tr)
		return
	}

	// Push existing lazy values
	st.pushLazy(vertex, tl, tr)

	tm := (tl + tr) / 2
	st.rangeUpdate(2*vertex, tl, tm, l, min(r, tm), val)
	st.rangeUpdate(2*vertex+1, tm+1, tr, max(l, tm+1), r, val)

	// Push lazy to children before combining
	st.pushLazy(2*vertex, tl, tm)
	st.pushLazy(2*vertex+1, tm+1, tr)

	st.tree[vertex] = st.combine(st.tree[2*vertex], st.tree[2*vertex+1])
}

// GetArray returns the current state of the array
func (st *SegmentTree) GetArray() []int {
	result := make([]int, st.size)
	for i := 0; i < st.size; i++ {
		result[i] = st.Query(i, i)
	}
	return result
}

// Size returns the size of the original array
func (st *SegmentTree) Size() int {
	return st.size
}

// Operation returns the operation type
func (st *SegmentTree) Operation() Operation {
	return st.operation
}

// Height returns the height of the segment tree
func (st *SegmentTree) Height() int {
	if st.size <= 1 {
		return 1
	}
	return int(math.Ceil(math.Log2(float64(st.size)))) + 1
}

// String returns a string representation of the segment tree
func (st *SegmentTree) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Segment Tree (%s):\n", st.operation))
	sb.WriteString(fmt.Sprintf("Size: %d, Height: %d\n", st.size, st.Height()))
	sb.WriteString("Array: ")
	arr := st.GetArray()
	for i, val := range arr {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", val))
	}
	sb.WriteString("\n")
	return sb.String()
}

// Validate checks if the segment tree is correctly constructed
func (st *SegmentTree) Validate() bool {
	return st.validate(1, 0, st.size-1)
}

// validate recursively checks tree structure
func (st *SegmentTree) validate(vertex, tl, tr int) bool {
	// Push lazy propagation before validation
	st.pushLazy(vertex, tl, tr)

	if tl == tr {
		// Leaf node - always valid
		return true
	}

	tm := (tl + tr) / 2

	// Check children first
	if !st.validate(2*vertex, tl, tm) || !st.validate(2*vertex+1, tm+1, tr) {
		return false
	}

	// Push lazy to children before checking combination
	st.pushLazy(2*vertex, tl, tm)
	st.pushLazy(2*vertex+1, tm+1, tr)

	// Check if current node value equals combination of children
	expected := st.combine(st.tree[2*vertex], st.tree[2*vertex+1])
	return st.tree[vertex] == expected
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Specialized Segment Trees

// SumSegmentTree is a specialized segment tree for range sum queries
type SumSegmentTree struct {
	*SegmentTree
}

// NewSumSegmentTree creates a new segment tree optimized for sum operations
func NewSumSegmentTree(arr []int) *SumSegmentTree {
	return &SumSegmentTree{
		SegmentTree: NewSegmentTree(arr, Sum),
	}
}

// RangeSum queries the sum of elements from l to r
func (st *SumSegmentTree) RangeSum(l, r int) int {
	return st.Query(l, r)
}

// MinSegmentTree is a specialized segment tree for range minimum queries
type MinSegmentTree struct {
	*SegmentTree
}

// NewMinSegmentTree creates a new segment tree optimized for minimum operations
func NewMinSegmentTree(arr []int) *MinSegmentTree {
	return &MinSegmentTree{
		SegmentTree: NewSegmentTree(arr, Min),
	}
}

// RangeMin queries the minimum element from l to r
func (st *MinSegmentTree) RangeMin(l, r int) int {
	return st.Query(l, r)
}

// MaxSegmentTree is a specialized segment tree for range maximum queries
type MaxSegmentTree struct {
	*SegmentTree
}

// NewMaxSegmentTree creates a new segment tree optimized for maximum operations
func NewMaxSegmentTree(arr []int) *MaxSegmentTree {
	return &MaxSegmentTree{
		SegmentTree: NewSegmentTree(arr, Max),
	}
}

// RangeMax queries the maximum element from l to r
func (st *MaxSegmentTree) RangeMax(l, r int) int {
	return st.Query(l, r)
}

// Advanced Operations

// RangeGCD computes GCD of elements in range using custom segment tree
func NewGCDSegmentTree(arr []int) *SegmentTree {
	gcd := func(a, b int) int {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}
	return NewCustomSegmentTree(arr, gcd, 0)
}

// RangeLCM computes LCM of elements in range using custom segment tree
func NewLCMSegmentTree(arr []int) *SegmentTree {
	gcd := func(a, b int) int {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}

	lcm := func(a, b int) int {
		if a == 0 || b == 0 {
			return 0
		}
		return (a * b) / gcd(a, b)
	}

	return NewCustomSegmentTree(arr, lcm, 1)
}
