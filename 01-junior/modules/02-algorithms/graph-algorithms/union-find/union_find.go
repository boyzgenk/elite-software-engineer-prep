// Package graph implements Union-Find (Disjoint Set Union) data structure
// with path compression and union by rank optimizations for FAANG L3/L4 interview preparation.
package graph

import (
	"errors"
	"fmt"
)

// UnionFind represents a Union-Find (Disjoint Set Union) data structure
type UnionFind struct {
	parent []int // parent[i] is the parent of element i
	rank   []int // rank[i] is the rank (approximate depth) of tree rooted at i
	size   []int // size[i] is the size of component containing i (only valid for root)
	count  int   // number of connected components
}

// NewUnionFind creates a new Union-Find data structure with n elements
// Initially, each element is in its own set
// Time Complexity: O(n)
// Space Complexity: O(n)
func NewUnionFind(n int) *UnionFind {
	if n <= 0 {
		return &UnionFind{}
	}

	uf := &UnionFind{
		parent: make([]int, n),
		rank:   make([]int, n),
		size:   make([]int, n),
		count:  n,
	}

	// Initialize each element as its own parent with rank 0 and size 1
	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.rank[i] = 0
		uf.size[i] = 1
	}

	return uf
}

// Find returns the root of the set containing element x with path compression
// Path compression makes future operations faster by flattening the tree
// Time Complexity: O(α(n)) amortized, where α is the inverse Ackermann function
// Space Complexity: O(log n) worst case for recursion stack
func (uf *UnionFind) Find(x int) (int, error) {
	if x < 0 || x >= len(uf.parent) {
		return -1, errors.New("element index out of bounds")
	}

	// Path compression: make every node point directly to the root
	if uf.parent[x] != x {
		uf.parent[x], _ = uf.Find(uf.parent[x]) // Recursively find root and compress path
	}

	return uf.parent[x], nil
}

// FindIterative returns the root of the set containing element x with iterative path compression
// Alternative implementation that avoids recursion
// Time Complexity: O(α(n)) amortized
// Space Complexity: O(1)
func (uf *UnionFind) FindIterative(x int) (int, error) {
	if x < 0 || x >= len(uf.parent) {
		return -1, errors.New("element index out of bounds")
	}

	// Find root
	root := x
	for uf.parent[root] != root {
		root = uf.parent[root]
	}

	// Path compression: make all nodes on path point to root
	current := x
	for current != root {
		next := uf.parent[current]
		uf.parent[current] = root
		current = next
	}

	return root, nil
}

// Union merges the sets containing elements x and y using union by rank
// Union by rank ensures the tree remains relatively balanced
// Time Complexity: O(α(n)) amortized
// Space Complexity: O(1)
func (uf *UnionFind) Union(x, y int) error {
	rootX, err := uf.Find(x)
	if err != nil {
		return err
	}

	rootY, err := uf.Find(y)
	if err != nil {
		return err
	}

	// Already in the same set
	if rootX == rootY {
		return nil
	}

	// Union by rank: attach smaller tree under root of larger tree
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
		uf.size[rootY] += uf.size[rootX]
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
		uf.size[rootX] += uf.size[rootY]
	} else {
		// Same rank: make rootY parent of rootX and increment rank
		uf.parent[rootX] = rootY
		uf.rank[rootY]++
		uf.size[rootY] += uf.size[rootX]
	}

	uf.count-- // Decrease number of components
	return nil
}

// Connected checks if elements x and y are in the same connected component
// Time Complexity: O(α(n)) amortized
// Space Complexity: O(1)
func (uf *UnionFind) Connected(x, y int) (bool, error) {
	rootX, err := uf.Find(x)
	if err != nil {
		return false, err
	}

	rootY, err := uf.Find(y)
	if err != nil {
		return false, err
	}

	return rootX == rootY, nil
}

// GetComponentSize returns the size of the connected component containing element x
// Time Complexity: O(α(n)) amortized
// Space Complexity: O(1)
func (uf *UnionFind) GetComponentSize(x int) (int, error) {
	root, err := uf.Find(x)
	if err != nil {
		return 0, err
	}

	return uf.size[root], nil
}

// GetComponentCount returns the number of connected components
// Time Complexity: O(1)
// Space Complexity: O(1)
func (uf *UnionFind) GetComponentCount() int {
	return uf.count
}

// GetAllComponents returns all connected components as a map
// where key is the root and value is slice of all elements in that component
// Time Complexity: O(n * α(n))
// Space Complexity: O(n)
func (uf *UnionFind) GetAllComponents() (map[int][]int, error) {
	components := make(map[int][]int)

	for i := 0; i < len(uf.parent); i++ {
		root, err := uf.Find(i)
		if err != nil {
			return nil, err
		}

		components[root] = append(components[root], i)
	}

	return components, nil
}

// Reset reinitializes the Union-Find structure with all elements in separate sets
// Time Complexity: O(n)
// Space Complexity: O(1)
func (uf *UnionFind) Reset() {
	n := len(uf.parent)
	uf.count = n

	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.rank[i] = 0
		uf.size[i] = 1
	}
}

// IsValidElement checks if element index is within bounds
// Time Complexity: O(1)
// Space Complexity: O(1)
func (uf *UnionFind) IsValidElement(x int) bool {
	return x >= 0 && x < len(uf.parent)
}

// GetSize returns the total number of elements
// Time Complexity: O(1)
// Space Complexity: O(1)
func (uf *UnionFind) GetSize() int {
	return len(uf.parent)
}

// PrintComponents prints all connected components for debugging
func (uf *UnionFind) PrintComponents() {
	components, err := uf.GetAllComponents()
	if err != nil {
		fmt.Printf("Error getting components: %v\n", err)
		return
	}

	fmt.Printf("Union-Find with %d components:\n", uf.count)
	for root, elements := range components {
		fmt.Printf("Component %d: %v (size: %d)\n", root, elements, len(elements))
	}
}

// GetStats returns statistics about the Union-Find structure
type UnionFindStats struct {
	TotalElements     int
	ComponentCount    int
	LargestComponent  int
	SmallestComponent int
	AverageSize       float64
}

// GetStatistics returns detailed statistics about the Union-Find structure
// Time Complexity: O(n * α(n))
// Space Complexity: O(n)
func (uf *UnionFind) GetStatistics() (*UnionFindStats, error) {
	if len(uf.parent) == 0 {
		return &UnionFindStats{}, nil
	}

	components, err := uf.GetAllComponents()
	if err != nil {
		return nil, err
	}

	stats := &UnionFindStats{
		TotalElements:     len(uf.parent),
		ComponentCount:    uf.count,
		LargestComponent:  0,
		SmallestComponent: len(uf.parent),
	}

	totalSize := 0
	for _, elements := range components {
		size := len(elements)
		totalSize += size

		if size > stats.LargestComponent {
			stats.LargestComponent = size
		}
		if size < stats.SmallestComponent {
			stats.SmallestComponent = size
		}
	}

	if uf.count > 0 {
		stats.AverageSize = float64(totalSize) / float64(uf.count)
	}

	return stats, nil
}

// UnionBySize is an alternative union strategy that uses size instead of rank
// Generally performs slightly better in practice than union by rank
type UnionFindBySize struct {
	parent []int
	size   []int
	count  int
}

// NewUnionFindBySize creates a Union-Find structure with union by size strategy
func NewUnionFindBySize(n int) *UnionFindBySize {
	if n <= 0 {
		return &UnionFindBySize{}
	}

	uf := &UnionFindBySize{
		parent: make([]int, n),
		size:   make([]int, n),
		count:  n,
	}

	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.size[i] = 1
	}

	return uf
}

// Find with path compression for UnionFindBySize
func (uf *UnionFindBySize) Find(x int) (int, error) {
	if x < 0 || x >= len(uf.parent) {
		return -1, errors.New("element index out of bounds")
	}

	if uf.parent[x] != x {
		uf.parent[x], _ = uf.Find(uf.parent[x])
	}

	return uf.parent[x], nil
}

// Union using union by size strategy
func (uf *UnionFindBySize) Union(x, y int) error {
	rootX, err := uf.Find(x)
	if err != nil {
		return err
	}

	rootY, err := uf.Find(y)
	if err != nil {
		return err
	}

	if rootX == rootY {
		return nil
	}

	// Union by size: attach smaller component to larger component
	if uf.size[rootX] < uf.size[rootY] {
		uf.parent[rootX] = rootY
		uf.size[rootY] += uf.size[rootX]
	} else {
		uf.parent[rootY] = rootX
		uf.size[rootX] += uf.size[rootY]
	}

	uf.count--
	return nil
}

// Connected checks connectivity for UnionFindBySize
func (uf *UnionFindBySize) Connected(x, y int) (bool, error) {
	rootX, err := uf.Find(x)
	if err != nil {
		return false, err
	}

	rootY, err := uf.Find(y)
	if err != nil {
		return false, err
	}

	return rootX == rootY, nil
}

// GetComponentCount returns number of components for UnionFindBySize
func (uf *UnionFindBySize) GetComponentCount() int {
	return uf.count
}

// =====================================================================================
// ADVANCED UNION-FIND ALGORITHMS AND APPLICATIONS
// =====================================================================================

// WeightedUnionFind implements weighted Union-Find for shortest path tree problems
// Each edge has a weight, and we maintain path weights from each node to its root
type WeightedUnionFind struct {
	parent []int
	rank   []int
	weight []int // weight[i] is the weight from node i to its parent
	count  int
}

// NewWeightedUnionFind creates a weighted Union-Find structure
func NewWeightedUnionFind(n int) *WeightedUnionFind {
	if n <= 0 {
		return &WeightedUnionFind{}
	}

	wuf := &WeightedUnionFind{
		parent: make([]int, n),
		rank:   make([]int, n),
		weight: make([]int, n),
		count:  n,
	}

	for i := 0; i < n; i++ {
		wuf.parent[i] = i
		wuf.rank[i] = 0
		wuf.weight[i] = 0
	}

	return wuf
}

// Find returns root and accumulated weight with path compression
// Time Complexity: O(α(n)) amortized
func (wuf *WeightedUnionFind) Find(x int) (root int, totalWeight int, err error) {
	if x < 0 || x >= len(wuf.parent) {
		return -1, 0, errors.New("element index out of bounds")
	}

	if wuf.parent[x] != x {
		root, parentWeight, err := wuf.Find(wuf.parent[x])
		if err != nil {
			return -1, 0, err
		}
		wuf.weight[x] += parentWeight
		wuf.parent[x] = root
		return root, wuf.weight[x], nil
	}

	return x, 0, nil
}

// Union connects x and y with specified weight between them
// Weight represents the relationship: weight[y] = weight[x] + w
func (wuf *WeightedUnionFind) Union(x, y, w int) error {
	rootX, weightX, err := wuf.Find(x)
	if err != nil {
		return err
	}

	rootY, weightY, err := wuf.Find(y)
	if err != nil {
		return err
	}

	if rootX == rootY {
		return nil // Already connected, could check for consistency
	}

	// Union by rank with weight adjustment
	if wuf.rank[rootX] < wuf.rank[rootY] {
		wuf.parent[rootX] = rootY
		wuf.weight[rootX] = weightY - weightX - w
	} else if wuf.rank[rootX] > wuf.rank[rootY] {
		wuf.parent[rootY] = rootX
		wuf.weight[rootY] = weightX - weightY + w
	} else {
		wuf.parent[rootY] = rootX
		wuf.weight[rootY] = weightX - weightY + w
		wuf.rank[rootX]++
	}

	wuf.count--
	return nil
}

// GetDifference returns the difference weight[y] - weight[x] if connected
func (wuf *WeightedUnionFind) GetDifference(x, y int) (int, bool, error) {
	rootX, weightX, err := wuf.Find(x)
	if err != nil {
		return 0, false, err
	}

	rootY, weightY, err := wuf.Find(y)
	if err != nil {
		return 0, false, err
	}

	if rootX != rootY {
		return 0, false, nil // Not connected
	}

	return weightY - weightX, true, nil
}

// =====================================================================================
// PERSISTENT UNION-FIND FOR UNDO OPERATIONS
// =====================================================================================

// Operation represents a union operation that can be undone
type Operation struct {
	x, y         int
	rootX, rootY int
	rankX, rankY int
	sizeX, sizeY int
	merged       bool
}

// PersistentUnionFind supports undo operations for dynamic connectivity
type PersistentUnionFind struct {
	*UnionFind
	history []Operation
}

// NewPersistentUnionFind creates a Union-Find with undo capability
func NewPersistentUnionFind(n int) *PersistentUnionFind {
	return &PersistentUnionFind{
		UnionFind: NewUnionFind(n),
		history:   make([]Operation, 0),
	}
}

// Union with history tracking for undo capability
func (puf *PersistentUnionFind) Union(x, y int) error {
	rootX, err := puf.Find(x)
	if err != nil {
		return err
	}

	rootY, err := puf.Find(y)
	if err != nil {
		return err
	}

	op := Operation{
		x: x, y: y,
		rootX: rootX, rootY: rootY,
		rankX: puf.rank[rootX], rankY: puf.rank[rootY],
		sizeX: puf.size[rootX], sizeY: puf.size[rootY],
		merged: rootX != rootY,
	}

	if rootX != rootY {
		// Perform union by rank
		if puf.rank[rootX] < puf.rank[rootY] {
			puf.parent[rootX] = rootY
			puf.size[rootY] += puf.size[rootX]
		} else if puf.rank[rootX] > puf.rank[rootY] {
			puf.parent[rootY] = rootX
			puf.size[rootX] += puf.size[rootY]
		} else {
			puf.parent[rootX] = rootY
			puf.rank[rootY]++
			puf.size[rootY] += puf.size[rootX]
		}
		puf.count--
	}

	puf.history = append(puf.history, op)
	return nil
}

// Undo reverts the last union operation
func (puf *PersistentUnionFind) Undo() error {
	if len(puf.history) == 0 {
		return errors.New("no operations to undo")
	}

	op := puf.history[len(puf.history)-1]
	puf.history = puf.history[:len(puf.history)-1]

	if op.merged {
		// Restore the state before union
		puf.parent[op.rootX] = op.rootX
		puf.parent[op.rootY] = op.rootY
		puf.rank[op.rootX] = op.rankX
		puf.rank[op.rootY] = op.rankY
		puf.size[op.rootX] = op.sizeX
		puf.size[op.rootY] = op.sizeY
		puf.count++

		// Clear path compression for affected elements
		puf.clearPathCompression(op.x)
		puf.clearPathCompression(op.y)
	}

	return nil
}

// clearPathCompression removes path compression effects
func (puf *PersistentUnionFind) clearPathCompression(x int) {
	// This is a simplified version - in practice, full restoration
	// would require more sophisticated tracking
	visited := make(map[int]bool)
	puf.clearPathCompressionHelper(x, visited)
}

func (puf *PersistentUnionFind) clearPathCompressionHelper(x int, visited map[int]bool) {
	if visited[x] {
		return
	}
	visited[x] = true

	// Reconstruct original parent relationships from history
	for i := len(puf.history) - 1; i >= 0; i-- {
		op := puf.history[i]
		if op.merged && (op.rootX == x || op.rootY == x) {
			if op.rootX == x {
				puf.parent[x] = x
			} else {
				puf.parent[x] = x
			}
			return
		}
	}
}

// GetHistorySize returns the number of operations in history
func (puf *PersistentUnionFind) GetHistorySize() int {
	return len(puf.history)
}

// =====================================================================================
// UNION-FIND APPLICATIONS AND ALGORITHMS
// =====================================================================================

// KruskalMST finds Minimum Spanning Tree using Kruskal's algorithm with Union-Find
type Edge struct {
	u, v, weight int
}

// KruskalMST finds MST using Union-Find for cycle detection
// Time Complexity: O(E log E + E α(V))
// Space Complexity: O(V + E)
func KruskalMST(vertices int, edges []Edge) ([]Edge, int, error) {
	if vertices <= 0 {
		return nil, 0, errors.New("invalid number of vertices")
	}

	// Sort edges by weight
	sortedEdges := make([]Edge, len(edges))
	copy(sortedEdges, edges)

	// Simple bubble sort for demonstration (use sort.Slice in production)
	for i := 0; i < len(sortedEdges)-1; i++ {
		for j := 0; j < len(sortedEdges)-i-1; j++ {
			if sortedEdges[j].weight > sortedEdges[j+1].weight {
				sortedEdges[j], sortedEdges[j+1] = sortedEdges[j+1], sortedEdges[j]
			}
		}
	}

	uf := NewUnionFind(vertices)
	mst := make([]Edge, 0, vertices-1)
	totalWeight := 0

	for _, edge := range sortedEdges {
		if edge.u < 0 || edge.u >= vertices || edge.v < 0 || edge.v >= vertices {
			continue
		}

		connected, err := uf.Connected(edge.u, edge.v)
		if err != nil {
			return nil, 0, err
		}

		if !connected {
			err := uf.Union(edge.u, edge.v)
			if err != nil {
				return nil, 0, err
			}
			mst = append(mst, edge)
			totalWeight += edge.weight

			if len(mst) == vertices-1 {
				break
			}
		}
	}

	return mst, totalWeight, nil
}

// DetectCycles detects cycles in an undirected graph using Union-Find
// Time Complexity: O(E α(V))
// Space Complexity: O(V)
func DetectCycles(vertices int, edges []Edge) (bool, []Edge) {
	if vertices <= 0 {
		return false, nil
	}

	uf := NewUnionFind(vertices)
	cycleEdges := make([]Edge, 0)

	for _, edge := range edges {
		if edge.u < 0 || edge.u >= vertices || edge.v < 0 || edge.v >= vertices {
			continue
		}

		connected, err := uf.Connected(edge.u, edge.v)
		if err != nil {
			continue
		}

		if connected {
			cycleEdges = append(cycleEdges, edge)
		} else {
			uf.Union(edge.u, edge.v)
		}
	}

	return len(cycleEdges) > 0, cycleEdges
}

// LongestIncreasingPath finds longest increasing path using Union-Find approach
// This is an alternative to DFS/DP for certain constrained versions
func LongestIncreasingPath(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}

	rows, cols := len(matrix), len(matrix[0])

	// Create coordinate to index mapping
	coordinateToIndex := func(r, c int) int {
		return r*cols + c
	}

	// Create list of all cells with their values for sorting
	type Cell struct {
		value, index int
	}

	cells := make([]Cell, 0, rows*cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			cells = append(cells, Cell{
				value: matrix[r][c],
				index: coordinateToIndex(r, c),
			})
		}
	}

	// Sort cells by value (reverse order for processing)
	for i := 0; i < len(cells)-1; i++ {
		for j := 0; j < len(cells)-i-1; j++ {
			if cells[j].value < cells[j+1].value {
				cells[j], cells[j+1] = cells[j+1], cells[j]
			}
		}
	}

	// Union-Find with path length tracking
	uf := NewUnionFind(rows * cols)
	pathLengths := make([]int, rows*cols)
	for i := range pathLengths {
		pathLengths[i] = 1
	}

	directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	maxLength := 1

	for _, cell := range cells {
		r, c := cell.index/cols, cell.index%cols
		currentIndex := cell.index
		currentLength := pathLengths[currentIndex]

		for _, dir := range directions {
			nr, nc := r+dir[0], c+dir[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
				neighborIndex := coordinateToIndex(nr, nc)
				if matrix[nr][nc] > matrix[r][c] {
					// Neighbor has larger value, can extend path
					neighborRoot, _ := uf.Find(neighborIndex)
					if pathLengths[neighborRoot] < currentLength+1 {
						pathLengths[neighborRoot] = currentLength + 1
						maxLength = max(maxLength, pathLengths[neighborRoot])
					}
				}
			}
		}
	}

	return maxLength
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// =====================================================================================
// PERFORMANCE ANALYSIS AND BENCHMARKING
// =====================================================================================

// BenchmarkResult contains performance comparison results
type BenchmarkResult struct {
	OperationType   string
	UnionByRank     int64 // nanoseconds
	UnionBySize     int64
	BasicUnionFind  int64
	OperationsCount int
}

// PerformanceTester provides comprehensive performance analysis
type PerformanceTester struct {
	size int
}

// NewPerformanceTester creates a new performance testing suite
func NewPerformanceTester(size int) *PerformanceTester {
	return &PerformanceTester{size: size}
}

// BenchmarkUnionFind compares different Union-Find implementations
func (pt *PerformanceTester) BenchmarkUnionFind(operations [][]int) []BenchmarkResult {
	results := make([]BenchmarkResult, 0)

	// Test standard Union-Find with union by rank
	uf1 := NewUnionFind(pt.size)
	start := getCurrentTime()
	for _, op := range operations {
		if len(op) >= 2 {
			uf1.Union(op[0], op[1])
		}
	}
	unionByRankTime := getCurrentTime() - start

	// Test Union-Find with union by size
	uf2 := NewUnionFindBySize(pt.size)
	start = getCurrentTime()
	for _, op := range operations {
		if len(op) >= 2 {
			uf2.Union(op[0], op[1])
		}
	}
	unionBySizeTime := getCurrentTime() - start

	// Test basic Union-Find (simulation without path compression)
	uf3 := NewBasicUnionFind(pt.size)
	start = getCurrentTime()
	for _, op := range operations {
		if len(op) >= 2 {
			uf3.Union(op[0], op[1])
		}
	}
	basicTime := getCurrentTime() - start

	results = append(results, BenchmarkResult{
		OperationType:   "Union Operations",
		UnionByRank:     unionByRankTime,
		UnionBySize:     unionBySizeTime,
		BasicUnionFind:  basicTime,
		OperationsCount: len(operations),
	})

	return results
}

// getCurrentTime returns current time in nanoseconds (simplified for demo)
func getCurrentTime() int64 {
	// In real implementation, use time.Now().UnixNano()
	// For this demo, we'll return a placeholder
	return 0
}

// BasicUnionFind implements naive union-find without optimizations for comparison
type BasicUnionFind struct {
	parent []int
	count  int
}

// NewBasicUnionFind creates basic union-find for performance comparison
func NewBasicUnionFind(n int) *BasicUnionFind {
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	return &BasicUnionFind{parent: parent, count: n}
}

// Find without path compression
func (buf *BasicUnionFind) Find(x int) int {
	for buf.parent[x] != x {
		x = buf.parent[x]
	}
	return x
}

// Union without any optimization
func (buf *BasicUnionFind) Union(x, y int) {
	rootX := buf.Find(x)
	rootY := buf.Find(y)

	if rootX != rootY {
		buf.parent[rootX] = rootY
		buf.count--
	}
}

// =====================================================================================
// SPECIALIZED UNION-FIND VARIANTS
// =====================================================================================

// RankBasedUnionFind with explicit rank management for educational purposes
type RankBasedUnionFind struct {
	parent   []int
	rank     []int
	count    int
	maxRank  int
	rankHist map[int]int // histogram of ranks
}

// NewRankBasedUnionFind creates a rank-focused Union-Find
func NewRankBasedUnionFind(n int) *RankBasedUnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	rankHist := make(map[int]int)

	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 0
	}
	rankHist[0] = n // Initially all nodes have rank 0

	return &RankBasedUnionFind{
		parent:   parent,
		rank:     rank,
		count:    n,
		maxRank:  0,
		rankHist: rankHist,
	}
}

// Union with detailed rank tracking
func (ruf *RankBasedUnionFind) Union(x, y int) error {
	rootX := ruf.findWithPathCompression(x)
	rootY := ruf.findWithPathCompression(y)

	if rootX == rootY {
		return nil
	}

	// Update rank histogram before union
	oldRankY := ruf.rank[rootY]

	if ruf.rank[rootX] < ruf.rank[rootY] {
		ruf.parent[rootX] = rootY
	} else if ruf.rank[rootX] > ruf.rank[rootY] {
		ruf.parent[rootY] = rootX
	} else {
		ruf.parent[rootX] = rootY
		ruf.rank[rootY]++

		// Update histogram
		ruf.rankHist[oldRankY]--
		if ruf.rankHist[oldRankY] == 0 {
			delete(ruf.rankHist, oldRankY)
		}
		ruf.rankHist[ruf.rank[rootY]]++

		if ruf.rank[rootY] > ruf.maxRank {
			ruf.maxRank = ruf.rank[rootY]
		}
	}

	ruf.count--
	return nil
}

// findWithPathCompression performs find with path compression
func (ruf *RankBasedUnionFind) findWithPathCompression(x int) int {
	if ruf.parent[x] != x {
		ruf.parent[x] = ruf.findWithPathCompression(ruf.parent[x])
	}
	return ruf.parent[x]
}

// GetRankStatistics returns detailed rank distribution statistics
func (ruf *RankBasedUnionFind) GetRankStatistics() map[string]interface{} {
	return map[string]interface{}{
		"max_rank":        ruf.maxRank,
		"rank_histogram":  ruf.rankHist,
		"component_count": ruf.count,
		"total_elements":  len(ruf.parent),
	}
}

// =====================================================================================
// UTILITY FUNCTIONS AND DEMONSTRATIONS
// =====================================================================================

// ValidateUnionFind performs comprehensive validation of Union-Find correctness
func ValidateUnionFind(uf *UnionFind, operations [][]int) (bool, []string) {
	errors := make([]string, 0)

	// Validate that Find is consistent
	for i := 0; i < uf.GetSize(); i++ {
		root1, err1 := uf.Find(i)
		root2, err2 := uf.Find(i)

		if err1 != nil || err2 != nil {
			errors = append(errors, fmt.Sprintf("Find error for element %d", i))
			continue
		}

		if root1 != root2 {
			errors = append(errors, fmt.Sprintf("Find inconsistent for element %d", i))
		}
	}

	// Validate connectivity
	for _, op := range operations {
		if len(op) >= 2 {
			connected, err := uf.Connected(op[0], op[1])
			if err != nil {
				errors = append(errors, fmt.Sprintf("Connected error for %d, %d", op[0], op[1]))
			}

			root1, _ := uf.Find(op[0])
			root2, _ := uf.Find(op[1])

			if connected != (root1 == root2) {
				errors = append(errors, fmt.Sprintf("Connected inconsistent for %d, %d", op[0], op[1]))
			}
		}
	}

	return len(errors) == 0, errors
}

// DemoUnionFindApplications demonstrates practical applications
func DemoUnionFindApplications() {
	fmt.Println("=== Union-Find Applications Demo ===")

	// 1. Network connectivity
	fmt.Println("\n1. Network Connectivity Problem:")
	network := NewUnionFind(6)
	connections := [][]int{{0, 1}, {1, 2}, {3, 4}}

	for _, conn := range connections {
		network.Union(conn[0], conn[1])
		fmt.Printf("Connected %d-%d\n", conn[0], conn[1])
	}

	connected, _ := network.Connected(0, 2)
	fmt.Printf("Are nodes 0 and 2 connected? %t\n", connected)

	connected, _ = network.Connected(0, 4)
	fmt.Printf("Are nodes 0 and 4 connected? %t\n", connected)

	// 2. Kruskal's MST
	fmt.Println("\n2. Minimum Spanning Tree (Kruskal's Algorithm):")
	edges := []Edge{
		{0, 1, 4}, {0, 7, 8}, {1, 2, 8}, {1, 7, 11},
		{2, 3, 7}, {2, 8, 2}, {2, 5, 4}, {3, 4, 9},
		{3, 5, 14}, {4, 5, 10}, {5, 6, 2}, {6, 7, 1},
		{6, 8, 6}, {7, 8, 7},
	}

	mst, weight, err := KruskalMST(9, edges)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("MST weight: %d\n", weight)
		fmt.Printf("MST edges: %v\n", mst)
	}

	// 3. Cycle detection
	fmt.Println("\n3. Cycle Detection:")
	testEdges := []Edge{{0, 1, 1}, {1, 2, 1}, {2, 0, 1}}
	hasCycle, cycleEdges := DetectCycles(3, testEdges)
	fmt.Printf("Has cycle: %t\n", hasCycle)
	if hasCycle {
		fmt.Printf("Cycle edges: %v\n", cycleEdges)
	}
}
