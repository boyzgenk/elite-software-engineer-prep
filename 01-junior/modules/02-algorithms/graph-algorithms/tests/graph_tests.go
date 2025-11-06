// Comprehensive Graph Algorithm Tests
// Tests for all graph algorithms implementations

package graph

import (
	"fmt"
	"testing"
)

// ========================================================================
// TEST STRUCTURES
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
// TRAVERSAL ALGORITHM TESTS
// ========================================================================

func TestDFS(t *testing.T) {
	fmt.Println("Testing DFS traversal...")

	graph := NewGraph(5, false)
	graph.AddEdge(0, 1, 1)
	graph.AddEdge(0, 2, 1)
	graph.AddEdge(1, 3, 1)
	graph.AddEdge(2, 4, 1)

	visited := make([]bool, graph.vertices)
	var result []int

	dfs(graph, 0, visited, &result)

	if len(result) != 5 {
		t.Errorf("Expected 5 vertices in DFS, got %d", len(result))
	}

	fmt.Printf("✅ DFS Result: %v\n", result)
}

func dfs(graph *Graph, vertex int, visited []bool, result *[]int) {
	visited[vertex] = true
	*result = append(*result, vertex)

	for _, edge := range graph.adjList[vertex] {
		if !visited[edge.To] {
			dfs(graph, edge.To, visited, result)
		}
	}
}

func TestBFS(t *testing.T) {
	fmt.Println("Testing BFS traversal...")

	graph := NewGraph(5, false)
	graph.AddEdge(0, 1, 1)
	graph.AddEdge(0, 2, 1)
	graph.AddEdge(1, 3, 1)
	graph.AddEdge(2, 4, 1)

	result := bfs(graph, 0)

	if len(result) != 5 {
		t.Errorf("Expected 5 vertices in BFS, got %d", len(result))
	}

	// Check BFS property: level-order traversal
	if result[0] != 0 {
		t.Errorf("Expected BFS to start with vertex 0")
	}

	fmt.Printf("✅ BFS Result: %v\n", result)
}

func bfs(graph *Graph, start int) []int {
	visited := make([]bool, graph.vertices)
	var result []int
	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		result = append(result, vertex)

		for _, edge := range graph.adjList[vertex] {
			if !visited[edge.To] {
				visited[edge.To] = true
				queue = append(queue, edge.To)
			}
		}
	}

	return result
}

// ========================================================================
// CYCLE DETECTION TESTS
// ========================================================================

func TestCycleDetectionUndirected(t *testing.T) {
	fmt.Println("Testing cycle detection in undirected graph...")

	// Graph without cycle
	graph1 := NewGraph(3, false)
	graph1.AddEdge(0, 1, 1)
	graph1.AddEdge(1, 2, 1)

	if hasCycleUndirected(graph1) {
		t.Errorf("Expected no cycle in acyclic graph")
	}

	// Graph with cycle
	graph2 := NewGraph(3, false)
	graph2.AddEdge(0, 1, 1)
	graph2.AddEdge(1, 2, 1)
	graph2.AddEdge(2, 0, 1)

	if !hasCycleUndirected(graph2) {
		t.Errorf("Expected cycle in cyclic graph")
	}

	fmt.Println("✅ Undirected cycle detection passed")
}

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
			return true
		}
	}
	return false
}

func TestCycleDetectionDirected(t *testing.T) {
	fmt.Println("Testing cycle detection in directed graph...")

	// DAG (no cycle)
	graph1 := NewGraph(3, true)
	graph1.AddEdge(0, 1, 1)
	graph1.AddEdge(1, 2, 1)

	if hasCycleDirected(graph1) {
		t.Errorf("Expected no cycle in DAG")
	}

	// Graph with cycle
	graph2 := NewGraph(3, true)
	graph2.AddEdge(0, 1, 1)
	graph2.AddEdge(1, 2, 1)
	graph2.AddEdge(2, 0, 1)

	if !hasCycleDirected(graph2) {
		t.Errorf("Expected cycle in cyclic directed graph")
	}

	fmt.Println("✅ Directed cycle detection passed")
}

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
			return true
		}
	}

	recStack[vertex] = false
	return false
}

// ========================================================================
// SHORTEST PATH ALGORITHM TESTS
// ========================================================================

func TestDijkstra(t *testing.T) {
	fmt.Println("Testing Dijkstra's algorithm...")

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

	distances, _ := dijkstra(graph, 0)

	// Expected distances from vertex 0
	expected := []int{0, 8, 9, 7, 5}

	for i, dist := range distances {
		if dist != expected[i] {
			t.Errorf("Expected distance to vertex %d: %d, got %d", i, expected[i], dist)
		}
	}

	fmt.Printf("✅ Dijkstra distances: %v\n", distances)
}

func dijkstra(graph *Graph, source int) ([]int, []int) {
	const INF = int(1e9)
	dist := make([]int, graph.vertices)
	parent := make([]int, graph.vertices)
	visited := make([]bool, graph.vertices)

	for i := range dist {
		dist[i] = INF
		parent[i] = -1
	}
	dist[source] = 0

	for count := 0; count < graph.vertices; count++ {
		u := minDistance(dist, visited)
		if u == -1 {
			break
		}
		visited[u] = true

		for _, edge := range graph.adjList[u] {
			v := edge.To
			weight := edge.Weight
			if !visited[v] && dist[u] != INF && dist[u]+weight < dist[v] {
				dist[v] = dist[u] + weight
				parent[v] = u
			}
		}
	}

	return dist, parent
}

func minDistance(dist []int, visited []bool) int {
	const INF = int(1e9)
	min := INF
	minIndex := -1

	for v := 0; v < len(dist); v++ {
		if !visited[v] && dist[v] <= min {
			min = dist[v]
			minIndex = v
		}
	}
	return minIndex
}

// ========================================================================
// TOPOLOGICAL SORT TESTS
// ========================================================================

func TestTopologicalSort(t *testing.T) {
	fmt.Println("Testing topological sort...")

	// Create DAG
	graph := NewGraph(6, true)
	graph.AddEdge(5, 2, 1)
	graph.AddEdge(5, 0, 1)
	graph.AddEdge(4, 0, 1)
	graph.AddEdge(4, 1, 1)
	graph.AddEdge(2, 3, 1)
	graph.AddEdge(3, 1, 1)

	result := topologicalSort(graph)

	if len(result) != graph.vertices {
		t.Errorf("Expected %d vertices in topological sort, got %d", graph.vertices, len(result))
	}

	// Verify topological order property
	position := make([]int, graph.vertices)
	for i, vertex := range result {
		position[vertex] = i
	}

	for u := 0; u < graph.vertices; u++ {
		for _, edge := range graph.adjList[u] {
			v := edge.To
			if position[u] >= position[v] {
				t.Errorf("Topological order violated: vertex %d should come before %d", u, v)
			}
		}
	}

	fmt.Printf("✅ Topological sort: %v\n", result)
}

func topologicalSort(graph *Graph) []int {
	visited := make([]bool, graph.vertices)
	var stack []int

	for i := 0; i < graph.vertices; i++ {
		if !visited[i] {
			topologicalSortDFS(graph, i, visited, &stack)
		}
	}

	// Reverse stack to get topological order
	result := make([]int, len(stack))
	for i, v := range stack {
		result[len(stack)-1-i] = v
	}

	return result
}

func topologicalSortDFS(graph *Graph, vertex int, visited []bool, stack *[]int) {
	visited[vertex] = true

	for _, edge := range graph.adjList[vertex] {
		if !visited[edge.To] {
			topologicalSortDFS(graph, edge.To, visited, stack)
		}
	}

	*stack = append(*stack, vertex)
}

// ========================================================================
// MINIMUM SPANNING TREE TESTS
// ========================================================================

func TestKruskalMST(t *testing.T) {
	fmt.Println("Testing Kruskal's MST algorithm...")

	graph := NewGraph(4, false)
	graph.AddEdge(0, 1, 10)
	graph.AddEdge(0, 2, 6)
	graph.AddEdge(0, 3, 5)
	graph.AddEdge(1, 3, 15)
	graph.AddEdge(2, 3, 4)

	mst, totalWeight := kruskalMST(graph)
	expectedWeight := 19 // 4 + 5 + 10

	if len(mst) != graph.vertices-1 {
		t.Errorf("Expected %d edges in MST, got %d", graph.vertices-1, len(mst))
	}

	if totalWeight != expectedWeight {
		t.Errorf("Expected MST weight %d, got %d", expectedWeight, totalWeight)
	}

	fmt.Printf("✅ Kruskal MST weight: %d, edges: %v\n", totalWeight, mst)
}

type MSTEdge struct {
	From, To, Weight int
}

func kruskalMST(graph *Graph) ([]MSTEdge, int) {
	edges := getAllEdges(graph)

	// Sort edges by weight
	for i := 0; i < len(edges)-1; i++ {
		for j := 0; j < len(edges)-i-1; j++ {
			if edges[j].Weight > edges[j+1].Weight {
				edges[j], edges[j+1] = edges[j+1], edges[j]
			}
		}
	}

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
			if graph.directed || u < edge.To {
				edges = append(edges, MSTEdge{u, edge.To, edge.Weight})
			}
		}
	}
	return edges
}

// Union-Find data structure
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
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX, rootY := uf.Find(x), uf.Find(y)
	if rootX == rootY {
		return false
	}

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

// ========================================================================
// GRAPH PROPERTIES TESTS
// ========================================================================

func TestConnectedComponents(t *testing.T) {
	fmt.Println("Testing connected components...")

	graph := NewGraph(5, false)
	graph.AddEdge(0, 1, 1)
	graph.AddEdge(1, 2, 1)
	graph.AddEdge(3, 4, 1)

	components := countConnectedComponents(graph)
	expected := 2

	if components != expected {
		t.Errorf("Expected %d connected components, got %d", expected, components)
	}

	fmt.Printf("✅ Connected components: %d\n", components)
}

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

func TestBipartiteCheck(t *testing.T) {
	fmt.Println("Testing bipartite graph check...")

	// Bipartite graph (even cycle)
	graph1 := NewGraph(4, false)
	graph1.AddEdge(0, 1, 1)
	graph1.AddEdge(1, 2, 1)
	graph1.AddEdge(2, 3, 1)
	graph1.AddEdge(3, 0, 1)

	if !isBipartite(graph1) {
		t.Errorf("Expected even cycle to be bipartite")
	}

	// Non-bipartite graph (odd cycle)
	graph2 := NewGraph(3, false)
	graph2.AddEdge(0, 1, 1)
	graph2.AddEdge(1, 2, 1)
	graph2.AddEdge(2, 0, 1)

	if isBipartite(graph2) {
		t.Errorf("Expected odd cycle to be non-bipartite")
	}

	fmt.Println("✅ Bipartite check passed")
}

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
				colors[v] = 1 - colors[u]
				queue = append(queue, v)
			} else if colors[v] == colors[u] {
				return false
			}
		}
	}
	return true
}

// ========================================================================
// PERFORMANCE TESTS
// ========================================================================

func TestGraphPerformance(t *testing.T) {
	fmt.Println("Testing graph algorithm performance...")

	sizes := []int{100, 500, 1000}

	for _, size := range sizes {
		fmt.Printf("\n--- Testing with %d vertices ---\n", size)

		// Create dense graph
		graph := NewGraph(size, false)
		edges := 0
		maxEdges := size * (size - 1) / 4 // 25% density

		for i := 0; i < size && edges < maxEdges; i++ {
			for j := i + 1; j < size && edges < maxEdges; j++ {
				if (i*j)%4 == 0 { // Add some edges
					graph.AddEdge(i, j, 1)
					edges++
				}
			}
		}

		fmt.Printf("Created graph with %d edges (%.1f%% density)\n",
			edges, float64(edges)*200.0/float64(size*(size-1)))

		// Test DFS
		visited := make([]bool, size)
		var result []int
		dfs(graph, 0, visited, &result)
		fmt.Printf("DFS visited %d vertices\n", len(result))

		// Test BFS
		bfsResult := bfs(graph, 0)
		fmt.Printf("BFS visited %d vertices\n", len(bfsResult))

		// Test connected components
		components := countConnectedComponents(graph)
		fmt.Printf("Found %d connected components\n", components)
	}
}

// ========================================================================
// MAIN TEST RUNNER
// ========================================================================

func runAllTests() {
	tests := []struct {
		name string
		fn   func(*testing.T)
	}{
		{"DFS", TestDFS},
		{"BFS", TestBFS},
		{"Cycle Detection Undirected", TestCycleDetectionUndirected},
		{"Cycle Detection Directed", TestCycleDetectionDirected},
		{"Dijkstra", TestDijkstra},
		{"Topological Sort", TestTopologicalSort},
		{"Kruskal MST", TestKruskalMST},
		{"Connected Components", TestConnectedComponents},
		{"Bipartite Check", TestBipartiteCheck},
		{"Performance", TestGraphPerformance},
	}

	fmt.Println("🧪 RUNNING COMPREHENSIVE GRAPH ALGORITHM TESTS")
	fmt.Println("===============================================")

	passed := 0
	total := len(tests)

	for _, test := range tests {
		fmt.Printf("\n--- %s ---\n", test.name)

		// Create a mock testing.T
		mockT := &testing.T{}

		// Run test and catch panics
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("❌ Test %s PANICKED: %v\n", test.name, r)
				}
			}()

			test.fn(mockT)

			if !mockT.Failed() {
				fmt.Printf("✅ Test %s PASSED\n", test.name)
				passed++
			} else {
				fmt.Printf("❌ Test %s FAILED\n", test.name)
			}
		}()
	}

	fmt.Printf("\n🎯 TEST SUMMARY: %d/%d tests passed (%.1f%%)\n",
		passed, total, float64(passed)*100.0/float64(total))

	if passed == total {
		fmt.Println("🏆 All tests passed! Graph algorithms are working correctly.")
	} else {
		fmt.Printf("⚠️  %d tests failed. Review implementations.\n", total-passed)
	}
}

func main() {
	runAllTests()

	fmt.Println("\n📚 Graph Algorithm Test Coverage:")
	fmt.Println("• Traversal algorithms (DFS, BFS)")
	fmt.Println("• Cycle detection (directed and undirected)")
	fmt.Println("• Shortest path algorithms (Dijkstra)")
	fmt.Println("• Topological sorting")
	fmt.Println("• Minimum spanning trees (Kruskal)")
	fmt.Println("• Graph properties (connectivity, bipartiteness)")
	fmt.Println("• Performance testing with different graph sizes")

	fmt.Println("\n🚀 Next Steps:")
	fmt.Println("1. Add more shortest path algorithms (Bellman-Ford, Floyd-Warshall)")
	fmt.Println("2. Implement advanced algorithms (SCC, bridges, articulation points)")
	fmt.Println("3. Add more MST algorithms (Prim's)")
	fmt.Println("4. Create benchmarks for algorithm comparison")
	fmt.Println("5. Test with real-world graph datasets")
}
