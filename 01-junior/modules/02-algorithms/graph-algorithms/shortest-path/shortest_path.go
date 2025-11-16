// Package graph implements comprehensive shortest path algorithms for weighted graphs
// including Dijkstra, Bellman-Ford, Floyd-Warshall, A*, Johnson's, and optimized variants
// for FAANG L3/L4 interview preparation.
package graph

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
	"sort"
)

// Edge represents a weighted edge in the graph
type Edge struct {
	To     int
	Weight int
}

// WeightedGraph represents a weighted directed graph using adjacency list
type WeightedGraph struct {
	vertices int
	adjList  map[int][]Edge
}

// NewWeightedGraph creates a new weighted graph with the specified number of vertices
func NewWeightedGraph(vertices int) *WeightedGraph {
	if vertices <= 0 {
		return &WeightedGraph{vertices: 0, adjList: make(map[int][]Edge)}
	}

	g := &WeightedGraph{
		vertices: vertices,
		adjList:  make(map[int][]Edge, vertices),
	}

	// Initialize adjacency lists for all vertices
	for i := 0; i < vertices; i++ {
		g.adjList[i] = make([]Edge, 0)
	}

	return g
}

// AddWeightedEdge adds a weighted directed edge from source to destination
// Time Complexity: O(1) amortized
// Space Complexity: O(1)
func (g *WeightedGraph) AddWeightedEdge(src, dest, weight int) error {
	if src < 0 || src >= g.vertices || dest < 0 || dest >= g.vertices {
		return errors.New("vertex index out of bounds")
	}

	edge := Edge{To: dest, Weight: weight}
	g.adjList[src] = append(g.adjList[src], edge)
	return nil
}

// AddUndirectedWeightedEdge adds an undirected weighted edge between two vertices
func (g *WeightedGraph) AddUndirectedWeightedEdge(u, v, weight int) error {
	if err := g.AddWeightedEdge(u, v, weight); err != nil {
		return err
	}
	return g.AddWeightedEdge(v, u, weight)
}

// PriorityQueueItem represents an item in the priority queue for Dijkstra's algorithm
type PriorityQueueItem struct {
	vertex   int
	distance int
	index    int // Index in the heap (needed for heap.Interface)
}

// PriorityQueue implements a min-heap for Dijkstra's algorithm
type PriorityQueue []*PriorityQueueItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PriorityQueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Dijkstra finds the shortest path from source to all other vertices
// Returns distances array and parent array for path reconstruction
// Time Complexity: O((V + E) log V) using binary heap
// Space Complexity: O(V) for distances, parent arrays, and priority queue
func (g *WeightedGraph) Dijkstra(source int) ([]int, []int, error) {
	if source < 0 || source >= g.vertices {
		return nil, nil, errors.New("source vertex out of bounds")
	}

	// Initialize distances and parent arrays
	distances := make([]int, g.vertices)
	parent := make([]int, g.vertices)
	visited := make([]bool, g.vertices)

	// Initialize all distances to infinity except source
	for i := 0; i < g.vertices; i++ {
		distances[i] = math.MaxInt32
		parent[i] = -1
	}
	distances[source] = 0

	// Create and initialize priority queue
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &PriorityQueueItem{vertex: source, distance: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*PriorityQueueItem)
		currentVertex := current.vertex

		// Skip if already visited (handles duplicate entries in PQ)
		if visited[currentVertex] {
			continue
		}
		visited[currentVertex] = true

		// Check all neighbors
		for _, edge := range g.adjList[currentVertex] {
			neighbor := edge.To
			weight := edge.Weight

			if !visited[neighbor] {
				newDistance := distances[currentVertex] + weight

				// Relaxation step
				if newDistance < distances[neighbor] {
					distances[neighbor] = newDistance
					parent[neighbor] = currentVertex
					heap.Push(pq, &PriorityQueueItem{vertex: neighbor, distance: newDistance})
				}
			}
		}
	}

	return distances, parent, nil
}

// DijkstraShortestPath finds the shortest path between source and destination
// Returns the path, total distance, and error if any
// Time Complexity: O((V + E) log V)
// Space Complexity: O(V)
func (g *WeightedGraph) DijkstraShortestPath(source, destination int) ([]int, int, error) {
	distances, parent, err := g.Dijkstra(source)
	if err != nil {
		return nil, -1, err
	}

	// Check if destination is reachable
	if distances[destination] == math.MaxInt32 {
		return nil, -1, nil // No path exists
	}

	// Reconstruct path using parent array
	path := make([]int, 0)
	current := destination

	for current != -1 {
		path = append([]int{current}, path...)
		current = parent[current]
	}

	return path, distances[destination], nil
}

// BellmanFord finds shortest paths from source to all vertices, handling negative weights
// Returns distances, parent array, and whether negative cycle exists
// Time Complexity: O(V * E)
// Space Complexity: O(V)
func (g *WeightedGraph) BellmanFord(source int) ([]int, []int, bool, error) {
	if source < 0 || source >= g.vertices {
		return nil, nil, false, errors.New("source vertex out of bounds")
	}

	// Initialize distances and parent arrays
	distances := make([]int, g.vertices)
	parent := make([]int, g.vertices)

	for i := 0; i < g.vertices; i++ {
		distances[i] = math.MaxInt32
		parent[i] = -1
	}
	distances[source] = 0

	// Relax all edges V-1 times
	for i := 0; i < g.vertices-1; i++ {
		for vertex := 0; vertex < g.vertices; vertex++ {
			if distances[vertex] != math.MaxInt32 {
				for _, edge := range g.adjList[vertex] {
					neighbor := edge.To
					weight := edge.Weight

					if distances[vertex]+weight < distances[neighbor] {
						distances[neighbor] = distances[vertex] + weight
						parent[neighbor] = vertex
					}
				}
			}
		}
	}

	// Check for negative cycles
	hasNegativeCycle := false
	for vertex := 0; vertex < g.vertices; vertex++ {
		if distances[vertex] != math.MaxInt32 {
			for _, edge := range g.adjList[vertex] {
				neighbor := edge.To
				weight := edge.Weight

				if distances[vertex]+weight < distances[neighbor] {
					hasNegativeCycle = true
					break
				}
			}
		}
		if hasNegativeCycle {
			break
		}
	}

	return distances, parent, hasNegativeCycle, nil
}

// FloydWarshall finds shortest paths between all pairs of vertices
// Returns distance matrix and next vertex matrix for path reconstruction
// Time Complexity: O(V^3)
// Space Complexity: O(V^2)
func (g *WeightedGraph) FloydWarshall() ([][]int, [][]int, error) {
	// Initialize distance and next matrices
	dist := make([][]int, g.vertices)
	next := make([][]int, g.vertices)

	for i := 0; i < g.vertices; i++ {
		dist[i] = make([]int, g.vertices)
		next[i] = make([]int, g.vertices)

		for j := 0; j < g.vertices; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = math.MaxInt32
			}
			next[i][j] = -1
		}
	}

	// Fill initial distances from adjacency list
	for vertex := 0; vertex < g.vertices; vertex++ {
		for _, edge := range g.adjList[vertex] {
			dist[vertex][edge.To] = edge.Weight
			next[vertex][edge.To] = edge.To
		}
	}

	// Floyd-Warshall algorithm
	for k := 0; k < g.vertices; k++ {
		for i := 0; i < g.vertices; i++ {
			for j := 0; j < g.vertices; j++ {
				if dist[i][k] != math.MaxInt32 && dist[k][j] != math.MaxInt32 {
					if dist[i][k]+dist[k][j] < dist[i][j] {
						dist[i][j] = dist[i][k] + dist[k][j]
						next[i][j] = next[i][k]
					}
				}
			}
		}
	}

	return dist, next, nil
}

// ReconstructFloydWarshallPath reconstructs path between two vertices using next matrix
func (g *WeightedGraph) ReconstructFloydWarshallPath(next [][]int, start, end int) []int {
	if next[start][end] == -1 {
		return nil // No path exists
	}

	path := []int{start}
	current := start

	for current != end {
		current = next[current][end]
		path = append(path, current)
	}

	return path
}

// GetVerticesCount returns the number of vertices in the weighted graph
func (g *WeightedGraph) GetVerticesCount() int {
	return g.vertices
}

// GetWeightedAdjacent returns the adjacency list for a given vertex
func (g *WeightedGraph) GetWeightedAdjacent(vertex int) ([]Edge, error) {
	if vertex < 0 || vertex >= g.vertices {
		return nil, errors.New("vertex index out of bounds")
	}

	// Return a copy to prevent external modification
	result := make([]Edge, len(g.adjList[vertex]))
	copy(result, g.adjList[vertex])
	return result, nil
}

// PrintWeightedGraph prints the weighted adjacency list representation
func (g *WeightedGraph) PrintWeightedGraph() {
	fmt.Printf("Weighted Graph with %d vertices:\n", g.vertices)
	for i := 0; i < g.vertices; i++ {
		fmt.Printf("Vertex %d: ", i)
		for _, edge := range g.adjList[i] {
			fmt.Printf("(%d, %d) ", edge.To, edge.Weight)
		}
		fmt.Println()
	}
}

// ========================================================================
// ADVANCED SHORTEST PATH ALGORITHMS
// ========================================================================

// Coordinate represents a 2D point for A* algorithm
type Coordinate struct {
	X, Y int
}

// AStarNode represents a node in A* algorithm with f, g, h scores
type AStarNode struct {
	vertex int
	coord  Coordinate
	g      int // Distance from start
	h      int // Heuristic (estimated distance to goal)
	f      int // Total score (g + h)
	parent *AStarNode
	index  int // Index in heap
}

// AStarPriorityQueue implements min-heap for A* algorithm
type AStarPriorityQueue []*AStarNode

func (pq AStarPriorityQueue) Len() int { return len(pq) }

func (pq AStarPriorityQueue) Less(i, j int) bool {
	if pq[i].f == pq[j].f {
		return pq[i].h < pq[j].h // Break ties by h score
	}
	return pq[i].f < pq[j].f
}

func (pq AStarPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *AStarPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*AStarNode)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *AStarPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

// ManhattanDistance calculates Manhattan distance heuristic
func ManhattanDistance(a, b Coordinate) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// EuclideanDistance calculates Euclidean distance heuristic
func EuclideanDistance(a, b Coordinate) int {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return int(math.Sqrt(dx*dx + dy*dy))
}

// AStar finds shortest path using A* algorithm with coordinates and heuristic
// Time Complexity: O(E) in best case, O(b^d) in worst case where b is branching factor
// Space Complexity: O(V)
func (g *WeightedGraph) AStar(start, goal int, coordinates map[int]Coordinate, heuristic func(Coordinate, Coordinate) int) ([]int, int, error) {
	if start < 0 || start >= g.vertices || goal < 0 || goal >= g.vertices {
		return nil, -1, errors.New("vertex index out of bounds")
	}

	startCoord, startExists := coordinates[start]
	goalCoord, goalExists := coordinates[goal]
	if !startExists || !goalExists {
		return nil, -1, errors.New("coordinates not provided for start or goal vertex")
	}

	openSet := &AStarPriorityQueue{}
	closedSet := make(map[int]bool)
	gScore := make(map[int]int)
	nodeMap := make(map[int]*AStarNode)

	// Initialize start node
	startNode := &AStarNode{
		vertex: start,
		coord:  startCoord,
		g:      0,
		h:      heuristic(startCoord, goalCoord),
		parent: nil,
	}
	startNode.f = startNode.g + startNode.h

	heap.Init(openSet)
	heap.Push(openSet, startNode)
	gScore[start] = 0
	nodeMap[start] = startNode

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*AStarNode)
		currentVertex := current.vertex

		if currentVertex == goal {
			// Reconstruct path
			path := make([]int, 0)
			node := current
			for node != nil {
				path = append([]int{node.vertex}, path...)
				node = node.parent
			}
			return path, current.g, nil
		}

		closedSet[currentVertex] = true

		// Check all neighbors
		for _, edge := range g.adjList[currentVertex] {
			neighbor := edge.To

			if closedSet[neighbor] {
				continue
			}

			tentativeG := current.g + edge.Weight

			if existingG, exists := gScore[neighbor]; !exists || tentativeG < existingG {
				neighborCoord := coordinates[neighbor]
				neighborNode := &AStarNode{
					vertex: neighbor,
					coord:  neighborCoord,
					g:      tentativeG,
					h:      heuristic(neighborCoord, goalCoord),
					parent: current,
				}
				neighborNode.f = neighborNode.g + neighborNode.h

				gScore[neighbor] = tentativeG
				nodeMap[neighbor] = neighborNode
				heap.Push(openSet, neighborNode)
			}
		}
	}

	return nil, -1, nil // No path found
}

// JohnsonAllPairs finds shortest paths between all pairs using Johnson's algorithm
// Efficient for sparse graphs with negative edges (but no negative cycles)
// Time Complexity: O(V^2 log V + VE)
// Space Complexity: O(V^2)
func (g *WeightedGraph) JohnsonAllPairs() ([][]int, error) {
	// Step 1: Add new vertex connected to all vertices with weight 0
	newVertex := g.vertices
	tempGraph := NewWeightedGraph(g.vertices + 1)

	// Copy original edges
	for v := 0; v < g.vertices; v++ {
		for _, edge := range g.adjList[v] {
			tempGraph.AddWeightedEdge(v, edge.To, edge.Weight)
		}
	}

	// Add edges from new vertex to all original vertices
	for v := 0; v < g.vertices; v++ {
		tempGraph.AddWeightedEdge(newVertex, v, 0)
	}

	// Step 2: Run Bellman-Ford from new vertex
	h, _, hasNegCycle, err := tempGraph.BellmanFord(newVertex)
	if err != nil {
		return nil, err
	}
	if hasNegCycle {
		return nil, errors.New("graph contains negative cycle")
	}

	// Step 3: Reweight all edges
	reweightedGraph := NewWeightedGraph(g.vertices)
	for v := 0; v < g.vertices; v++ {
		for _, edge := range g.adjList[v] {
			newWeight := edge.Weight + h[v] - h[edge.To]
			reweightedGraph.AddWeightedEdge(v, edge.To, newWeight)
		}
	}

	// Step 4: Run Dijkstra from each vertex on reweighted graph
	result := make([][]int, g.vertices)
	for v := 0; v < g.vertices; v++ {
		distances, _, err := reweightedGraph.Dijkstra(v)
		if err != nil {
			return nil, err
		}

		result[v] = make([]int, g.vertices)
		for u := 0; u < g.vertices; u++ {
			if distances[u] == math.MaxInt32 {
				result[v][u] = math.MaxInt32
			} else {
				// Restore original weights
				result[v][u] = distances[u] - h[v] + h[u]
			}
		}
	}

	return result, nil
}

// BidirectionalDijkstra finds shortest path using bidirectional search
// More efficient for single source-destination queries
// Time Complexity: O(E + V log V) - often faster in practice
// Space Complexity: O(V)
func (g *WeightedGraph) BidirectionalDijkstra(source, target int) ([]int, int, error) {
	if source < 0 || source >= g.vertices || target < 0 || target >= g.vertices {
		return nil, -1, errors.New("vertex index out of bounds")
	}

	if source == target {
		return []int{source}, 0, nil
	}

	// Forward search data structures
	fwdDist := make([]int, g.vertices)
	fwdParent := make([]int, g.vertices)
	fwdVisited := make([]bool, g.vertices)
	fwdPQ := &PriorityQueue{}

	// Backward search data structures
	bwdDist := make([]int, g.vertices)
	bwdParent := make([]int, g.vertices)
	bwdVisited := make([]bool, g.vertices)
	bwdPQ := &PriorityQueue{}

	// Build reverse graph for backward search
	reverseGraph := NewWeightedGraph(g.vertices)
	for v := 0; v < g.vertices; v++ {
		for _, edge := range g.adjList[v] {
			reverseGraph.AddWeightedEdge(edge.To, v, edge.Weight)
		}
	}

	// Initialize distances
	for i := 0; i < g.vertices; i++ {
		fwdDist[i] = math.MaxInt32
		bwdDist[i] = math.MaxInt32
		fwdParent[i] = -1
		bwdParent[i] = -1
	}

	fwdDist[source] = 0
	bwdDist[target] = 0

	heap.Init(fwdPQ)
	heap.Init(bwdPQ)
	heap.Push(fwdPQ, &PriorityQueueItem{vertex: source, distance: 0})
	heap.Push(bwdPQ, &PriorityQueueItem{vertex: target, distance: 0})

	meetingPoint := -1
	minPathLength := math.MaxInt32

	for fwdPQ.Len() > 0 || bwdPQ.Len() > 0 {
		// Forward search step
		if fwdPQ.Len() > 0 {
			current := heap.Pop(fwdPQ).(*PriorityQueueItem)
			v := current.vertex

			if !fwdVisited[v] {
				fwdVisited[v] = true

				// Check if we've met the backward search
				if bwdVisited[v] && fwdDist[v]+bwdDist[v] < minPathLength {
					minPathLength = fwdDist[v] + bwdDist[v]
					meetingPoint = v
				}

				// Process neighbors
				for _, edge := range g.adjList[v] {
					neighbor := edge.To
					if !fwdVisited[neighbor] {
						newDist := fwdDist[v] + edge.Weight
						if newDist < fwdDist[neighbor] {
							fwdDist[neighbor] = newDist
							fwdParent[neighbor] = v
							heap.Push(fwdPQ, &PriorityQueueItem{vertex: neighbor, distance: newDist})
						}
					}
				}
			}
		}

		// Backward search step
		if bwdPQ.Len() > 0 {
			current := heap.Pop(bwdPQ).(*PriorityQueueItem)
			v := current.vertex

			if !bwdVisited[v] {
				bwdVisited[v] = true

				// Check if we've met the forward search
				if fwdVisited[v] && fwdDist[v]+bwdDist[v] < minPathLength {
					minPathLength = fwdDist[v] + bwdDist[v]
					meetingPoint = v
				}

				// Process neighbors in reverse graph
				for _, edge := range reverseGraph.adjList[v] {
					neighbor := edge.To
					if !bwdVisited[neighbor] {
						newDist := bwdDist[v] + edge.Weight
						if newDist < bwdDist[neighbor] {
							bwdDist[neighbor] = newDist
							bwdParent[neighbor] = v
							heap.Push(bwdPQ, &PriorityQueueItem{vertex: neighbor, distance: newDist})
						}
					}
				}
			}
		}
	}

	if meetingPoint == -1 {
		return nil, -1, nil // No path found
	}

	// Reconstruct path
	path := make([]int, 0)

	// Forward path from source to meeting point
	curr := meetingPoint
	fwdPath := make([]int, 0)
	for curr != -1 {
		fwdPath = append([]int{curr}, fwdPath...)
		curr = fwdParent[curr]
	}

	// Backward path from meeting point to target
	curr = bwdParent[meetingPoint]
	bwdPath := make([]int, 0)
	for curr != -1 {
		bwdPath = append(bwdPath, curr)
		curr = bwdParent[curr]
	}

	path = append(path, fwdPath...)
	path = append(path, bwdPath...)

	return path, minPathLength, nil
}

// KShortestPaths finds k shortest paths from source to destination using Yen's algorithm
// Time Complexity: O(K * V * (E + V log V))
// Space Complexity: O(K * V)
func (g *WeightedGraph) KShortestPaths(source, destination, k int) ([][]int, []int, error) {
	if source < 0 || source >= g.vertices || destination < 0 || destination >= g.vertices {
		return nil, nil, errors.New("vertex index out of bounds")
	}
	if k <= 0 {
		return nil, nil, errors.New("k must be positive")
	}

	type PathInfo struct {
		path     []int
		distance int
	}

	// Find shortest path
	shortestPath, shortestDist, err := g.DijkstraShortestPath(source, destination)
	if err != nil {
		return nil, nil, err
	}
	if shortestPath == nil {
		return nil, nil, nil // No path exists
	}

	paths := [][]int{shortestPath}
	distances := []int{shortestDist}

	candidates := make([]PathInfo, 0)

	for len(paths) < k {
		if len(paths) == 0 {
			break
		}

		lastPath := paths[len(paths)-1]

		// Generate candidate paths
		for i := 0; i < len(lastPath)-1; i++ {
			spurNode := lastPath[i]
			rootPath := lastPath[:i+1]

			// Create modified graph by removing edges
			modifiedGraph := NewWeightedGraph(g.vertices)

			// Copy all edges except those to be removed
			for v := 0; v < g.vertices; v++ {
				for _, edge := range g.adjList[v] {
					shouldRemove := false

					// Remove edges used by previous paths at this point
					for _, prevPath := range paths {
						if len(prevPath) > i+1 && len(rootPath) > i {
							if prevPath[i] == v && prevPath[i+1] == edge.To {
								shouldRemove = true
								break
							}
						}
					}

					if !shouldRemove {
						modifiedGraph.AddWeightedEdge(v, edge.To, edge.Weight)
					}
				}
			}

			// Find shortest path from spur node to destination in modified graph
			spurPath, spurDist, err := modifiedGraph.DijkstraShortestPath(spurNode, destination)
			if err == nil && spurPath != nil {
				// Calculate total path and distance
				totalPath := append(append([]int{}, rootPath[:len(rootPath)-1]...), spurPath...)

				totalDist := spurDist
				for j := 0; j < len(rootPath)-1; j++ {
					// Find edge weight between consecutive nodes in root path
					for _, edge := range g.adjList[rootPath[j]] {
						if edge.To == rootPath[j+1] {
							totalDist += edge.Weight
							break
						}
					}
				}

				// Add to candidates if not duplicate
				isDuplicate := false
				for _, candidate := range candidates {
					if len(candidate.path) == len(totalPath) {
						match := true
						for idx := range candidate.path {
							if candidate.path[idx] != totalPath[idx] {
								match = false
								break
							}
						}
						if match {
							isDuplicate = true
							break
						}
					}
				}

				if !isDuplicate {
					candidates = append(candidates, PathInfo{path: totalPath, distance: totalDist})
				}
			}
		}

		if len(candidates) == 0 {
			break
		}

		// Sort candidates by distance
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].distance < candidates[j].distance
		})

		// Add best candidate to results
		paths = append(paths, candidates[0].path)
		distances = append(distances, candidates[0].distance)

		// Remove selected candidate
		candidates = candidates[1:]
	}

	return paths, distances, nil
}

// PrintShortestPathSummary prints analysis of shortest path algorithms performance
func (g *WeightedGraph) PrintShortestPathSummary() {
	fmt.Printf("=== Shortest Path Algorithms Summary ===\n")
	fmt.Printf("Graph: %d vertices\n\n", g.vertices)

	fmt.Printf("Available Algorithms:\n")
	fmt.Printf("1. Dijkstra: O((V+E)logV) - No negative weights\n")
	fmt.Printf("2. Bellman-Ford: O(VE) - Handles negative weights\n")
	fmt.Printf("3. Floyd-Warshall: O(V³) - All-pairs shortest paths\n")
	fmt.Printf("4. A*: O(E) best case - Uses heuristic for faster search\n")
	fmt.Printf("5. Johnson: O(V²logV + VE) - Sparse graphs with negative weights\n")
	fmt.Printf("6. Bidirectional: O(E+VlogV) - Single source-destination\n")
	fmt.Printf("7. K-Shortest: O(K*V*(E+VlogV)) - Multiple shortest paths\n\n")

	edgeCount := 0
	for v := 0; v < g.vertices; v++ {
		edgeCount += len(g.adjList[v])
	}

	fmt.Printf("Graph Properties:\n")
	fmt.Printf("- Vertices: %d\n", g.vertices)
	fmt.Printf("- Edges: %d\n", edgeCount)
	fmt.Printf("- Density: %.2f%%\n", float64(edgeCount)/float64(g.vertices*(g.vertices-1))*100)

	if edgeCount < g.vertices*g.vertices/4 {
		fmt.Printf("- Recommendation: Use Dijkstra or Johnson's for sparse graph\n")
	} else {
		fmt.Printf("- Recommendation: Use Floyd-Warshall for dense graph\n")
	}
}
