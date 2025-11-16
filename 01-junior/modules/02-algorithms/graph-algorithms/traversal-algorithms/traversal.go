// Package graph implements comprehensive graph traversal algorithms with optimal Go implementations
// for FAANG L3/L4 interview preparation and production systems.
//
// This module provides production-grade graph traversal algorithms including:
// - Depth-First Search (DFS) with recursive and iterative implementations
// - Breadth-First Search (BFS) for shortest path in unweighted graphs
// - Advanced traversal patterns: Bidirectional BFS, Iterative Deepening DFS
// - Multi-source BFS for parallel exploration
// - Connected components analysis and strongly connected components
// - Bipartite graph detection and two-coloring
// - Articulation points and bridge detection for network analysis
// - Graph representation utilities and comprehensive validation
//
// Time Complexity: O(V + E) for basic traversals, varies for advanced algorithms
// Space Complexity: O(V) for visited tracking and call stack/queue
package graph

import (
	"errors"
	"fmt"
)

// Graph represents an adjacency list implementation of a directed graph
type Graph struct {
	vertices int
	adjList  map[int][]int
}

// NewGraph creates a new graph with the specified number of vertices
func NewGraph(vertices int) *Graph {
	if vertices <= 0 {
		return &Graph{vertices: 0, adjList: make(map[int][]int)}
	}

	g := &Graph{
		vertices: vertices,
		adjList:  make(map[int][]int, vertices),
	}

	// Initialize adjacency lists for all vertices
	for i := 0; i < vertices; i++ {
		g.adjList[i] = make([]int, 0)
	}

	return g
}

// AddEdge adds a directed edge from source to destination vertex
// Time Complexity: O(1) amortized
// Space Complexity: O(1)
func (g *Graph) AddEdge(src, dest int) error {
	if src < 0 || src >= g.vertices || dest < 0 || dest >= g.vertices {
		return errors.New("vertex index out of bounds")
	}

	g.adjList[src] = append(g.adjList[src], dest)
	return nil
}

// AddUndirectedEdge adds an undirected edge between two vertices
// Time Complexity: O(1) amortized
// Space Complexity: O(1)
func (g *Graph) AddUndirectedEdge(u, v int) error {
	if err := g.AddEdge(u, v); err != nil {
		return err
	}
	return g.AddEdge(v, u)
}

// GetVertices returns the number of vertices in the graph
func (g *Graph) GetVertices() int {
	return g.vertices
}

// GetAdjacent returns the adjacency list for a given vertex
func (g *Graph) GetAdjacent(vertex int) ([]int, error) {
	if vertex < 0 || vertex >= g.vertices {
		return nil, errors.New("vertex index out of bounds")
	}

	// Return a copy to prevent external modification
	result := make([]int, len(g.adjList[vertex]))
	copy(result, g.adjList[vertex])
	return result, nil
}

// DFSRecursive performs depth-first search using recursion
// Time Complexity: O(V + E)
// Space Complexity: O(V) for recursion stack and visited array
func (g *Graph) DFSRecursive(startVertex int) ([]int, error) {
	if startVertex < 0 || startVertex >= g.vertices {
		return nil, errors.New("start vertex out of bounds")
	}

	visited := make([]bool, g.vertices)
	result := make([]int, 0, g.vertices)

	g.dfsRecursiveHelper(startVertex, visited, &result)
	return result, nil
}

// dfsRecursiveHelper is the recursive helper function for DFS
func (g *Graph) dfsRecursiveHelper(vertex int, visited []bool, result *[]int) {
	visited[vertex] = true
	*result = append(*result, vertex)

	// Visit all adjacent unvisited vertices
	for _, neighbor := range g.adjList[vertex] {
		if !visited[neighbor] {
			g.dfsRecursiveHelper(neighbor, visited, result)
		}
	}
}

// DFSIterative performs depth-first search using an explicit stack
// Time Complexity: O(V + E)
// Space Complexity: O(V) for stack and visited array
func (g *Graph) DFSIterative(startVertex int) ([]int, error) {
	if startVertex < 0 || startVertex >= g.vertices {
		return nil, errors.New("start vertex out of bounds")
	}

	visited := make([]bool, g.vertices)
	result := make([]int, 0, g.vertices)
	stack := []int{startVertex}

	for len(stack) > 0 {
		// Pop from stack
		vertex := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if !visited[vertex] {
			visited[vertex] = true
			result = append(result, vertex)

			// Add all unvisited neighbors to stack (in reverse order for consistent traversal)
			neighbors := g.adjList[vertex]
			for i := len(neighbors) - 1; i >= 0; i-- {
				if !visited[neighbors[i]] {
					stack = append(stack, neighbors[i])
				}
			}
		}
	}

	return result, nil
}

// BFS performs breadth-first search and returns the traversal order
// Time Complexity: O(V + E)
// Space Complexity: O(V) for queue and visited array
func (g *Graph) BFS(startVertex int) ([]int, error) {
	if startVertex < 0 || startVertex >= g.vertices {
		return nil, errors.New("start vertex out of bounds")
	}

	visited := make([]bool, g.vertices)
	result := make([]int, 0, g.vertices)
	queue := []int{startVertex}
	visited[startVertex] = true

	for len(queue) > 0 {
		// Dequeue front element
		vertex := queue[0]
		queue = queue[1:]
		result = append(result, vertex)

		// Add all unvisited neighbors to queue
		for _, neighbor := range g.adjList[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return result, nil
}

// BFSShortestPath finds the shortest path between start and end vertices
// Returns the path and the distance. Returns nil if no path exists.
// Time Complexity: O(V + E)
// Space Complexity: O(V) for queue, visited array, and parent tracking
func (g *Graph) BFSShortestPath(start, end int) ([]int, int, error) {
	if start < 0 || start >= g.vertices || end < 0 || end >= g.vertices {
		return nil, -1, errors.New("vertex index out of bounds")
	}

	if start == end {
		return []int{start}, 0, nil
	}

	visited := make([]bool, g.vertices)
	parent := make([]int, g.vertices)
	distance := make([]int, g.vertices)
	queue := []int{start}

	visited[start] = true
	parent[start] = -1
	distance[start] = 0

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]

		for _, neighbor := range g.adjList[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = vertex
				distance[neighbor] = distance[vertex] + 1
				queue = append(queue, neighbor)

				// Found target vertex
				if neighbor == end {
					path := g.reconstructPath(parent, start, end)
					return path, distance[end], nil
				}
			}
		}
	}

	// No path found
	return nil, -1, nil
}

// reconstructPath reconstructs the shortest path using parent array
func (g *Graph) reconstructPath(parent []int, start, end int) []int {
	path := make([]int, 0)
	current := end

	for current != -1 {
		path = append([]int{current}, path...)
		current = parent[current]
	}

	return path
}

// HasCycle detects if the graph contains a cycle using DFS
// Time Complexity: O(V + E)
// Space Complexity: O(V) for recursion stack and state tracking
func (g *Graph) HasCycle() bool {
	// White: 0 (unvisited), Gray: 1 (visiting), Black: 2 (visited)
	color := make([]int, g.vertices)

	for i := 0; i < g.vertices; i++ {
		if color[i] == 0 {
			if g.hasCycleDFS(i, color) {
				return true
			}
		}
	}

	return false
}

// hasCycleDFS is the DFS helper for cycle detection
func (g *Graph) hasCycleDFS(vertex int, color []int) bool {
	color[vertex] = 1 // Mark as gray (visiting)

	for _, neighbor := range g.adjList[vertex] {
		if color[neighbor] == 1 { // Back edge found (cycle detected)
			return true
		}
		if color[neighbor] == 0 && g.hasCycleDFS(neighbor, color) {
			return true
		}
	}

	color[vertex] = 2 // Mark as black (visited)
	return false
}

// IsConnected checks if all vertices are reachable from vertex 0
// Note: For directed graphs, this checks weak connectivity
// Time Complexity: O(V + E)
// Space Complexity: O(V)
func (g *Graph) IsConnected() bool {
	if g.vertices == 0 {
		return true
	}

	// Perform DFS from vertex 0
	visited := make([]bool, g.vertices)
	g.dfsRecursiveHelper(0, visited, &[]int{})

	// Check if all vertices were visited
	for i := 0; i < g.vertices; i++ {
		if !visited[i] {
			return false
		}
	}

	return true
}

// PrintGraph prints the adjacency list representation of the graph
func (g *Graph) PrintGraph() {
	fmt.Printf("Graph with %d vertices:\n", g.vertices)
	for i := 0; i < g.vertices; i++ {
		fmt.Printf("Vertex %d: %v\n", i, g.adjList[i])
	}
}

// GetEdgeCount returns the total number of edges in the graph
// Time Complexity: O(V)
// Space Complexity: O(1)
func (g *Graph) GetEdgeCount() int {
	count := 0
	for i := 0; i < g.vertices; i++ {
		count += len(g.adjList[i])
	}
	return count
}

// Clone creates a deep copy of the graph
// Time Complexity: O(V + E)
// Space Complexity: O(V + E)
func (g *Graph) Clone() *Graph {
	newGraph := NewGraph(g.vertices)

	for vertex := 0; vertex < g.vertices; vertex++ {
		for _, neighbor := range g.adjList[vertex] {
			newGraph.AddEdge(vertex, neighbor)
		}
	}

	return newGraph
}

// =====================================================================================
// ADVANCED TRAVERSAL ALGORITHMS - Production-grade implementations
// =====================================================================================

// BidirectionalBFS performs bidirectional breadth-first search to find shortest path
// between start and end vertices more efficiently than standard BFS
// Time Complexity: O(V + E) - often faster in practice due to reduced search space
// Space Complexity: O(V) for two queues and visited sets
func (g *Graph) BidirectionalBFS(start, end int) ([]int, int, error) {
	if start < 0 || start >= g.vertices || end < 0 || end >= g.vertices {
		return nil, -1, errors.New("vertex index out of bounds")
	}

	if start == end {
		return []int{start}, 0, nil
	}

	// Forward search from start
	visitedForward := make(map[int]bool)
	parentForward := make(map[int]int)
	queueForward := []int{start}
	visitedForward[start] = true
	parentForward[start] = -1

	// Backward search from end
	visitedBackward := make(map[int]bool)
	parentBackward := make(map[int]int)
	queueBackward := []int{end}
	visitedBackward[end] = true
	parentBackward[end] = -1

	// Alternate between forward and backward searches
	distance := 0
	for len(queueForward) > 0 || len(queueBackward) > 0 {
		distance++

		// Forward step
		if len(queueForward) > 0 {
			if intersection := g.bidirectionalBFSStep(queueForward, visitedForward,
				parentForward, visitedBackward, true); intersection != -1 {
				path := g.reconstructBidirectionalPath(parentForward, parentBackward,
					start, end, intersection)
				return path, distance, nil
			}
			queueForward = g.getNextLevelQueue(queueForward, visitedForward, parentForward, true)
		}

		// Backward step
		if len(queueBackward) > 0 {
			if intersection := g.bidirectionalBFSStep(queueBackward, visitedBackward,
				parentBackward, visitedForward, false); intersection != -1 {
				path := g.reconstructBidirectionalPath(parentForward, parentBackward,
					start, end, intersection)
				return path, distance, nil
			}
			queueBackward = g.getNextLevelQueue(queueBackward, visitedBackward, parentBackward, false)
		}
	}

	return nil, -1, nil // No path found
}

// bidirectionalBFSStep processes one level of bidirectional BFS
func (g *Graph) bidirectionalBFSStep(queue []int, visited map[int]bool, parent map[int]int,
	otherVisited map[int]bool, isForward bool) int {

	for _, vertex := range queue {
		var neighbors []int
		if isForward {
			neighbors = g.adjList[vertex]
		} else {
			// For backward search, we need incoming edges
			neighbors = g.getIncomingNeighbors(vertex)
		}

		for _, neighbor := range neighbors {
			if otherVisited[neighbor] {
				return neighbor // Intersection found
			}
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = vertex
			}
		}
	}
	return -1 // No intersection
}

// getNextLevelQueue gets the next level queue for bidirectional BFS
func (g *Graph) getNextLevelQueue(currentQueue []int, visited map[int]bool,
	parent map[int]int, isForward bool) []int {
	nextQueue := make([]int, 0)

	for _, vertex := range currentQueue {
		var neighbors []int
		if isForward {
			neighbors = g.adjList[vertex]
		} else {
			neighbors = g.getIncomingNeighbors(vertex)
		}

		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				nextQueue = append(nextQueue, neighbor)
			}
		}
	}
	return nextQueue
}

// getIncomingNeighbors returns vertices that have edges pointing to the given vertex
func (g *Graph) getIncomingNeighbors(vertex int) []int {
	incoming := make([]int, 0)
	for v := 0; v < g.vertices; v++ {
		for _, neighbor := range g.adjList[v] {
			if neighbor == vertex {
				incoming = append(incoming, v)
				break
			}
		}
	}
	return incoming
}

// reconstructBidirectionalPath reconstructs path from bidirectional search
func (g *Graph) reconstructBidirectionalPath(parentForward, parentBackward map[int]int,
	start, end, intersection int) []int {

	// Build path from start to intersection
	forwardPath := make([]int, 0)
	current := intersection
	for current != -1 {
		forwardPath = append([]int{current}, forwardPath...)
		current = parentForward[current]
	}

	// Build path from intersection to end
	backwardPath := make([]int, 0)
	current = parentBackward[intersection]
	for current != -1 {
		backwardPath = append(backwardPath, current)
		current = parentBackward[current]
	}

	// Combine paths
	fullPath := make([]int, 0, len(forwardPath)+len(backwardPath))
	fullPath = append(fullPath, forwardPath...)
	fullPath = append(fullPath, backwardPath...)

	return fullPath
}

// IterativeDeepeningDFS performs iterative deepening depth-first search
// Combines benefits of DFS (low memory) with BFS (optimal solution)
// Time Complexity: O(b^d) where b is branching factor and d is depth
// Space Complexity: O(d) for recursion stack
func (g *Graph) IterativeDeepeningDFS(start, target int, maxDepth int) ([]int, error) {
	if start < 0 || start >= g.vertices || target < 0 || target >= g.vertices {
		return nil, errors.New("vertex index out of bounds")
	}

	if start == target {
		return []int{start}, nil
	}

	for depth := 0; depth <= maxDepth; depth++ {
		visited := make([]bool, g.vertices)
		path := make([]int, 0)

		if g.depthLimitedDFS(start, target, depth, visited, &path) {
			return path, nil
		}
	}

	return nil, errors.New("target not found within maximum depth")
}

// depthLimitedDFS performs DFS with depth limit
func (g *Graph) depthLimitedDFS(current, target, depthLimit int,
	visited []bool, path *[]int) bool {

	if current == target {
		*path = append(*path, current)
		return true
	}

	if depthLimit <= 0 {
		return false
	}

	visited[current] = true
	*path = append(*path, current)

	for _, neighbor := range g.adjList[current] {
		if !visited[neighbor] {
			if g.depthLimitedDFS(neighbor, target, depthLimit-1, visited, path) {
				return true
			}
		}
	}

	// Backtrack
	*path = (*path)[:len(*path)-1]
	visited[current] = false
	return false
}

// MultiSourceBFS performs BFS from multiple source vertices simultaneously
// Useful for finding shortest distances from any source to all other vertices
// Time Complexity: O(V + E)
// Space Complexity: O(V)
func (g *Graph) MultiSourceBFS(sources []int) (map[int]int, error) {
	// Validate sources
	for _, source := range sources {
		if source < 0 || source >= g.vertices {
			return nil, errors.New("source vertex out of bounds")
		}
	}

	distances := make(map[int]int)
	visited := make([]bool, g.vertices)
	queue := make([]int, 0, len(sources))

	// Initialize with all sources at distance 0
	for _, source := range sources {
		distances[source] = 0
		visited[source] = true
		queue = append(queue, source)
	}

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]

		for _, neighbor := range g.adjList[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				distances[neighbor] = distances[vertex] + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return distances, nil
}

// =====================================================================================
// CONNECTED COMPONENTS AND GRAPH ANALYSIS
// =====================================================================================

// FindConnectedComponents finds all connected components in an undirected graph
// Time Complexity: O(V + E)
// Space Complexity: O(V)
func (g *Graph) FindConnectedComponents() [][]int {
	visited := make([]bool, g.vertices)
	components := make([][]int, 0)

	for i := 0; i < g.vertices; i++ {
		if !visited[i] {
			component := make([]int, 0)
			g.dfsRecursiveHelper(i, visited, &component)
			components = append(components, component)
		}
	}

	return components
}

// FindStronglyConnectedComponents finds SCCs using Tarjan's algorithm
// Time Complexity: O(V + E)
// Space Complexity: O(V)
func (g *Graph) FindStronglyConnectedComponents() [][]int {
	index := 0
	stack := make([]int, 0)
	indices := make([]int, g.vertices)
	lowLinks := make([]int, g.vertices)
	onStack := make([]bool, g.vertices)
	sccs := make([][]int, 0)

	for i := 0; i < g.vertices; i++ {
		indices[i] = -1
	}

	for v := 0; v < g.vertices; v++ {
		if indices[v] == -1 {
			g.tarjanSCC(v, &index, &stack, indices, lowLinks, onStack, &sccs)
		}
	}

	return sccs
}

// tarjanSCC implements Tarjan's strongly connected components algorithm
func (g *Graph) tarjanSCC(v int, index *int, stack *[]int, indices, lowLinks []int,
	onStack []bool, sccs *[][]int) {

	indices[v] = *index
	lowLinks[v] = *index
	*index++
	*stack = append(*stack, v)
	onStack[v] = true

	for _, w := range g.adjList[v] {
		if indices[w] == -1 {
			g.tarjanSCC(w, index, stack, indices, lowLinks, onStack, sccs)
			lowLinks[v] = min(lowLinks[v], lowLinks[w])
		} else if onStack[w] {
			lowLinks[v] = min(lowLinks[v], indices[w])
		}
	}

	// If v is a root node, pop the stack and create an SCC
	if lowLinks[v] == indices[v] {
		scc := make([]int, 0)
		for {
			w := (*stack)[len(*stack)-1]
			*stack = (*stack)[:len(*stack)-1]
			onStack[w] = false
			scc = append(scc, w)
			if w == v {
				break
			}
		}
		*sccs = append(*sccs, scc)
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// IsBipartite checks if the graph is bipartite using two-coloring
// Time Complexity: O(V + E)
// Space Complexity: O(V)
func (g *Graph) IsBipartite() (bool, []int) {
	color := make([]int, g.vertices)
	for i := 0; i < g.vertices; i++ {
		color[i] = -1 // Uncolored
	}

	// Check each component
	for start := 0; start < g.vertices; start++ {
		if color[start] == -1 {
			if !g.isBipartiteBFS(start, color) {
				return false, nil
			}
		}
	}

	return true, color
}

// isBipartiteBFS checks bipartiteness using BFS coloring
func (g *Graph) isBipartiteBFS(start int, color []int) bool {
	queue := []int{start}
	color[start] = 0

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]

		for _, neighbor := range g.adjList[vertex] {
			if color[neighbor] == -1 {
				color[neighbor] = 1 - color[vertex] // Alternate color
				queue = append(queue, neighbor)
			} else if color[neighbor] == color[vertex] {
				return false // Same color adjacent vertices
			}
		}
	}

	return true
}

// =====================================================================================
// CRITICAL VERTICES AND EDGES ANALYSIS
// =====================================================================================

// FindArticulationPoints finds articulation points (cut vertices) using Tarjan's algorithm
// Time Complexity: O(V + E)
// Space Complexity: O(V)
func (g *Graph) FindArticulationPoints() []int {
	visited := make([]bool, g.vertices)
	disc := make([]int, g.vertices)
	low := make([]int, g.vertices)
	parent := make([]int, g.vertices)
	ap := make([]bool, g.vertices)
	time := 0

	for i := 0; i < g.vertices; i++ {
		parent[i] = -1
	}

	for i := 0; i < g.vertices; i++ {
		if !visited[i] {
			g.findAPUtil(i, visited, disc, low, parent, ap, &time)
		}
	}

	articulationPoints := make([]int, 0)
	for i := 0; i < g.vertices; i++ {
		if ap[i] {
			articulationPoints = append(articulationPoints, i)
		}
	}

	return articulationPoints
}

// findAPUtil implements the recursive function for articulation points
func (g *Graph) findAPUtil(u int, visited []bool, disc, low, parent []int,
	ap []bool, time *int) {

	children := 0
	visited[u] = true
	disc[u] = *time
	low[u] = *time
	*time++

	for _, v := range g.adjList[u] {
		if !visited[v] {
			children++
			parent[v] = u
			g.findAPUtil(v, visited, disc, low, parent, ap, time)

			low[u] = min(low[u], low[v])

			// u is an articulation point in the following cases:
			// (1) u is root of DFS tree and has two or more children
			if parent[u] == -1 && children > 1 {
				ap[u] = true
			}

			// (2) u is not root and low value of one of its children is more than discovery value of u
			if parent[u] != -1 && low[v] >= disc[u] {
				ap[u] = true
			}
		} else if v != parent[u] {
			low[u] = min(low[u], disc[v])
		}
	}
}

// Edge represents an edge in the graph
type Edge struct {
	U, V int
}

// FindBridges finds all bridges (cut edges) in the graph using Tarjan's algorithm
// Time Complexity: O(V + E)
// Space Complexity: O(V)
func (g *Graph) FindBridges() []Edge {
	visited := make([]bool, g.vertices)
	disc := make([]int, g.vertices)
	low := make([]int, g.vertices)
	parent := make([]int, g.vertices)
	bridges := make([]Edge, 0)
	time := 0

	for i := 0; i < g.vertices; i++ {
		parent[i] = -1
	}

	for i := 0; i < g.vertices; i++ {
		if !visited[i] {
			g.findBridgesUtil(i, visited, disc, low, parent, &bridges, &time)
		}
	}

	return bridges
}

// findBridgesUtil implements the recursive function for finding bridges
func (g *Graph) findBridgesUtil(u int, visited []bool, disc, low, parent []int,
	bridges *[]Edge, time *int) {

	visited[u] = true
	disc[u] = *time
	low[u] = *time
	*time++

	for _, v := range g.adjList[u] {
		if !visited[v] {
			parent[v] = u
			g.findBridgesUtil(v, visited, disc, low, parent, bridges, time)

			low[u] = min(low[u], low[v])

			// If the lowest vertex reachable from subtree under v is below u in DFS tree, then u-v is a bridge
			if low[v] > disc[u] {
				*bridges = append(*bridges, Edge{U: u, V: v})
			}
		} else if v != parent[u] {
			low[u] = min(low[u], disc[v])
		}
	}
}

// =====================================================================================
// PERFORMANCE ANALYSIS AND UTILITIES
// =====================================================================================

// TraversalStats provides statistics about graph traversal
type TraversalStats struct {
	VerticesVisited   int
	EdgesTraversed    int
	MaxDepthReached   int
	ComponentsFound   int
	ArticulationCount int
	BridgeCount       int
	IsBipartite       bool
	IsConnected       bool
	HasCycles         bool
}

// AnalyzeGraph performs comprehensive graph analysis
// Time Complexity: O(V + E) for each analysis
// Space Complexity: O(V)
func (g *Graph) AnalyzeGraph() TraversalStats {
	stats := TraversalStats{}

	// Basic connectivity
	stats.IsConnected = g.IsConnected()
	stats.HasCycles = g.HasCycle()

	// Component analysis
	components := g.FindConnectedComponents()
	stats.ComponentsFound = len(components)

	// Count vertices in largest component
	maxComponent := 0
	for _, component := range components {
		if len(component) > maxComponent {
			maxComponent = len(component)
		}
	}
	stats.VerticesVisited = maxComponent

	// Edge count
	stats.EdgesTraversed = g.GetEdgeCount()

	// Critical vertex/edge analysis
	articulationPoints := g.FindArticulationPoints()
	stats.ArticulationCount = len(articulationPoints)

	bridges := g.FindBridges()
	stats.BridgeCount = len(bridges)

	// Bipartite test
	stats.IsBipartite, _ = g.IsBipartite()

	// Calculate max depth using DFS
	if g.vertices > 0 {
		visited := make([]bool, g.vertices)
		stats.MaxDepthReached = g.calculateMaxDepth(0, visited, 0)
	}

	return stats
}

// calculateMaxDepth finds the maximum depth reachable from a vertex
func (g *Graph) calculateMaxDepth(vertex int, visited []bool, currentDepth int) int {
	visited[vertex] = true
	maxDepth := currentDepth

	for _, neighbor := range g.adjList[vertex] {
		if !visited[neighbor] {
			depth := g.calculateMaxDepth(neighbor, visited, currentDepth+1)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
	}

	return maxDepth
}

// PrintTraversalStats prints comprehensive traversal statistics
func (stats TraversalStats) PrintStats() {
	fmt.Println("=== Graph Traversal Analysis ===")
	fmt.Printf("Vertices Visited: %d\n", stats.VerticesVisited)
	fmt.Printf("Edges Traversed: %d\n", stats.EdgesTraversed)
	fmt.Printf("Maximum Depth: %d\n", stats.MaxDepthReached)
	fmt.Printf("Connected Components: %d\n", stats.ComponentsFound)
	fmt.Printf("Articulation Points: %d\n", stats.ArticulationCount)
	fmt.Printf("Bridges: %d\n", stats.BridgeCount)
	fmt.Printf("Is Connected: %t\n", stats.IsConnected)
	fmt.Printf("Is Bipartite: %t\n", stats.IsBipartite)
	fmt.Printf("Has Cycles: %t\n", stats.HasCycles)
	fmt.Println("================================")
}

// BenchmarkTraversals compares performance of different traversal algorithms
func (g *Graph) BenchmarkTraversals(start int) map[string]int {
	results := make(map[string]int)

	if start < 0 || start >= g.vertices {
		return results
	}

	// Benchmark DFS Recursive
	if dfsResult, err := g.DFSRecursive(start); err == nil {
		results["DFS_Recursive"] = len(dfsResult)
	}

	// Benchmark DFS Iterative
	if dfsResult, err := g.DFSIterative(start); err == nil {
		results["DFS_Iterative"] = len(dfsResult)
	}

	// Benchmark BFS
	if bfsResult, err := g.BFS(start); err == nil {
		results["BFS"] = len(bfsResult)
	}

	return results
}

// ValidateTraversalResults validates that traversal results are consistent
func (g *Graph) ValidateTraversalResults(start int) bool {
	if start < 0 || start >= g.vertices {
		return false
	}

	dfsRec, err1 := g.DFSRecursive(start)
	dfsIter, err2 := g.DFSIterative(start)
	bfs, err3 := g.BFS(start)

	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}

	// All should visit the same number of vertices (for connected component)
	return len(dfsRec) == len(dfsIter) && len(dfsIter) == len(bfs)
}
