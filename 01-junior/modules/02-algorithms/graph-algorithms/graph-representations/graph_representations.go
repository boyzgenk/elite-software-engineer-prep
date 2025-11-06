// Graph Representations in Go
// Comprehensive implementation of different graph data structures and representations

package main

import (
	"fmt"
	"sort"
)

// ========================================================================
// ADJACENCY MATRIX REPRESENTATION
// ========================================================================

// AdjacencyMatrix represents a graph using a 2D array
// Best for: Dense graphs, quick edge existence checks O(1)
// Memory: O(V²) where V is number of vertices
type AdjacencyMatrix struct {
	vertices int
	matrix   [][]int // For weighted graphs, use weights; for unweighted, use 0/1
	directed bool
}

// NewAdjacencyMatrix creates a new adjacency matrix
func NewAdjacencyMatrix(vertices int, directed bool) *AdjacencyMatrix {
	matrix := make([][]int, vertices)
	for i := range matrix {
		matrix[i] = make([]int, vertices)
	}

	return &AdjacencyMatrix{
		vertices: vertices,
		matrix:   matrix,
		directed: directed,
	}
}

// AddEdge adds an edge between two vertices
func (am *AdjacencyMatrix) AddEdge(from, to, weight int) {
	if from < 0 || from >= am.vertices || to < 0 || to >= am.vertices {
		return
	}

	am.matrix[from][to] = weight
	if !am.directed {
		am.matrix[to][from] = weight
	}
}

// RemoveEdge removes an edge between two vertices
func (am *AdjacencyMatrix) RemoveEdge(from, to int) {
	if from < 0 || from >= am.vertices || to < 0 || to >= am.vertices {
		return
	}

	am.matrix[from][to] = 0
	if !am.directed {
		am.matrix[to][from] = 0
	}
}

// HasEdge checks if an edge exists between two vertices - O(1)
func (am *AdjacencyMatrix) HasEdge(from, to int) bool {
	if from < 0 || from >= am.vertices || to < 0 || to >= am.vertices {
		return false
	}
	return am.matrix[from][to] != 0
}

// GetWeight returns the weight of edge between two vertices
func (am *AdjacencyMatrix) GetWeight(from, to int) int {
	if from < 0 || from >= am.vertices || to < 0 || to >= am.vertices {
		return 0
	}
	return am.matrix[from][to]
}

// GetNeighbors returns all neighbors of a vertex - O(V)
func (am *AdjacencyMatrix) GetNeighbors(vertex int) []int {
	if vertex < 0 || vertex >= am.vertices {
		return nil
	}

	var neighbors []int
	for i := 0; i < am.vertices; i++ {
		if am.matrix[vertex][i] != 0 {
			neighbors = append(neighbors, i)
		}
	}
	return neighbors
}

// Print displays the adjacency matrix
func (am *AdjacencyMatrix) Print() {
	fmt.Println("Adjacency Matrix:")
	fmt.Print("   ")
	for i := 0; i < am.vertices; i++ {
		fmt.Printf("%3d", i)
	}
	fmt.Println()

	for i := 0; i < am.vertices; i++ {
		fmt.Printf("%2d ", i)
		for j := 0; j < am.vertices; j++ {
			fmt.Printf("%3d", am.matrix[i][j])
		}
		fmt.Println()
	}
}

// ========================================================================
// ADJACENCY LIST REPRESENTATION
// ========================================================================

// Edge represents a weighted edge to a neighbor
type Edge struct {
	To     int
	Weight int
}

// AdjacencyList represents a graph using adjacency lists
// Best for: Sparse graphs, efficient iteration over neighbors
// Memory: O(V + E) where V is vertices, E is edges
type AdjacencyList struct {
	vertices int
	adjList  [][]Edge
	directed bool
}

// NewAdjacencyList creates a new adjacency list
func NewAdjacencyList(vertices int, directed bool) *AdjacencyList {
	adjList := make([][]Edge, vertices)
	for i := range adjList {
		adjList[i] = make([]Edge, 0)
	}

	return &AdjacencyList{
		vertices: vertices,
		adjList:  adjList,
		directed: directed,
	}
}

// AddEdge adds an edge between two vertices
func (al *AdjacencyList) AddEdge(from, to, weight int) {
	if from < 0 || from >= al.vertices || to < 0 || to >= al.vertices {
		return
	}

	al.adjList[from] = append(al.adjList[from], Edge{To: to, Weight: weight})
	if !al.directed {
		al.adjList[to] = append(al.adjList[to], Edge{To: from, Weight: weight})
	}
}

// RemoveEdge removes an edge between two vertices
func (al *AdjacencyList) RemoveEdge(from, to int) {
	if from < 0 || from >= al.vertices || to < 0 || to >= al.vertices {
		return
	}

	// Remove edge from -> to
	for i, edge := range al.adjList[from] {
		if edge.To == to {
			al.adjList[from] = append(al.adjList[from][:i], al.adjList[from][i+1:]...)
			break
		}
	}

	// If undirected, remove edge to -> from
	if !al.directed {
		for i, edge := range al.adjList[to] {
			if edge.To == from {
				al.adjList[to] = append(al.adjList[to][:i], al.adjList[to][i+1:]...)
				break
			}
		}
	}
}

// HasEdge checks if an edge exists between two vertices - O(degree)
func (al *AdjacencyList) HasEdge(from, to int) bool {
	if from < 0 || from >= al.vertices || to < 0 || to >= al.vertices {
		return false
	}

	for _, edge := range al.adjList[from] {
		if edge.To == to {
			return true
		}
	}
	return false
}

// GetWeight returns the weight of edge between two vertices
func (al *AdjacencyList) GetWeight(from, to int) int {
	if from < 0 || from >= al.vertices || to < 0 || to >= al.vertices {
		return 0
	}

	for _, edge := range al.adjList[from] {
		if edge.To == to {
			return edge.Weight
		}
	}
	return 0
}

// GetNeighbors returns all neighbors of a vertex - O(degree)
func (al *AdjacencyList) GetNeighbors(vertex int) []Edge {
	if vertex < 0 || vertex >= al.vertices {
		return nil
	}
	return al.adjList[vertex]
}

// GetVertexCount returns the number of vertices
func (al *AdjacencyList) GetVertexCount() int {
	return al.vertices
}

// GetEdgeCount returns the total number of edges
func (al *AdjacencyList) GetEdgeCount() int {
	count := 0
	for _, neighbors := range al.adjList {
		count += len(neighbors)
	}

	if !al.directed {
		count /= 2 // Each undirected edge is counted twice
	}
	return count
}

// Print displays the adjacency list
func (al *AdjacencyList) Print() {
	fmt.Println("Adjacency List:")
	for i := 0; i < al.vertices; i++ {
		fmt.Printf("Vertex %d: ", i)
		for j, edge := range al.adjList[i] {
			if j > 0 {
				fmt.Print(" -> ")
			}
			if edge.Weight == 1 {
				fmt.Printf("%d", edge.To)
			} else {
				fmt.Printf("%d(w:%d)", edge.To, edge.Weight)
			}
		}
		fmt.Println()
	}
}

// ========================================================================
// EDGE LIST REPRESENTATION
// ========================================================================

// WeightedEdge represents an edge with source, destination and weight
type WeightedEdge struct {
	From   int
	To     int
	Weight int
}

// EdgeList represents a graph as a list of edges
// Best for: Algorithms that process edges (Kruskal's MST, sorting edges)
// Memory: O(E) where E is number of edges
type EdgeList struct {
	vertices int
	edges    []WeightedEdge
	directed bool
}

// NewEdgeList creates a new edge list
func NewEdgeList(vertices int, directed bool) *EdgeList {
	return &EdgeList{
		vertices: vertices,
		edges:    make([]WeightedEdge, 0),
		directed: directed,
	}
}

// AddEdge adds an edge to the edge list
func (el *EdgeList) AddEdge(from, to, weight int) {
	if from < 0 || from >= el.vertices || to < 0 || to >= el.vertices {
		return
	}

	el.edges = append(el.edges, WeightedEdge{From: from, To: to, Weight: weight})
}

// GetEdges returns all edges
func (el *EdgeList) GetEdges() []WeightedEdge {
	return el.edges
}

// GetVertexCount returns the number of vertices
func (el *EdgeList) GetVertexCount() int {
	return el.vertices
}

// GetEdgeCount returns the number of edges
func (el *EdgeList) GetEdgeCount() int {
	return len(el.edges)
}

// SortByWeight sorts edges by weight (useful for MST algorithms)
func (el *EdgeList) SortByWeight() {
	sort.Slice(el.edges, func(i, j int) bool {
		return el.edges[i].Weight < el.edges[j].Weight
	})
}

// Print displays the edge list
func (el *EdgeList) Print() {
	fmt.Println("Edge List:")
	for i, edge := range el.edges {
		fmt.Printf("Edge %d: %d -> %d (weight: %d)\n", i, edge.From, edge.To, edge.Weight)
	}
}

// ========================================================================
// COMPRESSED SPARSE ROW (CSR) REPRESENTATION
// ========================================================================

// CSRGraph represents a graph using Compressed Sparse Row format
// Best for: Large sparse graphs, memory-efficient representation
// Memory: O(V + E) with better cache locality than adjacency list
type CSRGraph struct {
	vertices int
	rowPtr   []int // Row pointers - indices into colInd/values
	colInd   []int // Column indices (destination vertices)
	values   []int // Edge weights
	directed bool
}

// NewCSRGraph creates a new CSR graph from an adjacency list
func NewCSRGraph(al *AdjacencyList) *CSRGraph {
	vertices := al.vertices
	rowPtr := make([]int, vertices+1)
	colInd := make([]int, 0)
	values := make([]int, 0)

	// Build CSR format
	rowPtr[0] = 0
	for i := 0; i < vertices; i++ {
		neighbors := al.GetNeighbors(i)
		for _, edge := range neighbors {
			colInd = append(colInd, edge.To)
			values = append(values, edge.Weight)
		}
		rowPtr[i+1] = len(colInd)
	}

	return &CSRGraph{
		vertices: vertices,
		rowPtr:   rowPtr,
		colInd:   colInd,
		values:   values,
		directed: al.directed,
	}
}

// GetNeighbors returns neighbors of a vertex - O(degree) with better cache locality
func (csr *CSRGraph) GetNeighbors(vertex int) []Edge {
	if vertex < 0 || vertex >= csr.vertices {
		return nil
	}

	start := csr.rowPtr[vertex]
	end := csr.rowPtr[vertex+1]

	neighbors := make([]Edge, end-start)
	for i := start; i < end; i++ {
		neighbors[i-start] = Edge{
			To:     csr.colInd[i],
			Weight: csr.values[i],
		}
	}

	return neighbors
}

// Print displays the CSR representation
func (csr *CSRGraph) Print() {
	fmt.Println("CSR Graph Representation:")
	fmt.Printf("Row Pointers: %v\n", csr.rowPtr)
	fmt.Printf("Column Indices: %v\n", csr.colInd)
	fmt.Printf("Values: %v\n", csr.values)

	fmt.Println("\nNeighbors per vertex:")
	for i := 0; i < csr.vertices; i++ {
		neighbors := csr.GetNeighbors(i)
		fmt.Printf("Vertex %d: ", i)
		for j, edge := range neighbors {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%d(w:%d)", edge.To, edge.Weight)
		}
		fmt.Println()
	}
}

// ========================================================================
// GRAPH REPRESENTATION COMPARISON
// ========================================================================

type GraphRepresentationComparison struct{}

func (grc *GraphRepresentationComparison) CompareRepresentations() {
	fmt.Println("🔍 GRAPH REPRESENTATION COMPARISON")
	fmt.Println("====================================")

	fmt.Println("\n=== ADJACENCY MATRIX ===")
	fmt.Println("✅ Pros:")
	fmt.Println("  • O(1) edge existence check")
	fmt.Println("  • O(1) edge weight retrieval")
	fmt.Println("  • Simple implementation")
	fmt.Println("  • Good for dense graphs")

	fmt.Println("❌ Cons:")
	fmt.Println("  • O(V²) space complexity")
	fmt.Println("  • O(V) to iterate neighbors")
	fmt.Println("  • Wasteful for sparse graphs")

	fmt.Println("\n=== ADJACENCY LIST ===")
	fmt.Println("✅ Pros:")
	fmt.Println("  • O(V + E) space complexity")
	fmt.Println("  • O(degree) to iterate neighbors")
	fmt.Println("  • Efficient for sparse graphs")
	fmt.Println("  • Dynamic size")

	fmt.Println("❌ Cons:")
	fmt.Println("  • O(degree) edge existence check")
	fmt.Println("  • More complex implementation")
	fmt.Println("  • Pointer overhead")

	fmt.Println("\n=== EDGE LIST ===")
	fmt.Println("✅ Pros:")
	fmt.Println("  • O(E) space complexity")
	fmt.Println("  • Simple for edge-based algorithms")
	fmt.Println("  • Easy to sort edges")
	fmt.Println("  • No vertex overhead")

	fmt.Println("❌ Cons:")
	fmt.Println("  • O(E) to find neighbors")
	fmt.Println("  • O(E) edge existence check")
	fmt.Println("  • Not suitable for traversals")

	fmt.Println("\n=== CSR (Compressed Sparse Row) ===")
	fmt.Println("✅ Pros:")
	fmt.Println("  • O(V + E) space like adj list")
	fmt.Println("  • Better cache locality")
	fmt.Println("  • Memory efficient")
	fmt.Println("  • Good for large sparse graphs")

	fmt.Println("❌ Cons:")
	fmt.Println("  • Static structure")
	fmt.Println("  • Complex to modify")
	fmt.Println("  • Setup overhead")
}

func (grc *GraphRepresentationComparison) ShowSpaceComplexity() {
	fmt.Println("\n📊 SPACE COMPLEXITY COMPARISON")
	fmt.Println("===============================")

	vertices := []int{10, 100, 1000}
	densities := []float64{0.1, 0.5, 0.9} // 10%, 50%, 90% of possible edges

	for _, v := range vertices {
		maxEdges := v * (v - 1) / 2 // Undirected graph
		fmt.Printf("\n--- %d vertices (max %d edges) ---\n", v, maxEdges)

		for _, density := range densities {
			edges := int(float64(maxEdges) * density)

			matrixSpace := v * v * 4          // 4 bytes per int
			listSpace := (v + 2*edges) * 8    // Rough estimate with pointers
			edgeSpace := edges * 12           // 3 ints per edge
			csrSpace := (v + 1 + 2*edges) * 4 // Row ptrs + col indices + values

			fmt.Printf("Density %.1f%% (%d edges):\n", density*100, edges)
			fmt.Printf("  Matrix: %6d bytes\n", matrixSpace)
			fmt.Printf("  List:   %6d bytes\n", listSpace)
			fmt.Printf("  Edge:   %6d bytes\n", edgeSpace)
			fmt.Printf("  CSR:    %6d bytes\n", csrSpace)
		}
	}
}

// ========================================================================
// DEMONSTRATION AND TESTING
// ========================================================================

func main() {
	fmt.Println("🗺️  GRAPH REPRESENTATIONS DEMO")
	fmt.Println("==============================")

	// Create a sample graph: 0-1-2
	//                        |   |
	//                        3---+

	fmt.Println("\n=== ADJACENCY MATRIX DEMO ===")
	matrix := NewAdjacencyMatrix(4, false)
	matrix.AddEdge(0, 1, 5)
	matrix.AddEdge(1, 2, 3)
	matrix.AddEdge(0, 3, 2)
	matrix.AddEdge(2, 3, 1)
	matrix.Print()

	fmt.Printf("Edge 0->1 exists: %t, weight: %d\n", matrix.HasEdge(0, 1), matrix.GetWeight(0, 1))
	fmt.Printf("Neighbors of vertex 0: %v\n", matrix.GetNeighbors(0))

	fmt.Println("\n=== ADJACENCY LIST DEMO ===")
	adjList := NewAdjacencyList(4, false)
	adjList.AddEdge(0, 1, 5)
	adjList.AddEdge(1, 2, 3)
	adjList.AddEdge(0, 3, 2)
	adjList.AddEdge(2, 3, 1)
	adjList.Print()

	fmt.Printf("Edge 0->1 exists: %t, weight: %d\n", adjList.HasEdge(0, 1), adjList.GetWeight(0, 1))
	fmt.Printf("Total edges: %d\n", adjList.GetEdgeCount())

	fmt.Println("\n=== EDGE LIST DEMO ===")
	edgeList := NewEdgeList(4, false)
	edgeList.AddEdge(0, 1, 5)
	edgeList.AddEdge(1, 2, 3)
	edgeList.AddEdge(0, 3, 2)
	edgeList.AddEdge(2, 3, 1)
	edgeList.Print()

	fmt.Println("\nSorted by weight:")
	edgeList.SortByWeight()
	edgeList.Print()

	fmt.Println("\n=== CSR DEMO ===")
	csrGraph := NewCSRGraph(adjList)
	csrGraph.Print()

	fmt.Println("\n=== COMPARISON ANALYSIS ===")
	comparison := &GraphRepresentationComparison{}
	comparison.CompareRepresentations()
	comparison.ShowSpaceComplexity()

	fmt.Println("\n🎯 Choose representation based on:")
	fmt.Println("• Dense graphs (>50% edges): Adjacency Matrix")
	fmt.Println("• Sparse graphs: Adjacency List or CSR")
	fmt.Println("• Edge-based algorithms: Edge List")
	fmt.Println("• Large static graphs: CSR")
	fmt.Println("• Dynamic graphs: Adjacency List")
}
