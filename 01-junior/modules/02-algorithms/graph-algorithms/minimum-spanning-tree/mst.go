// Package graph implements comprehensive Minimum Spanning Tree (MST) algorithms
// including Kruskal's, Prim's, Borůvka's, and advanced optimizations for FAANG L3/L4 interview preparation.
package graph

import (
	"container/heap"
	"errors"
	"fmt"
	"sort"
)

// WeightedGraph represents a weighted graph for MST algorithms
type WeightedGraph struct {
	vertices int
	adjList  map[int][]MSTGraphEdge
}

// MSTGraphEdge represents an edge in the weighted graph
type MSTGraphEdge struct {
	To     int
	Weight int
}

// NewWeightedGraph creates a new weighted graph
func NewWeightedGraph(vertices int) *WeightedGraph {
	g := &WeightedGraph{
		vertices: vertices,
		adjList:  make(map[int][]MSTGraphEdge, vertices),
	}
	for i := 0; i < vertices; i++ {
		g.adjList[i] = make([]MSTGraphEdge, 0)
	}
	return g
}

// AddWeightedEdge adds a weighted edge to the graph
func (g *WeightedGraph) AddWeightedEdge(src, dest, weight int) error {
	if src < 0 || src >= g.vertices || dest < 0 || dest >= g.vertices {
		return errors.New("vertex index out of bounds")
	}
	edge := MSTGraphEdge{To: dest, Weight: weight}
	g.adjList[src] = append(g.adjList[src], edge)
	return nil
}

// AddUndirectedWeightedEdge adds an undirected weighted edge
func (g *WeightedGraph) AddUndirectedWeightedEdge(u, v, weight int) error {
	if err := g.AddWeightedEdge(u, v, weight); err != nil {
		return err
	}
	return g.AddWeightedEdge(v, u, weight)
}

// MST Edge represents an edge in the context of MST algorithms
type MSTEdge struct {
	From   int
	To     int
	Weight int
}

// MST Result contains the result of MST computation
type MSTResult struct {
	Edges      []MSTEdge
	TotalCost  int
	IsComplete bool // false if graph is disconnected
}

// EdgeList is a slice of MSTEdge that implements sort.Interface
type EdgeList []MSTEdge

func (e EdgeList) Len() int           { return len(e) }
func (e EdgeList) Less(i, j int) bool { return e[i].Weight < e[j].Weight }
func (e EdgeList) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// ========================================================================
// UNION-FIND DATA STRUCTURE FOR MST ALGORITHMS
// ========================================================================

// UnionFind represents a Union-Find data structure optimized for MST
type UnionFind struct {
	parent []int
	rank   []int
	count  int
}

// NewUnionFind creates a new Union-Find structure
func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{
		parent: make([]int, n),
		rank:   make([]int, n),
		count:  n,
	}
	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.rank[i] = 0
	}
	return uf
}

// Find with path compression
func (uf *UnionFind) Find(x int) (int, error) {
	if x < 0 || x >= len(uf.parent) {
		return -1, errors.New("element index out of bounds")
	}
	if uf.parent[x] != x {
		root, err := uf.Find(uf.parent[x])
		if err != nil {
			return -1, err
		}
		uf.parent[x] = root
	}
	return uf.parent[x], nil
}

// Union by rank
func (uf *UnionFind) Union(x, y int) error {
	rootX, err := uf.Find(x)
	if err != nil {
		return err
	}
	rootY, err := uf.Find(y)
	if err != nil {
		return err
	}

	if rootX != rootY {
		if uf.rank[rootX] < uf.rank[rootY] {
			uf.parent[rootX] = rootY
		} else if uf.rank[rootX] > uf.rank[rootY] {
			uf.parent[rootY] = rootX
		} else {
			uf.parent[rootY] = rootX
			uf.rank[rootX]++
		}
		uf.count--
	}
	return nil
}

// Connected checks if two elements are connected
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

// GetComponentCount returns number of connected components
func (uf *UnionFind) GetComponentCount() int {
	return uf.count
}

// ========================================================================
// KRUSKAL'S ALGORITHM
// ========================================================================

// KruskalMST finds the Minimum Spanning Tree using Kruskal's algorithm
// Uses Union-Find data structure to detect cycles efficiently
// Time Complexity: O(E log E) due to sorting, where E is the number of edges
// Space Complexity: O(V + E) for Union-Find structure and edge storage
func (g *WeightedGraph) KruskalMST() (*MSTResult, error) {
	if g.vertices == 0 {
		return &MSTResult{IsComplete: true}, nil
	}

	// Collect all edges from the graph
	var edges EdgeList
	for vertex := 0; vertex < g.vertices; vertex++ {
		for _, edge := range g.adjList[vertex] {
			// For undirected graph, only add edge once (from smaller to larger vertex)
			if vertex < edge.To {
				edges = append(edges, MSTEdge{From: vertex, To: edge.To, Weight: edge.Weight})
			}
		}
	}

	// Sort edges by weight
	sort.Sort(edges)

	// Initialize Union-Find structure
	uf := NewUnionFind(g.vertices)

	result := &MSTResult{
		Edges:     make([]MSTEdge, 0, g.vertices-1),
		TotalCost: 0,
	}

	// Process edges in order of increasing weight
	for _, edge := range edges {
		// Check if adding this edge creates a cycle
		connected, err := uf.Connected(edge.From, edge.To)
		if err != nil {
			return nil, err
		}

		if !connected {
			// Add edge to MST
			result.Edges = append(result.Edges, edge)
			result.TotalCost += edge.Weight

			// Union the components
			err := uf.Union(edge.From, edge.To)
			if err != nil {
				return nil, err
			}

			// If we have V-1 edges, we're done
			if len(result.Edges) == g.vertices-1 {
				result.IsComplete = true
				break
			}
		}
	}

	// Check if graph is connected (MST should have V-1 edges)
	result.IsComplete = len(result.Edges) == g.vertices-1

	return result, nil
}

// ========================================================================
// PRIM'S ALGORITHM
// ========================================================================

// PrimEdge represents an edge for Prim's algorithm priority queue
type PrimEdge struct {
	From   int
	To     int
	Weight int
	index  int // Index in heap
}

// PrimPriorityQueue implements a min-heap for Prim's algorithm
type PrimPriorityQueue []*PrimEdge

func (pq PrimPriorityQueue) Len() int { return len(pq) }

func (pq PrimPriorityQueue) Less(i, j int) bool {
	return pq[i].Weight < pq[j].Weight
}

func (pq PrimPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PrimPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	edge := x.(*PrimEdge)
	edge.index = n
	*pq = append(*pq, edge)
}

func (pq *PrimPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	edge := old[n-1]
	old[n-1] = nil
	edge.index = -1
	*pq = old[0 : n-1]
	return edge
}

// PrimMST finds the Minimum Spanning Tree using Prim's algorithm
// Uses a priority queue to always select the minimum weight edge
// Time Complexity: O(E log V) using binary heap
// Space Complexity: O(V) for visited array and priority queue
func (g *WeightedGraph) PrimMST(startVertex int) (*MSTResult, error) {
	if startVertex < 0 || startVertex >= g.vertices {
		return nil, errors.New("start vertex out of bounds")
	}

	if g.vertices == 0 {
		return &MSTResult{IsComplete: true}, nil
	}

	visited := make([]bool, g.vertices)
	result := &MSTResult{
		Edges:     make([]MSTEdge, 0, g.vertices-1),
		TotalCost: 0,
	}

	// Priority queue for edges (min-heap by weight)
	pq := &PrimPriorityQueue{}
	heap.Init(pq)

	// Start with the given vertex
	visited[startVertex] = true
	visitedCount := 1

	// Add all edges from start vertex to priority queue
	for _, edge := range g.adjList[startVertex] {
		if !visited[edge.To] {
			heap.Push(pq, &PrimEdge{
				From:   startVertex,
				To:     edge.To,
				Weight: edge.Weight,
			})
		}
	}

	// Process edges until MST is complete or no more edges
	for pq.Len() > 0 && visitedCount < g.vertices {
		// Get minimum weight edge
		minEdge := heap.Pop(pq).(*PrimEdge)

		// Skip if destination is already visited
		if visited[minEdge.To] {
			continue
		}

		// Add edge to MST
		result.Edges = append(result.Edges, MSTEdge{
			From:   minEdge.From,
			To:     minEdge.To,
			Weight: minEdge.Weight,
		})
		result.TotalCost += minEdge.Weight

		// Mark destination as visited
		visited[minEdge.To] = true
		visitedCount++

		// Add all edges from the new vertex to unvisited vertices
		for _, edge := range g.adjList[minEdge.To] {
			if !visited[edge.To] {
				heap.Push(pq, &PrimEdge{
					From:   minEdge.To,
					To:     edge.To,
					Weight: edge.Weight,
				})
			}
		}
	}

	// Check if MST is complete
	result.IsComplete = visitedCount == g.vertices

	return result, nil
}

// ========================================================================
// ADVANCED MST ALGORITHMS
// ========================================================================

// DensePrimMST implements Prim's algorithm optimized for dense graphs
// Uses adjacency matrix instead of priority queue for O(V²) complexity
// Time Complexity: O(V²) - better than heap-based Prim for dense graphs
// Space Complexity: O(V²) for adjacency matrix
func (g *WeightedGraph) DensePrimMST(startVertex int) (*MSTResult, error) {
	if startVertex < 0 || startVertex >= g.vertices {
		return nil, errors.New("start vertex out of bounds")
	}

	if g.vertices == 0 {
		return &MSTResult{IsComplete: true}, nil
	}

	// Build adjacency matrix
	const maxInt = int(^uint(0) >> 1)
	adjMatrix := make([][]int, g.vertices)
	for i := range adjMatrix {
		adjMatrix[i] = make([]int, g.vertices)
		for j := range adjMatrix[i] {
			adjMatrix[i][j] = maxInt
		}
	}

	// Fill adjacency matrix from graph
	for vertex := 0; vertex < g.vertices; vertex++ {
		for _, edge := range g.adjList[vertex] {
			adjMatrix[vertex][edge.To] = edge.Weight
			adjMatrix[edge.To][vertex] = edge.Weight // Assume undirected
		}
	}

	// Prim's algorithm with adjacency matrix
	inMST := make([]bool, g.vertices)
	key := make([]int, g.vertices)
	parent := make([]int, g.vertices)

	// Initialize keys and parents
	for i := 0; i < g.vertices; i++ {
		key[i] = maxInt
		parent[i] = -1
	}

	key[startVertex] = 0

	result := &MSTResult{
		Edges:     make([]MSTEdge, 0, g.vertices-1),
		TotalCost: 0,
	}

	for count := 0; count < g.vertices-1; count++ {
		// Find minimum key vertex not yet in MST
		minKey := maxInt
		minIndex := -1

		for v := 0; v < g.vertices; v++ {
			if !inMST[v] && key[v] < minKey {
				minKey = key[v]
				minIndex = v
			}
		}

		if minIndex == -1 {
			break // Graph is disconnected
		}

		inMST[minIndex] = true

		// Add edge to MST (except for start vertex)
		if parent[minIndex] != -1 {
			result.Edges = append(result.Edges, MSTEdge{
				From:   parent[minIndex],
				To:     minIndex,
				Weight: key[minIndex],
			})
			result.TotalCost += key[minIndex]
		}

		// Update keys of adjacent vertices
		for v := 0; v < g.vertices; v++ {
			if !inMST[v] && adjMatrix[minIndex][v] != maxInt && adjMatrix[minIndex][v] < key[v] {
				key[v] = adjMatrix[minIndex][v]
				parent[v] = minIndex
			}
		}
	}

	result.IsComplete = len(result.Edges) == g.vertices-1
	return result, nil
}

// BoruvkaMST implements Borůvka's algorithm for MST
// Parallelizable algorithm that finds cheapest edge for each component
// Time Complexity: O(E log V)
// Space Complexity: O(V + E)
func (g *WeightedGraph) BoruvkaMST() (*MSTResult, error) {
	if g.vertices == 0 {
		return &MSTResult{IsComplete: true}, nil
	}

	uf := NewUnionFind(g.vertices)
	result := &MSTResult{
		Edges:     make([]MSTEdge, 0, g.vertices-1),
		TotalCost: 0,
	}

	// Continue until we have a single component or can't find more edges
	for uf.GetComponentCount() > 1 {
		// Find cheapest edge for each component
		const maxInt = int(^uint(0) >> 1)
		cheapest := make([]MSTEdge, g.vertices)
		hasEdge := make([]bool, g.vertices)

		// Initialize cheapest edges with max weight
		for i := 0; i < g.vertices; i++ {
			cheapest[i] = MSTEdge{Weight: maxInt}
		}

		// Find cheapest outgoing edge for each component
		for vertex := 0; vertex < g.vertices; vertex++ {
			root1, _ := uf.Find(vertex)

			for _, edge := range g.adjList[vertex] {
				root2, _ := uf.Find(edge.To)

				// If vertices are in different components
				if root1 != root2 {
					mstEdge := MSTEdge{From: vertex, To: edge.To, Weight: edge.Weight}

					// Update cheapest edge for component root1
					if edge.Weight < cheapest[root1].Weight {
						cheapest[root1] = mstEdge
						hasEdge[root1] = true
					}
				}
			}
		}

		// Add all cheapest edges
		edgesAdded := 0
		for i := 0; i < g.vertices; i++ {
			if hasEdge[i] {
				edge := cheapest[i]
				connected, _ := uf.Connected(edge.From, edge.To)

				if !connected {
					result.Edges = append(result.Edges, edge)
					result.TotalCost += edge.Weight
					uf.Union(edge.From, edge.To)
					edgesAdded++
				}
			}
		}

		// If no edges were added, graph is disconnected
		if edgesAdded == 0 {
			break
		}
	}

	result.IsComplete = uf.GetComponentCount() == 1
	return result, nil
}

// ========================================================================
// MST ANALYSIS AND UTILITIES
// ========================================================================

// ValidateMST validates that the given edges form a valid MST
// Checks if the edges form a tree (connected and acyclic) and verifies total weight
func (g *WeightedGraph) ValidateMST(mstEdges []MSTEdge) (bool, error) {
	if len(mstEdges) != g.vertices-1 {
		return false, nil // MST must have exactly V-1 edges
	}

	// Check if edges form a connected graph using Union-Find
	uf := NewUnionFind(g.vertices)

	for _, edge := range mstEdges {
		// Check if edge exists in original graph
		edgeExists := false
		for _, graphEdge := range g.adjList[edge.From] {
			if graphEdge.To == edge.To && graphEdge.Weight == edge.Weight {
				edgeExists = true
				break
			}
		}

		if !edgeExists {
			return false, errors.New("edge not found in original graph")
		}

		// Check for cycle
		connected, err := uf.Connected(edge.From, edge.To)
		if err != nil {
			return false, err
		}

		if connected {
			return false, nil // Cycle detected
		}

		err = uf.Union(edge.From, edge.To)
		if err != nil {
			return false, err
		}
	}

	// Check if all vertices are connected
	return uf.GetComponentCount() == 1, nil
}

// CompareMSTAlgorithms compares the results of different MST algorithms
func (g *WeightedGraph) CompareMSTAlgorithms() (*MSTComparison, error) {
	kruskalResult, err := g.KruskalMST()
	if err != nil {
		return nil, fmt.Errorf("Kruskal's algorithm failed: %v", err)
	}

	primResult, err := g.PrimMST(0)
	if err != nil {
		return nil, fmt.Errorf("Prim's algorithm failed: %v", err)
	}

	boruvkaResult, err := g.BoruvkaMST()
	if err != nil {
		return nil, fmt.Errorf("Borůvka's algorithm failed: %v", err)
	}

	comparison := &MSTComparison{
		KruskalResult: kruskalResult,
		PrimResult:    primResult,
		BoruvkaResult: boruvkaResult,
		AllAgree:      true,
	}

	// Check if all algorithms produce the same total cost
	if kruskalResult.TotalCost != primResult.TotalCost ||
		kruskalResult.TotalCost != boruvkaResult.TotalCost {
		comparison.AllAgree = false
	}

	// Check if all algorithms agree on completeness
	if kruskalResult.IsComplete != primResult.IsComplete ||
		kruskalResult.IsComplete != boruvkaResult.IsComplete {
		comparison.AllAgree = false
	}

	return comparison, nil
}

// MSTComparison contains results from different MST algorithms
type MSTComparison struct {
	KruskalResult *MSTResult
	PrimResult    *MSTResult
	BoruvkaResult *MSTResult
	AllAgree      bool
}

// MSTProfiling contains performance metrics for MST algorithms
type MSTProfiling struct {
	Algorithm       string
	EdgeCount       int
	VertexCount     int
	TotalCost       int
	IsComplete      bool
	TimeComplexity  string
	SpaceComplexity string
}

// PrintMST prints the MST result in a readable format
func (result *MSTResult) PrintMST() {
	fmt.Printf("Minimum Spanning Tree:\n")
	fmt.Printf("Complete: %v\n", result.IsComplete)
	fmt.Printf("Total Cost: %d\n", result.TotalCost)
	fmt.Printf("Edges (%d):\n", len(result.Edges))

	for i, edge := range result.Edges {
		fmt.Printf("  %d: %d -- %d (weight: %d)\n", i+1, edge.From, edge.To, edge.Weight)
	}
}

// GetMSTStats returns statistics about the MST
type MSTStats struct {
	EdgeCount     int
	TotalWeight   int
	MinEdgeWeight int
	MaxEdgeWeight int
	AvgEdgeWeight float64
	IsComplete    bool
}

// GetStatistics returns detailed statistics about the MST
func (result *MSTResult) GetStatistics() *MSTStats {
	if len(result.Edges) == 0 {
		return &MSTStats{IsComplete: result.IsComplete}
	}

	stats := &MSTStats{
		EdgeCount:     len(result.Edges),
		TotalWeight:   result.TotalCost,
		MinEdgeWeight: result.Edges[0].Weight,
		MaxEdgeWeight: result.Edges[0].Weight,
		IsComplete:    result.IsComplete,
	}

	totalWeight := 0
	for _, edge := range result.Edges {
		totalWeight += edge.Weight

		if edge.Weight < stats.MinEdgeWeight {
			stats.MinEdgeWeight = edge.Weight
		}
		if edge.Weight > stats.MaxEdgeWeight {
			stats.MaxEdgeWeight = edge.Weight
		}
	}

	stats.AvgEdgeWeight = float64(totalWeight) / float64(len(result.Edges))

	return stats
}

// PrintMSTAlgorithmGuide prints a guide for choosing MST algorithms
func (g *WeightedGraph) PrintMSTAlgorithmGuide() {
	fmt.Printf("=== MST Algorithm Selection Guide ===\n")
	fmt.Printf("Graph: %d vertices\n\n", g.vertices)

	edgeCount := 0
	for v := 0; v < g.vertices; v++ {
		edgeCount += len(g.adjList[v])
	}
	edgeCount /= 2

	density := float64(edgeCount) / float64(g.vertices*(g.vertices-1)/2) * 100

	fmt.Printf("Graph Properties:\n")
	fmt.Printf("- Vertices: %d\n", g.vertices)
	fmt.Printf("- Edges: %d\n", edgeCount)
	fmt.Printf("- Density: %.2f%%\n\n", density)

	fmt.Printf("Algorithm Recommendations:\n")

	if density < 25 {
		fmt.Printf("✓ SPARSE GRAPH (< 25%% density)\n")
		fmt.Printf("  1. Kruskal's Algorithm - Best for sparse graphs\n")
		fmt.Printf("     Time: O(E log E), Space: O(V + E)\n")
		fmt.Printf("  2. Prim's Algorithm - Alternative choice\n")
		fmt.Printf("     Time: O(E log V), Space: O(V)\n")
	} else if density < 75 {
		fmt.Printf("✓ MEDIUM DENSITY GRAPH (25-75%% density)\n")
		fmt.Printf("  1. Prim's Algorithm - Good balance\n")
		fmt.Printf("     Time: O(E log V), Space: O(V)\n")
		fmt.Printf("  2. Borůvka's Algorithm - Parallelizable\n")
		fmt.Printf("     Time: O(E log V), Space: O(V + E)\n")
	} else {
		fmt.Printf("✓ DENSE GRAPH (> 75%% density)\n")
		fmt.Printf("  1. Dense Prim's Algorithm - Optimal for dense graphs\n")
		fmt.Printf("     Time: O(V²), Space: O(V²)\n")
		fmt.Printf("  2. Standard Prim's - Alternative choice\n")
		fmt.Printf("     Time: O(E log V), Space: O(V)\n")
	}

	fmt.Printf("\nAvailable Algorithms:\n")
	fmt.Printf("- Kruskal: Edge-based, good for sparse graphs\n")
	fmt.Printf("- Prim: Vertex-based, good general purpose\n")
	fmt.Printf("- DensePrim: Matrix-based, optimal for dense graphs\n")
	fmt.Printf("- Borůvka: Component-based, parallelizable\n")
	fmt.Printf("- Validation: Verify MST correctness\n")
	fmt.Printf("- Comparison: Compare algorithm results\n")
}
