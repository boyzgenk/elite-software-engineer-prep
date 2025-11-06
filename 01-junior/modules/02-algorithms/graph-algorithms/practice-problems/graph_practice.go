// Graph Algorithm Practice Problems
// Comprehensive practice problems covering all major graph algorithms

package main

import (
	"fmt"
	"math"
	"sort"
)

// ========================================================================
// BASIC GRAPH STRUCTURES
// ========================================================================

type Edge struct {
	To     int
	Weight int
}

type Graph struct {
	vertices int
	adjList  [][]Edge
	directed bool
}

func NewGraph(vertices int, directed bool) *Graph {
	adjList := make([][]Edge, vertices)
	for i := range adjList {
		adjList[i] = make([]Edge, 0)
	}
	return &Graph{vertices: vertices, adjList: adjList, directed: directed}
}

func (g *Graph) AddEdge(from, to, weight int) {
	g.adjList[from] = append(g.adjList[from], Edge{To: to, Weight: weight})
	if !g.directed {
		g.adjList[to] = append(g.adjList[to], Edge{To: from, Weight: weight})
	}
}

// ========================================================================
// PROBLEM 1: GRAPH TRAVERSAL PROBLEMS
// ========================================================================

// Problem 1.1: Count Connected Components
// Given an undirected graph, count the number of connected components
func countConnectedComponents(graph *Graph) int {
	visited := make([]bool, graph.vertices)
	count := 0

	for i := 0; i < graph.vertices; i++ {
		if !visited[i] {
			dfsComponent(graph, i, visited)
			count++
		}
	}

	return count
}

func dfsComponent(graph *Graph, vertex int, visited []bool) {
	visited[vertex] = true
	for _, edge := range graph.adjList[vertex] {
		if !visited[edge.To] {
			dfsComponent(graph, edge.To, visited)
		}
	}
}

// Problem 1.2: Find All Paths Between Two Vertices
// Find all simple paths from source to destination
func findAllPaths(graph *Graph, start, end int) [][]int {
	var allPaths [][]int
	visited := make([]bool, graph.vertices)
	var currentPath []int

	findAllPathsDFS(graph, start, end, visited, currentPath, &allPaths)
	return allPaths
}

func findAllPathsDFS(graph *Graph, current, target int, visited []bool, path []int, allPaths *[][]int) {
	visited[current] = true
	path = append(path, current)

	if current == target {
		// Found a path, add copy to results
		pathCopy := make([]int, len(path))
		copy(pathCopy, path)
		*allPaths = append(*allPaths, pathCopy)
	} else {
		// Continue searching
		for _, edge := range graph.adjList[current] {
			if !visited[edge.To] {
				findAllPathsDFS(graph, edge.To, target, visited, path, allPaths)
			}
		}
	}

	// Backtrack
	visited[current] = false
}

// Problem 1.3: Detect Cycle in Undirected Graph
// Return true if the undirected graph contains a cycle
func hasCycleUndirected(graph *Graph) bool {
	visited := make([]bool, graph.vertices)

	for i := 0; i < graph.vertices; i++ {
		if !visited[i] {
			if hasCycleUndirectedDFS(graph, i, -1, visited) {
				return true
			}
		}
	}
	return false
}

func hasCycleUndirectedDFS(graph *Graph, vertex, parent int, visited []bool) bool {
	visited[vertex] = true

	for _, edge := range graph.adjList[vertex] {
		neighbor := edge.To
		if !visited[neighbor] {
			if hasCycleUndirectedDFS(graph, neighbor, vertex, visited) {
				return true
			}
		} else if neighbor != parent {
			return true // Back edge found (cycle)
		}
	}
	return false
}

// Problem 1.4: Detect Cycle in Directed Graph
// Return true if the directed graph contains a cycle
func hasCycleDirected(graph *Graph) bool {
	visited := make([]bool, graph.vertices)
	recStack := make([]bool, graph.vertices)

	for i := 0; i < graph.vertices; i++ {
		if !visited[i] {
			if hasCycleDirectedDFS(graph, i, visited, recStack) {
				return true
			}
		}
	}
	return false
}

func hasCycleDirectedDFS(graph *Graph, vertex int, visited, recStack []bool) bool {
	visited[vertex] = true
	recStack[vertex] = true

	for _, edge := range graph.adjList[vertex] {
		neighbor := edge.To
		if !visited[neighbor] {
			if hasCycleDirectedDFS(graph, neighbor, visited, recStack) {
				return true
			}
		} else if recStack[neighbor] {
			return true // Back edge in recursion stack
		}
	}

	recStack[vertex] = false
	return false
}

// ========================================================================
// PROBLEM 2: SHORTEST PATH PROBLEMS
// ========================================================================

// Problem 2.1: Dijkstra's Algorithm Implementation
// Find shortest path from source to all vertices
func dijkstra(graph *Graph, source int) ([]int, []int) {
	dist := make([]int, graph.vertices)
	parent := make([]int, graph.vertices)
	visited := make([]bool, graph.vertices)

	// Initialize distances
	for i := range dist {
		dist[i] = math.MaxInt32
		parent[i] = -1
	}
	dist[source] = 0

	for count := 0; count < graph.vertices; count++ {
		// Find minimum distance vertex
		u := minDistance(dist, visited)
		if u == -1 {
			break
		}
		visited[u] = true

		// Update distances of adjacent vertices
		for _, edge := range graph.adjList[u] {
			v := edge.To
			weight := edge.Weight
			if !visited[v] && dist[u] != math.MaxInt32 && dist[u]+weight < dist[v] {
				dist[v] = dist[u] + weight
				parent[v] = u
			}
		}
	}

	return dist, parent
}

func minDistance(dist []int, visited []bool) int {
	min := math.MaxInt32
	minIndex := -1

	for v := 0; v < len(dist); v++ {
		if !visited[v] && dist[v] <= min {
			min = dist[v]
			minIndex = v
		}
	}
	return minIndex
}

// Problem 2.2: Bellman-Ford Algorithm
// Find shortest paths and detect negative cycles
func bellmanFord(graph *Graph, source int) ([]int, []int, bool) {
	dist := make([]int, graph.vertices)
	parent := make([]int, graph.vertices)

	// Initialize distances
	for i := range dist {
		dist[i] = math.MaxInt32
		parent[i] = -1
	}
	dist[source] = 0

	// Relax edges V-1 times
	for i := 0; i < graph.vertices-1; i++ {
		for u := 0; u < graph.vertices; u++ {
			if dist[u] != math.MaxInt32 {
				for _, edge := range graph.adjList[u] {
					v := edge.To
					weight := edge.Weight
					if dist[u]+weight < dist[v] {
						dist[v] = dist[u] + weight
						parent[v] = u
					}
				}
			}
		}
	}

	// Check for negative cycles
	hasNegativeCycle := false
	for u := 0; u < graph.vertices; u++ {
		if dist[u] != math.MaxInt32 {
			for _, edge := range graph.adjList[u] {
				v := edge.To
				weight := edge.Weight
				if dist[u]+weight < dist[v] {
					hasNegativeCycle = true
					break
				}
			}
		}
		if hasNegativeCycle {
			break
		}
	}

	return dist, parent, hasNegativeCycle
}

// Problem 2.3: Floyd-Warshall Algorithm
// All-pairs shortest paths
func floydWarshall(graph *Graph) [][]int {
	dist := make([][]int, graph.vertices)
	for i := range dist {
		dist[i] = make([]int, graph.vertices)
		for j := range dist[i] {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = math.MaxInt32
			}
		}
	}

	// Fill in direct edges
	for u := 0; u < graph.vertices; u++ {
		for _, edge := range graph.adjList[u] {
			dist[u][edge.To] = edge.Weight
		}
	}

	// Floyd-Warshall main loop
	for k := 0; k < graph.vertices; k++ {
		for i := 0; i < graph.vertices; i++ {
			for j := 0; j < graph.vertices; j++ {
				if dist[i][k] != math.MaxInt32 && dist[k][j] != math.MaxInt32 {
					if dist[i][k]+dist[k][j] < dist[i][j] {
						dist[i][j] = dist[i][k] + dist[k][j]
					}
				}
			}
		}
	}

	return dist
}

// ========================================================================
// PROBLEM 3: MINIMUM SPANNING TREE PROBLEMS
// ========================================================================

// Edge for MST algorithms
type MSTEdge struct {
	From, To, Weight int
}

// Problem 3.1: Kruskal's Algorithm
// Find minimum spanning tree using Kruskal's algorithm
func kruskalMST(graph *Graph) ([]MSTEdge, int) {
	edges := getAllEdges(graph)
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Weight < edges[j].Weight
	})

	uf := NewUnionFind(graph.vertices)
	var mst []MSTEdge
	totalWeight := 0

	for _, edge := range edges {
		if uf.Union(edge.From, edge.To) {
			mst = append(mst, edge)
			totalWeight += edge.Weight
			if len(mst) == graph.vertices-1 {
				break
			}
		}
	}

	return mst, totalWeight
}

func getAllEdges(graph *Graph) []MSTEdge {
	var edges []MSTEdge
	for u := 0; u < graph.vertices; u++ {
		for _, edge := range graph.adjList[u] {
			if graph.directed || u < edge.To { // Avoid duplicates for undirected
				edges = append(edges, MSTEdge{u, edge.To, edge.Weight})
			}
		}
	}
	return edges
}

// Simple Union-Find implementation
type UnionFind struct {
	parent []int
	rank   []int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent, rank}
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // Path compression
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX, rootY := uf.Find(x), uf.Find(y)
	if rootX == rootY {
		return false // Already connected
	}

	// Union by rank
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}
	return true
}

// Problem 3.2: Prim's Algorithm
// Find minimum spanning tree using Prim's algorithm
func primMST(graph *Graph) ([]MSTEdge, int) {
	inMST := make([]bool, graph.vertices)
	key := make([]int, graph.vertices)
	parent := make([]int, graph.vertices)

	// Initialize keys
	for i := range key {
		key[i] = math.MaxInt32
		parent[i] = -1
	}
	key[0] = 0

	var mst []MSTEdge
	totalWeight := 0

	for count := 0; count < graph.vertices; count++ {
		u := minKey(key, inMST)
		inMST[u] = true

		if parent[u] != -1 {
			mst = append(mst, MSTEdge{parent[u], u, key[u]})
			totalWeight += key[u]
		}

		// Update keys of adjacent vertices
		for _, edge := range graph.adjList[u] {
			v := edge.To
			weight := edge.Weight
			if !inMST[v] && weight < key[v] {
				key[v] = weight
				parent[v] = u
			}
		}
	}

	return mst, totalWeight
}

func minKey(key []int, inMST []bool) int {
	min := math.MaxInt32
	minIndex := -1

	for v := 0; v < len(key); v++ {
		if !inMST[v] && key[v] < min {
			min = key[v]
			minIndex = v
		}
	}
	return minIndex
}

// ========================================================================
// PROBLEM 4: ADVANCED GRAPH PROBLEMS
// ========================================================================

// Problem 4.1: Strongly Connected Components (Kosaraju's Algorithm)
// Find all strongly connected components in a directed graph
func findSCCs(graph *Graph) [][]int {
	// Step 1: Get finish times using DFS
	visited := make([]bool, graph.vertices)
	var finishOrder []int

	for i := 0; i < graph.vertices; i++ {
		if !visited[i] {
			dfsFinishTime(graph, i, visited, &finishOrder)
		}
	}

	// Step 2: Create transpose graph
	transpose := createTranspose(graph)

	// Step 3: DFS on transpose in reverse finish order
	visited = make([]bool, graph.vertices)
	var sccs [][]int

	for i := len(finishOrder) - 1; i >= 0; i-- {
		vertex := finishOrder[i]
		if !visited[vertex] {
			var scc []int
			dfsCollectSCC(transpose, vertex, visited, &scc)
			sccs = append(sccs, scc)
		}
	}

	return sccs
}

func dfsFinishTime(graph *Graph, vertex int, visited []bool, finishOrder *[]int) {
	visited[vertex] = true
	for _, edge := range graph.adjList[vertex] {
		if !visited[edge.To] {
			dfsFinishTime(graph, edge.To, visited, finishOrder)
		}
	}
	*finishOrder = append(*finishOrder, vertex)
}

func createTranspose(graph *Graph) *Graph {
	transpose := NewGraph(graph.vertices, true)
	for u := 0; u < graph.vertices; u++ {
		for _, edge := range graph.adjList[u] {
			transpose.AddEdge(edge.To, u, edge.Weight)
		}
	}
	return transpose
}

func dfsCollectSCC(graph *Graph, vertex int, visited []bool, scc *[]int) {
	visited[vertex] = true
	*scc = append(*scc, vertex)
	for _, edge := range graph.adjList[vertex] {
		if !visited[edge.To] {
			dfsCollectSCC(graph, edge.To, visited, scc)
		}
	}
}

// Problem 4.2: Bridges in Graph (Tarjan's Algorithm)
// Find all bridges (critical edges) in an undirected graph
func findBridges(graph *Graph) []MSTEdge {
	visited := make([]bool, graph.vertices)
	disc := make([]int, graph.vertices)
	low := make([]int, graph.vertices)
	parent := make([]int, graph.vertices)
	var bridges []MSTEdge
	time := 0

	for i := range parent {
		parent[i] = -1
	}

	for i := 0; i < graph.vertices; i++ {
		if !visited[i] {
			bridgesDFS(graph, i, visited, disc, low, parent, &bridges, &time)
		}
	}

	return bridges
}

func bridgesDFS(graph *Graph, u int, visited []bool, disc, low, parent []int, bridges *[]MSTEdge, time *int) {
	visited[u] = true
	disc[u] = *time
	low[u] = *time
	*time++

	for _, edge := range graph.adjList[u] {
		v := edge.To
		if !visited[v] {
			parent[v] = u
			bridgesDFS(graph, v, visited, disc, low, parent, bridges, time)

			low[u] = min(low[u], low[v])

			// If low[v] > disc[u], then (u,v) is a bridge
			if low[v] > disc[u] {
				*bridges = append(*bridges, MSTEdge{u, v, edge.Weight})
			}
		} else if v != parent[u] {
			low[u] = min(low[u], disc[v])
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Problem 4.3: Articulation Points (Tarjan's Algorithm)
// Find all articulation points (critical vertices) in an undirected graph
func findArticulationPoints(graph *Graph) []int {
	visited := make([]bool, graph.vertices)
	disc := make([]int, graph.vertices)
	low := make([]int, graph.vertices)
	parent := make([]int, graph.vertices)
	ap := make([]bool, graph.vertices)
	time := 0

	for i := range parent {
		parent[i] = -1
	}

	for i := 0; i < graph.vertices; i++ {
		if !visited[i] {
			articulationDFS(graph, i, visited, disc, low, parent, ap, &time)
		}
	}

	var result []int
	for i := 0; i < graph.vertices; i++ {
		if ap[i] {
			result = append(result, i)
		}
	}
	return result
}

func articulationDFS(graph *Graph, u int, visited []bool, disc, low, parent []int, ap []bool, time *int) {
	children := 0
	visited[u] = true
	disc[u] = *time
	low[u] = *time
	*time++

	for _, edge := range graph.adjList[u] {
		v := edge.To
		if !visited[v] {
			children++
			parent[v] = u
			articulationDFS(graph, v, visited, disc, low, parent, ap, time)

			low[u] = min(low[u], low[v])

			// Check articulation point conditions
			if parent[u] == -1 && children > 1 {
				ap[u] = true // Root with multiple children
			}
			if parent[u] != -1 && low[v] >= disc[u] {
				ap[u] = true // Non-root with high low value
			}
		} else if v != parent[u] {
			low[u] = min(low[u], disc[v])
		}
	}
}

// ========================================================================
// PROBLEM 5: GRAPH COLORING
// ========================================================================

// Problem 5.1: Graph Coloring (Greedy)
// Color the graph using minimum number of colors
func greedyColoring(graph *Graph) ([]int, int) {
	colors := make([]int, graph.vertices)
	for i := range colors {
		colors[i] = -1
	}

	colors[0] = 0 // Color first vertex with color 0
	maxColor := 0

	for u := 1; u < graph.vertices; u++ {
		// Find available colors
		available := make([]bool, graph.vertices)
		for i := range available {
			available[i] = true
		}

		// Mark colors used by adjacent vertices as unavailable
		for _, edge := range graph.adjList[u] {
			v := edge.To
			if colors[v] != -1 {
				available[colors[v]] = false
			}
		}

		// Find first available color
		for color := 0; color < graph.vertices; color++ {
			if available[color] {
				colors[u] = color
				if color > maxColor {
					maxColor = color
				}
				break
			}
		}
	}

	return colors, maxColor + 1
}

// Problem 5.2: Check if Graph is Bipartite
// Return true if graph can be colored with 2 colors
func isBipartite(graph *Graph) bool {
	colors := make([]int, graph.vertices)
	for i := range colors {
		colors[i] = -1
	}

	for i := 0; i < graph.vertices; i++ {
		if colors[i] == -1 {
			if !isBipartiteBFS(graph, i, colors) {
				return false
			}
		}
	}
	return true
}

func isBipartiteBFS(graph *Graph, start int, colors []int) bool {
	queue := []int{start}
	colors[start] = 0

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		for _, edge := range graph.adjList[u] {
			v := edge.To
			if colors[v] == -1 {
				colors[v] = 1 - colors[u] // Alternate color
				queue = append(queue, v)
			} else if colors[v] == colors[u] {
				return false // Same color for adjacent vertices
			}
		}
	}
	return true
}

// ========================================================================
// DEMONSTRATION AND TESTING
// ========================================================================

func demonstrateTraversalProblems() {
	fmt.Println("=== TRAVERSAL PROBLEMS ===")

	// Create test graph
	graph := NewGraph(5, false)
	graph.AddEdge(0, 1, 1)
	graph.AddEdge(1, 2, 1)
	graph.AddEdge(3, 4, 1)

	fmt.Printf("Connected Components: %d\n", countConnectedComponents(graph))
	fmt.Printf("Has Cycle (Undirected): %t\n", hasCycleUndirected(graph))

	// Test directed cycle
	dirGraph := NewGraph(3, true)
	dirGraph.AddEdge(0, 1, 1)
	dirGraph.AddEdge(1, 2, 1)
	dirGraph.AddEdge(2, 0, 1)
	fmt.Printf("Has Cycle (Directed): %t\n", hasCycleDirected(dirGraph))
}

func demonstrateShortestPath() {
	fmt.Println("\n=== SHORTEST PATH PROBLEMS ===")

	graph := NewGraph(5, true)
	graph.AddEdge(0, 1, 10)
	graph.AddEdge(0, 4, 5)
	graph.AddEdge(1, 2, 1)
	graph.AddEdge(1, 4, 2)
	graph.AddEdge(2, 3, 4)
	graph.AddEdge(3, 2, 6)
	graph.AddEdge(4, 1, 3)
	graph.AddEdge(4, 2, 9)
	graph.AddEdge(4, 3, 2)

	dist, _ := dijkstra(graph, 0)
	fmt.Printf("Dijkstra distances from 0: %v\n", dist)

	dist2, _, hasNegCycle := bellmanFord(graph, 0)
	fmt.Printf("Bellman-Ford distances: %v, Negative cycle: %t\n", dist2, hasNegCycle)
}

func demonstrateMST() {
	fmt.Println("\n=== MINIMUM SPANNING TREE ===")

	graph := NewGraph(4, false)
	graph.AddEdge(0, 1, 10)
	graph.AddEdge(0, 2, 6)
	graph.AddEdge(0, 3, 5)
	graph.AddEdge(1, 3, 15)
	graph.AddEdge(2, 3, 4)

	mst, weight := kruskalMST(graph)
	fmt.Printf("Kruskal MST weight: %d, edges: %v\n", weight, mst)

	mst2, weight2 := primMST(graph)
	fmt.Printf("Prim MST weight: %d, edges: %v\n", weight2, mst2)
}

func demonstrateAdvanced() {
	fmt.Println("\n=== ADVANCED PROBLEMS ===")

	// SCC example
	dirGraph := NewGraph(5, true)
	dirGraph.AddEdge(1, 0, 1)
	dirGraph.AddEdge(0, 2, 1)
	dirGraph.AddEdge(2, 1, 1)
	dirGraph.AddEdge(0, 3, 1)
	dirGraph.AddEdge(3, 4, 1)

	sccs := findSCCs(dirGraph)
	fmt.Printf("Strongly Connected Components: %v\n", sccs)

	// Bridges example
	undirGraph := NewGraph(5, false)
	undirGraph.AddEdge(1, 0, 1)
	undirGraph.AddEdge(0, 2, 1)
	undirGraph.AddEdge(2, 1, 1)
	undirGraph.AddEdge(0, 3, 1)
	undirGraph.AddEdge(3, 4, 1)

	bridges := findBridges(undirGraph)
	fmt.Printf("Bridges: %v\n", bridges)

	ap := findArticulationPoints(undirGraph)
	fmt.Printf("Articulation Points: %v\n", ap)
}

func demonstrateColoring() {
	fmt.Println("\n=== GRAPH COLORING ===")

	graph := NewGraph(5, false)
	graph.AddEdge(0, 1, 1)
	graph.AddEdge(1, 2, 1)
	graph.AddEdge(2, 3, 1)
	graph.AddEdge(3, 4, 1)
	graph.AddEdge(4, 0, 1)
	graph.AddEdge(1, 3, 1)

	colors, numColors := greedyColoring(graph)
	fmt.Printf("Graph coloring: %v, Colors needed: %d\n", colors, numColors)

	fmt.Printf("Is bipartite: %t\n", isBipartite(graph))
}

func main() {
	fmt.Println("🧩 GRAPH ALGORITHM PRACTICE PROBLEMS")
	fmt.Println("=====================================")

	demonstrateTraversalProblems()
	demonstrateShortestPath()
	demonstrateMST()
	demonstrateAdvanced()
	demonstrateColoring()

	fmt.Println("\n🎯 Practice Suggestions:")
	fmt.Println("1. Implement each algorithm from scratch")
	fmt.Println("2. Test with different graph types (dense, sparse, directed, undirected)")
	fmt.Println("3. Analyze time and space complexity")
	fmt.Println("4. Practice on LeetCode graph problems")
	fmt.Println("5. Understand when to use each algorithm")
}
